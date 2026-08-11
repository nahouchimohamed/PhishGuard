package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phishguard/phishguard/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func TestHashTokenStable(t *testing.T) {
	a := HashToken("abc")
	b := HashToken("abc")
	c := HashToken("xyz")
	if a != b {
		t.Fatal("hash should be stable")
	}
	if a == c {
		t.Fatal("different inputs should differ")
	}
	if len(a) != 64 {
		t.Fatalf("expected sha256 hex length 64, got %d", len(a))
	}
}

func TestNewTokenAndNewClickToken(t *testing.T) {
	a, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewClickToken()
	if err != nil {
		t.Fatal(err)
	}
	if a == "" || b == "" {
		t.Fatal("tokens should be non-empty")
	}
	if a == b {
		t.Fatal("successive tokens should differ")
	}
	if HashToken(a) == a {
		t.Fatal("hash should differ from raw token")
	}
}

func TestSetSessionCookieAttributes(t *testing.T) {
	ttl := 2 * time.Hour
	svc := newService(nil, "phishguard_session", ttl)
	rec := httptest.NewRecorder()
	svc.SetSessionCookie(rec, "raw-token")

	res := rec.Result()
	cookies := res.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != "phishguard_session" {
		t.Fatalf("name=%q", c.Name)
	}
	if c.Value != "raw-token" {
		t.Fatalf("value=%q", c.Value)
	}
	if !c.HttpOnly {
		t.Fatal("expected HttpOnly")
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("SameSite=%v want Lax", c.SameSite)
	}
	if c.Path != "/" {
		t.Fatalf("Path=%q", c.Path)
	}
	if c.MaxAge != int(ttl.Seconds()) {
		t.Fatalf("MaxAge=%d want %d", c.MaxAge, int(ttl.Seconds()))
	}
}

func TestClearSessionCookieAttributes(t *testing.T) {
	svc := newService(nil, "phishguard_session", time.Hour)
	rec := httptest.NewRecorder()
	svc.ClearSessionCookie(rec)

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Value != "" {
		t.Fatalf("value should be empty, got %q", c.Value)
	}
	if !c.HttpOnly {
		t.Fatal("expected HttpOnly")
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("SameSite=%v want Lax", c.SameSite)
	}
	if c.Path != "/" {
		t.Fatalf("Path=%q", c.Path)
	}
	if c.MaxAge != -1 {
		t.Fatalf("MaxAge=%d want -1", c.MaxAge)
	}
}

type memUserStore struct {
	user     store.User
	sessions map[string]string // tokenHash -> userID
}

func (m *memUserStore) GetUserByEmail(ctx context.Context, email string) (store.User, error) {
	if email != m.user.Email {
		return store.User{}, store.ErrNotFound
	}
	return m.user, nil
}

func (m *memUserStore) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	if m.sessions == nil {
		m.sessions = map[string]string{}
	}
	m.sessions[tokenHash] = userID
	return nil
}

func (m *memUserStore) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	delete(m.sessions, tokenHash)
	return nil
}

func (m *memUserStore) GetUserBySessionTokenHash(ctx context.Context, tokenHash string) (store.User, error) {
	if m.sessions[tokenHash] != m.user.ID {
		return store.User{}, store.ErrNotFound
	}
	return m.user, nil
}

func TestLoginLogoutUserFromRequest(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	st := &memUserStore{
		user: store.User{
			ID:           "u1",
			OrgID:        "o1",
			Email:        "a@example.com",
			DisplayName:  "Ada",
			Role:         store.RoleAdmin,
			PasswordHash: string(hash),
		},
		sessions: map[string]string{},
	}
	svc := newService(st, "sess", time.Hour)

	_, _, err = svc.Login(context.Background(), "missing@example.com", "secret")
	if !errors.Is(err, store.ErrUnauthorized) {
		t.Fatalf("missing user: %v", err)
	}
	_, _, err = svc.Login(context.Background(), "a@example.com", "wrong")
	if !errors.Is(err, store.ErrUnauthorized) {
		t.Fatalf("bad password: %v", err)
	}

	user, raw, err := svc.Login(context.Background(), "a@example.com", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != "u1" || raw == "" {
		t.Fatalf("unexpected login result %#v %q", user, raw)
	}
	if st.sessions[HashToken(raw)] != "u1" {
		t.Fatal("session not stored")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "sess", Value: raw})
	got, err := svc.UserFromRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "u1" {
		t.Fatalf("got user %#v", got)
	}

	if err := svc.Logout(context.Background(), raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := st.sessions[HashToken(raw)]; ok {
		t.Fatal("session should be deleted")
	}
	_, err = svc.UserFromRequest(context.Background(), req)
	if !errors.Is(err, store.ErrUnauthorized) {
		t.Fatalf("after logout: %v", err)
	}
}

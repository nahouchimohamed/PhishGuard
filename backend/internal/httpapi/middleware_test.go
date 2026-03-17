package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/store"
)

type fakeAuth struct {
	user store.User
	err  error
}

func (f *fakeAuth) CookieName() string { return "sess" }
func (f *fakeAuth) Login(ctx context.Context, email, password string) (store.User, string, error) {
	return store.User{}, "", store.ErrUnauthorized
}
func (f *fakeAuth) Logout(ctx context.Context, rawToken string) error { return nil }
func (f *fakeAuth) UserFromRequest(ctx context.Context, r *http.Request) (store.User, error) {
	if f.err != nil {
		return store.User{}, f.err
	}
	return f.user, nil
}
func (f *fakeAuth) SetSessionCookie(w http.ResponseWriter, rawToken string) {}
func (f *fakeAuth) ClearSessionCookie(w http.ResponseWriter)                {}

type fakeStoreAPI struct {
	pingErr error
}

func (f *fakeStoreAPI) Ping(ctx context.Context) error { return f.pingErr }
func (f *fakeStoreAPI) ListEmployees(ctx context.Context, orgID string) ([]store.User, error) {
	return nil, nil
}
func (f *fakeStoreAPI) ListTemplates(ctx context.Context, orgID string) ([]store.Template, error) {
	return nil, nil
}
func (f *fakeStoreAPI) ListCampaigns(ctx context.Context, orgID string) ([]store.Campaign, error) {
	return nil, nil
}
func (f *fakeStoreAPI) GetCampaign(ctx context.Context, id, orgID string) (store.Campaign, error) {
	return store.Campaign{}, store.ErrNotFound
}
func (f *fakeStoreAPI) GetTemplate(ctx context.Context, id, orgID string) (store.Template, error) {
	return store.Template{}, store.ErrNotFound
}
func (f *fakeStoreAPI) ListRecipientRows(ctx context.Context, campaignID string) ([]store.RecipientRow, error) {
	return nil, nil
}
func (f *fakeStoreAPI) ListInbox(ctx context.Context, userID string) ([]store.InboxItem, error) {
	return nil, nil
}
func (f *fakeStoreAPI) GetInboxItem(ctx context.Context, recipientID, userID string) (store.InboxItem, error) {
	return store.InboxItem{}, store.ErrNotFound
}

type fakeCampaignAPI struct{}

func (f *fakeCampaignAPI) Launch(ctx context.Context, in campaign.LaunchInput) (store.Campaign, error) {
	return store.Campaign{}, store.ErrInvalidInput
}
func (f *fakeCampaignAPI) GetStats(ctx context.Context, campaignID string) (store.CampaignStats, error) {
	return store.CampaignStats{}, nil
}
func (f *fakeCampaignAPI) RetryLaunch(ctx context.Context, campaignID, orgID string) (store.Campaign, error) {
	return store.Campaign{}, store.ErrNotFound
}
func (f *fakeCampaignAPI) Open(ctx context.Context, recipientID, userID, orgID string) error {
	return nil
}
func (f *fakeCampaignAPI) Report(ctx context.Context, recipientID, userID, orgID string) error {
	return nil
}
func (f *fakeCampaignAPI) Click(ctx context.Context, token string) (store.Template, string, error) {
	return store.Template{}, "", store.ErrNotFound
}
func (f *fakeCampaignAPI) SimURL(token string) string { return "/sim/" + token }

type fakeRedis struct{ err error }

func (f *fakeRedis) Ping(ctx context.Context) error { return f.err }

func TestRequireAuthUnauthorized(t *testing.T) {
	auth := &fakeAuth{err: store.ErrUnauthorized}
	s := newServer(auth, &fakeStoreAPI{}, &fakeCampaignAPI{}, &fakeRedis{})
	h := s.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not run")
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/me", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestRequireAuthOK(t *testing.T) {
	auth := &fakeAuth{user: store.User{ID: "u1", Role: store.RoleAdmin}}
	s := newServer(auth, &fakeStoreAPI{}, &fakeCampaignAPI{}, &fakeRedis{})
	called := false
	h := s.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		called = true
		u := userFromCtx(r.Context())
		if u.ID != "u1" {
			t.Fatalf("user=%#v", u)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/me", nil))
	if !called || rec.Code != http.StatusNoContent {
		t.Fatalf("called=%v status=%d", called, rec.Code)
	}
}

func TestRequireRoleForbidden(t *testing.T) {
	auth := &fakeAuth{user: store.User{ID: "u1", Role: store.RoleEmployee}}
	s := newServer(auth, &fakeStoreAPI{}, &fakeCampaignAPI{}, &fakeRedis{})
	h := s.requireRole(store.RoleAdmin, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not run")
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/employees", nil))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestRequireRoleOK(t *testing.T) {
	auth := &fakeAuth{user: store.User{ID: "u1", Role: store.RoleAdmin}}
	s := newServer(auth, &fakeStoreAPI{}, &fakeCampaignAPI{}, &fakeRedis{})
	called := false
	h := s.requireRole(store.RoleAdmin, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/employees", nil))
	if !called || rec.Code != http.StatusNoContent {
		t.Fatalf("called=%v status=%d", called, rec.Code)
	}
}

func TestHealthzOK(t *testing.T) {
	s := newServer(&fakeAuth{}, &fakeStoreAPI{}, &fakeCampaignAPI{}, &fakeRedis{})
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" || body["db"] != "ok" || body["redis"] != "ok" {
		t.Fatalf("body=%v", body)
	}
}

func TestHealthzUnavailable(t *testing.T) {
	s := newServer(&fakeAuth{}, &fakeStoreAPI{pingErr: errors.New("db down")}, &fakeCampaignAPI{}, &fakeRedis{})
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "unavailable" || body["db"] != "down" {
		t.Fatalf("body=%v", body)
	}
}

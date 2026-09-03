package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/phishguard/phishguard/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// userStore is the subset of store.Store used by session auth.
type userStore interface {
	GetUserByEmail(ctx context.Context, email string) (store.User, error)
	CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
	GetUserBySessionTokenHash(ctx context.Context, tokenHash string) (store.User, error)
}

// dummyPasswordHash is a valid bcrypt hash used only to keep unknown-email
// login timing closer to the known-email path (password is never accepted).
var dummyPasswordHash = mustBcrypt("phishguard-dummy-password-not-used")

func mustBcrypt(pw string) []byte {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return h
}

type Service struct {
	store      userStore
	cookieName string
	ttl        time.Duration
	secure     bool
}

func New(st *store.Store, cookieName string, ttl time.Duration, secure bool) *Service {
	return &Service{store: st, cookieName: cookieName, ttl: ttl, secure: secure}
}

func newService(st userStore, cookieName string, ttl time.Duration, secure bool) *Service {
	return &Service{store: st, cookieName: cookieName, ttl: ttl, secure: secure}
}

func (s *Service) CookieName() string { return s.cookieName }

func (s *Service) Login(ctx context.Context, email, password string) (store.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return store.User{}, "", store.ValidationError{Msg: "email and password are required"}
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			_ = bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
			return store.User{}, "", store.ErrUnauthorized
		}
		return store.User{}, "", fmt.Errorf("lookup user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return store.User{}, "", store.ErrUnauthorized
	}

	raw, err := randomToken()
	if err != nil {
		return store.User{}, "", err
	}
	hash := HashToken(raw)
	if err := s.store.CreateSession(ctx, user.ID, hash, time.Now().UTC().Add(s.ttl)); err != nil {
		return store.User{}, "", fmt.Errorf("create session: %w", err)
	}
	return user, raw, nil
}

func (s *Service) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	return s.store.DeleteSessionByTokenHash(ctx, HashToken(rawToken))
}

func (s *Service) UserFromRequest(ctx context.Context, r *http.Request) (store.User, error) {
	c, err := r.Cookie(s.cookieName)
	if err != nil || c.Value == "" {
		return store.User{}, store.ErrUnauthorized
	}
	user, err := s.store.GetUserBySessionTokenHash(ctx, HashToken(c.Value))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return store.User{}, store.ErrUnauthorized
		}
		return store.User{}, fmt.Errorf("lookup session: %w", err)
	}
	return user, nil
}

func (s *Service) SetSessionCookie(w http.ResponseWriter, rawToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(s.ttl.Seconds()),
	})
}

func (s *Service) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// NewClickToken returns an unguessable token for the public simulation URL.
func NewClickToken() (string, error) {
	return randomToken()
}

// NewToken is an alias for NewClickToken for callers that prefer a shorter name.
func NewToken() (string, error) {
	return randomToken()
}

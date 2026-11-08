package httpapi

import (
	"context"
	"net/http"

	"github.com/phishguard/phishguard/internal/store"
)

type ctxKey int

const userKey ctxKey = 1

func userFromCtx(ctx context.Context) store.User {
	u, _ := ctx.Value(userKey).(store.User)
	return u
}

func (s *Server) requireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.auth.UserFromRequest(r.Context(), r)
		if err != nil {
			mapStoreErr(w, store.ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next(w, r.WithContext(ctx))
	})
}

func (s *Server) requireRole(role store.Role, next http.HandlerFunc) http.Handler {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		user := userFromCtx(r.Context())
		if user.Role != role {
			mapStoreErr(w, store.ErrForbidden)
			return
		}
		next(w, r)
	})
}

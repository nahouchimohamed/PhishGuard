package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/phishguard/phishguard/internal/store"
	"github.com/phishguard/phishguard/internal/validate"
)

type ctxKey int

const (
	userKey  ctxKey = 1
	reqIDKey ctxKey = 2
)

func userFromCtx(ctx context.Context) store.User {
	u, _ := ctx.Value(userKey).(store.User)
	return u
}

func requestIDFromCtx(ctx context.Context) string {
	id, _ := ctx.Value(reqIDKey).(string)
	return id
}

func (s *Server) requireAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.auth.UserFromRequest(r.Context(), r)
		if err != nil {
			mapStoreErr(w, err)
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

// recoverer logs panics and returns a JSON 500.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic", "err", rec, "request_id", requestIDFromCtx(r.Context()), "path", r.URL.Path)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), reqIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("150405.000000000")))
	}
	return hex.EncodeToString(b)
}

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusCapture{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(sw, r)
		slog.Info("request",
			"request_id", requestIDFromCtx(r.Context()),
			"method", r.Method,
			"path", redactPath(r.URL.Path),
			"status", sw.code,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusCapture struct {
	http.ResponseWriter
	code int
}

func (w *statusCapture) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// redactPath hides simulation tokens from access logs.
func redactPath(path string) string {
	const prefix = "/api/sim/"
	if strings.HasPrefix(path, prefix) && len(path) > len(prefix) {
		return prefix + "*"
	}
	return path
}

// withTimeout cancels the request context after d and returns JSON 504 if the
// handler has not yet written a response.
func withTimeout(d time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), d)
		defer cancel()

		tw := &timeoutWriter{ResponseWriter: w}
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(tw, r.WithContext(ctx))
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			tw.mu.Lock()
			if !tw.wroteHeader {
				tw.timedOut = true
				tw.wroteHeader = true
				writeError(w, http.StatusGatewayTimeout, "request timeout")
			}
			tw.mu.Unlock()
			<-done
		}
	})
}

type timeoutWriter struct {
	http.ResponseWriter
	mu          sync.Mutex
	wroteHeader bool
	timedOut    bool
}

func (w *timeoutWriter) WriteHeader(statusCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timedOut || w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *timeoutWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	if !w.wroteHeader {
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(b)
}

// pathUUID reads a path value and returns false (with 404) when it is not UUID-shaped.
func pathUUID(w http.ResponseWriter, r *http.Request, name string) (string, bool) {
	id := r.PathValue(name)
	if err := validate.UUID(id); err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return "", false
	}
	return id, true
}

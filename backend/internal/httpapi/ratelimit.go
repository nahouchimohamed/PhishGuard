package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// loginLimiter is a fixed-window counter keyed by normalized email.
type loginLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]loginWindow
	now     func() time.Time
}

type loginWindow struct {
	count int
	start time.Time
}

func newLoginLimiter(limit int, window time.Duration) *loginLimiter {
	return &loginLimiter{
		limit:   limit,
		window:  window,
		entries: map[string]loginWindow{},
		now:     time.Now,
	}
}

// allow reports whether the key may proceed and, if not, how long to wait.
func (l *loginLimiter) allow(key string) (ok bool, retryAfter time.Duration) {
	key = strings.TrimSpace(strings.ToLower(key))
	if key == "" {
		return true, 0
	}
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	e, exists := l.entries[key]
	if !exists || now.Sub(e.start) >= l.window {
		l.entries[key] = loginWindow{count: 1, start: now}
		return true, 0
	}
	if e.count >= l.limit {
		return false, l.window - now.Sub(e.start)
	}
	e.count++
	l.entries[key] = e
	return true, 0
}

func writeRetryAfter(w http.ResponseWriter, d time.Duration) {
	sec := int(d.Seconds())
	if sec < 1 {
		sec = 1
	}
	w.Header().Set("Retry-After", strconv.Itoa(sec))
	writeError(w, http.StatusTooManyRequests, "too many login attempts")
}

package httpserve

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// Serve runs srv on ln until ctx is cancelled, then Shutdown with grace.
// Callers should set ReadHeaderTimeout (and usually Read/Write/Idle timeouts)
// on srv before Serve; this helper only owns the serve/shutdown lifecycle.
func Serve(ctx context.Context, srv *http.Server, ln net.Listener, grace time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		err := <-errCh
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

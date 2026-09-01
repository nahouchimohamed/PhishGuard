package retry

import (
	"context"
	"time"
)

// Do runs fn up to attempts times. Between failures it waits base * 2^i,
// honouring ctx cancellation instead of an unconditional Sleep.
func Do(ctx context.Context, attempts int, base time.Duration, fn func(context.Context) error) error {
	if attempts < 1 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = fn(ctx)
		if err == nil {
			return nil
		}
		if i == attempts-1 {
			break
		}
		delay := base << i
		if err := Wait(ctx, delay); err != nil {
			return err
		}
	}
	return err
}

// Wait sleeps for d or until ctx is done.
func Wait(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

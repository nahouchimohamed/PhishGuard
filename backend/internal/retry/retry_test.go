package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoSucceedsAfterRetries(t *testing.T) {
	n := 0
	err := Do(context.Background(), 3, time.Millisecond, func(ctx context.Context) error {
		n++
		if n < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("n=%d", n)
	}
}

func TestDoHonoursCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Do(ctx, 5, time.Second, func(ctx context.Context) error {
		return errors.New("fail")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestWaitCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Wait(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestDoAttemptsTable(t *testing.T) {
	tests := []struct {
		name     string
		attempts int
		failFor  int
		wantN    int
		wantErr  bool
	}{
		{name: "first try", attempts: 3, failFor: 0, wantN: 1, wantErr: false},
		{name: "clamped min attempts", attempts: 0, failFor: 0, wantN: 1, wantErr: false},
		{name: "exhausts attempts", attempts: 2, failFor: 5, wantN: 2, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := 0
			err := Do(context.Background(), tt.attempts, time.Millisecond, func(ctx context.Context) error {
				n++
				if n <= tt.failFor {
					return errors.New("transient")
				}
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tt.wantErr)
			}
			if n != tt.wantN {
				t.Fatalf("n=%d want %d", n, tt.wantN)
			}
		})
	}
}

func TestWaitCompletes(t *testing.T) {
	start := time.Now()
	if err := Wait(context.Background(), 5*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 5*time.Millisecond {
		t.Fatal("Wait returned too early")
	}
}

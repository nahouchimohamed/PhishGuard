package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/config"
	"github.com/phishguard/phishguard/internal/httpserve"
	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/retry"
	"github.com/phishguard/phishguard/internal/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("store", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	dash, err := cache.New(cfg.RedisURL)
	if err != nil {
		slog.Error("redis", "err", err)
		os.Exit(1)
	}
	defer dash.Close()

	q, err := queue.New(ctx, cfg)
	if err != nil {
		slog.Error("queue", "err", err)
		os.Exit(1)
	}

	campSvc := campaign.New(st, q, dash, cfg.PublicAppURL)

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	metricsServer := &http.Server{
		Addr:              cfg.MetricsAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ln, err := net.Listen("tcp", cfg.MetricsAddr)
	if err != nil {
		slog.Error("metrics listen", "err", err)
		os.Exit(1)
	}
	go func() {
		slog.Info("worker metrics listening", "addr", cfg.MetricsAddr)
		if err := httpserve.Serve(ctx, metricsServer, ln, 5*time.Second); err != nil {
			slog.Error("metrics serve", "err", err)
			stop()
		}
	}()

	slog.Info("worker polling queue", "queue_url", cfg.SQSQueueURL)
	for {
		if ctx.Err() != nil {
			slog.Info("worker stopped")
			return
		}

		msgs, err := q.Receive(ctx, 5, 10)
		if err != nil {
			if ctx.Err() != nil {
				continue
			}
			slog.Error("receive", "err", err)
			_ = retry.Wait(ctx, 2*time.Second)
			continue
		}
		processBatch(ctx, msgs, campSvc.Deliver, q.Delete)
	}
}

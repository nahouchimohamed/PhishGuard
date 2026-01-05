package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/config"
	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/queue"
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
	}
	go func() {
		slog.Info("worker metrics listening", "addr", cfg.MetricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics listen", "err", err)
			stop()
		}
	}()

	slog.Info("worker polling queue", "queue_url", cfg.SQSQueueURL)
	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = metricsServer.Shutdown(shutdownCtx)
			cancel()
			slog.Info("worker stopped")
			return
		default:
		}

		msgs, err := q.Receive(ctx, 5, 10)
		if err != nil {
			if ctx.Err() != nil {
				continue
			}
			slog.Error("receive", "err", err)
			time.Sleep(2 * time.Second)
			continue
		}
		for _, m := range msgs {
			if err := campSvc.Deliver(ctx, m.Payload.RecipientID, m.Payload.CampaignID); err != nil {
				slog.Error("deliver", "err", err, "recipient_id", m.Payload.RecipientID)
				continue
			}
			if err := q.Delete(ctx, m.ReceiptHandle); err != nil {
				slog.Error("delete message", "err", err)
			}
		}
	}
}

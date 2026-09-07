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

	"github.com/phishguard/phishguard/internal/auth"
	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/config"
	"github.com/phishguard/phishguard/internal/httpapi"
	"github.com/phishguard/phishguard/internal/httpserve"
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

	if err := st.SeedIfEmpty(ctx, cfg.SeedPassword); err != nil {
		slog.Error("seed", "err", err)
		os.Exit(1)
	}

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

	authSvc := auth.New(st, cfg.SessionCookie, cfg.SessionTTL, cfg.SessionCookieSecure)
	campSvc := campaign.New(st, q, dash, cfg.PublicAppURL)
	srv := httpapi.New(authSvc, st, campSvc, dash)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		slog.Error("listen", "err", err)
		os.Exit(1)
	}
	slog.Info("api listening", "addr", cfg.HTTPAddr)
	if err := httpserve.Serve(ctx, httpServer, ln, 10*time.Second); err != nil {
		slog.Error("serve", "err", err)
		os.Exit(1)
	}
	slog.Info("api stopped")
}

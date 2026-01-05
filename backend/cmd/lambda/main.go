package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/config"
	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/store"
)

// Same Deliver function as the long-running worker, behind a Lambda adapter.
// Everyday docker compose uses cmd/worker against ElasticMQ. This binary is for
// LocalStack and Terraform deployments only — do not attach both consumers to one queue.
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	st, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("store", "err", err)
		os.Exit(1)
	}

	dash, err := cache.New(cfg.RedisURL)
	if err != nil {
		slog.Error("redis", "err", err)
		os.Exit(1)
	}

	q, err := queue.New(ctx, cfg)
	if err != nil {
		slog.Error("queue", "err", err)
		os.Exit(1)
	}

	campSvc := campaign.New(st, q, dash, cfg.PublicAppURL)

	lambda.Start(func(ctx context.Context, ev events.SQSEvent) error {
		for _, record := range ev.Records {
			msg, err := queue.ParseBody(record.Body)
			if err != nil {
				slog.Error("parse", "err", err)
				continue
			}
			if err := campSvc.Deliver(ctx, msg.RecipientID, msg.CampaignID); err != nil {
				return err
			}
		}
		return nil
	})
}

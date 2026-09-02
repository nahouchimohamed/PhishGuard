package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/queue"
)

type deliverFunc func(ctx context.Context, recipientID, campaignID string) error
type deleteFunc func(ctx context.Context, receiptHandle string) error

const perMessageDeadline = 20 * time.Second

// processBatch delivers each message under a per-message deadline.
// Invalid payloads are poisoned (deleted) and counted; deliver failures are counted and left for retry.
func processBatch(ctx context.Context, msgs []queue.Received, deliver deliverFunc, del deleteFunc) {
	for _, m := range msgs {
		if err := ctx.Err(); err != nil {
			return
		}
		msgCtx, cancel := context.WithTimeout(ctx, perMessageDeadline)
		if err := m.Payload.Validate(); err != nil {
			metrics.QueuePoison.Inc()
			slog.Warn("poison message", "err", err, "receipt", m.ReceiptHandle)
			if delErr := del(msgCtx, m.ReceiptHandle); delErr != nil {
				slog.Error("delete poison message", "err", delErr)
			}
			cancel()
			continue
		}
		if err := deliver(msgCtx, m.Payload.RecipientID, m.Payload.CampaignID); err != nil {
			metrics.DeliveryFailures.Inc()
			slog.Error("deliver", "err", err, "recipient_id", m.Payload.RecipientID)
			cancel()
			continue
		}
		if err := del(msgCtx, m.ReceiptHandle); err != nil {
			slog.Error("delete message", "err", err)
		}
		cancel()
	}
}

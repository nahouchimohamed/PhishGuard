package main

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/phishguard/phishguard/internal/queue"
)

func TestProcessBatchPoisonsInvalid(t *testing.T) {
	var deleted []string
	var delivered int
	msgs := []queue.Received{
		{Payload: queue.DeliverMessage{RecipientID: "", CampaignID: "c1"}, ReceiptHandle: "h1"},
		{Payload: queue.DeliverMessage{RecipientID: "r1", CampaignID: "c1"}, ReceiptHandle: "h2"},
	}
	processBatch(context.Background(), msgs,
		func(ctx context.Context, recipientID, campaignID string) error {
			delivered++
			return nil
		},
		func(ctx context.Context, receiptHandle string) error {
			deleted = append(deleted, receiptHandle)
			return nil
		},
	)
	if delivered != 1 {
		t.Fatalf("delivered=%d", delivered)
	}
	if len(deleted) != 2 || deleted[0] != "h1" || deleted[1] != "h2" {
		t.Fatalf("deleted=%v", deleted)
	}
}

func TestProcessBatchKeepsFailedForRetry(t *testing.T) {
	var deleted []string
	var mu sync.Mutex
	msgs := []queue.Received{
		{Payload: queue.DeliverMessage{RecipientID: "r1", CampaignID: "c1"}, ReceiptHandle: "h1"},
	}
	processBatch(context.Background(), msgs,
		func(ctx context.Context, recipientID, campaignID string) error {
			return errors.New("db down")
		},
		func(ctx context.Context, receiptHandle string) error {
			mu.Lock()
			deleted = append(deleted, receiptHandle)
			mu.Unlock()
			return nil
		},
	)
	if len(deleted) != 0 {
		t.Fatalf("should not delete on deliver failure: %v", deleted)
	}
}

func TestProcessBatchHonoursCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	processBatch(ctx, []queue.Received{
		{Payload: queue.DeliverMessage{RecipientID: "r1", CampaignID: "c1"}, ReceiptHandle: "h1"},
	}, func(ctx context.Context, recipientID, campaignID string) error {
		called = true
		return nil
	}, func(ctx context.Context, receiptHandle string) error { return nil })
	if called {
		t.Fatal("should not deliver after cancel")
	}
}

func TestProcessBatchEmpty(t *testing.T) {
	called := false
	processBatch(context.Background(), nil,
		func(ctx context.Context, recipientID, campaignID string) error {
			called = true
			return nil
		},
		func(ctx context.Context, receiptHandle string) error {
			called = true
			return nil
		},
	)
	if called {
		t.Fatal("empty batch should be a no-op")
	}
}

func TestProcessBatchPoisonsMissingCampaign(t *testing.T) {
	var deleted []string
	processBatch(context.Background(), []queue.Received{
		{Payload: queue.DeliverMessage{RecipientID: "r1", CampaignID: ""}, ReceiptHandle: "hx"},
	}, func(ctx context.Context, recipientID, campaignID string) error {
		t.Fatal("should not deliver poison message")
		return nil
	}, func(ctx context.Context, receiptHandle string) error {
		deleted = append(deleted, receiptHandle)
		return nil
	})
	if len(deleted) != 1 || deleted[0] != "hx" {
		t.Fatalf("deleted=%v", deleted)
	}
}

package campaign

import (
	"context"

	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/store"
)

// Store is the persistence surface used by campaign.Service.
type Store interface {
	GetTemplate(ctx context.Context, id, orgID string) (store.Template, error)
	CreateCampaignWithRecipients(ctx context.Context, orgID, createdBy, templateID, name string, userIDs, tokens []string) (store.Campaign, []string, error)
	MarkCampaignActive(ctx context.Context, campaignID string) error
	GetCampaign(ctx context.Context, id, orgID string) (store.Campaign, error)
	ListUndeliveredRecipientIDs(ctx context.Context, campaignID string) ([]string, error)
	MarkDelivered(ctx context.Context, recipientID string) (bool, error)
	GetRecipient(ctx context.Context, recipientID string) (store.Recipient, error)
	GetUserByID(ctx context.Context, id string) (store.User, error)
	InsertEvent(ctx context.Context, orgID, campaignID, userID, kind, detail string) error
	GetCampaignStats(ctx context.Context, campaignID string) (store.CampaignStats, error)
	GetInboxItem(ctx context.Context, recipientID, userID string) (store.InboxItem, error)
	MarkReported(ctx context.Context, recipientID, userID string) error
	MarkOpened(ctx context.Context, recipientID, userID string) error
	MarkClickedByToken(ctx context.Context, token string) (store.Recipient, store.Template, error)
}

// Publisher publishes deliver jobs.
type Publisher interface {
	PublishDeliver(ctx context.Context, msg queue.DeliverMessage) error
}

// StatsCache caches campaign stats and supports invalidation after mutations.
type StatsCache interface {
	GetStats(ctx context.Context, campaignID string) (store.CampaignStats, bool, error)
	SetStats(ctx context.Context, campaignID string, st store.CampaignStats) error
	Invalidate(ctx context.Context, campaignID string) error
}

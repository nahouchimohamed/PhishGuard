package campaign

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/phishguard/phishguard/internal/auth"
	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/store"
)

type Service struct {
	store  Store
	queue  Publisher
	cache  StatsCache
	appURL string
}

func New(st *store.Store, q *queue.Client, dash *cache.Dashboard, appURL string) *Service {
	return &Service{store: st, queue: q, cache: dash, appURL: appURL}
}

func newService(st Store, q Publisher, dash StatsCache, appURL string) *Service {
	return &Service{store: st, queue: q, cache: dash, appURL: appURL}
}

type LaunchInput struct {
	OrgID            string
	CreatedBy        string
	Name             string
	TemplateID       string
	RecipientUserIDs []string
}

// Launch writes the campaign + recipients, then publishes one SQS message per recipient.
// If publish fails partway, the campaign stays "launching" and RetryLaunch can finish it.
func (s *Service) Launch(ctx context.Context, in LaunchInput) (store.Campaign, error) {
	if in.Name == "" || in.TemplateID == "" || len(in.RecipientUserIDs) == 0 {
		return store.Campaign{}, store.ErrInvalidInput
	}
	if _, err := s.store.GetTemplate(ctx, in.TemplateID, in.OrgID); err != nil {
		return store.Campaign{}, err
	}

	tokens := make([]string, len(in.RecipientUserIDs))
	for i := range tokens {
		tok, err := auth.NewClickToken()
		if err != nil {
			return store.Campaign{}, err
		}
		tokens[i] = tok
	}

	campaign, recipientIDs, err := s.store.CreateCampaignWithRecipients(
		ctx, in.OrgID, in.CreatedBy, in.TemplateID, in.Name, in.RecipientUserIDs, tokens,
	)
	if err != nil {
		return store.Campaign{}, err
	}

	if err := s.publishAll(ctx, campaign.ID, recipientIDs); err != nil {
		slog.Error("publish after launch incomplete", "campaign_id", campaign.ID, "err", err)
		return campaign, fmt.Errorf("campaign saved but queue publish incomplete: %w", err)
	}

	if err := s.store.MarkCampaignActive(ctx, campaign.ID); err != nil {
		return campaign, err
	}
	campaign.Status = store.StatusActive
	metrics.CampaignLaunches.Inc()
	return campaign, nil
}

// RetryLaunch re-publishes undelivered recipients for a campaign stuck in launching.
func (s *Service) RetryLaunch(ctx context.Context, campaignID, orgID string) (store.Campaign, error) {
	c, err := s.store.GetCampaign(ctx, campaignID, orgID)
	if err != nil {
		return store.Campaign{}, err
	}
	ids, err := s.store.ListUndeliveredRecipientIDs(ctx, campaignID)
	if err != nil {
		return store.Campaign{}, err
	}
	if err := s.publishAll(ctx, campaignID, ids); err != nil {
		return c, err
	}
	if err := s.store.MarkCampaignActive(ctx, campaignID); err != nil {
		return c, err
	}
	c.Status = store.StatusActive
	return c, nil
}

func (s *Service) publishAll(ctx context.Context, campaignID string, recipientIDs []string) error {
	var firstErr error
	for _, rid := range recipientIDs {
		err := s.queue.PublishDeliver(ctx, queue.DeliverMessage{
			RecipientID: rid,
			CampaignID:  campaignID,
		})
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Deliver marks a recipient delivered. Safe to call twice for the same ID.
func (s *Service) Deliver(ctx context.Context, recipientID, campaignID string) error {
	changed, err := s.store.MarkDelivered(ctx, recipientID)
	if err != nil {
		return err
	}
	if !changed {
		slog.Info("duplicate deliver ignored", "recipient_id", recipientID)
		return nil
	}

	rec, err := s.store.GetRecipient(ctx, recipientID)
	if err != nil {
		return err
	}
	user, err := s.store.GetUserByID(ctx, rec.UserID)
	if err != nil {
		return err
	}
	_ = s.store.InsertEvent(ctx, user.OrgID, campaignID, user.ID, "delivered", recipientID)
	_ = s.cache.Invalidate(ctx, campaignID)
	metrics.Deliveries.Inc()
	slog.Info("delivered", "recipient_id", recipientID, "campaign_id", campaignID)
	return nil
}

func (s *Service) GetStats(ctx context.Context, campaignID string) (store.CampaignStats, error) {
	if st, ok, err := s.cache.GetStats(ctx, campaignID); err != nil {
		slog.Warn("cache get failed", "err", err)
	} else if ok {
		return st, nil
	}

	st, err := s.store.GetCampaignStats(ctx, campaignID)
	if err != nil {
		return store.CampaignStats{}, err
	}
	if err := s.cache.SetStats(ctx, campaignID, st); err != nil {
		slog.Warn("cache set failed", "err", err)
	}
	return st, nil
}

func (s *Service) Report(ctx context.Context, recipientID, userID, orgID string) error {
	item, err := s.store.GetInboxItem(ctx, recipientID, userID)
	if err != nil {
		return err
	}
	if err := s.store.MarkReported(ctx, recipientID, userID); err != nil {
		return err
	}
	_ = s.store.InsertEvent(ctx, orgID, item.CampaignID, userID, "reported", recipientID)
	_ = s.cache.Invalidate(ctx, item.CampaignID)
	metrics.Reports.Inc()
	return nil
}

func (s *Service) Open(ctx context.Context, recipientID, userID, orgID string) error {
	item, err := s.store.GetInboxItem(ctx, recipientID, userID)
	if err != nil {
		return err
	}
	if err := s.store.MarkOpened(ctx, recipientID, userID); err != nil {
		return err
	}
	_ = s.store.InsertEvent(ctx, orgID, item.CampaignID, userID, "opened", recipientID)
	_ = s.cache.Invalidate(ctx, item.CampaignID)
	return nil
}

func (s *Service) Click(ctx context.Context, token string) (store.Template, string, error) {
	rec, tmpl, err := s.store.MarkClickedByToken(ctx, token)
	if err != nil {
		return store.Template{}, "", err
	}
	user, err := s.store.GetUserByID(ctx, rec.UserID)
	if err != nil {
		return store.Template{}, "", err
	}
	_ = s.store.InsertEvent(ctx, user.OrgID, rec.CampaignID, user.ID, "clicked", rec.ID)
	_ = s.cache.Invalidate(ctx, rec.CampaignID)
	metrics.Clicks.Inc()
	return tmpl, rec.CampaignID, nil
}

func (s *Service) SimURL(token string) string {
	return fmt.Sprintf("%s/sim/%s", s.appURL, token)
}

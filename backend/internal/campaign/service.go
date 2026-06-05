package campaign

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/phishguard/phishguard/internal/auth"
	"github.com/phishguard/phishguard/internal/cache"
	"github.com/phishguard/phishguard/internal/metrics"
	"github.com/phishguard/phishguard/internal/queue"
	"github.com/phishguard/phishguard/internal/store"
	"github.com/phishguard/phishguard/internal/validate"
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
	name := strings.TrimSpace(in.Name)
	if err := validate.Name(name, 1, 120); err != nil {
		return store.Campaign{}, store.ValidationError{Msg: err.Error()}
	}
	if err := validate.UUID(in.TemplateID); err != nil {
		return store.Campaign{}, store.ValidationError{Msg: "template_id must be a UUID"}
	}
	ids, err := normalizeRecipientIDs(in.RecipientUserIDs)
	if err != nil {
		return store.Campaign{}, err
	}
	if _, err := s.store.GetTemplate(ctx, in.TemplateID, in.OrgID); err != nil {
		return store.Campaign{}, err
	}

	tokens := make([]string, len(ids))
	for i := range tokens {
		tok, err := auth.NewClickToken()
		if err != nil {
			return store.Campaign{}, err
		}
		tokens[i] = tok
	}

	campaign, recipientIDs, err := s.store.CreateCampaignWithRecipients(
		ctx, in.OrgID, in.CreatedBy, in.TemplateID, name, ids, tokens,
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

func normalizeRecipientIDs(ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, store.ValidationError{Msg: "at least one recipient is required"}
	}
	if len(ids) > 500 {
		return nil, store.ValidationError{Msg: "at most 500 recipients are allowed"}
	}
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if err := validate.UUID(id); err != nil {
			return nil, store.ValidationError{Msg: "recipient_user_ids must be UUIDs"}
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if len(out) == 0 {
		return nil, store.ValidationError{Msg: "at least one recipient is required"}
	}
	return out, nil
}

// RetryLaunch re-publishes undelivered recipients for a campaign stuck in launching.
func (s *Service) RetryLaunch(ctx context.Context, campaignID, orgID string) (store.Campaign, error) {
	c, err := s.store.GetCampaign(ctx, campaignID, orgID)
	if err != nil {
		return store.Campaign{}, err
	}
	if c.Status != store.StatusLaunching {
		return store.Campaign{}, store.ErrConflict
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
	metrics.CampaignLaunches.Inc()
	return c, nil
}

func (s *Service) publishAll(ctx context.Context, campaignID string, recipientIDs []string) error {
	var errs []error
	for _, rid := range recipientIDs {
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}
		err := s.queue.PublishDeliver(ctx, queue.DeliverMessage{
			RecipientID: rid,
			CampaignID:  campaignID,
		})
		if err != nil {
			metrics.QueuePublishErrors.Inc()
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Deliver marks a recipient delivered. Safe to call twice for the same ID.
// Lookups run before MarkDelivered so a failed lookup cannot leave a delivered
// row without its event/metric. Duplicate deliveries still invalidate the cache.
func (s *Service) Deliver(ctx context.Context, recipientID, campaignID string) error {
	rec, err := s.store.GetRecipient(ctx, recipientID)
	if err != nil {
		return err
	}
	user, err := s.store.GetUserByID(ctx, rec.UserID)
	if err != nil {
		return err
	}
	cid := rec.CampaignID

	changed, err := s.store.MarkDelivered(ctx, recipientID)
	if err != nil {
		return err
	}
	if !changed {
		slog.Info("duplicate deliver ignored", "recipient_id", recipientID)
		logInvalidate(ctx, s.cache, cid)
		return nil
	}

	logInsertEvent(ctx, s.store, user.OrgID, cid, user.ID, "delivered", recipientID)
	logInvalidate(ctx, s.cache, cid)
	metrics.Deliveries.Inc()
	slog.Info("delivered", "recipient_id", recipientID, "campaign_id", cid)
	return nil
}

func (s *Service) GetStats(ctx context.Context, campaignID string) (store.CampaignStats, error) {
	if st, ok, err := s.cache.GetStats(ctx, campaignID); err != nil {
		metrics.StatsCache.WithLabelValues("error").Inc()
		slog.Warn("cache get failed", "err", err)
	} else if ok {
		metrics.StatsCache.WithLabelValues("hit").Inc()
		return st, nil
	} else {
		metrics.StatsCache.WithLabelValues("miss").Inc()
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
	first, err := s.store.MarkReported(ctx, recipientID, userID)
	if err != nil {
		return err
	}
	if !first {
		return nil
	}
	logInsertEvent(ctx, s.store, orgID, item.CampaignID, userID, "reported", recipientID)
	logInvalidate(ctx, s.cache, item.CampaignID)
	metrics.Reports.Inc()
	return nil
}

func (s *Service) Open(ctx context.Context, recipientID, userID, orgID string) error {
	item, err := s.store.GetInboxItem(ctx, recipientID, userID)
	if err != nil {
		return err
	}
	first, err := s.store.MarkOpened(ctx, recipientID, userID)
	if err != nil {
		return err
	}
	if !first {
		return nil
	}
	logInsertEvent(ctx, s.store, orgID, item.CampaignID, userID, "opened", recipientID)
	logInvalidate(ctx, s.cache, item.CampaignID)
	return nil
}

func (s *Service) Click(ctx context.Context, token string) (store.Template, string, error) {
	rec, tmpl, first, err := s.store.MarkClickedByToken(ctx, token)
	if err != nil {
		return store.Template{}, "", err
	}
	if !first {
		return tmpl, rec.CampaignID, nil
	}
	user, err := s.store.GetUserByID(ctx, rec.UserID)
	if err != nil {
		return store.Template{}, "", err
	}
	logInsertEvent(ctx, s.store, user.OrgID, rec.CampaignID, user.ID, "clicked", rec.ID)
	logInvalidate(ctx, s.cache, rec.CampaignID)
	metrics.Clicks.Inc()
	return tmpl, rec.CampaignID, nil
}

func (s *Service) SimURL(token string) string {
	return fmt.Sprintf("%s/sim/%s", s.appURL, token)
}

func logInsertEvent(ctx context.Context, st Store, orgID, campaignID, userID, kind, detail string) {
	if err := st.InsertEvent(ctx, orgID, campaignID, userID, kind, detail); err != nil {
		slog.Error("insert event failed", "kind", kind, "campaign_id", campaignID, "err", err)
	}
}

func logInvalidate(ctx context.Context, c StatsCache, campaignID string) {
	if err := c.Invalidate(ctx, campaignID); err != nil {
		slog.Error("cache invalidate failed", "campaign_id", campaignID, "err", err)
	}
}

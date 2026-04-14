package httpapi

import (
	"context"
	"net/http"

	"github.com/phishguard/phishguard/internal/campaign"
	"github.com/phishguard/phishguard/internal/store"
)

// authService is the session/auth surface used by HTTP handlers and middleware.
type authService interface {
	CookieName() string
	Login(ctx context.Context, email, password string) (store.User, string, error)
	Logout(ctx context.Context, rawToken string) error
	UserFromRequest(ctx context.Context, r *http.Request) (store.User, error)
	SetSessionCookie(w http.ResponseWriter, rawToken string)
	ClearSessionCookie(w http.ResponseWriter)
}

// storeAPI is the persistence surface used by HTTP handlers.
type storeAPI interface {
	Ping(ctx context.Context) error
	ListEmployees(ctx context.Context, orgID string) ([]store.User, error)
	ListTemplates(ctx context.Context, orgID string) ([]store.Template, error)
	ListCampaigns(ctx context.Context, orgID string) ([]store.Campaign, error)
	GetCampaign(ctx context.Context, id, orgID string) (store.Campaign, error)
	GetTemplate(ctx context.Context, id, orgID string) (store.Template, error)
	ListRecipientRows(ctx context.Context, campaignID string) ([]store.RecipientRow, error)
	ListInbox(ctx context.Context, userID string) ([]store.InboxItem, error)
	GetInboxItem(ctx context.Context, recipientID, userID string) (store.InboxItem, error)
}

// campaignAPI is the campaign orchestration surface used by HTTP handlers.
type campaignAPI interface {
	Launch(ctx context.Context, in campaign.LaunchInput) (store.Campaign, error)
	GetStats(ctx context.Context, campaignID string) (store.CampaignStats, error)
	RetryLaunch(ctx context.Context, campaignID, orgID string) (store.Campaign, error)
	Open(ctx context.Context, recipientID, userID, orgID string) error
	Report(ctx context.Context, recipientID, userID, orgID string) error
	Click(ctx context.Context, token string) (store.Template, string, error)
	SimURL(token string) string
}

// redisPinger is optional Redis health for /healthz.
type redisPinger interface {
	Ping(ctx context.Context) error
}

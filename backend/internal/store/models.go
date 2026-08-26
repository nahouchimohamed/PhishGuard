package store

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type CampaignStatus string

const (
	StatusDraft     CampaignStatus = "draft"
	StatusLaunching CampaignStatus = "launching"
	StatusActive    CampaignStatus = "active"
	StatusClosed    CampaignStatus = "closed"
)

type User struct {
	ID           string    `json:"id"`
	OrgID        string    `json:"org_id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Template struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"org_id"`
	Name       string    `json:"name"`
	Subject    string    `json:"subject"`
	BodyText   string    `json:"body_text"`
	LessonText string    `json:"lesson_text"`
	CreatedAt  time.Time `json:"created_at"`
}

type Campaign struct {
	ID         string         `json:"id"`
	OrgID      string         `json:"org_id"`
	CreatedBy  string         `json:"created_by"`
	TemplateID string         `json:"template_id"`
	Name       string         `json:"name"`
	Status     CampaignStatus `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	LaunchedAt *time.Time     `json:"launched_at,omitempty"`
}

type Recipient struct {
	ID          string     `json:"id"`
	CampaignID  string     `json:"campaign_id"`
	UserID      string     `json:"user_id"`
	ClickToken  string     `json:"click_token"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	ClickedAt   *time.Time `json:"clicked_at,omitempty"`
	ReportedAt  *time.Time `json:"reported_at,omitempty"`
}

type InboxItem struct {
	RecipientID string     `json:"recipient_id"`
	CampaignID  string     `json:"campaign_id"`
	Subject     string     `json:"subject"`
	BodyText    string     `json:"body_text"`
	LessonText  string     `json:"lesson_text"`
	FromLabel   string     `json:"from_label"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	ClickedAt   *time.Time `json:"clicked_at,omitempty"`
	ReportedAt  *time.Time `json:"reported_at,omitempty"`
	ClickToken  string     `json:"click_token"`
	SimURL      string     `json:"sim_url,omitempty"`
}

type CampaignStats struct {
	Total     int `json:"total"`
	Delivered int `json:"delivered"`
	Opened    int `json:"opened"`
	Clicked   int `json:"clicked"`
	Reported  int `json:"reported"`
}

type CampaignDetail struct {
	Campaign   Campaign
	Template   Template
	Stats      CampaignStats
	Recipients []RecipientRow
}

type RecipientRow struct {
	RecipientID string     `json:"recipient_id"`
	UserID      string     `json:"user_id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	OpenedAt    *time.Time `json:"opened_at,omitempty"`
	ClickedAt   *time.Time `json:"clicked_at,omitempty"`
	ReportedAt  *time.Time `json:"reported_at,omitempty"`
}

type Event struct {
	ID         string    `json:"id"`
	OrgID      string    `json:"org_id"`
	CampaignID string    `json:"campaign_id"`
	UserID     string    `json:"user_id,omitempty"`
	Kind       string    `json:"kind"`
	Detail     string    `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

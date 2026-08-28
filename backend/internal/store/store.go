package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
)

// ValidationError wraps ErrInvalidInput with a clearer client-facing message.
type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	if e.Msg == "" {
		return ErrInvalidInput.Error()
	}
	return e.Msg
}

func (e ValidationError) Unwrap() error { return ErrInvalidInput }

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (s *Store) CreateOrg(ctx context.Context, name string) (string, error) {
	var id string
	err := s.pool.QueryRow(ctx, `INSERT INTO orgs (name) VALUES ($1) RETURNING id`, name).Scan(&id)
	return id, err
}

func (s *Store) CreateUser(ctx context.Context, orgID, email, displayName, role, passwordHash string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		INSERT INTO users (org_id, email, display_name, role, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, org_id, email, display_name, role, password_hash, created_at
	`, orgID, email, displayName, role, passwordHash).Scan(
		&u.ID, &u.OrgID, &u.Email, &u.DisplayName, &u.Role, &u.PasswordHash, &u.CreatedAt,
	)
	return u, err
}

func (s *Store) CreateTemplate(ctx context.Context, orgID, name, subject, body, lesson string) (Template, error) {
	var t Template
	err := s.pool.QueryRow(ctx, `
		INSERT INTO templates (org_id, name, subject, body_text, lesson_text)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, org_id, name, subject, body_text, lesson_text, created_at
	`, orgID, name, subject, body, lesson).Scan(
		&t.ID, &t.OrgID, &t.Name, &t.Subject, &t.BodyText, &t.LessonText, &t.CreatedAt,
	)
	return t, err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT id, org_id, email, display_name, role, password_hash, created_at
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.OrgID, &u.Email, &u.DisplayName, &u.Role, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) GetUserByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT id, org_id, email, display_name, role, password_hash, created_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.OrgID, &u.Email, &u.DisplayName, &u.Role, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) ListEmployees(ctx context.Context, orgID string) ([]User, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, org_id, email, display_name, role, password_hash, created_at
		FROM users
		WHERE org_id = $1 AND role = 'employee'
		ORDER BY display_name
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.OrgID, &u.Email, &u.DisplayName, &u.Role, &u.PasswordHash, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *Store) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (s *Store) GetUserBySessionTokenHash(ctx context.Context, tokenHash string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.org_id, u.email, u.display_name, u.role, u.password_hash, u.created_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now()
	`, tokenHash).Scan(&u.ID, &u.OrgID, &u.Email, &u.DisplayName, &u.Role, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}

func (s *Store) ListTemplates(ctx context.Context, orgID string) ([]Template, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, org_id, name, subject, body_text, lesson_text, created_at
		FROM templates WHERE org_id = $1 ORDER BY name
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Template
	for rows.Next() {
		var t Template
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.Subject, &t.BodyText, &t.LessonText, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) GetTemplate(ctx context.Context, id, orgID string) (Template, error) {
	var t Template
	err := s.pool.QueryRow(ctx, `
		SELECT id, org_id, name, subject, body_text, lesson_text, created_at
		FROM templates WHERE id = $1 AND org_id = $2
	`, id, orgID).Scan(&t.ID, &t.OrgID, &t.Name, &t.Subject, &t.BodyText, &t.LessonText, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Template{}, ErrNotFound
	}
	return t, err
}

func (s *Store) CreateCampaignWithRecipients(
	ctx context.Context,
	orgID, createdBy, templateID, name string,
	recipientUserIDs []string,
	clickTokens []string,
) (Campaign, []string, error) {
	if len(recipientUserIDs) == 0 {
		return Campaign{}, nil, ErrInvalidInput
	}
	if len(recipientUserIDs) != len(clickTokens) {
		return Campaign{}, nil, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Campaign{}, nil, err
	}
	defer tx.Rollback(ctx)

	var c Campaign
	err = tx.QueryRow(ctx, `
		INSERT INTO campaigns (org_id, created_by, template_id, name, status)
		VALUES ($1, $2, $3, $4, 'launching')
		RETURNING id, org_id, created_by, template_id, name, status, created_at, launched_at
	`, orgID, createdBy, templateID, name).Scan(
		&c.ID, &c.OrgID, &c.CreatedBy, &c.TemplateID, &c.Name, &c.Status, &c.CreatedAt, &c.LaunchedAt,
	)
	if err != nil {
		return Campaign{}, nil, err
	}

	recipientIDs := make([]string, 0, len(recipientUserIDs))
	for i, userID := range recipientUserIDs {
		var rid string
		err = tx.QueryRow(ctx, `
			INSERT INTO campaign_recipients (campaign_id, user_id, click_token)
			VALUES ($1, $2, $3)
			RETURNING id
		`, c.ID, userID, clickTokens[i]).Scan(&rid)
		if err != nil {
			return Campaign{}, nil, err
		}
		recipientIDs = append(recipientIDs, rid)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO events (org_id, campaign_id, user_id, kind, detail)
		VALUES ($1, $2, $3, 'campaign_created', $4)
	`, orgID, c.ID, createdBy, name); err != nil {
		return Campaign{}, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Campaign{}, nil, err
	}
	return c, recipientIDs, nil
}

func (s *Store) MarkCampaignActive(ctx context.Context, campaignID string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE campaigns
		SET status = 'active', launched_at = COALESCE(launched_at, now())
		WHERE id = $1
	`, campaignID)
	return err
}

func (s *Store) ListCampaigns(ctx context.Context, orgID string) ([]Campaign, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, org_id, created_by, template_id, name, status, created_at, launched_at
		FROM campaigns WHERE org_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Campaign
	for rows.Next() {
		var c Campaign
		if err := rows.Scan(&c.ID, &c.OrgID, &c.CreatedBy, &c.TemplateID, &c.Name, &c.Status, &c.CreatedAt, &c.LaunchedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetCampaign(ctx context.Context, id, orgID string) (Campaign, error) {
	var c Campaign
	err := s.pool.QueryRow(ctx, `
		SELECT id, org_id, created_by, template_id, name, status, created_at, launched_at
		FROM campaigns WHERE id = $1 AND org_id = $2
	`, id, orgID).Scan(&c.ID, &c.OrgID, &c.CreatedBy, &c.TemplateID, &c.Name, &c.Status, &c.CreatedAt, &c.LaunchedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Campaign{}, ErrNotFound
	}
	return c, err
}

func (s *Store) GetCampaignStats(ctx context.Context, campaignID string) (CampaignStats, error) {
	var st CampaignStats
	err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*)::int,
			COUNT(delivered_at)::int,
			COUNT(opened_at)::int,
			COUNT(clicked_at)::int,
			COUNT(reported_at)::int
		FROM campaign_recipients
		WHERE campaign_id = $1
	`, campaignID).Scan(&st.Total, &st.Delivered, &st.Opened, &st.Clicked, &st.Reported)
	return st, err
}

func (s *Store) ListRecipientRows(ctx context.Context, campaignID string) ([]RecipientRow, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT r.id, r.user_id, u.email, u.display_name,
		       r.delivered_at, r.opened_at, r.clicked_at, r.reported_at
		FROM campaign_recipients r
		JOIN users u ON u.id = r.user_id
		WHERE r.campaign_id = $1
		ORDER BY u.display_name
	`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RecipientRow
	for rows.Next() {
		var r RecipientRow
		if err := rows.Scan(
			&r.RecipientID, &r.UserID, &r.Email, &r.DisplayName,
			&r.DeliveredAt, &r.OpenedAt, &r.ClickedAt, &r.ReportedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// MarkDelivered is idempotent: a second call for the same recipient is a no-op.
func (s *Store) MarkDelivered(ctx context.Context, recipientID string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE campaign_recipients
		SET delivered_at = now()
		WHERE id = $1 AND delivered_at IS NULL
	`, recipientID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (s *Store) GetRecipient(ctx context.Context, recipientID string) (Recipient, error) {
	var r Recipient
	err := s.pool.QueryRow(ctx, `
		SELECT id, campaign_id, user_id, click_token, delivered_at, opened_at, clicked_at, reported_at
		FROM campaign_recipients WHERE id = $1
	`, recipientID).Scan(
		&r.ID, &r.CampaignID, &r.UserID, &r.ClickToken,
		&r.DeliveredAt, &r.OpenedAt, &r.ClickedAt, &r.ReportedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Recipient{}, ErrNotFound
	}
	return r, err
}

func (s *Store) ListInbox(ctx context.Context, userID string) ([]InboxItem, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT r.id, r.campaign_id, t.subject, t.body_text, t.lesson_text,
		       'IT Security Training <noreply@training.phishguard.local>',
		       r.delivered_at, r.opened_at, r.clicked_at, r.reported_at, r.click_token
		FROM campaign_recipients r
		JOIN campaigns c ON c.id = r.campaign_id
		JOIN templates t ON t.id = c.template_id
		WHERE r.user_id = $1 AND r.delivered_at IS NOT NULL
		ORDER BY r.delivered_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []InboxItem
	for rows.Next() {
		var item InboxItem
		if err := rows.Scan(
			&item.RecipientID, &item.CampaignID, &item.Subject, &item.BodyText, &item.LessonText,
			&item.FromLabel, &item.DeliveredAt, &item.OpenedAt, &item.ClickedAt, &item.ReportedAt, &item.ClickToken,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) GetInboxItem(ctx context.Context, recipientID, userID string) (InboxItem, error) {
	var item InboxItem
	err := s.pool.QueryRow(ctx, `
		SELECT r.id, r.campaign_id, t.subject, t.body_text, t.lesson_text,
		       'IT Security Training <noreply@training.phishguard.local>',
		       r.delivered_at, r.opened_at, r.clicked_at, r.reported_at, r.click_token
		FROM campaign_recipients r
		JOIN campaigns c ON c.id = r.campaign_id
		JOIN templates t ON t.id = c.template_id
		WHERE r.id = $1 AND r.user_id = $2 AND r.delivered_at IS NOT NULL
	`, recipientID, userID).Scan(
		&item.RecipientID, &item.CampaignID, &item.Subject, &item.BodyText, &item.LessonText,
		&item.FromLabel, &item.DeliveredAt, &item.OpenedAt, &item.ClickedAt, &item.ReportedAt, &item.ClickToken,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return InboxItem{}, ErrNotFound
	}
	return item, err
}

// MarkOpened records the first open. Returns (true, nil) when newly opened,
// (false, nil) when already opened, or ErrNotFound when the inbox row is missing.
// Unit-tested via campaign fakes; requires live Postgres for integration coverage.
func (s *Store) MarkOpened(ctx context.Context, recipientID, userID string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE campaign_recipients
		SET opened_at = now()
		WHERE id = $1 AND user_id = $2 AND delivered_at IS NOT NULL AND opened_at IS NULL
	`, recipientID, userID)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 1 {
		return true, nil
	}
	var exists bool
	err = s.pool.QueryRow(ctx, `
		SELECT true FROM campaign_recipients
		WHERE id = $1 AND user_id = $2 AND delivered_at IS NOT NULL
	`, recipientID, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

// MarkReported records the first report. Returns (true, nil) when newly reported.
// Unit-tested via campaign fakes; requires live Postgres for integration coverage.
func (s *Store) MarkReported(ctx context.Context, recipientID, userID string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE campaign_recipients
		SET reported_at = now(),
		    opened_at = COALESCE(opened_at, now())
		WHERE id = $1 AND user_id = $2 AND delivered_at IS NOT NULL AND reported_at IS NULL
	`, recipientID, userID)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 1 {
		return true, nil
	}
	var exists bool
	err = s.pool.QueryRow(ctx, `
		SELECT true FROM campaign_recipients
		WHERE id = $1 AND user_id = $2 AND delivered_at IS NOT NULL
	`, recipientID, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

// MarkClickedByToken records a click. first is true only on the first click for the token.
// Unit-tested via campaign fakes; requires live Postgres for integration coverage.
func (s *Store) MarkClickedByToken(ctx context.Context, token string) (Recipient, Template, bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Recipient{}, Template{}, false, err
	}
	defer tx.Rollback(ctx)

	var r Recipient
	var t Template
	err = tx.QueryRow(ctx, `
		SELECT r.id, r.campaign_id, r.user_id, r.click_token,
		       r.delivered_at, r.opened_at, r.clicked_at, r.reported_at,
		       t.id, t.org_id, t.name, t.subject, t.body_text, t.lesson_text, t.created_at
		FROM campaign_recipients r
		JOIN campaigns c ON c.id = r.campaign_id
		JOIN templates t ON t.id = c.template_id
		WHERE r.click_token = $1 AND r.delivered_at IS NOT NULL
	`, token).Scan(
		&r.ID, &r.CampaignID, &r.UserID, &r.ClickToken,
		&r.DeliveredAt, &r.OpenedAt, &r.ClickedAt, &r.ReportedAt,
		&t.ID, &t.OrgID, &t.Name, &t.Subject, &t.BodyText, &t.LessonText, &t.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Recipient{}, Template{}, false, ErrNotFound
	}
	if err != nil {
		return Recipient{}, Template{}, false, err
	}

	first := r.ClickedAt == nil
	if _, err := tx.Exec(ctx, `
		UPDATE campaign_recipients
		SET clicked_at = COALESCE(clicked_at, now()),
		    opened_at = COALESCE(opened_at, now())
		WHERE id = $1
	`, r.ID); err != nil {
		return Recipient{}, Template{}, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Recipient{}, Template{}, false, err
	}

	now := time.Now().UTC()
	if r.ClickedAt == nil {
		r.ClickedAt = &now
	}
	if r.OpenedAt == nil {
		r.OpenedAt = &now
	}
	return r, t, first, nil
}

func (s *Store) InsertEvent(ctx context.Context, orgID, campaignID, userID, kind, detail string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO events (org_id, campaign_id, user_id, kind, detail)
		VALUES ($1, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, $4, $5)
	`, orgID, campaignID, userID, kind, detail)
	return err
}

func (s *Store) ListUndeliveredRecipientIDs(ctx context.Context, campaignID string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id FROM campaign_recipients
		WHERE campaign_id = $1 AND delivered_at IS NULL
	`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Store) ListEvents(ctx context.Context, campaignID, orgID string) ([]Event, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT e.id, e.org_id, COALESCE(e.campaign_id::text, ''), COALESCE(e.user_id::text, ''),
		       e.kind, e.detail, e.created_at
		FROM events e
		JOIN campaigns c ON c.id = e.campaign_id
		WHERE e.campaign_id = $1 AND c.org_id = $2
		ORDER BY e.created_at ASC
		LIMIT 200
	`, campaignID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.OrgID, &e.CampaignID, &e.UserID, &e.Kind, &e.Detail, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	if out == nil {
		out = []Event{}
	}
	return out, rows.Err()
}

func (s *Store) CloseCampaign(ctx context.Context, id, orgID string) (Campaign, error) {
	var c Campaign
	err := s.pool.QueryRow(ctx, `
		UPDATE campaigns
		SET status = 'closed'
		WHERE id = $1 AND org_id = $2 AND status IN ('active', 'launching')
		RETURNING id, org_id, created_by, template_id, name, status, created_at, launched_at
	`, id, orgID).Scan(&c.ID, &c.OrgID, &c.CreatedBy, &c.TemplateID, &c.Name, &c.Status, &c.CreatedAt, &c.LaunchedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		existing, gerr := s.GetCampaign(ctx, id, orgID)
		if gerr != nil {
			return Campaign{}, gerr
		}
		if existing.Status == StatusClosed {
			return existing, nil
		}
		return Campaign{}, ErrConflict
	}
	return c, err
}

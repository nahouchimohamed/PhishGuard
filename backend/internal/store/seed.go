package store

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// SeedIfEmpty creates one org, four users, and three training templates when the DB is empty.
func (s *Store) SeedIfEmpty(ctx context.Context, password string) error {
	n, err := s.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}
	passwordHash := string(hash)

	orgID, err := s.CreateOrg(ctx, "Acme Training Co")
	if err != nil {
		return err
	}

	users := []struct {
		email, name, role string
	}{
		{"admin@phishguard.local", "Alex Admin", string(RoleAdmin)},
		{"alice@phishguard.local", "Alice Employee", string(RoleEmployee)},
		{"bob@phishguard.local", "Bob Employee", string(RoleEmployee)},
		{"cara@phishguard.local", "Cara Employee", string(RoleEmployee)},
	}
	for _, u := range users {
		if _, err := s.CreateUser(ctx, orgID, u.email, u.name, u.role, passwordHash); err != nil {
			return fmt.Errorf("seed user %s: %w", u.email, err)
		}
	}

	templates := []struct {
		name, subject, body, lesson string
	}{
		{
			name:    "Password expiry",
			subject: "[Training] Action required: your password expires today",
			body: `Hello,

Our security system detected that your password will expire in 2 hours.
Click the link below to renew it immediately and avoid account lockout.

Renew password: {{LINK}}

— IT Helpdesk`,
			lesson: `This was a simulated phishing message.

Red flags:
- Urgent deadline ("2 hours", "immediately")
- Generic greeting ("Hello") instead of your name
- Link that does not match the claimed company domain
- Pressure to click before thinking

Good habit: open the official IT portal yourself; never follow links from unexpected expiry warnings.`,
		},
		{
			name:    "Shared document",
			subject: "[Training] You have been granted access to Q4 payroll.xlsx",
			body: `Hi,

A confidential spreadsheet was shared with you.

Open document: {{LINK}}

Please review before end of day.

— Finance Ops`,
			lesson: `This was a simulated phishing message.

Red flags:
- Unexpected "shared document" about payroll (high-value bait)
- Vague sender ("Finance Ops") with no ticket number
- Link text that hides the real destination
- Request to act "before end of day"

Good habit: verify unexpected file shares in your real document system, not from email links.`,
		},
		{
			name:    "Package delivery",
			subject: "[Training] Delivery attempt failed — reschedule package #88421",
			body: `Dear customer,

We could not deliver your package today.

Reschedule delivery: {{LINK}}

Fees may apply if you do not respond within 24 hours.

— Shipping Desk`,
			lesson: `This was a simulated phishing message.

Red flags:
- Threat of fees or penalties
- Package number you did not expect
- Unknown "Shipping Desk" address style
- Single click to "reschedule"

Good habit: track packages only through the courier you already ordered from.`,
		},
	}

	for _, t := range templates {
		if _, err := s.CreateTemplate(ctx, orgID, t.name, t.subject, t.body, t.lesson); err != nil {
			return fmt.Errorf("seed template %s: %w", t.name, err)
		}
	}

	slog.Info("database seeded", "org_id", orgID, "users", len(users), "templates", len(templates))
	return nil
}

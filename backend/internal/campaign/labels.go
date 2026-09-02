package campaign

import "time"

// RecipientLabel turns timestamps into a single UI status.
// Report and click can both exist; report wins the label because that is the desired training outcome.
// Undelivered recipients are Pending (matches apps/web/src/lib/labels.ts).
func RecipientLabel(deliveredAt, openedAt, clickedAt, reportedAt *time.Time) string {
	switch {
	case reportedAt != nil:
		return "Reported"
	case clickedAt != nil:
		return "Clicked"
	case openedAt != nil:
		return "Opened"
	case deliveredAt != nil:
		return "Delivered"
	default:
		return "Pending"
	}
}

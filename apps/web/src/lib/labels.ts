/** Mirrors backend/internal/campaign.RecipientLabel, including Pending. */
export type RecipientLabel = 'Reported' | 'Clicked' | 'Opened' | 'Delivered' | 'Pending';

export function recipientLabel(item: {
	opened_at?: string | null;
	clicked_at?: string | null;
	reported_at?: string | null;
	delivered_at?: string | null;
}): RecipientLabel {
	if (item.reported_at) return 'Reported';
	if (item.clicked_at) return 'Clicked';
	if (item.opened_at) return 'Opened';
	if (item.delivered_at) return 'Delivered';
	return 'Pending';
}

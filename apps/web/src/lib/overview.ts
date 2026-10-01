import type { Campaign } from '$lib/types';

export function recentCampaigns(rows: Campaign[], limit = 5): Campaign[] {
	return [...rows].sort((a, b) => b.created_at.localeCompare(a.created_at)).slice(0, limit);
}

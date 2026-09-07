import type { Campaign } from '$lib/types';

export type CampaignFilter = 'all' | 'active' | 'closed' | 'launching' | 'draft';

export function filterCampaigns(rows: Campaign[], filter: CampaignFilter): Campaign[] {
	if (filter === 'all') return rows;
	return rows.filter((c) => c.status === filter);
}

export function searchCampaigns(rows: Campaign[], q: string): Campaign[] {
	const needle = q.trim().toLowerCase();
	if (!needle) return rows;
	return rows.filter((c) => c.name.toLowerCase().includes(needle));
}

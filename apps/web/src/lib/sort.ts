import type { Campaign } from '$lib/types';
import { reportRate } from '$lib/stats';

export type CampaignSort = 'newest' | 'report' | 'name';

export function sortCampaigns(rows: Campaign[], sort: CampaignSort): Campaign[] {
	const copy = [...rows];
	switch (sort) {
		case 'name':
			return copy.sort((a, b) => a.name.localeCompare(b.name));
		case 'report':
			return copy.sort((a, b) => {
				const ar = a.stats ? reportRate(a.stats) : -1;
				const br = b.stats ? reportRate(b.stats) : -1;
				return br - ar;
			});
		case 'newest':
		default:
			return copy.sort((a, b) => b.created_at.localeCompare(a.created_at));
	}
}

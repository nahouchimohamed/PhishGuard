import type { PageServerLoad } from './$types';
import { aggregateStats } from '$lib/stats';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { Campaign } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);
	const res = await apiFetch<Campaign[]>('/api/campaigns', { cookies });
	if (!res.ok || !res.data) {
		return {
			campaigns: [],
			summary: aggregateStats([]),
			error: friendlyError(res.status, res.error)
		};
	}
	const statsRows = res.data.map(
		(c) => c.stats ?? { total: 0, delivered: 0, opened: 0, clicked: 0, reported: 0 }
	);
	return { campaigns: res.data, summary: aggregateStats(statsRows), error: null };
};

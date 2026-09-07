import type { PageServerLoad } from './$types';
import type { CampaignFilter } from '$lib/filters';
import { filterCampaigns, searchCampaigns } from '$lib/filters';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { Campaign } from '$lib/types';

const allowed: CampaignFilter[] = ['all', 'active', 'closed', 'launching', 'draft'];

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);
	const raw = url.searchParams.get('status') ?? 'all';
	const statusFilter = (allowed.includes(raw as CampaignFilter) ? raw : 'all') as CampaignFilter;
	const q = url.searchParams.get('q') ?? '';
	const res = await apiFetch<Campaign[]>('/api/campaigns', { cookies });
	if (!res.ok || !res.data) {
		return { campaigns: [], statusFilter, q, error: friendlyError(res.status, res.error) };
	}
	const filtered = searchCampaigns(filterCampaigns(res.data, statusFilter), q);
	return { campaigns: filtered, statusFilter, q, error: null };
};

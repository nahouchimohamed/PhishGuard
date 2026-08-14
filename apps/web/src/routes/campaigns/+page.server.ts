import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { Campaign } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);

	const res = await apiFetch<Campaign[]>('/api/campaigns', { cookies });
	if (!res.ok) {
		return { campaigns: [] as Campaign[], error: friendlyError(res.status, res.error) };
	}
	return { campaigns: res.data ?? [], error: null };
};

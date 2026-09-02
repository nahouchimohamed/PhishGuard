import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { Template } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);
	const res = await apiFetch<Template[]>('/api/templates', { cookies });
	if (!res.ok || !res.data) {
		return { templates: [], error: friendlyError(res.status, res.error) };
	}
	return { templates: res.data, error: null };
};

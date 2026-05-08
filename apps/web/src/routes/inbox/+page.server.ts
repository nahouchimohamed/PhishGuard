import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { InboxItem } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'employee', url.pathname);

	const res = await apiFetch<InboxItem[]>('/api/inbox', { cookies });
	if (!res.ok) {
		return { items: [] as InboxItem[], error: friendlyError(res.status, res.error) };
	}
	return { items: res.data ?? [], error: null };
};

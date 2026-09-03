import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { User } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);
	const res = await apiFetch<User[]>('/api/employees', { cookies });
	if (!res.ok || !res.data) {
		return { employees: [], error: friendlyError(res.status, res.error) };
	}
	return { employees: res.data, error: null };
};

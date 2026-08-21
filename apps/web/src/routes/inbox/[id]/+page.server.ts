import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch, apiPath } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { InboxItem } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, params, url }) => {
	requireRole(locals.user, 'employee', url.pathname);

	const res = await apiFetch<InboxItem>(apiPath('/api/inbox', params.id), { cookies });
	if (!res.ok || !res.data) {
		return { item: null, error: friendlyError(res.status, res.error) };
	}
	return { item: res.data, error: null };
};

export const actions: Actions = {
	report: async ({ locals, cookies, params, url }) => {
		requireRole(locals.user, 'employee', url.pathname);
		const res = await apiFetch(apiPath('/api/inbox', params.id, 'report'), {
			method: 'POST',
			cookies
		});
		if (!res.ok) return fail(res.status || 400, { error: friendlyError(res.status, res.error) });
		throw redirect(303, `/inbox/${params.id}`);
	}
};

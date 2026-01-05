import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import type { InboxItem } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, params }) => {
	if (!locals.user) throw redirect(303, '/login');
	if (locals.user.role !== 'employee') throw redirect(303, '/campaigns');

	const res = await apiFetch<InboxItem>(`/api/inbox/${params.id}`, { cookies });
	if (!res.ok || !res.data) {
		return { item: null, error: res.error ?? 'Not found' };
	}
	return { item: res.data, error: null };
};

export const actions: Actions = {
	report: async ({ locals, cookies, params }) => {
		if (!locals.user || locals.user.role !== 'employee') throw redirect(303, '/login');
		const res = await apiFetch(`/api/inbox/${params.id}/report`, { method: 'POST', cookies });
		if (!res.ok) return fail(res.status, { error: res.error ?? 'Report failed' });
		throw redirect(303, `/inbox/${params.id}`);
	}
};

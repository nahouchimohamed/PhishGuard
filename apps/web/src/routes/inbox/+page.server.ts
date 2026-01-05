import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import type { InboxItem } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.user) throw redirect(303, '/login');
	if (locals.user.role !== 'employee') throw redirect(303, '/campaigns');

	const res = await apiFetch<InboxItem[]>('/api/inbox', { cookies });
	return { items: res.data ?? [], error: res.error };
};

import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import type { Campaign } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.user) throw redirect(303, '/login');
	if (locals.user.role !== 'admin') throw redirect(303, '/inbox');

	const res = await apiFetch<Campaign[]>('/api/campaigns', { cookies });
	if (!res.ok) {
		return { campaigns: [] as Campaign[], error: res.error };
	}
	return { campaigns: res.data ?? [], error: null };
};

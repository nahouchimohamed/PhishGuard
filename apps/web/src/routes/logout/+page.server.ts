import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch, COOKIE } from '$lib/server/api';

export const load: PageServerLoad = async () => {
	throw redirect(303, '/login');
};

export const actions: Actions = {
	default: async ({ cookies }) => {
		try {
			await apiFetch('/api/logout', { method: 'POST', cookies });
		} catch {
			// Always clear the local session even if the API is unreachable.
		}
		cookies.delete(COOKIE, { path: '/' });
		throw redirect(303, '/login');
	}
};

import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { apiFetch, COOKIE } from '$lib/server/api';

export const actions: Actions = {
	default: async ({ cookies }) => {
		await apiFetch('/api/logout', { method: 'POST', cookies });
		cookies.delete(COOKIE, { path: '/' });
		throw redirect(303, '/login');
	}
};

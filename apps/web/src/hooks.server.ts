import type { Handle } from '@sveltejs/kit';
import { apiFetch, COOKIE } from '$lib/server/api';
import type { User } from '$lib/types';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get(COOKIE);
	if (token) {
		const me = await apiFetch<User>('/api/me', { cookies: event.cookies });
		event.locals.user = me.ok ? me.data : null;
		if (!me.ok) {
			event.cookies.delete(COOKIE, { path: '/' });
		}
	} else {
		event.locals.user = null;
	}
	return resolve(event);
};

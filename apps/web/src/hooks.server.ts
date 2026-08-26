import { sequence } from '@sveltejs/kit/hooks';
import type { Handle, HandleServerError } from '@sveltejs/kit';
import { apiFetch } from '$lib/server/api';
import { clearSessionCookie, COOKIE } from '$lib/server/session';
import { securityHeaders } from '$lib/server/headers';
import type { User } from '$lib/types';
import { randomUUID } from 'node:crypto';

const sessionHandle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get(COOKIE);
	if (token) {
		const me = await apiFetch<User>('/api/me', { cookies: event.cookies });
		event.locals.user = me.ok ? me.data : null;
		if (!me.ok && me.status === 401) {
			clearSessionCookie(event.cookies);
		}
	} else {
		event.locals.user = null;
	}
	return resolve(event);
};

export const handle = sequence(sessionHandle, securityHeaders);

export const handleError: HandleServerError = ({ error }) => {
	const id = randomUUID();
	console.error('app error', id, error);
	return {
		message: 'Something went wrong.',
		id
	};
};

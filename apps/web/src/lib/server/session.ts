import type { Cookies } from '@sveltejs/kit';

export const COOKIE = 'phishguard_session';
export const SESSION_MAX_AGE_SECONDS = 60 * 60 * 24 * 7;

export function setSessionCookie(cookies: Cookies, token: string) {
	cookies.set(COOKIE, token, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		maxAge: SESSION_MAX_AGE_SECONDS
	});
}

export function clearSessionCookie(cookies: Cookies) {
	cookies.delete(COOKIE, { path: '/' });
}

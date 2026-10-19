import type { Handle } from '@sveltejs/kit';

/** Baseline browser hardening headers for the local training UI. */
export const securityHeaders: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);
	response.headers.set('X-Content-Type-Options', 'nosniff');
	response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');
	response.headers.set('X-Frame-Options', 'DENY');
	if (event.locals.user) {
		response.headers.set('Cache-Control', 'no-store');
	}
	return response;
};

import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { setSessionCookie } from '$lib/server/session';
import { homePathFor } from '$lib/roles';
import { safeRedirectPath } from '$lib/redirect';
import type { LoginResponse } from '$lib/types';

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.user) throw redirect(303, homePathFor(locals.user.role));
	return {};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const form = await request.formData();
		const email = String(form.get('email') ?? '').trim().toLowerCase();
		const password = String(form.get('password') ?? '');
		const redirectTo = safeRedirectPath(form.get('redirectTo'), '');

		const res = await apiFetch<LoginResponse>('/api/login', {
			method: 'POST',
			body: { email, password }
		});
		if (!res.ok || !res.data?.session_token || !res.data.user) {
			const status = res.status >= 400 ? res.status : 400;
			return fail(status, {
				error: friendlyError(res.status, res.error),
				email
			});
		}

		setSessionCookie(cookies, res.data.session_token);
		throw redirect(303, redirectTo || homePathFor(res.data.user.role));
	}
};

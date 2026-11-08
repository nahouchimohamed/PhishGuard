import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch, COOKIE } from '$lib/server/api';
import type { User } from '$lib/types';

type LoginResponse = {
	user: User;
	session_token: string;
};

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.user?.role === 'admin') throw redirect(303, '/campaigns');
	if (locals.user?.role === 'employee') throw redirect(303, '/inbox');
	return {};
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const form = await request.formData();
		const email = String(form.get('email') ?? '');
		const password = String(form.get('password') ?? '');

		const res = await apiFetch<LoginResponse>('/api/login', {
			method: 'POST',
			body: { email, password }
		});
		if (!res.ok || !res.data?.session_token || !res.data.user) {
			return fail(400, { error: res.error ?? 'Login failed', email });
		}

		cookies.set(COOKIE, res.data.session_token, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 60 * 60 * 24 * 7
		});

		if (res.data.user.role === 'admin') throw redirect(303, '/campaigns');
		throw redirect(303, '/inbox');
	}
};

import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import type { Campaign, Template, User } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies }) => {
	if (!locals.user) throw redirect(303, '/login');
	if (locals.user.role !== 'admin') throw redirect(303, '/inbox');

	const [templates, employees] = await Promise.all([
		apiFetch<Template[]>('/api/templates', { cookies }),
		apiFetch<User[]>('/api/employees', { cookies })
	]);

	return {
		templates: templates.data ?? [],
		employees: employees.data ?? [],
		error: templates.error || employees.error || null
	};
};

export const actions: Actions = {
	default: async ({ request, cookies, locals }) => {
		if (!locals.user || locals.user.role !== 'admin') throw redirect(303, '/login');

		const form = await request.formData();
		const name = String(form.get('name') ?? '').trim();
		const template_id = String(form.get('template_id') ?? '');
		const recipient_user_ids = form.getAll('recipient_user_ids').map(String);

		const values = { name, template_id, recipient_user_ids };

		if (!name || !template_id || recipient_user_ids.length === 0) {
			return fail(400, {
				error: 'Name, template, and at least one employee are required.',
				values
			});
		}

		const res = await apiFetch<Campaign | { campaign: Campaign; warning?: string }>('/api/campaigns', {
			method: 'POST',
			body: { name, template_id, recipient_user_ids },
			cookies
		});

		if (!res.ok && res.status !== 202) {
			return fail(res.status || 400, { error: res.error ?? 'Launch failed', values });
		}

		const camp =
			res.data && 'id' in (res.data as Campaign)
				? (res.data as Campaign)
				: (res.data as { campaign: Campaign }).campaign;

		throw redirect(303, `/campaigns/${camp.id}`);
	}
};

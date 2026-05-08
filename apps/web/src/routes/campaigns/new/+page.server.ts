import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import { parseLaunchResponse } from '$lib/server/launch';
import { validateCampaignForm } from '$lib/validate';
import type { Campaign, LaunchResponse, Template, User } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	requireRole(locals.user, 'admin', url.pathname);

	const [templates, employees] = await Promise.all([
		apiFetch<Template[]>('/api/templates', { cookies }),
		apiFetch<User[]>('/api/employees', { cookies })
	]);

	return {
		templates: templates.data ?? [],
		employees: employees.data ?? [],
		error:
			templates.error || employees.error
				? friendlyError(templates.status || employees.status, templates.error || employees.error)
				: null
	};
};

export const actions: Actions = {
	default: async ({ request, cookies, locals, url }) => {
		requireRole(locals.user, 'admin', url.pathname);

		const form = await request.formData();
		const validated = validateCampaignForm({
			name: form.get('name'),
			template_id: form.get('template_id'),
			recipient_user_ids: form.getAll('recipient_user_ids')
		});

		if (!validated.ok) {
			return fail(400, {
				error: 'Fix the highlighted fields and try again.',
				errors: validated.errors,
				values: validated.values
			});
		}

		const { name, template_id, recipient_user_ids } = validated.values;
		const res = await apiFetch<Campaign | LaunchResponse>('/api/campaigns', {
			method: 'POST',
			body: { name, template_id, recipient_user_ids },
			cookies
		});

		if (!res.ok && res.status !== 202) {
			return fail(res.status || 400, {
				error: friendlyError(res.status, res.error),
				values: validated.values
			});
		}

		const parsed = parseLaunchResponse(res.status, res.data);
		if (!parsed.ok) {
			return fail(502, { error: parsed.error, values: validated.values });
		}

		const notice = parsed.incomplete ? '?notice=launch-incomplete' : '';
		throw redirect(303, `/campaigns/${parsed.campaign.id}${notice}`);
	}
};

import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';
import type { Campaign, CampaignStats, RecipientRow, Template } from '$lib/types';

type Detail = {
	campaign: Campaign;
	template: Template;
	stats: CampaignStats;
	recipients: RecipientRow[];
};

export const load: PageServerLoad = async ({ locals, cookies, params }) => {
	if (!locals.user) throw redirect(303, '/login');
	if (locals.user.role !== 'admin') throw redirect(303, '/inbox');

	const res = await apiFetch<Detail>(`/api/campaigns/${params.id}`, { cookies });
	if (!res.ok || !res.data) {
		return { detail: null, error: res.error ?? 'Not found' };
	}
	return { detail: res.data, error: null };
};

export const actions: Actions = {
	retry: async ({ locals, cookies, params }) => {
		if (!locals.user || locals.user.role !== 'admin') throw redirect(303, '/login');
		const res = await apiFetch(`/api/campaigns/${params.id}/retry`, {
			method: 'POST',
			cookies
		});
		if (!res.ok) return fail(res.status, { error: res.error ?? 'Retry failed' });
		throw redirect(303, `/campaigns/${params.id}`);
	}
};

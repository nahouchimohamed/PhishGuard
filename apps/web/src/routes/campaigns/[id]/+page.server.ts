import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { apiFetch, apiPath } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import { requireRole } from '$lib/server/guard';
import type { CampaignDetail } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, params, url }) => {
	requireRole(locals.user, 'admin', url.pathname);

	const res = await apiFetch<CampaignDetail>(apiPath('/api/campaigns', params.id), { cookies });
	if (!res.ok || !res.data) {
		return { detail: null, error: friendlyError(res.status, res.error), notice: null };
	}
	const notice = url.searchParams.get('notice');
	return {
		detail: res.data,
		error: null,
		notice: notice === 'launch-incomplete' ? 'Campaign saved, but queue publish was incomplete.' : null
	};
};

export const actions: Actions = {
	retry: async ({ locals, cookies, params, url }) => {
		requireRole(locals.user, 'admin', url.pathname);
		const res = await apiFetch(apiPath('/api/campaigns', params.id, 'retry'), {
			method: 'POST',
			cookies
		});
		if (!res.ok) return fail(res.status || 400, { error: friendlyError(res.status, res.error) });
		throw redirect(303, `/campaigns/${params.id}`);
	}
};

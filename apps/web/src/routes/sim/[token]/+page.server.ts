import type { PageServerLoad } from './$types';
import { apiFetch, apiPath } from '$lib/server/api';
import { friendlyError } from '$lib/server/errors';
import type { SimLesson } from '$lib/types';

export const load: PageServerLoad = async ({ params }) => {
	const res = await apiFetch<SimLesson>(apiPath('/api/sim', params.token));
	if (!res.ok || !res.data) {
		return {
			lesson: null,
			error: friendlyError(res.status, res.error || 'This training link is invalid or expired.')
		};
	}
	return { lesson: res.data, error: null };
};

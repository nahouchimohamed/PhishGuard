import type { PageServerLoad } from './$types';
import { apiFetch } from '$lib/server/api';

type SimResult = {
	campaign_id: string;
	subject: string;
	lesson_text: string;
	name: string;
};

export const load: PageServerLoad = async ({ params }) => {
	const res = await apiFetch<SimResult>(`/api/sim/${params.token}`);
	if (!res.ok || !res.data) {
		return { lesson: null, error: res.error ?? 'This training link is invalid or expired.' };
	}
	return { lesson: res.data, error: null };
};

import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { homePathFor } from '$lib/roles';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(303, '/login');
	}
	throw redirect(303, homePathFor(locals.user.role));
};

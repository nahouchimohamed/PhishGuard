import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/private';
import { requireRole } from '$lib/server/guard';
import { COOKIE } from '$lib/server/session';

export const GET: RequestHandler = async ({ locals, cookies, params, url }) => {
	requireRole(locals.user, 'admin', url.pathname);
	const base = env.API_BASE_URL || 'http://localhost:8080';
	const token = cookies.get(COOKIE);
	const res = await fetch(`${base}/api/campaigns/${encodeURIComponent(params.id)}/export`, {
		headers: {
			Accept: 'text/csv',
			...(token ? { Cookie: `${COOKIE}=${token}` } : {})
		},
		signal: AbortSignal.timeout(5000)
	});
	if (!res.ok) throw error(res.status, 'Export failed');
	const body = await res.text();
	return new Response(body, {
		headers: {
			'Content-Type': 'text/csv; charset=utf-8',
			'Content-Disposition': `attachment; filename="campaign-${params.id}.csv"`
		}
	});
};

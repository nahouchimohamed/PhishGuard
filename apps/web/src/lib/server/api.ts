import { env } from '$env/dynamic/private';
import type { Cookies } from '@sveltejs/kit';
import { COOKIE } from '$lib/server/session';

function baseURL(): string {
	return env.API_BASE_URL || 'http://localhost:8080';
}

export type ApiResult<T> = {
	ok: boolean;
	status: number;
	data: T | null;
	error?: string;
};

export function apiPath(base: string, ...parts: string[]): string {
	const encoded = parts.map((p) => encodeURIComponent(p));
	return [base.replace(/\/$/, ''), ...encoded].join('/');
}

export async function apiFetch<T>(
	path: string,
	opts: {
		method?: string;
		body?: unknown;
		cookies?: Cookies;
	} = {}
): Promise<ApiResult<T>> {
	const headers: Record<string, string> = {
		Accept: 'application/json'
	};
	if (opts.body !== undefined) {
		headers['Content-Type'] = 'application/json';
	}

	const cookie = opts.cookies?.get(COOKIE) ?? null;
	if (cookie) {
		headers.Cookie = `${COOKIE}=${cookie}`;
	}

	let res: Response;
	try {
		res = await fetch(`${baseURL()}${path}`, {
			method: opts.method ?? (opts.body !== undefined ? 'POST' : 'GET'),
			headers,
			body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
			signal: AbortSignal.timeout(5000)
		});
	} catch {
		return { ok: false, status: 503, data: null, error: 'API unreachable' };
	}

	let data: T | null = null;
	let error: string | undefined;
	const text = await res.text();
	if (text) {
		try {
			const parsed = JSON.parse(text);
			if (!res.ok) {
				error = parsed.error ?? text;
			} else {
				data = parsed as T;
			}
		} catch {
			if (!res.ok) error = text;
			else error = 'Invalid JSON from API';
		}
	}

	return { ok: res.ok, status: res.status, data, error };
}

export { COOKIE };

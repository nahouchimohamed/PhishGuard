import { env } from '$env/dynamic/private';
import type { Cookies } from '@sveltejs/kit';

const COOKIE = 'phishguard_session';

function baseURL(): string {
	return env.API_BASE_URL || 'http://localhost:8080';
}

export type ApiResult<T> = {
	ok: boolean;
	status: number;
	data: T | null;
	error?: string;
	setCookie?: string | null;
};

export async function apiFetch<T>(
	path: string,
	opts: {
		method?: string;
		body?: unknown;
		cookies?: Cookies;
		cookieHeader?: string | null;
	} = {}
): Promise<ApiResult<T>> {
	const headers: Record<string, string> = {
		Accept: 'application/json'
	};
	if (opts.body !== undefined) {
		headers['Content-Type'] = 'application/json';
	}

	const cookie = opts.cookieHeader ?? opts.cookies?.get(COOKIE) ?? null;
	if (cookie) {
		headers.Cookie = `${COOKIE}=${cookie}`;
	}

	const res = await fetch(`${baseURL()}${path}`, {
		method: opts.method ?? (opts.body !== undefined ? 'POST' : 'GET'),
		headers,
		body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined
	});

	const setCookie = res.headers.getSetCookie?.() ?? [];
	let rawSet: string | null = null;
	for (const line of setCookie) {
		if (line.startsWith(`${COOKIE}=`)) {
			rawSet = line;
			break;
		}
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
		}
	}

	return { ok: res.ok, status: res.status, data, error, setCookie: rawSet };
}

export function applySessionCookie(cookies: Cookies, setCookie: string | null | undefined) {
	if (!setCookie) return;
	// Parse Max-Age / empty value from Go Set-Cookie.
	const parts = setCookie.split(';').map((p) => p.trim());
	const [kv] = parts;
	const eq = kv.indexOf('=');
	const value = eq >= 0 ? kv.slice(eq + 1) : '';
	const maxAgePart = parts.find((p) => p.toLowerCase().startsWith('max-age='));
	const maxAge = maxAgePart ? Number(maxAgePart.split('=')[1]) : undefined;

	if (!value || maxAge === -1) {
		cookies.delete(COOKIE, { path: '/' });
		return;
	}
	cookies.set(COOKIE, value, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		maxAge: Number.isFinite(maxAge) ? maxAge : 60 * 60 * 24 * 7
	});
}

export { COOKIE };

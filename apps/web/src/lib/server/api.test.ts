import { beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('$env/dynamic/private', () => ({
	env: { API_BASE_URL: 'http://api.test' }
}));

import { apiFetch, apiPath } from './api';

describe('apiPath', () => {
	it('encodes path segments', () => {
		expect(apiPath('/api/sim', 'a/b')).toBe('/api/sim/a%2Fb');
		expect(apiPath('/api/campaigns', 'id', 'retry')).toBe('/api/campaigns/id/retry');
	});
});

describe('apiFetch', () => {
	beforeEach(() => {
		vi.unstubAllGlobals();
	});

	it('posts JSON and parses success bodies', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () =>
				new Response(JSON.stringify({ ok: true }), {
					status: 200,
					headers: { 'Content-Type': 'application/json' }
				})
			)
		);
		const res = await apiFetch<{ ok: boolean }>('/api/x', { method: 'POST', body: { a: 1 } });
		expect(res.ok).toBe(true);
		expect(res.data).toEqual({ ok: true });
		expect(fetch).toHaveBeenCalledWith(
			'http://api.test/api/x',
			expect.objectContaining({ method: 'POST' })
		);
	});

	it('returns 503 when fetch throws', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => {
				throw new Error('offline');
			})
		);
		const res = await apiFetch('/api/x');
		expect(res).toEqual({ ok: false, status: 503, data: null, error: 'API unreachable' });
	});

	it('maps JSON error bodies', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(
				async () =>
					new Response(JSON.stringify({ error: 'unauthorized' }), {
						status: 401,
						headers: { 'Content-Type': 'application/json' }
					})
			)
		);
		const res = await apiFetch('/api/me');
		expect(res.ok).toBe(false);
		expect(res.status).toBe(401);
		expect(res.error).toBe('unauthorized');
	});

	it('handles empty success bodies', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () => new Response(null, { status: 204 }))
		);
		const res = await apiFetch('/api/x', { method: 'DELETE' });
		expect(res.status).toBe(204);
		expect(res.ok).toBe(true);
		expect(res.data).toBeNull();
		expect(res.error).toBeUndefined();
	});
});

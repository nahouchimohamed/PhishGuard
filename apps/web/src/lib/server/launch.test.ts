import { describe, expect, it } from 'vitest';
import type { Campaign } from '$lib/types';
import { parseLaunchResponse } from './launch';

const camp = {
	id: 'c1',
	org_id: 'o1',
	created_by: 'u1',
	template_id: 't1',
	name: 'Drill',
	status: 'active',
	created_at: '2026-01-01T00:00:00Z'
};

describe('parseLaunchResponse', () => {
	it('handles 201 campaign bodies', () => {
		const res = parseLaunchResponse(201, camp);
		expect(res.ok).toBe(true);
		if (res.ok) {
			expect(res.campaign.id).toBe('c1');
			expect(res.incomplete).toBe(false);
		}
	});

	it('handles 202 wrapped bodies with warning', () => {
		const res = parseLaunchResponse(202, { campaign: camp, warning: 'partial' });
		expect(res.ok).toBe(true);
		if (res.ok) {
			expect(res.incomplete).toBe(true);
			expect(res.warning).toBe('partial');
		}
	});

	it('rejects null bodies', () => {
		expect(parseLaunchResponse(201, null).ok).toBe(false);
	});

	it('rejects unrecognized response shapes', () => {
		expect(parseLaunchResponse(201, { error: 'bad' } as unknown as Campaign).ok).toBe(false);
		expect(parseLaunchResponse(202, { warning: 'partial' } as unknown as Campaign).ok).toBe(false);
	});
});

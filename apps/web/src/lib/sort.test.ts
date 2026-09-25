import { describe, expect, it } from 'vitest';
import { sortCampaigns } from './sort';
import type { Campaign } from './types';

function camp(id: string, name: string, created_at: string, reported = 0, delivered = 1): Campaign {
	return {
		id,
		name,
		created_at,
		org_id: 'o',
		created_by: 'u',
		template_id: 't',
		status: 'active',
		stats: { total: delivered, delivered, opened: 0, clicked: 0, reported }
	};
}

describe('sortCampaigns', () => {
	it('orders by report rate descending', () => {
		const rows = [
			camp('a', 'A', '2026-08-01T00:00:00Z', 1, 4),
			camp('b', 'B', '2026-08-02T00:00:00Z', 3, 4)
		];
		expect(sortCampaigns(rows, 'report').map((c) => c.id)).toEqual(['b', 'a']);
	});
});

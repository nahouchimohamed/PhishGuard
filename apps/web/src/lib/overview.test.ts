import { describe, expect, it } from 'vitest';
import { recentCampaigns } from './overview';
import type { Campaign } from './types';

function camp(id: string, created_at: string): Campaign {
	return {
		id,
		created_at,
		name: id,
		org_id: 'o',
		created_by: 'u',
		template_id: 't',
		status: 'active'
	};
}

describe('recentCampaigns', () => {
	it('returns newest first up to the limit', () => {
		const rows = [camp('a', '2026-08-01T00:00:00Z'), camp('b', '2026-09-01T00:00:00Z')];
		expect(recentCampaigns(rows, 1).map((c) => c.id)).toEqual(['b']);
	});
});

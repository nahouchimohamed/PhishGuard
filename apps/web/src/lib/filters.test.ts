import { describe, expect, it } from 'vitest';
import { filterCampaigns, searchCampaigns } from './filters';
import type { Campaign } from './types';

const base = {
	org_id: 'o',
	created_by: 'u',
	template_id: 't',
	created_at: '2026-08-01T00:00:00Z'
} as const;

function camp(id: string, status: Campaign['status'], name = id): Campaign {
	return { id, status, name, ...base };
}

describe('campaign filters', () => {
	it('filters by status', () => {
		const rows = [camp('1', 'active'), camp('2', 'closed')];
		expect(filterCampaigns(rows, 'closed').map((c) => c.id)).toEqual(['2']);
	});

	it('searches by name', () => {
		const rows = [camp('1', 'active', 'Payroll'), camp('2', 'active', 'VPN')];
		expect(searchCampaigns(rows, 'pay').map((c) => c.id)).toEqual(['1']);
	});
});

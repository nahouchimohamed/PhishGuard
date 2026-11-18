import { describe, expect, it } from 'vitest';
import { countByLabel } from './recipients';

describe('countByLabel', () => {
	it('tallies outcomes', () => {
		expect(
			countByLabel([
				{ recipient_id: '1', user_id: 'u', email: 'a', display_name: 'A', reported_at: '1' },
				{ recipient_id: '2', user_id: 'u', email: 'b', display_name: 'B' }
			])
		).toMatchObject({ Reported: 1, Pending: 1 });
	});
});

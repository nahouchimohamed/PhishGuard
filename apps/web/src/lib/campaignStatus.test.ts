import { describe, expect, it } from 'vitest';
import { isCloseable } from './campaignStatus';

describe('isCloseable', () => {
	it('allows active and launching only', () => {
		expect(isCloseable('active')).toBe(true);
		expect(isCloseable('launching')).toBe(true);
		expect(isCloseable('closed')).toBe(false);
	});
});

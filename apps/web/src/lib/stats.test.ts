import { describe, expect, it } from 'vitest';
import { clickRate, percent, reportRate } from './stats';

describe('percent', () => {
	it('returns 0 for empty or invalid denominators', () => {
		expect(percent(1, 0)).toBe(0);
		expect(percent(1, -2)).toBe(0);
		expect(percent(Number.NaN, 10)).toBe(0);
	});

	it('rounds to one decimal place', () => {
		expect(percent(1, 3)).toBe(33.3);
		expect(percent(2, 3)).toBe(66.7);
	});

	it('returns 0 for non-finite numerators', () => {
		expect(percent(Number.POSITIVE_INFINITY, 10)).toBe(0);
		expect(percent(Number.NEGATIVE_INFINITY, 10)).toBe(0);
	});
});

describe('rates', () => {
	it('uses delivered as the denominator', () => {
		expect(reportRate({ reported: 1, delivered: 4 })).toBe(25);
		expect(clickRate({ clicked: 1, delivered: 5 })).toBe(20);
	});

	it('returns 0 when nothing was delivered', () => {
		expect(reportRate({ reported: 2, delivered: 0 })).toBe(0);
		expect(clickRate({ clicked: 3, delivered: 0 })).toBe(0);
	});

	it('handles a full report rate', () => {
		expect(reportRate({ reported: 5, delivered: 5 })).toBe(100);
	});
});

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
});

describe('rates', () => {
	it('uses delivered as the denominator', () => {
		expect(reportRate({ reported: 1, delivered: 4 })).toBe(25);
		expect(clickRate({ clicked: 1, delivered: 5 })).toBe(20);
	});
});

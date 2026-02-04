import { describe, expect, it } from 'vitest';
import { formatDateTime } from './format';

describe('formatDateTime', () => {
	it('returns an em dash for empty or invalid values', () => {
		expect(formatDateTime(null)).toBe('—');
		expect(formatDateTime(undefined)).toBe('—');
		expect(formatDateTime('')).toBe('—');
		expect(formatDateTime('not-a-date')).toBe('—');
	});

	it('formats valid timestamps in fixed UTC', () => {
		expect(formatDateTime('2026-08-01T10:05:00.000Z')).toBe('2026-08-01 10:05 UTC');
	});
});

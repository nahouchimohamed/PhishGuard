import { describe, expect, it } from 'vitest';
import { clickTone, reportTone, toneClass } from './rateTone';

describe('rate tones', () => {
	it('treats high report rates as good', () => {
		expect(reportTone(80)).toBe('good');
		expect(reportTone(40)).toBe('warn');
		expect(reportTone(10)).toBe('bad');
	});

	it('treats low click rates as good', () => {
		expect(clickTone(5)).toBe('good');
		expect(clickTone(20)).toBe('warn');
		expect(clickTone(40)).toBe('bad');
	});

	it('maps tones to tailwind text classes', () => {
		expect(toneClass('good')).toContain('emerald');
		expect(toneClass('bad')).toContain('rose');
	});
});

import { describe, expect, it } from 'vitest';
import { eventKindLabel, eventToneClass } from './events';

describe('event presentation', () => {
	it('labels and tints known kinds', () => {
		expect(eventKindLabel('reported')).toBe('Reported');
		expect(eventToneClass('clicked')).toContain('rose');
		expect(eventToneClass('reported')).toContain('emerald');
	});
});

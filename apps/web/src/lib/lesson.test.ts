import { describe, expect, it } from 'vitest';
import { parseLesson } from './lesson';

describe('parseLesson', () => {
	it('returns a plain-text fallback when markers are missing', () => {
		expect(parseLesson('Just a short tip.')).toEqual({
			intro: 'Just a short tip.',
			redFlags: [],
			goodHabit: ''
		});
	});

	it('splits intro, red flags, and good habit', () => {
		const parsed = parseLesson(
			'Watch for urgency.\nRed flags:\n- odd sender\n- short link\nGood habit: Report first.'
		);
		expect(parsed.intro).toBe('Watch for urgency.');
		expect(parsed.redFlags).toEqual(['odd sender', 'short link']);
		expect(parsed.goodHabit).toBe('Report first.');
	});

	it('handles empty input', () => {
		expect(parseLesson('')).toEqual({ intro: '', redFlags: [], goodHabit: '' });
	});

	it('parses a good habit without red flags', () => {
		const parsed = parseLesson('Stay calm.\nGood habit: Pause before clicking.');
		expect(parsed.intro).toBe('Stay calm.');
		expect(parsed.redFlags).toEqual([]);
		expect(parsed.goodHabit).toBe('Pause before clicking.');
	});

	it('accepts semicolon-separated red flags', () => {
		const parsed = parseLesson('Tip.\nRed flags:\nodd domain; urgent tone\nGood habit: Ask IT.');
		expect(parsed.redFlags).toEqual(['odd domain', 'urgent tone']);
		expect(parsed.goodHabit).toBe('Ask IT.');
	});
});

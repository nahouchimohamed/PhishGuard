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
});

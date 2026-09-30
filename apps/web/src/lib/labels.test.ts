import { describe, expect, it } from 'vitest';
import { recipientLabel } from './labels';

describe('recipientLabel', () => {
	it('prefers reported over clicked', () => {
		expect(
			recipientLabel({
				delivered_at: '1',
				opened_at: '1',
				clicked_at: '1',
				reported_at: '1'
			})
		).toBe('Reported');
	});

	it('returns clicked when not reported', () => {
		expect(recipientLabel({ delivered_at: '1', clicked_at: '1' })).toBe('Clicked');
	});

	it('returns opened', () => {
		expect(recipientLabel({ delivered_at: '1', opened_at: '1' })).toBe('Opened');
	});

	it('returns delivered', () => {
		expect(recipientLabel({ delivered_at: '1' })).toBe('Delivered');
	});

	it('returns pending when undelivered', () => {
		expect(recipientLabel({})).toBe('Pending');
		expect(recipientLabel({ delivered_at: null, opened_at: '' })).toBe('Pending');
	});

	it('treats empty strings as absent timestamps', () => {
		expect(recipientLabel({ delivered_at: '1', opened_at: '', clicked_at: '', reported_at: '' })).toBe(
			'Delivered'
		);
	});
});

import { describe, expect, it } from 'vitest';
import { friendlyError } from './errors';

describe('friendlyError', () => {
	it('maps known API strings and status codes', () => {
		expect(friendlyError(401, 'unauthorized')).toBe('Please sign in again.');
		expect(friendlyError(404, 'not found')).toBe('That item was not found.');
		expect(friendlyError(503, 'anything')).toBe('The API is temporarily unreachable.');
		expect(friendlyError(429)).toBe('Too many attempts. Wait a moment and try again.');
	});

	it('falls back safely', () => {
		expect(friendlyError(400, 'custom')).toBe('custom');
		expect(friendlyError(500)).toBe('Something went wrong. Try again shortly.');
	});

	it('maps forbidden and conflict API messages', () => {
		expect(friendlyError(403, 'forbidden')).toBe('You do not have access to that.');
		expect(friendlyError(409, 'conflict')).toBe('That action conflicts with the current state.');
	});
});

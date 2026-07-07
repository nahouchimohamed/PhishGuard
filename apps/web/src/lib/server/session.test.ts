import { describe, expect, it } from 'vitest';
import { COOKIE, SESSION_MAX_AGE_SECONDS } from './session';

describe('session constants', () => {
	it('uses a stable cookie name', () => {
		expect(COOKIE).toBe('phishguard_session');
	});

	it('keeps SESSION_MAX_AGE_SECONDS at one week', () => {
		expect(SESSION_MAX_AGE_SECONDS).toBe(60 * 60 * 24 * 7);
		expect(SESSION_MAX_AGE_SECONDS).toBe(604800);
	});
});

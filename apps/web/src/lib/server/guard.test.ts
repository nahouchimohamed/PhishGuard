import { describe, expect, it } from 'vitest';
import { isRedirect } from '@sveltejs/kit';
import { requireRole, requireUser } from './guard';
import type { User } from '$lib/types';

const admin: User = {
	id: '1',
	org_id: 'o',
	email: 'a@x',
	display_name: 'A',
	role: 'admin'
};

describe('requireUser', () => {
	it('redirects anonymous users to login', () => {
		try {
			requireUser(null, '/campaigns');
			throw new Error('expected redirect');
		} catch (e) {
			expect(isRedirect(e)).toBe(true);
			if (isRedirect(e)) {
				expect(e.location).toContain('/login');
				expect(e.location).toContain('redirectTo');
			}
		}
	});

	it('returns the user when present', () => {
		expect(requireUser(admin)).toBe(admin);
	});
});

describe('requireRole', () => {
	it('sends the wrong role home', () => {
		try {
			requireRole(admin, 'employee');
			throw new Error('expected redirect');
		} catch (e) {
			expect(isRedirect(e)).toBe(true);
			if (isRedirect(e)) expect(e.location).toBe('/campaigns');
		}
	});
});

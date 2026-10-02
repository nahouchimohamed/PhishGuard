import { describe, expect, it } from 'vitest';
import { homePathFor } from './roles';

describe('homePathFor', () => {
	it('sends admins to campaigns', () => {
		expect(homePathFor('admin')).toBe('/overview');
	});

	it('sends employees to inbox', () => {
		expect(homePathFor('employee')).toBe('/inbox');
	});

	it('keeps admin and employee homes distinct', () => {
		expect(homePathFor('admin')).not.toBe(homePathFor('employee'));
	});
});

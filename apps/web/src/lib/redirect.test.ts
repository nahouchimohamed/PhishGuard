import { describe, expect, it } from 'vitest';
import { safeRedirectPath } from './redirect';

describe('safeRedirectPath', () => {
	it('accepts same-site relative paths', () => {
		expect(safeRedirectPath('/campaigns')).toBe('/campaigns');
		expect(safeRedirectPath('/inbox/abc')).toBe('/inbox/abc');
	});

	it('rejects open redirects and schemes', () => {
		expect(safeRedirectPath('//evil.example')).toBe('/');
		expect(safeRedirectPath('/\\evil')).toBe('/');
		expect(safeRedirectPath('https://evil.example')).toBe('/');
		expect(safeRedirectPath('javascript:alert(1)')).toBe('/');
		expect(safeRedirectPath(null)).toBe('/');
	});

	it('rejects empty and whitespace-only values', () => {
		expect(safeRedirectPath('')).toBe('/');
		expect(safeRedirectPath('   ')).toBe('/');
	});

	it('honours a custom fallback path', () => {
		expect(safeRedirectPath('https://evil.example', '/inbox')).toBe('/inbox');
		expect(safeRedirectPath(undefined, '/campaigns')).toBe('/campaigns');
	});
});

import { describe, expect, it } from 'vitest';
import { truncate } from './text';

describe('truncate', () => {
	it('ellipsizes long names', () => {
		expect(truncate('abcdefghij', 8)).toBe('abcdefg…');
	});
});

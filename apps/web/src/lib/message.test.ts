import { describe, expect, it } from 'vitest';
import { splitMessageBody } from './message';

describe('splitMessageBody', () => {
	it('returns a single text segment when no link marker exists', () => {
		expect(splitMessageBody('hello', 'https://x')).toEqual([{ type: 'text', value: 'hello' }]);
	});

	it('inserts a link segment for each {{LINK}} marker', () => {
		expect(splitMessageBody('Click {{LINK}} now', 'https://sim/x')).toEqual([
			{ type: 'text', value: 'Click ' },
			{ type: 'link', href: 'https://sim/x', label: 'https://sim/x' },
			{ type: 'text', value: ' now' }
		]);
	});
});

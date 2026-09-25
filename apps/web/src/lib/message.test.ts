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

	it('supports multiple link markers', () => {
		expect(splitMessageBody('A {{LINK}} B {{LINK}} C', 'https://sim/y')).toEqual([
			{ type: 'text', value: 'A ' },
			{ type: 'link', href: 'https://sim/y', label: 'https://sim/y' },
			{ type: 'text', value: ' B ' },
			{ type: 'link', href: 'https://sim/y', label: 'https://sim/y' },
			{ type: 'text', value: ' C' }
		]);
	});

	it('treats a leading marker as a link-first body', () => {
		expect(splitMessageBody('{{LINK}} after', 'https://sim/z')).toEqual([
			{ type: 'link', href: 'https://sim/z', label: 'https://sim/z' },
			{ type: 'text', value: ' after' }
		]);
	});
});

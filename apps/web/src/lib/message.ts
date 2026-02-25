export type MessageSegment =
	| { type: 'text'; value: string }
	| { type: 'link'; href: string; label: string };

/** Split a template body on {{LINK}} without using {@html}. */
export function splitMessageBody(body: string, href: string): MessageSegment[] {
	const parts = (body ?? '').split('{{LINK}}');
	if (parts.length === 1) {
		return [{ type: 'text', value: body ?? '' }];
	}
	const out: MessageSegment[] = [];
	parts.forEach((part, i) => {
		if (part) out.push({ type: 'text', value: part });
		if (i < parts.length - 1) {
			out.push({ type: 'link', href, label: href });
		}
	});
	return out;
}

import type { RecipientRow } from '$lib/types';
import { recipientLabel } from '$lib/labels';

export function countByLabel(rows: RecipientRow[]): Record<string, number> {
	const out: Record<string, number> = {
		Reported: 0,
		Clicked: 0,
		Opened: 0,
		Delivered: 0,
		Pending: 0
	};
	for (const row of rows) out[recipientLabel(row)] += 1;
	return out;
}

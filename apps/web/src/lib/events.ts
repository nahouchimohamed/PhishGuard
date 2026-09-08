export function eventKindLabel(kind: string): string {
	const map: Record<string, string> = {
		delivered: 'Delivered',
		opened: 'Opened',
		clicked: 'Clicked',
		reported: 'Reported',
		closed: 'Closed',
		launched: 'Launched'
	};
	return map[kind] ?? kind;
}

export function eventToneClass(kind: string): string {
	switch (kind) {
		case 'reported':
			return 'border-emerald-200 bg-emerald-50';
		case 'clicked':
			return 'border-rose-200 bg-rose-50';
		case 'closed':
			return 'border-slate-300 bg-slate-100';
		default:
			return 'border-slate-200 bg-white';
	}
}

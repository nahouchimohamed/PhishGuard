/** Map a percentage into a traffic-light tone for tables and cards. */
export type RateTone = 'good' | 'warn' | 'bad' | 'neutral';

export function reportTone(rate: number): RateTone {
	if (!Number.isFinite(rate)) return 'neutral';
	if (rate >= 60) return 'good';
	if (rate >= 30) return 'warn';
	return 'bad';
}

export function clickTone(rate: number): RateTone {
	if (!Number.isFinite(rate)) return 'neutral';
	if (rate <= 10) return 'good';
	if (rate <= 25) return 'warn';
	return 'bad';
}

export function toneClass(tone: RateTone): string {
	switch (tone) {
		case 'good':
			return 'text-emerald-700';
		case 'warn':
			return 'text-amber-700';
		case 'bad':
			return 'text-rose-700';
		default:
			return 'text-slate-700';
	}
}

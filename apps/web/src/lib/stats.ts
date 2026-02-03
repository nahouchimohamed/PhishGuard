export function percent(part: number, whole: number): number {
	if (!Number.isFinite(part) || !Number.isFinite(whole) || whole <= 0) return 0;
	return Math.round((part / whole) * 1000) / 10;
}

export function reportRate(stats: { reported: number; delivered: number }): number {
	return percent(stats.reported, stats.delivered);
}

export function clickRate(stats: { clicked: number; delivered: number }): number {
	return percent(stats.clicked, stats.delivered);
}

/** Round a part/whole ratio to one decimal percent; invalid inputs yield 0. */
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

export function openRate(stats: { opened: number; delivered: number }): number {
	return percent(stats.opened, stats.delivered);
}

export type RateStats = {
	total: number;
	delivered: number;
	opened: number;
	clicked: number;
	reported: number;
};

export function aggregateStats(rows: RateStats[]): RateStats {
	const out: RateStats = { total: 0, delivered: 0, opened: 0, clicked: 0, reported: 0 };
	for (const row of rows) {
		out.total += row.total;
		out.delivered += row.delivered;
		out.opened += row.opened;
		out.clicked += row.clicked;
		out.reported += row.reported;
	}
	return out;
}

export function deliveryRate(stats: { delivered: number; total: number }): number {
	return percent(stats.delivered, stats.total);
}

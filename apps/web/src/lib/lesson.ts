export type ParsedLesson = {
	intro: string;
	redFlags: string[];
	goodHabit: string;
};

/** Parse seeded awareness copy stored on campaign templates. */
export function parseLesson(raw: string): ParsedLesson {
	const text = (raw ?? '').trim();
	if (!text) {
		return { intro: '', redFlags: [], goodHabit: '' };
	}

	const redIdx = text.search(/Red flags:/i);
	const habitIdx = text.search(/Good habit:/i);

	if (redIdx < 0 && habitIdx < 0) {
		return { intro: text, redFlags: [], goodHabit: '' };
	}

	const introEnd = redIdx >= 0 ? redIdx : habitIdx;
	const intro = text.slice(0, introEnd).trim();

	let redFlags: string[] = [];
	if (redIdx >= 0) {
		const start = redIdx + 'Red flags:'.length;
		const end = habitIdx >= 0 ? habitIdx : text.length;
		redFlags = text
			.slice(start, end)
			.split(/\n|;/)
			.map((s) => s.replace(/^[-*•\s]+/, '').trim())
			.filter(Boolean);
	}

	let goodHabit = '';
	if (habitIdx >= 0) {
		goodHabit = text.slice(habitIdx + 'Good habit:'.length).trim();
	}

	return { intro, redFlags, goodHabit };
}

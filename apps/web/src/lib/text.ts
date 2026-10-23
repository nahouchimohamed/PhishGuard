export function truncate(value: string, max = 80): string {
	const trimmed = value.trim();
	if (trimmed.length <= max) return trimmed;
	return `${trimmed.slice(0, Math.max(0, max - 1)).trimEnd()}…`;
}

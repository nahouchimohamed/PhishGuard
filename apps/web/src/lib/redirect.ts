/** Allow only same-site relative paths used after login. */
export function safeRedirectPath(raw: unknown, fallback = '/'): string {
	if (typeof raw !== 'string') return fallback;
	const value = raw.trim();
	if (!value.startsWith('/')) return fallback;
	if (value.startsWith('//')) return fallback;
	if (value.includes('\\')) return fallback;
	if (/^[a-zA-Z][a-zA-Z0-9+.-]*:/.test(value)) return fallback;
	return value;
}

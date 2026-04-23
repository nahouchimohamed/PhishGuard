const API_MESSAGES: Record<string, string> = {
	unauthorized: 'Please sign in again.',
	forbidden: 'You do not have access to that.',
	'not found': 'That item was not found.',
	'invalid input': 'Check the form and try again.',
	conflict: 'That action conflicts with the current state.',
	'internal error': 'Something went wrong. Try again shortly.'
};

export function friendlyError(status: number, apiMessage?: string | null): string {
	if (status === 503) return 'The API is temporarily unreachable.';
	if (status === 429) return 'Too many attempts. Wait a moment and try again.';
	if (status === 401) return 'Please sign in again.';
	if (status === 403) return 'You do not have access to that.';
	if (status === 404) return 'That item was not found.';
	const key = (apiMessage ?? '').trim().toLowerCase();
	if (key && API_MESSAGES[key]) return API_MESSAGES[key];
	if (status >= 500) return 'Something went wrong. Try again shortly.';
	return apiMessage?.trim() || 'Request failed.';
}

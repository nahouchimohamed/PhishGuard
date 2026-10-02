import type { Role } from '$lib/types';

export function homePathFor(role: Role): string {
	return role === 'admin' ? '/overview' : '/inbox';
}

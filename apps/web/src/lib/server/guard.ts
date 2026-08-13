import { redirect } from '@sveltejs/kit';
import type { User, Role } from '$lib/types';
import { homePathFor } from '$lib/roles';

export function requireUser(user: User | null | undefined, redirectTo?: string): User {
	if (!user) {
		const q = redirectTo ? `?redirectTo=${encodeURIComponent(redirectTo)}` : '';
		throw redirect(303, `/login${q}`);
	}
	return user;
}

export function requireRole(
	user: User | null | undefined,
	role: Role,
	redirectTo?: string
): User {
	const u = requireUser(user, redirectTo);
	if (u.role !== role) {
		throw redirect(303, homePathFor(u.role));
	}
	return u;
}

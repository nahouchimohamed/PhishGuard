import type { CampaignStatus } from '$lib/types';

export function isCloseable(status: CampaignStatus | string): boolean {
	return status === 'active' || status === 'launching';
}

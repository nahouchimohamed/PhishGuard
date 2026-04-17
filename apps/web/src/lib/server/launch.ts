import type { Campaign, LaunchResponse } from '$lib/types';

export type ParsedLaunch =
	| { ok: true; campaign: Campaign; warning?: string; incomplete: boolean }
	| { ok: false; error: string };

export function parseLaunchResponse(
	status: number,
	data: Campaign | LaunchResponse | null | undefined
): ParsedLaunch {
	if (!data) return { ok: false, error: 'Launch returned an empty body.' };

	if ('campaign' in data && data.campaign && typeof data.campaign === 'object' && 'id' in data.campaign) {
		return {
			ok: true,
			campaign: data.campaign,
			warning: data.warning,
			incomplete: status === 202 || Boolean(data.warning)
		};
	}

	if ('id' in data && typeof (data as Campaign).id === 'string') {
		return { ok: true, campaign: data as Campaign, incomplete: status === 202 };
	}

	return { ok: false, error: 'Launch response was not recognized.' };
}

/** Shared API shapes used by the SvelteKit BFF and UI. */
export type Role = 'admin' | 'employee';

export type User = {
	id: string;
	org_id: string;
	email: string;
	display_name: string;
	role: Role;
};

export type Template = {
	id: string;
	org_id: string;
	name: string;
	subject: string;
	body_text: string;
	lesson_text: string;
};

export type CampaignStats = {
	total: number;
	delivered: number;
	opened: number;
	clicked: number;
	reported: number;
};

export type CampaignStatus = 'draft' | 'launching' | 'active';

export type Campaign = {
	id: string;
	org_id: string;
	created_by: string;
	template_id: string;
	name: string;
	status: CampaignStatus | string;
	created_at: string;
	launched_at?: string | null;
	stats?: CampaignStats;
};

export type CampaignDetail = {
	campaign: Campaign;
	template: Template;
	stats: CampaignStats;
	recipients: RecipientRow[];
};

export type LoginResponse = {
	user: User;
	session_token: string;
};

export type SimLesson = {
	campaign_id: string;
	subject: string;
	lesson_text: string;
	name: string;
};

export type LaunchResponse = {
	campaign: Campaign;
	warning?: string;
};

export type InboxItem = {
	recipient_id: string;
	campaign_id: string;
	subject: string;
	body_text: string;
	lesson_text: string;
	from_label: string;
	delivered_at?: string | null;
	opened_at?: string | null;
	clicked_at?: string | null;
	reported_at?: string | null;
	click_token: string;
	sim_url?: string;
};

export type RecipientRow = {
	recipient_id: string;
	user_id: string;
	email: string;
	display_name: string;
	delivered_at?: string | null;
	opened_at?: string | null;
	clicked_at?: string | null;
	reported_at?: string | null;
};

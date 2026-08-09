const UUID_RE =
	/^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

export type CampaignFormInput = {
	name: string;
	template_id: string;
	recipient_user_ids: string[];
};

export type CampaignFormErrors = {
	name?: string;
	template_id?: string;
	recipient_user_ids?: string;
};

export type CampaignFormResult =
	| { ok: true; values: CampaignFormInput }
	| { ok: false; errors: CampaignFormErrors; values: CampaignFormInput };

export function validateCampaignForm(input: {
	name: unknown;
	template_id: unknown;
	recipient_user_ids: unknown[];
}): CampaignFormResult {
	const name = String(input.name ?? '').trim();
	const template_id = String(input.template_id ?? '').trim();
	const rawIds = input.recipient_user_ids.map((v) => String(v));
	const recipient_user_ids = [...new Set(rawIds.filter(Boolean))];
	const values = { name, template_id, recipient_user_ids };
	const errors: CampaignFormErrors = {};

	if (!name) errors.name = 'Name is required.';
	else if (name.length > 120) errors.name = 'Name must be 120 characters or fewer.';

	if (!template_id) errors.template_id = 'Template is required.';
	else if (!UUID_RE.test(template_id)) errors.template_id = 'Template id looks invalid.';

	if (recipient_user_ids.length === 0) {
		errors.recipient_user_ids = 'Pick at least one employee.';
	} else if (recipient_user_ids.some((id) => !UUID_RE.test(id))) {
		errors.recipient_user_ids = 'One or more employee ids look invalid.';
	}

	if (Object.keys(errors).length > 0) {
		return { ok: false, errors, values };
	}
	return { ok: true, values };
}

import { describe, expect, it } from 'vitest';
import { validateCampaignForm } from './validate';

const tid = '11111111-1111-4111-8111-111111111111';
const uid = '22222222-2222-4222-8222-222222222222';

describe('validateCampaignForm', () => {
	it('accepts a trimmed valid payload and dedupes recipients', () => {
		const res = validateCampaignForm({
			name: '  Drill  ',
			template_id: tid,
			recipient_user_ids: [uid, uid]
		});
		expect(res.ok).toBe(true);
		if (res.ok) {
			expect(res.values.name).toBe('Drill');
			expect(res.values.recipient_user_ids).toEqual([uid]);
		}
	});

	it('returns field errors for missing and oversized input', () => {
		const res = validateCampaignForm({
			name: 'x'.repeat(121),
			template_id: '',
			recipient_user_ids: []
		});
		expect(res.ok).toBe(false);
		if (!res.ok) {
			expect(res.errors.name).toMatch(/120/);
			expect(res.errors.template_id).toBeTruthy();
			expect(res.errors.recipient_user_ids).toBeTruthy();
		}
	});

	it('rejects a non-uuid template id', () => {
		const res = validateCampaignForm({
			name: 'Drill',
			template_id: 'not-a-uuid',
			recipient_user_ids: [uid]
		});
		expect(res.ok).toBe(false);
		if (!res.ok) {
			expect(res.errors.template_id).toMatch(/invalid/i);
		}
	});

	it('rejects invalid recipient ids', () => {
		const res = validateCampaignForm({
			name: 'Drill',
			template_id: tid,
			recipient_user_ids: ['bad-id']
		});
		expect(res.ok).toBe(false);
		if (!res.ok) {
			expect(res.errors.recipient_user_ids).toMatch(/invalid/i);
		}
	});

	it('rejects a blank name after trim', () => {
		const res = validateCampaignForm({
			name: '   ',
			template_id: tid,
			recipient_user_ids: [uid]
		});
		expect(res.ok).toBe(false);
		if (!res.ok) {
			expect(res.errors.name).toMatch(/required/i);
		}
	});
});

<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';

	let { data, form } = $props();
	let pending = $state(false);
</script>

<PageShell title="New campaign" user={data.user} width="max-w-2xl">
	{#snippet banner()}Training mode — templates are seeded and stay read-only.{/snippet}

	<p class="mb-2 text-sm">
		<a class="text-accent hover:underline" href="/templates">Browse templates</a>
		· <a class="text-accent hover:underline" href="/employees">Browse employees</a>
	</p>
	<h1 class="page-title">New campaign</h1>
	<p class="page-subtitle mt-1">
		Saves recipients, then enqueues one inbox delivery per person.
	</p>

	{#if form?.error || data.error}
		<div class="mt-4">
			<Alert>{form?.error ?? data.error}</Alert>
		</div>
	{/if}

	{#if data.templates.length === 0 || data.employees.length === 0}
		<div class="card mt-6 text-slate-600">
			{#if data.templates.length === 0}
				No templates are available yet.
			{:else}
				No employees are available to target.
			{/if}
		</div>
	{:else}
		<form
			method="POST"
			class="card mt-6 space-y-5 shadow-panel"
			use:enhance={() => {
				pending = true;
				return async ({ update }) => {
					await update();
					pending = false;
				};
			}}
		>
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Campaign name</span>
				<input
					class="input"
					name="name"
					required
					maxlength="120"
					placeholder="Q1 password-expiry drill"
					value={form?.values?.name ?? ''}
					aria-invalid={form?.errors?.name ? 'true' : undefined}
					aria-describedby={form?.errors?.name ? 'name-error' : undefined}
				/>
				{#if form?.errors?.name}
					<span id="name-error" class="mt-1 block text-xs text-red-700">{form.errors.name}</span>
				{/if}
			</label>

			<label class="block text-sm">
				<span class="mb-1 block font-medium">Template</span>
				<select
					class="input"
					name="template_id"
					required
					aria-invalid={form?.errors?.template_id ? 'true' : undefined}
					aria-describedby={form?.errors?.template_id ? 'template-error' : undefined}
				>
					<option value="">Select a template</option>
					{#each data.templates as t}
						<option value={t.id} selected={t.id === data.templatePrefill} selected={form?.values?.template_id === t.id}
							>{t.name} — {t.subject}</option
						>
					{/each}
				</select>
				{#if form?.errors?.template_id}
					<span id="template-error" class="mt-1 block text-xs text-red-700"
						>{form.errors.template_id}</span
					>
				{/if}
			</label>

			<fieldset>
				<legend class="mb-2 text-sm font-medium">Employees</legend>
				<div class="space-y-2">
					{#each data.employees as e}
						<label class="flex items-center gap-2 text-sm">
							<input
								type="checkbox"
								name="recipient_user_ids"
								value={e.id}
								checked={form?.values?.recipient_user_ids
									? form.values.recipient_user_ids.includes(e.id)
									: true}
							/>
							<span>{e.display_name} ({e.email})</span>
						</label>
					{/each}
				</div>
				{#if form?.errors?.recipient_user_ids}
					<span class="mt-1 block text-xs text-red-700">{form.errors.recipient_user_ids}</span>
				{/if}
			</fieldset>

			<div class="flex gap-3">
				<SubmitButton label="Launch" pendingLabel="Launching…" {pending} />
				<a class="btn-secondary" href="/campaigns">Cancel</a>
			</div>
		</form>
	{/if}
</PageShell>

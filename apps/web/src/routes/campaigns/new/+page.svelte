<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	let { data, form } = $props();
</script>

<div class="banner">Training mode — templates are seeded and read-only.</div>
<Nav user={data.user!} />

<main class="mx-auto max-w-2xl px-4 py-8">
	<h1 class="text-2xl font-semibold">New campaign</h1>
	<p class="mt-1 text-sm text-slate-600">
		Commit-then-publish: the API saves recipients, then enqueues one SQS message each.
	</p>

	{#if form?.error || data.error}
		<p class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">
			{form?.error ?? data.error}
		</p>
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
		<form method="POST" class="card mt-6 space-y-5">
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Campaign name</span>
				<input
					class="input"
					name="name"
					required
					placeholder="Q1 password-expiry drill"
					value={form?.values?.name ?? ''}
				/>
			</label>

			<label class="block text-sm">
				<span class="mb-1 block font-medium">Template</span>
				<select class="input" name="template_id" required>
					<option value="">Select a template</option>
					{#each data.templates as t}
						<option value={t.id} selected={form?.values?.template_id === t.id}
							>{t.name} — {t.subject}</option
						>
					{/each}
				</select>
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
			</fieldset>

			<div class="flex gap-3">
				<button class="btn" type="submit">Launch</button>
				<a class="btn-secondary" href="/campaigns">Cancel</a>
			</div>
		</form>
	{/if}
</main>

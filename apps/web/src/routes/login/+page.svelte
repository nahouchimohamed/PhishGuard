<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';

	let { form } = $props();
	let pending = $state(false);
</script>

<PageShell title="Sign in" width="max-w-md">
	{#snippet banner()}PhishGuard is a local training simulator — no real email is ever sent.{/snippet}

	<div class="card">
		<h1 class="text-2xl font-semibold">Sign in</h1>
		<p class="mt-1 text-sm text-slate-600">
			Seeded users share password <code class="rounded bg-slate-100 px-1">phishguard-dev</code>
		</p>

		{#if form?.error}
			<div class="mt-4"><Alert>{form.error}</Alert></div>
		{/if}

		<form
			method="POST"
			class="mt-6 space-y-4"
			use:enhance={() => {
				pending = true;
				return async ({ update }) => {
					await update();
					pending = false;
				};
			}}
		>
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Email</span>
				<input
					class="input"
					name="email"
					type="email"
					required
					autocomplete="username"
					value={form?.email ?? 'admin@phishguard.local'}
				/>
			</label>
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Password</span>
				<input
					class="input"
					name="password"
					type="password"
					required
					autocomplete="current-password"
					value="phishguard-dev"
				/>
			</label>
			<SubmitButton class="btn w-full" label="Continue" pendingLabel="Signing in…" {pending} />
		</form>

		<ul class="mt-6 space-y-1 text-xs text-slate-500">
			<li>Admin: admin@phishguard.local</li>
			<li>Employees: alice@ / bob@ / cara@phishguard.local</li>
		</ul>
	</div>
</PageShell>

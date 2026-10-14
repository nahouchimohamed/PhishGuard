<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';

	let { form } = $props();
	let pending = $state(false);
</script>

<PageShell title="Sign in" width="max-w-md">
	{#snippet banner()}Training only — credentials never leave this workstation.{/snippet}

	<div class="card shadow-md ring-1 ring-slate-200/60">
		<h1 class="page-title">Sign in</h1>
		<p class="page-subtitle mt-2">
			Local seed password:
			<code class="rounded bg-slate-100 px-1.5 py-0.5 font-mono text-xs">phishguard-dev</code>
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
			<SubmitButton class="btn w-full" label="Sign in" pendingLabel="Signing in…" {pending} />
		</form>

		<ul class="mt-6 space-y-1 text-xs text-slate-500">
			<li>Admin: admin@phishguard.local</li>
			<li>Employees: alice@ / bob@ / cara@phishguard.local</li>
		</ul>
	</div>
</PageShell>

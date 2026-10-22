<script lang="ts">
	import PageShell from '$lib/components/PageShell.svelte';

	let { data } = $props();
</script>

<PageShell title="Account" user={data.user} width="max-w-2xl">
	{#snippet banner()}Local simulation account — cookie session only.{/snippet}

	<h1 class="page-title mb-6">Account</h1>
	<div class="card grid gap-4 sm:grid-cols-2">
		<div>
			<div class="text-xs uppercase tracking-wide text-slate-500">Name</div>
			<div class="mt-1 font-medium">{data.account.display_name}</div>
		</div>
		<div>
			<div class="text-xs uppercase tracking-wide text-slate-500">Email</div>
			<div class="mt-1 font-medium">{data.account.email}</div>
		</div>
		<div>
			<div class="text-xs uppercase tracking-wide text-slate-500">Role</div>
			<div class="mt-1">
				<span
					class={`inline-flex rounded-md px-2 py-0.5 text-xs font-medium capitalize ${data.account.role === 'admin' ? 'bg-ink text-white' : 'bg-teal-50 text-teal-900'}`}
					>{data.account.role}</span
				>
			</div>
		</div>
		<div>
			<div class="text-xs uppercase tracking-wide text-slate-500">Organization</div>
			<div class="mt-1 font-mono text-xs">{data.account.org_id}</div>
		</div>
	</div>
</PageShell>

<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	import { recipientLabel } from '$lib/labels';
	let { data, form } = $props();
	const d = $derived(data.detail);
</script>

<div class="banner">Training mode — watch Grafana at localhost:3001 for the same counters.</div>
<Nav user={data.user!} />

<main class="mx-auto max-w-5xl px-4 py-8">
	{#if !d}
		<p class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{data.error}</p>
	{:else}
		<div class="mb-6 flex flex-wrap items-start justify-between gap-4">
			<div>
				<p class="text-sm text-slate-500">
					<a class="hover:underline" href="/campaigns">Campaigns</a> /
				</p>
				<h1 class="text-2xl font-semibold">{d.campaign.name}</h1>
				<p class="text-sm text-slate-600">
					Status <span class="capitalize">{d.campaign.status}</span> · Template {d.template.name}
				</p>
			</div>
			{#if d.campaign.status === 'launching'}
				<form method="POST" action="?/retry">
					<button class="btn" type="submit">Retry queue publish</button>
				</form>
			{/if}
		</div>

		{#if form?.error}
			<p class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{form.error}</p>
		{/if}

		<div class="mb-6 grid grid-cols-2 gap-3 sm:grid-cols-5">
			{#each [
				['Total', d.stats.total],
				['Delivered', d.stats.delivered],
				['Opened', d.stats.opened],
				['Reported', d.stats.reported],
				['Clicked', d.stats.clicked]
			] as [label, value]}
				<div class="card text-center">
					<div class="text-2xl font-semibold">{value}</div>
					<div class="text-xs uppercase tracking-wide text-slate-500">{label}</div>
				</div>
			{/each}
		</div>

		<div class="card">
			<h2 class="mb-3 font-medium">Recipients</h2>
			<table class="min-w-full text-left text-sm">
				<thead class="text-slate-500">
					<tr>
						<th class="pb-2 pr-4">Name</th>
						<th class="pb-2 pr-4">Email</th>
						<th class="pb-2">Status</th>
					</tr>
				</thead>
				<tbody>
					{#each d.recipients as r}
						<tr class="border-t border-slate-100">
							<td class="py-2 pr-4">{r.display_name}</td>
							<td class="py-2 pr-4">{r.email}</td>
							<td class="py-2">{recipientLabel(r)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</main>

<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';
	import { formatDateTime } from '$lib/format';
	import { recipientLabel } from '$lib/labels';
	import { clickTone, reportTone } from '$lib/rateTone';
	import { clickRate, openRate, reportRate } from '$lib/stats';

	let { data, form } = $props();
	const d = $derived(data.detail);
	let pending = $state(false);
</script>

<PageShell title={d?.campaign.name ?? 'Campaign'} user={data.user} width="max-w-5xl">
	{#snippet banner()}
		Training mode — watch Grafana at
		<a class="underline" href="http://localhost:3002">localhost:3002</a> for the same counters.
	{/snippet}

	{#if !d}
		<Alert>{data.error}</Alert>
	{:else}
		<div class="mb-6 flex flex-wrap items-start justify-between gap-4">
			<div>
				<nav class="text-sm text-slate-500" aria-label="Breadcrumb">
					<a class="hover:underline" href="/campaigns">Campaigns</a>
					<span aria-hidden="true"> / </span>
					<span class="text-slate-700">{d.campaign.name}</span>
				</nav>
				<h1 class="text-2xl font-semibold">{d.campaign.name}</h1>
				<p class="text-sm text-slate-600">
					Status <span class="capitalize">{d.campaign.status}</span> · Template {d.template.name}
				</p>
				<p class="mt-1 text-xs text-slate-500">
					Created <time datetime={d.campaign.created_at}>{formatDateTime(d.campaign.created_at)}</time>
					{#if d.campaign.launched_at}
						· Launched
						<time datetime={d.campaign.launched_at}>{formatDateTime(d.campaign.launched_at)}</time>
					{/if}
				</p>
			</div>
			{#if d.campaign.status === 'launching'}
				<form
					method="POST"
					action="?/retry"
					use:enhance={() => {
						pending = true;
						return async ({ update }) => {
							await update();
							pending = false;
						};
					}}
				>
					<SubmitButton label="Retry queue publish" pendingLabel="Retrying…" {pending} />
				</form>
			{/if}
		</div>

		{#if form?.error}
			<div class="mb-4"><Alert>{form.error}</Alert></div>
		{/if}
		{#if data.notice}
			<div class="mb-4"><Alert tone="warning">{data.notice}</Alert></div>
		{/if}

		<div class="mb-6 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
			<StatCard label="Total" value={d.stats.total} />
			<StatCard label="Delivered" value={d.stats.delivered} />
			<StatCard label="Opened" value={d.stats.opened} />
			<StatCard label="Open %" value={`${openRate(d.stats)}%`} />
			<StatCard label="Report %" value={`${reportRate(d.stats)}%`} tone={reportTone(reportRate(d.stats))} />
			<StatCard label="Click %" value={`${clickRate(d.stats)}%`} tone={clickTone(clickRate(d.stats))} />
		</div>

		<div class="card overflow-x-auto">
			<h2 class="mb-3 font-medium">Recipients</h2>
			{#if d.recipients.length === 0}
				<p class="text-sm text-slate-600">No recipients on this campaign.</p>
			{:else}
				<table class="min-w-full text-left text-sm">
					<caption class="sr-only">Campaign recipients</caption>
					<thead class="text-slate-500">
						<tr>
							<th class="pb-2 pr-4" scope="col">Name</th>
							<th class="pb-2 pr-4" scope="col">Email</th>
							<th class="pb-2 pr-4" scope="col">Status</th>
							<th class="pb-2" scope="col">Delivered</th>
						</tr>
					</thead>
					<tbody>
						{#each d.recipients as r}
							<tr class="border-t border-slate-100">
								<td class="py-2 pr-4">{r.display_name}</td>
								<td class="py-2 pr-4">{r.email}</td>
								<td class="py-2 pr-4"><StatusBadge label={recipientLabel(r)} /></td>
								<td class="py-2">{formatDateTime(r.delivered_at)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
	{/if}
</PageShell>

<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import CampaignStatusPill from '$lib/components/CampaignStatusPill.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';
	import { isCloseable } from '$lib/campaignStatus';
	import { eventKindLabel, eventToneClass } from '$lib/events';
	import { formatDateTime } from '$lib/format';
	import { recipientLabel } from '$lib/labels';
	import { clickTone, reportTone } from '$lib/rateTone';
	import { countByLabel } from '$lib/recipients';
	import { clickRate, openRate, reportRate } from '$lib/stats';

	let { data, form } = $props();
	const d = $derived(data.detail);
	const tallies = $derived(d ? countByLabel(d.recipients) : null);
	let outcomeFilter = $state<'All' | 'Reported' | 'Clicked' | 'Opened' | 'Delivered' | 'Pending'>('All');
	const visibleRecipients = $derived(
		!d
			? []
			: outcomeFilter === 'All'
				? d.recipients
				: d.recipients.filter((r) => recipientLabel(r) === outcomeFilter)
	);
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
					<a class="hover:underline" href="/overview">Overview</a>
					<span aria-hidden="true"> / </span>
					<a class="hover:underline" href="/campaigns">Campaigns</a>
					<span aria-hidden="true"> / </span>
					<span class="text-slate-700">{d.campaign.name}</span>
				</nav>
				<div class="mt-1 flex flex-wrap items-center gap-3">
					<h1 class="page-title">{d.campaign.name}</h1>
					<CampaignStatusPill status={d.campaign.status} />
				</div>
				<p class="page-subtitle mt-1">
					Template {d.template.name}
					{#if d.campaign.status === 'launching'}
						· Queue publish in progress
					{/if}
				</p>
				<p class="mt-1 text-xs text-slate-500">
					Created <time datetime={d.campaign.created_at}>{formatDateTime(d.campaign.created_at)}</time>
					{#if d.campaign.launched_at}
						· Launched
						<time datetime={d.campaign.launched_at}>{formatDateTime(d.campaign.launched_at)}</time>
					{/if}
				</p>
			</div>
			<div class="flex flex-wrap gap-2">
				<a
					class="btn-secondary"
					class:pointer-events-none={d.recipients.length === 0}
					class:opacity-50={d.recipients.length === 0}
					href={`/campaigns/${d.campaign.id}/export`}>Download CSV</a
				>
				{#if isCloseable(d.campaign.status)}
					<form
						method="POST"
						action="?/close"
						use:enhance={() => {
							if (!confirm('Close this campaign? Recipients keep their history.')) {
								return async () => {};
							}
							pending = true;
							return async ({ update }) => {
								await update();
								pending = false;
							};
						}}
					>
						<button class="btn-danger" type="submit" disabled={pending}
							>{pending ? 'Closing…' : 'Close campaign'}</button
						>
					</form>
				{/if}
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
			<StatCard
				label="Report %"
				value={`${reportRate(d.stats)}%`}
				tone={reportTone(reportRate(d.stats))}
			/>
			<StatCard
				label="Click %"
				value={`${clickRate(d.stats)}%`}
				tone={clickTone(clickRate(d.stats))}
			/>
		</div>

		<div class="card overflow-x-auto">
			<h2 class="mb-2 font-medium">Recipients</h2>
			<div class="mb-3 flex flex-wrap items-center justify-between gap-3">
				{#if tallies}
					<p class="text-xs text-slate-500">
						Reported {tallies.Reported} · Clicked {tallies.Clicked} · Opened {tallies.Opened} · Delivered
						{tallies.Delivered} · Pending {tallies.Pending}
					</p>
				{/if}
				<label class="text-xs text-slate-600">
					Outcome
					<select class="input ml-2 w-36" bind:value={outcomeFilter}>
						{#each ['All', 'Reported', 'Clicked', 'Opened', 'Delivered', 'Pending'] as opt}
							<option value={opt}>{opt}</option>
						{/each}
					</select>
				</label>
			</div>
			{#if visibleRecipients.length === 0}
				<p class="text-sm text-slate-600">
					{#if outcomeFilter === 'All'}
						No recipients on this campaign.
					{:else}
						No recipients with status {outcomeFilter}.
					{/if}
				</p>
			{:else}
				<table class="min-w-full text-left text-sm">
					<thead class="text-slate-500">
						<tr>
							<th class="pb-2 pr-4" scope="col">Name</th>
							<th class="pb-2 pr-4" scope="col">Email</th>
							<th class="pb-2 pr-4" scope="col">Status</th>
							<th class="pb-2" scope="col">Delivered</th>
						</tr>
					</thead>
					<tbody>
						{#each visibleRecipients as r}
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

		<div class="card mt-6">
			<h2 class="mb-3 font-medium">Activity</h2>
			{#if !d.events || d.events.length === 0}
				<p class="text-sm text-slate-600">No events recorded yet.</p>
			{:else}
				<ol class="space-y-2">
					{#each d.events as ev}
						<li class={`rounded-lg border px-3 py-2 text-sm ${eventToneClass(ev.kind)}`}>
							<div class="flex flex-wrap gap-2">
								<span class="font-medium">{eventKindLabel(ev.kind)}</span>
								<time class="text-slate-500" datetime={ev.created_at}>{formatDateTime(ev.created_at)}</time>
							</div>
							{#if ev.detail}
								<p class="mt-1 text-slate-700">{ev.detail}</p>
							{/if}
						</li>
					{/each}
				</ol>
			{/if}
		</div>
	{/if}
</PageShell>

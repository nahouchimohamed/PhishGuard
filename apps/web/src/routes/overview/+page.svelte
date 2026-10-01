<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import CampaignStatusPill from '$lib/components/CampaignStatusPill.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import { recentCampaigns } from '$lib/overview';
	import { clickTone, reportTone } from '$lib/rateTone';
	import { clickRate, deliveryRate, openRate, reportRate } from '$lib/stats';

	let { data } = $props();
	const s = $derived(data.summary);
	const activeCount = $derived(data.campaigns.filter((c) => c.status === 'active').length);
	const recent = $derived(recentCampaigns(data.campaigns, 5));
</script>

<PageShell title="Overview" user={data.user} width="max-w-5xl">
	{#snippet banner()}Org-wide drill outcomes across every campaign.{/snippet}

	<div class="mb-6 flex flex-wrap items-center justify-between gap-4">
		<div>
			<h1 class="page-title">Overview</h1>
			<p class="page-subtitle">{data.campaigns.length} campaigns · {activeCount} active</p>
		</div>
		<div class="flex gap-2">
			<a class="btn-secondary" href="/campaigns">All campaigns</a>
			<a class="btn" href="/campaigns/new">New campaign</a>
		</div>
	</div>

	{#if data.error}
		<Alert>{data.error}</Alert>
	{:else}
		<div class="mb-8 grid grid-cols-2 gap-3 lg:grid-cols-5">
			<StatCard label="Recipients" value={s.total} />
			<StatCard label="Delivery %" value={`${deliveryRate(s)}%`} />
			<StatCard label="Open %" value={`${openRate(s)}%`} />
			<StatCard label="Report %" value={`${reportRate(s)}%`} tone={reportTone(reportRate(s))} />
			<StatCard label="Click %" value={`${clickRate(s)}%`} tone={clickTone(clickRate(s))} />
		</div>

		<div class="card">
			<h2 class="mb-3 font-medium">Recent campaigns</h2>
			<ul class="divide-y divide-slate-100">
				{#each recent as c}
					<li class="flex items-center justify-between gap-3 py-2 text-sm">
						<div class="flex min-w-0 items-center gap-2">
							<a class="truncate font-medium text-accent hover:underline" href={`/campaigns/${c.id}`}
								>{c.name}</a
							>
							<CampaignStatusPill status={c.status} />
						</div>
						<span class="shrink-0 text-slate-500"
							>{c.stats ? `${reportRate(c.stats)}% reported` : '—'}</span
						>
					</li>
				{:else}
					<li class="py-4 text-slate-500">
						No campaigns yet.
						<a class="ml-2 text-accent hover:underline" href="/campaigns/new">Launch one</a>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</PageShell>

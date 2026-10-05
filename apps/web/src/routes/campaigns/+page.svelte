<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import CampaignStatusPill from '$lib/components/CampaignStatusPill.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import ProgressBar from '$lib/components/ProgressBar.svelte';
	import { formatDateTime } from '$lib/format';
	import { clickTone, reportTone, toneClass } from '$lib/rateTone';
	import { clickRate, openRate, reportRate } from '$lib/stats';

	let { data } = $props();
</script>

<PageShell title="Campaigns" user={data.user} width="max-w-5xl">
	{#snippet banner()}Training mode — campaigns deliver into the in-app inbox only (no SMTP).{/snippet}

	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="page-title">Campaigns</h1>
			<p class="page-subtitle">Launch a drill, then watch report vs click rates.</p>
		</div>
		<a class="btn" href="/campaigns/new">New campaign</a>
	</div>

	<div class="mb-4 flex flex-wrap items-center gap-3">
		<div class="flex flex-wrap gap-2 text-sm" role="navigation" aria-label="Campaign status filter">
			{#each [['all', 'All'], ['active', 'Active'], ['launching', 'Launching'], ['closed', 'Closed'], ['draft', 'Draft']] as [value, label]}
				<a
					class={`rounded-full px-3 py-1 ${data.statusFilter === value ? 'bg-ink text-white' : 'bg-white text-slate-700 ring-1 ring-slate-200 hover:bg-slate-50'}`}
					href={value === 'all' ? `/campaigns${data.q ? `?q=${encodeURIComponent(data.q)}` : ''}` : `/campaigns?status=${value}${data.q ? `&q=${encodeURIComponent(data.q)}` : ''}`}
					aria-current={data.statusFilter === value ? 'page' : undefined}>{label}</a
				>
			{/each}
		</div>
		<form class="ml-auto flex flex-wrap items-center gap-2" method="GET" action="/campaigns">
			{#if data.statusFilter !== 'all'}
				<input type="hidden" name="status" value={data.statusFilter} />
			{/if}
			<select class="input w-40" name="sort" value={data.sort}>
				<option value="newest">Newest</option>
				<option value="report">Best report %</option>
				<option value="name">Name</option>
			</select>
			<input class="input w-56" name="q" value={data.q} placeholder="Search campaigns" />
			<button class="btn-secondary px-3" type="submit">Apply</button>
		</form>
	</div>

	{#if data.error}
		<Alert>{data.error}</Alert>
	{:else if data.campaigns.length === 0}
		<div class="card-muted">
			{#if data.statusFilter !== 'all' || data.q}
				No campaigns match these filters.
			{:else}
				No campaigns yet. Create one to start a drill.
			{/if}
		</div>
	{:else}
		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white">
			<table class="min-w-full text-left text-sm">
				<caption class="sr-only">Campaign list</caption>
				<thead class="bg-slate-50 text-slate-600">
					<tr>
						<th class="px-4 py-3 font-medium" scope="col">Name</th>
						<th class="px-4 py-3 font-medium" scope="col">Status</th>
						<th class="px-4 py-3 font-medium" scope="col">Created</th>
						<th class="px-4 py-3 font-medium" scope="col">Delivered</th>
						<th class="px-4 py-3 font-medium" scope="col">Open %</th>
						<th class="px-4 py-3 font-medium" scope="col">Report %</th>
						<th class="px-4 py-3 font-medium" scope="col">Click %</th>
					</tr>
				</thead>
				<tbody>
					{#each data.campaigns as c}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3">
								<a class="font-medium text-accent hover:underline" href={`/campaigns/${c.id}`}
									>{c.name}</a
								>
							</td>
							<td class="px-4 py-3"><CampaignStatusPill status={c.status} /></td>
							<td class="px-4 py-3">{formatDateTime(c.created_at)}</td>
							<td class="px-4 py-3">
								<div class="mb-1 text-xs text-slate-500">{c.stats?.delivered ?? 0}/{c.stats?.total ?? 0}</div>
								<ProgressBar value={c.stats?.delivered ?? 0} max={c.stats?.total ?? 0} />
							</td>
							<td class="px-4 py-3">{c.stats ? openRate(c.stats) : 0}%</td>
							<td class={`px-4 py-3 font-medium ${c.stats ? toneClass(reportTone(reportRate(c.stats))) : ''}`}
								>{c.stats ? reportRate(c.stats) : 0}%</td
							>
							<td class={`px-4 py-3 font-medium ${c.stats ? toneClass(clickTone(clickRate(c.stats))) : ''}`}
								>{c.stats ? clickRate(c.stats) : 0}%</td
							>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</PageShell>

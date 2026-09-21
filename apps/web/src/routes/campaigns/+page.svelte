<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import { formatDateTime } from '$lib/format';
	import { clickRate, reportRate } from '$lib/stats';

	let { data } = $props();
</script>

<PageShell title="Campaigns" user={data.user} width="max-w-5xl">
	{#snippet banner()}Training mode — campaigns deliver into the in-app inbox only (no SMTP).{/snippet}

	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-semibold">Campaigns</h1>
			<p class="text-sm text-slate-600">Launch a drill, then watch report vs click rates.</p>
		</div>
		<a class="btn" href="/campaigns/new">New campaign</a>
	</div>

	{#if data.error}
		<Alert>{data.error}</Alert>
	{:else if data.campaigns.length === 0}
		<div class="card text-slate-600">No campaigns yet. Create one to start a drill.</div>
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
							<td class="px-4 py-3 capitalize">{c.status}</td>
							<td class="px-4 py-3">{formatDateTime(c.created_at)}</td>
							<td class="px-4 py-3">{c.stats?.delivered ?? 0}/{c.stats?.total ?? 0}</td>
							<td class="px-4 py-3">{c.stats ? reportRate(c.stats) : 0}%</td>
							<td class="px-4 py-3">{c.stats ? clickRate(c.stats) : 0}%</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</PageShell>

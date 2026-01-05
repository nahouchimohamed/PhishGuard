<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	let { data } = $props();
</script>

<div class="banner">Training mode — campaigns deliver into the in-app inbox only.</div>
<Nav user={data.user!} />

<main class="mx-auto max-w-5xl px-4 py-8">
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-semibold">Campaigns</h1>
			<p class="text-sm text-slate-600">Launch a drill, then watch report vs click rates.</p>
		</div>
		<a class="btn" href="/campaigns/new">New campaign</a>
	</div>

	{#if data.error}
		<p class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{data.error}</p>
	{/if}

	{#if data.campaigns.length === 0}
		<div class="card text-slate-600">No campaigns yet. Create one to start a drill.</div>
	{:else}
		<div class="overflow-hidden rounded-xl border border-slate-200 bg-white">
			<table class="min-w-full text-left text-sm">
				<thead class="bg-slate-50 text-slate-600">
					<tr>
						<th class="px-4 py-3 font-medium">Name</th>
						<th class="px-4 py-3 font-medium">Status</th>
						<th class="px-4 py-3 font-medium">Delivered</th>
						<th class="px-4 py-3 font-medium">Reported</th>
						<th class="px-4 py-3 font-medium">Clicked</th>
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
							<td class="px-4 py-3">{c.stats?.delivered ?? 0}/{c.stats?.total ?? 0}</td>
							<td class="px-4 py-3">{c.stats?.reported ?? 0}</td>
							<td class="px-4 py-3">{c.stats?.clicked ?? 0}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</main>

<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';

	let { data } = $props();
	let q = $state('');
	const filtered = $derived(
		data.employees.filter((e) => {
			const needle = q.trim().toLowerCase();
			if (!needle) return true;
			return (
				e.display_name.toLowerCase().includes(needle) || e.email.toLowerCase().includes(needle)
			);
		})
	);
</script>

<PageShell title="Employees" user={data.user} width="max-w-5xl">
	{#snippet banner()}Recipient directory for local drills.{/snippet}

	<div class="mb-6 flex flex-wrap items-end justify-between gap-4">
		<div>
			<h1 class="page-title">Employees</h1>
			<p class="page-subtitle">{data.employees.length} people in this org</p>
		</div>
		<input class="input w-64" bind:value={q} placeholder="Filter by name or email" />
	</div>

	{#if data.error}
		<Alert>{data.error}</Alert>
	{:else}
		<div class="overflow-hidden rounded-xl border border-slate-200 bg-white">
			<table class="min-w-full text-left text-sm">
				<thead class="bg-slate-50 text-slate-600">
					<tr>
						<th class="px-4 py-3 font-medium" scope="col">Name</th>
						<th class="px-4 py-3 font-medium" scope="col">Email</th>
					</tr>
				</thead>
				<tbody>
					{#each filtered as e}
						<tr class="border-t border-slate-100 odd:bg-white even:bg-slate-50/70 hover:bg-teal-50/40">
							<td class="px-4 py-3 font-medium">{e.display_name}</td>
							<td class="max-w-xs break-all px-4 py-3 text-slate-600">{e.email}</td>
						</tr>
					{:else}
						<tr><td class="px-4 py-6 text-slate-500" colspan="2">No matches.</td></tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</PageShell>

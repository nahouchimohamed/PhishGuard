<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';

	let { data } = $props();
	let q = $state('');
	const filtered = $derived(
		data.templates.filter((t) => {
			const needle = q.trim().toLowerCase();
			if (!needle) return true;
			return (
				t.name.toLowerCase().includes(needle) || t.subject.toLowerCase().includes(needle)
			);
		})
	);
</script>

<PageShell title="Templates" user={data.user} width="max-w-5xl">
	{#snippet banner()}Seeded training templates — browse before you launch.{/snippet}

	<div class="mb-6 flex flex-wrap items-end justify-between gap-4">
		<div>
			<h1 class="page-title">Templates</h1>
			<p class="page-subtitle">Subjects stay tagged [Training] · {filtered.length} shown</p>
		</div>
		<label class="block text-sm">
			<span class="mb-1 block font-medium text-slate-600">Search</span>
			<input class="input w-64" bind:value={q} placeholder="Name or subject" autofocus />
		</label>
	</div>

	{#if data.error}
		<Alert>{data.error}</Alert>
	{:else if filtered.length === 0}
		<div class="card text-slate-600">No templates match.</div>
	{:else}
		<div class="grid gap-3 sm:grid-cols-2">
			{#each filtered as t}
				<a
					class="card block transition hover:-translate-y-0.5 hover:border-accent/40 hover:shadow-panel"
					href={`/templates/${t.id}`}
				>
					<div class="flex items-start justify-between gap-2">
						<div class="font-medium text-ink">{t.name}</div>
						<span class="text-xs text-accent">View</span>
					</div>
					<p class="mt-2 line-clamp-2 text-sm leading-5 text-slate-600">{t.subject}</p>
				</a>
			{/each}
		</div>
	{/if}
</PageShell>

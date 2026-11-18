<script lang="ts">
	let { status }: { status: string } = $props();

	const tone = $derived(
		(
			{
				draft: 'bg-slate-100 text-slate-700',
				launching: 'bg-amber-50 text-amber-900 ring-1 ring-amber-200',
				active: 'bg-sky-50 text-sky-900',
				closed: 'bg-slate-700 text-slate-100'
			} as Record<string, string>
		)[status] ?? 'bg-slate-100 text-slate-700'
	);

	const label = $derived(
		({ draft: 'Draft', launching: 'Launching', active: 'Active', closed: 'Closed' } as Record<
			string,
			string
		>)[status] ?? status
	);
</script>

<span class={`inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium ${tone}`}
	>{#if status === 'launching'}<span class="h-1.5 w-1.5 animate-pulse rounded-full bg-amber-600"></span>{/if}{label}</span
>

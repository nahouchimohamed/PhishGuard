<script lang="ts">
	let {
		tone = 'danger',
		children
	}: {
		tone?: 'danger' | 'warning' | 'success' | 'info';
		children: import('svelte').Snippet;
	} = $props();

	const cls = $derived(
		{
			danger: 'border-rose-200 bg-rose-50 text-rose-900',
			warning: 'border-amber-200 bg-amber-50 text-amber-950',
			success: 'border-emerald-200 bg-emerald-50 text-emerald-900',
			info: 'border-sky-200 bg-sky-50 text-sky-950'
		}[tone]
	);
</script>

<div class={`rounded-lg border px-3 py-2 text-sm ${cls}`} role="alert">
	{@render children()}
</div>

<script lang="ts">
	let {
		tone = 'error',
		children
	}: {
		tone?: 'error' | 'status' | 'warning';
		children: import('svelte').Snippet;
	} = $props();

	const role = $derived(tone === 'error' ? 'alert' : 'status');
	const classes = $derived(
		tone === 'error'
			? 'rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700'
			: tone === 'warning'
				? 'rounded-lg bg-amber-50 px-3 py-2 text-sm text-amber-900'
				: 'rounded-lg bg-slate-50 px-3 py-2 text-sm text-slate-700'
	);
</script>

<p class={classes} {role}>{@render children()}</p>

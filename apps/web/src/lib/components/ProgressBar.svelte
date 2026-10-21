<script lang="ts">
	let { value, max = 100, label = '' }: { value: number; max?: number; label?: string } = $props();
	const pct = $derived(max > 0 ? Math.min(100, Math.round((value / max) * 1000) / 10) : 0);
</script>

<div class="w-full">
	{#if label}
		<div class="mb-1 flex justify-between text-xs text-slate-500">
			<span>{label}</span>
			<span>{pct}%</span>
		</div>
	{/if}
	<div class="h-2 overflow-hidden rounded-full bg-slate-100">
		<div
			class="h-full rounded-full bg-accent transition-all duration-300 ease-out"
			style={`width: ${pct}%`}
		></div>
	</div>
</div>

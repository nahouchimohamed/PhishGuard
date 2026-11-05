<script lang="ts">
	let {
		label,
		value,
		tone = 'neutral'
	}: {
		label: string;
		value: string | number;
		tone?: 'good' | 'warn' | 'bad' | 'neutral';
	} = $props();

	const valueClass = $derived(
		{
			good: 'text-emerald-700',
			warn: 'text-amber-700',
			bad: 'text-rose-700',
			neutral: 'text-ink'
		}[tone]
	);
</script>

<div class="card px-3 py-4 text-center">
	<div class={`text-2xl font-semibold tabular-nums ${valueClass}`}>{value}</div>
	<div class="mt-1 text-[11px] uppercase tracking-wide text-slate-500">{label}</div>
</div>

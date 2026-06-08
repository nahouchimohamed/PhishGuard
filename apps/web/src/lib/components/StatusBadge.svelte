<script lang="ts">
	import type { RecipientLabel } from '$lib/labels';

	let { label }: { label: RecipientLabel } = $props();

	const tone = $derived(
		{
			Reported: 'bg-emerald-50 text-emerald-800',
			Clicked: 'bg-rose-50 text-rose-800',
			Opened: 'bg-sky-50 text-sky-800',
			Delivered: 'bg-slate-100 text-slate-700',
			Pending: 'bg-amber-50 text-amber-900'
		}[label]
	);
</script>

<span class={`inline-flex rounded-md px-2 py-0.5 text-xs font-medium ${tone}`}>{label}</span>

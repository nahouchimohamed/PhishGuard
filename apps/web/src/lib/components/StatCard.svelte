<script lang="ts">
	let { label, value }: { label: string; value: string | number } = $props();
</script>

<div class="card text-center">
	<div class="text-2xl font-semibold">{value}</div>
	<div class="text-xs uppercase tracking-wide text-slate-500">{label}</div>
</div>

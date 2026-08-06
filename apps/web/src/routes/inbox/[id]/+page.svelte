<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	import { recipientLabel } from '$lib/labels';

	let { data, form } = $props();

	const renderedBody = $derived.by(() => {
		const item = data.item;
		if (!item) return '';
		const link = item.sim_url ?? '#';
		return item.body_text.replaceAll('{{LINK}}', link);
	});
</script>

<div class="banner">SIMULATION — reporting is the goal. The link opens a lesson page, never a password form.</div>
<Nav user={data.user!} />

<main class="mx-auto max-w-3xl px-4 py-8">
	{#if !data.item}
		<p class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{data.error}</p>
	{:else}
		<p class="mb-4 text-sm">
			<a class="text-accent hover:underline" href="/inbox">← Inbox</a>
			<span class="ml-3 text-slate-500">{recipientLabel(data.item)}</span>
		</p>

		{#if form?.error}
			<p class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{form.error}</p>
		{/if}

		<article class="card">
			<div class="text-xs text-slate-500">{data.item.from_label}</div>
			<h1 class="mt-1 text-xl font-semibold">{data.item.subject}</h1>
			<pre
				class="mt-4 break-words whitespace-pre-wrap font-sans text-sm leading-6 text-slate-800"
			>{renderedBody}</pre>

			{#if data.item.sim_url}
				<p class="mt-4 text-sm">
					Training link:
					<a
						class="break-all text-accent underline"
						href={data.item.sim_url}
						data-sveltekit-preload-data="off">{data.item.sim_url}</a
					>
				</p>
			{/if}
		</article>

		<div class="mt-4 flex flex-wrap gap-3">
			{#if !data.item.reported_at}
				<form method="POST" action="?/report">
					<button class="btn" type="submit">Report as suspicious</button>
				</form>
			{:else}
				<span class="rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-800">Reported — nice work</span>
			{/if}
		</div>
	{/if}
</main>

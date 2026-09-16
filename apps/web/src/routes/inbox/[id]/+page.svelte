<script lang="ts">
	import { enhance } from '$app/forms';
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import SubmitButton from '$lib/components/SubmitButton.svelte';
	import { recipientLabel } from '$lib/labels';
	import { parseLesson } from '$lib/lesson';
	import { splitMessageBody } from '$lib/message';

	let { data, form } = $props();
	let pending = $state(false);

	const segments = $derived.by(() => {
		const item = data.item;
		if (!item) return [];
		return splitMessageBody(item.body_text, item.sim_url ?? '#');
	});
	const lesson = $derived(data.item?.reported_at ? parseLesson(data.item.lesson_text) : null);
</script>

<PageShell title={data.item?.subject ?? 'Message'} user={data.user} width="max-w-3xl">
	{#snippet banner()}
		SIMULATION — report first. After you report, the lesson unlocks below. Links never ask for a password.
	{/snippet}

	{#if !data.item}
		<Alert>{data.error}</Alert>
	{:else}
		<p class="mb-4 flex flex-wrap items-center gap-3 text-sm">
			<a class="text-accent hover:underline" href="/inbox">← Inbox</a>
			<StatusBadge label={recipientLabel(data.item)} />
		</p>

		{#if form?.error}
			<div class="mb-4"><Alert>{form.error}</Alert></div>
		{/if}

		<article class="card">
			<div class="text-xs text-slate-500">{data.item.from_label}</div>
			<h1 class="mt-1 text-xl font-semibold">{data.item.subject}</h1>
			<pre class="mt-4 break-words whitespace-pre-wrap font-sans text-sm leading-6 text-slate-800"
				>{#each segments as seg}{#if seg.type === 'text'}{seg.value}{:else}<a
							class="break-all text-accent underline"
							href={seg.href}
							data-sveltekit-preload-data="off">{seg.label}</a
						>{/if}{/each}</pre
			>
		</article>

		<div class="mt-4 flex flex-wrap gap-3">
			{#if !data.item.reported_at}
				<form
					method="POST"
					action="?/report"
					use:enhance={() => {
						pending = true;
						return async ({ update }) => {
							await update();
							pending = false;
						};
					}}
				>
					<SubmitButton label="Report as suspicious" pendingLabel="Reporting…" {pending} />
				</form>
			{:else}
				<span class="rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-800"
					>Reported — nice work</span
				>
			{/if}
		</div>

		{#if lesson}
			<section class="card mt-6 border-emerald-200 bg-emerald-50/50">
				<h2 class="mb-2 font-medium">What to watch for</h2>
				<p class="text-sm text-slate-800">{lesson.intro}</p>
				{#if lesson.redFlags.length}
					<ul class="mt-3 list-disc space-y-1 pl-5 text-sm text-slate-700">
						{#each lesson.redFlags as flag}<li>{flag}</li>{/each}
					</ul>
				{/if}
				{#if lesson.goodHabit}
					<p class="mt-3 text-sm"><span class="font-medium">Good habit:</span> {lesson.goodHabit}</p>
				{/if}
			</section>
		{/if}
	{/if}
</PageShell>

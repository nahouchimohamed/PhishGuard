<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import { parseLesson } from '$lib/lesson';

	let { data } = $props();
	const lesson = $derived(data.template ? parseLesson(data.template.lesson_text) : null);
</script>

<PageShell title={data.template?.name ?? 'Template'} user={data.user} width="max-w-3xl">
	{#snippet banner()}Preview only — launching still happens from New campaign.{/snippet}

	{#if !data.template}
		<Alert>{data.error}</Alert>
	{:else}
		<p class="mb-4 text-sm"><a class="text-accent hover:underline" href="/templates">← Templates</a></p>
		<h1 class="page-title">{data.template.name}</h1>
		<p class="page-subtitle mt-1">{data.template.subject}</p>

		<section class="card mt-6">
			<h2 class="mb-2 font-medium">Message body</h2>
			<pre class="whitespace-pre-wrap font-sans text-sm leading-6 text-slate-800">{data.template.body_text}</pre>
		</section>

		{#if lesson}
			<section class="card mt-4 border-accent/20 bg-teal-50/40">
				<h2 class="mb-2 font-medium">Awareness lesson</h2>
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

		<div class="mt-6">
			<a class="btn" href={`/campaigns/new?template=${encodeURIComponent(data.template.id)}`}
				>Use in new campaign</a
			>
		</div>
	{/if}
</PageShell>

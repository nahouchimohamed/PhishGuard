<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import { parseLesson } from '$lib/lesson';

	let { data } = $props();
	const parsed = $derived(data.lesson ? parseLesson(data.lesson.lesson_text) : null);
</script>

<PageShell title="Simulation result" user={data.user} width="max-w-2xl">
	{#snippet banner()}Simulated message. This page never asks for a password.{/snippet}

	{#if !data.lesson || !parsed}
		<Alert>{data.error}</Alert>
	{:else}
		<article class="card border-amber-200">
			<p class="text-xs uppercase tracking-wide text-amber-800">Awareness notice</p>
			<h1 class="mt-1 text-2xl font-semibold">{data.lesson.subject}</h1>
			<p class="mt-1 text-sm text-slate-600">Campaign: {data.lesson.name}</p>

			{#if parsed.intro}
				<p class="mt-4 text-sm leading-6 text-slate-800">{parsed.intro}</p>
			{/if}

			{#if parsed.redFlags.length}
				<div class="mt-4">
					<h2 class="text-sm font-semibold">Red flags</h2>
					<ul class="mt-2 list-disc space-y-1 pl-5 text-sm text-slate-800">
						{#each parsed.redFlags as flag}
							<li>{flag}</li>
						{/each}
					</ul>
				</div>
			{/if}

			{#if parsed.goodHabit}
				<div class="mt-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-900">
					<strong>Good habit:</strong>
					{parsed.goodHabit}
				</div>
			{/if}
		</article>

		<p class="mt-6">
			{#if data.user}
				<a class="btn" href="/inbox">Back to inbox</a>
			{:else}
				<a class="btn" href="/login">Sign in</a>
			{/if}
		</p>
	{/if}
</PageShell>

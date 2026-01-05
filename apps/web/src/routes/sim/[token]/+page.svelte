<script lang="ts">
	let { data } = $props();
</script>

<div class="banner">This was a PhishGuard simulation. There is no login form on purpose.</div>

<main class="mx-auto max-w-2xl px-4 py-10">
	{#if !data.lesson}
		<div class="card">
			<h1 class="text-xl font-semibold">Link not found</h1>
			<p class="mt-2 text-sm text-slate-600">{data.error}</p>
		</div>
	{:else}
		<div class="card border-amber-200">
			<p class="text-xs font-medium uppercase tracking-wide text-amber-700">Training complete</p>
			<h1 class="mt-2 text-2xl font-semibold">You clicked a simulated phishing link</h1>
			<p class="mt-1 text-sm text-slate-600">Template: {data.lesson.name}</p>
			<p class="mt-1 text-sm text-slate-500">{data.lesson.subject}</p>
			<pre class="mt-6 whitespace-pre-wrap font-sans text-sm leading-6 text-slate-800"
				>{data.lesson.lesson_text}</pre
			>
			<a class="btn mt-6 inline-flex" href="/inbox">Back to inbox</a>
		</div>
	{/if}
</main>

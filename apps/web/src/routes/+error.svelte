<script lang="ts">
	import { page } from '$app/state';

	const message = $derived(page.error?.message || 'Something went wrong.');
	const id = $derived(page.error && 'id' in page.error ? String((page.error as { id?: string }).id ?? '') : '');
</script>

<main class="mx-auto flex min-h-[60vh] max-w-lg flex-col justify-center px-4">
	<div class="card">
		<h1 class="text-2xl font-semibold">Error {page.status}</h1>
		<p class="mt-2 text-sm text-slate-700" role="alert">{message}</p>
		{#if id}
			<p class="mt-2 text-xs text-slate-500">Reference: {id}</p>
		{/if}
		<p class="mt-6">
			<a class="btn" href="/">Back home</a>
		</p>
	</div>
</main>

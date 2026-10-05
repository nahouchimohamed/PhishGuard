<script lang="ts">
	import { page } from '$app/state';
</script>

<main class="mx-auto flex min-h-[60vh] max-w-lg flex-col justify-center px-4 py-16 text-center">
	<p class="text-sm font-medium uppercase tracking-wide text-slate-500">Something went wrong</p>
	<h1 class="page-title mt-2">{page.status}</h1>
	<p class="page-subtitle mt-2">{page.error?.message ?? 'Unexpected error'}</p>
	<a class="btn mx-auto mt-8" href="/">Back home</a>
</main>

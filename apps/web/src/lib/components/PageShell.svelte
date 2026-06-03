<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	import Banner from '$lib/components/Banner.svelte';
	import type { User } from '$lib/types';

	let {
		title,
		user = null,
		width = 'max-w-5xl',
		banner,
		children
	}: {
		title: string;
		user?: User | null;
		width?: 'max-w-md' | 'max-w-2xl' | 'max-w-3xl' | 'max-w-5xl';
		banner: import('svelte').Snippet;
		children: import('svelte').Snippet;
	} = $props();
</script>

<svelte:head>
	<title>{title} · PhishGuard</title>
</svelte:head>

<a
	class="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded focus:bg-white focus:px-3 focus:py-2"
	href="#main">Skip to content</a
>

<Banner>{@render banner()}</Banner>
{#if user}
	<Nav {user} />
{/if}

<main id="main" class={`mx-auto ${width} px-4 py-8`}>
	{@render children()}
</main>

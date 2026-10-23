<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<svelte:head>
	<title>PhishGuard</title>
</svelte:head>

<div class="min-h-screen">
	{@render children()}
</div>

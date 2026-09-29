<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<svelte:head>
	<title>PhishGuard</title>
</svelte:head>

{@render children()}

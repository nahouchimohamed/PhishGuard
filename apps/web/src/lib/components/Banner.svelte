<script lang="ts">
	let { children }: { children: import('svelte').Snippet } = $props();
</script>

<div class="banner sticky top-0 z-40">
	{@render children()}
</div>

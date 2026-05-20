<script lang="ts">
	let { children }: { children: import('svelte').Snippet } = $props();
</script>

<div class="banner" role="note">
	{@render children()}
</div>

<script lang="ts">
	let { children }: { children: import('svelte').Snippet } = $props();
</script>

<!-- Persistent training banner — keep visible on every authenticated page. -->
<div class="banner" role="note">
	{@render children()}
</div>

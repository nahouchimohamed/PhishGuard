<script lang="ts">
	import type { User } from '$lib/types';

	let { user }: { user: User } = $props();
</script>

<header class="border-b border-slate-200 bg-white">
	<div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
		<a href="/" class="font-semibold tracking-tight">PhishGuard</a>
		<div class="flex flex-wrap items-center gap-4 text-sm">
			<span class="text-slate-600">{user.display_name} · {user.role}</span>
			{#if user.role === 'admin'}
				<a class="hover:text-accent" href="/campaigns">Campaigns</a>
			{:else}
				<a class="hover:text-accent" href="/inbox">Inbox</a>
			{/if}
			<form method="POST" action="/logout">
				<button class="btn-secondary" type="submit">Log out</button>
			</form>
		</div>
	</div>
</header>

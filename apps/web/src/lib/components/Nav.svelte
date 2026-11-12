<script lang="ts">
	import { page } from '$app/state';
	import type { User } from '$lib/types';

	let { user }: { user: User } = $props();

	function current(href: string): 'page' | undefined {
		return page.url.pathname === href || page.url.pathname.startsWith(href + '/')
			? 'page'
			: undefined;
	}
</script>

<header class="border-b border-slate-200/80 bg-white/90 backdrop-blur">
	<div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
		<div class="flex items-center gap-2">
			<a href="/" class="text-lg font-semibold tracking-tight text-accent">PhishGuard</a>
			<span class="hidden rounded-full bg-amber-50 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-amber-900 sm:inline"
				>Training</span
			>
		</div>
		<nav class="flex flex-wrap items-center gap-4 text-sm" aria-label="Primary">
			<span class="text-slate-600">{user.display_name} · {user.role}</span>
			{#if user.role === 'admin'}
				<a class="hover:text-accent" href="/overview" aria-current={current('/overview')}
					>Overview</a
				>
				<a class="hover:text-accent" href="/campaigns" aria-current={current('/campaigns')}
					>Campaigns</a
				>
				<a class="hover:text-accent" href="/templates" aria-current={current('/templates')}
					>Templates</a
				>
				<a class="hover:text-accent" href="/employees" aria-current={current('/employees')}
					>Employees</a
				>
			{:else}
				<a class="hover:text-accent" href="/inbox" aria-current={current('/inbox')}>Inbox</a>
			{/if}
			<a class="hover:text-accent" href="/account" aria-current={current('/account')}>Account</a>
			<form method="POST" action="/logout">
				<button class="btn-secondary px-3 py-1.5 text-xs" type="submit">Log out</button>
			</form>
		</nav>
	</div>
</header>

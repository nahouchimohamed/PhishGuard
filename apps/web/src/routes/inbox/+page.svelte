<script lang="ts">
	import Nav from '$lib/components/Nav.svelte';
	import { recipientLabel } from '$lib/labels';
	let { data } = $props();
</script>

<div class="banner">SIMULATION INBOX — these messages are training drills, not real mail.</div>
<Nav user={data.user!} />

<main class="mx-auto max-w-3xl px-4 py-8">
	<h1 class="text-2xl font-semibold">Inbox</h1>
	<p class="mt-1 text-sm text-slate-600">Open a message, then Report it or follow the training link.</p>

	{#if data.error}
		<p class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{data.error}</p>
	{:else if data.items.length === 0}
		<div class="card mt-6 text-slate-600">
			No delivered messages yet. Ask an admin to launch a campaign, then refresh.
		</div>
	{:else}
		<ul class="mt-6 divide-y divide-slate-100 overflow-hidden rounded-xl border border-slate-200 bg-white">
			{#each data.items as item}
				<li>
					<a
						class="flex items-start justify-between gap-4 px-4 py-3 hover:bg-slate-50"
						href={`/inbox/${item.recipient_id}`}
						data-sveltekit-preload-data="off"
					>
						<div>
							<div class="font-medium">{item.subject}</div>
							<div class="text-xs text-slate-500">{item.from_label}</div>
						</div>
						<span class="shrink-0 text-xs font-medium text-slate-600">{recipientLabel(item)}</span>
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</main>

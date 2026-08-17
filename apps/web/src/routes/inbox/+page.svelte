<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import PageShell from '$lib/components/PageShell.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import { formatDateTime } from '$lib/format';
	import { recipientLabel } from '$lib/labels';

	let { data } = $props();
	const unread = $derived(data.items.filter((i) => !i.opened_at).length);
</script>

<PageShell title="Inbox" user={data.user} width="max-w-3xl">
	{#snippet banner()}SIMULATION INBOX — these messages are training drills, never real mail.{/snippet}

	<h1 class="page-title">Inbox</h1>
	<p class="page-subtitle mt-1">
		{data.items.length} message{data.items.length === 1 ? '' : 's'}
		{#if unread > 0}
			· {unread} unread
		{/if}
	</p>
	<p class="mt-1 text-sm text-slate-600">Open a message, then Report it or follow the training link.</p>

	{#if data.error}
		<div class="mt-4"><Alert>{data.error}</Alert></div>
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
						class:bg-teal-50={!item.opened_at}
						href={`/inbox/${item.recipient_id}`}
						data-sveltekit-preload-data="off"
					>
						<div>
							<div class="font-medium" class:font-semibold={!item.opened_at}>{item.subject}</div>
							<div class="text-xs text-slate-500">
								{item.from_label}
								{#if item.delivered_at}
									· {formatDateTime(item.delivered_at)}
								{/if}
							</div>
						</div>
						<StatusBadge label={recipientLabel(item)} />
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</PageShell>

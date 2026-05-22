<script lang="ts">
	let {
		label = 'Submit',
		pendingLabel = 'Working…',
		pending = false,
		class: className = 'btn',
		formaction
	}: {
		label?: string;
		pendingLabel?: string;
		pending?: boolean;
		class?: string;
		formaction?: string;
	} = $props();
</script>

<button class={className} type="submit" disabled={pending} {formaction}>
	{pending ? pendingLabel : label}
</button>

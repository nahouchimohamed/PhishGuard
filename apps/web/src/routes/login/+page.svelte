<script lang="ts">
	let { form } = $props();
</script>

<div class="banner">PhishGuard is a local training simulator. No real email is sent.</div>

<main class="mx-auto flex min-h-[80vh] max-w-md flex-col justify-center px-4">
	<div class="card">
		<h1 class="text-2xl font-semibold">Sign in</h1>
		<p class="mt-1 text-sm text-slate-600">
			Seeded users share password <code class="rounded bg-slate-100 px-1">phishguard-dev</code>
		</p>

		{#if form?.error}
			<p class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{form.error}</p>
		{/if}

		<form method="POST" class="mt-6 space-y-4">
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Email</span>
				<input
					class="input"
					name="email"
					type="email"
					required
					value={form?.email ?? 'admin@phishguard.local'}
				/>
			</label>
			<label class="block text-sm">
				<span class="mb-1 block font-medium">Password</span>
				<input class="input" name="password" type="password" required value="phishguard-dev" />
			</label>
			<button class="btn w-full" type="submit">Continue</button>
		</form>

		<ul class="mt-6 space-y-1 text-xs text-slate-500">
			<li>Admin: admin@phishguard.local</li>
			<li>Employees: alice@ / bob@ / cara@phishguard.local</li>
		</ul>
	</div>
</main>

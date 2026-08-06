import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';
import { describe, expect, it } from 'vitest';

const srcRoot = join(__dirname);

function listSvelteFiles(dir: string): string[] {
	const out: string[] = [];
	for (const name of readdirSync(dir)) {
		const full = join(dir, name);
		const st = statSync(full);
		if (st.isDirectory()) {
			if (name === 'node_modules' || name === '.svelte-kit') continue;
			out.push(...listSvelteFiles(full));
		} else if (name.endsWith('.svelte')) {
			out.push(full);
		}
	}
	return out;
}

describe('training safety scan', () => {
	const files = listSvelteFiles(srcRoot);

	it('keeps password inputs only on the login page', () => {
		const offenders = files.filter((f) => {
			const rel = relative(srcRoot, f).replace(/\\/g, '/');
			if (rel === 'routes/login/+page.svelte') return false;
			return /type\s*=\s*["']password["']/.test(readFileSync(f, 'utf8'));
		});
		expect(offenders).toEqual([]);
	});

	it('never uses {@html} in Svelte templates', () => {
		const offenders = files.filter((f) => /\{@html\b/.test(readFileSync(f, 'utf8')));
		expect(offenders).toEqual([]);
	});

	it('keeps the simulation page free of forms and inputs', () => {
		const sim = files.filter((f) => relative(srcRoot, f).replace(/\\/g, '/').includes('routes/sim/'));
		expect(sim.length).toBeGreaterThan(0);
		for (const f of sim) {
			const text = readFileSync(f, 'utf8');
			expect(text).not.toMatch(/<form\b/i);
			expect(text).not.toMatch(/<input\b/i);
			expect(text).not.toMatch(/<textarea\b/i);
		}
	});
});

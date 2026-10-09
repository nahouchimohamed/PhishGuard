/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				ink: '#0f172a',
				panel: '#f8fafc',
				accent: '#0f766e',
				surface: '#ffffff'
			},
			boxShadow: {
				panel: '0 1px 2px rgb(15 23 42 / 0.06), 0 8px 24px rgb(15 23 42 / 0.04)'
			}
		}
	},
	plugins: []
};

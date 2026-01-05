/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				ink: '#0f172a',
				panel: '#f8fafc',
				accent: '#0f766e'
			}
		}
	},
	plugins: []
};

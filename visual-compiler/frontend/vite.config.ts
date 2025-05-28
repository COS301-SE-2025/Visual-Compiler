import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	test: {
		globals:true,
		clearMocks:true,
		workspace: [
			{
				extends: './vite.config.ts',
				plugins: [svelteTesting()],
				test: {
					name: 'client',
					environment: 'jsdom',
					clearMocks: true,
					include: ['tests/**/*.{test,spec}.{js,ts}'],
					exclude: ['src/lib/server/**', 'node_modules/**', 'tests/**/*.server.{test,spec}.{js,ts}'],
					setupFiles: ['./vitest-setup-client.ts']
				}
			},
			{
				extends: './vite.config.ts',
				test: {
					name: 'server',
					environment: 'node',
					include: ['tests/**/*.server.{test,spec}.{js,ts}', 'src/**/*.server.{test,spec}.{js,ts}'],
					exclude: ['tests/**/*.{test,spec}.{js,ts}', 'node_modules/**'],
					setupFiles: ['./vitest-setup-server.ts']
				}
			}
		]
	}
});

import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig(({ mode }) => ({
	plugins: [sveltekit()],
	test: {
		globals: true,
		clearMocks: true,
		coverage: {
			exclude: [
				// Configuration files
				'svelte.config.js',
				'vite.config.ts',
				'vitest.config.ts',
				'**/*.config.{js,ts}',
				'eslint.config.js',
				'cypress.config.js',
				
				// Build artifacts and generated code
				'.svelte-kit/**',
				'build/**',
				'dist/**',
				'coverage/**',
				'static/**',
				
				// Type definitions (no runtime logic)
				'**/*.d.ts',
				'src/app.d.ts',
				'src/lib/types.ts',
				'src/lib/index.ts',
				
				// Test utilities and mocks
				'tests/**/*.svelte',
				'**/*{.,_}{test,spec,mock}.{js,ts,svelte}',
				'tests/unit-tests/mock-*.{js,ts,svelte}',
				'vitest-setup*.ts',
				
				// Dependencies
				'node_modules/**'
			]
		},
		workspace: [
			{
				extends: './vite.config.ts',
				plugins: [svelteTesting()],
				test: {
					name: 'client',
					environment: 'jsdom',
					clearMocks: true,
					include: ['tests/**/*.{test,spec}.{js,ts}'],
					exclude: [
						'src/lib/server/**',
						'node_modules/**',
						'tests/**/*.server.{test,spec}.{js,ts}'
					],
					setupFiles: ['./vitest-setup-client.ts']
				},
				resolve: {
					conditions: ['browser'] 
					},
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
	},
}));

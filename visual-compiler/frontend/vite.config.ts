import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig(({ mode }) => ({
	...(mode === 'development' && {
		server: {
			allowedHosts: [
				'visual-compiler-alb-1542446286.eu-north-1.elb.amazonaws.com',
				'visual-compiler.co.za',
				'www.visual-compiler.co.za'
			],
			port: 5173
		}
	}),
	plugins: [sveltekit()],
	test: {
		globals: true,
		clearMocks: true,
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
}));

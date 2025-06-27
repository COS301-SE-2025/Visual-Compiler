import svelte from 'eslint-plugin-svelte';
import prettier from 'eslint-plugin-prettier';
import unicorn from 'eslint-plugin-unicorn';
import svelteParser from 'svelte-eslint-parser';
import ts from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';

/** @type {import("eslint").Linter.FlatConfig[]} */
export default [
	{
		ignores: [
			'.svelte-kit/',
			'node_modules/',
			'src/lib/components/landing/Features.svelte',
			'src/lib/components/landing/Hero.svelte',
			'src/lib/components/landing/Walkthrough.svelte',
			'src/lib/components/main/Toolbox.svelte'
		]
	},
	{
		files: ['**/*.js', '**/*.ts'],
		languageOptions: {
			parser: tsParser,
			parserOptions: {
				sourceType: 'module',
				tsconfigRootDir: process.cwd()
			}
		},
		plugins: {
			'@typescript-eslint': ts,
			prettier,
			unicorn
		},
		rules: {
			'prettier/prettier': 'warn',
			'unicorn/filename-case': [
				'error',
				{
					case: 'kebabCase'
				}
			]
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				parser: tsParser,
				ecmaVersion: 'latest',
				sourceType: 'module',
				extraFileExtensions: ['.svelte']
			}
		},
		plugins: {
			svelte,
			unicorn,
			prettier
		},
		rules: {
			'prettier/prettier': 'warn',
			'unicorn/filename-case': [
				'error',
				{
					case: 'kebabCase'
				}
			]
		}
	}
];

import svelte from 'eslint-plugin-svelte';
import prettier from 'eslint-plugin-prettier';
import unicorn from 'eslint-plugin-unicorn';
import svelteParser from 'svelte-eslint-parser';

/** @type {import("eslint").Linter.Config[]} */
export default [
	{
		files: ['**/*.js', '**/*.ts', '**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				ecmaVersion: 'latest',
				sourceType: 'module',
				extraFileExtensions: ['.svelte']
			}
		},
		plugins: {
			svelte,
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
	}
];

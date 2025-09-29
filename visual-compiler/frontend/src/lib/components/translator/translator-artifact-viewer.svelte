<script lang="ts">
	import { AddToast } from '$lib/stores/toast';
	import { lexerState } from '$lib/stores/lexer';

	// Prop to receive the translated code
	export let translated_code: string[] = [];
	export let translationError: any = null;
	/**
	 * copy_to_clipboard
	 * @description Copies the generated code to the user's clipboard.
	 * @param {void}
	 * @returns {void}
	 */
	function copy_to_clipboard() {
		if (translated_code && translated_code.length > 0) {
			const code_as_string = translated_code.join('\n');
			navigator.clipboard
				.writeText(code_as_string)
				.then(() => {
					AddToast('Code copied! Your translated code is now in the clipboard', 'success');
				})
				.catch((err) => {
					console.error('Failed to copy text: ', err);
					AddToast('Copy failed: Unable to copy code to clipboard', 'error');
				});
		}
	}

	// Add reactive statement for translator data
	$: if ($lexerState?.translator_data?.code) {
		translated_code = $lexerState.translator_data.code;
	}
</script>

<div class="artifact-container">
	<div class="artifact-header">
		<h2 class="artifact-title">Translator Artefact</h2>
		{#if translated_code && translated_code.length > 0}
			<button class="copy-button" on:click={copy_to_clipboard} title="Copy to clipboard">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="18"
					height="18"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
					<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
				</svg>
				<span>Copy</span>
			</button>
		{/if}
	</div>

	<div class="artifact-viewer">
		{#if translationError}
			<div class="error-state">
				<div class="error-icon">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="48"
						height="48"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10" />
						<line x1="12" y1="8" x2="12" y2="12" />
						<line x1="12" y1="16" x2="12.01" y2="16" />
					</svg>
				</div>
				<h4>Translation Failed</h4>
				<p class="error-message">
					The translation could not be completed with the provided rules. Please check your input and
					try again.
				</p>
				<pre class="error-details">{translationError.message || String(translationError)}</pre>
			</div>
		{:else if translated_code && translated_code.length > 0}
			<div class="artifact-subheader">
				<h3>Target Code</h3>
			</div>
			<div class="code-wrapper">
				<div class="line-numbers">
					{#each translated_code as _, i}
						<span>{i + 1}</span>
					{/each}
				</div>
				<pre
					class="code-block"><code>{#each translated_code as line}{line}{'\n'}{/each}</code></pre>
			</div>
		{:else}
			<div class="empty-state">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="48"
					height="48"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
					<polyline points="14 2 14 8 20 8" />
					<line x1="16" y1="13" x2="8" y2="13" />
					<line x1="16" y1="17" x2="8" y2="17" />
					<polyline points="10 9 9 9 8 9" />
				</svg>
				<p>Generated code will appear here...</p>
			</div>
		{/if}
	</div>
</div>

<style>
	.artifact-container {
		padding: 1.5rem;
		background-color: #f8f9fa;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.artifact-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid #e5e7eb;
		padding-bottom: 0.75rem;
	}

	.artifact-viewer .artifact-subheader {
		border-bottom: none;
		padding-bottom: 0;
		margin-bottom: 1.5rem;
		text-align: center;
		color: white
	}

	.artifact-title {
		margin: 0;
		color: #001a6e;
		font-family: 'Times New Roman', serif;
		font-size: 1.25rem;
	}

	h3 {
		color: #001a6e;
		font-size: 1.5rem;
		margin: 0;
		font-family: 'Times New Roman', serif;
	}

	.copy-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		background-color: #eef2f7;
		color: #001a6e;
		border: 1px solid #c0c7d3;
		padding: 0.4rem 0.8rem;
		border-radius: 6px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.copy-button:hover {
		background-color: #dbe4f0;
		border-color: #001a6e;
	}

	.code-wrapper {
		display: flex;
		background-color: #fdfdfd;
		border-radius: 4px;
		border: 1px solid #e5e7eb;
		overflow: hidden;
	}

	.line-numbers {
		padding: 1rem 0.75rem;
		font-family: 'Fira Code', monospace;
		font-size: 0.9rem;
		text-align: right;
		color: #9ca3af;
		background-color: #f8f9fa;
		border-right: 1px solid #e5e7eb;
		user-select: none;
	}

	.line-numbers span {
		display: block;
		line-height: 1.5;
	}

	.code-block {
		margin: 0;
		padding: 1rem;
		flex-grow: 1;
		white-space: pre-wrap;
		word-wrap: break-word;
		font-family: 'Fira Code', monospace;
		font-size: 0.9rem;
		line-height: 1.5;
		color: #333;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 1rem;
		color: #6b7280;
		width: 100%;
		text-align: center;
	}

	.empty-state svg {
		margin-bottom: 1rem;
		opacity: 0.5;
	}

	/* --- Dark Mode --- */
	:global(html.dark-mode) .artifact-container {
		 background: #1a2a4a;
	}

	:global(html.dark-mode) .artifact-header {
		border-bottom-color: #4a5568;
	}
	:global(html.dark-mode) .artifact-title,
	:global(html.dark-mode) h3 {
		color: #ebeef1;
	}
  
	:global(html.dark-mode) .artifact-header {
		border-bottom-color: #4a5568;
	}
	:global(html.dark-mode) .artifact-title {
		color: #ebeef1;
	}
	:global(html.dark-mode) .copy-button {
		background-color: #4a5568;
		color: #e0e7ff;
		border-color: #6b7280;
	}
	:global(html.dark-mode) .copy-button:hover {
		background-color: #374151;
		border-color: #9ca3af;
	}
	:global(html.dark-mode) .code-wrapper {
		background-color: #1a202c;
		border-color: #4a5568;
	}
	:global(html.dark-mode) .line-numbers {
		color: #6b7280;
		background-color: #2d3748;
		border-right-color: #4a5568;
	}
	:global(html.dark-mode) .code-block {
		color: #d1d5db;
	}
	:global(html.dark-mode) .empty-state {
		color: #9ca3af;
	}

	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		padding: 2rem;
		flex-grow: 1;
		background-color: #fff5f5;
		border: 1px solid #e53e3e;
		border-radius: 8px;
		color: #9b2c2c;
	}
	.error-icon {
		color: #e53e3e;
		margin-bottom: 1rem;
	}
	.error-state h4 {
		margin: 0 0 0.5rem 0;
		font-size: 1.25rem;
		color: #c53030;
	}
	.error-message {
		margin: 0 0 1rem 0;
		max-width: 450px;
		line-height: 1.6;
	}
	.error-details {
		background-color: #fed7d7;
		padding: 0.75rem 1rem;
		border-radius: 6px;
		font-family: 'Fira Code', monospace;
		font-size: 0.85rem;
		white-space: pre-wrap;
		word-wrap: break-word;
		max-width: 100%;
		text-align: left;
		color: #742a2a;
	}
	:global(html.dark-mode) .error-state {
		background-color: #4a2d2d;
		border-color: #e53e3e;
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-icon {
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-state h4 {
		color: #fc8181;
	}
	:global(html.dark-mode) .error-details {
		background-color: #4a2d2d;
		color: #fed7d7;
	}
</style>


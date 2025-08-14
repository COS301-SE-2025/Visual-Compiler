<script lang="ts">
	import { AddToast } from '$lib/stores/toast';
	import { confirmedSourceCode } from '$lib/stores/source-code';
	import { tick } from 'svelte';
	import { onMount } from 'svelte';

	let code_text = '';
	
	let previous_code_text = '';
	let isDefaultInput = false;
	let isConfirmed = false;
	let textareaEl: HTMLTextAreaElement;
	export let onCodeSubmitted: (code: string) => void = () => {};

	// --- NEW MOCK DATA FOR PROJECTS ---
	const mockProjects = [
		{ name: 'Select a project...', code: '' },
		{ name: 'Lexer Project', code: '// Mock translated code for Lexer Project\n\nfunction lexer(code) {\n  // ...\n}' },
		{ name: 'Parser Project', code: '// Mock translated code for Parser Project\n\nfunction parser(tokens) {\n  // ...\n}' },
		{ name: 'Translator Project', code: '// Mock translated code for Translator Project\n\nfunction translator(ast) {\n  // ...\n}' }
	];
	let selectedProject = mockProjects[0];

	// Sync with global store
	let confirmed_code = '';
	const unsubscribe = confirmedSourceCode.subscribe(value => {
		confirmed_code = value;
		if (!code_text) code_text = value; // Auto fill input if empty
	});

	onMount(() => {
		return () => unsubscribe(); // clean up store subscription
	});

	function handleDefaultInput() {
		if (!isDefaultInput) {
			previous_code_text = code_text;
			code_text = 'int blue = 13 + 22;';
			isDefaultInput = true;
		} else {
			code_text = previous_code_text;
			isDefaultInput = false;
		}
	}

	function handleFileChange(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		if (!file.name.toLowerCase().endsWith('.txt')) {
			AddToast('Only .txt files are allowed. Please upload a valid plain text file.', 'error');
			input.value = '';
			return;
		}

		const reader = new FileReader();
		reader.onload = () => {
			code_text = reader.result as string;
			AddToast('File uploaded successfully!', 'success');
		};
		reader.onerror = () => {
			AddToast('Failed to read file.', 'error');
		};
		reader.readAsText(file);
	}

	async function submitCode() {
		if (!code_text.trim()) return;
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		try {
			const res = await fetch('http://localhost:8080/api/lexing/code', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					users_id: user_id,
					source_code: code_text
				})
			});
			if (!res.ok) throw new Error();
			AddToast('Code confirmed and saved!', 'success');
			confirmedSourceCode.set(code_text);
			isConfirmed = true;
			onCodeSubmitted(code_text);
			await tick();
		} catch {
			AddToast('Failed to save source code', 'error');
		}
	}

	// --- NEW FUNCTION TO HANDLE PROJECT SELECTION ---
	function handleProjectSelect() {
		if (selectedProject) {
			code_text = selectedProject.code;
			isDefaultInput = false; 
		}
	}

	// Reset confirmation flag if user changes the text
	$: isConfirmed = code_text === confirmed_code && !!code_text;
	$: displayed_text = isConfirmed ? `Current source code: ${code_text}` : code_text;
</script>

<div class="code-input-container">
	<div class="code-input-header-row">
		<h2 class="code-input-header">Source Code Input</h2>
		<button
			type="button"
			class="default-source-btn"
			title={isDefaultInput ? 'Restore your input' : 'Insert default source code'}
			aria-label={isDefaultInput ? 'Restore your input' : 'Insert default source code'}
			on:click={handleDefaultInput}
		>
			{#if isDefaultInput} 🧹 {:else} 🪄 {/if}
		</button>
	</div>

	{#if isConfirmed}
		<p class="current"><strong>Current source code:</strong></p>
	{/if}

	<textarea
		bind:value={code_text}
		rows="10"
		placeholder="Paste or type your source code here…"
	></textarea>

	<div class="controls-grid">
		<div class="control-item">
			<label class="upload-btn">
				Upload File
				<input type="file" accept=".txt" on:change={handleFileChange} />
			</label>
		</div>

		<div class="control-item project-selector">
			<label for="project-select">Import from Project</label>
			<select id="project-select" bind:value={selectedProject} on:change={handleProjectSelect}>
				{#each mockProjects as project}
					<option value={project}>{project.name}</option>
				{/each}
			</select>
		</div>
	</div>


	<div class="controls">
		<button
			type="button"
			class="confirm-btn"
			on:click={submitCode}
			disabled={!code_text.trim()}
		>
			
			{#if isConfirmed}
				Code Confirmed
				<span class="tick">✔</span>
			{:else}
				Confirm Code
			{/if}
		</button>
	</div>

	
</div>

<style>
	.current {
		margin-top: 0;
		margin-bottom: 0.25rem;
	}

	.code-input-container {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		margin: 1rem;
	}

	textarea {
		resize: vertical;
		padding: 0.4rem;
		font-family: Menlo, Monaco, 'Courier New', monospace;
		font-size: 0.9rem;
		line-height: 1.2;
		border: 1px solid #ccc;
		border-radius: 4px;
		box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.1);
		height: 100px;
		margin-bottom: 0.5rem;
	}

	.controls {
		display: flex;
		gap: 0.5rem;
		margin: 0.3rem 0 0.7rem;
		justify-content: center;
	}

	.controls-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
		align-items: flex-end; /* Align items to the bottom */
		margin-bottom: 1rem;
	}

	.control-item {
		display: flex;
		flex-direction: column;
	}

	.project-selector label {
		font-size: 0.9rem;
		font-weight: 500;
		margin-bottom: 0.5rem;
		color: #333;
	}

	.project-selector select {
		padding: 0.5rem 0.8rem;
		border-radius: 4px;
		border: 1px solid #ccc;
		background-color: white;
		font-size: 0.95rem;
		cursor: pointer;
		-webkit-appearance: none;
		appearance: none;
		background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right 0.7rem center;
		background-size: 1em;
	}
	
	.project-selector select:hover {
		border-color: #888;
	}
	
	.upload-btn {
		position: relative;
		overflow: hidden;
		display: inline-block;
		padding: 0.5rem 1rem;
		background: #646468;
		color: white;
		border-radius: 4px;
		font-size: 0.95rem;
		cursor: pointer;
		text-align: center;
	}

	.upload-btn input[type='file'] {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		opacity: 0;
		cursor: pointer;
	}

	.upload-btn:hover {
		background: #838386;
	}

	.confirm-btn {
		padding: 0.5rem 1.5rem;
		background: #BED2E6;
		color: 000000;
		border: none;
		border-radius: 4px;
		font-size: 0.95rem;
		cursor: pointer;
		transition: background 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.3rem;
	}

	.confirm-btn:disabled {
		background: #e6e6e6;
		color: #666666;
		cursor: not-allowed;
	}

	.confirm-btn:not(:disabled):hover {
		background: #a8bdd1;
	}

	.tick {
		font-size: 1.2rem;
		line-height: 1;
		margin-right: -0.45rem;
		margin-left: 0.25rem;
	}

	.default-source-btn {
		background: #e0e7ff;
		color: #1e40af;
		border: none;
		border-radius: 50%;
		width: 2.2rem;
		height: 2.2rem;
		font-size: 1.3rem;
		cursor: pointer;
		transition: background 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.default-source-btn:hover,
	.default-source-btn:focus {
		background: #d0d9ff;
	}

	.code-input-header-row {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		margin-bottom: 0.3rem;
		justify-content: center;
	}

	.code-input-header {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 700;
	}

	:global(html.dark-mode) textarea {
		background: #2d3748;
		color: #ccc;
		border-color: #4a5568;
	}

	:global(html.dark-mode) .project-selector label {
		color: #ccc;
	}

	:global(html.dark-mode) .project-selector select {
		background: #2d3748;
		color: #ccc;
		border-color: #4a5568;
		background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23cccccc' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
	}

	:global(html.dark-mode) .project-selector select:hover {
		border-color: #6b7280;
	}
	
	:global(html.dark-mode) select option {
		background-color: #2d3748;
		color: #ccc;
	}
	
	:global(html.dark-mode) select option:hover {
		background-color: #4a5568;
	}

	 :global(html.dark-mode) .default-source-btn  {
        background-color: #2d3748;
        border-color: #4a5568;
        color: #d1d5db;
    }

	 :global(html.dark-mode) .default-source-btn.selected {
        background-color: #001a6e;
        border-color: #60a5fa;
        color: #e0e7ff;
    }

    :global(html.dark-mode) .default-source-btn:not(.selected):hover {
        background-color: #374151;
        border-color: #6b7280;
    }

	:global(html.dark-mode) .confirm-btn {
		background: #001A6E;
		color: #ffffff;
		border: 1px solid #374151;
	}

	:global(html.dark-mode) .confirm-btn:not(:disabled):hover {
		background: #002a8e;
	}

	:global(html.dark-mode) .confirm-btn:disabled {
		background: #2d3748;
		color: #a0aec0;
		border-color: #4a5568;
	}
</style>
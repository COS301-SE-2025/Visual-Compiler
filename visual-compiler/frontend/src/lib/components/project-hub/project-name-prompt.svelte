<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly } from 'svelte/transition';
	import { AddToast } from '$lib/stores/toast';

	export let show: boolean;
	export let recentProjects: { name: string }[] = [];
	let projectName = '';

	const dispatch = createEventDispatcher();

	// Reset project name when modal opens
	$: if (show) {
		projectName = '';
	}

	function handleConfirm() {
		if (projectName.trim()) {
			// Check for duplicates before proceeding
			const isDuplicate = recentProjects.some(p => p.name === projectName.trim());
			
			if (isDuplicate) {
				AddToast('Project name already exists. Please choose a different name.', 'error');
				return;
			}
			
			dispatch('confirm', projectName.trim());
		}
	}

	function handleCancel() {
		dispatch('cancel');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleConfirm();
		} else if (event.key === 'Escape') {
			handleCancel();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
	<div class="backdrop" transition:fly={{ y: -50, duration: 300, opacity: 0.5 }}>
		<div class="prompt-modal" on:click|stopPropagation>
			<h3 class="prompt-heading">Create New Project</h3>
			<p class="prompt-subheading">Please enter a name for your new blank project.</p>

			<div class="input-container">
				<input
					type="text"
					bind:value={projectName}
					placeholder="e.g., My First Compiler"
					class="project-name-input"
					id="project-name-input"
					aria-label="Project Name"
				/>
			</div>

			<div class="button-container">
				<button class="cancel-button" on:click={handleCancel}>Cancel</button>
				<button
					class="create-button"
					on:click={handleConfirm}
					disabled={!projectName.trim()}
					id="confirm-project-name"
				>
					Create Project
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.6);
		-webkit-backdrop-filter: blur(5px);
		backdrop-filter: blur(5px);
		z-index: 10000;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.prompt-modal {
		background-color: #f1f5f9;
		border-radius: 0.75rem;
		padding: 2rem;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
		text-align: center;
	}

	.prompt-heading {
		font-size: 1.5rem;
		font-weight: 700;
		color: #1f2937;
		margin-bottom: 0.5rem;
	}

	.prompt-subheading {
		font-size: 1rem;
		color: #4b5567;
		margin-bottom: 1.5rem;
	}

	.project-name-input {
		width: 100%;
		padding: 0.75rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		font-size: 1rem;
		box-sizing: border-box;
		margin-bottom: 1.5rem;
	}

	.project-name-input:focus {
		outline: none;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.4);
		border-color: #3b82f6;
	}

	.button-container {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
	}

	.cancel-button,
	.create-button {
		padding: 0.6rem 1.2rem;
		border-radius: 0.375rem;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		border: none;
		transition: background-color 0.2s ease;
	}

	.cancel-button {
		background-color: #e5e7eb;
		color: #374151;
	}

	.cancel-button:hover {
		background-color: #d1d5db;
	}

	.create-button {
		background-color: #3b82f6;
		color: white;
	}

	.create-button:hover {
		background-color: #2563eb;
	}

	.create-button:disabled {
		background-color: #9ca3af;
		cursor: not-allowed;
	}

	:global(html.dark) .prompt-modal {
		background: #2d3748;
	}
	:global(html.dark) .prompt-heading {
		color: #edf2f7;
	}
	:global(html.dark) .prompt-subheading {
		color: #a0aec0;
	}
	:global(html.dark) .project-name-input {
		background-color: #4a5568;
		border-color: #718096;
		color: #edf2f7;
	}
	:global(html.dark) .cancel-button {
		background-color: #4a5568;
		color: #e2e8f0;
	}
	:global(html.dark) .cancel-button:hover {
		background-color: #2d3748;
	}
	:global(html.dark) .create-button {
		background-color: #4299e1;
	}
	:global(html.dark) .create-button:hover {
		background-color: #3182ce;
	}
	:global(html.dark) .create-button:disabled {
		background-color: #718096;
	}
</style>

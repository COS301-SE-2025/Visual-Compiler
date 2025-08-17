<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly } from 'svelte/transition';

	export let show: boolean;
	export let projectName: string;

	const dispatch = createEventDispatcher();

	function handleConfirm() {
		dispatch('confirm');
	}

	function handleCancel() {
		dispatch('cancel');
	}
</script>

{#if show}
	<div class="backdrop" transition:fly={{ y: -50, duration: 300, opacity: 0.5 }}>
		<div class="prompt-modal" on:click|stopPropagation>
			<h3 class="prompt-heading">Confirm Deletion</h3>
			<p class="prompt-subheading">
				Are you sure you want to delete the project "<strong>{projectName}</strong>"? This action cannot
				be undone.
			</p>

			<div class="button-container">
				<button class="cancel-button" on:click={handleCancel}>Cancel</button>
				<button class="delete-confirm-button" on:click={handleConfirm}> Delete </button>
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
		z-index: 10001; /* Higher z-index to appear above the project hub */
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
		line-height: 1.5;
	}

	.button-container {
		display: flex;
		justify-content: center;
		gap: 0.75rem;
	}

	.cancel-button,
	.delete-confirm-button {
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

	.delete-confirm-button {
		background-color: #ef4444;
		color: white;
	}

	.delete-confirm-button:hover {
		background-color: #dc2626;
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
	:global(html.dark) .cancel-button {
		background-color: #4a5568;
		color: #e2e8f0;
	}
	:global(html.dark) .cancel-button:hover {
		background-color: #2d3748;
	}
	:global(html.dark) .delete-confirm-button {
		background-color: #f87171;
	}
	:global(html.dark) .delete-confirm-button:hover {
		background-color: #ef4444;
	}
</style>

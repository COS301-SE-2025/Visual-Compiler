<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fly } from 'svelte/transition';

	export let show = false;
	export let hasNodes = false;
	export let projectName = '';

	const dispatch = createEventDispatcher();

	function confirm() {
		dispatch('confirm');
	}

	function cancel() {
		dispatch('cancel');
	}
</script>

{#if show}
	<div class="backdrop" transition:fly={{ y: -50, duration: 300, opacity: 0.5 }} on:click={cancel}>
		<div class="modal" on:click|stopPropagation>
			<h3>Clear Canvas</h3>

			{#if hasNodes && projectName}
				<p>This will clear all nodes and connections from the canvas.</p>
				<p><strong>Note:</strong> Since your canvas has nodes, this will also delete the project "{projectName}" and recreate it fresh to ensure all data is completely cleared from the database.</p>
				<p class="warning-text">This action cannot be undone!</p>
			{:else if hasNodes}
				<p>This will clear all nodes and connections from the canvas.</p>
				<p class="warning-text">This action cannot be undone!</p>
			{:else}
				<p>The canvas is already empty.</p>
			{/if}

			<div class="button-container">
				<button class="cancel-btn" on:click={cancel}>Cancel</button>
				{#if hasNodes}
					<button class="confirm-btn" on:click={confirm}>Clear Canvas</button>
				{:else}
					<button class="confirm-btn" on:click={confirm}>OK</button>
				{/if}
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
        z-index: 10001; /* Higher z-index to appear above the main workspace */
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .modal {
        background-color: #f1f5f9;
        border-radius: 0.75rem;
        padding: 2rem;
        width: 100%;
        max-width: 400px;
        box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.2);
        text-align: center;
    }

    h3 {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1f2937;
        margin-bottom: 0.5rem;
    }

    p {
        font-size: 1rem;
        color: #4b5567;
        margin-bottom: 1.5rem;
        line-height: 1.5;
    }

    .warning-text {
        color: #dc2626;
        font-weight: 500;
        font-style: italic;
    }

    .button-container {
        display: flex;
        justify-content: center;
        gap: 0.75rem;
    }

    .cancel-btn,
    .confirm-btn {
        padding: 0.6rem 1.2rem;
        border-radius: 0.375rem;
        font-size: 0.9rem;
        font-weight: 600;
        cursor: pointer;
        border: none;
        transition: background-color 0.2s ease;
    }

    .cancel-btn {
        background-color: #e5e7eb;
        color: #374151;
    }

    .cancel-btn:hover {
        background-color: #d1d5db;
    }

    .confirm-btn {
        background-color: #dc2626; /* Changed from blue to red */
        color: white;
    }

    .confirm-btn:hover {
        background-color: #b91c1c; /* Darker red on hover */
    }

    :global(html.dark-mode) .modal {
        background: #1f2937; /* Dark gray background instead of white */
    }

    :global(html.dark-mode) h3 {
        color: #f9fafb; /* White text for heading */
    }

    :global(html.dark-mode) p {
        color: #d1d5db; /* Light gray text for subheading */
    }

    :global(html.dark-mode) .cancel-btn {
        background-color: #4a5568; /* Dark gray button */
        color: #f9fafb; /* White text */
    }

    :global(html.dark-mode) .cancel-btn:hover {
        background-color: #2d3748; /* Darker gray on hover */
    }

    :global(html.dark-mode) .confirm-btn {
        background-color: #dc2626; /* Red background for clear action */
        color: #ffffff; /* White text */
    }

    :global(html.dark-mode) .confirm-btn:hover {
        background-color: #b91c1c; /* Darker red on hover */
    }

    /* ADD: Dark mode backdrop for better contrast */
    :global(html.dark-mode) .backdrop {
        background: rgba(0, 0, 0, 0.8); /* Darker backdrop in dark mode */
    }
</style>

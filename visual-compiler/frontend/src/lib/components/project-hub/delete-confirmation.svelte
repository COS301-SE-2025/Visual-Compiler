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
        background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
        border-radius: 0.75rem;
        padding: 2rem;
        width: 100%;
        max-width: 400px;
        box-shadow: 
            0 25px 50px -12px rgba(0, 0, 0, 0.25),
            0 10px 25px -5px rgba(0, 0, 0, 0.1),
            0 0 0 1px rgba(148, 163, 184, 0.1);
        text-align: center;
        border: 1px solid rgba(203, 213, 225, 0.5);
        transition: all 0.3s ease;
    }

    .prompt-heading {
        font-size: 1.5rem;
        font-weight: 700;
        color: #1f2937;
        margin-bottom: 0.5rem;
        background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        text-shadow: 0 2px 4px rgba(220, 38, 38, 0.1);
    }

    .prompt-subheading {
        font-size: 1rem;
        color: #64748b;
        margin-bottom: 1.5rem;
        line-height: 1.6;
        font-weight: 500;
        opacity: 0.9;
    }

    .prompt-subheading strong {
        color: #1f2937;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .button-container {
        display: flex;
        justify-content: center;
        gap: 0.75rem;
    }

    .cancel-button,
    .delete-confirm-button {
        padding: 0.75rem 1.5rem;
        border-radius: 0.5rem;
        font-size: 0.9rem;
        font-weight: 600;
        cursor: pointer;
        border: none;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
    }

    .cancel-button {
        background: linear-gradient(145deg, #f1f5f9 0%, #e2e8f0 100%);
        color: #475569;
        border: 2px solid #cbd5e1;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .cancel-button:hover {
        background: linear-gradient(145deg, #e2e8f0 0%, #cbd5e1 100%);
        border-color: #94a3b8;
        color: #334155;
        transform: translateY(-2px);
        box-shadow: 
            0 4px 12px rgba(0, 0, 0, 0.15),
            0 2px 4px rgba(71, 85, 105, 0.1);
    }

    .cancel-button:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .delete-confirm-button {
        background: linear-gradient(145deg, #ef4444 0%, #dc2626 100%);
        border: 2px solid #b91c1c;
        color: #ffffff;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        box-shadow: 
            0 4px 12px rgba(239, 68, 68, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .delete-confirm-button:hover {
        background: linear-gradient(145deg, #dc2626 0%, #b91c1c 100%);
        border-color: #991b1b;
        transform: translateY(-2px);
        box-shadow: 
            0 8px 25px rgba(239, 68, 68, 0.4),
            0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .delete-confirm-button:active {
        transform: translateY(-1px);
        box-shadow: 
            0 4px 12px rgba(239, 68, 68, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    /* Focus states for accessibility */
    .cancel-button:focus,
    .delete-confirm-button:focus {
        outline: 3px solid rgba(59, 130, 246, 0.3);
        outline-offset: 2px;
    }

    /* Mobile responsiveness */
    @media (max-width: 640px) {
        .prompt-modal {
            margin: 1rem;
            padding: 1.5rem;
        }
        
        .button-container {
            flex-direction: column;
            gap: 0.5rem;
        }
        
        .cancel-button,
        .delete-confirm-button {
            width: 100%;
            padding: 0.75rem 1rem;
        }
    }

    /* Enhanced Dark Mode Styles */
    :global(html.dark-mode) .backdrop {
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(8px);
        -webkit-backdrop-filter: blur(8px);
    }

    :global(html.dark-mode) .prompt-modal {
        background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
        border: 1px solid #334155;
        color: #f1f5f9;
        box-shadow:
            0 25px 50px -12px rgba(0, 0, 0, 0.6),
            0 0 0 1px rgba(148, 163, 184, 0.1);
    }

    :global(html.dark-mode) .prompt-heading {
        color: #f8fafc;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
        background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }

    :global(html.dark-mode) .prompt-subheading {
        color: #cbd5e1;
        opacity: 0.9;
        line-height: 1.6;
    }

    :global(html.dark-mode) .prompt-subheading strong {
        color: #f8fafc;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .cancel-button {
        background: linear-gradient(145deg, #475569 0%, #334155 100%);
        border: 2px solid #64748b;
        color: #f1f5f9;
        font-weight: 600;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .cancel-button:hover {
        background: linear-gradient(145deg, #64748b 0%, #475569 100%);
        border-color: #94a3b8;
        color: #ffffff;
        transform: translateY(-2px);
        box-shadow: 
            0 4px 12px rgba(0, 0, 0, 0.3),
            0 2px 4px rgba(100, 116, 139, 0.2);
    }

    :global(html.dark-mode) .cancel-button:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .delete-confirm-button {
        background: linear-gradient(145deg, #dc2626 0%, #b91c1c 100%);
        border: 2px solid #ef4444;
        color: #ffffff;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
        box-shadow: 
            0 4px 12px rgba(220, 38, 38, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    :global(html.dark-mode) .delete-confirm-button:hover {
        background: linear-gradient(145deg, #b91c1c 0%, #991b1b 100%);
        border-color: #f87171;
        transform: translateY(-2px);
        box-shadow: 
            0 8px 25px rgba(220, 38, 38, 0.4),
            0 4px 12px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .delete-confirm-button:active {
        transform: translateY(-1px);
        box-shadow: 
            0 4px 12px rgba(220, 38, 38, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }
</style>

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
        background: linear-gradient(145deg, #ffffff 0%, #fdfdfd 100%);
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
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .prompt-subheading {
        font-size: 1rem;
        color: #64748b;
        margin-bottom: 1.5rem;
        line-height: 1.6;
        font-weight: 500;
        opacity: 0.9;
    }

    .input-container {
        margin-bottom: 1.5rem;
    }

    .project-name-input {
        width: 100%;
        padding: 0.875rem 1rem;
        border: 2px solid #e2e8f0;
        border-radius: 0.5rem;
        font-size: 1rem;
        box-sizing: border-box;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        font-weight: 500;
        background: linear-gradient(145deg, #ffffff 0%, #fefefe 100%);
        color: #1f2937;
        box-shadow: 
            inset 0 1px 2px rgba(0, 0, 0, 0.05),
            0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .project-name-input::placeholder {
        color: #9ca3af;
        opacity: 0.7;
        font-style: italic;
        font-weight: 400;
    }

    .project-name-input:focus {
        outline: none;
        border-color: #3b82f6;
        background: linear-gradient(145deg, #fefefe 0%, #ffffff 100%);
        color: #111827;
        box-shadow: 
            0 0 0 4px rgba(59, 130, 246, 0.15),
            inset 0 1px 2px rgba(0, 0, 0, 0.05),
            0 4px 12px rgba(59, 130, 246, 0.1);
        transform: translateY(-1px);
    }

    .button-container {
        display: flex;
        justify-content: center;
        gap: 0.75rem;
    }

    .cancel-button,
    .create-button {
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
        background: linear-gradient(145deg, #f8fafc 0%, #f1f5f9 100%);
        color: #475569;
        border: 2px solid #cbd5e1;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .cancel-button:hover {
        background: linear-gradient(145deg, #f1f5f9 0%, #e2e8f0 100%);
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

    .create-button {
        background: linear-gradient(145deg, #3b82f6 0%, #2563eb 100%);
        border: 2px solid #1d4ed8;
        color: #ffffff;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        box-shadow: 
            0 4px 12px rgba(59, 130, 246, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .create-button:hover:not(:disabled) {
        background: linear-gradient(145deg, #2563eb 0%, #1e40af 100%);
        border-color: #1e40af;
        transform: translateY(-2px);
        box-shadow: 
            0 8px 25px rgba(59, 130, 246, 0.4),
            0 4px 12px rgba(0, 0, 0, 0.15);
    }

    .create-button:active:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 
            0 4px 12px rgba(59, 130, 246, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .create-button:disabled {
        background: linear-gradient(145deg, #d1d5db 0%, #c4c4c4 100%);
        border-color: #9ca3af;
        color: #6b7280;
        cursor: not-allowed;
        transform: none;
        box-shadow: none;
        opacity: 0.7;
        text-shadow: none;
    }

    /* Focus states for accessibility */
    .cancel-button:focus,
    .create-button:focus {
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
        .create-button {
            width: 100%;
            padding: 0.75rem 1rem;
        }
        
        .project-name-input {
            font-size: 16px; /* Prevents zoom on iOS */
        }
    }

    /* Enhanced Dark Mode Styles - Reduced gradients */
    :global(html.dark-mode) .backdrop {
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(8px);
        -webkit-backdrop-filter: blur(8px);
    }

    :global(html.dark-mode) .prompt-modal {
        background: linear-gradient(145deg, #1e293b 0%, #1a202c 100%);
        border: 1px solid #334155;
        color: #f1f5f9;
        box-shadow:
            0 25px 50px -12px rgba(0, 0, 0, 0.6),
            0 0 0 1px rgba(148, 163, 184, 0.1);
    }

    :global(html.dark-mode) .prompt-heading {
        color: #f8fafc;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }

    :global(html.dark-mode) .prompt-subheading {
        color: #cbd5e1;
        opacity: 0.9;
        line-height: 1.6;
    }

    :global(html.dark-mode) .project-name-input {
        background: linear-gradient(145deg, #374151 0%, #2d3748 100%);
        border: 2px solid #4b5563;
        color: #f9fafb;
        font-weight: 600;
        box-shadow: 
            inset 0 1px 2px rgba(0, 0, 0, 0.1),
            0 1px 3px rgba(0, 0, 0, 0.1);
    }

    :global(html.dark-mode) .project-name-input::placeholder {
        color: #9ca3af;
        opacity: 0.8;
        font-style: italic;
    }

    :global(html.dark-mode) .project-name-input:focus {
        background: linear-gradient(145deg, #4b5563 0%, #374151 100%);
        border-color: #3b82f6;
        color: #ffffff;
        box-shadow: 
            0 0 0 4px rgba(59, 130, 246, 0.15),
            inset 0 1px 2px rgba(0, 0, 0, 0.1),
            0 4px 12px rgba(59, 130, 246, 0.1);
        outline: none;
        transform: translateY(-1px);
    }

    :global(html.dark-mode) .cancel-button {
        background: linear-gradient(145deg, #475569 0%, #3f4754 100%);
        border: 2px solid #64748b;
        color: #f1f5f9;
        font-weight: 600;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .cancel-button:hover {
        background: linear-gradient(145deg, #64748b 0%, #546070 100%);
        border-color: #94a3b8;
        color: #ffffff;
        transform: translateY(-2px);
        box-shadow: 
            0 4px 12px rgba(0, 0, 0, 0.3),
            0 2px 4px rgba(100, 116, 139, 0.2);
    }

    :global(html.dark-mode) .create-button {
        background: linear-gradient(145deg, #1e40af 0%, #1a365d 100%);
        border: 2px solid #3b82f6;
        color: #ffffff;
        font-weight: 700;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
        box-shadow: 
            0 4px 12px rgba(30, 64, 175, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    :global(html.dark-mode) .create-button:hover:not(:disabled) {
        background: linear-gradient(145deg, #2563eb 0%, #2b77e7 100%);
        border-color: #60a5fa;
        transform: translateY(-2px);
        box-shadow: 
            0 8px 25px rgba(37, 99, 235, 0.4),
            0 4px 12px rgba(0, 0, 0, 0.15);
    }

    :global(html.dark-mode) .create-button:disabled {
        background: linear-gradient(145deg, #6b7280 0%, #5a6173 100%);
        border-color: #6b7280;
        color: #9ca3af;
        cursor: not-allowed;
        transform: none;
        box-shadow: none;
        opacity: 0.7;
    }

    :global(html.dark-mode) .cancel-button:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    }

    :global(html.dark-mode) .create-button:active:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 
            0 4px 12px rgba(30, 64, 175, 0.3),
            0 2px 4px rgba(0, 0, 0, 0.1);
    }

    /* Focus states for accessibility */
    :global(html.dark-mode) .cancel-button:focus,
    :global(html.dark-mode) .create-button:focus {
        outline: 3px solid rgba(59, 130, 246, 0.5);
        outline-offset: 2px;
    }
</style>

<!-- src/lib/components/main/project-hub.svelte -->
<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { projectName } from '$lib/stores/project';
	import ProjectNamePrompt from './project-name-prompt.svelte';

	export let show = false;

	let userName = '';
	let showProjectNamePrompt = false;

	let recentProjects = [
		{ name: 'Lexer', dateModified: 'August 4, 2025' },
		{ name: 'Pascal Parser', dateModified: 'July 28, 2025' },
		{ name: 'x86 Generator', dateModified: 'July 15, 2025' },
		{ name: 'Semantic Rules', dateModified: 'June 30, 2025' },
		{ name: 'AST Viewer', dateModified: 'June 21, 2025' },
		{ name: 'Java-Python', dateModified: 'June 15, 2025' }
	];

	const dispatch = createEventDispatcher();

	function handleClose() {
		dispatch('close');
	}

	function createNewProject() {
		showProjectNamePrompt = true;
	}

	function handleProjectNameConfirm(event: CustomEvent<string>) {
		const newProjectName = event.detail;
		projectName.set(newProjectName);
		console.log(`Project created with name: ${newProjectName}`); // For verification
		showProjectNamePrompt = false;
		handleClose(); // This closes the ProjectHub overlay
	}

	onMount(() => {
		const storedUserId = localStorage.getItem('user_id');
		userName = storedUserId || 'Guest';
	});
</script>

{#if show}
	<div class="backdrop" transition:fly={{ y: -50, duration: 300, opacity: 0.5 }} on:click={handleClose}>
		<div class="modal" on:click|stopPropagation>
			<button class="close-button" on:click={handleClose}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="close-icon"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>

			<h3 class="section-heading">Start a new project</h3>
			<div class="start-project-buttons">
				<button class="project-button" on:click={createNewProject}>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="plus-icon"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
					</svg>
					<span class="button-label">New Blank</span>
				</button>
				<button class="project-button">
					<span class="default-project-text">Default</span>
					<span class="button-label">Project</span>
				</button>
			</div>

			<div class="search-bar-container">
				<div class="search-icon-wrapper">
					<svg
						class="search-icon"
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 20 20"
						fill="currentColor"
					>
						<path
							fill-rule="evenodd"
							d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
							clip-rule="evenodd"
						/>
					</svg>
				</div>
				<input type="text" placeholder="Search projects..." class="search-input" />
			</div>

			<h3 class="section-heading">Recent Projects</h3>
			<div class="project-list-container">
				<div class="project-grid">
					{#each recentProjects as project}
						<div class="project-block">
							<p class="project-name">{project.name}</p>
							<p class="project-date">{project.dateModified}</p>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
{/if}

<ProjectNamePrompt
	bind:show={showProjectNamePrompt}
	on:confirm={handleProjectNameConfirm}
	on:cancel={() => (showProjectNamePrompt = false)}
/>

<style>
    .backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        -webkit-backdrop-filter: blur(4px);
        backdrop-filter: blur(4px);
        z-index: 9999;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 1rem;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
    }

    .modal {
        position: relative; 
        width: 100%;
        max-width: 700px; 
        max-height: fit-content;
        aspect-ratio: 1 / 1;
        background-color: #f1f5f9;
        border-radius: 0.75rem;
        box-shadow:
            0 10px 15px -3px rgba(0, 0, 0, 0.1),
            0 4px 6px -2px rgba(0, 0, 0, 0.05);
        padding: 2rem;
        display: flex;
        flex-direction: column;
    }

    .close-button {
        position: absolute;
        top: 0.75rem;
        right: 0.75rem;
        background: transparent;
        border: none;
        cursor: pointer;
        padding: 0.5rem;
        border-radius: 50%;
        transition: background-color 0.2s ease;
    }

    .close-button:hover {
        background-color: rgba(0, 0, 0, 0.1);
    }

    .close-icon {
        width: 1.5rem;
        height: 1.5rem;
        color: #4b5567;
    }

    .greeting-header {
        font-size: 1.875rem;
        font-weight: 700;
        color: #1f2937;
        text-align: center;
        flex-shrink: 0;
        margin-top: 1rem;
    }

    .section-heading {
        font-size: 1.125rem;
        font-weight: 600;
        color: #374151;
        margin-top: 1.5rem;
        margin-bottom: 0.75rem;
        flex-shrink: 0;
    }

    .start-project-buttons {
        display: flex;
        gap: 1rem;
        flex-shrink: 0;
    }

    .project-button {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        width: 8rem;
        height: 8rem;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        background-color: #fff;
        cursor: pointer;
        transition: background-color 0.2s ease-in-out;
    }

    .project-button:hover {
        background-color: #e5e7eb;
    }

    .plus-icon {
        height: 3rem;
        width: 3rem;
        color: #3b82f6;
    }

    .button-label {
        margin-top: 0.5rem;
        font-size: 0.875rem;
        font-weight: 500;
        color: #4b5567;
    }

    .default-project-text {
        font-size: 1.25rem;
        font-weight: 700;
        color: #374151;
    }

    .search-bar-container {
        position: relative;
        margin-top: 2rem;
        flex-shrink: 0;
    }

    .search-icon-wrapper {
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0.75rem;
        display: flex;
        align-items: center;
        pointer-events: none;
    }

    .search-icon {
        height: 1.25rem;
        width: 1.25rem;
        color: #9ca3af;
    }

    .search-input {
        padding: 0.75rem;
        padding-left: 2.5rem;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        width: 100%;
        box-sizing: border-box;
    }

    .search-input:focus {
        outline: none;
        box-shadow: 0 0 0 2px #3b82f6;
    }

    .project-list-container {
        flex-grow: 1;
        overflow-y: auto;
        min-height: 0;
        margin-right: -1rem;
        padding-right: 1rem;
    }

    .project-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        gap: 1rem;
    }

    .project-block {
        background-color: #fff;
        border: 1px solid #d1d5db;
        border-radius: 0.375rem;
        padding: 1rem;
        cursor: pointer;
        transition:
            transform 0.2s ease,
            box-shadow 0.2s ease;
    }

    .project-block:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }

    .project-name {
        font-size: 1rem;
        font-weight: 600;
        color: #1f2937;
        margin-bottom: 0.25rem;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .project-date {
        font-size: 0.875rem;
        color: #6b7281;
    }

    :global(html.dark) .modal {
        background: #2d3748;
    }

    :global(html.dark) .greeting-header,
    :global(html.dark) .default-project-text,
    :global(html.dark) .project-name {
        color: #edf2f7;
    }

    :global(html.dark) .section-heading,
    :global(html.dark) .button-label,
    :global(html.dark) .project-date {
        color: #a0aec0;
    }

    :global(html.dark) .project-button {
        background-color: #4a5568;
        border-color: #718096;
    }

    :global(html.dark) .project-button:hover {
        background-color: #2d3748;
    }

    :global(html.dark) .search-input {
        background-color: #4a5568;
        border-color: #718096;
        color: #edf2f7;
    }

    :global(html.dark) .project-block {
        background-color: #4a5568;
        border-color: #718096;
    }

    :global(html.dark) .close-icon {
        color: #a0aec0;
	}
</style>

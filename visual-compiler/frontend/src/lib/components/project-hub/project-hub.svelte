<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import { fly } from 'svelte/transition';
	import { projectName } from '$lib/stores/project';
	import { pipelineStore, resetPipeline, type Pipeline } from '$lib/stores/pipeline';
	import { confirmedSourceCode, resetSourceCode } from '$lib/stores/source-code';
	import ProjectNamePrompt from './project-name-prompt.svelte';
	import DeleteConfirmPrompt from './delete-confirmation.svelte'; 
	import { AddToast } from '$lib/stores/toast';
	import { 
        updateLexerStateFromProject,
        resetLexerState
    } from '$lib/stores/lexer';
    import { 
        updateParserStateFromProject, 
        resetParserState 
    } from '$lib/stores/parser'; 
    import { 
        updateAnalyserStateFromProject, 
        resetAnalyserState 
    } from '$lib/stores/analyser';
    import { 
        updateTranslatorStateFromProject, 
        resetTranslatorState 
    } from '$lib/stores/translator';
    
    import { phase_completion_status } from '$lib/stores/pipeline';
	import { triggerTutorialForNewProject } from '$lib/stores/tutorial';

	const dispatch = createEventDispatcher();

	export let show = false;

	let userName = '';
	let showProjectNamePrompt = false;
	let showDeleteConfirmPrompt = false; // State for the delete confirmation
	let projectToDelete = ''; // State to hold the name of the project to be deleted
	let hasExistingProject = false; // Track if a project is already loaded
	let currentProjectName = ''; // Track the current project name
	let searchQuery = ''; // Search functionality

	interface Project {
		name: string;
		dateModified: string;
	}

	let recentProjects: Project[] = [];

	// Reactive statement for filtered projects based on search
	$: filteredProjects = recentProjects.filter(project => 
		project.name.toLowerCase().includes(searchQuery.toLowerCase())
	);

	// Clear search when modal opens
	$: if (show) {
		searchQuery = '';
		fetchProjects();
	}

	async function fetchProjects() {
		const userId = sessionStorage.getItem('user_id');
		if (!userId) return;

		try {
			const response = await fetch(`https://www.visual-compiler.co.za/api/users/getProjects?users_id=${userId}`, {
				method: 'GET',
				headers: {
					'accept': 'application/json'
				}
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();
			
			if (data.all_projects) {
				// Transform projects into the required format
				recentProjects = data.all_projects.map((projectName: string) => ({
					name: projectName,
					dateModified: new Date().toLocaleDateString('en-US', {
						year: 'numeric',
						month: 'long',
						day: 'numeric'
					})
				}));
			}
		} catch (error) {
			console.error('Error fetching projects:', error);
			recentProjects = []; // Reset to empty array on error
		}
	}

	function handleClose() {
		if (hasExistingProject) {
			dispatch('close');
		}
	}

	function createNewProject() {
		showProjectNamePrompt = true;
	}

	function clearSearch() {
		searchQuery = '';
	}

	function connectNode(pipeline: Pipeline) {

		const node_connections = pipeline.connections;

		node_connections.forEach(conn => {

			const start_node = document.getElementById(conn.sourceAnchor);
			const end_node = document.getElementById(conn.targetAnchor);

			start_node?.dispatchEvent(new MouseEvent("mousedown", { bubbles: true, button: 0 }));
			end_node?.dispatchEvent(new MouseEvent("mouseup", { bubbles: true, button: 0 }));
		});

	}

	async function selectProject(selectedProjectName: string) {
        const userId = sessionStorage.getItem('user_id');
        if (!userId) {
            AddToast('Please log in to select a project', 'error');
            return;
        }

        try {
            console.log('=== LOADING PROJECT:', selectedProjectName, '===');
            
            const response = await fetch(
                `https://www.visual-compiler.co.za/api/users/getProject?project_name=${selectedProjectName}&users_id=${userId}`, {
                method: 'GET',
                headers: {
                    'accept': 'application/json'
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            if (data.message === "Retrieved users project details") {
                if (data.results) {
                    // FIX: Set project name FIRST, then load data
                    projectName.set(selectedProjectName);
                    
                    // Small delay to ensure project name change is processed
                    await new Promise(resolve => setTimeout(resolve, 50));
                    
                    // FIX: Handle source code FIRST before loading other phase states
                    if (data.results.lexing?.code) {
                        console.log('Loading source code:', data.results.lexing.code);
                        confirmedSourceCode.set(data.results.lexing.code);
                    }
                    
                    // Load all phase states from project data
                    console.log('Loading project data into stores...');
                    updateLexerStateFromProject(data.results);
                    updateParserStateFromProject(data.results);
                    updateAnalyserStateFromProject(data.results);
                    updateTranslatorStateFromProject(data.results);

                    // Update phase completion status based on actual artifacts
                    const hasTokens = !!(data.results.lexing?.tokens && data.results.lexing.tokens.length > 0);
                    const hasCode = !!data.results.lexing?.code;
                    const hasTree = !!(data.results.parsing?.tree);
                    const hasSymbolTable = !!(data.results.analysing?.symbol_table_artefact?.symbolscopes);
                    const hasTranslation = !!(data.results.translating?.code && data.results.translating.code.length > 0);
                    
                    phase_completion_status.set({
                        source: hasCode,
                        lexer: hasTokens,
                        parser: hasTree,
                        analyser: hasSymbolTable,
                        translator: hasTranslation
                    });

                    // Handle pipeline data
                    if (data.results && data.results.pipeline) {
                        pipelineStore.set(data.results.pipeline);
                        console.log('Restored pipeline:', data.results.pipeline);
                        
                        await tick();
                        connectNode(data.results.pipeline); 
                    } else {
                        pipelineStore.set({
                            nodes: [],
                            connections: [],
                            lastSaved: null
                        });
                        console.log('No saved pipeline found, initialized empty state');
                    }

                    console.log('=== PROJECT LOADED SUCCESSFULLY ===');
                    handleClose();
                    AddToast('Project loaded successfully', 'success');
                }
            } else {
                console.error('Failed to retrieve project details');
                AddToast('Failed to retrieve project details', 'error');
            }
        } catch (error: any) {
            console.error('Error loading project:', error);
            const errorMessage = error.message || 'Unknown error occurred';
            AddToast(`Error loading project: ${errorMessage}`, 'error');
        }
    }

	async function handleProjectNameConfirm(event: CustomEvent<string>) {
		const newProjectName = event.detail;
		const userId = sessionStorage.getItem('user_id');
		if (!userId) return;

		try {
            // IMPORTANT: Reset project name FIRST to trigger component resets
            projectName.set(''); // Clear project name first
            
            // Then reset all stores
            console.log('=== RESETTING ALL STORES FOR NEW PROJECT ===');
            resetLexerState();
            resetParserState();
            resetAnalyserState();
            resetTranslatorState();
            resetSourceCode();
            resetPipeline();
            
            phase_completion_status.set({
                source: false,
                lexer: false,
                parser: false,
                analyser: false,
                translator: false
            });

            const response = await fetch('https://www.visual-compiler.co.za/api/users/save', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    project_name: newProjectName,
                    users_id: userId
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Only AFTER everything is reset and API call succeeds, set the new project name
            projectName.set(newProjectName);
            
            showProjectNamePrompt = false;
            await fetchProjects();
            
            setTimeout(() => {
                triggerTutorialForNewProject();
            }, 500);
            
            handleClose();
            console.log('=== NEW PROJECT CREATED WITH CLEAN STATE ===');
        } catch (error) {
            console.error('Error saving project:', error);
        }
	}

	/**
	 * @function handleDeleteClick
	 * @description Opens the confirmation prompt before deleting a project.
	 * @param {string} projName - The name of the project to delete.
	 * @param {MouseEvent} event - The mouse event, used to stop propagation.
	 */
	function handleDeleteClick(projName: string, event: MouseEvent) {
		event.stopPropagation(); // Prevents the main project block from being clicked.
		projectToDelete = projName;
		showDeleteConfirmPrompt = true;
	}

	/**
	 * @function confirmDelete
	 * @description Proceeds with the deletion after user confirmation.
	 */
	async function confirmDelete() {
		console.log(`Confirmed deletion of project: ${projectToDelete}`);
		const userId = sessionStorage.getItem('user_id');
		if (!userId) return;

		try {
			const response = await fetch('https://www.visual-compiler.co.za/api/users/deleteProject', {
				method: 'DELETE',
				headers: {
					'accept': 'application/json',
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					project_name: projectToDelete,
					users_id: userId
				})
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();
			console.log('Project deleted successfully:', data);

			// Remove from UI list after successful deletion
			recentProjects = recentProjects.filter((p) => p.name !== projectToDelete);
			
			// If the deleted project is the currently loaded project, clear it
			if (projectToDelete === currentProjectName) {
				projectName.set('');
				currentProjectName = '';
				hasExistingProject = false;
			}
			
			showDeleteConfirmPrompt = false;
			projectToDelete = '';
		} catch (error) {
			console.error('Error deleting project:', error);
			// Still close the modal even if delete fails
			showDeleteConfirmPrompt = false;
			projectToDelete = '';
		}
	}

	onMount(() => {
		const storedUserId = sessionStorage.getItem('user_id');
		userName = storedUserId || 'Guest';
		
		// Check if there's already a project loaded
		const unsubscribe = projectName.subscribe(value => {
			hasExistingProject = value.trim() !== '';
			currentProjectName = value;
		});
		
		fetchProjects(); // Fetch projects when component mounts
		
		// Clean up subscription when component is destroyed
		return () => {
			unsubscribe();
		};
	});
</script>

{#if show}
	<div class="backdrop" transition:fly={{ y: -50, duration: 300, opacity: 0.5 }} on:click={hasExistingProject ? handleClose : undefined}>
		<div class="modal" on:click|stopPropagation>
			{#if hasExistingProject}
			<button class="close-button" on:click={handleClose} aria-label="Close project hub">
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
			{/if}

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
				<input 
					type="text" 
					placeholder="Search projects..." 
					class="search-input" 
					bind:value={searchQuery}
				/>
				{#if searchQuery}
					<button class="clear-search-btn" on:click={clearSearch} title="Clear search">
						<svg viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
						</svg>
					</button>
				{/if}
			</div>

			<h3 class="section-heading">
				{searchQuery ? `Search Results (${filteredProjects.length})` : 'Recent Projects'}
			</h3>
			<div class="project-list-container">
				{#if filteredProjects.length > 0}
					<div class="project-grid">
						{#each filteredProjects as project}
							<div class="project-block" on:click={() => selectProject(project.name)}>
								<button
									class="delete-button"
									on:click={(event) => handleDeleteClick(project.name, event)}
									aria-label="Delete project"
									title="Delete project"
								>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="16"
									height="16"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
								>
									<polyline points="3 6 5 6 21 6" />
									<path
										d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
									/>
								</svg>
							</button>
							<p class="project-name">{project.name}</p>
							<p class="project-date">{project.dateModified}</p>
						</div>
					{/each}
				</div>
			{:else}
				<div class="no-results">
					{#if searchQuery}
						<div class="no-results-content">
							<svg class="no-results-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<circle cx="11" cy="11" r="8"></circle>
								<path d="m21 21-4.35-4.35"></path>
								<path d="M8 11h6"></path>
							</svg>
							<h4>No projects found</h4>
							<p>No projects match "{searchQuery}". Try a different search term.</p>
						</div>
					{:else}
						<div class="no-results-content">
							<svg class="no-results-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
								<line x1="16" y1="2" x2="16" y2="6"></line>
								<line x1="8" y1="2" x2="8" y2="6"></line>
								<line x1="3" y1="10" x2="21" y2="10"></line>
							</svg>
							<h4>No projects yet</h4>
							<p>You haven't created any projects yet. Create your first project to get started!</p>
						</div>
					{/if}
				</div>
			{/if}
			</div>
		</div>
	</div>
{/if}

<ProjectNamePrompt
	bind:show={showProjectNamePrompt}
	on:confirm={handleProjectNameConfirm}
	on:cancel={() => (showProjectNamePrompt = false)}
	recentProjects={recentProjects}
/>

<DeleteConfirmPrompt
	bind:show={showDeleteConfirmPrompt}
	projectName={projectToDelete}
	on:confirm={confirmDelete}
	on:cancel={() => (showDeleteConfirmPrompt = false)}
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
        background: #ffffff;
        border-radius: 0.75rem;
        box-shadow:
            0 20px 40px -12px rgba(0, 0, 0, 0.2),
            0 8px 20px -5px rgba(0, 0, 0, 0.08),
            0 0 0 1px rgba(203, 213, 225, 0.2);
        border: 1px solid rgba(203, 213, 225, 0.4);
        padding: 2rem;
        display: flex;
        flex-direction: column;
    }

    .close-button {
        position: absolute;
        top: 0.75rem;
        right: 0.75rem;
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        cursor: pointer;
        padding: 0.5rem;
        border-radius: 8px;
        transition: all 0.2s ease;
    }

    .close-button:hover {
        background: #fee2e2;
        border-color: #fecaca;
        transform: scale(1.02);
    }

    .close-icon {
        width: 1.5rem;
        height: 1.5rem;
        color: #64748b;
        transition: color 0.2s ease;
    }

    .close-button:hover .close-icon {
        color: #ef4444;
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
        font-weight: 700;
        color: #374151;
        margin-top: 1.5rem;
        margin-bottom: 0.75rem;
        flex-shrink: 0;
        border-bottom: 1px solid rgba(203, 213, 225, 0.4);
        padding-bottom: 0.5rem;
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
        border: 2px solid #d1d5db;
        border-radius: 0.5rem;
        background: linear-gradient(145deg, #001a6e 0%, #001556 100%);
        cursor: pointer;
        transition: all 0.2s ease;
        box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
    }

    .project-button:hover {
        background: linear-gradient(145deg, #002080 0%, #001a6e 100%);
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 26, 110, 0.2);
    }

    .project-button:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
    }

    .plus-icon {
        height: 3rem;
        width: 3rem;
        color: #ffffff;
        filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
    }

    .button-label {
        margin-top: 0.5rem;
        font-size: 0.875rem;
        font-weight: 600;
        color: #ffffff;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
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
        padding: 0.875rem 1rem;
        padding-left: 2.75rem;
        padding-right: 2.75rem;
        border: 2px solid #e2e8f0;
        border-radius: 0.5rem;
        width: 100%;
        box-sizing: border-box;
        background: #ffffff;
        color: #1f2937;
        font-weight: 500;
        transition: all 0.2s ease;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .search-input::placeholder {
        color: #9ca3af;
        opacity: 0.7;
        font-style: italic;
        font-weight: 400;
    }

    .search-input:focus {
        outline: none;
        border-color: #6366f1;
        background: #fefefe;
        color: #111827;
        box-shadow: 0 2px 4px rgba(99, 102, 241, 0.1);
    }

    .clear-search-btn {
        position: absolute;
        right: 8px;
        top: 50%;
        transform: translateY(-50%);
        background: #f8fafc;
        border: 1px solid #e2e8f0;
        padding: 6px;
        cursor: pointer;
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s ease;
    }

    .clear-search-btn:hover {
        background: #fee2e2;
        border-color: #fecaca;
        transform: translateY(-50%) scale(1.05);
    }

    .clear-search-btn svg {
        width: 16px;
        height: 16px;
        color: #6b7280;
        transition: color 0.2s ease;
    }

    .clear-search-btn:hover svg {
        color: #ef4444;
    }

    .project-list-container {
        flex-grow: 1;
        overflow-y: auto;
        min-height: 0;
        margin-right: -1rem;
        padding-right: 1rem;
        scrollbar-width: thin;
        scrollbar-color: #cbd5e1 #f1f5f9;
    }

    .project-list-container::-webkit-scrollbar {
        width: 8px;
    }

    .project-list-container::-webkit-scrollbar-track {
        background: #f1f5f9;
        border-radius: 4px;
    }

    .project-list-container::-webkit-scrollbar-thumb {
        background: #cbd5e1;
        border-radius: 4px;
        border: 1px solid #f1f5f9;
    }

    .project-list-container::-webkit-scrollbar-thumb:hover {
        background: #94a3b8;
    }

    .project-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
        gap: 1rem;
    }

    .project-block {
        position: relative;
        background: #ffffff;
        border: 2px solid #e2e8f0;
        border-radius: 0.5rem;
        padding: 1rem;
        cursor: pointer;
        transition: all 0.2s ease;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        overflow: hidden;
        margin-top: 0.4rem;
    }

    .project-block::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 2px;
        background: linear-gradient(90deg, 
            transparent 0%, 
            rgba(99, 102, 241, 0.3) 50%, 
            transparent 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    .project-block:hover {
        background: #fefefe;
        border-color: #cbd5e1;
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(4, 26, 71, 0.1);
    }

    .project-block:hover::before {
        opacity: 1;
    }

    .delete-button {
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        background: #f8fafc;
        border: 1px solid rgba(239, 68, 68, 0.2);
        cursor: pointer;
        padding: 0.375rem;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #9ca3af;
        opacity: 0;
        transition: all 0.2s ease;
        backdrop-filter: blur(4px);
    }

    .project-block:hover .delete-button {
        opacity: 1;
        background: #fefefe;
        border-color: rgba(239, 68, 68, 0.3);
    }

    .delete-button:hover {
        background: #fee2e2;
        border-color: #fecaca;
        color: #ef4444;
        transform: scale(1.05);
        box-shadow: 0 2px 4px rgba(239, 68, 68, 0.2);
    }

    .project-name {
        font-size: 1rem;
        font-weight: 700;
        color: #1f2937;
        margin-bottom: 0.25rem;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        transition: color 0.2s ease;
     
    }

    .project-block:hover .project-name {
        color: #111827;
    }

    .project-date {
        font-size: 0.875rem;
        color: #64748b;
        font-weight: 500;
        opacity: 0.8;
        transition: all 0.2s ease;
    }

    .project-block:hover .project-date {
        color: #475569;
        opacity: 1;
    }

    /* No Results Styles */
    .no-results {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 200px;
        padding: 2rem;
        background: #f8fafc;
        border-radius: 1rem;
        border: 1px solid rgba(203, 213, 225, 0.4);
    }

    .no-results-content {
        text-align: center;
        max-width: 300px;
    }

    .no-results-icon {
        width: 3rem;
        height: 3rem;
        color: #9ca3af;
        margin: 0 auto 1rem;
        display: block;
    }

    .no-results-content h4 {
        font-size: 1.125rem;
        font-weight: 700;
        color: #374151;
        margin: 0 0 0.5rem 0;
    }

    .no-results-content p {
        font-size: 0.95rem;
        color: #6b7280;
        line-height: 1.6;
        margin: 0;
        opacity: 0.9;
    }

    /* Focus states for accessibility */
    .project-button:focus,
    .project-block:focus {
        outline: 3px solid rgba(99, 102, 241, 0.3);
        outline-offset: 2px;
    }

    .delete-button:focus,
    .close-button:focus,
    .clear-search-btn:focus {
        outline: 2px solid rgba(99, 102, 241, 0.4);
        outline-offset: 2px;
    }

    /* Mobile responsiveness */
    @media (max-width: 640px) {
        .modal {
            margin: 1rem;
            padding: 1.5rem;
        }
        
        .project-button {
            width: 6rem;
            height: 6rem;
        }
        
        .plus-icon {
            height: 2rem;
            width: 2rem;
        }
    }

    /* Enhanced Dark Mode Styles */
    :global(html.dark-mode) .backdrop {
        background: rgba(0, 0, 0, 0.8);
        backdrop-filter: blur(8px);
        -webkit-backdrop-filter: blur(8px);
    }

    :global(html.dark-mode) .modal {
        background: #1e293b;
        border: 1px solid #334155;
        color: #f1f5f9;
        box-shadow:
            0 20px 40px -12px rgba(0, 0, 0, 0.6),
            0 0 0 1px rgba(148, 163, 184, 0.1);
    }

    :global(html.dark-mode) .close-button {
        background: #334155;
        border: 1px solid #475569;
        transition: all 0.2s ease;
    }

    :global(html.dark-mode) .close-button:hover {
        background: rgba(239, 68, 68, 0.1);
        border-color: rgba(239, 68, 68, 0.2);
        transform: scale(1.05);
    }

    :global(html.dark-mode) .close-icon {
        color: #cbd5e1;
    }

    :global(html.dark-mode) .close-button:hover .close-icon {
        color: #fca5a5;
    }

    :global(html.dark-mode) .greeting-header {
        color: #f8fafc;
    }

    :global(html.dark-mode) .section-heading {
        color: #e2e8f0;
        border-bottom: 1px solid rgba(148, 163, 184, 0.3);
    }

    :global(html.dark-mode) .project-button {
        background: linear-gradient(145deg, #4338ca 0%, #3730a3 100%);
        border: 2px solid #6366f1;
        color: #ffffff;
        box-shadow: 0 2px 4px rgba(67, 56, 202, 0.1);
    }

    :global(html.dark-mode) .project-button:hover {
        background: linear-gradient(145deg, #6366f1 0%, #4f46e5 100%);
        box-shadow: 0 4px 8px rgba(99, 102, 241, 0.2);
        transform: translateY(-2px);
    }

    :global(html.dark-mode) .plus-icon {
        color: #ffffff;
    }

    :global(html.dark-mode) .button-label {
        color: #f1f5f9;
    }

    :global(html.dark-mode) .default-project-text {
        color: #f8fafc;
    }

    :global(html.dark-mode) .search-input {
        background: #374151;
        border: 2px solid #4b5563;
        color: #f9fafb;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    :global(html.dark-mode) .search-input::placeholder {
        color: #9ca3af;
        opacity: 0.8;
    }

    :global(html.dark-mode) .search-input:focus {
        background: #4b5563;
        border-color: #6366f1;
        color: #ffffff;
        box-shadow: 0 2px 4px rgba(99, 102, 241, 0.1);
    }

    :global(html.dark-mode) .search-icon {
        color: #e5e7eb;
    }

    :global(html.dark-mode) .clear-search-btn {
        background: #334155;
        border: 1px solid #475569;
    }

    :global(html.dark-mode) .clear-search-btn:hover {
        background: rgba(239, 68, 68, 0.1);
        border-color: rgba(239, 68, 68, 0.2);
        transform: translateY(-50%) scale(1.05);
    }

    :global(html.dark-mode) .clear-search-btn svg {
        color: #d1d5db;
    }

    :global(html.dark-mode) .clear-search-btn:hover svg {
        color: #fca5a5;
    }

    :global(html.dark-mode) .project-block {
        background: #374151;
        border: 2px solid #4b5563;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    :global(html.dark-mode) .project-block::before {
        background: linear-gradient(90deg, 
            transparent 0%, 
            rgba(99, 102, 241, 0.3) 50%, 
            transparent 100%);
    }

    :global(html.dark-mode) .project-block:hover {
        background: #4b5563;
        border-color: #6b7280;
        box-shadow: 0 4px 8px rgba(4, 26, 71, 0.15);
        transform: translateY(-2px);
    }

    :global(html.dark-mode) .project-name {
        color: #f8fafc;
    }

    :global(html.dark-mode) .project-block:hover .project-name {
        color: #ffffff;
    }

    :global(html.dark-mode) .project-date {
        color: #cbd5e1;
    }

    :global(html.dark-mode) .project-block:hover .project-date {
        color: #e2e8f0;
    }

    :global(html.dark-mode) .delete-button {
        background: #334155;
        border: 1px solid rgba(239, 68, 68, 0.2);
        color: #cbd5e1;
    }

    :global(html.dark-mode) .project-block:hover .delete-button {
        background: #475569;
        border-color: rgba(239, 68, 68, 0.3);
    }

    :global(html.dark-mode) .delete-button:hover {
        background: rgba(239, 68, 68, 0.1);
        border-color: rgba(239, 68, 68, 0.4);
        color: #fca5a5;
        transform: scale(1.05);
        box-shadow: 0 2px 4px rgba(239, 68, 68, 0.2);
    }

    :global(html.dark-mode) .no-results {
        background: #334155;
        border: 1px solid #475569;
    }

    :global(html.dark-mode) .no-results-icon {
        color: #64748b;
    }

    :global(html.dark-mode) .no-results-content h4 {
        color: #f1f5f9;
    }

    :global(html.dark-mode) .no-results-content p {
        color: #cbd5e1;
    }

    :global(html.dark-mode) .project-list-container {
        scrollbar-width: thin;
        scrollbar-color: #475569 #1e293b;
    }

    :global(html.dark-mode) .project-list-container::-webkit-scrollbar-track {
        background: #1e293b;
    }

    :global(html.dark-mode) .project-list-container::-webkit-scrollbar-thumb {
        background: #475569;
        border: 1px solid #1e293b;
    }

    :global(html.dark-mode) .project-list-container::-webkit-scrollbar-thumb:hover {
        background: #64748b;
    }
</style>

<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import { fly } from 'svelte/transition';
	import { projectName } from '$lib/stores/project';
	import { pipelineStore, resetPipeline, type Pipeline } from '$lib/stores/pipeline';
	import { confirmedSourceCode, resetSourceCode } from '$lib/stores/source-code';
	import ProjectNamePrompt from './project-name-prompt.svelte';
	import DeleteConfirmPrompt from './delete-confirmation.svelte'; 
	import { AddToast } from '$lib/stores/toast';
	import { updateLexerStateFromProject, resetLexerState } from '$lib/stores/lexer';
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
		const userId = localStorage.getItem('user_id');
		if (!userId) return;

		try {
			const response = await fetch(`http://localhost:8080/api/users/getProjects?users_id=${userId}`, {
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
		const userId = localStorage.getItem('user_id');
		if (!userId) {
			AddToast('Please log in to select a project', 'error');
			return;
		}

		// Clear existing state before loading new project
		resetLexerState();
		resetSourceCode();
		
		try {
			const response = await fetch(
				`http://localhost:8080/api/users/getProject?project_name=${selectedProjectName}&users_id=${userId}`, {
				method: 'GET',
				headers: {
					'accept': 'application/json'
				}
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();
			
			// Log detailed project data to console
			console.log('Project Data:', {
				message: data.message,
				projectName: selectedProjectName,
				hasResults: !!data.results,
				lexingData: data.results?.lexing,
				hasSourceCode: !!data.results?.lexing?.code,
				pipelineData: data.results?.pipeline,
				fullResponse: data
			});

			if (data.message === "Retrieved users project details") {
				if (data.results) {
					updateLexerStateFromProject({
						lexing: data.results.lexing,
						parsing: data.results.parsing,
						analysing: data.results.analysing,
						translating: {
							code: data.results.translating?.code || [],
							translating_rules: data.results.translating?.translating_rules?.map(rule => ({
								sequence: Array.isArray(rule.sequence) ? rule.sequence : [],
								translation: Array.isArray(rule.translation) ? rule.translation : []
							})) || []  
						}
					});

					// Update phase completion status
					const hasTokens = !!(data.results.lexing?.tokens && data.results.lexing.tokens.length > 0);
					const hasCode = !!data.results.lexing?.code;
					const hasTree = !!(data.results.parsing?.tree?.root);
					const hasSymbolTable = !!(data.results.analysing?.symbol_table_artefact?.symbolscopes);
					const hasTranslation = !!(data.results.translating?.code && data.results.translating.code.length > 0);
					
					phase_completion_status.set({
						source: hasCode,
						lexer: hasTokens,
						parser: hasTree,
						analyser: hasSymbolTable,
						translator: hasTranslation
					});

					// Handle source code if it exists
					if (hasCode) {
						confirmedSourceCode.set(data.results.lexing.code);
					}

					// Check if there's pipeline data
					if (data.results && data.results.pipeline) {
						// Update the pipeline store with the saved pipeline data
						pipelineStore.set(data.results.pipeline);
						console.log('Restored pipeline:', data.results.pipeline);
						AddToast('Pipeline restored successfully', 'success');

						await tick();
						connectNode(data.results.pipeline); 
					} else {
						// If no pipeline data, initialize with empty state
						pipelineStore.set({
							nodes: [],
							connections: [],
							lastSaved: null
						});
						console.log('No saved pipeline found, initialized empty state');
					}

					// Store the selected project name in the store
					projectName.set(selectedProjectName);
					
					// Close the project hub modal
					handleClose();
				}
			} else {
				console.error('Failed to retrieve project details');
				AddToast('Failed to retrieve project details', 'error');
			}
		} catch (error: any) {
			console.error('Error verifying project selection:', error);
			const errorMessage = error.message || 'Unknown error occurred';
			AddToast(`Error loading project: ${errorMessage}`, 'error');
		}
	}

	async function handleProjectNameConfirm(event: CustomEvent<string>) {
		const newProjectName = event.detail;
		const userId = localStorage.getItem('user_id');
		if (!userId) return;

		try {
			const response = await fetch('http://localhost:8080/api/users/save', {
				method: 'POST',
				headers: {
					'accept': 'application/json',
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

			const data = await response.json();

			// Update UI and close modal
			projectName.set(newProjectName);
			resetPipeline(); // Clear the canvas for the new blank project
			resetLexerState(); // Clear all saved input data for a fresh start
			resetSourceCode(); // Clear the source code for a fresh start
			
			// Reset phase completion status for new project
			phase_completion_status.set({
				source: false,
				lexer: false,
				parser: false,
				analyser: false,
				translator: false
			});
			
			showProjectNamePrompt = false;
			await fetchProjects(); // Refresh the project list
			
			// Trigger tutorial for new projects
			setTimeout(() => {
				triggerTutorialForNewProject();
			}, 500); // Small delay to ensure the workspace is loaded
			
			handleClose();
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
		const userId = localStorage.getItem('user_id');
		if (!userId) return;

		try {
			const response = await fetch('http://localhost:8080/api/users/deleteProject', {
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
		const storedUserId = localStorage.getItem('user_id');
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
		padding-right: 2.5rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		width: 100%;
		box-sizing: border-box;
	}

	.search-input:focus {
		outline: none;
		box-shadow: 0 0 0 2px #3b82f6;
	}

	.clear-search-btn {
		position: absolute;
		right: 8px;
		top: 50%;
		transform: translateY(-50%);
		background: none;
		border: none;
		padding: 4px;
		cursor: pointer;
		border-radius: 4px;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: background-color 0.2s ease;
	}

	.clear-search-btn:hover {
		background-color: #f3f4f6;
	}

	.clear-search-btn svg {
		width: 16px;
		height: 16px;
		color: #6b7280;
	}

	.clear-search-btn:hover svg {
		color: #374151;
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
		position: relative; /* Needed for absolute positioning of the delete button */
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

	.delete-button {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		background: transparent;
		border: none;
		cursor: pointer;
		padding: 0.25rem;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #9ca3af;
		opacity: 0; /* Hidden by default */
		transition:
			opacity 0.2s ease,
			background-color 0.2s ease;
	}

	.project-block:hover .delete-button {
		opacity: 1; /* Show on hover */
	}

	.delete-button:hover {
		background-color: #e5e7eb;
		color: #ef4444;
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

	/* No Results Styles */
	.no-results {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 200px;
		padding: 2rem;
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
		font-weight: 600;
		color: #374151;
		margin: 0 0 0.5rem 0;
	}

	.no-results-content p {
		font-size: 0.95rem;
		color: #6b7280;
		line-height: 1.5;
		margin: 0;
	}

	/* Dark Mode Styles */
	:global(html.dark-mode) .backdrop {
		background: rgba(0, 0, 0, 0.7);
	}

	:global(html.dark-mode) .modal {
		background-color: #1a2a4a;
		color: #f0f0f0;
		box-shadow:
			0 10px 15px -3px rgba(0, 0, 0, 0.3),
			0 4px 6px -2px rgba(0, 0, 0, 0.2);
	}

	:global(html.dark-mode) .close-button:hover {
		background-color: rgba(74, 85, 104, 0.2);
	}

	:global(html.dark-mode) .close-icon {
		color: #a0aec0;
	}

	:global(html.dark-mode) .greeting-header,
	:global(html.dark-mode) .default-project-text,
	:global(html.dark-mode) .project-name {
		color: #ebeef1;
	}

	:global(html.dark-mode) .section-heading,
	:global(html.dark-mode) .button-label,
	:global(html.dark-mode) .project-date {
		color: #a0aec0;
	}

	:global(html.dark-mode) .project-button {
		background-color: #041a47;
		border: 2px solid #041a47;
		color: #ffffff;
		box-shadow: 0 2px 8px rgba(0, 26, 110, 0.3);
	}

	:global(html.dark-mode) .project-button:hover {
		background-color: #002a8e;
		border-color: #	;
		box-shadow: 0 4px 12px rgba(0, 26, 110, 0.4);
		transform: translateY(-1px);
	}

	:global(html.dark-mode) .plus-icon {
		color: #ffffff;
	}

	:global(html.dark-mode) .search-input {
		background-color: #374151;
		border: 2px solid #4b5563;
		color: #ffffff;
		font-weight: 500;
	}

	:global(html.dark-mode) .search-input::placeholder {
		color: #9ca3af;
		opacity: 1;
	}

	:global(html.dark-mode) .search-input:focus {
		border-color: #001A6E;
		background-color: #4b5563;
		box-shadow: 0 0 0 3px rgba(0, 26, 110, 0.2);
	}

	:global(html.dark-mode) .search-icon {
		color: #d1d5db;
	}

	:global(html.dark-mode) .clear-search-btn {
		color: #d1d5db;
	}

	:global(html.dark-mode) .clear-search-btn:hover {
		background-color: #6b7280;
		color: #ffffff;
	}

	:global(html.dark-mode) .project-block {
		background-color: #374151;
		border: 2px solid #4b5563;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
	}

	:global(html.dark-mode) .project-block:hover {
		background-color: #4b5563;
		border-color: #6b7280;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		transform: translateY(-1px);
	}

	:global(html.dark-mode) .delete-button {
		color: #9ca3af;
	}

	:global(html.dark-mode) .delete-button:hover {
		background-color: #4a5568;
		color: #f87171;
	}

	:global(html.dark-mode) .no-results-content h4 {
		color: #ebeef1;
	}

	:global(html.dark-mode) .no-results-content p {
		color: #a0aec0;
	}

	:global(html.dark-mode) .no-results-icon {
		color: #718096;
	}
</style>

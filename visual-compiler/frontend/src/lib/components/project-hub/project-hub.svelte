<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import { fly } from 'svelte/transition';
	import { projectName } from '$lib/stores/project';
	import { pipelineStore, resetPipeline, type Pipeline } from '$lib/stores/pipeline';
	import { confirmedSourceCode } from '$lib/stores/source-code';
	import ProjectNamePrompt from './project-name-prompt.svelte';
	import DeleteConfirmPrompt from './delete-confirmation.svelte'; 
	import { AddToast } from '$lib/stores/toast';
	import { updateLexerStateFromProject } from '$lib/stores/lexer';
	import { phase_completion_status } from '$lib/stores/pipeline';

	const dispatch = createEventDispatcher();

	export let show = false;

	// Watch for changes to the 'show' prop and refresh projects when modal opens
	$: if (show) {
		fetchProjects();
	}

	let userName = '';
	let showProjectNamePrompt = false;
	let showDeleteConfirmPrompt = false; // State for the delete confirmation
	let projectToDelete = ''; // State to hold the name of the project to be deleted
	let hasExistingProject = false; // Track if a project is already loaded
	let currentProjectName = ''; // Track the current project name

	interface Project {
		name: string;
		dateModified: string;
	}

	let recentProjects: Project[] = [];
	
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
				if (data.results?.lexing) {
					// Update lexer state
					updateLexerStateFromProject(data.results);

					// Update phase completion status directly using the store
					const hasTokens = !!(data.results.lexing.tokens && data.results.lexing.tokens.length > 0);
					const hasCode = !!data.results.lexing.code;
					
					phase_completion_status.set({
						source: hasCode,
						lexer: hasTokens,
						parser: !!(data.results.lexing.tree && !data.results.lexing.parsing_error),
						analyser: !!(data.results.lexing.symbol_table && !data.results.lexing.analyser_error),
						translator: !!(data.results.lexing.translated_code && data.results.lexing.translated_code.length > 0)
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
			showProjectNamePrompt = false;
			await fetchProjects(); // Refresh the project list
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

	:global(html.dark) .delete-button {
		color: #9ca3af;
	}

	:global(html.dark) .delete-button:hover {
		background-color: #4a5568;
		color: #f87171;
	}

	:global(html.dark) .close-icon {
		color: #a0aec0;
	}
</style>

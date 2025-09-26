<script lang="ts">
	import { AddToast } from '$lib/stores/toast';
	import { confirmedSourceCode } from '$lib/stores/source-code';
	import { tick } from 'svelte';
	import { onMount, onDestroy } from 'svelte';
	import { projectName } from '$lib/stores/project';
	import { get } from 'svelte/store';  
	import { activePhase, setActivePhase } from '$lib/stores/pipeline'; 

	let code_text = '';
	
	let previous_code_text = '';
	let isDefaultInput = false;
	let isConfirmed = false;
	let textareaEl: HTMLTextAreaElement;
	export let onCodeSubmitted: (code: string) => void = () => {};
	
	// Add flag to control window closing behavior
	let shouldCloseWindow = true;

	// Guest user detection
	let isGuestUser = false;

	// --- REAL PROJECTS DATA ---
	let projects: Array<{ name: string; code: string }> = [
		{ name: 'Select a project...', code: '' }
	];
	let selectedProject = projects[0];

	// Fetch projects from backend
	async function fetchProjects() {
		// Skip fetching projects for guest users
		if (isGuestUser) {
			return;
		}

		const userId = localStorage.getItem('user_id');
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
			
			if (data.all_projects && Array.isArray(data.all_projects)) {
				// Transform project names into the required format
				const fetchedProjects = data.all_projects.map((projectName: string) => ({
					name: projectName,
					code: '' // We'll fetch individual project code when selected
				}));
				
				// Keep the default "Select a project..." option at the beginning
				projects = [
					{ name: 'Select a project...', code: '' },
					...fetchedProjects
				];
				
				// Reset selected project to default
				selectedProject = projects[0];
			}
		} catch (error) {
			console.error('Error fetching projects:', error);
			AddToast('Failed to load projects. Please try again later.', 'error');
		}
	}

	// Sync with global store
	let confirmed_code = '';
	const unsubscribe = confirmedSourceCode.subscribe(value => {
		confirmed_code = value;
		// Always update code_text when confirmed code changes
		code_text = value;
		// Set isConfirmed based on whether there is confirmed code
		isConfirmed = !!value;
	});

	
	// Add event listener for AI-generated source code
	let aiEventListener: (event: CustomEvent) => void;

	// Add another event listener for AI-submitted source code
	let aiSubmittedEventListener: (event: CustomEvent) => void;

	onMount(() => {
		// Check if user is a guest
		const accessToken = sessionStorage.getItem('access_token');
		isGuestUser = accessToken === 'guestuser';
		
		fetchProjects();
		
		// Listen for AI-generated source code
		aiEventListener = (event: CustomEvent) => {
			if (event.detail && event.detail.code) {
				// Store the previous code before replacing
				previous_code_text = code_text;
				// Replace the textarea content with AI-generated code
				code_text = event.detail.code;
				// Reset the default input flag since this is AI-generated
				isDefaultInput = false;
			}
		};

		// Listen for AI-submitted source code
		aiSubmittedEventListener = (event: CustomEvent) => {
			if (event.detail && event.detail.code) {
				// Update the confirmed source code store
				confirmedSourceCode.set(event.detail.code);
				isConfirmed = true;
				console.log('AI code submitted and confirmed, window staying open');
			}
		};

		window.addEventListener('ai-source-generated', aiEventListener);
		window.addEventListener('ai-source-submitted', aiSubmittedEventListener);
	});

	onDestroy(() => {
		setActivePhase(null);
		unsubscribe();
		
		if (aiEventListener) {
			window.removeEventListener('ai-source-generated', aiEventListener);
		}
		if (aiSubmittedEventListener) {
			window.removeEventListener('ai-source-submitted', aiSubmittedEventListener);
		}
	});

	function handleDefaultInput() {
		if (!isDefaultInput) {
			previous_code_text = code_text;
			code_text = 'int blue = 13;\n\n';
			code_text += 'int new(int red)\n';
			code_text += '{\n';
			code_text += '    red = red + 1;\n';
			code_text += '    return red;\n';
			code_text += '}\n\n';
			code_text += 'int _i = 0;\n\n';
			code_text += 'for _i range(12)\n';
			code_text += '{\n';
			code_text += '    blue = new(blue);\n';
			code_text += '    print(blue);\n';
			code_text += '}\n';
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
			AddToast('Invalid file type: Only .txt files are supported. Please upload a plain text file', 'error');
			input.value = '';
			return;
		}

		const reader = new FileReader();
		reader.onload = () => {
			code_text = reader.result as string;
			AddToast('File uploaded successfully! Your source code is ready to use', 'success');
		};
		reader.onerror = () => {
			AddToast('Upload failed: Unable to read the selected file. Please try again', 'error');
		};
		reader.readAsText(file);
	}

	async function submitCode() {
		if (!code_text.trim()) return;
		const user_id = localStorage.getItem('user_id');
		const project = get(projectName);
		if (!user_id) {
			AddToast('Authentication required: Please log in to save source code', 'error');
			return;
		}
		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}

    try {
        const res = await fetch('https://www.visual-compiler.co.za/api/lexing/code', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`
            },
            body: JSON.stringify({
                project_name: project,
                source_code: code_text
            })
        });
        
        if (!res.ok) {
            const errorData = await res.json();
            console.error('API Error:', errorData); // Debug log
            
            if (res.status === 401) {
                AddToast('Authentication failed: Please log in again', 'error');
            } else {
                AddToast(`Save failed: ${errorData.error || 'Unable to save source code'}`, 'error');
            }
            return;
        }
        
        const data = await res.json();
        console.log('Success response:', data); // Debug log
        
        AddToast('Source code saved successfully! Ready to begin lexical analysis', 'success');
        confirmedSourceCode.set(code_text);
        isConfirmed = true;
        
        // Only call onCodeSubmitted if we should close the window (manual submission)
        if (shouldCloseWindow) {
            onCodeSubmitted(code_text);
        }
        
        await tick();
    } catch (error) {
        console.error('Request failed:', error); // Debug log
        AddToast(`Save failed: ${(error as Error).message}. Please check your connection and try again`, 'error');
    }
}
	// --- NEW FUNCTION TO HANDLE PROJECT SELECTION ---
	async function handleProjectSelect() {
		if (!selectedProject || selectedProject.name === 'Select a project...') {
			code_text = '';
			isDefaultInput = false;
			return;
		}

		const userId = localStorage.getItem('user_id');
		if (!userId) {
			AddToast('Authentication required: Please log in to load project data', 'error');
			return;
		}

		try {
			const response = await fetch(
				`https://www.visual-compiler.co.za/api/users/getProject?project_name=${encodeURIComponent(selectedProject.name)}&users_id=${userId}`, 
				{
					method: 'GET',
					headers: {
						'accept': 'application/json'
					}
				}
			);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();
			
			if (data.message === "Retrieved users project details" && data.results?.translating?.code) {
				// Populate the text editor with the translation code
				code_text = data.results.translating.code.join('\n'); // Join array of code lines into string
				isDefaultInput = false;
				AddToast('Project code loaded successfully!', 'success');
			} else {
				// If no translation code exists, just clear the text area
				code_text = '';
				AddToast('Project loaded, but no translator code found', 'info');
			}
		} catch (error) {
			console.error('Error fetching project details:', error);
			AddToast('Failed to load project code. Please try again.', 'error');
			// Reset to empty if there's an error
			code_text = '';
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

	<div class="controls-grid" class:guest-mode={isGuestUser}>
		<div class="control-item">
			<label class="upload-btn">
				Upload File
				<input type="file" accept=".txt" on:change={handleFileChange} />
			</label>
		</div>

		{#if !isGuestUser}
			<div class="control-item project-selector">
				<label for="project-select">Import from Project</label>
				<select id="project-select" bind:value={selectedProject} on:change={handleProjectSelect}>
					{#each projects as project}
						<option value={project}>{project.name}</option>
					{/each}
				</select>
			</div>
		{/if}
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

	.controls-grid.guest-mode {
		grid-template-columns: 1fr;
		justify-items: center;
	}

	.controls-grid.guest-mode .upload-btn {
		width: 150px;
		padding: 0.5rem 1.5rem;
		justify-self: center;
	}

	.controls-grid.guest-mode + .controls .confirm-btn {
		width: 150px;
		justify-content: center;
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
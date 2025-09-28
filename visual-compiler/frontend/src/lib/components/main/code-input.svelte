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

	// Custom dropdown state
	let isDropdownOpen = false;

	// --- REAL PROJECTS DATA ---
	let projects: Array<{ name: string; code: string }> = [];
	let selectedProject = { name: 'Select a project...', code: '' };

	// Fetch projects from backend
	async function fetchProjects() {
		// Skip fetching projects for guest users
		if (isGuestUser) {
			console.log('Skipping project fetch - guest user'); // Debug log
			return;
		}

		const userId = sessionStorage.getItem('user_id');
		console.log('Fetching projects for userId:', userId); // Debug log
		
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
			
			if (data.all_projects && Array.isArray(data.all_projects)) {
				// Transform project names into the required format
				projects = data.all_projects.map((projectName: string) => ({
					name: projectName,
					code: '' // We'll fetch individual project code when selected
				}));
				
				// Keep selected project as placeholder
				selectedProject = { name: 'Select a project...', code: '' };
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

		// Close dropdown when clicking outside
		function handleClickOutside(event: MouseEvent) {
			const target = event.target as Element;
			if (!target.closest('.custom-dropdown')) {
				isDropdownOpen = false;
			}
		}

		window.addEventListener('ai-source-generated', aiEventListener);
		window.addEventListener('ai-source-submitted', aiSubmittedEventListener);
		document.addEventListener('click', handleClickOutside);

		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
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
    const project = get(projectName);
    
		// Check sessionStorage for auth token
		const accessToken = sessionStorage.getItem('access_token') || 
					   sessionStorage.getItem('authToken') || 
					   sessionStorage.getItem('token');    if (!accessToken) {
        AddToast('Authentication required: Please log in to save source code', 'error');
        return;
    }
    if (!project) {
        AddToast('No project selected: Please select or create a project first', 'error');
        return;
    }

    console.log('Using access token:', accessToken.substring(0, 20) + '...'); // Debug log

    try {
        const res = await fetch('http://localhost:8080/api/lexing/code', {
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

		const userId = sessionStorage.getItem('user_id');
		if (!userId) {
			AddToast('Authentication required: Please log in to load project data', 'error');
			return;
		}

		try {
			const response = await fetch(
				`http://localhost:8080/api/users/getProject?project_name=${encodeURIComponent(selectedProject.name)}&users_id=${userId}`, 
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

	// Custom dropdown functions
	function toggleDropdown() {
		isDropdownOpen = !isDropdownOpen;
	}

	function selectProject(project: { name: string; code: string }) {
		selectedProject = project;
		isDropdownOpen = false;
		handleProjectSelect();
	}

	function closeDropdown() {
		isDropdownOpen = false;
	}

	// Position dropdown menu dynamically
	function positionDropdown(node: HTMLElement) {
		if (!node) return;
		
		const rect = node.getBoundingClientRect();
		const viewportHeight = window.innerHeight;
		const spaceBelow = viewportHeight - rect.bottom;
		const spaceAbove = rect.top;
		const dropdownHeight = 150; // max-height of dropdown
		
		// If not enough space below but enough above, position upward
		if (spaceBelow < dropdownHeight && spaceAbove > dropdownHeight) {
			node.style.top = 'auto';
			node.style.bottom = '100%';
			node.style.borderTop = '2px solid #e2e8f0';
			node.style.borderBottom = 'none';
			node.style.borderRadius = '10px 10px 0 0';
		} else {
			node.style.top = '100%';
			node.style.bottom = 'auto';
			node.style.borderTop = 'none';
			node.style.borderBottom = '2px solid #e2e8f0';
			node.style.borderRadius = '0 0 10px 10px';
		}
	}

	// Reset confirmation flag if user changes the text
	$: isConfirmed = code_text === confirmed_code && !!code_text;
	$: displayed_text = isConfirmed ? `Current source code: ${code_text}` : code_text;

</script>

<div class="code-input-container">
	<!-- Header Section -->
	<div class="modal-header">
		<div class="header-content">
			<div class="header-icon">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
					<path d="M14 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V8L14 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
					<polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
				</svg>
			</div>
			<div class="header-text">
				<h2 class="modal-title">Source Code Input</h2>
				<p class="modal-subtitle">Enter or import your source code to begin compilation</p>
			</div>
		</div>
	</div>

	<!-- Code Editor Section -->
	<div class="editor-section">

		<div class="editor-wrapper">
			<textarea
				bind:this={textareaEl}
				bind:value={code_text}
				placeholder="Enter your source code here..."
				spellcheck="false"
			></textarea>
		</div>
	</div>

	<!-- Controls Section -->
	<div class="controls-section">
		<div class="options-stack">
			<div class="separator">
					<span>or</span>
				</div>
			<!-- Show Example Button -->
			<button
				type="button"
				class="option-btn example-btn"
				title={isDefaultInput ? 'Restore your input' : 'Show example code'}
				on:click={handleDefaultInput}
			>
				{#if isDefaultInput}
					Restore Input
				{:else}
					Show Example
				{/if}
			</button>

			<!-- OR Separator -->
			<div class="separator">
				<span>or</span>
			</div>

			<!-- Upload File Button -->
			<label class="option-btn upload-btn">
				<span>Upload File</span>
				<input type="file" accept=".txt" on:change={handleFileChange} />
			</label>

			<!-- Import from Project (only for non-guest users) -->
			{#if !isGuestUser && projects.length > 0}
				<!-- OR Separator -->
				<div class="separator">
					<span>or</span>
				</div>

				<div class="project-import-option">
					<div class="project-label">
						<span>Import from saved project:</span>
					</div>
					<div class="custom-dropdown">
						<button
							type="button"
							class="dropdown-button"
							class:open={isDropdownOpen}
							on:click={toggleDropdown}
						>
							<span class="dropdown-text" class:placeholder={selectedProject.name === 'Select a project...'}>{selectedProject.name}</span>
							<svg class="dropdown-arrow" class:rotated={isDropdownOpen} width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
								<polyline points="6 9 12 15 18 9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
							</svg>
						</button>
						{#if isDropdownOpen}
							<div class="dropdown-menu" use:positionDropdown on:click|stopPropagation>
								{#each projects.filter(p => p.name !== 'Select a project...') as project}
									<button
										type="button"
										class="dropdown-option"
										class:selected={selectedProject.name === project.name}
										on:click={() => selectProject(project)}
									>
										{project.name}
									</button>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<!-- Confirm Button -->
		<div class="confirm-section">
			<button
				type="button"
				class="confirm-btn"
				class:confirmed={isConfirmed}
				on:click={submitCode}
				disabled={!code_text.trim()}
			>
				{#if isConfirmed}
					Code Confirmed
				{:else}
					Confirm Code
				{/if}
			</button>
		</div>
	</div>
</div>

<style>
	.code-input-container {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: linear-gradient(145deg, #ffffff, #f8fafc);
		overflow: hidden;
		position: relative;
	}

	/* Header Section */
	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 2rem 2rem 1.5rem;
		border-bottom: 1px solid rgba(226, 232, 240, 0.8);
		background: linear-gradient(135deg, #f8fafc, #ffffff);
	}

	.header-content {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.header-icon {
		width: 48px;
		height: 48px;
		background: linear-gradient(135deg, #041a47, #0f2b5b);
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		box-shadow: 0 4px 12px rgba(4, 26, 71, 0.25);
	}

	.header-text {
		flex: 1;
	}

	.modal-title {
		margin: 0;
		font-size: 1.75rem;
		font-weight: 700;
		color: #1e293b;
		line-height: 1.2;
	}

	.modal-subtitle {
		margin: 0.25rem 0 0 0;
		font-size: 1rem;
		color: #64748b;
		font-weight: 500;
	}

	/* Editor Section */
	.editor-section {
		flex: 1;
		padding: 1.5rem 2rem 0.2rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		overflow: hidden;
		align-items: center;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 600;
		align-self: flex-start;
	}

	.status-badge.confirmed {
		background: linear-gradient(135deg, #10b981, #34d399);
		color: white;
		box-shadow: 0 2px 8px rgba(16, 185, 129, 0.25);
	}

	.editor-wrapper {
		flex: 1;
		position: relative;
		overflow: hidden;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		border: 2px solid #e2e8f0;
		transition: border-color 0.2s ease;
		width: 75%;
		max-width: 800px;
	}

	.editor-wrapper:focus-within {
		border-color: #041a47;
		box-shadow: 0 4px 12px rgba(4, 26, 71, 0.15);
	}

	textarea {
		width: 100%;
		height: 100%;
		min-height: 350px;
		padding: 1.5rem;
		border: none;
		outline: none;
		resize: none;
		font-family: 'JetBrains Mono', 'Fira Code', 'Monaco', 'Menlo', 'Consolas', monospace;
		font-size: 0.925rem;
		line-height: 1.6;
		color: #1e293b;
		background: #ffffff;
		tab-size: 2;
		white-space: pre;
		overflow-wrap: normal;
		overflow-x: auto;
		scrollbar-width: thin;
		scrollbar-color: #cbd5e1 transparent;
	}

	textarea::placeholder {
		color: #94a3b8;
		font-style: italic;
	}

	textarea::-webkit-scrollbar {
		width: 8px;
		height: 8px;
	}

	textarea::-webkit-scrollbar-track {
		background: transparent;
	}

	textarea::-webkit-scrollbar-thumb {
		background: #cbd5e1;
		border-radius: 4px;
	}

	textarea::-webkit-scrollbar-thumb:hover {
		background: #94a3b8;
	}

	/* Controls Section */
	.controls-section {
		padding: 0.5rem 2rem 1.5rem;
		background: linear-gradient(135deg, #f8fafc, #ffffff);
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
		position: relative;
		overflow: visible;
	}

	.options-stack {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		width: 100%;
		max-width: 400px;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.option-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.6rem 1rem;
		background: linear-gradient(135deg, #64748b, #748299);
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		overflow: hidden;
		box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
		text-decoration: none;
		width: 100%;
		max-width: 220px;
		justify-content: center;
	}

	.option-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
	}

	.example-btn {
		background: #1e40af;	
	}

	.example-btn:hover {
		box-shadow: 0 4px 12px rgba(30, 64, 175, 0.3);
	}

	.separator {
		display: flex;
		align-items: center;
		color: #64748b;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		position: relative;
		width: 100%;
		max-width: 280px;
		justify-content: center;
		margin: 0.125rem 0;
	}

	.separator::before,
	.separator::after {
		content: '';
		flex: 1;
		height: 1px;

	}

	.separator span {
		padding: 0 1rem;
		background: linear-gradient(135deg, #f8fafc, #ffffff);
		border-radius: 4px;
	}

	.upload-btn input[type="file"] {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		opacity: 0;
		cursor: pointer;
	}

	.project-import-option {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		width: 100%;
		max-width: 220px;
		align-items: center;
	}

	.project-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		font-weight: 600;
		color: #475569;
		text-align: center;
	}

	.custom-dropdown {
		position: relative;
		width: 100%;
		z-index: 10;
	}

	.dropdown-button {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		padding: 0.625rem 0.875rem;
		border: 2px solid #e2e8f0;
		border-radius: 10px;
		background: white;
		font-size: 0.925rem;
		font-weight: 500;
		color: #1e293b;
		cursor: pointer;
		transition: all 0.2s ease;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		text-align: left;
	}

	.dropdown-button:hover {
		border-color: #041a47;
		box-shadow: 0 0 0 3px rgba(4, 26, 71, 0.1);
	}

	.dropdown-button:focus {
		outline: none;
		border-color: #041a47;
		box-shadow: 0 0 0 3px rgba(4, 26, 71, 0.2);
	}

	.dropdown-button.open {
		border-color: #041a47;
		box-shadow: 0 0 0 3px rgba(4, 26, 71, 0.1);
	}

	.dropdown-text.placeholder {
		color: #64748b;
		font-style: italic;
	}

	.dropdown-arrow {
		transition: transform 0.2s ease;
		color: #64748b;
	}

	.dropdown-arrow.rotated {
		transform: rotate(180deg);
	}

	.dropdown-menu {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		z-index: 1000;
		background: white;
		border: 2px solid #e2e8f0;
		border-top: none;
		border-radius: 0 0 10px 10px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		max-height: 150px;
		overflow-y: auto;
		scrollbar-width: thin;
		scrollbar-color: #cbd5e1 #f1f5f9;
	}

	.dropdown-option {
		display: block;
		width: 100%;
		padding: 0.5rem 0.875rem;
		background: transparent;
		color: #1e293b;
		font-size: 0.925rem;
		font-weight: 500;
		border: none;
		cursor: pointer;
		transition: all 0.15s ease;
		text-align: left;
		border-radius: 0;
	}

	.dropdown-option:hover:not(.selected) {
		background: #e2e8f0 !important;
		color: #1e293b !important;
	}

	.dropdown-option:focus:not(.selected) {
		outline: none;
		background: #e2e8f0 !important;
		color: #1e293b !important;
	}

	.dropdown-option.selected {
		background: linear-gradient(135deg, #041a47, #0f2b5b) !important;
		color: white !important;
		font-weight: 600;
	}

	.dropdown-option.selected:hover {
		background: linear-gradient(135deg, #0f2b5b, #1e3a8a) !important;
		color: white !important;
	}

	.dropdown-option.placeholder {
		color: #64748b !important;
		font-style: italic;
	}

	.dropdown-option.placeholder:hover {
		background: #f8fafc !important;
		color: #475569 !important;
	}

	.dropdown-menu::-webkit-scrollbar {
		width: 8px;
	}

	.dropdown-menu::-webkit-scrollbar-track {
		background: #f1f5f9;
		border-radius: 0 0 8px 0;
	}

	.dropdown-menu::-webkit-scrollbar-thumb {
		background: #cbd5e1;
		border-radius: 4px;
		border: 1px solid #e2e8f0;
	}

	.dropdown-menu::-webkit-scrollbar-thumb:hover {
		background: #94a3b8;
	}

	.confirm-section {
		display: flex;
		justify-content: center;
	}

	.confirm-btn {
		display: flex;
		align-items: center;
		gap: 0.5em;
		padding: 0.7rem 1rem;
		background: #041a47;
		color: white;
		border: none;
		border-radius: 12px;
		font-size: 1rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.2s ease;
		box-shadow: 0 4px 16px rgba(4, 26, 71, 0.25);
		min-width: 180px;
		justify-content: center;
	}

	.confirm-btn:not(:disabled):hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(4, 26, 71, 0.35);
	}

	.confirm-btn.confirmed {

		box-shadow: 0 4px 16px rgba(16, 185, 129, 0.25);
	}

	.confirm-btn.confirmed:hover {
		box-shadow: 0 6px 20px rgba(16, 185, 129, 0.35);
	}

	.confirm-btn:disabled {
		background: #e2e8f0;
		color: #94a3b8;
		cursor: not-allowed;
		transform: none;
		box-shadow: none;
	}

	/* Dark Mode Support */
	:global(html.dark-mode) .code-input-container {
		background: linear-gradient(145deg, #1a2a4a, #2d3748);
	}

	:global(html.dark-mode) .modal-header {
		border-bottom-color: rgba(74, 85, 104, 0.8);
		background: linear-gradient(135deg, #1a2a4a, #2d3748);
	}

	:global(html.dark-mode) .modal-title {
		color: #f1f5f9;
	}

	:global(html.dark-mode) .modal-subtitle {
		color: #94a3b8;
	}

	:global(html.dark-mode) .header-icon {
		background: linear-gradient(135deg, #1e40af, #3b82f6);
	}

	:global(html.dark-mode) .editor-wrapper {
		border-color: #4a5568;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
	}

	:global(html.dark-mode) .editor-wrapper:focus-within {
		border-color: #3b82f6;
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
	}

	:global(html.dark-mode) textarea {
		background: #2d3748;
		color: #f1f5f9;
		scrollbar-color: #4a5568 transparent;
	}

	:global(html.dark-mode) textarea::placeholder {
		color: #64748b;
	}

	:global(html.dark-mode) textarea::-webkit-scrollbar-thumb {
		background: #4a5568;
	}

	:global(html.dark-mode) textarea::-webkit-scrollbar-thumb:hover {
		background: #64748b;
	}

	:global(html.dark-mode) .controls-section {
		border-top-color: rgba(74, 85, 104, 0.8);
		background: linear-gradient(135deg, #1a2a4a, #2d3748);
	}

	:global(html.dark-mode) .separator {
		color: #94a3b8;
	}

	:global(html.dark-mode) .separator::before,
	:global(html.dark-mode) .separator::after {
		background: linear-gradient(to right, transparent, #4a5568, transparent);
	}

	:global(html.dark-mode) .separator span {
		background: linear-gradient(135deg, #1a2a4a, #2d3748);
	}

	:global(html.dark-mode) .example-btn {
		background: linear-gradient(135deg, #1d4ed8, #2563eb);
	}

	:global(html.dark-mode) .project-label {
		color: #94a3b8;
	}

	:global(html.dark-mode) .dropdown-button {
		background: #2d3748;
		color: #f1f5f9;
		border-color: #4a5568;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.25);
	}

	:global(html.dark-mode) .dropdown-button:hover {
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	:global(html.dark-mode) .dropdown-button:focus,
	:global(html.dark-mode) .dropdown-button.open {
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
	}

	:global(html.dark-mode) .dropdown-text.placeholder {
		color: #94a3b8;
	}

	:global(html.dark-mode) .dropdown-arrow {
		color: #94a3b8;
	}

	:global(html.dark-mode) .dropdown-menu {
		background: #2d3748;
		border-color: #4a5568;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
		scrollbar-color: #4a5568 #2d3748;
	}

	:global(html.dark-mode) .dropdown-option {
		background: transparent;
		color: #f1f5f9;
	}

	:global(html.dark-mode) .dropdown-option:hover:not(.selected) {
		background: #4a5568 !important;
		color: #f1f5f9 !important;
	}

	:global(html.dark-mode) .dropdown-option:focus:not(.selected) {
		background: #4a5568 !important;
		color: #f1f5f9 !important;
	}

	:global(html.dark-mode) .dropdown-option.selected {
		background: linear-gradient(135deg, #1e40af, #3b82f6) !important;
		color: white !important;
	}

	:global(html.dark-mode) .dropdown-option.selected:hover {
		background: linear-gradient(135deg, #1d4ed8, #2563eb) !important;
		color: white !important;
	}

	:global(html.dark-mode) .dropdown-option.placeholder {
		color: #94a3b8 !important;
	}

	:global(html.dark-mode) .dropdown-option.placeholder:hover {
		background: #374151 !important;
		color: #cbd5e1 !important;
	}

	:global(html.dark-mode) .dropdown-menu::-webkit-scrollbar-track {
		background: #2d3748;
	}

	:global(html.dark-mode) .dropdown-menu::-webkit-scrollbar-thumb {
		background: #4a5568;
		border-color: #374151;
	}

	:global(html.dark-mode) .dropdown-menu::-webkit-scrollbar-thumb:hover {
		background: #64748b;
	}

	:global(html.dark-mode) .confirm-btn {
		background: linear-gradient(135deg, #1e40af, #3b82f6);
	}

	:global(html.dark-mode) .confirm-btn.confirmed {
		background: linear-gradient(135deg, #059669, #10b981);
	}

	:global(html.dark-mode) .confirm-btn:disabled {
		background: #374151;
		color: #6b7280;
	}

	/* Responsive Design */
	@media (max-width: 768px) {
		.modal-header {
			padding: 1.5rem 1.5rem 1rem;
		}

		.header-content {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.75rem;
		}

		.editor-section {
			padding: 1rem 1.5rem ;
		}

		.editor-wrapper {
			width: 100%;
		}

		.controls-section {
			padding: 1rem 1.5rem 1.5rem;
		}

		.options-stack {
			max-width: 100%;
		}

		.option-btn,
		.project-import-option {
			max-width: 100%;
		}

		.separator {
			max-width: 100%;
		}

		textarea {
			min-height: 250px;
			font-size: 0.875rem;
		}
	}

	@media (max-width: 480px) {
		.modal-header {
			padding: 1rem;
		}

		.modal-title {
			font-size: 1.5rem;
		}

		.editor-section {
			padding: 1rem;
		}

		.controls-section {
			padding: 1rem;
		}

		.options-stack {
			gap: 0.75rem;
		}

		.option-btn {
			padding: 0.875rem 1.25rem;
			font-size: 0.875rem;
		}

		textarea {
			min-height: 200px;
			padding: 1rem;
		}
	}
</style>
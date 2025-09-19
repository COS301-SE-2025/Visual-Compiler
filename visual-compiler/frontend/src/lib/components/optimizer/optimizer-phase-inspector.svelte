<script lang="ts">
    import { onMount } from 'svelte';
    import { optimizerState } from '$lib/stores/optimizer';
    import { projectName } from '$lib/stores/project';
    import { AddToast } from '$lib/stores/toast';
    import { get } from 'svelte/store';

    let selectedLanguage = 'Go';
    let selectedTechniques: string[] = [];
    let inputCode = `package main

func main() {
	
}`;

    // Default button functionality variables
    let showDefault = false;
    let userSelectedLanguage = 'Go';
    let userSelectedTechniques: string[] = [];
    let userInputCode = `package main

func main() {
	
}`;

    const languages = ['Java', 'Python', 'Go'];
    const techniques = ['Constant Folding', 'Dead Code Elimination', 'Loop Unrolling'];

    // Default constants
    const DEFAULT_LANGUAGE = 'Go';
    const DEFAULT_TECHNIQUES = ['Constant Folding', 'Dead Code Elimination', 'Loop Unrolling'];
    const DEFAULT_INPUT_CODE = `package main

import "fmt"

func main() {
	for blue := 1; blue < 5; blue++ {
		fmt.Println("[Block " + blue + "]")
		for red := 1; red <= blue; red++ {
			fmt.Println("	-> " + red)
			if blue == red {
				fmt.Println(blue + red)
			}
		}
	}
}
	
func nothing() (int) {
	var random int = 13
}`;

    // Initialize the store with default values on mount
    onMount(() => {
        updateStore();
    });

    // Default button functions
    function insertDefault() {
        // Save current user selections
        userSelectedLanguage = selectedLanguage;
        userSelectedTechniques = [...selectedTechniques];
        userInputCode = inputCode;
        
        // Set to default values
        selectedLanguage = DEFAULT_LANGUAGE;
        selectedTechniques = [...DEFAULT_TECHNIQUES];
        inputCode = DEFAULT_INPUT_CODE;
        showDefault = true;
        
        updateStore();
    }

    function removeDefault() {
        // Restore user selections
        selectedLanguage = userSelectedLanguage;
        selectedTechniques = [...userSelectedTechniques];
        inputCode = userInputCode;
        showDefault = false;
        
        updateStore();
    }

    function toggleTechnique(technique: string) {
        if (selectedTechniques.includes(technique)) {
            selectedTechniques = selectedTechniques.filter(t => t !== technique);
        } else {
            selectedTechniques = [...selectedTechniques, technique];
        }
        updateStore();
    }

    function updateStore() {
        optimizerState.update(state => ({
            ...state,
            selectedLanguage: selectedLanguage as 'Java' | 'Python' | 'Go',
            selectedTechniques,
            inputCode
        }));
    }

    async function handleSubmit() {
        if (!inputCode.trim() || selectedTechniques.length === 0) {
            AddToast('Please enter code and select at least one optimisation technique', 'error');
            return;
        }

        // Set optimizing state
        optimizerState.update(state => ({
            ...state,
            isOptimizing: true,
            optimizationError: null,
            optimizedCode: null
        }));

        try {
            const currentProjectName = get(projectName);
            
            // Try multiple ways to get the user ID
            let usersId = localStorage.getItem('users_id') || 
                         localStorage.getItem('user_id') || 
                         localStorage.getItem('userId') ||
                         sessionStorage.getItem('users_id') ||
                         sessionStorage.getItem('user_id');

            // If still not found, try to get from cookies or use a default
            if (!usersId) {
                // Check if there's a user session or auth token
                const authToken = localStorage.getItem('authToken') || localStorage.getItem('token');
                if (authToken) {
                    // You might need to decode the token to get user ID
                    // For now, let's use a default user ID for testing
                    usersId = 'default_user'; // Replace with actual user ID logic
                } else {
                    throw new Error('Please log in to use the optimizer');
                }
            }


            // Step 1: Store the source code
            const storeResponse = await fetch('http://localhost:8080/api/optimising/source_code', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    project_name: currentProjectName,
                    source_code: inputCode,
                    users_id: usersId
                })
            });

            if (!storeResponse.ok) {
                const errorData = await storeResponse.json();
                throw new Error(errorData.error || 'Failed to store source code');
            }

            const storeData = await storeResponse.json();

            // Step 2: Optimize the code
            const optimizePayload = {
                project_name: currentProjectName,
                users_id: usersId,
                constant_folding: selectedTechniques.includes('Constant Folding'),
                dead_code: selectedTechniques.includes('Dead Code Elimination'),
                loop_unrolling: selectedTechniques.includes('Loop Unrolling')
            };


            const optimizeResponse = await fetch('http://localhost:8080/api/optimising/optimise', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(optimizePayload)
            });

            if (!optimizeResponse.ok) {
                const errorData = await optimizeResponse.json();
                throw new Error(errorData.details || 'optimisation failed');
            }

            const optimizeData = await optimizeResponse.json();

            // Process the optimized code
            const optimizedCodeLines = optimizeData.optimised_code.split('\n');
            
            // Update store with results
            optimizerState.update(state => ({
                ...state,
                isOptimizing: false,
                optimizedCode: {
                    optimized: optimizedCodeLines,
                    language: selectedLanguage,
                    techniques: selectedTechniques,
                    performanceGains: {
                        executionTime: "Improved",
                        memoryUsage: "Reduced",
                        codeSize: "Optimized"
                    }
                }
            }));

            AddToast('Code optimisation completed successfully!', 'success');

        } catch (error) {
            console.error('Optimisation error:', error);
            
            optimizerState.update(state => ({
                ...state,
                isOptimizing: false,
                optimizationError: error
            }));

            AddToast(`Optimisation failed: ${error.message}`, 'error');
        }
    }

    // React to store changes
    $: selectedLanguage, updateStore();
    $: inputCode, updateStore();

    // Store user input when not showing default
    $: if (!showDefault && inputCode) {
        userInputCode = inputCode;
    }

    // Initialize from store, but keep default code if store is empty
    $: if ($optimizerState) {
        if (!showDefault) {
            selectedLanguage = $optimizerState.selectedLanguage;
            selectedTechniques = $optimizerState.selectedTechniques;
            // Only update inputCode from store if it's not empty, otherwise keep default
            if ($optimizerState.inputCode.trim()) {
                inputCode = $optimizerState.inputCode;
            }
        }
    }
</script>

<div class="inspector-container">
    <div class="inspector-header">
        <h2 class="inspector-title">Optimising</h2>
    </div>

    <div class="instructions-section">
        <div class="instructions-content">
            <h4 class="instructions-header">Instructions</h4>
            <p class="instructions-text">
                Select one or more optimisation techniques and enter your Go source code.
            </p>
        </div>
    </div>

    

    <div class="technique-section">
        <div class="technique-header">
            <h3 class="section-heading">Optimising Techniques</h3>
            <div class="default-section">
                {#if !showDefault}
                    <button
                        class="default-toggle-btn"
                        on:click={insertDefault}
                        type="button"
                        aria-label="Insert default input"
                        title="Insert default input"
                    >
                        <span class="icon">🪄</span>
                    </button>
                {/if}
                {#if showDefault}
                    <button
                        class="default-toggle-btn selected"
                        on:click={removeDefault}
                        type="button"
                        aria-label="Remove default input"
                        title="Remove default input"
                    >
                        <span class="icon">🧹</span>
                    </button>
                {/if}
            </div>
        </div>
        <div class="technique-buttons">
            {#each techniques as technique}
                <button 
                    class="technique-btn {selectedTechniques.includes(technique) ? 'selected' : ''}"
                    on:click={() => toggleTechnique(technique)}
                >
                    {technique}
                </button>
            {/each}
        </div>
    </div>

    <div class="input-section">
        <h3 class="section-heading">Source Code Input (Go)</h3>
        <div class="input-area">
            <textarea 
                bind:value={inputCode}
                placeholder="Enter your code here..."
                rows="18"
            ></textarea>
        </div>
    </div>

    <div class="submit-section">
        <button 
            class="submit-button"
            on:click={handleSubmit}
            disabled={!inputCode.trim() || selectedTechniques.length === 0 || $optimizerState.isOptimizing}
        >
            {#if $optimizerState.isOptimizing}
               
                Optimising...
            {:else}
                
                Optimise Code
            {/if}
        </button>
    </div>
</div>

<style>
    .inspector-container {
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        height: 100%;
        font-family: Arial, sans-serif;
    }

    .inspector-header {
        padding-bottom: 0.75rem;
        margin-bottom: -1rem;
    }

    .inspector-title {
        margin: 0;
    color: #8451C7;
        font-size: 2rem;
        text-align: center;
        font-weight: bold;
    }

    .section-heading {
        margin: 0 0 1rem 0;
        color: #000;
        font-size: 1.5rem;
        font-weight: bold;
    }

    /* Technique Section */
    .technique-section {
        display: flex;
        flex-direction: column;
    }

    .technique-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .technique-buttons {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .technique-btn {
        padding: 0.75rem 1rem;
        border: 2px solid #AFA2D7;
        background: #fff;
        color: #000;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.2s ease;
        text-align: left;
        width: 100%;
    }

    .technique-btn:hover {
        background: #f0f0f0;
    }

    .technique-btn.selected {
        background: #AFA2D7;
        color: #fff;
        border-color: #AFA2D7;
    }

    /* Default Section */
    .default-section {
        display: flex;
        align-items: center;
    }

    .default-toggle-btn {
        padding: 0.4rem 0.7rem;
        border-radius: 50%;
        border: 2px solid #AFA2D7;
        background: white;
        color: #8451C7;
        cursor: pointer;
        transition: all 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 2.3rem;
        width: 2.3rem;
    }

    .default-toggle-btn.selected {
        background: #F3E8FF;
        border-color: #8451C7;
    }

    .default-toggle-btn:hover,
    .default-toggle-btn:focus {
        background: #F3E8FF;
        border-color: #8451C7;
        transform: scale(1.05);
    }

    /* Input Section */
    .input-section {
        display: flex;
        flex-direction: column;
        flex: 1;
    }

    .input-area {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .input-area textarea {
        width: 100%;
        min-height: 350px;
        padding: 0.75rem;
        border: 2px solid #AFA2D7;
        border-radius: 4px;
        font-family: 'Courier New', monospace;
        font-size: 0.9rem;
        resize: vertical;
        background: #fff;
        color: #000;
        flex: 1;
    }

    .input-area textarea:focus {
        outline: none;
        border-color: #AFA2D7;
        box-shadow: 0 0 0 2px rgba(175, 162, 215, 0.2);
    }

    .input-area textarea::placeholder {
        color: #999;
        font-style: italic;
    }

    /* Submit Section */
    .submit-section {
        display: flex;
        justify-content: center;
        padding-top: 1rem;
        margin-top: auto;
    }

    .submit-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.75rem 2rem;
        background: linear-gradient(135deg, #AFA2D7, #AFA2D7);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
        box-shadow: 0 2px 4px rgba(175, 162, 215, 0.2);
        min-width: 160px;
    }

    .submit-button:hover:not(:disabled) {
        background: linear-gradient(135deg, #AFA2D7, #8451C7);
        transform: translateY(-1px);
        box-shadow: 0 4px 8px rgba(175, 162, 215, 0.3);
    }

    .submit-button:active:not(:disabled) {
        transform: translateY(0);
    }

    .submit-button:disabled {
        background: #9CA3AF;
        cursor: not-allowed;
        transform: none;
        box-shadow: none;
    }

    /* Dark Mode Styles */
    :global(html.dark-mode) .inspector-container {
        background: #1a2a4a;
    }

    :global(html.dark-mode) .inspector-title {
        color: #8451C7;
    }

    :global(html.dark-mode) .section-heading {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .technique-btn {
        background: #2d3748;
        color: #E5E7EB;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .technique-btn:hover {
        background: #374151;
    }

    :global(html.dark-mode) .technique-btn.selected {
        background: #8451C7;
        color: #fff;
        border-color: #8451C7;
    }

    :global(html.dark-mode) .input-area textarea {
        background: #2d3748;
        color: #E5E7EB;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .input-area textarea:focus {
        border-color: #8451C7;
    }

    :global(html.dark-mode) .input-area textarea::placeholder {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .submit-button:disabled {
        background: #4B5563;
        color: #9CA3AF;
    }

    :global(html.dark-mode) .default-toggle-btn {
        background: transparent;
        color: #AFA2D7;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .default-toggle-btn.selected {
        background: #2d3748;
        border-color: #8451C7;
        color: #AFA2D7;
    }

    :global(html.dark-mode) .default-toggle-btn:hover {
        background: rgba(45, 55, 72, 0.5);
        border-color: #8451C7;
    }

      .instructions-section {
        margin: 1.5rem 0 0 0;
        background: #F3E8FF;
        border-radius: 8px;
        border-left: 4px solid #AFA2D7;
        transition: background-color 0.3s ease, border-color 0.3s ease;
    }
    .instructions-content {
        padding: 1.25rem 1.5rem;
    }
    .instructions-header {
        margin: 0 0 0.75rem 0;
        font-size: 1rem;
        font-weight: 600;
        color: #8451C7;
        transition: color 0.3s ease;
    }
    .instructions-text {
        margin: 0;
        font-size: 0.9rem;
        line-height: 1.5;
        color: #4B3A6A;
        transition: color 0.3s ease;
    }
    :global(html.dark-mode) .instructions-section {
        background: #2D1B69;
        border-left-color: #8451C7;
    }
    :global(html.dark-mode) .instructions-header {
        color: #AFA2D7;
    }
    :global(html.dark-mode) .instructions-text {
        color: #C4B5FD;
    }
</style>

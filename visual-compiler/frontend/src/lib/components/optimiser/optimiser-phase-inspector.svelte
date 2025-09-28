<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { optimiserState } from '$lib/stores/optimiser';
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

    let aiOptimiserEventListener: (event: CustomEvent) => void;

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

        aiOptimiserEventListener = (event: CustomEvent) => {
            if (event.detail && event.detail.code) {
                console.log('Received AI optimiser code:', event.detail.code);

                // Save current user input if not showing default
                if (!showDefault) {
                    userInputCode = inputCode;
                }
                
                // Insert AI-generated code into the textarea
                inputCode = event.detail.code;
                
                // Reset default state since this is AI-generated
                showDefault = false;
                
                // Update the store with the new code
                updateStore();
                
                AddToast('AI optimiser code inserted into code input area!', 'success');
                
                console.log('Final optimiser input code:', inputCode);
            }
        };

        // Listen for the event with "z" spelling
        window.addEventListener('ai-optimiser-generated', aiOptimiserEventListener);
    });

    onDestroy(() => {
        if (aiOptimiserEventListener) {
            window.removeEventListener('ai-optimiser-generated', aiOptimiserEventListener);
        }
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
        optimiserState.update(state => ({
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

        optimiserState.update(state => ({
            ...state,
            isOptimising: true,
            optimisationError: null,
            optimisedCode: null
        }));

        try {
            const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
            const currentProjectName = get(projectName);
            
            if (!accessToken) {
                AddToast('Authentication required: Please log in to use the optimiser', 'error');
                return;
            }
            if (!currentProjectName) {
                AddToast('No project selected: Please select or create a project first', 'error');
                return;
            }

            const storeResponse = await fetch('http://localhost:8080/api/optimising/source_code', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({
                    project_name: currentProjectName,
                    source_code: inputCode
                })
            });

            if (!storeResponse.ok) {
                const errorData = await storeResponse.json();
                throw new Error(errorData.error || 'Failed to store source code');
            }

            const optimisePayload = {
                project_name: currentProjectName,
                constant_folding: selectedTechniques.includes('Constant Folding'),
                dead_code: selectedTechniques.includes('Dead Code Elimination'),
                loop_unrolling: selectedTechniques.includes('Loop Unrolling')
            };

            const optimiseResponse = await fetch('http://localhost:8080/api/optimising/optimise', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(optimisePayload)
            });

            if (!optimiseResponse.ok) {
                const errorData = await optimiseResponse.json();
                throw new Error(errorData.details || 'optimisation failed');
            }

            const optimiseData = await optimiseResponse.json();

            const optimisedCodeLines = optimiseData.optimised_code.split('\n');

            optimiserState.update(state => ({
                ...state,
                isOptimising: false,
                optimisedCode: {
                    optimised: optimisedCodeLines,
                    language: selectedLanguage,
                    techniques: selectedTechniques,
                    performanceGains: {
                        executionTime: "Improved",
                        memoryUsage: "Reduced",
                        codeSize: "Optimised"
                    }
                }
            }));

            AddToast('Code optimisation completed successfully!', 'success');

        } catch (error) {
            console.error('Optimisation error:', error);

            optimiserState.update(state => ({
                ...state,
                isOptimising: false,
                optimisationError: error
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
    $: if ($optimiserState) {
        if (!showDefault) {
            selectedLanguage = $optimiserState.selectedLanguage;
            selectedTechniques = $optimiserState.selectedTechniques;
            // Only update inputCode from store if it's not empty, otherwise keep default
            if ($optimiserState.inputCode.trim()) {
                inputCode = $optimiserState.inputCode;
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
                    <button
                        class="option-btn example-btn"
                        class:selected={showDefault}
                        on:click={showDefault ? removeDefault : insertDefault}
                        type="button"
                        aria-label={showDefault ? 'Restore Input' : 'Show Example'}
                        title={showDefault ? 'Restore your input' : 'Show example code'}
                    >
                        {showDefault ? 'Restore Input' : 'Show Example'}
                    </button>
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
            disabled={!inputCode.trim() || selectedTechniques.length === 0 || $optimiserState.isOptimising}
        >
            {#if $optimiserState.isOptimising}

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
        width: 150px;
        min-width: 150px;
        justify-content: center;
        margin-left: 1rem;
    }

    .example-btn {
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        box-shadow: 0 2px 8px rgba(132, 81, 199, 0.2);
    }

    .example-btn:hover {
        box-shadow: 0 4px 12px rgba(132, 81, 199, 0.3);
    }

    .option-btn:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
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

    :global(html.dark-mode) .example-btn {
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
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

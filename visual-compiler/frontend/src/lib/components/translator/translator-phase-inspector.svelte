<script lang="ts">
    export let source_code: string;
    import { AddToast } from '$lib/stores/toast';
    import { createEventDispatcher, onMount, onDestroy } from 'svelte';
    import { projectName } from '$lib/stores/project';
    import { get } from 'svelte/store'; 
    import { lexerState } from '$lib/stores/lexer';
    import { translatorState, updateTranslatorInputs, markTranslatorSubmitted, updateTranslatorArtifacts } from '$lib/stores/translator';

    const dispatch = createEventDispatcher();

    let rules = [{ tokenSequence: '', lines: [''] }];
    let isSubmitted = false;
    let translationSuccessful = false;
    let show_default_rules = false;

    // FIX: Add project change tracking
    let hasInitialized = false;
    let currentProjectName = '';

    // FIX: Force reinitialization when project changes
    $: if ($projectName !== currentProjectName) {
        console.log('Translator: Project changed from', currentProjectName, 'to', $projectName);
        
        // RESET component state immediately when project changes
        if (currentProjectName !== '' && $projectName !== currentProjectName) {
            rules = [{ tokenSequence: '', lines: [''] }];
            isSubmitted = false;
            translationSuccessful = false;
            show_default_rules = false;

            // FIX: Clear artifacts from parent component by dispatching empty data
            dispatch('translationreceived', []);
        }
        
        hasInitialized = false;
        currentProjectName = $projectName;
    }

    // FIX: Update store subscription to handle project changes and artifacts
    $: if ($translatorState && $projectName && (!hasInitialized || $projectName !== currentProjectName)) {
        console.log('Translator component initializing/reinitializing with state:', $translatorState);
        
        rules = [...$translatorState.rules];
        show_default_rules = $translatorState.show_default_rules || false;
        isSubmitted = $translatorState.isSubmitted;
        translationSuccessful = $translatorState.translationSuccessful || false;
        
        // Ensure we always have at least one empty rule
        if (rules.length === 0) {
            rules = [{ tokenSequence: '', lines: [''] }];
        }
        
        // FIX: Handle artifacts - only show if there is actual translated code
        if ($translatorState.hasTranslatedCode && $translatorState.translatedCode && $translatorState.translatedCode.length > 0) {
            const translatedCode = [...$translatorState.translatedCode];
            
            // Dispatch the translated code to parent component
            dispatch('translationreceived', translatedCode);
            
            console.log('Restored translated code with', translatedCode.length, 'lines');
        } else {
            // FIX: Explicitly clear artifacts if no translated code exists
            dispatch('translationreceived', []);
            
            console.log('No translated code to restore - cleared artifacts');
        }
        
        hasInitialized = true;
        console.log('Translator component initialized with:', { 
            rules: rules.length, 
            isSubmitted, 
            hasTranslatedCode: $translatorState.hasTranslatedCode,
            translatedLines: $translatorState.hasTranslatedCode ? $translatorState.translatedCode?.length : 0
        });
    }

    // Force initial rule to be visible
    $: if (rules.length === 0) {
        rules = [{ tokenSequence: '', lines: [''] }];
    }

	// --- DEFAULT RULES DATA ---
	const DEFAULT_TRANSLATION_RULES = [
		{
			tokenSequence: 'KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER, SEPARATOR',
			lines: ['add     rax, {INTEGER}', 'mov     [{IDENTIFIER}], rax']
		},
		{
			tokenSequence: 'KEYWORD, IDENTIFIER, OPEN_BRACKET, KEYWORD, IDENTIFIER, CLOSE_BRACKET, OPEN_SCOPE, IDENTIFIER, ASSIGNMENT, IDENTIFIER, OPERATOR, INTEGER, SEPARATOR, KEYWORD, IDENTIFIER, SEPARATOR, CLOSE_SCOPE',
			lines: ['func {IDENTIFIER}:', '     mov     rbx, [{IDENTIFIER}]', '     add     rbx, {INTEGER}', '     mov     [{IDENTIFIER}], rbx', '     return']
		},
		{
			tokenSequence: 'CONTROL, IDENTIFIER, CONTROL, OPEN_BRACKET, INTEGER, CLOSE_BRACKET, OPEN_SCOPE, IDENTIFIER, ASSIGNMENT, IDENTIFIER, OPEN_BRACKET, IDENTIFIER, CLOSE_BRACKET, SEPARATOR, KEYWORD, OPEN_BRACKET, IDENTIFIER, CLOSE_BRACKET, SEPARATOR, CLOSE_SCOPE',
			lines: ['func {CONTROL}:', '     jump    [{IDENTIFIER}], {INTEGER}', '     param   rcx, [{IDENTIFIER}]', '     call    {IDENTIFIER}', '     print   [{IDENTIFIER}]']
		}
	];

    /**
	 * addRule
	 * @description Adds a new, empty translation rule block to the list.
	 * @param {void}
	 * @returns {void}
	 */
	function addRule() {
		rules = [...rules, { tokenSequence: '', lines: [''] }];
		isSubmitted = false;
		handleRulesChange();
	}

	/**
	 * removeRule
	 * @description Removes a translation rule block by its index.
	 * @param {number} ruleIndex - The index of the rule to remove.
	 * @returns {void}
	 */
	function removeRule(ruleIndex: number) {
		if (rules.length > 1) {
			rules = rules.filter((_, i) => i !== ruleIndex);
			isSubmitted = false;
		} else {
			// If removing the last rule, replace it with an empty one
			rules = [{ tokenSequence: '', lines: [''] }];
			isSubmitted = false;
		}
		handleRulesChange();
	}

	/**
	 * addLine
	 * @description Adds a new, empty line to a specific rule.
	 * @param {number} ruleIndex - The index of the rule to add a line to.
	 * @returns {void}
	 */
	function addLine(ruleIndex: number) {
		rules[ruleIndex].lines = [...rules[ruleIndex].lines, ''];
		rules = rules;
		isSubmitted = false;
		handleRulesChange();
	}

	/**
	 * removeLine
	 * @description Removes a specific line from a rule.
	 * @param {number} ruleIndex - The index of the rule containing the line.
	 * @param {number} lineIndex - The index of the line to remove.
	 * @returns {void}
	 */
	function removeLine(ruleIndex: number, lineIndex: number) {
		if (rules[ruleIndex].lines.length > 1) {
			rules[ruleIndex].lines = rules[ruleIndex].lines.filter((_, i) => i !== lineIndex);
			rules = rules;
			isSubmitted = false;
			handleRulesChange();
		}
	}

	/**
	 * insertDefaultRules
	 * @description Populates the editor with pre-defined example translation rules.
	 * @param {void}
	 * @returns {void}
	 */
	function insertDefaultRules() {
		rules = JSON.parse(JSON.stringify(DEFAULT_TRANSLATION_RULES));
		show_default_rules = true;
		isSubmitted = false;
		updateTranslatorInputs({
			rules: [...rules],
			show_default_rules: true,
			isSubmitted: false
		});
	}

	/**
	 * removeDefaultRules
	 * @description Clears the editor and removes the example rules.
	 * @param {void}
	 * @returns {void}
	 */
	function removeDefaultRules() {
		rules = [{ tokenSequence: '', lines: [''] }];
		show_default_rules = false;
		isSubmitted = false;
		updateTranslatorInputs({
			rules: [...rules],
			show_default_rules: false,
			isSubmitted: false
		});
	}

	/**
	 * handleSubmit
	 * @description Validates the user's input and submits the translation rules to the backend.
	 * @param {void}
	 * @returns {Promise<void>}
	 */
	async function handleSubmit() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName);
		
		if (!accessToken) {
			AddToast('Authentication required: Please log in to save translation rules', 'error');
			return;
		}
		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}

		const isValid = rules.every(
			(rule) => rule.tokenSequence.trim() !== '' && rule.lines.every((line) => line.trim() !== '')
		);

		if (!isValid) {
			AddToast('Incomplete rules: Please fill in all token sequence and translation line fields', 'error');
			return;
		}

		const apiPayload = {
			project_name: project,
			translation_rules: rules.map((rule) => ({
				sequence: [rule.tokenSequence],
				translation: rule.lines
			}))
		};

		if (apiPayload.translation_rules.length === 0) {
			AddToast('No translation rules: Please add at least one translation rule before submitting', 'error');
			return;
		}

		console.log('Submitting translation input:', JSON.stringify(apiPayload, null, 2));

		try {
			const response = await fetch('http://localhost:8080/api/translating/readRules', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify(apiPayload)
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.details || 'Failed to submit translation rules.');
			}

			const result = await response.json();
			AddToast('Translation rules saved successfully! Ready to translate your code', 'success');
			isSubmitted = true;

			// Mark as submitted in store
			markTranslatorSubmitted();
			updateTranslatorInputs({
				isSubmitted: true
			});
			
		} catch (error: any) {
			console.error('Rule submission Error:', error);
			AddToast('Rule submission failed: ' + (error.message || 'Please check your connection and try again'), 'error');
		}
	}

	/**
	 * handleTranslate
	 * @description Initiates the final translation process using the submitted rules.
	 * @param {void}
	 * @returns {Promise<void>}
	 */
	async function handleTranslate() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName);
		
		if (!accessToken) {
			AddToast('Authentication required: Please log in to perform translation', 'error');
			return;
		}
		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}

		console.log('Requesting final translation from backend...');

		try {
			const response = await fetch('http://localhost:8080/api/translating/translate', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify({ project_name: project })
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.details || 'Failed to translate code.');
			}

			const result = await response.json();
			updateTranslatorArtifacts(result.code || []);
			AddToast('Translation complete! Your code has been successfully translated', 'success');
			translationSuccessful = true;

			// Update store
			updateTranslatorInputs({
				translationSuccessful: true
			});

			// Dispatch the translated code to the parent component
			dispatch('translationreceived', result.code);
		} catch (error: any) {
			console.error('Translation Error:', error);
			// Dispatch a new event for the error
			dispatch('translationerror', error);
			AddToast('Translation failed: ' + (error.message || 'Unable to translate code. Please check your rules and try again'), 'error');
		}
	}

    // Add event listener for AI-generated translator rules
    let aiTranslatorEventListener: (event: CustomEvent) => void;

    onMount(() => {
		// Ensure we always have at least one empty rule on mount
		if (rules.length === 0) {
			rules = [{ tokenSequence: '', lines: [''] }];
		}
		
		// Listen for AI-generated translator rules
		aiTranslatorEventListener = (event: CustomEvent) => {
			if (event.detail && event.detail.rules && Array.isArray(event.detail.rules)) {
				console.log('Received AI translator rules:', event.detail.rules);
				
				// Clear existing rules and populate with AI-generated rules
				rules = event.detail.rules.map(rule => ({
					tokenSequence: rule.sequence || '',
					lines: Array.isArray(rule.translation) ? rule.translation : ['']
				}));
				
				// Reset submission states
				isSubmitted = false;
				translationSuccessful = false;
				show_default_rules = false;
				
				// Force reactivity
				rules = [...rules];
				
				AddToast('AI translator rules inserted into translation editor!', 'success');
				
				console.log('Final translator rules:', rules);
			}
		};

		window.addEventListener('ai-translator-generated', aiTranslatorEventListener);
	});

	onDestroy(() => {
		if (aiTranslatorEventListener) {
			window.removeEventListener('ai-translator-generated', aiTranslatorEventListener);
		}
	});

	function clearAllInputs() {
        // Reset translation rules
        rules = [{ tokenSequence: '', lines: [''] }];
        
        // Reset states
        isSubmitted = false;
        translationSuccessful = false;
        show_default_rules = false;

        // Update store
        updateTranslatorInputs({
            rules: [...rules],
            show_default_rules: false,
            isSubmitted: false,
            translationSuccessful: false
        });
        
        AddToast('All translator inputs cleared successfully!', 'success');
    }

    // ADD: Missing input change tracking function
    function handleRulesChange() {
        if (hasInitialized) {
            updateTranslatorInputs({
                rules: [...rules],
                isSubmitted: false
            });
        }
        isSubmitted = false;
    }

    // Add reactive statement to track input changes
    $: if (hasInitialized) {
        // Watch for changes in rules and update store
        updateTranslatorInputs({
            rules: [...rules]
        });
    }
</script>

<div class="inspector-container">
	<div class="header-container">
		<h1 class="heading">TRANSLATING</h1>
	</div>

	<div class="section">
		<h3 class="section-heading1">Source Code</h3>
		<div class="code-block-wrapper">
			<pre class="code-block">{source_code || 'No source code available.'}</pre>
		</div>
	</div>

	<div class="instructions-section">
		<div class="instructions-content">
			<h4 class="instructions-header">Instructions</h4>
			<p class="instructions-text">
				Enter translation rules with a sequence of tokens and their corresponding target code. Link the token types in the target code using &#123; and &#125; on either side.
			</p>
		</div>
	</div>

	<div class="section">
		<div class="section-header">
			<h2 class="section-heading">Translation Rules</h2>

			<div class="button-group">
				<button
					class="option-btn example-btn"
					class:selected={show_default_rules}
					on:click={show_default_rules ? removeDefaultRules : insertDefaultRules}
					type="button"
					aria-label={show_default_rules ? 'Restore your input' : 'Show context-free grammar example'}
					title={show_default_rules ? 'Restore your input' : 'Show context-free grammar example'}
				>
					{show_default_rules ? 'Restore Input' : 'Show Example'}
			</button>
				<button
					class="clear-toggle-btn"
					on:click={clearAllInputs}
					type="button"
					aria-label="Clear all inputs" 
					title="Clear all inputs"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<polyline points="3 6 5 6 21 6" />
						<path d="m19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6m3 0V4c0-1 1-2 2-2h4c0-1 1-2 2-2v2" />
						<line x1="10" y1="11" x2="10" y2="17" />
						<line x1="14" y1="11" x2="14" y2="17" />
					</svg>
					</button>
			</div>

		</div>
		<div class="rules-container">
			{#each rules as rule, ruleIndex}
				<div class="rule-block">
					<div class="form-group">
						<div class="rule-header">
							<label class="form-label" for="token-seq-{ruleIndex}">Token Sequence</label>
							<button
								class="remove-btn"
								on:click={() => removeRule(ruleIndex)}
								title="Remove Rule">
								✕
							</button>
						</div>
						<input
							type="text"
							class="input-field"
							id="token-seq-{ruleIndex}"
							bind:value={rule.tokenSequence}
							on:input={handleRulesChange}
							placeholder="KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER..."
						/>
					</div>
					<div style="display: flex; justify-content: center; align-items: center; margin: 0.1rem 0;">
						<span style="font-size: 2rem; color: #888;">&#8595;</span>
					</div>

					{#each rule.lines as line, lineIndex}
						<div class="line-group">
							<input
								type="text"
								class="input-field"
								id="line-{ruleIndex}-{lineIndex}"
								bind:value={rules[ruleIndex].lines[lineIndex]}
								on:input={handleRulesChange}
								placeholder="Line {lineIndex + 1}"
							/>
							<button
								class="remove-line-btn"
								on:click={() => removeLine(ruleIndex, lineIndex)}
								disabled={rule.lines.length <= 1}
								title="Remove Line"
								aria-label="Remove Line"
								><svg
									xmlns="http://www.w3.org/2000/svg"
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									><polyline points="3 6 5 6 21 6" /><path
										d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
									/><line x1="10" y1="11" x2="10" y2="17" /><line
										x1="14"
										y1="11"
										x2="14"
										y2="17"
									/></svg
								></button
							>
						</div>
					{/each}
					<button class="add-line" on:click={() => addLine(ruleIndex)}>+ Add Line</button>
				</div>
			{/each}
		</div>
		
		<div class="button-container">
			<button class="add-rule-btn" on:click={addRule}>+ Add New Rule</button>
			<div class="action-buttons">
				<button class="action-btn submit" on:click={handleSubmit}>Submit Rules</button>
				<button 
					class="action-btn submit" 
					class:disabled={!isSubmitted}
					disabled={!isSubmitted}
					on:click={handleTranslate}
					title={isSubmitted ? "Translate code using submitted rules" : "Submit rules first"}
				>
					Translate Code
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	:root {
		--bg-primary: white;
		--bg-secondary: #eee;
		--bg-tertiary: #4a5568;
		--text-primary: black;
		--text-secondary: #a0aec0;
		--accent-blue: #3182ce;
		--accent-green: #38a169;
		--accent-red: #e53e3e;
		--accent-purple: #805ad5;
		--accent-orange: #6c757d;
		--border-color: #4a5568;
		--font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial,
			sans-serif;
	}

	.inspector-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		padding: 1.5rem;
		font-family: var(--font-sans);
		background-color: var(--bg-primary);
		color: var(--text-primary);
	}

	.header-container {
		display: flex;
		justify-content: center;
		align-items: center;
		position: relative;
	}

    .instructions-section {
		margin-top: 0.8rem;
		margin-bottom: 0.5rem;
		background: #f8f9fa;
		border-radius: 8px;
		border-left: 4px solid #bed2e6;
		transition: background-color 0.3s ease, border-color 0.3s ease;
	}

	.instructions-content {
		padding: 1.25rem 1.5rem;
	}

	.instructions-header {
		margin: 0 0 0.75rem 0;
		font-size: 1rem;
		font-weight: 600;
		color: #333;
		transition: color 0.3s ease;
	}

	.instructions-text {
		margin: 0;
		font-size: 0.9rem;
		line-height: 1.5;
		color: #555;
		transition: color 0.3s ease;
	}


	.default-toggle-btn {
		right: 0;
		background: white;
		border: 2px solid #e5e7eb;
		color: #001a6e;
		font-size: 1.2rem;
		cursor: pointer;
		transition: background 0.2s, border-color 0.2s;

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
		max-width: 150px;
		justify-content: center;
		margin-left: 1rem;
	}

	.example-btn {
		background: linear-gradient(135deg, #1e40af, #3b82f6);
	}

	.example-btn:hover {
		box-shadow: 0 4px 12px rgba(30, 64, 175, 0.3);
	}

	.option-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
	}

	.form-label {
		color: #1a2a4a;
		font-family: 'Times New Roman';
		font-weight: 500;
	}
	.section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.add-line {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		width: 100%;
		background-color: #eef2f7;
		color: #001a6e;
		border: 1px dashed #c0c7d3;
		border-color: #c0c7d3;
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.add-line:hover,
	.add-rule-btn:hover {
		border-color: #001a6e;
		transform: translateY(-2px);
	}
	.section-heading {
		color: #001a6e;
		margin-bottom: 0;
		margin-top: 0;
		font-family: 'Times New Roman';
		font-size: 1.25rem;
		font-weight: 600;
		text-align: center;
	}

	.section-heading1 {
		color: #444;
		margin-bottom: 0;
		margin-top: 0;
		font-family: 'Times New Roman';
	}

	.action-btn.translate {
		background-color: var(--accent-orange);
		color: white;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.25rem;
		cursor: pointer;
		font-size: 0.9rem;
		font-weight: 600;
		transition: background-color 0.2s, transform 0.1s;
		display: block;
		margin: 0 auto;
	}
	.heading {
		color: black;
		margin-bottom: 0;
		margin-top: 0;
		font-family: 'Times New Roman';
		text-align: center;
		flex-grow: 1;
	}

	.code-block-wrapper {
		background-color: #eee;
		border-radius: 0.375rem;
		padding: 1rem;
	}

	.code-block {
		font-family: 'Fira Code', monospace;
		font-size: 0.9rem;
		white-space: pre-wrap;
		word-wrap: break-word;
		margin: 0;
		max-height: 260px;
		overflow: auto;
	}

	.rules-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.button-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		margin-top: 0.5rem;
	}

	.action-buttons {
		display: flex;
		justify-content: center;
		gap: 1rem;
	}

	.rule-block {
		background-color: #f5f5f5;
		border-radius: 0.5rem;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		position: relative;
	}

	.rule-header {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.form-group,
	.line-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.line-group {
		flex-direction: row;
		align-items: center;
	}

	label {
		font-weight: 500;
		color: var(--text-secondary);
	}

	.input-field {
		width: 100%;
		padding: 0.75rem;
		background-color: var(--bg-primary);
		border: 1px solid #ddd;
		border-radius: 0.25rem;
		color: var(--text-primary);
		font-size: 1rem;
		transition: border-color 0.2s, box-shadow 0.2s;
	}

	.input-field:focus {
		outline: none;
		border-color: var(--accent-blue);
		box-shadow: 0 0 0 2px rgba(49, 130, 206, 0.5);
	}

	.action-btn {
		padding: 0.6rem 1.5rem;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		font-size: 0.9rem;
		font-weight: 500;
		transition: background-color 0.2s, transform 0.2s;
		width: 200px;
		max-width: 100%;
	}

	.add-rule-btn {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 0.5rem;
		background-color: #eef2f7;
		color: #001a6e;
		border: 1px dashed #c0c7d3;
		padding: 0.6rem 1.5rem;
		border-radius: 6px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		width: 200px;
		max-width: 100%;
		margin-bottom:0.5rem
	}

	.action-btn:hover,
	.submit:hover {
		transform: translateY(-2px);
	}

	.submit:hover {
		background-color: #a8bdd1;
	}

	.remove-btn,
	.remove-line-btn {
		background: transparent;
		color: var(--text-secondary);
		border: none;
		border-radius: 50%;
		cursor: pointer;
		font-size: 1.2rem;
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: background-color 0.2s, color 0.2s;
	}

	.remove-btn:hover,
	.remove-line-btn:hover {
		color: red;
	}

	.remove-line-btn {
		font-size: 1.5rem;
		margin-left: 0.5rem;
	}

	.remove-btn:disabled,
	.remove-line-btn:disabled {
		color: #4a5568;
		cursor: not-allowed;
		background-color: transparent;
	}

	.submit {
		background: #BED2E6;
		color: #000000;
		border: none;
		border-radius: 6px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s ease, transform 0.2s ease, opacity 0.2s ease;
	}

	.submit:hover:not(:disabled) {
		background: #a8bdd1;
		transform: translateY(-2px);
	}

	.submit:disabled,
	.translate:disabled,
	.translate.disabled {
		background: #d6d8db;
		color: #6c757d;
		cursor: not-allowed;
		opacity: 0.6;
		transform: none;
	}

	.submit:disabled:hover,
	.translate:disabled:hover,
	.translate.disabled:hover {
		background: #d6d8db;
		color: #6c757d;
		transform: none;
	}

	.button-group {
		position: absolute;
        right: 0;
        display: flex;
        gap: 0.5rem;
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
		width: 140px;
		min-width: 140px;
		justify-content: center;
		margin-left: 1rem;
	}

	.example-btn {
		background: #1e40af;
	}

	.example-btn:hover {
		box-shadow: 0 4px 12px rgba(30, 64, 175, 0.3);
	}

	.option-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
	}

	.option-btn.selected {
		background: #1e40af;
	}


    .clear-toggle-btn {
        background: white;
        border: 2px solid #e5e7eb;
        color: #ef4444;
        font-size: 1.2rem;
        cursor: pointer;
        transition: background 0.2s, border-color 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 2.2rem;
        width: 2.2rem;
        border-radius: 50%;
    }

    .clear-toggle-btn:hover {
        background: #fff5f5;
        border-color: #ef4444;
    }

	/* --- Dark Mode --- */
	:global(html.dark-mode) .inspector-container {
		background: #1a2a4a;
	}
	:global(html.dark-mode) .code-block-wrapper,
	:global(html.dark-mode) .rule-block {
		background: #2d3748;
	}
	:global(html.dark-mode) .heading,
	:global(html.dark-mode) .section-heading,
	:global(html.dark-mode) .section-heading1,
	:global(html.dark-mode) .form-label,
	:global(html.dark-mode) .code-block {
		color: #ebeef1;
	}
	:global(html.dark-mode) .input-field {
		background-color: #2d3748;
		border-color: #4b5563;
		color: #f0f0f0;
	}
	:global(html.dark-mode) .input-field:focus {
		border-color: #60a5fa;
		box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.4);
	}
	:global(html.dark-mode) .add-line,
	:global(html.dark-mode) .add-rule-btn {
		border-color: #60a5fa;
		color: #60a5fa;
		background: transparent;
	}
	:global(html.dark-mode) .example-btn {
		background: linear-gradient(135deg, #1d4ed8, #2563eb);
	}
	:global(html.dark-mode) .submit,
	:global(html.dark-mode) .translate {
		background-color: #001A6E;
		color: #ffffff;
	}

	:global(html.dark-mode) .submit:hover:not(:disabled),
	:global(html.dark-mode) .translate:hover:not(:disabled) {
		background-color: #002a8e;
	}

	:global(html.dark-mode) .submit:disabled,
	:global(html.dark-mode) .translate:disabled,
	:global(html.dark-mode) .translate.disabled {
		background-color: #495057;
		color: #6c757d;
		cursor: not-allowed;
		opacity: 0.6;
		transform: none;
	}
	:global(html.dark-mode) .clear-toggle-btn {
        background-color: #2d3748;
        border-color: #4a5568;
        color: #ef4444;
    }

    :global(html.dark-mode) .clear-toggle-btn:hover {
        background-color: #7f1d1d;
		border-color: #ef4444;
    }

	:global(html.dark-mode) .submit:disabled:hover,
	:global(html.dark-mode) .translate:disabled:hover,
	:global(html.dark-mode) .translate.disabled:hover {
		background-color: #495057;
		color: #6c757d;
		transform: none;
	}

	
	
	/* Add section-header style */
	.section-header {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		margin-bottom: 1rem;
		position: relative;
	}

	:global(html.dark-mode) .submit {
		background-color: #001A6E; 
		color: #ffffff;            
	}

	:global(html.dark-mode) .submit:hover {
		background-color: #002a8e;
	}

	:global(html.dark-mode) .submit:disabled {
		background-color: #2d3748;
		color: #9ca3af;
		border-color: #4a5568;
	}
</style>
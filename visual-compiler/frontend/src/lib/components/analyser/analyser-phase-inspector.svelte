<script lang="ts">
    import { slide, fade } from 'svelte/transition';
    import { AddToast } from '$lib/stores/toast';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    export let source_code: string = '';

    // --- INTERFACES ---
    interface ScopeRule {
        id: number;
        start: string;
        end: string;
    }

    interface TypeRule {
        id: number;
        resultType: string;
        assignment: string;
        lhs: string;
        operators: string;
        rhs: string;
    }

    interface GrammarRule {
        variable: string;
        type: string;
        functionDeclaration: string;
        parameter: string;
        assignment: string;
        operator: string;
        term: string;
    }

    // --- STATE ---

    // Scope Rules State
    let scopeRules: ScopeRule[] = [{ id: 0, start: '', end: '' }];
    let nextScopeId = 1;
    let showStartTooltip = false;
    let showEndTooltip = false;
    let submittedScopeRules: ScopeRule[] = []; // Renamed for clarity

    // Type Rules State
    let typeRules: TypeRule[] = [{ id: 0, resultType: '', assignment: '', lhs: '', operators: '', rhs: '' }];
    let nextTypeId = 1;
    let submittedTypeRules: TypeRule[] = [];

    // Grammar Rules State (single object as per screenshot, user only inputs 1 value for each)
    let grammarRules: GrammarRule = {
        variable: '',
        type: '',
        functionDeclaration: '',
        parameter: '',
        assignment: '',
        operator: '',
        term: ''
    };
    let submittedGrammarRules: GrammarRule = { ...grammarRules }; // Initialize with empty values

    // Overall Submission State
    let rulesSubmitted = false; // This flag now controls the overall submission state

    const DEFAULT_SCOPE_RULES = [
        { id: 0, start: '{', end: '}' }
    ];

    const DEFAULT_TYPE_RULES = [
        { id: 0, resultType: 'int', assignment: '=', lhs: 'INTEGER', operators: '', rhs: '' }
    ];

    const DEFAULT_GRAMMAR_RULES: GrammarRule = {
        variable: 'Identifier',
        type: 'TYPE',
        functionDeclaration: 'FUNCTION',
        parameter: 'PARAMETER',
        assignment: 'ASSIGNMENT',
        operator: 'OPERATOR',
        term: 'TERM'
    };
    let show_default_rules = false;

    // --- LOGIC ---

    // Scope Rules Logic
    function addScopeRow() {
        scopeRules = [...scopeRules, { id: nextScopeId++, start: '', end: '' }];
        rulesSubmitted = false;
    }

    function removeScopeRow(index: number) {
        scopeRules.splice(index, 1);
        scopeRules = scopeRules; // Trigger reactivity
        rulesSubmitted = false;
    }

    function handleScopeRuleInput() {
        rulesSubmitted = false;
    }

    // Type Rules Logic
    function addTypeRow() {
        typeRules = [...typeRules, { id: nextTypeId++, resultType: '', assignment: '', lhs: '', operators: '', rhs: '' }];
        rulesSubmitted = false;
    }

    function removeTypeRow(index: number) {
        typeRules.splice(index, 1);
        typeRules = typeRules; // Trigger reactivity
        rulesSubmitted = false;
    }

    function handleTypeRuleInput() {
        rulesSubmitted = false;
    }

    // Grammar Rules Logic (inputs are directly bound, no add/remove)
    function handleGrammarRuleInput() {
        rulesSubmitted = false;
    }

    // Universal Submission Logic
    function handleSubmit() {
        // Validate Scope Rules
        if (scopeRules.some((rule) => rule.start.trim() === '' || rule.end.trim() === '')) {
            AddToast('Please fill out all Scope Rule fields before submitting.', 'error');
            return;
        }

        // Validate Type Rules
        if (typeRules.some((rule) =>
            rule.resultType.trim() === '' ||
            rule.assignment.trim() === '' ||
            rule.lhs.trim() === '' ||
            rule.operators.trim() === '' ||
            rule.rhs.trim() === ''
        )) {
            AddToast('Please fill out all Type Rule fields before submitting.', 'error');
            return;
        }

        // Validate Grammar Rules (check if any field is empty)
        const grammarRuleFields = Object.values(grammarRules);
        if (grammarRuleFields.some(field => typeof field === 'string' && field.trim() === '')) {
            AddToast('Please fill out all Grammar Rule fields before submitting.', 'error');
            return;
        }

        // If all validations pass
        submittedScopeRules = JSON.parse(JSON.stringify(scopeRules));
        submittedTypeRules = JSON.parse(JSON.stringify(typeRules));
        submittedGrammarRules = JSON.parse(JSON.stringify(grammarRules));

        rulesSubmitted = true;
        AddToast('All rules submitted successfully!', 'success');
        console.log('Submitted Scope Rules:', submittedScopeRules);
        console.log('Submitted Type Rules:', submittedTypeRules);
        console.log('Submitted Grammar Rules:', submittedGrammarRules);
    }

    // Universal Reset Logic
    function resetState() {
        scopeRules = [{ id: 0, start: '', end: '' }];
        nextScopeId = 1;
        submittedScopeRules = [];

        typeRules = [{ id: 0, resultType: '', assignment: '', lhs: '', operators: '', rhs: '' }];
        nextTypeId = 1;
        submittedTypeRules = [];

        grammarRules = {
            variable: '',
            type: '',
            functionDeclaration: '',
            parameter: '',
            assignment: '',
            operator: '',
            term: ''
        };
        submittedGrammarRules = { ...grammarRules };

        rulesSubmitted = false;
        dispatch('reset');
        AddToast('All rules have been reset.', 'info');
    }

    function handleGenerate() {
        dispatch('generate');
        AddToast('Symbol table generated', 'success');
    }

    function insertDefaultRules() {
        scopeRules = DEFAULT_SCOPE_RULES.map(r => ({ ...r }));
        nextScopeId = 1;
        typeRules = DEFAULT_TYPE_RULES.map(r => ({ ...r }));
        nextTypeId = 1;
        grammarRules = { ...DEFAULT_GRAMMAR_RULES };
        show_default_rules = true;
        rulesSubmitted = false;
    }

    function removeDefaultRules() {
        scopeRules = [{ id: 0, start: '', end: '' }];
        nextScopeId = 1;
        typeRules = [{ id: 0, resultType: '', assignment: '', lhs: '', operators: '', rhs: '' }];
        nextTypeId = 1;
        grammarRules = {
            variable: '',
            type: '',
            functionDeclaration: '',
            parameter: '',
            assignment: '',
            operator: '',
            term: ''
        };
        show_default_rules = false;
        rulesSubmitted = false;
    }

    // --- REACTIVE STATEMENTS ---

    // Scope Rules Completeness
    $: lastScopeRowComplete =
        scopeRules.length > 0 &&
        scopeRules[scopeRules.length - 1].start.trim() !== '' &&
        scopeRules[scopeRules.length - 1].end.trim() !== '';

    $: allScopeRowsComplete = scopeRules.every(
        (rule) => rule.start.trim() !== '' && rule.end.trim() !== ''
    );

    // Type Rules Completeness
    $: lastTypeRowComplete =
        typeRules.length > 0 &&
        typeRules[typeRules.length - 1].resultType.trim() !== '' &&
        typeRules[typeRules.length - 1].assignment.trim() !== '' &&
        typeRules[typeRules.length - 1].lhs.trim() !== '' &&
        typeRules[typeRules.length - 1].operators.trim() !== '' &&
        typeRules[typeRules.length - 1].rhs.trim() !== '';

    $: allTypeRowsComplete = typeRules.every(
        (rule) =>
            rule.resultType.trim() !== '' &&
            rule.assignment.trim() !== '' &&
            rule.lhs.trim() !== '' &&
            rule.operators.trim() !== '' &&
            rule.rhs.trim() !== ''
    );

    // Grammar Rules Completeness
    $: allGrammarRulesComplete = Object.values(grammarRules).every(
        (field) => typeof field === 'string' && field.trim() !== ''
    );

    // Overall Completeness for Submit Button
    $: overallAllRowsComplete = allScopeRowsComplete && allTypeRowsComplete && allGrammarRulesComplete;
</script>

<div class="panel-container">
    <div class="default-toggle-wrapper">
        <button
            class="default-toggle-btn"
            class:selected={show_default_rules}
            on:click={show_default_rules ? removeDefaultRules : insertDefaultRules}
            type="button"
            aria-label={show_default_rules ? 'Remove default rules' : 'Insert default rules'}
            title={show_default_rules ? 'Remove default rules' : 'Insert default rules'}
        >
            <span class="icon">{show_default_rules ? '🧹' : '🪄'}</span>
        </button>
    </div>
    <div class="header">
        {#if rulesSubmitted}
            <div class="reset-wrapper" transition:fade={{ duration: 150 }}>
                <button class="reset-button" on:click={resetState} aria-label="Reset rules">
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
                        <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                        <path d="M3 3v5h5" />
                    </svg>
                </button>
            </div>
        {/if}
        <h1 class="analyser-heading">ANALYSING</h1>
    </div>

    <!-- Source Code Section -->
    <div class="source-code-section">
        <h3 class="source-code-header">Source Code</h3>
        <pre class="source-display">{source_code || 'No source code available'}</pre>
    </div>

    <div class="scrollable-content">
        <!-- Scope Rules Box -->
        <div class="analyser-box">
            <div class="analyser-box-header">
                <h2 class="heading2">Scope Rules</h2>
            </div>
            <div class="rules-list">
                {#each scopeRules as rule, i (rule.id)}
                    <div class="rule-row" in:slide={{ duration: 250 }}>
                        <div class="input-tooltip-wrapper">
                            <input
                                type="text"
                                bind:value={rule.start}
                                placeholder="Start Delimiter"
                                aria-label="Start Delimiter"
                                class="scope-input"
                                on:focus={() => i === 0 && (showStartTooltip = true)}
                                on:blur={() => i === 0 && (showStartTooltip = false)}
                                on:input={handleScopeRuleInput}
                            />
                            {#if i === 0 && showStartTooltip}
                                <div class="tooltip">Examples include <code>{'{'}</code></div>
                            {/if}
                        </div>
                        <span class="arrow-separator">→</span>
                        <div class="input-tooltip-wrapper">
                            <input
                                type="text"
                                bind:value={rule.end}
                                placeholder="End Delimiter"
                                aria-label="End Delimiter"
                                class="scope-input"
                                on:focus={() => i === 0 && (showEndTooltip = true)}
                                on:blur={() => i === 0 && (showEndTooltip = false)}
                                on:input={handleScopeRuleInput}
                            />
                            {#if i === 0 && showEndTooltip}
                                <div class="tooltip">Examples include <code>{'}'}</code></div>
                            {/if}
                        </div>
                        <button
                            class="icon-button delete-button"
                            on:click={() => removeScopeRow(i)}
                            aria-label="Remove scope rule row"
                            disabled={scopeRules.length === 1}
                        >
                            <svg
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
                                /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg
                        >
                    </div>
                {/each}
            </div>
            {#if !rulesSubmitted && lastScopeRowComplete}
                <div class="add-button-wrapper" transition:slide={{ duration: 150 }}>
                    <button class="add-button" on:click={addScopeRow} aria-label="Add new scope delimiter row">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        ><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg
                        >
                        <span>Add Scope Rule</span>
                    </button>
                </div>
            {/if}
        </div>

        <!-- Type Rules Box -->
        <div class="analyser-box">
            <div class="analyser-box-header">
                <h2 class="heading2">Type Rules</h2>
            </div>
            <div class="rules-list type-rules-list">
                {#each typeRules as rule, i (rule.id)}
                    <div class="rule-row" in:slide={{ duration: 250 }}>
                        <div class="type-rule-inputs">
                            <input type="text" bind:value={rule.resultType} placeholder="Result Type" aria-label="Result type" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.assignment} placeholder="Assignment" aria-label="Assignment" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.lhs} placeholder="LHS" aria-label="LHS" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.operators} placeholder="Operators" aria-label="Operators" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.rhs} placeholder="RHS" aria-label="RHS" class="scope-input" on:input={handleTypeRuleInput} />
                        </div>
                        <button
                            class="icon-button delete-button"
                            on:click={() => removeTypeRow(i)}
                            aria-label="Remove type rule row"
                            disabled={typeRules.length === 1}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg>
                        </button>
                    </div>
                {/each}
            </div>
            {#if !rulesSubmitted && lastTypeRowComplete}
                <div class="add-button-wrapper" transition:slide={{ duration: 150 }}>
                    <button class="add-button" on:click={addTypeRow} aria-label="Add new type rule row">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2.5"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        ><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg
                        >
                        <span>Add Type Rule</span>
                    </button>
                </div>
            {/if}
        </div>

        <!-- Linking Grammar Rules Box -->
        <div class="analyser-box">
            <div class="analyser-box-header">
                <h2 class="heading2">Linking Grammar Rules</h2>
            </div>
            <div class="grammar-rules-grid">
                <div class="grammar-rule-item">
                    <label for="grammar-variable" class="grammar-label">Variable</label>
                    <input id="grammar-variable" type="text" bind:value={grammarRules.variable} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-type" class="grammar-label">Type</label>
                    <input id="grammar-type" type="text" bind:value={grammarRules.type} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-function-declaration" class="grammar-label">Function Declaration</label>
                    <input id="grammar-function-declaration" type="text" bind:value={grammarRules.functionDeclaration} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-parameter" class="grammar-label">Parameter</label>
                    <input id="grammar-parameter" type="text" bind:value={grammarRules.parameter} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-assignment" class="grammar-label">Assignment</label>
                    <input id="grammar-assignment" type="text" bind:value={grammarRules.assignment} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-operator" class="grammar-label">Operator</label>
                    <input id="grammar-operator" type="text" bind:value={grammarRules.operator} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-term" class="grammar-label">Term</label>
                    <input id="grammar-term" type="text" bind:value={grammarRules.term} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
            </div>
        </div>
    </div>
    <div class="actions-container">
        {#if !rulesSubmitted}
            <button class="submit-button" on:click={handleSubmit} disabled={!overallAllRowsComplete}>
                Submit All Rules
            </button>
        {:else}
            <div class="generate-wrapper" transition:slide|local={{ duration: 250 }}>
                <button class="generate-button" on:click={handleGenerate}> Generate Symbol Table </button>
            </div>
        {/if}
    </div>
</div>

<style>
    /* General Heading Styles (existing) */
    .heading {
        font-weight: 600;
        margin-bottom: 2rem;
        text-align: center;
        margin-top: 0;
    }
    .heading2 {
        font-size: 1.25rem;
        font-weight: 600;
        margin: 0 0 0.25rem 0;
        color: #1a2a4a; /* Adjusted for consistency */
    }
    /* New style for spacing between sections */
    .new-section-heading {
        margin-top: 1.5rem;
        margin-bottom: 0.75rem; /* Adjusted for better spacing with inputs */
    }

    /* Panel Layout */
    .panel-container {
        display: flex;
        flex-direction: column;
        padding: 1.5rem 1rem;
        height: 100%; /* Important for flexbox to work correctly */
    }

    .header {
        position: relative;
        flex-shrink: 0; /* Prevents header from shrinking */
    }

    /* Scrollable Content Area */
    .scrollable-content {
        flex-grow: 1; /* Allows this div to take up available vertical space */
        overflow-y: auto; /* Enables vertical scrolling */
        padding-right: 0.5rem; /* Space for the scrollbar */
        margin-right: -0.5rem; /* Compensate for scrollbar padding */
        padding-bottom: 1rem; /* Padding before the actions container */
    }

    /* Rules List (for Scope & Type rules) */
    .rules-list {
        margin-bottom: 0.5rem; /* Space before add button */
    }
    .rule-row {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 0.75rem;
    }
    .input-tooltip-wrapper {
        position: relative;
        flex: 1;
        display: flex;
    }
    .scope-input {
        flex: 1;
        padding: 0.6rem 0.8rem;
        font-size: 0.9rem;
        border: 1px solid #d0d7e0;
        border-radius: 6px;
        background-color: #ffffff;
        transition:
            border-color 0.2s ease,
            box-shadow 0.2s ease;
    }
    .scope-input:focus {
        outline: none;
        border-color: #074799;
        box-shadow: 0 0 0 3px rgba(7, 71, 153, 0.15);
    }
    .arrow-separator {
        color: #7a8aa3;
        font-size: 1.2rem;
        font-weight: 600;
        flex-shrink: 0;
    }

   
    .type-rules-list .rule-row {
        align-items: flex-start; 
    }
    .type-rule-inputs {
        flex: 1;
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(80px, 1fr)); 
        gap: 0.75rem;
    }
 
    .type-rule-inputs .scope-input {
        width: auto; 
    }


   
    .grammar-rules-grid {
        display: grid;
        grid-template-columns: auto 1fr; 
        gap: 0.75rem 1rem;
        align-items: center;
        padding-top: 0.5rem; 
		max-width: 90%;
    }
    .grammar-rule-item {
        display: contents; 
    }
    .grammar-label {
        font-weight: 500;
        color: #1a2a4a;
         text-align: left;
        padding-right: 0.5rem;
    }


    /* Icon Button Styles (existing) */
    .icon-button {
        display: flex;
        align-items: center;
        justify-content: center;
        background: none;
        border: none;
        padding: 0.25rem;
        border-radius: 50%;
        cursor: pointer;
        color: #7a8aa3;
        transition:
            background-color 0.2s ease,
            color 0.2s ease;
    }
    .icon-button:hover:not(:disabled) {
        background-color: #eef2f7;
    }
    .delete-button:hover:not(:disabled) {
        color: #d9534f;
    }
    .icon-button:disabled {
        color: #c0c7d3;
        cursor: not-allowed;
    }
    .tooltip {
        position: absolute;
        top: 110%;
        left: 0;
        z-index: 10;
        background: #fff;
        color: #1a2a4a;
        border: 1px solid #d0d7e0;
        border-radius: 6px;
        padding: 0.5rem 0.75rem;
        font-size: 0.85rem;
        box-shadow: 0 2px 8px rgba(30, 40, 80, 0.08);
        white-space: nowrap;
        margin-top: 0.2rem;
        transition: opacity 0.15s;
    }
    .tooltip code {
        background: #eef2f7;
        padding: 0.1em 0.3em;
        border-radius: 4px;
        font-size: 0.92em;
    }
    .actions-container {
        padding-top: 1rem;
        margin-top: auto; /* Pushes to the bottom */
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        flex-shrink: 0; /* Prevents action container from shrinking */
    }
    .add-button-wrapper,
    .generate-wrapper {
        display: flex;
        justify-content: center;
    }
    .add-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        width: 100%;
        background-color: #eef2f7;
        color: #001a6e;
        border: 1px dashed #c0c7d3;
        padding: 0.5rem 1rem;
        border-radius: 6px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
    }
    .add-button:hover {
        background-color: #e1e8f0;
        border-color: #001a6e;
    }
    .submit-button,
    .generate-button {
        width: 100%;
        padding: 0.7rem 1rem;
        font-size: 0.9rem;
        font-weight: 600;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        transition:
            background-color 0.2s ease,
            opacity 0.2s ease;
    }
    .submit-button {
        background-color: #001a6e;
        color: white;
    }
    .submit-button:hover:not(:disabled) {
        background-color: #074799;
    }
    .submit-button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    .generate-button {
        background-color: #6c757d;
        color: white;
    }
    .generate-button:hover {
        background-color: #565e64;
    }

    .reset-wrapper {
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        z-index: 10;
    }

    .reset-button {
        background: none;
        border: none;
        padding: 0;
        margin: 0;
        cursor: pointer;
        color: #7a8aa3;
        transition: color 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
    }

    .reset-button:hover {
        color: #001a6e;
    }

    /* --- DARK MODE STYLES --- */
    :global(html.dark-mode) .heading,
    :global(html.dark-mode) .heading2,
    :global(html.dark-mode) .grammar-label {
        color: #e0e8f0;
    }

    :global(html.dark-mode) .scope-input {
        background-color: #253352;
        border-color: #3a4a6a;
        color: #e0e8f0;
    }
    :global(html.dark-mode) .scope-input::placeholder {
        color: #6a7b99;
    }
    :global(html.dark-mode) .scope-input:focus {
        border-color: #3a7bd5;
        box-shadow: 0 0 0 3px rgba(58, 123, 213, 0.2);
    }
    :global(html.dark-mode) .arrow-separator {
        color: #8a9bb3;
    }
    :global(html.dark-mode) .icon-button {
        color: #8a9bb3;
    }
    :global(html.dark-mode) .icon-button:hover:not(:disabled) {
        background-color: #2a4a8a;
    }
    :global(html.dark-mode) .delete-button:hover:not(:disabled) {
        color: #ff8a8a;
    }
    :global(html.dark-mode) .icon-button:disabled {
        color: #5a6a8a;
    }
    :global(html.dark-mode) .tooltip {
        background: #253352;
        color: #e0e8f0;
        border: 1px solid #3a4a6a;
        box-shadow: 0 2px 8px rgba(30, 40, 80, 0.18);
    }
    :global(html.dark-mode) .tooltip code {
        background: #2a4a8a;
        color: #e0e8f0;
    }
    :global(html.dark-mode) .add-button {
        background-color: transparent;
        color: #9ab0d3;
        border-color: #3a4a6a;
    }
    :global(html.dark-mode) .add-button:hover {
        background-color: #2a4a8a;
        color: #e0e8f0;
        border-color: #3a7bd5;
    }
    :global(html.dark-mode) .submit-button {
        background-color: #3a7bd5;
    }
    :global(html.dark-mode) .submit-button:hover:not(:disabled) {
        background-color: #4a8ce5;
    }
    :global(html.dark-mode) .generate-button {
        background-color: #4a5568;
    }
    :global(html.dark-mode) .generate-button:hover {
        background-color: #3b4453;
    }
    :global(html.dark-mode) .reset-button {
        color: #8a9bb3;
    }
    :global(html.dark-mode) .reset-button:hover {
        color: #e0e8f0;
    }

	   ::-webkit-scrollbar {
        width: 11px;
    }
    ::-webkit-scrollbar-track {
        background: #f1f1f1;
    }
    ::-webkit-scrollbar-thumb {
        background-color: #888;
        border-radius: 10px;
    }
    ::-webkit-scrollbar-thumb:hover {
        background: #555;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-track {
        background: #2d3748;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb {
        background-color: #4a5568; 
        border-color: #2d3748; 
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb:hover {
        background: #616e80;
    }

    .default-toggle-wrapper {
        position: absolute;
        top: 1rem;
        right: 1rem;
        z-index: 20;
        margin-bottom: 0;
    }
    .panel-container {
        position: relative;
    }
    .default-toggle-btn {
        background: white;
        border: 2px solid #e5e7eb;
        color: #001a6e;
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
    .default-toggle-btn.selected {
        background: #d0e2ff;
        border-color: #003399;
    }
    .default-toggle-btn:hover {
        background: #f5f8fd;
        border-color: #7da2e3;
    }
    .icon {
        font-size: 1.3rem;
        line-height: 1;
        pointer-events: none;
    }

    .analyser-heading {
        color: #001a6e;
        text-align: center;
        margin-top: 0;
        font-family: 'Times New Roman';
        margin-bottom: 0;
        font-size: 2rem;
        font-weight: bold;
    }
    .source-code-section {
        margin-bottom: 1rem;
    }
    .source-code-header {
        color: #444;
        margin-bottom: 0.5rem;
        font-family: 'Times New Roman';
    }
    .source-display {
        background: #f5f5f5;
        padding: 1rem;
        border-radius: 4px;
        overflow-x: auto;
        font-family: monospace;
        white-space: pre-wrap;
        margin: 0;
    }
    .analyser-box {
        background: #f8f9fa;
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        margin-bottom: 1.5rem;
    }
    .analyser-box-header {
        display: flex;
        align-items: center;
        margin-bottom: 1rem;
    }
    .analyser-box .heading2 {
        margin: 0;
        color: #001a6e;
        font-family: 'Times New Roman';
    }
    @media (max-width: 700px) {
        .analyser-box {
            padding: 1rem;
        }
    }
    :global(html.dark-mode) .analyser-heading,
    :global(html.dark-mode) .analyser-box .heading2,
    :global(html.dark-mode) .source-code-header {
        color: #ebeef1;
    }
    :global(html.dark-mode) .analyser-box,
    :global(html.dark-mode) .source-display {
        background: #2d3748;
    }
</style>
<script lang="ts">
    import { slide, fade } from 'svelte/transition';
    import { AddToast } from '$lib/stores/toast';
    import { createEventDispatcher } from 'svelte';
    import type { SymbolTable } from '$lib/types';
    import { projectName } from '$lib/stores/project';
	import { get } from 'svelte/store'; 
	import { error } from '@sveltejs/kit';
    import { lexerState } from '$lib/stores/lexer';

    const dispatch = createEventDispatcher();

    export let source_code: string = '';
    let show_symbol_table = false;
    let symbol_table: SymbolTable = { symbols: [] };
    let is_loading = false;
    export let onGenerateSymbolTable: (data: {
		symbol_table : Symbol[];
        analyser_error: string;
        analyser_error_details: string;
	}) => void = () => {};

    // --- INTERFACES ---
    interface ScopeRule {
        id: number;
        Start: string;
        End: string;
    }

    interface TypeRule {
        id: number;
        ResultData: string;
        Assignment: string;
        LHSData: string;
        Operator: string[];
        RHSData: string;
    }

    interface GrammarRule {
        VariableRule: string;
        TypeRule: string;
        FunctionRule: string;
        ParameterRule: string;
        AssignmentRule: string;
        OperatorRule: string;
        TermRule: string;
    }

    // --- STATE ---

    // Scope Rules State
    let scope_rules: ScopeRule[] = [{ id: 0, Start: '', End: '' }];
    let next_scope_id = 1;
    let show_start_tooltip = false;
    let show_end_tooltip = false;
    let submitted_scope_rules: ScopeRule[] = []; // Renamed for clarity

    // Type Rules State
    let type_rules: TypeRule[] = [{ id: 0, ResultData: '',Assignment: '', LHSData: '', Operator: [''], RHSData: '' }];
    let next_type_id = 1;
    let submitted_type_rules: TypeRule[] = [];

    // Grammar Rules State (single object as per screenshot, user only inputs 1 value for each)
    let grammar_rules: GrammarRule = {
        VariableRule: '',
        TypeRule: '',
        FunctionRule: '',
        ParameterRule: '',
        AssignmentRule: '',
        OperatorRule: '',
        TermRule: ''
    };
    let submitted_grammar_rules: GrammarRule = { ...grammar_rules }; // Initialize with empty values

    // Overall Submission State
    let rules_submitted = false; // This flag now controls the overall submission state
    let hasInitialized = false; // Flag to prevent reactive statement from re-running

    const DEFAULT_SCOPE_RULES = [
        { id: 0, Start: '{', End: '}' }
    ];

    const DEFAULT_TYPE_RULES = [
        { id: 0, ResultData: 'int', Assignment: '=', LHSData: 'INTEGER', Operator: [], RHSData: '' },
        { id: 1, ResultData: 'int', Assignment: '=', LHSData: 'int', Operator: [], RHSData: '' },
        { id: 2, ResultData: 'int', Assignment: '=', LHSData: 'int', Operator: ['+'], RHSData: 'INTEGER' }
    ];

    const DEFAULT_GRAMMAR_RULES: GrammarRule = {
        VariableRule: 'IDENTIFIER',
        TypeRule: 'TYPE',
        FunctionRule: 'FUNCTION_DEFINITION',
        ParameterRule: 'PARAMETER',
        AssignmentRule: 'ASSIGNMENT',
        OperatorRule: 'OPERATOR',
        TermRule: 'ELEMENT'
    };
    let show_default_rules = false;

    // --- LOGIC ---

    // Scope Rules Logic
    function addScopeRow() {
        scope_rules = [...scope_rules, { id: next_scope_id++, Start: '', End: '' }];
        rules_submitted = false;
    }

    function removeScopeRow(index: number) {
        scope_rules.splice(index, 1);
        scope_rules = scope_rules; // Trigger reactivity
        rules_submitted = false;
    }

    function handleScopeRuleInput() {
        rules_submitted = false;
    }

    // Type Rules Logic
    function addTypeRow() {
        type_rules = [...type_rules, { id: next_type_id++, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }];
        rules_submitted = false;
    }

    function removeTypeRow(index: number) {
        type_rules.splice(index, 1);
        type_rules = type_rules; // Trigger reactivity
        rules_submitted = false;
    }

    function updateTypeOperator(rule: TypeRule, input: string) {
        const trimmedValue = input.trim();
        if (trimmedValue === '') {
            rule.Operator = [''];
        } else {
            rule.Operator = trimmedValue
                .split(',')
                .map(op => op.trim())
                .filter(op => op.length > 0);
        }
        handleTypeRuleInput();
    }

    function handleTypeRuleInput() {
        rules_submitted = false;
    }

    // Grammar Rules Logic (inputs are directly bound, no add/remove)
    function handleGrammarRuleInput() {
        rules_submitted = false;
    }

    // Universal Submission Logic
    function handleSubmit() {
        // Validate Scope Rules
        if (scope_rules.some((rule) => rule.Start.trim() === '' || rule.End.trim() === '')) {
            AddToast('Incomplete scope rules: Please fill in all Start and End fields for scope analysis', 'error');
            return;
        }

        // Validate Type Rules
        if (type_rules.some((rule) =>
            rule.ResultData.trim() === '' ||
            rule.Assignment.trim() === '' ||
            rule.LHSData.trim() === ''
        )) {
            AddToast('Incomplete type rules: Please fill in all Result, Assignment, and LHS fields for type checking', 'error');
            return;
        }

        // Validate Grammar Rules (check if any field is empty)
        const grammarRuleFields = Object.values(grammar_rules);
        if (grammarRuleFields.some(field => typeof field === 'string' && field.trim() === '')) {
            AddToast('Incomplete grammar rules: Please fill in all grammar constraint fields', 'error');
            return;
        }

        // If all validations pass
        submitted_scope_rules = JSON.parse(JSON.stringify(scope_rules));
        submitted_type_rules = JSON.parse(JSON.stringify(type_rules));
        submitted_grammar_rules = JSON.parse(JSON.stringify(grammar_rules));

        rules_submitted = true;
        AddToast('Semantic rules saved successfully! Ready to generate symbol table and perform analysis', 'success');
    }

    // Universal Reset Logic
    function resetState() {
        scope_rules = [{ id: 0, Start: '', End: '' }];
        next_scope_id = 1;
        submitted_scope_rules = [];

        type_rules = [{ id: 0, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }];
        next_type_id = 1;
        submitted_type_rules = [];

        grammar_rules = {
            VariableRule: '',
            TypeRule: '',
            FunctionRule: '',
            ParameterRule: '',
            AssignmentRule: '',
            OperatorRule: '',
            TermRule: ''
        };
        submitted_grammar_rules = { ...grammar_rules };

        rules_submitted = false;
        dispatch('reset');
        AddToast('Rules reset: All semantic analysis rules have been cleared', 'info');
    }

    function insertDefaultRules() {
        scope_rules = DEFAULT_SCOPE_RULES.map(r => ({ ...r }));
        next_scope_id = 1;
        type_rules = DEFAULT_TYPE_RULES.map(r => ({ ...r }));
        next_type_id = 3;
        grammar_rules = { ...DEFAULT_GRAMMAR_RULES };
        show_default_rules = true;
        rules_submitted = false;
    }

    function removeDefaultRules() {
        scope_rules = [{ id: 0, Start: '', End: '' }];
        next_scope_id = 1;
        type_rules = [{ id: 0, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }];
        next_type_id = 1;
        grammar_rules = {
            VariableRule: '',
            TypeRule: '',
            FunctionRule: '',
            ParameterRule: '',
            AssignmentRule: '',
            OperatorRule: '',
            TermRule: ''
        };
        show_default_rules = false;
        rules_submitted = false;
    }

    async function handleGenerate() {
    try {
        const user_id = localStorage.getItem('user_id');
		const project = get(projectName);
		if (!user_id) {
			AddToast('Authentication required: Please log in to generate symbol table', 'error');
			return;
		}
		if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }
        is_loading = true;
        
        const requestData = {
            users_id: user_id,
            scope_rules: submitted_scope_rules,
            grammar_rules: submitted_grammar_rules,
            type_rules: submitted_type_rules,
            project_name: project
        };


        const response = await fetch('http://localhost:8080/api/analysing/analyse', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        });
        
        const result = await response.json();
        
        // Transform the data with proper null checks
        const symbols = result.symbol_table?.SymbolScopes?.map((s: any) => ({
            name: s.Name || s.name || 'unknown',
            type: s.Type || s.type || 'unknown',
            scope: s.Scope || s.scope || 0
        })) || [];

        onGenerateSymbolTable({
				symbol_table: symbols,
                analyser_error: result.error,
                analyser_error_details: result.details
			});

        
        symbol_table =  symbols ;
        show_symbol_table = true;

        if (result.error) {
            AddToast('Semantic error detected! Check the analysis results for details', 'error');
            dispatch('generate',{
                symbol_table: symbols,
                analyser_error: true,
                analyser_error_details: result.details
            });
            console.log(result)
        }else {
            AddToast('Semantic analysis complete! Symbol table generated successfully', 'success');
            dispatch('generate',{
                symbol_table: symbols
            });
        }
    } catch (error) {
         const err = error as { 
            response?: { 
                data?: any; 
                status?: number 
            }; 
            message?: string 
        };
        console.error('Error details:', {
            error: err,
            response: err.response?.data,
            status: err.response?.status
        });
        console.error('Error generating symbol table:', error);
        AddToast('Analysis failed: ' + ('Unable to generate symbol table. Please check your connection'), 'error');
    } finally {
        is_loading = false;
    }
}

    // --- REACTIVE STATEMENTS ---

    // Scope Rules Completeness
    $: lastScopeRowComplete =
        scope_rules.length > 0 &&
        scope_rules[scope_rules.length - 1].Start.trim() !== '' &&
        scope_rules[scope_rules.length - 1].End.trim() !== '';

    $: allScopeRowsComplete = scope_rules.every(
        (rule) => rule.Start.trim() !== '' && rule.End.trim() !== ''
    );

    // Type Rules Completeness
    $: lastTypeRowComplete =
        type_rules.length > 0 &&
        type_rules[type_rules.length - 1].ResultData.trim() !== '' &&
        type_rules[type_rules.length - 1].Assignment.trim() !== '' &&
        type_rules[type_rules.length - 1].LHSData.trim() !== ''

    $: allTypeRowsComplete = type_rules.every(
        (rule) =>
            rule.ResultData.trim() !== '' &&
            rule.Assignment.trim() !== '' &&
            rule.LHSData.trim() !== ''
    );

    // Grammar Rules Completeness
    $: allGrammarRulesComplete = Object.values(grammar_rules).every(
        (field) => typeof field === 'string' && field.trim() !== ''
    );

    // Overall Completeness for Submit Button
    $: overallAllRowsComplete = allScopeRowsComplete && allTypeRowsComplete && allGrammarRulesComplete;

    // Update rules when project data is loaded
    $: if ($lexerState?.analyzer_data && !hasInitialized) {
        const analyzerData = $lexerState.analyzer_data;
        
        // Update scope rules
        scope_rules = analyzerData.scope_rules.map((rule, index) => ({
            id: index,
            Start: rule.start,
            End: rule.end
        }));
        next_scope_id = scope_rules.length;

        // Update type rules
        type_rules = analyzerData.type_rules.map((rule, index) => ({
            id: index,
            ResultData: rule.resultdata,
            Assignment: rule.assignment,
            LHSData: rule.lhsdata,
            Operator: rule.operator || [''],
            RHSData: rule.rhsdata
        }));
        next_type_id = type_rules.length;

        // Update grammar rules
        grammar_rules = {
            VariableRule: analyzerData.grammar_rules.variablerule,
            TypeRule: analyzerData.grammar_rules.typerule,
            FunctionRule: analyzerData.grammar_rules.functionrule,
            ParameterRule: analyzerData.grammar_rules.parameterrule,
            AssignmentRule: analyzerData.grammar_rules.assignmentrule,
            OperatorRule: analyzerData.grammar_rules.operatorrule,
            TermRule: analyzerData.grammar_rules.termrule
        };

        // Set submitted state
        submitted_scope_rules = [...scope_rules];
        submitted_type_rules = [...type_rules];
        submitted_grammar_rules = { ...grammar_rules };
        rules_submitted = true;
        hasInitialized = true;
        
        // Force reactivity by triggering array updates
        scope_rules = [...scope_rules];
        type_rules = [...type_rules];
    }
</script>

<div class="panel-container">
    <div class="header-container">
        <h1 class="analyser-heading">ANALYSING</h1>
    </div>

    <!-- Source Code Section -->
    <div class="source-code-section">
        <h3 class="source-code-header">Source Code</h3>
        <pre class="source-display">{source_code || 'No source code available'}</pre>
    </div>

    <div class="instructions-section">
		<div class="instructions-content">
			<h4 class="instructions-header">Instructions</h4>
			<p class="instructions-text">
				Enter scope rules and type rules as well as link your grammar variables to the applicable programming constructs.
			</p>
		</div>
	</div>

    <div class="scrollable-content">
        <!-- Scope Rules Box -->
        <div class="analyser-box">
            <div class="analyser-box-header">
                <h2 class="heading2">Scope Rules</h2>
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
            <div class="rules-list">
                {#each scope_rules as rule, i (rule.id)}
                    <div class="rule-row" in:slide={{ duration: 250 }}>
                        <div class="input-tooltip-wrapper">
                            <input
                                type="text"
                                bind:value={rule.Start}
                                placeholder="Start Delimiter"
                                aria-label="Start Delimiter"
                                class="scope-input"
                                on:focus={() => i === 0 && (show_start_tooltip = true)}
                                on:blur={() => i === 0 && (show_start_tooltip = false)}
                                on:input={handleScopeRuleInput}
                            />
                            {#if i === 0 && show_start_tooltip}
                                <div class="tooltip">Examples include <code>{'{'}</code></div>
                            {/if}
                        </div>
                        <span class="arrow-separator">→</span>
                        <div class="input-tooltip-wrapper">
                            <input
                                type="text"
                                bind:value={rule.End}
                                placeholder="End Delimiter"
                                aria-label="End Delimiter"
                                class="scope-input"
                                on:focus={() => i === 0 && (show_end_tooltip = true)}
                                on:blur={() => i === 0 && (show_end_tooltip = false)}
                                on:input={handleScopeRuleInput}
                            />
                            {#if i === 0 && show_end_tooltip}
                                <div class="tooltip">Examples include <code>{'}'}</code></div>
                            {/if}
                        </div>
                        <button
                            class="icon-button delete-button"
                            on:click={() => removeScopeRow(i)}
                            aria-label="Remove scope rule row"
                            disabled={scope_rules.length === 1}
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
                        </button>
                    </div>
                {/each}
            </div>
            {#if lastScopeRowComplete}
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
                {#each type_rules as rule, i (rule.id)}
                    <div class="rule-row" in:slide={{ duration: 250 }}>
                        <div class="type-rule-inputs">
                            <input type="text" bind:value={rule.ResultData} placeholder="Result Type" aria-label="Result type" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.Assignment} placeholder="Assignment" aria-label="Assignment" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" bind:value={rule.LHSData} placeholder="LHS" aria-label="LHS" class="scope-input" on:input={handleTypeRuleInput} />
                            <input type="text" value={rule.Operator.join(',')} on:input={(e) =>{if (!e.target) return; updateTypeOperator(rule, (e.target as HTMLInputElement).value)}} placeholder="Operator" aria-label="Operator(s)" class="scope-input"/>
                            <input type="text" bind:value={rule.RHSData} placeholder="RHS" aria-label="RHS" class="scope-input" on:input={handleTypeRuleInput} />
                        </div>
                        <button
                            class="icon-button delete-button"
                            on:click={() => removeTypeRow(i)}
                            aria-label="Remove type rule row"
                            disabled={type_rules.length === 1}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg>
                        </button>
                    </div>
                {/each}
            </div>
            {#if lastTypeRowComplete}
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
                    <input id="grammar-variable" type="text" bind:value={grammar_rules.VariableRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-type" class="grammar-label">Type</label>
                    <input id="grammar-type" type="text" bind:value={grammar_rules.TypeRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-function-declaration" class="grammar-label">Function Declaration</label>
                    <input id="grammar-function-declaration" type="text" bind:value={grammar_rules.FunctionRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-parameter" class="grammar-label">Parameter</label>
                    <input id="grammar-parameter" type="text" bind:value={grammar_rules.ParameterRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-assignment" class="grammar-label">Assignment</label>
                    <input id="grammar-assignment" type="text" bind:value={grammar_rules.AssignmentRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-operator" class="grammar-label">Operator</label>
                    <input id="grammar-operator" type="text" bind:value={grammar_rules.OperatorRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
                <div class="grammar-rule-item">
                    <label for="grammar-term" class="grammar-label">Term</label>
                    <input id="grammar-term" type="text" bind:value={grammar_rules.TermRule} class="scope-input" on:input={handleGrammarRuleInput} />
                </div>
            </div>
        </div>
    </div>
    <div class="actions-container">
        {#if !rules_submitted}
            <button class="submit-button" on:click={handleSubmit} disabled={!overallAllRowsComplete}>
                Submit All Rules
            </button>
        {:else}
            <div class="generate-wrapper" transition:slide|local={{ duration: 250 }}>
                <button 
                    class="generate-button" 
                    on:click={handleGenerate} 
                    disabled={is_loading}
                >
                    {#if is_loading}
                        Generating...
                    {:else}
                        Generate Symbol Table
                    {/if}
                </button>
            </div>
        {/if}
    </div>
</div>

<style>
	/* General Heading Styles (existing) */
	.heading2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0 0 0.25rem 0;
		color: #1a2a4a; /* Adjusted for consistency */
	}
	/* Panel Layout */
	.panel-container {
		display: flex;
		flex-direction: column;
		padding: 1.5rem 1rem;
		height: auto;
	}

	.header {
		position: relative;
		flex-shrink: 0; /* Prevents header from shrinking */
	}

    .instructions-section {
		margin: 1.5rem 0 2rem 0;
        margin-top: 1rem;
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

	/* Scrollable Content Area */
	.scrollable-content {
        flex: 0 0 auto;
        overflow-y: visible;
        padding-right: 0;
        margin-right: 0;
        padding-bottom: 1rem;
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
		flex: 0;
        width: 180px; 
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
		grid-template-columns: repeat(5, auto);
		gap: 0.75rem;
	}

	.type-rule-inputs .scope-input {
		width: 70px;
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
		padding-top: 0rem;
		margin-top: 0rem;
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
		background-color: #BED2E6;
		color: #000000;
	}
	.submit-button:hover:not(:disabled) {
		background-color: #a8bdd1;
        transform: translateY(-2px);
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
	:global(html.dark-mode) .analyser-heading,
	:global(html.dark-mode) .heading2,
	:global(html.dark-mode) .grammar-label,
	:global(html.dark-mode) .source-code-header {
		color: #ebeef1;
	}

    :global(html.dark-mode) .panel-container{
        background-color: #1a2a4a;
    }
	
	:global(html.dark-mode) .analyser-box,
	:global(html.dark-mode) .source-display {
		background-color: #2d3748;
	}

	:global(html.dark-mode) .scope-input {
		background-color: #2d3748;
		border-color: #4b5563;
		color: #f0f0f0;
	}
	:global(html.dark-mode) .scope-input:focus {
		border-color: #60a5fa;
		box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.4);
	}
	:global(html.dark-mode) .arrow-separator {
		color: #9ca3af;
	}
	:global(html.dark-mode) .icon-button {
		color: #9ca3af;
	}
	:global(html.dark-mode) .icon-button:hover:not(:disabled) {
		background-color: #374151;
	}
	:global(html.dark-mode) .delete-button:hover:not(:disabled) {
		color: #ff8a8a;
	}
	:global(html.dark-mode) .tooltip {
		background: #1f2937;
		color: #d1d5db;
		border-color: #4b5563;
	}
	:global(html.dark-mode) .tooltip code {
		background: #374151;
	}
	:global(html.dark-mode) .add-button {
		background-color: transparent;
		color: #60a5fa;
		border-color: #4b5563;
	}
	:global(html.dark-mode) .add-button:hover {
		background-color: #374151;
		border-color: #60a5fa;
	}
	:global(html.dark-mode) .submit-button,
	:global(html.dark-mode) .generate-button {
		background-color: #001A6E;
		color: #ffffff;
	}
    :global(html.dark-mode) .submit-button:hover:not(:disabled) {
        background-color: #002a8e;
    }
	:global(html.dark-mode) .reset-button {
		color: #9ca3af;
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
		max-height: 260px;
		overflow: auto;
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
		justify-content: space-between;
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
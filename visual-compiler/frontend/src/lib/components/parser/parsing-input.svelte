<script lang="ts">
    import { AddToast } from '$lib/stores/toast';
    import { projectName } from '$lib/stores/project';
    import { get } from 'svelte/store';
    import { lexerState } from '$lib/stores/lexer';
    import { onMount, createEventDispatcher, onDestroy } from 'svelte';
    import { parserState, updateParserInputs, markParserSubmitted } from '$lib/stores/parser';
    import { updateParserArtifacts } from '$lib/stores/parser';

    export let source_code = '';

    const dispatch = createEventDispatcher();

    let parseTree = null;
    let hasParseTree = false;

    // FIX: Data structure to use productions consistently
    interface Rule {
        id: number;
        nonTerminal: string;
        productions: string; // Single comma-separated string for multiple productions
    }

    // --- STATE MANAGEMENT ---
    let translation_id_counter = 1;
    let rule_id_counter = 1;

    // FIX: Start with empty array instead of pre-populated rule
    let grammar_rules: Rule[] = [];
    let variables_string = '';
    let terminals_string = '';
    let show_default_grammar = false;
    let is_grammar_submitted = false;
    let tokens = null;
    let tokens_unidentified = null;

    // FIX: Add loading states for buttons
    let isSubmittingGrammar = false;
    let isGeneratingTree = false;

    let user_grammar_rules: Rule[] = [];
    let user_variables_string = '';
    let user_terminals_string = '';
    let user_rule_id_counter = 0;

    // --- DEFAULT GRAMMAR DATA ---
    const DEFAULT_GRAMMAR = {
        variables: 'PROGRAM, STATEMENT, FUNCTION, ITERATION, DECLARATION, ELEMENT, TYPE, EXPRESSION, FUNCTION_DEFINITION, FUNCTION_BLOCK, RETURN, ITERATION_DEFINITION, ITERATION_BLOCK, PARAMETER, PRINT',
        terminals: 'KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER, OPERATOR, SEPARATOR, OPEN_BRACKET, CLOSE_BRACKET, OPEN_SCOPE, CLOSE_SCOPE, CONTROL',
        rules: [
            { nonTerminal: 'PROGRAM', productions: 'STATEMENT FUNCTION STATEMENT ITERATION' },
            { nonTerminal: 'STATEMENT', productions: 'DECLARATION SEPARATOR' },
            { nonTerminal: 'DECLARATION', productions: 'TYPE IDENTIFIER ASSIGNMENT ELEMENT' },
            { nonTerminal: 'DECLARATION', productions: 'IDENTIFIER ASSIGNMENT EXPRESSION'},
            { nonTerminal: 'DECLARATION', productions: 'IDENTIFIER ASSIGNMENT IDENTIFIER PARAMETER' },
            { nonTerminal: 'TYPE', productions: 'KEYWORD' },
            { nonTerminal: 'EXPRESSION', productions: 'ELEMENT OPERATOR ELEMENT' },
            { nonTerminal: 'ELEMENT', productions: 'INTEGER' },
            { nonTerminal: 'ELEMENT', productions: 'IDENTIFIER' },
            { nonTerminal: 'FUNCTION', productions: 'FUNCTION_DEFINITION FUNCTION_BLOCK' },
            { nonTerminal: 'FUNCTION_DEFINITION', productions: 'TYPE IDENTIFIER PARAMETER' },
            { nonTerminal: 'FUNCTION_BLOCK', productions: 'OPEN_SCOPE STATEMENT RETURN CLOSE_SCOPE' },
            { nonTerminal: 'RETURN', productions: 'KEYWORD ELEMENT SEPARATOR' },
            { nonTerminal: 'ITERATION', productions: 'ITERATION_DEFINITION ITERATION_BLOCK' },
            { nonTerminal: 'ITERATION_DEFINITION', productions: 'CONTROL IDENTIFIER CONTROL PARAMETER' },
            { nonTerminal: 'ITERATION_BLOCK', productions: 'OPEN_SCOPE STATEMENT PRINT CLOSE_SCOPE' },
            { nonTerminal: 'PARAMETER', productions: 'OPEN_BRACKET ELEMENT CLOSE_BRACKET' }, 
            { nonTerminal: 'PARAMETER', productions: 'OPEN_BRACKET TYPE IDENTIFIER CLOSE_BRACKET' },
            { nonTerminal: 'PRINT', productions: 'KEYWORD OPEN_BRACKET ELEMENT CLOSE_BRACKET SEPARATOR' }
        ]
    };

    let aiParserEventListener: (event: CustomEvent) => void;

    // FIX: Remove the automatic addNewRule call in onMount
    onMount(async () => {
        // DON'T call addNewRule() here - let store restoration handle it
        await fetchTokens();

        // Listen for AI-generated parser grammar
        aiParserEventListener = (event: CustomEvent) => {
            if (event.detail && event.detail.grammar) {
                console.log('Received AI parser grammar:', event.detail.grammar);
                
                // Save current state as user input before AI insertion
                if (!show_default_grammar) {
                    saveCurrentAsUserInput();
                }
                
                const grammar = event.detail.grammar;
                
                // Clear existing rules first
                grammar_rules = [];
                rule_id_counter = 0;
                translation_id_counter = 0;
                
                // Populate variables and terminals
                variables_string = grammar.variables || '';
                terminals_string = grammar.terminals || '';
                
                // Handle rules
                if (grammar.rules && Array.isArray(grammar.rules)) {
                    grammar_rules = (grammar.rules || []).map((rule) => {
                        rule_id_counter++;

                        let productions_string = (typeof rule.output === 'string') ? rule.output.trim() : '';

                        return {
                            id: rule_id_counter,
                            nonTerminal: rule.input || '',
                            productions: productions_string
                        };
                    });
                }
                
                // If no rules were generated, add a blank one
                if (grammar_rules.length === 0) {
                    addNewRule();
                }
                

                is_grammar_submitted = false;
                
                // Force reactivity update
                grammar_rules = [...grammar_rules];
                variables_string = variables_string;
                terminals_string = terminals_string;
                
                AddToast('AI parser grammar inserted into grammar editor!', 'success');
                
                console.log('Final grammar_rules:', grammar_rules);
                console.log('Final variables_string:', variables_string);
                console.log('Final terminals_string:', terminals_string);
            }
        };

        if (typeof window !== 'undefined') {
            window.addEventListener('ai-parser-generated', aiParserEventListener);
        }
    });

    onDestroy(() => {
        if (aiParserEventListener && typeof window !== 'undefined') {
            window.removeEventListener('ai-parser-generated', aiParserEventListener);
        }
    });

    // Function to save current state as user input
    function saveCurrentAsUserInput() {
        user_grammar_rules = JSON.parse(JSON.stringify(grammar_rules));
        user_variables_string = variables_string;
        user_terminals_string = terminals_string;
        user_rule_id_counter = rule_id_counter;
    }

    // handleGrammarChange
    let hasInitialized = false;
    let currentProjectName = '';

    // FIX: Simplify project change detection
    $: if ($projectName !== currentProjectName) {
        console.log('Parser: Project changed from', currentProjectName, 'to', $projectName);
        hasInitialized = false;
        currentProjectName = $projectName;
    }

    // FIX: Initialize from store when project name is set and store has data
    $: if ($parserState && $projectName && (!hasInitialized || $projectName !== currentProjectName)) {
        console.log('Parser component initializing with project:', $projectName, 'and state:', $parserState);
        
        // Convert from store format to component format
        if ($parserState.grammar_rules && $parserState.grammar_rules.length > 0) {
            // Check if we have meaningful rules (not just empty default)
            const meaningfulRules = $parserState.grammar_rules.filter(rule => 
                rule.nonTerminal.trim() !== '' || 
                (rule.translations && rule.translations.some(t => t.value.trim() !== ''))
            );
            
            if (meaningfulRules.length > 0) {
                grammar_rules = meaningfulRules.map(rule => ({
                    id: rule.id,
                    nonTerminal: rule.nonTerminal,
                    productions: rule.translations
                        .filter(t => t.value.trim() !== '')
                        .map(t => t.value)
                        .join(', ')
                }));
                rule_id_counter = Math.max(...meaningfulRules.map(r => r.id)) + 1;
            } else {
                // No meaningful rules, start with empty
                grammar_rules = [{ id: 1, nonTerminal: '', productions: '' }];
                rule_id_counter = 1;
            }
        } else {
            grammar_rules = [{ id: 1, nonTerminal: '', productions: '' }];
            rule_id_counter = 1;
        }
        variables_string = $parserState.variables_string || '';
        terminals_string = $parserState.terminals_string || '';
        translation_id_counter = $parserState.translation_id_counter || 1;
        show_default_grammar = $parserState.show_default_grammar || false;
        is_grammar_submitted = $parserState.is_grammar_submitted || false;
        
        // Handle parse tree artifacts
        if ($parserState.hasParseTree && $parserState.parseTree) {
            parseTree = $parserState.parseTree;
            hasParseTree = true;
            dispatch('treereceived', parseTree);
        } else {
            parseTree = null;
            hasParseTree = false;
        }
        
        hasInitialized = true;
        console.log('Parser component initialized with:', { 
            grammar_rules: grammar_rules.length, 
            hasParseTree,
            is_grammar_submitted,
            meaningfulRules: grammar_rules.filter(r => r.nonTerminal.trim() !== '' || r.productions.trim() !== '').length
        });
    }

    // FIX: Only update store if there are meaningful changes
    function handleGrammarChange() {
        is_grammar_submitted = false;
        
        if (!show_default_grammar) {
            saveCurrentAsUserInput();
        }

        if (hasInitialized) {
            // FIX: Filter out completely empty rules before storing
            const meaningful_rules = grammar_rules.filter(rule => 
                rule.nonTerminal.trim() !== '' || rule.productions.trim() !== ''
            );
            
            const storeFormatRules = meaningful_rules.map(rule => ({
                id: rule.id,
                nonTerminal: rule.nonTerminal,
                translations: rule.productions.split(',').map((prod, index) => ({
                    id: index + 1,
                    value: prod.trim()
                })).filter(t => t.value !== '')
            }));

            // Only add empty rule if there are no meaningful rules
            if (storeFormatRules.length === 0) {
                storeFormatRules.push({
                    id: 1,
                    nonTerminal: '',
                    translations: [{ id: 1, value: '' }]
                });
            }

            updateParserInputs({
                grammar_rules: storeFormatRules,
                variables_string,
                terminals_string,
                rule_id_counter,
                translation_id_counter,
                show_default_grammar,
                is_grammar_submitted: false
            });
        }
    }

    // FIX: Only add new rule if the last rule has content
    function addNewRule() {
        // Check if we should add a rule
        const last_rule = grammar_rules[grammar_rules.length - 1];
        if (last_rule && last_rule.nonTerminal.trim() === '' && last_rule.productions.trim() === '') {
            // Don't add if the last rule is already empty
            return;
        }

        rule_id_counter++; 
        translation_id_counter++;
        grammar_rules = [...grammar_rules, {
            id: rule_id_counter,
            nonTerminal: '',
            productions: ''
        }];
        handleGrammarChange();
    }

    // insertDefaultGrammar
    function insertDefaultGrammar() {
        handleGrammarChange();
        
        // Save current user input before switching to default
        saveCurrentAsUserInput();
        
        // Set default values
        show_default_grammar = true;
        variables_string = DEFAULT_GRAMMAR.variables;
        terminals_string = DEFAULT_GRAMMAR.terminals;

        let new_rule_id = 0;

        grammar_rules = DEFAULT_GRAMMAR.rules.map((r) => {
            new_rule_id++;
            return {
                id: new_rule_id,
                nonTerminal: r.nonTerminal,
                productions: r.productions
            };
        });
        rule_id_counter = new_rule_id;
    }

    // removeDefaultGrammar
    function removeDefaultGrammar() {
        handleGrammarChange();
        show_default_grammar = false;
        
        grammar_rules = JSON.parse(JSON.stringify(user_grammar_rules));
        variables_string = user_variables_string;
        terminals_string = user_terminals_string;
        rule_id_counter = user_rule_id_counter;
        
        // FIX: Only add new rule if grammar_rules is completely empty
        if (grammar_rules.length === 0) {
            grammar_rules = [{ id: 1, nonTerminal: '', productions: '' }];
            rule_id_counter = 1;
        }
    }

    function clearAllInputs() {
        // Reset counters first
        rule_id_counter = 1;
        translation_id_counter = 1;
        
        // FIX: Create rule with productions
        grammar_rules = [{
            id: 1,
            nonTerminal: '',
            productions: ''
        }];
        
        // Reset variables and terminals
        variables_string = '';
        terminals_string = '';
        
        // Reset states
        show_default_grammar = false;
        is_grammar_submitted = false;
        
        // Update store if using it
        if (hasInitialized) {
            const storeFormatRules = [{
                id: 1,
                nonTerminal: '',
                translations: [{ id: 1, value: '' }]
            }];

            updateParserInputs({
                grammar_rules: storeFormatRules,
                variables_string: '',
                terminals_string: '',
                rule_id_counter: 1,
                translation_id_counter: 1,
                show_default_grammar: false,
                is_grammar_submitted: false
            });
        }
        
        AddToast('All grammar inputs cleared successfully!', 'success');
    }

    // handleSubmitGrammar
    async function handleSubmitGrammar() {
        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName);

        if (!accessToken) {
            AddToast('Authentication required: Please log in to save grammar rules', 'error');
            return;
        }

        if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }

        // FIX: Filter out rules where the non-terminal is empty
        const cleaned_rules = grammar_rules.filter((rule) => rule.nonTerminal.trim() !== '');
        grammar_rules = cleaned_rules;

        if (cleaned_rules.length === 0) {
            AddToast('Empty grammar: Please define at least one production rule to continue', 'error');
            addNewRule();
            return;
        }

        const variable_list = variables_string
            .split(',')
            .map((v) => v.trim())
            .filter((v) => v);
        const terminal_list = terminals_string
            .split(',')
            .map((t) => t.trim())
            .filter((t) => t);
        const defined_symbols = new Set([...variable_list, ...terminal_list]);
        const start_variable = cleaned_rules[0]?.nonTerminal.trim() || '';

        if (!variable_list.includes(start_variable)) {
            AddToast(
                `The start symbol '${start_variable}' must be included in the Variables list.`,
                'error'
            );
            return;
        }

        for (const rule of cleaned_rules) {
            const non_terminal = rule.nonTerminal.trim();
            // This check is slightly redundant due to the filter, but good for safety.
            if (!non_terminal) {
                AddToast('Missing non-terminal: Every production rule must have a left-hand side symbol', 'error');
                return;
            }
            if (!variable_list.includes(non_terminal)) {
                AddToast(
                    `The symbol '${non_terminal}' used in a rule must be defined in the Variables list.`,
                    'error'
                );
                return;
            }

            const has_at_least_one_production = rule.productions.trim() !== '';
            if (!has_at_least_one_production) {
                AddToast(`Empty production: Rule for '${non_terminal}' needs at least one production on the right-hand side`, 'error');
                return;
            }

            // Parse productions: split by comma first, then by spaces
            const production_alternatives = rule.productions.split(',').map(p => p.trim()).filter(p => p);
            const all_symbols = production_alternatives.flatMap(prod => prod.split(' ').map(s => s.trim()).filter(s => s));
            
            for (const symbol of all_symbols) {
                if (symbol && !defined_symbols.has(symbol)) {
                    AddToast(
                        `Invalid symbol '${symbol}' in rule for '${non_terminal}'. It must be defined as a Variable or Terminal.`,
                        'error'
                    );
                    return;
                }
            }
        }

        const formatted_rules = cleaned_rules
            .flatMap((rule) => {
                // Split by comma to get separate production alternatives
                const productions = rule.productions.split(',').map(p => p.trim()).filter(p => p !== '');
                
                // If no productions, create one with empty string (epsilon production)
                if (productions.length === 0) {
                    return [{
                        input: rule.nonTerminal.trim(),
                        output: [""]
                    }];
                }
                
                // Create a separate rule for each production alternative
                return productions.map(production => {
                    // Handle empty production (epsilon)
                    if (production.trim() === '') {
                        return {
                            input: rule.nonTerminal.trim(),
                            output: [""]
                        };
                    }
                    
                    return {
                        input: rule.nonTerminal.trim(),
                        output: production.split(' ').map(s => s.trim()).filter(s => s)
                    };
                });
            })
            .filter((rule) => rule.input && rule.output.length > 0);

        const final_json_output = {
            project_name: project, 
            variables: variable_list,
            terminals: terminal_list,
            start: start_variable,
            rules: formatted_rules
        };

        console.log('Sending grammar data to backend:', JSON.stringify(final_json_output, null, 2));

        // FIX: Set loading state
        isSubmittingGrammar = true;

        try {
            // FIX: Add 1-second delay before API call
            await new Promise(resolve => setTimeout(resolve, 1000));

            const response = await fetch('http://localhost:8080/api/parsing/grammar', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}` 
                },
                body: JSON.stringify(final_json_output)
            });

            if (!response.ok) {
                const error_data = await response.json();
                throw new Error(error_data.error || error_data.details || 'Failed to submit grammar');
            }

            const result = await response.json();
            AddToast('Grammar saved successfully! Your parsing rules are ready for syntax analysis', 'success');
            
            const submittedRules = grammar_rules.map(rule => ({
                id: rule.id,
                nonTerminal: rule.nonTerminal,
                translations: rule.productions.split(',').map((prod, index) => ({
                    id: index + 1,
                    value: prod.trim()
                })).filter(t => t.value !== '')
            }));

            updateParserInputs({
                grammar_rules: submittedRules,
                variables_string,
                terminals_string,
                rule_id_counter,
                translation_id_counter,
                show_default_grammar,
                is_grammar_submitted: true // Mark as submitted
            });
            // Mark as submitted in store
            markParserSubmitted();
            is_grammar_submitted = true;
            
        } catch (error) {
            console.error('Submit Grammar Error:', error);
            AddToast('Grammar save failed: ' + String(error), 'error');
            is_grammar_submitted = false;
        } finally {
            isSubmittingGrammar = false;
        }
    }

    // fetchTokens
    async function fetchTokens() {
        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName);

        if (!accessToken || !project) return null;

        try {
            const response = await fetch('http://localhost:8080/api/lexing/lexer', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({
                    project_name: project
                })
            });

            if (!response.ok) {
                throw new Error('Failed to fetch tokens');
            }

            const data = await response.json();
            tokens = data.tokens;
            tokens_unidentified = data.tokens_unidentified;
            return data.tokens;
        } catch (error) {
            console.error('Error fetching tokens:', error);
            return null;
        }
    }

    // generateSyntaxTree
    async function generateSyntaxTree() {
        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName);

        if (!accessToken || !project) {
            AddToast('Authentication error: Missing user credentials or project information', 'error');
            return;
        }

        // FIX: Set loading state
        isGeneratingTree = true;

        try {
            // FIX: Add 1-second delay before processing
            await new Promise(resolve => setTimeout(resolve, 1000));

            // First verify tokens exist
            console.log('Verifying tokens exist...');
            const tokensResponse = await fetch(`http://localhost:8080/api/lexing/lexer`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}` 
                },
                body: JSON.stringify({
                    project_name: project 
                })
            });

            if (!tokensResponse.ok) {
                console.log('Tokens not found, trying DFA tokenization...');
                const dfa_token_response = await fetch(`http://localhost:8080/api/lexing/dfaToTokens`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${accessToken}` 
                    },
                    body: JSON.stringify({
                        project_name: project 
                    })
                });

                if (!dfa_token_response.ok) {
                    AddToast('Tokens required: Please complete lexical analysis and generate tokens first', 'error');
                    dispatch('parsingerror', {parsing_error: true, parsing_error_details: 'Tokens required: Please complete lexical analysis and generate tokens first'});
                    return;
                }
            }

            // Now try to generate syntax tree
            const response = await fetch('http://localhost:8080/api/parsing/tree', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}` 
                },
                body: JSON.stringify({
                    project_name: project 
                })
            });

            const data = await response.json();
            console.log('Syntax tree response:', data); 

            if (!response.ok) {
                if (data.error) {
                    throw new Error(data.details || data.error);
                }
                throw new Error('Failed to generate syntax tree');
            }

            if (data.tree) {
                updateParserArtifacts(data.tree);
                AddToast('Parse tree generated successfully! Your syntax analysis is complete', 'success');
                dispatch('treereceived', data.tree);
            } else {
                dispatch('parsingerror', {parsing_error: true, parsing_error_details: 'No tree data in response'});
                throw new Error('No tree data in response');
            }
        } catch (error) {
            dispatch('parsingerror', {parsing_error: true, parsing_error_details: error.message});
            console.error('Full error details:', error); 
            AddToast(
                `Failed to generate syntax tree: ${error.message}. Please ensure tokens and grammar are valid.`,
                'error'
            );
        } finally {
            isGeneratingTree = false;
        }
    }

    
</script>

<div class="phase-inspector">
    <div class="parser-heading">
        <h1 class="parser-heading-h1">PARSING</h1>
    </div>

    <div class="source-code-section">
        <h3 class="source-code-header">Source Code</h3>
        <pre class="source-display">{source_code || 'No source code available'}</pre>
    </div>

    <div class="instructions-section">
		<div class="instructions-content">
			<h4 class="instructions-header">Instructions</h4>
			<p class="instructions-text">
				Enter a grammar to parse your tokens and link the terminals of your grammar to your token types.
			</p>
		</div>
	</div>

    <div class="grammar-editor">
        <div class="grammar-header">
            <h3>Context-Free Grammar</h3>

            <div class="button-group">
                <!-- Show example button first (to the left) -->
                <button
                class="option-btn example-btn"
                class:selected={show_default_grammar}
                on:click={show_default_grammar ? removeDefaultGrammar : insertDefaultGrammar}
                type="button"
                aria-label={show_default_grammar ? 'Restore your input' : 'Show context-free grammar example'}
                title={show_default_grammar ? 'Restore your input' : 'Show context-free grammar example'}
            >
                {show_default_grammar ? 'Restore Input' : 'Show Example'}
            </button>
                
                <!-- Clear button second (to the right) -->
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

        <div class="top-inputs">
            <div class="input-group">
                <label for="variables">Variables</label>
                <input
                    id="variables"
                    type="text"
                    placeholder="S, A, B"
                    bind:value={variables_string}
                    on:input={handleGrammarChange}
                />
            </div>
            <div class="input-group">
                <label for="terminals">Terminals</label>
                <input
                    id="terminals"
                    type="text"
                    placeholder="x, y, z"
                    bind:value={terminals_string}
                    on:input={handleGrammarChange}
                />
            </div>
        </div>

        <div class="start-label-container">
            <span class="start-label">Start</span>
        </div>

        <div class="rules-container">
            {#each grammar_rules as rule, i (rule.id)}
                <div class="rule-row">
                    <div class="rule-inputs">
                        <input
                            type="text"
                            class="non-terminal-input"
                            placeholder="VARIABLE"
                            bind:value={rule.nonTerminal}
                            on:input={handleGrammarChange}
                        />
                        <span class="arrow">→</span>
                        <div class="productions-container">
                            <textarea
                                class="productions-textarea"
                                placeholder="VARIABLE TERMINAL VARIABLE TERMINAL..."
                                bind:value={rule.productions}
                                on:input={handleGrammarChange}
                                rows="2"
                            ></textarea>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
        <!-- FIX: Only show Add Rule button when the last rule has content -->
        {#if grammar_rules.length === 0 || (grammar_rules[grammar_rules.length - 1]?.nonTerminal.trim() !== '' || grammar_rules[grammar_rules.length - 1]?.productions.trim() !== '')}
            <button class="add-rule-btn" on:click={addNewRule}>+ Add New Rule</button>
        {/if}
    </div>

    <div class="button-container">
        <!-- FIX: Add loading state to Submit Grammar button -->
        <button 
            class="submit-button" 
            on:click={handleSubmitGrammar}
            disabled={isSubmittingGrammar}
        >
            <div class="button-content">
                {#if isSubmittingGrammar}
                    <div class="loading-spinner"></div>
                    Submitting...
                {:else}
                    Submit Grammar
                {/if}
            </div>
        </button>
        
        <!-- FIX: Add loading state to Generate Syntax Tree button -->
        <button 
            class="submit-button generate-button" 
            class:disabled={!is_grammar_submitted || isGeneratingTree}
            disabled={!is_grammar_submitted || isGeneratingTree}
            on:click={generateSyntaxTree}
            title={is_grammar_submitted ? 
                (isGeneratingTree ? "Generating parse tree..." : "Generate syntax tree from submitted grammar") : 
                "Submit grammar first"}
        >
            <div class="button-content">
                {#if isGeneratingTree}
                    <div class="loading-spinner"></div>
                    Generating...
                {:else}
                    Generate Syntax Tree
                {/if}
            </div>
        </button>
    </div>
</div>

<style>
    .phase-inspector {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 0rem;
    }
    .parser-heading-h1 {
        color: black;
        text-align: center;
        margin-top: 0;
        font-family: 'Times New Roman';
        margin-bottom: 0;
    }
    .source-code-section {
        margin-top: 0rem;
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
    
    .instructions-section {
		margin: 1.5rem 0 2rem 0;
        margin-top: 2rem;
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
    
    .grammar-editor {
        background: #f8f9fa;
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }
    .grammar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }
    .grammar-editor h3 {
        margin: 0;
        color: #001a6e;
        font-family: 'Times New Roman';
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
		background: #BED2E6;
		color: black;
	}

	.example-btn:hover {
		background: #a8bdd1;
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
		transform: translateY(-2px);
	}

    :global(html.dark-mode) .example-btn {
        background: #001A6E;
        color: #ffffff;
    }

    :global(html.dark-mode) .example-btn:hover {
        background: #002a8e;
        box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
    }

	.option-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
	}

	

    .top-inputs {
        display: flex;
        gap: 1.5rem;
        margin-bottom: 1.5rem;
    }
    .input-group {
        flex: 1;
        display: flex;
        flex-direction: column;
    }
    .input-group label {
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #333;
        font-size: 0.9rem;
    }
    .input-group input {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-family: monospace;
        transition: border-color 0.2s, box-shadow 0.2s;
    }
    .input-group input:focus {
        outline: none;
        border-color: #007bff; 
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
    }


    .rules-container {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .start-label-container {
        margin-bottom: 0.5rem;
        padding-left: 0.25rem;
    }

    .rule-row {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    
    .start-label {
        font-weight: bold;
        color: #001a6e;
        font-family: monospace;
        font-size: 1rem;
    }
    
    .rule-inputs {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        flex-grow: 1;
        overflow: hidden;
        width: 100%;
    }
    .non-terminal-input {
        flex: 0 0 120px;
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-family: monospace;
        text-align: center;
        width: 2.8rem;
        transition: border-color 0.2s, box-shadow 0.2s;
    }
    .non-terminal-input:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
    }
    .arrow {
        font-size: 1.2rem;
        color: #555;
        font-weight: bold;
    }
    .productions-container {
        display: flex;
        align-items: flex-start;
        flex: 1;
        overflow: hidden;
    }
    .productions-textarea {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-family: monospace;
        width: 100%;
        flex: 1;
        height: 38px;
        resize: none;
        overflow-y: hidden;
        overflow-x: auto;
        word-wrap: break-word;
        white-space: nowrap;
        transition: border-color 0.2s, box-shadow 0.2s;
        line-height: 1.4;
    }
    .productions-textarea:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
    }

    .add-rule-btn {
        display: block;
        margin: 1.5rem auto 0;
        padding: 0.5rem 1rem;
        border: 1px dashed #001a6e;
        background: transparent;
        color: #001a6e;
        border-radius: 4px;
        cursor: pointer;
        font-weight: 500;
    }
    .button-container {
        display: flex;
        justify-content: center;
        gap: 1rem;
        margin-top: 1rem;
    }
    .submit-button,
    .generate-button {
        padding: 0.6rem 1.5rem;
        background: #BED2E6;
        color: #000000;
        border: none;
        border-radius: 6px;
        font-size: 0.9rem;
        font-weight: 500;
        cursor: pointer;
        transition: background-color 0.2s ease, transform 0.2s ease, opacity 0.2s ease;
        margin-top: 1rem;
        position: relative;
        overflow: hidden;
    }
    .button-content {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .loading-spinner {
        width: 16px;
        height: 16px;
        border: 2px solid transparent;
        border-top: 2px solid currentColor;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

    .submit-button:hover:not(:disabled),
    .generate-button:hover:not(:disabled) {
        background: #a8bdd1;
        box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
        transform: translateY(-2px);
    }

    .submit-button:disabled,
    .generate-button:disabled,
    .generate-button.disabled {
        background: #d6d8db;
        color: #6c757d;
        cursor: not-allowed;
        opacity: 0.6;
        transform: none;
        pointer-events: none;
    }

    .submit-button:disabled:hover,
    .generate-button:disabled:hover,
    .generate-button.disabled:hover {
        background: #d6d8db;
        color: #6c757d;
        transform: none;
    }
        .submit-button:disabled {
        cursor: wait;
        opacity: 0.8;
    }

    .generate-button:disabled {
        cursor: wait;
        opacity: 0.8;
    }

    .submit-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

    .generate-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }


    /* Dark Mode Styles */
    :global(html.dark-mode) .parser-heading-h1,
    :global(html.dark-mode) .source-code-header,
    :global(html.dark-mode) .grammar-editor h3,
    :global(html.dark-mode) .start-label,
    :global(html.dark-mode) .input-group label {
        color: #ebeef1;
    }
    :global(html.dark-mode) .source-display,
    :global(html.dark-mode) .grammar-editor {
        background: #2d3748;
    }
    :global(html.dark-mode) .source-display {
        color: #e5e7eb;
    }

    :global(html.dark-mode) .input-group input,
    :global(html.dark-mode) .non-terminal-input,
    :global(html.dark-mode) .productions-textarea {
        background-color: #2d3748;
        border-color: #4b5563;
        color: #f0f0f0;
    }
    :global(html.dark-mode) .input-group input:focus,
    :global(html.dark-mode) .non-terminal-input:focus,
    :global(html.dark-mode) .productions-textarea:focus {
        border-color: #60a5fa;
        box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.4); 
    }

    :global(html.dark-mode) .arrow {
        color: #9ca3af;
    }

    :global(html.dark-mode) .add-rule-btn {
        border-color: #60a5fa;
        color: #60a5fa;
    }
    :global(html.dark-mode) .submit-button,
    :global(html.dark-mode) .generate-button {
        background-color: #001A6E;
        color: #ffffff;
    }
    :global(html.dark-mode) .submit-button:hover:not(:disabled),
    :global(html.dark-mode) .generate-button:hover:not(:disabled) {
        background-color: #002a8e;
        box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
        transform: translateY(-2px);
    }

    
    :global(html.dark-mode) .submit-button:disabled,
    :global(html.dark-mode) .generate-button:disabled,
    :global(html.dark-mode) .generate-button.disabled {
        background: #495057;
        color: #6c757d;
        cursor: wait;
        opacity: 0.8;
        transform: none;
    }

    :global(html.dark-mode) .submit-button:disabled:hover,
    :global(html.dark-mode) .generate-button:disabled:hover,
    :global(html.dark-mode) .generate-button.disabled:hover {
        background: #495057;
        color: #6c757d;
        transform: none;
    }

        :global(html.dark-mode) .submit-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

    :global(html.dark-mode) .generate-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

    :global(html.dark-mode) .submit-button .loading-spinner {
        border-top-color: #ffffff;
    }

    :global(html.dark-mode) .generate-button .loading-spinner {
        border-top-color: #ffffff;
    }

    :global(html.dark-mode) .default-toggle-btn {
		background-color: #2d3748;
		border-color: #4a5568;
		color: #d1d5db;
	}
	:global(html.dark-mode) .default-toggle-btn.selected {
		background-color: #001a6e;
		border-color: #60a5fa;
		color: #e0e7ff;
	}
    :global(html.dark-mode) .default-toggle-btn:not(.selected):hover {
        background-color: #001a6e;
        border-color: #60a5fa;
    }

    :global(html.dark-mode) .instructions-section {
		background: #2d3748;
		border-left-color: #4da9ff;
	}

	:global(html.dark-mode) .instructions-header {
		color: #e2e8f0;
	}

	:global(html.dark-mode) .instructions-text {
		color: #cbd5e0;
	}

    .grammar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }

    .button-group {
        display: flex;
        gap: 0.5rem;
        align-items: center;
    }

    .clear-toggle-btn {
        background: white;
        border: 2px solid #e5e7eb;
        color: #ef4444;
        font-size: 1.2rem;
        cursor: pointer;
        transition: all 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 2.2rem;
        width: 2.2rem;
        border-radius: 50%;
    }

    .clear-toggle-btn:hover,
    .clear-toggle-btn:focus {
        background: #fff5f5;
        border-color: #ef4444;
    }

    /* Dark mode for clear button */
    :global(html.dark-mode) .clear-toggle-btn {
        background: transparent;
        border-color: #4a5568;
        color: #ef4444;
    }

    :global(html.dark-mode) .clear-toggle-btn:hover,
    :global(html.dark-mode) .clear-toggle-btn:focus {
        background: rgba(45, 55, 72, 0.5);
        border-color: #ef4444;
    }

    ::-webkit-scrollbar {
        width: 11px;
        height: 7px;
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
</style>
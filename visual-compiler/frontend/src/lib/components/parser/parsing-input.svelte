<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { AddToast } from '$lib/stores/toast';
    import { projectName } from '$lib/stores/project';
    import { get } from 'svelte/store';

    export let source_code = '';

    const dispatch = createEventDispatcher();

    // --- DATA STRUCTURE for CFG ---
    interface Translation {
        id: number;
        value: string;
    }

    interface Rule {
        id: number;
        nonTerminal: string;
        translations: Translation[];
    }

    // --- STATE MANAGEMENT ---
    let grammar_rules: Rule[] = [];
    let rule_id_counter = 0;
    let translation_id_counter = 0;
    let variables_string = '';
    let terminals_string = '';
    let show_default_grammar = false;
    let is_grammar_submitted = false;
    let tokens = null;
    let tokens_unidentified = null;

    // --- DEFAULT GRAMMAR DATA ---
    const DEFAULT_GRAMMAR = {
        variables: 'PROGRAM, STATEMENT, FUNCTION, ITERATION, DECLARATION, ELEMENT, TYPE, EXPRESSION, FUNCTION_DEFINITION, FUNCTION_BLOCK, RETURN, ITERATION_DEFINITION, ITERATION_BLOCK, PARAMETER, PRINT',
        terminals: 'KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER, OPERATOR, DELIMITER, OPEN_BRACKET, CLOSE_BRACKET, OPEN_SCOPE, CLOSE_SCOPE, CONTROL',
        rules: [
            { nonTerminal: 'PROGRAM', translations: ['STATEMENT', 'FUNCTION', 'STATEMENT', 'ITERATION'] },
            { nonTerminal: 'STATEMENT', translations: ['DECLARATION', 'DELIMITER'] },
            { nonTerminal: 'DECLARATION', translations: ['TYPE', 'IDENTIFIER', 'ASSIGNMENT', 'ELEMENT'] },
            { nonTerminal: 'DECLARATION', translations: ['IDENTIFIER', 'ASSIGNMENT', 'EXPRESSION'] },
            { nonTerminal: 'DECLARATION', translations: ['IDENTIFIER', 'ASSIGNMENT', 'IDENTIFIER', 'PARAMETER'] },
            { nonTerminal: 'TYPE', translations: ['KEYWORD'] },
            { nonTerminal: 'EXPRESSION', translations: ['ELEMENT', 'OPERATOR', 'ELEMENT'] },
            { nonTerminal: 'ELEMENT', translations: ['INTEGER'] },
            { nonTerminal: 'ELEMENT', translations: ['IDENTIFIER'] },
            { nonTerminal: 'FUNCTION', translations: ['FUNCTION_DEFINITION', 'FUNCTION_BLOCK'] },
            { nonTerminal: 'FUNCTION_DEFINITION', translations: ['TYPE', 'IDENTIFIER', 'PARAMETER'] },
            { nonTerminal: 'FUNCTION_BLOCK', translations: ['OPEN_SCOPE', 'STATEMENT', 'RETURN', 'CLOSE_SCOPE'] },
            { nonTerminal: 'RETURN', translations: ['KEYWORD', 'ELEMENT', 'DELIMITER'] },
            { nonTerminal: 'ITERATION', translations: ['ITERATION_DEFINITION', 'ITERATION_BLOCK'] },
            { nonTerminal: 'ITERATION_DEFINITION', translations: ['CONTROL', 'IDENTIFIER', 'CONTROL', 'PARAMETER'] },
            { nonTerminal: 'ITERATION_BLOCK', translations: ['OPEN_SCOPE', 'STATEMENT', 'PRINT', 'CLOSE_SCOPE'] },
            { nonTerminal: 'PARAMETER', translations: ['OPEN_BRACKET', 'ELEMENT', 'CLOSE_BRACKET'] },
            { nonTerminal: 'PARAMETER', translations: ['OPEN_BRACKET', 'TYPE', 'IDENTIFIER', 'CLOSE_BRACKET'] },
            { nonTerminal: 'PRINT', translations: ['KEYWORD', 'OPEN_BRACKET', 'ELEMENT', 'CLOSE_BRACKET', 'DELIMITER'] }
        ]
    };

    // onMount
    onMount(async () => {
        addNewRule();
        await fetchTokens(); // Try to fetch existing tokens
    });

    // handleGrammarChange
    function handleGrammarChange() {
        is_grammar_submitted = false;
    }

    // addNewRule
    function addNewRule() {
        handleGrammarChange();
        rule_id_counter++;
        translation_id_counter++;
        grammar_rules = [
            ...grammar_rules,
            {
                id: rule_id_counter,
                nonTerminal: '',
                translations: [{ id: translation_id_counter, value: '' }]
            }
        ];
    }

    // addTranslation
    function addTranslation(rule_id: number) {
        handleGrammarChange();
        translation_id_counter++;
        grammar_rules = grammar_rules.map((rule) => {
            if (rule.id === rule_id) {
                return {
                    ...rule,
                    translations: [...rule.translations, { id: translation_id_counter, value: '' }]
                };
            }
            return rule;
        });
    }

    // insertDefaultGrammar
    function insertDefaultGrammar() {
        handleGrammarChange();
        show_default_grammar = true;
        variables_string = DEFAULT_GRAMMAR.variables;
        terminals_string = DEFAULT_GRAMMAR.terminals;

        let new_rule_id = 0;
        let new_translation_id = 0;

        grammar_rules = DEFAULT_GRAMMAR.rules.map((r) => {
            new_rule_id++;
            return {
                id: new_rule_id,
                nonTerminal: r.nonTerminal,
                translations: r.translations.map((t) => {
                    new_translation_id++;
                    return {
                        id: new_translation_id,
                        value: t
                    };
                })
            };
        });
        rule_id_counter = new_rule_id;
        translation_id_counter = new_translation_id;
    }

    // removeDefaultGrammar
    function removeDefaultGrammar() {
        handleGrammarChange();
        show_default_grammar = false;
        variables_string = '';
        terminals_string = '';
        grammar_rules = [];
        rule_id_counter = 0;
        translation_id_counter = 0;
        addNewRule();
    }

    // handleSubmitGrammar
    async function handleSubmitGrammar() {
        const user_id = localStorage.getItem('user_id');
        const project = get(projectName);

        if (!user_id) {
            AddToast('Authentication required: Please log in to save grammar rules', 'error');
            return;
        }

        if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }

        // FIX: Filter out rules where the non-terminal is empty. This is the main check for a valid rule.
        const cleaned_rules = grammar_rules.filter((rule) => rule.nonTerminal.trim() !== '');

        // Update the UI to visually remove the empty rows, providing immediate feedback.
        grammar_rules = cleaned_rules;

        if (cleaned_rules.length === 0) {
            // If all rules were cleared, show an error and add a new blank rule for convenience.
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

            const has_at_least_one_production = rule.translations.some((t) => t.value.trim() !== '');
            if (!has_at_least_one_production) {
                AddToast(`Empty production: Rule for '${non_terminal}' needs at least one production on the right-hand side`, 'error');
                return;
            }

            const output_array = rule.translations
                .flatMap((t) => t.value.trim().split(' '))
                .filter((v) => v);
            for (const symbol of output_array) {
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
            .map((rule) => ({
                input: rule.nonTerminal.trim(),
                output: rule.translations
                    .flatMap((t) => t.value.trim().split(' '))
                    .filter((v) => v)
            }))
            .filter((rule) => rule.input && rule.output.length > 0);

        const final_json_output = {
            users_id: user_id,
            project_name: project,
            variables: variable_list,
            terminals: terminal_list,
            start: start_variable,
            rules: formatted_rules
        };

        try {
            const response = await fetch('http://localhost:8080/api/parsing/grammar', {
                method: 'POST',
                headers: { 
                    'accept': 'application/json',
                    'Content-Type': 'application/json' 
                },
                body: JSON.stringify(final_json_output)
            });

            if (!response.ok) {
                const error_data = await response.json();
                throw new Error(error_data.error || 'Failed to submit grammar');
            }

            const result = await response.json();
            AddToast('Grammar saved successfully! Your parsing rules are ready for syntax analysis', 'success');
            is_grammar_submitted = true;
        } catch (error) {
            console.error('Submit Grammar Error:', error);
            AddToast('Grammar save failed: ' + String(error), 'error');
            is_grammar_submitted = false;
        }
    }

    // fetchTokens
    async function fetchTokens() {
        const user_id = localStorage.getItem('user_id');
        const project = get(projectName);

        if (!user_id || !project) return null;

        try {
            const response = await fetch('http://localhost:8080/api/lexing/lexer', {
                method: 'POST',
                headers: {
                    'accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    users_id: user_id,
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
        const user_id = localStorage.getItem('user_id');
        const project = get(projectName);

        if (!user_id || !project) {
            AddToast('Authentication error: Missing user credentials or project information', 'error');
            return;
        }

        try {
            // First verify tokens exist
            const tokensResponse = await fetch(`http://localhost:8080/api/lexing/lexer`, {
                method: 'POST',
                headers: {
                    'accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    users_id: user_id,
                    project_name: project
                })
            });

            if (!tokensResponse.ok) {

                 const dfa_token_response = await fetch(`http://localhost:8080/api/lexing/dfaToTokens`, {
                    method: 'POST',
                    headers: {
                        'accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        users_id: user_id,
                        project_name: project
                    })
                });

                if (!dfa_token_response) {
                    AddToast('Tokens required: Please complete lexical analysis and generate tokens first', 'error');
                    dispatch('parsingerror', {parsing_error: true,parsing_error_details: 'Tokens required: Please complete lexical analysis and generate tokens first'});
                    return;
                }
            }

            // Now try to generate syntax tree
            const response = await fetch('http://localhost:8080/api/parsing/tree', {
                method: 'POST',
                headers: {
                    'accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    users_id: user_id,
                    project_name: project
                })
            });

            const data = await response.json();
            console.log('Syntax tree response:', data); // Debug log

            if (!response.ok) {
                if (data.error) {
                    throw new Error(data.details);
                }
                throw new Error('Failed to generate syntax tree');
            }

            if (data.tree) {
                AddToast('Parse tree generated successfully! Your syntax analysis is complete', 'success');
                dispatch('treereceived', data.tree);
            } else {
                dispatch('parsingerror',{parsing_error: true,parsing_error_details: 'A'});
                throw new Error('No tree data in response');
            }
        } catch (error) {
            dispatch('parsingerror',{parsing_error: true,parsing_error_details: error});
            console.error('Full error details:', error); // Debug log
            AddToast(
                `Failed to generate syntax tree: ${error.message}. Please ensure tokens and grammar are valid.`,
                'error'
            );
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

    <div class="grammar-editor">
        <div class="grammar-header">
            <h3>Context-Free Grammar</h3>
            <button
                class="default-toggle-btn"
                class:selected={show_default_grammar}
                on:click={show_default_grammar ? removeDefaultGrammar : insertDefaultGrammar}
                type="button"
                aria-label={show_default_grammar ? 'Remove default grammar' : 'Insert default grammar'}
                title={show_default_grammar ? 'Remove default grammar' : 'Insert default grammar'}
            >
                <span class="icon">{show_default_grammar ? '🧹' : '🪄'}</span>
            </button>
        </div>

        <div class="top-inputs">
            <div class="input-group">
                <label for="variables">Variables</label>
                <input
                    id="variables"
                    type="text"
                    placeholder="S, Y, Z"
                    bind:value={variables_string}
                    on:input={handleGrammarChange}
                />
            </div>
            <div class="input-group">
                <label for="terminals">Terminals</label>
                <input
                    id="terminals"
                    type="text"
                    placeholder="a, b, c"
                    bind:value={terminals_string}
                    on:input={handleGrammarChange}
                />
            </div>
        </div>

        <div class="rules-container">
            {#each grammar_rules as rule, i (rule.id)}
                <div class="rule-row">
                    <div class="rule-label">
                        {#if i === 0}
                            <span class="start-label">Start</span>
                        {/if}
                    </div>

                    <div class="rule-inputs">
                        <input
                            type="text"
                            class="non-terminal-input"
                            placeholder="LHS"
                            bind:value={rule.nonTerminal}
                            on:input={handleGrammarChange}
                        />
                        <span class="arrow">→</span>
                        <div class="translations-container">
                            {#each rule.translations as translation (translation.id)}
                                <div class="translation-wrapper">
                                    <input
                                        type="text"
                                        class="translation-input"
                                        placeholder="RHS"
                                        bind:value={translation.value}
                                        on:input={handleGrammarChange}
                                    />
                                </div>
                            {/each}
                            <button
                                class="add-translation-btn"
                                title="Add translation"
                                on:click={() => addTranslation(rule.id)}
                            >
                                +
                            </button>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
        <button class="add-rule-btn" on:click={addNewRule}>+ Add New Rule</button>
    </div>

    <div class="button-container">
        <button class="submit-button" on:click={handleSubmitGrammar}>Submit Grammar</button>
        {#if is_grammar_submitted}
            <button class="submit-button" on:click={generateSyntaxTree}>Generate Syntax Tree</button>
        {/if}
    </div>
</div>

<style>
    .phase-inspector {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }
    .parser-heading-h1 {
        color: black;
        text-align: center;
        margin-top: 0;
        font-family: 'Times New Roman';
        margin-bottom: 0;
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
        max-height: 250px;
		overflow: auto;
        font-family: monospace;
        white-space: pre-wrap;
        margin: 0;
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
    .default-toggle-btn {
        background: white;
        border: 2px solid #e5e7eb;
        color: #001a6e;
        font-size: 1.2rem;
        cursor: pointer;
        transition:
            background 0.2s,
            border-color 0.2s;
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
    .rule-row {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .rule-label {
        width: 60px;
        flex-shrink: 0;
        text-align: right;
        padding-right: 0.75rem;
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
    .translations-container {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex: 1;
        overflow-x: auto;
        padding-bottom: 0.5rem;
    }
    .translation-wrapper {
        position: relative;
        display: flex;
    }
    .translation-input {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-family: monospace;
        min-width: 100px;
        transition: border-color 0.2s, box-shadow 0.2s;
        width: 140px;
    }
    .translation-input:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
    }
    .add-translation-btn {
        background: #e0e0e0;
        color: #333;
        border: none;
        width: 30px;
        height: 30px;
        border-radius: 50%;
        cursor: pointer;
        font-size: 1.2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background-color 0.2s;
        flex-shrink: 0;
    }
    .add-translation-btn:hover {
        background-color: #ccc;
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
    .submit-button {
        padding: 0.6rem 1.5rem;
        background: #BED2E6;
        color: 000000;
        border: none;
        border-radius: 6px;
        font-size: 0.9rem;
        font-weight: 500;
        cursor: pointer;
    }

    .submit-button:hover {
        background: #a8bdd1;
        transform: translateY(-2px);
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
    :global(html.dark-mode) .translation-input {
        background-color: #2d3748;
        border-color: #4b5563;
        color: #f0f0f0;
    }
    :global(html.dark-mode) .input-group input:focus,
    :global(html.dark-mode) .non-terminal-input:focus,
    :global(html.dark-mode) .translation-input:focus {
        border-color: #60a5fa;
        box-shadow: 0 0 0 2px rgba(96, 165, 250, 0.4); 
    }

    :global(html.dark-mode) .arrow {
        color: #9ca3af;
    }
    :global(html.dark-mode) .add-translation-btn {
        background-color: #4b5563;
        color: #f0f0f0;
    }
    :global(html.dark-mode) .add-rule-btn {
        border-color: #60a5fa;
        color: #60a5fa;
    }
    :global(html.dark-mode) .submit-button {
        background-color: #001A6E;
        color: #ffffff;
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
        background-color: #374151;
        border-color: #6b7280;
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
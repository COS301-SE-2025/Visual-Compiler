<script lang="ts">
	// @ts-nocheck
	import type { Token } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { onMount, onDestroy } from 'svelte';
	import { DataSet, Network } from 'vis-network/standalone';
	import { fade, scale } from 'svelte/transition';
	import { projectName } from '$lib/stores/project';
	import { get } from 'svelte/store'; 
	import { lexerState, updateLexerInputs, updateAutomataInputs, markLexerSubmitted, updateLexerArtifacts } from '$lib/stores/lexer';
	import { confirmedSourceCode } from '$lib/stores/source-code';

	export let source_code = '';
	export let onGenerateTokens: (data: {
		tokens: Token[];
		unexpected_tokens: string[];
	}) => void = () => {};

	let tokens: Token[] = [];
    let show_tokens = false;
    let isSubmitted = false;
	let hasLocalChanges = false;

    // FIX: Add loading states for buttons
    let isSubmittingRules = false;
    let isGeneratingTokens = false;

	let inputRows = [{ type: '', regex: '', error: '' }];
	let userSourceCode = '';
	let userInputRows = [{ type: '', regex: '', error: '' }];
	let formError = '';
	let submissionStatus = { show: false, success: false, message: '' };
	let showGenerateButton = false;
	let showRegexActionButtons = false;
	let regexRulesSubmitted = false;
	let showAutomataVisOnly = false;

	// New state for the expandable modal
	let isExpanded = false;
	let currentAutomatonForModal = null; // To store the data for the modal
	let expandedVisContainer: HTMLElement; // Container for the expanded view
	let networkInstance = null; // To hold the active vis-network instance


	function addNewRow() {
        if (showDefault) {
            editableDefaultRows = [...editableDefaultRows, { type: '', regex: '', error: '' }];
        } else {
            userInputRows = [...userInputRows, { type: '', regex: '', error: '' }];
            hasLocalChanges = true; // Mark local changes when adding new row
            isSubmitted = false; // Reset submission status
        }
    }

	let hasInitialized = false;
	let currentProjectName = '';

	// FIX: Simplify project change detection
	$: if ($projectName !== currentProjectName) {
		console.log('Lexer: Project changed from', currentProjectName, 'to', $projectName);
		hasInitialized = false;
		currentProjectName = $projectName;
	}

	// FIX: Initialize from store when project name is set and store has data
	 $: if ($lexerState && $projectName && (!hasInitialized || $projectName !== currentProjectName)) {
        console.log('Lexer component initializing with project:', $projectName, 'and state:', $lexerState);
        
        // Restore inputs from store
        selectedType = $lexerState.selectedType || 'REGEX';
        userInputRows = $lexerState.userInputRows && $lexerState.userInputRows.length > 0 
            ? [...$lexerState.userInputRows] 
            : [{ type: '', regex: '', error: '' }];
        
        // Restore automata inputs
        if ($lexerState.automataInputs) {
            states = $lexerState.automataInputs.states || '';
            startState = $lexerState.automataInputs.startState || '';
            acceptedStates = $lexerState.automataInputs.acceptedStates || '';
            transitions = $lexerState.automataInputs.transitions || '';
        }
        
        // FIX: Restore tokens if they exist
        if ($lexerState.hasTokens && $lexerState.tokens) {
            tokens = [...$lexerState.tokens];
            show_tokens = true;
        } else {
            tokens = [];
            show_tokens = false;
        }
        
        source_code = $lexerState.sourceCode || $confirmedSourceCode || '';

        if (!hasLocalChanges) {
            isSubmitted = $lexerState.isSubmitted || false;
            
            // Set up action buttons based on submission state
            if (isSubmitted && selectedType === 'REGEX') {
                showRegexActionButtons = true;
            } else if (isSubmitted && selectedType === 'AUTOMATA') {
                showGenerateButton = true;
            }
        }
        
        hasInitialized = true;
        console.log('Lexer component initialized with:', { 
            selectedType, 
            userInputRows: userInputRows.length, 
            tokens: tokens.length,
            isSubmitted
        });
    }

	function validateRegex(pattern: string): boolean {
		try {
			new RegExp(pattern);
			return true;
		} catch (e) {
			return false;
		}
	}

	// FIX: Enhanced handleSubmit with loading animation
    async function handleSubmit() {
        formError = '';
        let hasErrors = false;
        let nonEmptyRows = [];
        const rowsToValidate = showDefault ? editableDefaultRows : userInputRows;

        for (const row of rowsToValidate) {
            if (!row.type && !row.regex) continue;
            row.error = '';
            if (!row.type || !row.regex) {
                row.error = 'Please fill in both Type and Regular Expression';
                hasErrors = true;
            } else if (!validateRegex(row.regex)) {
                row.error = 'Invalid regular expression pattern';
                hasErrors = true;
            }
            nonEmptyRows.push(row);
        }

        if (rowsToValidate.length === 1 && !rowsToValidate[0].type && !rowsToValidate[0].regex) {
            AddToast('Please fill in both Type and Regular Expression', 'error');
            return;
        }

        if (showDefault) {
            editableDefaultRows = nonEmptyRows;
        } else {
            userInputRows = nonEmptyRows;
        }

        if (hasErrors) {
            AddToast('Please fix the errors before submitting', 'error');
            return;
        }

        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName);
        
        if (!accessToken) {
            AddToast('Authentication required: Please log in to save lexical rules', 'error');
            return;
        }
        if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }

        // FIX: Set loading state
        isSubmittingRules = true;

        if (selectedType === 'REGEX') {
            const requestData = {
                project_name: project,
                pairs: nonEmptyRows.map((row) => ({
                    Type: row.type.toUpperCase(),
                    Regex: row.regex
                }))
            };
            try {
                // FIX: Add 1-second delay before API call
                await new Promise(resolve => setTimeout(resolve, 1000));

                const res = await fetch('http://localhost:8080/api/lexing/rules', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${accessToken}`
                    },
                    body: JSON.stringify(requestData)
                });
                if (!res.ok) {
                    const errorText = await res.text();
                    throw new Error(errorText);
                }
                
                AddToast('Rules stored successfully!', 'success');
                showRegexActionButtons = true;
                
                isSubmitted = true;
                hasLocalChanges = false; 
				
                markLexerSubmitted();

            } catch (error) {
                AddToast('Save failed: Unable to store lexical rules. Please check your connection and try again', 'error');
                // FIX: Don't set isSubmitted on error
                isSubmitted = false;
            } finally {
                // FIX: Reset loading state
                isSubmittingRules = false;
            }
            return;
        }

        // For AUTOMATA type
        try {
            // FIX: Add 1-second delay before processing
            await new Promise(resolve => setTimeout(resolve, 1000));

            AddToast('Ready for tokenization!', 'success');
            showGenerateButton = true;
            isSubmitted = true;
            hasLocalChanges = false; 
            markLexerSubmitted();
        } catch (error) {
            AddToast('Failed to submit automata rules. Please try again.', 'error');
			isSubmitted = false;
        } finally {
            // FIX: Reset loading state
            isSubmittingRules = false;
        }
    }

    // FIX: Enhanced generateTokens with loading animation
    async function generateTokens() {
        // FIX: Set loading state
        isGeneratingTokens = true;

        try {
            // FIX: Add 1-second delay before API call
            await new Promise(resolve => setTimeout(resolve, 1000));

            const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
            const project = get(projectName);
            
            if (!accessToken) {
                AddToast('Authentication required: Please log in to generate tokens', 'error');
                return;
            }
            if (!project) {
                AddToast('No project selected: Please select or create a project first', 'error');
                return;
            }

            const requestData = {
                project_name: project
            };

            const response = await fetch('http://localhost:8080/api/lexing/lexer', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Server error (${response.status}): ${errorText}`);
            }
            const data = await response.json();
            if (!Array.isArray(data.tokens)) {
                throw new Error('Expected tokens array in response');
            }

            onGenerateTokens({
                tokens: data.tokens,
                unexpected_tokens: data.tokens_unidentified
            });

            showGenerateButton = false;
            
            updateLexerArtifacts(data.tokens, source_code);
            
            AddToast(data.message || 'Tokens generated successfully!', 'success');
		 } catch (error) {
            console.error('Generate tokens error:', error);
            AddToast('Tokenization failed: Unable to generate tokens from your lexical rules', 'error');
        } finally {
            isGeneratingTokens = false;
        }
	}

	// Helper: Save DFA to backend
	async function saveDfaToBackend() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName); 

		if (!accessToken) {
			AddToast('Authentication required: Please log in to save automata data', 'error');
			return;
		}

		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}

		const dfa = {
			states: states
				.split(',')
				.map((s) => s.trim())
				.filter(Boolean),
			transitions: [],
			start_state: startState.trim(),
			accepting_states: parseAcceptedStates(acceptedStates),
			project_name: project
		};

		// Parse transitions
		const transitionLines = transitions
			.split('\n')
			.map((line) => line.trim())
			.filter(Boolean);
		for (const line of transitionLines) {
			const match = line.match(/^(.+),(.+)->(.+)$/);
			if (match) {
				const [_, from, label, to] = match;
				dfa.transitions.push({ from: from.trim(), label: label.trim(), to: to.trim() });
			}
		}

		// Log the DFA JSON being sent
		console.log('DFA sent to backend:', JSON.stringify(dfa, null, 2));

		try {
			const response = await fetch('http://localhost:8080/api/lexing/dfa', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify(dfa)
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Save failed: Unable to save DFA - ' + errorText, 'error');
				return false;
			}
			return true;
		} catch (error) {
			AddToast('Save error: Failed to save DFA - ' + error, 'error');
			return false;
		}
	}

	 async function handleTokenisation() {
        // FIX: Set loading state
        isGeneratingTokens = true;

        try {
            // FIX: Add 1-second delay before processing
            await new Promise(resolve => setTimeout(resolve, 1000));

            const saved = await saveDfaToBackend();
            if (!saved) return;

            const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
            const project = get(projectName); 

            if (!accessToken) {
                AddToast('Authentication required: Please log in to perform tokenization', 'error');
                return;
            }

            if (!project) {
                AddToast('No project selected: Please select or create a project first', 'error');
                return;
            }
			
            const body = { project_name: project };

            const response = await fetch('http://localhost:8080/api/lexing/dfaToTokens', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(body)
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }
            const data = await response.json();
                        
            onGenerateTokens({
                tokens: data.tokens,
                unexpected_tokens: data.tokens_unidentified
            });

            updateLexerArtifacts(data.tokens, source_code);
            
            AddToast('Tokenization complete! Your source code has been successfully tokenized', 'success');
        } catch (error) {
            AddToast('Tokenization failed: ' + error, 'error');
        } finally {
            // FIX: Reset loading state
            isGeneratingTokens = false;
        }
    }

	let previousInputs: typeof userInputRows = [];
	function handleInputChange() {
        hasLocalChanges = true; // Mark that user made local changes
        showGenerateButton = false;
        showRegexActionButtons = false;
        
        // FIX: Reset submission status when input changes
        isSubmitted = false;
        show_tokens = false;
        tokens = [];
        
        // Update the store with current inputs but don't mark as submitted
        updateLexerInputs(userInputRows);
    }

	// Update automata input change handler
	function handleAutomataInputChange() {
        hasLocalChanges = true; // Mark that user made local changes
        
        // FIX: Reset submission status when automata input changes
        isSubmitted = false;
        showGenerateButton = false;
        show_tokens = false;
        tokens = [];
        
        updateAutomataInputs({
            states,
            startState,
            acceptedStates,
            transitions
        });
    }


	let selectedType: 'AUTOMATA' | 'REGEX' | null = 'REGEX';
	let showDefault = false;
	let states = '';
	let startState = '';
	let acceptedStates = '';
	let transitions = '';

	function parseAutomaton() {
		const stateList = states
			.split(',')
			.map((s) => s.trim())
			.filter(Boolean);
		const start = startState.trim();
		// Extract only the state names from acceptedStates
		const accepted = parseAcceptedStates(acceptedStates).map((a) => a.state);
		const transitionLines = transitions
			.split('\n')
			.map((line) => line.trim())
			.filter(Boolean);
		const transitionObj: Record<string, Record<string, string[]>> = {};
		let alphabetSet = new Set<string>();
		for (const line of transitionLines) {
			const match = line.match(/^(.+),(.+)->(.+)$/);
			if (match) {
				const [_, from, symbol, to] = match;
				alphabetSet.add(symbol);
				if (!transitionObj[from]) transitionObj[from] = {};
				if (!transitionObj[from][symbol]) transitionObj[from][symbol] = [];
				transitionObj[from][symbol].push(to);
			}
		}
		return {
			states: stateList,
			startState: start,
			acceptedStates: accepted, // <-- now an array of state names
			transitions: transitionObj,
			alphabet: Array.from(alphabetSet)
		};
	}

	function nfaToDfa(nfa: any) {
		const { alphabet, startState, acceptedStates, transitions } = nfa;
		const dfaStates: string[] = [];
		const dfaTransitions: Record<string, Record<string, string>> = {};
		const dfaAcceptedStates: string[] = [];

		function stateSetToName(set: string[]) {
			return set.sort().join(',');
		}

		let unmarked: string[][] = [[startState]];

		while (unmarked.length > 0) {
			const currentSet = unmarked.pop()!;
			const name = stateSetToName(currentSet);
			if (!dfaStates.includes(name)) dfaStates.push(name);

			dfaTransitions[name] = {};

			for (const symbol of alphabet) {
				let nextSet: string[] = [];
				for (const state of currentSet) {
					if (transitions[state] && transitions[state][symbol]) {
						nextSet = nextSet.concat(transitions[state][symbol]);
					}
				}
				nextSet = Array.from(new Set(nextSet));
				if (nextSet.length === 0) continue;
				const nextName = stateSetToName(nextSet);
				dfaTransitions[name][symbol] = nextName;
				if (
					!dfaStates.includes(nextName) &&
					!unmarked.some((s) => stateSetToName(s) === nextName)
				) {
					unmarked.push(nextSet);
				}
			}
		}

		for (const dfaState of dfaStates) {
			const nfaStatesInDfa = dfaState.split(',');
			if (nfaStatesInDfa.some((s) => acceptedStates.includes(s))) {
				dfaAcceptedStates.push(dfaState);
			}
		}

		return {
			states: dfaStates,
			startState: stateSetToName([startState]),
			acceptedStates: dfaAcceptedStates,
			transitions: dfaTransitions,
			alphabet
		};
	}

	let nfaContainer: HTMLElement;
	let dfaContainer: HTMLElement;
	let showNfaVis = false;
	let showDfaVis = false;

	function safeStateId(name: string) {
		return name.replace(/[^a-zA-Z0-9_]/g, '_');
	}

	async function showNfaDiagram() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName); 

		if (!accessToken) {
			AddToast('Authentication required: Please log in to generate NFA', 'error');
			return;
		}

		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}

		const saved = await saveDfaToBackend();
		if (!saved) return;

		try {
			// Convert DFA to Regex
			const dfaToRegexRes = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify({ project_name: project })
			});
			if (!dfaToRegexRes.ok) {
				const errorText = await dfaToRegexRes.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}

			// Convert Regex to NFA
			const regexToNfaRes = await fetch('http://localhost:8080/api/lexing/regexToNFA', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify({ project_name: project })
			});
			if (!regexToNfaRes.ok) {
				const errorText = await regexToNfaRes.text();
				AddToast('Regex→NFA failed: ' + errorText, 'error');
				return;
			}
			const nfaData = await regexToNfaRes.json();
			regexNfa = adaptAutomatonForVis(nfaData.nfa);
			currentAutomatonForModal = { data: regexNfa, isDfa: false };
			showNfaVis = true;
			showDfaVis = false;
			automataDisplay = 'NFA';
			setTimeout(() => renderRegexAutomatonVis(nfaContainer, regexNfa, false), 0);
			AddToast('NFA generated from Regex and displayed!', 'success');
		} catch (error) {
			AddToast('Failed to generate NFA from Regex: ' + error, 'error');
		}
	}

	function resetInputs() {
		states = '';
		startState = '';
		acceptedStates = '';
		transitions = '';
		showDefault = false;
	}

	function selectType(type: 'AUTOMATA' | 'REGEX') {
		selectedType = type;
		showRegexActionButtons = false;
		regexRulesSubmitted = false;
		
		// Reset visualization states only
		showNfaVis = false;
		showDfaVis = false;
		showRegexNfaVis = false;
		showRegexDfaVis = false;
		showRegexOutput = false;
		automataDisplay = null;
		currentAutomatonForModal = null;

		isSubmitted = false;
        showGenerateButton = false;

		if (type === 'REGEX') {
			// Use stored inputs if they exist
			if (userInputRows.length === 0) {
                userInputRows = $lexerState.userInputRows.length > 0 
                    ? [...$lexerState.userInputRows]
                    : [{ type: '', regex: '', error: '' }];
            }
		}
	}

	function insertDefault() {
		showDefault = true;
		editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
		inputRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
		
		states = 'S0, S1, S2, S3, S4, S5, S6, S7, S8, S9, S10, S11, S12, S13, S14, S15, S16, S17, S18, S19, S20, S21, S22, S23, S24, S25, S26, S27, S28, S29, S30';
		startState = 'S0';
		acceptedStates = 'S1->ASSIGNMENT, S2->SEPARATOR, S4->OPEN_BRACKET, S6->IDENTIFIER, S7->INTEGER, S8->CLOSE_SCOPE, S10->OPERATOR, S11->CLOSE_BRACKET, S12->OPEN_SCOPE, S19->KEYWORD, S23->CONTROL, S27->KEYWORD, S29->CONTROL, S30->KEYWORD';
		transitions = '';
		transitions += 'S0,i->S3\n';
		transitions += 'S0,r->S9\n';
		transitions += 'S0,f->S13\n';
		transitions += 'S0,p->S5\n';
		transitions += 'S0,[a-zA-Z_]->S6\n';
		transitions += 'S0,[0-9]->S7\n';
		transitions += 'S0,+/%->S10\n';
		transitions += 'S0,=->S1\n';
		transitions += 'S0,;->S2\n';
		transitions += 'S0,(->S4\n';
		transitions += 'S0,)->S11\n';
		transitions += 'S0,{->S12\n';
		transitions += 'S0,}->S8\n';
		transitions += 'S3,n->S14\n';
		transitions += 'S5,r->S15\n';
		transitions += 'S6,[a-zA-Z_]->S6\n';
		transitions += 'S7,[0-9]->S7\n';
		transitions += 'S9,e->S16\n';
		transitions += 'S9,a->S17\n';
		transitions += 'S13,o->S18\n';
		transitions += 'S14,t->S19\n';
		transitions += 'S15,i->S20\n';
		transitions += 'S16,t->S21\n';
		transitions += 'S17,n->S22\n';
		transitions += 'S18,r->S23\n';
		transitions += 'S20,n->S24\n';
		transitions += 'S21,u->S25\n';
		transitions += 'S22,g->S26\n';
		transitions += 'S24,t->S27\n';
		transitions += 'S25,r->S28\n';
		transitions += 'S26,e->S29\n';
		transitions += 'S28,n->S30\n';
	}

	function removeDefault() {
		showDefault = false;
		inputRows = [...userInputRows];
		resetInputs();
	}

	const DEFAULT_INPUT_ROWS = [
		{ type: 'KEYWORD', regex: 'int|return|print', error: '' },
		{ type: 'CONTROL', regex: 'for|range', error: '' },
		{ type: 'IDENTIFIER', regex: '[a-zA-Z_]+', error: '' },
		{ type: 'INTEGER', regex: '[0-9]+', error: '' },
		{ type: 'ASSIGNMENT', regex: '=', error: '' },
		{ type: 'OPERATOR', regex: '[+\\-*/%]', error: '' },
		{ type: 'SEPARATOR', regex: ';', error: '' },
        { type: 'OPEN_BRACKET', regex: '\\(', error: '' },
        { type: 'CLOSE_BRACKET', regex: '\\)', error: '' },
        { type: 'OPEN_SCOPE', regex: '\\{', error: '' },
        { type: 'CLOSE_SCOPE', regex: '\\}', error: '' },
		{ type: 'STRING', regex: '"[a-zA-Z\\s]+"', error: '' }
	];

	let editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
	const DEFAULT_SOURCE_CODE = 'int blue = 13 + 22 ;';

	$: if (!showDefault && source_code) {
		userSourceCode = source_code;
	}

	function parseAcceptedStates(input: string): { state: string; token_type: string }[] {
		return input
			.split(',')
			.map((s) => s.trim())
			.filter(Boolean)
			.map((pair) => {
				const [state, token_type] = pair.split('->').map((x) => x.trim());
				if (state && token_type) {
					return { state, token_type };
				}
				// fallback: if only state is provided, treat type as empty string
				if (state) {
					return { state, token_type: '' };
				}
				return null;
			})
			.filter((x): x is { state: string; token_type: string } => !!x);
	}

	// Show DFA button handler
	async function handleShowDfa() {
        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName); 

        if (!accessToken) {
            AddToast('Authentication required: Please log in to generate DFA', 'error');
            return;
        }

        if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }

        const saved = await saveDfaToBackend();
        if (!saved) return;

        try {
            // Convert DFA to Regex
            const regexRes = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({ project_name: project })
            });
            if (!regexRes.ok) {
                const errorText = await regexRes.text();
                AddToast('DFA→Regex failed: ' + errorText, 'error');
                return;
            }

			// Convert Regex to DFA
            const dfaRes = await fetch('http://localhost:8080/api/lexing/regexToDFA', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({ project_name: project })
            });

            if (!dfaRes.ok) {
                const errorText = await dfaRes.text();
                AddToast('Regex→DFA failed: ' + errorText, 'error');
                return;
            }
            const dfaData = await dfaRes.json();
            regexDfa = adaptAutomatonForVis(dfaData.dfa);
            currentAutomatonForModal = { data: regexDfa, isDfa: true }; 
            
            // FIX: Set states to show only visualization like regex section
            showDfaVis = true;
            showNfaVis = false;
            showAutomataVisOnly = true; // This will hide the input and show only visualization
            automataDisplay = 'DFA'; // Set display type for consistency
            
            setTimeout(() => renderAutomatonVis(dfaContainer, regexDfa, true), 0);
            AddToast('DFA generated from automata and displayed!', 'success');
        } catch (error) {
            AddToast('Failed to generate DFA from automata: ' + error, 'error');
        }
    }

	let regexRules: { token_type?: string; Type?: string; Regex?: string; regex?: string }[] = [];
	let showRegexOutput = false;

	 async function handleConvertToRegex() {
        const saved = await saveDfaToBackend();
        if (!saved) return;

        const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
        const project = get(projectName); 

        if (!accessToken) {
            AddToast('Authentication required: Please log in to convert DFA', 'error');
            return;
        }

        if (!project) {
            AddToast('No project selected: Please select or create a project first', 'error');
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({ project_name: project })
            });
            if (!response.ok) {
                const errorText = await response.text();
                AddToast('DFA→Regex failed: Please check your DFA input.','error');
                return;
            }
            const data = await response.json();
            regexRules = Array.isArray(data.rules) ? data.rules : [];
            
            // FIX: Set states to show only RE display
            showRegexOutput = true;
			 showAutomataVisOnly = true; // This will hide the input and show only RE display
            automataDisplay = 'RE'; // Set display type for consistency
            
            AddToast('DFA converted to Regex successfully!', 'success');
        } catch (error) {
            AddToast('DFA→Regex failed: Please check your internet connection.', 'error');
        }
	}

	function handleBackFromAutomataVis() {
        showAutomataVisOnly = false;
        showDfaVis = false;
        showNfaVis = false;
        showRegexOutput = false;
        automataDisplay = null; // Reset display type
        currentAutomatonForModal = null; // Clear modal data
    }

	let regexNfa = null;
	let regexDfa = null;
	let showRegexNfaVis = false;
	let showRegexDfaVis = false;
	let regexNfaContainer: HTMLElement;
	let regexDfaContainer: HTMLElement;

	async function handleRegexToNFA() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName); 

		if (!accessToken) {
			AddToast('Authentication required: Please log in to convert Regex', 'error');
			return;
		}

		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}
		try {
			const response = await fetch('http://localhost:8080/api/lexing/regexToNFA', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify({ project_name: project })
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Regex→NFA failed: Please check your regex rules', 'error');
				return;
			}
			const data = await response.json();
			regexNfa = adaptAutomatonForVis(data.nfa);
			currentAutomatonForModal = { data: regexNfa, isDfa: false }; 
            
            // FIX: Set states to show only visualization like automata section
            showRegexNfaVis = true;
            showRegexDfaVis = false;
            showRegexVisOnly = true; // This will hide the input and show only visualization
            automataDisplay = 'NFA'; // Set display type for consistency
            
			AddToast('Regex converted to NFA!', 'success');
			setTimeout(() => renderRegexAutomatonVis(regexNfaContainer, regexNfa, false), 0);
		} catch (error) {
			AddToast('Regex→NFA failed: ' + error, 'error');
		}
	}

	async function handleRegexToDFA() {
		const accessToken = sessionStorage.getItem('access_token') || sessionStorage.getItem('authToken');
		const project = get(projectName); 

		if (!accessToken) {
			AddToast('Authentication required: Please log in to convert Regex', 'error');
			return;
		}

		if (!project) {
			AddToast('No project selected: Please select or create a project first', 'error');
			return;
		}
		try {
			const response = await fetch('http://localhost:8080/api/lexing/regexToDFA', {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${accessToken}`
				},
				body: JSON.stringify({ project_name: project })
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Regex→DFA failed: ' + errorText, 'error');
				return;
			}
			const data = await response.json();
			console.log('DFA from backend:', JSON.stringify(data.dfa, null, 2));
			regexDfa = adaptAutomatonForVis(data.dfa);
			currentAutomatonForModal = { data: regexDfa, isDfa: true };
            
            // FIX: Set states to show only visualization like automata section
            showRegexDfaVis = true;
            showRegexNfaVis = false;
            showRegexVisOnly = true; // This will hide the input and show only visualization
            automataDisplay = 'DFA'; // Set display type for consistency
            
			AddToast('Regex converted to DFA!', 'success');
			setTimeout(() => renderRegexAutomatonVis(regexDfaContainer, regexDfa, true), 0);
		} catch (error) {
			AddToast('Regex→DFA failed: ' + error, 'error');
		}
	}

	// FIX: Update the back function to restore the regex input view
	function handleBackFromRegexVis() {
		showRegexVisOnly = false;
		showRegexNfaVis = false;
		showRegexDfaVis = false;
		automataDisplay = null; // Reset display type
		currentAutomatonForModal = null; // Clear modal data
		
		// This will show the regex input rows again
	}

	// Helper function to render the automaton visualization
	function renderAutomatonVis(container: HTMLElement, automaton: any, isDfa: boolean) {
		if (!automaton || !container) return;

		if (networkInstance) {
			networkInstance.destroy();
			networkInstance = null;
		}

		const nodeIds: Record<string, string> = {};
		automaton.states.forEach((state: string) => {
			nodeIds[state] = state.replace(/[^a-zA-Z0-9_]/g, '_');
		});
		const nodes = new DataSet(
			automaton.states.map((state: string) => ({
				id: nodeIds[state],
				label: state,
				shape: 'circle',
				color: automaton.acceptedStates.includes(state)
					? '#D2FFD2'
					: state === automaton.startState
					? '#D2E5FF'
					: '#FFD2D2',
				borderWidth: automaton.acceptedStates.includes(state) ? 3 : 1
			}))
		);
		const edgesArr: any[] = [];
		for (const from of automaton.states) {
			for (const symbol of automaton.alphabet) {
				const tos = isDfa
					? [automaton.transitions[from]?.[symbol]].filter(Boolean)
					: automaton.transitions[from]?.[symbol] || [];
				for (const to of tos) {
					edgesArr.push({ from: nodeIds[from], to: nodeIds[to], label: symbol, arrows: 'to' });
				}
			}
		}
		const START_NODE_ID = '__start__';
		nodes.add({
			id: START_NODE_ID,
			label: '',
			shape: 'circle',
			color: 'rgba(0,0,0,0)',
			borderWidth: 0,
			size: 1,
			font: { size: 1 }
		});
		edgesArr.push({
			from: START_NODE_ID,
			to: nodeIds[automaton.startState],
			arrows: { to: { enabled: true, scaleFactor: 0.6 } },
			color: { color: '#222', opacity: 1 },
			width: 1.75,
			label: 'start',
			font: { size: 13, color: '#222', vadjust: -18, align: 'top' },
			smooth: { enabled: true, type: 'curvedCCW', roundness: 0.18 },
			length: 1,
			physics: false
		});
		const edges = new DataSet(edgesArr);
		networkInstance = new Network(
			container,
			{ nodes, edges },
			{
				nodes: { shape: 'circle', font: { size: 16 }, margin: 10 },
				edges: {
					smooth: { enabled: true, type: 'curvedCW', roundness: 0.3 },
					font: { size: 14, strokeWidth: 0 }
				},
				physics: false
			}
		);
	}

	function adaptAutomatonForVis(automaton: any) {
		let transitionsObj: Record<string, Record<string, string[]>> = {};
		if (Array.isArray(automaton.transitions)) {
			for (const t of automaton.transitions) {
				if (!transitionsObj[t.from]) transitionsObj[t.from] = {};
				if (!transitionsObj[t.from][t.label]) transitionsObj[t.from][t.label] = [];
				transitionsObj[t.from][t.label].push(t.to);
			}
		} else {
			transitionsObj = automaton.transitions || {};
		}

		let alphabet: string[] = automaton.alphabet || [];
		if ((!alphabet || alphabet.length === 0) && Object.keys(transitionsObj).length > 0) {
			const symbols = new Set<string>();
			for (const from in transitionsObj) {
				for (const symbol in transitionsObj[from]) {
					symbols.add(symbol);
				}
			}
			alphabet = Array.from(symbols);
		}

		return {
			states: automaton.states || Object.keys(transitionsObj),
			startState: automaton.start_state || automaton.startState || automaton.start || '',
			acceptedStates: (
				automaton.accepting_states ||
				automaton.acceptedStates ||
				automaton.accepting ||
				[]
			).map((a: any) => (typeof a === 'string' ? a : a.state || a)),
			transitions: transitionsObj,
			alphabet
		};
	}

	function renderRegexAutomatonVis(container: HTMLElement, automaton: any, isDfa = false) {
		if (!automaton || !container) return;

		if (networkInstance) {
			networkInstance.destroy();
			networkInstance = null;
		}

		const nodeIds: Record<string, string> = {};
		automaton.states.forEach((state: string) => {
			nodeIds[state] = state.replace(/[^a-zA-Z0-9_]/g, '_');
		});
		const nodes = new DataSet(
			automaton.states.map((state: string) => ({
				id: nodeIds[state],
				label: state,
				shape: 'circle',
				color: automaton.acceptedStates.includes(state)
					? '#D2FFD2'
					: state === automaton.startState
					? '#D2E5FF'
					: '#FFD2D2',
				borderWidth: automaton.acceptedStates.includes(state) ? 3 : 1
			}))
		);
		const edgesArr: any[] = [];
		for (const from of automaton.states) {
			for (const symbol of automaton.alphabet) {
				const tos = isDfa
					? [automaton.transitions[from]?.[symbol]].filter(Boolean)
					: automaton.transitions[from]?.[symbol] || [];
				for (const to of tos) {
					edgesArr.push({ from: nodeIds[from], to: nodeIds[to], label: symbol, arrows: 'to' });
				}
			}
		}
		const START_NODE_ID = '__start__';
		nodes.add({
			id: START_NODE_ID,
			label: '',
			shape: 'circle',
			color: 'rgba(0,0,0,0)',
			borderWidth: 0,
			size: 1,
			font: { size: 1 }
		});
		edgesArr.push({
			from: START_NODE_ID,
			to: nodeIds[automaton.startState],
			arrows: { to: { enabled: true, scaleFactor: 0.6 } },
			color: { color: '#222', opacity: 1 },
			width: 1.75,
			label: 'start',
			font: { size: 13, color: '#222', vadjust: -18, align: 'top' },
			smooth: { enabled: true, type: 'curvedCCW', roundness: 0.18 },
			length: 1,
			physics: false
		});
		const edges = new DataSet(edgesArr);
		networkInstance = new Network(
			container,
			{ nodes, edges },
			{
				nodes: { shape: 'circle', font: { size: 16 }, margin: 10 },
				edges: {
					smooth: { enabled: true, type: 'curvedCW', roundness: 0.3 },
					font: { size: 14, strokeWidth: 0 }
				},
				physics: false
			}
		);
	}

	let automataDisplay: 'NFA' | 'DFA' | 'RE' | null = null;

	let showRegexVisOnly = false;

	const toggleExpand = () => {
		isExpanded = !isExpanded;
		if (isExpanded && currentAutomatonForModal) {
			setTimeout(() => {
				renderRegexAutomatonVis(
					expandedVisContainer,
					currentAutomatonForModal.data,
					currentAutomatonForModal.isDfa
				);
			}, 50); 
		} else if (!isExpanded && currentAutomatonForModal) {
			let originalContainer = null;
			if (showRegexVisOnly) {
				if (showRegexNfaVis) originalContainer = regexNfaContainer;
				if (showRegexDfaVis) originalContainer = regexDfaContainer;
			} else {
				if (automataDisplay === 'NFA') originalContainer = nfaContainer;
				if (automataDisplay === 'DFA') originalContainer = dfaContainer;
			}

			if (originalContainer) {
				setTimeout(() => {
					renderRegexAutomatonVis(
						originalContainer,
						currentAutomatonForModal.data,
						currentAutomatonForModal.isDfa
					);
				}, 50);
			}
		}
	};

	const handleKeydown = (event: KeyboardEvent) => {
		if (isExpanded && event.key === 'Escape') {
			toggleExpand();
		}
	};

	function fitGraphToModal() {
		if (networkInstance) {
			networkInstance.fit();
		}
	}

	// Zoom functionality for expanded view
	function zoomIn() {
		if (networkInstance && isExpanded) {
			const scale = networkInstance.getScale();
			networkInstance.moveTo({
				scale: Math.min(scale * 1.2, 3) // Max zoom level of 3x
			});
		}
	}

	function zoomOut() {
		if (networkInstance && isExpanded) {
			const scale = networkInstance.getScale();
			networkInstance.moveTo({
				scale: Math.max(scale * 0.8, 0.1) // Min zoom level of 0.1x
			});
		}
	}

	function resetZoom() {
		if (networkInstance && isExpanded) {
			networkInstance.fit({
				animation: {
					duration: 500,
					easingFunction: 'easeInOutCubic'
				}
			});
		}
	}

	// Add a new function to format transitions
	function formatTransitions(transitions: Array<{from: string, to: string, label: string}>): string {
	    return transitions.map(t => `${t.from},${t.label}->${t.to}`).join('\n');
	}

	// Add a new function to format accepting states
	function formatAcceptingStates(accepting: Array<{state: string, type: string}>): string {
	    return accepting.map(a => `${a.state}->${a.type}`).join(', ');
	}

	// Add a watch on projectName store to update inputs when project changes
	$: $projectName && updateInputsFromProject();

	// Add function to update inputs from project data
	async function updateInputsFromProject() {
	    const userId = sessionStorage.getItem('user_id');
	    const project = get(projectName);
	    
	    if (!userId || !project) return;

	    try {
	        const response = await fetch(
	            `http://localhost:8080/api/users/getProject?project_name=${project}&users_id=${userId}`,
	            {
	                method: 'GET',
	                headers: { 'accept': 'application/json' }
	            }
	        );

	        if (!response || !response.ok) return;

	        const data = await response.json();
	        
	        if (data.results?.lexing) {
	            // Update automata inputs if DFA exists
	            if (data.results.lexing.dfa) {
	                states = data.results.lexing.dfa.states.join(', ');
	                startState = data.results.lexing.dfa.start;
	                acceptedStates = formatAcceptingStates(data.results.lexing.dfa.accepting);
	                transitions = formatTransitions(data.results.lexing.dfa.transitions);
	            }

	            // Update regex inputs if rules exist
	            if (data.results.lexing.rules) {
	                userInputRows = data.results.lexing.rules.map(rule => ({
	                    type: rule.type,
	                    regex: rule.regex,
	                    error: ''
	                }));
	            }

	            // If there's no input rows, add one empty row
	            if (userInputRows.length === 0) {
	                userInputRows = [{ type: '', regex: '', error: '' }];
	            }
	        }
	    } catch (error) {
	        console.error('Error updating inputs:', error);
	    }
	}

	// Add event listener for AI-generated lexer rules
    let aiLexerEventListener: (event: CustomEvent) => void;

    onMount(() => {
        // Listen for AI-generated lexer rules
        aiLexerEventListener = (event: CustomEvent) => {
            if (event.detail && event.detail.rules && Array.isArray(event.detail.rules)) {
                console.log('Received AI lexer rules:', event.detail.rules);
                
                // Clear existing user input rows and populate with AI-generated rules
                userInputRows = event.detail.rules.map(rule => ({
                    type: rule.type,
                    regex: rule.regex,
                    error: ''
                }));
                
                // Ensure we're in regex mode and not showing default
                selectedType = 'REGEX';
                showDefault = false;
                
                // Reset other states
                showRegexActionButtons = false;
                regexRulesSubmitted = false;
                showGenerateButton = false;
                submissionStatus = { show: false, success: false, message: '' };
                
                // Update the lexer state
                lexerState.update(state => ({
                    ...state,
                    userInputRows: [...userInputRows]
                }));
                
                AddToast('AI lexer rules inserted into input rows!', 'success');
            }
        };

        window.addEventListener('ai-lexer-generated', aiLexerEventListener);
    });

    onDestroy(() => {
        if (aiLexerEventListener) {
            window.removeEventListener('ai-lexer-generated', aiLexerEventListener);
        }
    });

	function clearAllInputs() {
		if (selectedType === 'REGEX') {
			// Reset regex inputs only
			userInputRows = [{ type: '', regex: '', error: '' }];
			editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
			showRegexVisOnly = false;
			showRegexNfaVis = false;
			showRegexDfaVis = false;
		} else if (selectedType === 'AUTOMATA') {
			// Reset automata inputs only
			states = '';
			startState = '';
			acceptedStates = '';
			transitions = '';
			showAutomataVisOnly = false; // FIX: Reset automata vis only state
			showDfaVis = false;
			showNfaVis = false;
		}
		
		 // Reset common states
        showDefault = false;
        showGenerateButton = false;
        showRegexActionButtons = false;
        showRegexVisOnly = false;
        showRegexNfaVis = false;
        showRegexDfaVis = false;
        showNfaVis = false;
        showDfaVis = false;
        showRegexOutput = false;
        showAutomataVisOnly = false; 
        automataDisplay = null;
        currentAutomatonForModal = null;
        
        // Update lexer state only for regex
        if (selectedType === 'REGEX') {
            lexerState.update(state => ({
                ...state,
                userInputRows: [...userInputRows]
            }));
        }
        
        AddToast(`All ${selectedType.toLowerCase()} inputs cleared successfully!`, 'success');
    }

    let currentProject = '';
    
    // Watch for project changes and reset component state
    $: if ($projectName !== currentProject) {
        if (currentProject !== '' && $projectName !== currentProject) {
            // Project changed - reset component state
            console.log('Project changed from', currentProject, 'to', $projectName);
            
            // Reset component state to initial values
            userInputRows = [{ type: '', regex: '', error: '' }];
            states = '';
            startState = '';
            acceptedStates = '';
            transitions = '';
            showDefault = false;
            showGenerateButton = false;
            showRegexActionButtons = false;
            hasInitialized = false;
            
            // Reset visualization states
            showNfaVis = false;
            showDfaVis = false;
            showRegexNfaVis = false;
            showRegexDfaVis = false;
            showRegexOutput = false;
            automataDisplay = null;
        }
        currentProject = $projectName;
    }

	 $: if ($lexerState?.sourceCode) {
        source_code = $lexerState.sourceCode;

    }
	$: if ($confirmedSourceCode) {
        source_code = $confirmedSourceCode;
    }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="phase-inspector">
	<div class="source-code-section">
		<div class="lexor-heading">
			<h1 class="lexor-heading-h1">LEXING</h1>
		</div>
		<h3 class="source-code-header">Source Code</h3>
		<pre class="source-display">{source_code || 'no source code available'}</pre>
	</div>

	<div class="instructions-section">
		<div class="instructions-content">
			<h4 class="instructions-header">Instructions</h4>
			<p class="instructions-text">
				Enter regular expression rules in order of priority or an automata to tokenise your source code. The token types are used throughout the other phases.
			</p>
		</div>
	</div>

	<div class="automaton-btn-row">
		<div class="tabs-and-example-group">
			<button
				class="automaton-btn {selectedType === 'REGEX' ? 'selected' : ''}"
				on:click={() => selectType('REGEX')}
				type="button"
			>
				Regular Expression
			</button>
			<button
				class="automaton-btn {selectedType === 'AUTOMATA' ? 'selected' : ''}"
				on:click={() => {
					selectType('AUTOMATA');
					automataDisplay = null;
				}}
				type="button"
			>
				Automata
			</button>
		</div>

		<div class="clear-button-container">
			{#if selectedType && !showDefault}
			<button
				class="option-btn example-btn"
				on:click={insertDefault}
				type="button"
				aria-label="Show example code"
				title="Show example code"
			>
				Show Example
			</button>
			{/if}
			{#if selectedType && showDefault}
				<button
					class="option-btn example-btn selected"
					on:click={removeDefault}
					type="button"
					aria-label="Restore your input"
					title="Restore your input"
				>
					Restore Input
				</button>
			{/if}
			{#if selectedType}
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
			{/if}
		</div>
	</div>

	 {#if selectedType === 'REGEX'}
		{#if showRegexVisOnly}
			<!-- Show only NFA or DFA visualization with back button -->
			{#if showRegexNfaVis && regexNfa}
				<div class="automata-container pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">NFA Visualization (from REGEX)</span>
					</div>
					<button on:click={toggleExpand} class="expand-btn" title="Expand view">
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
							<path d="M1.5 1a.5.5 0 0 0-.5.5v4a.5.5 0 0 1-1 0v-4A1.5 1.5 0 0 1 1.5 0h4a.5.5 0 0 1 0 1zM10 .5a.5.5 0 0 1 .5-.5h4A1.5 1.5 0 0 1 16 1.5v4a.5.5 0 0 1-1 0v-4a.5.5 0 0 0-.5-.5h-4a.5.5 0 0 1-.5-.5M.5 10a.5.5 0 0 1 .5.5v4a.5.5 0 0 0 .5.5h4a.5.5 0 0 1 0 1h-4A1.5 1.5 0 0 1 0 14.5v-4a.5.5 0 0 1 .5-.5m15 0a.5.5 0 0 1 .5.5v4a1.5 1.5 0 0 1-1.5 1.5h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 0 .5-.5v-4a.5.5 0 0 1 .5-.5"/>
						</svg>
					</button>
					<div bind:this={regexNfaContainer} class="vis-graph-area"></div>
				</div>
				<button
					class="submit-button"
					style="align-self: flex-start; margin-top: 1.5rem;"
					on:click={handleBackFromRegexVis}
				>
					← Back
				</button>
			{/if}
			
			{#if showRegexDfaVis && regexDfa}
				<div class="automata-container pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">DFA Visualization (from REGEX)</span>
					</div>
					<button on:click={toggleExpand} class="expand-btn" title="Expand view">
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
							<path d="M1.5 1a.5.5 0 0 0-.5.5v4a.5.5 0 0 1-1 0v-4A1.5 1.5 0 0 1 1.5 0h4a.5.5 0 0 1 0 1zM10 .5a.5.5 0 0 1 .5-.5h4A1.5 1.5 0 0 1 16 1.5v4a.5.5 0 0 1-1 0v-4a.5.5 0 0 0-.5-.5h-4a.5.5 0 0 1-.5-.5M.5 10a.5.5 0 0 1 .5.5v4a.5.5 0 0 0 .5.5h4a.5.5 0 0 1 0 1h-4A1.5 1.5 0 0 1 0 14.5v-4a.5.5 0 0 1 .5-.5m15 0a.5.5 0 0 1 .5.5v4a1.5 1.5 0 0 1-1.5 1.5h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 0 .5-.5v-4a.5.5 0 0 1 .5-.5"/>
						</svg>
					</button>
					<div bind:this={regexDfaContainer} class="vis-graph-area"></div>
				</div>
				<button
					class="submit-button"
					style="align-self: flex-start; margin-top: 1.5rem;"
					on:click={handleBackFromRegexVis}
				>
					← Back
				</button>
			{/if}
		{:else}
			<!-- Show regex input form -->
			<div>
				<div class="shared-block">
					<div class="block-headers">
						<div class="header-section">
							<h3>Type</h3>
						</div>
						<div class="header-section">
							<h3>Regular Expression</h3>
						</div>
					</div>
					<div class="input-rows">
						{#each showDefault ? editableDefaultRows : userInputRows as row, i}
							<div class="input-row">
								<div class="input-block">
									<input
										type="text"
										bind:value={row.type}
										    on:input={() => {
											hasLocalChanges = true;
											isSubmitted = false;
											show_tokens = false;
											tokens = [];
										}}
										placeholder="Enter type..."
										class:error={row.error}
									/>
								</div>
								<div class="input-block">
									<input
										type="text"
										bind:value={row.regex}
										    on:input={() => {
											hasLocalChanges = true;
											isSubmitted = false;
											show_tokens = false;
											tokens = [];
										}}
										placeholder="Enter regex pattern..."
										class:error={row.error}
									/>
								</div>
								{#if row.error}
									<div class="error-message">{row.error}</div>
								{/if}
							</div>
						{/each}
					</div>
					{#if (showDefault ? editableDefaultRows[editableDefaultRows.length - 1] : userInputRows[userInputRows.length - 1]).type && 
                        (showDefault ? editableDefaultRows[editableDefaultRows.length - 1] : userInputRows[userInputRows.length - 1]).regex}
                        <button class="add-rule-btn" on:click={addNewRow}>+ Add New Rule</button>
                    {/if}
				</div>
				{#if formError}
					<div class="form-error">{formError}</div>
				{/if}
				<div class="button-stack">
					<!-- Submit button -->
					<button 
						class="submit-button" 
						on:click={handleSubmit}
						disabled={isSubmittingRules}
					>
						<div class="button-content">
							{#if isSubmittingRules}
								<div class="loading-spinner"></div>
								Submitting...
							{:else}
								Submit
							{/if}
						</div>
					</button>

					<!-- FIX: Always show REGEX-specific buttons but disable them until submitted -->
					<div class="regex-action-buttons">
						<button 
							class="generate-button" 
							class:disabled={!isSubmitted || isGeneratingTokens}
							disabled={!isSubmitted || isGeneratingTokens}
							on:click={generateTokens}
							title={isSubmitted ? "Generate tokens from submitted rules" : "Submit regex rules first"}
						>
							<div class="button-content">
								{#if isGeneratingTokens}
									<div class="loading-spinner"></div>
									Generating...
								{:else}
									Generate Tokens
								{/if}
							</div>
						</button>
						<button
							class="generate-button"
							class:disabled={!isSubmitted}
							disabled={!isSubmitted}
							on:click={handleRegexToNFA}
							title={isSubmitted ? "Convert Regular Expression to a NFA" : "Submit regex rules first"}
						>NFA</button>
						<button
							class="generate-button"
							class:disabled={!isSubmitted}
							disabled={!isSubmitted}
							on:click={handleRegexToDFA}
							title={isSubmitted ? "Convert Regular Expression to a DFA" : "Submit regex rules first"}
						>DFA</button>
					</div>
				</div>
			</div>
		{/if}
	{:else if selectedType === 'AUTOMATA'}
		{#if showAutomataVisOnly}
			<!-- Show only DFA visualization or RE display with back button -->
			{#if automataDisplay === 'DFA' && showDfaVis && regexDfa}
				<div class="automata-container pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">DFA Visualization (from Automata)</span>
					</div>
					<button on:click={toggleExpand} class="expand-btn" title="Expand view">
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
							<path d="M1.5 1a.5.5 0 0 0-.5.5v4a.5.5 0 0 1-1 0v-4A1.5 1.5 0 0 1 1.5 0h4a.5.5 0 0 1 0 1zM10 .5a.5.5 0 0 1 .5-.5h4A1.5 1.5 0 0 1 16 1.5v4a.5.5 0 0 1-1 0v-4a.5.5 0 0 0-.5-.5h-4a.5.5 0 0 1-.5-.5M.5 10a.5.5 0 0 1 .5.5v4a.5.5 0 0 0 .5.5h4a.5.5 0 0 1 0 1h-4A1.5 1.5 0 0 1 0 14.5v-4a.5.5 0 0 1 .5-.5m15 0a.5.5 0 0 1 .5.5v4a1.5 1.5 0 0 1-1.5 1.5h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 0 .5-.5v-4a.5.5 0 0 1 .5-.5"/>
						</svg>
					</button>
					<div bind:this={dfaContainer} class="vis-graph-area"></div>
				</div>
				<button
					class="submit-button"
					style="align-self: flex-start; margin-top: 1.5rem;"
					on:click={handleBackFromAutomataVis}
				>
					← Back
				</button>
			{:else if automataDisplay === 'RE' && showRegexOutput && regexRules.length > 0}
				<div class="regex-display-container pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">Generated Regular Expressions (from Automata)</span>
					</div>
					<table class="regex-table">
						<thead>
							<tr>
								<th>Type</th>
								<th>Regular Expression</th>
							</tr>
						</thead>
						<tbody>
							{#each regexRules as rule}
								<tr>
									<td class="regex-type">{rule.token_type || rule.Type || rule.type || '-'}</td>
									<td class="regex-pattern"><code>{rule.regex || rule.Regex || '-'}</code></td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
				<button
					class="submit-button"
					style="align-self: flex-start; margin-top: 1.5rem;"
					on:click={handleBackFromAutomataVis}
				>
					← Back
				</button>
			{/if}
		{:else}
			<!-- Show automata input form -->
			<div class="automaton-section">
				<div class="automaton-left">
					<label>
						States:
						<input 
							class="automaton-input" 
							bind:value={states} 
							on:input={() => {
							hasLocalChanges = true;
							isSubmitted = false;
							showGenerateButton = false;
						}}
							placeholder="e.g. q0,q1,q2" 
						/>
					</label>
					<label>
						Start State:
						<input 
							class="automaton-input" 
							bind:value={startState} 
							on:input={() => {
								hasLocalChanges = true;
								isSubmitted = false;
								showGenerateButton = false;
							}}
							placeholder="e.g. q0" 
						/>
					</label>
					<label>
						Accepted States:
						<input
							class="automaton-input"
							bind:value={acceptedStates}
							on:input={() => {
								hasLocalChanges = true;
								isSubmitted = false;
								showGenerateButton = false;
							}}
							placeholder="e.g. q2->int, q1->string"
						/>
					</label>
				</div>
				<div class="automaton-right">
					<label>
						Transitions:
						<textarea
							class="automaton-input automaton-transitions"
							bind:value={transitions}
							on:input={() => {
								hasLocalChanges = true;
								isSubmitted = false;
								showGenerateButton = false;
							}}
							placeholder="e.g. q0,a->q1&#10;q1,b->q2"
						></textarea>
					</label>
				</div>
				
				<!-- FIX: Keep automata buttons ONLY in AUTOMATA section -->
				<div class="automata-action-row" style="grid-column: span 2;">
					<button
						class="action-btn"
						type="button"
						on:click={handleShowDfa}
					>Show Automata</button>
					<button
						class="action-btn"
						type="button"
						disabled={isGeneratingTokens}
						on:click={() => {
							handleTokenisation();
							automataDisplay = null;
						}}
					>
						<div class="button-content">
							{#if isGeneratingTokens}
								<div class="loading-spinner"></div>
								Tokenising...
							{:else}
								Tokenisation
							{/if}
						</div>
					</button>
					<button
						class="action-btn"
						type="button"
						on:click={handleConvertToRegex}
						title="Convert to Regular Expression"
					>RE</button>
				</div>
			</div>
		{/if}

		<!-- Remove the old conditional displays that were showing below the input form -->
		<!-- These blocks are no longer needed since we're using showAutomataVisOnly -->
{/if}
</div>


{#if isExpanded}
	<div 
		class="modal-backdrop" 
		on:click={toggleExpand} 
		on:keydown={(e) => e.key === 'Escape' && toggleExpand()}
		role="presentation"
		tabindex="-1"
		transition:fade={{ duration: 200 }}
	>
		<div
			class="modal-content"
			on:click|stopPropagation
			on:keydown|stopPropagation
			role="dialog"
			aria-modal="true"
			aria-labelledby="modal-title"
			tabindex="0"
			transition:scale={{ duration: 250, start: 0.95 }}
			on:introend={fitGraphToModal}
		>
			<div class="modal-header">
				<div class="header-left">
					<h3 id="modal-title">Expanded Automaton View</h3>
				</div>
				<div class="header-center">
					<div class="zoom-controls">
						<button 
							on:click={zoomOut} 
							class="zoom-btn" 
							title="Zoom Out"
							aria-label="Zoom out"
						>
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<circle cx="11" cy="11" r="8"/>
								<path d="m21 21-4.35-4.35"/>
								<line x1="8" y1="11" x2="14" y2="11"/>
							</svg>
						</button>
						<button 
							on:click={resetZoom} 
							class="zoom-btn reset-btn" 
							title="Reset Zoom"
							aria-label="Reset zoom to fit"
						>
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
								<circle cx="12" cy="12" r="3"/>
							</svg>
						</button>
						<button 
							on:click={zoomIn} 
							class="zoom-btn" 
							title="Zoom In"
							aria-label="Zoom in"
						>
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<circle cx="11" cy="11" r="8"/>
								<path d="m21 21-4.35-4.35"/>
								<line x1="11" y1="8" x2="11" y2="14"/>
								<line x1="8" y1="11" x2="14" y2="11"/>
							</svg>
						</button>
					</div>
				</div>
				<div class="header-right">
					<button on:click={toggleExpand} class="modal-close-btn" aria-label="Close expanded view">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="24"
							height="24"
							fill="currentColor"
							viewBox="0 0 16"
						>
							<path
								d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8z"
							/>
						</svg>
					</button>
				</div>
			</div>
			<div class="modal-body" bind:this={expandedVisContainer}>

			</div>
		</div>
	</div>
{/if}

<style>
	.phase-inspector {
		flex: 1.2;
		padding: 0.2rem 2rem 2rem 2rem;
		background: #fff;
		transition: background-color 0.3s ease;
	}

	.source-code-header {
			color: #444;
		transition: color 0.3s ease;
	}

	.lexor-heading-h1 {
		transition: color 0.3s ease;
	}

	.source-code-section {
		margin-bottom: 2rem;
	}

	.instructions-section {
		margin: 1.5rem 0 2rem 0;
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

	.source-display {
		background: #f5f5f5;
		padding: 1rem;
		border-radius: 4px;
		max-height: 260px;
		overflow: auto;
		font-family: monospace;
		white-space: pre-wrap;
		margin: 0;
		transition: background-color 0.3s ease, color 0.3s ease;
	}

	.shared-block {
		background: #f5f5f5;
		padding: 1.2rem;
		border-radius: 8px;
		position: relative;
		transition: background-color 0.3s ease, box-shadow 0.3s ease;
	}

	.lexor-heading {
		justify-items: center;
	}

	.block-headers {
		display: flex;
		gap: 2rem;
		margin-bottom: 1rem;
		padding-bottom: 0.8rem;
		border-bottom: 1px solid #e0e0e0;
		transition: border-color 0.3s ease;
	}

	.header-section {
		flex: 1;
	}

	.header-section h3 {
		margin: 0;
		color: #041a47;
		font-size: 1rem;
		font-weight: 600;
		transition: color 0.3s ease;
	}

	.input-rows {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding-top: 0.5rem;
	}

	.input-row {
		display: flex;
		gap: 2rem;
		position: relative;
	}

	.input-block {
		flex: 1;
		}

	.input-block input {
		width: 100%;
		padding: 0.8rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.9rem;
		background: white;
		color: #333;
		box-sizing: border-box;
		transition: border-color 0.2s, box-shadow 0.2s, background-color 0.3s ease, color 0.3s ease;
	}

	.input-block input:focus {
		outline: none;
		border-color: #041a47;
		box-shadow: 0 0 0 2px rgba(4, 26, 71, 0.1);
	}

	.input-block input.error {
		border-color: #dc3545;
	}

	.input-block input.error:focus {
		box-shadow: 0 0 0 2px rgba(220, 53, 69, 0.25);
	}

	.error-message {
		color: #dc3545;
		font-size: 0.8rem;
		margin-top: 0.25rem;
		position: absolute;
		bottom: -1.2rem;
	}

	.form-error {
		color: #dc3545;
		text-align: center;
		margin: 0.5rem 0;
		font-size: 0.9rem;
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
		box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
		margin: 0 auto;
		display: block;
		transition: background-color 0.2s, color 0.2s, transform 0.2s;
	}

	.submit-button:hover {
		background: #a8bdd1;
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
		transform: translateY(-2px);
	}

	.button-stack {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		margin-top: 1.5rem;
	}

	.generate-button {
		padding: 0.6rem 1.5rem;
		background: #BED2E6;
		color: 000000;
		border: none;
		border-radius: 6px;
		font-size: 0.9rem;
		font-weight: 500;
			cursor: pointer;
		box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
		transition: background-color 0.2s ease, transform 0.2s;
	}

	.generate-button:hover:not(:disabled) {
		background: #a8bdd1;
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
		transform: translateY(-2px);
	}

	.generate-button:disabled,
	.generate-button.disabled {
		background: #d6d8db;
		color: #6c757d;
		cursor: not-allowed;
		opacity: 0.6;
		transform: none;
	}

	.generate-button:disabled:hover,
	.generate-button.disabled:hover {
		background: #d6d8db;
		color: #6c757d;
		transform: none;
	}

	.automaton-btn-row {
		display: flex;
		gap: 0.7rem;
		margin: 2rem 0 1.5rem 0;
		align-items: center;
		justify-content: space-between;
	}

	.tabs-and-example-group {
		display: flex;
		gap: 0.7rem;
		align-items: center;
	}

	.clear-button-container {
		display: flex;
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
        height: 2.3rem;
        width: 2.3rem;
        border-radius: 50%;
    }

    .clear-toggle-btn:hover,
    .clear-toggle-btn:focus {
        background: #fff5f5;
        border-color: #ef4444;
    }

	.automaton-btn {
		padding: 0.4rem 1rem;
		border: 2px solid #e5e7eb;
		border-radius: 1.2rem;
		background: white;
		color: #041a47;
		font-weight: 500;
		font-size: 0.95rem;
		cursor: pointer;
		transition: all 0.2s ease;
		outline: none;
	}

	.automaton-btn.selected {
		border-color: #041a47;
		background: #e6edfa;
		color: #041a47;
	}

	.automaton-btn:not(.selected):hover {
		border-color: #7da2e3;
		background: #f5f8fd;
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
		max-width: 140px;
		justify-content: center;
		margin-right: 0.5rem;
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



	.automaton-section {
		margin-top: 0;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.2rem 1.5rem;
		max-width: 600px;
	}

	.automaton-section label {
		font-weight: 500;
		color: #041a47;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-size: 1.05rem;
		transition: color 0.3s ease;
		margin-top: 0.5rem;
	}

	.automaton-input {
		padding: 0.7rem 1rem;
		border: 1.5px solid #b6c6e3;
		border-radius: 0.7rem;
		font-size: 1rem;
		background: #f8fafc;
		color: #041a47;
		transition: all 0.2s ease;
		width: 100%;
		box-sizing: border-box;
	}

	.automaton-input:focus {
		border-color: #041a47;
		outline: none;
		background: #e6edfa;
	}

	.automaton-transitions {
		min-height: 13rem;
		font-family: 'Fira Mono', monospace;
		resize: vertical;
	}

	.automata-action-row {
		grid-column: span 2;
		display: flex;
		gap: 1.25rem;
		margin-top: 1.75rem;
		justify-content: flex-start;
		align-items: center;
		flex-wrap: wrap;
	}

	.action-btn {
		font-weight: 500;
		cursor: pointer;
		transition: all 0.25s ease;
		border: none;
		outline: none;
		padding: 0.625rem 1.25rem;
		border-radius: 0.75rem;
		font-size: 0.9375rem;
		background: #BED2E6;
		color: #000000;
		box-shadow: 0 2px 8px rgba(100, 116, 139, 0.2);
	}

	.action-btn:hover,
	.action-btn:focus {
		background: #a8bdd1;
		box-shadow: 0 4px 12px rgba(100, 116, 139, 0.3);
		transform: translateY(-2px);
	}

	.regex-action-buttons {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.pretty-vis-box,
	.regex-display-container {
		background: #f8fafc;
		border-radius: 14px;
		box-shadow: 0 2px 12px rgba(30, 64, 175, 0.07), 0 1.5px 6px rgba(0, 0, 0, 0.04);
		padding: 1.5rem;
		margin-top: 2.2rem;
		margin-bottom: 1.5rem;
		max-width: 750px;
		width: 100%;
		transition: background-color 0.3s ease, box-shadow 0.3s ease;
		position: relative; 
	}

	.vis-heading {
		margin-bottom: 1.1rem;
		font-size: 1.25rem;
		font-weight: 600;
		color: #1e40af;
		transition: color 0.3s ease;
	}

	.vis-title {
		font-size: 1.18rem;
		font-weight: 600;
		color: #1e40af;
		transition: color 0.3s ease;
	}

	.vis-graph-area {
		width: 100%;
		max-width: 450px;
		height: 350px;
		border: 1.5px solid #e0e7ef;
		border-radius: 10px;
		background: #fff;
		margin: 0 auto;
		box-shadow: 0 1px 4px rgba(30, 64, 175, 0.04);
		transition: background-color 0.3s ease, border-color 0.3s ease;
	}

	.regex-table {
		width: 100%;
		border-collapse: collapse;
		margin-top: 1rem;
	}

	.regex-table th,
	.regex-table td {
		padding: 0.7rem 1rem;
		border-bottom: 1px solid #e0e7ef;
		text-align: left;
		transition: border-color 0.3s ease;
	}

	.regex-table th {
		background: #e6edfa;
		color: #041a47;
		font-size: 1.05rem;
		transition: background-color 0.3s ease, color 0.3s ease;
	}

	.regex-type {
		font-weight: 600;
		color: #1e40af;
		letter-spacing: 0.03em;
		transition: color 0.3s ease;
	}

	.regex-pattern code {
		background: #e0e7ff;
		color: #0a2540;
		padding: 0.2em 0.5em;
		border-radius: 5px;
		font-family: 'Fira Mono', monospace;
		font-size: 1.01em;
		word-break: break-all;
		display: inline-block;
		transition: background-color 0.3s ease, color 0.3s ease;
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
	.submit-button,
    .generate-button,
    .action-btn {
        /* ... existing button styles ... */
        position: relative;
        overflow: hidden;
    }

    /* FIX: Add button content wrapper */
    .button-content {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    /* FIX: Add loading spinner styles */
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

	.submit-button:disabled,
    .generate-button:disabled,
    .action-btn:disabled {
        background: #d6d8db;
        color: #6c757d;
        cursor: wait;
        opacity: 0.8;
        transform: none;
		pointer-events: none;
    }

    .submit-button:disabled:hover,
    .generate-button:disabled:hover,
    .action-btn:disabled:hover {
        background: #d6d8db;
        color: #6c757d;
        transform: none;
    }

    .submit-button:disabled .loading-spinner,
    .generate-button:disabled .loading-spinner,
    .action-btn:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

	 .generate-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

	/* --- Dark Mode Styles --- */
	:global(html.dark-mode) .phase-inspector {
		background: #1a2a4a;
	}

	:global(html.dark-mode) .source-code-header,
	:global(html.dark-mode) .lexor-heading-h1,
	:global(html.dark-mode) .automaton-section label {
		color: #e2e8f0;
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

	:global(html.dark-mode) .source-display,
	:global(html.dark-mode) .shared-block,
	:global(html.dark-mode) .pretty-vis-box,
	:global(html.dark-mode) .regex-display-container {
		background: #2d3748;
	}

	:global(html.dark-mode) .source-display {
		color: #e2e8f0;
	}

	:global(html.dark-mode) .generate-button:disabled,
    :global(html.dark-mode) .generate-button.disabled {
        background: #495057;
        color: #6c757d;
        cursor: not-allowed;
        opacity: 0.6;
        transform: none;
    }

    :global(html.dark-mode) .generate-button:disabled:hover,
    :global(html.dark-mode) .generate-button.disabled:hover {
        background: #495057;
        color: #6c757d;
        transform: none;
    }

    :global(html.dark-mode) .generate-button:disabled .loading-spinner {
        border-top-color: #6c757d;
    }

	:global(html.dark-mode) .block-headers {
		border-bottom-color: #4a5568;
	}

	:global(html.dark-mode) .header-section h3 {
		color: #90cdf4;
	}

	:global(html.dark-mode) .input-block input,
	:global(html.dark-mode) .automaton-input {
		background: #2d3748;
		color: #e2e8f0;
		border-color: #4a5568;
	}

	:global(html.dark-mode) .input-block input::placeholder,
	:global(html.dark-mode) .automaton-input::placeholder {
		color: #a0aec0;
	}

	:global(html.dark-mode) .input-block input:focus,
	:global(html.dark-mode) .automaton-input:focus {
		border-color: #60a5fa;
		background-color: #2d3748;
		box-shadow: 0 0 0 2px rgba(99, 179, 237, 0.2);
	}

	:global(html.dark-mode) .submit-button {
		background: #001A6E;
		color: #ffffff;
		font-weight: 600;
	}

	:global(html.dark-mode) .submit-button:hover {
		background: #002a8e;
		box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
	}

	:global(html.dark-mode) .generate-button,
	:global(html.dark-mode) .action-btn {
		background: #001A6E;
		color: #ffffff;
	}

	:global(html.dark-mode) .generate-button:hover:not(:disabled),
	:global(html.dark-mode) .action-btn:hover:not(:disabled) {
		background: #002a8e;
		box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
	}

	:global(html.dark-mode) .generate-button:disabled,
	:global(html.dark-mode) .generate-button.disabled {
		background: #495057;
		color: #6c757d;
		cursor: not-allowed;
		opacity: 0.6;
		transform: none;
	}

	:global(html.dark-mode) .generate-button:disabled:hover,
	:global(html.dark-mode) .generate-button.disabled:hover {
		background: #495057;
		color: #6c757d;
		transform: none;
	}

	:global(html.dark-mode) .automaton-btn {
		background: transparent;
		color: #cbd5e1;
		border-color: #4a5568;
	}

	:global(html.dark-mode) .automaton-btn.selected {
		background: #2d3748;
		color: #ffffff;
		border-color: #63b3ed;
	}

	:global(html.dark-mode) .automaton-btn:not(.selected):hover {
		background: #002a8e;
        box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
	}

	:global(html.dark-mode) .example-btn {
		background: #001A6E;
        color: #ffffff;
	}
	:global(html.dark-mode) .example-btn:hover {
        background: #002a8e;
        box-shadow: 0 4px 12px rgba(0, 26, 110, 0.3);
    }

	:global(html.dark-mode) .clear-toggle-btn {
        background: transparent;
		color: #ef4444;
		border-color: #4a5568;
    }

    :global(html.dark-mode) .clear-toggle-btn:hover {
        background: #7f1d1d;
		border-color: #ef4444;
    }

	:global(html.dark-mode) .vis-heading,
	:global(html.dark-mode) .vis-title {
		color: #90cdf4;
	}

	:global(html.dark-mode) .vis-graph-area {
		background: #1a202c;
		border-color: #4a5568;
	}

	:global(html.dark-mode) .regex-table th {
		background: #1a202c;
		color: #e2e8f0;
	}

	:global(html.dark-mode) .regex-table td {
		border-color: #4a5568;
	}

	:global(html.dark-mode) .regex-type {
		color: #90cdf4;
	}

	:global(html.dark-mode) .regex-pattern code {
		background: #4a5568;
		color: #e2e8f0;
	}

	:global(html.dark-mode) .add-rule-btn {
        border-color: #60a5fa;
        color: #60a5fa;
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

	/* Back button after NFA/DFA visualization */
	button[style*="align-self: flex-start"] {
		background: #BED2E6 !important; 
		color: #000000 !important; 
	}

	:global(html.dark-mode) button[style*="align-self: flex-start"] {
		background: #001A6E !important; 
		color: #ffffff !important; 
	}


	.expand-btn {
		background: none;
		border: none;
		cursor: pointer;
		color: #1e40af;
		padding: 0.25rem;
		border-radius: 50%;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		transition: background-color 0.2s;
		position: absolute;
		top: 1.5rem;
		right: 1.5rem;
	}
	.expand-btn:hover {
		background-color: rgba(0, 0, 0, 0.1);
	}
	:global(html.dark-mode) .expand-btn {
		color: #90cdf4;
	}
	:global(html.dark-mode) .expand-btn:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.6);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
	}

	.modal-content {
		background: #fff;
		padding: 1.5rem;
		border-radius: 12px;
		box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
		width: 80vw;
		height: 80vh;
		display: flex;
		flex-direction: column;
	}
	:global(html.dark-mode) .modal-content {
		background: #2d3748;
	}

	.modal-header {
		display: flex;
		align-items: center;
		border-bottom: 1px solid #e0e0e0;
		padding-bottom: 1rem;
		margin-bottom: 1rem;
	}
	:global(html.dark-mode) .modal-header {
		border-bottom-color: #4a5568;
	}

	.header-left,
	.header-right {
		flex: 1;
		display: flex;
		align-items: center;
	}

	.header-left {
		justify-content: flex-start;
	}

	.header-center {
		flex: 0 0 auto;
		display: flex;
		justify-content: center;
	}

	.header-right {
		justify-content: flex-end;
	}

	.zoom-controls {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		background: #f8f9fa;
		border: 1px solid #dee2e6;
		border-radius: 8px;
		padding: 0.25rem;
	}
	:global(html.dark-mode) .zoom-controls {
		background: #4a5568;
		border-color: #6b7280;
	}

	.zoom-btn {
		background: none;
		border: none;
		cursor: pointer;
		color: #374151;
		padding: 0.5rem;
		border-radius: 4px;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		min-width: 32px;
		min-height: 32px;
	}
	.zoom-btn:hover {
		background-color: #e9ecef;
		color: #1f2937;
	}
	.zoom-btn:active {
		transform: scale(0.95);
	}
	:global(html.dark-mode) .zoom-btn {
		color: #f9fafb;
	}
	:global(html.dark-mode) .zoom-btn:hover {
		background-color: #5a6578;
		color: #ffffff;
	}

	.reset-btn {
		border-left: 1px solid #dee2e6;
		border-right: 1px solid #dee2e6;
		margin: 0 0.25rem;
		border-radius: 4px;
	}
	:global(html.dark-mode) .reset-btn {
		border-left-color: #6b7280;
		border-right-color: #6b7280;
	}

	.modal-header h3 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: #041a47;
	}
	:global(html.dark-mode) .modal-header h3 {
		color: #e2e8f0;
	}

	.modal-close-btn {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		color: #666;
		padding: 0.5rem;
		line-height: 1;
		border-radius: 50%;
		transition: background-color 0.2s;
	}
	.modal-close-btn:hover {
		background-color: #f0f0f0;
	}
	:global(html.dark-mode) .modal-close-btn {
		color: #a0aec0;
	}
	:global(html.dark-mode) .modal-close-btn:hover {
		background-color: #4a5568;
	}

	.modal-body {
		flex-grow: 1;
		width: 100%;
		height: 100%;
		overflow: hidden;
		position: relative;
	}

</style>

<script lang="ts">
	import type { Token } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { onMount } from 'svelte';
	import { DataSet, Network } from 'vis-network/standalone';

	export let source_code = '';
	export let onGenerateTokens: (data: {
		tokens: Token[];
		unexpected_tokens: string[];
	}) => void = () => {};

	let inputRows = [{ type: '', regex: '', error: '' }];
	let userSourceCode = '';
	let userInputRows = [{ type: '', regex: '', error: '' }];
	let formError = '';
	let submissionStatus = { show: false, success: false, message: '' };
	let showGenerateButton = false;
	let showRegexActionButtons = false;

	function addNewRow() {
		userInputRows = [...userInputRows, { type: '', regex: '', error: '' }];
	}

	function validateRegex(pattern: string): boolean {
		try {
			new RegExp(pattern);
			return true;
		} catch (e) {
			return false;
		}
	}

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
			submissionStatus = {
				show: true,
				success: false,
				message: 'Please fill in both Type and Regular Expression'
			};
			return;
		}

		if (showDefault) {
			editableDefaultRows = nonEmptyRows;
		} else {
			userInputRows = nonEmptyRows;
		}

		if (hasErrors) {
			submissionStatus = {
				show: true,
				success: false,
				message: 'Please fix the errors before submitting'
			};
			return;
		}

		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		if (selectedType === 'REGEX') {
			const requestData = {
				users_id: user_id,
				pairs: nonEmptyRows.map((row) => ({
					Type: row.type.toUpperCase(),
					Regex: row.regex
				}))
			};
			try {
				const res = await fetch('http://51.21.245.160:8080/api/lexing/rules', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(requestData)
				});
				if (!res.ok) {
					const errorText = await res.text();
					throw new Error(errorText);
				}
				submissionStatus = { show: true, success: true, message: 'Rules stored successfully!' };
				showRegexActionButtons = true;
			} catch (error) {
				AddToast('Failed to save rules', 'error');
			}
			return;
		}

		const requestData = {
			users_id: user_id,
			source_code: showDefault ? DEFAULT_SOURCE_CODE : userSourceCode,
			pairs: nonEmptyRows.map((row) => ({
				Type: row.type.toUpperCase(),
				Regex: row.regex
			}))
		};

		try {
			const storeResponse = await fetch('http://51.21.245.160:8080/api/lexing/code', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(requestData)
			});
			if (!storeResponse.ok) {
				const errorText = await storeResponse.text();
				throw new Error(`Server error (${storeResponse.status}): ${errorText}`);
			}
			submissionStatus = { show: true, success: true, message: 'Code stored successfully!' };
			showGenerateButton = true;
		} catch (error) {
			console.error('Store error:', error);
			AddToast('Cannot connect to server. Please ensure the backend is running.', 'error');
		}

		// Show regex action buttons after successful submit in REGEX mode
		if (selectedType === 'REGEX' && !hasErrors) {
			showRegexActionButtons = true;
		}
	}

	async function generateTokens() {
		try {
			const user_id = localStorage.getItem('user_id');
			if (!user_id) {
				AddToast('User not logged in.', 'error');
				return;
			}
			const requestData = {
				users_id: user_id,
				source_code: source_code,
				pairs: (showDefault ? editableDefaultRows : userInputRows).map((row) => ({
					Type: row.type.toUpperCase(),
					Regex: row.regex
				}))
			};

			const response = await fetch('http://51.21.245.160:8080/api/lexing/lexer', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
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
			submissionStatus = {
				show: true,
				success: true,
				message: data.message || 'Tokens generated successfully!'
			};
		} catch (error) {
			console.error('Generate tokens error:', error);
			AddToast('Error generating tokens', 'error');
		}
	}

	// Helper: Save DFA to backend
	async function saveDfaToBackend() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return false;
		}

		const dfa = {
			states: states
				.split(',')
				.map((s) => s.trim())
				.filter(Boolean),
			transitions: [],
			start_state: startState.trim(),
			accepting_states: parseAcceptedStates(acceptedStates),
			users_id: user_id
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
			const response = await fetch('http://51.21.245.160:8080/api/lexing/dfa', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(dfa)
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Failed to save DFA: ' + errorText, 'error');
				return false;
			}
			return true;
		} catch (error) {
			AddToast('Failed to save DFA: ' + error, 'error');
			return false;
		}
	}

	async function handleTokenisation() {
		const saved = await saveDfaToBackend();
		if (!saved) return;

		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		// Only need to send users_id, backend loads DFA from DB
		const body = { users_id: user_id };

		try {
			const response = await fetch('http://51.21.245.160:8080/api/lexing/dfaToTokens', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
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
			AddToast('Tokenisation successful!', 'success');
		} catch (error) {
			AddToast('Tokenisation failed: ' + error, 'error');
		}
	}

	let previousInputs: typeof userInputRows = [];
	function handleInputChange() {
		showGenerateButton = false;
		showRegexActionButtons = false;
		submissionStatus = { show: false, success: false, message: '' };
	}

	$: {
		const inputsChanged =
			userInputRows.length !== previousInputs.length ||
			userInputRows.some((row, index) => {
				const prevRow = previousInputs[index];
				return !prevRow || row.type !== prevRow.type || row.regex !== prevRow.regex;
			});
		if (inputsChanged) {
			handleInputChange();
			previousInputs = [...userInputRows];
		}
	}

	let selectedType: 'AUTOMATA' | 'REGEX' | null = null;
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
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		// 1. Save DFA to backend (if not already saved)
		const saved = await saveDfaToBackend();
		if (!saved) return;

		try {
			// 2. Convert DFA to Regex
			const dfaToRegexRes = await fetch('http://51.21.245.160:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!dfaToRegexRes.ok) {
				const errorText = await dfaToRegexRes.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}

			// 3. Convert Regex to NFA
			const regexToNfaRes = await fetch('http://51.21.245.160:8080/api/lexing/regexToNFA', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!regexToNfaRes.ok) {
				const errorText = await regexToNfaRes.text();
				AddToast('Regex→NFA failed: ' + errorText, 'error');
				return;
			}
			const nfaData = await regexToNfaRes.json();
			regexNfa = adaptAutomatonForVis(nfaData.nfa);
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
		resetInputs();
		showRegexActionButtons = false;
		if (type === 'REGEX') {
			showDefault = false;
			userInputRows = [{ type: '', regex: '', error: '' }];
			inputRows = [{ type: '', regex: '', error: '' }];
			editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
			formError = '';
			submissionStatus = { show: false, success: false, message: '' };
			showGenerateButton = false;
		}
	}

	function insertDefault() {
		showDefault = true;
		editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
		inputRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));
		states = 'S0, S1, S2, S3, S4, S5, S6, S7, S8';
		startState = 'S0';
		acceptedStates =
			'S3->KEYWORD, S4->IDENTIFIER, S5->ASSIGNMENT, S6->OPERATOR, S7->INTEGER, S8->SEPARATOR';
		transitions =
			'S0,i->S1\nS1,n->S2\nS2,t->S3\nS0,[a-zA-Z_]->S4\nS4,[a-zA-Z_]->S4\nS0,=->S5\nS0,[+-*/%]->S6\nS0,[0-9]->S7\nS7,[0-9]->S7\nS0,;->S8';
	}

	function removeDefault() {
		showDefault = false;
		inputRows = [...userInputRows];
		resetInputs();
	}

	const DEFAULT_INPUT_ROWS = [
		{ type: 'keyword', regex: 'int|str|if', error: '' },
		{ type: 'identifier', regex: '[a-zA-Z]+', error: '' },
		{ type: 'integer', regex: '[0-9]+', error: '' },
		{ type: 'assignment', regex: '=', error: '' },
		{ type: 'operator', regex: '[+\\-*/%]', error: '' },
		{ type: 'separator', regex: ';', error: '' }
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
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		// 1. Save DFA to backend (if not already saved)
		const saved = await saveDfaToBackend();
		if (!saved) return;

		try {
			// 2. Convert DFA to Regex
			const regexRes = await fetch('http://51.21.245.160:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!regexRes.ok) {
				const errorText = await regexRes.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}
			// Optionally, you can use the rules here if you want to display them

			// 3. Convert Regex to DFA
			const dfaRes = await fetch('http://51.21.245.160:8080/api/lexing/regexToDFA', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!dfaRes.ok) {
				const errorText = await dfaRes.text();
				AddToast('Regex→DFA failed: ' + errorText, 'error');
				return;
			}
			const dfaData = await dfaRes.json();
			regexDfa = adaptAutomatonForVis(dfaData.dfa);
			showDfaVis = true;
			showNfaVis = false;
			automataDisplay = 'DFA';
			setTimeout(() => renderRegexAutomatonVis(dfaContainer, regexDfa, true), 0);
			AddToast('DFA generated from Regex and displayed!', 'success');
		} catch (error) {
			AddToast('Failed to generate DFA from Regex: ' + error, 'error');
		}
	}

	let regexRules: { token_type?: string; Type?: string; Regex?: string; regex?: string }[] = [];
	let showRegexOutput = false;

	async function handleConvertToRegex() {
		const saved = await saveDfaToBackend();
		if (!saved) return;

		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		try {
			const response = await fetch('http://51.21.245.160:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}
			const data = await response.json();
			console.log('dfaToRegex response:', JSON.stringify(data, null, 2));
			regexRules = Array.isArray(data.rules) ? data.rules : [];
			showRegexOutput = true;
			AddToast('DFA converted to Regex successfully!', 'success');
		} catch (error) {
			AddToast('DFA→Regex failed: ' + error, 'error');
		}
	}

	let regexNfa = null;
	let regexDfa = null;
	let showRegexNfaVis = false;
	let showRegexDfaVis = false;
	let regexNfaContainer: HTMLElement;
	let regexDfaContainer: HTMLElement;

	async function handleRegexToNFA() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}
		try {
			const response = await fetch('http://51.21.245.160:8080/api/lexing/regexToNFA', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Regex→NFA failed: ' + errorText, 'error');
				return;
			}
			const data = await response.json();
			console.log('NFA from backend:', JSON.stringify(data.nfa, null, 2));
			regexNfa = adaptAutomatonForVis(data.nfa);
			showRegexNfaVis = true;
			showRegexDfaVis = false;
			showRegexVisOnly = true;
			AddToast('Regex converted to NFA!', 'success');
			setTimeout(() => renderRegexAutomatonVis(regexNfaContainer, regexNfa, false), 0);
		} catch (error) {
			AddToast('Regex→NFA failed: ' + error, 'error');
		}
	}

	async function handleRegexToDFA() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}
		try {
			const response = await fetch('http://51.21.245.160:8080/api/lexing/regexToDFA', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!response.ok) {
				const errorText = await response.text();
				AddToast('Regex→DFA failed: ' + errorText, 'error');
				return;
			}
			const data = await response.json();
			console.log('DFA from backend:', JSON.stringify(data.dfa, null, 2));
			regexDfa = adaptAutomatonForVis(data.dfa);
			showRegexDfaVis = true;
			showRegexNfaVis = false;
			showRegexVisOnly = true;
			AddToast('Regex converted to DFA!', 'success');
			setTimeout(() => renderRegexAutomatonVis(regexDfaContainer, regexDfa, true), 0);
		} catch (error) {
			AddToast('Regex→DFA failed: ' + error, 'error');
		}
	}

	function adaptAutomatonForVis(automaton: any) {
		// Convert transitions array to nested object if needed
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

		// Try to infer alphabet if not present
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
		new Network(
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

	function handleBackFromRegexVis() {
		showRegexVisOnly = false;
		showRegexNfaVis = false;
		showRegexDfaVis = false;
	}
</script>

<div class="phase-inspector">
	<div class="source-code-section">
		<div class="lexor-heading">
			<h1 class="lexor-heading-h1">LEXING</h1>
		</div>
		<h3 class="source-code-header">Source Code</h3>
		<pre class="source-display">{source_code || 'no source code available'}</pre>
	</div>

	<div class="automaton-btn-row">
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

		{#if selectedType && !showDefault}
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
		{#if selectedType && showDefault}
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

	{#if selectedType === 'REGEX'}
		{#if showRegexVisOnly}
			<!-- Only show the NFA/DFA display and back button -->
			{#if showRegexNfaVis && regexNfa}
				<div class="automata-container pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">NFA Visualization (from REGEX)</span>
					</div>
					<div bind:this={regexNfaContainer} class="vis-graph-area" />
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
					<div bind:this={regexDfaContainer} class="vis-graph-area" />
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
			<!-- Show the regular expression input, submit, and action buttons as before -->
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
										on:input={handleInputChange}
										placeholder="Enter type..."
										class:error={row.error}
									/>
								</div>
								<div class="input-block">
									<input
										type="text"
										bind:value={row.regex}
										on:input={handleInputChange}
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
					{#if userInputRows[userInputRows.length - 1].type && userInputRows[userInputRows.length - 1].regex}
						<button class="add-button" on:click={addNewRow}>
							<span>+</span>
						</button>
					{/if}
				</div>
				{#if formError}
					<div class="form-error">{formError}</div>
				{/if}
				<div class="button-stack">
					<button class="submit-button" on:click={handleSubmit}> Submit </button>
					{#if showRegexActionButtons}
						<div class="regex-action-buttons">
							<button class="generate-button" on:click={generateTokens}>Generate Tokens</button>
							<button
								class="generate-button"
								on:click={handleRegexToNFA}
								title="Convert Regular Expression to a NFA">NFA</button
							>
							<button
								class="generate-button"
								on:click={handleRegexToDFA}
								title="Convert Regular Expression to a DFA">DFA</button
							>
						</div>
					{/if}
				</div>
				{#if submissionStatus.show}
					<div
						class="status-message"
						class:success={submissionStatus.success === true}
						class:info={submissionStatus.message === 'info'}
					>
						{submissionStatus.message}
					</div>
				{/if}
			</div>
		{/if}
	{:else if selectedType === 'AUTOMATA'}
		<div class="automaton-section">
			<div class="automaton-left">
				<label>
					States:
					<input class="automaton-input" bind:value={states} placeholder="e.g. q0,q1,q2" />
				</label>
				<label>
					Start State:
					<input class="automaton-input" bind:value={startState} placeholder="e.g. q0" />
				</label>
				<label>
					Accepted States:
					<input
						class="automaton-input"
						bind:value={acceptedStates}
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
						placeholder="e.g. q0,a->q1&#10;q1,b->q2"
					></textarea>
				</label>
			</div>
			<div class="automata-action-row" style="grid-column: span 2;">
				<button
					class="action-btn"
					type="button"
					on:click={() => {
						showNfaDiagram();
						automataDisplay = 'NFA';
					}}>Show NFA</button
				>
				<button
					class="action-btn"
					type="button"
					on:click={() => {
						handleShowDfa();
						automataDisplay = 'DFA';
					}}>Show DFA</button
				>
				<button
					class="action-btn"
					type="button"
					on:click={() => {
						handleTokenisation();
						automataDisplay = null;
					}}>Tokenisation</button
				>
				<button
					class="action-btn"
					type="button"
					on:click={() => {
						handleConvertToRegex();
						automataDisplay = 'RE';
					}}
					title="Convert to Regular Expression">RE</button
				>
			</div>
		</div>

		{#if automataDisplay === 'NFA' && showNfaVis}
			<div class="automata-container pretty-vis-box">
				<div class="vis-heading">
					<span class="vis-title">NFA Visualization</span>
				</div>
				<div bind:this={nfaContainer} class="vis-graph-area" />
			</div>
		{:else if automataDisplay === 'DFA' && showDfaVis}
			<div class="automata-container pretty-vis-box">
				<div class="vis-heading">
					<span class="vis-title">DFA Visualization</span>
				</div>
				<div bind:this={dfaContainer} class="vis-graph-area" />
			</div>
		{:else if automataDisplay === 'RE' && showRegexOutput && regexRules.length > 0}
			<div class="regex-display-container pretty-vis-box">
				<div class="vis-heading">
					<span class="vis-title">Generated Regular Expressions</span>
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
		{/if}

		{#if showRegexNfaVis && regexNfa}
			<div class="automata-container pretty-vis-box">
				<div class="vis-heading">
					<span class="vis-title">NFA Visualization (from REGEX)</span>
				</div>
				<div bind:this={regexNfaContainer} class="vis-graph-area" />
			</div>
		{/if}
		{#if showRegexDfaVis && regexDfa}
			<div class="automata-container pretty-vis-box">
				<div class="vis-heading">
					<span class="vis-title">DFA Visualization (from REGEX)</span>
				</div>
				<div bind:this={regexDfaContainer} class="vis-graph-area" />
			</div>
		{/if}
	{/if}
</div>

<style>
	.source-code-header {
		color: #444;
	}
	.phase-inspector {
		flex: 1.2;
		padding-left: 2rem;
		padding-right: 2rem;
		padding-top: 0.2rem;
		padding-bottom: 2rem;
		background: #fff;
	}
	.source-code-section {
		margin-bottom: 2rem;
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
	.shared-block {
		background: #f5f5f5;
		padding: 1.2rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		position: relative;
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
	}
	.header-section {
		flex: 1;
	}
	.header-section h3 {
		margin: 0;
		color: #001a6e;
		font-size: 1rem;
		font-weight: 600;
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
		box-sizing: border-box;
		transition:
			border-color 0.2s,
			box-shadow 0.2s;
	}
	.input-block input:focus {
		outline: none;
		border-color: #001a6e;
		box-shadow: 0 0 0 2px rgba(32, 87, 129, 0.1);
	}
	.input-block input.error {
		border-color: #dc3545;
	}
	.input-block input.error:focus {
		border-color: #dc3545;
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
	.add-button {
		position: absolute;
		right: -16px;
		bottom: -16px;
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: #001a6e;
		color: white;
		border: none;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.2rem;
		transition: background 0.2s;
	}
	.add-button:hover {
		background: #27548a;
	}
	.button-container {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-top: 0.75rem;
	}
	.submit-button {
		padding: 0.6rem 1.5rem;
		background: #001a6e;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
		margin: 0 auto;
		display: block;
		/* Remove transform and margin-top for centering */
	}
	.button-stack {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem; /* Small gap between submit and action buttons */
		margin-top: 0.75rem;
	}
	.generate-button {
		padding: 0.6rem 1.5rem;
		background: #666;
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s ease;
		box-shadow: 0 2px 4px rgba(40, 167, 69, 0.1);
		margin-top: 0.4rem;
		height: 100%;
		line-height: 1.1;
	}
	.submit-button:hover {
		background: #27548a;
	}
	:global(html.dark-mode) .submit-button {
		background: #cccccc;
		color: #041a47;
		transition: transform 0.2s ease;
		margin-top: 1rem;
	}
	.generate-button:hover {
		background: #3d3d3d;
	}
	.status-message {
		text-align: center;
		padding: 0.5rem 1rem;
		margin-top: 0.75rem;
		border-radius: 4px;
		font-size: 0.9rem;
		background: #dc3545;
		color: white;
		opacity: 0;
		animation: fadeInOut 3s ease-in-out;
	}
	.status-message.success {
		background: #28a745;
	}
	.status-message.info {
		background: #0096c7;
	}
	@keyframes fadeInOut {
		0% {
			opacity: 0;
		}
		10% {
			opacity: 1;
		}
		90% {
			opacity: 1;
		}
		100% {
			opacity: 0;
		}
	}
	:global(html.dark-mode) .source-code-header {
		color: #ebeef1;
	}

	:global(html.dark-mode) .lexor-heading-h1 {
		color: #ebeef1;
	}
	:global(html.dark-mode) .source-display {
		color: black;
	}
	.automaton-btn-row {
		display: flex;
		gap: 0.7rem;
		margin: 2rem 0 1.5rem 0;
		align-items: center;
	}
	.automaton-btn {
		padding: 0.4rem 1rem;
		border: 2px solid #e5e7eb;
		border-radius: 1.2rem;
		background: white;
		color: #001a6e;
		font-weight: 500;
		font-size: 0.95rem;
		cursor: pointer;
		transition:
			border-color 0.2s,
			background 0.2s,
			color 0.2s;
		outline: none;
	}
	.automaton-btn.selected {
		border-color: #001a6e;
		background: #e6edfa;
		color: #001a6e;
	}
	.automaton-btn:not(.selected):hover {
		border-color: #7da2e3;
		background: #f5f8fd;
	}
	.default-toggle-btn {
		margin-left: 1.2rem;
		padding: 0.4rem 0.7rem;
		border-radius: 1.2rem;
		border: 2px solid #e5e7eb; /* match unselected automaton-btn */
		background: white; /* match unselected automaton-btn */
		color: #001a6e;
		font-size: 1.2rem;
		cursor: pointer;
		transition:
			background 0.2s,
			border-color 0.2s;
		display: flex;
		align-items: center;
		height: 2.2rem;
		width: 2.2rem;
		justify-content: center;
		position: relative;
	}
	.default-toggle-btn.selected {
		background: #d0e2ff;
		border-color: #003399;
	}
	.default-toggle-btn:hover,
	.default-toggle-btn:focus {
		background: #f5f8fd;
		border-color: #7da2e3;
	}
	.icon {
		font-size: 1.3rem;
		line-height: 1;
		pointer-events: none;
	}
	.automaton-section {
		margin-top: 0rem;
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.2rem 1.5rem;
		max-width: 600px;
	}
	.automaton-section label {
		margin-top: 0.75rem;
		font-weight: 500;
		color: #001a6e;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-size: 1.05rem;
	}
	.automaton-input {
		padding: 0.7rem 1rem;
		border: 1.5px solid #b6c6e3;
		border-radius: 0.7rem;
		font-size: 1rem;
		background: #f8fafc;
		color: #001a6e;
		transition: border-color 0.2s;
		width: 100%;
		box-sizing: border-box;
	}
	.automaton-input:focus {
		border-color: #001a6e;
		outline: none;
		background: #e6edfa;
	}
	.automaton-transitions {
		min-height: 13rem;
		min-width: 100%;
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
		display: inline-flex;
		align-items: center;
		justify-content: center;
		font-family:
			'Inter',
			-apple-system,
			BlinkMacSystemFont,
			sans-serif;
		font-weight: 500;
		text-align: center;
		white-space: nowrap;
		cursor: pointer;
		transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
		user-select: none;
		border: none;
		outline: none;
		position: relative;
		overflow: hidden;
		padding: 0.625rem 1.25rem;
		border-radius: 0.75rem;
		font-size: 0.9375rem;
		line-height: 1.5;
		background: #e0e7ff;
		color: #1e40af;
		box-shadow: 0 1px 2px 0 rgba(30, 64, 175, 0.05);
	}
	.action-btn:hover,
	.action-btn:focus {
		background: #d0d9ff;
		color: #1e3a8a;
		box-shadow:
			0 1px 3px 0 rgba(30, 64, 175, 0.1),
			0 1px 2px 0 rgba(30, 64, 175, 0.06);
		transform: translateY(-1px);
	}
	.action-btn:active {
		transform: translateY(0);
		box-shadow: inset 0 2px 4px 0 rgba(30, 64, 175, 0.06);
	}
	.action-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none !important;
	}
	.action-btn[title] {
		position: relative;
	}
	.action-btn[title]:hover::after {
		content: attr(title);
		position: absolute;
		bottom: 100%;
		left: 50%;
		transform: translateX(-50%);
		background: #1e3a8a;
		color: white;
		padding: 0.5rem 0.75rem;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		white-space: nowrap;
		margin-bottom: 0.5rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}
	.regex-output-section {
		margin-top: 2rem;
		padding: 1.5rem;
		border: 1px solid #d1e7dd;
		border-radius: 0.5rem;
		background: #f8f9fa;
	}
	.regex-output-section h3 {
		margin: 0 0 1rem 0;
		color: #0f5132;
		font-size: 1.1rem;
		font-weight: 500;
	}
	.regex-rules-container {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
	}
	.regex-rule {
		display: flex;
		justify-content: space-between;
		padding: 0.8rem;
		border: 1px solid #cfe2ff;
		border-radius: 0.5rem;
		background: #f0f9ff;
	}
	.token-type {
		font-weight: 500;
		color: #084298;
	}
	.regex-pattern {
		font-family: 'Fira Mono', monospace;
		color: #333;
	}
	.regex-display-container {
		margin-top: 2rem;
		background: #f8fafc;
		border-radius: 10px;
		padding: 1.5rem;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.07);
		max-width: 700px;
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
	}
	.regex-table th {
		background: #e6edfa;
		color: #001a6e;
		font-size: 1.05rem;
	}
	.regex-type {
		font-weight: 600;
		color: #1e40af;
		letter-spacing: 0.03em;
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
	}

	.regex-action-buttons {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-top: 0.25rem; /* Minimal space below submit */
	}

	.pretty-vis-box {
		background: #f8fafc;
		border-radius: 14px;
		box-shadow:
			0 2px 12px rgba(30, 64, 175, 0.07),
			0 1.5px 6px rgba(0, 0, 0, 0.04);
		padding: 1.5rem 1.5rem 1.2rem 1.5rem;
		margin-top: 2.2rem;
		margin-bottom: 1.5rem;
		max-width: 750px;
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	.vis-heading {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		margin-bottom: 1.1rem;
		font-size: 1.25rem;
		font-weight: 600;
		color: #1e40af;
		letter-spacing: 0.01em;
	}

	.vis-title {
		font-size: 1.18rem;
		font-weight: 600;
		color: #1e40af;
	}
	.vis-graph-area {
		width: 450px;
		height: 350px;
		border: 1.5px solid #e0e7ef;
		border-radius: 10px;
		background: #fff;
		margin: 0 auto;
		box-shadow: 0 1px 4px rgba(30, 64, 175, 0.04);
	}
</style>

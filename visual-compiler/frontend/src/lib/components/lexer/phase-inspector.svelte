<script lang="ts">
	import type { Token } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { onMount } from 'svelte';
	import { DataSet, Network, type Node, type Edge, type Options, type IdType } from 'vis-network/standalone';

	export let source_code = '';
	export let onGenerateTokens: (data: {
		tokens: Token[];
		unexpected_tokens: string[];
	}) => void = () => {};

	const DEFAULT_INPUT_ROWS = [
		{ type: 'keyword', regex: 'int|str|if', error: '' },
		{ type: 'identifier', regex: '[a-zA-Z]+', error: '' },
		{ type: 'integer', regex: '[0-9]+', error: '' },
		{ type: 'assignment', regex: '=', error: '' },
		{ type: 'operator', regex: '[+\\-*/%]', error: '' },
		{ type: 'separator', regex: ';', error: '' }
	];

	const DEFAULT_SOURCE_CODE = 'int blue = 13 + 22 ;';

	let inputRows = [{ type: '', regex: '', error: '' }];
	let userSourceCode = '';
	let userInputRows = [{ type: '', regex: '', error: '' }];
	let formError = '';
	let submissionStatus = { show: false, success: false, message: '' };
	let showGenerateButton = false;
	let showRegexActionButtons = false;
	let editableDefaultRows = DEFAULT_INPUT_ROWS.map((row) => ({ ...row }));

	// Define a type for our automaton structure for better type safety
	interface Automaton {
		states: string[];
		startState: string;
		acceptedStates: string[];
		transitions: Record<string, Record<string, string[] | string>>;
		alphabet: string[];
	}

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
			submissionStatus = { show: true, success: false, message: 'Please fill in both Type and Regular Expression' };
			return;
		}

		if (showDefault) {
			editableDefaultRows = nonEmptyRows;
		} else {
			userInputRows = nonEmptyRows;
		}

		if (hasErrors) {
			submissionStatus = { show: true, success: false, message: 'Please fix the errors before submitting' };
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
				const res = await fetch('http://localhost:8080/api/lexing/rules', {
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

		// This block seems to be for a different purpose than the REGEX one,
		// so it is kept separate. It handles storing source code along with regex pairs.
		const requestData = {
			users_id: user_id,
			source_code: showDefault ? DEFAULT_SOURCE_CODE : userSourceCode,
			pairs: nonEmptyRows.map((row) => ({
				Type: row.type.toUpperCase(),
				Regex: row.regex
			}))
		};

		try {
			const storeResponse = await fetch('http://localhost:8080/api/lexing/code', {
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

			const response = await fetch('http://localhost:8080/api/lexing/lexer', {
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
			submissionStatus = { show: true, success: true, message: data.message || 'Tokens generated successfully!' };
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

		const dfaTransitions: { from: string; label: string; to: string }[] = [];
		const transitionLines = transitions.split('\n').map((line) => line.trim()).filter(Boolean);
		for (const line of transitionLines) {
			const match = line.match(/^(.+),(.+)->(.+)$/);
			if (match) {
				const [_, from, label, to] = match;
				dfaTransitions.push({ from: from.trim(), label: label.trim(), to: to.trim() });
			}
		}

		const dfa = {
			states: states.split(',').map((s) => s.trim()).filter(Boolean),
			transitions: dfaTransitions,
			start_state: startState.trim(),
			accepting_states: parseAcceptedStates(acceptedStates),
			users_id: user_id
		};

		// Log the DFA JSON being sent
		console.log('DFA sent to backend:', JSON.stringify(dfa, null, 2));

		try {
			const response = await fetch('http://localhost:8080/api/lexing/dfa', {
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
			AddToast('Failed to save DFA: ' + String(error), 'error');
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
			const response = await fetch('http://localhost:8080/api/lexing/dfaToTokens', {
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
			AddToast('Tokenisation failed: ' + String(error), 'error');
		}
	}

	let previousInputs: typeof userInputRows = [];
	function handleInputChange() {
		showGenerateButton = false;
		showRegexActionButtons = false;
		submissionStatus = { show: false, success: false, message: '' };
	}

	$: {
		const currentRows = showDefault ? editableDefaultRows : userInputRows;
		const inputsChanged =
			currentRows.length !== previousInputs.length ||
			currentRows.some((row, index) => {
				const prevRow = previousInputs[index];
				return !prevRow || row.type !== prevRow.type || row.regex !== prevRow.regex;
			});
		if (inputsChanged) {
			handleInputChange();
			previousInputs = JSON.parse(JSON.stringify(currentRows));
		}
	}

	let selectedType: 'AUTOMATA' | 'REGEX' | null = null;
	let showDefault = false;
	let states = '';
	let startState = '';
	let acceptedStates = '';
	let transitions = '';

	function parseAutomaton(): Automaton {
		const stateList = states.split(',').map((s) => s.trim()).filter(Boolean);
		const start = startState.trim();
		const accepted = parseAcceptedStates(acceptedStates).map((a) => a.state);
		const transitionLines = transitions.split('\n').map((line) => line.trim()).filter(Boolean);
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
			acceptedStates: accepted,
			transitions: transitionObj,
			alphabet: Array.from(alphabetSet)
		};
	}

	function nfaToDfa(nfa: Automaton): Automaton {
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
					const stateTransitions = transitions[state];
					if (stateTransitions && stateTransitions[symbol]) {
						nextSet = nextSet.concat(stateTransitions[symbol]);
					}
				}
				nextSet = Array.from(new Set(nextSet));
				if (nextSet.length === 0) continue;
				const nextName = stateSetToName(nextSet);
				dfaTransitions[name][symbol] = nextName;
				if (!dfaStates.includes(nextName) && !unmarked.some((s) => stateSetToName(s) === nextName)) {
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

	async function renderNfaVis() {
		const nfa = parseAutomaton();
		const nodeIds: Record<string, string> = {};
		nfa.states.forEach((state) => {
			nodeIds[state] = safeStateId(state);
		});
		const nodes = new DataSet<Node>(
			nfa.states.map((state) => ({
				id: nodeIds[state],
				label: state,
				shape: 'circle',
				color: nfa.acceptedStates.includes(state) ? '#D2FFD2' : state === nfa.startState ? '#D2E5FF' : '#FFD2D2',
				borderWidth: nfa.acceptedStates.includes(state) ? 3 : 1
			}))
		);
		const edgesArr: Edge[] = [];
		for (const from of nfa.states) {
			for (const symbol of nfa.alphabet) {
				const tos = (nfa.transitions[from]?.[symbol] || []) as string[];
				for (const to of tos) {
					edgesArr.push({ from: nodeIds[from], to: nodeIds[to], label: symbol, arrows: 'to' });
				}
			}
		}
		const START_NODE_ID = '__start__';
		nodes.add({ id: START_NODE_ID, label: '', shape: 'text' });
		edgesArr.push({
			from: START_NODE_ID,
			to: nodeIds[nfa.startState],
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
		const options: Options = {
			nodes: { shape: 'circle', font: { size: 16 }, margin: { top: 10, right: 10, bottom: 10, left: 10 } },
			edges: { smooth: { enabled: true, type: 'curvedCW', roundness: 0.3 }, font: { size: 14, strokeWidth: 0 } },
			physics: false
		};
		new Network(nfaContainer, { nodes, edges }, options);
	}

	async function renderDfaVis() {
		const dfa = nfaToDfa(parseAutomaton());
		const nodeIds: Record<string, string> = {};
		dfa.states.forEach((state) => {
			nodeIds[state] = state.replace(/[^a-zA-Z0-9_]/g, '_');
		});
		const nodes = new DataSet<Node>(
			dfa.states.map((state) => ({
				id: nodeIds[state],
				label: state,
				shape: 'circle',
				color: dfa.acceptedStates.includes(state) ? '#D2FFD2' : state === dfa.startState ? '#D2E5FF' : '#FFD2D2',
				borderWidth: dfa.acceptedStates.includes(state) ? 3 : 1
			}))
		);
		const edgesArr: Edge[] = [];
		for (const from of dfa.states) {
			for (const symbol of dfa.alphabet) {
				const to = dfa.transitions[from]?.[symbol] as string;
				if (to) {
					edgesArr.push({ from: nodeIds[from], to: nodeIds[to], label: symbol, arrows: 'to' });
				}
			}
		}
		const START_NODE_ID = '__start__';
		nodes.add({ id: START_NODE_ID, label: '', shape: 'text' });
		edgesArr.push({
			from: START_NODE_ID,
			to: nodeIds[dfa.startState],
			arrows: { to: { enabled: true, scaleFactor: 0.6 } },
			color: { color: '#222', opacity: 1 },
			width: 1.75,
			label: 'start',
			font: { size: 13, color: '#222', vadjust: -18, align: 'top' },
			smooth: { enabled: true, type: 'curvedCCW', roundness: 0.18 },
			length: 5,
			physics: false
		});
		const edges = new DataSet(edgesArr);
		const options: Options = {
			nodes: { shape: 'circle', font: { size: 16 }, margin: { top: 10, right: 10, bottom: 10, left: 10 } },
			edges: { smooth: { enabled: true, type: 'curvedCW', roundness: 0.3 }, font: { size: 14, strokeWidth: 0 } },
			physics: false
		};
		new Network(dfaContainer, { nodes, edges }, options);
	}

	async function showNfaDiagram() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		const saved = await saveDfaToBackend();
		if (!saved) return;

		try {
			const nfaToDfaRes = await fetch('http://localhost:8080/api/lexing/nfaToDFA', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!nfaToDfaRes.ok) {
				const errorText = await nfaToDfaRes.text();
				AddToast('NFA→DFA failed: ' + errorText, 'error');
				return;
			}

			const dfaToRegexRes = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!dfaToRegexRes.ok) {
				const errorText = await dfaToRegexRes.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}

			const regexToNfaRes = await fetch('http://localhost:8080/api/lexing/regexToNFA', {
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
			setTimeout(() => {
				if (regexNfa) renderRegexAutomatonVis(nfaContainer, regexNfa, false);
			}, 0);
			AddToast('NFA generated from Regex and displayed!', 'success');
		} catch (error) {
			AddToast('Failed to generate NFA from Regex: ' + String(error), 'error');
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
		states = 'q0,q1,q2';
		startState = 'q0';
		acceptedStates = 'q2->int';
		transitions = 'q0,0->q0\nq0,0->q1\nq1,0->q2\nq1,1->q0\nq2,0->q1\nq2,1->q2';
	}

	function removeDefault() {
		showDefault = false;
		inputRows = [...userInputRows];
		resetInputs();
	}

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
				if (state) {
					return { state, token_type: '' };
				}
				return null;
			})
			.filter((x): x is { state: string; token_type: string } => !!x);
	}

	async function handleShowDfa() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}

		const saved = await saveDfaToBackend();
		if (!saved) return;

		try {
			const regexRes = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ users_id: user_id })
			});
			if (!regexRes.ok) {
				const errorText = await regexRes.text();
				AddToast('DFA→Regex failed: ' + errorText, 'error');
				return;
			}

			const dfaRes = await fetch('http://localhost:8080/api/lexing/regexToDFA', {
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
			setTimeout(() => {
				if (regexDfa) renderRegexAutomatonVis(dfaContainer, regexDfa, true);
			}, 0);
			AddToast('DFA generated from Regex and displayed!', 'success');
		} catch (error) {
			AddToast('Failed to generate DFA from Regex: ' + String(error), 'error');
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
			const response = await fetch('http://localhost:8080/api/lexing/dfaToRegex', {
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
			AddToast('DFA→Regex failed: ' + String(error), 'error');
		}
	}

	let regexNfa: Automaton | null = null;
	let regexDfa: Automaton | null = null;
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
			const response = await fetch('http://localhost:8080/api/lexing/regexToNFA', {
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
			AddToast('Regex converted to NFA!', 'success');
			setTimeout(() => {
				if (regexNfa) renderRegexAutomatonVis(regexNfaContainer, regexNfa, false);
			}, 0);
		} catch (error) {
			AddToast('Regex→NFA failed: ' + String(error), 'error');
		}
	}

	async function handleRegexToDFA() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('User not logged in.', 'error');
			return;
		}
		try {
			const response = await fetch('http://localhost:8080/api/lexing/regexToDFA', {
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
			AddToast('Regex converted to DFA!', 'success');
			setTimeout(() => {
				if (regexDfa) renderRegexAutomatonVis(regexDfaContainer, regexDfa, true);
			}, 0);
		} catch (error) {
			AddToast('Regex→DFA failed: ' + String(error), 'error');
		}
	}

	function adaptAutomatonForVis(automaton: any): Automaton {
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
			acceptedStates: (automaton.accepting_states || automaton.acceptedStates || automaton.accepting || []).map(
				(a: any) => (typeof a === 'string' ? a : a.state || a)
			),
			transitions: transitionsObj,
			alphabet
		};
	}

	function renderRegexAutomatonVis(container: HTMLElement, automaton: Automaton, isDfa = false) {
		if (!automaton || !container) return;
		const nodeIds: Record<string, IdType> = {};
		automaton.states.forEach((state: string) => {
			nodeIds[state] = state.replace(/[^a-zA-Z0-9_]/g, '_');
		});
		const nodes = new DataSet<Node>(
			automaton.states.map((state: string) => ({
				id: nodeIds[state],
				label: state,
				shape: 'circle',
				color: automaton.acceptedStates.includes(state) ? '#D2FFD2' : state === automaton.startState ? '#D2E5FF' : '#FFD2D2',
				borderWidth: automaton.acceptedStates.includes(state) ? 3 : 1
			}))
		);
		const edgesArr: Edge[] = [];
		for (const from of automaton.states) {
			for (const symbol of automaton.alphabet) {
				const tos = isDfa
					? ([automaton.transitions[from]?.[symbol]].filter(Boolean) as string[])
					: ((automaton.transitions[from]?.[symbol] as string[]) || []);
				for (const to of tos) {
					edgesArr.push({ from: nodeIds[from], to: nodeIds[to], label: symbol, arrows: 'to' });
				}
			}
		}
		const START_NODE_ID = '__start__';
		nodes.add({ id: START_NODE_ID, label: '', shape: 'text' });
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
		const options: Options = {
			nodes: { shape: 'circle', font: { size: 16 }, margin: 10 },
			edges: { smooth: { enabled: true, type: 'curvedCW', roundness: 0.3 }, font: { size: 14, strokeWidth: 0 } },
			physics: false
		};
		new Network(container, { nodes, edges }, options);
	}

	let automataDisplay: 'NFA' | 'DFA' | 'RE' | null = null;
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
		<button class="automaton-btn {selectedType === 'REGEX' ? 'selected' : ''}" on:click={() => selectType('REGEX')} type="button">
			Regular Expression
		</button>
		<button class="automaton-btn {selectedType === 'AUTOMATA' ? 'selected' : ''}" on:click={() => { selectType('AUTOMATA'); automataDisplay = null; }} type="button">
			Automata
		</button>

		{#if selectedType && !showDefault}
			<button class="default-toggle-btn" on:click={insertDefault} type="button" aria-label="Insert default input" title="Insert default input">
				<span class="icon">🪄</span>
			</button>
		{/if}
		{#if selectedType && showDefault}
			<button class="default-toggle-btn selected" on:click={removeDefault} type="button" aria-label="Remove default input" title="Remove default input">
				<span class="icon">🧹</span>
			</button>
		{/if}
	</div>

	{#if selectedType === 'REGEX'}
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
					{#each (showDefault ? editableDefaultRows : userInputRows) as row, i}
						<div class="input-row">
							<div class="input-block">
								<input type="text" bind:value={row.type} on:input={handleInputChange} placeholder="Enter type..." class:error={row.error} />
							</div>
							<div class="input-block">
								<input type="text" bind:value={row.regex} on:input={handleInputChange} placeholder="Enter regex pattern..." class:error={row.error} />
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
						<button class="generate-button" on:click={handleRegexToNFA} title="Convert Regular Expression to a NFA">NFA</button>
						<button class="generate-button" on:click={handleRegexToDFA} title="Convert Regular Expression to a DFA">DFA</button>
					</div>
				{/if}
			</div>

			{#if submissionStatus.show}
				<div class="status-message" class:success={submissionStatus.success === true} class:info={submissionStatus.message === 'info'}>
					{submissionStatus.message}
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
		</div>
	{:else if selectedType === 'AUTOMATA'}
		<div class="automaton-section">
			<div class="automaton-left">
				<label>
					States
					<input class="automaton-input" bind:value={states} placeholder="e.g. q0, q1, q2" />
				</label>
				<label>
					Start State
					<input class="automaton-input" bind:value={startState} placeholder="e.g. q0" />
				</label>
				<label>
					Accepted States
					<input class="automaton-input" bind:value={acceptedStates} placeholder="e.g. q2->int, q1->str" />
				</label>
			</div>
			<div class="automaton-right">
				<label>
					Transitions
					<textarea class="automaton-input automaton-transitions" bind:value={transitions} placeholder="e.g. q0,a->q1&#10;q1,b->q2"></textarea>
				</label>
			</div>
			<div class="automata-action-row">
				<button class="action-btn" type="button" on:click={() => { showNfaDiagram(); automataDisplay = 'NFA'; }}>Show NFA</button>
				<button class="action-btn" type="button" on:click={() => { renderDfaVis(); showDfaVis = true; automataDisplay = 'DFA'; }}>Show DFA</button>
				<button class="action-btn" type="button" on:click={() => { handleTokenisation(); automataDisplay = null; }}>Tokenisation</button>
				<button class="action-btn" type="button" on:click={() => { handleConvertToRegex(); automataDisplay = 'RE'; }} title="Convert to Regular Expression">RE</button>
			</div>
		</div>

		{#if automataDisplay === 'NFA' && showNfaVis}
			<div class="automata-container">
				<div class="pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">NFA Visualization</span>
					</div>
					<div bind:this={nfaContainer} class="vis-graph-area" />
				</div>
			</div>
		{:else if automataDisplay === 'DFA' && showDfaVis}
			<div class="automata-container">
				<div class="pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">DFA Visualization</span>
					</div>
					<div bind:this={dfaContainer} class="vis-graph-area" />
				</div>
			</div>
		{:else if automataDisplay === 'RE' && showRegexOutput && regexRules.length > 0}
			<div class="regex-display-container">
				<div class="vis-heading">
					<h3 class="title">Generated Regular Expressions</h3>
				</div>
				<div class="table-wrapper">
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
									<td class="regex-type">{rule.token_type || rule.Type || '-'}</td>
									<td class="regex-pattern">
										<code>{rule.regex || rule.Regex || '-'}</code>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}

		{#if showRegexNfaVis && regexNfa}
			<div class="automata-container">
				<div class="pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">NFA Visualization (from REGEX)</span>
					</div>
					<div bind:this={regexNfaContainer} class="vis-graph-area" />
				</div>
			</div>
		{/if}
		{#if showRegexDfaVis && regexDfa}
			<div class="automata-container">
				<div class="pretty-vis-box">
					<div class="vis-heading">
						<span class="vis-title">DFA Visualization (from REGEX)</span>
					</div>
					<div bind:this={regexDfaContainer} class="vis-graph-area" />
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	.title {
		font-size: 1.2rem;
		font-weight: 600;
		margin-bottom: 0;
	}
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
		transition: border-color 0.2s, box-shadow 0.2s;
	}
	.input-block input:focus {
		outline: none;
		border-color: #001a6e;
		box-shadow: 0 0 0 2px rgba(0, 26, 110, 0.1);
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
	}
	.button-stack {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem;
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
		color: #f0f0f0;
		background-color: #1f2937;
		border: 1px solid #4a5568;
	}

	:global(html.dark-mode) .header-section h3 {
		color: #f0f0f0;
	}

	:global(html.dark-mode) .shared-block {
		color: #f0f0f0;
		background-color: #1f2937;
		border: 2px solid #374151;
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
		transition: border-color 0.2s, background 0.2s, color 0.2s;
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
		border: 2px solid #e5e7eb;
		background: white;
		color: #001a6e;
		font-size: 1.2rem;
		cursor: pointer;
		transition: background 0.2s, border-color 0.2s;
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

	.token-type {
		font-weight: 500;
		color: #084298;
	}

	.regex-action-buttons {
		display: flex;
		justify-content: center;
		gap: 0.5rem;
		margin-top: 0.25rem;
	}

	.pretty-vis-box {
		background: #f8fafc;
		border-radius: 14px;
		border: 1px solid #e5e7eb;
		box-shadow:
			0 2px 12px rgba(30, 64, 175, 0.07),
			0 1.5px 6px rgba(0, 0, 0, 0.04);
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		align-items: stretch;
		width: 100%;
		max-width: 750px;
	}

	.vis-heading {
		display: flex;
		align-items: center;
		gap: 0.7rem;
		margin-bottom: 1.1rem;
		font-size: 1.25rem;
		font-weight: 600;
		color: #001a6e;
		letter-spacing: 0.01em;
		padding-bottom: 0.75rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.vis-title {
		font-size: 1.18rem;
		font-weight: 600;
		color: #001a6e;
	}
	.vis-graph-area {
		width: 100%;
		height: 350px;
		border: 1.5px solid #e0e7ef;
		border-radius: 10px;
		background: #fff;
		box-shadow: 0 1px 4px rgba(30, 64, 175, 0.04);
	}

	.regex-display-container {
		margin-top: 2.5rem;
		background: #ffffff;
		border: 1px solid #e5e7eb;
		border-radius: 12px;
		padding: 1rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		max-width: 700px;
	}

	.regex-display-container h4 {
		font-size: 1.1rem;
		font-weight: 600;
		color: #111827;
		padding: 0.5rem 0.5rem;
		margin: 0 0 1rem 0;
		border-bottom: 1px solid #e5e7eb;
	}

	.table-wrapper {
		border-radius: 8px;
		overflow: hidden;
		border: 1px solid #e5e7eb;
	}

	.regex-table {
		width: 100%;
		border-collapse: collapse;
	}

	.regex-table th,
	.regex-table td {
		padding: 0.8rem 1.25rem;
		text-align: left;
		border-bottom: 1px solid #e5e7eb;
	}

	.regex-table tbody tr:last-child td {
		border-bottom: none;
	}

	.regex-table th {
		background: #f9fafb;
		color: #4b5563;
		font-size: 0.8rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.regex-table tbody tr:nth-child(even) {
		background-color: #f9fafb;
	}

	.regex-type {
		font-weight: 500;
		color: #001a6e;
		font-size: 0.9rem;
	}

	.regex-pattern code {
		background: #eef2ff;
		color: #001a6e;
		padding: 0.3em 0.6em;
		border-radius: 6px;
		font-family: 'Fira Mono', monospace;
		font-size: 0.9rem;
		word-break: break-all;
		display: inline-block;
	}

	:global(html.dark-mode) .regex-display-container {
		background: #1f2937;
		border-color: #374151;
	}

	:global(html.dark-mode) .regex-display-container h4 {
		color: #f3f4f6;
		border-bottom-color: #374151;
	}

	:global(html.dark-mode) .table-wrapper {
		border-color: #374151;
	}

	:global(html.dark-mode) .regex-table th,
	:global(html.dark-mode) .regex-table td {
		border-bottom-color: #374151;
	}

	:global(html.dark-mode) .regex-table th {
		background: #374151;
		color: #9ca3af;
	}

	:global(html.dark-mode) .regex-table tbody tr:nth-child(even) {
		background-color: #1f2937;
	}

	:global(html.dark-mode) .regex-table tbody tr {
		background-color: #273142;
	}

	:global(html.dark-mode) .regex-type {
		color: #a5b4fc;
	}

	:global(html.dark-mode) .regex-pattern code {
		background: #222c5e;
		color: #e0e7ff;
	}

	.automaton-section {
		margin-top: 1rem;
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 12px;
		padding: 1.5rem;
		display: grid;
		grid-template-columns: 1fr 1fr;
		grid-template-rows: auto;
		gap: 1.5rem 2rem;
		max-width: 800px;
	}

	.automaton-left {
		display: flex;
		flex-direction: column;
		gap: 1.25rem; /* Increased spacing */
	}

	.automaton-right {
		display: flex;
	}

	.automaton-right label {
		display: flex;
		flex-direction: column;
		flex-grow: 1;
	}

	.automaton-section label {
		font-weight: 500;
		color: #374151;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-size: 0.9rem;
	}

	.automaton-input {
		padding: 0.75rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 6px;
		font-size: 0.9rem;
		background: #ffffff;
		color: #111827;
		transition: border-color 0.2s, box-shadow 0.2s;
		width: 100%;
		box-sizing: border-box;
	}

	.automaton-input:focus {
		outline: none;
		border-color: #001a6e;
		box-shadow: 0 0 0 3px rgba(0, 26, 110, 0.1);
	}

	.automaton-transitions {
		font-family: 'Fira Mono', monospace;
		resize: vertical;
		flex-grow: 1;
	}

	.automata-action-row {
		grid-column: span 2;
		display: flex;
		gap: 0.75rem;
		justify-content: flex-start;
		align-items: center;
		flex-wrap: wrap;
		padding-top: 1rem;
	}

	.action-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		font-family: 'Inter', sans-serif;
		font-weight: 500;
		text-align: center;
		white-space: nowrap;
		cursor: pointer;
		transition: all 0.2s ease-in-out;
		user-select: none;
		border: 1px solid #dbe7ff;
		outline: none;
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-size: 0.9rem;
		line-height: 1.25;
		background: #e6f0ff;
		color: #001a6e;
	}

	.action-btn:hover {
		background: #dbe7ff;
		border-color: #caddff;
	}

	.action-btn:focus {
		box-shadow: 0 0 0 3px rgba(0, 26, 110, 0.2);
	}

	:global(html.dark-mode) .automaton-section {
		background: #1f2937;
		border-color: #374151;
	}

	:global(html.dark-mode) .automaton-section label {
		color: #e0e7ff; /* Brighter label color for dark mode */
	}

	:global(html.dark-mode) .automaton-input {
		background-color: #374151;
		border-color: #4b5563;
		color: #f3f4f6;
	}

	:global(html.dark-mode) .automaton-input::placeholder {
		color: #9ca3af;
	}

	:global(html.dark-mode) .automaton-input:focus {
		border-color: #60a5fa;
		box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.2);
	}

	:global(html.dark-mode) .action-btn {
		background: #2c3a6e;
		color: #e0e7ff;
		border-color: #3b4a8e;
	}

	:global(html.dark-mode) .action-btn:hover {
		background: #3b4a8e;
	}

	:global(html.dark-mode) .pretty-vis-box {
		background: #1f2937;
		border-color: #374151;
		box-shadow:
			0 2px 12px rgba(0, 0, 0, 0.1),
			0 1.5px 6px rgba(0, 0, 0, 0.1);
	}
	:global(html.dark-mode) .vis-heading,
	:global(html.dark-mode) .vis-title {
		color: #d1d5db;
	}
	:global(html.dark-mode) .vis-graph-area {
		background-color: #273142;
		border-color: #374151;
	}

	.automata-container {
		display: flex;
		justify-content: center;
		width: 100%;
		margin-top: 2.2rem;
		margin-bottom: 1.5rem;
	}

	:global(html.dark-mode) .automaton-btn {
		background-color: #2d3748;
		border-color: #4a5568;
		color: #d1d5db;
	}
	:global(html.dark-mode) .automaton-btn.selected {
		background-color: #001a6e;
		border-color: #60a5fa;
		color: #e0e7ff;
	}
	:global(html.dark-mode) .automaton-btn:not(.selected):hover {
		background-color: #374151;
		border-color: #6b7280;
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
</style>

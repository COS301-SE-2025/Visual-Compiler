<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { writable, get, type Writable } from 'svelte/store';
	import type { NodeType, Token, SyntaxTree, NodeConnection } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { theme } from '../../lib/stores/theme';
	import { projectName } from '$lib/stores/project';
	import { pipelineStore, setActivePhase } from '$lib/stores/pipeline';
	import { confirmedSourceCode } from '$lib/stores/source-code';
	import NavBar from '$lib/components/main/nav-bar.svelte';
	import Toolbox from '$lib/components/main/Toolbox.svelte';
	import CodeInput from '$lib/components/main/code-input.svelte';
	import DrawerCanvas from '$lib/components/main/drawer-canvas.svelte';
	import WelcomeOverlay from '$lib/components/project-hub/project-hub.svelte';
	import ClearCanvasConfirmation from '$lib/components/main/clear-canvas-confirmation.svelte';
	import AiAssistant from '$lib/components/main/ai-assistant.svelte';
	import { phase_completion_status } from '$lib/stores/pipeline';

	// --- CANVAS STATE ---
	interface CanvasNode {
		id: string;
		type: NodeType;
		label: string;
		position: { x: number; y: number };
	}
	const nodes: Writable<CanvasNode[]> = writable([]);
	let node_counter = 0;

	// --- CONNECTION TRACKING STATE ---
	let physicalConnections: NodeConnection[] = [];

	// --- COMPONENT STATE ---
	let LexerPhaseTutorial: any;
	let LexerPhaseInspector: any;
	let LexerArtifactViewer: any;
	let ParserPhaseTutorial: any;
	let ParserPhaseInspector: any;
	let ParserArtifactViewer: any;
	let AnalyserPhaseTutorial: any;
	let AnalyserPhaseInspector: any;
	let AnalyserArtifactViewer: any;
	let TranslatorPhaseTutorial: any;
	let TranslatorPhaseInspector: any;
	let TranslatorArtifactViewer: any;
	let OptimizerPhaseTutorial: any;
	let OptimizerPhaseInspector: any;
	let OptimizerArtifactViewer: any;

	let showWelcomeOverlay = false;
	let workspace_el: HTMLElement;
	let show_drag_tip = false;
	let showClearCanvasModal = false;

	// --- UNSAVED CHANGES TRACKING ---
	let lastSavedState: string | null = null;

	// Function to handle beforeunload event
	const handleBeforeUnload = (event: BeforeUnloadEvent) => {
		// Get current pipeline state
		const currentNodes = get(nodes);
		const currentState = JSON.stringify({
			nodes: currentNodes,
			connections: physicalConnections
		});

		// Compare with last saved state
		if (lastSavedState && currentState !== lastSavedState) {
			// There are unsaved changes
			event.preventDefault();
			event.returnValue = '';
			return '';
		}
	};

	// Subscribe to the project name store
	let currentProjectName = '';
	projectName.subscribe((value) => {
		currentProjectName = value;
	});

	// Subscribe to pipeline store changes
	const unsubscribePipeline = pipelineStore.subscribe(pipeline => {
		if (pipeline) {
			if (Array.isArray(pipeline.nodes)) {
				// Update the canvas nodes with the restored pipeline nodes
				nodes.set(pipeline.nodes);
				// Reset the node counter to be higher than any existing node ID
				node_counter = pipeline.nodes.reduce((maxId, node) => {
					const idNum = parseInt(node.id.split('-')[1]) || 0;
					return Math.max(maxId, idNum);
				}, 0);

				// Filter out connections that reference non-existent nodes
				if (Array.isArray(pipeline.connections)) {
					const validConnections = pipeline.connections.filter(conn => {
						const sourceExists = pipeline.nodes.some(node => node.id === conn.sourceNodeId);
						const targetExists = pipeline.nodes.some(node => node.id === conn.targetNodeId);
						return sourceExists && targetExists;
					});
					physicalConnections = validConnections;
				}

				// Update last saved state when project is loaded
				lastSavedState = JSON.stringify({
					nodes: pipeline.nodes,
					connections: physicalConnections
				});
			}
		}
	});

	onMount(async () => {
		// Load dynamic components
		LexerPhaseTutorial = (await import('$lib/components/lexer/lexer-phase-tutorial.svelte')).default;
		LexerPhaseInspector = (await import('$lib/components/lexer/phase-inspector.svelte')).default;
		LexerArtifactViewer = (await import('$lib/components/lexer/lexer-artifact-viewer.svelte'))
			.default;
		ParserPhaseTutorial = (await import('$lib/components/parser/parser-phase-tutorial.svelte'))
			.default;
		ParserPhaseInspector = (await import('$lib/components/parser/parsing-input.svelte')).default;
		ParserArtifactViewer = (
			await import('$lib/components/parser/parser-artifact-viewer.svelte')
		).default;
		AnalyserPhaseTutorial = (
			await import('$lib/components/analyser/analyser-phase-tutorial.svelte')
		).default;
		AnalyserPhaseInspector = (
			await import('$lib/components/analyser/analyser-phase-inspector.svelte')
		).default;
		AnalyserArtifactViewer = (
			await import('$lib/components/analyser/analyser-artifact-viewer.svelte')
		).default;
		TranslatorPhaseTutorial = (
			await import('$lib/components/translator/translator-phase-tutorial.svelte')
		).default;
		TranslatorPhaseInspector = (
			await import('$lib/components/translator/translator-phase-inspector.svelte')
		).default;
		TranslatorArtifactViewer = (
			await import('$lib/components/translator/translator-artifact-viewer.svelte')
		).default;
		OptimizerPhaseTutorial = (
			await import('$lib/components/optimizer/optimizer-phase-tutorial.svelte')
		).default;
		OptimizerPhaseInspector = (
			await import('$lib/components/optimizer/optimizer-phase-inspector.svelte')
		).default;
		OptimizerArtifactViewer = (
			await import('$lib/components/optimizer/optimizer-artifact-viewer.svelte')
		).default;

		// Setup theme and UI state
		document.documentElement.classList.toggle('dark-mode', $theme === 'dark');
		if (!localStorage.getItem('hasSeenDragTip')) {
			show_drag_tip = true;
		}

		if (sessionStorage.getItem('showWelcomeOverlay') === 'true') {
			showWelcomeOverlay = true; // Trigger the overlay to show.
		}

		// --- UNSAVED CHANGES PROTECTION ---
		// Only add event listener if we're in the browser
		if (typeof window !== 'undefined') {
			window.addEventListener('beforeunload', handleBeforeUnload);

			// Initialize lastSavedState for blank canvas
			if (!lastSavedState) {
				lastSavedState = JSON.stringify({
					nodes: [],
					connections: []
				});
			}
		}

		// Return cleanup function
		return () => {
			// Cleanup subscriptions
			unsubscribePipeline();
			
			// Cleanup event listener if in browser
			if (typeof window !== 'undefined') {
				window.removeEventListener('beforeunload', handleBeforeUnload);
			}
		};
	});

	// Use onDestroy as additional cleanup
	onDestroy(() => {
		// This ensures cleanup even if the onMount return function doesn't run
		// Only remove event listener if we're in the browser
		if (typeof window !== 'undefined') {
			window.removeEventListener('beforeunload', handleBeforeUnload);
		}
	});

	function handleWelcomeClose() {
		showWelcomeOverlay = false;
	}

	// --- CANVAS STATE ---
	let selected_phase: NodeType | null = null;
	let show_code_input = false;
	let source_code = '';
	let show_tokens = false;
	let syntaxTreeData: SyntaxTree | null = null;
	let translationError: any = null;
	let savedProjectData: object | null = null;


	// Handle physical connection changes from canvas
	function handleConnectionChange(connections: NodeConnection[]) {
		physicalConnections = connections;
		console.log('Physical connections updated:', physicalConnections);
	}

	// Check if there's a physical connection between two node types
	function hasPhysicalConnection(sourceType: NodeType, targetType: NodeType): boolean {
		return physicalConnections.some(conn => {
			const sourceNode = findNodeByType(sourceType);
			const targetNode = findNodeByType(targetType);
			
			if (!sourceNode || !targetNode) return false;
			
			return (conn.sourceNodeId === sourceNode.id && conn.targetNodeId === targetNode.id) ||
			       (conn.sourceNodeId === targetNode.id && conn.targetNodeId === sourceNode.id);
		});
	}

	// --- CONNECTION VALIDATION FUNCTIONS ---
	function findNodeByType(nodeType: NodeType): CanvasNode | null {
		const currentNodes = get(nodes);
		return currentNodes.find(node => node.type === nodeType) || null;
	}



	// Add subscription to phase completion status
	let completion_status;
	phase_completion_status.subscribe(value => {
	    completion_status = value;
	});

	function validateNodeAccess(nodeType: NodeType): boolean {
		const currentNodes = get(nodes);
		const confirmedCode = get(confirmedSourceCode);
		source_code = confirmedCode; // Keep local variable in sync
		
		switch (nodeType) {
			case 'source':
				return true; // Source is always accessible
			
			case 'lexer':
				// Check if source node exists
				const sourceNode = findNodeByType('source');
				if (!sourceNode) {
					AddToast('Missing Source Code: Add a Source Code node from the toolbox to begin lexical analysis', 'error');
					return false;
				}
				// Check if source code has been submitted using the confirmedSourceCode store
				if (!confirmedCode.trim()) {
					AddToast('No source code provided: Please enter and submit your source code before proceeding to lexical analysis', 'error');
					return false;
				}
				// Check for physical connection between source and lexer
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing connection: Connect the Source Code node to the Lexer node to establish data flow', 'error');
					return false;
				}
				return true;
			
			case 'parser':
				// First check if source node exists
				const sourceNodeForParser = findNodeByType('source');
				if (!sourceNodeForParser) {
					AddToast('Missing Source Code: Add a Source Code node from the toolbox to begin parsing', 'error');
					return false;
				}
				// Check if source code has been submitted
				if (!source_code.trim()) {
					AddToast('No source code provided: Please enter and submit your source code before accessing the Parser', 'error');
					return false;
				}
				// Check if lexer node exists
				const lexerNode = findNodeByType('lexer');
				if (!lexerNode) {
					AddToast('Missing Lexer: Add and complete lexical analysis before parsing your code', 'error');
					return false;
				}
				// Check for physical connection between source and lexer
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing Source→Lexer connection: Connect these nodes to enable data flow', 'error');
					return false;
				}
				// Check if lexer phase has been completed using the store value
				if (!completion_status.lexer) {
					AddToast('Lexical analysis incomplete: Complete tokenization in the Lexer before parsing', 'error');
					return false;
				}
				// Check for physical connection between lexer and parser
				if (!hasPhysicalConnection('lexer', 'parser')) {
					AddToast('Missing Lexer→Parser connection: Connect these nodes to enable parsing', 'error');
					return false;
				}
				return true;
			
			case 'analyser':
				// First check if source node exists
				const sourceNodeForAnalyser = findNodeByType('source');
				if (!sourceNodeForAnalyser) {
					AddToast('Missing Source Code: Add a Source Code node from the toolbox to begin semantic analysis', 'error');
					return false;
				}
				// Check if source code has been submitted
				if (!source_code.trim()) {
					AddToast('No source code provided: Please enter and submit your source code before accessing the Analyser', 'error');
					return false;
				}
				// Check if lexer node exists
				const lexerNodeForAnalyser = findNodeByType('lexer');
				if (!lexerNodeForAnalyser) {
					AddToast('Missing Lexer: Complete lexical analysis before running semantic analysis', 'error');
					return false;
				}
				// Check for physical connection between source and lexer
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing Source→Lexer connection: Connect these nodes to establish data flow', 'error');
					return false;
				}
				// Check if lexer phase has been completed
				if (!completion_status.lexer) {
					AddToast('Lexical analysis incomplete: Complete tokenization before semantic analysis', 'error');
					return false;
				}
				// Check if parser node exists
				const parserNode = findNodeByType('parser');
				if (!parserNode) {
					AddToast('Missing Parser: Add and complete parsing before semantic analysis', 'error');
					return false;
				}
				// Check for physical connection between lexer and parser
				if (!hasPhysicalConnection('lexer', 'parser')) {
					AddToast('Missing Lexer→Parser connection: Connect these nodes to enable parsing', 'error');
					return false;
				}
				// Check if parser phase has been completed
				if (!completion_status.parser) {
					AddToast('Parsing incomplete: Complete syntax analysis before semantic analysis', 'error');
					return false;
				}
				// Check for physical connection between parser and analyser
				if (!hasPhysicalConnection('parser', 'analyser')) {
					AddToast('Missing Parser→Analyser connection: Connect these nodes to enable semantic analysis', 'error');
					return false;
				}
				return true;
			
			case 'translator':
				// First check if source node exists
				const sourceNodeForTranslator = findNodeByType('source');
				if (!sourceNodeForTranslator) {
					AddToast('Missing Source Code: Add a Source Code node from the toolbox to begin translation', 'error');
					return false;
				}
				// Check if source code has been submitted
				if (!source_code.trim()) {
					AddToast('No source code provided: Please enter and submit your source code before accessing the Translator', 'error');
					return false;
				}
				// Check if lexer node exists
				const lexerNodeForTranslator = findNodeByType('lexer');
				if (!lexerNodeForTranslator) {
					AddToast('Missing Lexer: Complete lexical analysis before code translation', 'error');
					return false;
				}
				// Check for physical connection between source and lexer
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing Source→Lexer connection: Connect these nodes to establish data flow', 'error');
					return false;
				}
				// Check if lexer phase has been completed
				if (!completion_status.lexer) {
					AddToast('Lexical analysis incomplete: Complete tokenization before translation', 'error');
					return false;
				}
				// Check if parser node exists
				const parserNodeForTranslator = findNodeByType('parser');
				if (!parserNodeForTranslator) {
					AddToast('Missing Parser: Complete parsing before code translation', 'error');
					return false;
				}
				// Check for physical connection between lexer and parser
				if (!hasPhysicalConnection('lexer', 'parser')) {
					AddToast('Missing Lexer→Parser connection: Connect these nodes to enable parsing', 'error');
					return false;
				}
				// Check if parser phase has been completed
				if (!completion_status.parser) {
					AddToast('Parsing incomplete: Complete syntax analysis before translation', 'error');
					return false;
				}
				// Check if analyser node exists
				const analyserNode = findNodeByType('analyser');
				if (!analyserNode) {
					AddToast('Missing Analyser: Complete semantic analysis before code translation', 'error');
					return false;
				}
				// Check for physical connection between parser and analyser
				if (!hasPhysicalConnection('parser', 'analyser')) {
					AddToast('Missing Parser→Analyser connection: Connect these nodes to enable analysis', 'error');
					return false;
				}
				// Check if analyser phase has been completed
				if (!completion_status.analyser) {
					AddToast('Analysis incomplete: Complete semantic analysis before translation', 'error');
					return false;
				}
				// Check for physical connection between analyser and translator
				if (!hasPhysicalConnection('analyser', 'translator')) {
					AddToast('Missing Analyser→Translator connection: Connect these nodes to enable translation', 'error');
					return false;
				}
				return true;
			
			default:
				return false;
		}
	}

	// --- SAVE PROJECT FUNCTIONALITY ---
	async function saveProject() {
		const user_id = localStorage.getItem('user_id');
		if (!user_id) {
			AddToast('Please log in to save your project.', 'error');
			return;
		}

		if (!currentProjectName) {
			AddToast('Please select a project first.', 'error');
			return;
		}

		// Prepare the pipeline data
		const canvasNodes = get(nodes);

		const pipeline = {
			nodes: canvasNodes,
			connections: physicalConnections,
			lastSaved: new Date().toISOString(),
		};

		// Update the pipeline store to keep it in sync
		pipelineStore.set(pipeline);

		console.log('Pipeline saved with anchors:', get(pipelineStore));

		try {
			const response = await fetch('http://localhost:8080/api/users/savePipeline', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					users_id: user_id,
					project_name: currentProjectName,
					pipeline: pipeline
				})
			});

			if (!response.ok) {
				const error = await response.text();
				throw new Error(error);
			}

			const data = await response.json();
			console.log('Project Saved:', data);
			savedProjectData = pipeline;
			
			// Update last saved state for unsaved changes tracking
			lastSavedState = JSON.stringify({
				nodes: canvasNodes,
				connections: physicalConnections
			});
			
			AddToast(`Project "${currentProjectName}" saved successfully!`, 'success');
		} catch (error) {
			console.error('Failed to save project:', error);
			AddToast(`Failed to save project: ${error.message}`, 'error');
		}
	}

	// --- CLEAR CANVAS FUNCTIONALITY ---
	function showClearCanvasConfirmation() {
		showClearCanvasModal = true;
	}

	function handleClearCanvasConfirm() {
		// Clear all nodes and connections
		nodes.set([]);
		physicalConnections = [];
		
		// Reset the pipeline store
		pipelineStore.update(pipeline => ({
			...pipeline,
			nodes: [],
			connections: []
		}));

		// Reset node counter
		node_counter = 0;

		// Reset the toolbox created nodes (we need to access the Toolbox component's internal state)
		// We'll trigger a custom event that the Toolbox component will listen to
		const event = new CustomEvent('resetToolbox');
		document.dispatchEvent(event);

		// Reset last saved state to reflect the cleared canvas
		lastSavedState = JSON.stringify({
			nodes: [],
			connections: []
		});

		showClearCanvasModal = false;
		AddToast('Canvas cleared successfully!', 'success');
	}

	function handleClearCanvasCancel() {
		showClearCanvasModal = false;
	}

	// --- TOOLTIPS AND LABELS ---
	const tooltips: Record<NodeType, string> = {
		source: 'Start here. Add source code to begin compilation.',
		lexer: 'Converts source code into tokens for processing.',
		parser: 'Analyzes the token stream to build a syntax tree.',
		analyser: 'Performs semantic analysis on the syntax tree.',
		translator: 'Translates the syntax tree into target code.',
		optimizer: 'Advanced optimisation techniques for code enhancement.'
	};

	const node_labels: Record<NodeType, string> = {
		source: 'Source Code',
		lexer: 'Lexer',
		parser: 'Parser',
		analyser: 'Analyser',
		translator: 'Translator',
		optimizer: 'Optimiser'
	};

	function handleCreateNode(type: NodeType) {
		node_counter++;
		nodes.update((curr) => {
			const start_x = 100;
			const start_y = 100;
			const x_offset = 300;
			const y_offset = 150;
			const nodes_per_row = 3;
			const new_node_index = curr.length;
			const new_position = {
				x: start_x + (new_node_index % nodes_per_row) * x_offset,
				y: start_y + Math.floor(new_node_index / nodes_per_row) * y_offset
			};
			const new_node = {
				id: `${type}-${node_counter}`,
				type,
				label: node_labels[type] || type[0].toUpperCase() + type.slice(1),
				position: new_position
			};
			return [...curr, new_node];
		});
		workspace_el?.focus();
	}

	function dismissDragTip() {
		localStorage.setItem('hasSeenDragTip', 'true');
		show_drag_tip = false;
	}

	function handlePhaseSelect(type: NodeType) {
		// Validate node access before proceeding
		if (type !== 'optimizer' && !validateNodeAccess(type)) {
			return; // Toast message already shown by validateNodeAccess
		}

		syntaxTreeData = null;
		parsing_error = false;
		translationError = null;
		if (type === 'source') {
			show_code_input = true;
			// Update active phase for AI assistant
			setActivePhase('source');
		} else {
			selected_phase = type;
			// Update active phase for AI assistant
			setActivePhase(type);
			// Only check for source code on non-optimizer phases
			if (type !== 'optimizer') {
				const confirmedCode = get(confirmedSourceCode);
				if (!confirmedCode.trim()) {
					AddToast('Source code required: Please add source code to begin the compilation process', 'error');
					selected_phase = null;
					// Clear active phase for AI assistant when phase selection fails
					setActivePhase(null);
					return;
				}
			}
		}
	}

	function handleTranslationError(event: CustomEvent) {
		translationError = event.detail;
		translated_code = [];
	}

	let show_symbol_table = false;
	let symbol_table: Symbol[] = [];
	let analyser_error = false;
	let analyser_error_details = '';

	function handleReset() {
		show_symbol_table = false;
		symbol_table = [];
		analyser_error = false;
		analyser_error_details = '';
	}

	function handleSymbolGeneration(data: { symbol_table: Symbol[]; analyser_error?: boolean; analyser_error_details?: string }) {
		if (data.symbol_table && data.symbol_table.length > 0) {
			show_symbol_table = true;
			symbol_table = data.symbol_table;
			analyser_error = false;
			analyser_error_details = '';
			// Mark analyser phase as complete when symbol table is generated successfully
			phase_completion_status.analyser = true;
			completion_status.analyser = true;
		} else {
			show_symbol_table = false;
			analyser_error = true;
			analyser_error_details = data.analyser_error_details || '';
		}
	}

	function handleTranslationReceived(event: CustomEvent<string[]>) {
		translated_code = event.detail;
		translationError = null;
		// Mark translator phase as complete when translation is received
		if (event.detail && event.detail.length > 0) {
			phase_completion_status.translator = true;
			completion_status.translator = true;
		}
	}

	function returnToCanvas() {
		selected_phase = null;
		show_code_input = false;
		// Clear active phase for AI assistant
		setActivePhase(null);
	}

	function handleCodeSubmit(code: string) {
		show_tokens = false;
		source_code = code;
		confirmedSourceCode.set(code); // Update the store
		show_code_input = false;
		// Mark source phase as complete when code is submitted
		phase_completion_status.source = true;
		completion_status.source = true;
	}

	function handleTokenGeneration(data: { tokens: Token[]; unexpected_tokens: string[] }) {
		show_tokens = true;
		tokens = data.tokens;
		unexpected_tokens = data.unexpected_tokens;
		// Mark lexer phase as complete when tokens are generated successfully
		if (data.tokens.length > 0) {
			phase_completion_status.lexer = true;
			completion_status.lexer = true;
		}
	}

	function handleTreeReceived(event: CustomEvent<SyntaxTree>) {
		syntaxTreeData = event.detail;
		artifactData = event.detail;
		// Mark parser phase as complete when syntax tree is received
		phase_completion_status.parser = true;
		completion_status.parser = true;
		parsing_error = false;
	}

	let tokens: Token[] = [];
	let unexpected_tokens: string[] = [];
	let translated_code: string[] = [];

	let artifactData: SyntaxTree | null = null;
	let parsing_error: boolean=false;
	let parsing_error_details: string="";

	function handleParsingError(data: { parsing_error?: boolean; parsing_error_details?: string }) {
		parsing_error = true;
		parsing_error_details = data.parsing_error_details || '';
		syntaxTreeData = null;
	}
</script>

<NavBar />

<div class="main">
	<WelcomeOverlay bind:show={showWelcomeOverlay} on:close={handleWelcomeClose} />

	<Toolbox {handleCreateNode} {tooltips} nodes={$nodes} />
	<div class="workspace" bind:this={workspace_el} tabindex="-1">
		{#if currentProjectName}
			<div class="project-header">
				<div class="project-info-bar">
					<span class="project-name">{currentProjectName}</span>
					<div class="separator"></div>
					<button class="save-button" on:click={saveProject} aria-label="Save Project" title="Save Project">
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
							<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
							<polyline points="17 21 17 13 7 13 7 21" />
							<polyline points="7 3 7 8 15 8" />
						</svg>
					</button>
					<button class="clear-button" on:click={showClearCanvasConfirmation} aria-label="Clear Canvas" title="Clear Canvas">
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
		{/if}
		<DrawerCanvas {nodes} initialConnections={physicalConnections} onPhaseSelect={handlePhaseSelect} onConnectionChange={handleConnectionChange} />

		{#if show_drag_tip}
			<div class="help-tip">
				<span
					><b>Pro-Tip:</b> For the smoothest experience, click to select a node before dragging
					it.</span
				>
				<button on:click={dismissDragTip} class="dismiss-tip-btn" aria-label="Dismiss tip">
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
						><line x1="18" y1="6" x2="6" y2="18" /><line x1="6" y1="6" x2="18" y2="18"
						/></svg
					>
				</button>
			</div>
		{/if}
	</div>

	{#if selected_phase}
		<div class="analysis-overlay">
			<div class="analysis-view">
				<div class="three-column-layout">
					{#if selected_phase === 'lexer' && LexerPhaseTutorial}
						<svelte:component this={LexerPhaseTutorial} />
						<svelte:component
							this={LexerPhaseInspector}
							{source_code}
							onGenerateTokens={handleTokenGeneration}
						/>
						<svelte:component
							this={LexerArtifactViewer}
							phase={selected_phase}
							{tokens}
							{unexpected_tokens}
							{show_tokens}
						/>
					{/if}

					{#if selected_phase === 'parser' && ParserPhaseTutorial}
						<svelte:component this={ParserPhaseTutorial} />
						<svelte:component
							this={ParserPhaseInspector}
							{source_code}
							on:treereceived={handleTreeReceived}
							on:parsingerror={handleParsingError}
						/>
						<svelte:component this={ParserArtifactViewer} syntaxTree={syntaxTreeData} {parsing_error}{parsing_error_details} />
					{/if}

					{#if selected_phase === 'analyser' && AnalyserPhaseTutorial}
						<svelte:component this={AnalyserPhaseTutorial} />
						<svelte:component
							this={AnalyserPhaseInspector}
							{source_code}
							onGenerateSymbolTable={handleSymbolGeneration}
						/>
						<svelte:component
							this={AnalyserArtifactViewer}
							phase={selected_phase}
							{symbol_table}
							{analyser_error}
							{analyser_error_details}
							{show_symbol_table}
						/>
					{/if}

					{#if selected_phase === 'translator' && TranslatorPhaseTutorial}
						<svelte:component this={TranslatorPhaseTutorial} />
						<svelte:component
							this={TranslatorPhaseInspector}
							{source_code}
							on:translationreceived={handleTranslationReceived}
							on:translationerror={handleTranslationError}
						/>
						<svelte:component
							this={TranslatorArtifactViewer}
							{translated_code}
							{translationError}
						/>
					{/if}

					{#if selected_phase === 'optimizer' && OptimizerPhaseTutorial}
						<svelte:component this={OptimizerPhaseTutorial} />
						<svelte:component this={OptimizerPhaseInspector} />
						<svelte:component this={OptimizerArtifactViewer} />
					{/if}
				</div>
				<button on:click={returnToCanvas} class="return-button"> ← Return to Canvas </button>
			</div>
		</div>
	{/if}

	{#if show_code_input}
		<div class="code-input-overlay">
			<div class="code-input-modal">
				<CodeInput onCodeSubmitted={handleCodeSubmit} />
				<button class="close-btn" on:click={() => (show_code_input = false)}>✕</button>
			</div>
		</div>
	{/if}
</div>

<!-- Clear Canvas Confirmation Modal -->
<ClearCanvasConfirmation 
	bind:show={showClearCanvasModal} 
	on:confirm={handleClearCanvasConfirm}
	on:cancel={handleClearCanvasCancel}
/>

<!-- AI Assistant Component -->
<AiAssistant />

<style>
	:global(html, body) {
		margin: 0;
		padding: 0;
		height: 100%;
		overflow: hidden;
	}
	:global(*) {
		box-sizing: border-box;
	}

	.main {
		display: flex;
		height: calc(100vh - 3.5rem);
		overflow: hidden;
		background-color: #f0f2f5;
		padding: 1rem;
		gap: 1rem;
		transition: background-color 0.3s ease;
	}
	.workspace {
		flex: 1;
		display: flex;
		flex-direction: column;
		position: relative;
		outline: none;
	}
	.project-header {
		position: absolute;
		top: 1rem;
		left: 1rem;
		z-index: 10;
	}
	.project-info-bar {
		display: flex;
		align-items: center;
		background-color: rgba(255, 255, 255, 0.9);
		color: #041a47;
		padding: 0.25rem 0.25rem 0.25rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}
	.project-name {
		padding-right: 0.75rem;
	}
	.separator {
		width: 1px;
		height: 20px;
		background-color: #d1d5db;
		margin-right: 0.5rem;
	}
	.save-button {
		background-color: transparent;
		color: #041a47;
		border: none;
		border-radius: 50%;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: background-color 0.2s ease;
	}
	.save-button:hover {
		background-color: #eef2f7;
	}
	.clear-button {
		background-color: transparent;
		color: #dc2626;
		border: none;
		border-radius: 50%;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: background-color 0.2s ease;
		margin-left: 0.25rem;
	}
	.clear-button:hover {
		background-color: #fee2e2;
	}
	.analysis-overlay {
		position: fixed;
		top: 3.5rem;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(255, 255, 255, 0.95);
		z-index: 100;
	}
	.analysis-view {
		height: 100%;
		padding: 1rem;
		background: #f5f5f5;
		display: flex;
		flex-direction: column;
	}
	.three-column-layout {
		display: flex;
		flex: 1;
		gap: 1rem;
		height: calc(100vh - 6rem);
	}
	.three-column-layout > :global(*) {
		flex: 1;
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		overflow-y: auto;
	}
	.return-button {
		position: fixed;
		bottom: 20px;
		right: 20px;
		padding: 0.5rem 1rem;
		background: #bed2e6;
		color: black;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		z-index: 1000;
		margin-right: 1rem;
		margin-bottom: 1rem;
		font-weight: bold;
	}
	.return-button:hover {
		background: #a8bdd1;
	}
	.code-input-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 2000;
	}
	.code-input-modal {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		max-width: 500px;
		width: 100%;
		position: relative;
	}

	.close-btn {
		position: absolute;
		top: 1rem;
		right: 1rem;
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
	}

	.help-tip {
		position: absolute;
		bottom: 20px;
		left: 50%;
		transform: translateX(-50%);
		background-color: rgba(4, 26, 71, 0.95);
		color: white;
		padding: 10px 15px 10px 20px;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
		display: flex;
		align-items: center;
		gap: 1rem;
		z-index: 50;
		font-size: 0.9rem;
	}

	.dismiss-tip-btn {
		background: none;
		border: none;
		color: white;
		opacity: 0.7;
		cursor: pointer;
		padding: 5px;
		line-height: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: opacity 0.2s ease;
	}
	.dismiss-tip-btn:hover {
		opacity: 1;
	}

	:global(html.dark-mode) .main {
		background-color: #161823;
	}
	:global(html.dark-mode) .project-info-bar {
		background-color: rgba(26, 32, 44, 0.9);
		color: #e2e8f0;
	}
	:global(html.dark-mode) .separator {
		background-color: #4a5568;
	}
	:global(html.dark-mode) .save-button {
		color: #e2e8f0;
	}
	:global(html.dark-mode) .save-button:hover {
		background-color: #2d3748;
	}
	:global(html.dark-mode) .clear-button {
		color: #f87171;
	}
	:global(html.dark-mode) .clear-button:hover {
		background-color: #2d1b1b;
	}
	:global(html.dark-mode) .analysis-overlay {
		background: rgba(10, 26, 58, 0.95);
	}
	:global(html.dark-mode) .analysis-view {
		background: #0a1a3a;
	}
	:global(html.dark-mode) .three-column-layout > :global(*) {
		background: #1a2a4a;
		color: #f0f0f0;
	}
	:global(html.dark-mode) .code-input-modal {
		background: #1a2a4a;
		color: #f0f0f0;
	}

	:global(html.dark-mode) .close-btn {
		color: #f0f0f0;
	}
	:global(html.dark-mode) .return-button {
		background: #1a3a7a;
		margin-right: 1rem;
		color: #cccccc;
	}
	:global(html.dark-mode) .return-button:hover {
		background: #2a4a8a;
		margin-right: 1rem;
	}
</style>

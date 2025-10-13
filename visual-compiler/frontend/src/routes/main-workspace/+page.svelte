<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { writable, get, type Writable } from 'svelte/store';
	import type { NodeType, Token, SyntaxTree, NodeConnection } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { theme } from '../../lib/stores/theme';
	import { pipelineStore, setActivePhase} from '$lib/stores/pipeline';
	import { confirmedSourceCode } from '$lib/stores/source-code';
	import { aiAssistantOpen } from '$lib/stores/ai-assistant';
	import NavBar from '$lib/components/main/nav-bar.svelte';
	import Toolbox from '$lib/components/main/Toolbox.svelte';
	import CodeInput from '$lib/components/main/code-input.svelte';
	import DrawerCanvas from '$lib/components/main/drawer-canvas.svelte';
	import WelcomeOverlay from '$lib/components/project-hub/project-hub.svelte';
	import ClearCanvasConfirmation from '$lib/components/main/clear-canvas-confirmation.svelte';
	import AiAssistant from '$lib/components/main/ai-assistant.svelte';
	import CanvasTutorial from '$lib/components/main/canvas-tutorial.svelte';
	import GuestWelcomePopup from '$lib/components/main/guest-welcome-popup.svelte';
	import { phase_completion_status } from '$lib/stores/pipeline';
	import { tutorialStore, checkTutorialStatus, hideCanvasTutorial, showCanvasTutorial } from '$lib/stores/tutorial';
	import { projectName, clearProject } from '$lib/stores/project';
	import { resetLexerState } from '$lib/stores/lexer';
	import { resetParserState } from '$lib/stores/parser';
	import { resetAnalyserState } from '$lib/stores/analyser';
	import { resetTranslatorState } from '$lib/stores/translator';
	import { resetSourceCode } from '$lib/stores/source-code';

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
	let invalid_connections: NodeConnection[] = [];

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
	let OptimiserPhaseTutorial: any;
	let OptimiserPhaseInspector: any;
	let OptimiserArtifactViewer: any;

	let showWelcomeOverlay = false;
	let workspace_el: HTMLElement;
	let show_drag_tip = false;
	let showClearCanvasModal = false;

	// --- GUEST USER STATE ---
	let showGuestWelcomePopup = false;
	let isGuestUser = false;

	// --- TUTORIAL STATE ---
	let canvasTutorialVisible = false;

	// --- RECENTER STATE ---
	let isRecentering = false;

	// Subscribe to tutorial store
	tutorialStore.subscribe(state => {
		canvasTutorialVisible = state.showCanvasTutorial;
	});

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
	let previousProjectName = '';
	let isInitialLoad = true;
	let projectChangeTimeout: NodeJS.Timeout | null = null;
	
	// Use reactive statement for more reliable project change detection
	$: if ($projectName !== currentProjectName && !isInitialLoad) {
		// Clear any pending project change
		if (projectChangeTimeout) {
			clearTimeout(projectChangeTimeout);
		}
		
		// Debounce project changes to prevent rapid switches
		projectChangeTimeout = setTimeout(() => {
			handleProjectChange($projectName, currentProjectName);
		}, 50);
	}
	
	// Initialize project name tracking
	$: if (isInitialLoad && $projectName) {
		currentProjectName = $projectName;
		previousProjectName = $projectName;
		isInitialLoad = false;
		console.log('Initial project loaded:', $projectName);
	}
	
	// Handle project changes with proper timing
	async function handleProjectChange(newProjectName: string, oldProjectName: string) {
		if (newProjectName && oldProjectName && newProjectName !== oldProjectName) {
			console.log('Project changed from', oldProjectName, 'to', newProjectName);
			
			// Store the phase state before changing
			const wasInPhase = selected_phase || show_code_input;
			
			if (wasInPhase) {
				console.log('Returning to canvas due to project change');
				
				// Return to canvas immediately
				returnToCanvas();
				
				// Small delay for data loading, then show toast
				setTimeout(() => {
					AddToast(`Switched to project: ${newProjectName}`, 'success');
				}, 200);
			}
		}
		
		// Update tracking variables
		previousProjectName = currentProjectName;
		currentProjectName = newProjectName;
	}

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
		OptimiserPhaseTutorial = (
			await import('$lib/components/optimiser/optimiser-phase-tutorial.svelte')
		).default;
		OptimiserPhaseInspector = (
			await import('$lib/components/optimiser/optimiser-phase-inspector.svelte')
		).default;
		OptimiserArtifactViewer = (
			await import('$lib/components/optimiser/optimiser-artifact-viewer.svelte')
		).default;

		// Setup theme and UI state
		document.documentElement.classList.toggle('dark-mode', $theme === 'dark');
		if (!sessionStorage.getItem('hasSeenDragTip')) {
			show_drag_tip = true;
		}

		if (sessionStorage.getItem('showWelcomeOverlay') === 'true') {
			showWelcomeOverlay = true; // Trigger the overlay to show.
		}

		// --- GUEST USER CHECK ---
		// Check if user is a guest and show guest welcome popup
		const accessToken = sessionStorage.getItem('access_token');
		if (accessToken === 'guestuser') {
			showGuestWelcomePopup = true;
			isGuestUser = true;
			
			// Ensure all phase states are completely reset for guest users
			// Clear any remaining phase data that might persist
			show_tokens = false;
			tokens = [];
			unexpected_tokens = [];
			syntaxTreeData = null;
			artifactData = null;
			parsing_error = false;
			parsing_error_details = '';
			show_symbol_table = false;
			symbol_table = [];
			analyser_error = false;
			analyser_error_details = '';
			translated_code = [];
			translationError = null;
		}

		// --- TUTORIAL INITIALIZATION ---
		// Check tutorial status on mount
		checkTutorialStatus();

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

		// FIX: Update event listeners to be more robust
		const handleProjectTokensLoaded = (event: CustomEvent) => {
			console.log('Project tokens loaded event received:', event.detail);
			if (event.detail?.tokens) {
				tokens = event.detail.tokens;
				show_tokens = true;
				handleTokenGeneration({ tokens: event.detail.tokens, unexpected_tokens: [] });
			}
		};

		const handleProjectTreeLoaded = (event: CustomEvent) => {
			console.log('Project tree loaded event received:', event.detail);
			if (event.detail?.tree) {
				syntaxTreeData = event.detail.tree;
				artifactData = event.detail.tree;
				handleTreeReceived({ detail: event.detail.tree });
			}
		};

		const handleProjectSymbolsLoaded = (event: CustomEvent) => {
			console.log('Project symbols loaded event received:', event.detail);
			if (event.detail?.symbols) {
				const symbols = event.detail.symbols.map((s: any) => ({
					name: s.Name || s.name || 'unknown',
					type: s.Type || s.type || 'unknown',
					scope: s.Scope || s.scope || 0
				}));
				symbol_table = symbols;
				show_symbol_table = true;
				handleSymbolGeneration({ symbol_table: symbols });
			}
		};

		const handleProjectTranslationLoaded = (event: CustomEvent) => {
			console.log('Project translation loaded event received:', event.detail);
			if (event.detail?.code) {
				translated_code = event.detail.code;
				handleTranslationReceived({ detail: event.detail.code });
			}
		};

		// Add event listeners
		window.addEventListener('project-tokens-loaded', handleProjectTokensLoaded);
		window.addEventListener('project-tree-loaded', handleProjectTreeLoaded);
		window.addEventListener('project-symbols-loaded', handleProjectSymbolsLoaded);
		window.addEventListener('project-translation-loaded', handleProjectTranslationLoaded);

		// Return cleanup function
		return () => {
			// Cleanup subscriptions
			unsubscribePipeline();
			
			// Cleanup event listener if in browser
			if (typeof window !== 'undefined') {
				window.removeEventListener('beforeunload', handleBeforeUnload);
			}

			// Clean up new event listeners
			window.removeEventListener('project-tokens-loaded', handleProjectTokensLoaded);
			window.removeEventListener('project-tree-loaded', handleProjectTreeLoaded);
			window.removeEventListener('project-symbols-loaded', handleProjectSymbolsLoaded);
			window.removeEventListener('project-translation-loaded', handleProjectTranslationLoaded);
		};
	});

	// Use onDestroy as additional cleanup
	onDestroy(() => {
		// This ensures cleanup even if the onMount return function doesn't run
		// Only remove event listener if we're in the browser
		if (typeof window !== 'undefined') {
			window.removeEventListener('beforeunload', handleBeforeUnload);
		}
		
		// Clear any pending project change timeout
		if (projectChangeTimeout) {
			clearTimeout(projectChangeTimeout);
		}
	});

	function handleWelcomeClose() {
		showWelcomeOverlay = false;
	}

	// Handle guest welcome popup close
	function handleGuestWelcomeClose() {
		showGuestWelcomePopup = false;
		// Show canvas tutorial after guest welcome popup is closed
		showCanvasTutorial();
	}

	// Handle tutorial close
	function handleTutorialClose() {
		hideCanvasTutorial();
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
	function hasPhysicalReceivingConnection(sourceType: NodeType, targetType: NodeType): boolean {
		return physicalConnections.some(conn => {
			const sourceNode = findNodeByType(sourceType);
			const targetNode = findNodeByType(targetType);
			
			if (!sourceNode || !targetNode) return false;
			console.log(conn.sourceNodeId, sourceNode.id)
			
			return (conn.sourceNodeId === sourceNode.id && conn.targetNodeId === targetNode.id)
		});
	}
	function getSpecificConnection(sourceType: NodeType, targetType: NodeType): NodeConnection|null {
		
			const sourceNode = findNodeByType(sourceType);
			const targetNode = findNodeByType(targetType);
			
			if (!sourceNode || !targetNode) {
				 return null;
			}
			
		return physicalConnections.find(conn => 
			(conn.sourceNodeId === sourceNode.id && conn.targetNodeId === targetNode.id)
		) || null;
	}
	function showInvalidConnection() {
		document.querySelectorAll('g[id^="edge-"]').forEach(el => {
			el.querySelector('path')?.classList.remove('invalid');
		});

		invalid_connections.forEach(conn => {
		
			const path_id = `${conn.targetAnchor}+${conn.sourceAnchor}`; 
			const path_id2 = `${conn.sourceAnchor}+${conn.targetAnchor}`; 
			console.log(path_id);
			console.log(path_id2);
			const edgePath = document.getElementById(path_id);
			if (edgePath) {
				edgePath.classList.add('invalid');
			} 
				const edgePath2 = document.getElementById(path_id2);
				if (edgePath2) {
					edgePath2.classList.add('invalid');
				}
			
		});
	}

	function checkInvalidConnections(): Boolean {
		let is_invalid = false;

		if (checkInvalidSourceCodeConnection()) {
			is_invalid=true;
		}
		if (checkInvalidLexerConnection()) {
			is_invalid=true;
		}
		if (checkInvalidParserConnection()) {
			is_invalid=true;
		}
		if (checkInvalidAnalyserConnection()) {
			is_invalid=true;
		}
		
		return is_invalid;
	}

	function checkInvalidSourceCodeConnection(): Boolean {
		let is_invalid = false;
		if (hasPhysicalConnection('source', 'parser')) {
					let conn;
					conn = getSpecificConnection('source','parser');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					conn = getSpecificConnection('parser','source');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
				if (hasPhysicalConnection('source', 'analyser')) {
					let conn;
					conn = getSpecificConnection('source','analyser');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					conn = getSpecificConnection('analyser','source');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
				if (hasPhysicalConnection('source', 'translator')) {
					let conn;
					conn = getSpecificConnection('source','translator');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					conn = getSpecificConnection('translator','source');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
		return is_invalid;
	}
	function checkInvalidLexerConnection() : Boolean {
		let is_invalid = false;
		if (hasPhysicalReceivingConnection('parser', 'lexer')) {
					let conn;
					conn = getSpecificConnection('parser','lexer');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
				if (hasPhysicalConnection('analyser', 'lexer')) {
					let conn;
					conn = getSpecificConnection('analyser','lexer');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}

			return is_invalid;
	}

	function checkInvalidParserConnection() : Boolean {
		let is_invalid = false;
		if (hasPhysicalConnection('lexer', 'analyser')) {
					let conn;
					conn = getSpecificConnection('lexer','analyser');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					conn = getSpecificConnection('analyser','lexer');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;;
				}
				if (hasPhysicalConnection('lexer', 'translator')) {
					let conn;
					conn = getSpecificConnection('lexer','translator');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
				if (hasPhysicalReceivingConnection('analyser', 'parser')) {
					let conn;
					conn = getSpecificConnection('analyser','parser');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
		return is_invalid;
	}
	function checkInvalidAnalyserConnection() : Boolean {
		let is_invalid = false;
		if (hasPhysicalConnection('parser', 'translator')) {
					let conn;
					conn = getSpecificConnection('parser','translator');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
				if (hasPhysicalReceivingConnection('analyser', 'parser')) {
					let conn;
					conn = getSpecificConnection('analyser','parser');
					if (conn!=null)
					{
						invalid_connections.push(conn);
					}
					showInvalidConnection();
					is_invalid = true;
				}
		return is_invalid;
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
		
		var is_invalid;
		is_invalid = false;

		
		switch (nodeType) {
			case 'source':
				checkInvalidConnections();
				showInvalidConnection();
				return true; // Source is always accessible
			
			case 'lexer':
				checkInvalidConnections();
				showInvalidConnection();
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
				if (checkInvalidSourceCodeConnection()) {
					is_invalid=true;
				}
				if (checkInvalidLexerConnection()) {
					is_invalid=true;
				}
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing Source Code → Lexer connection: Connect the Source Code node to the Lexer node to establish data flow', 'error');
					is_invalid = true;
				}
				if (is_invalid)
				{
					AddToast('Incorrect connections present (Highlighted in red)', 'error');
					return false;
				}
				
				return true;
			
			case 'parser':
				checkInvalidConnections();
				showInvalidConnection();
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
				if (checkInvalidSourceCodeConnection()) {
					is_invalid=true;
				}
				if (checkInvalidLexerConnection()) {
					is_invalid=true;
				}
				// Check for physical connection between source and lexer
				if (!hasPhysicalConnection('source', 'lexer')) {
					AddToast('Missing Source→Lexer connection: Connect these nodes to establish data flow', 'error');
					return false;
				}
				// Check if lexer phase has been completed using the store value
				if (!completion_status.lexer) {
					AddToast('Lexical analysis incomplete: Complete tokenization in the Lexer before parsing', 'error');
					return false;
				}
				if (checkInvalidParserConnection()) {
					is_invalid=true;
				}
				// Check for physical connection between lexer and parser
				if (!hasPhysicalConnection('lexer', 'parser')) {
					AddToast('Missing Lexer→Parser connection: Connect these nodes to enable parsing', 'error');
					return false;
				}
				if (is_invalid)
				{
					AddToast('Incorrect connections present (Highlighted in red)', 'error');
					return false;
				}
				return true;
			
			case 'analyser':
				checkInvalidConnections();
				showInvalidConnection();
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
				if (checkInvalidSourceCodeConnection()) {
					is_invalid=true;
				}
				if (checkInvalidLexerConnection()) {
					is_invalid=true;
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
				if (checkInvalidParserConnection()) {
					is_invalid=true;
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
				
				if (checkInvalidAnalyserConnection()) {
					is_invalid=true;
				}
				// Check for physical connection between parser and analyser
				if (!hasPhysicalConnection('parser', 'analyser')) {
					AddToast('Missing Parser→Analyser connection: Connect these nodes to enable semantic analysis', 'error');
					return false;
				}
				if (is_invalid)
				{
					AddToast('Incorrect connections present (Highlighted in red)', 'error');
					return false;
				}
				return true;
			
			case 'translator':
				checkInvalidConnections();
				showInvalidConnection();
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
				if (checkInvalidSourceCodeConnection()) {
					is_invalid=true;
				}
				if (checkInvalidLexerConnection()) {
					is_invalid=true;
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
				if (checkInvalidParserConnection()) {
					is_invalid=true;
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
				if (checkInvalidAnalyserConnection()) {
					is_invalid=true;
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
				if (is_invalid)
				{
					AddToast('Incorrect connections present (Highlighted in red)', 'error');
					return false;
				}
				return true;
			
			default:
				return false;
		}
	}

	// --- SAVE PROJECT FUNCTIONALITY ---
	async function saveProject() {
		const user_id = sessionStorage.getItem('user_id');
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

	// --- RECENTER CANVAS FUNCTIONALITY ---
	async function recenterCanvas() {
		// Save current state
		const currentNodes = get(nodes);
		const currentConnections = [...physicalConnections];
		
		if (currentNodes.length === 0) {
			AddToast('No nodes to recenter', 'info');
			return;
		}

		// Start recentering state
		isRecentering = true;

		// Create temporary pipeline data
		const tempPipeline = {
			nodes: currentNodes,
			connections: currentConnections
		};

		// Clear the canvas
		nodes.set([]);
		physicalConnections = [];

		// Reset node counter to match the highest existing node ID
		node_counter = currentNodes.reduce((maxId, node) => {
			const idNum = parseInt(node.id.split('-')[1]) || 0;
			return Math.max(maxId, idNum);
		}, 0);

		// Longer delay to ensure canvas is fully cleared and ready
		setTimeout(() => {
			// Recreate nodes with original positions
			nodes.set(tempPipeline.nodes);
			
			// Restore connections data immediately
			if (Array.isArray(tempPipeline.connections)) {
				const validConnections = tempPipeline.connections.filter(conn => {
					const sourceExists = tempPipeline.nodes.some(node => node.id === conn.sourceNodeId);
					const targetExists = tempPipeline.nodes.some(node => node.id === conn.targetNodeId);
					return sourceExists && targetExists;
				});
				physicalConnections = validConnections;
			}

			// Wait for nodes to be fully rendered before restoring visual connections
			setTimeout(() => {
				// Restore visual connections using DOM events (like project hub)
				tempPipeline.connections.forEach(conn => {
					const start_node = document.getElementById(conn.sourceAnchor);
					const end_node = document.getElementById(conn.targetAnchor);

					if (start_node && end_node) {
						start_node.dispatchEvent(new MouseEvent("mousedown", { bubbles: true, button: 0 }));
						end_node.dispatchEvent(new MouseEvent("mouseup", { bubbles: true, button: 0 }));
					}
				});
				
				// Short delay to ensure connections are visually established
				setTimeout(() => {
					isRecentering = false;
					AddToast('Canvas recentered successfully!', 'success');
				}, 150);
			}, 300); // Increased delay for better node rendering
		}, 150);
	}

	    async function handleClearCanvasConfirm() {
        const userId = sessionStorage.getItem('user_id');
        const project = get(projectName);
        const currentNodes = get(nodes);
        
        // Check if canvas is completely empty (no nodes)
        if (currentNodes.length === 0) {
            console.log('Canvas is already empty, no need to clear project');
            showClearCanvasModal = false;
            AddToast('Canvas is already empty!', 'info');
            return;
        }
        
        // Clear local canvas first
        clearLocalCanvas();
        resetAllStores();
        
        // If user is logged in and project exists, clear the project on server
        if (userId && project && !isGuestUser) {
            try {
                console.log('Clearing project on server by delete and recreate...');
                await clearProject(project, userId);
                
                AddToast(`Project "${project}" cleared and recreated successfully!`, 'success');
            } catch (error) {
                console.error('Error clearing project on server:', error);
                AddToast(`Failed to clear project on server: ${error.message}. Local canvas cleared only.`, 'warning');
            }
        } else {
            // Just show success for local clearing (guest users or no project)
            AddToast('Canvas cleared successfully!', 'success');
        }
    }

	function clearLocalCanvas() {
        // Clear all nodes and connections
        nodes.set([]);
        physicalConnections = [];
        invalid_connections = [];
        
        // Reset the pipeline store
        pipelineStore.update(pipeline => ({
            ...pipeline,
            nodes: [],
            connections: []
        }));

        // Reset node counter
        node_counter = 0;

        // Reset the toolbox created nodes
        const event = new CustomEvent('resetToolbox');
        document.dispatchEvent(event);

        // Reset last saved state to reflect the cleared canvas
        lastSavedState = JSON.stringify({
            nodes: [],
            connections: []
        });

        showClearCanvasModal = false;
    }

	function resetAllStores() {
        console.log('Resetting all phase stores...');
        
        // Reset all phase stores
        resetLexerState();
        resetParserState(); 
        resetAnalyserState();
        resetTranslatorState();
        resetSourceCode();
        
        // Reset phase completion status
        phase_completion_status.set({
            source: false,
            lexer: false,
            parser: false,
            analyser: false,
            translator: false
        });

        // Reset local artifact variables
        show_tokens = false;
        tokens = [];
        unexpected_tokens = [];
        syntaxTreeData = null;
        artifactData = null;
        parsing_error = false;
        parsing_error_details = '';
        show_symbol_table = false;
        symbol_table = [];
        analyser_error = false;
        analyser_error_details = '';
        translated_code = [];
        translationError = null;
        source_code = '';
        
        console.log('All stores reset successfully');
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
		optimiser: 'Advanced optimisation techniques for code enhancement.'
	};

	const node_labels: Record<NodeType, string> = {
		source: 'Source Code',
		lexer: 'Lexer',
		parser: 'Parser',
		analyser: 'Analyser',
		translator: 'Translator',
		optimiser: 'Optimiser'
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
		sessionStorage.setItem('hasSeenDragTip', 'true');
		show_drag_tip = false;
	}

	function handlePhaseSelect(type: NodeType) {
		// Validate node access before proceeding
		if (type !== 'optimiser' && !validateNodeAccess(type)) {
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
			// Only check for source code on non-optimiser phases
			if (type !== 'optimiser') {
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

	function handleSymbolGeneration(data: { 
        symbol_table: any[], 
        analyser_error?: string, 
        analyser_error_details?: string 
    }) {
        console.log('Symbol table generation received:', data);
        
        if (data.symbol_table && data.symbol_table.length > 0) {
            symbol_table = data.symbol_table;
            show_symbol_table = true;
            analyser_error = false;
            analyser_error_details = '';
            
            // Mark analyser phase as complete when symbol table is generated successfully
            phase_completion_status.update(status => ({
                ...status,
                analyser: true
            }));
            completion_status.analyser = true;
            
        } else {
            // FIX: Clear symbol table if empty array is received
            symbol_table = [];
            show_symbol_table = false;
            analyser_error = !!data.analyser_error;
            analyser_error_details = data.analyser_error_details || '';
            console.log('Symbol table cleared');
        }
    }

	function handleTranslationReceived(event: { detail: string[] }) {
        console.log('Translation received:', event.detail);
        
        if (event.detail && event.detail.length > 0) {
            translated_code = event.detail;
            translationError = null;
            
            // Mark translator phase as complete when translation is received
            phase_completion_status.update(status => ({
                ...status,
                translator: true
            }));
            completion_status.translator = true;
            
        } else {
            // FIX: Clear translated code if empty array is received
            translated_code = [];
            translationError = null;
            console.log('Translated code cleared');
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
					{#if !isGuestUser}
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
					{/if}
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
		<DrawerCanvas {nodes} initialConnections={physicalConnections} {tooltips} onPhaseSelect={handlePhaseSelect} onConnectionChange={handleConnectionChange} />

		<!-- Recenter Loading Overlay -->
		{#if isRecentering}
			<div class="recenter-loading-overlay">
				<div class="recenter-loading-content">
					<div class="recenter-spinner"></div>
					<span class="recenter-loading-text">Recentering canvas...</span>
				</div>
			</div>
		{/if}

		<!-- Recenter Button -->
		{#if $nodes.length > 0 && !isRecentering}
			   <button class="recenter-button icon-only" on:click={recenterCanvas} aria-label="Recenter Canvas" title="Recenter all nodes and connections">
				   <svg width="32" height="32" viewBox="0 0 296.991 296.991" fill="#ffffff" xmlns="http://www.w3.org/2000/svg" stroke="#ffffff">
					   <path d="M281.991,133.496h-13.445c-6.8-54.74-50.31-98.25-105.05-105.05V15c0-8.284-6.716-15-15-15c-8.284,0-15,6.716-15,15v13.446 c-54.74,6.8-98.25,50.31-105.05,105.05H15.001c-8.284,0-15,6.716-15,15c0,8.284,6.716,15,15,15h13.445 c6.8,54.74,50.31,98.25,105.05,105.05v13.445c0,8.284,6.716,15,15,15c8.284,0,15-6.716,15-15v-13.445 c54.74-6.8,98.25-50.31,105.05-105.05h13.445c8.284,0,15-6.716,15-15C296.991,140.212,290.275,133.496,281.991,133.496z M163.496,238.232V191c0-8.284-6.716-15-15-15c-8.284,0-15,6.716-15,15v47.232c-38.172-6.36-68.376-36.564-74.736-74.736h47.231 c8.284,0,15-6.716,15-15c0-8.284-6.716-15-15-15H58.76c6.36-38.172,36.564-68.376,74.736-74.736v47.231c0,8.284,6.716,15,15,15 c8.284,0,15-6.716,15-15V58.76c38.172,6.36,68.376,36.564,74.736,74.736h-47.231c-8.284,0-15,6.716-15,15c0,8.284,6.716,15,15,15 h47.231C231.872,201.668,201.667,231.872,163.496,238.232z"></path>
				   </svg>
			   </button>
		{/if}

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

					{#if selected_phase === 'optimiser' && OptimiserPhaseTutorial}
						<svelte:component this={OptimiserPhaseTutorial} />
						<svelte:component this={OptimiserPhaseInspector} />
						<svelte:component this={OptimiserArtifactViewer} />
					{/if}
				</div>
			</div>
		</div>
	{/if}

	<!-- Return to Canvas Button (positioned to the left of AI Assistant) -->
	{#if selected_phase}
		<button on:click={returnToCanvas} class="return-button" class:ai-open={$aiAssistantOpen} aria-label="Return to Canvas" title="Return to Canvas">
			<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path d="M19 12H5M5 12L12 19M5 12L12 5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
			</svg>
		</button>
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
	hasNodes={$nodes.length > 0}
    projectName={currentProjectName}
/>

<!-- Guest Welcome Popup -->
<GuestWelcomePopup 
	bind:show={showGuestWelcomePopup} 
	on:close={handleGuestWelcomeClose}
/>

<!-- AI Assistant Component -->
<AiAssistant />

<!-- Canvas Tutorial Modal -->
<CanvasTutorial 
	bind:show={canvasTutorialVisible} 
	on:close={handleTutorialClose}
/>


<style>

	:global(path.edge.invalid),
	:global(path.target.invalid) {
	stroke: red !important;
	stroke-width: 3px !important;
	}

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
		transition: all 0.2s ease;
	}
	.clear-button:hover {
		background-color: #fef2f2;
		color: #b91c1c;
		transform: scale(1.1);
	}
	   .recenter-button {
		   position: absolute;
		   bottom: 1rem;
		   left: 1rem;
		   width: 56px;
		   height: 56px;
		   background: linear-gradient(135deg, #001A6E 0%, #041a47 100%);
		   border: none;
		   border-radius: 12px;
		   cursor: pointer;
		   box-shadow: 0 2px 8px rgba(4, 26, 71, 0.15);
		   transition: all 0.2s ease;
		   z-index: 100;
		   display: flex;
		   align-items: center;
		   justify-content: center;
		   padding: 0;
	   }
   .recenter-button.icon-only svg {
	   display: block;
	   margin: auto;
	   width: 32px;
	   height: 32px;
   }
	.recenter-button {
		position: absolute;
		bottom: 1rem;
		left: 1rem;
		background: #001A6E;
		color: white;
		border: none;
		padding: 0.75rem 1rem;
		border-radius: 8px;
		cursor: pointer;
		font-size: 0.875rem;
		font-weight: 500;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		box-shadow: 0 2px 8px rgba(4, 26, 71, 0.15);
		transition: all 0.2s ease;
		z-index: 100;
	}

	.recenter-button:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(4, 26, 71, 0.25);
		background: #001A6E;
	}

	.recenter-button svg {
		flex-shrink: 0;
	}

	.recenter-loading-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(240, 242, 245, 0.8);
		backdrop-filter: blur(2px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 200;
		border-radius: 12px;
	}

	.recenter-loading-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		background: white;
		padding: 2rem;
		border-radius: 12px;
		box-shadow: 0 4px 20px rgba(4, 26, 71, 0.15);
	}

	.recenter-spinner {
		width: 32px;
		height: 32px;
		border: 3px solid #e5e7eb;
		border-top: 3px solid #001A6E;
		border-radius: 50%;
		animation: recenter-spin 1s linear infinite;
	}

	.recenter-loading-text {
		color: #041a47;
		font-weight: 500;
		font-size: 0.875rem;
	}

	@keyframes recenter-spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
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
		bottom: 2rem;
		right: 6.3rem;
		width: 56px;
		height: 56px;
		background: #bed2e6;
		color: #041a47;
		border: none;
		border-radius: 12px;
		cursor: pointer;
		z-index: 2000;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 2px 8px rgba(4, 26, 71, 0.15);
		transition: all 0.2s ease;
	}
	.return-button.ai-open {
		right: calc(6.3rem + 330px); /* Move left by AI assistant width (380px) + some margin */
	}
	.return-button:hover {
		background: #a8bdd1;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(4, 26, 71, 0.2);
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
		padding: 0;
		border-radius: 16px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
		max-width: 800px;
		width: 85vw;
		max-height: 85vh;
		height: auto;
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		border: 1px solid rgba(226, 232, 240, 0.8);
	}

	.close-btn {
		position: absolute;
		top: 1rem;
		right: 1rem;
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		z-index: 1;
		padding: 0.5rem;
		border-radius: 4px;
		transition: background-color 0.2s ease;
	}

	.close-btn:hover {
		background: rgba(0, 0, 0, 0.1);
	}

	/* Responsive design for code input modal */
	@media (max-width: 768px) {
		.code-input-modal {
			max-width: 95vw;
			width: 95vw;
			max-height: 85vh;
			padding: 1.5rem;
		}
	}

	@media (max-width: 480px) {
		.code-input-modal {
			max-width: 98vw;
			width: 98vw;
			max-height: 90vh;
			padding: 1rem;
			border-radius: 8px;
		}
	}

	.help-tip {
		position: absolute;
		top: 20px;
		right: 20px;
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
		white-space: nowrap;
		max-width: max-content;
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
		flex-shrink: 0;
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
		color: #fca5a5;
		transform: scale(1.1);
	}

	:global(html.dark-mode) .recenter-button {
		background: #1a2a4a;
		color: #e2e8f0;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
	}

	:global(html.dark-mode) .recenter-button:hover {
		background: linear-gradient(135deg, #2a3a5a 0%, #1a2a4a 100%);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
	}

	:global(html.dark-mode) .recenter-loading-overlay {
		background: rgba(22, 24, 35, 0.8);
	}

	:global(html.dark-mode) .recenter-loading-content {
		background: #1a2a4a;
		color: #e2e8f0;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
	}

	:global(html.dark-mode) .recenter-spinner {
		border: 3px solid #4a5568;
		border-top: 3px solid #64b5f6;
	}

	:global(html.dark-mode) .recenter-loading-text {
		color: #e2e8f0;
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
	
	:global(html.dark-mode) .close-btn:hover {
		background: rgba(255, 255, 255, 0.1);
	}
	:global(html.dark-mode) .return-button {
		background: #1a3a7a;
		color: #cccccc;
		box-shadow: 0 2px 8px rgba(26, 58, 122, 0.3);
	}
	:global(html.dark-mode) .return-button.ai-open {
		right: calc(6.3rem + 330px); /* Move left by AI assistant width (380px) + some margin */
	}
	:global(html.dark-mode) .return-button:hover {
		background: #2a4a8a;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(26, 58, 122, 0.4);
	}
</style>
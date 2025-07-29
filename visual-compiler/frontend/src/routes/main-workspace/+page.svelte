<script lang="ts">
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import type { NodeType, Token, SyntaxTree } from '$lib/types';
	import { AddToast } from '$lib/stores/toast';
	import { theme } from '../../lib/stores/theme';
	import NavBar from '$lib/components/main/nav-bar.svelte';
	import Toolbox from '$lib/components/main/toolbox.svelte';
	import CodeInput from '$lib/components/main/code-input.svelte';
	import DrawerCanvas from '$lib/components/main/drawer-canvas.svelte';

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

	let workspace_el: HTMLElement;
	let show_drag_tip = false;

	onMount(async () => {
		LexerPhaseTutorial = (await import('$lib/components/lexer/lexer-phase-tutorial.svelte')).default;
		LexerPhaseInspector = (await import('$lib/components/lexer/phase-inspector.svelte')).default;
		LexerArtifactViewer = (await import('$lib/components/lexer/lexer-artifact-viewer.svelte')).default;
		ParserPhaseTutorial = (await import('$lib/components/parser/parser-phase-tutorial.svelte')).default;
		ParserPhaseInspector = (await import('$lib/components/parser/parsing-input.svelte')).default;
		ParserArtifactViewer = (await import('$lib/components/parser/parser-artifact-viewer.svelte')).default;
		AnalyserPhaseTutorial = (await import('$lib/components/analyser/analyser-phase-tutorial.svelte')).default;
		AnalyserPhaseInspector = (await import('$lib/components/analyser/analyser-phase-inspector.svelte')).default;
		AnalyserArtifactViewer = (await import('$lib/components/analyser/analyser-artifact-viewer.svelte')).default;
		TranslatorPhaseTutorial = (await import('$lib/components/translator/translator-phase-tutorial.svelte')).default;
		TranslatorPhaseInspector = (await import('$lib/components/translator/translator-phase-inspector.svelte')).default;
		TranslatorArtifactViewer = (await import('$lib/components/translator/translator-artifact-viewer.svelte')).default;

		document.documentElement.classList.toggle('dark-mode', $theme === 'dark');
		if (!localStorage.getItem('hasSeenDragTip')) {
			show_drag_tip = true;
		}
	});

	// --- CANVAS STATE ---
	interface CanvasNode {
		id: string;
		type: NodeType;
		label: string;
		position: { x: number; y: number };
	}
	const nodes = writable<CanvasNode[]>([]);
	let node_counter = 0;
	let selected_phase: NodeType | null = null;
	let show_code_input = false;
	let source_code = '';
	let show_tokens = false;
	let syntaxTreeData: SyntaxTree | null = null;

	// --- TOOLTIPS AND LABELS ---
	const tooltips: Record<NodeType, string> = {
		source: 'Start here. Add source code to begin compilation.',
		lexer: 'Converts source code into tokens for processing.',
		parser: 'Analyzes the token stream to build a syntax tree.',
		analyser: 'Performs semantic analysis on the syntax tree.',
		translator: 'Translates the syntax tree into target code.'
	};

	const node_labels: Record<NodeType, string> = {
		source: 'Source Code',
		lexer: 'Lexer',
		parser: 'Parser',
		analyser: 'Analyser',
		translator: 'Translator'
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
		syntaxTreeData = null; // Reset tree on new phase selection
		if (type === 'source') {
			show_code_input = true;
		} else {
			selected_phase = type;
			if (!source_code.trim()) {
				AddToast('Please enter source code before proceeding', 'error');
				selected_phase = null;
				return;
			}
		}
	}

	let showSymbolTable = false;

	 function handleReset() {
    showSymbolTable = false;
  }

	function returnToCanvas() {
		selected_phase = null;
		show_code_input = false;
	}

	function handleCodeSubmit(code: string) {
		show_tokens = false;
		source_code = code;
		show_code_input = false;
	}

	function handleTokenGeneration(data: { tokens: Token[]; unexpected_tokens: string[] }) {
		show_tokens = true;
		tokens = data.tokens;
		unexpected_tokens = data.unexpected_tokens;
	}

	function handleTreeReceived(event: CustomEvent<SyntaxTree>) {
		syntaxTreeData = event.detail;
	}

	let tokens: Token[] = [];
	let unexpected_tokens: string[] = [];
</script>

<NavBar />

<div class="main">
	<!-- svelte-ignore component_name_lowercase -->
	<!-- svelte-ignore element_invalid_self_closing_tag -->
	<Toolbox {handleCreateNode} {tooltips} />
	<div class="workspace" bind:this={workspace_el} tabindex="-1">
		<DrawerCanvas {nodes} onPhaseSelect={handlePhaseSelect} />

		{#if show_drag_tip}
			<div class="help-tip">
				<span
					><b>Pro-Tip:</b> For the smoothest experience, click to select a node before dragging it.</span
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
						><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"
						></line></svg
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
					/>
					<svelte:component this={ParserArtifactViewer} syntaxTree={syntaxTreeData} />
				{/if}

				{#if selected_phase === 'analyser' && AnalyserPhaseTutorial}
					<svelte:component this={AnalyserPhaseTutorial} />
					 <AnalyserPhaseInspector
						bind:showSymbolTable
						on:generate={() => showSymbolTable = true}
						on:reset={handleReset}
					/>
					<svelte:component this={AnalyserArtifactViewer} {showSymbolTable} on:close={() => showSymbolTable = false} />
				{/if}

				{#if selected_phase === 'translator' && TranslatorPhaseTutorial}
					<svelte:component this={TranslatorPhaseTutorial} />
					<svelte:component this={TranslatorPhaseInspector} {source_code} />
					<svelte:component this={TranslatorArtifactViewer} />
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
		background: #001a6e;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		z-index: 1000;
		margin-right: 1rem;
		margin-bottom: 1rem;
	}
	.return-button:hover {
		background: #074799;
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

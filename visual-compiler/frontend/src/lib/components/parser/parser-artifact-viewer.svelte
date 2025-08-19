<script lang="ts">
	// @ts-nocheck
	import { onMount, afterUpdate } from 'svelte';
	import { lexerState } from '$lib/stores/lexer';
	import {
		Network,
		DataSet,
		type Node,
		type Edge,
		type Options,
		type IdType
	} from 'vis-network/standalone';
	import type { SyntaxTree, SyntaxTreeNode } from '$lib/types';
	import { fade, scale } from 'svelte/transition';

	export let parsing_error: boolean ;
	export let parsing_error_details: string;

	// Exported prop for the syntax tree data, following PascalCase convention.
	export let syntaxTree: SyntaxTree | null = null;

	let tree_container: HTMLElement;
	let network_instance: Network | null = null;

	// New state for the expandable modal
	let isExpanded = false;
	let expandedVisContainer: HTMLElement;

	// afterUpdate
	// Return type: void
	// Parameter type(s): none
	// Svelte lifecycle function that runs after the component's props have been updated.
	// It ensures the tree is re-rendered whenever the SyntaxTree prop changes.
	afterUpdate(() => {
		if (syntaxTree && tree_container) {
			renderTree();
		} else if (!syntaxTree && network_instance) {
			network_instance.destroy();
			network_instance = null;
		}
	});

	// onMount
	// Return type: void
	// Parameter type(s): none
	// Svelte lifecycle function that runs when the component is first created.
	// It renders the initial tree if data is available.
	onMount(() => {
		if (syntaxTree) {
			renderTree();
		}
	});

	// renderTree
	// Return type: void
	// Parameter type(s): none
	// Renders the syntax tree visualization using the vis-network library.
	function renderTree() {
		if (!syntaxTree || !syntaxTree.root) {
			if (network_instance) network_instance.destroy();
			return;
		}

		const nodes_dataset = new DataSet<Node>();
		const edges_dataset = new DataSet<Edge>();
		let node_id_counter = 0;

		// buildNodesAndEdges
		// Return type: void
		// Parameter type(s): node (SyntaxTreeNode), parent_id (IdType | null)
		// Recursively traverses the tree data to build nodes and edges for vis-network.
		function buildNodesAndEdges(node: any, parent_id: IdType | null) {
			const current_id = node_id_counter++;

			let label_text = node.symbol;
			if (node.value && node.value !== '') {
				label_text += `\n(${node.value})`;
			}

			nodes_dataset.add({
				id: current_id,
				label: label_text,
				shape: 'box',
				color: node.children ? '#E6F0FF' : '#D2FFD2',
				borderWidth: 1,
				font: {
					multi: true,
					align: 'center'
				}
			});

			if (parent_id !== null) {
				edges_dataset.add({
					from: parent_id,
					to: current_id
				});
			}

			if (node.children) {
				for (const child of node.children) {
					buildNodesAndEdges(child, current_id);
				}
			}
		}

		buildNodesAndEdges(syntaxTree.root, null);

		const options: Options = {
			layout: {
				hierarchical: {
					direction: 'UD',
					sortMethod: 'directed',
					shakeTowards: 'roots',
					levelSeparation: 80,
					nodeSpacing: 150
				}
			},
			edges: {
				smooth: {
					enabled: true,
					type: 'cubicBezier',
					forceDirection: 'vertical',
					roundness: 0.4
				}
			},
			physics: false,
			interaction: {
				dragNodes: true,
				dragView: true,
				zoomView: true
			}
		};

		if (network_instance) {
			network_instance.destroy();
		}

		const container = isExpanded ? expandedVisContainer : tree_container;
		if (container) {
			network_instance = new Network(
				container,
				{ nodes: nodes_dataset, edges: edges_dataset },
				options
			);
		}
	}

	// Functions for the expandable modal
	const toggleExpand = () => {
		isExpanded = !isExpanded;
		if (isExpanded) {
			// Re-render the tree in the modal after a short delay for the animation
			setTimeout(() => {
				renderTree();
			}, 50);
		} else {
			// Re-render the tree in the original container
			setTimeout(() => {
				renderTree();
			}, 50);
		}
	};

	const handleKeydown = (event: KeyboardEvent) => {
		if (isExpanded && event.key === 'Escape') {
			toggleExpand();
		}
	};

	function fitGraphToModal() {
		if (network_instance) {
			network_instance.fit();
		}
	}

	// Zoom functionality for expanded view
	function zoomIn() {
		if (network_instance && isExpanded) {
			const scale = network_instance.getScale();
			network_instance.moveTo({
				scale: Math.min(scale * 1.2, 3) // Max zoom level of 3x
			});
		}
	}

	function zoomOut() {
		if (network_instance && isExpanded) {
			const scale = network_instance.getScale();
			network_instance.moveTo({
				scale: Math.max(scale * 0.8, 0.1) // Min zoom level of 0.1x
			});
		}
	}

	function resetZoom() {
		if (network_instance && isExpanded) {
			network_instance.fit({
				animation: {
					duration: 500,
					easingFunction: 'easeInOutCubic'
				}
			});
		}
	}

	// Add a subscription to lexerState
	lexerState.subscribe(state => {
		if (state.parser_data?.tree) {
			syntaxTree = state.parser_data.tree;
			parsing_error = false;
			parsing_error_details = "";
			// Force re-render of the tree
			if (tree_container) {
				renderTree();
			}
		}
	});
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="artifact-container">
	<div class="artifact-header">
		<h2 class="artifact-title">Parser Artefact</h2>
	</div>

	<div class="artifact-viewer">
		{#if parsing_error}
			<!-- Error state is shown, no header needed here -->
			<div class="error-state">
				<div class="error-icon">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="48"
						height="48"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10" />
						<line x1="12" y1="8" x2="12" y2="12" />
						<line x1="12" y1="16" x2="12.01" y2="16" />
					</svg>
				</div>
				<h4>Syntax Error</h4>
				<p class="error-message">
					The source code could not be parsed with the provided tokens and grammar. Please check your
					input again.
				</p>
			</div>
		{:else}
			<!-- No error, so show the header and then the tree or empty state -->
			<div class="artifact-header">
				<h3>Syntax Tree</h3>
			</div>

			{#if syntaxTree && syntaxTree.root}
				<div class="tree-display-wrapper">
					<button on:click={toggleExpand} class="expand-btn" title="Expand view">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							fill="currentColor"
							viewBox="0 0 16 16"
						>
							<path
								d="M1.5 1a.5.5 0 0 0-.5.5v4a.5.5 0 0 1-1 0v-4A1.5 1.5 0 0 1 1.5 0h4a.5.5 0 0 1 0 1zM10 .5a.5.5 0 0 1 .5-.5h4A1.5 1.5 0 0 1 16 1.5v4a.5.5 0 0 1-1 0v-4a.5.5 0 0 0-.5-.5h-4a.5.5 0 0 1-.5-.5M.5 10a.5.5 0 0 1 .5.5v4a.5.5 0 0 0 .5.5h4a.5.5 0 0 1 0 1h-4A1.5 1.5 0 0 1 0 14.5v-4a.5.5 0 0 1 .5-.5m15 0a.5.5 0 0 1 .5.5v4a1.5 1.5 0 0 1-1.5 1.5h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 0 .5-.5v-4a.5.5 0 0 1 .5-.5"
							/>
						</svg>
					</button>
					<div bind:this={tree_container} class="tree-display" data-testid="tree-display-container" />
				</div>
			{:else}
				<div class="empty-state">
					The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.
				</div>
			{/if}
		{/if}
	</div>
</div>

<!-- Expandable Modal -->
{#if isExpanded}
	<div class="modal-backdrop" on:click={toggleExpand} transition:fade={{ duration: 200 }}>
		<div
			class="modal-content"
			on:click|stopPropagation
			transition:scale={{ duration: 250, start: 0.95 }}
			on:introend={fitGraphToModal}
		>
			<div class="modal-header">
				<div class="header-left">
					<h3>Expanded Syntax Tree</h3>
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
							viewBox="0 0 16 16"
						>
							<path
								d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8z"
							/>
						</svg>
					</button>
				</div>
			</div>
			<div class="modal-body" bind:this={expandedVisContainer}>
				<!-- Expanded vis-network graph will be rendered here -->
			</div>
		</div>
	</div>
{/if}

<style>
	.artifact-container {
		padding: 1.5rem;
		background-color: #f8f9fa;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		gap: 1rem;
		height: 100%;
	}

	.artifact-header {
		border-bottom: 1px solid #e5e7eb;
		padding-bottom: 0.75rem;
	}

	.artifact-viewer .artifact-header {
		border-bottom: none;
		padding-bottom: 0;
		margin-bottom: 1.5rem;
		text-align: center;
	}

	.artifact-title {
		margin: 0;
		color: #001a6e;
		font-family: 'Times New Roman', serif;
		font-size: 1.25rem;
	}

	h3 {
		color: #001a6e;
		font-size: 1.5rem;
		margin: 0;
		font-family: 'Times New Roman', serif;
	}

	.tree-display-wrapper {
		position: relative;
		width: 100%;
		height: 600px;
	}

	.tree-display {
		width: 100%;
		height: 100%;
		border: 1px solid #ddd;
		border-radius: 8px;
		background: #fdfdfd;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 1rem;
		color: #6b7280;
		height: 100%;
		flex-grow: 1;
		text-align: center;
		font-style: italic;
	}

	.empty-state svg {
		margin-bottom: 1rem;
		opacity: 0.5;
	}

	/* --- Error Styling --- */
	.error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		padding: 2rem;
		flex-grow: 1;
		background-color: #fff5f5;
		border: 1px solid #e53e3e;
		border-radius: 8px;
		color: #9b2c2c;
	}
	.error-icon {
		color: #e53e3e;
		margin-bottom: 1rem;
	}
	.error-state h4 {
		margin: 0 0 0.5rem 0;
		font-size: 1.25rem;
		color: #c53030;
	}
	.error-message {
		margin: 0 0 1rem 0;
		max-width: 450px;
		line-height: 1.6;
	}
	.error-details {
		background-color: #fed7d7;
		padding: 0.75rem 1rem;
		border-radius: 6px;
		font-family: 'Fira Code', monospace;
		font-size: 0.85rem;
		white-space: pre-wrap;
		word-wrap: break-word;
		max-width: 100%;
		text-align: left;
		color: #742a2a;
	}

	/* --- Dark Mode --- */
	:global(html.dark-mode) .artifact-container {
		background: #2d3748;
	}
	:global(html.dark-mode) .artifact-header {
		border-bottom-color: #4a5568;
	}
	:global(html.dark-mode) .artifact-title,
	:global(html.dark-mode) h3 {
		color: #ebeef1;
	}
	:global(html.dark-mode) .tree-display {
		border-color: #4b5563;
		background: #2d3748;
	}
	:global(html.dark-mode) .empty-state {
		color: #9ca3af;
	}
	:global(html.dark-mode) .error-state {
		background-color: #2d3748;
		border-color: #e53e3e;
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-icon {
		color: #fca5a5;
	}
	:global(html.dark-mode) .error-state h4 {
		color: #fc8181;
	}
	:global(html.dark-mode) .error-details {
		background-color: #4a2d2d;
		color: #fed7d7;
	}

	/* New styles for modal and expand button */
	.expand-btn {
		background: #ffffff;
		border: 1px solid #e5e7eb;
		cursor: pointer;
		color: #374151;
		padding: 0.25rem;
		border-radius: 50%;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s;
		position: absolute;
		top: 0.75rem;
		right: 0.75rem;
		z-index: 10;
	}
	.expand-btn:hover {
		background-color: #f3f4f6;
		border-color: #d1d5db;
	}
	:global(html.dark-mode) .expand-btn {
		background: #374151;
		border-color: #4b5563;
		color: #d1d5db;
	}
	:global(html.dark-mode) .expand-btn:hover {
		background-color: #4b5563;
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
		z-index: 1;
	}
</style>

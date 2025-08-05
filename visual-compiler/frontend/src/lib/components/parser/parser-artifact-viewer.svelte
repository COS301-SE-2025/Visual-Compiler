<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import {
		Network,
		DataSet,
		type Node,
		type Edge,
		type Options,
		type IdType
	} from 'vis-network/standalone';
	import type { SyntaxTree, SyntaxTreeNode } from '$lib/types';

	export let parsingError: any = null;

	// Exported prop for the syntax tree data, following PascalCase convention.
	export let syntaxTree: SyntaxTree | null = null;

	let tree_container: HTMLElement;
	let network_instance: Network | null = null;

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
		function buildNodesAndEdges(node: SyntaxTreeNode, parent_id: IdType | null) {
			const current_id = node_id_counter++;

			let label_text = node.symbol;
			if (node.value) {
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
		network_instance = new Network(
			tree_container,
			{ nodes: nodes_dataset, edges: edges_dataset },
			options
		);
	}
</script>

<div class="artifact-container">
	<div class="artifact-header">
		<h2 class="artifact-title">Parser Artifact</h2>
	</div>

	<div class="artifact-viewer">
		{#if parsingError}
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
				<h4>Parsing Failed</h4>
				<p class="error-message">
					The source code could not be parsed with the provided tokens and grammar. Please check your input again.
				</p>
				<pre class="error-details">{parsingError.message || String(parsingError)}</pre>
			</div>
		{:else}
			<!-- No error, so show the header and then the tree or empty state -->
			<div class="artifact-header">
				<h3>Syntax Tree</h3>
			</div>

			{#if syntaxTree && syntaxTree.root}
				<div bind:this={tree_container} class="tree-display" data-testid="tree-display-container" />
			{:else}
				<div class="empty-state">
					The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.
				</div>
			{/if}
		{/if}
	</div>
</div>

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

	.tree-display {
		width: 100%;
		height: 600px;
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
</style>

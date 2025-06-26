<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import { Network, DataSet, type Node, type Edge, type Options, type IdType } from 'vis-network/standalone';
	import type { SyntaxTree, SyntaxTreeNode } from '$lib/types';

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
		network_instance = new Network(tree_container, { nodes: nodes_dataset, edges: edges_dataset }, options);
	}
</script>

<div class="artifact-viewer">
	<div class="artifact-header">
		<h3>Artifacts: Syntax Tree</h3>
	</div>

	{#if syntaxTree && syntaxTree.root}
		<div bind:this={tree_container} class="tree-display"></div>
	{:else}
		<div class="empty-state">The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.</div>
	{/if}
</div>

<style>
	.artifact-viewer {
		padding: 1.5rem;
		flex-grow: 1;
		display: flex;
		flex-direction: column;
	}
	.artifact-header {
		display: flex;
		justify-content: center;
		margin-bottom: 1.5rem;
	}
	h3 {
		color: #001a6e;
		font-size: 1.5rem;
		margin: 0;
		font-family: 'Times New Roman', serif;
	}
	.empty-state {
		color: #666;
		font-style: italic;
		text-align: center;
		padding: 2rem;
		flex-grow: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.tree-display {
		width: 100%;
		height: 600px;
		border: 1px solid #ddd;
		border-radius: 8px;
		background: #fdfdfd;
	}

	
	:global(html.dark-mode) h3 {
		color: #ebeef1;
	}
	:global(html.dark-mode) .empty-state {
		color: #9ca3af;
	}
	:global(html.dark-mode) .tree-display {
		border-color: #4b5563;
		background: #2d3748;
	}
</style>

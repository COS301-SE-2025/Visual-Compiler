<script lang="ts">
	import { Svelvet, Node } from 'svelvet';
	import type { Writable } from 'svelte/store';
	import type { NodeType, NodeConnection } from '$lib/types';
	import { theme } from '../../stores/theme';

	interface CanvasNode {
		id: string;
		type: NodeType;
		label: string;
		position: { x: number; y: number };
	}

	export let nodes: Writable<CanvasNode[]>;
	export let initialConnections: NodeConnection[] = [];

	export let onPhaseSelect: (type: NodeType) => void = () => {};
	export let onConnectionChange: (connections: NodeConnection[]) => void = () => {};

	// Track physical connections between nodes
	let nodeConnections: NodeConnection[] = [...initialConnections];

	// Make nodeConnections reactive to changes in initialConnections
	$: nodeConnections = [...initialConnections];

	// Function to restore nodes to their original saved positions
	function restoreOriginalPositions(nodesList: CanvasNode[]): CanvasNode[] {
		// Simply return the nodes as they are - with their saved positions
		// The displayNodes computed property will handle the override
		console.log('Restoring nodes to original positions:', nodesList.map(n => ({ id: n.id, x: n.position.x, y: n.position.y })));
		return nodesList;
	}

	// Track if we need to reposition nodes (to avoid infinite updates)
	let hasRepositioned = false;
	let lastNodeCount = 0;
	let lastNodeIds: string[] = [];

	// Create a computed store that always returns properly positioned nodes
	$: displayNodes = (() => {
		const currentNodeIds = $nodes.map(n => n.id).sort();
		const nodeIdsChanged = JSON.stringify(currentNodeIds) !== JSON.stringify(lastNodeIds);
		
		if ($nodes.length > 0 && nodeIdsChanged) {
			console.log('Nodes changed, ensuring positions are respected by Svelvet:', {
				nodeCount: $nodes.length,
				nodeIdsChanged,
				currentPositions: $nodes.map(n => ({ id: n.id, x: n.position.x, y: n.position.y }))
			});
			
			const restoredNodes = restoreOriginalPositions($nodes);
			console.log('Restoring positions for Svelvet:', restoredNodes.map(n => ({ id: n.id, x: n.position.x, y: n.position.y })));
			
			// Update the tracking variables
			hasRepositioned = true;
			lastNodeCount = $nodes.length;
			lastNodeIds = currentNodeIds;
			
			// Also update the original store to keep it in sync
			setTimeout(() => {
				nodes.set(restoredNodes);
			}, 10);
			
			return restoredNodes;
		} else if ($nodes.length === 0) {
			hasRepositioned = false;
			lastNodeCount = 0;
			lastNodeIds = [];
		}
		
		return $nodes;
	})();

	// Create a key that changes when display nodes change
	$: canvasKey = displayNodes.length > 0 ? `${displayNodes.length}-${displayNodes.map(n => `${n.id}-${n.position.x}-${n.position.y}`).join('|')}` : 'empty';

	let canvas_el: any;
	let last_click = -Infinity;
	const DOUBLE_CLICK_MILLISECONDS = 300;



	// Handle new edge connections
	function handleConnection(event: CustomEvent) {
		const { sourceNode, targetNode } = event.detail;

		const sourceNodeId = sourceNode.id.replace("N-", "");
		const targetNodeId = targetNode.id.replace("N-", "");

		//const sourceNodeId = sourceNode.id.replace("N-", "");
		const sourceCanvasNode = displayNodes.find(node => node.id === sourceNodeId);
		const targetCanvasNode = displayNodes.find(node => node.id === targetNodeId);

		if (sourceCanvasNode && targetCanvasNode) {
			const newConnection: NodeConnection = {
				id: `${sourceNodeId}-${targetNodeId}`,
				sourceNodeId: sourceCanvasNode.id,
				targetNodeId: targetCanvasNode.id,
				sourceType: sourceCanvasNode.type,
				targetType: targetCanvasNode.type,
				sourceAnchor: event.detail.sourceAnchor.id,
				targetAnchor: event.detail.targetAnchor.id,
			};

			nodeConnections = [...nodeConnections, newConnection];
			onConnectionChange(nodeConnections);
		}
	}

	// Handle edge disconnections
	function handleDisconnection(event: CustomEvent) {
		console.log('Connection removed:', event.detail);
		const { sourceNode, targetNode } = event.detail;
		
		// Extract node IDs from the node objects
		const sourceNodeId = sourceNode.id.replace('N-', '');
		const targetNodeId = targetNode.id.replace('N-', '');
		
		// Find the canvas nodes to get their actual IDs
		const sourceCanvasNode = displayNodes.find(node => node.id === sourceNodeId);
		const targetCanvasNode = displayNodes.find(node => node.id === targetNodeId);
		
		if (sourceCanvasNode && targetCanvasNode) {
			nodeConnections = nodeConnections.filter(conn => 
				!(conn.sourceNodeId === sourceCanvasNode.id && conn.targetNodeId === targetCanvasNode.id) &&
				!(conn.sourceNodeId === targetCanvasNode.id && conn.targetNodeId === sourceCanvasNode.id)
			);
			onConnectionChange(nodeConnections);
			console.log('Updated connections after removal:', nodeConnections);
		}
	};

	// onNodeClick
	// Return type: void
	// Parameter type(s): NodeType
	// Dispatches a 'phaseSelect' event if a node is double-clicked within a time limit.
	function onNodeClick(type: NodeType) {
		const now = performance.now();

		// Add these three lines for debugging
		console.log('--------------------------------');
		console.log(`Click registered! now: ${now}, last_click: ${last_click}`);
		console.log(`DIFFERENCE: ${now - last_click}`);

		if (now - last_click < DOUBLE_CLICK_MILLISECONDS) {
			// Add this line for debugging
			console.log('Double-click condition was MET!');
			onPhaseSelect(type);
		}
		last_click = now;
	}
</script>

<div class="drawer-canvas">
	<div class="canvas-container" class:dark-mode={$theme === 'dark'}>
		{#key canvasKey}
		<Svelvet 
			bind:this={canvas_el} 
			theme={'custom-theme'}
			on:connection={handleConnection}
			on:disconnection={handleDisconnection}
			on:nodeMove={(event) => {
				const { node, position } = event.detail;
				const nodeId = node.id.replace('N-', '');
				const updatedNodes = $nodes.map(n => 
					n.id === nodeId 
						? { ...n, position: { x: position.x, y: position.y } }
						: n
				);
				nodes.set(updatedNodes);
			}}
		>
			{#each displayNodes as node (node.id)}
				<Node
					id={`N-${node.id}`}
					position={node.position}
					drop="center"
					bgColor={node.type === 'optimizer' ? ($theme === 'dark' ? '#8451C7' : '#AFA2D7') : ($theme === 'dark' ? '#001A6E' : '#BED2E6')}
					textColor={node.type === 'optimizer' ? '#fff' : ($theme === 'dark' ? '#ffffff' : '#000000')}
					borderColor={node.type === 'optimizer' ? ($theme === 'dark' ? '#374151' : '#FFFFFF') : ($theme === 'dark' ? '#374151' : '#FFFFFF')}
					label={node.label}
					editable={false}
					inputs={node.type !== 'source' && node.type !== 'optimizer' ? 1 : 0}
					outputs={node.type !== 'translator' && node.type !== 'optimizer' ? 1 : 0}
					on:nodeClicked={() => onNodeClick(node.type)}
				/>
			{/each}
		</Svelvet>
		{/key}
	</div>
</div>

<style>
	:root[svelvet-theme='custom-theme'] {
		--background-color: transparent;
		--dot-color: transparent;
		--node-color: #BED2E6;
		--node-text-color: #000000;
		--node-border-color: #FFFFFF;
		--node-selection-color: #3b82f6;
	}

	.drawer-canvas {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.canvas-container {
		flex: 1;
		position: relative;
		overflow: hidden;
		border-radius: 12px;
		border: 1px solid #e0e0e0;
		--grid-size: 20px;
		--grid-line-color: #e5e7eb;
		background-color: #f9fafb;
		background-image:
			linear-gradient(to right, var(--grid-line-color) 1px, transparent 1px),
			linear-gradient(to bottom, var(--grid-line-color) 1px, transparent 1px);
		background-size: var(--grid-size) var(--grid-size);
		--edge-color: #374151;
		transition:
			background-color 0.3s ease,
			border-color 0.3s ease;
	}

	/* --- Styles for Dark Mode --- */
	.canvas-container.dark-mode {
		border-color: #374151;
		--grid-line-color: #2c2f40;
		background-color: #1b1d2a;
		--edge-color: #ffffff;
	}

	.canvas-container :global(.svelvet) {
		width: 100% !important;
		height: 100% !important;
	}

	:global(.svelvet-edge path) {
		stroke: var(--edge-color) !important;
		stroke-width: 3px !important;
		transition: stroke 0.3s ease;
	}

	:global(g[id^='N-'] rect) {
		transition: filter 0.2s ease;
	}

	:global(g.selected[id^='N-'] rect) {
		filter: drop-shadow(0 0 12px white);
	}

	:global(g[id^='N-source'] .handle-left),
	:global(g[id^='source'] .handle-left) {
		display: none !important;
	}
</style>

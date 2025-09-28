<script lang="ts">
	import { Svelvet, Node, Anchor } from 'svelvet';
	import type { Writable } from 'svelte/store';
	import type { NodeType, NodeConnection } from '$lib/types';
	import { theme } from '../../stores/theme';
	import { AddToast } from '../../stores/toast';
	import { onDestroy } from 'svelte';

	interface CanvasNode {
		id: string;
		type: NodeType;
		label: string;
		position: { x: number; y: number };
	}

	export let nodes: Writable<CanvasNode[]>;
	export let initialConnections: NodeConnection[] = [];
	export let tooltips: Record<NodeType, string> = {};

	export let onPhaseSelect: (type: NodeType) => void = () => {};
	export let onConnectionChange: (connections: NodeConnection[]) => void = () => {};

	// Track physical connections between nodes
	let nodeConnections: NodeConnection[] = [...initialConnections];
	
	// Track if we're in the initial loading phase
	let isLoading = true;
	let hasInitialized = false;

	// Make nodeConnections reactive to changes in initialConnections
	$: {
		if (!hasInitialized && initialConnections.length > 0) {
			nodeConnections = [...initialConnections];
			isLoading = true;
			// Set a timeout to mark loading as complete after connections are established
			setTimeout(() => {
				isLoading = false;
				hasInitialized = true;
			}, 500);
		} else if (initialConnections.length === 0 && !hasInitialized) {
			// No initial connections, not loading
			isLoading = false;
			hasInitialized = true;
		}
	}

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

	// Create a key that only changes when we truly need to recreate the canvas
	// Avoid recreation when just adding nodes - only recreate when clearing canvas
	$: canvasKey = displayNodes.length === 0 ? 'empty-canvas' : 'active-canvas';

	// Function to add tooltips to Svelvet nodes
	function addTooltipsToNodes() {
		if (typeof document !== 'undefined') {
			displayNodes.forEach(node => {
				const nodeId = `N-${node.id}`;
				// Try multiple selectors
				let nodeElement = document.querySelector(`g[id="${nodeId}"]`);
				if (!nodeElement) {
					nodeElement = document.querySelector(`[id="${nodeId}"]`);
				}
				if (!nodeElement) {
					// Try finding by text content
					nodeElement = Array.from(document.querySelectorAll('g')).find(g => 
						g.textContent?.includes(node.label)
					);
				}
				
				if (nodeElement && tooltips[node.type]) {
					// Remove existing tooltip listeners
					nodeElement.removeEventListener('mouseenter', handleNodeMouseEnter);
					nodeElement.removeEventListener('mouseleave', handleNodeMouseLeave);
					
					// Add new tooltip listeners
					nodeElement.addEventListener('mouseenter', (e) => {
						handleNodeMouseEnter(e, tooltips[node.type]);
					});
					nodeElement.addEventListener('mouseleave', (e) => {
						handleNodeMouseLeave(e);
					});
				}
			});
		}
	}

	// Tooltip element
	let tooltipElement: HTMLDivElement | null = null;

	function createTooltipElement() {
		if (tooltipElement) return tooltipElement;
		
		tooltipElement = document.createElement('div');
		tooltipElement.className = 'canvas-node-tooltip';
		document.body.appendChild(tooltipElement);
		return tooltipElement;
	}

	function handleNodeMouseEnter(event: MouseEvent, tooltipText: string) {
		const tooltip = createTooltipElement();
		if (!tooltip) return;

		tooltip.textContent = tooltipText;
		tooltip.style.visibility = 'visible';
		tooltip.style.opacity = '1';

		// Position tooltip above the mouse cursor
		const updateTooltipPosition = (e: MouseEvent) => {
			tooltip.style.left = `${e.pageX - tooltip.offsetWidth / 2}px`;
			tooltip.style.top = `${e.pageY - tooltip.offsetHeight - 10}px`;
		};

		updateTooltipPosition(event);
		
		// Update position as mouse moves
		const target = event.currentTarget as Element;
		const mouseMoveHandler = (e: MouseEvent) => updateTooltipPosition(e);
		target.addEventListener('mousemove', mouseMoveHandler);
		
		// Store handler for cleanup
		(target as any)._tooltipMouseMoveHandler = mouseMoveHandler;
	}

	function handleNodeMouseLeave(event: MouseEvent) {
		if (tooltipElement) {
			tooltipElement.style.visibility = 'hidden';
			tooltipElement.style.opacity = '0';
		}
		
		// Clean up mousemove handler
		const target = event.currentTarget as Element;
		if ((target as any)._tooltipMouseMoveHandler) {
			target.removeEventListener('mousemove', (target as any)._tooltipMouseMoveHandler);
			delete (target as any)._tooltipMouseMoveHandler;
		}
	}

	// Add tooltips after nodes are rendered
	$: if (displayNodes.length > 0 && Object.keys(tooltips).length > 0) {
		// Use a small delay to ensure nodes are fully rendered
		setTimeout(() => addTooltipsToNodes(), 300);
	}

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
			// Skip validation if we're still loading initial connections
			if (isLoading) {
				console.log('Skipping validation during initial loading phase');
				return;
			}
			
			// Check if source node already has an outgoing connection
			const hasOutgoingConnection = nodeConnections.some(conn => 
				conn.sourceNodeId === sourceCanvasNode.id
			);
			
			// Check if target node already has an incoming connection
			const hasIncomingConnection = nodeConnections.some(conn => 
				conn.targetNodeId === targetCanvasNode.id
			);

			if (hasOutgoingConnection) {
				AddToast('Connection limit: Each node can only have one outgoing connection', 'error');
				// Manually remove the invalid connection using DOM manipulation
				setTimeout(() => {
					nodeConnections = nodeConnections.filter(
						conn => conn.id !== `${sourceNodeId}-${targetNodeId}`
					);
					//canvas_el?.cancelEdgeDrag?.();
					const startAnchor = document.getElementById(event.detail.sourceAnchor.id);
					const endAnchor = document.getElementById(event.detail.targetAnchor.id);
					// Simulate disconnection by triggering mouse events
					endAnchor?.dispatchEvent(new MouseEvent("mousedown", { bubbles: true, button: 0 }));
					startAnchor?.dispatchEvent(new MouseEvent("mouseup", { bubbles: true, button: 0 }));
				}, 50); // Small delay to ensure the connection is fully rendered before removing
				return; // Prevent the connection from being stored
			}

			if (hasIncomingConnection) {
				AddToast('Connection limit: Each node can only have one incoming connection', 'error');
			}

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

	// Cleanup tooltip element when component is destroyed
	onDestroy(() => {
		if (tooltipElement && document.body.contains(tooltipElement)) {
			document.body.removeChild(tooltipElement);
			tooltipElement = null;
		}
	});
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
					bgColor={node.type === 'optimiser' ? ($theme === 'dark' ? '#8451C7' : '#AFA2D7') : ($theme === 'dark' ? '#001A6E' : '#BED2E6')}
					textColor={node.type === 'optimiser' ? ($theme === 'dark' ? '#fff' : '#000') : ($theme === 'dark' ? '#ffffff' : '#000000')}
					borderColor={node.type === 'optimiser' ? ($theme === 'dark' ? '#374151' : '#FFFFFF') : ($theme === 'dark' ? '#374151' : '#FFFFFF')}
					label={node.label}
					editable={false}
					inputs={node.type !== 'source' && node.type !== 'optimiser' ? 1 : 0}
					outputs={node.type !== 'translator' && node.type !== 'optimiser' ? 1 : 0}
					on:nodeClicked={() => onNodeClick(node.type)}
				/>
				{/each}
						
		</Svelvet>
		{/key}
	</div>
</div>

<style>
	:global(.svelvet-anchor.disabled) {
		pointer-events: none;   
	}

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

	/* Canvas Node Tooltip Styling */
	:global(.canvas-node-tooltip) {
		position: absolute;
		visibility: hidden;
		opacity: 0;
		background-color: #333;
		color: #fff;
		font-size: 0.75rem;
		padding: 0.4rem 0.8rem;
		border-radius: 4px;
		white-space: nowrap;
		z-index: 1000;
		transition: opacity 0.2s ease;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		pointer-events: none;
		font-family: 'Times New Roman', Times, serif;
	}
</style>

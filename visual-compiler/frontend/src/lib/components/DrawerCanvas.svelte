<script lang="ts">
  import { Node, Svelvet } from 'svelvet';
  import { writable } from 'svelte/store';
  import { onMount } from 'svelte';
  import type { NodeType } from '$lib/types';

  export let sourceCode: string;

  interface CanvasNode {
    id: string;
    type: NodeType;
    label: string;
    position: { x: number; y: number };
  }

  const nodes = writable<CanvasNode[]>([]);
  let nodeCounter = 0;
  let canvas: any;

  onMount(() => {
    const handleCreateNode = (event: Event) => {
      const detail = (event as CustomEvent).detail;
      if (!detail || !detail.type) return;
      
      const type: NodeType = detail.type;
      nodeCounter += 1;
      
      nodes.update(current => [
        ...current,
        {
          id: `${type}-${nodeCounter}`,
          type,
          label: type.charAt(0).toUpperCase() + type.slice(1),
          position: { 
            x: Math.floor(Math.random() * 300) + 100, 
            y: Math.floor(Math.random() * 300) + 100 
          }
        }
      ]);
    };

    window.addEventListener('createCanvasNode', handleCreateNode);
    
    return () => {
      window.removeEventListener('createCanvasNode', handleCreateNode);
    };
  });
</script>



<div class="drawer-canvas">
  <div class="code-display">
    {sourceCode}
  </div>
  <div class="canvas-container">
    <Svelvet bind:this={canvas} >
      {#each $nodes as node (node.id)}
        <Node
          id={node.id}
          label={node.label}
          position={node.position}
          useDefaults
        />
      {/each}
    </Svelvet>
  </div>
</div>

<style>
  .drawer-canvas {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  
  .code-display {
    padding: 1rem;
    background: #f8f8f8;
    border-bottom: 1px solid #ddd;
    min-height: 100px;
  }
  
  .canvas-container {
    flex: 1;
    position: relative;
  }
</style>
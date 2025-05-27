<script lang="ts">
  import { Svelvet, Node } from 'svelvet';
  import { createEventDispatcher } from 'svelte';
  import type { Writable } from 'svelte/store';
  import type { NodeType } from '$lib/types';

  interface CanvasNode {
    id: string;
    type: NodeType;
    label: string;
    position: { x: number; y: number };
  }

  export let nodes: Writable<CanvasNode[]>;

  const dispatch = createEventDispatcher<{ phaseSelect: NodeType }>();
  let canvasEl: any;
  let lastClick = 0;
  const DOUBLE_CLICK_MS = 300;

  function onNodeClick(type: NodeType) {
    const now = performance.now();
    if (now - lastClick < DOUBLE_CLICK_MS) {
      dispatch('phaseSelect', type);
    }
    lastClick = now;
  }
</script>

<div class="drawer-canvas">
  <div class="canvas-container">
    <Svelvet bind:this={canvasEl}>
      {#each $nodes as node (node.id)}
        <Node
          id={node.id}
          label={node.label}
          position={node.position}
          drop="center"
          useDefaults
          bgColor="#041a47"
          textColor="#fff"
          on:nodeClicked={() => onNodeClick(node.type)}
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
    overflow: hidden;
  }
  .canvas-container {
    flex: 1;
    position: relative;
    overflow: auto;
  }
  .canvas-container :global(.svelvet) {
    width: 100% !important;
    height: 100% !important;
  }
  :global(.svelvet-edge path) {
    stroke: #fff !important;
    stroke-width: 3px !important;
  }
  :global(g[id^="N-"] rect) {
    transition: filter 0.2s ease;
  }
  :global(g.selected[id^="N-"] rect) {
    filter: drop-shadow(0 0 12px #fff);
  }

  :global(g[id^="N-source"] .handle-left),
:global(g[id^="source"] .handle-left) {
  display: none !important;
}

</style>
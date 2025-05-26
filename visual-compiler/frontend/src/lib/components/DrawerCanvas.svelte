<script lang="ts">
  import { Svelvet, Node } from 'svelvet';
  import { get } from 'svelte/store';
  import type { Writable } from 'svelte/store';
  import type { NodeType } from '$lib/types';
  
  export let sourceCode: string;
  export let nodes: Writable<
    { id: string; type: NodeType; label: string; position: {x:number,y:number} }[]
  >;
  
  let canvasEl: any;
</script>

<div class="drawer-canvas">
  <div class="code-display">{sourceCode}</div>
  <div class="canvas-container">
    <Svelvet bind:this={canvasEl}>
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
  .drawer-canvas { flex:1; display:flex; flex-direction:column; }
  .code-display  {
    padding:1rem; background:#f8f8f8; border-bottom:1px solid #ddd;
    min-height:100px;
  }
  .canvas-container { flex:1; position:relative; }
</style>

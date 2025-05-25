<script context="module" lang="ts">
  export type NodeType = 'lexing' | 'parsing' | 'analysis';
</script>

<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher<{
    blockClick: { type: NodeType };
  }>();

  const nodeTypes: {id: NodeType, label: string}[] = [
    { id: 'lexing', label: 'Lexing' },
    { id: 'parsing', label: 'Parsing' },
    { id: 'analysis', label: 'Analysis' }
  ];

  function createNode(type: NodeType) {
    dispatch('blockClick', { type });
  }
</script>

<div class="toolbox">
  {#each nodeTypes as type}
    <button on:click={() => createNode(type.id)} class="node-button">
      {type.label}
    </button>
  {/each}
</div>


<style>
  .toolbox {
    width: 200px;
    background: #f4f5f7;
    border-right: 1px solid #ddd;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .node-button {
    padding: 0.5rem;
    cursor: pointer;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    text-align: center;
  }
  
  .node-button:hover {
    background: #f0f0f0;
  }
</style>
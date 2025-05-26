<script lang="ts">
  import type { NodeType } from '$lib/types';
  // receive the parent’s function
  export let handleCreateNode: (type: NodeType) => void;

  const nodeTypes: { id: NodeType; label: string }[] = [
    { id: 'lexing',  label: 'Lexing'  },
    { id: 'parsing', label: 'Parsing' },
    { id: 'analysis',label: 'Analysis'}
  ];
</script>

<div class="toolbox">
  {#each nodeTypes as n}
    <button on:click={() => handleCreateNode(n.id)}>
      {n.label}
    </button>
  {/each}
</div>

<style>
  .toolbox {
    width: 200px; background:#f4f5f7; border-right:1px solid #ddd;
    padding:1rem; display:flex; flex-direction:column; gap:1rem;
  }
  button {
    padding:.5rem; border:1px solid #ddd; background:white; cursor:pointer;
    border-radius:4px; text-align:center;
  }
  button:hover { background:#f0f0f0 }
</style>
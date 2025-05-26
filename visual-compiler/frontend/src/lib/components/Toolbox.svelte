<script lang="ts">
  import type { NodeType } from '$lib/types';
  export let handleCreateNode: (type: NodeType) => void;

  const nodeTypes: { id: NodeType; label: string }[] = [
    { id: 'lexing',  label: 'Lexing'  },
    { id: 'parsing', label: 'Parsing' },
    { id: 'analysis',label: 'Analysis'}
  ];
</script>

<aside class="toolbox">
  <h2 class="toolbox-heading">Phases</h2>
  {#each nodeTypes as n}
    <button
      class="phase-btn"
      on:click={() => handleCreateNode(n.id)}
    >
      {n.label}
    </button>
  {/each}
</aside>

<style>
  .toolbox {
     flex: 0 0 clamp(240px, 20%, 300px);
    background: #fafafa;
    border-right: 1px solid #e0e0e0;
    box-shadow: 2px 0 5px rgba(0,0,0,0.05);
    padding: 1.5rem 1rem;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    gap: 1rem; 
  }

  .toolbox-heading {
    margin: 0 0 1.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    text-align: center;
    color: #041a47;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .phase-btn {
    padding: 0.75rem 1rem;
    font-size: 1rem;
    font-weight: 500;
    color: #222;
    background: #fff;
    border: 1px solid #ddd;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: 
      background 0.2s ease, 
      transform 0.1s ease, 
      box-shadow 0.2s ease;
  }

  .phase-btn:hover {
    background: #f0f4ff;
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  }

  .phase-btn:active {
    transform: translateY(0);
    background: #e8edf8;
    box-shadow: 0 2px 3px rgba(0,0,0,0.08);
  }
</style>

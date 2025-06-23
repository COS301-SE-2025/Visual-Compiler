<script lang="ts">
  import type { NodeType } from '$lib/types';
  import { theme } from '../../stores/theme';

  export let handleCreateNode: (type: NodeType) => void;
  export let tooltips: Record<NodeType, string>;

  const node_types: { id: NodeType; label: string }[] = [
    { id: 'source', label: 'Source Code' },
    { id: 'lexer', label: 'Lexer' },
    { id: 'parser', label: 'Parser' }
  ];

  // createNode
  // Return type: void
  // Parameter type(s): NodeType
  // A wrapper function that calls the handleCreateNode prop passed from the parent.
  function createNode(type: NodeType) {
    handleCreateNode(type);
  }
</script>

<aside class="toolbox">
  <h2 class="toolbox-heading">Blocks</h2>
  <h2 class="toolbox-instruction">Click a block to add it to the canvas.</h2>
  {#each node_types as n}
    <button class="phase-btn" on:click={() => createNode(n.id)}>
      {n.label}
      <span class="custom-tooltip">{tooltips[n.id]}</span>
    </button>
  {/each}
</aside>

<style>
  .toolbox {
    flex: 0 0 clamp(240px, 20%, 300px);
    background: #ffffff;
    padding: 1.5rem 1rem;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    align-items: center;
    transition: all 0.3s ease;
    margin: 0;
    border-radius: 12px;
    border: 1px solid #e0e0e0;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
  }

  .toolbox-heading {
    margin: 0 0 1.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    text-align: center;
    color: #041a47;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    transition: color 0.3s ease;
  }

  .toolbox-instruction {
    font-size: 1rem;
    font-weight: 400;
    text-align: center;
    color: #555;
    margin: -1rem 0 1rem 0;
    padding: 0 0.5rem;
    transition: color 0.3s ease;
  }

  .phase-btn {
    height: 90px;
    width: 200px;
    background-color: #041a47;
    color: white;
    border-radius: 8px;
    padding: 8px;
    font-size: 0.9rem;
    font-weight: 500;
    text-align: center;
    border: none;
    cursor: pointer;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
    position: relative;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .phase-btn:hover {
    background-color: #05276f;
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .phase-btn:active {
    transform: translateY(0);
    background: #e8edf8;
    box-shadow: 0 2px 3px rgba(0, 0, 0, 0.08);
  }

  .custom-tooltip {
    visibility: hidden;
    opacity: 0;
    position: absolute;
    bottom: 110%;
    left: 50%;
    transform: translateX(-50%);
    background: #1e1e1e;
    color: #fff;
    font-size: 0.75rem;
    padding: 6px 10px;
    border-radius: 15px;
    white-space: nowrap;
    z-index: 10;
    transition: opacity 0.2s ease;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    pointer-events: none;
  }

  .custom-tooltip::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: #1e1e1e transparent transparent transparent;
  }

  .phase-btn:hover .custom-tooltip {
    visibility: visible;
    opacity: 1;
  }

  /* Dark Mode Styles */
  :global(html.dark-mode) .toolbox {
    background: #1b1d2a;
    border-color: #374151;
    box-shadow: none;
  }

  :global(html.dark-mode) .toolbox-heading {
    color: #d3d3d3;
  }

  :global(html.dark-mode) .toolbox-instruction {
    color: #b0b0b0;
  }

  :global(html.dark-mode) .phase-btn {
    background-color: #041a47;
    color: #f0f0f0;
    border: 1px solid #374151;
  }

  :global(html.dark-mode) .phase-btn:hover {
    background-color: #2a4a8a;
  }

  :global(html.dark-mode) .phase-btn:active {
    background: #3a5a9a;
  }

  :global(html.dark-mode) .custom-tooltip {
    background: #333;
  }

  :global(html.dark-mode) .custom-tooltip::after {
    border-color: #333 transparent transparent transparent;
  }
</style>
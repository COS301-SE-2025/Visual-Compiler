<script lang="ts">
  import NavBar          from '$lib/components/NavBar.svelte';
  import Toolbox         from '$lib/components/Toolbox.svelte';
  import CodeInput       from '$lib/components/CodeInput.svelte';
  import DrawerCanvas    from '$lib/components/DrawerCanvas.svelte';
  import PhaseTutorial   from '$lib/components/PhaseTutorial.svelte';
  import PhaseInspector  from '$lib/components/PhaseInspector.svelte';
  import ArtifactViewer  from '$lib/components/ArtifactViewer.svelte';
  import type { NodeType } from '$lib/types';
  import { writable }    from 'svelte/store';

  // --- CANVAS NODES STATE ---
  interface CanvasNode {
    id:       string;
    type:     NodeType;
    label:    string;
    position: { x: number; y: number };
  }
  const nodes = writable<CanvasNode[]>([]);
  let nodeCounter = 0;

  function handleCreateNode(type: NodeType) {
    nodeCounter += 1;
    nodes.update(curr => [
      ...curr,
      {
        id:    `${type}-${nodeCounter}`,
        type,
        label: type[0].toUpperCase() + type.slice(1),
        position: {
          x: Math.random() * 300 + 100,
          y: Math.random() * 300 + 100
        }
      }
    ]);
  }

  // --- PHASE SELECTION (double-click) ---
  let selectedPhase: NodeType | null = null;
  function handlePhaseSelect(e: CustomEvent<NodeType>) {
    selectedPhase = e.detail;
  }

  // Back-to-canvas
  function returnToCanvas() {
    selectedPhase = null;
  }
</script>

<NavBar />

<div class="main">
  {#if !selectedPhase}
    <!-- edit view -->
    <Toolbox {handleCreateNode} />

    <div class="workspace">
      <CodeInput />
      <!-- listen here for your dblclick → phaseSelect event -->
      <DrawerCanvas {nodes} on:phaseSelect={handlePhaseSelect} />
    </div>
  {:else}
    <!-- analysis view -->
    <div class="analysis-view">
      <div class="three-column-layout">
        <PhaseTutorial  phase={selectedPhase} />
        <PhaseInspector phase={selectedPhase} />
        <ArtifactViewer phase={selectedPhase} />
      </div>
      <button on:click={returnToCanvas} class="return-button">
        ← Return to Canvas
      </button>
    </div>
  {/if}
</div>

<style>
  :global(html, body) {
    margin: 0; padding: 0; height: 100%; overflow: hidden;
  }
  .main {
    display: flex;
    height: calc(100vh - 3.5rem);
    overflow: hidden;
  }
  .workspace {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .analysis-view {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .three-column-layout {
    display: flex;
    flex: 1;
  }
  .return-button {
    position: fixed;
    bottom: 20px; right: 20px;
    padding: 0.5rem 1rem;
    background: #001A6E; color: white;
    border: none; border-radius: 4px;
    cursor: pointer; z-index: 1000;
  }
  .return-button:hover {
    background: #074799;
  }
</style>

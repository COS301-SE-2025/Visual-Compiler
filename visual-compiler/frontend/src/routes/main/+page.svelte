<script lang="ts">
  import NavBar          from '$lib/components/NavBar.svelte';
  import Toolbox         from '$lib/components/Toolbox.svelte';
  import CodeInput       from '$lib/components/CodeInput.svelte';
  import DrawerCanvas    from '$lib/components/DrawerCanvas.svelte';
  import PhaseTutorial   from '$lib/components/PhaseTutorial.svelte';
  import PhaseInspector  from '$lib/components/PhaseInspector.svelte';
  import ArtifactViewer  from '$lib/components/ArtifactViewer.svelte';
  import type { NodeType, CanvasEdge } from '$lib/types';
  import { writable }    from 'svelte/store';
  import { toasts }      from '$lib/stores/toast';

  // --- CANVAS STATE ---
  interface CanvasNode {
    id:       string;
    type:     NodeType;
    label:    string;
    position: { x: number; y: number };
  }
  
  const edges = writable<CanvasEdge[]>([]);
  const nodes = writable<CanvasNode[]>([]);
  let nodeCounter = 0;
  let selectedPhase: NodeType | null = null;
  let showCodeInput = false;
  let sourceCode = '';

  // --- TOOLTIPS AND LABELS ---
  const tooltips: Record<NodeType, string> = {
    source: 'Start here. Add source code to begin compilation.',
    lexer: 'Converts source code into tokens for processing.'
  };

  const nodeLabels: Record<NodeType, string> = {
    source: 'Source Code',
    lexer: 'Lexer'
  };

  // --- NODE CREATION ---
  function handleCreateNode(type: NodeType) {
    nodeCounter++;
    nodes.update(curr => [
      ...curr,
      {
        id: `${type}-${nodeCounter}`,
        type,
        label: nodeLabels[type] || type[0].toUpperCase() + type.slice(1),
        position: {
          x: Math.random() * 300 + 100,
          y: Math.random() * 300 + 100
        }
      }
    ]);
  }

  // --- PHASE SELECTION ---
  function handlePhaseSelect(e: CustomEvent<NodeType>) {
    const type = e.detail;
    if (type === 'source') {
      showCodeInput = true;
    } else {
      selectedPhase = type;
      // Only proceed if source code exists
      if (!sourceCode.trim()) {
        toasts.add({
          type: 'error',
          message: 'Please enter source code before proceeding'
        });
        selectedPhase = null;
        return;
      }
    }
  }

  function returnToCanvas() {
    selectedPhase = null;
    showCodeInput = false;
  }

  // --- CODE HANDLING ---
  function handleCodeSubmit(event: CustomEvent<string>) {
    sourceCode = event.detail;
    showCodeInput = false;
  }

  // --- TOKEN GENERATION ---
  let tokens = [];

  function handleTokenGeneration(event: CustomEvent<{ patterns: Array<{ type: string; pattern: string }> }>) {
    const patterns = event.detail.patterns;
    tokens = generateTokens(sourceCode, patterns);
  }

  function generateTokens(source: string, patterns: any[]) {
    // This is a simple implementation - you might want to make it more robust
    const tokens = [];
    let remainingText = source;

    while (remainingText.length > 0) {
      let matched = false;
      
      for (const { type, pattern } of patterns) {
        const regex = new RegExp(`^${pattern}`);
        const match = remainingText.match(regex);
        
        if (match) {
          tokens.push({ type, value: match[0] });
          remainingText = remainingText.slice(match[0].length);
          matched = true;
          break;
        }
      }

      if (!matched) {
        // Skip invalid character
        remainingText = remainingText.slice(1);
      }
    }

    return tokens;
  }
</script>

<NavBar />

<div class="main">
  <!-- Always show canvas and toolbox -->
  <Toolbox {handleCreateNode} {tooltips} />
  <div class="workspace">
    <DrawerCanvas {nodes} {edges} on:phaseSelect={handlePhaseSelect} />
  </div>

  <!-- Analysis overlay -->
  {#if selectedPhase}
    <div class="analysis-overlay">
      <div class="analysis-view">
        <div class="three-column-layout">
          <PhaseTutorial phase={selectedPhase} />
          <PhaseInspector 
            phase={selectedPhase}
            {sourceCode}
            on:generateTokens={handleTokenGeneration}
          />
          <ArtifactViewer 
            phase={selectedPhase}
            {tokens}
          />
        </div>
        <button on:click={returnToCanvas} class="return-button">
          ← Return to Canvas
        </button>
      </div>
    </div>
  {/if}

  <!-- Code input modal -->
  {#if showCodeInput}
    <div class="code-input-overlay">
      <div class="code-input-modal">
        <h2 class="modal-title">Enter Source Code</h2>
        <CodeInput on:codeSubmitted={handleCodeSubmit} />
        <button class="close-btn" on:click={() => (showCodeInput = false)}>✕</button>
      </div>
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
  .analysis-overlay {
    position: fixed;
    top: 3.5rem; /* Account for navbar */
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.95);
    z-index: 100;
  }
  .analysis-view {
    height: 100%;
    padding: 1rem;
    background: #f5f5f5;
    display: flex;
    flex-direction: column;
  }
  .three-column-layout {
    display: flex;
    flex: 1;
    gap: 1rem;
    height: calc(100vh - 6rem); /* Account for navbar and padding */
  }
  .three-column-layout > :global(*) {
    flex: 1;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    overflow-y: auto;
  }
  .return-button {
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 0.5rem 1rem;
    background: #001A6E;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    z-index: 1000;
  }
  .return-button:hover {
    background: #074799;
  }

  /* New modal styles */
  .code-input-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
  }

  .code-input-modal {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    max-width: 500px;
    width: 100%;
    position: relative;
  }

  .modal-title {
    margin: 0 0 1rem 0;
    font-size: 1.5rem;
    color: #333;
  }

  .close-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
  }
</style>

<script lang="ts">
  import NavBar from '$lib/components/main/NavBar.svelte';
  import Toolbox from '$lib/components/main/Toolbox.svelte';
  import CodeInput from '$lib/components/main/CodeInput.svelte';
  import DrawerCanvas from '$lib/components/main/DrawerCanvas.svelte';
  import type { NodeType, Token } from '$lib/types';
  import { writable } from 'svelte/store';
  import { addToast } from '$lib/stores/toast';
  import { theme } from '../../lib/stores/theme';
  import { onMount } from 'svelte';
  
  import LexerPhaseTutorial from '$lib/components/lexor/PhaseTutorial.svelte';
  import LexerPhaseInspector from '$lib/components/lexor/PhaseInspector.svelte';
  import LexerArtifactViewer from '$lib/components/lexor/ArtifactViewer.svelte';
  import ParserPhaseTutorial from '$lib/components/parser/PhaseTutorial.svelte';
  import ParserPhaseInspector from '$lib/components/parser/ParsingInput.svelte';
  import ParserArtifactViewer from '$lib/components/parser/ArtifactViewer.svelte';


  let workspaceEl: HTMLElement;

  // --- NEW: State for the one-time help tip ---
  let showDragTip = false;

  // Initialize theme
  onMount(() => {
    document.documentElement.classList.toggle('dark-mode', $theme === 'dark');

    // --- NEW: Check if the user has seen the tip before ---
    if (!localStorage.getItem('hasSeenDragTip')) {
      showDragTip = true;
    }
  });

  // --- CANVAS STATE ---
  interface CanvasNode {
    id: string;
    type: NodeType;
    label: string;
    position: { x: number; y: number };
  }
  
  const nodes = writable<CanvasNode[]>([]);
  let nodeCounter = 0;
  let selectedPhase: NodeType | null = null;
  let showCodeInput = false;
  let sourceCode = '';
  let showTokens = false;

  // --- TOOLTIPS AND LABELS ---
  const tooltips: Record<NodeType, string> = {
    source: 'Start here. Add source code to begin compilation.',
    lexer: 'Converts source code into tokens for processing.',
    parser: 'Analyzes the token stream to build a syntax tree.'
  };

  const nodeLabels: Record<NodeType, string> = {
    source: 'Source Code',
    lexer: 'Lexer',
    parser: 'Parser'
  };

  // --- NODE CREATION ---
  function handleCreateNode(type: NodeType) {
    nodeCounter++;
    nodes.update(curr => {
      // Define the layout parameters for a clean grid
      const start_x = 100; // Initial X position
      const start_y = 100; // Initial Y position
      const x_offset = 300; // Horizontal space between nodes
      const y_offset = 150; // Vertical space between rows
      const nodes_per_row = 3; // How many nodes before starting a new row

      // Calculate position based on the number of existing nodes
      const new_node_index = curr.length;
      const new_position = {
        x: start_x + (new_node_index % nodes_per_row) * x_offset,
        y: start_y + Math.floor(new_node_index / nodes_per_row) * y_offset
      };

      const new_node = {
        id: `${type}-${nodeCounter}`,
        type,
        label: nodeLabels[type] || type[0].toUpperCase() + type.slice(1),
        position: new_position
      };

      return [...curr, new_node];
    });

    // --- FIX: Programmatically focus the workspace after creating a node ---
    // This prevents the "jumping" bug by ensuring the canvas is the active element.
    workspaceEl?.focus();
  }
  
  // --- NEW: Function to dismiss the help tip permanently ---
  function dismissDragTip() {
    localStorage.setItem('hasSeenDragTip', 'true');
    showDragTip = false;
  }

  // --- PHASE SELECTION ---
  function handlePhaseSelect(e: CustomEvent<NodeType>) {
    const type = e.detail;
    if (type === 'source') {
      showCodeInput = true;
    } else {
      selectedPhase = type;
      if (!sourceCode.trim()) {
        addToast('Please enter source code before proceeding', "error");
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
    showTokens = false;
    sourceCode = event.detail;
    showCodeInput = false;
  }

  // --- TOKEN GENERATION ---
  let tokens: Token[] = [];
  let unexpectedTokens: string[] = [];

  function handleTokenGeneration(event: CustomEvent<{ tokens: Token[], unexpected_tokens: string[] }>) {
    showTokens = true;
    console.log('Received event:', event.detail); // Debug log
    tokens = event.detail.tokens;
    unexpectedTokens = event.detail.unexpected_tokens;
  }
</script>

<NavBar />

<div class="main">
  <Toolbox {handleCreateNode} {tooltips} />
  <!-- FIX: Bind the div to the variable and add tabindex="-1" to make it focusable -->
  <div class="workspace" bind:this={workspaceEl} tabindex="-1">
    <DrawerCanvas {nodes} on:phaseSelect={handlePhaseSelect} />

    <!-- UPDATED: Help tip with an 'X' close icon -->
    {#if showDragTip}
      <div class="help-tip">
        <span><b>Pro-Tip:</b> For the smoothest experience, click to select a node before dragging it.</span>
        <button on:click={dismissDragTip} class="dismiss-tip-btn" aria-label="Dismiss tip">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
        </button>
      </div>
    {/if}
  </div>

  <!-- This block handles rendering for DIFFERENT phases -->
  {#if selectedPhase}
    <div class="analysis-overlay">
      <div class="analysis-view">
        <div class="three-column-layout">

          {#if selectedPhase === 'lexer'}
            <LexerPhaseTutorial />
            <LexerPhaseInspector 
              {sourceCode}
              on:generateTokens={handleTokenGeneration}
            />
            <LexerArtifactViewer 
              phase={selectedPhase} 
              {tokens}
              {unexpectedTokens}
              {showTokens}
            />
          {/if}

          {#if selectedPhase === 'parser'}
            <ParserPhaseTutorial />
            <ParserPhaseInspector {sourceCode} />
            <ParserArtifactViewer />
          {/if}

        </div>
        <button on:click={returnToCanvas} class="return-button">
          ← Return to Canvas
        </button>
      </div>
    </div>
  {/if}

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
  :global(*) {
    box-sizing: border-box;
  }
  .main {
    display: flex;
    height: calc(100vh - 3.5rem);
    overflow: hidden;
    background-color: #f0f2f5;
    padding: 1rem;
    gap: 1rem;
    transition: background-color 0.3s ease;
  }
  .workspace {
    flex: 1;
    display: flex;
    flex-direction: column;
    position: relative;
    outline: none;
  }
  .analysis-overlay {
    position: fixed;
    top: 3.5rem; 
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
    height: calc(100vh - 6rem);
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
    margin-right: 1rem;
  }
  .return-button:hover {
    background: #074799;
  }
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
    text-align: center; 
    width: 100%; 
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

  /* UPDATED: Help tip styles */
  .help-tip {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    background-color: rgba(4, 26, 71, 0.95);
    color: white;
    padding: 10px 15px 10px 20px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
    display: flex;
    align-items: center;
    gap: 1rem;
    z-index: 50;
    font-size: 0.9rem;
  }

  .dismiss-tip-btn {
    background: none;
    border: none;
    color: white;
    opacity: 0.7;
    cursor: pointer;
    padding: 5px;
    line-height: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: opacity 0.2s ease;
  }
  .dismiss-tip-btn:hover {
    opacity: 1;
  }
  
  /* Dark Mode Styles */
  :global(html.dark-mode) .main {
    background-color: #161823;
  }
  :global(html.dark-mode) .analysis-overlay {
    background: rgba(10, 26, 58, 0.95);
  }
  :global(html.dark-mode) .analysis-view {
    background: #0a1a3a;
  }
  :global(html.dark-mode) .three-column-layout > :global(*) {
    background: #1a2a4a;
    color: #f0f0f0;
  }
  :global(html.dark-mode) .code-input-modal {
    background: #1a2a4a;
    color: #f0f0f0;
  }
  :global(html.dark-mode) .modal-title {
    color: #f0f0f0;
  }
  :global(html.dark-mode) .close-btn {
    color: #f0f0f0;
  }
  :global(html.dark-mode) .return-button {
    background: #1a3a7a;
  }
  :global(html.dark-mode) .return-button:hover {
    background: #2a4a8a;
  }

  :global(html.dark-mode) .return-button {
    background-color: #cccccc;
    margin-right: 1rem;
    color:#041a47;
  }

  :global(html.dark-mode) .return-button:hover {
    background-color: #5f636b;
    margin-right: 1rem;
  }
</style>

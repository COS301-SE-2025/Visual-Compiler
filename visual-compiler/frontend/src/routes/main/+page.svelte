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

  // --- Import ALL phase components with specific names ---
  import LexerPhaseTutorial from '$lib/components/lexor/PhaseTutorial.svelte';
  import LexerPhaseInspector from '$lib/components/lexor/PhaseInspector.svelte';
  import LexerArtifactViewer from '$lib/components/lexor/ArtifactViewer.svelte';
  import ParserPhaseTutorial from '$lib/components/parser/PhaseTutorial.svelte';
  import ParserPhaseInspector from '$lib/components/parser/ParsingInput.svelte';
  import ParserArtifactViewer from '$lib/components/parser/ArtifactViewer.svelte';

  // Initialize theme
  onMount(() => {
    document.documentElement.setAttribute('svelvet-theme', $theme);
    document.documentElement.classList.toggle('dark-mode', $theme === 'dark');
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
  <div class="workspace">
    <DrawerCanvas {nodes} on:phaseSelect={handlePhaseSelect} />
  </div>

  <!-- This block now handles rendering for DIFFERENT phases -->
  {#if selectedPhase}
    <div class="analysis-overlay">
      <div class="analysis-view">
        <div class="three-column-layout">

          <!-- Conditional Rendering for Lexer Phase -->
          {#if selectedPhase === 'lexer'}
            <LexerPhaseTutorial />
            <LexerPhaseInspector 
              {sourceCode}
              on:generateTokens={handleTokenGeneration}
            />
            <LexerArtifactViewer 
              {tokens}
              {unexpectedTokens}
              {showTokens}
            />
          {/if}

          <!-- Conditional Rendering for Parser Phase -->
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

  /* Dark mode styles */
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

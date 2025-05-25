<script lang="ts">
  import NavBar       from '$lib/components/NavBar.svelte';
  import Toolbox      from '$lib/components/Toolbox.svelte';
  import CodeInput    from '$lib/components/CodeInput.svelte';
  import DrawerCanvas from '$lib/components/DrawerCanvas.svelte';
  import type { NodeType } from '$lib/types';

  let sourceCode = '';

  function onCodeSubmitted(event: CustomEvent<string>) {
    sourceCode = event.detail;
  }

  // Accept a plain string (the event.detail), then cast it
  function forwardNodeEvent(type: string) {
    const nodeType = type as NodeType;
    window.dispatchEvent(
      new CustomEvent('createCanvasNode', { detail: { type: nodeType } })
    );
  }
</script>

<NavBar />

<div class="main-layout">
  <!-- now matches (e: string) => void exactly -->
  <Toolbox on:blockClick={forwardNodeEvent} />

  <div class="workspace">
    <CodeInput on:codeSubmitted={onCodeSubmitted} />
    <DrawerCanvas {sourceCode} />
  </div>
</div>

<style>
  :global(html, body) {
    margin: 0;
    padding: 0;
  }

  .main-layout {
    display: flex;
    height: calc(100vh - 3rem);
  }

  .workspace {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
</style>
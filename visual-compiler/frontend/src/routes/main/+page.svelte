<script lang="ts">
  import NavBar from '$lib/components/NavBar.svelte';
  import Toolbox from '$lib/components/Toolbox.svelte';
  import CodeInput from '$lib/components/CodeInput.svelte';
  import DrawerCanvas from '$lib/components/DrawerCanvas.svelte';
  import type { NodeType } from '$lib/components/Toolbox.svelte';

  let sourceCode = '';

  function onCodeSubmitted(event: CustomEvent<string>) {
    sourceCode = event.detail;
  }

  function handleBlockClick(event: CustomEvent<{ type: NodeType }>) {
    window.dispatchEvent(new CustomEvent('globalCreateNode', { 
      detail: { type: event.detail.type } 
    }));
  }
</script>

<NavBar />

<div class="main-layout">
  <Toolbox on:blockClick={handleBlockClick} />
  
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
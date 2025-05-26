<script lang="ts">
  import NavBar     from '$lib/components/NavBar.svelte';
  import Toolbox    from '$lib/components/Toolbox.svelte';
  import CodeInput  from '$lib/components/CodeInput.svelte';
  import Canvas     from '$lib/components/Canvas.svelte';

  import PhaseTutorial  from '$lib/components/PhaseTutorial.svelte';
  import PhaseInspector from '$lib/components/PhaseInspector.svelte';
  import ArtifactViewer from '$lib/components/ArtifactViewer.svelte';

  // CodeInput
   let sourceCode = '';
   let codeSubmitted = false;

  function onCodeSubmitted(event: CustomEvent<string>) {
    sourceCode = event.detail;
    codeSubmitted = true;
  }

  // Function to go back to the code input view
  function returnToEditor() {
    codeSubmitted = false;
  }
</script>

<NavBar />

<div class="main">
  {#if !codeSubmitted}
    <Toolbox />
    <div class="workspace">
      <CodeInput on:codeSubmitted={onCodeSubmitted} />
      <Canvas />
    </div>
  {:else}
    <div class="analysis-view">
      <div class="three-column-layout">
        <PhaseTutorial />
        <PhaseInspector sourceCode={sourceCode} />
        <ArtifactViewer />
      </div>
      <button on:click={returnToEditor} class="return-button">Return to Editor</button>
    </div>
  {/if}
</div>

<style>
  :global(html, body) {
    margin: 0;
    padding: 0;
  }

  .main {
    display: flex;
    height: calc(100vh - 3rem);
  }

  .workspace {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .three-column-layout {
    display: flex;
    flex: 1;
  }

  .analysis-view {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .return-button {
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 10px 20px;
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
</style>
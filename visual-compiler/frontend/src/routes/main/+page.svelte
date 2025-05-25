<script lang="ts">
  import NavBar     from '$lib/components/NavBar.svelte';
  import Toolbox    from '$lib/components/Toolbox.svelte';
  import CodeInput  from '$lib/components/CodeInput.svelte';
  import Canvas     from '$lib/components/Canvas.svelte';

  // import PhaseTutorial  from '$lib/components/PhaseTutorial.svelte';
  // import PhaseInspector from '$lib/components/PhaseInspector.svelte';
  // import ArtifactViewer from '$lib/components/ArtifactViewer.svelte';
</script>

<NavBar />

<div class="main">
  <Toolbox />

  <div class="workspace">
    <CodeInput />
    <Canvas />
  </div>
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
</style>

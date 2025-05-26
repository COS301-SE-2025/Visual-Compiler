<script lang="ts">
  import NavBar     from '$lib/components/NavBar.svelte';
  import Toolbox    from '$lib/components/Toolbox.svelte';
  import CodeInput  from '$lib/components/CodeInput.svelte';
  import DrawerCanvas from '$lib/components/DrawerCanvas.svelte';
  import type { NodeType } from '$lib/types';
  import { writable } from 'svelte/store';

  // --- SOURCE CODE STATE ---
  let sourceCode = '';

  function onCodeSubmitted(e: CustomEvent<string>) {
    sourceCode = e.detail;
  }

  // --- CANVAS NODES STATE ---
  interface CanvasNode {
    id:     string;
    type:   NodeType;
    label:  string;
    position: { x: number; y: number };
  }

  const nodes = writable<CanvasNode[]>([]);
  let nodeCounter = 0;

  // called directly by Toolbox
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
</script>

<NavBar />

<div class="main-layout">
  <!-- PASS the handler directly; Toolbox just calls it -->
  <Toolbox {handleCreateNode} />

  <div class="workspace">
    <CodeInput     on:codeSubmitted={onCodeSubmitted} />
    <!-- GIVE Canvas the nodes store & sourceCode -->
    <DrawerCanvas {sourceCode} {nodes} />
  </div>
</div>

<style>
  :global(html, body) { margin:0; padding:0; }
  .main-layout { display:flex; height:calc(100vh - 3rem); }
  .workspace   { flex:1; display:flex; flex-direction:column; }
</style>
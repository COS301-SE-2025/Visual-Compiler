<script lang="ts">

</script>

<div class="artifact-viewer">
  <h3>Artifacts: Syntax Tree</h3>
  <div class="empty-state">
    The generated Abstract Syntax Tree (AST) will be visualized here.
  </div>
</div>

<style>

  .artifact-viewer {
    padding: 1.5rem;
  }
  h3 {
    color: #001A6E;
    font-size: 1.25rem;
    margin: 0 0 1.5rem 0;
    padding-bottom: 1rem;
    border-bottom: 2px solid #e5e7eb;
  }
  .empty-state {
    color: #666;
    font-style: italic;
    text-align: center;
    padding: 2rem;
  }

  /* Dark Mode Styles */
  :global(html.dark-mode) h3 { color: #ebeef1; border-bottom-color: #4b5563; }
  :global(html.dark-mode) .empty-state { color: #9ca3af; }
</style>

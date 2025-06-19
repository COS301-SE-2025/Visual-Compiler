<script lang="ts">
  import type { Token } from '$lib/types';
  
  export let phase: string;
  export let tokens: Token[] = [];
  export let unexpectedTokens: string[] = [];
  export let showTokens = false;

  $: if (phase) showTokens = false;

  $: console.log('ArtifactViewer phase:', phase);
  $: console.log('ArtifactViewer tokens:', tokens);
  $: console.log('ArtifactViewer unexpected tokens:', unexpectedTokens);
</script>

<div class="artifact-viewer">
  {#if phase === 'lexer' && showTokens}
    <div class="tokens-heading">
      <h3>Tokens</h3>
    </div>
    {#if tokens && tokens.length > 0}
      <table class="token-table">
        <thead>
          <tr>
            <th>Type</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {#each tokens as token}
            <tr>
              <td>{token.Type}</td>
              <td>{token.Value}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    {:else}
      <div class="no-tokens">No tokens generated yet</div>
    {/if}

    {#if unexpectedTokens && unexpectedTokens.length > 0}
      <div class="unexpected-tokens-container">
        <h4>Unidentified Input</h4>
        <div class="token-list">
          {#each unexpectedTokens as token}
            <span class="unexpected-token">{token}</span>
          {/each}
        </div>
      </div>
    {/if}
  {:else if phase === 'lexer'}
    <div class="empty-state">
      Tokens will appear here after generation
    </div>
  {/if}
</div>

<style>
  .artifact-viewer {
    padding: 2rem 1.5rem 1.5rem;
  }

  h3 {
    color: #001A6E;
    font-size: 1.25rem;
    margin: 0 0 1.5rem 0;
    padding-bottom: 1rem;
    border-bottom: 2px solid #e5e7eb;
  }

  .token-table {
    width: 100%;
    table-layout: fixed;
    border-collapse: collapse;
    margin: 1.5rem 0 2rem 0;
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .token-table th,
  .token-table td {
    width: 50%;
    padding: 0.75rem 1rem;
    border: 1px solid #e0e0e0;
    text-align: left;
  }

  .token-table th {
    background: #001A6E;
    color: white;
    font-weight: 500;
  }

  .unexpected-tokens-container {
    margin-top: 1.5rem;
    padding: 1.2rem;
    background: #f8f9fa;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .unexpected-tokens-container h4 {
    color: #001A6E;
    font-size: 1.1rem;
    margin: 0 0 1.5rem 0;
    padding-bottom: 0.6rem;
    border-bottom: 1px solid #e5e7eb;
  }

  .token-list {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }

  .tokens-heading {
    color: #001A6E;
  }
  
  .unexpected-token {
    padding: 0.5rem 1rem;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    color: #001A6E;
    font-family: monospace;
    font-size: 0.9rem;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .empty-state {
    color: #666;
    font-style: italic;
    text-align: center;
    padding: 2rem;
  }

  .no-tokens {
    color: #666;
    font-style: italic;
    text-align: center;
    margin: 1rem 0;
  }

  /* Dark mode styles */
  :global(html.dark-mode) .artifact-viewer {
    background: #1a2a4a;
  }

  :global(html.dark-mode) h3,
  :global(html.dark-mode) .tokens-heading {
    color: #ebeef1;
    border-bottom-color: #4b5563;
  }

  :global(html.dark-mode) .token-table {
    background: #2d3748;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  }

  :global(html.dark-mode) .token-table th,
  :global(html.dark-mode) .token-table td {
    border-color: #4b5563;
    color: #f0f0f0;
  }

  :global(html.dark-mode) .token-table th {
    background: #1e3a8a;
  }

  :global(html.dark-mode) .unexpected-tokens-container {
    background: #1f2937;
    border-color: #4b5563;
  }

  :global(html.dark-mode) .unexpected-tokens-container h4 {
    color: #ebeef1;
    border-bottom-color: #4b5563;
  }

  :global(html.dark-mode) .unexpected-token {
    background: #374151;
    border-color: #4b5563;
    color: #e5e7eb;
  }

  :global(html.dark-mode) .empty-state,
  :global(html.dark-mode) .no-tokens {
    color: #9ca3af;
  }
</style>
<script lang="ts">
  export let phase: string;
  export let tokens: Array<{ type: string, value: string }> = [];
</script>

<div class="artifact-viewer">
  {#if phase === 'lexer' && tokens.length > 0}
    <h3>Tokens</h3>
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
            <td>{token.type}</td>
            <td>{token.value}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <div class="placeholder">
      <p>Token table will appear here after generation</p>
    </div>
  {/if}
</div>

<style>
  .artifact-viewer {
    flex: 1;
    padding: 2rem;
    background: #fff;
    overflow: auto;
    
  }

  .artifact-viewer h3 {
    margin-top: 0;
    color: #001A6E;
    font-size: 1.25rem;
    font-weight: 600;
    border-bottom: 2px solid #001A6E;
    padding-bottom: 0.5rem;
  }

  .token-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
    font-family: sans-serif;
    
  }

  .token-table th,
  .token-table td {
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    text-align: left;
    font-size: 0.95rem;
    
  }

  .token-table th {
    background: #001A6E;
    color: #fff;
    font-weight: 500;
   
  }

  .token-table tr:nth-child(even) {
    background: #f9f9f9;
  }

  .placeholder {
    display: flex;
    height: 100%;
    align-items: center;
    justify-content: center;
    color: #666;
    font-style: italic;
  }
</style>

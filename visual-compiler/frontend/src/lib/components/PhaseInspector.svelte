<script lang="ts">
  export let sourceCode = 'Source code will be displayed here';
  
  // Array to store multiple type-regex pairs
  let inputRows = [{ type: '', regex: '' }];

  function addNewRow() {
    inputRows = [...inputRows, { type: '', regex: '' }];
  }
</script>

<div class="phase-inspector">
  <div class="source-code-section">
    <h3>Source Code</h3>
    <pre class="source-display">{sourceCode}</pre>
  </div>
  
  <div class="shared-block">
    <div class="block-headers">
      <div class="header-section">
        <h3>Type</h3>
      </div>
      <div class="header-section">
        <h3>Regular Expression</h3>
      </div>
    </div>
    <div class="input-rows">
      {#each inputRows as row, i}
        <div class="input-row">
          <div class="input-block">
            <input 
              type="text" 
              bind:value={row.type} 
              placeholder="Enter type..."
            />
          </div>
          <div class="input-block">
            <input 
              type="text" 
              bind:value={row.regex} 
              placeholder="Enter regex pattern..."
            />
          </div>
        </div>
      {/each}
    </div>

    {#if inputRows[inputRows.length - 1].type && inputRows[inputRows.length - 1].regex}
      <button class="add-button" on:click={addNewRow}>
        <span>+</span>
      </button>
    {/if}
  </div>
  
  <div class="button-container">
    <button class="submit-button">
      Submit
    </button>
  </div>
</div>

<style>
  .phase-inspector {
    flex: 1.2;
    padding: 2rem;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    background: #fff;
  }

  .source-code-section {
    margin-bottom: 2rem;
  }

  .source-display {
    background: #f5f5f5;
    padding: 1rem;
    border-radius: 4px;
    overflow-x: auto;
    font-family: monospace;
    white-space: pre-wrap;
    margin: 0;
  }

  .shared-block {
    background: #f5f5f5;
    padding: 1.2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
    position: relative;
  }

  .block-headers {
    display: flex;
    gap: 2rem;
    margin-bottom: 1rem;
    padding-bottom: 0.8rem;
    border-bottom: 1px solid #e0e0e0;
  }

  .header-section {
    flex: 1;
  }

  .header-section h3 {
    margin: 0;
    color: #001A6E;
    font-size: 1rem;
    font-weight: 600;
  }

  .input-rows {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding-top: 0.5rem;
  }

  .input-row {
    display: flex;
    gap: 2rem;
  }

  .input-block {
    flex: 1;
  }

  .input-block input {
    width: 100%;
    padding: 0.8rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 0.9rem;
    background: white;
    box-sizing: border-box;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .input-block input:focus {
    outline: none;
    border-color: #001A6E;
    box-shadow: 0 0 0 2px rgba(32,87,129,0.1);
  }

  .add-button {
    position: absolute;
    right: -16px;
    bottom: -16px;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #001A6E;
    color: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.2rem;
    transition: background 0.2s;
  }

  .add-button:hover {
    background: #27548A;
  }

  .button-container {
    display: flex;
    justify-content: center;
    margin-top: 1rem;
  }

  .submit-button {
    padding: 0.6rem 1.75rem;
    background: #001A6E;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
  }

  .submit-button:hover {
    background: #27548A;
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 26, 110, 0.2);
  }

  .submit-button:active {
    transform: translateY(0);
    box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
  }
</style>

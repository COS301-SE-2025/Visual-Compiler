<script lang="ts">
  export let sourceCode = 'Source code will be displayed here';
  
  let inputRows = [{ type: '', regex: '', error: '' }];
  let formError = '';
  let submissionStatus = { show: false, success: false, message: '' };

  function addNewRow() {
    inputRows = [...inputRows, { type: '', regex: '', error: '' }];
  }

  function validateRegex(pattern: string): boolean {
    try {
      new RegExp(pattern);
      return true;
    } catch (e) {
      return false;
    }
  }

  function handleSubmit() {
    formError = '';
    let hasErrors = false;

    inputRows = inputRows.map(row => {
      // Reset previous error
      row.error = '';

      // Check if both fields are empty
      if (!row.type && !row.regex) {
        return null; // Will be filtered out
      }

      // Check if one field is empty
      if (!row.type || !row.regex) {
        row.error = 'Please fill in both Type and Regular Expression';
        hasErrors = true;
        return row;
      }

      // Validate regex
      if (!validateRegex(row.regex)) {
        row.error = 'Invalid regular expression pattern';
        hasErrors = true;
        return row;
      }

      return row;
    }).filter(row => row !== null);

    if (hasErrors) {
      submissionStatus = {
        show: true,
        success: false,
        message: 'Please fix the errors before submitting'
      };
    } else if (inputRows.length === 0) {
      submissionStatus = {
        show: true,
        success: false,
        message: 'At least one row with both fields is required'
      };
      inputRows = [{ type: '', regex: '', error: '' }];
    } else {
      // Process valid submission
      submissionStatus = {
        show: true,
        success: true,
        message: 'Successfully submitted!'
      };
      console.log('Valid submission:', inputRows);
    }

    // Hide the status message after 3 seconds
    setTimeout(() => {
      submissionStatus = { show: false, success: false, message: '' };
    }, 3000);
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
              class:error={row.error}
            />
          </div>
          <div class="input-block">
            <input 
              type="text" 
              bind:value={row.regex} 
              placeholder="Enter regex pattern..."
              class:error={row.error}
            />
          </div>
          {#if row.error}
            <div class="error-message">{row.error}</div>
          {/if}
        </div>
      {/each}
    </div>

    {#if inputRows[inputRows.length - 1].type && inputRows[inputRows.length - 1].regex}
      <button class="add-button" on:click={addNewRow}>
        <span>+</span>
      </button>
    {/if}
  </div>
  
  {#if formError}
    <div class="form-error">{formError}</div>
  {/if}

  {#if submissionStatus.show}
    <div class="status-message" class:success={submissionStatus.success}>
      {submissionStatus.message}
    </div>
  {/if}
  
  <div class="button-container">
    <button class="submit-button" on:click={handleSubmit}>
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
    position: relative;
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

  .input-block input.error {
    border-color: #dc3545;
  }

  .input-block input.error:focus {
    border-color: #dc3545;
    box-shadow: 0 0 0 2px rgba(220, 53, 69, 0.25);
  }

  .error-message {
    color: #dc3545;
    font-size: 0.8rem;
    margin-top: 0.25rem;
    position: absolute;
    bottom: -1.2rem;
  }

  .form-error {
    color: #dc3545;
    text-align: center;
    margin: 0.5rem 0;
    font-size: 0.9rem;
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

  .status-message {
    text-align: center;
    padding: 0.5rem 1rem;
    margin: 0.5rem 0;
    border-radius: 4px;
    font-size: 0.9rem;
    background: #dc3545;
    color: white;
    opacity: 0;
    animation: fadeInOut 3s ease-in-out;
  }

  .status-message.success {
    background: #28a745;
  }

  @keyframes fadeInOut {
    0% { opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { opacity: 0; }
  }
</style>

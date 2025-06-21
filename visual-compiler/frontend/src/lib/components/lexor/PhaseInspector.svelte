<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Token } from '$lib/types';
  import type { NodeType } from '$lib/types';
	import { addToast } from '$lib/stores/toast';
  
  export let sourceCode = '';  
  
  let inputRows = [{ type: '', regex: '', error: '' }];
  let formError = '';
  let submissionStatus = { show: false, success: false, message: '' };
  let showGenerateButton = false;
  const dispatch = createEventDispatcher<{ 
    generateTokens: { 
      tokens: Token[] 
    } 
  }>();

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

  async function handleSubmit() {
    formError = '';
    let hasErrors = false;
    let nonEmptyRows = [];

    // Handle validation and filtering in one pass
    for (const row of inputRows) {
      if (!row.type && !row.regex) continue;

      row.error = '';

      if (!row.type || !row.regex) {
        row.error = 'Please fill in both Type and Regular Expression';
        hasErrors = true;
      } else if (!validateRegex(row.regex)) {
        row.error = 'Invalid regular expression pattern';
        hasErrors = true;
      }

      nonEmptyRows.push(row);
    }

    // First check if we have a single empty row
    if (inputRows.length === 1 && !inputRows[0].type && !inputRows[0].regex) {
      submissionStatus = {
        show: true,
        success: false,
        message: 'Please fill in both Type and Regular Expression'
      };
      return;
    }

    // Update inputRows with filtered rows
    inputRows = nonEmptyRows;

    if (hasErrors) {
      submissionStatus = {
        show: true,
        success: false,
        message: 'Please fix the errors before submitting'
      };
      return;
    }

    // Prepare data for API call
    const requestData = {
      source_code: sourceCode,
      pairs: inputRows.map(row => ({
        Type: row.type.toUpperCase(),
        Regex: row.regex
      }))
    };

    try {
      // First, store the source code and regex pairs
      const storeResponse = await fetch('http://localhost:8080/api/lexing/code', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData)
      });

      if (!storeResponse.ok) {
        const errorText = await storeResponse.text();
        throw new Error(`Server error (${storeResponse.status}): ${errorText}`);
      }

      // Show success message
      submissionStatus = {
        show: true,
        success: true,
        message: 'Code stored successfully!'
      };

      // Show generate button
      showGenerateButton = true;

    } catch (error) {
      console.error('Store error:', error);
      addToast("Cannot connect to server. Please ensure the backend is running.", "error");
    }
  }

  async function generateTokens() {
    try {
      const requestData = {
        source_code: sourceCode,
        pairs: inputRows.map(row => ({
          Type: row.type.toUpperCase(),
          Regex: row.regex
        }))
      };

      console.log('Sending request with data:', requestData); // Debug log

      const response = await fetch('http://localhost:8080/api/lexing/lexer', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData)
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Server error (${response.status}): ${errorText}`);
      }

      const data = await response.json();
      console.log('Raw response:', data); // Debug entire response
      
      if (!Array.isArray(data.tokens)) {
        throw new Error('Expected tokens array in response');
      }

      // Dispatch token data to parent component
      dispatch('generateTokens', {
        tokens: data.tokens,
        unexpected_tokens: data.tokens_unidentified// Make sure this matches the backend response
      });

      showGenerateButton = false;
      submissionStatus = {
        show: true,
        success: true,
        message: data.message || 'Tokens generated successfully!'
      };

    } catch (error) {
      console.error('Generate tokens error:', error);
      addToast("Error generating tokens","error");
    }
  }

  // Track previous input values
  let previousInputs: typeof inputRows = [];

  // Reset generate button when inputs change
  function handleInputChange() {
    showGenerateButton = false;
    // Reset submission status
    submissionStatus = { 
      show: false, 
      success: false, 
      message: '' 
    };
  }

  // Watch for changes in inputRows array
  $: {
    // Check if array length changed or values changed
    const inputsChanged = inputRows.length !== previousInputs.length ||
      inputRows.some((row, index) => {
        const prevRow = previousInputs[index];
        return !prevRow || row.type !== prevRow.type || row.regex !== prevRow.regex;
      });

    if (inputsChanged) {
      handleInputChange();
      previousInputs = [...inputRows];
    }
  }

  let selectedAutomaton: 'NFA' | 'DFA' | 'REGEX' | null = null;
  let showDefault = false;

  // NFA/DFA input state
  let nfaStates = '';
  let nfaAlphabets = '';
  let nfaStartState = '';
  let nfaAcceptedStates = '';
  let nfaTransitions = '';

  let dfaStates = '';
  let dfaAlphabets = '';
  let dfaStartState = '';
  let dfaAcceptedStates = '';
  let dfaTransitions = '';

  function resetInputs() {
    nfaStates = '';
    nfaAlphabets = '';
    nfaStartState = '';
    nfaAcceptedStates = '';
    nfaTransitions = '';
    dfaStates = '';
    dfaAlphabets = '';
    dfaStartState = '';
    dfaAcceptedStates = '';
    dfaTransitions = '';
    showDefault = false;
  }

  function selectAutomaton(type: 'NFA' | 'DFA' | 'REGEX') {
    selectedAutomaton = type;
    resetInputs();
  }

  function insertDefault() {
    showDefault = true;
    if (selectedAutomaton === 'NFA') {
      nfaStates = 'q0,q1,q2';
      nfaAlphabets = 'a,b';
      nfaStartState = 'q0';
      nfaAcceptedStates = 'q2';
      nfaTransitions = 'q0,a->q1\nq1,b->q2'; // <-- Make sure this is set
    } else if (selectedAutomaton === 'DFA') {
      dfaStates = 'A,B';
      dfaAlphabets = '0,1';
      dfaStartState = 'A';
      dfaAcceptedStates = 'B';
      dfaTransitions = 'A,0->A\nA,1->B\nB,0->A\nB,1->B';
    }
    // No default for REGEX yet
  }

  function removeDefault() {
    showDefault = false;
    nfaStates = '';
    nfaAlphabets = '';
    nfaStartState = '';
    nfaAcceptedStates = '';
    nfaTransitions = '';
    dfaStates = '';
    dfaAlphabets = '';
    dfaStartState = '';
    dfaAcceptedStates = '';
    dfaTransitions = '';
  }
</script>

<div class="phase-inspector">
  <div class="source-code-section">
    <div class="lexor-heading">
       <h1 class="lexor-heading-h1">LEXING</h1>
    </div>
    <h3 class="source-code-header">Source Code</h3>
    <pre class="source-display">{sourceCode || 'No source code available'}</pre>
  </div>
  
  <!-- Automaton selection buttons -->
<div class="automaton-btn-row">
  <button
    class="automaton-btn {selectedAutomaton === 'NFA' ? 'selected' : ''}"
    on:click={() => selectAutomaton('NFA')}
    type="button"
  >NFA</button>
  <button
    class="automaton-btn {selectedAutomaton === 'DFA' ? 'selected' : ''}"
    on:click={() => selectAutomaton('DFA')}
    type="button"
  >DFA</button>
  <button
    class="automaton-btn {selectedAutomaton === 'REGEX' ? 'selected' : ''}"
    on:click={() => selectAutomaton('REGEX')}
    type="button"
  >Regular Expression</button>
  {#if selectedAutomaton && !showDefault}
    <button
      class="default-toggle-btn"
      on:click={insertDefault}
      type="button"
      aria-label="Insert default input"
      title="Insert default input"
    >
      <span class="icon">🪄</span>
    </button>
  {/if}
  {#if selectedAutomaton && showDefault}
    <button
      class="default-toggle-btn selected"
      on:click={removeDefault}
      type="button"
      aria-label="Remove default input"
      title="Remove default input"
    >
      <span class="icon">🧹</span>
    </button>
  {/if}
</div>

  <!-- Conditional content rendering -->
  {#if selectedAutomaton === 'REGEX'}
    <!-- Place your existing lexing/regex UI here -->
    <div>
      <!-- BEGIN LEXING UI -->
      <!-- ...everything from your shared-block and below, up to the closing </style>... -->
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
              on:input={handleInputChange}
              placeholder="Enter type..."
              class:error={row.error}
            />
          </div>
          <div class="input-block">
            <input 
              type="text" 
              bind:value={row.regex} 
              on:input={handleInputChange}
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

  <div class="button-container">
    <button 
      class="submit-button" 
      class:shifted={showGenerateButton} 
      on:click={handleSubmit}
    >
      Submit
    </button>
    
    {#if showGenerateButton}
      <button class="generate-button" on:click={generateTokens}>
        Generate Tokens
      </button>
    {/if}
  </div>

  {#if submissionStatus.show}
    <div 
      class="status-message" 
      class:success={submissionStatus.success === true}
      class:info={submissionStatus.message === 'info'}
    >
      {submissionStatus.message}
    </div>
  {/if}
      <!-- END LEXING UI -->
    </div>
  {:else if selectedAutomaton === 'NFA'}
    <div class="automaton-section">
      <label>
        States:
        <input class="automaton-input" bind:value={nfaStates} placeholder="e.g. q0,q1,q2" />
      </label>
      <label>
        Alphabets:
        <input class="automaton-input" bind:value={nfaAlphabets} placeholder="e.g. a,b" />
      </label>
      <label>
        Start State:
        <input class="automaton-input" bind:value={nfaStartState} placeholder="e.g. q0" />
      </label>
      <label>
        Accepted States:
        <input class="automaton-input" bind:value={nfaAcceptedStates} placeholder="e.g. q2" />
      </label>
      <label class="automaton-transitions-label">
        Transitions:
        <textarea
          class="automaton-input automaton-transitions"
          bind:value={nfaTransitions}
          placeholder="e.g. q0,a->q1&#10;q1,b->q2"
        ></textarea>
      </label>
    </div>
  {:else if selectedAutomaton === 'DFA'}
    <div class="automaton-section">
      <label>
        States:
        <input class="automaton-input" bind:value={dfaStates} placeholder="e.g. A,B" />
      </label>
      <label>
        Alphabets:
        <input class="automaton-input" bind:value={dfaAlphabets} placeholder="e.g. 0,1" />
      </label>
      <label>
        Start State:
        <input class="automaton-input" bind:value={dfaStartState} placeholder="e.g. A" />
      </label>
      <label>
        Accepted States:
        <input class="automaton-input" bind:value={dfaAcceptedStates} placeholder="e.g. B" />
      </label>
      <label class="automaton-transitions-label">
        Transitions:
        <textarea class="automaton-input automaton-transitions" bind:value={dfaTransitions} placeholder="e.g. A,0->A&#10;A,1->B"></textarea>
      </label>
    </div>
  {/if}
</div>

<style>

  .source-code-header{
    color: #444;
  }
  .phase-inspector {
    flex: 1.2;
    padding-left: 2rem;
    padding-right: 2rem;
    padding-top: 0.2rem;
    padding-bottom: 2rem;
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

  .lexor-heading{
    justify-items:center;
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
    gap: 0.5rem;
    margin-top: 0.75rem;
  }

  .submit-button {
    padding: 0.6rem 1.5rem;
    background: #001A6E;
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: transform 0.2s ease; 
    box-shadow: 0 2px 4px rgba(0, 26, 110, 0.1);
    position: relative; 
    margin-top: 1rem;
  }

  .submit-button.shifted {
    transform: translateX(-1rem);
  }


  .submit-button:not(.shifted) {
    transform: translateX(0);
  }

  .generate-button {
  padding: 0.6rem 1.5rem; 
  background: #666;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  box-shadow: 0 2px 4px rgba(40, 167, 69, 0.1);
  margin-top: 1rem; 
  height: 100%; 
  line-height: 1.1; 
}

  .submit-button:hover {
    background: #27548A;
    
  } 

   :global(html.dark-mode) .submit-button {
    background: #cccccc;
    color:#041a47;
    transition: transform 0.2s ease; 
    margin-top: 1rem;
  }

  .generate-button:hover {
    background: #3d3d3d;
  }

  .status-message {
    text-align: center;
    padding: 0.5rem 1rem;
    margin-top: 0.75rem;
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

  .status-message.info {
    background: #0096c7; 
  }

  @keyframes fadeInOut {
    0% { opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { opacity: 0; }
  }
  :global(html.dark-mode) .source-code-header {
    color: #ebeef1;
  }
 
  
  :global(html.dark-mode) .lexor-heading-h1 {
    color: #ebeef1;
  }

  :global(html.dark-mode)  .source-display{
    color : black;
  }

  .automaton-btn-row {
    display: flex;
    gap: 0.7rem;
    margin: 2rem 0 1.5rem 0;
    align-items: center;
  }
  .automaton-btn {
    padding: 0.4rem 1rem;
    border: 2px solid #e5e7eb;
    border-radius: 1.2rem;
    background: white;
    color: #001A6E;
    font-weight: 500;
    font-size: 0.95rem;
    cursor: pointer;
    transition: border-color 0.2s, background 0.2s, color 0.2s;
    outline: none;
  }
  .automaton-btn.selected {
    border-color: #001A6E;
    background: #e6edfa;
    color: #001A6E;
  }
  .automaton-btn:not(.selected):hover {
    border-color: #7da2e3;
    background: #f5f8fd;
  }
  .default-toggle-btn {
  margin-left: 1.2rem;
  padding: 0.4rem 0.7rem;
  border-radius: 1.2rem;
  border: 2px solid #e5e7eb;      /* match unselected automaton-btn */
  background: white;              /* match unselected automaton-btn */
  color: #001A6E;
  font-size: 1.2rem;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
  display: flex;
  align-items: center;
  height: 2.2rem;
  width: 2.2rem;
  justify-content: center;
  position: relative;
}
  .default-toggle-btn.selected {
    background: #d0e2ff;
    border-color: #003399;
  }
  .default-toggle-btn:hover,
  .default-toggle-btn:focus {
    background: #f5f8fd;
    border-color: #7da2e3;
  }
  .icon {
    font-size: 1.3rem;
    line-height: 1;
    pointer-events: none;
  }

  .automaton-section {
    margin-top: 2rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.2rem 1.5rem;
    max-width: 600px;
  }
  .automaton-section label {
    font-weight: 500;
    color: #001A6E;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    font-size: 1rem;
  }
  .automaton-input {
    padding: 0.7rem 1rem;
    border: 1.5px solid #b6c6e3;
    border-radius: 0.7rem;
    font-size: 1rem;
    background: #f8fafc;
    color: #001A6E;
    transition: border-color 0.2s;
    width: 100%;
    box-sizing: border-box;
  }
  .automaton-input:focus {
    border-color: #001A6E;
    outline: none;
    background: #e6edfa;
  }
  .automaton-transitions-label {
    grid-column: span 2;
  }
  .automaton-transitions {
    min-height: 6rem;
    min-width: 100%;
    font-family: 'Fira Mono', monospace;
    resize: vertical;
  }
</style>

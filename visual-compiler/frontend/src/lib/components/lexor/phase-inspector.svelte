<script lang="ts">
  import type { Token } from '$lib/types';
  import type { NodeType } from '$lib/types';
  import { AddToast } from '$lib/stores/toast';

  export let source_code = '';

  // This prop will be called instead of dispatching an event
  export let onGenerateTokens: (data: {
    tokens: Token[];
    unexpected_tokens: string[];
  }) => void = () => {};

  let input_rows = [{ type: '', regex: '', error: '' }];
  let form_error = '';
  let submission_status = { show: false, success: false, message: '' };
  let show_generate_button = false;

  // addNewRow
  // Return type: void
  // Parameter type(s): none
  // Adds a new, empty row to the regex input form.
  function addNewRow() {
    input_rows = [...input_rows, { type: '', regex: '', error: '' }];
  }

  // validateRegex
  // Return type: boolean
  // Parameter type(s): string
  // Checks if a given string is a valid regular expression pattern.
  function validateRegex(pattern: string): boolean {
    try {
      new RegExp(pattern);
      return true;
    } catch (e) {
      return false;
    }
  }

  // handleSubmit
  // Return type: Promise<void>
  // Parameter type(s): none
  // Validates the form and sends the source code and regex pairs to the backend.
  async function handleSubmit() {
    form_error = '';
    let has_errors = false;
    let non_empty_rows = [];

    // Handle validation and filtering in one pass
    for (const row of input_rows) {
      if (!row.type && !row.regex) continue;

      row.error = '';

      if (!row.type || !row.regex) {
        row.error = 'Please fill in both Type and Regular Expression';
        has_errors = true;
      } else if (!validateRegex(row.regex)) {
        row.error = 'Invalid regular expression pattern';
        has_errors = true;
      }

      non_empty_rows.push(row);
    }

    if (input_rows.length === 1 && !input_rows[0].type && !input_rows[0].regex) {
      submission_status = {
        show: true,
        success: false,
        message: 'Please fill in both Type and Regular Expression'
      };
      return;
    }

    input_rows = non_empty_rows;

    if (has_errors) {
      submission_status = {
        show: true,
        success: false,
        message: 'Please fix the errors before submitting'
      };
      return;
    }

    const request_data = {
      source_code: source_code,
      pairs: input_rows.map((row) => ({
        Type: row.type.toUpperCase(),
        Regex: row.regex
      }))
    };

    try {
      const store_response = await fetch('http://localhost:8080/api/lexing/code', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(request_data)
      });

      if (!store_response.ok) {
        const error_text = await store_response.text();
        throw new Error(`Server error (${store_response.status}): ${error_text}`);
      }

      submission_status = {
        show: true,
        success: true,
        message: 'Code stored successfully!'
      };

      show_generate_button = true;
    } catch (error) {
      console.error('Store error:', error);
      AddToast('Cannot connect to server. Please ensure the backend is running.', 'error');
    }
  }

  // generateTokens
  // Return type: Promise<void>
  // Parameter type(s): none
  // Requests token generation from the backend and dispatches the result.
  async function generateTokens() {
    try {
      const request_data = {
        source_code: source_code,
        pairs: input_rows.map((row) => ({
          Type: row.type.toUpperCase(),
          Regex: row.regex
        }))
      };

      console.log('Sending request with data:', request_data); // Debug log

      const response = await fetch('http://localhost:8080/api/lexing/lexer', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(request_data)
      });

      if (!response.ok) {
        const error_text = await response.text();
        throw new Error(`Server error (${response.status}): ${error_text}`);
      }

      const data = await response.json();
      console.log('Raw response:', data); // Debug entire response

      if (!Array.isArray(data.tokens)) {
        throw new Error('Expected tokens array in response');
      }

      // Call the prop function instead of dispatching
      onGenerateTokens({
        tokens: data.tokens,
        unexpected_tokens: data.unexpected_tokens
      });

      show_generate_button = false;
      submission_status = {
        show: true,
        success: true,
        message: data.message || 'Tokens generated successfully!'
      };
    } catch (error) {
      console.error('Generate tokens error:', error);
      AddToast('Error generating tokens', 'error');
    }
  }

  // handleInputChange
  // Return type: void
  // Parameter type(s): none
  // Resets the generate button and submission status when form inputs change.
  function handleInputChange() {
    show_generate_button = false;
    submission_status = {
      show: false,
      success: false,
      message: ''
    };
  }

  // resetInputs
  // Return type: void
  // Parameter type(s): none
  // Clears all NFA/DFA input fields and hides the default view.
  function resetInputs() {
    nfa_states = '';
    nfa_alphabets = '';
    nfa_start_state = '';
    nfa_accepted_states = '';
    nfa_transitions = '';
    dfa_states = '';
    dfa_alphabets = '';
    dfa_start_state = '';
    dfa_accepted_states = '';
    dfa_transitions = '';
    show_default = false;
  }

  // selectAutomaton
  // Return type: void
  // Parameter type(s): 'NFA' | 'DFA' | 'REGEX'
  // Sets the currently selected automaton type and resets inputs.
  function selectAutomaton(type: 'NFA' | 'DFA' | 'REGEX') {
    selected_automaton = type;
    resetInputs();
  }

  // insertDefault
  // Return type: void
  // Parameter type(s): none
  // Populates the NFA or DFA input fields with default example data.
  function insertDefault() {
    show_default = true;
    if (selected_automaton === 'NFA') {
      nfa_states = 'q0,q1,q2';
      nfa_alphabets = 'a,b';
      nfa_start_state = 'q0';
      nfa_accepted_states = 'q2';
      nfa_transitions = 'q0,a->q1\nq1,b->q2';
    } else if (selected_automaton === 'DFA') {
      dfa_states = 'A,B';
      dfa_alphabets = '0,1';
      dfa_start_state = 'A';
      dfa_accepted_states = 'B';
      dfa_transitions = 'A,0->A\nA,1->B\nB,0->A\nB,1->B';
    }
  }

  // removeDefault
  // Return type: void
  // Parameter type(s): none
  // Clears all NFA/DFA input fields from the default example data.
  function removeDefault() {
    show_default = false;
    resetInputs();
  }

  let previous_inputs: typeof input_rows = [];

  $: {
    const inputs_changed =
      input_rows.length !== previous_inputs.length ||
      input_rows.some((row, index) => {
        const prev_row = previous_inputs[index];
        return !prev_row || row.type !== prev_row.type || row.regex !== prev_row.regex;
      });

    if (inputs_changed) {
      handleInputChange();
      previous_inputs = [...input_rows];
    }
  }

  let selected_automaton: 'NFA' | 'DFA' | 'REGEX' | null = null;
  let show_default = false;

  let nfa_states = '';
  let nfa_alphabets = '';
  let nfa_start_state = '';
  let nfa_accepted_states = '';
  let nfa_transitions = '';

  let dfa_states = '';
  let dfa_alphabets = '';
  let dfa_start_state = '';
  let dfa_accepted_states = '';
  let dfa_transitions = '';
</script>

<div class="phase-inspector">
  <div class="source-code-section">
    <div class="lexor-heading">
      <h1 class="lexor-heading-h1">LEXING</h1>
    </div>
    <h3 class="source-code-header">Source Code</h3>
    <pre class="source-display">{source_code || 'No source code available'}</pre>
  </div>

  <div class="automaton-btn-row">
    <button
      class="automaton-btn {selected_automaton === 'NFA' ? 'selected' : ''}"
      on:click={() => selectAutomaton('NFA')}
      type="button">NFA</button
    >
    <button
      class="automaton-btn {selected_automaton === 'DFA' ? 'selected' : ''}"
      on:click={() => selectAutomaton('DFA')}
      type="button">DFA</button
    >
    <button
      class="automaton-btn {selected_automaton === 'REGEX' ? 'selected' : ''}"
      on:click={() => selectAutomaton('REGEX')}
      type="button">Regular Expression</button
    >
    {#if selected_automaton && !show_default}
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
    {#if selected_automaton && show_default}
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

  {#if selected_automaton === 'REGEX'}
    <div>
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
          {#each input_rows as row, i}
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

        {#if input_rows[input_rows.length - 1].type && input_rows[input_rows.length - 1].regex}
          <button class="add-button" on:click={addNewRow}>
            <span>+</span>
          </button>
        {/if}
      </div>

      {#if form_error}
        <div class="form-error">{form_error}</div>
      {/if}

      <div class="button-container">
        <button
          class="submit-button"
          class:shifted={show_generate_button}
          on:click={handleSubmit}
        >
          Submit
        </button>

        {#if show_generate_button}
          <button class="generate-button" on:click={generateTokens}> Generate Tokens </button>
        {/if}
      </div>

      {#if submission_status.show}
        <div
          class="status-message"
          class:success={submission_status.success === true}
          class:info={submission_status.message === 'info'}
        >
          {submission_status.message}
        </div>
      {/if}
    </div>
  {:else if selected_automaton === 'NFA'}
    <div class="automaton-section">
      <label>
        States:
        <input class="automaton-input" bind:value={nfa_states} placeholder="e.g. q0,q1,q2" />
      </label>
      <label>
        Alphabets:
        <input class="automaton-input" bind:value={nfa_alphabets} placeholder="e.g. a,b" />
      </label>
      <label>
        Start State:
        <input class="automaton-input" bind:value={nfa_start_state} placeholder="e.g. q0" />
      </label>
      <label>
        Accepted States:
        <input
          class="automaton-input"
          bind:value={nfa_accepted_states}
          placeholder="e.g. q2"
        />
      </label>
      <label class="automaton-transitions-label">
        Transitions:
        <textarea
          class="automaton-input automaton-transitions"
          bind:value={nfa_transitions}
          placeholder="e.g. q0,a->q1&#10;q1,b->q2"
        ></textarea>
      </label>
    </div>
  {:else if selected_automaton === 'DFA'}
    <div class="automaton-section">
      <label>
        States:
        <input class="automaton-input" bind:value={dfa_states} placeholder="e.g. A,B" />
      </label>
      <label>
        Alphabets:
        <input class="automaton-input" bind:value={dfa_alphabets} placeholder="e.g. 0,1" />
      </label>
      <label>
        Start State:
        <input class="automaton-input" bind:value={dfa_start_state} placeholder="e.g. A" />
      </label>
      <label>
        Accepted States:
        <input
          class="automaton-input"
          bind:value={dfa_accepted_states}
          placeholder="e.g. B"
        />
      </label>
      <label class="automaton-transitions-label">
        Transitions:
        <textarea
          class="automaton-input automaton-transitions"
          bind:value={dfa_transitions}
          placeholder="e.g. A,0->A&#10;A,1->B"></textarea
        >
      </label>
    </div>
  {/if}
</div>

<style>
  .source-code-header {
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
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    position: relative;
  }

  .lexor-heading {
    justify-items: center;
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
    color: #001a6e;
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
    border-color: #001a6e;
    box-shadow: 0 0 0 2px rgba(32, 87, 129, 0.1);
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
    background: #001a6e;
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
    background: #27548a;
  }

  .button-container {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin-top: 0.75rem;
  }

  .submit-button {
    padding: 0.6rem 1.5rem;
    background: #001a6e;
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
    background: #27548a;
  }

  :global(html.dark-mode) .submit-button {
    background: #cccccc;
    color: #041a47;
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
    0% {
      opacity: 0;
    }
    10% {
      opacity: 1;
    }
    90% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
  :global(html.dark-mode) .source-code-header {
    color: #ebeef1;
  }

  :global(html.dark-mode) .lexor-heading-h1 {
    color: #ebeef1;
  }

  :global(html.dark-mode) .source-display {
    color: black;
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
    color: #001a6e;
    font-weight: 500;
    font-size: 0.95rem;
    cursor: pointer;
    transition: border-color 0.2s, background 0.2s, color 0.2s;
    outline: none;
  }
  .automaton-btn.selected {
    border-color: #001a6e;
    background: #e6edfa;
    color: #001a6e;
  }
  .automaton-btn:not(.selected):hover {
    border-color: #7da2e3;
    background: #f5f8fd;
  }
  .default-toggle-btn {
    margin-left: 1.2rem;
    padding: 0.4rem 0.7rem;
    border-radius: 1.2rem;
    border: 2px solid #e5e7eb;
    background: white; 
    color: #001a6e;
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
    color: #001a6e;
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
    color: #001a6e;
    transition: border-color 0.2s;
    width: 100%;
    box-sizing: border-box;
  }
  .automaton-input:focus {
    border-color: #001a6e;
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
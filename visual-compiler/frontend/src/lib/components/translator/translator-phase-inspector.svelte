<script lang="ts">
  // Prop to receive the source code from the main page
  export let source_code: string;

  // --- Reactive State ---
  // Holds all the translation rules. Starts with one empty rule.
  let rules = [{ tokenSequence: '', lines: [''] }];
  // Flags to control UI state
  let isSubmitted = false;
  let translationSuccessful = false;

  // --- Functions ---

  /**
   * Adds a new, empty translation rule block.
   */
  function addRule() {
    rules = [...rules, { tokenSequence: '', lines: [''] }];
  }

  /**
   * Removes a translation rule block by its index.
   */
  function removeRule(ruleIndex: number) {
    // Prevent removing the very last rule block
    if (rules.length > 1) {
      rules = rules.filter((_, i) => i !== ruleIndex);
    }
  }

  /**
   * Adds a new line to a specific rule.
   */
  function addLine(ruleIndex: number) {
    rules[ruleIndex].lines = [...rules[ruleIndex].lines, ''];
    rules = rules; // Trigger reactivity
  }

  /**
   * Removes a specific line from a rule.
   */
  function removeLine(ruleIndex: number, lineIndex: number) {
    // Prevent removing the last line
    if (rules[ruleIndex].lines.length > 1) {
      rules[ruleIndex].lines = rules[ruleIndex].lines.filter((_, i) => i !== lineIndex);
      rules = rules; // Trigger reactivity
    }
  }

  /**
   * Handles the submission of the translation rules.
   */
  async function handleSubmit() {
    const isValid = rules.every(
      (rule) => rule.tokenSequence.trim() !== '' && rule.lines.every((line) => line.trim() !== '')
    );

    if (!isValid) {
      alert('Please fill out all token sequences and lines before submitting.');
      return;
    }

    const apiPayload = {
      rules: rules.map((rule) => ({
        token_sequence: rule.tokenSequence,
        translation_lines: rule.lines,
      })),
    };

    console.log('API Payload:', JSON.stringify(apiPayload, null, 2));

    // SIMULATED API CALL
    // In a real application, you would make your API call here.
    // const response = await fetch('/api/translate/rules', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(apiPayload)
    // });
    // if(response.ok) {
    //   isSubmitted = true;
    //   alert('Submission successful! You can now perform the translation.');
    // } else {
    //   alert('Submission failed. Please check the console.');
    // }

    // For demonstration, we'll just assume it passed
    isSubmitted = true;
  }

  /**
   * Placeholder for the final translation action.
   */
  function handleTranslate() {
    console.log('Performing translation with the submitted rules...');
    // Add your logic to call the translation endpoint here
    translationSuccessful = true;
    alert('Code translated successfully!');
  }
</script>

<div class="inspector-container">
    <h1 class="heading">TRANSLATING </h1>
  <div class="section">
    <h3 class="section-heading">Source Code</h3>
    <div class="code-block-wrapper">
      <pre class="code-block">{source_code || 'No source code available.'}</pre>
    </div>
  </div>

  <div class="section">
    <h3 class="section-heading">Translator Input</h3>
    <div class="rules-container">
      {#each rules as rule, ruleIndex}
        <div class="rule-block">
          <div class="rule-header">
            <h4 class="rule-title">Translation Rule {ruleIndex + 1}</h4>
            <button
              class="remove-btn"
              on:click={() => removeRule(ruleIndex)}
              disabled={rules.length <= 1}
              title="Remove Rule"
            >
              ✕
            </button>
          </div>
          
          <div class="form-group">
            <label for="token-seq-{ruleIndex}">Token Sequence</label>
            <input
              type="text"
              class="input-field"
              id="token-seq-{ruleIndex}"
              bind:value={rule.tokenSequence}
              placeholder="e.g., IDENTIFIER EQUALS NUMBER"
            />
          </div>
    
          {#each rule.lines as line, lineIndex}
            <div class="line-group">
              
              <input
                type="text"
                class="input-field"
                id="line-{ruleIndex}-{lineIndex}"
                bind:value={rules[ruleIndex].lines[lineIndex]}
                placeholder="Line {lineIndex + 1}"
              />
              <button
                class="remove-line-btn"
                on:click={() => removeLine(ruleIndex, lineIndex)}
                disabled={rule.lines.length <= 1}
                title="Remove Line"
              >−</button>
            </div>
              
          {/each}
         <button class="add-btn" on:click={() => addLine(ruleIndex)}>+ Add Line</button>
        </div>
        
      {/each}
       <div>
                <button class="add-rule-btn" on:click={addRule}>+ Add New Rule</button>
              
                <button class="action-btn submit" on:click={handleSubmit} disabled={isSubmitted}>
                {#if isSubmitted}✓ Submitted{:else}Submit Rules{/if}</button>
                 
          </div>
       
    </div>
  
  </div>
  
  <div class="actions">
 

    {#if isSubmitted}
      <button class="action-btn translate" on:click={handleTranslate}>
        Translate Code
      </button>
    {/if}
  </div>

  {#if translationSuccessful}
    <p class="success-message">Translation complete! View the result in the Artifact Viewer.</p>
  {/if}
</div>

<style>
  :root {
    --bg-primary: white;
    --bg-secondary: #eee;
    --bg-tertiary: #4a5568;
    --text-primary: black;
    --text-secondary: #a0aec0;
    --accent-blue: #3182ce;
    --accent-green: #38a169;
    --accent-red: #e53e3e;
    --accent-purple: #805ad5;
    --accent-orange: #dd6b20;
    --border-color: #4a5568;
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  }

  .inspector-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1.5rem;
    font-family: var(--font-sans);
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .section-heading {
    color: #444;
    margin-bottom: 0;
    margin-top: 0;
    font-family: 'Times New Roman';
  }

  .heading{
    color: black;
    margin-bottom: 0;
    margin-top: 0;
    font-family: 'Times New Roman';
    text-align: center;
  }
  
  .code-block-wrapper {
    background-color: #eee;
    border-radius: 0.375rem;
  
    padding: 1rem;
  }

  .code-block {
    font-family: 'Fira Code', monospace;
    font-size: 0.9rem;
    white-space: pre-wrap;
    word-wrap: break-word;
    margin: 0;
    
  }

  .rules-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .rule-block {
    background-color: #f5f5f5;
    border-radius: 0.5rem;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .rule-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .rule-title {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
  }

  .form-group, .line-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .line-group {
    flex-direction: row;
    align-items: center;
  }

  label {
    font-weight: 500;
    color: var(--text-secondary);
  }

  .input-field {
    width: 100%;
    padding: 0.75rem;
    background-color: var(--bg-primary);
    border: 1px solid #ddd;
    border-radius: 0.25rem;
    color: var(--text-primary);
    font-size: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .input-field:focus {
    outline: none;
    border-color: var(--accent-blue);
    box-shadow: 0 0 0 2px rgba(49, 130, 206, 0.5);
  }

  .add-btn, .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 600;
    transition: background-color 0.2s, transform 0.1s;
     width: 45%;
  }
  
  .add-btn {
    background-color: #001a6e;
    color: white;
    align-self: flex-start;
  }

   .add-rule-btn {
    justify-content: center;
    gap: 0.5rem;
    background-color: #eef2f7;
    color: #001a6e;
    border: 1px dashed #c0c7d3;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    width: 45%;
    margin-right: 1rem;
    margin-left: 0.8rem;
  }

  .add-btn:hover, .action-btn:hover, .submit:hover { background-color: #1a317d; }
  .add-btn:active { transform: scale(0.98); }

  .add-rule {
    align-self: flex-start;
    margin-top: 1rem;
  }

  .remove-btn, .remove-line-btn {
    background: transparent;
    color: var(--text-secondary);
    border: none;
    border-radius: 50%;
    cursor: pointer;
    font-size: 1.2rem;
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s, color 0.2s;
  }

  .remove-btn:hover, .remove-line-btn:hover {
    background-color: var(--accent-red);
    color: white;
  }
  
  .remove-line-btn {
    font-size: 1.5rem;
    margin-left: 0.5rem;
  }

  .remove-btn:disabled, .remove-line-btn:disabled {
    color: #4a5568;
    cursor: not-allowed;
    background-color: transparent;
  }

  .actions {
    display: flex;
    gap: 1rem;
    padding-top: 1rem;
  }
  

  .submit {
    background-color: #001a6e;
    color: white;
  }
 

  .submit:disabled {
    background-color: var(--accent-green);
    cursor: default;
  }

  .translate {
    background-color: var(--accent-orange);
    color: white;
  }
  .translate:hover { background-color: #c05621; }

  .success-message {
    color: var(--accent-green);
    font-weight: bold;
    text-align: center;
  }
</style>
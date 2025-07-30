<script lang="ts">
  export let source_code: string;
    import { AddToast } from '$lib/stores/toast';

  let rules = [{ tokenSequence: '', lines: [''] }];

  let isSubmitted = false;
  let translationSuccessful = false;



  function addRule() {
    rules = [...rules, { tokenSequence: '', lines: [''] }];
  }

 
  function removeRule(ruleIndex: number) {

    if (rules.length > 1) {
      rules = rules.filter((_, i) => i !== ruleIndex);
    }
  }


  function addLine(ruleIndex: number) {
    rules[ruleIndex].lines = [...rules[ruleIndex].lines, ''];
    rules = rules; 
  }

  function removeLine(ruleIndex: number, lineIndex: number) {

    if (rules[ruleIndex].lines.length > 1) {
      rules[ruleIndex].lines = rules[ruleIndex].lines.filter((_, i) => i !== lineIndex);
      rules = rules;
    }
  }

  
  async function handleSubmit() {

    const user_id = localStorage.getItem('user_id');
        if (!user_id) {
            AddToast('User not logged in.', 'error');
            return;
        }

    const isValid = rules.every(
      (rule) => rule.tokenSequence.trim() !== '' && rule.lines.every((line) => line.trim() !== '')
    );

    if (!isValid) {
      AddToast('Please fill out all token sequences and lines before submitting.');
      return;
    }

    const apiPayload = {
      rules: rules.map((rule) => ({
        token_sequence: rule.tokenSequence,
        translation_lines: rule.lines,
      })),
    };

    if (apiPayload.rules.some(rule => rule.token_sequence.trim() === '' || rule.translation_lines.length === 0)) {
      AddToast('All rules must have a token sequence and at least one translation line.');
      return;
    }

    if (apiPayload.rules.length === 0) {
      AddToast('No rules to submit. Please add at least one rule.');
      return;
    }
   

     try {
            const response = await fetch('http://localhost:8080/api/translation//readRules', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(apiPayload)
            });

            if (!response.ok) 
            {
                const error_data = await response.json();
                throw new Error(error_data.details || 'Failed to submit translation input');
            }

            const result = await response.json();
            AddToast('Translation input received successfully!', 'success');

        } catch (error) 
        {
            console.error('Translation input Error:', error);
            AddToast(String(error), 'error');
        }

    
    isSubmitted = true;
  }

 
  function handleTranslate() {
    console.log('Performing translation with the submitted rules...');


    translationSuccessful = true;
    AddToast('Code translated successfully!');
  }
</script>

<div class="inspector-container">
    <h1 class="heading">TRANSLATING </h1>
  <div class="section">
    <h3 class="section-heading1">Source Code</h3>
    <div class="code-block-wrapper">
      <pre class="code-block">{source_code || 'No source code available.'}</pre>
    </div>
  </div>

  <div class="section">
    <h2 class="section-heading">Translation Rules</h2>
    <div class="rules-container">
      {#each rules as rule, ruleIndex}
        <div class="rule-block">
          <div class="form-group">
            <div class="rule-header" >
                <label class="form-label" for="token-seq-{ruleIndex}" style="margin-right: auto;">Token Sequence</label>
                <button
                    class="remove-btn"
                    on:click={() => removeRule(ruleIndex)}
                    disabled={rules.length <= 1}
                    title="Remove Rule"
                    style="margin-left: auto;"
                >
                    ✕
                </button>
            </div>

      
            
            <input
              type="text"
              class="input-field"
              id="token-seq-{ruleIndex}"
              bind:value={rule.tokenSequence}
              placeholder="Enter token sequence"
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
                aria-label="Remove Line"
              ><svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg></button>
            </div>
              
          {/each}
         <button class="add-line" on:click={() => addLine(ruleIndex)}>+ Add Line</button>
        </div>
        
      {/each}
       <div>
                <button class="add-rule-btn" on:click={addRule}>+ Add New Rule</button>
              
                <button class="action-btn submit" on:click={handleSubmit} >
                  Submit Rules
                </button>

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
    --accent-orange: #6c757d;
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

  .form-label{
    color:#1a2a4a;
    font-family: 'Times New Roman';
    font-weight: 500;
  }
  .section {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .add-line {
       display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 100%;
    background-color: #eef2f7;
    color: #001a6e;
    border: 1px dashed #c0c7d3;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .add-line:hover {
    border-color: #001a6e;
  }
  .section-heading {
    color: #001a6e;
    margin-bottom: 0;
    margin-top: 0;
    font-family: 'Times New Roman';
    font-size: 1.25rem;
    font-weight: 600;
    text-align: center;
  }

  .section-heading1{
    color: #444;
    margin-bottom: 0;
    margin-top: 0;
    font-family: 'Times New Roman';
  }

  .action-btn.translate {
    background-color: var(--accent-orange);
    color: white;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 600;
    transition: background-color 0.2s, transform 0.1s;
    display: block;
    margin: 0 auto;
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
    position: relative;
}

.rule-header {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    margin-bottom: 0.5rem;
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

   .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    font-size: 0.9rem;
    font-weight: 600;
    transition: background-color 0.2s, transform 0.1s;
     width: 45%;
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

   .action-btn:hover, .submit:hover { background-color: #1a317d; }



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
    color: red;
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
    background-color: #1a317d;
    cursor: default;
  }

  .translate {
    background-color: var(--accent-orange);
    color: white;
  }
  .translate:hover { background-color: rgb(98,102,109) }

  .success-message {
    color: var(--accent-green);
    font-weight: bold;
    text-align: center;
  }
</style>
<script lang="ts">
  import { onMount } from 'svelte';

  export let sourceCode = '';

  // --- DATA STRUCTURE for CFG ---
  // A single translation for a rule (e.g., "expr + term")
  interface Translation {
    id: number;
    value: string;
  }

  // A full grammar rule (e.g., "expr -> expr + term | term")
  interface Rule {
    id: number;
    nonTerminal: string;
    translations: Translation[];
  }

  // --- STATE MANAGEMENT ---
  let grammar: Rule[] = [];
  let rule_id_counter = 0;
  let translation_id_counter = 0;

  onMount(() => {
    addNewRule();
  });

  function addNewRule() {
    rule_id_counter++;
    translation_id_counter++;
    grammar = [
      ...grammar,
      {
        id: rule_id_counter,
        nonTerminal: '',
        translations: [{ id: translation_id_counter, value: '' }]
      }
    ];
  }

  function addTranslation(ruleId: number) {
    translation_id_counter++;
    grammar = grammar.map(rule => {
      if (rule.id === ruleId) {
        return {
          ...rule,
          translations: [
            ...rule.translations,
            { id: translation_id_counter, value: '' }
          ]
        };
      }
      return rule;
    });
  }

  function handleSubmitGrammar() {
    console.log('Submitting Grammar:', grammar);
    // API call would go here
  }
</script>

<div class="phase-inspector">
  <div class="parser-heading">
    <h1 class="parser-heading-h1">PARSING</h1>
  </div>

  <div class="source-code-section">
    <h3 class="source-code-header">Source Code</h3>
    <pre class="source-display">{sourceCode || 'No source code available'}</pre>
  </div>
  
  <div class="grammar-editor">
    <h3>Context-Free Grammar</h3>
    <div class="rules-container">
      {#each grammar as rule, i (rule.id)}
        <div class="rule-row">
          <div class="rule-label">
            {#if i === 0}
              <span class="start-label">Start →</span>
            {/if}
          </div>

          <div class="rule-inputs">
            <input 
              type="text" 
              class="non-terminal-input" 
              placeholder="expr"
              bind:value={rule.nonTerminal}
            />
            <span class="arrow">→</span>
            <div class="translations-container">
              {#each rule.translations as translation, j (translation.id)}
                <input 
                  type="text" 
                  class="translation-input" 
                  placeholder="expr + term"
                  bind:value={translation.value}
                />
                <!-- THE CHANGE IS HERE: The '|' separator has been removed -->
              {/each}
              <button class="add-translation-btn" on:click={() => addTranslation(rule.id)}>
                +
              </button>
            </div>
          </div>
        </div>
      {/each}
    </div>
    <button class="add-rule-btn" on:click={addNewRule}>
      + Add New Rule
    </button>
  </div>
  
  <div class="button-container">
    <button class="submit-button" on:click={handleSubmitGrammar}>
      Submit Grammar
    </button>
  </div>
</div>

<style>
  .phase-inspector {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .parser-heading-h1 {
    color: #001A6E;
    text-align: center;
    margin-top: 0;
  }
  .source-code-section {
    margin-bottom: 1rem;
  }
  .source-code-header {
    color: #444;
    margin-bottom: 0.5rem;
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
  .grammar-editor {
    background: #f8f9fa;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }
  .grammar-editor h3 {
    margin-top: 0;
    margin-bottom: 1.5rem;
    color: #001A6E;
  }
  .rules-container {
    display: flex;
    flex-direction: column;
    gap: 1rem; 
  }
  .rule-row {
    display: flex;
    align-items: center;
  }
  .rule-label {
    width: 80px; 
    flex-shrink: 0;
    text-align: right;
    padding-right: 0.75rem;
  }
  .start-label {
    font-weight: bold;
    color: #001A6E;
    font-family: monospace;
    font-size: 1.1rem;
  }
  .rule-inputs {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    /* This prevents the container from overflowing its parent */
    overflow: hidden; 
  }
  .non-terminal-input {
    flex: 0 0 60px; 
    padding: 0.6rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: monospace;
    text-align: center;
    width : 90%;
  }
  .arrow {
    font-size: 1.2rem;
    color: #555;
    font-weight: bold;
  }
  .translations-container {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex: 1;
    overflow-x: auto; 
    padding-bottom: 0.5rem;
  }
  .translation-input {
    padding: 0.6rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: monospace;
    width: 5rem; 
  }
  .add-translation-btn {
    background: #e0e0e0;
    color: #333;
    border: none;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    cursor: pointer;
    font-size: 1.2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
    flex-shrink: 0;
  }
  .add-translation-btn:hover {
    background-color: #ccc;
  }
  .add-rule-btn {
    display: block;
    margin: 1.5rem auto 0;
    padding: 0.5rem 1rem;
    border: 1px dashed #001A6E;
    background: transparent;
    color: #001A6E;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
  }
  .button-container {
    text-align: center;
    margin-top: 1rem;
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
  }

  /* Dark Mode Styles */
  :global(html.dark-mode) .parser-heading-h1 { color: #ebeef1; }
  :global(html.dark-mode) .source-code-header { color: #ebeef1; }
  :global(html.dark-mode) .source-display { color: black; }
  :global(html.dark-mode) .grammar-editor { background: #1f2937; }
  :global(html.dark-mode) .grammar-editor h3 { color: #ebeef1; }
  :global(html.dark-mode) .start-label { color: #60a5fa; }
  :global(html.dark-mode) .non-terminal-input, 
  :global(html.dark-mode) .translation-input {
    background-color: #2d3748;
    border-color: #4b5563;
    color: #f0f0f0;
  }
  :global(html.dark-mode) .arrow { color: #9ca3af; }
  :global(html.dark-mode) .add-translation-btn {
    background-color: #4b5563;
    color: #f0f0f0;
  }
  :global(html.dark-mode) .add-rule-btn {
    border-color: #60a5fa;
    color: #60a5fa;
  }
  :global(html.dark-mode) .submit-button {
    background-color: #cccccc;
    color: #041a47;
  }
</style>

<script lang="ts">
  import type { NodeType } from '$lib/types';

  let current_step = 1;
  const total_steps = 4;

  // nextStep
  // Return type: void
  // Parameter type(s): none
  // Increments the current step of the tutorial if not on the last step.
  function nextStep() {
    if (current_step < total_steps) current_step++;
  }

  // prevStep
  // Return type: void
  // Parameter type(s): none
  // Decrements the current step of the tutorial if not on the first step.
  function prevStep() {
    if (current_step > 1) current_step--;
  }

  const sample_code = `if (count &gt; 0) { return true; }`;
  const token_result = `[KEYWORD:if] [PUNCTUATION:(] [IDENTIFIER:count] [OPERATOR:>] [NUMBER:0] [PUNCTUATION:)] [PUNCTUATION:{] [KEYWORD:return] [BOOLEAN:true] [PUNCTUATION:;] [PUNCTUATION:}]`;

  const regex_basics = [
    { type: 'Literal Match', pattern: 'hello', example: 'Matches "hello" in text' },
    { type: 'Any Character', pattern: 'h.llo', example: 'Matches "hello", "hallo"' },
    { type: 'Digits', pattern: '\\d+', example: 'Matches "42", "123"' },
    { type: 'Letters', pattern: '[a-zA-Z]+', example: 'Matches "hello", "Test"' }
  ];

  const regex_examples = [
    { type: 'Keyword', pattern: '(int|string|)', example: 'int, string' },
    { type: 'Integer', pattern: '\\d+', example: '42, 123, 0' },
    { type: 'Identifier', pattern: '[a-zA-Z_]+', example: 'count, total_sum' },
    { type: 'Operator', pattern: '[+\\-*/]', example: '+, -, *, /' },
    { type: 'Assignment', pattern: '=', example: '=' }
  ];

  const tokens_by_type = [
    {
      type: 'Keyword',
      tokens: ['if', 'return']
    },
    {
      type: 'Punctuation',
      tokens: ['(', ')', '{', '}', ';']
    },
    {
      type: 'Identifier',
      tokens: ['count']
    },
    {
      type: 'Operator',
      tokens: ['>']
    },
    {
      type: 'Number',
      tokens: ['0']
    },
    {
      type: 'Boolean',
      tokens: ['true']
    }
  ];
</script>

<div class="phase-tutorial">
  <div class="tutorial-content">
    <h2>What is Lexing?</h2>

    <p class="description">
      A lexer breaks down source code into tokens. Each token has a type and value. Tokens are
      identified using regular expressions that match specific patterns.
    </p>
    <div class="separator"></div>
    <h3>Steps to Follow:</h3>
    <div class="tutorial-container">
      <div class="tutorial-content-area">
        {#if current_step === 1}
          <div class="tutorial-step">
            <h3>1. Source Code Input</h3>
            <p>Start by entering your source code in the input box:</p>
            <div class="code-sample">
              <code>{@html sample_code}</code>
            </div>
          </div>
        {:else if current_step === 2}
          <div class="tutorial-step">
            <h3>2. Regular Expression Basics</h3>
            <p>Regular expressions use special characters to match patterns:</p>
            <div class="pattern-example">
              {#each regex_basics as { type, pattern, example }}
                <div class="type-pair">
                  <span class="type">{type}</span>
                  <code class="pattern-code">{pattern}</code>
                  <span class="example">{example}</span>
                </div>
              {/each}
            </div>
          </div>
        {:else if current_step === 3}
          <div class="tutorial-step">
            <h3>3. Common Token Patterns</h3>
            <p>Here are some common patterns used in lexical analysis:</p>
            <div class="type-example">
              {#each regex_examples as { type, pattern, example }}
                <div class="type-pair">
                  <span class="type">{type}</span>
                  <code class="pattern-code">{pattern}</code>
                  <span class="example">Examples: {example}</span>
                </div>
              {/each}
            </div>
          </div>
        {:else if current_step === 4}
          <div class="tutorial-step">
            <h3>4. Tokenisation Result</h3>
            <div class="token-table">
              <div class="table-header">
                <span class="header-type">Type</span>
                <span class="header-tokens">Tokens</span>
              </div>
              {#each tokens_by_type as { type, tokens }}
                <div class="table-row">
                  <span class="token-type">{type}</span>
                  <span class="token-values">{tokens.join(', ')}</span>
                </div>
              {/each}
            </div>
          </div>
        {/if}
        <div class="separator"></div>
      </div>

      <div class="navigation-container">
        <div class="navigation">
          <button class="nav-button" disabled={current_step === 1} on:click={prevStep}>
            ← Previous
          </button>
          <span class="step-counter">{current_step} / {total_steps}</span>
          <button class="nav-button" disabled={current_step === total_steps} on:click={nextStep}>
            Next →
          </button>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .phase-tutorial {
    padding: 1rem;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .tutorial-content {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  h2 {
    color: #001a6e;
    margin-bottom: 1rem;
    font-size: 1.5rem;
  }

  h3 {
    color: #444;
    margin: 1.5rem 0 0.5rem 0;
  }

  .description {
    line-height: 1.6;
    color: #333;
    margin: 1rem 0;
    font-size: 0.95rem;
  }

  .separator {
    border-bottom: 1px solid #e5e7eb;
    margin: 1rem 0;
    width: 100%;
  }

  .tutorial-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 1rem;
  }

  .tutorial-content-area {
    flex: 1;
    overflow-y: auto;
    padding-bottom: 2rem;
  }

  .tutorial-step {
    margin-bottom: 2rem;
  }

  .code-sample {
    background: #f1f3f5;
    padding: 0.75rem;
    border-radius: 4px;
    margin: 0.5rem 0;
    font-family: monospace;
  }

  .pattern-example,
  .type-example {
    background: #f8f9fa;
    padding: 1rem;
    border-radius: 8px;
    margin: 1rem 0;
    max-width: 650px;
  }

  .pattern-code {
    font-family: 'Fira Code', monospace;
    font-size: 0.9rem;
    background: #f1f3f5;
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    color: #001a6e;
    width: fit-content;
  }

  .example {
    color: #666;
    font-size: 0.9rem;
  }

  .navigation-container {
    margin-top: auto;
    margin-bottom: 1rem;
    padding: 1rem 0;
    display: flex;
    justify-content: center;
  }

  .navigation {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    max-width: 500px;
  }

  .nav-button {
    padding: 0.5rem 1rem;
    background: #001a6e;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    min-width: 100px;
    font-size: 0.9rem;
  }

  .nav-button:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }

  .step-counter {
    color: #666;
    font-size: 0.9rem;
  }

  .token-table {
    background: white;
    border-radius: 8px;
    overflow: hidden;
    margin: 1rem 0;
    max-width: 500px;
    border: 1px solid #e5e7eb;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .table-header {
    display: grid;
    grid-template-columns: 90px 1fr;
    background: #001a6e;
    color: white;
    padding: 0.5rem 0.75rem;
    font-weight: 500;
    font-size: 0.9rem;
    border-bottom: 2px solid #e5e7eb;
  }

  .table-row {
    display: grid;
    grid-template-columns: 90px 1fr;
    padding: 0.5rem 0.75rem;
    border-bottom: 1px solid #e5e7eb;
    background: white;
  }

  .table-row:last-child {
    border-bottom: none;
  }

  .table-row:hover {
    background: #f8f9fa;
  }

  .token-type {
    color: #001a6e;
    font-weight: 500;
    font-family: 'Fira Code', monospace;
    font-size: 0.85rem;
  }

  .header-tokens,
  .token-values {
    padding-left: 3rem;
  }

  .token-values {
    font-family: 'Fira Code', monospace;
    color: #666;
    font-size: 0.85rem;
  }

  .type-pair {
    margin: 0.75rem 0;
    display: grid;
    grid-template-columns: 90px 130px 1fr;
    align-items: center;
    gap: 0.75rem;
  }

  .type {
    color: #001a6e;
    font-weight: 500;
    font-family: 'Fira Code', monospace;
    font-size: 0.9rem;
  }

  /* Dark Mode Styles */
  :global(html.dark-mode) .phase-tutorial {
    background: #1a2a4a;
  }

  :global(html.dark-mode) h2 {
    color: #ebeef1;
  }

  :global(html.dark-mode) h3 {
    color: #f0f0f0;
  }

  :global(html.dark-mode) .description {
    color: #d1d5db;
  }

  :global(html.dark-mode) .separator {
    border-bottom-color: #374151;
  }

  :global(html.dark-mode) .code-sample {
    background: #2d3748;
    color: #f0f0f0;
  }

  :global(html.dark-mode) .pattern-example,
  :global(html.dark-mode) .type-example {
    background: #2d3748;
  }

  :global(html.dark-mode) .pattern-code {
    background: #4a5568;
    color: #f0f0f0;
  }

  :global(html.dark-mode) .example {
    color: #d1d5db;
  }

  :global(html.dark-mode) .navigation-container {
    background: #1a2a4a;
    border-top-color: #374151;
  }

  :global(html.dark-mode) .token-table {
    background: #2d3748;
    border-color: #4a5568;
  }

  :global(html.dark-mode) .table-row {
    background: #2d3748;
    border-bottom-color: #4a5568;
  }

  :global(html.dark-mode) .table-row:hover {
    background: #3c4a62;
  }

  :global(html.dark-mode) .token-type {
    color: #60a5fa;
  }

  :global(html.dark-mode) .token-values {
    color: #d1d5db;
  }

  :global(html.dark-mode) .type {
    color: #60a5fa;
  }

  :global(html.dark-mode) .nav-button {
    background: #cccccc;
    color: #041a47;
  }

  :global(html.dark-mode) .nav-button:disabled {
    background: #4b5563;
    color: #041a47;
  }
</style>
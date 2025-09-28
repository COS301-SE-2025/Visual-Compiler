<script lang="ts">
    let current_step = 1;
    const total_steps = 5;

    // Reactive declaration to ensure reactivity
    $: reactive_step = current_step;

    function nextStep() {
        if (current_step < total_steps) {
            current_step++;
            console.log("Next Step:", current_step); // Debug log
        }
    }

    function prevStep() {
        if (current_step > 1) {
            current_step--;
            console.log("Previous Step:", current_step); // Debug log
        }
    }

    const example_rules = [
        {
            sequence: 'INTEGER, OPERATOR, INTEGER',
            translation: 'mov   rax, {INTEGER}\nadd   rax, {INTEGER}',
            source: '12 + 13',
            output: 'mov   rax, 12\nadd   rax, 13'
        },
        {
            sequence: 'KEYWORD, IDENTIFIER, ASSIGNMENT, BOOLEAN',
            translation: 'mov   rax, [{IDENTIFIER}]\nmov   rax, {BOOLEAN}',
            source: 'bool found = true;',
            output: 'mov   rax, [found]\nmov   rax, 1'
        }
    ];

    const critical_donts = [
        { mistake: 'Not using token types', wrong: 'int, result, =, 1', right: 'KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER' },
        { mistake: 'Skipping placeholders', wrong: 'mov [result], 1', right: 'MOV [{IDENTIFIER}], {INTEGER}' }
    ];
</script>

<div class="phase-tutorial">
    <div class="tutorial-content">
        <h2>What is Translating?</h2>
        <p class="description">
            This tool converts your source code (like Go) into target code (like Assembly) using rules you define. Think of it as teaching the computer how to rewrite your code in a different language.
        </p>
        <div class="separator"></div>
        <h3>Steps to Follow:</h3>
        <div class="tutorial-container">
            <div class="tutorial-content-area">
                {#if reactive_step === 1}
                    <div class="tutorial-step">
                        <h3>1. Complete Lexing and Parsing</h3>
                        <div class="step-content">
                            <p> Before using the translator, you must complete the lexing and parsing phases:</p>
                            <div class="token-list">
                                <div class="token-item">
                                    <span class="token-type">Lexer</span>
                                    <span class="token-desc">Converts your source code into a token stream</span>
                                </div>
                                <div class="token-item">
                                    <span class="token-type">Parser</span>
                                    <span class="token-desc">Converts your token stream into a syntax tree</span>
                                </div>
                            </div>
                            <div class="info-block">
                                <p>Once the syntax tree is ready, you can proceed to define translation rules and generate target code. Note that your token types are used in translation. ✨</p>
                            </div>
                        </div>
                    </div>
                {:else if reactive_step === 2}
                    <div class="tutorial-step">
                        <h3>2. Define Translation Rules</h3>
                        <div class="step-content">
                            <p>Click "+ Add New Rule" for each pattern you want to translate:</p>
                            <div class="example-block">
                                <div class="example-header">
                                    <span>Example Rule for Addition</span>
                                </div>
                                <div class="example-content">
                                    <div class="type-pair">
                                        <span class="type">Sequence:</span>
                                        <code class="pattern-code">IDENTIFIER, OPERATOR, INTEGER</code>
                                    </div>
                                    <div class="type-pair">
                                        <span class="type">Translation:</span>
                                        <code class="pattern-code">add [&#123;IDENTIFIER&#125;], &#123;INTEGER&#125;</code>
                                    </div>
                                    <div class="output-section">
                                        <p>This converts <code> var + 13 </code> to <code> add [var], 13 </code></p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                {:else if reactive_step === 3}
                    <div class="tutorial-step">
                        <h3>3. Submit and Generate</h3>
                        <div class="step-content">
                            <p> Click "Submit Rules" to generate your output:</p>
                            <div class="code-sample">
                                <div class="code-header">
                                    <span>Example Output for <code>result = 12</code>:</span>
                                </div>
                                <pre class="code-output">
Target Code Output:

mov    rax, [result]
add    rax, 12
                </pre>
                            </div>
                        </div>
                    </div>
                {:else if reactive_step === 4}
                    <div class="tutorial-step">
                        <h3>4. Ready-to-Use Examples</h3>
                        <div class="step-content">
                            <p> Here are some examples to get you started:</p>
                            {#each example_rules as { sequence, translation, source, output }}
                                <div class="example-block">
                                    <div class="example-header">
                                        <span>Source Code: <code>{source}</code></span>
                                    </div>
                                    <div class="example-content">
                                        <div class="type-pair">
                                            <span class="type">Sequence:</span>
                                            <code class="pattern-code">{sequence}</code>
                                        </div>
                                        <div class="type-pair">
                                            <span class="type">Translation:</span>
                                            <code class="pattern-code">{translation}</code>
                                        </div>
                                        <div class="output-section">
                                            <span>Target Code:</span>
                                            <pre class="code-output">{output}</pre>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {:else if reactive_step === 5}
                    <div class="tutorial-step">
                        <h3>5. Important Warnings</h3>
                        <div class="step-content">
                            <p>Avoid these common mistakes:</p>
                            {#each critical_donts as { mistake, wrong, right }}
                                <div class="example-block">
                                    <div class="example-header">
                                        <span>{mistake}</span>
                                    </div>
                                    <div class="example-content">
                                        <div class="type-pair">
                                            <span class="type">❌ Wrong:</span>
                                            <code class="pattern-code">{wrong}</code>
                                        </div>
                                        <div class="type-pair">
                                            <span class="type">✅ Right:</span>
                                            <code class="pattern-code">{right}</code>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    </div>
                {/if}
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
    color: var(--text-color, #333333); /* Light mode default */
}

.tutorial-content {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
}

h2 {
    color: var(--heading-color, #001a6e); /* Light mode default */
    margin-bottom: 1rem;
    font-size: 1.5rem;
}

h3 {
    color: var(--subheading-color, #444444); /* Light mode default */
    margin: 1.5rem 0 0.5rem 0;
}

.description {
    line-height: 1.6;
    margin: 1rem 0;
    font-size: 0.95rem;
}

.separator {
    border-bottom: 1px solid var(--separator-color, #e5e7eb); /* Light mode default */
    margin: 1rem 0;
    width: 100%;
}

.tutorial-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 0;
    padding: 1rem;
}

.tutorial-content-area {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding-bottom: 2rem;
}

.tutorial-step {
    margin-bottom: 2rem;
    background-color: var(--block-background, #ffffff); /* Light mode default */
}

.code-sample {
    background: var(--code-background, #f1f3f5); /* Light mode default */
    padding: 0.75rem;
    border-radius: 4px;
    margin: 0.5rem 0;
    font-family: monospace;
    color: var(--code-text-color, #333333); /* Light mode default */
    border: 1px solid var(--separator-color, #e5e7eb);
}

.type-pair {
    margin: 0.75rem 0;
    display: grid;
    grid-template-columns: 150px 1fr; /* Increased from 90px to 150px */
    align-items: center;
    gap: 0.75rem;
}

.pattern-code {
    font-family: 'Fira Code', monospace;
    font-size: 0.9rem;
    background: var(--code-background, #f1f3f5); /* Light mode default */
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    color: var(--code-highlight-color, #001a6e); /* Light mode default */
    width: fit-content;
    border: 1px solid var(--separator-color, #e5e7eb);
}

.navigation-container {
    position: sticky;
    margin-top: 1rem;
    bottom: 0;
    background: var(--navigation-background, #ffffff); /* Light mode default */
    z-index: 10;
    padding: 0rem 0;
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
    background: var(--button-background, #BED2E6); /* Light mode default */
    color: var(--button-text-color, #000000); /* Light mode default */
    border: none;
    border-radius: 4px;
    cursor: pointer;
    min-width: 100px;
    font-size: 0.9rem;
    transition: background-color 0.2s ease, transform 0.2s ease;
}

.nav-button:hover {
    background: var(--button-hover-background, #a8bdd1); /* Light mode hover */
    transform: translateY(-2px);
}

.nav-button:disabled {
    background: var(--button-disabled-background, #cccccc); /* Light mode default */
    cursor: not-allowed;
}

.step-counter {
    color: var(--step-counter-color, #666666); /* Light mode default */
    font-size: 0.9rem;
}

.step-content {
    padding: 0.5rem;
}

.code-header {
    margin-bottom: 0.5rem;
    font-weight: bold;
}

.code-output {
    margin: 0;
    padding: 0.5rem;
    background: var(--code-background, #f1f3f5);
    border-radius: 4px;
}

.token-list {
    display: grid;
    gap: 0.75rem;
    padding: 0.5rem;
}

.token-item {
    display: grid;
    grid-template-columns: 180px 1fr; /* Increased from 120px to 180px */
    gap: 1rem;
    align-items: center;
}

/* Add specific styling for sequence and translation labels */
.type-pair .type {
    white-space: nowrap; /* Prevents text from wrapping */
    min-width: 150px; /* Ensures minimum width for the label */
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.token-type {
    font-weight: bold;
    color: var(--code-highlight-color, #001a6e);
    white-space: nowrap; /* Prevents text from wrapping */
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.example-block {
    background: var(--code-background, #f1f3f5);
    border-radius: 6px;
    margin: 1rem 0;
    padding: 1rem;
    border: 1px solid var(--separator-color, #e5e7eb);
}

.example-header {
    margin-bottom: 1rem;
    font-weight: bold;
}

.example-content {
    padding: 0.5rem;
}

.output-section {
    margin-top: 1rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--separator-color, #e5e7eb);
}

.info-block {
    background: var(--code-background, #f1f3f5);
    border-radius: 6px;
    margin: 1rem 0;
    padding: 1rem;
    border: 1px solid var(--separator-color, #e5e7eb);
}

/* Dark Mode Styles */
:global(html.dark-mode) {
    --background-color: #1a202c;
    --text-color: #e2e8f0;
    --heading-color: #90cdf4;
    --subheading-color: #a0aec0;
    --separator-color: #4a5568;
    --block-background: #1a2a4a;
    --code-background: #2d3748;
    --code-text-color: #e2e8f0;
    --code-highlight-color: #63b3ed;
    --navigation-background: #1a2a4a;
    --button-background: #3182ce;
    --button-hover-background: #2b6cb0;
    --button-text-color: #e2e8f0;
    --button-disabled-background: #718096;
    --step-counter-color: #a0aec0;
}

:global(html.dark-mode) .nav-button{
    background: #001A6E;  
    color: #ffffff;
}
:global(html.dark-mode) .nav-button:hover {
    background: #002a8e; 
}

:global(html.dark-mode) .nav-button:disabled {
    background: #2d3748;
    color: #9ca3af;
}
</style>
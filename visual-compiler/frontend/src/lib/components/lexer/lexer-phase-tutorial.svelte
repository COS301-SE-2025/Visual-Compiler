<script lang="ts">
    import type { NodeType } from '$lib/types';

    let current_step = 1;
    const total_steps = 6;

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

    const sample_code = `if (count > 0) { return true; }`;
    const token_result = `[KEYWORD:if] [PUNCTUATION:(] [IDENTIFIER:count] [OPERATOR:>] [NUMBER:0] [PUNCTUATION:)] [PUNCTUATION:{] [KEYWORD:return] [BOOLEAN:true] [PUNCTUATION:;] [PUNCTUATION:}]`;

    const regex_basics = [
        { type: 'Exact Match', pattern: 'hello', example: 'Matches "hello"' },
        { type: 'Range', pattern: '[0-9]', example: 'Matches "1" and "3"' },
        { type: 'Alternation', pattern: '(a|b)', example: 'Matches "a" or "b"' },
        { type: 'Star (0 or More)', pattern: '2*', example: 'Matches "", "2", "22"' },
        { type: 'Plus (1 or More)', pattern: '2+', example: 'Matches "2" and "22"' }
    ];

    const regex_examples = [
        { type: 'Keyword', pattern: '(int|return)', example: 'int, return' },
        { type: 'Identifier', pattern: '[a-zA-Z_]+', example: 'count, total' },
        { type: 'Comparison', pattern: '(<|>|=)', example: '<, >, =' },
        { type: 'Integer', pattern: '[0-9]+', example: '42, 123, 0' },
        { type: 'Boolean', pattern: '(true|false)', example: 'true, false' }
    ];

    const tokens_by_type = [
        {
            type: 'Keyword',
            tokens: ['if']
        },
        {
            type: 'Identifier',
            tokens: ['count']
        },
        {
            type: 'Comparison',
            tokens: ['>']
        },
        {
            type: 'Integer',
            tokens: ['0']
        },
        {
            type: 'Keyword',
            tokens: ['return']
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
                        <h3>2.1 Regular Expression Basics</h3>
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
                        <a
                            class="yt-link"
                            href="https://www.youtube.com/watch?v=nIp604p0M8M&ab_channel=lydia"
                            target="_blank"
                            rel="noopener"
                        >
                            📺 Watch: Regular Expressions explained (YouTube)
                        </a>
                    </div>
                {:else if current_step === 3}
                    <div class="tutorial-step">
                        <h3>2.2 Deterministic Finite Automata (DFA)</h3>
                        <div class="automata-explainer">
                            <div class="automata-example">
                                <div class="automata-label">
                                    Example: DFA accepting strings ending with <code>1</code>
                                </div>
                                <div class="dfa-nfa-table">
                                    <table>
                                        <thead>
                                            <tr>
                                                <th>State</th>
                                                <th>Input 0</th>
                                                <th>Input 1</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            <tr>
                                                <td><b>q0</b> <span class="state-badge start">start</span></td>
                                                <td>q0</td>
                                                <td>q1</td>
                                            </tr>
                                            <tr>
                                                <td><b>q1</b> <span class="state-badge accept">accept</span></td>
                                                <td>q0</td>
                                                <td>q1</td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                                <div class="automata-how">
                                    <b>How it works:</b>
                                    <ul>
                                        <li>Starts at <b>q0</b>.</li>
                                        <li>Every <code>0</code> keeps it in <b>q0</b>.</li>
                                        <li>Every <code>1</code> moves to <b>q1</b>.</li>
                                        <li>If the string ends in <b>q1</b>, it is accepted.</li>
                                    </ul>
                                </div>
                                <div class="test-cases">
                                    <b>Test Cases:</b>
                                    <ul>
                                        <li><code>0101</code> → Ends in 1 → <span class="accept">Accept</span></li>
                                        <li><code>1000</code> → Ends in 0 → <span class="reject">Reject</span></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <a
                            class="yt-link"
                            href="https://www.youtube.com/watch?v=PK3wL7DXuuw&ab_channel=lydia"
                            target="_blank"
                            rel="noopener"
                        >
                            📺 Watch: DFA explained (YouTube)
                        </a>
                    </div>
                {:else if current_step === 4}
                    <div class="tutorial-step">
                        <h3>2.3 Non-deterministic Finite Automata (NFA)</h3>
                        <div class="automata-explainer">
                            <div class="automata-example">
                                <div class="automata-label">
                                    Example: NFA accepting strings ending with <code>1</code>
                                </div>
                                <div class="dfa-nfa-table">
                                    <table>
                                        <thead>
                                            <tr>
                                                <th>State</th>
                                                <th>Input 0</th>
                                                <th>Input 1</th>
                                                <th>ε</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            <tr>
                                                <td><b>q0</b> <span class="state-badge start">start</span></td>
                                                <td>q0</td>
                                                <td>q0, q1</td>
                                                <td>∅</td>
                                            </tr>
                                            <tr>
                                                <td><b>q1</b> <span class="state-badge accept">accept</span></td>
                                                <td>∅</td>
                                                <td>∅</td>
                                                <td>∅</td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                                <div class="automata-how">
                                    <b>How it works:</b>
                                    <ul>
                                        <li>
                                            At any <code>1</code>, the NFA can choose to stay in <b>q0</b> or move to
                                            <b>q1</b>.
                                        </li>
                                        <li>If any possible path ends in <b>q1</b>, the string is accepted.</li>
                                    </ul>
                                </div>
                                <div class="test-cases">
                                    <b>Test Cases:</b>
                                    <ul>
                                        <li>
                                            <code>0101</code> → Possible path: q0→q0→q0→q1 →
                                            <span class="accept">Accept</span>
                                        </li>
                                        <li>
                                            <code>1000</code> → No path reaches q1 → <span class="reject">Reject</span>
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <a
                            class="yt-link"
                            href="https://www.youtube.com/watch?v=W8Uu0inPmU8&ab_channel=lydia"
                            target="_blank"
                            rel="noopener"
                        >
                            📺 Watch: NFA explained (YouTube)
                        </a>
                    </div>
                {:else if current_step === 5}
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
                {:else if current_step === 6}
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
                        <a
                            class="yt-link"
                            href="https://www.youtube.com/watch?v=MZ9NZdZteG4&ab_channel=NesoAcademy"
                            target="_blank"
                            rel="noopener"
                        >
                            📺 Watch: Tokenisation in Lexical Analysis (YouTube)
                        </a>
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
        height: 100%; /* or 100% if inside a fixed-height parent */
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }

    .tutorial-content {
        display: flex;
        flex-direction: column;
        height: 100%;
        min-height: 0;
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
        position: sticky;
        margin-top: 1rem;
        bottom: 0;
        background: #fff;
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
        background: #BED2E6;
        color: 000000;
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

    .dfa-nfa-table {
        margin: 0.7rem 0 0.7rem 0;
    }

    .dfa-nfa-table table {
        border-collapse: collapse;
        width: 100%;
        max-width: 400px;
        background: #f1f3f5;
        border-radius: 6px;
        overflow: hidden;
        font-size: 0.97rem;
    }

    .dfa-nfa-table th,
    .dfa-nfa-table td {
        border: 1px solid #e5e7eb;
        padding: 0.4rem 0.7rem;
        text-align: center;
    }

    .dfa-nfa-table th {
        background: #e6edfa;
        color: #001a6e;
    }

    .accept {
        color: #16a34a;
        font-weight: 600;
    }

    .reject {
        color: #dc2626;
        font-weight: 600;
    }

    .automata-explainer {
        background: #f8f9fa;
        border-radius: 8px;
        padding: 0.7rem 1.5rem 0.7rem 1.5rem;
        margin: 1.2rem 0 1.5rem 0;
        box-shadow: 0 1px 4px rgba(30, 64, 175, 0.04);
        max-width: 650px;
        margin-top: 0.5rem;
    }

    .automata-example {
        margin-top: 1.2rem;
    }

    .automata-label {
        font-weight: 600;
        color: #1e40af;
        margin-bottom: 0.5rem;
        font-size: 1.08rem;
        margin-top: 0;
    }

    .state-badge {
        display: inline-block;
        font-size: 0.75em;
        font-weight: 600;
        background: #e0e7ff;
        color: #1e40af;
        border-radius: 6px;
        padding: 0.1em 0.6em;
        margin-left: 0.4em;
        vertical-align: middle;
    }
    .state-badge.accept {
        background: #d1fae5;
        color: #047857;
    }
    .state-badge.start {
        background: #fef9c3;
        color: #b45309;
    }

    .automata-how {
        margin: 1rem 0 0.5rem 0;
        font-size: 0.98rem;
    }

    .automata-how ul {
        margin: 0.3rem 0 0 1.2rem;
        padding: 0;
    }

    .test-cases {
        margin-top: 0.7rem;
        font-size: 0.97rem;
    }

    h3 {
        color: #444;
        margin: 0.3rem 0 0rem 0;
    }

    .yt-link {
        display: inline-block;
        margin-top: 1rem;
        color: #1e40af;
        font-weight: 500;
        text-decoration: none;
        font-size: 0.98rem;
        transition: color 0.2s;
    }
    .yt-link:hover {
        color: #0a2540;
        text-decoration: underline;
    }

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
    }

    :global(html.dark-mode) .nav-button {
        background: #001A6E;
        color: #ffffff;
    }

    :global(html.dark-mode) .nav-button:disabled {
        background: #4b5563;
        color: #9ca3af; 
    }

    :global(html.dark-mode) .step-counter {
        color: #d1d5db;
    }

    :global(html.dark-mode) .token-table {
        background: #2d3748;
        border-color: #4a5568;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3); 
    }

    :global(html.dark-mode) .table-header {
        background: #003366;
        border-bottom-color: #4a5568;
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

    :global(html.dark-mode) .dfa-nfa-table table {
        background: #2d3748;
    }

    :global(html.dark-mode) .dfa-nfa-table th,
    :global(html.dark-mode) .dfa-nfa-table td {
        border-color: #4a5568;
        color: #d1d5db;
    }

    :global(html.dark-mode) .dfa-nfa-table th {
        background: #003366; 
        color: #f0f0f0;
    }

    :global(html.dark-mode) .accept {
        color: #34d399; 
    }

    :global(html.dark-mode) .reject {
        color: #f87171;
    }

    :global(html.dark-mode) .automata-explainer {
        background: #2d3748;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.3); 
    }

    :global(html.dark-mode) .automata-label {
        color: #60a5fa;
    }

    :global(html.dark-mode) .state-badge {
        background: #1c3d5a;
        color: #90cdf4;
    }
    :global(html.dark-mode) .state-badge.accept {
        background: #14532d;
        color: #6ee7b7;
    }
    :global(html.dark-mode) .state-badge.start {
        background: #78350f;
        color: #fcd34d;
    }

    :global(html.dark-mode) .automata-how,
    :global(html.dark-mode) .test-cases {
        color: #d1d5db;
    }

    :global(html.dark-mode) .yt-link {
        color: #60a5fa;
    }
    :global(html.dark-mode) .yt-link:hover {
        color: #90cdf4;
    }


    ::-webkit-scrollbar {
        width: 11px;
    }
    ::-webkit-scrollbar-track {
        background: #f1f1f1;
    }
    ::-webkit-scrollbar-thumb {
        background-color: #888;
        border-radius: 10px;
    }
    ::-webkit-scrollbar-thumb:hover {
        background: #555;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-track {
        background: #2d3748;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb {
        background-color: #4a5568; 
        border-color: #2d3748; 
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb:hover {
        background: #616e80;
    }


</style>
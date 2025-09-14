<script lang="ts">
    let current_step = 1;
    const total_steps = 5;

    function nextStep() {
        if (current_step < total_steps) current_step++;
    }
    function prevStep() {
        if (current_step > 1) current_step--;
    }
</script>

<div class="phase-tutorial">
    <div class="tutorial-content">
        <h2>Optimisation</h2>
        <p class="description">
            Optimisation is the process of automatically improving the intermediate code generated from your source code. The goal is to produce a final program that runs faster, uses less memory, or results in a smaller binary file, all <b>without changing the program's intended behavior or output</b>.
        </p>
        <div class="separator"></div>
        <div class="tutorial-container">
            <div class="tutorial-content-area">
                {#if current_step === 1}
                    <div class="tutorial-step">
                        <h3>1. What is optimisation?</h3>
                        <p>
                            Optimisations are performed on an <b>Intermediate Representation (IR)</b> of the code, such as an Abstract Syntax Tree (AST), which is a structured, tree-like model that the compiler can easily analyze and manipulate. The code you provided performs source-to-source transformation, meaning it takes Go code, optimizes its AST, and outputs new, improved Go code.
                        </p>
                        <div class="opt-benefits">
                            <b>Think of it like a master editor reviewing a draft of a document:</b>
                            <ul>
                                <li>The core message (the program's function) must remain unchanged.</li>
                                <li>The editor's job is to make the text more concise, efficient, and clear to follow (faster execution).</li>
                                <li>They might remove redundant phrases (dead code), replace long descriptions with simpler terms (constant folding), or reorganize paragraphs for better flow (loop unrolling).</li>
                            </ul>
                        </div>
                    </div>
                {:else if current_step === 2}
                    <div class="tutorial-step">
                        <h3>2. Benefits of Optimisation</h3>
                        <div class="opt-benefits">
                            <p><b>Increased Execution Speed:</b> By eliminating unnecessary operations and streamlining logic, the CPU can execute the program's instructions much faster. For large-scale applications or frequently run code, these micro-savings add up to significant performance gains.</p>
                            
                            <p><b>Reduced Memory Footprint:</b> Optimisations can remove unused variables, functions, and data structures, leading to a program that consumes less RAM during execution.</p>
                            
                            <p><b>Smaller Binary Size:</b> When dead code is eliminated and loops are unrolled (sometimes), the final compiled executable file (.exe, .bin) can become smaller, which is critical for embedded systems with limited storage.</p>
                            
                            <p><b>Improved Energy Efficiency:</b> Faster, leaner programs require less CPU time to complete their tasks, which directly translates to lower power consumption and longer battery life on mobile devices.</p>
                        </div>
                    </div>
                {:else if current_step === 3}
                    <div class="tutorial-step">
                        <h3>3. Constant Folding</h3>
                        <p>
                            <b>Concept:</b> Constant folding is the process of evaluating constant expressions during compile time rather than waiting to evaluate them at runtime. The compiler identifies sub-expressions within the code that are composed entirely of constants (literal values like 5 or variables declared as const), calculates their result once, and replaces the entire expression with the computed constant value.
                        </p>
                        <div class="opt-example">
                            <div class="code-block">
                                <b>Before Constant Folding:</b>
                                <pre><code>const radius = 10
circumference := 2 * 3.14159 * radius</code></pre>
                            </div>
                            <p>At runtime, the CPU would have to load the values and perform the multiplication <span class="opt-highlight">2 * 3.14159 * 10</span>.</p>
                            <div class="code-block">
                                <b>After Constant Folding:</b>
                                <pre><code>const radius = 10
circumference := 62.8318</code></pre>
                            </div>
                        </div>
                        <p>The compiler pre-computed <span class="opt-highlight">2 * 3.14159 * 10 = 62.8318</span>. Now, at runtime, the program simply loads the value 62.8318, which is much faster.</p>
                    </div>
                {:else if current_step === 4}
                    <div class="tutorial-step">
                        <h3>4. Dead Code Elimination</h3>
                        <p>
                            <b>Concept:</b> Dead code elimination (DCE) is the removal of code that has no effect on the program's output. This code is "dead" because its execution is either impossible or its results are never used. There are two main types:
                        </p>
                        <ul>
                            <li><b>Unreachable Code:</b> Code that can never be executed, such as statements after a return or inside an <code>if false</code> block.</li>
                            <li><b>Unused Code:</b> Code whose result is computed but never influences the program's final state, such as writing to a variable that is never read.</li>
                        </ul>
                        <div class="opt-example">
                            <div class="code-block">
                                <b>Before Dead Code Elimination:</b>
                                <pre><code>func calculate(a int) int &#123;
    result := a * 2   // This value is calculated...
    debug := true
    if debug &#123;        // ...but this condition is always true...
        log.Println("Debug is on")
    &#125;
    if false &#123;        // ...this block is never executed.
        log.Println("This is dead code")
    &#125;
    return a + 10     // ...and 'result' is never used here.
&#125;</code></pre>
                            </div>
                            <div class="code-block">
                                <b>After Dead Code Elimination:</b>
                                <pre><code>func calculate(a int) int &#123;
    return a + 10
&#125;</code></pre>
                            </div>
                        </div>
                        <p>The compiler removes the unused <code>result</code> variable, the always-true <code>if</code> block (though the log call might remain due to side effects), and the unreachable <code>if false</code> block.</p>
                    </div>
                {:else if current_step === 5}
                    <div class="tutorial-step">
                        <h3>5. Loop Unrolling</h3>
                        <p>
                            <b>Concept:</b> Loop unrolling is an optimisation that reduces the overhead of a loop's control mechanism by explicitly repeating the loop body multiple times per iteration. Loop control overhead includes the instructions for incrementing the counter, checking the loop condition, and jumping back to the start.
                        </p>
                        <div class="opt-example">
                            <div class="code-block">
                                <b>Before Loop Unrolling (High Overhead):</b>
                                <pre><code>// 3 iterations, 3 condition checks, 3 increments
for i := 0; i &lt; 3; i++ &#123;
    process(i) // The actual work
&#125;</code></pre>
                            </div>
                            <div class="code-block">
                                <b>After Loop Unrolling (Low Overhead):</b>
                                <pre><code>// 0 iterations, 0 condition checks, 0 increments
process(0) // Just the work
process(1)
process(2)</code></pre>
                            </div>
                        </div>
                        <p>The CPU now executes three <code>process</code> calls in a straight line, which is much faster than managing the loop counter.</p>
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
    }

    .tutorial-content {
        display: flex;
        flex-direction: column;
        height: 100%;
        min-height: 0;
    }

    h2 {
        color: #AFA2D7;
        margin-bottom: 1rem;
        font-size: 1.5rem;
    }

    h3 {
        color: #AFA2D7;
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

    .opt-block {
        background: #F3E8FF;
        border-left: 4px solid #AFA2D7;
        border-radius: 8px;
        padding: 1rem 1.2rem;
        margin: 1rem 0;
        font-size: 0.97rem;
        box-shadow: 0 1px 4px rgba(175, 162, 215, 0.04);
    }

    .opt-block ul {
        margin: 0.5rem 0 0 1.2rem;
    }

    .opt-benefits {
        background: #F8FAFC;
        border-radius: 8px;
        padding: 1rem 1.2rem;
        margin: 1rem 0;
        font-size: 0.97rem;
        border-left: 4px solid #AFA2D7;
    }

    .opt-benefits ul {
        margin: 0.5rem 0 0 1.2rem;
    }

    .opt-benefits p {
        margin: 0.8rem 0;
        line-height: 1.6;
    }

    .opt-example {
        border-radius: 8px;
        margin: 1rem 0;
        font-size: 0.97rem;
    }

    .code-block {
        background: #F1F5F9;
        border-left: 4px solid #AFA2D7;
        border-radius: 8px;
        padding: 1rem 1.2rem;
        margin: 0.5rem 0;
    }

    .opt-highlight {
        color: #8451C7;
        font-weight: bold;
    }

    pre, code {
        font-family: 'Fira Code', monospace;
        font-size: 0.97rem;
        background: #F3E8FF;
        border-radius: 4px;
        padding: 0.2rem 0.5rem;
        margin: 0.2rem 0;
        color: #1a2a4a;
    }

    .code-block pre, .code-block code {
        background: transparent;
        padding: 0;
        margin: 0.5rem 0 0 0;
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
        box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.03);
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
        background: #AFA2D7;
        color: #fff;
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

    /* Dark Mode Styles */
    :global(html.dark-mode) .phase-tutorial {
        background: #1a2a4a;
    }

    :global(html.dark-mode) h2 {
        color: #8451C7;
    }

    :global(html.dark-mode) h3 {
        color: #8451C7;
    }

    :global(html.dark-mode) .description {
        color: #d1d5db;
    }

    :global(html.dark-mode) .separator {
        border-bottom-color: #374151;
    }

    :global(html.dark-mode) .opt-block {
        background: #2D1B69;
        border-left-color: #8451C7;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .opt-benefits {
        background: #2d3748;
        border-left-color: #8451C7;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .code-block {
        background: #2d3748;
        border-left-color: #8451C7;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .opt-highlight {
        color: #AFA2D7;
    }

    :global(html.dark-mode) pre,
    :global(html.dark-mode) code {
        background: #2D1B69;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .code-block pre,
    :global(html.dark-mode) .code-block code {
        background: transparent;
        color: #E0E7FF;
    }

    :global(html.dark-mode) .navigation-container {
        background: #1a2a4a;
    }

    :global(html.dark-mode) .nav-button {
        background: #8451C7;
        color: #F3E8FF;
    }

    :global(html.dark-mode) .nav-button:disabled {
        background: #4b5563;
        color: #9ca3af;
    }

    :global(html.dark-mode) .step-counter {
        color: #d1d5db;
    }
</style>

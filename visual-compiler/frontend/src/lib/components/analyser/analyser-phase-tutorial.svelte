<script lang="ts">
	let current_step = 1;
	const total_steps = 5;

	function nextStep() {
		if (current_step < total_steps) current_step++;
	}
	function prevStep() {
		if (current_step > 1) current_step--;
	}

	const sample_code = `int x = 5 + 3;`;
	const sample_tokens = `[KEYWORD:int] [IDENTIFIER:x] [ASSIGNMENT:=] [INTEGER:5] [OPERATOR:+] [INTEGER:3] [SEPARATOR:;]`;

	// Semantic-specific data
	const scope_rules = [
		{ rule: 'Scope starts at { and ends at }' , example: '{ int y = x; }' },
		{ rule: 'Variables must be declared before use', example: 'x = 5; ❌ (if x not declared)' },
		{ rule: 'Inner scopes can access outer variables', example: '{ int x; { int y = x; } }' }
	];

	const type_rules = [
		{ rule: 'int = int + int → ✅', example: 'int x = 5 + 3;' },
		{ rule: 'int = float → ⚠️ (may truncate)', example: 'int x = 3.14;' },
		{ rule: 'bool = int > int → ✅', example: 'bool valid = (x > 5);' }
	];

	const symbol_table = [
		{ name: 'x', type: 'int', scope: '0 (global)' },
		{ name: 'y', type: 'int', scope: '1 (block)' }
	];
</script>

<div class="phase-tutorial">
	<div class="tutorial-content">
		<h2>What is Semantic Analysis?</h2>
		<p class="description">
			Semantic analysis ensures the <b>meaning</b> of a program is correct according to the language rules.
			It checks declarations, types, and scopes after parsing. Scope checks and type checks are performed on the syntax tree.
		</p>
		<div class="separator"></div>
		<h3>Steps to Follow:</h3>
		<div class="tutorial-container">
			<div class="tutorial-content-area">
				{#if current_step === 1}
					<div class="tutorial-step">
						<h3>1. Input: Syntax Tree from Parser</h3>
						<p>The parser generates a syntax tree, which becomes the input for semantic analysis.</p>
						<div class="syntax-tree">
							<pre>
STATEMENT
  |
  +-- DECLARATION
  |     |
  |     +-- TYPE: KEYWORD("int")
  |     +-- IDENTIFIER: ID("x")
  |     +-- ASSIGNMENT: OP("=")
  |     +-- EXPRESSION
  |           |
  |           +-- TERM: INTEGER("5")
  |           +-- OPERATOR: OP("+")
  |           +-- TERM: INTEGER("3")
  +-- SEPARATOR(";")</pre>
						</div>
					</div>
				{:else if current_step === 2}
					<div class="tutorial-step">
						<h3>2. Scope Checking</h3>
						<p> Scope rules can be 1 or more rules that indicate how a new scope starts and ends.
                            It is used to perform scope checks to ensure that variables and functions called have been declared, and are within the relevant scope 
                        </p>
						<div class="cfg-section">
                            <div class="code-sample">
                                <pre><code>int x = 5;
&#123;
    int y = x;  <span class="comment-valid">// ✅ x is in outer scope</span>
&#125;
y = 3;      <span class="comment-error">// ❌ y is out of scope</span></code></pre>
                            </div>
							<div class="cfg-block">
                                <b>Scope Rules:</b>
                                <ul class="scope-rules-list">
                                    {#each scope_rules as r}
                                    <li>
                                        <span class="cfg-rule">{r.rule}</span>:
                                        <div class="rule-example"><code>{r.example}</code></div>
                                    </li>
                                    {/each}
                                </ul>
                            </div>
							
						</div>
					</div>
				{:else if current_step === 3}
					<div class="tutorial-step">
						<h3>3. Type Checking</h3>
						<p> Ensure types are compatible in expressions, assignments, and function calls.
                            Rules to indicate which variable types are allowed to be assigned to which.
                            These rules are used to perform type checks, ensuring valid types are assigned to variables

                        </p>
						<div class="cfg-rules">
							<table>
								<thead>
									<tr>
										<th>Type Rule</th>
										<th>Example</th>
									</tr>
								</thead>
								<tbody>
									{#each type_rules as r}
										<tr>
											<td><span class="cfg-rule">{r.rule}</span></td>
											<td><span class="cfg-example">{r.example}</span></td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
                {:else if current_step === 4}  
                <div class="tutorial-step">
                    <h3>4. Link Grammar Rules to Semantic Checks</h3>
                    
                    <div class="link-grammar-section">
                    <div class="link-grammar-block">
                        <h4>🧩 Why Link Grammar Rules?</h4>
                        <p>To apply semantic checks (scope and type), the compiler must connect specific parser rules to relevant checks.</p>
                        
                        <div class="link-steps">
                        <div class="step-content">
                            <b>Steps:</b>
                            <ol>
                            <li>Identify grammar rules that represent:
                                <ul>
                                <li>Variable declarations</li>
                                <li>Assignments</li>
                                <li>Function calls</li>
                                </ul>
                            </li>
                            <li>Attach semantic checks to those rules</li>
                            </ol>
                        </div>
                        </div>
                    </div>
                    
                    <div class="link-example">
                        <h4>✅ Example:</h4>
                        <div class="grammar-rule">
                        <span class="rule-name">Declaration</span> → 
                        <span class="rule-part">TYPE</span> 
                        <span class="rule-part">ID</span> 
                        <span class="rule-part">SEMICOLON</span>
                        </div>
                        <div class="link-checks">
                        <p>You link:</p>
                        <ul>
                            <li>A scope check: <span class="check-detail">Is ID being redeclared?</span></li>
                            <li>A symbol table update: <span class="check-detail">Add ID with its type and scope</span></li>
                        </ul>
                        </div>
                    </div>
                    </div>
                </div>
				{:else if current_step === 5}  <!-- Assuming this is step 5 -->
                <div class="tutorial-step">
                    <h3>5. Scope Table</h3>
                    
                    <div class="scope-table-section">
                    <div class="scope-table-intro">
                        <h4>📋 What Is a Scope Table?</h4>
                        <p>A scope table (also called a symbol table) tracks:</p>
                        <ul>
                        <li>Each declared identifier (variable or function)</li>
                        <li>Its type</li>
                        <li>Its scope level</li>
                        </ul>
                    </div>

                    <div class="scope-table-demo">
                        <div class="scope-table-structure">
                        <h4>Structure:</h4>
                        <div class="symbol-table">
                            <div class="table-header">
                            <span>Type</span>
                            <span>Name</span>
                            <span>Scope</span>
                            </div>
                            <div class="table-row">
                            <span>int</span>
                            <span>blue</span>
                            <span>0</span>
                            </div>
                        </div>
                        <p class="table-caption">This means:<br>
                        • A variable named <strong>blue</strong><br>
                        • Of type <strong>int</strong><br>
                        • Declared in scope <strong>0</strong> (typically global)</p>
                        </div>
                    </div>
                    <div class="scope-table-usage">
                        <h4>How It's Used:</h4>
                        <p>During semantic analysis, every variable use is looked up in the scope table.</p>
                        <div class="usage-example">
                            <div class="valid-case">
                            <span class="icon">✅</span> Found in current/enclosing scope → Valid
                            </div>
                            <div class="error-case">
                            <span class="icon">❌</span> Not found or out of scope → Semantic error
                            </div>
                        </div>
                        </div>
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

.token-table {
    margin-top: 2rem;
    display: table;
    width: 100%;
    border-collapse: collapse;
    margin: 0.5rem 0;
}
.table-header,
.table-row {
    display: table-row;
}
.table-header {
    background: #e6edfa;
    color: #001a6e;
    font-weight: 600;
}
.cfg-section {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin: 1rem 0;
}
.cfg-block {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 1rem 1.2rem;
    max-width: 450px;
    font-size: 0.97rem;
    box-shadow: 0 1px 4px rgba(30, 64, 175, 0.04);
}
.cfg-nt {
    color: #1e40af;
    font-weight: 600;
    font-family: 'Fira Code', monospace;
}
.cfg-t {
    color: #047857;
    font-weight: 600;
    font-family: 'Fira Code', monospace;
}
.cfg-rules {
    margin: 1rem 0;
}
.cfg-rules table {
    border-collapse: collapse;
    width: 100%;
    max-width: 600px;
    background: #f1f3f5;
    border-radius: 6px;
    overflow: hidden;
    font-size: 0.97rem;
}
.cfg-rules th,
.cfg-rules td {
    border: 1px solid #e5e7eb;
    padding: 0.4rem 0.7rem;
    text-align: left;
}
.cfg-rules th {
    background: #e6edfa;
    color: #001a6e;
    font-size: 1.05rem;
}
.cfg-rule {
    font-family: 'Fira Code', monospace;
    color: #001a6e;
}
.cfg-example {
    color: #666;
    font-size: 0.95rem;
}
.syntax-tree {
    background: #f8fafc;
    border-radius: 8px;
    padding: 1rem 1.2rem;
    font-family: 'Fira Code', monospace;
    font-size: 0.97rem;
    margin: 1rem 0;
    color: #1e293b;
    overflow-x: auto;
}
.step-breakdown {
    margin: 1rem 0 0 1.2rem;
    font-size: 0.97rem;
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

.code-sample {
  background: #f8f9fa;
  border-radius: 6px;
  padding: 1rem;
  margin: 1rem 0;
  overflow-x: auto;
}

.code-sample pre {
  margin: 0;
  padding: 0;
  white-space: pre;
  overflow-x: auto;
  font-family: monospace;
  background: #f8f9fa;
  border-radius: 4px;
  padding: 12px;
}

.code-sample code {
  display: block;
  padding: 0.5rem;
}

.comment-valid {
  color: #059669; /* Green for valid */
}

.comment-error {
  color: #dc2626; /* Red for errors */
}

.scope-rules-list {
  margin: 0.5rem 0;
  padding-left: 1.2rem;
}

.scope-rules-list li {
  margin-bottom: 1rem;
}

.rule-example {
  margin-top: 0.3rem;
  padding-left: 1rem;
}

.rule-example code {
  background: #f1f3f5;
  padding: 0.3rem 0.5rem;
  border-radius: 4px;
  display: inline-block;
}
/* Scope Table Styles */
.scope-table-section {
  margin: 1rem 0;
}

.scope-table-intro ul {
  margin: 0.5rem 0 1rem 0.5rem;
}

.scope-table-demo {
  display: flex;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.symbol-table {
  width: 100%; /* Full width */
  max-width: none; /* Remove previous max-width constraint */
  min-width: 500px; /* Minimum width */
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  overflow: hidden;
  margin: 0.5rem 0;
}

.table-header, .table-row {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr; /* Wider middle column */
  padding: 0.7rem;
  text-align: center;
}

.table-header span,
.table-row span {
  flex: 1;
  text-align: center;
}

.table-caption {
  margin: 0.5rem 0;
  font-size: 0.9rem;
  color: #4b5563;
}

.usage-example {
  margin-top: 1rem;
}

.valid-case, .error-case {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  margin: 0.3rem 0;
  border-radius: 4px;
}

.valid-case {
  background: #f0fdf4;
  color: #166534;
}

.error-case {
  background: #fef2f2;
  color: #991b1b;
}

.icon {
  font-size: 1.1rem;
}

/* Dark Mode */
:global(html.dark-mode) .symbol-table {
  border-color: #4a5568;
}

:global(html.dark-mode) .table-header {
  background: #4a5568;
}

:global(html.dark-mode) .table-row {
  border-color: #4a5568;
}

:global(html.dark-mode) .table-caption {
  color: #d1d5db;
}

:global(html.dark-mode) .valid-case {
  background: #1a4531;
  color: #34d399;
}

:global(html.dark-mode) .error-case {
  background: #5c1a1a;
  color: #f87171;
}

/* Dark mode */
:global(html.dark-mode) .rule-example code {
  background: #2d3748;
}

/* Dark mode support */
:global(html.dark-mode) .code-sample {
  background: #2d3748;
}

/* Dark mode support */
:global(html.dark-mode) .code-sample pre {
  background: #2d3748;
}

:global(html.dark-mode) .comment-valid {
  color: #34d399;
}

:global(html.dark-mode) .comment-error {
  color: #f87171;
}

@media (max-width: 700px) {
    .cfg-section {
        flex-direction: column;
        gap: 1rem;
    }
    .cfg-block {
        max-width: 100%;
    }
}

/* Dark Mode Styles */
:global(html.dark-mode) h2 {
    color: #ebeef1;
}

:global(html.dark-mode) h3 {
    color: #ebeef1;
}

:global(html.dark-mode) .description {
    color: #ebeef1;
}

:global(html.dark-mode) .code-sample {
    background: #2d3748;
    color: #ebeef1; 
}

:global(html.dark-mode) .token-table .table-header {
    background: #4a5568;
    color: #ebeef1;
}

:global(html.dark-mode) .token-table .header-type,
:global(html.dark-mode) .token-table .header-tokens,
:global(html.dark-mode) .token-table .token-type,
:global(html.dark-mode) .token-table .token-values {
    border-color: #4a5568; 
    background: #2d3748; 
    color: #ebeef1;
}

:global(html.dark-mode) .cfg-block {
    background: #2d3748;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    color: #ebeef1;
}

:global(html.dark-mode) .cfg-nt {
    color: #90cdf4; 
}

:global(html.dark-mode) .cfg-t {
    color: #68d391; 
}

:global(html.dark-mode) .cfg-rules table {
    background: #2d3748;
}

:global(html.dark-mode) .cfg-rules th {
    background: #4a5568;
    color: #ebeef1;
}

:global(html.dark-mode) .cfg-rules th,
:global(html.dark-mode) .cfg-rules td {
    border-color: #4a5568;
    color: #ebeef1;
}

:global(html.dark-mode) .cfg-rule {
    color: #90cdf4;
}

:global(html.dark-mode) .cfg-example {
    color: #a0aec0; 
}

:global(html.dark-mode) .syntax-tree {
    background: #2d3748;
    color: #ebeef1;
}

:global(html.dark-mode) .step-breakdown {
    color: #ebeef1;
}

:global(html.dark-mode) .navigation-container {
    background: #1a2a4a;
}

:global(html.dark-mode) .nav-button {
    background: #001A6E; 
    color: #ffffff;
}

:global(html.dark-mode) .nav-button:disabled {
    background: #718096; 
}

:global(html.dark-mode) .step-counter {
    color: #a0aec0;
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
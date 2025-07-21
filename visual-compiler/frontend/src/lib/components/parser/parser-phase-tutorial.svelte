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

	const cfg_nonterminals = [
		{ name: 'STATEMENT', desc: 'A complete instruction (e.g., int x = 5;)' },
		{ name: 'DECLARATION', desc: 'Variable creation (e.g., int x = 5)' },
		{ name: 'EXPRESSION', desc: 'Math/logic operations (e.g., 5 + 3)' },
		{ name: 'TYPE', desc: 'Data type (e.g., int)' },
		{ name: 'TERM', desc: 'Basic value in expressions (e.g., 5)' }
	];
	const cfg_terminals = [
		{ name: 'KEYWORD', desc: 'Reserved words (e.g., int, str)' },
		{ name: 'IDENTIFIER', desc: 'Variable names (e.g., x, count)' },
		{ name: 'ASSIGNMENT', desc: '=' },
		{ name: 'INTEGER', desc: 'Numbers (e.g., 42)' },
		{ name: 'OPERATOR', desc: 'Math symbols (e.g., +, *)' },
		{ name: 'SEPARATOR', desc: ';' }
	];
	const cfg_rules = [
		{ rule: 'STATEMENT → DECLARATION SEPARATOR', example: 'int x = 5;' },
		{ rule: 'DECLARATION → TYPE IDENTIFIER ASSIGNMENT EXPRESSION', example: 'int x = 5 + 3' },
		{ rule: 'EXPRESSION → TERM OPERATOR TERM', example: '5 + 3' },
		{ rule: 'TERM → INTEGER', example: '5' },
		{ rule: 'TYPE → KEYWORD', example: 'int' }
	];
</script>

<div class="phase-tutorial">
	<div class="tutorial-content">
		<h2>What is Parsing?</h2>
		<p class="description">
			Parsing is the process of analyzing a sequence of tokens (words/symbols) to determine its
			structure according to a set of grammar rules. After the lexer (tokenizer) breaks down source
			code into tokens, the parser takes over to understand the structure and meaning of those
			tokens.
		</p>
		<div class="separator"></div>
		<h3>Steps to Follow:</h3>
		<div class="tutorial-container">
			<div class="tutorial-content-area">
				{#if current_step === 1}
					<div class="tutorial-step">
						<h3>1. Enter Source Code and Generate Tokens</h3>
						<p>Start by entering your source code and generating tokens in the lexer phase:</p>
						<div class="code-sample">
							<code>{@html sample_code}</code>
						</div>
						<b><br />Tokens:</b>
						<div class="token-table">
							<div class="table-header">
								<span class="header-type">Type</span>
								<span class="header-tokens">Token</span>
							</div>
							<div class="table-row">
								<span class="token-type">KEYWORD</span>
								<span class="token-values">int</span>
							</div>
							<div class="table-row">
								<span class="token-type">IDENTIFIER</span>
								<span class="token-values">x</span>
							</div>
							<div class="table-row">
								<span class="token-type">ASSIGNMENT</span>
								<span class="token-values">=</span>
							</div>
							<div class="table-row">
								<span class="token-type">INTEGER</span>
								<span class="token-values">5</span>
							</div>
							<div class="table-row">
								<span class="token-type">OPERATOR</span>
								<span class="token-values">+</span>
							</div>
							<div class="table-row">
								<span class="token-type">INTEGER</span>
								<span class="token-values">3</span>
							</div>
							<div class="table-row">
								<span class="token-type">SEPARATOR</span>
								<span class="token-values">;</span>
							</div>
						</div>
					</div>
				{:else if current_step === 2}
					<div class="tutorial-step">
						<h3>2. Context-Free Grammar (CFG)</h3>
						<p>
							A <b>CFG</b> is a set of rules that define how sentences in a formal language can be constructed.
						</p>
						<div class="cfg-section">
							<div class="cfg-block">
								<b>Variables (Non-Terminals):</b>
								<ul>
									{#each cfg_nonterminals as nt}
										<li><span class="cfg-nt">{nt.name}</span>: {nt.desc}</li>
									{/each}
								</ul>
								(Represent groups of terminals/non-terminals, like statements, expressions, etc.)
							</div>
							<div class="cfg-block">
								<b>Terminals (Tokens):</b>
								<ul>
									{#each cfg_terminals as t}
										<li><span class="cfg-t">{t.name}</span>: {t.desc}</li>
									{/each}
								</ul>
								(Basic symbols that cannot be broken down further, like keywords, identifiers, etc.)
							</div>
						</div>
					</div>
				{:else if current_step === 3}
					<div class="tutorial-step">
						<h3>3. Production Rules</h3>
						<p>Production rules describe how non-terminals expand into other symbols:</p>
						<div class="cfg-rules">
							<table>
								<thead>
									<tr>
										<th>Rule</th>
										<th>Example</th>
									</tr>
								</thead>
								<tbody>
									{#each cfg_rules as r}
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
						<h3>4. Example Syntax Tree</h3>
						<p>For the statement <code>int x = 5 + 3;</code>:</p>
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
  +-- SEPARATOR(";")
              </pre>
						</div>
					</div>
				{:else if current_step === 5}
					<div class="tutorial-step">
						<h3>5. Step-by-Step Breakdown</h3>
						<ol class="step-breakdown">
							<li>
								<b>Root (<span class="cfg-nt">STATEMENT</span>)</b>: Must expand to
								<span class="cfg-nt">DECLARATION</span> <span class="cfg-t">SEPARATOR</span>.
							</li>
							<li>
								<b>DECLARATION</b>: Expands to
								<span class="cfg-nt">TYPE IDENTIFIER ASSIGNMENT EXPRESSION</span>:
								<ul>
									<li><b>TYPE</b> → <span class="cfg-t">KEYWORD("int")</span></li>
									<li><b>IDENTIFIER</b> → <span class="cfg-t">"x"</span></li>
									<li><b>ASSIGNMENT</b> → <span class="cfg-t">"="</span></li>
									<li><b>EXPRESSION</b> → <span class="cfg-nt">TERM OPERATOR TERM</span></li>
								</ul>
							</li>
							<li>
								<b>EXPRESSION</b>: Expands to <span class="cfg-t">5 + 3</span>:
								<ul>
									<li>First <b>TERM</b> → <span class="cfg-t">INTEGER("5")</span></li>
									<li><b>OPERATOR</b> → <span class="cfg-t">"+"</span></li>
									<li>Second <b>TERM</b> → <span class="cfg-t">INTEGER("3")</span></li>
								</ul>
							</li>
							<li>
								<b>Termination</b>: <span class="cfg-t">SEPARATOR</span> →
								<span class="cfg-t">";"</span> ends the statement.
							</li>
						</ol>
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
h3 {
    color: #444;
    margin: 0.7rem 0 0.3rem 0;
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
.header-type,
.header-tokens,
.token-type,
.token-values {
    display: table-cell;
    border: 1px solid #e5e7eb;
    padding: 0.6rem 0.8rem;
    text-align: left;
}
.table-header {
    background: #e6edfa;
    color: #001a6e;
    font-weight: 600;
}
.cfg-section {
    display: flex;
    gap: 2.5rem;
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
    background: #003366; 
    color: #ebeef1;
}

:global(html.dark-mode) .nav-button:disabled {
    background: #718096; 
}

:global(html.dark-mode) .step-counter {
    color: #a0aec0;
}
</style>
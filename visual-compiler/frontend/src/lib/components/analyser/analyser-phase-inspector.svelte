<script lang="ts">
	import { slide } from 'svelte/transition';

	// Interface for a single scope rule
	interface ScopeRule {
		id: number;
		start: string;
		end: string;
	}

	// Reactive array to hold the scope delimiter rules
	let scopeRules: ScopeRule[] = [{ id: 0, start: '', end: '' }];
	let nextId = 1;

	// Function to add a new, empty row for a scope rule
	function addScopeRow() {
		scopeRules = [...scopeRules, { id: nextId++, start: '', end: '' }];
	}

	// Function to remove a scope rule by its index
	function removeScopeRow(index: number) {
		scopeRules.splice(index, 1);
		scopeRules = scopeRules; // Trigger reactivity
	}

	// Reactive statement to determine if the "Add" button should be visible
	$: showAddButton =
		scopeRules.length > 0 &&
		scopeRules[scopeRules.length - 1].start.trim() !== '' &&
		scopeRules[scopeRules.length - 1].end.trim() !== '';

    let showStartTooltip = false;
    let showEndTooltip = false;
</script>

<div class="panel-container">
	<div class="header">
        <h1 class="heading">ANALYSING</h1>
		<h2 class="heading2">Scope Rules</h2>
		<p>Define the pairs of tokens or characters that mark the beginning and end of a new scope.</p>
	</div>

    <div class="rules-list">
        {#each scopeRules as rule, i (rule.id)}
            <div class="rule-row" in:slide={{ duration: 250 }}>
                <div class="input-tooltip-wrapper">
                    <input
                        type="text"
                        bind:value={rule.start}
                        placeholder="Start Delimiter "
                        aria-label="Start Delimiter"
                        class="scope-input"
                        on:focus={() => i === 0 && (showStartTooltip = true)}
                        on:blur={() => i === 0 && (showStartTooltip = false)}
                    />
                    {#if i === 0 && showStartTooltip}
						<div class="tooltip">Examples include <code>&#123;</code></div>
                    {/if}
                </div>
                <span class="arrow-separator">→</span>
                <div class="input-tooltip-wrapper">
                    <input
                        type="text"
                        bind:value={rule.end}
                        placeholder="End Delimiter "
                        aria-label="End Delimiter"
                        class="scope-input"
                        on:focus={() => i === 0 && (showEndTooltip = true)}
                        on:blur={() => i === 0 && (showEndTooltip = false)}
                    />
                    {#if i === 0 && showEndTooltip}
                        <div class="tooltip">Examples include <code>}</code></div>
                    {/if}
                </div>
                <button
                    class="icon-button delete-button"
                    on:click={() => removeScopeRow(i)}
                    aria-label="Remove row"
                    disabled={scopeRules.length === 1}
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><polyline points="3 6 5 6 21 6" /><path
                            d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                        /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg
                    >
                </button>
            </div>
        {/each}
    </div>



	{#if showAddButton}
		<div class="add-button-container" transition:slide={{ duration: 150 }}>
			<button class="add-button" on:click={addScopeRow} aria-label="Add new scope delimiter row">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2.5"
					stroke-linecap="round"
					stroke-linejoin="round"
					><line x1="12" y1="5" x2="12" y2="19" /><line x1="5" y1="12" x2="19" y2="12" /></svg
				>
				<span>Add Delimiter</span>
			</button>
		</div>
	{/if}
</div>

<style>

    .heading2{
        color: #444;
    }
    .heading {
        font-weight: 600;
        margin-bottom: 2rem;
        text-align: center;
        margin-top: 0;
    }
	.panel-container {
		display: flex;
		flex-direction: column;
		padding: 1.5rem 1rem;
		height: 100%;
	}

	.header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0 0 0.25rem 0;
		color: #1a2a4a;
	}

	.header p {
		font-size: 0.9rem;
		color: #5a6a8a;
		margin: 0 0 1.5rem 0;
	}

	.rules-list {
		flex-grow: 1;
		overflow-y: auto;
		padding-right: 0.5rem;
	}

	.rule-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 0.75rem;
	}

	.scope-input {
		flex: 1;
		padding: 0.6rem 0.8rem;
		font-size: 0.9rem;
		border: 1px solid #d0d7e0;
		border-radius: 6px; /* Matching border radius */
		background-color: #ffffff;
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease;
	}

	.scope-input:focus {
		outline: none;
		border-color: #074799;
		box-shadow: 0 0 0 3px rgba(7, 71, 153, 0.15);
	}

	.arrow-separator {
		color: #7a8aa3;
		font-size: 1.2rem;
	}

	.icon-button {
		display: flex;
		align-items: center;
		justify-content: center;
		background: none;
		border: none;
		padding: 0.25rem;
		border-radius: 50%;
		cursor: pointer;
		color: #7a8aa3;
		transition:
			background-color 0.2s ease,
			color 0.2s ease;
	}

	.icon-button:hover:not(:disabled) {
		background-color: #eef2f7;
		color: #d9534f;
	}
	
	.icon-button:disabled {
		color: #c0c7d3;
		cursor: not-allowed;
	}

	.add-button-container {
		display: flex;
		justify-content: center;
		padding-top: 1rem;
		margin-top: auto;
	}

	.add-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		background-color: #eef2f7;
		color: #001a6e;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s ease;
	}

	.add-button:hover {
		background-color: #e1e8f0;
	}

	/* --- DARK MODE STYLES --- */
	:global(html.dark-mode) .header h2 {
		color: #e0e8f0;
	}
	:global(html.dark-mode) .header p {
		color: #8a9bb3;
	}
	:global(html.dark-mode) .scope-input {
		background-color: #253352;
		border-color: #3a4a6a;
		color: #e0e8f0;
	}
	:global(html.dark-mode) .scope-input::placeholder {
		color: #6a7b99;
	}
	:global(html.dark-mode) .scope-input:focus {
		border-color: #3a7bd5;
		box-shadow: 0 0 0 3px rgba(58, 123, 213, 0.2);
	}
	:global(html.dark-mode) .arrow-separator {
		color: #8a9bb3;
	}
	:global(html.dark-mode) .icon-button {
		color: #8a9bb3;
	}
	:global(html.dark-mode) .icon-button:hover:not(:disabled) {
		background-color: #2a4a8a;
		color: #ff8a8a;
	}
	:global(html.dark-mode) .icon-button:disabled {
		color: #5a6a8a;
	}
	:global(html.dark-mode) .add-button {
		background-color: #2a4a8a;
		color: #e0e8f0;
	}
	:global(html.dark-mode) .add-button:hover {
		background-color: #3a5a9a;
	}

     .input-tooltip-wrapper {
        position: relative;
        display: flex;
        align-items: center;
    }

    .tooltip {
        position: absolute;
        top: 110%;
        left: 0;
        z-index: 10;
        background: #fff;
        color: #1a2a4a;
        border: 1px solid #d0d7e0;
        border-radius: 6px;
        padding: 0.5rem 0.75rem;
        font-size: 0.92rem;
        box-shadow: 0 2px 8px rgba(30,40,80,0.08);
        white-space: nowrap;
        margin-top: 0.2rem;
        transition: opacity 0.15s;
    }

    .tooltip code {
        background: #eef2f7;
        padding: 0.1em 0.3em; 
        border-radius: 4px;
        font-size: 0.92em;
    }

    :global(html.dark-mode) .tooltip {
        background: #253352;
        color: #e0e8f0;
        border: 1px solid #3a4a6a;
        box-shadow: 0 2px 8px rgba(30,40,80,0.18);
    }

    :global(html.dark-mode) .tooltip code {
        background: #2a4a8a;
        color: #e0e8f0;
    }
</style>
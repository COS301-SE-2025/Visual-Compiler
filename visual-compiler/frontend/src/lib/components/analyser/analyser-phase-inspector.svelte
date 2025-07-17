<script lang="ts">
	import { slide, fade } from 'svelte/transition';
	import { AddToast } from '$lib/stores/toast';
    import { createEventDispatcher } from 'svelte';
    
    const dispatch = createEventDispatcher();
	// --- STATE ---
	interface ScopeRule {
		id: number;
		start: string;
		end: string;
	}

	let scopeRules: ScopeRule[] = [{ id: 0, start: '', end: '' }];
	let nextId = 1;
	let showStartTooltip = false;
	let showEndTooltip = false;
	let rulesSubmitted = false;
	let submittedRules: ScopeRule[] = [];

	// --- LOGIC ---
	function addScopeRow() {
		scopeRules = [...scopeRules, { id: nextId++, start: '', end: '' }];
		rulesSubmitted = false;
	}

	function removeScopeRow(index: number) {
		scopeRules.splice(index, 1);
		scopeRules = scopeRules;
		rulesSubmitted = false;
	}

	function handleRuleInput() {
		rulesSubmitted = false;
	}

	function handleSubmit() {
		if (scopeRules.some((rule) => rule.start.trim() === '' || rule.end.trim() === '')) {
			AddToast('Please fill out all delimiter fields before submitting.', 'error');
			return;
		}
		submittedRules = JSON.parse(JSON.stringify(scopeRules));
		rulesSubmitted = true;
		AddToast('Scope rules submitted successfully!', 'success');
		console.log('Submitted Scope Rules:', submittedRules);
	}

	function resetState() {
		scopeRules = [{ id: 0, start: '', end: '' }];
		nextId = 1;
		submittedRules = [];
		rulesSubmitted = false;
        dispatch('reset');
		AddToast('Scope rules have been reset.', 'info');
	}

    function handleGenerate() {
        dispatch('generate');
        AddToast('Symbol table generated', 'success');
    }
	// --- REACTIVE STATEMENTS ---
	$: lastRowComplete =
		scopeRules.length > 0 &&
		scopeRules[scopeRules.length - 1].start.trim() !== '' &&
		scopeRules[scopeRules.length - 1].end.trim() !== '';

	$: allRowsComplete = scopeRules.every(
		(rule) => rule.start.trim() !== '' && rule.end.trim() !== ''
	);


</script>

<div class="panel-container">
	<div class="header">
		{#if rulesSubmitted}
			<div class="reset-wrapper" transition:fade={{ duration: 150 }}>
                <button class="reset-button" on:click={resetState} aria-label="Reset rules">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
					<path d="M3 3v5h5" />
				</svg>
			</button>
			</div>
		{/if}

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
						placeholder="Start Delimiter"
						aria-label="Start Delimiter"
						class="scope-input"
						on:focus={() => i === 0 && (showStartTooltip = true)}
						on:blur={() => i === 0 && (showStartTooltip = false)}
						on:input={handleRuleInput}
					/>
					{#if i === 0 && showStartTooltip}
						<div class="tooltip">Examples include <code>{'{'}</code></div>
					{/if}
				</div>

				<span class="arrow-separator">→</span>

				<div class="input-tooltip-wrapper">
					<input
						type="text"
						bind:value={rule.end}
						placeholder="End Delimiter"
						aria-label="End Delimiter"
						class="scope-input"
						on:focus={() => i === 0 && (showEndTooltip = true)}
						on:blur={() => i === 0 && (showEndTooltip = false)}
						on:input={handleRuleInput}
					/>
					{#if i === 0 && showEndTooltip}
						<div class="tooltip">Examples include <code>{'}'}</code></div>
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

	<div class="actions-container">
		{#if !rulesSubmitted}
			{#if lastRowComplete}
				<div class="add-button-wrapper" transition:slide={{ duration: 150 }}>
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
			<button class="submit-button" on:click={handleSubmit} disabled={!allRowsComplete}>
				Submit Rules
			</button>
		{:else}
			<div class="generate-wrapper" transition:slide|local={{ duration: 250 }}>
				<button class="generate-button" on:click={handleGenerate}> Generate Symbol Table </button>
			</div>
		{/if}
	</div>
</div>

<style>
	.heading2 {
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
	.header {
		position: relative; 
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
		margin-right: -0.5rem;
	}
	.rule-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 0.75rem;
	}
	.input-tooltip-wrapper {
		position: relative;
		flex: 1;
		display: flex;
	}
	.scope-input {
		flex: 1;
		padding: 0.6rem 0.8rem;
		font-size: 0.9rem;
		border: 1px solid #d0d7e0;
		border-radius: 6px;
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
		font-weight: 600;
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
	}
	.delete-button:hover:not(:disabled) {
		color: #d9534f;
	}
	.icon-button:disabled {
		color: #c0c7d3;
		cursor: not-allowed;
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
		font-size: 0.85rem;
		box-shadow: 0 2px 8px rgba(30, 40, 80, 0.08);
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
	.actions-container {
		padding-top: 1rem;
		margin-top: auto;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
	.add-button-wrapper,
	.generate-wrapper {
		display: flex;
		justify-content: center;
	}
	.add-button {
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
	.add-button:hover {
		background-color: #e1e8f0;
		border-color: #001a6e;
	}
	.submit-button,
	.generate-button {
		width: 100%;
		padding: 0.7rem 1rem;
		font-size: 0.9rem;
		font-weight: 600;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition:
			background-color 0.2s ease,
			opacity 0.2s ease;
	}
	.submit-button {
		background-color: #001a6e;
		color: white;
	}
	.submit-button:hover:not(:disabled) {
		background-color: #074799;
	}
	.submit-button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	.generate-button {
		background-color: #10b981;
		color: white;
	}
	.generate-button:hover {
		background-color: #059669;
	}

	
	.reset-wrapper {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    z-index: 10;
}

.reset-button {
    background: none;
    border: none;
    padding: 0;
    margin: 0;
    cursor: pointer;
    color: #7a8aa3;
    transition: color 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
}

.reset-button:hover {
    color: #001a6e;
}

/* Dark mode styles */
:global(html.dark-mode) .reset-button {
    color: #8a9bb3;
}

:global(html.dark-mode) .reset-button:hover {
    color: #e0e8f0;
}

	/* --- DARK MODE STYLES --- */
	:global(html.dark-mode) .header h2,
	:global(html.dark-mode) .heading,
	:global(html.dark-mode) .heading2 {
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
	}
	:global(html.dark-mode) .delete-button:hover:not(:disabled) {
		color: #ff8a8a;
	}
	:global(html.dark-mode) .icon-button:disabled {
		color: #5a6a8a;
	}
	:global(html.dark-mode) .tooltip {
		background: #253352;
		color: #e0e8f0;
		border: 1px solid #3a4a6a;
		box-shadow: 0 2px 8px rgba(30, 40, 80, 0.18);
	}
	:global(html.dark-mode) .tooltip code {
		background: #2a4a8a;
		color: #e0e8f0;
	}
	:global(html.dark-mode) .add-button {
		background-color: transparent;
		color: #9ab0d3;
		border-color: #3a4a6a;
	}
	:global(html.dark-mode) .add-button:hover {
		background-color: #2a4a8a;
		color: #e0e8f0;
		border-color: #3a7bd5;
	}
	:global(html.dark-mode) .submit-button {
		background-color: #3a7bd5;
	}
	:global(html.dark-mode) .submit-button:hover:not(:disabled) {
		background-color: #4a8ce5;
	}
	:global(html.dark-mode) .generate-button {
		background-color: #10b981;
	}
	:global(html.dark-mode) .generate-button:hover {
		background-color: #059669;
	}
	:global(html.dark-mode) .reset-button:hover:not(:disabled) {
		color: #4a8ce5;
	}
</style>
<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { showCanvasTutorial } from '$lib/stores/tutorial';

	const dispatch = createEventDispatcher();

	let active_question_index = -1;

	const faqs = [
		{
			question: 'How do I add a new node to the canvas?',
			answer:
				'Click on a block in the "Blocks" panel on the left. This will add a new node of that type to your workspace.'
		},
		{
			question: 'How do I move a node?',
			answer:
				'Click a node to select it (it will be highlighted), then drag it to a new position.'
		},
		{
			question: 'How do I connect two nodes?',
			answer:
				'Click and hold on an output anchor (the grey circle on the right of a node) and drag to an input anchor (the circle on the left of another node).'
		},
		{
			question: 'How do I delete a connection?',
			answer:
				'Click and drag one of its connection points away from its anchor and release it onto the empty canvas.'
		},
		{
			question: 'How do I import code from a saved project?',
			answer:
				'Use the "Import from saved project" dropdown in the code input area.'
		},
		{
			question: 'How do I upload code from a file?',
			answer:
				'Click the "Upload File" button in the code input area and select a .txt file from your computer.'
		},
		{
			question: 'How do I use the AI Assistant?',
			answer:
				'Click the purple AI button in the bottom-right corner. You can ask questions about compiler concepts or generate input for any phase.'
		},
		{
			question: 'Why can’t I enter a node or phase?',
			answer:
				'You must complete the previous phases of the compiler before you can enter the next phase.'
		},
		{
			question: 'How do I get a code example?',
			answer:
				'Click the "Show Example" button in the code input area to insert a sample code snippet.'
		},
	];

	function toggleQuestion(index: number) {
		if (active_question_index === index) {
			active_question_index = -1;
		} else {
			active_question_index = index;
		}
	}

	function startTutorial() {
		showCanvasTutorial();
		dispatch('close');
	}
</script>

<div class="help-popup" role="dialog" aria-label="Help Menu">
	<div class="help-header">
		<h2>Help & FAQ</h2>
		<button class="close-btn" on:click={() => dispatch('close')}>✕</button>
	</div>
	<div class="faq-list">
		{#each faqs as faq, index}
			<div class="faq-item">
				<button class="faq-question" on:click={() => toggleQuestion(index)}>
					<span>{faq.question}</span>
					<span class="arrow">{active_question_index === index ? '▲' : '▼'}</span>
				</button>
				{#if active_question_index === index}
					<div class="faq-answer">
						<p>{faq.answer}</p>
					</div>
				{/if}
			</div>
		{/each}
	</div>
	
	<div class="tutorial-section">
		<button class="tutorial-btn" on:click={startTutorial}>
			<span>Take Canvas Tutorial</span>
		</button>
	</div>
</div>

<style>
	.help-popup {
		position: absolute;
		top: calc(100% + 12px);
		right: 0;
		width: 380px;
		max-height: 70vh;
		background: #ffffff;
		border-radius: 8px;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
		border: 1px solid #e0e0e0;
		z-index: 4000;
		display: flex;
		flex-direction: column;
	}

	.help-popup::before {
		content: '';
		position: absolute;
		top: -10px;
		right: 15px;
		width: 0;
		height: 0;
		border-left: 10px solid transparent;
		border-right: 10px solid transparent;
		border-bottom: 10px solid #ffffff;
		filter: drop-shadow(0 -1px 1px rgba(0, 0, 0, 0.05));
	}

	.help-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		border-bottom: 1px solid #e5e5e5;
	}

	.help-header h2 {
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		color: #333;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 1.25rem;
		cursor: pointer;
		color: #888;
		padding: 0;
		line-height: 1;
	}
	.close-btn:hover {
		color: #333;
	}

	.faq-list {
		overflow-y: auto;
		padding: 0.5rem;
	}

	.faq-item {
		border-bottom: 1px solid #e5e5e5;
	}
	.faq-item:last-child {
		border-bottom: none;
	}

	.faq-question {
		width: 100%;
		background: none;
		border: none;
		text-align: left;
		padding: 1rem 1.25rem;
		font-size: 0.95rem;
		cursor: pointer;
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-weight: 500;
		color: #333;
		border-radius: 4px;
	}
	.faq-question:hover {
		background-color: #f9f9f9;
	}

	.faq-answer {
		padding: 0 1.25rem 1rem 1.25rem;
		color: #555;
		line-height: 1.6;
		font-size: 0.9rem;
	}
	.faq-answer p {
		margin: 0;
	}

	.arrow {
		font-size: 0.8rem;
	}

	/* Tutorial Section Styles */
	.tutorial-section {
		padding: 1rem 1.5rem;
		border-top: 1px solid #e0e0e0;
		background: linear-gradient(135deg, #f8fafc, #f1f5f9);
	}

	.tutorial-btn {
		width: 100%;
		padding: 0.75rem 1rem;
		background: linear-gradient(135deg, #041a47, #0c2d5e);
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 0.95rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		box-shadow: 0 2px 4px rgba(4, 26, 71, 0.3);
	}

	.tutorial-btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 8px rgba(4, 26, 71, 0.4);
		background: linear-gradient(135deg, #052759, #0e3a75);
	}

	:global(html.dark-mode) .tutorial-btn {
		background: linear-gradient(135deg, #001A6E, #002a8e);
		box-shadow: 0 2px 4px rgba(0, 26, 110, 0.4);
	}

	:global(html.dark-mode) .tutorial-btn:hover {
		background: linear-gradient(135deg, #002a8e, #003bb3);
		box-shadow: 0 4px 8px rgba(0, 26, 110, 0.5);
	}

	:global(html.dark-mode) .help-popup {
		background: #1a2a4a;
		color: #f0f0f0;
		border-color: #2a3a5a;
	}
	:global(html.dark-mode) .help-popup::before {
		border-bottom-color: #1a2a4a;
		filter: drop-shadow(0 -1px 1px rgba(0, 0, 0, 0.2));
	}
	:global(html.dark-mode) .help-header h2,
	:global(html.dark-mode) .faq-question,
	:global(html.dark-mode) .close-btn {
		color: #f0f0f0;
	}
	:global(html.dark-mode) .close-btn:hover {
		color: #bbb;
	}
	:global(html.dark-mode) .help-header {
		border-bottom: 1px solid #2a3a5a;
	}
	:global(html.dark-mode) .faq-item {
		border-bottom: 1px solid #2a3a5a;
	}
	:global(html.dark-mode) .faq-question:hover {
		background-color: #203557;
	}
	:global(html.dark-mode) .faq-answer {
		color: #ccc;
	}
	
	:global(html.dark-mode) .tutorial-section {
		border-top-color: #2a3a5a;
		background: #1a2a4a ;
	}


</style>

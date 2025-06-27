<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	let active_question_index = -1;

	const faqs = [
		{
			question: 'How do I add a new node to the canvas?',
			answer:
				'Simply click on one of the blocks in the "Blocks" panel on the left. This will add a new node of that type to the workspace.'
		},
		{
			question: 'How do I move a node?',
			answer:
				'For the smoothest experience, click a node once to select it (it will be highlighted), and then click and drag it to a new position.'
		},
		{
			question: 'How do I connect two nodes?',
			answer:
				'Click and hold on an output anchor (the grey circle on the right side of a node) and drag the line to an input anchor (the circle on the left side of another node).'
		},
		{
			question: 'How do I delete a node or connection?',
			answer:
				'To delete a node, right-click on it to open its configuration modal, where you will find a delete option. To delete a connection, click and drag one of its connection points away from its anchor and release it onto the empty canvas.'
		},
		{
			question: 'Why can’t I configure a Lexer or Parser node?',
			answer:
				'You must first add and submit code using a "Source Code" node before you can configure the subsequent phases of the compiler.'
		}
	];

	function toggleQuestion(index: number) {
		if (active_question_index === index) {
			active_question_index = -1;
		} else {
			active_question_index = index;
		}
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
</style>

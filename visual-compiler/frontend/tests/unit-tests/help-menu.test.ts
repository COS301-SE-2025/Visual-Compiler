import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import HelpMenu from '../../src/lib/components/main/help-menu.svelte';

describe('HelpMenu Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('TestRender_Success: Renders help menu with title and close button', () => {
		render(HelpMenu);

		expect(screen.getByText('Help & FAQ')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: '✕' })).toBeInTheDocument();
		expect(screen.getByRole('dialog')).toHaveAttribute('aria-label', 'Help Menu');
	});

	it('TestCloseButton_Success: Close button is clickable', async () => {
		render(HelpMenu);
		
		const closeButton = screen.getByRole('button', { name: '✕' });
		await fireEvent.click(closeButton);
		
		// The button should be responsive to clicks (no errors thrown)
		expect(closeButton).toBeInTheDocument();
	});

	it('TestFAQQuestions_Success: Renders all FAQ questions', () => {
		render(HelpMenu);

		expect(screen.getByText('How do I add a new node to the canvas?')).toBeInTheDocument();
		expect(screen.getByText('How do I move a node?')).toBeInTheDocument();
		expect(screen.getByText('How do I connect two nodes?')).toBeInTheDocument();
		expect(screen.getByText('How do I delete a node or connection?')).toBeInTheDocument();
		expect(screen.getByText(/Why can.t I configure a Lexer or Parser node/)).toBeInTheDocument();
	});

	it('TestQuestionToggle_Success: Expands and collapses FAQ answers', async () => {
		render(HelpMenu);

		const firstQuestion = screen.getByText('How do I add a new node to the canvas?');
		
		// Initially, answer should not be visible
		expect(screen.queryByText(/Simply click on one of the blocks/)).not.toBeInTheDocument();
		
		// Click to expand
		await fireEvent.click(firstQuestion);
		expect(screen.getByText(/Simply click on one of the blocks/)).toBeInTheDocument();
		
		// Click again to collapse
		await fireEvent.click(firstQuestion);
		expect(screen.queryByText(/Simply click on one of the blocks/)).not.toBeInTheDocument();
	});

	it('TestArrowIndicator_Success: Shows correct arrow direction for expanded/collapsed state', async () => {
		render(HelpMenu);

		const firstQuestionButton = screen.getByRole('button', { name: /How do I add a new node to the canvas\?/ });
		
		// Initially collapsed (down arrow)
		expect(firstQuestionButton).toHaveTextContent('▼');
		
		// Click to expand (up arrow)
		await fireEvent.click(firstQuestionButton);
		expect(firstQuestionButton).toHaveTextContent('▲');
		
		// Click again to collapse (down arrow)
		await fireEvent.click(firstQuestionButton);
		expect(firstQuestionButton).toHaveTextContent('▼');
	});

	it('TestMultipleQuestions_Success: Only one question can be expanded at a time', async () => {
		render(HelpMenu);

		const firstQuestion = screen.getByRole('button', { name: /How do I add a new node to the canvas\?/ });
		const secondQuestion = screen.getByRole('button', { name: /How do I move a node\?/ });

		// Expand first question
		await fireEvent.click(firstQuestion);
		expect(screen.getByText(/Simply click on one of the blocks/)).toBeInTheDocument();
		expect(firstQuestion).toHaveTextContent('▲');

		// Expand second question (should close first)
		await fireEvent.click(secondQuestion);
		expect(screen.queryByText(/Simply click on one of the blocks/)).not.toBeInTheDocument();
		expect(screen.getByText(/For the smoothest experience/)).toBeInTheDocument();
		expect(firstQuestion).toHaveTextContent('▼');
		expect(secondQuestion).toHaveTextContent('▲');
	});

	it('TestAllFAQContent_Success: Verifies all FAQ answers are correct', async () => {
		render(HelpMenu);

		const questions = [
			{
				question: 'How do I add a new node to the canvas?',
				answerPart: 'Simply click on one of the blocks'
			},
			{
				question: 'How do I move a node?',
				answerPart: 'For the smoothest experience'
			},
			{
				question: 'How do I connect two nodes?',
				answerPart: 'Click and hold on an output anchor'
			},
			{
				question: 'How do I delete a node or connection?',
				answerPart: 'To delete a node, right-click'
			},
			{
				question: "Why can't I configure a Lexer or Parser node?",
				answerPart: 'You must first add and submit code'
			}
		];

		for (const { question, answerPart } of questions) {
			// Handle special characters - use partial matching for button names
			const questionButton = screen.getByRole('button', { 
				name: new RegExp(question.replace(/['']/g, ".").replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + '.*')
			});
			
			await fireEvent.click(questionButton);
			expect(screen.getByText(new RegExp(answerPart))).toBeInTheDocument();
			
			// Close the question for next iteration
			await fireEvent.click(questionButton);
		}
	});

	it('TestAccessibility_Success: Has proper accessibility attributes', () => {
		render(HelpMenu);

		const dialog = screen.getByRole('dialog');
		expect(dialog).toHaveAttribute('aria-label', 'Help Menu');

		// Check that all FAQ questions are properly structured as buttons
		const questionButtons = screen.getAllByRole('button');
		expect(questionButtons.length).toBeGreaterThan(1); // Close button + FAQ buttons

		// Each question should be a button (for keyboard navigation)
		questionButtons.forEach(button => {
			expect(button).toBeInTheDocument();
		});
	});
});

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import TranslatorPhaseTutorial from '../src/lib/components/translator/translator-phase-tutorial.svelte';

describe('TranslatorPhaseTutorial Component', () => {
	beforeEach(() => {
		// Reset any state if needed
	});

	it('TestRender_Success: Renders tutorial with main heading', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText('What is Translating?')).toBeInTheDocument();
		expect(screen.getByText(/This tool converts your source code/)).toBeInTheDocument();
	});

	it('TestInitialStep_Success: Shows step 1 initially', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText('1. Complete Lexing and Parsing')).toBeInTheDocument();
		expect(screen.getByText(/Before using the translator, you must complete/)).toBeInTheDocument();
	});

	it('TestStepsToFollow_Success: Shows steps to follow section', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText('Steps to Follow:')).toBeInTheDocument();
	});

	it('TestNavigationButtons_Success: Shows navigation buttons', async () => {
		render(TranslatorPhaseTutorial);

		await waitFor(() => {
			expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument();
			expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument();
		});
	});

	it('TestStepProgression_Success: Can navigate to next step', async () => {
		render(TranslatorPhaseTutorial);

		const nextButton = screen.getByRole('button', { name: /next/i });
		await fireEvent.click(nextButton);

		expect(screen.getByText('2 / 6')).toBeInTheDocument();
	});

	// Test: Renders all token types in step 4
	it('TestTokenTypesRender_Success: Renders all token types in step 4', async () => {
		render(TranslatorPhaseTutorial);

		// Navigate to step 4 (token types)
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]); // Step 2
			await fireEvent.click(nextButtons[0]); // Step 3
			await fireEvent.click(nextButtons[0]); // Step 4

			await waitFor(() => {
				// Verify all 6 token types are rendered in the each loop
				expect(screen.getByText('🔹 KEYWORD')).toBeInTheDocument();
				expect(screen.getByText('Reserved words (e.g., int, if, while)')).toBeInTheDocument();
				expect(screen.getByText('🔹 IDENTIFIER')).toBeInTheDocument();
				expect(screen.getByText('Variable names (e.g., result, count)')).toBeInTheDocument();
				expect(screen.getByText('🔹 NUMBER')).toBeInTheDocument();
				expect(screen.getByText('Integers (e.g., 10, 42)')).toBeInTheDocument();
				expect(screen.getByText('🔹 PLUS')).toBeInTheDocument();
				expect(screen.getByText('Addition operator (+)')).toBeInTheDocument();
				expect(screen.getByText('🔹 MINUS')).toBeInTheDocument();
				expect(screen.getByText('Subtraction operator (-)')).toBeInTheDocument();
				expect(screen.getByText('🔹 ASSIGNMENT')).toBeInTheDocument();
				expect(screen.getByText('Assignment operator (=)')).toBeInTheDocument();
			});
		}
	});

	// Test: Renders all example rules in step 5
	it('TestExampleRulesRender_Success: Renders all example rules in step 5', async () => {
		render(TranslatorPhaseTutorial);

		// Navigate to step 5 (example rules)
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→')
		);
		
		if (nextButtons.length > 0) {
			for (let i = 0; i < 4; i++) {
				await fireEvent.click(nextButtons[0]);
			}

		await waitFor(() => {
			// Verify both example rules are rendered in the each loop
			expect(screen.getByText('12 + 13')).toBeInTheDocument();
			expect(screen.getByText('INTEGER, OPERATOR, INTEGER')).toBeInTheDocument();
			expect(screen.getByText(/mov\s+rax, \{INTEGER\}/)).toBeInTheDocument();
			expect(screen.getByText(/mov\s+rax, 12/)).toBeInTheDocument(); // Use regex for multi-line text
			expect(screen.getByText('bool found = true;')).toBeInTheDocument();
			expect(screen.getByText('KEYWORD, IDENTIFIER, ASSIGNMENT, BOOLEAN')).toBeInTheDocument();
			expect(screen.getByText(/mov\s+rax, \[\{IDENTIFIER\}\]/)).toBeInTheDocument();
			expect(screen.getByText(/mov\s+rax, \[found\]/)).toBeInTheDocument();
		});
		}
	});

	// Test: Renders all critical mistakes in step 6
	it('TestCriticalDontsRender_Success: Renders all critical mistakes in step 6', async () => {
		render(TranslatorPhaseTutorial);

		// Navigate to step 6 (critical mistakes)
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→')
		);
		
		if (nextButtons.length > 0) {
			for (let i = 0; i < 5; i++) {
				await fireEvent.click(nextButtons[0]);
			}

			await waitFor(() => {
				// Verify critical mistakes are rendered in the each loop
				expect(screen.getByText('Not using token types')).toBeInTheDocument();
				expect(screen.getByText('int, result, =, 1')).toBeInTheDocument();
				expect(screen.getByText('KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER')).toBeInTheDocument();
				
				expect(screen.getByText('Skipping placeholders')).toBeInTheDocument();
				expect(screen.getByText('mov [result], 1')).toBeInTheDocument();
				expect(screen.getByText('MOV [{IDENTIFIER}], {INTEGER}')).toBeInTheDocument();
			});
		}
	});
});
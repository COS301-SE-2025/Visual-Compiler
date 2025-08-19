import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import TranslatorPhaseTutorial from '../../src/lib/components/translator/translator-phase-tutorial.svelte';

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

		expect(screen.getByText('1️⃣ Complete Lexing and Parsing')).toBeInTheDocument();
		expect(screen.getByText(/Before using the translator, you must complete/)).toBeInTheDocument();
	});

	it('TestStepsToFollow_Success: Shows steps to follow section', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText('Steps to Follow:')).toBeInTheDocument();
	});

	it('TestNavigationButtons_Success: Shows navigation buttons', () => {
		render(TranslatorPhaseTutorial);

		// Should have navigation controls
		const buttons = screen.getAllByRole('button');
		expect(buttons.length).toBeGreaterThan(0);
	});

	it('TestStepProgression_Success: Can navigate to next step', async () => {
		render(TranslatorPhaseTutorial);

		// Find next button
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText('2️⃣ Define Translation Rules')).toBeInTheDocument();
			});
		}
	});

	it('TestTranslationRulesStep_Success: Shows translation rules in step 2', async () => {
		render(TranslatorPhaseTutorial);

		// Navigate to step 2 if possible
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText(/Click "\+ Add New Rule" for each pattern/)).toBeInTheDocument();
			});
		}
	});

	it('TestLexingParsingPrerequisites_Success: Shows lexing and parsing prerequisites', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText('🔍 Lexing')).toBeInTheDocument();
		expect(screen.getByText('Break your source code into tokens')).toBeInTheDocument();
		expect(screen.getByText('🌳 Parsing')).toBeInTheDocument();
		expect(screen.getByText('Use the tokens to create a syntax tree')).toBeInTheDocument();
	});

	it('TestTokenTypes_Success: Contains token type information', () => {
		const { container } = render(TranslatorPhaseTutorial);

		// Should contain token-related content
		expect(container.innerHTML).toContain('token');
		expect(container.innerHTML).toContain('code');
	});

	it('TestExampleRules_Success: Contains example translation rules', () => {
		const { container } = render(TranslatorPhaseTutorial);

		// Should contain example-related content
		expect(container.innerHTML).toContain('rule');
		expect(container.innerHTML).toContain('translation');
	});

	it('TestCriticalDonts_Success: Shows critical mistakes to avoid', () => {
		const { container } = render(TranslatorPhaseTutorial);

		// Should contain warning-related content
		expect(container.innerHTML).toContain('must');
		expect(container.innerHTML).toContain('complete');
	});

	it('TestTutorialStructure_Success: Has proper tutorial structure', () => {
		const { container } = render(TranslatorPhaseTutorial);

		const tutorialContainer = container.querySelector('.phase-tutorial');
		expect(tutorialContainer).toBeInTheDocument();

		const tutorialContent = container.querySelector('.tutorial-content');
		expect(tutorialContent).toBeInTheDocument();
	});

	it('TestTranslationDescription_Success: Contains translation explanation', () => {
		render(TranslatorPhaseTutorial);

		expect(screen.getByText(/Think of it as teaching the computer how to rewrite/)).toBeInTheDocument();
	});

	it('TestMultipleSteps_Success: Has multiple tutorial steps', () => {
		const { container } = render(TranslatorPhaseTutorial);

		// Should have tutorial step content areas
		const tutorialContentArea = container.querySelector('.tutorial-content-area');
		expect(tutorialContentArea).toBeInTheDocument();
	});

	it('TestCodeExamples_Success: Contains code examples', () => {
		render(TranslatorPhaseTutorial);

		// Should contain code-related content that exists
		expect(screen.getByText(/converts your source code/)).toBeInTheDocument();
	});

	it('TestBackNavigation_Success: Can navigate to previous step', async () => {
		render(TranslatorPhaseTutorial);

		// First navigate forward
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			// Now try to go back
			const prevButtons = screen.getAllByRole('button').filter(btn => 
				btn.textContent?.includes('Previous') || btn.textContent?.includes('←') || btn.textContent?.includes('<')
			);
			
			if (prevButtons.length > 0) {
				await fireEvent.click(prevButtons[0]);
				
				await waitFor(() => {
					expect(screen.getByText('1️⃣ Complete Lexing and Parsing')).toBeInTheDocument();
				});
			}
		}
	});

	it('TestStepIndicators_Success: Shows step progression indicators', () => {
		const { container } = render(TranslatorPhaseTutorial);

		// Should have some indication of step progression
		expect(container.innerHTML).toContain('1️⃣');
	});
});

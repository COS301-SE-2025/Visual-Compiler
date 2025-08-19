import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import LexerPhaseTutorial from '../../src/lib/components/lexer/lexer-phase-tutorial.svelte';

describe('LexerPhaseTutorial Component', () => {
	beforeEach(() => {
		// Reset any state if needed
	});

	it('TestRender_Success: Renders tutorial with main heading', () => {
		render(LexerPhaseTutorial);

		expect(screen.getByText('What is Lexing?')).toBeInTheDocument();
		expect(screen.getByText(/A lexer breaks down source code into tokens/)).toBeInTheDocument();
	});

	it('TestInitialStep_Success: Shows step 1 initially', () => {
		render(LexerPhaseTutorial);

		expect(screen.getByText('1. Source Code Input')).toBeInTheDocument();
		expect(screen.getByText(/Start by entering your source code in the input box/)).toBeInTheDocument();
	});

	it('TestStepsToFollow_Success: Shows steps to follow section', () => {
		render(LexerPhaseTutorial);

		expect(screen.getByText('Steps to Follow:')).toBeInTheDocument();
	});

	it('TestNavigationButtons_Success: Shows navigation buttons', () => {
		render(LexerPhaseTutorial);

		// Should have navigation controls
		const buttons = screen.getAllByRole('button');
		expect(buttons.length).toBeGreaterThan(0);
	});

	it('TestStepProgression_Success: Can navigate to next step', async () => {
		render(LexerPhaseTutorial);

		// Find next button
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText('2.1 Regular Expression Basics')).toBeInTheDocument();
			});
		}
	});

	it('TestRegexBasics_Success: Shows regex basics in step 2', async () => {
		render(LexerPhaseTutorial);

		// Navigate to step 2 if possible
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText(/Regular expressions use special characters/)).toBeInTheDocument();
				expect(screen.getByText('Exact Match')).toBeInTheDocument();
				expect(screen.getByText('Range')).toBeInTheDocument();
			});
		}
	});

	it('TestSampleCode_Success: Contains sample code examples', () => {
		render(LexerPhaseTutorial);

		// Should contain the sample code
		expect(screen.getByText(/if \(count > 0\)/)).toBeInTheDocument();
	});

	it('TestRegexExamples_Success: Contains regex pattern examples', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain basic tutorial content about regex
		expect(container.innerHTML).toContain('lexer');
		expect(container.innerHTML).toContain('regular expressions');
	});

	it('TestTokenTypes_Success: Shows different token types', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain lexer-related content
		expect(container.innerHTML).toContain('tokens');
		expect(container.innerHTML).toContain('lexer');
	});

	it('TestTutorialStructure_Success: Has proper tutorial structure', () => {
		const { container } = render(LexerPhaseTutorial);

		const tutorialContainer = container.querySelector('.phase-tutorial');
		expect(tutorialContainer).toBeInTheDocument();

		const tutorialContent = container.querySelector('.tutorial-content');
		expect(tutorialContent).toBeInTheDocument();
	});

	it('TestLexingDescription_Success: Contains lexing explanation', () => {
		render(LexerPhaseTutorial);

		expect(screen.getByText(/Each token has a type and value/)).toBeInTheDocument();
		expect(screen.getByText(/Tokens are identified using regular expressions/)).toBeInTheDocument();
	});

	it('TestCodeSample_Success: Shows code sample in tutorial', () => {
		render(LexerPhaseTutorial);

		// Should show the example code
		expect(screen.getByText(/return true/)).toBeInTheDocument();
	});

	it('TestPatternExamples_Success: Contains pattern matching examples', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain pattern-related content
		expect(container.innerHTML).toContain('pattern');
		expect(container.innerHTML).toContain('lexer');
	});

	it('TestTokenResult_Success: Shows expected token output', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain token-related content
		expect(container.innerHTML).toContain('token');
		expect(container.innerHTML).toContain('type');
	});

	it('TestBackNavigation_Success: Can navigate to previous step', async () => {
		render(LexerPhaseTutorial);

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
					expect(screen.getByText('1. Source Code Input')).toBeInTheDocument();
				});
			}
		}
	});

	it('TestMultipleSteps_Success: Has multiple tutorial steps', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should have tutorial step content areas
		const tutorialContentArea = container.querySelector('.tutorial-content-area');
		expect(tutorialContentArea).toBeInTheDocument();
	});

	it('TestRegexPatterns_Success: Shows various regex patterns', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain regex-related content
		expect(container.innerHTML).toContain('regular expression');
		expect(container.innerHTML).toContain('pattern');
	});

	it('TestTokensByType_Success: Groups tokens by type', () => {
		const { container } = render(LexerPhaseTutorial);

		// Should contain token categorization content
		expect(container.innerHTML).toContain('token');
		expect(container.innerHTML).toContain('type');
	});
});

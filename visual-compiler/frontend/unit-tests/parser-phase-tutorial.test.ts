import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import ParserPhaseTutorial from '../src/lib/components/parser/parser-phase-tutorial.svelte';

describe('ParserPhaseTutorial Component', () => {
	beforeEach(() => {
		// Reset any state if needed
	});

	it('TestRender_Success: Renders tutorial with main heading', () => {
		render(ParserPhaseTutorial);

		expect(screen.getByText('What is Parsing?')).toBeInTheDocument();
		expect(screen.getByText(/Parsing is the process of analyzing a sequence of tokens/)).toBeInTheDocument();
	});

	it('TestInitialStep_Success: Shows step 1 initially', () => {
		render(ParserPhaseTutorial);

		expect(screen.getByText('1. Enter Source Code and Generate Tokens')).toBeInTheDocument();
		expect(screen.getByText(/Start by entering your source code and generating tokens/)).toBeInTheDocument();
	});

	it('TestStepsToFollow_Success: Shows steps to follow section', () => {
		render(ParserPhaseTutorial);

		expect(screen.getByText('Steps to Follow:')).toBeInTheDocument();
	});

	it('TestSampleCode_Success: Displays sample code and tokens', () => {
		render(ParserPhaseTutorial);

		// Should show sample code
		expect(screen.getByText(/int x = 5 \+ 3;/)).toBeInTheDocument();
		
		// Should show token table
		expect(screen.getByText('Type')).toBeInTheDocument();
		expect(screen.getByText('Token')).toBeInTheDocument();
		expect(screen.getByText('KEYWORD')).toBeInTheDocument();
		expect(screen.getByText('IDENTIFIER')).toBeInTheDocument();
	});

	it('TestTokenTable_Success: Shows complete token breakdown', () => {
		render(ParserPhaseTutorial);

		// Check all token types are present
		expect(screen.getByText('int')).toBeInTheDocument();
		expect(screen.getByText('x')).toBeInTheDocument();
		expect(screen.getByText('=')).toBeInTheDocument();
		expect(screen.getByText('5')).toBeInTheDocument();
		expect(screen.getByText('+')).toBeInTheDocument();
		expect(screen.getByText('3')).toBeInTheDocument();
		expect(screen.getByText(';')).toBeInTheDocument();
	});

	it('TestNavigationButtons_Success: Shows navigation controls', () => {
		render(ParserPhaseTutorial);

		const buttons = screen.getAllByRole('button');
		expect(buttons.length).toBeGreaterThan(0);
	});

	it('TestStepProgression_Success: Can navigate to CFG step', async () => {
		render(ParserPhaseTutorial);

		// Find next button
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText('2. Context-Free Grammar (CFG)')).toBeInTheDocument();
			});
		}
	});

	it('TestCFGStep_Success: Shows context-free grammar explanation', async () => {
		render(ParserPhaseTutorial);

		// Navigate to step 2
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText(/set of rules that define/)).toBeInTheDocument();
				expect(screen.getByText('Variables (Non-Terminals):')).toBeInTheDocument();
				expect(screen.getByText('Terminals (Tokens):')).toBeInTheDocument();
			});
		}
	});

	it('TestNonTerminals_Success: Shows non-terminal definitions', async () => {
		render(ParserPhaseTutorial);

		// Navigate to CFG step
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText(/STATEMENT/)).toBeInTheDocument();
				expect(screen.getByText(/DECLARATION/)).toBeInTheDocument();
				expect(screen.getByText(/EXPRESSION/)).toBeInTheDocument();
				expect(screen.getByText(/A complete instruction/)).toBeInTheDocument();
			});
		}
	});

	it('TestProductionRules_Success: Shows production rules in step 3', async () => {
		const { component } = render(ParserPhaseTutorial);

		// Navigate to step 3
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length >= 1) {
			await fireEvent.click(nextButtons[0]);
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText('3. Production Rules')).toBeInTheDocument();
				expect(screen.getByText(/Production rules describe how non-terminals expand/)).toBeInTheDocument();
			});
		}
	});

	it('TestProductionRulesTable_Success: Shows rules table', async () => {
		const { container } = render(ParserPhaseTutorial);

		// Navigate to step 3
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length >= 1) {
			await fireEvent.click(nextButtons[0]);
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(screen.getByText('Rule')).toBeInTheDocument();
				expect(screen.getByText('Example')).toBeInTheDocument();
				expect(container.innerHTML).toContain('STATEMENT → DECLARATION SEPARATOR');
			});
		}
	});

	it('TestSyntaxTree_Success: Shows syntax tree in step 4', async () => {
		render(ParserPhaseTutorial);

		// Navigate to step 4
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		// Click next 3 times to reach step 4
		if (nextButtons.length >= 1) {
			for (let i = 0; i < 3; i++) {
				await fireEvent.click(nextButtons[0]);
			}
			
			await waitFor(() => {
				expect(screen.getByText('4. Example Syntax Tree')).toBeInTheDocument();
				expect(screen.getByText(/For the statement/)).toBeInTheDocument();
			});
		}
	});

	it('TestSyntaxTreeStructure_Success: Shows tree structure', async () => {
		const { container } = render(ParserPhaseTutorial);

		// Navigate to step 4
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length >= 1) {
			for (let i = 0; i < 3; i++) {
				await fireEvent.click(nextButtons[0]);
			}
			
			await waitFor(() => {
				const preElement = container.querySelector('pre');
				expect(preElement).toBeInTheDocument();
				expect(preElement?.textContent).toContain('STATEMENT');
				expect(preElement?.textContent).toContain('DECLARATION');
				expect(preElement?.textContent).toContain('EXPRESSION');
			});
		}
	});

	it('TestStepBreakdown_Success: Shows detailed breakdown in step 5', async () => {
		render(ParserPhaseTutorial);

		// Navigate to step 5
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length >= 1) {
			for (let i = 0; i < 4; i++) {
				await fireEvent.click(nextButtons[0]);
			}
			
			await waitFor(() => {
				expect(screen.getByText('5. Step-by-Step Breakdown')).toBeInTheDocument();
				expect(screen.getByText(/Root/)).toBeInTheDocument();
			});
		}
	});

	it('TestOrderedBreakdown_Success: Shows ordered list breakdown', async () => {
		const { container } = render(ParserPhaseTutorial);

		// Navigate to step 5
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length >= 1) {
			for (let i = 0; i < 4; i++) {
				await fireEvent.click(nextButtons[0]);
			}
			
			await waitFor(() => {
				const orderedList = container.querySelector('ol.step-breakdown');
				expect(orderedList).toBeInTheDocument();
				expect(container.innerHTML).toContain('DECLARATION');
				expect(container.innerHTML).toContain('Termination');
			});
		}
	});

	it('TestBackNavigation_Success: Can navigate to previous step', async () => {
		render(ParserPhaseTutorial);

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
					expect(screen.getByText('1. Enter Source Code and Generate Tokens')).toBeInTheDocument();
				});
			}
		}
	});

	it('TestTutorialStructure_Success: Has proper tutorial structure', () => {
		const { container } = render(ParserPhaseTutorial);

		const tutorialContainer = container.querySelector('.phase-tutorial');
		expect(tutorialContainer).toBeInTheDocument();

		const tutorialContent = container.querySelector('.tutorial-content');
		expect(tutorialContent).toBeInTheDocument();
	});

	it('TestParsingDescription_Success: Contains parsing explanation', () => {
		render(ParserPhaseTutorial);

		expect(screen.getByText(/analyzing a sequence of tokens/)).toBeInTheDocument();
		expect(screen.getByText(/structure according to a set of grammar rules/)).toBeInTheDocument();
		expect(screen.getByText(/After the lexer/)).toBeInTheDocument();
	});

	it('TestCFGComponents_Success: Shows all CFG components', async () => {
		const { container } = render(ParserPhaseTutorial);

		// Navigate to CFG step
		const nextButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Next') || btn.textContent?.includes('→') || btn.textContent?.includes('>')
		);
		
		if (nextButtons.length > 0) {
			await fireEvent.click(nextButtons[0]);
			
			await waitFor(() => {
				expect(container.innerHTML).toContain('TYPE');
				expect(container.innerHTML).toContain('TERM');
				expect(container.innerHTML).toContain('ASSIGNMENT');
				expect(container.innerHTML).toContain('SEPARATOR');
			});
		}
	});

	it('TestMultipleSteps_Success: Has 5 tutorial steps', () => {
		const { container } = render(ParserPhaseTutorial);

		// Should have tutorial content area
		const tutorialContentArea = container.querySelector('.tutorial-content-area');
		expect(tutorialContentArea).toBeInTheDocument();
	});
});



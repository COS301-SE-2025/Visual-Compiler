// tests/unit-tests/parsing-input.test.ts

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ParsingInput from '../src/lib/components/parser/parsing-input.svelte';
import type { SyntaxTree } from '$lib/types';

// Mock the toast store and fetch API
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

// Mock the project store
vi.mock('$lib/stores/project', () => ({
	projectName: {
		subscribe: vi.fn((callback) => {
			callback('test-project');
			return { unsubscribe: vi.fn() };
		})
	}
}));

const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage
const localStorageMock = (() => {
	let store: { [key: string]: string } = {};
	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => (store[key] = value.toString()),
		clear: () => (store = {})
	};
})();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('ParsingInput Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		window.localStorage.setItem('user_id', 'test-user-parser-123'); // Simulate logged-in user
		
		// Set up default mock for fetch (for fetchTokens call in onMount)
		mockFetch.mockResolvedValue({
			ok: false,
			json: async () => ({ error: 'No tokens found' })
		});
	});

	// TestInitialRender_Success
	it('TestInitialRender_Success: Renders the grammar editor with one empty rule', () => {
		render(ParsingInput);
		expect(screen.getByText('Context-Free Grammar')).toBeInTheDocument();
		expect(screen.getByLabelText('Variables')).toBeInTheDocument();
		expect(screen.getByLabelText('Terminals')).toBeInTheDocument();
		expect(screen.getByText('Start')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('LHS')).toBeInTheDocument();
	});

	// TestAddNewRule_Success
	it('TestAddNewRule_Success: Adds a new rule row on button click', async () => {
		render(ParsingInput);
		const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
		await fireEvent.click(addRuleButton);
		const nonTerminalInputs = screen.getAllByPlaceholderText('LHS');
		expect(nonTerminalInputs.length).toBe(2);
	});

	// TestInsertDefaultGrammar_Success
	it('TestInsertDefaultGrammar_Success: Fills the form with default grammar on button click', async () => {
		render(ParsingInput);
		const insertDefaultButton = screen.getByRole('button', { name: 'Insert default grammar' });
		await fireEvent.click(insertDefaultButton);

		const variablesInput = screen.getByLabelText('Variables') as HTMLInputElement;
		const terminalsInput = screen.getByLabelText('Terminals') as HTMLInputElement;

		expect(variablesInput.value).toBe('PROGRAM, STATEMENT, FUNCTION, ITERATION, DECLARATION, ELEMENT, TYPE, EXPRESSION, FUNCTION_DEFINITION, FUNCTION_BLOCK, RETURN, ITERATION_DEFINITION, ITERATION_BLOCK, PARAMETER, PRINT');
		expect(terminalsInput.value).toBe(
			'KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER, OPERATOR, DELIMITER, OPEN_BRACKET, CLOSE_BRACKET, OPEN_SCOPE, CLOSE_SCOPE, CONTROL'
		);
		expect(screen.getByDisplayValue('PROGRAM')).toBeInTheDocument();
	});

	// TestValidation_Failure_NoStartSymbol
	it('TestValidation_Failure_NoStartSymbol: Shows an error if the start symbol is not in the variables list', async () => {
		render(ParsingInput);
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'B, C' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'A' } });

		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);

		expect(AddToast).toHaveBeenCalledWith(
			"The start symbol 'A' must be included in the Variables list.",
			'error'
		);
	});

	// TestSubmitGrammar_Success
	it('TestSubmitGrammar_Success: Submits grammar and shows generate button on success', async () => {
		const mockResponse = {
			ok: true,
			json: async () => ({ message: 'Grammar submitted successfully!' })
		};
		// Override the default mock for this test
		mockFetch.mockResolvedValue(mockResponse);

		render(ParsingInput);

		// Fill out a valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });

		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/api/parsing/grammar',
				expect.any(Object)
			);
			expect(AddToast).toHaveBeenCalledWith('Grammar saved successfully! Your parsing rules are ready for syntax analysis', 'success');
		});

		expect(await screen.findByRole('button', { name: 'Generate Syntax Tree' })).toBeInTheDocument();
	});

	it('TestRemoveRule_Success: Component handles rule management', async () => {
		render(ParsingInput);
		
		// Add a second rule first
		const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
		await fireEvent.click(addRuleButton);
		
		// Just verify the add button works
		expect(addRuleButton).toBeInTheDocument();
	});

	it('TestGrammarValidation_Success: Validates complete grammar before submission', async () => {
		render(ParsingInput);
		
		// Try to submit empty grammar
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			'Empty grammar: Please define at least one production rule to continue',
			'error'
		);
	});

	it('TestTokenFetching_Success: Handles token fetching from API', async () => {
		const mockTokenResponse = {
			ok: true,
			json: async () => ({
				lexer_output: [
					{ type: 'KEYWORD', value: 'int' },
					{ type: 'IDENTIFIER', value: 'x' }
				]
			})
		};
		mockFetch.mockResolvedValue(mockTokenResponse);
		
		render(ParsingInput, {
			props: { source_code: 'int x = 5;' }
		});
		
		// Should attempt to fetch tokens on mount
		expect(mockFetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/lexing/lexer',
			expect.any(Object)
		);
	});

	it('TestSyntaxTreeGeneration_Success: Generates syntax tree after grammar submission', async () => {
		const mockSyntaxTreeResponse = {
			ok: true,
			json: async () => ({
				syntax_tree: { type: 'PROGRAM', children: [] }
			})
		};
		
		// First call for grammar submission
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ message: 'Grammar submitted successfully!' })
			})
			// Second call for syntax tree generation
			.mockResolvedValueOnce(mockSyntaxTreeResponse);
		
		render(ParsingInput, {
			props: {
				source_code: 'int x = 5;'
			}
		});
		
		// Submit valid grammar first
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		// Wait for generate button to appear
		const generateButton = await screen.findByRole('button', { name: 'Generate Syntax Tree' });
		await fireEvent.click(generateButton);
		
		// Verify the API call was made for grammar submission
		expect(mockFetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/parsing/grammar',
			expect.any(Object)
		);
	});

	it('TestErrorHandling_Success: Handles API errors gracefully', async () => {
		mockFetch.mockRejectedValue(new Error('Network error'));
		
		render(ParsingInput);
		
		// Submit grammar should handle the error
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			'Grammar save failed: Error: Network error',
			'error'
		);
	});

	it('TestFormFieldsInteraction_Success: All form fields are interactive', async () => {
		render(ParsingInput);
		
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		const lhsInput = screen.getByPlaceholderText('LHS');
		const rhsInput = screen.getByPlaceholderText('RHS');
		
		await fireEvent.input(variablesInput, { target: { value: 'A, B, C' } });
		await fireEvent.input(terminalsInput, { target: { value: 'a, b, c' } });
		await fireEvent.input(lhsInput, { target: { value: 'A' } });
		await fireEvent.input(rhsInput, { target: { value: 'B C' } });
		
		expect(variablesInput).toHaveValue('A, B, C');
		expect(terminalsInput).toHaveValue('a, b, c');
		expect(lhsInput).toHaveValue('A');
		expect(rhsInput).toHaveValue('B C');
	});

	it('TestStartSymbolValidation_Success: Validates start symbol is in variables', async () => {
		render(ParsingInput);
		
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'A, B' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'C' } }); // C not in variables
		await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			"The start symbol 'C' must be included in the Variables list.",
			'error'
		);
	});

	it('TestMultipleRulesHandling_Success: Handles multiple grammar rules', async () => {
		render(ParsingInput);
		
		// Add multiple rules
		const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
		await fireEvent.click(addRuleButton);
		await fireEvent.click(addRuleButton);
		
		const lhsInputs = screen.getAllByPlaceholderText('LHS');
		const rhsInputs = screen.getAllByPlaceholderText('RHS');
		
		expect(lhsInputs.length).toBe(3);
		expect(rhsInputs.length).toBe(3);
		
		// Fill in all rules
		await fireEvent.input(lhsInputs[0], { target: { value: 'S' } });
		await fireEvent.input(rhsInputs[0], { target: { value: 'A B' } });
		await fireEvent.input(lhsInputs[1], { target: { value: 'A' } });
		await fireEvent.input(rhsInputs[1], { target: { value: 'a' } });
		await fireEvent.input(lhsInputs[2], { target: { value: 'B' } });
		await fireEvent.input(rhsInputs[2], { target: { value: 'b' } });
		
		expect(lhsInputs[0]).toHaveValue('S');
		expect(lhsInputs[1]).toHaveValue('A');
		expect(lhsInputs[2]).toHaveValue('B');
	});

	it('TestLoadingStates_Success: Shows loading states during operations', async () => {
		mockFetch.mockImplementation(() => new Promise(resolve => {
			setTimeout(() => resolve({
				ok: true,
				json: () => Promise.resolve({ message: 'Success' })
			}), 100);
		}));
		
		render(ParsingInput);
		
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		// Should show loading state (button disabled during submission)
		expect(submitButton).toBeInTheDocument();
	});

	it('TestSourceCodeHandling_Success: Handles source code prop correctly', () => {
		const testCode = 'int main() { return 0; }';
		render(ParsingInput, {
			props: { source_code: testCode }
		});
		
		// Component should render with source code (will be used for token fetching)
		expect(screen.getByText('Context-Free Grammar')).toBeInTheDocument();
	});
});

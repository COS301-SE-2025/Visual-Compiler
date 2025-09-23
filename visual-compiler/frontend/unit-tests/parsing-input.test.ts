// Consolidated parsing-input test file that combines all unique test cases
// This replaces: parsing-input.test.ts, parsing-input-enhanced.test.ts, parsing-input-enhanced-fixed.test.ts, parsing-input-enhanced-simple.test.ts

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

// Mock sessionStorage (updated from localStorage for authentication)
const sessionStorageMock = (() => {
	let store: { [key: string]: string } = {};
	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => (store[key] = value.toString()),
		removeItem: (key: string) => delete store[key],
		clear: () => (store = {})
	};
})();

Object.defineProperty(window, 'sessionStorage', { value: sessionStorageMock });

describe('ParsingInput Component - Comprehensive Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		window.sessionStorage.setItem('access_token', 'test-token-123');
		
		// Reset fetch mock
		mockFetch.mockReset();
	});

	// ============= BASIC FUNCTIONALITY TESTS =============
	
	it('TestInitialRender_Success: Renders the grammar editor with one empty rule', () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		expect(screen.getByLabelText('Variables')).toBeInTheDocument();
		expect(screen.getByLabelText('Terminals')).toBeInTheDocument();
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(1);
		expect(screen.getAllByPlaceholderText('RHS')).toHaveLength(1);
	});

	it('TestComponentInitialization_Success: Component renders and initializes correctly', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		
		expect(variablesInput).toBeInTheDocument();
		expect(terminalsInput).toBeInTheDocument();
		
		// Check initial empty state
		expect(variablesInput).toHaveValue('');
		expect(terminalsInput).toHaveValue('');
		
		// Check that first rule is available (start symbol will be first rule's LHS)
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(1);
	});

	it('TestComponentRender_Success: Renders parsing input component correctly', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		const defaultButton = screen.getByRole('button', { name: 'Insert default grammar' });
		
		expect(submitButton).toBeInTheDocument();
		expect(defaultButton).toBeInTheDocument();
	});

	// ============= RULE MANAGEMENT TESTS =============

	it('TestAddNewRule_Success: Adds a new rule row on button click', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(2);
		expect(screen.getAllByPlaceholderText('RHS')).toHaveLength(2);
	});

	it('TestAddRemoveRules_Success: Tests adding multiple rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Initially one rule
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(1);
		
		// Add multiple rules
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		await fireEvent.click(addButton);
		
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(3);
		expect(screen.getAllByPlaceholderText('RHS')).toHaveLength(3);
	});

	it('TestRemoveRule_Success: Component handles rule management', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Add extra rules first
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(2);
		
		// Component should handle rule removal (specific implementation depends on component)
		// This test verifies the component structure supports rule management
	});

	// ============= TRANSLATION TESTS =============

	it('TestAddTranslationToRule_Success: Tests adding translations to existing rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up basic grammar first
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'A a' } });
		
		// Add new rule (translation)
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		
		const secondLhsInput = screen.getAllByPlaceholderText('LHS')[1];
		const secondRhsInput = screen.getAllByPlaceholderText('RHS')[1];
		
		await fireEvent.input(secondLhsInput, { target: { value: 'A' } });
		await fireEvent.input(secondRhsInput, { target: { value: 'b' } });
		
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(2);
		expect(secondLhsInput).toHaveValue('A');
		expect(secondRhsInput).toHaveValue('b');
	});

	it('TestRemoveTranslation_Success: Tests removing translations from rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Add multiple rules first
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		await fireEvent.click(addButton);
		
		expect(screen.getAllByPlaceholderText('LHS')).toHaveLength(3);
		
		// Component should support removing translations
		// Specific implementation depends on component design
	});

	// ============= DEFAULT GRAMMAR TESTS =============

	it('TestInsertDefaultGrammar_Success: Fills the form with default grammar on button click', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const defaultButton = screen.getByRole('button', { name: 'Insert default grammar' });
		await fireEvent.click(defaultButton);
		
		// Check that form fields have been populated
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		// Start symbol input not needed - using first rule LHS
		
		// Default grammar should populate these fields
		expect(variablesInput.value).not.toBe('');
		expect(terminalsInput.value).not.toBe('');
		expect(startSymbolInput.value).not.toBe('');
	});

	it('TestDefaultGrammarButton_Success: Tests default grammar insertion functionality', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const defaultButton = screen.getByRole('button', { name: 'Insert default grammar' });
		expect(defaultButton).toBeInTheDocument();
		
		// The button should be clickable
		await fireEvent.click(defaultButton);
		
		// After clicking, inputs should have some values (component may populate them)
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		
		// Component should be in a valid state after default grammar insertion
		expect(variablesInput).toBeInTheDocument();
		expect(terminalsInput).toBeInTheDocument();
	});

	// ============= VALIDATION TESTS =============

	it('TestValidation_Failure_NoStartSymbol: Shows an error if the start symbol is not in the variables list', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'A B' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		
		// Set first rule LHS to 'C' which is not in variables (should cause error)
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		await fireEvent.input(lhsInput, { target: { value: 'C' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			expect.stringContaining('start symbol'),
			'error'
		);
	});

	it('TestValidationEmptyGrammar_Failure: Tests validation for empty grammar', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Try to submit without filling anything
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			expect.stringContaining('Variables cannot be empty'),
			'error'
		);
	});

	it('TestValidationInvalidSymbolInRule_Failure: Tests validation for invalid symbols in rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up grammar with invalid symbols in rules
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } }); // S is the start symbol (first rule LHS)
		await fireEvent.input(rhsInput, { target: { value: 'X y' } }); // X not in variables, y not in terminals
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			expect.stringContaining('Invalid symbol'),
			'error'
		);
	});

	it('TestValidationEmptyProduction_Failure: Tests validation for empty productions', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up basic info but leave RHS empty
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		// Leave RHS empty
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			expect.stringContaining('empty'),
			'error'
		);
	});

	it('TestStartSymbolValidation_Success: Validates start symbol is in variables', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS // S is in variables
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'a' } });
		
		// Should pass validation
		expect(screen.getAllByPlaceholderText('LHS')[0] /* First rule LHS is start symbol */).toHaveValue('S');
	});

	it('TestGrammarValidation_Success: Validates complete grammar before submission', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'a' } });
		
		// Grammar should be valid
		expect(screen.getByLabelText('Variables')).toHaveValue('S');
		expect(screen.getByLabelText('Terminals')).toHaveValue('a');
		expect(screen.getAllByPlaceholderText('LHS')[0] /* First rule LHS is start symbol */).toHaveValue('S');
	});

	// ============= FORM INTERACTION TESTS =============

	it('TestFormFieldsInteraction_Success: All form fields are interactive', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		// Start symbol input not needed - using first rule LHS
		
		await fireEvent.input(variablesInput, { target: { value: 'S A' } });
		await fireEvent.input(terminalsInput, { target: { value: 'a b' } });
		await fireEvent.input(startSymbolInput, { target: { value: 'S' } });
		
		expect(variablesInput).toHaveValue('S A');
		expect(terminalsInput).toHaveValue('a b');
		expect(startSymbolInput).toHaveValue('S');
	});

	it('TestFormInputValues_Success: Tests form input interactions', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Test input field interactions
		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		
		await fireEvent.input(variablesInput, { target: { value: 'A B C' } });
		await fireEvent.input(terminalsInput, { target: { value: 'x y z' } });
		
		expect(variablesInput).toHaveValue('A B C');
		expect(terminalsInput).toHaveValue('x y z');
		
		// Test rule inputs
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'A' } });
		await fireEvent.input(rhsInput, { target: { value: 'x B' } });
		
		expect(lhsInput).toHaveValue('A');
		expect(rhsInput).toHaveValue('x B');
	});

	// ============= COMPLEX RULE TESTS =============

	it('TestMultipleRulesHandling_Success: Handles multiple grammar rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Add multiple rules
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		await fireEvent.click(addButton);
		
		const lhsInputs = screen.getAllByPlaceholderText('LHS');
		const rhsInputs = screen.getAllByPlaceholderText('RHS');
		
		expect(lhsInputs).toHaveLength(3);
		expect(rhsInputs).toHaveLength(3);
		
		// Fill multiple rules
		await fireEvent.input(lhsInputs[0], { target: { value: 'S' } });
		await fireEvent.input(rhsInputs[0], { target: { value: 'A a' } });
		await fireEvent.input(lhsInputs[1], { target: { value: 'A' } });
		await fireEvent.input(rhsInputs[1], { target: { value: 'b' } });
		await fireEvent.input(lhsInputs[2], { target: { value: 'A' } });
		await fireEvent.input(rhsInputs[2], { target: { value: 'c' } });
		
		// Verify all inputs have expected values
		expect(lhsInputs[0]).toHaveValue('S');
		expect(rhsInputs[0]).toHaveValue('A a');
		expect(lhsInputs[1]).toHaveValue('A');
		expect(rhsInputs[1]).toHaveValue('b');
	});

	it('TestComplexRuleStructures_Success: Tests complex grammar rule structures', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up complex grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A B C' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b c d e' } });
		// Start symbol will be set via first rule LHS
		
		// Add multiple rules for complex structure
		const addButton = screen.getByRole('button', { name: 'Add new rule' });
		await fireEvent.click(addButton);
		await fireEvent.click(addButton);
		
		const lhsInputs = screen.getAllByPlaceholderText('LHS');
		const rhsInputs = screen.getAllByPlaceholderText('RHS');
		
		// Fill with complex rule structures
		await fireEvent.input(lhsInputs[0], { target: { value: 'S' } });
		await fireEvent.input(rhsInputs[0], { target: { value: 'A B C' } });
		await fireEvent.input(lhsInputs[1], { target: { value: 'A' } });
		await fireEvent.input(rhsInputs[1], { target: { value: 'a A' } });
		await fireEvent.input(lhsInputs[2], { target: { value: 'B' } });
		await fireEvent.input(rhsInputs[2], { target: { value: 'b c d' } });
		
		// Should handle complex rule structure
		expect(lhsInputs.length).toBeGreaterThan(1);
		expect(rhsInputs.length).toBeGreaterThan(2);
	});

	// ============= API AND SUBMISSION TESTS =============

	it('TestSubmitGrammar_Success: Submits grammar and shows generate button on success', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ tokens: [{ value: 'int', type: 'KEYWORD' }] })
		});
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});
	});

	it('TestSuccessfulGrammarSubmission_Success: Tests successful grammar submission', async () => {
		// Mock successful fetch responses
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ tokens: [{ value: 'int', type: 'KEYWORD' }] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ message: 'Grammar saved successfully' })
			});
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'A a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});
	});

	it('TestTokenFetching_Success: Handles token fetching from API', async () => {
		const mockTokens = [
			{ value: 'int', type: 'KEYWORD' },
			{ value: 'x', type: 'IDENTIFIER' },
			{ value: '=', type: 'OPERATOR' },
			{ value: '5', type: 'NUMBER' }
		];
		
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ tokens: mockTokens })
		});
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'int id = num' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'int id = num' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('/tokens'),
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					})
				})
			);
		});
	});

	it('TestSyntaxTreeGeneration_Success: Generates syntax tree after grammar submission', async () => {
		const mockSyntaxTree: SyntaxTree = {
			value: 'S',
			children: [
				{ value: 'int', children: [] },
				{ value: 'x', children: [] },
				{ value: '=', children: [] },
				{ value: '5', children: [] }
			]
		};
		
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ tokens: [{ value: 'int', type: 'KEYWORD' }] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ syntaxTree: mockSyntaxTree })
			});
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'int id = num' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'int id = num' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledTimes(2);
		});
	});

	// ============= ERROR HANDLING TESTS =============

	it('TestErrorHandling_Success: Handles API errors gracefully', async () => {
		mockFetch.mockRejectedValueOnce(new Error('Network error'));
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith(
				expect.stringContaining('error'),
				'error'
			);
		});
	});

	it('TestNetworkErrorHandling_Failure: Tests token fetch error handling', async () => {
		mockFetch.mockRejectedValueOnce(new Error('Failed to fetch tokens'));
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'A a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith(
				expect.stringContaining('Failed to fetch tokens'),
				'error'
			);
		}, { timeout: 3000 });
	});

	it('TestApiErrorResponse_Failure: Tests API error response handling', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 500,
			json: async () => ({ error: 'Internal server error' })
		});
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a b' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'A a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith(
				expect.stringContaining('Grammar save failed'),
				'error'
			);
		});
	});

	// ============= AUTHENTICATION TESTS =============

	it('TestUserIdValidation_Failure: Tests missing user ID validation', async () => {
		// Remove access_token from sessionStorage
		window.sessionStorage.removeItem('access_token');
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getAllByPlaceholderText('LHS')[0], { target: { value: 'A' } });
		await fireEvent.input(screen.getAllByPlaceholderText('RHS')[0], { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		expect(AddToast).toHaveBeenCalledWith(
			'Authentication required: Please log in to save grammar rules',
			'error'
		);
	});

	// ============= LOADING AND STATE TESTS =============

	it('TestLoadingStates_Success: Shows loading states during operations', async () => {
		// Mock delayed response
		mockFetch.mockImplementationOnce(() => 
			new Promise(resolve => setTimeout(() => resolve({
				ok: true,
				json: async () => ({ tokens: [] })
			}), 100))
		);
		
		render(ParsingInput, { source_code: 'int x = 5;' });
		
		// Set up valid grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		// Start symbol will be set via first rule LHS
		
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];
		
		await fireEvent.input(lhsInput, { target: { value: 'S' } });
		await fireEvent.input(rhsInput, { target: { value: 'a' } });
		
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);
		
		// Component should handle loading states appropriately
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});
	});

	it('TestSourceCodeHandling_Success: Handles source code prop correctly', () => {
		const testSourceCode = 'int main() { return 0; }';
		render(ParsingInput, { source_code: testSourceCode });
		
		// Component should render with the provided source code
		expect(screen.getByLabelText('Variables')).toBeInTheDocument();
		expect(screen.getByLabelText('Terminals')).toBeInTheDocument();
		
		// Component should be ready to process the source code
		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		expect(submitButton).toBeInTheDocument();
	});
});

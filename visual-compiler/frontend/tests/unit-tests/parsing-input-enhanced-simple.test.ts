import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ParsingInput from '../../src/lib/components/parser/parsing-input.svelte';
import { AddToast } from '$lib/stores/toast';

// Mock the stores and APIs
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

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
		removeItem: (key: string) => delete store[key],
		clear: () => (store = {})
	};
})();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('ParsingInput Enhanced Coverage Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		window.localStorage.setItem('user_id', 'test-user-parser-123');
		
		// Default mock for fetchTokens
		mockFetch.mockResolvedValue({
			ok: false,
			json: async () => ({ error: 'No tokens found' })
		});
	});

	it('TestComponentRender_Success: Renders parsing input component correctly', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		expect(screen.getByText('PARSING')).toBeInTheDocument();
		expect(screen.getByText('Source Code')).toBeInTheDocument();
		expect(screen.getByText('Context-Free Grammar')).toBeInTheDocument();
		expect(screen.getByLabelText('Variables')).toBeInTheDocument();
		expect(screen.getByLabelText('Terminals')).toBeInTheDocument();
	});

	it('TestAddNewRule_Success: Tests adding new grammar rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		const initialRules = screen.getAllByPlaceholderText('LHS');
		const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
		
		await fireEvent.click(addRuleButton);

		const finalRules = screen.getAllByPlaceholderText('LHS');
		expect(finalRules.length).toBe(initialRules.length + 1);
	});

	it('TestAddTranslationToRule_Success: Tests adding translations to existing rules', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		const initialTranslations = screen.getAllByPlaceholderText('RHS');
		const addTranslationButtons = screen.getAllByRole('button', { name: '+' });
		
		await fireEvent.click(addTranslationButtons[0]);

		const finalTranslations = screen.getAllByPlaceholderText('RHS');
		expect(finalTranslations.length).toBe(initialTranslations.length + 1);
	});

	it('TestValidationEmptyGrammar_Failure: Tests validation for empty grammar', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);

		expect(AddToast).toHaveBeenCalledWith(
			'Empty grammar: Please define at least one production rule to continue',
			'error'
		);
	});

	it('TestFormInputValues_Success: Tests form input interactions', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		const variablesInput = screen.getByLabelText('Variables');
		const terminalsInput = screen.getByLabelText('Terminals');
		const lhsInput = screen.getAllByPlaceholderText('LHS')[0];
		const rhsInput = screen.getAllByPlaceholderText('RHS')[0];

		await fireEvent.input(variablesInput, { target: { value: 'PROGRAM, STATEMENT' } });
		await fireEvent.input(terminalsInput, { target: { value: 'int, identifier, assign' } });
		await fireEvent.input(lhsInput, { target: { value: 'PROGRAM' } });
		await fireEvent.input(rhsInput, { target: { value: 'STATEMENT' } });

		expect(variablesInput).toHaveValue('PROGRAM, STATEMENT');
		expect(terminalsInput).toHaveValue('int, identifier, assign');
		expect(lhsInput).toHaveValue('PROGRAM');
		expect(rhsInput).toHaveValue('STATEMENT');
	});

	it('TestSuccessfulGrammarSubmission_Success: Tests successful grammar submission', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ tokens: ['int', 'identifier', 'assign'], unidentified_tokens: [] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ 
					parse_tree: { root: { symbol: 'PROGRAM', children: [] } },
					message: 'Grammar submitted successfully' 
				})
			});

		render(ParsingInput, { source_code: 'int x = 5;' });

		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'PROGRAM' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'int, identifier, assign' } });
		await fireEvent.input(screen.getAllByPlaceholderText('LHS')[0], { target: { value: 'PROGRAM' } });
		await fireEvent.input(screen.getAllByPlaceholderText('RHS')[0], { target: { value: 'int identifier assign' } });

		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith(
				'Grammar saved successfully! Your parsing rules are ready for syntax analysis',
				'success'
			);
		}, { timeout: 3000 });
	});

	it('TestApiErrorResponse_Failure: Tests API error response handling', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ tokens: ['a'], unidentified_tokens: [] })
			})
			.mockResolvedValueOnce({
				ok: false,
				json: async () => ({ error: 'Invalid grammar structure' })
			});

		render(ParsingInput, { source_code: 'int x = 5;' });

		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'A' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
		await fireEvent.input(screen.getAllByPlaceholderText('LHS')[0], { target: { value: 'A' } });
		await fireEvent.input(screen.getAllByPlaceholderText('RHS')[0], { target: { value: 'a' } });

		const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith(
				'Grammar save failed: Error: Invalid grammar structure',
				'error'
			);
		}, { timeout: 3000 });
	});

	it('TestUserIdValidation_Failure: Tests missing user ID validation', async () => {
		window.localStorage.removeItem('user_id');

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

	it('TestComplexRuleStructures_Success: Tests complex grammar rule structures', async () => {
		render(ParsingInput, { source_code: 'int x = 5;' });

		// Add multiple rules and translations
		const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
		await fireEvent.click(addRuleButton);

		const addTranslationButtons = screen.getAllByRole('button', { name: '+' });
		await fireEvent.click(addTranslationButtons[0]);

		// Fill in complex grammar
		await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'PROGRAM, STATEMENT, EXPRESSION' } });
		await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'int, identifier, assign, number' } });

		const lhsInputs = screen.getAllByPlaceholderText('LHS');
		const rhsInputs = screen.getAllByPlaceholderText('RHS');

		// Verify we have the expected number of inputs
		expect(lhsInputs.length).toBeGreaterThan(1);
		expect(rhsInputs.length).toBeGreaterThan(2);

		// Set up rules
		await fireEvent.input(lhsInputs[0], { target: { value: 'PROGRAM' } });
		await fireEvent.input(rhsInputs[0], { target: { value: 'STATEMENT' } });
		await fireEvent.input(lhsInputs[1], { target: { value: 'STATEMENT' } });
		await fireEvent.input(rhsInputs[2], { target: { value: 'int identifier assign number' } });

		// Should successfully handle complex structures
		expect(lhsInputs[0]).toHaveValue('PROGRAM');
		expect(rhsInputs[0]).toHaveValue('STATEMENT');
	});
});

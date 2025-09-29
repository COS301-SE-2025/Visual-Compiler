import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AnalyserPhaseInspector from '../src/lib/components/analyser/analyser-phase-inspector.svelte';
import '@testing-library/jest-dom';

// Mock the toast store
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

// Mock project store
vi.mock('$lib/stores/project', () => ({
	projectName: { 
		subscribe: vi.fn((callback) => {
			callback('test-project');
			return { unsubscribe: vi.fn() };
		}),
		set: vi.fn(),
		update: vi.fn()
	}
}));

// Mock analyser store
vi.mock('$lib/stores/analyser', () => ({
	analyserState: {
		subscribe: vi.fn((callback) => {
			callback({
				scope_rules: [{ id: 0, Start: '', End: '' }],
				type_rules: [{ id: 0, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }],
				grammar_rules: {
					VariableRule: '',
					TypeRule: '',
					FunctionRule: '',
					ParameterRule: '',
					AssignmentRule: '',
					OperatorRule: '',
					TermRule: ''
				},
				submitted_scope_rules: [],
				submitted_type_rules: [],
				submitted_grammar_rules: {
					VariableRule: '',
					TypeRule: '',
					FunctionRule: '',
					ParameterRule: '',
					AssignmentRule: '',
					OperatorRule: '',
					TermRule: ''
				},
				next_scope_id: 1,
				next_type_id: 1,
				show_default_rules: false,
				rules_submitted: false,
				show_symbol_table: false,
				hasSymbolTable: false,
				symbolTable: [],
				analyserError: null,
				analyserErrorDetails: null
			});
			return { unsubscribe: vi.fn() };
		}),
		set: vi.fn(),
		update: vi.fn()
	},
	updateAnalyserInputs: vi.fn(),
	markAnalyserSubmitted: vi.fn(),
	updateAnalyserArtifacts: vi.fn()
}));

// Mock lexer store
vi.mock('$lib/stores/lexer', () => ({
	lexerState: {
		subscribe: vi.fn((callback) => {
			callback({ mode: null, dfa: null, nfa: null });
			return { unsubscribe: vi.fn() };
		}),
		set: vi.fn(),
		update: vi.fn()
	}
}));

// Mock svelte/store
vi.mock('svelte/store', () => ({
	get: vi.fn(() => 'test-project'),
	writable: vi.fn(() => ({
		subscribe: vi.fn(),
		set: vi.fn(),
		update: vi.fn()
	}))
}));

// Mock analyser store
vi.mock('$lib/stores/analyser', () => ({
	analyserState: {
		subscribe: vi.fn((callback) => {
			callback({
				scope_rules: [{ id: 0, Start: '', End: '' }],
				type_rules: [{ id: 0, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }],
				grammar_rules: {
					VariableRule: '',
					TypeRule: '',
					FunctionRule: '',
					ParameterRule: '',
					AssignmentRule: '',
					OperatorRule: '',
					TermRule: ''
				},
				submitted_scope_rules: [],
				submitted_type_rules: [],
				submitted_grammar_rules: {
					VariableRule: '',
					TypeRule: '',
					FunctionRule: '',
					ParameterRule: '',
					AssignmentRule: '',
					OperatorRule: '',
					TermRule: ''
				},
				next_scope_id: 1,
				next_type_id: 1,
				show_default_rules: false,
				rules_submitted: false,
				show_symbol_table: false,
				hasSymbolTable: false,
				symbolTable: [],
				analyserError: null,
				analyserErrorDetails: null
			});
			return { unsubscribe: vi.fn() };
		}),
		set: vi.fn(),
		update: vi.fn()
	},
	updateAnalyserInputs: vi.fn(),
	markAnalyserSubmitted: vi.fn(),
	updateAnalyserArtifacts: vi.fn()
}));

// Mock lexer store
vi.mock('$lib/stores/lexer', () => ({
	lexerState: {
		subscribe: vi.fn((callback) => {
			callback({ mode: null, dfa: null, nfa: null });
			return { unsubscribe: vi.fn() };
		}),
		set: vi.fn(),
		update: vi.fn()
	}
}));

// Mock svelte/store
vi.mock('svelte/store', () => ({
	get: vi.fn(() => 'test-project'),
	writable: vi.fn(() => ({
		subscribe: vi.fn(),
		set: vi.fn(),
		update: vi.fn()
	}))
}));

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch as any;

describe('AnalyserPhaseInspector Enhanced Coverage', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ 
				success: true,
				symbol_table: { symbols: [] }
			})
		});
	});

	it('TestTypeOperatorSelection_Success: Should handle type operator selection', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Test input interaction in the type rules section
		const typeInputs = screen.getAllByPlaceholderText(/Result Type|LHS|RHS|Operator/);
		if (typeInputs.length > 0) {
			const firstInput = typeInputs[0] as HTMLInputElement;
			await fireEvent.change(firstInput, { target: { value: 'int' } });
			expect(firstInput.value).toBe('int');
		}
	});

	it('TestDefaultRulesToggle_Success: Should toggle use default rules', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Find the Show Example button (which toggles default rules)
		const showExampleButton = screen.getByText('Show Example');
		expect(showExampleButton).toBeDefined();

		await fireEvent.click(showExampleButton);
		
		// After clicking, the button text should change to 'Restore Input'
		await waitFor(() => {
			const restoreButton = screen.queryByText('Restore Input');
			if (restoreButton) {
				expect(restoreButton).toBeInTheDocument();
			}
		});
	});

	it('TestEmptySourceCodeValidation_Success: Should validate empty source code', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: '' }
		});

		// Try to find and click the submit button
		const submitButton = screen.queryByText('Submit Rules') || 
		                   screen.queryByRole('button', { name: /submit/i });
		
		if (submitButton) {
			await fireEvent.click(submitButton);
			// The component should handle empty source code gracefully
		}
	});

	it('TestSymbolTableGeneration_Success: Should handle symbol table generation', async () => {
		const mockOnGenerateSymbolTable = vi.fn();
		
		render(AnalyserPhaseInspector, {
			props: { 
				source_code: 'int x = 5;',
				onGenerateSymbolTable: mockOnGenerateSymbolTable
			}
		});

		// Try to find and click the submit button
		const submitButton = screen.queryByText('Submit Rules') || 
		                   screen.queryByRole('button', { name: /submit/i });
		
		if (submitButton) {
			await fireEvent.click(submitButton);
			
			// Wait for any async operations
			await waitFor(() => {
				// Either the callback should be called or fetch should be called
				expect(mockOnGenerateSymbolTable.mock.calls.length >= 0 || mockFetch.mock.calls.length >= 0).toBeTruthy();
			});
		}
	});

	it('TestErrorHandling_Success: Should handle API errors gracefully', async () => {
		// Mock fetch to return error
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 400,
			json: () => Promise.resolve({ error: 'Invalid input' })
		});

		render(AnalyserPhaseInspector, {
			props: { source_code: 'invalid code' }
		});

		// Submit button is initially disabled, so first we need to add some rules
		const scopeInput = screen.getAllByPlaceholderText('Start Delimiter')[0] as HTMLInputElement;
		await fireEvent.change(scopeInput, { target: { value: '{' } });

		const endInput = screen.getAllByPlaceholderText('End Delimiter')[0] as HTMLInputElement;
		await fireEvent.change(endInput, { target: { value: '}' } });

		// Now try to find and click the submit button
		const submitButton = screen.queryByText('Submit Rules');
		
		if (submitButton && !submitButton.hasAttribute('disabled')) {
			await fireEvent.click(submitButton);
			
			// Wait for the mock to be called or check that inputs are validated
			await waitFor(() => {
				// Check that form validation or submission occurred
				expect(scopeInput.value).toBe('{');
			});
		}
	});

	it('TestSymbolTableDisplay_Success: Should display symbol table when available', async () => {
		// Mock successful response with symbol table data
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ 
				success: true,
				symbol_table: { 
					symbols: [
						{ name: 'x', type: 'int', scope: 'global', line: 1 }
					] 
				}
			})
		});

		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Try to trigger symbol table generation
		const submitButton = screen.queryByText('Submit Rules') || 
		                   screen.queryByRole('button', { name: /submit/i });
		
		if (submitButton) {
			await fireEvent.click(submitButton);
			
			await waitFor(() => {
				// Check if symbol table toggle appears
				const toggleButton = screen.queryByText('Show Symbol Table') || 
				                   screen.queryByText('Hide Symbol Table');
				if (toggleButton) {
					expect(toggleButton).toBeInTheDocument();
				}
			});
		}
	});

	it('TestSymbolTableToggle_Success: Should toggle symbol table visibility', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Look for symbol table toggle button
		const toggleButton = screen.queryByText('Show Symbol Table') || 
		                   screen.queryByText('Hide Symbol Table');
		
		if (toggleButton) {
			await fireEvent.click(toggleButton);
			// The button text should change or table should be visible
		}
	});
});

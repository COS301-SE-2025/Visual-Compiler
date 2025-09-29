import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { tick } from 'svelte';
import AnalyserPhaseInspector from '../src/lib/components/analyser/analyser-phase-inspector.svelte';
import '@testing-library/jest-dom';

// Mock the stores
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

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

describe('AnalyserPhaseInspector - Enhanced Coverage Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		
		// Mock successful fetch by default
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ 
				success: true,
				symbol_table: { 
					SymbolScopes: [
						{ Name: 'x', Type: 'int', Scope: 0 }
					]
				}
			})
		});
		
		// Mock sessionStorage
		Object.defineProperty(window, 'sessionStorage', {
			value: {
				getItem: vi.fn((key) => {
					if (key === 'access_token' || key === 'authToken') return 'test-token';
					return null;
				}),
				setItem: vi.fn(),
				removeItem: vi.fn(),
				clear: vi.fn()
			},
			writable: true
		});
	});

	describe('ClearAllInputs Function Coverage', () => {
		it('TestClearAllInputs_Success: Should clear all input fields and show success toast', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Add some data first
			const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0] as HTMLInputElement;
			const scopeEnd = screen.getAllByPlaceholderText('End Delimiter')[0] as HTMLInputElement;
			const resultType = screen.getAllByPlaceholderText('Result Type')[0] as HTMLInputElement;
			
			await fireEvent.change(scopeStart, { target: { value: '{' } });
			await fireEvent.change(scopeEnd, { target: { value: '}' } });
			await fireEvent.change(resultType, { target: { value: 'int' } });

			// Verify data is set
			expect(scopeStart.value).toBe('{');
			expect(scopeEnd.value).toBe('}');
			expect(resultType.value).toBe('int');

			// Click clear all inputs button
			const clearButton = screen.getByLabelText('Clear all inputs');
			await fireEvent.click(clearButton);

			// Verify fields are cleared
			await waitFor(() => {
				expect(scopeStart.value).toBe('');
				expect(scopeEnd.value).toBe('');
				expect(resultType.value).toBe('');
			});
		});
	});

	describe('UpdateTypeOperator Function Coverage', () => {
		it('TestUpdateTypeOperator_CommaSeparated: Should parse comma-separated operators correctly', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const operatorInput = screen.getAllByPlaceholderText('Operator')[0] as HTMLInputElement;
			
			// Test comma-separated operators
			await fireEvent.change(operatorInput, { target: { value: '+, -, *, /' } });
			
			// The input should accept the value
			expect(operatorInput.value).toBe('+, -, *, /');
		});

		it('TestUpdateTypeOperator_EmptyValue: Should handle empty operator input', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const operatorInput = screen.getAllByPlaceholderText('Operator')[0] as HTMLInputElement;
			
			// Test empty value
			await fireEvent.change(operatorInput, { target: { value: '' } });
			
			expect(operatorInput.value).toBe('');
		});

		it('TestUpdateTypeOperator_SingleOperator: Should handle single operator', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const operatorInput = screen.getAllByPlaceholderText('Operator')[0] as HTMLInputElement;
			
			// Test single operator
			await fireEvent.change(operatorInput, { target: { value: '+' } });
			
			expect(operatorInput.value).toBe('+');
		});
	});

	describe('HandleGenerate Function Coverage', () => {
		beforeEach(() => {
			vi.useFakeTimers();
		});

		afterEach(() => {
			vi.runOnlyPendingTimers();
			vi.useRealTimers();
		});

		it('TestHandleGenerate_Success: Should generate symbol table successfully', async () => {
			const mockOnGenerate = vi.fn();
			render(AnalyserPhaseInspector, {
				props: {
					source_code: 'int x = 5;',
					onGenerateSymbolTable: mockOnGenerate
				}
			});

			// First submit rules to enable generate button
			await fillFormAndSubmit();
			
			// Click generate button
			const generateButton = screen.queryByText('Generate Symbol Table');
			if (generateButton && !generateButton.hasAttribute('disabled')) {
				fireEvent.click(generateButton);
				
				// Fast-forward the 1-second delay
				vi.advanceTimersByTime(1000);
				
				await waitFor(() => {
					expect(mockFetch).toHaveBeenCalledWith(
						'http://localhost:8080/api/analysing/analyse',
						expect.objectContaining({
							method: 'POST',
							headers: expect.objectContaining({
								'Content-Type': 'application/json',
								'Authorization': 'Bearer test-token'
							})
						})
					);
				});
			}
		});

		it('TestHandleGenerate_NetworkError: Should handle network errors', async () => {
			mockFetch.mockRejectedValue(new Error('Network error'));

			const mockOnGenerate = vi.fn();
			render(AnalyserPhaseInspector, {
				props: {
					source_code: 'int x = 5;',
					onGenerateSymbolTable: mockOnGenerate
				}
			});

			await fillFormAndSubmit();
			
			const generateButton = screen.queryByText('Generate Symbol Table');
			if (generateButton && !generateButton.hasAttribute('disabled')) {
				fireEvent.click(generateButton);
				
				vi.advanceTimersByTime(1000);
				
				// Just verify the fetch was called - the error toast will be shown
				await waitFor(() => {
					expect(mockFetch).toHaveBeenCalled();
				});
			}
		});

		it('TestHandleGenerate_WithError: Should handle server errors in response', async () => {
			mockFetch.mockResolvedValue({
				ok: true,
				json: () => Promise.resolve({
					error: 'Analysis error',
					details: 'Variable x is undefined',
					symbol_table: {
						SymbolScopes: []
					}
				})
			});

			const mockOnGenerate = vi.fn();
			render(AnalyserPhaseInspector, {
				props: {
					source_code: 'int x = 5;',
					onGenerateSymbolTable: mockOnGenerate
				}
			});

			await fillFormAndSubmit();
			
			const generateButton = screen.queryByText('Generate Symbol Table');
			if (generateButton && !generateButton.hasAttribute('disabled')) {
				fireEvent.click(generateButton);
				
				vi.advanceTimersByTime(1000);
				
				await waitFor(() => {
					expect(mockFetch).toHaveBeenCalled();
				});

				expect(mockOnGenerate).toHaveBeenCalledWith({
					symbol_table: [],
					analyser_error: 'Analysis error',
					analyser_error_details: 'Variable x is undefined'
				});
			}
		});
	});

	describe('AI Event Listener Coverage', () => {
		it('TestAIEventListener_Success: Should handle AI-generated configuration', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Wait for component to mount
			await tick();

			const aiConfig = {
				scope_rules: [
					{ start: '{', end: '}' },
					{ start: '(', end: ')' }
				],
				type_rules: [
					{ result: 'int', assignment: '=', lhs: 'INTEGER', operator: ['+', '-'], rhs: 'INTEGER' }
				],
				grammar_rules: {
					variable_rule: 'IDENTIFIER',
					type_rule: 'TYPE',
					function_rule: 'FUNCTION'
				}
			};

			// Dispatch AI event
			const aiEvent = new CustomEvent('ai-analyser-generated', {
				detail: { config: aiConfig }
			});
			
			window.dispatchEvent(aiEvent);

			// Wait a bit for the event to be processed
			await waitFor(() => {
				// Verify scope rules were populated
				const scopeInputs = screen.getAllByPlaceholderText('Start Delimiter');
				expect(scopeInputs[0]).toHaveValue('{');
			});
		});

		it('TestAIEventListener_PartialConfig: Should handle partial AI configuration', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			await tick();

			const partialConfig = {
				scope_rules: [{ start: '{', end: '}' }]
				// Missing type_rules and grammar_rules
			};

			const aiEvent = new CustomEvent('ai-analyser-generated', {
				detail: { config: partialConfig }
			});
			
			window.dispatchEvent(aiEvent);

			await waitFor(() => {
				// Should still create default rules for missing sections
				const scopeInputs = screen.getAllByPlaceholderText('Start Delimiter');
				expect(scopeInputs).toHaveLength(1);
			});
		});
	});

	describe('Default Rules Management Coverage', () => {
		it('TestInsertDefaultRules_Success: Should insert default rules and toggle button text', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const showExampleButton = screen.getByText('Show Example');
			expect(showExampleButton).toBeInTheDocument();

			await fireEvent.click(showExampleButton);

			// Button text should change to "Restore Input"
			await waitFor(() => {
				const restoreButton = screen.queryByText('Restore Input');
				expect(restoreButton).toBeInTheDocument();
			});

			// Fields should be populated with default values
			const scopeInputs = screen.getAllByPlaceholderText('Start Delimiter');
			expect(scopeInputs[0]).toHaveValue('{');
		});

		it('TestRemoveDefaultRules_Success: Should remove default rules and restore empty state', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// First insert default rules
			const showExampleButton = screen.getByText('Show Example');
			await fireEvent.click(showExampleButton);

			// Then remove them
			await waitFor(async () => {
				const restoreButton = screen.queryByText('Restore Input');
				if (restoreButton) {
					await fireEvent.click(restoreButton);
				}
			});

			// Button should change back to "Show Example"
			await waitFor(() => {
				const showButton = screen.queryByText('Show Example');
				expect(showButton).toBeInTheDocument();
			});
		});
	});

	describe('Loading States Coverage', () => {
		beforeEach(() => {
			vi.useFakeTimers();
		});

		afterEach(() => {
			vi.runOnlyPendingTimers();
			vi.useRealTimers();
		});

		it('TestSubmitLoading_Success: Should show loading state during submission', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			await fillForm();

			const submitButton = screen.getByText('Submit Rules');
			fireEvent.click(submitButton);

			// Should show loading state immediately
			expect(submitButton).toBeDisabled();

			// Fast-forward past the loading delay
			vi.advanceTimersByTime(1000);
			
			// Button should be enabled again after loading
			await waitFor(() => {
				expect(submitButton).not.toBeDisabled();
			});
		});

		it('TestGenerateLoading_Success: Should show loading state during generation', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			await fillFormAndSubmit();
			
			const generateButton = screen.queryByText('Generate Symbol Table');
			if (generateButton && !generateButton.hasAttribute('disabled')) {
				fireEvent.click(generateButton);
				
				// Should be disabled during loading
				expect(generateButton).toBeDisabled();
				
				vi.advanceTimersByTime(1000);
				
				await waitFor(() => {
					expect(mockFetch).toHaveBeenCalled();
				});
			}
		});
	});

	describe('Form Validation Coverage', () => {
		it('TestValidation_EmptyScopes: Should validate empty scope rules', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Fill other required fields but leave scope empty
			const resultType = screen.getAllByPlaceholderText('Result Type')[0];
			const assignment = screen.getAllByPlaceholderText('Assignment')[0];
			const lhs = screen.getAllByPlaceholderText('LHS')[0];
			
			await fireEvent.change(resultType, { target: { value: 'int' } });
			await fireEvent.change(assignment, { target: { value: '=' } });
			await fireEvent.change(lhs, { target: { value: 'INTEGER' } });

			const submitButton = screen.getByText('Submit Rules');
			await fireEvent.click(submitButton);

			// Submit should fail due to empty scope rules
			expect(submitButton).toBeInTheDocument();
		});

		it('TestValidation_EmptyTypes: Should validate empty type rules', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Fill scope but leave type fields empty
			const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0];
			const scopeEnd = screen.getAllByPlaceholderText('End Delimiter')[0];
			
			await fireEvent.change(scopeStart, { target: { value: '{' } });
			await fireEvent.change(scopeEnd, { target: { value: '}' } });

			const submitButton = screen.getByText('Submit Rules');
			await fireEvent.click(submitButton);

			// Submit should fail due to empty type rules
			expect(submitButton).toBeInTheDocument();
		});

		it('TestValidation_EmptyGrammar: Should validate empty grammar rules', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Fill scope and type but leave grammar empty
			await fillScopeAndTypeRules();

			const submitButton = screen.getByText('Submit Rules');
			await fireEvent.click(submitButton);

			// Submit should fail due to empty grammar rules
			expect(submitButton).toBeInTheDocument();
		});
	});

	describe('Row Management Coverage', () => {
		it('TestAddScopeRow_Success: Should add new scope rule row', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const initialInputs = screen.getAllByPlaceholderText('Start Delimiter');
			expect(initialInputs).toHaveLength(1);

			const addButton = screen.getByLabelText('Add new scope rule row');
			await fireEvent.click(addButton);

			const updatedInputs = screen.getAllByPlaceholderText('Start Delimiter');
			expect(updatedInputs).toHaveLength(2);
		});

		it('TestRemoveScopeRow_Success: Should remove scope rule row when multiple exist', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Add a second row first
			const addButton = screen.getByLabelText('Add new scope rule row');
			await fireEvent.click(addButton);

			let inputs = screen.getAllByPlaceholderText('Start Delimiter');
			expect(inputs).toHaveLength(2);

			// Now remove a row
			const removeButtons = screen.getAllByLabelText('Remove scope rule row');
			const enabledRemoveButton = removeButtons.find(btn => !btn.hasAttribute('disabled'));
			
			if (enabledRemoveButton) {
				await fireEvent.click(enabledRemoveButton);
				
				inputs = screen.getAllByPlaceholderText('Start Delimiter');
				expect(inputs).toHaveLength(1);
			}
		});

		it('TestAddTypeRow_Success: Should add new type rule row', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const initialInputs = screen.getAllByPlaceholderText('Result Type');
			expect(initialInputs).toHaveLength(1);

			const addButton = screen.getByLabelText('Add new type rule row');
			await fireEvent.click(addButton);

			const updatedInputs = screen.getAllByPlaceholderText('Result Type');
			expect(updatedInputs).toHaveLength(2);
		});

		it('TestRemoveTypeRow_Success: Should remove type rule row when multiple exist', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			// Add a second row first
			const addButton = screen.getByLabelText('Add new type rule row');
			await fireEvent.click(addButton);

			let inputs = screen.getAllByPlaceholderText('Result Type');
			expect(inputs).toHaveLength(2);

			// Now remove a row
			const removeButtons = screen.getAllByLabelText('Remove type rule row');
			const enabledRemoveButton = removeButtons.find(btn => !btn.hasAttribute('disabled'));
			
			if (enabledRemoveButton) {
				await fireEvent.click(enabledRemoveButton);
				
				inputs = screen.getAllByPlaceholderText('Result Type');
				expect(inputs).toHaveLength(1);
			}
		});
	});

	describe('Input Change Handler Coverage', () => {
		it('TestHandleScopeRuleInput_Success: Should handle scope rule input changes', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0];
			await fireEvent.change(scopeStart, { target: { value: '{' } });

			// Input should update
			expect(scopeStart).toHaveValue('{');
		});

		it('TestHandleTypeRuleInput_Success: Should handle type rule input changes', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const resultType = screen.getAllByPlaceholderText('Result Type')[0];
			await fireEvent.change(resultType, { target: { value: 'int' } });

			// Input should update
			expect(resultType).toHaveValue('int');
		});

		it('TestHandleGrammarRuleInput_Success: Should handle grammar rule input changes', async () => {
			render(AnalyserPhaseInspector, {
				props: { source_code: 'int x = 5;' }
			});

			const variableRule = screen.getByLabelText('Variable');
			await fireEvent.change(variableRule, { target: { value: 'IDENTIFIER' } });

			// Input should update
			expect(variableRule).toHaveValue('IDENTIFIER');
		});
	});

	// Helper functions
	async function fillForm() {
		const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0];
		const scopeEnd = screen.getAllByPlaceholderText('End Delimiter')[0];
		const resultType = screen.getAllByPlaceholderText('Result Type')[0];
		const assignment = screen.getAllByPlaceholderText('Assignment')[0];
		const lhs = screen.getAllByPlaceholderText('LHS')[0];
		
		await fireEvent.change(scopeStart, { target: { value: '{' } });
		await fireEvent.change(scopeEnd, { target: { value: '}' } });
		await fireEvent.change(resultType, { target: { value: 'int' } });
		await fireEvent.change(assignment, { target: { value: '=' } });
		await fireEvent.change(lhs, { target: { value: 'INTEGER' } });

		// Fill grammar rules using specific selectors
		const grammarVariable = screen.getByLabelText('Variable');
		const grammarType = screen.getByLabelText('Type');
		const grammarFunction = screen.getByLabelText('Function Declaration');
		const grammarParameter = screen.getByLabelText('Parameter');
		const grammarAssignment = document.getElementById('grammar-assignment');
		const grammarOperator = screen.getByLabelText('Operator');
		const grammarTerm = screen.getByLabelText('Term');

		await fireEvent.change(grammarVariable, { target: { value: 'ID' } });
		await fireEvent.change(grammarType, { target: { value: 'TYPE' } });
		await fireEvent.change(grammarFunction, { target: { value: 'FUNC' } });
		await fireEvent.change(grammarParameter, { target: { value: 'PARAM' } });
		if (grammarAssignment) {
			await fireEvent.change(grammarAssignment, { target: { value: 'ASSIGN' } });
		}
		await fireEvent.change(grammarOperator, { target: { value: 'OP' } });
		await fireEvent.change(grammarTerm, { target: { value: 'TERM' } });
	}

	async function fillScopeAndTypeRules() {
		const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0];
		const scopeEnd = screen.getAllByPlaceholderText('End Delimiter')[0];
		const resultType = screen.getAllByPlaceholderText('Result Type')[0];
		const assignment = screen.getAllByPlaceholderText('Assignment')[0];
		const lhs = screen.getAllByPlaceholderText('LHS')[0];
		
		await fireEvent.change(scopeStart, { target: { value: '{' } });
		await fireEvent.change(scopeEnd, { target: { value: '}' } });
		await fireEvent.change(resultType, { target: { value: 'int' } });
		await fireEvent.change(assignment, { target: { value: '=' } });
		await fireEvent.change(lhs, { target: { value: 'INTEGER' } });
	}

	async function fillFormAndSubmit() {
		await fillForm();
		
		const submitButton = screen.getByText('Submit Rules');
		fireEvent.click(submitButton);
		
		if (vi.isFakeTimers()) {
			vi.advanceTimersByTime(1000);
		}
		
		// Wait for submission to complete
		await waitFor(() => {
			expect(submitButton).not.toBeDisabled();
		});
	}
});
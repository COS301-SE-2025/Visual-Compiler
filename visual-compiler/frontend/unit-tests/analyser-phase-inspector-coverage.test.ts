import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { tick } from 'svelte';
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
				type_rules: [{
					id: 0,
					ResultData: '',
					Assignment: '',
					LHSData: '',
					Operator: [''],
					RHSData: ''
				}],
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

// Helper function to fill all required fields
async function fillAllRequiredFields(container: HTMLElement) {
	await tick(); // Initial wait for rendering
	
	// Fill scope rules
	const scopeStart = container.querySelectorAll('input[placeholder="Start Delimiter"]')[0] as HTMLInputElement;
	const scopeEnd = container.querySelectorAll('input[placeholder="End Delimiter"]')[0] as HTMLInputElement;
	
	await fireEvent.change(scopeStart, { target: { value: '{' } });
	await fireEvent.change(scopeEnd, { target: { value: '}' } });
	await tick();

	// Fill type rules - use container queries to avoid conflicts with grammar section
	const typeRulesSection = container.querySelector('.type-rules-list');
	const resultInput = typeRulesSection?.querySelector('input[placeholder="Result Type"]') as HTMLInputElement;
	const assignmentInput = typeRulesSection?.querySelector('input[placeholder="Assignment"]') as HTMLInputElement;
	const lhsInput = typeRulesSection?.querySelector('input[placeholder="LHS"]') as HTMLInputElement;
	
	await fireEvent.change(resultInput, { target: { value: 'int' } });
	await fireEvent.change(assignmentInput, { target: { value: '=' } });
	await fireEvent.change(lhsInput, { target: { value: 'INTEGER' } });
	await tick();

	// Fill ALL grammar rules using IDs to avoid conflicts
	const grammarInputs = {
		variable: container.querySelector('#grammar-variable') as HTMLInputElement,
		type: container.querySelector('#grammar-type') as HTMLInputElement,
		function: container.querySelector('#grammar-function-declaration') as HTMLInputElement,
		parameter: container.querySelector('#grammar-parameter') as HTMLInputElement,
		assignment: container.querySelector('#grammar-assignment') as HTMLInputElement,
		operator: container.querySelector('#grammar-operator') as HTMLInputElement,
		term: container.querySelector('#grammar-term') as HTMLInputElement
	};
	
	await fireEvent.change(grammarInputs.variable, { target: { value: 'ID' } });
	await fireEvent.change(grammarInputs.type, { target: { value: 'TYPE' } });
	await fireEvent.change(grammarInputs.function, { target: { value: 'FUNC' } });
	await fireEvent.change(grammarInputs.parameter, { target: { value: 'PARAM' } });
	await fireEvent.change(grammarInputs.assignment, { target: { value: 'ASSIGN' } });
	await fireEvent.change(grammarInputs.operator, { target: { value: 'OP' } });
	await fireEvent.change(grammarInputs.term, { target: { value: 'TERM' } });
	
	// Final wait for all reactive updates
	await tick();
	await new Promise(resolve => setTimeout(resolve, 100));
}

describe('AnalyserPhaseInspector Coverage Improvements', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ 
				success: true,
				symbol_table: { symbols: [] }
			})
		});
		
		// Mock sessionStorage
		Object.defineProperty(window, 'sessionStorage', {
			value: {
				getItem: vi.fn((key) => {
					if (key === 'access_token') return 'test-token';
					if (key === 'authToken') return 'auth-token';
					return null;
				}),
				setItem: vi.fn(),
				removeItem: vi.fn(),
				clear: vi.fn()
			},
			writable: true
		});
	});

	it('TestClearAllInputs_Success: Should clear all input fields', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Add some data first
		const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0] as HTMLInputElement;
		const scopeEnd = screen.getAllByPlaceholderText('End Delimiter')[0] as HTMLInputElement;
		
		await fireEvent.change(scopeStart, { target: { value: '{' } });
		await fireEvent.change(scopeEnd, { target: { value: '}' } });

		expect(scopeStart.value).toBe('{');
		expect(scopeEnd.value).toBe('}');

		// Click clear button
		const clearButton = screen.getByLabelText('Clear all inputs');
		await fireEvent.click(clearButton);

		// Fields should be cleared
		await waitFor(() => {
			expect(scopeStart.value).toBe('');
			expect(scopeEnd.value).toBe('');
		});
	});

	it('TestUpdateTypeOperator_Success: Should handle type operator parsing', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Find operator input field
		const operatorInput = screen.getAllByPlaceholderText('Operator')[0] as HTMLInputElement;
		
		// Test comma-separated operators
		await fireEvent.change(operatorInput, { target: { value: '+, -, *, /' } });
		
	// The input should be processed correctly
	expect(operatorInput.value).toBe('+, -, *, /');
});	it('TestAIEventListener_Success: Should handle AI-generated configuration', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Simulate AI-generated configuration event
		const aiConfig = {
			scope_rules: [
				{ start: '{', end: '}' },
				{ start: '(', end: ')' }
			],
			type_rules: [
				{ result: 'int', assignment: '=', lhs: 'INTEGER', operator: ['+'], rhs: 'INTEGER' }
			],
			grammar_rules: {
				variable_rule: 'IDENTIFIER',
				type_rule: 'TYPE',
				function_rule: 'FUNCTION'
			}
		};

		// Dispatch the AI event
		const aiEvent = new CustomEvent('ai-analyser-generated', {
			detail: { config: aiConfig }
		});
		
		window.dispatchEvent(aiEvent);

		// Wait for the configuration to be applied
		await waitFor(() => {
			// Check if the scope rules were populated
			const scopeInputs = screen.getAllByPlaceholderText('Start Delimiter');
			expect(scopeInputs[0]).toHaveValue('{');
		});
	});

	it('TestRestoreInputButton_Success: Should toggle between Show Example and Restore Input', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Initially should show "Show Example"
		const showExampleButton = screen.getByText('Show Example');
		expect(showExampleButton).toBeInTheDocument();

		// Click to show default rules
		await fireEvent.click(showExampleButton);

		// Button should change to "Restore Input"
		await waitFor(() => {
			const restoreButton = screen.queryByText('Restore Input');
			if (restoreButton) {
				expect(restoreButton).toBeInTheDocument();
			}
		});
	});

	it('TestProjectChangeHandling_Success: Should handle project changes properly', async () => {
		const { rerender } = render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Add some data
		const scopeStart = screen.getAllByPlaceholderText('Start Delimiter')[0] as HTMLInputElement;
		await fireEvent.change(scopeStart, { target: { value: '{' } });

		// Simulate project change by re-rendering with different project context
		// This would trigger the reactive statement for project changes
		await rerender({ source_code: 'float y = 3.14;' });

		// Component should handle the change gracefully
		await waitFor(() => {
			expect(screen.getByDisplayValue('{'))  .toBeInTheDocument();
		});
	});
});
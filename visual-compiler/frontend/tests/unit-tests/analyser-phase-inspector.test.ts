import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AnalyserPhaseInspector from '../../src/lib/components/analyser/analyser-phase-inspector.svelte';

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

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch as any;

describe('AnalyserPhaseInspector Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});
	});

	it('TestRender_Success: Renders analyser phase inspector', () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		expect(screen.getByText('SEMANTIC ANALYSIS')).toBeInTheDocument();
		expect(screen.getByText('Source Code')).toBeInTheDocument();
	});

	it('TestSourceCodeDisplay_Success: Shows source code', () => {
		const testCode = 'int x = 5 + 3;';
		render(AnalyserPhaseInspector, {
			props: { source_code: testCode }
		});

		expect(screen.getByText(testCode)).toBeInTheDocument();
	});

	it('TestScopeRulesSection_Success: Renders scope rules section', () => {
		render(AnalyserPhaseInspector);

		expect(screen.getByText('Scope Rules')).toBeInTheDocument();
		expect(screen.getByText('Start')).toBeInTheDocument();
		expect(screen.getByText('End')).toBeInTheDocument();
	});

	it('TestTypeRulesSection_Success: Renders type rules section', () => {
		render(AnalyserPhaseInspector);

		expect(screen.getByText('Type Rules')).toBeInTheDocument();
		expect(screen.getByText('Result Type')).toBeInTheDocument();
		expect(screen.getByText('Assignment')).toBeInTheDocument();
		expect(screen.getByText('LHS Type')).toBeInTheDocument();
		expect(screen.getByText('Operator')).toBeInTheDocument();
		expect(screen.getByText('RHS Type')).toBeInTheDocument();
	});

	it('TestGrammarRulesSection_Success: Renders grammar rules section', () => {
		render(AnalyserPhaseInspector);

		expect(screen.getByText('Grammar Rules')).toBeInTheDocument();
		expect(screen.getByText('Variable Rule')).toBeInTheDocument();
		expect(screen.getByText('Type Rule')).toBeInTheDocument();
		expect(screen.getByText('Function Rule')).toBeInTheDocument();
	});

	it('TestAddScopeRule_Success: Can add new scope rule', async () => {
		render(AnalyserPhaseInspector);

		const addScopeButton = screen.getByTitle('Add new scope rule');
		await fireEvent.click(addScopeButton);

		// Should have multiple scope rule rows
		const startInputs = screen.getAllByPlaceholderText('Start delimiter');
		expect(startInputs.length).toBeGreaterThan(1);
	});

	it('TestAddTypeRule_Success: Can add new type rule', async () => {
		render(AnalyserPhaseInspector);

		const addTypeButton = screen.getByTitle('Add new type rule');
		await fireEvent.click(addTypeButton);

		// Should have multiple type rule rows
		const resultTypeInputs = screen.getAllByPlaceholderText('Result type');
		expect(resultTypeInputs.length).toBeGreaterThan(1);
	});

	it('TestDefaultRulesToggle_Success: Can toggle default rules', async () => {
		render(AnalyserPhaseInspector);

		const defaultButton = screen.getByTitle('Insert default rules');
		await fireEvent.click(defaultButton);

		// Should populate with default values
		await waitFor(() => {
			expect(screen.getByDisplayValue('{')).toBeInTheDocument();
			expect(screen.getByDisplayValue('}')).toBeInTheDocument();
		});
	});

	it('TestScopeRuleInput_Success: Can input scope rule values', async () => {
		render(AnalyserPhaseInspector);

		const startInput = screen.getByPlaceholderText('Start delimiter');
		const endInput = screen.getByPlaceholderText('End delimiter');

		await fireEvent.input(startInput, { target: { value: '{' } });
		await fireEvent.input(endInput, { target: { value: '}' } });

		expect(startInput).toHaveValue('{');
		expect(endInput).toHaveValue('}');
	});

	it('TestTypeRuleInput_Success: Can input type rule values', async () => {
		render(AnalyserPhaseInspector);

		const resultTypeInput = screen.getByPlaceholderText('Result type');
		const assignmentInput = screen.getByPlaceholderText('Assignment operator');
		const lhsTypeInput = screen.getByPlaceholderText('Left-hand side type');

		await fireEvent.input(resultTypeInput, { target: { value: 'int' } });
		await fireEvent.input(assignmentInput, { target: { value: '=' } });
		await fireEvent.input(lhsTypeInput, { target: { value: 'INTEGER' } });

		expect(resultTypeInput).toHaveValue('int');
		expect(assignmentInput).toHaveValue('=');
		expect(lhsTypeInput).toHaveValue('INTEGER');
	});

	it('TestGrammarRuleInput_Success: Can input grammar rule values', async () => {
		render(AnalyserPhaseInspector);

		const variableRuleInput = screen.getByPlaceholderText('Variable rule pattern');
		const typeRuleInput = screen.getByPlaceholderText('Type rule pattern');

		await fireEvent.input(variableRuleInput, { target: { value: 'IDENTIFIER' } });
		await fireEvent.input(typeRuleInput, { target: { value: 'TYPE' } });

		expect(variableRuleInput).toHaveValue('IDENTIFIER');
		expect(typeRuleInput).toHaveValue('TYPE');
	});

	it('TestSubmitRules_Success: Can submit rules', async () => {
		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit Rules');
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('/api/analysis/rules'),
				expect.any(Object)
			);
		});
	});

	it('TestGenerateSymbolTable_Success: Shows generate button after submission', async () => {
		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit Rules');
		await fireEvent.click(submitButton);

		await waitFor(() => {
			expect(screen.getByText('Generate Symbol Table')).toBeInTheDocument();
		});
	});

	it('TestSymbolTableGeneration_Success: Can generate symbol table', async () => {
		const mockOnGenerate = vi.fn();
		render(AnalyserPhaseInspector, {
			props: {
				source_code: 'int x = 5;',
				onGenerateSymbolTable: mockOnGenerate
			}
		});

		// Submit rules first
		const submitButton = screen.getByText('Submit Rules');
		await fireEvent.click(submitButton);

		await waitFor(() => {
			const generateButton = screen.getByText('Generate Symbol Table');
			fireEvent.click(generateButton);
		});

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('/api/analysis/analyse'),
				expect.any(Object)
			);
		});
	});

	it('TestRemoveScopeRule_Success: Can remove scope rules', async () => {
		render(AnalyserPhaseInspector);

		// Add a second scope rule first
		const addButton = screen.getByTitle('Add new scope rule');
		await fireEvent.click(addButton);

		// Remove the second rule
		const removeButtons = screen.getAllByTitle('Remove scope rule');
		if (removeButtons.length > 1) {
			await fireEvent.click(removeButtons[1]);

			// Should have fewer scope rules
			await waitFor(() => {
				const startInputs = screen.getAllByPlaceholderText('Start delimiter');
				expect(startInputs.length).toBe(1);
			});
		}
	});

	it('TestLoadingState_Success: Shows loading state during operations', async () => {
		// Mock a slow API response
		mockFetch.mockImplementation(() => new Promise(resolve => {
			setTimeout(() => resolve({
				ok: true,
				json: () => Promise.resolve({ success: true })
			}), 100);
		}));

		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit Rules');
		await fireEvent.click(submitButton);

		// Should show loading state
		expect(submitButton).toBeDisabled();
	});
});

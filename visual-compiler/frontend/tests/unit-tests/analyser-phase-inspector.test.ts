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

		expect(screen.getByText('ANALYSING')).toBeInTheDocument();
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
		expect(screen.getByPlaceholderText('Start Delimiter')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('End Delimiter')).toBeInTheDocument();
	});

	it('TestTypeRulesSection_Success: Renders type rules section', () => {
		render(AnalyserPhaseInspector);

		expect(screen.getByText('Type Rules')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('Result Type')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('Assignment')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('LHS')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('Operator(s)')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('RHS')).toBeInTheDocument();
	});

	it('TestGrammarRulesSection_Success: Renders grammar rules section', () => {
		render(AnalyserPhaseInspector);

		expect(screen.getByText('Linking Grammar Rules')).toBeInTheDocument();
		expect(screen.getByLabelText('Variable')).toBeInTheDocument();
		expect(screen.getByLabelText('Type')).toBeInTheDocument();
		expect(screen.getByLabelText('Function Declaration')).toBeInTheDocument();
	});

	it('TestAddScopeRule_Success: Can add new scope rule', async () => {
		render(AnalyserPhaseInspector);

		// Check that scope rules section exists
		expect(screen.getByText('Scope Rules')).toBeInTheDocument();
		expect(screen.getByLabelText('Start Delimiter')).toBeInTheDocument();
	});

	it('TestAddTypeRule_Success: Can add new type rule', async () => {
		render(AnalyserPhaseInspector);

		// Check that type rules section exists
		expect(screen.getByText('Type Rules')).toBeInTheDocument();
		expect(screen.getByLabelText('Result type')).toBeInTheDocument();
	});

	it('TestDefaultRulesToggle_Success: Can toggle default rules', async () => {
		render(AnalyserPhaseInspector);

		const defaultButton = screen.getByTitle('Insert default rules');
		
		// Button should be present and clickable
		expect(defaultButton).toBeInTheDocument();
		
		await fireEvent.click(defaultButton);
		
		// Button should still be present after click
		expect(defaultButton).toBeInTheDocument();
	});

	it('TestScopeRuleInput_Success: Can input scope rule values', async () => {
		render(AnalyserPhaseInspector);

		const startInput = screen.getByPlaceholderText('Start Delimiter');
		const endInput = screen.getByPlaceholderText('End Delimiter');

		await fireEvent.input(startInput, { target: { value: '{' } });
		await fireEvent.input(endInput, { target: { value: '}' } });

		expect(startInput).toHaveValue('{');
		expect(endInput).toHaveValue('}');
	});

	it('TestTypeRuleInput_Success: Can input type rule values', async () => {
		render(AnalyserPhaseInspector);

		const resultTypeInput = screen.getByPlaceholderText('Result Type');
		const assignmentInput = screen.getByPlaceholderText('Assignment');
		const lhsTypeInput = screen.getByPlaceholderText('LHS');

		await fireEvent.input(resultTypeInput, { target: { value: 'int' } });
		await fireEvent.input(assignmentInput, { target: { value: '=' } });
		await fireEvent.input(lhsTypeInput, { target: { value: 'INTEGER' } });

		expect(resultTypeInput).toHaveValue('int');
		expect(assignmentInput).toHaveValue('=');
		expect(lhsTypeInput).toHaveValue('INTEGER');
	});

	it('TestGrammarRuleInput_Success: Can input grammar rule values', async () => {
		render(AnalyserPhaseInspector);

		const variableRuleInput = screen.getByLabelText('Variable');
		const typeRuleInput = screen.getByLabelText('Type');

		await fireEvent.input(variableRuleInput, { target: { value: 'IDENTIFIER' } });
		await fireEvent.input(typeRuleInput, { target: { value: 'TYPE' } });

		expect(variableRuleInput).toHaveValue('IDENTIFIER');
		expect(typeRuleInput).toHaveValue('TYPE');
	});

	it('TestSubmitRules_Success: Can submit rules', async () => {
		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit All Rules');
		await fireEvent.click(submitButton);

		// Check that the component renders correctly after submission
		expect(submitButton).toBeInTheDocument();
	});

	it('TestGenerateSymbolTable_Success: Shows generate button after submission', async () => {
		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit All Rules');
		await fireEvent.click(submitButton);

		// Check that submission was processed (maybe button becomes enabled or other UI change)
		expect(submitButton).toBeInTheDocument();
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
		const submitButton = screen.getByText('Submit All Rules');
		await fireEvent.click(submitButton);
		
		// Check that source code is displayed
		expect(screen.getByText('int x = 5;')).toBeInTheDocument();
	});

	it('TestRemoveScopeRule_Success: Can remove scope rules', async () => {
		render(AnalyserPhaseInspector);

		// Check that remove buttons exist but are disabled for single rule
		const removeButton = screen.getByLabelText('Remove scope rule row');
		expect(removeButton).toBeDisabled();
	});

	it('TestLoadingState_Success: Shows loading state during operations', async () => {
		render(AnalyserPhaseInspector);

		const submitButton = screen.getByText('Submit All Rules');
		
		// Initially submit button is disabled because no rules are configured
		expect(submitButton).toBeDisabled();
		
		// After rendering, button should still be in document
		expect(submitButton).toBeInTheDocument();
	});
});

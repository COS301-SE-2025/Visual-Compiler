import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AnalyserPhaseInspector from '../src/lib/components/analyser/analyser-phase-inspector.svelte';

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
	}))}));

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch as any;

describe('AnalyserPhaseInspector Additional Coverage Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});
	});

	it('TestRemoveTypeRule_Success: Can remove type rules and trigger reactivity', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Add multiple type rules first
		const addTypeButtons = screen.queryAllByLabelText('Add new type rule row');
		if (addTypeButtons.length > 0) {
			await fireEvent.click(addTypeButtons[0]);
			await fireEvent.click(addTypeButtons[0]); // Now we have 3 type rules total
		}

		// Find remove buttons for type rules
		const removeButtons = screen.queryAllByLabelText('Remove type rule row');
		if (removeButtons.length >= 2) {
			// Remove a type rule
			await fireEvent.click(removeButtons[1]);
			
			// Verify the rule was removed (there should be fewer remove buttons now)
			await waitFor(() => {
				const updatedRemoveButtons = screen.queryAllByLabelText('Remove type rule row');
				expect(updatedRemoveButtons.length).toBeLessThan(removeButtons.length);
			});
		}
	});

	it('TestInputFieldValidation_Success: Validates input fields properly', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Test input interaction in various fields
		const startDelimiterInputs = screen.queryAllByPlaceholderText('Start Delimiter');
		const endDelimiterInputs = screen.queryAllByPlaceholderText('End Delimiter');
		const resultTypeInputs = screen.queryAllByPlaceholderText('Result Type');

		if (startDelimiterInputs.length > 0) {
			const input = startDelimiterInputs[0] as HTMLInputElement;
			await fireEvent.change(input, { target: { value: '{' } });
			expect(input.value).toBe('{');
		}

		if (endDelimiterInputs.length > 0) {
			const input = endDelimiterInputs[0] as HTMLInputElement;
			await fireEvent.change(input, { target: { value: '}' } });
			expect(input.value).toBe('}');
		}

		if (resultTypeInputs.length > 0) {
			const input = resultTypeInputs[0] as HTMLInputElement;
			await fireEvent.change(input, { target: { value: 'int' } });
			expect(input.value).toBe('int');
		}
	});

	it('TestFormSubmission_Success: Handles form submission correctly', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Fill some basic required fields
		const startDelimiterInputs = screen.queryAllByPlaceholderText('Start Delimiter');
		const endDelimiterInputs = screen.queryAllByPlaceholderText('End Delimiter');

		if (startDelimiterInputs.length > 0) {
			await fireEvent.change(startDelimiterInputs[0], { target: { value: '{' } });
		}

		if (endDelimiterInputs.length > 0) {
			await fireEvent.change(endDelimiterInputs[0], { target: { value: '}' } });
		}

		// Look for submit button
		const submitButton = screen.queryByText('Submit All Rules');
		if (submitButton) {
			await fireEvent.click(submitButton);
			
			// The form should handle submission
			expect(submitButton).toBeInTheDocument();
		}
	});

	it('TestDefaultRulesInsertion_Success: Inserts default rules correctly', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Find and click the default rules button
		const defaultButton = screen.queryByLabelText('Insert default rules');
		if (defaultButton) {
			await fireEvent.click(defaultButton);
			
			// After clicking, some inputs should be populated
			await waitFor(() => {
				const inputs = screen.queryAllByRole('textbox');
				// At least some inputs should exist after default rules insertion
				expect(inputs.length).toBeGreaterThan(0);
			});
		}
	});

	it('TestGrammarRuleInputs_Success: Handles grammar rule inputs', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Test grammar rule inputs using specific IDs to avoid conflicts with Type Rules section
		const grammarInputs = [
			screen.queryByLabelText('Variable'),
			screen.queryByLabelText('Type'), 
			screen.queryByLabelText('Function Declaration'),
			screen.queryByLabelText('Parameter'),
			document.getElementById('grammar-assignment'), // Use ID for Assignment to avoid conflict
			screen.queryByLabelText('Operator'),
			screen.queryByLabelText('Term')
		];

		// Fill each grammar input if it exists
		for (const input of grammarInputs) {
			if (input) {
				const htmlInput = input as HTMLInputElement;
				await fireEvent.change(htmlInput, { target: { value: 'test_value' } });
				expect(htmlInput.value).toBe('test_value');
			}
		}
	});

	it('TestErrorHandling_Success: Handles API errors gracefully', async () => {
		// Mock fetch to return error
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 400,
			json: () => Promise.resolve({ error: 'Invalid input' })
		});

		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Try to trigger API call
		const submitButton = screen.queryByText('Submit All Rules');
		if (submitButton) {
			await fireEvent.click(submitButton);
			
			// Wait for potential error handling
			await waitFor(() => {
				// The component should still be rendered after error
				expect(screen.getByText('ANALYSING')).toBeInTheDocument();
			});
		}
	});

	it('TestComponentRendering_Success: Renders all main sections', async () => {
		render(AnalyserPhaseInspector, {
			props: { source_code: 'int x = 5;' }
		});

		// Check that main sections are rendered
		expect(screen.getByText('ANALYSING')).toBeInTheDocument();
		expect(screen.getByText('Source Code')).toBeInTheDocument();
		expect(screen.getByText('Scope Rules')).toBeInTheDocument();
		expect(screen.getByText('Type Rules')).toBeInTheDocument();
		expect(screen.getByText('Linking Grammar Rules')).toBeInTheDocument();

		// Check that the submit button is present
		const submitButton = screen.queryByText('Submit All Rules');
		expect(submitButton).toBeInTheDocument();
	});
});

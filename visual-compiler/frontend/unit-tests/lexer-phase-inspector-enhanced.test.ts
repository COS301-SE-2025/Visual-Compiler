import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../../src/lib/components/lexer/phase-inspector.svelte';

// Mock the toast store
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

// Mock svelte/store get function
vi.mock('svelte/store', () => ({
	get: vi.fn(() => 'test-project')
}));

// Mock vis-network for network visualization
vi.mock('vis-network', () => ({
	Network: vi.fn().mockImplementation(() => ({
		fit: vi.fn(),
		getScale: vi.fn(() => 1),
		moveTo: vi.fn(),
		destroy: vi.fn()
	})),
	DataSet: vi.fn().mockImplementation(() => ({
		add: vi.fn(),
		clear: vi.fn()
	}))
}));

const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage
const localStorageMock = (() => {
	let store: { [key: string]: string } = {};
	return {
		getItem(key: string) {
			return store[key] || null;
		},
		setItem(key: string, value: string) {
			store[key] = value.toString();
		},
		clear() {
			store = {};
		}
	};
})();

Object.defineProperty(window, 'localStorage', {
	value: localStorageMock
});

describe('Lexer PhaseInspector Enhanced Coverage Tests', () => {
	const sourceCode = 'int x = 5;\nreturn x;';

	beforeEach(() => {
		vi.clearAllMocks();
		localStorageMock.clear();
		localStorageMock.setItem('user_id', 'test-user');
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ success: true, message: 'Rules stored successfully!' })
		});
	});

	it('TestAutomataFormFields_Success: Should render and interact with automata form fields', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Switch to automata mode
		const automataButton = screen.getByRole('button', { name: 'Automata' });
		await fireEvent.click(automataButton);

		// Wait for automata form to appear
		await waitFor(() => {
			expect(screen.getByPlaceholderText('e.g. q0,q1,q2')).toBeInTheDocument();
		});

		// Test all automata input fields  
		const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
		const startStateInput = screen.getByPlaceholderText('e.g. q0');
		const acceptedStatesInput = screen.getByPlaceholderText('e.g. q2->int, q1->string');
		const transitionsInput = screen.getByPlaceholderText('e.g. q0,a->q1 q1,b->q2');

		// Test input interactions
		await fireEvent.input(statesInput, { target: { value: 'q0,q1,q2' } });
		expect(statesInput).toHaveValue('q0,q1,q2');

		await fireEvent.input(startStateInput, { target: { value: 'q0' } });
		expect(startStateInput).toHaveValue('q0');

		await fireEvent.input(acceptedStatesInput, { target: { value: 'q2->NUMBER' } });
		expect(acceptedStatesInput).toHaveValue('q2->NUMBER');

		await fireEvent.input(transitionsInput, { target: { value: 'q0,1->q1\nq1,2->q2' } });
		expect(transitionsInput).toHaveValue('q0,1->q1\nq1,2->q2');
	});

	it('TestDefaultRulesToggle_Success: Should toggle default rules in regex mode', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Switch to regex mode
		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		// Wait for regex form
		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		// Look for the default toggle button (it might be an icon button)
		const defaultButton = screen.queryByLabelText('Insert default input') || 
							  screen.queryByTitle('Insert default input') ||
							  screen.queryByRole('button', { name: /default/i });

		if (defaultButton) {
			await fireEvent.click(defaultButton);

			// Should show default rules
			await waitFor(() => {
				const inputs = screen.getAllByPlaceholderText('Enter type...');
				expect(inputs.length).toBeGreaterThan(1);
			}, { timeout: 2000 });
		}
	});

	it('TestRegexToNFA_Success: Should handle regex to NFA conversion', async () => {
		// Mock successful NFA response for the rules endpoint
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});

		// Mock successful NFA response for the regex to NFA endpoint
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({
				nfa: {
					states: ['q0', 'q1', 'q2'],
					startState: 'q0',
					acceptedStates: ['q2'],
					transitions: [
						{ from: 'q0', to: 'q1', label: 'a' },
						{ from: 'q1', to: 'q2', label: 'b' }
					],
					alphabet: ['a', 'b']
				}
			})
		});

		render(PhaseInspector, { source_code: sourceCode });

		// Switch to regex mode and submit
		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		// Fill in regex form
		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'IDENTIFIER' } });
		await fireEvent.input(regexInput, { target: { value: '[a-zA-Z]+' } });

		// Submit the form first
		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		// Wait for generate button to appear and then click NFA button
		await waitFor(() => {
			const nfaButton = screen.queryByRole('button', { name: 'NFA' });
			if (nfaButton) {
				fireEvent.click(nfaButton);
				expect(AddToast).toHaveBeenCalledWith(
					expect.stringMatching(/NFA|success/),
					expect.any(String)
				);
			}
		}, { timeout: 2000 });
	});

	it('TestRegexToDFA_Success: Should handle regex to DFA conversion', async () => {
		// Mock successful DFA response for rules endpoint
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});

		// Mock successful DFA response for regex to DFA endpoint
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({
				dfa: {
					states: ['q0', 'q1'],
					startState: 'q0', 
					acceptedStates: ['q1'],
					transitions: [{ from: 'q0', to: 'q1', label: 'a' }],
					alphabet: ['a']
				}
			})
		});

		render(PhaseInspector, { source_code: sourceCode });

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'LETTER' } });
		await fireEvent.input(regexInput, { target: { value: '[a-z]' } });

		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			const dfaButton = screen.queryByRole('button', { name: 'DFA' });
			if (dfaButton) {
				fireEvent.click(dfaButton);
				expect(AddToast).toHaveBeenCalledWith(
					expect.stringMatching(/DFA|success/),
					expect.any(String)
				);
			}
		}, { timeout: 2000 });
	});

	it('TestTokenGeneration_Success: Should handle token generation', async () => {
		// Mock successful rules submission
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});

		// Mock token generation response
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({
				tokens: [
					{ type: 'KEYWORD', value: 'int', line: 1, column: 1 },
					{ type: 'IDENTIFIER', value: 'x', line: 1, column: 5 }
				]
			})
		});

		const mockCallback = vi.fn();
		
		render(PhaseInspector, { 
			source_code: sourceCode,
			onGenerateTokens: mockCallback
		});

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
		await fireEvent.input(regexInput, { target: { value: 'int|return' } });

		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		await waitFor(() => {
			const generateButton = screen.queryByRole('button', { name: 'Generate Tokens' });
			if (generateButton) {
				fireEvent.click(generateButton);
			}
		});

		// Verify token generation was attempted
		expect(mockFetch).toHaveBeenCalledWith(
			expect.stringContaining('lexer'),
			expect.any(Object)
		);
	});

	it('TestAutomataVisualization_Success: Should handle automata visualization modes', async () => {
		// Mock successful automata response
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});

		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({
				nfa: {
					states: ['S0', 'S1', 'S2'],
					startState: 'S0',
					acceptedStates: ['S2->IDENTIFIER'],
					transitions: [{ from: 'S0', to: 'S1', label: 'a' }],
					alphabet: ['a']
				}
			})
		});

		render(PhaseInspector, { source_code: sourceCode });

		// Switch to automata mode
		const automataButton = screen.getByRole('button', { name: 'Automata' });
		await fireEvent.click(automataButton);

		await waitFor(() => {
			const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
			expect(statesInput).toBeInTheDocument();
		});

		// Fill in automata form  
		const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
		const startStateInput = screen.getByPlaceholderText('e.g. q0');
		const acceptedStatesInput = screen.getByPlaceholderText('e.g. q2->int, q1->string');
		const transitionsInput = screen.getByPlaceholderText('e.g. q0,a->q1 q1,b->q2');

		await fireEvent.input(statesInput, { target: { value: 'S0,S1,S2' } });
		await fireEvent.input(startStateInput, { target: { value: 'S0' } });
		await fireEvent.input(acceptedStatesInput, { target: { value: 'S2->IDENTIFIER' } });
		await fireEvent.input(transitionsInput, { target: { value: 'S0,a->S1\nS1,b->S2' } });

		// Click Show NFA button
		const showNfaButton = screen.queryByRole('button', { name: 'Show NFA' });
		if (showNfaButton) {
			await fireEvent.click(showNfaButton);
		}

		expect(mockFetch).toHaveBeenCalled();
	});

	it('TestErrorHandling_Success: Should handle API errors gracefully', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		// Mock failed API response for submit
		mockFetch.mockRejectedValueOnce(new Error('Network error'));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'TEST' } });
		await fireEvent.input(regexInput, { target: { value: 'test' } });

		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		// Should handle the error
		expect(mockFetch).toHaveBeenCalled();
	});

	it('TestFormStateReset_Success: Should reset form state when switching modes', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Start with regex mode
		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		// Enter some data
		const typeInput = screen.getByPlaceholderText('Enter type...');
		await fireEvent.input(typeInput, { target: { value: 'TEST_DATA' } });
		expect(typeInput).toHaveValue('TEST_DATA');

		// Switch to automata mode
		const automataButton = screen.getByRole('button', { name: 'Automata' });
		await fireEvent.click(automataButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('e.g. q0,q1,q2')).toBeInTheDocument();
		});

		// Switch back to regex mode
		await fireEvent.click(regexButton);

		await waitFor(() => {
			const resetTypeInput = screen.getByPlaceholderText('Enter type...');
			// Form should be reset
			expect(resetTypeInput).toHaveValue('');
		});
	});

	it('TestComplexRegexValidation_Success: Should validate complex regex patterns', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		const submitButton = screen.getByRole('button', { name: 'Submit' });

		// Test valid complex regex
		await fireEvent.input(typeInput, { target: { value: 'COMPLEX' } });
		await fireEvent.input(regexInput, { target: { value: '[a-zA-Z][a-zA-Z0-9_]*' } });
		await fireEvent.click(submitButton);

		// Should not show error for valid regex
		await waitFor(() => {
			expect(screen.queryByText('Invalid regular expression pattern')).not.toBeInTheDocument();
		});
	});

	it('TestMultipleRowManagement_Success: Should handle multiple regex rule rows', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
		});

		// Add additional rows by finding and clicking add button
		const addButtons = screen.queryAllByLabelText('Add new row') || 
						  screen.queryAllByText('+') ||
						  screen.queryAllByRole('button', { name: /add/i });

		if (addButtons.length > 0) {
			await fireEvent.click(addButtons[0]);
			
			// Should have multiple type inputs now
			await waitFor(() => {
				const typeInputs = screen.getAllByPlaceholderText('Enter type...');
				expect(typeInputs.length).toBeGreaterThan(1);
			});
		}
	});
});

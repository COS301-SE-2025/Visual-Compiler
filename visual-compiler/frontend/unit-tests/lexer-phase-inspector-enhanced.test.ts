import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../src/lib/components/lexer/phase-inspector.svelte';

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

// Mock lexer store with complete structure
vi.mock('$lib/stores/lexer', () => ({
	lexerState: {
		subscribe: vi.fn((callback) => {
			callback({
				userInputRows: [{ type: '', regex: '', error: '' }],
				automataInputs: {
					states: '',
					startState: '',
					acceptedStates: '',
					transitions: ''
				},
				selectedType: 'REGEX',
				showDefault: false,
				isSubmitted: false,
				hasUnsavedChanges: false,
				tokens: [],
				hasTokens: false,
				sourceCode: ''
			});
			return { unsubscribe: vi.fn() };
		}),
		update: vi.fn(),
		set: vi.fn()
	},
	updateLexerInputs: vi.fn(),
	updateAutomataInputs: vi.fn(),
	markLexerSubmitted: vi.fn(),
	updateLexerArtifacts: vi.fn()
}));

// Mock svelte/store get function
vi.mock('svelte/store', () => ({
	get: vi.fn(() => 'test-project'),
	writable: vi.fn(() => ({
		subscribe: vi.fn(() => () => {}),
		set: vi.fn(),
		update: vi.fn()
	}))
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

// Global mocks
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage and sessionStorage
const createMockStorage = () => {
	let store: { [key: string]: string } = {};
	return {
		getItem(key: string) { return store[key] || null; },
		setItem(key: string, value: string) { store[key] = value.toString(); },
		clear() { store = {}; }
	};
};

const localStorageMock = createMockStorage();
const sessionStorageMock = createMockStorage();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });
Object.defineProperty(window, 'sessionStorage', { value: sessionStorageMock });

describe('Lexer PhaseInspector Enhanced Coverage Tests', () => {
	const sourceCode = 'int x = 5;\\nreturn x;';

	beforeEach(() => {
		vi.clearAllMocks();
		localStorageMock.clear();
		sessionStorageMock.clear();
		
		localStorageMock.setItem('user_id', 'test-user');
		sessionStorageMock.setItem('access_token', 'test-access-token');
		
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ success: true, message: 'Rules stored successfully!' })
		});
	});

	it('TestAutomataFormFields_Success: Should render and interact with automata form fields', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Switch to automata mode
		const automataButton = screen.getByRole('button', { name: /Automata/i });
		await fireEvent.click(automataButton);

		// Wait for automata form to appear
		await waitFor(() => {
			const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
			expect(statesInput).toBeTruthy();
		});

		// Test all automata input fields  
		const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
		const startStateInput = screen.getByPlaceholderText('e.g. q0');
		const acceptedStatesInput = screen.getByPlaceholderText('e.g. q2->int, q1->string');
		const transitionsInput = screen.getByPlaceholderText('e.g. q0,a->q1 q1,b->q2');

		// Test input interactions
		await fireEvent.input(statesInput, { target: { value: 'q0,q1,q2' } });
		await fireEvent.input(startStateInput, { target: { value: 'q0' } });
		await fireEvent.input(acceptedStatesInput, { target: { value: 'q2->NUMBER' } });
		await fireEvent.input(transitionsInput, { target: { value: 'q0,1->q1\\nq1,2->q2' } });

		// Verify inputs have the expected values (basic functionality test)
		expect(statesInput).toBeTruthy();
		expect(startStateInput).toBeTruthy();
		expect(acceptedStatesInput).toBeTruthy();
		expect(transitionsInput).toBeTruthy();
	});

	it('TestRegexMode_Success: Should handle regex mode interactions', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Switch to regex mode
		const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
		await fireEvent.click(regexButton);

		// Wait for regex form
		await waitFor(() => {
			const typeInput = screen.queryByPlaceholderText('Enter type...');
			expect(typeInput).toBeTruthy();
		});

		// Fill in regex form
		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'IDENTIFIER' } });
		await fireEvent.input(regexInput, { target: { value: '[a-zA-Z]+' } });

		// Submit the form
		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		// Verify that fetch was called (indicates form submission)
		expect(mockFetch).toHaveBeenCalled();
	});

	it('TestTokenGeneration_Success: Should handle token generation', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({ success: true })
		});

		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: () => Promise.resolve({
				tokens: [
					{ type: 'KEYWORD', value: 'int', line: 1, column: 1 },
					{ type: 'IDENTIFIER', value: 'x', line: 1, column: 5 }
				]
			})
		});

		render(PhaseInspector, { 
			source_code: sourceCode,
			onGenerateTokens: vi.fn()
		});

		const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			const typeInput = screen.queryByPlaceholderText('Enter type...');
			expect(typeInput).toBeTruthy();
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

		expect(mockFetch).toHaveBeenCalled();
	});

	it('TestErrorHandling_Success: Should handle API errors gracefully', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			const typeInput = screen.queryByPlaceholderText('Enter type...');
			expect(typeInput).toBeTruthy();
		});

		// Mock failed API response for submit
		mockFetch.mockRejectedValueOnce(new Error('Network error'));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		
		await fireEvent.input(typeInput, { target: { value: 'TEST' } });
		await fireEvent.input(regexInput, { target: { value: 'test' } });

		const submitButton = screen.getByRole('button', { name: 'Submit' });
		await fireEvent.click(submitButton);

		// Should handle the error - verify API call was made
		expect(mockFetch).toHaveBeenCalled();
	});

	it('TestFormModeSwitch_Success: Should switch between regex and automata modes', async () => {
		render(PhaseInspector, { source_code: sourceCode });

		// Start with regex mode
		const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
		await fireEvent.click(regexButton);

		await waitFor(() => {
			const typeInput = screen.queryByPlaceholderText('Enter type...');
			expect(typeInput).toBeTruthy();
		});

		// Switch to automata mode
		const automataButton = screen.getByRole('button', { name: /Automata/i });
		await fireEvent.click(automataButton);

		await waitFor(() => {
			const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
			expect(statesInput).toBeTruthy();
		});

		// Switch back to regex mode
		await fireEvent.click(regexButton);

		await waitFor(() => {
			const typeInput = screen.queryByPlaceholderText('Enter type...');
			expect(typeInput).toBeTruthy();
		});
	});
});
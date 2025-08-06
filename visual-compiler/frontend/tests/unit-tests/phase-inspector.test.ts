// tests/unit-tests/phase-inspector.test.ts

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../../src/lib/components/lexer/phase-inspector.svelte';
import type { Token } from '$lib/types';

// Mock the toast store and fetch API
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

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

describe('PhaseInspector Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Set a user_id for each test to simulate a logged-in user
		window.localStorage.setItem('user_id', 'test-user-123');
	});

	const sourceCode = 'let x = 1;';

	it('TestInitialRender_Success: Renders correctly and shows form on button click', async () => {
		render(PhaseInspector, { source_code: sourceCode });
		expect(screen.getByText(sourceCode)).toBeInTheDocument();

		const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
		await fireEvent.click(regexButton);

		expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
	});

	it('TestValidation_Failure_InvalidRegex: Shows an error for an invalid regex pattern', async () => {
		render(PhaseInspector, { source_code: sourceCode });
		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		const submitButton = screen.getByRole('button', { name: 'Submit' });

		await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
		await fireEvent.input(regexInput, { target: { value: '[' } }); // Invalid regex
		await fireEvent.click(submitButton);

		// Check for the specific error message next to the input
		expect(await screen.findByText('Invalid regular expression pattern')).toBeInTheDocument();
		expect(mockFetch).not.toHaveBeenCalled();
	});

	it('TestValidation_Failure_IncompleteRow: Shows an error for an incomplete row', async () => {
		render(PhaseInspector, { source_code: sourceCode });
		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const submitButton = screen.getByRole('button', { name: 'Submit' });

		await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
		await fireEvent.click(submitButton);

		expect(
			await screen.findByText('Please fill in both Type and Regular Expression')
		).toBeInTheDocument();
	});

	it('TestAddRow_Success: Adds a new row when the plus button is clicked', async () => {
		render(PhaseInspector, { source_code: sourceCode });
		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');

		await fireEvent.input(typeInput, { target: { value: 'TYPE1' } });
		await fireEvent.input(regexInput, { target: { value: 'REGEX1' } });

		// The add button should now be visible based on the component's logic
		const addButton = await screen.findByRole('button', { name: '+' });
		await fireEvent.click(addButton);

		expect(screen.getAllByPlaceholderText('Enter type...').length).toBe(2);
	});

	it('TestHandleSubmit_Success: Calls fetch and shows generate button on valid submission', async () => {
		mockFetch.mockResolvedValue({ ok: true, json: () => Promise.resolve({ message: 'Success' }) });

		render(PhaseInspector, { source_code: sourceCode });
		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		const typeInput = screen.getByPlaceholderText('Enter type...');
		const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
		await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
		await fireEvent.input(regexInput, { target: { value: 'if|else' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

		await waitFor(() => {
			// The endpoint is /rules when selectedType is REGEX
			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/api/lexing/rules',
				expect.any(Object)
			);
		});

		expect(await screen.findByRole('button', { name: 'Generate Tokens' })).toBeInTheDocument();
	});

	it('TestHandleSubmit_Failure_ServerError: Shows error toast on fetch failure', async () => {
		mockFetch.mockResolvedValue({
			ok: false,
			status: 500,
			text: () => Promise.resolve('Internal Server Error')
		});

		render(PhaseInspector, { source_code: sourceCode });
		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		await fireEvent.input(screen.getByPlaceholderText('Enter type...'), {
			target: { value: 'KEYWORD' }
		});
		await fireEvent.input(screen.getByPlaceholderText('Enter regex pattern...'), {
			target: { value: 'if' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

		await waitFor(() => {
			// The component shows this specific error on fetch failure
			expect(AddToast).toHaveBeenCalledWith('Failed to save rules', 'error');
		});
	});

	it('TestGenerateTokens_Success: Calls generate endpoint and triggers correct callback', async () => {
		const mockServerResponse = {
			tokens: [{ Type: 'KEYWORD', Value: 'if' }],
			tokens_unidentified: ['@']
		};

		const expectedCallbackPayload = {
			tokens: mockServerResponse.tokens,
			unexpected_tokens: mockServerResponse.tokens_unidentified
		};

		// First fetch for handleSubmit, second for generateTokens
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ message: 'Rules stored successfully!' })
			})
			.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(mockServerResponse) });

		const handleGenerate = vi.fn();

		render(PhaseInspector, {
			source_code: sourceCode,
			onGenerateTokens: handleGenerate
		});

		await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

		await fireEvent.input(screen.getByPlaceholderText('Enter type...'), {
			target: { value: 'KEYWORD' }
		});
		await fireEvent.input(screen.getByPlaceholderText('Enter regex pattern...'), {
			target: { value: 'if' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

		const generateButton = await screen.findByRole('button', { name: 'Generate Tokens' });
		await fireEvent.click(generateButton);

		await waitFor(() => {
			// Check that the second fetch call was for token generation
			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/api/lexing/lexer',
				expect.any(Object)
			);
			expect(handleGenerate).toHaveBeenCalledWith(expectedCallbackPayload);
		});
	});
});

import { render, fireEvent, waitFor, screen } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import CodeInput from '../../src/lib/components/main/code-input.svelte';

// Mock the toast store
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn(),
	toasts: []
}));
import { AddToast } from '$lib/stores/toast';

// Mock the project store
vi.mock('$lib/stores/project', () => ({
	projectName: {
		subscribe: vi.fn((callback) => {
			callback('test-project');
			return vi.fn(); // Return the unsubscribe function directly
		})
	}
}));

// Mock the source code store
vi.mock('$lib/stores/source-code', () => ({
	confirmedSourceCode: {
		set: vi.fn(),
		update: vi.fn(),
		subscribe: vi.fn((callback) => {
			callback('');
			return vi.fn(); // Return the unsubscribe function directly
		})
	}
}));

// Mock the global fetch function
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

describe('CodeInput Component', () => {
	beforeEach(() => {
		// Clear mocks and set a user_id before each test
		vi.clearAllMocks();
		window.localStorage.setItem('user_id', 'test-user-123');
	});

	it('TestEmptyState_Success: Renders with the confirm button disabled', () => {
		render(CodeInput);
		const confirm_button = screen.getByText(/confirm code/i);
		expect(confirm_button).toBeDisabled();
	});

	it('TestButtonEnable_Success: Enables the confirm button when text is entered', async () => {
		render(CodeInput);
		const textarea = screen.getByPlaceholderText(/paste or type your source code here…/i);
		const confirm_button = screen.getByText(/confirm code/i);

		await fireEvent.input(textarea, { target: { value: 'some code' } });
		expect(confirm_button).toBeEnabled();
	});

	it('TestFileUpload_Failure: Shows an error when a non-.txt file is uploaded', async () => {
		render(CodeInput);
		const file_input = screen.getByLabelText(/upload file/i) as HTMLInputElement;

		const file = new File(['content'], 'test.pdf', { type: 'application/pdf' });
		await fireEvent.change(file_input, { target: { files: [file] } });
		expect(AddToast).toHaveBeenCalledWith(
			'Only .txt files are allowed. Please upload a valid plain text file.',
			'error'
		);
	});

	it('TestFileUpload_Success: Shows a success message for a .txt file upload', async () => {
		render(CodeInput);
		const file_input = screen.getByLabelText(/upload file/i) as HTMLInputElement;

		const file = new File(['file content'], 'test.txt', { type: 'text/plain' });
		await fireEvent.change(file_input, { target: { files: [file] } });

		await waitFor(() => {
			expect(AddToast).toHaveBeenCalledWith('File uploaded successfully!', 'success');
		});
	});

	it('TestEventDispatch_Success: Dispatches the codeSubmitted event with the correct text', async () => {
		const mockHandler = vi.fn();
		const test_code = 'let x = 10';

		// Mock a successful fetch response
		const mockResponse = {
			ok: true
		};
		mockFetch.mockResolvedValue(mockResponse);

		render(CodeInput, { onCodeSubmitted: mockHandler });

		const textarea = screen.getByPlaceholderText(/paste or type your source code here…/i);
		const confirm_button = screen.getByText(/confirm code/i);

		await fireEvent.input(textarea, { target: { value: test_code } });
		await fireEvent.click(confirm_button);

		// Wait for the fetch promise to resolve and the event to be called
		await waitFor(() => {
			expect(mockHandler).toHaveBeenCalledWith(test_code);
		});
	});
});

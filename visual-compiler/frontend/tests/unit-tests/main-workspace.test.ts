import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, fireEvent, screen, within } from '@testing-library/svelte';
import MainWorkspace from '../../src/routes/main-workspace/+page.svelte';

// Mock SvelteKit runtime
(globalThis as any).__SVELTEKIT_PAYLOAD__ = {
	data: {},
	errors: {}
};

// Mock SvelteKit navigation
vi.mock('$app/navigation', () => ({
	goto: vi.fn(),
	invalidateAll: vi.fn()
}));

vi.mock('$app/stores', () => ({
	page: {
		subscribe: vi.fn(() => ({ unsubscribe: vi.fn() }))
	}
}));

const localStorageMock = (() => {
	let store: { [key: string]: string } = {};
	return {
		getItem: (key: string) => store[key] || null,
		setItem: (key: string, value: string) => (store[key] = value.toString()),
		removeItem: (key: string) => delete store[key],
		clear: () => (store = {})
	};
})();

Object.defineProperty(window, 'localStorage', {
	value: localStorageMock
});

describe('MainWorkspace Component', () => {
	afterEach(() => {
		localStorage.clear();
		vi.clearAllMocks();
	});

	it('TestInitialRender_Success: Renders all main sections of the page', () => {
		render(MainWorkspace);
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();

		const toolbox = screen.getByTestId('toolbox');
		// FIX: Use getByText for a more direct query inside the toolbox
		expect(within(toolbox).getByText('Source Code')).toBeInTheDocument();
		expect(within(toolbox).getByText('Lexer')).toBeInTheDocument();
	});

	it('TestNodeCreation_Success: Creates a Source Code node when its button is clicked', async () => {
		render(MainWorkspace);

		const toolbox = screen.getByTestId('toolbox');
		// FIX: Use getByText to unambiguously find the button
		const source_code_button = within(toolbox).getByText('Source Code');

		await fireEvent.click(source_code_button);

		const all_nodes = screen.getAllByText('Source Code');
		expect(all_nodes.length).toBe(2);
	});

	it('TestMultipleNodeCreation_Success: Creates multiple unique nodes on subsequent clicks', async () => {
		render(MainWorkspace);

		const toolbox = screen.getByTestId('toolbox');
		// FIX: Use getByText for unambiguous queries
		const source_code_button = within(toolbox).getByText('Source Code');
		const lexer_button = within(toolbox).getByText('Lexer');

		await fireEvent.click(source_code_button);
		await fireEvent.click(lexer_button);

		const node_labels = screen.getAllByText(/Source Code|Lexer/);
		expect(node_labels.length).toBe(4);
	});

	it('TestProTip_Success: Shows the pro-tip on first visit and dismisses it on click', async () => {
		// Clear localStorage and ensure hasSeenDragTip is not set
		localStorage.clear();
		expect(localStorage.getItem('hasSeenDragTip')).toBeNull();
		
		render(MainWorkspace);

		// First check if the dismiss button is present (which would indicate the tip is showing)
		const dismiss_button = await screen.findByRole('button', { name: 'Dismiss tip' }, { timeout: 3000 });
		expect(dismiss_button).toBeInTheDocument();

		// Click the dismiss button
		await fireEvent.click(dismiss_button);

		// Verify the button is no longer present (tip dismissed)
		expect(screen.queryByRole('button', { name: 'Dismiss tip' })).toBeNull();
	});
});

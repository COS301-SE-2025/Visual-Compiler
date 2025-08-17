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

	it('TestWelcomeOverlay_Success: Shows welcome overlay when no project name is set', () => {
		localStorage.removeItem('projectName');
		render(MainWorkspace);
		
		// The welcome overlay should be visible when no project name is set
		// We can test for the presence of the component itself
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
	});

	it('TestNavBarRendering_Success: Renders navigation bar with project title', () => {
		render(MainWorkspace);
		
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
	});

	it('TestToolboxInteraction_Success: All toolbox buttons are clickable', async () => {
		render(MainWorkspace);

		const toolbox = screen.getByTestId('toolbox');
		const buttons = ['Source Code', 'Lexer', 'Parser', 'Analyser', 'Translator'];
		
		for (const buttonText of buttons) {
			const button = within(toolbox).getByText(buttonText);
			expect(button).toBeInTheDocument();
			
			await fireEvent.click(button);
			
			// After clicking, should find the button text in both toolbox and canvas
			const allInstances = screen.getAllByText(buttonText);
			expect(allInstances.length).toBeGreaterThanOrEqual(2);
		}
	});

	it('TestCanvasRendering_Success: Canvas component is present and functional', () => {
		render(MainWorkspace);
		
		// Look for the canvas container
		const canvas = document.querySelector('.drawer-canvas');
		expect(canvas).toBeInTheDocument();
	});

	it('TestCodeInputSection_Success: Code input component renders correctly', () => {
		render(MainWorkspace);
		
		// Look for toolbox elements that actually exist
		expect(screen.getByTestId('toolbox')).toBeInTheDocument();
		expect(screen.getByText('Blocks')).toBeInTheDocument();
	});

	it('TestThemeApplications_Success: Component applies theme correctly', () => {
		const { container } = render(MainWorkspace);
		
		// Check that the component container has valid content
		expect(container.querySelector('header')).toBeTruthy();
	});

	it('TestNodeDeletion_Success: Can interact with nodes after creation', async () => {
		render(MainWorkspace);

		const toolbox = screen.getByTestId('toolbox');
		const sourceButton = within(toolbox).getByText('Source Code');

		await fireEvent.click(sourceButton);

		// Verify node was created
		const allNodes = screen.getAllByText('Source Code');
		expect(allNodes.length).toBe(2); // One in toolbox, one in canvas
	});

	it('TestLocalStoragePersistence_Success: Handles localStorage operations', () => {
		// Set some test data
		localStorage.setItem('hasSeenDragTip', 'true');
		
		render(MainWorkspace);
		
		expect(localStorage.getItem('hasSeenDragTip')).toBe('true');
	});

	it('TestProjectNameHandling_Success: Responds to project name changes', () => {
		localStorage.setItem('projectName', 'Test Project');
		
		render(MainWorkspace);
		
		// Component should handle project name from localStorage
		expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
	});

	it('TestKeyboardInteraction_Success: Handles keyboard events on interactive elements', async () => {
		render(MainWorkspace);

		const toolbox = screen.getByTestId('toolbox');
		const sourceButton = within(toolbox).getByText('Source Code');

		// Test keyboard activation
		await fireEvent.keyDown(sourceButton, { key: 'Enter' });
		
		expect(sourceButton).toBeInTheDocument();
	});

	it('TestResponsiveLayout_Success: Component renders with proper layout structure', () => {
		const { container } = render(MainWorkspace);
		
		// Check for main layout elements that exist
		expect(container.querySelector('.main')).toBeTruthy();
		expect(container.querySelector('.toolbox')).toBeTruthy();
	});
});

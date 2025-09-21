import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Toolbox from '../src/lib/components/main/Toolbox.svelte';

// Mock the toast store at the top level
vi.mock('../src/lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

// Mock the theme store
vi.mock('../src/lib/stores/theme', () => ({
	theme: { subscribe: vi.fn() }
}));

// Import the mocked function
import { AddToast } from '../src/lib/stores/toast';

// Mock tooltips to simplify the test
const mockTooltips = {
	source: 'Source tooltip',
	lexer: 'Lexer tooltip',
	parser: 'Parser tooltip',
	analyser: 'Analyser tooltip',
	translator: 'Translator tooltip'
};

describe('Toolbox Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		// Clean up any event listeners
		document.removeEventListener('resetToolbox', expect.any(Function));
	});

	it('TestInitialRender_Success: Renders the heading and all phase buttons', () => {
		render(Toolbox, {
			props: {
				handleCreateNode: vi.fn(),
				tooltips: mockTooltips
			}
		});

		expect(screen.getByText('Blocks')).toBeInTheDocument();
		expect(screen.getByText('Click a block to add it to the canvas.')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Source Code/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Lexer/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Parser/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Analyser/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Translator/i })).toBeInTheDocument();
	});

	it('TestNodeCreation_Success: Calls the handler function with the correct type on click', async () => {
		const mockHandler = vi.fn();
		render(Toolbox, {
			props: {
				handleCreateNode: mockHandler,
				tooltips: mockTooltips
			}
		});

		// Test clicking the 'Lexer' button wrapper
		const lexerButton = screen.getByRole('button', { name: /Lexer/i });
		const lexerButtonWrapper = lexerButton.parentElement;
		if (lexerButtonWrapper) {
			await fireEvent.click(lexerButtonWrapper);
		}
		expect(mockHandler).toHaveBeenCalledWith('lexer');

		// Test clicking the 'Parser' button wrapper
		const parserButton = screen.getByRole('button', { name: /Parser/i });
		const parserButtonWrapper = parserButton.parentElement;
		if (parserButtonWrapper) {
			await fireEvent.click(parserButtonWrapper);
		}
		expect(mockHandler).toHaveBeenCalledWith('parser');
	});

	it('TestNodeCreationViaButton_Success: Calls handler when clicking button directly', async () => {
		const mockHandler = vi.fn();
		render(Toolbox, {
			props: {
				handleCreateNode: mockHandler,
				tooltips: mockTooltips
			}
		});

		// Test clicking directly on the button wrapper (which handles the click)
		const sourceButton = screen.getByRole('button', { name: /Source Code/i });
		const sourceButtonWrapper = sourceButton.parentElement;
		if (sourceButtonWrapper) {
			await fireEvent.click(sourceButtonWrapper);
		}
		expect(mockHandler).toHaveBeenCalledWith('source');
	});

	it('TestDisabledNodeHandling_Success: Shows toast when trying to create duplicate node', async () => {
		const mockHandler = vi.fn();
		const mockNodes = [
			{ id: '1', type: 'source' as const, label: 'Source', position: { x: 0, y: 0 } }
		];

		render(Toolbox, {
			props: {
				handleCreateNode: mockHandler,
				tooltips: mockTooltips,
				nodes: mockNodes
			}
		});

		// Wait for reactive statement to process nodes
		await waitFor(() => {
			const sourceButton = screen.getByRole('button', { name: /Source Code/i });
			expect(sourceButton).toBeDisabled();
		});

		// Try to click the disabled source button (via wrapper)
		const sourceButton = screen.getByRole('button', { name: /Source Code/i });
		const sourceButtonWrapper = sourceButton.parentElement;
		if (sourceButtonWrapper) {
			await fireEvent.click(sourceButtonWrapper);
		}

		// Should show toast and not call handler
		expect(AddToast).toHaveBeenCalledWith(
			'Duplicate node: Only one node of each type is allowed in the pipeline',
			'info'
		);
		expect(mockHandler).not.toHaveBeenCalled();
	});

	it('TestToastOnlyShownOnce_Success: Toast is only shown once for duplicate attempts', async () => {
		const mockHandler = vi.fn();
		const mockNodes = [
			{ id: '1', type: 'lexer' as const, label: 'Lexer', position: { x: 0, y: 0 } }
		];

		render(Toolbox, {
			props: {
				handleCreateNode: mockHandler,
				tooltips: mockTooltips,
				nodes: mockNodes
			}
		});

		// Wait for reactive statement to process nodes
		await waitFor(() => {
			const lexerButton = screen.getByRole('button', { name: /Lexer/i });
			expect(lexerButton).toBeDisabled();
		});

		// Click disabled button multiple times
		const lexerButton = screen.getByRole('button', { name: /Lexer/i });
		const lexerButtonWrapper = lexerButton.parentElement;
		
		if (lexerButtonWrapper) {
			await fireEvent.click(lexerButtonWrapper);
			await fireEvent.click(lexerButtonWrapper);
			await fireEvent.click(lexerButtonWrapper);
		}

		// Toast should only be called once
		expect(AddToast).toHaveBeenCalledTimes(1);
	});

	it('TestResetToolboxEvent_Success: Responds to reset event and clears created node types', async () => {
		const mockHandler = vi.fn();
		const mockNodes = [
			{ id: '1', type: 'parser' as const, label: 'Parser', position: { x: 0, y: 0 } }
		];

		const { component, rerender } = render(Toolbox, {
			handleCreateNode: mockHandler,
			tooltips: mockTooltips,
			nodes: mockNodes
		});

		// Wait for reactive statement to process nodes
		await waitFor(() => {
			const parserButton = screen.getByRole('button', { name: /Parser/i });
			expect(parserButton).toBeDisabled();
		});

		// Dispatch reset event
		const resetEvent = new CustomEvent('resetToolbox');
		document.dispatchEvent(resetEvent);

		// Re-render with empty nodes to simulate reset
		await rerender({
			handleCreateNode: mockHandler,
			tooltips: mockTooltips,
			nodes: []
		});

		// Wait for reactive statement to process empty nodes
		await waitFor(() => {
			const parserButton = screen.getByRole('button', { name: /Parser/i });
			expect(parserButton).not.toBeDisabled();
		});

		// Should be able to create node again
		const parserButton = screen.getByRole('button', { name: /Parser/i });
		const parserButtonWrapper = parserButton.parentElement;
		if (parserButtonWrapper) {
			await fireEvent.click(parserButtonWrapper);
		}
		expect(mockHandler).toHaveBeenCalledWith('parser');
	});

	it('TestReactiveNodesSyncing_Success: Buttons disable when nodes are added', async () => {
		const mockHandler = vi.fn();
		
		const { rerender } = render(Toolbox, {
			handleCreateNode: mockHandler,
			tooltips: mockTooltips,
			nodes: []
		});

		// Initially all buttons should be enabled
		expect(screen.getByRole('button', { name: /Source Code/i })).not.toBeDisabled();
		expect(screen.getByRole('button', { name: /Lexer/i })).not.toBeDisabled();

		// Add nodes
		const mockNodesWithMultiple = [
			{ id: '1', type: 'source' as const, label: 'Source', position: { x: 0, y: 0 } },
			{ id: '2', type: 'analyser' as const, label: 'Analyser', position: { x: 100, y: 0 } }
		];

		await rerender({
			handleCreateNode: mockHandler,
			tooltips: mockTooltips,
			nodes: mockNodesWithMultiple
		});

		// Wait for reactive statement to process
		await waitFor(() => {
			expect(screen.getByRole('button', { name: /Source Code/i })).toBeDisabled();
			expect(screen.getByRole('button', { name: /Analyser/i })).toBeDisabled();
			expect(screen.getByRole('button', { name: /Lexer/i })).not.toBeDisabled();
		});
	});

	it('TestTooltipDisplay_Success: Tooltips are rendered for each button', () => {
		render(Toolbox, {
			props: {
				handleCreateNode: vi.fn(),
				tooltips: mockTooltips
			}
		});

		// Check tooltips are rendered (they have the custom-tooltip class)
		const tooltips = document.querySelectorAll('.custom-tooltip');
		expect(tooltips).toHaveLength(6); // One for each node type
		
		// Check specific tooltip content
		expect(screen.getByText('Source tooltip')).toBeInTheDocument();
		expect(screen.getByText('Lexer tooltip')).toBeInTheDocument();
		expect(screen.getByText('Parser tooltip')).toBeInTheDocument();
	});

	it('TestButtonNumbers_Success: Buttons show correct numbers', () => {
		render(Toolbox, {
			props: {
				handleCreateNode: vi.fn(),
				tooltips: mockTooltips
			}
		});

		// Check button numbers are displayed
		const numberElements = document.querySelectorAll('.button-number-corner');
		expect(numberElements).toHaveLength(5);
		
		expect(numberElements[0]?.textContent).toBe('1'); // Source
		expect(numberElements[1]?.textContent).toBe('2'); // Lexer
		expect(numberElements[2]?.textContent).toBe('3'); // Parser
		expect(numberElements[3]?.textContent).toBe('4'); // Analyser
		expect(numberElements[4]?.textContent).toBe('5'); // Translator
	});

	it('TestComponentCleanup_Success: Event listeners are properly removed', () => {
		const removeEventListenerSpy = vi.spyOn(document, 'removeEventListener');
		
		const { unmount } = render(Toolbox, {
			props: {
				handleCreateNode: vi.fn(),
				tooltips: mockTooltips
			}
		});

		// Unmount the component
		unmount();

		// Check that removeEventListener was called
		expect(removeEventListenerSpy).toHaveBeenCalledWith('resetToolbox', expect.any(Function));
		
		removeEventListenerSpy.mockRestore();
	});
});



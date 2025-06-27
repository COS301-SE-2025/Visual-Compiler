import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import toolbox from '../../src/lib/components/main/toolbox.svelte';

// Mock tooltips to simplify the test
const mockTooltips = {
	source: 'Source tooltip',
	lexer: 'Lexer tooltip',
	parser: 'Parser tooltip'
};

describe('Toolbox Component', () => {
	it('TestInitialRender_Success: Renders the heading and all phase buttons', () => {
		// FIX D: Add the required 'handleCreateNode' prop
		render(toolbox, {
			props: {
				handleCreateNode: vi.fn(), // Provide a mock function
				tooltips: mockTooltips
			}
		});

		expect(screen.getByText('Blocks')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Source Code/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Lexer/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /Parser/i })).toBeInTheDocument();
	});

	it('TestNodeCreation_Success: Calls the handler function with the correct type on click', async () => {
		const mockHandler = vi.fn();
		render(toolbox, {
			props: {
				handleCreateNode: mockHandler,
				tooltips: mockTooltips
			}
		});

		// Test clicking the 'Lexer' button
		const lexerButton = screen.getByRole('button', { name: /Lexer/i });
		await fireEvent.click(lexerButton);
		expect(mockHandler).toHaveBeenCalledWith('lexer');

		// Test clicking the 'Parser' button
		const parserButton = screen.getByRole('button', { name: /Parser/i });
		await fireEvent.click(parserButton);
		expect(mockHandler).toHaveBeenCalledWith('parser');
	});
});

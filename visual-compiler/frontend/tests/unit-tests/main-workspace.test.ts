import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, fireEvent, screen, within } from '@testing-library/svelte';
import MainWorkspace from '../../src/routes/main-workspace/+page.svelte';

const localStorageMock = (() => {
  let store: { [key: string]: string } = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => (store[key] = value.toString()),
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
    render(MainWorkspace);

    // Use a custom function to find text across multiple elements
    const pro_tip = await screen.findByText((content, element) => {
        if (!element) {
            return false;
        }
        return element.textContent === 'Pro-Tip: For the smoothest experience, click to select a node before dragging it.';
    });
    
    expect(pro_tip).toBeInTheDocument();

    const dismiss_button = screen.getByRole('button', { name: 'Dismiss tip' });
    await fireEvent.click(dismiss_button);

    expect(screen.queryByText(/Pro-Tip:/)).toBeNull();
});
});
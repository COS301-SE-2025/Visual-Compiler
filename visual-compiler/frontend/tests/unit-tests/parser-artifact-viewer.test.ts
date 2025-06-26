import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import ParserArtifactViewer from '../../src/lib/components/parser/parser-artifact-viewer.svelte';
import type { SyntaxTree } from '$lib/types';

// Mock the vis-network library to prevent errors in the JSDOM test environment.
const mockNetwork = {
  destroy: vi.fn()
};

vi.mock('vis-network/standalone', () => ({
  Network: vi.fn(() => mockNetwork),
  // DataSet needs to be a class that returns an object with an 'add' method.
  DataSet: vi.fn(() => ({
    add: vi.fn()
  }))
}));

describe('ParserArtifactViewer Component', () => {
  // TestInitialState_Success
  it('TestInitialState_Success: Renders the empty state message when no tree is provided', () => {
    render(ParserArtifactViewer, { syntaxTree: null });
    expect(
      screen.getByText('The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.')
    ).toBeInTheDocument();
  });

  // TestTreeRender_Success
  it('TestTreeRender_Success: Renders the tree container when a syntax tree is provided', () => {
    const mock_tree: SyntaxTree = {
      root: {
        symbol: 'STATEMENT',
        value: '',
        children: [{ symbol: 'DECLARATION', value: 'let', children: null }]
      }
    };

    render(ParserArtifactViewer, { syntaxTree: mock_tree });

    expect(
      screen.queryByText('The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.')
    ).toBeNull();
    expect(screen.getByText('Artifacts: Syntax Tree')).toBeInTheDocument();
    
    expect(screen.getByTestId('tree-display-container')).toBeInTheDocument();
  });

  // TestTreeUpdate_Success
  it('TestTreeUpdate_Success: Updates the view when the syntaxTree prop changes', async () => {
    const { rerender } = render(ParserArtifactViewer, { syntaxTree: null });

    expect(
      screen.getByText('The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.')
    ).toBeInTheDocument();

    const mock_tree: SyntaxTree = {
      root: { symbol: 'PROGRAM', value: '', children: [] }
    };
    await rerender({ syntaxTree: mock_tree });

    expect(
      screen.queryByText('The generated Abstract Syntax Tree (AST) will be visualized here once it is generated.')
    ).toBeNull();
    expect(screen.getByTestId('tree-display-container')).toBeInTheDocument();
  });

  // TestTreeDestroy_Success
  it('TestTreeDestroy_Success: Cleans up the network instance when the tree is removed', async () => {
    const mock_tree: SyntaxTree = {
      root: { symbol: 'PROGRAM', value: '', children: [] }
    };
    const { rerender } = render(ParserArtifactViewer, { syntaxTree: mock_tree });
    
    vi.clearAllMocks();
    await rerender({ syntaxTree: null });

    // Check if the destroy method was called on the mockNetwork instance.
    expect(mockNetwork.destroy).toHaveBeenCalled();
  });
});

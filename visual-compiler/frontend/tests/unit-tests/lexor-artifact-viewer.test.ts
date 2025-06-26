import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import ArtifactViewer from '../../src/lib/components/lexer/lexer-artifact-viewer.svelte';
import ToastContainer from '../../src/lib/components/toast-conatiner.svelte';
import type { Token } from '$lib/types';

describe('ArtifactViewer Component', () => {
  it('TestEmptyState_Success: Renders the initial empty message for the lexer phase', () => {

    render(ArtifactViewer, { props: { phase: 'lexer', show_tokens: false } });
    expect(
      screen.getByText('Tokens will appear here after generation')
    ).toBeInTheDocument();
  });

  it('TestWrongPhase_Success: Does not render the lexer empty message for other phases', () => {

    render(ArtifactViewer, { props: { phase: 'parser' } });
    expect(
      screen.queryByText('Tokens will appear here after generation')
    ).toBeNull();
  });

  it('TestTokenDisplay_Success: Renders a table of tokens when provided', () => {
    const mock_tokens: Token[] = [
      { Type: 'KEYWORD', Value: 'if' },
      { Type: 'PUNCTUATION', Value: '(' }
    ];


    render(ArtifactViewer, {
      props: {
        phase: 'lexer',
        show_tokens: true,
        tokens: mock_tokens,
        unexpected_tokens: []
      }
    });

    expect(screen.getByText('Type')).toBeInTheDocument();
    expect(screen.getByText('if')).toBeInTheDocument();
  });

  it('TestUnidentifiedTokenDisplay_Success: Renders a list of unidentified tokens when provided', () => {
    const mock_unexpected_tokens = ['@', '#'];


    render(ArtifactViewer, {
      props: {
        phase: 'lexer',
        show_tokens: true,
        tokens: [],
        unexpected_tokens: mock_unexpected_tokens
      }
    });

    expect(screen.getByText('Unidentified Input')).toBeInTheDocument();
    expect(screen.getByText('@')).toBeInTheDocument();
  });

  it('TestNoTokensMessage_Success: Renders a message when token generation results in an empty array', () => {

    render(ArtifactViewer, {
      props: {
        phase: 'lexer',
        show_tokens: true,
        tokens: [],
        unexpected_tokens: []
      }
    });

    expect(screen.getByText('No tokens generated yet')).toBeInTheDocument();
  });
});
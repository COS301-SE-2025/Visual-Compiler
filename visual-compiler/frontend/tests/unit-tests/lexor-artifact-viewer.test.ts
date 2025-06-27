import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import ArtifactViewer from '../../src/lib/components/lexer/lexer-artifact-viewer.svelte';
import type { Token } from '$lib/types';

describe('ArtifactViewer Component', () => {
	// TestEmptyState_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies the component shows the correct initial message when no tokens are available.
	it('TestEmptyState_Success: Renders the initial empty message for the lexer phase', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: false
		});
		expect(screen.getByText('Tokens will appear here after generation')).toBeInTheDocument();
	});

	// TestWrongPhase_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies that the component does not render the lexer-specific message for other phases.
	it('TestWrongPhase_Success: Does not render the lexer empty message for other phases', () => {
		render(ArtifactViewer, { phase: 'parser' });
		expect(screen.queryByText('Tokens will appear here after generation')).toBeNull();
	});

	// TestTokenDisplay_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies that a table of tokens is rendered correctly from the props.
	it('TestTokenDisplay_Success: Renders a table of tokens when provided', () => {
		// FIX: Changed mock data properties to lowercase (`type`, `value`) to match the component's expectation.
		const mock_tokens: Token[] = [
			{ type: 'KEYWORD', value: 'if' },
			{ type: 'PUNCTUATION', value: '(' }
		];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: mock_tokens,
			unexpected_tokens: []
		});

		// FIX: Changed the assertion to look for the correctly capitalized table header "Type".
		expect(screen.getByText('Type')).toBeInTheDocument();
		expect(screen.getByText('if')).toBeInTheDocument();
		expect(screen.getByText('(')).toBeInTheDocument();
	});

	// TestUnidentifiedTokenDisplay_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies that a list of unidentified tokens is rendered correctly.
	it('TestUnidentifiedTokenDisplay_Success: Renders a list of unidentified tokens when provided', () => {
		const mock_unexpected_tokens = ['@', '#'];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [],
			unexpected_tokens: mock_unexpected_tokens
		});

		expect(screen.getByText('Unidentified Input')).toBeInTheDocument();
		expect(screen.getByText('@')).toBeInTheDocument();
		expect(screen.getByText('#')).toBeInTheDocument();
	});

	// TestNoTokensMessage_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies that the correct message is shown when token generation results in an empty array.
	it('TestNoTokensMessage_Success: Renders a message when token generation results in an empty array', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [],
			unexpected_tokens: []
		});

		expect(screen.getByText('No tokens generated yet')).toBeInTheDocument();
	});
});

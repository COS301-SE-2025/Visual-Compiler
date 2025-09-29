import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import ArtifactViewer from '../src/lib/components/lexer/lexer-artifact-viewer.svelte';
import type { Token } from '$lib/types';

// Mock the lexer store
const mockUnsubscribe = vi.fn();
vi.mock('$lib/stores/lexer', () => ({
	lexerState: {
		subscribe: vi.fn((callback) => {
			// Call with empty state initially
			callback({ tokens: null, tokens_unidentified: null });
			return mockUnsubscribe;
		})
	}
}));

// Mock svelte/store get function
vi.mock('svelte/store', () => ({
	get: vi.fn(() => ({ tokens: null, tokens_unidentified: null }))
}));

describe('ArtifactViewer Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});
	// TestEmptyState_Success
	// Return type: void
	// Parameter type(s): none
	// Verifies the component shows the correct initial message when no tokens are available.
	it('TestEmptyState_Success: Renders the initial empty message for the lexer phase', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: false
		});
		// The component shows this message when phase is lexer but show_tokens is false
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
		expect(screen.getByText('Value')).toBeInTheDocument();
		expect(screen.getByText('KEYWORD')).toBeInTheDocument();
		expect(screen.getByText('if')).toBeInTheDocument();
		expect(screen.getByText('PUNCTUATION')).toBeInTheDocument();
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

		expect(screen.getByText('Unidentified Token(s) Found')).toBeInTheDocument();
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

	// TestStoreSubscription_Success
	// Tests that the component correctly subscribes to lexer state and renders tokens from props
	it('TestStoreSubscription_Success: Component renders tokens from props when show_tokens is true', () => {
		const mockTokens: Token[] = [
			{ type: 'IDENTIFIER', value: 'variable' },
			{ type: 'OPERATOR', value: '=' }
		];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: mockTokens,
			unexpected_tokens: []
		});

		// Verify tokens are displayed
		expect(screen.getByText('IDENTIFIER')).toBeInTheDocument();
		expect(screen.getByText('variable')).toBeInTheDocument();
		expect(screen.getByText('OPERATOR')).toBeInTheDocument();
		expect(screen.getByText('=')).toBeInTheDocument();
	});

	// TestStoreSubscriptionWithUnidentifiedTokens_Success
	// Tests component with unidentified tokens through props
	it('TestStoreSubscriptionWithUnidentifiedTokens_Success: Renders unidentified tokens from props', () => {
		const mockUnexpectedTokens = ['$', '%', '&'];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [],
			unexpected_tokens: mockUnexpectedTokens
		});

		expect(screen.getByText('Unidentified Token(s) Found')).toBeInTheDocument();
		expect(screen.getByText('$')).toBeInTheDocument();
		expect(screen.getByText('%')).toBeInTheDocument();
		expect(screen.getByText('&')).toBeInTheDocument();
	});

	// TestOnMountStateCheck_Success
	// Tests onMount lifecycle checking initial lexer state
	it('TestOnMountStateCheck_Success: Component mounts and checks lexer store state', () => {
		// Test that component can mount without errors and use store
		const { container } = render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: false
		});

		// Verify component mounted successfully
		expect(container.firstChild).toBeInTheDocument();
		expect(screen.getByText('Tokens will appear here after generation')).toBeInTheDocument();
	});

	// TestComponentCleanup_Success  
	// Tests that onDestroy properly unsubscribes from store
	it('TestComponentCleanup_Success: Unsubscribes from store on destroy', () => {
		const { unmount } = render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: false
		});

		// Unmount component to trigger cleanup
		unmount();

		// Verify unsubscribe was called during cleanup
		expect(mockUnsubscribe).toHaveBeenCalled();
	});

	// TestMixedTokensAndErrors_Success
	// Tests rendering when both tokens and unidentified tokens exist
	it('TestMixedTokensAndErrors_Success: Renders both tokens and unidentified tokens', () => {
		const mockTokens: Token[] = [
			{ type: 'NUMBER', value: '42' }
		];
		const mockUnexpected = ['@unknown'];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: mockTokens,
			unexpected_tokens: mockUnexpected
		});

		// Verify both sections are rendered
		expect(screen.getByText('Unidentified Token(s) Found')).toBeInTheDocument();
		expect(screen.getByText('@unknown')).toBeInTheDocument();
		expect(screen.getByText('NUMBER')).toBeInTheDocument();
		expect(screen.getByText('42')).toBeInTheDocument();
	});

	// TestTableStructure_Success
	// Tests the table structure and accessibility
	it('TestTableStructure_Success: Renders accessible table structure', () => {
		const mockTokens: Token[] = [
			{ type: 'STRING', value: '"hello"' },
			{ type: 'SEMICOLON', value: ';' }
		];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: mockTokens,
			unexpected_tokens: []
		});

		// Check table structure
		const table = screen.getByRole('table');
		expect(table).toBeInTheDocument();
		expect(table).toHaveClass('token-table');

		// Check headers
		const headers = screen.getAllByRole('columnheader');
		expect(headers).toHaveLength(2);
		expect(headers[0]).toHaveTextContent('Type');
		expect(headers[1]).toHaveTextContent('Value');

		// Check rows
		const rows = screen.getAllByRole('row');
		expect(rows).toHaveLength(3); // header + 2 data rows
	});

	// TestComponentStructure_Success
	// Tests the overall component structure and headings
	it('TestComponentStructure_Success: Renders correct component structure', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [],
			unexpected_tokens: []
		});

		// Check main container
		const container = document.querySelector('.artifact-viewer');
		expect(container).toBeInTheDocument();

		// Check headings
		expect(screen.getByText('Lexer Artefact')).toBeInTheDocument();
		expect(screen.getByText('Tokens')).toBeInTheDocument();

		// Check heading elements structure
		const artifactTitle = screen.getByText('Lexer Artefact');
		expect(artifactTitle).toHaveClass('artifact-title');
		
		const tokensHeading = screen.getByRole('heading', { name: 'Tokens' });
		expect(tokensHeading.tagName).toBe('H3');
	});

	// TestErrorStateStructure_Success
	// Tests the error state component structure
	it('TestErrorStateStructure_Success: Renders complete error state structure', () => {
		const mockUnexpected = ['invalid', 'tokens'];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [],
			unexpected_tokens: mockUnexpected
		});

		// Check error state container
		const errorState = document.querySelector('.error-state');
		expect(errorState).toBeInTheDocument();

		// Check error icon
		const errorIcon = document.querySelector('.error-icon');
		expect(errorIcon).toBeInTheDocument();

		// Check error components
		expect(screen.getByText('Unidentified Token(s) Found')).toBeInTheDocument();
		expect(screen.getByText(/The source code could not be completely tokenised/)).toBeInTheDocument();

		// Check token list container
		const tokenList = document.querySelector('.token-list');
		expect(tokenList).toBeInTheDocument();

		// Check individual error details
		const errorDetails = document.querySelectorAll('.error-details');
		expect(errorDetails).toHaveLength(2);
		expect(errorDetails[0]).toHaveTextContent('invalid');
		expect(errorDetails[1]).toHaveTextContent('tokens');
	});

	// TestEmptyUnexpectedTokensArray_Success
	// Tests behavior with empty unexpected tokens array
	it('TestEmptyUnexpectedTokensArray_Success: Handles empty unexpected tokens array', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [{ type: 'TEST', value: 'test' }],
			unexpected_tokens: []
		});

		// Error state should not be shown
		expect(screen.queryByText('Unidentified Token(s) Found')).not.toBeInTheDocument();
		
		// Token table should be shown
		expect(screen.getByText('TEST')).toBeInTheDocument();
		expect(screen.getByText('test')).toBeInTheDocument();
	});

	// TestNullUnexpectedTokens_Success  
	// Tests behavior when unexpected_tokens is null/undefined
	it('TestNullUnexpectedTokens_Success: Handles null unexpected tokens', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: [{ type: 'NULL_TEST', value: 'null_test' }],
			unexpected_tokens: undefined as any
		});

		// Error state should not be shown
		expect(screen.queryByText('Unidentified Token(s) Found')).not.toBeInTheDocument();
		
		// Token table should be shown
		expect(screen.getByText('NULL_TEST')).toBeInTheDocument();
	});

	// TestStoreStateWithNullValues_Success
	// Tests store subscription with null/undefined values through empty props
	it('TestStoreStateWithNullValues_Success: Handles empty state gracefully', () => {
		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: false,
			tokens: [],
			unexpected_tokens: []
		});

		// Should show empty state message
		expect(screen.getByText('Tokens will appear here after generation')).toBeInTheDocument();
	});

	// TestStoreStateWithUndefinedUnidentified_Success
	// Tests store handling when tokens_unidentified is undefined through props
	it('TestStoreStateWithUndefinedUnidentified_Success: Handles undefined tokens_unidentified gracefully', () => {
		const mockTokens: Token[] = [{ type: 'STORE_TEST', value: 'store_test' }];

		render(ArtifactViewer, {
			phase: 'lexer',
			show_tokens: true,
			tokens: mockTokens,
			unexpected_tokens: undefined as any
		});

		expect(screen.getByText('STORE_TEST')).toBeInTheDocument();
		expect(screen.queryByText('Unidentified Token(s) Found')).not.toBeInTheDocument();
	});

	// TestPhaseOtherThanLexer_Success
	// Tests that component doesn't render lexer content for other phases
	it('TestPhaseOtherThanLexer_Success: Does not render lexer-specific content for non-lexer phases', () => {
		render(ArtifactViewer, {
			phase: 'parser',
			show_tokens: true,
			tokens: [{ type: 'SHOULD_NOT_SHOW', value: 'should_not_show' }],
			unexpected_tokens: ['should_not_show_error']
		});

		// Headers are always rendered, but content should not be shown for non-lexer phases
		expect(screen.getByText('Lexer Artefact')).toBeInTheDocument(); // Always rendered
		expect(screen.getByText('Tokens')).toBeInTheDocument(); // Always rendered
		
		// These should NOT be rendered for non-lexer phases
		expect(screen.queryByText('Tokens will appear here after generation')).not.toBeInTheDocument();
		expect(screen.queryByText('SHOULD_NOT_SHOW')).not.toBeInTheDocument();
		expect(screen.queryByText('Unidentified Token(s) Found')).not.toBeInTheDocument();
	});
});



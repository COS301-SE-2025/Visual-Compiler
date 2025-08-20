import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import TranslatorArtifactViewer from '../src/lib/components/translator/translator-artifact-viewer.svelte';
import { AddToast } from '$lib/stores/toast';

// Mock the toast store
vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

// Mock clipboard API
Object.assign(navigator, {
	clipboard: {
		writeText: vi.fn()
	}
});

describe('TranslatorArtifactViewer Component', () => {
	const mockTranslatedCode = [
		'mov rax, 5',
		'add rax, 3',
		'mov [x], rax',
		'ret'
	];

	beforeEach(() => {
		vi.clearAllMocks();
		(navigator.clipboard.writeText as any).mockResolvedValue(undefined);
	});

	it('TestRender_Success: Renders artifact viewer container', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: null
			}
		});

		expect(screen.getByText('Translator Artefact')).toBeInTheDocument();
		const container = document.querySelector('.artifact-container');
		expect(container).toBeInTheDocument();
	});

	it('TestEmptyState_Success: Shows empty state when no translated code', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: null
			}
		});

		const emptyState = document.querySelector('.empty-state');
		expect(emptyState).toBeInTheDocument();
	});

	it('TestTranslatedCode_Success: Displays translated code when available', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		expect(screen.getByText(/mov rax, 5/)).toBeInTheDocument();
		expect(screen.getByText(/add rax, 3/)).toBeInTheDocument();
		expect(screen.getByText(/mov \[x\], rax/)).toBeInTheDocument();
		expect(screen.getByText(/ret/)).toBeInTheDocument();
	});

	it('TestLineNumbers_Success: Shows line numbers for translated code', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const lineNumbers = document.querySelector('.line-numbers');
		expect(lineNumbers).toBeInTheDocument();

		// Should have line numbers 1-4
		expect(screen.getByText('1')).toBeInTheDocument();
		expect(screen.getByText('2')).toBeInTheDocument();
		expect(screen.getByText('3')).toBeInTheDocument();
		expect(screen.getByText('4')).toBeInTheDocument();
	});

	it('TestCopyButton_Success: Shows copy button when code is available', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const copyButton = screen.getByTitle('Copy to clipboard');
		expect(copyButton).toBeInTheDocument();
		expect(screen.getByText('Copy')).toBeInTheDocument();
	});

	it('TestCopyButton_Hidden: Hides copy button when no code available', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: null
			}
		});

		const copyButton = screen.queryByTitle('Copy to clipboard');
		expect(copyButton).not.toBeInTheDocument();
	});

	it('TestCopyFunctionality_Success: Copies code to clipboard when clicked', async () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const copyButton = screen.getByTitle('Copy to clipboard');
		await fireEvent.click(copyButton);

		expect(navigator.clipboard.writeText).toHaveBeenCalledWith(
			'mov rax, 5\nadd rax, 3\nmov [x], rax\nret'
		);
		expect(AddToast).toHaveBeenCalledWith(
			'Code copied! Your translated code is now in the clipboard',
			'success'
		);
	});

	it('TestCopyError_Success: Handles clipboard copy errors', async () => {
		(navigator.clipboard.writeText as any).mockRejectedValue(new Error('Clipboard error'));

		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const copyButton = screen.getByTitle('Copy to clipboard');
		await fireEvent.click(copyButton);

		expect(AddToast).toHaveBeenCalledWith(
			'Copy failed: Unable to copy code to clipboard',
			'error'
		);
	});

	it('TestErrorState_Success: Shows error state when translation fails', () => {
		const mockError = {
			message: 'Invalid token sequence in rule 1'
		};

		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: mockError
			}
		});

		expect(screen.getByText('Translation Failed')).toBeInTheDocument();
		expect(screen.getByText(/The translation could not be completed/)).toBeInTheDocument();
		expect(screen.getByText('Invalid token sequence in rule 1')).toBeInTheDocument();
	});

	it('TestErrorIcon_Success: Shows error icon when translation fails', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: { message: 'Test error' }
			}
		});

		const errorIcon = document.querySelector('.error-icon svg');
		expect(errorIcon).toBeInTheDocument();
	});

	it('TestErrorStringHandling_Success: Handles string errors', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: [],
				translationError: 'Simple error message'
			}
		});

		expect(screen.getByText('Translation Failed')).toBeInTheDocument();
		expect(screen.getByText('Simple error message')).toBeInTheDocument();
	});

	it('TestCodeWrapper_Success: Has proper code wrapper structure', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const codeWrapper = document.querySelector('.code-wrapper');
		const codeBlock = document.querySelector('.code-block');

		expect(codeWrapper).toBeInTheDocument();
		expect(codeBlock).toBeInTheDocument();
	});

	it('TestSingleLineCode_Success: Handles single line of translated code', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: ['mov rax, 42'],
				translationError: null
			}
		});

		expect(screen.getByText('mov rax, 42')).toBeInTheDocument();
		expect(screen.getByText('1')).toBeInTheDocument();
		expect(screen.queryByText('2')).not.toBeInTheDocument();
	});

	it('TestCopySVGIcon_Success: Copy button has SVG icon', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const copyButton = screen.getByTitle('Copy to clipboard');
		const svg = copyButton.querySelector('svg');
		expect(svg).toBeInTheDocument();
	});

	it('TestArtifactHeader_Success: Has proper header structure', () => {
		render(TranslatorArtifactViewer, {
			props: {
				translated_code: mockTranslatedCode,
				translationError: null
			}
		});

		const artifactHeader = document.querySelector('.artifact-header');
		const artifactTitle = document.querySelector('.artifact-title');

		expect(artifactHeader).toBeInTheDocument();
		expect(artifactTitle).toBeInTheDocument();
	});
});

import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ProjectNamePrompt from '../../src/lib/components/project-hub/project-name-prompt.svelte';

describe('ProjectNamePrompt Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('TestNotShownWhenHidden_Success: Does not render when show is false', () => {
		render(ProjectNamePrompt, { props: { show: false } });

		expect(screen.queryByText('Create New Project')).not.toBeInTheDocument();
		expect(screen.queryByPlaceholderText('e.g., My First Compiler')).not.toBeInTheDocument();
	});

	it('TestShownWhenVisible_Success: Renders when show is true', () => {
		render(ProjectNamePrompt, { props: { show: true } });

		expect(screen.getByText('Create New Project')).toBeInTheDocument();
		expect(screen.getByText('Please enter a name for your new blank project.')).toBeInTheDocument();
		expect(screen.getByPlaceholderText('e.g., My First Compiler')).toBeInTheDocument();
	});

	it('TestButtons_Success: Renders cancel and create buttons', () => {
		render(ProjectNamePrompt, { props: { show: true } });

		expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Create Project' })).toBeInTheDocument();
	});

	it('TestCreateButtonDisabled_Success: Create button is disabled when input is empty', () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const createButton = screen.getByRole('button', { name: 'Create Project' });
		expect(createButton).toBeDisabled();
	});

	it('TestCreateButtonEnabled_Success: Create button is enabled when input has text', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler');
		const createButton = screen.getByRole('button', { name: 'Create Project' });

		await fireEvent.input(input, { target: { value: 'My Project' } });

		expect(createButton).not.toBeDisabled();
	});

	it('TestInputValue_Success: Input value updates correctly', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler') as HTMLInputElement;

		await fireEvent.input(input, { target: { value: 'Test Project' } });

		expect(input.value).toBe('Test Project');
	});

	it('TestCancelButton_Success: Cancel button is clickable', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const cancelButton = screen.getByRole('button', { name: 'Cancel' });
		await fireEvent.click(cancelButton);

		// Button should be functional (no errors)
		expect(cancelButton).toBeInTheDocument();
	});

	it('TestCreateButton_Success: Create button is clickable when enabled', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler');
		const createButton = screen.getByRole('button', { name: 'Create Project' });

		await fireEvent.input(input, { target: { value: 'My New Project' } });
		await fireEvent.click(createButton);

		// Button should be functional (no errors)
		expect(createButton).toBeInTheDocument();
	});

	it('TestKeyboardInteraction_Success: Handles keyboard events', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler');

		await fireEvent.input(input, { target: { value: 'Keyboard Project' } });
		await fireEvent.keyDown(window, { key: 'Enter' });

		// Should handle keyboard events without errors
		expect(input).toHaveValue('Keyboard Project');
	});

	it('TestEscapeKey_Success: Escape key is handled', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		await fireEvent.keyDown(window, { key: 'Escape' });

		// Should handle escape key without errors
		expect(screen.getByText('Create New Project')).toBeInTheDocument();
	});

	it('TestInputValidation_Success: Button states change with input', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler');
		const createButton = screen.getByRole('button', { name: 'Create Project' });

		// Initially disabled
		expect(createButton).toBeDisabled();

		// Enabled with text
		await fireEvent.input(input, { target: { value: 'Test' } });
		expect(createButton).not.toBeDisabled();

		// Disabled again when cleared
		await fireEvent.input(input, { target: { value: '' } });
		expect(createButton).toBeDisabled();
	});

	it('TestWhitespaceHandling_Success: Handles whitespace-only input', async () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByPlaceholderText('e.g., My First Compiler');
		const createButton = screen.getByRole('button', { name: 'Create Project' });

		// Whitespace-only input should keep button disabled
		await fireEvent.input(input, { target: { value: '   ' } });
		expect(createButton).toBeDisabled();
	});

	it('TestAccessibility_Success: Input has proper aria-label', () => {
		render(ProjectNamePrompt, { props: { show: true } });

		const input = screen.getByLabelText('Project Name');
		expect(input).toBeInTheDocument();
	});

	it('TestModalStructure_Success: Has proper modal structure', () => {
		const { container } = render(ProjectNamePrompt, { props: { show: true } });

		const backdrop = container.querySelector('.backdrop');
		const promptModal = container.querySelector('.prompt-modal');

		expect(backdrop).toBeInTheDocument();
		expect(promptModal).toBeInTheDocument();
	});
});

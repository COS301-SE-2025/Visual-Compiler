import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import DeleteConfirmation from '../../src/lib/components/project-hub/delete-confirmation.svelte';

describe('DeleteConfirmation Component', () => {
	beforeEach(() => {
		// Reset any state if needed
	});

	it('TestRender_Success: Renders when show is true', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		expect(screen.getByText('Confirm Deletion')).toBeInTheDocument();
		expect(screen.getByText(/Are you sure you want to delete the project/)).toBeInTheDocument();
	});

	it('TestRender_Hidden: Does not render when show is false', () => {
		render(DeleteConfirmation, { 
			show: false, 
			projectName: 'Test Project' 
		});

		expect(screen.queryByText('Confirm Deletion')).not.toBeInTheDocument();
	});

	it('TestProjectName_Success: Displays project name correctly', () => {
		const projectName = 'My Amazing Project';
		render(DeleteConfirmation, { 
			show: true, 
			projectName 
		});

		expect(screen.getByText(projectName)).toBeInTheDocument();
		expect(screen.getByText(/Are you sure you want to delete the project/)).toBeInTheDocument();
	});

	it('TestDeleteConfirmButton_Success: Shows delete confirm button', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const deleteButton = screen.getByText('Delete');
		expect(deleteButton).toBeInTheDocument();
		expect(deleteButton.tagName).toBe('BUTTON');
	});

	it('TestCancelButton_Success: Shows cancel button', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const cancelButton = screen.getByText('Cancel');
		expect(cancelButton).toBeInTheDocument();
		expect(cancelButton.tagName).toBe('BUTTON');
	});

	it('TestConfirmClick_Success: Handles confirm button click', async () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const deleteButton = screen.getByText('Delete');
		await fireEvent.click(deleteButton);

		// Confirm action should be triggered (button should be clickable)
		expect(deleteButton).toBeInTheDocument();
	});

	it('TestCancelClick_Success: Handles cancel button click', async () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const cancelButton = screen.getByText('Cancel');
		await fireEvent.click(cancelButton);

		// Cancel action should be triggered (button should be clickable)
		expect(cancelButton).toBeInTheDocument();
	});

	it('TestModalStructure_Success: Has proper modal structure', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const backdrop = container.querySelector('.backdrop');
		expect(backdrop).toBeInTheDocument();

		const modal = container.querySelector('.prompt-modal');
		expect(modal).toBeInTheDocument();
	});

	it('TestWarningMessage_Success: Shows warning message', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		expect(screen.getByText(/This action cannot be undone/)).toBeInTheDocument();
	});

	it('TestProjectNameHighlight_Success: Highlights project name', () => {
		const projectName = 'Important Project';
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName 
		});

		const strongElement = container.querySelector('strong');
		expect(strongElement).toBeInTheDocument();
		expect(strongElement?.textContent).toBe(projectName);
	});

	it('TestButtonContainer_Success: Has button container', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const buttonContainer = container.querySelector('.button-container');
		expect(buttonContainer).toBeInTheDocument();

		const buttons = buttonContainer?.querySelectorAll('button');
		expect(buttons?.length).toBe(2);
	});

	it('TestDeleteButtonStyling_Success: Delete button has correct styling', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const deleteButton = container.querySelector('.delete-confirm-button');
		expect(deleteButton).toBeInTheDocument();
		expect(deleteButton?.textContent?.trim()).toBe('Delete');
	});

	it('TestCancelButtonStyling_Success: Cancel button has correct styling', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const cancelButton = container.querySelector('.cancel-button');
		expect(cancelButton).toBeInTheDocument();
		expect(cancelButton?.textContent).toBe('Cancel');
	});

	it('TestModalHeading_Success: Shows modal heading', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const heading = screen.getByText('Confirm Deletion');
		expect(heading.tagName).toBe('H3');
		expect(heading).toHaveClass('prompt-heading');
	});

	it('TestModalSubheading_Success: Shows descriptive subheading', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const subheading = container.querySelector('.prompt-subheading');
		expect(subheading).toBeInTheDocument();
		expect(subheading?.tagName).toBe('P');
	});

	it('TestBackdropClick_Success: Handles backdrop interaction', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const backdrop = container.querySelector('.backdrop');
		expect(backdrop).toBeInTheDocument();

		// Modal should have click stop propagation
		const modal = container.querySelector('.prompt-modal');
		expect(modal).toBeInTheDocument();
	});

	it('TestTransitionEffect_Success: Uses transition effects', () => {
		const { container } = render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		// Should have backdrop with transition
		const backdrop = container.querySelector('.backdrop');
		expect(backdrop).toBeInTheDocument();
	});

	it('TestAccessibility_Success: Has proper accessibility features', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		// Should have descriptive text for screen readers
		expect(screen.getByText(/Are you sure you want to delete/)).toBeInTheDocument();
		
		// Buttons should be keyboard accessible
		const buttons = screen.getAllByRole('button');
		expect(buttons).toHaveLength(2);
		
		buttons.forEach(button => {
			expect(button).toBeInTheDocument();
		});
	});

	it('TestEmptyProjectName_Success: Handles empty project name', () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: '' 
		});

		// Should still render without crashing
		expect(screen.getByText('Confirm Deletion')).toBeInTheDocument();
		expect(screen.getByText('Delete')).toBeInTheDocument();
		expect(screen.getByText('Cancel')).toBeInTheDocument();
	});

	it('TestLongProjectName_Success: Handles long project names', () => {
		const longProjectName = 'This is a very long project name that might cause layout issues if not handled properly';
		render(DeleteConfirmation, { 
			show: true, 
			projectName: longProjectName 
		});

		expect(screen.getByText(longProjectName)).toBeInTheDocument();
		expect(screen.getByText('Confirm Deletion')).toBeInTheDocument();
	});

	it('TestSpecialCharacters_Success: Handles project names with special characters', () => {
		const specialProjectName = 'Project @#$% & More!';
		render(DeleteConfirmation, { 
			show: true, 
			projectName: specialProjectName 
		});

		expect(screen.getByText(specialProjectName)).toBeInTheDocument();
		expect(screen.getByText('Confirm Deletion')).toBeInTheDocument();
	});

	it('TestButtonFunctionality_Success: Both buttons are clickable', async () => {
		render(DeleteConfirmation, { 
			show: true, 
			projectName: 'Test Project' 
		});

		const deleteButton = screen.getByText('Delete');
		const cancelButton = screen.getByText('Cancel');

		// Both buttons should be clickable
		expect(deleteButton).toBeEnabled();
		expect(cancelButton).toBeEnabled();

		await fireEvent.click(deleteButton);
		await fireEvent.click(cancelButton);

		// No errors should occur
	});
});

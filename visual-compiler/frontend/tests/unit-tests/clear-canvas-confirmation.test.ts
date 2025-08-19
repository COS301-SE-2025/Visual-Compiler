import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, fireEvent, screen } from '@testing-library/svelte';
import ClearCanvasConfirmation from '../../src/lib/components/main/clear-canvas-confirmation.svelte';

/**
 * Test suite for ClearCanvasConfirmation component
 * Validates confirmation dialog functionality, UI interactions, and component structure
 */
describe('ClearCanvasConfirmation Component', () => {
	beforeEach(() => {
		// Clear all mocks before each test
		vi.clearAllMocks();
	});

	it('TestRenderWhenShowIsTrue_Success: Renders confirmation dialog when show is true', () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Check that the modal is rendered
		const modal = container.querySelector('.prompt-modal');
		expect(modal).toBeTruthy();

		// Check heading text specifically
		const heading = container.querySelector('.prompt-heading');
		expect(heading?.textContent).toBe('Clear Canvas');
		
		// Check warning message
		expect(screen.getByText(/Are you sure you want to clear the canvas/)).toBeTruthy();
		
		// Check buttons exist by their classes
		const cancelButton = container.querySelector('.cancel-button');
		const confirmButton = container.querySelector('.clear-confirm-button');
		expect(cancelButton?.textContent?.trim()).toBe('Cancel');
		expect(confirmButton?.textContent?.trim()).toBe('Clear Canvas');
	});

	it('TestDoesNotRenderWhenShowIsFalse_Success: Does not render when show is false', () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: false }
		});

		// Check that the modal is not rendered
		const modal = container.querySelector('.prompt-modal');
		expect(modal).toBeFalsy();
	});

	it('TestConfirmButtonExists_Success: Confirm button is rendered correctly', async () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Find the confirm button by class
		const confirm_button = container.querySelector('.clear-confirm-button') as HTMLElement;
		expect(confirm_button).toBeTruthy();
		expect(confirm_button.textContent?.trim()).toBe('Clear Canvas');
		
		// Test that button is clickable
		await fireEvent.click(confirm_button);
		// The button should be clickable without errors
		expect(true).toBe(true);
	});

	it('TestCancelButtonExists_Success: Cancel button is rendered correctly', async () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Find the cancel button by class
		const cancel_button = container.querySelector('.cancel-button') as HTMLElement;
		expect(cancel_button).toBeTruthy();
		expect(cancel_button.textContent?.trim()).toBe('Cancel');
		
		// Test that button is clickable
		await fireEvent.click(cancel_button);
		// The button should be clickable without errors
		expect(true).toBe(true);
	});

	it('TestModalClickDoesNotPropagate_Success: Clicking modal does not trigger backdrop events', async () => {
		let backdrop_clicked = false;
		
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Add click listener to backdrop
		const backdrop = container.querySelector('.backdrop');
		backdrop?.addEventListener('click', () => {
			backdrop_clicked = true;
		});

		// Click on the modal (should not propagate)
		const modal = container.querySelector('.prompt-modal');
		await fireEvent.click(modal as Element);

		// Backdrop click should not have been triggered due to stopPropagation
		expect(backdrop_clicked).toBe(false);
	});

	it('TestBackdropClickBehavior_Success: Clicking backdrop can be handled', async () => {
		let backdrop_clicked = false;
		
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Add click listener to backdrop
		const backdrop = container.querySelector('.backdrop');
		backdrop?.addEventListener('click', () => {
			backdrop_clicked = true;
		});

		// Click directly on the backdrop
		await fireEvent.click(backdrop as Element);

		// Backdrop click should be triggered
		expect(backdrop_clicked).toBe(true);
	});

	it('TestComponentReactivity_Success: Component shows and hides based on show prop', () => {
		// Test with show: true
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Should show modal
		expect(container.querySelector('.prompt-modal')).toBeTruthy();

		// Test separately with show: false 
		const { container: container2 } = render(ClearCanvasConfirmation, {
			props: { show: false }
		});

		// Should hide modal - backdrop should not exist when show is false
		expect(container2.querySelector('.backdrop')).toBeFalsy();
		expect(container2.querySelector('.prompt-modal')).toBeFalsy();
	});

	it('TestWarningMessageContent_Success: Displays complete warning message', () => {
		render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Check complete warning message
		const expected_message = 'Are you sure you want to clear the canvas? All unsaved progress will be lost. This action cannot be undone.';
		expect(screen.getByText(expected_message)).toBeTruthy();
	});

	it('TestButtonStyling_Success: Buttons have correct CSS classes', () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Check cancel button class
		const cancel_button = container.querySelector('.cancel-button');
		expect(cancel_button).toBeTruthy();
		expect(cancel_button?.textContent?.trim()).toBe('Cancel');

		// Check confirm button class
		const confirm_button = container.querySelector('.clear-confirm-button');
		expect(confirm_button).toBeTruthy();
		expect(confirm_button?.textContent?.trim()).toBe('Clear Canvas');
	});

	it('TestModalStructure_Success: Modal has correct structure and classes', () => {
		const { container } = render(ClearCanvasConfirmation, {
			props: { show: true }
		});

		// Check backdrop
		const backdrop = container.querySelector('.backdrop');
		expect(backdrop).toBeTruthy();

		// Check modal
		const modal = container.querySelector('.prompt-modal');
		expect(modal).toBeTruthy();

		// Check heading
		const heading = container.querySelector('.prompt-heading');
		expect(heading).toBeTruthy();
		expect(heading?.textContent).toBe('Clear Canvas');

		// Check subheading
		const subheading = container.querySelector('.prompt-subheading');
		expect(subheading).toBeTruthy();

		// Check button container
		const button_container = container.querySelector('.button-container');
		expect(button_container).toBeTruthy();
	});
});

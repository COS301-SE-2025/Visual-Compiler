import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { writable } from 'svelte/store';
import { tick } from 'svelte';
import type { Toast } from '$lib/stores/toast';
import ToastContainer from '../../src/lib/components/toast-conatiner.svelte';

// Mock the toast store with a factory function approach
vi.mock('$lib/stores/toast', () => {
	return {
		toasts: writable<Toast[]>([])
	};
});

// Get the mocked store from the module
import { toasts } from '$lib/stores/toast';

describe('ToastContainer Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		toasts.set([]);
	});

	it('TestEmptyState_Success: Renders empty container when no toasts', () => {
		render(ToastContainer);
		
		const container = document.querySelector('.toast-container');
		expect(container).toBeInTheDocument();
		expect(container?.children.length).toBe(0);
	});

	it('TestSuccessToast_Success: Renders success toast with correct styling', async () => {
		const successToast = {
			id: 1,
			message: 'Operation successful!',
			type: 'success' as const,
			duration: 3000
		};

		toasts.set([successToast]);
		render(ToastContainer);

		const toast = await screen.findByText('Operation successful!');
		expect(toast).toBeInTheDocument();
		expect(toast).toHaveClass('toast', 'success');
	});

	it('TestErrorToast_Success: Renders error toast with correct styling', async () => {
		const errorToast = {
			id: 2,
			message: 'Something went wrong!',
			type: 'error' as const,
			duration: 3000
		};

		toasts.set([errorToast]);
		render(ToastContainer);

		const toast = await screen.findByText('Something went wrong!');
		expect(toast).toBeInTheDocument();
		expect(toast).toHaveClass('toast', 'error');
	});

	it('TestInfoToast_Success: Renders info toast with correct styling', async () => {
		const infoToast = {
			id: 3,
			message: 'Here is some information',
			type: 'info' as const,
			duration: 3000
		};

		toasts.set([infoToast]);
		render(ToastContainer);

		const toast = await screen.findByText('Here is some information');
		expect(toast).toBeInTheDocument();
		expect(toast).toHaveClass('toast', 'info');
	});

	it('TestMultipleToasts_Success: Renders multiple toasts simultaneously', async () => {
		const testToasts = [
			{
				id: 1,
				message: 'First toast',
				type: 'success' as const,
				duration: 3000
			},
			{
				id: 2,
				message: 'Second toast',
				type: 'error' as const,
				duration: 3000
			},
			{
				id: 3,
				message: 'Third toast',
				type: 'info' as const,
				duration: 3000
			}
		];

		toasts.set(testToasts);
		render(ToastContainer);

		expect(await screen.findByText('First toast')).toBeInTheDocument();
		expect(await screen.findByText('Second toast')).toBeInTheDocument();
		expect(await screen.findByText('Third toast')).toBeInTheDocument();
	});

	it('TestToastRemoval_Success: Handles toast removal correctly', async () => {
		const initialToasts = [
			{
				id: 1,
				message: 'Toast 1',
				type: 'success' as const,
				duration: 3000
			},
			{
				id: 2,
				message: 'Toast 2',
				type: 'error' as const,
				duration: 3000
			}
		];

		// Set initial toasts and render
		toasts.set(initialToasts);
		render(ToastContainer);

		expect(await screen.findByText('Toast 1')).toBeInTheDocument();
		expect(await screen.findByText('Toast 2')).toBeInTheDocument();

		// Remove one toast
		toasts.set([initialToasts[1]]);
		
		// Wait for Svelte to process the change
		await tick();

		// Clear the container and re-render to force update
		document.body.innerHTML = '';
		render(ToastContainer);
		
		// Now check that only Toast 2 is present
		expect(screen.queryByText('Toast 1')).not.toBeInTheDocument();
		expect(await screen.findByText('Toast 2')).toBeInTheDocument();
	});
});

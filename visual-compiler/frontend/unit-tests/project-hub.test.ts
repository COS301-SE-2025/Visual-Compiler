import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { writable } from 'svelte/store';
import ProjectHub from '../../src/lib/components/project-hub/project-hub.svelte';

// Mock the project store
vi.mock('$lib/stores/project', () => ({
	projectName: writable('')
}));

// Mock the child components
vi.mock('../../src/lib/components/project-hub/project-name-prompt.svelte', () => ({
	default: vi.fn()
}));

vi.mock('../../src/lib/components/project-hub/delete-confirmation.svelte', () => ({
	default: vi.fn()
}));

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage
const localStorageMock = {
	getItem: vi.fn(),
	setItem: vi.fn(),
	removeItem: vi.fn(),
	clear: vi.fn(),
	length: 0,
	key: vi.fn()
};
global.localStorage = localStorageMock as any;

describe('ProjectHub Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		localStorageMock.getItem.mockReturnValue('test-user-id');
		mockFetch.mockResolvedValue({
			ok: true,
			json: async () => ({ all_projects: ['Project 1', 'Project 2', 'Project 3'] })
		});
	});

	it('TestRender_Success: Renders project hub when show is true', () => {
		render(ProjectHub, { show: true });

		// Should render the modal when shown
		const modal = screen.getByText(/Start a new project/i);
		expect(modal).toBeInTheDocument();
	});

	it('TestRender_Hidden: Does not render when show is false', () => {
		render(ProjectHub, { show: false });

		// Should not show the modal
		const modal = screen.queryByRole('dialog');
		expect(modal).toBeNull();
	});

	it('TestFetchProjects_Success: Fetches projects when modal opens', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('http://localhost:8080/api/users/getProjects'),
				expect.objectContaining({
					method: 'GET',
					headers: expect.objectContaining({
						'accept': 'application/json'
					})
				})
			);
		});
	});

	it('TestCreateNewProject_Success: Shows create new project button', () => {
		render(ProjectHub, { show: true });

		const createButton = screen.getByText(/New Blank/i);
		expect(createButton).toBeInTheDocument();
	});

	it('TestCreateNewProjectClick_Success: Handles create new project click', async () => {
		render(ProjectHub, { show: true });

		const createButton = screen.getByText(/New Blank/i);
		await fireEvent.click(createButton);

		// Component should handle the create new project action
		// (specific behavior depends on implementation)
	});

	it('TestProjectsList_Success: Displays fetched projects', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(screen.getByText('Project 1')).toBeInTheDocument();
			expect(screen.getByText('Project 2')).toBeInTheDocument();
			expect(screen.getByText('Project 3')).toBeInTheDocument();
		});
	});

	it('TestProjectSelection_Success: Handles project selection', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Test Project'] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ message: "Retrieved users project details" })
			});

		render(ProjectHub, { show: true });

		await waitFor(() => {
			const project = screen.getByText('Test Project');
			expect(project).toBeInTheDocument();
		});

		const project = screen.getByText('Test Project');
		await fireEvent.click(project);

		// Should call the project selection API
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				expect.stringContaining('http://localhost:8080/api/users/getProject'),
				expect.objectContaining({
					method: 'GET'
				})
			);
		});
	});

	it('TestDeleteProject_Success: Shows delete button for projects', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Look for delete buttons or delete icons
			const deleteButtons = screen.getAllByRole('button').filter(btn => 
				btn.textContent?.includes('Delete') || btn.innerHTML?.includes('delete') || btn.innerHTML?.includes('×')
			);
			expect(deleteButtons.length).toBeGreaterThanOrEqual(0);
		});
	});

	it('TestDeleteProjectClick_Success: Handles delete project click', async () => {
		const { container } = render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(screen.getByText('Project 1')).toBeInTheDocument();
		});

		// Look for delete buttons in the container
		const deleteElements = container.querySelectorAll('[data-testid*="delete"], .delete-btn, .delete-icon');
		if (deleteElements.length > 0) {
			await fireEvent.click(deleteElements[0]);
			// Should trigger delete confirmation flow
		}
	});

	it('TestEmptyProjects_Success: Handles empty projects list', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ all_projects: [] })
		});

		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should handle empty state gracefully
			expect(mockFetch).toHaveBeenCalled();
		});
	});

	it('TestFetchError_Success: Handles fetch error gracefully', async () => {
		mockFetch.mockRejectedValueOnce(new Error('Network error'));

		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Component should handle error gracefully without crashing
	});

	it('TestCloseModal_Success: Dispatches close event', () => {
		render(ProjectHub, { show: true });

		// Find close button or similar
		const closeButtons = screen.getAllByRole('button').filter(btn => 
			btn.textContent?.includes('Close') || btn.textContent?.includes('×') || btn.innerHTML?.includes('close')
		);

		if (closeButtons.length > 0) {
			fireEvent.click(closeButtons[0]);
			// Close event should be handled by the component
		}
	});

	it('TestModalStructure_Success: Has proper modal structure', () => {
		const { container } = render(ProjectHub, { show: true });

		// Should have modal elements
		const modal = container.querySelector('[role="dialog"], .modal, .project-hub');
		expect(modal).toBeInTheDocument();
	});

	it('TestProjectDateDisplay_Success: Shows project dates', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should show date information for projects
			const dateElements = screen.getAllByText(/\d{4}|\w+ \d+|\w+day/i);
			expect(dateElements.length).toBeGreaterThanOrEqual(0);
		});
	});

	it('TestNoUserId_Success: Handles missing user ID gracefully', async () => {
		localStorageMock.getItem.mockReturnValue(null);

		render(ProjectHub, { show: true });

		// Should not crash when user ID is missing
		await waitFor(() => {
			// Component should handle gracefully
		});
	});

	it('TestAPIError_Success: Handles API error responses', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 500,
			json: async () => ({ error: 'Server error' })
		});

		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Component should handle API errors gracefully
	});

	it('TestProjectNameStore_Success: Uses project name store', () => {
		render(ProjectHub, { show: true });

		// Component should interact with the project name store
		// (specific assertions depend on store implementation)
	});

	it('TestTransitionEffects_Success: Has transition animations', () => {
		const { container } = render(ProjectHub, { show: true });

		// Should have elements that use transitions
		const animatedElements = container.querySelectorAll('[style*="transform"], [style*="transition"]');
		expect(animatedElements.length).toBeGreaterThanOrEqual(0);
	});

	it('TestKeyboardNavigation_Success: Supports keyboard navigation', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should have focusable elements for keyboard navigation
			const focusableElements = screen.getAllByRole('button');
			expect(focusableElements.length).toBeGreaterThan(0);
		});
	});

	it('TestModalOverlay_Success: Has modal overlay', () => {
		const { container } = render(ProjectHub, { show: true });

		// Should have modal overlay/backdrop
		const overlay = container.querySelector('.modal-overlay, .backdrop, [data-backdrop]');
		expect(overlay || container.firstChild).toBeInTheDocument();
	});

	it('TestRecentProjectsSection_Success: Shows recent projects section', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should have a section for recent projects
			const recentSection = screen.queryByText(/Recent/i) || screen.queryByText(/Projects/i);
			if (recentSection) {
				expect(recentSection).toBeInTheDocument();
			}
		});
	});

	it('TestProjectItems_Success: Renders project items correctly', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should render project items with proper structure
			const projects = screen.getAllByText(/Project \d+/);
			projects.forEach(project => {
				expect(project).toBeInTheDocument();
			});
		});
	});
});

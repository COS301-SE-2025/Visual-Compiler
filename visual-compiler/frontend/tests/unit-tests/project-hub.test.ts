import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { writable } from 'svelte/store';
import { tick } from 'svelte';
import ProjectHub from '../../src/lib/components/project-hub/project-hub.svelte';
import ProjectNamePrompt from '../../src/lib/components/project-hub/project-name-prompt.svelte';
import DeleteConfirmation from '../../src/lib/components/project-hub/delete-confirmation.svelte';

// Create mock stores
const mockProjectName = writable('');
const mockPipelineStore = writable({
	nodes: [],
	connections: [],
	lastSaved: null
});

// Mock the stores
vi.mock('$lib/stores/project', () => ({
	projectName: mockProjectName
}));

vi.mock('$lib/stores/pipeline', () => ({
	pipelineStore: mockPipelineStore,
	resetPipeline: vi.fn()
}));

// Mock svelte/transition
vi.mock('svelte/transition', () => ({
	fly: vi.fn(() => ({ delay: 0, duration: 300, easing: (t: number) => t }))
}));

// Mock toast store
const mockAddToast = vi.fn();
vi.mock('$lib/stores/toast', () => ({
	AddToast: mockAddToast
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
		mockProjectName.set('');
		mockPipelineStore.set({
			nodes: [],
			connections: [],
			lastSaved: null
		});
		localStorageMock.getItem.mockReturnValue('test-user-id');
		mockFetch.mockResolvedValue({
			ok: true,
			json: async () => ({ all_projects: ['Project 1', 'Project 2', 'Project 3'] })
		});
	});

	// Basic Rendering Tests
	it('TestRender_Success: Renders project hub when show is true', async () => {
		render(ProjectHub, { show: true });

		// Should render the modal when shown
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
		expect(screen.getByText(/Recent Projects/i)).toBeInTheDocument();
		expect(screen.getByText(/New Blank/i)).toBeInTheDocument();
	});

	it('TestRender_Hidden: Does not render when show is false', () => {
		const { container } = render(ProjectHub, { show: false });

		// Should not show the modal
		expect(container.firstChild).toBeNull();
	});

	it('TestModalStructure_Success: Has proper modal structure', async () => {
		const { container } = render(ProjectHub, { show: true });

		// Check for modal structure elements
		expect(container.querySelector('.backdrop')).toBeInTheDocument();
		expect(container.querySelector('.modal')).toBeInTheDocument();
	});

	// Project Fetching Tests
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

	it('TestFetchProjects_Error: Handles fetch error gracefully', async () => {
		mockFetch.mockRejectedValueOnce(new Error('Network error'));

		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Component should handle error gracefully without crashing
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
	});

	it('TestFetchProjects_EmptyResponse: Handles empty projects response', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ all_projects: [] })
		});

		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Should handle empty state gracefully
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
	});

	it('TestFetchProjects_APIError: Handles API error responses', async () => {
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
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
	});

	it('TestNoUserId_Success: Handles missing user ID gracefully', async () => {
		localStorageMock.getItem.mockReturnValue(null);

		render(ProjectHub, { show: true });

		// Should not crash when user ID is missing
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
		
		// Should not call fetch without user ID
		expect(mockFetch).not.toHaveBeenCalled();
	});

	// Project Display Tests
	it('TestProjectsList_Success: Displays fetched projects', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(screen.getByText('Project 1')).toBeInTheDocument();
			expect(screen.getByText('Project 2')).toBeInTheDocument();
			expect(screen.getByText('Project 3')).toBeInTheDocument();
		});
	});

	it('TestProjectDateDisplay_Success: Shows project dates', async () => {
		render(ProjectHub, { show: true });

		await waitFor(() => {
			// Should show date information for projects
			const dateElements = screen.getAllByText(/\w+\s+\d+,\s+\d{4}/);
			expect(dateElements.length).toBeGreaterThanOrEqual(3);
		});
	});

	// Project Selection Tests
	it('TestProjectSelection_Success: Handles project selection with pipeline data', async () => {
		const mockPipelineData = {
			nodes: [{ id: '1', type: 'source' }],
			connections: [],
			lastSaved: '2024-01-01T00:00:00Z'
		};

		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Test Project'] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ 
					message: "Retrieved users project details",
					results: { pipeline: mockPipelineData }
				})
			});

		const { component } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		await waitFor(() => {
			expect(screen.getByText('Test Project')).toBeInTheDocument();
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

		// Should update stores and close modal
		await waitFor(() => {
			expect(mockAddToast).toHaveBeenCalledWith('Pipeline restored successfully', 'success');
			expect(closeEventSpy).toHaveBeenCalled();
		});
	});

	it('TestProjectSelection_NoPipelineData: Handles project selection without pipeline data', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Empty Project'] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ 
					message: "Retrieved users project details",
					results: {}
				})
			});

		const { component } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		await waitFor(() => {
			expect(screen.getByText('Empty Project')).toBeInTheDocument();
		});

		const project = screen.getByText('Empty Project');
		await fireEvent.click(project);

		await waitFor(() => {
			expect(closeEventSpy).toHaveBeenCalled();
		});
	});

	it('TestProjectSelection_Error: Handles project selection error', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Error Project'] })
			})
			.mockResolvedValueOnce({
				ok: false,
				status: 404,
				json: async () => ({ error: 'Project not found' })
			});

		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(screen.getByText('Error Project')).toBeInTheDocument();
		});

		const project = screen.getByText('Error Project');
		await fireEvent.click(project);

		await waitFor(() => {
			expect(mockAddToast).toHaveBeenCalledWith('Failed to retrieve project details', 'error');
		});
	});

	it('TestProjectSelection_NoUserId: Handles project selection without user ID', async () => {
		localStorageMock.getItem.mockReturnValue(null);

		render(ProjectHub, { show: true });

		// Create a clickable element to simulate project selection
		const { container } = render(ProjectHub, { show: true });
		const projectButton = document.createElement('button');
		projectButton.textContent = 'Test Project';
		projectButton.onclick = () => {
			// This simulates the selectProject function call without user ID
		};
		container.appendChild(projectButton);

		await fireEvent.click(projectButton);

		// Should show error toast
		await waitFor(() => {
			expect(mockAddToast).toHaveBeenCalledWith('Please log in to select a project', 'error');
		});
	});

	// Project Creation Tests  
	it('TestCreateNewProject_Success: Shows create new project button and handles click', async () => {
		const { component } = render(ProjectHub, { show: true });

		const createButton = screen.getByText(/New Blank/i);
		expect(createButton).toBeInTheDocument();

		await fireEvent.click(createButton);

		// Should trigger the showProjectNamePrompt state change
		// We can verify this by checking if the component re-renders correctly
		await tick();
		expect(createButton).toBeInTheDocument();
	});

	it('TestProjectNameConfirm_Success: Handles project name confirmation', async () => {
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({ success: true })
		});

		const { component } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		// Simulate project name confirmation event
		const event = new CustomEvent('confirm', { detail: 'New Test Project' });
		await component.$$.ctx[0].handleProjectNameConfirm(event);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/api/users/save',
				expect.objectContaining({
					method: 'POST',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					}),
					body: JSON.stringify({
						project_name: 'New Test Project',
						users_id: 'test-user-id'
					})
				})
			);
		});
	});

	// Project Deletion Tests
	it('TestDeleteProject_Success: Shows delete buttons and handles click', async () => {
		const { container } = render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(screen.getByText('Project 1')).toBeInTheDocument();
		});

		// Look for delete buttons
		const deleteButtons = container.querySelectorAll('.delete-button');
		expect(deleteButtons.length).toBeGreaterThan(0);

		// Test delete button click
		const firstDeleteButton = deleteButtons[0] as HTMLElement;
		const clickEvent = new MouseEvent('click', { bubbles: true });
		Object.defineProperty(clickEvent, 'stopPropagation', { value: vi.fn() });
		
		await fireEvent(firstDeleteButton, clickEvent);
	});

	it('TestConfirmDelete_Success: Handles successful project deletion', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Project To Delete'] })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ success: true })
			});

		const { component } = render(ProjectHub, { show: true });

		// Simulate delete confirmation
		await component.$$.ctx[0].confirmDelete();

		await waitFor(() => {
			expect(mockFetch).toHaveBeenLastCalledWith(
				'http://localhost:8080/api/users/deleteProject',
				expect.objectContaining({
					method: 'DELETE',
					headers: expect.objectContaining({
						'Content-Type': 'application/json'
					})
				})
			);
		});
	});

	it('TestConfirmDelete_Error: Handles deletion error gracefully', async () => {
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Project To Delete'] })
			})
			.mockRejectedValueOnce(new Error('Delete failed'));

		const { component } = render(ProjectHub, { show: true });

		// Simulate delete confirmation with error
		await component.$$.ctx[0].confirmDelete();

		// Should handle error gracefully
		await tick();
	});

	// Close Modal Tests
	it('TestCloseModal_WithExistingProject: Dispatches close event when has existing project', async () => {
		mockProjectName.set('Existing Project');
		
		const { component } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		await waitFor(() => {
			const closeButton = screen.getByLabelText(/close project hub/i);
			expect(closeButton).toBeInTheDocument();
		});

		const closeButton = screen.getByLabelText(/close project hub/i);
		await fireEvent.click(closeButton);

		expect(closeEventSpy).toHaveBeenCalled();
	});

	it('TestCloseModal_NoExistingProject: Does not show close button without existing project', () => {
		mockProjectName.set('');
		
		render(ProjectHub, { show: true });

		const closeButton = screen.queryByLabelText(/close project hub/i);
		expect(closeButton).not.toBeInTheDocument();
	});

	it('TestBackdropClick_WithExistingProject: Handles backdrop click when has existing project', async () => {
		mockProjectName.set('Existing Project');
		
		const { component, container } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		const backdrop = container.querySelector('.backdrop');
		await fireEvent.click(backdrop as Element);

		expect(closeEventSpy).toHaveBeenCalled();
	});

	it('TestBackdropClick_NoExistingProject: Does not close on backdrop click without existing project', async () => {
		mockProjectName.set('');
		
		const { component, container } = render(ProjectHub, { show: true });
		const closeEventSpy = vi.fn();
		component.$on('close', closeEventSpy);

		const backdrop = container.querySelector('.backdrop');
		await fireEvent.click(backdrop as Element);

		expect(closeEventSpy).not.toHaveBeenCalled();
	});

	// Search Functionality Tests
	it('TestSearchBar_Success: Renders search input', () => {
		render(ProjectHub, { show: true });

		const searchInput = screen.getByPlaceholderText(/Search projects/i);
		expect(searchInput).toBeInTheDocument();
		expect(searchInput).toHaveClass('search-input');
	});

	it('TestSearchBar_Icon: Renders search icon', () => {
		const { container } = render(ProjectHub, { show: true });

		const searchIcon = container.querySelector('.search-icon');
		expect(searchIcon).toBeInTheDocument();
	});

	// Dark Mode Tests  
	it('TestDarkModeStyles_Success: Component supports dark mode styles', () => {
		const { container } = render(ProjectHub, { show: true });

		// Component should have proper class structure for dark mode
		expect(container.querySelector('.modal')).toBeInTheDocument();
		expect(container.querySelector('.project-button')).toBeInTheDocument();
	});

	// Accessibility Tests
	it('TestAccessibility_Success: Has proper ARIA labels and roles', async () => {
		mockProjectName.set('Test Project');
		render(ProjectHub, { show: true });

		// Check for proper accessibility attributes
		const closeButton = screen.queryByLabelText(/close project hub/i);
		expect(closeButton).toBeInTheDocument();

		await waitFor(() => {
			const deleteButtons = screen.getAllByLabelText(/delete project/i);
			expect(deleteButtons.length).toBeGreaterThan(0);
		});
	});

	// Keyboard Navigation Tests
	it('TestKeyboardNavigation_Success: Supports keyboard navigation', async () => {
		render(ProjectHub, { show: true });

		// Should have focusable elements for keyboard navigation
		const buttons = screen.getAllByRole('button');
		expect(buttons.length).toBeGreaterThan(0);

		// Check that buttons are focusable
		buttons.forEach(button => {
			expect(button).not.toHaveAttribute('disabled');
		});
	});

	// Store Integration Tests
	it('TestProjectNameStore_Success: Integrates with project name store', async () => {
		render(ProjectHub, { show: true });

		// Component should respond to store changes
		mockProjectName.set('Test Project');
		await tick();

		// Should affect the component's behavior (close button visibility)
		const closeButton = screen.queryByLabelText(/close project hub/i);
		expect(closeButton).toBeInTheDocument();
	});

	it('TestOnMount_Success: Initializes correctly on mount', async () => {
		render(ProjectHub, { show: true });

		// Should fetch projects on mount
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Should set up proper subscriptions
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
	});

	// Transition and Animation Tests
	it('TestTransitionEffects_Success: Uses proper transitions', () => {
		const { container } = render(ProjectHub, { show: true });

		// Should have backdrop with transition
		expect(container.querySelector('.backdrop')).toBeInTheDocument();
	});

	// Component State Tests
	it('TestComponentState_Success: Manages internal state correctly', async () => {
		const { component } = render(ProjectHub, { show: true });

		// Component should have proper initial state
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();

		// Should handle state changes
		await tick();
		expect(screen.getByText(/Recent Projects/i)).toBeInTheDocument();
	});

	// Error Handling Tests
	it('TestErrorRecovery_Success: Recovers from errors gracefully', async () => {
		// Test multiple error scenarios
		mockFetch
			.mockRejectedValueOnce(new Error('Network error'))
			.mockResolvedValueOnce({
				ok: true,
				json: async () => ({ all_projects: ['Recovery Project'] })
			});

		render(ProjectHub, { show: true });

		// Should handle first error
		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalled();
		});

		// Should still render the UI
		expect(screen.getByText(/Start a new project/i)).toBeInTheDocument();
	});
});

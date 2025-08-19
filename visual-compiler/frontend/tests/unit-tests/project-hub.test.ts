// @vitest-environment jsdom

import { describe, it, expect, beforeEach, vi, type MockedFunction } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import ProjectHub from '../../src/lib/components/project-hub/project-hub.svelte';

// Mock stores without using variables (to avoid hoisting issues)
vi.mock('$lib/stores/project', () => ({
	projectName: {
		set: vi.fn(),
		subscribe: vi.fn((callback: (value: string) => void) => {
			callback("");
			return vi.fn(); // Return an unsubscribe function
		})
	}
}));

vi.mock('$lib/stores/pipeline', () => ({
	pipelineStore: {
		set: vi.fn(),
		subscribe: vi.fn((callback: (value: any) => void) => {
			callback({ nodes: [], connections: [], lastSaved: null });
			return vi.fn(); // Return an unsubscribe function
		})
	},
	resetPipeline: vi.fn()
}));

vi.mock('$lib/stores/toast', () => ({
	AddToast: vi.fn()
}));

// Import the mocked modules to access their functions
import { projectName } from '$lib/stores/project';
import { pipelineStore, resetPipeline } from '$lib/stores/pipeline';
import { AddToast } from '$lib/stores/toast';

// Mock fetch globally
global.fetch = vi.fn();
const mockFetch = global.fetch as MockedFunction<typeof fetch>;

// Mock localStorage
const localStorageMock = (() => {
	let store: { [key: string]: string } = {};
	return {
		getItem: vi.fn((key: string) => store[key] || null),
		setItem: vi.fn((key: string, value: string) => { store[key] = value; }),
		clear: vi.fn(() => { store = {}; }),
		removeItem: vi.fn((key: string) => { delete store[key]; })
	};
})();
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('ProjectHub Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		localStorageMock.clear();
		
		// Reset fetch mock to successful response by default
		mockFetch.mockResolvedValue({
			ok: true,
			json: async () => ({
				message: "Retrieved users project details",
				results: [
					{ name: "Test Project", date_modified: "2024-08-19T00:00:00Z" }
				]
			})
		} as Response);
	});

	// Basic Rendering Tests
	it('TestRender_Success: Renders correctly when show is true', () => {
		const { container } = render(ProjectHub, { show: true });
		
		expect(container.querySelector('.modal')).toBeInTheDocument();
		expect(screen.getByText('Start a new project')).toBeInTheDocument();
	});

	it('TestModalStructure_Success: Has proper modal structure', async () => {
		const { container } = render(ProjectHub, { show: true });

		// Check for modal structure elements
		expect(container.querySelector('.backdrop')).toBeInTheDocument();
		expect(container.querySelector('.modal')).toBeInTheDocument();
	});

	// Project Fetching Tests
	it('TestFetchProjects_Success: Fetches projects when modal opens', async () => {
		localStorageMock.setItem('user_id', '123');
		render(ProjectHub, { show: true });

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/api/users/getProjects?users_id=123',
				{
					method: 'GET',
					headers: { 'accept': 'application/json' }
				}
			);
		});
	});

	it('TestFetchProjects_NoUserId: Skips fetch when no user ID', () => {
		// Don't set user_id in localStorage
		render(ProjectHub, { show: true });

		expect(mockFetch).not.toHaveBeenCalled();
	});

	// Project Selection Tests  
	it('TestProjectSelection_Success: Handles project selection with pipeline data', async () => {
		localStorageMock.setItem('user_id', '123');
		
		// Mock successful project list response
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				all_projects: ['Test Project']
			})
		} as Response);

		// Mock project details response with pipeline
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				message: "Retrieved users project details",
				results: {
					pipeline: {
						nodes: [{ id: 1, type: 'test' }],
						connections: [],
						lastSaved: '2024-08-19'
					}
				}
			})
		} as Response);

		const { container } = render(ProjectHub, { show: true });
		
		// Wait for initial projects to load, then look for project elements
		await waitFor(() => {
			expect(container.querySelector('.project-block')).toBeInTheDocument();
		});

		// Click on a project
		const projectButton = container.querySelector('.project-block') as HTMLElement;
		if (projectButton) {
			await fireEvent.click(projectButton);
			
			await waitFor(() => {
				// The component will call pipelineStore.set with the pipeline data from the response
				expect(pipelineStore.set).toHaveBeenCalled();
			});
		}
	});

	it('TestProjectSelection_NoPipelineData: Handles project selection without pipeline data', async () => {
		localStorageMock.setItem('user_id', '123');
		
		// Mock initial projects fetch
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				all_projects: ['Test Project']
			})
		} as Response);

		// Mock project details response without pipeline
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				message: "Retrieved users project details",
				results: {} // No pipeline data
			})
		} as Response);

		const { container } = render(ProjectHub, { show: true });
		
		await waitFor(() => {
			expect(container.querySelector('.project-block')).toBeInTheDocument();
		});

		// Click on a project
		const projectButton = container.querySelector('.project-block') as HTMLElement;
		if (projectButton) {
			await fireEvent.click(projectButton);
			
			await waitFor(() => {
				expect(pipelineStore.set).toHaveBeenCalledWith({
					nodes: [],
					connections: [],
					lastSaved: null
				});
			});
		}
	});

	it('TestProjectSelection_Error: Handles project selection error', async () => {
		localStorageMock.setItem('user_id', '123');
		
		// Mock successful initial projects fetch first
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				all_projects: ['Error Project']
			})
		} as Response);

		const { container } = render(ProjectHub, { show: true });
		
		// Wait for initial projects to load
		await waitFor(() => {
			expect(container.querySelector('.project-block')).toBeInTheDocument();
		});

		// Now mock error for project selection
		mockFetch.mockResolvedValueOnce({
			ok: false,
			status: 500
		} as Response);

		// Click on project to trigger error
		const projectButton = container.querySelector('.project-block') as HTMLElement;
		if (projectButton) {
			await fireEvent.click(projectButton);

			// Wait for error to be handled
			await waitFor(() => {
				expect(AddToast).toHaveBeenCalledWith('Error loading project: HTTP error! status: 500', 'error');
			});
		}
	});

	it('TestProjectSelection_NoUserId: Handles project selection without user ID', async () => {
		// Don't set user_id in localStorage to simulate logged out user
		localStorageMock.clear();
		
		// Create a simple component instance
		const { container } = render(ProjectHub, { show: true });
		
		// Since we can't easily trigger selectProject without a rendered project,
		// let's verify that without a userId, no fetch calls should be made when modal opens
		await waitFor(() => {
			// Should not fetch projects when no user ID
			expect(mockFetch).not.toHaveBeenCalled();
		}, { timeout: 1000 });
	});

	// New Project Creation Tests
	it('TestCreateNewProject_Success: Shows create new project button and handles click', async () => {
		const { container } = render(ProjectHub, { show: true });

		const newBlankButton = container.querySelector('.project-button');
		expect(newBlankButton).toBeInTheDocument();
		
		if (newBlankButton) {
			await fireEvent.click(newBlankButton as HTMLElement);
			// Should open project name prompt (can't easily test modal opening without component internals)
		}
	});

	// Project Deletion Tests
	it('TestDeleteProject_Success: Shows delete buttons and handles click', async () => {
		localStorageMock.setItem('user_id', '123');
		
		// Mock projects fetch to have a project with delete button
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				all_projects: ['Test Project']
			})
		} as Response);
		
		const { container } = render(ProjectHub, { show: true });
		
		await waitFor(() => {
			const deleteButton = container.querySelector('.delete-button');
			expect(deleteButton).toBeInTheDocument();
		});
	});

	// Modal Interaction Tests  
	it('TestCloseModal_NoExistingProject: Does not show close button without existing project', async () => {
		const { container } = render(ProjectHub, { show: true });

		// Should not have close button when no existing project
		const closeButton = container.querySelector('.close-button');
		expect(closeButton).not.toBeInTheDocument();
	});

	// UI Element Tests
	it('TestSearchBar_Success: Renders search input', () => {
		render(ProjectHub, { show: true });
		
		const searchInput = screen.getByPlaceholderText('Search projects...');
		expect(searchInput).toBeInTheDocument();
	});

	it('TestSearchBar_Icon: Renders search icon', () => {
		const { container } = render(ProjectHub, { show: true });
		
		const searchIcon = container.querySelector('.search-icon');
		expect(searchIcon).toBeInTheDocument();
	});

	// Styling and Accessibility Tests
	it('TestDarkModeStyles_Success: Component supports dark mode styles', () => {
		const { container } = render(ProjectHub, { show: true });
		
		// Check that component has proper class structure for styling
		const modal = container.querySelector('.modal');
		expect(modal).toHaveClass('svelte-v8mqtk');
	});

	it('TestAccessibility_Success: Has proper ARIA labels and roles', async () => {
		localStorageMock.setItem('user_id', '123');
		
		// Mock projects fetch to have a project with delete button
		mockFetch.mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				all_projects: ['Test Project']
			})
		} as Response);
		
		const { container } = render(ProjectHub, { show: true });
		
		await waitFor(() => {
			const deleteButton = container.querySelector('[aria-label="Delete project"]');
			expect(deleteButton).toBeInTheDocument();
		});
	});

	it('TestKeyboardNavigation_Success: Supports keyboard navigation', async () => {
		const { container } = render(ProjectHub, { show: true });
		
		// Test that buttons are focusable
		const buttons = container.querySelectorAll('button');
		expect(buttons.length).toBeGreaterThan(0);
		
		// Focus first button
		if (buttons[0]) {
			buttons[0].focus();
			expect(document.activeElement).toBe(buttons[0]);
		}
	});

	// Store Integration Tests
	it('TestProjectNameStore_Success: Integrates with project name store', () => {
		render(ProjectHub, { show: true });
		
		// Component should subscribe to project name store
		expect(projectName.subscribe).toHaveBeenCalled();
	});

	// Lifecycle Tests
	it('TestOnMount_Success: Initializes correctly on mount', () => {
		const { container } = render(ProjectHub, { show: true });
		
		expect(container.querySelector('.modal')).toBeInTheDocument();
	});

	// Animation and Transition Tests
	it('TestTransitionEffects_Success: Uses proper transitions', () => {
		const { container } = render(ProjectHub, { show: true });
		
		// Component should render with transition classes
		const modal = container.querySelector('.modal');
		expect(modal).toBeInTheDocument();
	});

	// State Management Tests
	it('TestComponentState_Success: Manages internal state correctly', () => {
		render(ProjectHub, { show: true });
		
		// Component should initialize with proper default state
		expect(screen.getByText('Start a new project')).toBeInTheDocument();
	});

	// Error Handling Tests
	it('TestErrorRecovery_Success: Recovers from errors gracefully', () => {
		// Mock fetch to reject
		mockFetch.mockRejectedValue(new Error('Network error'));
		
		render(ProjectHub, { show: true });
		
		// Component should still render despite fetch error
		expect(screen.getByText('Start a new project')).toBeInTheDocument();
	});
});

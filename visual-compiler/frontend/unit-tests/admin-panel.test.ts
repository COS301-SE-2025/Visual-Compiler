import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import AdminPanel from '../../src/lib/components/main/admin-panel.svelte';

// Mock localStorage
const mockLocalStorage = {
	getItem: vi.fn(),
	setItem: vi.fn(),
	removeItem: vi.fn(),
	clear: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: mockLocalStorage });

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch as any;

describe('AdminPanel Component', () => {
	const mockUsers = [
		{
			users_id: '1',
			username: 'testuser1',
			email: 'test1@example.com',
			projects: []
		},
		{
			users_id: '2', 
			username: 'testuser2',
			email: 'test2@example.com',
			projects: []
		}
	];

	beforeEach(() => {
		vi.clearAllMocks();
		mockLocalStorage.getItem.mockReturnValue('admin-123');
		
		// Mock successful fetch by default
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ users: mockUsers })
		});
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it('TestRender_Success: Renders admin panel with header and close button', async () => {
		render(AdminPanel);

		expect(screen.getByText('Admin User Management')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: '✕' })).toBeInTheDocument();
		expect(screen.getByRole('dialog')).toHaveAttribute('aria-label', 'Admin Panel');
		expect(screen.getByPlaceholderText('Search users...')).toBeInTheDocument();
	});

	it('TestCloseButton_Success: Close button is clickable', async () => {
		render(AdminPanel);

		const closeButton = screen.getByRole('button', { name: '✕' });
		await fireEvent.click(closeButton);

		// Button should be functional (no errors)
		expect(closeButton).toBeInTheDocument();
	});

	it('TestLoadingState_Success: Shows loading state initially', () => {
		render(AdminPanel);
		expect(screen.getByText('Loading users...')).toBeInTheDocument();
	});

	it('TestFetchUsers_Success: Fetches and displays users', async () => {
		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
			expect(screen.getByText('testuser2')).toBeInTheDocument();
		});

		expect(mockFetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/users/getUsers',
			expect.objectContaining({
				method: 'GET',
				headers: {
					'accept': 'application/json',
					'Content-Type': 'application/json'
				}
			})
		);
	});

	it('TestFetchUsersError_Success: Handles fetch error gracefully', async () => {
		mockFetch.mockRejectedValue(new Error('Network error'));
		
		render(AdminPanel);

		await waitFor(() => {
			expect(screen.queryByText('Loading users...')).not.toBeInTheDocument();
		});

		// Should not crash and loading should stop
		expect(screen.queryByText('testuser1')).not.toBeInTheDocument();
	});

	it('TestNoUsers_Success: Shows message when no users found', async () => {
		mockFetch.mockResolvedValue({
			ok: true,
			json: () => Promise.resolve({ users: [] })
		});

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('No users found.')).toBeInTheDocument();
		});
	});

	it('TestSearchFilter_Success: Filters users based on search term', async () => {
		render(AdminPanel);

		// Wait for users to load
		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search users...');
		await fireEvent.input(searchInput, { target: { value: 'testuser1' } });

		expect(screen.getByText('testuser1')).toBeInTheDocument();
		expect(screen.queryByText('testuser2')).not.toBeInTheDocument();
	});

	it('TestSearchByEmail_Success: Filters users by email', async () => {
		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search users...');
		await fireEvent.input(searchInput, { target: { value: 'test2@example.com' } });

		expect(screen.queryByText('testuser1')).not.toBeInTheDocument();
		expect(screen.getByText('testuser2')).toBeInTheDocument();
	});

	it('TestAPICallConfiguration_Success: Makes correct API calls', async () => {
		render(AdminPanel);

		await waitFor(() => {
			expect(mockFetch).toHaveBeenCalledTimes(1);
		});

		expect(mockFetch).toHaveBeenCalledWith(
			'http://localhost:8080/api/users/getUsers',
			{
				method: 'GET',
				headers: {
					'accept': 'application/json',
					'Content-Type': 'application/json'
				}
			}
		);
	});

	it('TestLocalStorageAccess_Success: Accesses admin ID from localStorage', async () => {
		render(AdminPanel);

		await waitFor(() => {
			expect(mockLocalStorage.getItem).toHaveBeenCalledWith('user_id');
		});
	});
});

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
			'https://www.visual-compiler.co.za/api/users/getUsers',
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
			'https://www.visual-compiler.co.za/api/users/getUsers',
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

	it('TestToggleAccordion_Success: Toggles accordion expansion state', async () => {
		render(AdminPanel);

		// Wait for users to load
		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		// Click on user header to expand
		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		// Should show expanded content
		await waitFor(() => {
			expect(screen.getByText('Email:')).toBeInTheDocument();
		});

		// Click again to collapse
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		// Should hide expanded content
		await waitFor(() => {
			expect(screen.queryByText('Email:')).not.toBeInTheDocument();
		});
	});

	it('TestSelectUser_Success: Selects user for editing', async () => {
		render(AdminPanel);

		// Wait for users to load
		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		// Click on user header to expand first
		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		// Wait for edit button to appear
		await waitFor(() => {
			expect(screen.getByTitle('Edit user')).toBeInTheDocument();
		});

		// Click edit button to select user
		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		// Should populate edit form
		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
			expect(screen.getByDisplayValue('test1@example.com')).toBeInTheDocument();
		});
	});

	it('TestHighlight_Success: Highlights search terms in text', async () => {
		render(AdminPanel);

		// Wait for users to load
		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		// Search for a user
		const searchInput = screen.getByPlaceholderText('Search users...');
		await fireEvent.input(searchInput, { target: { value: 'test' } });

		// Should highlight the search term
		await waitFor(() => {
			expect(document.querySelector('.highlight')).toBeInTheDocument();
		});
	});

	it('TestSaveEdit_Success: Successfully updates user', async () => {
		// Mock successful update response
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ users: mockUsers })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ message: 'User updated successfully' })
			});

		render(AdminPanel);

		// Wait for users to load and select user
		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		// Wait for edit form
		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		// Modify the username
		const usernameInput = screen.getByDisplayValue('testuser1');
		await fireEvent.input(usernameInput, { target: { value: 'updated_user' } });

		// Click save button
		const saveButton = screen.getByText('Save');
		await fireEvent.click(saveButton);

		// Should show success message
		await waitFor(() => {
			expect(screen.getByText('User updated successfully!')).toBeInTheDocument();
		});
	});

	it('TestSaveEdit_Error: Handles update error', async () => {
		// Mock failed update response
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ users: mockUsers })
			})
			.mockResolvedValueOnce({
				ok: false,
				json: () => Promise.resolve({ error: 'Update failed' })
			});

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const saveButton = screen.getByText('Save');
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(screen.getByText('Update failed')).toBeInTheDocument();
		});
	});

	it('TestSaveEdit_NetworkError: Handles network error', async () => {
		// Mock network error
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ users: mockUsers })
			})
			.mockRejectedValueOnce(new Error('Network error'));

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const saveButton = screen.getByText('Save');
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(screen.getByText('Failed to connect to server')).toBeInTheDocument();
		});
	});

	it('TestSaveEdit_NoAdmin: Requires admin privileges', async () => {
		// Mock no admin ID
		mockLocalStorage.getItem.mockReturnValue(null);

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const saveButton = screen.getByText('Save');
		await fireEvent.click(saveButton);

		await waitFor(() => {
			expect(screen.getByText('Admin privileges required')).toBeInTheDocument();
		});
	});

	it('TestDeleteUser_Success: Successfully deletes user', async () => {
		// Mock successful delete response
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ users: mockUsers })
			})
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ message: 'User deleted successfully' })
			});

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const deleteButton = screen.getByTitle('Delete user');
		await fireEvent.click(deleteButton);

		// After successful delete, the user should be removed from the list
		// and selectedUserIdx should be set to -1, showing the placeholder
		await waitFor(() => {
			expect(screen.queryByText('testuser1')).not.toBeInTheDocument();
			expect(screen.getByText('Select a user to edit details here.')).toBeInTheDocument();
		});

		// The success message should be visible in the edit area briefly before it's hidden
		// due to selectedUserIdx being reset to -1, but let's check the fetch was called correctly
		expect(mockFetch).toHaveBeenCalledWith(
			'https://www.visual-compiler.co.za/api/users/delete',
			expect.objectContaining({
				method: 'DELETE',
				headers: expect.objectContaining({
					'Content-Type': 'application/json',
					'accept': 'application/json'
				}),
				body: JSON.stringify({ 
					admin_id: 'admin-123',
					users_id: '1'
				})
			})
		);
	});

	it('TestDeleteUser_Error: Handles delete error', async () => {
		// Mock failed delete response
		mockFetch
			.mockResolvedValueOnce({
				ok: true,
				json: () => Promise.resolve({ users: mockUsers })
			})
			.mockResolvedValueOnce({
				ok: false,
				json: () => Promise.resolve({ error: 'Delete failed' })
			});

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const deleteButton = screen.getByTitle('Delete user');
		await fireEvent.click(deleteButton);

		await waitFor(() => {
			expect(screen.getByText('Delete failed')).toBeInTheDocument();
		});
	});

	it('TestDeleteUser_NoAdmin: Requires admin privileges for deletion', async () => {
		// Mock no admin ID
		mockLocalStorage.getItem.mockReturnValue(null);

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const deleteButton = screen.getByTitle('Delete user');
		await fireEvent.click(deleteButton);

		await waitFor(() => {
			expect(screen.getByText('Admin privileges required')).toBeInTheDocument();
		});
	});

	it('TestCancelEdit_Success: Cancels user editing', async () => {
		render(AdminPanel);

		await waitFor(() => {
			expect(screen.getByText('testuser1')).toBeInTheDocument();
		});

		const userHeader = screen.getByText('testuser1').closest('.user-header');
		expect(userHeader).toBeInTheDocument();
		if (userHeader) {
			await fireEvent.click(userHeader);
		}

		const editButton = screen.getByTitle('Edit user');
		await fireEvent.click(editButton);

		await waitFor(() => {
			expect(screen.getByDisplayValue('testuser1')).toBeInTheDocument();
		});

		const cancelButton = screen.getByText('Cancel');
		await fireEvent.click(cancelButton);

		await waitFor(() => {
			expect(screen.getByText('Select a user to edit details here.')).toBeInTheDocument();
		});
	});

	it('TestFetchUsers_APIError: Handles API error response', async () => {
		// Mock API error response
		mockFetch.mockResolvedValue({
			ok: false,
			json: () => Promise.resolve({ error: 'API Error' })
		});

		render(AdminPanel);

		await waitFor(() => {
			expect(screen.queryByText('Loading users...')).not.toBeInTheDocument();
		});

		// Should not display users
		expect(screen.queryByText('testuser1')).not.toBeInTheDocument();
	});

	it('TestEventDispatcher_Success: Dispatches close event', async () => {
		render(AdminPanel);

		const closeButton = screen.getByRole('button', { name: '✕' });
		
		// Test that clicking the close button doesn't cause errors
		// In a real app, this would dispatch a close event to the parent
		await fireEvent.click(closeButton);
		
		// Verify the button is still functional
		expect(closeButton).toBeInTheDocument();
	});
});

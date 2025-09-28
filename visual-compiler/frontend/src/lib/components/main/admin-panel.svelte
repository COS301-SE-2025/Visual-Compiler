<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';

    const dispatch = createEventDispatcher();

    let users = [];
    let searchTerm = '';
    let expandedIndex = -1;
    let loading = true;
    let selectedUserIdx: number = -1;
    let editUsername = '';
    let editEmail = '';
    let originalUsername = ''; // Store original values for comparison
    let originalEmail = '';
    let editError = '';
    let editSuccess = '';

    let adminId = sessionStorage.getItem('user_id');

    interface User {
        users_id: string;
        username: string;
        email: string;
        projects: any[];
    }

    async function fetchUsers() {
        loading = true;
        try {
            const res = await fetch('http://localhost:8080/api/users/getUsers', {
                method: 'GET',
                headers: {
                    'accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            });

            if (!res.ok) {
                const errorData = await res.json();
                throw new Error(errorData.error || 'Failed to fetch users');
            }

            const data = await res.json();
            users = (data.users || []).map((u: User) => ({
                ...u,
                ID: u.users_id, // Map the MongoDB ObjectID string to ID
                username: u.username,
                email: u.email,
                projects: u.projects || []
            }));
            loading = false;
        } catch (error) {
            console.error('Error fetching users:', error);
            loading = false;
        }
    }

    onMount(fetchUsers);

    function toggleAccordion(index: number) {
        expandedIndex = expandedIndex === index ? -1 : index;
    }

    function selectUser(idx: number) {
        selectedUserIdx = idx;
        editUsername = users[idx].username;
        editEmail = users[idx].email;
        originalUsername = users[idx].username; // Store original values
        originalEmail = users[idx].email;
        editError = '';
        editSuccess = '';
        console.log('Selected user ID:', users[idx].ID); 
    }

    $: filteredUsers = users.filter(u =>
        u.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
        u.email.toLowerCase().includes(searchTerm.toLowerCase())
    );

    function highlight(text: string) {
        if (!searchTerm) return text;
        const re = new RegExp(`(${searchTerm.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
        return text.replace(re, '<span class="highlight">$1</span>');
    }

    async function saveEdit() {
        editError = '';
        editSuccess = '';
        
        if (!adminId) {
            editError = 'Admin privileges required';
            return;
        }

        // Check if any fields have actually changed
        const usernameChanged = editUsername !== originalUsername;
        const emailChanged = editEmail !== originalEmail;
        
        if (!usernameChanged && !emailChanged) {
            editError = 'No changes were made';
            return;
        }

        const user = users[selectedUserIdx];
        
        // Only include fields that have actually changed
        const updateData: any = {
            admin_id: adminId,
            users_id: user.ID
        };
        
        if (usernameChanged) {
            updateData.username = editUsername;
        }
        
        if (emailChanged) {
            updateData.email = editEmail;
        }

        try {
            const res = await fetch('http://localhost:8080/api/users/update', {
                method: 'PATCH',
                headers: { 
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                },
                body: JSON.stringify(updateData)
            });

            const data = await res.json();
            if (res.ok) {
                editSuccess = 'User updated successfully!';
                // Update the local state with new values
                if (usernameChanged) {
                    users[selectedUserIdx].username = editUsername;
                    originalUsername = editUsername; // Update original value
                }
                if (emailChanged) {
                    users[selectedUserIdx].email = editEmail;
                    originalEmail = editEmail; // Update original value
                }
            } else {
                editError = data.error || 'Failed to update user.';
            }
        } catch (error) {
            editError = 'Failed to connect to server';
            console.error('Error updating user:', error);
        }
    }

    async function deleteUser() {
        editError = '';
        editSuccess = '';

        if (!adminId) {
            editError = 'Admin privileges required';
            return;
        }

        const user = users[selectedUserIdx];
        try {
            const res = await fetch('http://localhost:8080/api/users/delete', {
                method: 'DELETE',
                headers: { 
                    'Content-Type': 'application/json',
                    'accept': 'application/json'
                },
                body: JSON.stringify({ 
                    admin_id: adminId,
                    users_id: user.ID 
                })
            });

            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.error || 'Failed to delete user');
            }

            editSuccess = 'User deleted successfully!';
            users = users.filter((_, idx) => idx !== selectedUserIdx);
            selectedUserIdx = -1;
        } catch (error) {
            editError = error.message || 'Failed to delete user';
            console.error('Error deleting user:', error);
        }
    }
</script>

<div class="admin-panel-bg">
    <div class="admin-panel-window" role="dialog" aria-label="Admin Panel">
        <div class="admin-panel-header">
            <h2>Admin User Management</h2>
            <button class="close-btn" on:click={() => dispatch('close')}>✕</button>
        </div>
        <div class="admin-panel-content">
            <!-- Left: User List -->
            <div class="user-list-col">
                <div class="search-bar">
                    <input
                        type="text"
                        placeholder="Search users..."
                        bind:value={searchTerm}
                    />
                </div>
                {#if loading}
                    <p class="loading">Loading users...</p>
                {:else if filteredUsers.length === 0}
                    <p class="no-users">No users found.</p>
                {:else}
                    <div class="user-accordion">
                        {#each filteredUsers as user, idx}
                            <div class="user-card {selectedUserIdx === users.indexOf(user) ? 'selected' : ''}">
                                <div class="user-header" on:click={() => toggleAccordion(idx)}>
                                    <span>{@html highlight(user.username)}</span>
                                    <span class="arrow">{expandedIndex === idx ? '▲' : '▼'}</span>
                                </div>
                                {#if expandedIndex === idx}
                                    <div class="user-details">
                                        <div><strong>Email:</strong> {@html highlight(user.email)}</div>
                                        <button class="edit-btn" title="Edit user" on:click={() => selectUser(users.indexOf(user))}>
                                            <svg width="18" height="18" fill="none" stroke="#1976d2" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
                                                <path d="M12 20h9" />
                                                <path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19.5 3 21l1.5-4L16.5 3.5z" />
                                            </svg>
                                        </button>
                                    </div>
                                {/if}
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>
            <!-- Right: Editing column -->
            <div class="edit-col">
                {#if selectedUserIdx !== -1}
                    <h3>Edit User</h3>
                    <div class="edit-form">
                        <label>
                            Username
                            <input type="text" bind:value={editUsername} />
                        </label>
                        <label>
                            Email
                            <input type="email" bind:value={editEmail} />
                        </label>
                        <div class="edit-actions">
                            <button class="save-btn" on:click={saveEdit}>Save</button>
                            <button class="cancel-btn" on:click={() => selectedUserIdx = -1}>Cancel</button>
                            <button class="delete-btn" on:click={deleteUser} title="Delete user">
                                <svg width="18" height="18" fill="none" stroke="#d32f2f" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24">
                                    <polyline points="3 6 5 6 21 6"></polyline>
                                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"></path>
                                    <line x1="10" y1="11" x2="10" y2="17"></line>
                                    <line x1="14" y1="11" x2="14" y2="17"></line>
                                </svg>
                            </button>
                        </div>
                        {#if editError}
                            <p class="edit-error">{editError}</p>
                        {/if}
                        {#if editSuccess}
                            <p class="edit-success">{editSuccess}</p>
                        {/if}
                    </div>
                {:else}
                    <p class="edit-placeholder">Select a user to edit details here.</p>
                {/if}
            </div>
        </div>
    </div>
</div>

<style>
    .admin-panel-bg {
        position: fixed;
        top: 0; left: 0; right: 0; bottom: 0;
        background: rgba(0,0,0,0.25);
        z-index: 5000;
        display: flex;
        justify-content: flex-end;
        align-items: flex-start;
    }
    .admin-panel-window {
        margin: 2rem;
        width: 900px;
        max-width: 90vw;
        height: 80vh;
        background: #f9fafc;
        border-radius: 16px;
        box-shadow: 0 8px 32px rgba(0,0,0,0.18);
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }
    .admin-panel-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 2rem 1rem 2rem;
        border-bottom: 1px solid #e5e5e5;
        background: #f1f5fa;
    }
    .admin-panel-header h2 {
        margin: 0;
        font-size: 1.4rem;
        font-weight: 700;
        letter-spacing: 0.02em;
        color: #041a47;
    }
    .close-btn {
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        color: #888;
        transition: color 0.2s;
    }
    .close-btn:hover {
        color: #041a47;
    }
    .admin-panel-content {
        display: flex;
        flex: 1;
        height: 100%;
    }
    .user-list-col {
        width: 45%;
        border-right: 1px solid #e5e5e5;
        padding: 1.5rem 1rem 1.5rem 2rem;
        overflow-y: auto;
        background: #f7fafd;
    }
    .search-bar {
        margin-bottom: 1.2rem;
    }
    .search-bar input {
        width: 100%;
        padding: 0.6rem 1rem;
        border-radius: 8px;
        border: 1.5px solid #b0c4de;
        font-size: 1rem;
        background: #fff;
        transition: border 0.2s;
    }
    .search-bar input:focus {
        border: 1.5px solid #041a47;
        outline: none;
    }
    .user-accordion {
        display: flex;
        flex-direction: column;
        gap: 0.7rem;
    }
    .user-card {
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 1px 4px rgba(25, 118, 210, 0.07);
        padding: 0.7rem 1.1rem;
        transition: box-shadow 0.2s, border 0.2s;
        border: 2px solid transparent;
    }
    .user-card.selected {
        border: 2px solid #1976d2;
        box-shadow: 0 2px 8px rgba(25, 118, 210, 0.13);
    }
    .user-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        cursor: pointer;
        font-weight: 600;
        font-size: 1.07rem;
        color: #041a47;
    }
    .user-details {
        margin-top: 0.5rem;
        font-size: 0.97rem;
        color: #444;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    .edit-btn {
        background: #e3eafc;
        border: none;
        border-radius: 50%;
        padding: 0.35rem;
        cursor: pointer;
        transition: background 0.2s;
        margin-left: 1rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .edit-btn:hover {
        background: #1976d2;
    }
    .edit-btn svg {
        stroke: #1976d2;
        transition: stroke 0.2s;
    }
    .edit-btn:hover svg {
        stroke: #fff;
    }
    .arrow {
        font-size: 0.9rem;
    }
    .edit-col {
        width: 55%;
        padding: 2.5rem 2rem;
        background: #fff;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
    }
    .edit-form {
        display: flex;
        flex-direction: column;
        gap: 1.2rem;
        max-width: 350px;
    }
    .edit-form label {
        font-weight: 500;
        color: #041a47;
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
    }
    .edit-form input {
        padding: 0.5rem 0.8rem;
        border-radius: 6px;
        border: 1.5px solid #b0c4de;
        font-size: 1rem;
        background: #f7fafd;
        transition: border 0.2s;
    }
    .edit-form input:focus {
        border: 1.5px solid #041a47;
        outline: none;
    }
    .edit-actions {
        display: flex;
        gap: 1rem;
    }
    .save-btn {
        background: #041a47;
        color: #fff;
        border: none;
        border-radius: 5px;
        padding: 0.5rem 1.3rem;
        font-size: 1rem;
        cursor: pointer;
        transition: background 0.2s;
    }
    .save-btn:hover {
        background: #041a47;
    }
    .cancel-btn {
        background: #e0e0e0;
        color: #444;
        border: none;
        border-radius: 5px;
        padding: 0.5rem 1.3rem;
        font-size: 1rem;
        cursor: pointer;
        transition: background 0.2s;
    }
    .cancel-btn:hover {
        background: #bdbdbd;
    }
    .delete-btn {
        background: #ffeaea;
        border: none;
        border-radius: 50%;
        padding: 0.35rem;
        cursor: pointer;
        transition: background 0.2s;
        margin-left: 1rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .delete-btn:hover {
        background: #d32f2f;
    }
    .delete-btn svg {
        stroke: #d32f2f;
        transition: stroke 0.2s;
    }
    .delete-btn:hover svg {
        stroke: #fff;
    }
    .edit-error {
        color: #d32f2f;
        font-size: 0.97rem;
        margin-top: 0.5rem;
    }
    .edit-success {
        color: #388e3c;
        font-size: 0.97rem;
        margin-top: 0.5rem;
    }
    .edit-placeholder {
        color: #888;
        font-size: 1.1rem;
        margin-top: 2rem;
    }
    .loading, .no-users {
        color: #888;
        font-size: 1.05rem;
        margin-top: 2rem;
    }
    .highlight {
        background: #ffe082;
        color: #222;
        border-radius: 3px;
        padding: 0 2px;
    }
    :global(html.dark-mode) .admin-panel-window {
        background: #1a2233;
    }
    :global(html.dark-mode) .admin-panel-header {
        background: #232c43;
    }
    :global(html.dark-mode) .user-list-col {
        background: #232c43;
    }
    :global(html.dark-mode) .edit-col {
        background: #1a2233;
    }
    :global(html.dark-mode) .user-card {
        background: #232c43;
        color: #90caf9;
    }
    :global(html.dark-mode) .user-header {
        color: #90caf9;
    }
    :global(html.dark-mode) .edit-form label {
        color: #90caf9;
    }
    :global(html.dark-mode) .edit-form input {
        background: #1a2233;
        color: #fff;
    }
    :global(html.dark-mode) .save-btn {
        background: #90caf9;
        color: #1a2233;
    }
    :global(html.dark-mode) .edit-btn {
        background: #90caf9;
        color: #1a2233;
    }
    :global(html.dark-mode) .edit-btn:hover,
    :global(html.dark-mode) .save-btn:hover {
        background: #1976d2;
        color: #fff;
    }

    
    ::-webkit-scrollbar {
        width: 11px;
    }
    ::-webkit-scrollbar-track {
        background: #f1f1f1;
    }
    ::-webkit-scrollbar-thumb {
        background-color: #888;
        border-radius: 10px;
    }
    ::-webkit-scrollbar-thumb:hover {
        background: #555;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-track {
        background: #2d3748;
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb {
        background-color: #4a5568; 
        border-color: #2d3748; 
    }

    :global(html.dark-mode) ::-webkit-scrollbar-thumb:hover {
        background: #616e80;
    }

</style>

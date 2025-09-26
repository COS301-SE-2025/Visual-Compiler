import { writable } from 'svelte/store';

export const projectName = writable<string>('');

export const saveProject = async (name: string, userId: string) => {
    try {
        const response = await fetch('https://www.visual-compiler.co.za/api/users/save', {
            method: 'POST',
            headers: {
                'accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                project_name: name,
                users_id: userId
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        if (data.message) {
            projectName.set(name);
            return true;
        }
        return false;
    } catch (error) {
        console.error('Error saving project:', error);
        return false;
    }
};

export const getProjects = async (userId: string) => {
    try {
        const response = await fetch(`https://www.visual-compiler.co.za/api/users/getProjects?users_id=${userId}`, {
            method: 'GET',
            headers: {
                'accept': 'application/json'
            }
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        if (data.all_projects) {
            return data.all_projects;
        }
        return [];
    } catch (error) {
        console.error('Error fetching projects:', error);
        return [];
    }
};

export const deleteProject = async (name: string, userId: string) => {
    try {
        const response = await fetch('http://localhost:8080/api/users/deleteProject', {
            method: 'DELETE',
            headers: {
                'accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                project_name: name,
                users_id: userId
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data.message || 'Project deleted successfully';
    } catch (error) {
        console.error('Error deleting project:', error);
        throw error;
    }
};

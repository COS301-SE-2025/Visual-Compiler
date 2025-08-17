import { writable } from 'svelte/store';

export const projectName = writable<string>('');

export const saveProject = async (name: string, userId: string) => {
    try {
        const response = await fetch('http://localhost:8080/api/users/save', {
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
        const response = await fetch(`http://localhost:8080/api/users/getProjects?users_id=${userId}`, {
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

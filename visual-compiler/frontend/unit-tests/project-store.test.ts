import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { get } from 'svelte/store';
import { projectName, saveProject, getProjects } from '../../src/lib/stores/project';

// Mock fetch globally
const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

describe('Project Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset the store
        projectName.set('');
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('projectName store', () => {
        it('should have initial empty string value', () => {
            const currentValue = get(projectName);
            expect(currentValue).toBe('');
        });

        it('should update projectName value', () => {
            projectName.set('My Test Project');
            const currentValue = get(projectName);
            expect(currentValue).toBe('My Test Project');
        });

        it('should support subscription to changes', () => {
            const values: string[] = [];
            
            const unsubscribe = projectName.subscribe((value) => {
                values.push(value);
            });

            projectName.set('Project 1');
            projectName.set('Project 2');

            expect(values).toEqual(['', 'Project 1', 'Project 2']);
            
            unsubscribe();
        });
    });

    describe('saveProject function', () => {
        it('should successfully save project and update store', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ message: 'Project saved successfully' })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await saveProject('Test Project', 'user123');

            expect(result).toBe(true);
            expect(get(projectName)).toBe('Test Project');
            expect(mockFetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/users/save',
                {
                    method: 'POST',
                    headers: {
                        'accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        project_name: 'Test Project',
                        users_id: 'user123'
                    })
                }
            );
        });

        it('should handle HTTP error responses', async () => {
            const mockResponse = {
                ok: false,
                status: 400,
                json: vi.fn().mockResolvedValue({ error: 'Bad request' })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await saveProject('Test Project', 'user123');

            expect(result).toBe(false);
            expect(get(projectName)).toBe(''); // Should not update store on failure
            expect(consoleSpy).toHaveBeenCalledWith(
                'Error saving project:', 
                expect.any(Error)
            );

            consoleSpy.mockRestore();
        });

        it('should handle network errors', async () => {
            mockFetch.mockRejectedValue(new Error('Network error'));

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await saveProject('Test Project', 'user123');

            expect(result).toBe(false);
            expect(get(projectName)).toBe(''); // Should not update store on failure
            expect(consoleSpy).toHaveBeenCalledWith(
                'Error saving project:',
                expect.any(Error)
            );

            consoleSpy.mockRestore();
        });

        it('should handle response without message property', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ success: true }) // No message property
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await saveProject('Test Project', 'user123');

            expect(result).toBe(false);
            expect(get(projectName)).toBe(''); // Should not update store when no message
        });

        it('should handle empty project name', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ message: 'Project saved' })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await saveProject('', 'user123');

            expect(result).toBe(true);
            expect(get(projectName)).toBe('');
            expect(mockFetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/users/save',
                expect.objectContaining({
                    body: JSON.stringify({
                        project_name: '',
                        users_id: 'user123'
                    })
                })
            );
        });
    });

    describe('getProjects function', () => {
        it('should successfully fetch projects list', async () => {
            const mockProjects = [
                { id: '1', name: 'Project 1', created_at: '2023-01-01' },
                { id: '2', name: 'Project 2', created_at: '2023-01-02' }
            ];
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ all_projects: mockProjects })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await getProjects('user123');

            expect(result).toEqual(mockProjects);
            expect(mockFetch).toHaveBeenCalledWith(
                'http://localhost:8080/api/users/getProjects?users_id=user123',
                {
                    method: 'GET',
                    headers: {
                        'accept': 'application/json'
                    }
                }
            );
        });

        it('should handle HTTP error responses', async () => {
            const mockResponse = {
                ok: false,
                status: 404,
                json: vi.fn().mockResolvedValue({ error: 'User not found' })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await getProjects('user123');

            expect(result).toEqual([]);
            expect(consoleSpy).toHaveBeenCalledWith(
                'Error fetching projects:',
                expect.any(Error)
            );

            consoleSpy.mockRestore();
        });

        it('should handle network errors', async () => {
            mockFetch.mockRejectedValue(new Error('Network error'));

            const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

            const result = await getProjects('user123');

            expect(result).toEqual([]);
            expect(consoleSpy).toHaveBeenCalledWith(
                'Error fetching projects:',
                expect.any(Error)
            );

            consoleSpy.mockRestore();
        });

        it('should handle response without all_projects property', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ projects: [] }) // Wrong property name
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await getProjects('user123');

            expect(result).toEqual([]);
        });

        it('should handle empty projects list', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ all_projects: [] })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await getProjects('user123');

            expect(result).toEqual([]);
        });

        it('should handle null all_projects', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ all_projects: null })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await getProjects('user123');

            expect(result).toEqual([]);
        });

        it('should handle undefined all_projects', async () => {
            const mockResponse = {
                ok: true,
                status: 200,
                json: vi.fn().mockResolvedValue({ all_projects: undefined })
            };
            mockFetch.mockResolvedValue(mockResponse);

            const result = await getProjects('user123');

            expect(result).toEqual([]);
        });
    });
});

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { get } from 'svelte/store';

describe('Theme Store', () => {
    let mockLocalStorage: Record<string, string>;
    let mockDocument: {
        documentElement: {
            classList: {
                add: ReturnType<typeof vi.fn>;
                toggle: ReturnType<typeof vi.fn>;
            };
        };
    };

    beforeEach(() => {
        // Reset mocks
        mockLocalStorage = {};
        mockDocument = {
            documentElement: {
                classList: {
                    add: vi.fn(),
                    toggle: vi.fn()
                }
            }
        };

        // Mock localStorage
        Object.defineProperty(global, 'localStorage', {
            value: {
                getItem: vi.fn((key: string) => mockLocalStorage[key] || null),
                setItem: vi.fn((key: string, value: string) => {
                    mockLocalStorage[key] = value;
                }),
                removeItem: vi.fn((key: string) => {
                    delete mockLocalStorage[key];
                }),
                clear: vi.fn(() => {
                    mockLocalStorage = {};
                })
            },
            writable: true,
            configurable: true
        });

        // Mock document
        Object.defineProperty(global, 'document', {
            value: mockDocument,
            writable: true,
            configurable: true
        });

        vi.clearAllMocks();
        vi.resetModules();
    });

    afterEach(() => {
        vi.restoreAllMocks();
        vi.resetModules();
    });

    describe('theme store initialization', () => {
        it('should initialize with light theme when no stored theme', async () => {
            const { theme } = await import('../../src/lib/stores/theme');
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });

        it('should initialize with stored dark theme', async () => {
            mockLocalStorage['vc-theme'] = 'dark';
            
            const { theme } = await import('../../src/lib/stores/theme');
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('dark');
        });

        it('should initialize with light theme when stored theme is not dark', async () => {
            mockLocalStorage['vc-theme'] = 'light';
            
            const { theme } = await import('../../src/lib/stores/theme');
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });

        it('should initialize with light theme when stored theme is invalid', async () => {
            mockLocalStorage['vc-theme'] = 'invalid-theme';
            
            const { theme } = await import('../../src/lib/stores/theme');
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });
    });

    describe('ToggleTheme function', () => {
        it('should toggle from light to dark theme', async () => {
            const { theme, ToggleTheme } = await import('../../src/lib/stores/theme');
            
            expect(get(theme)).toBe('light');
            
            ToggleTheme();
            
            expect(get(theme)).toBe('dark');
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', true);
        });

        it('should toggle from dark to light theme', async () => {
            mockLocalStorage['vc-theme'] = 'dark';
            
            const { theme, ToggleTheme } = await import('../../src/lib/stores/theme');
            
            expect(get(theme)).toBe('dark');
            
            ToggleTheme();
            
            expect(get(theme)).toBe('light');
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', false);
        });

        it('should handle multiple toggles correctly', async () => {
            const { theme, ToggleTheme } = await import('../../src/lib/stores/theme');
            
            expect(get(theme)).toBe('light'); // Initial state
            
            ToggleTheme(); // light -> dark
            expect(get(theme)).toBe('dark');
            
            ToggleTheme(); // dark -> light
            expect(get(theme)).toBe('light');
            
            ToggleTheme(); // light -> dark
            expect(get(theme)).toBe('dark');
            
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledTimes(3);
        });

        it('should update DOM class correctly', async () => {
            const { ToggleTheme } = await import('../../src/lib/stores/theme');
            
            ToggleTheme(); // light -> dark
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', true);
            
            ToggleTheme(); // dark -> light
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', false);
        });
    });

    describe('theme store subscription', () => {
        it('should support subscription to theme changes', async () => {
            const { theme, ToggleTheme } = await import('../../src/lib/stores/theme');
            
            const themes: string[] = [];
            
            const unsubscribe = theme.subscribe((currentTheme) => {
                themes.push(currentTheme);
            });
            
            ToggleTheme(); // light -> dark
            ToggleTheme(); // dark -> light
            
            expect(themes).toEqual(['light', 'dark', 'light']);
            
            unsubscribe();
        });

        it('should handle manual theme updates', async () => {
            const { theme } = await import('../../src/lib/stores/theme');
            
            theme.set('dark');
            expect(get(theme)).toBe('dark');
            
            theme.set('light');
            expect(get(theme)).toBe('light');
        });
    });

    describe('environment compatibility', () => {
        it('should handle missing localStorage gracefully', async () => {
            // Mock localStorage as undefined for SSR environment
            Object.defineProperty(global, 'localStorage', {
                value: undefined,
                writable: true,
                configurable: true
            });
            
            const { theme } = await import('../../src/lib/stores/theme');
            expect(get(theme)).toBe('light'); // Should default to light
        });

        it('should handle ToggleTheme function correctly', async () => {
            const { theme, ToggleTheme } = await import('../../src/lib/stores/theme');
            
            // Should work with proper document mock
            expect(() => ToggleTheme()).not.toThrow();
            expect(get(theme)).toBe('dark'); // Should toggle to dark
        });
    });
});

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { get } from 'svelte/store';
import { theme, ToggleTheme } from '../src/lib/stores/theme';

describe('Theme Store', () => {
    let mockSessionStorage: Record<string, string>;
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
        mockSessionStorage = {};
        mockDocument = {
            documentElement: {
                classList: {
                    add: vi.fn(),
                    toggle: vi.fn()
                }
            }
        };

        // Mock sessionStorage
        Object.defineProperty(global, 'sessionStorage', {
            value: {
                getItem: vi.fn((key: string) => mockSessionStorage[key] || null),
                setItem: vi.fn((key: string, value: string) => {
                    mockSessionStorage[key] = value;
                }),
                removeItem: vi.fn((key: string) => {
                    delete mockSessionStorage[key];
                }),
                clear: vi.fn(() => {
                    mockSessionStorage = {};
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
        it('should initialize with light theme when no stored theme', () => {
            // Using imported theme           
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });

        it('should initialize with stored dark theme', async () => {
            // Set up sessionStorage mock to return 'dark' before module import
            mockSessionStorage['vc-theme'] = 'dark';
            
            // Re-import the module to test initialization with mocked sessionStorage
            vi.resetModules();
            const { theme: freshTheme } = await import('../src/lib/stores/theme.ts');
            
            const currentTheme = get(freshTheme);
            expect(currentTheme).toBe('dark');
        });

        it('should initialize with light theme when stored theme is not dark', () => {
            mockSessionStorage['vc-theme'] = 'light';
            
            // Using imported theme
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });

        it('should initialize with light theme when stored theme is invalid', () => {
            mockSessionStorage['vc-theme'] = 'invalid-theme';
            
            // Using imported theme
            
            const currentTheme = get(theme);
            expect(currentTheme).toBe('light');
        });
    });

    describe('ToggleTheme function', () => {
        it('should toggle from light to dark theme', () => {

            // Using imported theme and ToggleTheme
            
            expect(get(theme)).toBe('light');
            
            ToggleTheme();
            
            expect(get(theme)).toBe('dark');
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', true);
        });

        it('should toggle from dark to light theme', () => {
            mockSessionStorage['vc-theme'] = 'dark';
            
            // Using imported theme and ToggleTheme
            
            expect(get(theme)).toBe('dark');
            
            ToggleTheme();
            
            expect(get(theme)).toBe('light');
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', false);
        });

        it('should handle multiple toggles correctly', () => {

            // Using imported theme and ToggleTheme
            
            expect(get(theme)).toBe('light'); // Initial state
            
            ToggleTheme(); // light -> dark
            expect(get(theme)).toBe('dark');
            
            ToggleTheme(); // dark -> light
            expect(get(theme)).toBe('light');
            
            ToggleTheme(); // light -> dark
            expect(get(theme)).toBe('dark');
            
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledTimes(3);
        });

        it('should update DOM class correctly', () => {
            // Ensure we start with light theme
            theme.set('light');
            
            // Clear any previous mock calls
            mockDocument.documentElement.classList.toggle.mockClear();
            
            ToggleTheme(); // light -> dark
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', true);
            
            ToggleTheme(); // dark -> light
            expect(mockDocument.documentElement.classList.toggle).toHaveBeenCalledWith('dark-mode', false);
        });
    });

    describe('theme store subscription', () => {
        it('should support subscription to theme changes', () => {
            const themes: string[] = [];
            
            const unsubscribe = theme.subscribe((currentTheme) => {
                themes.push(currentTheme);
            });
            
            ToggleTheme(); // light -> dark
            ToggleTheme(); // dark -> light
            
            expect(themes).toEqual(['light', 'dark', 'light']);
            
            unsubscribe();
        });

        it('should handle manual theme updates', () => {
          
            // Using imported theme
            
            theme.set('dark');
            expect(get(theme)).toBe('dark');
            
            theme.set('light');
            expect(get(theme)).toBe('light');
        });
    });

    describe('environment compatibility', () => {
        it('should handle missing sessionStorage gracefully', () => {
            // Mock sessionStorage as undefined for SSR environment
            Object.defineProperty(global, 'sessionStorage', {
                value: undefined,
                writable: true,
                configurable: true
            });
            

            // Using imported theme

            expect(get(theme)).toBe('light'); // Should default to light
        });

        it('should handle ToggleTheme function correctly', () => {

            // Using imported theme and ToggleTheme
            
            // Should work with proper document mock
            expect(() => ToggleTheme()).not.toThrow();
            expect(get(theme)).toBe('dark'); // Should toggle to dark
        });
    });
});

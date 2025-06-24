import { render, screen, fireEvent} from '@testing-library/svelte';
import { tick } from 'svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import NavBar from '../../src/lib/components/main/nav-bar.svelte';

// --- Mock SvelteKit's navigation (this pattern is fine) ---
const mockedGoto = vi.fn();
vi.mock('$app/navigation', () => ({
  goto: (path: string) => mockedGoto(path)
}));

// --- Mock the Theme Store (Correct Hoisting-Safe Pattern) ---
vi.mock('../../src/lib/stores/theme', async (importOriginal) => {
  const { writable } = await import('svelte/store');
  const mockThemeStore = writable<'light' | 'dark'>('light');
  return {
    // This is the mock that the component will use
    theme: mockThemeStore,
    ToggleTheme: vi.fn(() => {
      mockThemeStore.update((current) => (current === 'light' ? 'dark' : 'light'));
    })
  };
});

// We need to get a reference to the mocks that were created above.
const { theme: mockThemeStore, ToggleTheme: mockedToggleTheme } = await import(
  '../../src/lib/stores/theme'
);

describe('NavBar Component', () => {
  beforeEach(() => {
    // Reset our mock references before each test
    vi.clearAllMocks();
    mockThemeStore.set('light');
  });

  it('TestInitialRender_Success: Renders the title and buttons', () => {
    render(NavBar);
    expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Toggle theme' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Logout' })).toBeInTheDocument();
  });

  it('TestLogoutButton_Success: Calls goto with the correct path on click', async () => {
    render(NavBar);
    const logoutButton = screen.getByRole('button', { name: 'Logout' });
    await fireEvent.click(logoutButton);
    expect(mockedGoto).toHaveBeenCalledWith('/auth-page');
  });

  it('TestThemeToggle_Success: Calls the ToggleTheme function on click', async () => {
    render(NavBar);
    const themeButton = screen.getByRole('button', { name: 'Toggle theme' });
    await fireEvent.click(themeButton);
    expect(mockedToggleTheme).toHaveBeenCalled();
  });

  it('TestThemeIconChange_Success: Displays the correct theme icon when the store changes', async () => {
    render(NavBar);

    // Initially, check for the light mode icon
    expect(screen.getByAltText('Light mode')).toBeInTheDocument();
    expect(screen.queryByAltText('Dark mode')).toBeNull();

    // Manually update our mock store's value to 'dark'
    mockThemeStore.set('dark');

    // Wait for Svelte to update the DOM
    await tick();

    // Now, the dark mode icon should be visible and the light one should be gone
    expect(screen.getByAltText('Dark mode')).toBeInTheDocument();
    expect(screen.queryByAltText('Light mode')).toBeNull();
  });
});
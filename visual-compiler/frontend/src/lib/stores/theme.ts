import { writable } from 'svelte/store';

const stored_theme = typeof localStorage !== 'undefined' ? localStorage.getItem('vc-theme') : null;
const initial_theme = stored_theme === 'dark' ? 'dark' : 'light';

export const theme = writable<'light' | 'dark'>(initial_theme);

// ToggleTheme
// Return type: void
// Parameter type(s): none
// Toggles the application theme between 'light' and 'dark' and persists the choice.
export function ToggleTheme() {
  theme.update(current => {
    const new_theme = current === 'light' ? 'dark' : 'light';
    // This toggles the class that our custom CSS will use.
    document.documentElement.classList.toggle('dark-mode', new_theme === 'dark');
    localStorage.setItem('vc-theme', new_theme);
    return new_theme;
  });
}

// Initialize on load
if (typeof document !== 'undefined') {
  // We just make sure the correct class is present on page load.
  if (initial_theme === 'dark') {
    document.documentElement.classList.add('dark-mode');
  }
}
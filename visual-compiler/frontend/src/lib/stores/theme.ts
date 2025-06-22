import { writable } from 'svelte/store';

const storedTheme = typeof localStorage !== 'undefined' ? localStorage.getItem('vc-theme') : null;
const initialTheme = storedTheme === 'dark' ? 'dark' : 'light';

export const theme = writable<'light' | 'dark'>(initialTheme);

export function toggleTheme() {
  theme.update(current => {
    const newTheme = current === 'light' ? 'dark' : 'light';
    // This line is the only part that should be inside the toggle function.
    // It toggles the class that our custom CSS will use.
    document.documentElement.classList.toggle('dark-mode', newTheme === 'dark');
    localStorage.setItem('vc-theme', newTheme);
    return newTheme;
  });
}

// Initialize on load
if (typeof document !== 'undefined') {
  // We no longer set the svelvet-theme attribute here.
  // We just make sure the correct class is present on page load.
  if (initialTheme === 'dark') {
    document.documentElement.classList.add('dark-mode');
  }
}

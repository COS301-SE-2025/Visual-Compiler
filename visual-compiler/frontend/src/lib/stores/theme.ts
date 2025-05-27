import { writable } from 'svelte/store';

const storedTheme = typeof localStorage !== 'undefined' ? localStorage.getItem('vc-theme') : null;
const initialTheme = storedTheme === 'dark' ? 'dark' : 'light';

export const theme = writable<'light' | 'dark'>(initialTheme);

export function toggleTheme() {
  theme.update(current => {
    const newTheme = current === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('svelvet-theme', newTheme);
    document.documentElement.classList.toggle('dark-mode', newTheme === 'dark');
    localStorage.setItem('vc-theme', newTheme);
    return newTheme;
  });
}

// Initialize on load
if (typeof document !== 'undefined') {
  document.documentElement.setAttribute('svelvet-theme', initialTheme);
  if (initialTheme === 'dark') {
    document.documentElement.classList.add('dark-mode');
  }
}
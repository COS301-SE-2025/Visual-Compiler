import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info';

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
  duration?: number;
}

export const toasts = writable<Toast[]>([]);

let counter = 0;

// AddToast
// Return type: void
// Parameter type(s): string, ToastType, number
// Adds a new toast notification to the store and sets a timer to remove it.
export function AddToast(message: string, type: ToastType = 'info', duration = 2000) {
  const id = counter++;
  toasts.update(all => [...all, { id, message, type, duration }]);

  setTimeout(() => {
    toasts.update(all => all.filter(t => t.id !== id));
  }, duration);
}
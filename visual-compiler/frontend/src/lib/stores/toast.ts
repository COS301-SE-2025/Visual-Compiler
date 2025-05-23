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

export function addToast(message: string, type: ToastType = 'info', duration = 3000) {
  const id = counter++;
  toasts.update(all => [...all, { id, message, type, duration }]);

  setTimeout(() => {
    toasts.update(all => all.filter(t => t.id !== id));
  }, duration);
}

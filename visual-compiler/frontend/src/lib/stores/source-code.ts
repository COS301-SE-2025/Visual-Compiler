import { writable } from 'svelte/store';

export const confirmedSourceCode = writable('');

export function resetSourceCode() {
    confirmedSourceCode.set('');
}
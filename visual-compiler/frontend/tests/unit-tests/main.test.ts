import {describe,expect,vi } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen,waitFor } from '@testing-library/svelte';
import page_comp from '../../src/routes/main/+page.svelte';
import toasts from '../../src/lib/components/ToastContainer.svelte';
import { theme } from '$lib/stores/theme';
import { addToast } from '$lib/stores/toast';

// Mock stores and dependencies
vi.mock('$lib/stores/theme', () => ({
  theme: {
    subscribe: vi.fn((fn) => {
      fn('light');
      return { unsubscribe: vi.fn() };
    })
  }
}));

vi.mock('$lib/stores/toast', () => ({
  addToast: vi.fn()
}));

const mockSuccessResponse = {
    ok:true,
    json: async()=> ({message:"Success"})
};
const mockFailedResponse = {
    ok: false,
    json: async () => ({error: "Failed"})
};

describe("Testing main page", () => {
    it("Page rendering", () => {
        render(page_comp);
        expect(screen.getByRole('banner')).toBeInTheDocument();
        expect(screen.getByText('Visual Compiler')).toBeInTheDocument();
        expect(screen.getByRole('img', {name: 'Visual Compiler logo'})).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Toggle theme'})).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Logout'})).toBeInTheDocument();
        expect(screen.getByText('Blocks')).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Source Code Start here. Add source code to begin compilation.'})).toBeInTheDocument();
        expect(screen.getByRole('button', {name: 'Lexer Converts source code into tokens for processing.'})).toBeInTheDocument();
    })

    /*it("Test CreateNode-SourceCode",async() => {
        render(page_comp);
    
        const button = screen.getByRole('button', {name: 'Source Code Start here. Add source code to begin compilation.'});
        await fireEvent.click(button);
       expect(screen.getByRole('button', { name: /Source Code Node/i })).toBeInTheDocument();

    })
    it("Test CreateNode-Lexer",async() => {
        render(page_comp);
    
        const button = screen.getByRole('button', {name: 'Lexer Converts source code into tokens for processing.'});
        await fireEvent.click(button);
       expect(screen.getByRole('button', { name: /Source Code Node/i })).toBeInTheDocument();

    })*/
});
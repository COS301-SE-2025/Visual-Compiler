import {describe,expect,vi } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen} from '@testing-library/svelte';
import page_comp from '../../src/routes/main/+page.svelte';

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

    it("Test CreateNode-SourceCode",async() => {
        render(page_comp);
    
        const button = screen.getByRole('button', {name: 'Source Code Start here. Add source code to begin compilation.'});
        await fireEvent.click(button);
        expect(document.getElementById('N-source-1')).toBeInTheDocument();

    })

    it("Test CreateNode-Lexer",async() => {
        render(page_comp);
    
        const button = screen.getByRole('button', {name: 'Lexer Converts source code into tokens for processing.'});
        await fireEvent.click(button);
        expect(document.getElementById('N-lexer-1')).toBeInTheDocument();

    })

    it("Test Create Multiple Nodes", async() => {
        render(page_comp);
    
        const button = screen.getByRole('button', {name: 'Source Code Start here. Add source code to begin compilation.'});
        await fireEvent.click(button);
        expect(document.getElementById('N-source-1')).toBeInTheDocument();

        const lbutton = screen.getByRole('button', {name: 'Lexer Converts source code into tokens for processing.'});
        await fireEvent.click(lbutton);
        expect(document.getElementById('N-lexer-2')).toBeInTheDocument();
    })
});
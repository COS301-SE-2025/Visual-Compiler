// tests/unit-tests/parsing-input.test.ts

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ParsingInput from '../../src/lib/components/parser/parsing-input.svelte';
import type { SyntaxTree } from '$lib/types';

// Mock the toast store and fetch API
vi.mock('$lib/stores/toast', () => ({
  AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage
const localStorageMock = (() => {
  let store: { [key: string]: string } = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => (store[key] = value.toString()),
    clear: () => (store = {})
  };
})();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('ParsingInput Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.localStorage.setItem('user_id', 'test-user-parser-123'); // Simulate logged-in user
  });

  // TestInitialRender_Success
  it('TestInitialRender_Success: Renders the grammar editor with one empty rule', () => {
    render(ParsingInput);
    expect(screen.getByText('Context-Free Grammar')).toBeInTheDocument();
    expect(screen.getByLabelText('Variables')).toBeInTheDocument();
    expect(screen.getByLabelText('Terminals')).toBeInTheDocument();
    expect(screen.getByText('Start')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('LHS')).toBeInTheDocument();
  });

  // TestAddNewRule_Success
  it('TestAddNewRule_Success: Adds a new rule row on button click', async () => {
    render(ParsingInput);
    const addRuleButton = screen.getByRole('button', { name: '+ Add New Rule' });
    await fireEvent.click(addRuleButton);
    const nonTerminalInputs = screen.getAllByPlaceholderText('LHS');
    expect(nonTerminalInputs.length).toBe(2);
  });

  // TestInsertDefaultGrammar_Success
  it('TestInsertDefaultGrammar_Success: Fills the form with default grammar on button click', async () => {
    render(ParsingInput);
    const insertDefaultButton = screen.getByRole('button', { name: 'Insert default grammar' });
    await fireEvent.click(insertDefaultButton);

    const variablesInput = screen.getByLabelText('Variables') as HTMLInputElement;
    const terminalsInput = screen.getByLabelText('Terminals') as HTMLInputElement;

    expect(variablesInput.value).toBe('STATEMENT, DECLARATION, EXPRESSION, TYPE, TERM');
    expect(terminalsInput.value).toBe('KEYWORD, IDENTIFIER, ASSIGNMENT, INTEGER, OPERATOR, SEPARATOR');
    expect(screen.getByDisplayValue('STATEMENT')).toBeInTheDocument();
  });

  // TestValidation_Failure_NoStartSymbol
  it('TestValidation_Failure_NoStartSymbol: Shows an error if the start symbol is not in the variables list', async () => {
    render(ParsingInput);
    await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'B, C' } });
    await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'A' } });
    
    const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
    await fireEvent.click(submitButton);

    expect(AddToast).toHaveBeenCalledWith("The start symbol 'A' must be included in the Variables list.", 'error');
  });

  // TestSubmitGrammar_Success
  it('TestSubmitGrammar_Success: Submits grammar and shows generate button on success', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ message: 'Grammar submitted successfully!' })
    });

    render(ParsingInput);

    // Fill out a valid grammar
    await fireEvent.input(screen.getByLabelText('Variables'), { target: { value: 'S' } });
    await fireEvent.input(screen.getByLabelText('Terminals'), { target: { value: 'a' } });
    await fireEvent.input(screen.getByPlaceholderText('LHS'), { target: { value: 'S' } });
    await fireEvent.input(screen.getByPlaceholderText('RHS'), { target: { value: 'a' } });

    const submitButton = screen.getByRole('button', { name: 'Submit Grammar' });
    await fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/api/parsing/grammar', expect.any(Object));
      expect(AddToast).toHaveBeenCalledWith('Grammar submitted successfully!', 'success');
    });

    expect(await screen.findByRole('button', { name: 'Generate Syntax Tree' })).toBeInTheDocument();
  });

  
});

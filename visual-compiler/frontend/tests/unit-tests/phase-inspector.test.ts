// tests/unit-tests/phase-inspector.test.ts

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../../src/lib/components/lexer/phase-inspector.svelte';
import type { Token } from '$lib/types';

// Mock the toast store and fetch API
vi.mock('$lib/stores/toast', () => ({
  AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

describe('PhaseInspector Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const sourceCode = 'let x = 1;';

  // This test now passes the 'source_code' prop correctly.
  it('TestInitialRender_Success: Renders correctly and shows form on button click', async () => {
    render(PhaseInspector, { source_code: sourceCode });
    expect(screen.getByText(sourceCode)).toBeInTheDocument();

    const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
    await fireEvent.click(regexButton);

    expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
  });

  // This test and all subsequent tests are restored.
  it('TestValidation_Failure_InvalidRegex: Shows an error for an invalid regex pattern', async () => {
    render(PhaseInspector, { source_code: sourceCode });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    const submitButton = screen.getByRole('button', { name: 'Submit' });

    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: '[' } });
    await fireEvent.click(submitButton);

    expect(await screen.findByText('Invalid regular expression pattern')).toBeInTheDocument();
    expect(mockFetch).not.toHaveBeenCalled();
  });

  it('TestValidation_Failure_IncompleteRow: Shows an error for an incomplete row', async () => {
    render(PhaseInspector, { source_code: sourceCode });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const submitButton = screen.getByRole('button', { name: 'Submit' });

    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.click(submitButton);

    expect(await screen.findByText('Please fill in both Type and Regular Expression')).toBeInTheDocument();
  });

  it('TestAddRow_Success: Adds a new row when the plus button is clicked', async () => {
    render(PhaseInspector, { source_code: sourceCode });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'TYPE1' } });
    await fireEvent.input(regexInput, { target: { value: 'REGEX1' } });
    
    const addButton = await screen.findByRole('button', { name: '+' });
    await fireEvent.click(addButton);
    
    expect(screen.getAllByPlaceholderText('Enter type...').length).toBe(2);
  });

  it('TestHandleSubmit_Success: Calls fetch and shows generate button on valid submission', async () => {
    mockFetch.mockResolvedValue({ ok: true, json: () => Promise.resolve({ message: 'Success' }) });
    
    render(PhaseInspector, { source_code: sourceCode });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'if|else' } });
    await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/api/lexing/code', expect.any(Object));
    });
    
    expect(await screen.findByRole('button', { name: 'Generate Tokens' })).toBeInTheDocument();
  });
  
  it('TestHandleSubmit_Failure_ServerError: Shows error toast on fetch failure', async () => {
    mockFetch.mockResolvedValue({ ok: false, status: 500, text: () => Promise.resolve('Internal Server Error') });
    
    render(PhaseInspector, { source_code: sourceCode });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    await fireEvent.input(screen.getByPlaceholderText('Enter type...'), { target: { value: 'KEYWORD' } });
    await fireEvent.input(screen.getByPlaceholderText('Enter regex pattern...'), { target: { value: 'if' } });
    await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));
    
    await waitFor(() => {
      expect(AddToast).toHaveBeenCalledWith('Cannot connect to server. Please ensure the backend is running.', 'error');
    });
  });

  // This test now uses the correct data structure and event handling logic.
  it('TestGenerateTokens_Success: Calls generate endpoint and triggers correct callback', async () => {
    // The mock server now provides the 'tokens_unidentified' key.
    const mockServerResponse = {
      tokens: [{ Type: 'KEYWORD', Value: 'if' }],
      tokens_unidentified: ['@']
    };
    
    // The component should map 'tokens_unidentified' to 'unexpected_tokens'.
    const expectedCallbackPayload = {
      tokens: mockServerResponse.tokens,
      unexpected_tokens: mockServerResponse.tokens_unidentified
    };

    mockFetch
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ message: 'Success' }) })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(mockServerResponse) });

    const handleGenerate = vi.fn();
    
    render(PhaseInspector, {
      source_code: sourceCode,
      onGenerateTokens: handleGenerate
    });
    
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    await fireEvent.input(screen.getByPlaceholderText('Enter type...'), { target: { value: 'KEYWORD' } });
    await fireEvent.input(screen.getByPlaceholderText('Enter regex pattern...'), { target: { value: 'if' } });
    await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

    const generateButton = await screen.findByRole('button', { name: 'Generate Tokens' });
    await fireEvent.click(generateButton);

    await waitFor(() => {
      expect(handleGenerate).toHaveBeenCalledWith(expectedCallbackPayload);
    });
  });
});
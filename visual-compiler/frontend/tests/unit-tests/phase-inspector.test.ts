import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../../src/lib/components/lexor/phase-inspector.svelte';
import type { Token } from '$lib/types';

// Mock the toast store and fetch API
vi.mock('$lib/stores/toast', () => ({
  AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

const mockFetch = vi.fn();
vi.stubGlobal('fetch', mockFetch);

describe('PhaseInspector Component', () => {
  // Reset mocks before each test
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const sourceCode = 'let x = 1;';

  it('TestInitialRender_Success: Renders correctly and shows form on button click', async () => {
    render(PhaseInspector, { props: { source_code: sourceCode } });
    expect(screen.getByText(sourceCode)).toBeInTheDocument();
    
    // Initially, the inputs should NOT be visible
    expect(screen.queryByPlaceholderText('Enter type...')).toBeNull();

    // Find and click the button to show the form
    const regexButton = screen.getByRole('button', { name: 'Regular Expression' });
    await fireEvent.click(regexButton);

    // NOW, the inputs should be visible
    expect(screen.getByPlaceholderText('Enter type...')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Enter regex pattern...')).toBeInTheDocument();
  });

  it('TestValidation_Failure_InvalidRegex: Shows an error for an invalid regex pattern', async () => {
    render(PhaseInspector, { props: { source_code: sourceCode } });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    const submitButton = screen.getByRole('button', { name: 'Submit' });

    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: '[' } }); // Invalid regex
    await fireEvent.click(submitButton);

    expect(await screen.findByText('Invalid regular expression pattern')).toBeInTheDocument();
    expect(mockFetch).not.toHaveBeenCalled();
  });

  it('TestValidation_Failure_IncompleteRow: Shows an error for an incomplete row', async () => {
    render(PhaseInspector, { props: { source_code: sourceCode } });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const submitButton = screen.getByRole('button', { name: 'Submit' });

    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.click(submitButton);

    expect(await screen.findByText('Please fill in both Type and Regular Expression')).toBeInTheDocument();
    expect(mockFetch).not.toHaveBeenCalled();
  });

  it('TestAddRow_Success: Adds a new row when the plus button is clicked', async () => {
    render(PhaseInspector, { props: { source_code: sourceCode } });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));
    
    expect(screen.getAllByPlaceholderText('Enter type...').length).toBe(1);

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'TYPE1' } });
    await fireEvent.input(regexInput, { target: { value: 'REGEX1' } });
    
    const addButton = await screen.findByRole('button', { name: '+' });
    await fireEvent.click(addButton);
    
    expect(screen.getAllByPlaceholderText('Enter type...').length).toBe(2);
  });

  it('TestHandleSubmit_Success: Calls fetch and shows generate button on valid submission', async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ message: 'Success' })
    });

    render(PhaseInspector, { props: { source_code: sourceCode } });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    const submitButton = screen.getByRole('button', { name: 'Submit' });

    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'if|else' } });
    await fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/api/lexing/code', expect.any(Object));
    });
    
    expect(await screen.findByRole('button', { name: 'Generate Tokens' })).toBeInTheDocument();
  });
  
  it('TestHandleSubmit_Failure_ServerError: Shows error toast on fetch failure', async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 500,
      text: () => Promise.resolve('Internal Server Error')
    });
    
    render(PhaseInspector, { props: { source_code: sourceCode } });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'if|else' } });
    
    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);
    
    await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith('Cannot connect to server. Please ensure the backend is running.', 'error');
    });
  });

  it('TestGenerateTokens_Success: Calls generate endpoint and triggers onGenerateTokens event', async () => {
    // UPDATED THE KEY NAME IN THIS MOCK OBJECT
    const mockTokenPayload = {
      tokens: [{ Type: 'KEYWORD', Value: 'if' }],
      unexpected_tokens: ['@'] 
    };
    
    // Mock first fetch for handleSubmit, second for generateTokens
    mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ message: 'Success' }) })
             .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(mockTokenPayload) });

    const handleGenerate = vi.fn();
    render(PhaseInspector, {
      props: {
        source_code: sourceCode,
        onGenerateTokens: handleGenerate
      }
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Regular Expression' }));

    // Step 1: Submit valid data to make the 'Generate Tokens' button appear
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'if' } });
    await fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

    // Step 2: Click the 'Generate Tokens' button
    const generateButton = await screen.findByRole('button', { name: 'Generate Tokens' });
    await fireEvent.click(generateButton);

    // Assert that the onGenerateTokens prop function was called with the payload
    await waitFor(() => {
        expect(handleGenerate).toHaveBeenCalledWith(mockTokenPayload);
    });

    // Assert fetch was called twice, the second time to the lexer endpoint
    expect(mockFetch).toHaveBeenCalledTimes(2);
    expect(mockFetch).toHaveBeenLastCalledWith('http://localhost:8080/api/lexing/lexer', expect.any(Object));
});
});
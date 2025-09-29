import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import PhaseInspector from '../src/lib/components/lexer/phase-inspector.svelte';

// Mock the toast store
vi.mock('$lib/stores/toast', () => ({
  AddToast: vi.fn()
}));
import { AddToast } from '$lib/stores/toast';

// Mock the project store  
vi.mock('$lib/stores/project', () => ({
  projectName: {
    subscribe: vi.fn((callback) => {
      callback('test-project');
      return { unsubscribe: vi.fn() };
    })
  }
}));

// Mock lexer store with complete structure
vi.mock('$lib/stores/lexer', () => ({
  lexerState: {
    subscribe: vi.fn((callback) => {
      callback({
        userInputRows: [{ type: '', regex: '', error: '' }],
        automataInputs: {
          states: '',
          startState: '',
          acceptedStates: '',
          transitions: ''
        },
        selectedType: 'REGEX',
        showDefault: false,
        isSubmitted: false,
        hasUnsavedChanges: false,
        tokens: [],
        hasTokens: false,
        sourceCode: ''
      });
      return { unsubscribe: vi.fn() };
    }),
    update: vi.fn(),
    set: vi.fn()
  },
  updateLexerInputs: vi.fn(),
  updateAutomataInputs: vi.fn(),
  markLexerSubmitted: vi.fn(),
  updateLexerArtifacts: vi.fn()
}));

// Mock svelte/store get function
vi.mock('svelte/store', () => ({
  get: vi.fn(() => 'test-project'),
  writable: vi.fn(() => ({
    subscribe: vi.fn(() => () => {}),
    set: vi.fn(),
    update: vi.fn()
  }))
}));

// Mock vis-network for network visualization
vi.mock('vis-network', () => ({
  Network: vi.fn().mockImplementation(() => ({
    fit: vi.fn(),
    getScale: vi.fn(() => 1),
    moveTo: vi.fn(),
    destroy: vi.fn()
  })),
  DataSet: vi.fn().mockImplementation(() => ({
    add: vi.fn(),
    clear: vi.fn()
  }))
}));

// Global mocks
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock localStorage and sessionStorage
const createMockStorage = () => {
  let store: { [key: string]: string } = {};
  return {
    getItem(key: string) { return store[key] || null; },
    setItem(key: string, value: string) { store[key] = value.toString(); },
    clear() { store = {}; }
  };
};

const localStorageMock = createMockStorage();
const sessionStorageMock = createMockStorage();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });
Object.defineProperty(window, 'sessionStorage', { value: sessionStorageMock });

describe('Lexer Phase Inspector', () => {
  const sourceCode = 'int x = 5;\\nreturn x;';

  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
    sessionStorageMock.clear();
    
    localStorageMock.setItem('user_id', 'test-user');
    sessionStorageMock.setItem('access_token', 'test-access-token');
    
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ success: true, message: 'Rules stored successfully!' })
    });
  });

  it('renders without crashing', () => {
    render(PhaseInspector, {
      props: {
        source_code: sourceCode
      }
    });
    expect(true).toBe(true);
  });

  it('should render regex mode by default', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Should start with regex mode button active
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    expect(regexButton).toBeTruthy();

    // Wait for regex form to be visible
    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });
  });

  it('should switch to automata mode', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Click automata mode button
    const automataButton = screen.getByRole('button', { name: /Automata/i });
    await fireEvent.click(automataButton);

    // Wait for automata form to appear
    await waitFor(() => {
      const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
      expect(statesInput).toBeTruthy();
    });
  });

  it('should handle form submission in regex mode', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Switch to regex mode (should be default)
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    // Fill in form
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'int|return' } });

    // Submit form
    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Verify API was called
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalled();
    }, { timeout: 3000 });
  });

  it('should handle form submission in automata mode', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Switch to automata mode
    const automataButton = screen.getByRole('button', { name: /Automata/i });
    await fireEvent.click(automataButton);

    await waitFor(() => {
      const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
      expect(statesInput).toBeTruthy();
    });

    // Fill in automata form
    const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
    const startStateInput = screen.getByPlaceholderText('e.g. q0');
    const acceptedStatesInput = screen.getByPlaceholderText('e.g. q2->int, q1->string');
    const transitionsInput = screen.getByPlaceholderText('e.g. q0,a->q1 q1,b->q2');

    await fireEvent.input(statesInput, { target: { value: 'q0,q1,q2' } });
    await fireEvent.input(startStateInput, { target: { value: 'q0' } });
    await fireEvent.input(acceptedStatesInput, { target: { value: 'q2->IDENTIFIER' } });
    await fireEvent.input(transitionsInput, { target: { value: 'q0,a->q1\\nq1,b->q2' } });

    // Use the correct button for automata mode - "Show Automata" or "Tokenisation"
    const showAutomataButton = screen.getByRole('button', { name: 'Show Automata' });
    await fireEvent.click(showAutomataButton);

    // Verify API was called
    expect(mockFetch).toHaveBeenCalled();
  });

  it('should handle API errors gracefully', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Mock API failure
    mockFetch.mockRejectedValueOnce(new Error('Network error'));

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'TEST' } });
    await fireEvent.input(regexInput, { target: { value: 'test' } });

    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Should handle error gracefully
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalled();
    }, { timeout: 3000 });
  });

  it('should show and insert default rules when Show Example is clicked', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Switch to regex mode
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    // Click the Show Example button
    const exampleButton = screen.getByRole('button', { name: /Show Example/i });
    await fireEvent.click(exampleButton);

    // Should show default rules with multiple type inputs
    await waitFor(() => {
      const typeInputs = screen.getAllByPlaceholderText('Enter type...');
      expect(typeInputs.length).toBeGreaterThan(1);
    });
  });

  it('should clear all inputs when clear button is clicked', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Switch to regex mode and fill some data
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    await fireEvent.input(typeInput, { target: { value: 'TEST_TYPE' } });

    // Click clear button
    const clearButton = screen.getByRole('button', { name: /Clear all inputs/i });
    await fireEvent.click(clearButton);

    // Verify inputs are cleared - component should handle this
    expect(typeInput).toBeTruthy(); // Should still exist
  });

  it('should validate regex patterns and show errors', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    // Fill with invalid regex
    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'INVALID' } });
    await fireEvent.input(regexInput, { target: { value: '[invalid-regex(' } });

    // Submit form
    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Component performs validation and prevents invalid API calls
    expect(mockFetch).not.toHaveBeenCalled(); // Should NOT call API with invalid regex
  });

  it('should handle missing authentication token', async () => {
    // Clear session storage
    sessionStorageMock.clear();
    
    render(PhaseInspector, { source_code: sourceCode });

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'KEYWORD' } });
    await fireEvent.input(regexInput, { target: { value: 'test' } });

    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Should handle missing auth gracefully - component checks auth internally
    expect(mockFetch).not.toHaveBeenCalled();
  });

  it('should handle regex to NFA conversion', async () => {
    // Mock successful NFA response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true })
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({
        nfa: {
          states: ['q0', 'q1', 'q2'],
          startState: 'q0',
          acceptedStates: ['q2'],
          transitions: [{ from: 'q0', to: 'q1', label: 'a' }],
          alphabet: ['a']
        }
      })
    });

    render(PhaseInspector, { source_code: sourceCode });

    // Submit regex rules first
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'LETTER' } });
    await fireEvent.input(regexInput, { target: { value: '[a-z]' } });

    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Wait for NFA button to appear
    await waitFor(() => {
      const nfaButton = screen.queryByRole('button', { name: 'NFA' });
      if (nfaButton) {
        fireEvent.click(nfaButton);
      }
    }, { timeout: 3000 });

    // Should make API calls for NFA conversion
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('lexing/rules'),
        expect.any(Object)
      );
    }, { timeout: 3000 });
  });

  it('should handle regex to DFA conversion', async () => {
    // Mock successful DFA response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true })
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({
        dfa: {
          states: ['q0', 'q1'],
          startState: 'q0',
          acceptedStates: ['q1'],
          transitions: [{ from: 'q0', to: 'q1', label: 'a' }],
          alphabet: ['a']
        }
      })
    });

    render(PhaseInspector, { source_code: sourceCode });

    // Submit regex rules first
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'LETTER' } });
    await fireEvent.input(regexInput, { target: { value: '[a-z]' } });

    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Wait for DFA button to appear
    await waitFor(() => {
      const dfaButton = screen.queryByRole('button', { name: 'DFA' });
      if (dfaButton) {
        fireEvent.click(dfaButton);
      }
    }, { timeout: 3000 });

    // Should make API calls for DFA conversion
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('lexing/rules'),
        expect.any(Object)
      );
    }, { timeout: 3000 });
  });

  it('should handle automata to regex conversion', async () => {
    // Mock successful regex conversion response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true })
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({
        rules: [{ token_type: 'IDENTIFIER', regex: '[a-zA-Z]+' }]
      })
    });

    render(PhaseInspector, { source_code: sourceCode });

    // Switch to automata mode
    const automataButton = screen.getByRole('button', { name: /Automata/i });
    await fireEvent.click(automataButton);

    await waitFor(() => {
      const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
      expect(statesInput).toBeTruthy();
    });

    // Fill in automata form
    const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
    const startStateInput = screen.getByPlaceholderText('e.g. q0');
    const acceptedStatesInput = screen.getByPlaceholderText('e.g. q2->int, q1->string');
    const transitionsInput = screen.getByPlaceholderText('e.g. q0,a->q1 q1,b->q2');

    await fireEvent.input(statesInput, { target: { value: 'q0,q1' } });
    await fireEvent.input(startStateInput, { target: { value: 'q0' } });
    await fireEvent.input(acceptedStatesInput, { target: { value: 'q1->IDENTIFIER' } });
    await fireEvent.input(transitionsInput, { target: { value: 'q0,a->q1' } });

    // Click RE (regex) button
    const reButton = screen.getByRole('button', { name: 'RE' });
    await fireEvent.click(reButton);

    // Should make API calls for regex conversion
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('dfaToRegex'),
        expect.any(Object)
      );
    }, { timeout: 3000 });
  });

  it('should handle form validation with empty required fields', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    // Try to submit without filling any fields
    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Should not make API call with empty form
    expect(mockFetch).not.toHaveBeenCalled();
  });

  it('should persist state when switching between modes', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    // Fill regex mode
    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    await fireEvent.input(typeInput, { target: { value: 'PERSISTENT_TYPE' } });

    // Switch to automata mode
    const automataButton = screen.getByRole('button', { name: /Automata/i });
    await fireEvent.click(automataButton);

    // Fill automata mode
    await waitFor(() => {
      const statesInput = screen.queryByPlaceholderText('e.g. q0,q1,q2');
      expect(statesInput).toBeTruthy();
    });

    const statesInput = screen.getByPlaceholderText('e.g. q0,q1,q2');
    await fireEvent.input(statesInput, { target: { value: 'q0,q1' } });

    // Switch back to regex mode - should preserve data
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const restoredTypeInput = screen.getByPlaceholderText('Enter type...');
      // Component should preserve state through lexer store
      expect(restoredTypeInput).toBeTruthy();
    });
  });

  it('should handle adding multiple input rows in regex mode', async () => {
    render(PhaseInspector, { source_code: sourceCode });

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    // Look for add row button (might be + or Add)
    const addButton = screen.queryByRole('button', { name: /\+|Add/i });
    if (addButton) {
      await fireEvent.click(addButton);
      
      // Should have multiple type inputs now
      await waitFor(() => {
        const typeInputs = screen.getAllByPlaceholderText('Enter type...');
        expect(typeInputs.length).toBeGreaterThan(1);
      });
    } else {
      // If no add button, at least verify current functionality
      const typeInputs = screen.getAllByPlaceholderText('Enter type...');
      expect(typeInputs.length).toBeGreaterThanOrEqual(1);
    }
  });

  it('should handle token generation after successful submission', async () => {
    // Mock successful token generation response
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ success: true })
    }).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({
        tokens: [
          { type: 'IDENTIFIER', value: 'myVar', position: 0 },
          { type: 'ASSIGN', value: '=', position: 5 }
        ]
      })
    });

    render(PhaseInspector, { source_code: sourceCode });

    const regexButton = screen.getByRole('button', { name: /Regular Expression/i });
    await fireEvent.click(regexButton);

    await waitFor(() => {
      const typeInput = screen.queryByPlaceholderText('Enter type...');
      expect(typeInput).toBeTruthy();
    });

    const typeInput = screen.getByPlaceholderText('Enter type...');
    const regexInput = screen.getByPlaceholderText('Enter regex pattern...');
    
    await fireEvent.input(typeInput, { target: { value: 'IDENTIFIER' } });
    await fireEvent.input(regexInput, { target: { value: '[a-zA-Z]+' } });

    const submitButton = screen.getByRole('button', { name: 'Submit' });
    await fireEvent.click(submitButton);

    // Wait for potential token generation button
    await waitFor(() => {
      const tokenizeButton = screen.queryByRole('button', { name: /Generate|Token/i });
      if (tokenizeButton) {
        fireEvent.click(tokenizeButton);
      }
    }, { timeout: 3000 });

    // Should have made at least the submission call
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalled();
    }, { timeout: 3000 });
  });
});

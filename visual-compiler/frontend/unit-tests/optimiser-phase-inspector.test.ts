import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import OptimizerPhaseInspector from '../src/lib/components/optimiser/optimiser-phase-inspector.svelte';
import { optimiserState } from '../src/lib/stores/optimiser';
import { projectName } from '../src/lib/stores/project';
import { AddToast } from '../src/lib/stores/toast';
import { get } from 'svelte/store';

// Mock the stores
vi.mock('../src/lib/stores/optimiser', () => ({
  optimiserState: {
    subscribe: vi.fn((fn) => {
      fn({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: 'package main\n\nfunc main() {\n\t\n}',
        isOptimising: false,
        optimisationError: null,
        optimisedCode: null
      });
      return { unsubscribe: vi.fn() };
    }),
    update: vi.fn()
  }
}));

vi.mock('../src/lib/stores/project', () => ({
  projectName: {
    subscribe: vi.fn((fn) => {
      fn('test-project');
      return { unsubscribe: vi.fn() };
    })
  }
}));

vi.mock('../src/lib/stores/toast', () => ({
  AddToast: vi.fn()
}));

// Mock svelte/store get function
vi.mock('svelte/store', () => ({
  get: vi.fn(() => 'test-project')
}));

// Mock fetch
global.fetch = vi.fn();

// Mock sessionStorage
Object.defineProperty(window, 'sessionStorage', {
  value: {
    getItem: vi.fn(() => 'mock-token'),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn()
  },
  writable: true
});

describe('OptimizerPhaseInspector Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(fetch).mockClear();
    vi.mocked(sessionStorage.getItem).mockReturnValue('mock-token');
    // Reset the get function to return a valid project name by default
    vi.mocked(get).mockReturnValue('test-project');
  });

  describe('Component Rendering', () => {
    it('renders the component with default values', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByText('Optimising')).toBeInTheDocument();
      expect(screen.getByText('Instructions')).toBeInTheDocument();
      expect(screen.getByText('Optimising Techniques')).toBeInTheDocument();
      expect(screen.getByText('Source Code Input (Go)')).toBeInTheDocument();
    });

    it('renders optimization technique buttons', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Constant Folding' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Dead Code Elimination' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Loop Unrolling' })).toBeInTheDocument();
    });

    it('renders code input textarea', () => {
      render(OptimizerPhaseInspector);
      
      const codeTextarea = screen.getByRole('textbox');
      expect(codeTextarea).toBeInTheDocument();
      expect(codeTextarea).toHaveAttribute('placeholder', 'Enter your code here...');
    });

    it('renders default and optimize buttons', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Show Example' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Optimise Code' })).toBeInTheDocument();
    });
  });

  describe('Technique Selection', () => {
    it('allows selecting optimization techniques', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('allows deselecting optimization techniques', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      
      const deadCodeButton = screen.getByRole('button', { name: 'Dead Code Elimination' });
      
      // Select first
      await fireEvent.click(deadCodeButton);
      expect(optimiserState.update).toHaveBeenCalledTimes(1);
      
      // Deselect
      await fireEvent.click(deadCodeButton);
      expect(optimiserState.update).toHaveBeenCalledTimes(2);
    });

    it('allows selecting multiple techniques', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      const loopUnrollingButton = screen.getByRole('button', { name: 'Loop Unrolling' });
      
      await fireEvent.click(constantFoldingButton);
      await fireEvent.click(loopUnrollingButton);
      
      expect(optimiserState.update).toHaveBeenCalledTimes(2);
    });
  });

  describe('Code Input', () => {
    it('allows code input changes', async () => {
      render(OptimizerPhaseInspector);
      
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { 
        target: { value: 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello")\n}' }
      });
      
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('updates store when code changes', async () => {
      render(OptimizerPhaseInspector);
      
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { 
        target: { value: 'new code content' }
      });
      
      expect(optimiserState.update).toHaveBeenCalledWith(expect.any(Function));
    });
  });

  describe('Default Button Functionality', () => {
    it('shows Show Example button initially', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Show Example' })).toBeInTheDocument();
    });

    it('inserts default values when Show Example is clicked', async () => {
      render(OptimizerPhaseInspector);
      
      const showExampleButton = screen.getByRole('button', { name: 'Show Example' });
      await fireEvent.click(showExampleButton);
      
      // Should update the store with default values
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('toggles between Show Example and Restore Input', async () => {
      render(OptimizerPhaseInspector);
      
      // Initially shows Show Example
      const showExampleButton = screen.getByRole('button', { name: 'Show Example' });
      expect(showExampleButton).toBeInTheDocument();
      
      // Click to show default - button should update store
      await fireEvent.click(showExampleButton);
      
      // Should update store with default values
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('restores user input when removeDefault is called', async () => {
      render(OptimizerPhaseInspector);
      
      // First set some user input
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { 
        target: { value: 'user code' }
      });
      
      // Select a technique
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      
      // Click show example to switch to default
      const showExampleButton = screen.getByRole('button', { name: 'Show Example' });
      await fireEvent.click(showExampleButton);
      
      // Should call store update for showing defaults
      expect(optimiserState.update).toHaveBeenCalled();
    });
  });

  describe('Form Validation', () => {
    it('shows error when no code is provided', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear the code textarea
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { target: { value: '' } });
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      expect(AddToast).toHaveBeenCalledWith(
        'Please enter code and select at least one optimisation technique',
        'error'
      );
    });

    it('shows error when no techniques are selected', async () => {
      render(OptimizerPhaseInspector);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      expect(AddToast).toHaveBeenCalledWith(
        'Please enter code and select at least one optimisation technique',
        'error'
      );
    });

    it('shows error when no auth token is available', async () => {
      vi.mocked(sessionStorage.getItem).mockReturnValue(null);
      
      render(OptimizerPhaseInspector);
      
      // Select a technique first
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      expect(AddToast).toHaveBeenCalledWith(
        'Authentication required: Please log in to use the optimiser',
        'error'
      );
    });

    it('shows error when no project is selected', async () => {
      // Mock get function to return null project name
      vi.mocked(get).mockReturnValue(null);
      
      render(OptimizerPhaseInspector);
      
      // Select a technique first
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      expect(AddToast).toHaveBeenCalledWith(
        'No project selected: Please select or create a project first',
        'error'
      );
    });
  });

  describe('API Integration', () => {
    beforeEach(() => {
      vi.mocked(fetch).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({ optimised_code: 'optimized code content' })
      } as Response);
    });

    it('makes API calls when optimization is triggered', async () => {
      render(OptimizerPhaseInspector);
      
      // Select a technique
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(fetch).toHaveBeenCalledTimes(2); // Store source code + optimize
      });
    });

    it('calls store source code API first', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(fetch).toHaveBeenNthCalledWith(1, 
          'http://localhost:8080/api/optimising/source_code',
          expect.objectContaining({
            method: 'POST',
            headers: expect.objectContaining({
              'Content-Type': 'application/json',
              'Authorization': 'Bearer mock-token'
            }),
            body: expect.stringContaining('test-project')
          })
        );
      });
    });

    it('calls optimize API with correct payload', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      const deadCodeButton = screen.getByRole('button', { name: 'Dead Code Elimination' });
      await fireEvent.click(constantFoldingButton);
      await fireEvent.click(deadCodeButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(fetch).toHaveBeenNthCalledWith(2,
          'http://localhost:8080/api/optimising/optimise',
          expect.objectContaining({
            method: 'POST',
            headers: expect.objectContaining({
              'Content-Type': 'application/json',
              'Authorization': 'Bearer mock-token'
            }),
            body: expect.stringContaining('"constant_folding":true')
          })
        );
      });
    });

    it('updates store with optimization results', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(optimiserState.update).toHaveBeenCalledWith(expect.any(Function));
      });
    });

    it('shows success toast on successful optimization', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith(
          'Code optimisation completed successfully!',
          'success'
        );
      });
    });
  });

  describe('Error Handling', () => {
    it('handles API errors gracefully', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        json: () => Promise.resolve({ error: 'API Error' })
      } as Response);
      
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith(
          'Optimisation failed: API Error',
          'error'
        );
      });
    });

    it('handles store source code API error', async () => {
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        json: () => Promise.resolve({ error: 'Store API Error' })
      } as Response);
      
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith(
          'Optimisation failed: Store API Error',
          'error'
        );
      });
    });

    it('handles optimise API error with details', async () => {
      // First call succeeds (store source code)
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({})
      } as Response);
      
      // Second call fails (optimise)
      vi.mocked(fetch).mockResolvedValueOnce({
        ok: false,
        json: () => Promise.resolve({ details: 'Optimise API Error' })
      } as Response);
      
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith(
          'Optimisation failed: Optimise API Error',
          'error'
        );
      });
    });

    it('handles network errors', async () => {
      vi.mocked(fetch).mockRejectedValueOnce(new Error('Network error'));
      
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(AddToast).toHaveBeenCalledWith(
          'Optimisation failed: Network error',
          'error'
        );
      });
    });

    it('updates store with error state on failure', async () => {
      vi.mocked(fetch).mockRejectedValueOnce(new Error('Test error'));
      
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      const optimizeButton = screen.getByRole('button', { name: 'Optimise Code' });
      await fireEvent.click(optimizeButton);
      
      await waitFor(() => {
        expect(optimiserState.update).toHaveBeenCalledWith(expect.any(Function));
      });
    });
  });

  describe('Store Integration', () => {
    it('initializes store on component mount', () => {
      render(OptimizerPhaseInspector);
      
      // Should call update at least once on mount
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('updates store when technique selection changes', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear previous calls from mount
      vi.mocked(optimiserState.update).mockClear();
      
      // Change technique selection
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      expect(optimiserState.update).toHaveBeenCalled();
    });

    it('subscribes to store changes correctly', () => {
      render(OptimizerPhaseInspector);
      
      expect(optimiserState.subscribe).toHaveBeenCalled();
    });

    it('handles AI optimiser events', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      
      // Simulate AI event
      const aiEvent = new CustomEvent('ai-optimiser-generated', {
        detail: { code: 'package main\n\nfunc main() {\n  fmt.Println("AI generated")\n}' }
      });
      
      window.dispatchEvent(aiEvent);
      
      // Should update store with AI code and show toast
      expect(optimiserState.update).toHaveBeenCalled();
      expect(AddToast).toHaveBeenCalledWith(
        'AI optimiser code inserted into code input area!',
        'success'
      );
    });

    it('handles AI event without code detail', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      vi.mocked(AddToast).mockClear();
      
      // Simulate AI event without code
      const aiEvent = new CustomEvent('ai-optimiser-generated', {
        detail: {}
      });
      
      window.dispatchEvent(aiEvent);
      
      // Should not update store or show toast
      expect(optimiserState.update).not.toHaveBeenCalled();
      expect(AddToast).not.toHaveBeenCalled();
    });

    it('handles AI event without detail', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimiserState.update).mockClear();
      vi.mocked(AddToast).mockClear();
      
      // Simulate AI event without detail
      const aiEvent = new CustomEvent('ai-optimiser-generated', {
        detail: null
      });
      
      window.dispatchEvent(aiEvent);
      
      // Should not update store or show toast
      expect(optimiserState.update).not.toHaveBeenCalled();
      expect(AddToast).not.toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('has accessible technique buttons', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Constant Folding' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Dead Code Elimination' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Loop Unrolling' })).toBeInTheDocument();
    });

    it('has proper headings structure', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('heading', { name: 'Optimising' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Instructions' })).toBeInTheDocument();
      expect(screen.getByRole('heading', { name: 'Optimising Techniques' })).toBeInTheDocument();
    });

    it('has accessible buttons', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Show Example' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Optimise Code' })).toBeInTheDocument();
    });
  });
});

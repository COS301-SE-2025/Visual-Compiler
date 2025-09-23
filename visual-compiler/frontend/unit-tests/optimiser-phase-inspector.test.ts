import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import OptimizerPhaseInspector from '../src/lib/components/optimizer/optimizer-phase-inspector.svelte';
import { optimizerState } from '../src/lib/stores/optimizer';
import { projectName } from '../src/lib/stores/project';
import { AddToast } from '../src/lib/stores/toast';

// Mock the stores
vi.mock('../src/lib/stores/optimizer', () => ({
  optimizerState: {
    subscribe: vi.fn((fn) => {
      fn({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: 'package main\n\nfunc main() {\n\t\n}',
        isOptimizing: false,
        optimizationError: null,
        optimizedCode: null
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
      
      expect(screen.getByRole('button', { name: 'Insert default input' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Optimise Code' })).toBeInTheDocument();
    });
  });

  describe('Technique Selection', () => {
    it('allows selecting optimization techniques', async () => {
      render(OptimizerPhaseInspector);
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      expect(optimizerState.update).toHaveBeenCalled();
    });

    it('allows deselecting optimization techniques', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimizerState.update).mockClear();
      
      const deadCodeButton = screen.getByRole('button', { name: 'Dead Code Elimination' });
      
      // Select first
      await fireEvent.click(deadCodeButton);
      expect(optimizerState.update).toHaveBeenCalledTimes(1);
      
      // Deselect
      await fireEvent.click(deadCodeButton);
      expect(optimizerState.update).toHaveBeenCalledTimes(2);
    });

    it('allows selecting multiple techniques', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear mount calls
      vi.mocked(optimizerState.update).mockClear();
      
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      const loopUnrollingButton = screen.getByRole('button', { name: 'Loop Unrolling' });
      
      await fireEvent.click(constantFoldingButton);
      await fireEvent.click(loopUnrollingButton);
      
      expect(optimizerState.update).toHaveBeenCalledTimes(2);
    });
  });

  describe('Code Input', () => {
    it('allows code input changes', async () => {
      render(OptimizerPhaseInspector);
      
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { 
        target: { value: 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello")\n}' }
      });
      
      expect(optimizerState.update).toHaveBeenCalled();
    });

    it('updates store when code changes', async () => {
      render(OptimizerPhaseInspector);
      
      const codeTextarea = screen.getByRole('textbox');
      await fireEvent.input(codeTextarea, { 
        target: { value: 'new code content' }
      });
      
      expect(optimizerState.update).toHaveBeenCalledWith(expect.any(Function));
    });
  });

  describe('Default Button Functionality', () => {
    it('shows Insert Default button initially', () => {
      render(OptimizerPhaseInspector);
      
      expect(screen.getByRole('button', { name: 'Insert default input' })).toBeInTheDocument();
    });

    it('inserts default values when Insert Default is clicked', async () => {
      render(OptimizerPhaseInspector);
      
      const insertDefaultButton = screen.getByRole('button', { name: 'Insert default input' });
      await fireEvent.click(insertDefaultButton);
      
      // Should update the store with default values
      expect(optimizerState.update).toHaveBeenCalled();
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
        'Authentication required: Please log in to use the optimizer',
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
        expect(optimizerState.update).toHaveBeenCalledWith(expect.any(Function));
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
        expect(optimizerState.update).toHaveBeenCalledWith(expect.any(Function));
      });
    });
  });

  describe('Store Integration', () => {
    it('initializes store on component mount', () => {
      render(OptimizerPhaseInspector);
      
      // Should call update at least once on mount
      expect(optimizerState.update).toHaveBeenCalled();
    });

    it('updates store when technique selection changes', async () => {
      render(OptimizerPhaseInspector);
      
      // Clear previous calls from mount
      vi.mocked(optimizerState.update).mockClear();
      
      // Change technique selection
      const constantFoldingButton = screen.getByRole('button', { name: 'Constant Folding' });
      await fireEvent.click(constantFoldingButton);
      
      expect(optimizerState.update).toHaveBeenCalled();
    });

    it('subscribes to store changes correctly', () => {
      render(OptimizerPhaseInspector);
      
      expect(optimizerState.subscribe).toHaveBeenCalled();
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
      
      expect(screen.getByRole('button', { name: 'Insert default input' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Optimise Code' })).toBeInTheDocument();
    });
  });
});
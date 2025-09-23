import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import OptimizerArtifactViewer from '../src/lib/components/optimizer/optimizer-artifact-viewer.svelte';
import { optimizerState } from '../src/lib/stores/optimizer';
import { AddToast } from '../src/lib/stores/toast';

// Mock the stores
vi.mock('../src/lib/stores/toast', () => ({
  AddToast: vi.fn()
}));

vi.mock('../src/lib/stores/optimizer', () => ({
  optimizerState: {
    subscribe: vi.fn((fn) => {
      fn({
        optimizedCode: null,
        optimizationError: null
      });
      return { unsubscribe: vi.fn() };
    })
  }
}));

// Mock clipboard API
Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn()
  }
});

describe('OptimizerArtifactViewer Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset clipboard mock
    vi.mocked(navigator.clipboard.writeText).mockResolvedValue(undefined);
  });

  describe('Empty State', () => {
    it('renders empty state when no optimized code is provided', () => {
      render(OptimizerArtifactViewer);
      
      expect(screen.getByText('Optimiser Artefact')).toBeInTheDocument();
      expect(screen.getByText('Optimised code will appear here...')).toBeInTheDocument();
      expect(screen.queryByText('Copy')).not.toBeInTheDocument();
    });

    it('shows empty state with correct icon', () => {
      render(OptimizerArtifactViewer);
      
      const emptyStateDiv = screen.getByText('Optimised code will appear here...').closest('.empty-state');
      expect(emptyStateDiv).toBeInTheDocument();
      
      const svg = emptyStateDiv?.querySelector('svg');
      expect(svg).toBeInTheDocument();
    });
  });

  describe('Optimized Code Display', () => {
    it('displays optimized code when provided', () => {
      const mockOptimizedCode = {
        optimized: ['package main', 'func main() {', '  fmt.Println("Hello")', '}'],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      expect(screen.getByText('Optimised Code')).toBeInTheDocument();
      expect(screen.getByText('Copy')).toBeInTheDocument();
      
      // Check if code lines are displayed - use flexible text matching
      expect(screen.getByText((content, element) => content.includes('package main'))).toBeInTheDocument();
      expect(screen.getByText((content, element) => content.includes('func main()'))).toBeInTheDocument();
    });

    it('displays line numbers correctly', () => {
      const mockOptimizedCode = {
        optimized: ['line1', 'line2', 'line3'],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      // Check line numbers
      expect(screen.getByText('1')).toBeInTheDocument();
      expect(screen.getByText('2')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
    });

    it('handles empty optimized array', () => {
      const mockOptimizedCode = {
        optimized: [],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      expect(screen.getByText('Optimised code will appear here...')).toBeInTheDocument();
      expect(screen.queryByText('Copy')).not.toBeInTheDocument();
    });
  });

  describe('Copy Functionality', () => {
    it('copies optimized code to clipboard successfully', async () => {
      const mockOptimizedCode = {
        optimized: ['package main', 'func main() {}'],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      const copyButton = screen.getByText('Copy');
      await fireEvent.click(copyButton);

      expect(navigator.clipboard.writeText).toHaveBeenCalledWith('package main\nfunc main() {}');
      expect(AddToast).toHaveBeenCalledWith('Code copied! Your optimised code is now in the clipboard', 'success');
    });

    it('handles clipboard copy failure', async () => {
      const mockOptimizedCode = {
        optimized: ['package main', 'func main() {}'],
        language: 'Go'
      };

      // Mock clipboard failure
      vi.mocked(navigator.clipboard.writeText).mockRejectedValue(new Error('Clipboard error'));
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      const copyButton = screen.getByText('Copy');
      await fireEvent.click(copyButton);

      expect(consoleSpy).toHaveBeenCalledWith('Failed to copy text: ', expect.any(Error));
      expect(AddToast).toHaveBeenCalledWith('Copy failed: Unable to copy code to clipboard', 'error');
      
      consoleSpy.mockRestore();
    });

    it('does not show copy button when no optimized code', () => {
      render(OptimizerArtifactViewer);
      
      expect(screen.queryByText('Copy')).not.toBeInTheDocument();
    });

    it('does not show copy button when optimized array is empty', () => {
      const mockOptimizedCode = {
        optimized: [],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      expect(screen.queryByText('Copy')).not.toBeInTheDocument();
    });
  });

  describe('Error Handling', () => {
    it('displays optimization error when provided', () => {
      const mockError = {
        message: 'Optimization failed due to syntax error'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      expect(screen.getByText('Optimisation Failed')).toBeInTheDocument();
      expect(screen.getByText('The optimisation could not be completed with the provided code. Please check your input and try again.')).toBeInTheDocument();
      expect(screen.getByText('Optimization failed due to syntax error')).toBeInTheDocument();
    });

    it('handles specific Go syntax error message', () => {
      const mockError = {
        message: 'source code is not valid go syntax'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      expect(screen.getByText('The provided code contains syntax errors. Please check your Go code syntax and try again.')).toBeInTheDocument();
    });

    it('handles user ID error message', () => {
      const mockError = {
        message: 'User ID not found in session'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      expect(screen.getByText('Authentication error. Please refresh the page and try again.')).toBeInTheDocument();
    });

    it('handles error without message property', () => {
      const mockError = 'Simple error string';

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      expect(screen.getByText('Simple error string')).toBeInTheDocument();
    });

    it('displays error icon when there is an error', () => {
      const mockError = { message: 'Test error' };

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      const errorIcon = document.querySelector('svg');
      expect(errorIcon).toBeInTheDocument();
    });
  });

  describe('Store Integration', () => {
    it('displays optimized code when passed as props', () => {
      const mockOptimizedCode = {
        optimized: ['store code line 1', 'store code line 2'],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      expect(screen.getByText((content) => content.includes('store code line 1'))).toBeInTheDocument();
      expect(screen.getByText((content) => content.includes('store code line 2'))).toBeInTheDocument();
    });

    it('displays error when passed as props', () => {
      const mockError = { message: 'Store error message' };

      render(OptimizerArtifactViewer, {
        props: {
          optimizationError: mockError
        }
      });

      expect(screen.getByText('Optimisation Failed')).toBeInTheDocument();
      expect(screen.getByText('Store error message')).toBeInTheDocument();
    });
  });

  describe('Component Rendering', () => {
    it('renders correctly with mixed states', () => {
      const mockOptimizedCode = {
        optimized: ['test line'],
        language: 'Go'
      };

      render(OptimizerArtifactViewer, {
        props: {
          optimizedCode: mockOptimizedCode
        }
      });

      expect(screen.getByText('Optimiser Artefact')).toBeInTheDocument();
      expect(screen.getByText('test line')).toBeInTheDocument();
    });
  });
});
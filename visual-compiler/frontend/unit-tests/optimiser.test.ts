import { describe, it, expect, beforeEach, vi } from 'vitest';
import { get } from 'svelte/store';
import {
  optimiserState,
  resetOptimiserState,
  updateOptimiserStateFromProject,
  type OptimiserState
} from '../src/lib/stores/optimiser';

describe('Optimizer Store', () => {
  beforeEach(() => {
    resetOptimiserState();
  });

  describe('Initial State', () => {
    it('has correct default state', () => {
      const state = get(optimiserState);
      
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimising: false,
        optimisedCode: null,
        optimisationError: null
      });
    });

    it('initializes with Go as default language', () => {
      const state = get(optimiserState);
      
      expect(state.selectedLanguage).toBe('Go');
    });

    it('initializes with empty arrays and strings', () => {
      const state = get(optimiserState);
      
      expect(state.inputCode).toBe('');
      expect(state.selectedTechniques).toEqual([]);
      expect(state.optimisedCode).toBeNull();
      expect(state.optimisationError).toBeNull();
      expect(state.isOptimising).toBe(false);
    });
  });

  describe('Reset Functionality', () => {
    it('resets state to initial values', () => {
      // Modify state first using the store's update method
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Java',
        selectedTechniques: ['technique1'],
        inputCode: 'some code',
        isOptimising: true,
        optimisedCode: {
          optimized: ['optimized code'],
          language: 'Java',
          techniques: ['technique1']
        },
        optimisationError: 'some error'
      }));
      
      // Reset and verify
      resetOptimiserState();
      const state = get(optimiserState);
      
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimising: false,
        optimisedCode: null,
        optimisationError: null
      });
    });

    it('can be called multiple times safely', () => {
      resetOptimiserState();
      resetOptimiserState();
      resetOptimiserState();
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Go');
      expect(state.selectedTechniques).toEqual([]);
    });
  });

  describe('Direct State Updates', () => {
    it('allows updating selected language', () => {
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Java'
      }));
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Java');
    });

    it('allows updating selected techniques', () => {
      const techniques = ['Dead Code Elimination', 'Loop Unrolling'];
      optimiserState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimiserState);
      expect(state.selectedTechniques).toEqual(techniques);
    });

    it('allows updating input code', () => {
      const testCode = 'func main() { fmt.Println("Hello") }';
      optimiserState.update(state => ({
        ...state,
        inputCode: testCode
      }));
      
      const state = get(optimiserState);
      expect(state.inputCode).toBe(testCode);
    });

    it('allows updating optimization status', () => {
      optimiserState.update(state => ({
        ...state,
        isOptimising: true
      }));
      
      const state = get(optimiserState);
      expect(state.isOptimising).toBe(true);
    });

    it('allows updating optimized code', () => {
      const optimisedCode = {
        optimized: ['optimized function() { return result; }'],
        language: 'Go',
        techniques: ['Dead Code Elimination'],
        performanceGains: {
          executionTime: '50% faster',
          memoryUsage: '30% less',
          codeSize: '20% smaller'
        }
      };

      optimiserState.update(state => ({
        ...state,
        optimisedCode: optimisedCode
      }));
      
      const state = get(optimiserState);
      expect(state.optimisedCode).toEqual(optimisedCode);
    });

    it('allows updating optimization error', () => {
      const error = 'Compilation failed';
      optimiserState.update(state => ({
        ...state,
        optimisationError: error
      }));
      
      const state = get(optimiserState);
      expect(state.optimisationError).toBe(error);
    });
  });

  describe('updateOptimiserStateFromProject', () => {
    it('updates state from valid project data', () => {
      const projectData = {
        optimising: {
          language: 'Java',
          techniques: ['Constant Folding', 'Dead Code Elimination'],
          input_code: 'public class Test { public static void main(String[] args) {} }',
          optimised_code: {
            optimised: ['public class Test { public static void main(String[] args) {} }'],
            language: 'Java',
            techniques: ['Constant Folding'],
            performance_gains: {
              executionTime: '10% faster',
              memoryUsage: '5% less',
              codeSize: '15% smaller'
            }
          },
          error: null
        }
      };

      updateOptimiserStateFromProject(projectData);
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Java');
      expect(state.selectedTechniques).toEqual(['Constant Folding', 'Dead Code Elimination']);
      expect(state.inputCode).toBe('public class Test { public static void main(String[] args) {} }');
      expect(state.optimisedCode).toEqual({
        optimised: ['public class Test { public static void main(String[] args) {} }'],
        language: 'Java',
        techniques: ['Constant Folding'],
        performanceGains: {
          executionTime: '10% faster',
          memoryUsage: '5% less',
          codeSize: '15% smaller'
        }
      });
      expect(state.optimisationError).toBeNull();
      expect(state.isOptimising).toBe(false);
    });

    it('handles project data with error', () => {
      const projectData = {
        optimising: {
          language: 'Python',
          techniques: ['Loop Unrolling'],
          input_code: 'def test(): pass',
          optimised_code: null,
          error: 'Syntax error on line 1'
        }
      };

      updateOptimiserStateFromProject(projectData);
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Python');
      expect(state.selectedTechniques).toEqual(['Loop Unrolling']);
      expect(state.inputCode).toBe('def test(): pass');
      expect(state.optimisedCode).toBeNull();
      expect(state.optimisationError).toBe('Syntax error on line 1');
    });

    it('handles project data with missing optimizing section', () => {
      const projectData = {
        someOtherData: 'value'
      };

      const initialState = get(optimiserState);
      updateOptimiserStateFromProject(projectData);
      
      const state = get(optimiserState);
      expect(state).toEqual(initialState);
    });

    it('handles null project data', () => {
      const initialState = get(optimiserState);
      updateOptimiserStateFromProject(null);
      
      const state = get(optimiserState);
      expect(state).toEqual(initialState);
    });

    it('handles undefined project data', () => {
      const initialState = get(optimiserState);
      updateOptimiserStateFromProject(undefined);
      
      const state = get(optimiserState);
      expect(state).toEqual(initialState);
    });

    it('uses default values for missing fields', () => {
      const projectData = {
        optimising: {
          // Missing most fields
        }
      };

      updateOptimiserStateFromProject(projectData);
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Go'); // default
      expect(state.selectedTechniques).toEqual([]); // default
      expect(state.inputCode).toBe(''); // default
      expect(state.optimisedCode).toBeNull(); // default
      expect(state.optimisationError).toBeNull(); // default
    });

    it('handles partial optimized_code data', () => {
      const projectData = {
        optimising: {
          language: 'Go',
          optimised_code: {
            optimised: ['some code'],
            // Missing language, techniques, performance_gains
          }
        }
      };

      updateOptimiserStateFromProject(projectData);
      
      const state = get(optimiserState);
      expect(state.optimisedCode).toEqual({
        optimised: ['some code'],
        language: 'Go', // default
        techniques: [], // default
        performanceGains: undefined // missing
      });
    });
  });

  describe('Store Reactivity', () => {
    it('notifies subscribers when state changes', () => {
      const subscriber = vi.fn();
      const unsubscribe = optimiserState.subscribe(subscriber);
      
      optimiserState.update(state => ({
        ...state,
        inputCode: 'test'
      }));
      
      expect(subscriber).toHaveBeenCalled();
      unsubscribe();
    });

    it('provides current state to new subscribers', () => {
      optimiserState.update(state => ({
        ...state,
        inputCode: 'existing code'
      }));
      
      const subscriber = vi.fn();
      const unsubscribe = optimiserState.subscribe(subscriber);
      
      expect(subscriber).toHaveBeenCalledWith(
        expect.objectContaining({
          inputCode: 'existing code'
        })
      );
      
      unsubscribe();
    });

    it('handles multiple subscribers', () => {
      const subscriber1 = vi.fn();
      const subscriber2 = vi.fn();
      
      const unsubscribe1 = optimiserState.subscribe(subscriber1);
      const unsubscribe2 = optimiserState.subscribe(subscriber2);
      
      optimiserState.update(state => ({
        ...state,
        inputCode: 'test'
      }));
      
      expect(subscriber1).toHaveBeenCalled();
      expect(subscriber2).toHaveBeenCalled();
      
      unsubscribe1();
      unsubscribe2();
    });
  });

  describe('Language Support', () => {
    it('supports Java language', () => {
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Java'
      }));
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Java');
    });

    it('supports Python language', () => {
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Python'
      }));
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Python');
    });

    it('supports Go language', () => {
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Go'
      }));
      
      const state = get(optimiserState);
      expect(state.selectedLanguage).toBe('Go');
    });
  });

  describe('Complex Workflow Scenarios', () => {
    it('handles complete optimization workflow', () => {
      // Initial setup
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Go',
        selectedTechniques: ['Dead Code Elimination', 'Constant Folding'],
        inputCode: 'func main() { var x = 5; fmt.Println(x) }'
      }));

      // Start optimization
      optimiserState.update(state => ({
        ...state,
        isOptimising: true
      }));

      let state = get(optimiserState);
      expect(state.isOptimising).toBe(true);

      // Complete optimization successfully
      optimiserState.update(state => ({
        ...state,
        isOptimising: false,
        optimisedCode: {
          optimized: ['func main() { fmt.Println(5) }'],
          language: 'Go',
          techniques: ['Dead Code Elimination', 'Constant Folding'],
          performanceGains: {
            executionTime: '15% faster',
            memoryUsage: '10% less',
            codeSize: '5% smaller'
          }
        },
        optimisationError: null
      }));

      state = get(optimiserState);
      expect(state.isOptimising).toBe(false);
      expect(state.optimisedCode).toBeDefined();
      expect(state.optimisationError).toBeNull();
    });

    it('handles optimization error workflow', () => {
      // Setup and start optimization
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Java',
        inputCode: 'invalid syntax here',
        isOptimising: true
      }));

      // Optimization fails
      optimiserState.update(state => ({
        ...state,
        isOptimising: false,
        optimisedCode: null,
        optimisationError: 'Compilation failed: syntax error'
      }));

      const state = get(optimiserState);
      expect(state.isOptimising).toBe(false);
      expect(state.optimisedCode).toBeNull();
      expect(state.optimisationError).toBe('Compilation failed: syntax error');
    });

    it('handles state reset during optimization', () => {
      // Setup optimization in progress
      optimiserState.update(state => ({
        ...state,
        selectedLanguage: 'Python',
        inputCode: 'def test(): pass',
        isOptimising: true,
        selectedTechniques: ['Loop Unrolling']
      }));

      // Reset state
      resetOptimiserState();

      const state = get(optimiserState);
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimising: false,
        optimisedCode: null,
        optimisationError: null
      });
    });
  });

  describe('Performance Gains Handling', () => {
    it('handles complete performance gains data', () => {
      const performanceGains = {
        executionTime: '25% faster',
        memoryUsage: '15% less',
        codeSize: '10% smaller'
      };

      optimiserState.update(state => ({
        ...state,
        optimisedCode: {
          optimized: ['optimized code'],
          language: 'Go',
          techniques: ['Optimization'],
          performanceGains: performanceGains
        }
      }));

      const state = get(optimiserState);
      expect(state.optimisedCode?.performanceGains).toEqual(performanceGains);
    });

    it('handles missing performance gains', () => {
      optimiserState.update(state => ({
        ...state,
        optimisedCode: {
          optimized: ['optimized code'],
          language: 'Go',
          techniques: ['Optimization']
          // No performanceGains property
        }
      }));

      const state = get(optimiserState);
      expect(state.optimisedCode?.performanceGains).toBeUndefined();
    });
  });

  describe('Type Safety and Edge Cases', () => {
    it('maintains correct TypeScript types', () => {
      const state: optimiserState = get(optimiserState);
      
      expect(typeof state.selectedLanguage).toBe('string');
      expect(Array.isArray(state.selectedTechniques)).toBe(true);
      expect(typeof state.inputCode).toBe('string');
      expect(typeof state.isOptimising).toBe('boolean');
      expect(state.optimisedCode === null || typeof state.optimisedCode === 'object').toBe(true);
    });

    it('handles empty techniques array', () => {
      optimiserState.update(state => ({
        ...state,
        selectedTechniques: []
      }));
      
      const state = get(optimiserState);
      expect(state.selectedTechniques).toEqual([]);
    });

    it('handles multiple techniques', () => {
      const techniques = [
        'Dead Code Elimination',
        'Constant Folding',
        'Loop Unrolling',
        'Function Inlining',
        'Common Sub-expression Elimination'
      ];
      
      optimiserState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimiserState);
      expect(state.selectedTechniques).toEqual(techniques);
    });

    it('preserves technique order', () => {
      const techniques = ['C', 'A', 'B'];
      optimiserState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimiserState);
      expect(state.selectedTechniques).toEqual(['C', 'A', 'B']);
    });

    it('handles multiline input code', () => {
      const multilineCode = `func main() {
        fmt.Println("Hello")
        fmt.Println("World")
      }`;
      
      optimiserState.update(state => ({
        ...state,
        inputCode: multilineCode
      }));
      
      const state = get(optimiserState);
      expect(state.inputCode).toBe(multilineCode);
    });

    it('handles special characters in input code', () => {
      const codeWithSymbols = 'func test() { return "special chars: !@#$%^&*()_+" }';
      optimiserState.update(state => ({
        ...state,
        inputCode: codeWithSymbols
      }));
      
      const state = get(optimiserState);
      expect(state.inputCode).toBe(codeWithSymbols);
    });
  });
});

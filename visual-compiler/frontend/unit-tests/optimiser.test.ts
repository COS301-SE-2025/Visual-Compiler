import { describe, it, expect, beforeEach, vi } from 'vitest';
import { get } from 'svelte/store';
import {
  optimizerState,
  resetOptimizerState,
  updateOptimizerStateFromProject,
  type OptimizerState
} from '../src/lib/stores/optimizer';

describe('Optimizer Store', () => {
  beforeEach(() => {
    resetOptimizerState();
  });

  describe('Initial State', () => {
    it('has correct default state', () => {
      const state = get(optimizerState);
      
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimizing: false,
        optimizedCode: null,
        optimizationError: null
      });
    });

    it('initializes with Go as default language', () => {
      const state = get(optimizerState);
      
      expect(state.selectedLanguage).toBe('Go');
    });

    it('initializes with empty arrays and strings', () => {
      const state = get(optimizerState);
      
      expect(state.inputCode).toBe('');
      expect(state.selectedTechniques).toEqual([]);
      expect(state.optimizedCode).toBeNull();
      expect(state.optimizationError).toBeNull();
      expect(state.isOptimizing).toBe(false);
    });
  });

  describe('Reset Functionality', () => {
    it('resets state to initial values', () => {
      // Modify state first using the store's update method
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Java',
        selectedTechniques: ['technique1'],
        inputCode: 'some code',
        isOptimizing: true,
        optimizedCode: {
          optimized: ['optimized code'],
          language: 'Java',
          techniques: ['technique1']
        },
        optimizationError: 'some error'
      }));
      
      // Reset and verify
      resetOptimizerState();
      const state = get(optimizerState);
      
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimizing: false,
        optimizedCode: null,
        optimizationError: null
      });
    });

    it('can be called multiple times safely', () => {
      resetOptimizerState();
      resetOptimizerState();
      resetOptimizerState();
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Go');
      expect(state.selectedTechniques).toEqual([]);
    });
  });

  describe('Direct State Updates', () => {
    it('allows updating selected language', () => {
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Java'
      }));
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Java');
    });

    it('allows updating selected techniques', () => {
      const techniques = ['Dead Code Elimination', 'Loop Unrolling'];
      optimizerState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimizerState);
      expect(state.selectedTechniques).toEqual(techniques);
    });

    it('allows updating input code', () => {
      const testCode = 'func main() { fmt.Println("Hello") }';
      optimizerState.update(state => ({
        ...state,
        inputCode: testCode
      }));
      
      const state = get(optimizerState);
      expect(state.inputCode).toBe(testCode);
    });

    it('allows updating optimization status', () => {
      optimizerState.update(state => ({
        ...state,
        isOptimizing: true
      }));
      
      const state = get(optimizerState);
      expect(state.isOptimizing).toBe(true);
    });

    it('allows updating optimized code', () => {
      const optimizedCode = {
        optimized: ['optimized function() { return result; }'],
        language: 'Go',
        techniques: ['Dead Code Elimination'],
        performanceGains: {
          executionTime: '50% faster',
          memoryUsage: '30% less',
          codeSize: '20% smaller'
        }
      };

      optimizerState.update(state => ({
        ...state,
        optimizedCode: optimizedCode
      }));
      
      const state = get(optimizerState);
      expect(state.optimizedCode).toEqual(optimizedCode);
    });

    it('allows updating optimization error', () => {
      const error = 'Compilation failed';
      optimizerState.update(state => ({
        ...state,
        optimizationError: error
      }));
      
      const state = get(optimizerState);
      expect(state.optimizationError).toBe(error);
    });
  });

  describe('updateOptimizerStateFromProject', () => {
    it('updates state from valid project data', () => {
      const projectData = {
        optimizing: {
          language: 'Java',
          techniques: ['Constant Folding', 'Dead Code Elimination'],
          input_code: 'public class Test { public static void main(String[] args) {} }',
          optimized_code: {
            optimized: ['public class Test { public static void main(String[] args) {} }'],
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

      updateOptimizerStateFromProject(projectData);
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Java');
      expect(state.selectedTechniques).toEqual(['Constant Folding', 'Dead Code Elimination']);
      expect(state.inputCode).toBe('public class Test { public static void main(String[] args) {} }');
      expect(state.optimizedCode).toEqual({
        optimized: ['public class Test { public static void main(String[] args) {} }'],
        language: 'Java',
        techniques: ['Constant Folding'],
        performanceGains: {
          executionTime: '10% faster',
          memoryUsage: '5% less',
          codeSize: '15% smaller'
        }
      });
      expect(state.optimizationError).toBeNull();
      expect(state.isOptimizing).toBe(false);
    });

    it('handles project data with error', () => {
      const projectData = {
        optimizing: {
          language: 'Python',
          techniques: ['Loop Unrolling'],
          input_code: 'def test(): pass',
          optimized_code: null,
          error: 'Syntax error on line 1'
        }
      };

      updateOptimizerStateFromProject(projectData);
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Python');
      expect(state.selectedTechniques).toEqual(['Loop Unrolling']);
      expect(state.inputCode).toBe('def test(): pass');
      expect(state.optimizedCode).toBeNull();
      expect(state.optimizationError).toBe('Syntax error on line 1');
    });

    it('handles project data with missing optimizing section', () => {
      const projectData = {
        someOtherData: 'value'
      };

      const initialState = get(optimizerState);
      updateOptimizerStateFromProject(projectData);
      
      const state = get(optimizerState);
      expect(state).toEqual(initialState);
    });

    it('handles null project data', () => {
      const initialState = get(optimizerState);
      updateOptimizerStateFromProject(null);
      
      const state = get(optimizerState);
      expect(state).toEqual(initialState);
    });

    it('handles undefined project data', () => {
      const initialState = get(optimizerState);
      updateOptimizerStateFromProject(undefined);
      
      const state = get(optimizerState);
      expect(state).toEqual(initialState);
    });

    it('uses default values for missing fields', () => {
      const projectData = {
        optimizing: {
          // Missing most fields
        }
      };

      updateOptimizerStateFromProject(projectData);
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Go'); // default
      expect(state.selectedTechniques).toEqual([]); // default
      expect(state.inputCode).toBe(''); // default
      expect(state.optimizedCode).toBeNull(); // default
      expect(state.optimizationError).toBeNull(); // default
    });

    it('handles partial optimized_code data', () => {
      const projectData = {
        optimizing: {
          language: 'Go',
          optimized_code: {
            optimized: ['some code'],
            // Missing language, techniques, performance_gains
          }
        }
      };

      updateOptimizerStateFromProject(projectData);
      
      const state = get(optimizerState);
      expect(state.optimizedCode).toEqual({
        optimized: ['some code'],
        language: 'Go', // default
        techniques: [], // default
        performanceGains: undefined // missing
      });
    });
  });

  describe('Store Reactivity', () => {
    it('notifies subscribers when state changes', () => {
      const subscriber = vi.fn();
      const unsubscribe = optimizerState.subscribe(subscriber);
      
      optimizerState.update(state => ({
        ...state,
        inputCode: 'test'
      }));
      
      expect(subscriber).toHaveBeenCalled();
      unsubscribe();
    });

    it('provides current state to new subscribers', () => {
      optimizerState.update(state => ({
        ...state,
        inputCode: 'existing code'
      }));
      
      const subscriber = vi.fn();
      const unsubscribe = optimizerState.subscribe(subscriber);
      
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
      
      const unsubscribe1 = optimizerState.subscribe(subscriber1);
      const unsubscribe2 = optimizerState.subscribe(subscriber2);
      
      optimizerState.update(state => ({
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
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Java'
      }));
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Java');
    });

    it('supports Python language', () => {
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Python'
      }));
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Python');
    });

    it('supports Go language', () => {
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Go'
      }));
      
      const state = get(optimizerState);
      expect(state.selectedLanguage).toBe('Go');
    });
  });

  describe('Complex Workflow Scenarios', () => {
    it('handles complete optimization workflow', () => {
      // Initial setup
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Go',
        selectedTechniques: ['Dead Code Elimination', 'Constant Folding'],
        inputCode: 'func main() { var x = 5; fmt.Println(x) }'
      }));

      // Start optimization
      optimizerState.update(state => ({
        ...state,
        isOptimizing: true
      }));

      let state = get(optimizerState);
      expect(state.isOptimizing).toBe(true);

      // Complete optimization successfully
      optimizerState.update(state => ({
        ...state,
        isOptimizing: false,
        optimizedCode: {
          optimized: ['func main() { fmt.Println(5) }'],
          language: 'Go',
          techniques: ['Dead Code Elimination', 'Constant Folding'],
          performanceGains: {
            executionTime: '15% faster',
            memoryUsage: '10% less',
            codeSize: '5% smaller'
          }
        },
        optimizationError: null
      }));

      state = get(optimizerState);
      expect(state.isOptimizing).toBe(false);
      expect(state.optimizedCode).toBeDefined();
      expect(state.optimizationError).toBeNull();
    });

    it('handles optimization error workflow', () => {
      // Setup and start optimization
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Java',
        inputCode: 'invalid syntax here',
        isOptimizing: true
      }));

      // Optimization fails
      optimizerState.update(state => ({
        ...state,
        isOptimizing: false,
        optimizedCode: null,
        optimizationError: 'Compilation failed: syntax error'
      }));

      const state = get(optimizerState);
      expect(state.isOptimizing).toBe(false);
      expect(state.optimizedCode).toBeNull();
      expect(state.optimizationError).toBe('Compilation failed: syntax error');
    });

    it('handles state reset during optimization', () => {
      // Setup optimization in progress
      optimizerState.update(state => ({
        ...state,
        selectedLanguage: 'Python',
        inputCode: 'def test(): pass',
        isOptimizing: true,
        selectedTechniques: ['Loop Unrolling']
      }));

      // Reset state
      resetOptimizerState();

      const state = get(optimizerState);
      expect(state).toEqual({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimizing: false,
        optimizedCode: null,
        optimizationError: null
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

      optimizerState.update(state => ({
        ...state,
        optimizedCode: {
          optimized: ['optimized code'],
          language: 'Go',
          techniques: ['Optimization'],
          performanceGains: performanceGains
        }
      }));

      const state = get(optimizerState);
      expect(state.optimizedCode?.performanceGains).toEqual(performanceGains);
    });

    it('handles missing performance gains', () => {
      optimizerState.update(state => ({
        ...state,
        optimizedCode: {
          optimized: ['optimized code'],
          language: 'Go',
          techniques: ['Optimization']
          // No performanceGains property
        }
      }));

      const state = get(optimizerState);
      expect(state.optimizedCode?.performanceGains).toBeUndefined();
    });
  });

  describe('Type Safety and Edge Cases', () => {
    it('maintains correct TypeScript types', () => {
      const state: OptimizerState = get(optimizerState);
      
      expect(typeof state.selectedLanguage).toBe('string');
      expect(Array.isArray(state.selectedTechniques)).toBe(true);
      expect(typeof state.inputCode).toBe('string');
      expect(typeof state.isOptimizing).toBe('boolean');
      expect(state.optimizedCode === null || typeof state.optimizedCode === 'object').toBe(true);
    });

    it('handles empty techniques array', () => {
      optimizerState.update(state => ({
        ...state,
        selectedTechniques: []
      }));
      
      const state = get(optimizerState);
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
      
      optimizerState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimizerState);
      expect(state.selectedTechniques).toEqual(techniques);
    });

    it('preserves technique order', () => {
      const techniques = ['C', 'A', 'B'];
      optimizerState.update(state => ({
        ...state,
        selectedTechniques: techniques
      }));
      
      const state = get(optimizerState);
      expect(state.selectedTechniques).toEqual(['C', 'A', 'B']);
    });

    it('handles multiline input code', () => {
      const multilineCode = `func main() {
        fmt.Println("Hello")
        fmt.Println("World")
      }`;
      
      optimizerState.update(state => ({
        ...state,
        inputCode: multilineCode
      }));
      
      const state = get(optimizerState);
      expect(state.inputCode).toBe(multilineCode);
    });

    it('handles special characters in input code', () => {
      const codeWithSymbols = 'func test() { return "special chars: !@#$%^&*()_+" }';
      optimizerState.update(state => ({
        ...state,
        inputCode: codeWithSymbols
      }));
      
      const state = get(optimizerState);
      expect(state.inputCode).toBe(codeWithSymbols);
    });
  });
});
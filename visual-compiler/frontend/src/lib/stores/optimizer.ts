import { writable } from 'svelte/store';

export interface OptimizedCode {
    original: string[];
    optimized: string[];
    language: 'Java' | 'Python' | 'Go';
    techniques: string[];
    performanceGains?: {
        executionTime: string;
        memoryUsage: string;
        codeSize: string;
    };
}

export interface OptimizerState {
    selectedLanguage: 'Java' | 'Python' | 'Go';
    selectedTechniques: string[];
    inputCode: string;
    optimizedCode: OptimizedCode | null;
    optimizationError: any;
    isOptimizing: boolean;
}

export const optimizerState = writable<OptimizerState>({
    selectedLanguage: 'Java',
    selectedTechniques: [],
    inputCode: '',
    optimizedCode: null,
    optimizationError: null,
    isOptimizing: false
});

export function resetOptimizerState() {
    optimizerState.set({
        selectedLanguage: 'Java',
        selectedTechniques: [],
        inputCode: '',
        optimizedCode: null,
        optimizationError: null,
        isOptimizing: false
    });
}

export function updateOptimizerStateFromProject(projectData: any) {
    if (!projectData?.optimizing) return;

    const optimizingData = projectData.optimizing;

    optimizerState.update(state => ({
        ...state,
        selectedLanguage: optimizingData.language || 'Java',
        selectedTechniques: optimizingData.techniques || [],
        inputCode: optimizingData.input_code || '',
        optimizedCode: optimizingData.optimized_code ? {
            original: optimizingData.optimized_code.original || [],
            optimized: optimizingData.optimized_code.optimized || [],
            language: optimizingData.optimized_code.language || 'Java',
            techniques: optimizingData.optimized_code.techniques || [],
            performanceGains: optimizingData.optimized_code.performance_gains
        } : null,
        optimizationError: optimizingData.error || null,
        isOptimizing: false
    }));
}

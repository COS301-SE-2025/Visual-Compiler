import { writable } from 'svelte/store';

export interface OptimizerState {
    selectedLanguage: 'Java' | 'Python' | 'Go';
    selectedTechniques: string[];
    inputCode: string;
    isOptimizing: boolean;
    optimizedCode: {
        optimized: string[];
        language: string;
        techniques?: string[];
        performanceGains?: {
            executionTime: string;
            memoryUsage: string;
            codeSize: string;
        };
    } | null;
    optimizationError: any;
}

export const optimizerState = writable<OptimizerState>({
    selectedLanguage: 'Go',
    selectedTechniques: [],
    inputCode: '',
    isOptimizing: false,
    optimizedCode: null,
    optimizationError: null
});

export function resetOptimizerState() {
    optimizerState.set({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimizing: false,
        optimizedCode: null,
        optimizationError: null
    });
}

export function updateOptimizerStateFromProject(projectData: any) {
    if (!projectData?.optimizing) return;

    const optimizingData = projectData.optimizing;

    optimizerState.update(state => ({
        ...state,
        selectedLanguage: optimizingData.language || 'Go',
        selectedTechniques: optimizingData.techniques || [],
        inputCode: optimizingData.input_code || '',
        optimizedCode: optimizingData.optimized_code ? {
            optimized: optimizingData.optimized_code.optimized || [],
            language: optimizingData.optimized_code.language || 'Go',
            techniques: optimizingData.optimized_code.techniques || [],
            performanceGains: optimizingData.optimized_code.performance_gains
        } : null,
        optimizationError: optimizingData.error || null,
        isOptimizing: false
    }));
}

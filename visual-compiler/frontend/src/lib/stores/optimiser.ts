import { writable } from 'svelte/store';

export interface OptimiserState {
    selectedLanguage: 'Java' | 'Python' | 'Go';
    selectedTechniques: string[];
    inputCode: string;
    isOptimising: boolean;
    optimisedCode: {
        optimised: string[];
        language: string;
        techniques?: string[];
        performanceGains?: {
            executionTime: string;
            memoryUsage: string;
            codeSize: string;
        };
    } | null;
    optimisationError: any;
}

export const optimiserState = writable<OptimiserState>({
    selectedLanguage: 'Go',
    selectedTechniques: [],
    inputCode: '',
    isOptimising: false,
    optimisedCode: null,
    optimisationError: null
});

export function resetOptimiserState() {
    optimiserState.set({
        selectedLanguage: 'Go',
        selectedTechniques: [],
        inputCode: '',
        isOptimising: false,
        optimisedCode: null,
        optimisationError: null
    });
}

export function updateOptimiserStateFromProject(projectData: any) {
    if (!projectData?.optimising) return;

    const optimisingData = projectData.optimising;

    optimiserState.update(state => ({
        ...state,
        selectedLanguage: optimisingData.language || 'Go',
        selectedTechniques: optimisingData.techniques || [],
        inputCode: optimisingData.input_code || '',
        optimisedCode: optimisingData.optimised_code ? {
            optimised: optimisingData.optimised_code.optimised || [],
            language: optimisingData.optimised_code.language || 'Go',
            techniques: optimisingData.optimised_code.techniques || [],
            performanceGains: optimisingData.optimised_code.performance_gains
        } : null,
        optimisationError: optimisingData.error || null,
        isOptimising: false
    }));
}

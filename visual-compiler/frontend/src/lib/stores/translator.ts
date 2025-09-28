import { writable } from 'svelte/store';
import { get } from 'svelte/store';
import { projectName } from './project';

interface TranslationRule {
    tokenSequence: string;
    lines: string[];
}

interface TranslatorState {
    rules: TranslationRule[];
    show_default_rules: boolean;
    isSubmitted: boolean;
    translationSuccessful: boolean;
    hasUnsavedChanges: boolean;
    // ADD: Artifact storage
    translatedCode: string[];
    hasTranslatedCode: boolean;
    translationError: string | null;
}

// Define initial state as a constant
const INITIAL_TRANSLATOR_STATE: TranslatorState = {
    rules: [{ tokenSequence: '', lines: [''] }],
    isSubmitted: false,
    translatedCode: [],
    hasTranslatedCode: false,
    translationSuccessful: false,
    translationError: null,
    hasUnsavedChanges: false
};

// Create store with initial state
export const translatorState = writable<TranslatorState>(INITIAL_TRANSLATOR_STATE);

// Helper functions
export const updateTranslatorInputs = (data: Partial<TranslatorState>) => {
    translatorState.update(state => ({
        ...state,
        ...data,
        hasUnsavedChanges: true
    }));
};

export const markTranslatorSubmitted = () => {
    translatorState.update(state => ({
        ...state,
        isSubmitted: true,
        hasUnsavedChanges: false
    }));
};

// FIX: Proper reset function
export const resetTranslatorState = () => {
    translatorState.set(JSON.parse(JSON.stringify(INITIAL_TRANSLATOR_STATE)));
    console.log('=== TRANSLATOR STATE RESET COMPLETE ===');
};

// ADD: Function to update artifacts
export const updateTranslatorArtifacts = (translatedCode: string[], error: string | null = null) => {
    translatorState.update(state => ({
        ...state,
        translatedCode: [...translatedCode],
        hasTranslatedCode: translatedCode.length > 0,
        translationSuccessful: translatedCode.length > 0 && !error,
        translationError: error
    }));
};

// Project loading function
export const updateTranslatorStateFromProject = (projectData: any) => {
    console.log('Updating translator state from project data:', projectData?.translating);
    if (!projectData?.translating) return;

    const currentProject = get(projectName);
    if (!currentProject) return;

    const updates: Partial<TranslatorState> = {};

    // Handle translation rules
    if (projectData.translating.translating_rules && Array.isArray(projectData.translating.translating_rules)) {
        updates.rules = projectData.translating.translating_rules.map((rule: any) => ({
            tokenSequence: Array.isArray(rule.sequence) ? rule.sequence.join(', ') : (rule.tokenSequence || ''),
            lines: Array.isArray(rule.translation) ? rule.translation : (Array.isArray(rule.lines) ? rule.lines : [''])
        }));
    }

    // FIX: Handle artifacts - translated code
    if (projectData.translating.code && Array.isArray(projectData.translating.code)) {
        updates.translatedCode = [...projectData.translating.code];
        updates.hasTranslatedCode = true;
        updates.translationSuccessful = true;
        updates.translationError = null;
        
        console.log('Loaded translated code with', updates.translatedCode.length, 'lines');
    }

    // Handle translation errors
    if (projectData.translating.error) {
        updates.translationError = projectData.translating.error;
        updates.translationSuccessful = false;
    }

    // Mark as submitted if there's data
    if (projectData.translating && Object.keys(projectData.translating).length > 0) {
        updates.isSubmitted = true;
        updates.hasUnsavedChanges = false;
    }

    if (Object.keys(updates).length > 0) {
        console.log('Updating translator state with:', updates);
        translatorState.update(state => ({
            ...state,
            ...updates
        }));
    }
};
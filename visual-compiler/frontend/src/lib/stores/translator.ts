import { writable } from 'svelte/store';

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

const initialState: TranslatorState = {
    rules: [{ tokenSequence: '', lines: [''] }],
    show_default_rules: false,
    isSubmitted: false,
    translationSuccessful: false,
    hasUnsavedChanges: false,
    // ADD: Artifact fields
    translatedCode: [],
    hasTranslatedCode: false,
    translationError: null
};

export const translatorState = writable<TranslatorState>(initialState);

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

export const resetTranslatorState = () => {
    translatorState.set(initialState);
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
    if (!projectData?.translating) return;

    const updates: Partial<TranslatorState> = {};

    // Handle translation rules
    if (projectData.translating.translating_rules && Array.isArray(projectData.translating.translating_rules)) {
        updates.rules = projectData.translating.translating_rules.map((rule: any) => ({
            tokenSequence: Array.isArray(rule.sequence) ? rule.sequence.join(', ') : (rule.tokenSequence || ''),
            lines: Array.isArray(rule.translation) ? rule.translation : (Array.isArray(rule.lines) ? rule.lines : [''])
        }));
    }

    // Handle artifacts - translated code
    if (projectData.translating.code && Array.isArray(projectData.translating.code)) {
        updates.translatedCode = [...projectData.translating.code];
        updates.hasTranslatedCode = true;
        updates.translationSuccessful = true;
        updates.translationError = null;
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
        translatorState.update(state => ({
            ...state,
            ...updates
        }));
    }
};
import { writable } from 'svelte/store';
import { get } from 'svelte/store';
import { projectName } from './project';

interface LexerInputRow {
    type: string;
    regex: string;
    error: string;
}

interface LexerState {
    userInputRows: LexerInputRow[];
    automataInputs: {
        states: string;
        startState: string;
        acceptedStates: string;
        transitions: string;
    };
    selectedType: 'AUTOMATA' | 'REGEX' | null;
    showDefault: boolean;
    isSubmitted: boolean;
    hasUnsavedChanges: boolean;
    // ADD: Artifact storage
    tokens: any[];
    hasTokens: boolean;
    sourceCode: string;
}

// Define initial state as a constant
const INITIAL_LEXER_STATE: LexerState = {
    selectedType: 'REGEX',
    userInputRows: [{ type: '', regex: '', error: '' }],
    automataInputs: {
        states: '',
        startState: '',
        acceptedStates: '',
        transitions: ''
    },
    showDefault: false,  
    tokens: [],
    hasTokens: false,
    sourceCode: '',
    isSubmitted: false,
    hasUnsavedChanges: false
};

// Create store with initial state
export const lexerState = writable<LexerState>(INITIAL_LEXER_STATE);

// Helper functions
export const updateLexerInputs = (inputs: LexerInputRow[]) => {
    lexerState.update(state => ({
        ...state,
        userInputRows: [...inputs],
        hasUnsavedChanges: true
    }));
};

export const updateAutomataInputs = (inputs: LexerState['automataInputs']) => {
    lexerState.update(state => ({
        ...state,
        automataInputs: { ...inputs },
        hasUnsavedChanges: true
    }));
};

export const markLexerSubmitted = () => {
    lexerState.update(state => ({
        ...state,
        isSubmitted: true,
        hasUnsavedChanges: false
    }));
};

// ADD: Export all reset functions
export const resetLexerState = () => {
    lexerState.set(JSON.parse(JSON.stringify(INITIAL_LEXER_STATE))); // Deep copy
    console.log('Lexer state reset to:', INITIAL_LEXER_STATE);
};

export const updateLexerArtifacts = (tokens: any[], sourceCode: string = '') => {
    lexerState.update(state => ({
        ...state,
        tokens: [...tokens],
        hasTokens: tokens.length > 0,
        sourceCode
    }));
};

// UPDATE: Project loading function to include artifacts
export const updateLexerStateFromProject = (projectData: any) => {
    console.log('Updating lexer state from project data:', projectData?.lexing);
    if (!projectData?.lexing) return;

    const updates: Partial<LexerState> = {};
    
    // Handle automata data
    if (projectData.lexing.automata) {
        const automataData = projectData.lexing.automata;
        updates.automataInputs = {
            states: automataData.states?.join(', ') || '',
            startState: automataData.start_state || '',
            acceptedStates: automataData.accepted_states?.join(', ') || '',
            transitions: Array.isArray(automataData.transitions) 
                ? automataData.transitions.map(t => `${t.from || ''}-${t.symbol || ''}-${t.to || ''}`).join(', ')
                : automataData.transitions || ''
        };
        updates.selectedType = 'AUTOMATA';
    }

    // Handle regex rules
    if (projectData.lexing.regex_rules && Array.isArray(projectData.lexing.regex_rules)) {
        updates.userInputRows = projectData.lexing.regex_rules.map((rule: any) => ({
            type: rule.token_type || rule.type || '',
            regex: rule.regex_pattern || rule.regex || '',
            error: ''
        }));
        updates.selectedType = 'REGEX';
    }

    // Handle artifacts - tokens
    if (projectData.lexing.tokens && Array.isArray(projectData.lexing.tokens)) {
        updates.tokens = projectData.lexing.tokens.map((token: any) => ({
            value: token.Value || token.value || '',
            type: token.Type || token.type || 'UNKNOWN'
        }));
        updates.hasTokens = true;
    }

    // Handle source code
    if (projectData.lexing.code) {
        updates.sourceCode = projectData.lexing.code;
    }

    // Mark as submitted if there's data
    if (projectData.lexing && Object.keys(projectData.lexing).length > 0) {
        updates.isSubmitted = true;
        updates.hasUnsavedChanges = false;
    }

    if (Object.keys(updates).length > 0) {
        console.log('Updating lexer state with:', updates);
        lexerState.update(state => ({
            ...state,
            ...updates
        }));
    }
};

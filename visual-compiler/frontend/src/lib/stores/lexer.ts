import { writable } from 'svelte/store';

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

const initialState: LexerState = {
    userInputRows: [{ type: '', regex: '', error: '' }],
    automataInputs: {
        states: '',
        startState: '',
        acceptedStates: '',
        transitions: ''
    },
    selectedType: 'REGEX',
    showDefault: false,
    isSubmitted: false,
    hasUnsavedChanges: false,
    // ADD: Artifact fields
    tokens: [],
    hasTokens: false,
    sourceCode: ''
};

export const lexerState = writable<LexerState>(initialState);

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

export const resetLexerState = () => {
    lexerState.set(initialState);
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
    if (!projectData?.lexing) return;

    const updates: Partial<LexerState> = {};

    // Handle DFA automata data - convert to correct format
    if (projectData.lexing.dfa) {
        const dfa = projectData.lexing.dfa;
        
        // Format states
        const statesStr = Array.isArray(dfa.states) ? dfa.states.join(', ') : '';
        
        // Format accepting states
        const acceptingStr = Array.isArray(dfa.accepting) 
            ? dfa.accepting.map(a => `${a.state}->${a.type}`).join(', ')
            : '';
        
        // Format transitions
        const transitionsStr = Array.isArray(dfa.transitions)
            ? dfa.transitions.map(t => `${t.from},${t.label}->${t.to}`).join('\n')
            : '';

        updates.automataInputs = {
            states: statesStr,
            startState: dfa.start || '',
            acceptedStates: acceptingStr,
            transitions: transitionsStr
        };
        updates.selectedType = 'AUTOMATA';
    }

    // Handle regex rules from project data  
    if (projectData.lexing.rules && Array.isArray(projectData.lexing.rules)) {
        updates.userInputRows = projectData.lexing.rules.map((rule: any) => ({
            type: rule.type || '',
            regex: rule.regex || '',
            error: ''
        }));
        if (!updates.selectedType) {
            updates.selectedType = 'REGEX';
        }
    }

    // Handle artifacts - tokens
    if (projectData.lexing.tokens && Array.isArray(projectData.lexing.tokens)) {
        updates.tokens = [...projectData.lexing.tokens];
        updates.hasTokens = true;
    }

    // Handle source code
    if (projectData.lexing.code) {
        updates.sourceCode = projectData.lexing.code;
    }

    // Mark as submitted if there's data
    if (projectData.lexing && (projectData.lexing.rules || projectData.lexing.dfa)) {
        updates.isSubmitted = true;
        updates.hasUnsavedChanges = false;
    }

    if (Object.keys(updates).length > 0) {
        lexerState.update(state => ({
            ...state,
            ...updates
        }));
    }
};

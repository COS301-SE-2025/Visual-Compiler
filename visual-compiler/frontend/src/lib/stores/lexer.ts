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
    hasUnsavedChanges: false
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


export const updateLexerStateFromProject = (projectData: any) => {
    if (!projectData) return;

    const updates: Partial<LexerState> = {};

    // Handle regex rules from project data
    if (projectData.lexer_data?.regex_rules) {
        updates.userInputRows = projectData.lexer_data.regex_rules.map((rule: any) => ({
            type: rule.Type || rule.type || '',
            regex: rule.Regex || rule.regex || '',
            error: ''
        }));
    }

    // Handle automata data from project data
    if (projectData.lexer_data?.automata) {
        updates.automataInputs = {
            states: projectData.lexer_data.automata.states || '',
            startState: projectData.lexer_data.automata.start_state || '',
            acceptedStates: projectData.lexer_data.automata.accepted_states || '',
            transitions: projectData.lexer_data.automata.transitions || ''
        };
    }

    // Set selected type based on what data is available
    if (projectData.lexer_data?.regex_rules && projectData.lexer_data.regex_rules.length > 0) {
        updates.selectedType = 'REGEX';
    } else if (projectData.lexer_data?.automata) {
        updates.selectedType = 'AUTOMATA';
    }

    // Mark as submitted if there's data
    if (projectData.lexer_data) {
        updates.isSubmitted = true;
        updates.hasUnsavedChanges = false;
    }

    lexerState.update(state => ({
        ...state,
        ...updates
    }));
};

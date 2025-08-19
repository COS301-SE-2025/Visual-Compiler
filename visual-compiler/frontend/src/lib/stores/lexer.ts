import { writable } from 'svelte/store';

export interface DFAState {
    accepting: Array<{
        state: string;
        type: string;  
    }>;
    start: string;
    states: string[];
    transitions: Array<{
        from: string;
        to: string;
        label: string;  
    }>;
}

export interface Rule {
    type: string;
    regex: string;
}

export interface LexerState {
    mode: 'AUTOMATA' | 'REGEX' | null;
    dfa: DFAState | null;
    tokens: Array<{
        type: string;
        value: string;
    }> | null;
    tokens_unidentified: string[] | null;
    inputs: Array<{
        type: string;
        regex: string;
        error?: string;
    }>;
    states: string;
    startState: string;
    acceptedStates: string;
    transitions: string;
}

export const lexerState = writable<LexerState>({
    mode: null,
    dfa: null,
    tokens: null,
    tokens_unidentified: null,
    inputs: [],
    states: '',
    startState: '',
    acceptedStates: '',
    transitions: ''
});

export function updateLexerStateFromProject(projectData: any) {
    if (!projectData?.lexing) return;

    const lexingData = projectData.lexing;
    lexerState.update(state => ({
        ...state,
        mode: lexingData.dfa ? 'AUTOMATA' : (lexingData.pairs ? 'REGEX' : null),
        dfa: lexingData.dfa ? {
            accepting: lexingData.dfa.accepting_states || [],
            start: lexingData.dfa.start_state || '',
            states: lexingData.dfa.states || [],
            transitions: lexingData.dfa.transitions || []
        } : null,
        inputs: lexingData.pairs || [],
        states: lexingData.dfa?.states?.join(', ') || '',
        startState: lexingData.dfa?.start_state || '',
        acceptedStates: lexingData.dfa?.accepting_states?.map((s: any) => 
            `${s.state}->${s.token_type}`
        ).join(', ') || '',
        transitions: lexingData.dfa?.transitions?.map((t: any) => 
            `${t.from},${t.label}->${t.to}`
        ).join('\n') || '',
        rules: lexingData.rules || null,
        tokens: lexingData.tokens || null,
        tokens_unidentified: lexingData.tokens_unidentified || null
    }));
}

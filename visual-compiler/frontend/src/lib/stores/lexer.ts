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

export interface AnalyzerData {
    scope_rules: Array<{
        start: string;
        end: string;
    }>;
    type_rules: Array<{
        resultdata: string;
        assignment: string;
        lhsdata: string;
        operator: string[];
        rhsdata: string;
    }>;
    grammar_rules: {
        variablerule: string;
        typerule: string;
        functionrule: string;
        parameterrule: string;
        assignmentrule: string;
        operatorrule: string;
        termrule: string;
    };
}

export interface TranslatorData {
    code: string[];
    translating_rules: Array<{
        sequence: string[];
        translation: string[];
    }>;
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
    userInputRows: Array<{ 
        type: string;
        regex: string;
        error: string;
    }>;
    parser_data?: {
        tree?: {
            root: {
                symbol: string;
                value: string;
                children: any[];
            }
        }
    };
    analyzer_data?: AnalyzerData;
    symbol_table?: Array<{
        name: string;
        type: string;
        scope: number;
    }>;
    translator_data?: TranslatorData;
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
    transitions: '',
    userInputRows: [{ type: '', regex: '', error: '' }]
});

export function updateLexerStateFromProject(projectData: any) {
    if (!projectData?.lexing) return;

    const lexingData = projectData.lexing;
    const parsingData = projectData.parsing;
    const analysingData = projectData.analysing;
    const translatingData = projectData.translating;

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
        tokens: lexingData.tokens || null,
        tokens_unidentified: lexingData.tokens_unidentified || [],
        parser_data: parsingData ? {
            tree: parsingData.tree,
            grammar: parsingData.grammar
        } : undefined,
        analyzer_data: analysingData ? {
            scope_rules: analysingData.scope_rules.map((rule: any) => ({
                start: rule.start || '',
                end: rule.end || ''
            })),
            type_rules: analysingData.type_rules.map((rule: any) => ({
                resultdata: rule.resultdata || '',
                assignment: rule.assignment || '',
                lhsdata: rule.lhsdata || '',
                operator: rule.operator || [],
                rhsdata: rule.rhsdata || ''
            })),
            grammar_rules: analysingData.grammar_rules || {}
        } : undefined,
        symbol_table: analysingData?.symbol_table_artefact?.symbolscopes?.map(symbol => ({
            name: symbol.name,
            type: symbol.type,
            scope: symbol.scope
        })) || [],
        translator_data: translatingData ? {
            code: translatingData.code || [],
            translating_rules: (translatingData.translating_rules || [])
                .map(rule => ({
                    sequence: Array.isArray(rule.sequence) ? rule.sequence.join(', ') : '',
                    translation: rule.translation || []
                }))
        } : undefined
    }));
}

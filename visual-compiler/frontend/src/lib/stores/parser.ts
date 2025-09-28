// lib/stores/parser.ts
import { writable } from 'svelte/store';

interface Translation {
    id: number;
    value: string;
}

interface Rule {
    id: number;
    nonTerminal: string;
    translations: Translation[];
}

interface ParserState {
    grammar_rules: Rule[];
    variables_string: string;
    terminals_string: string;
    rule_id_counter: number;
    translation_id_counter: number;
    show_default_grammar: boolean;
    is_grammar_submitted: boolean;
    hasUnsavedChanges: boolean;
}

const initialState: ParserState = {
    grammar_rules: [{ id: 1, nonTerminal: '', translations: [{ id: 1, value: '' }] }],
    variables_string: '',
    terminals_string: '',
    rule_id_counter: 1,
    translation_id_counter: 1,
    show_default_grammar: false,
    is_grammar_submitted: false,
    hasUnsavedChanges: false
};

export const parserState = writable<ParserState>(initialState);

export const updateParserInputs = (data: Partial<ParserState>) => {
    parserState.update(state => ({
        ...state,
        ...data,
        hasUnsavedChanges: true
    }));
};

export const markParserSubmitted = () => {
    parserState.update(state => ({
        ...state,
        is_grammar_submitted: true,
        hasUnsavedChanges: false
    }));
};

export const resetParserState = () => {
    parserState.set(initialState);
};
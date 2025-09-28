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
    // ADD: Artifact storage
    parseTree: any;
    hasParseTree: boolean;
    parsingError: string | null;
}

const initialState: ParserState = {
    grammar_rules: [{ id: 1, nonTerminal: '', translations: [{ id: 1, value: '' }] }],
    variables_string: '',
    terminals_string: '',
    rule_id_counter: 1,
    translation_id_counter: 1,
    show_default_grammar: false,
    is_grammar_submitted: false,
    hasUnsavedChanges: false,
    // ADD: Artifact fields
    parseTree: null,
    hasParseTree: false,
    parsingError: null
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

// ADD: Function to update artifacts
export const updateParserArtifacts = (parseTree: any, error: string | null = null) => {
    parserState.update(state => ({
        ...state,
        parseTree,
        hasParseTree: !!parseTree,
        parsingError: error
    }));
};

// UPDATE: Project loading function to include artifacts
export const updateParserStateFromProject = (projectData: any) => {
    if (!projectData?.parsing) return;

    const updates: Partial<ParserState> = {};

    // Handle grammar from project data
    if (projectData.parsing.grammar) {
        const grammarData = projectData.parsing.grammar;
        
        if (grammarData.rules && Array.isArray(grammarData.rules)) {
            updates.grammar_rules = grammarData.rules.map((rule: any, index: number) => ({
                id: index + 1,
                nonTerminal: rule.input || '',
                translations: Array.isArray(rule.output) 
                    ? rule.output.map((output: string, tIndex: number) => ({
                        id: tIndex + 1,
                        value: output || ''
                    }))
                    : [{ id: 1, value: rule.output || '' }]
            }));
            
            updates.rule_id_counter = grammarData.rules.length + 1;
        }

        if (grammarData.variables) {
            updates.variables_string = Array.isArray(grammarData.variables) 
                ? grammarData.variables.join(', ') 
                : grammarData.variables || '';
        }
        
        if (grammarData.terminals) {
            updates.terminals_string = Array.isArray(grammarData.terminals) 
                ? grammarData.terminals.join(', ') 
                : grammarData.terminals || '';
        }
    }

    // Handle artifacts - parse tree
    if (projectData.parsing.tree) {
        updates.parseTree = projectData.parsing.tree;
        updates.hasParseTree = true;
        updates.parsingError = null;
    }

    // Mark as submitted if there's data
    if (projectData.parsing && Object.keys(projectData.parsing).length > 0) {
        updates.is_grammar_submitted = true;
        updates.hasUnsavedChanges = false;
    }

    if (Object.keys(updates).length > 0) {
        parserState.update(state => ({
            ...state,
            ...updates
        }));
    }
};
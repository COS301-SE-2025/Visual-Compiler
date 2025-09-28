import { writable } from 'svelte/store';

interface ScopeRule {
    id: number;
    Start: string;
    End: string;
}

interface TypeRule {
    id: number;
    ResultData: string;
    Assignment: string;
    LHSData: string;
    Operator: string[];
    RHSData: string;
}

interface GrammarRule {
    VariableRule: string;
    TypeRule: string;
    FunctionRule: string;
    ParameterRule: string;
    AssignmentRule: string;
    OperatorRule: string;
    TermRule: string;
}

interface AnalyserState {
    scope_rules: ScopeRule[];
    type_rules: TypeRule[];
    grammar_rules: GrammarRule;
    submitted_scope_rules: ScopeRule[];
    submitted_type_rules: TypeRule[];
    submitted_grammar_rules: GrammarRule;
    next_scope_id: number;
    next_type_id: number;
    show_default_rules: boolean;
    rules_submitted: boolean;
    show_symbol_table: boolean;
    hasUnsavedChanges: boolean;
    // ADD: Artifact storage
    symbolTable: any[];
    hasSymbolTable: boolean;
    analyserError: string | null;
    analyserErrorDetails: string | null;
}

const initialState: AnalyserState = {
    scope_rules: [{ id: 0, Start: '', End: '' }],
    type_rules: [{ id: 0, ResultData: '', Assignment: '', LHSData: '', Operator: [''], RHSData: '' }],
    grammar_rules: {
        VariableRule: '',
        TypeRule: '',
        FunctionRule: '',
        ParameterRule: '',
        AssignmentRule: '',
        OperatorRule: '',
        TermRule: ''
    },
    submitted_scope_rules: [],
    submitted_type_rules: [],
    submitted_grammar_rules: {
        VariableRule: '',
        TypeRule: '',
        FunctionRule: '',
        ParameterRule: '',
        AssignmentRule: '',
        OperatorRule: '',
        TermRule: ''
    },
    next_scope_id: 1,
    next_type_id: 1,
    show_default_rules: false,
    rules_submitted: false,
    show_symbol_table: false,
    hasUnsavedChanges: false,
    // ADD: Artifact fields
    symbolTable: [],
    hasSymbolTable: false,
    analyserError: null,
    analyserErrorDetails: null
};

export const analyserState = writable<AnalyserState>(initialState);

// Helper functions
export const updateAnalyserInputs = (data: Partial<AnalyserState>) => {
    analyserState.update(state => ({
        ...state,
        ...data,
        hasUnsavedChanges: true
    }));
};

export const markAnalyserSubmitted = () => {
    analyserState.update(state => ({
        ...state,
        rules_submitted: true,
        hasUnsavedChanges: false
    }));
};

export const resetAnalyserState = () => {
    analyserState.set(initialState);
};

// ADD: Function to update artifacts
export const updateAnalyserArtifacts = (symbolTable: any[], error: string | null = null, errorDetails: string | null = null) => {
    analyserState.update(state => ({
        ...state,
        symbolTable: [...symbolTable],
        hasSymbolTable: symbolTable.length > 0,
        show_symbol_table: symbolTable.length > 0,
        analyserError: error,
        analyserErrorDetails: errorDetails
    }));
};

// Project loading function
export const updateAnalyserStateFromProject = (projectData: any) => {
    if (!projectData?.analysing) return;

    const updates: Partial<AnalyserState> = {};

    // Handle scope rules
    if (projectData.analysing.scope_rules && Array.isArray(projectData.analysing.scope_rules)) {
        updates.scope_rules = projectData.analysing.scope_rules.map((rule: any, index: number) => ({
            id: index,
            Start: rule.start || rule.Start || '',
            End: rule.end || rule.End || ''
        }));
        updates.next_scope_id = updates.scope_rules.length;
        updates.submitted_scope_rules = [...updates.scope_rules];
    }

    // Handle type rules
    if (projectData.analysing.type_rules && Array.isArray(projectData.analysing.type_rules)) {
        updates.type_rules = projectData.analysing.type_rules.map((rule: any, index: number) => ({
            id: index,
            ResultData: rule.resultdata || rule.ResultData || '',
            Assignment: rule.assignment || rule.Assignment || '',
            LHSData: rule.lhsdata || rule.LHSData || '',
            Operator: Array.isArray(rule.operator) ? rule.operator : (Array.isArray(rule.Operator) ? rule.Operator : ['']),
            RHSData: rule.rhsdata || rule.RHSData || ''
        }));
        updates.next_type_id = updates.type_rules.length;
        updates.submitted_type_rules = [...updates.type_rules];
    }

    // Handle grammar rules
    if (projectData.analysing.grammar_rules) {
        const grammarData = projectData.analysing.grammar_rules;
        updates.grammar_rules = {
            VariableRule: grammarData.variablerule || grammarData.VariableRule || '',
            TypeRule: grammarData.typerule || grammarData.TypeRule || '',
            FunctionRule: grammarData.functionrule || grammarData.FunctionRule || '',
            ParameterRule: grammarData.parameterrule || grammarData.ParameterRule || '',
            AssignmentRule: grammarData.assignmentrule || grammarData.AssignmentRule || '',
            OperatorRule: grammarData.operatorrule || grammarData.OperatorRule || '',
            TermRule: grammarData.termrule || grammarData.TermRule || ''
        };
        updates.submitted_grammar_rules = { ...updates.grammar_rules };
    }

    // Handle artifacts - symbol table
    if (projectData.analysing.symbol_table_artefact?.symbolscopes) {
        const symbols = projectData.analysing.symbol_table_artefact.symbolscopes.map((s: any) => ({
            name: s.Name || s.name || 'unknown',
            type: s.Type || s.type || 'unknown', 
            scope: s.Scope || s.scope || 0
        }));
        updates.symbolTable = symbols;
        updates.hasSymbolTable = true;
        updates.show_symbol_table = true;
        updates.analyserError = null;
        updates.analyserErrorDetails = null;
    }

    // Mark as submitted if there's data
    if (projectData.analysing && Object.keys(projectData.analysing).length > 0) {
        updates.rules_submitted = true;
        updates.hasUnsavedChanges = false;
    }

    if (Object.keys(updates).length > 0) {
        analyserState.update(state => ({
            ...state,
            ...updates
        }));
    }
};
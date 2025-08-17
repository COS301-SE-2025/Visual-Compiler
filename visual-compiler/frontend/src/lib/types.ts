export type NodeType = 'source' | 'lexer' | 'parser' | 'analyser' | 'translator' ;

export interface Token {
	type: string;
	value: string;
}

export interface SyntaxTreeNode {
	symbol: string;
	value: string;
	children: SyntaxTreeNode[] | null;
}

export interface SyntaxTree {
	root: SyntaxTreeNode;
}

export interface Symbol {
    name: string;
    type: string;
    scope: number; 
}

export interface SymbolTable {
    symbols: Symbol[];
}

export interface NodeConnection {
    id: string;
    sourceNodeId: string;
    targetNodeId: string;
    sourceType: NodeType;
    targetType: NodeType;
    sourceAnchor?: string;
    targetAnchor?: string;
}


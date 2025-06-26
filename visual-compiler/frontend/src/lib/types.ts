export type NodeType = 'source' | 'lexer' | 'parser';

export interface Token {
  Type: string;  
  Value: string; 
}

export interface SyntaxTreeNode {
  symbol: string;
  value: string;
  children: SyntaxTreeNode[] | null;
}

export interface SyntaxTree {
  root: SyntaxTreeNode;
}
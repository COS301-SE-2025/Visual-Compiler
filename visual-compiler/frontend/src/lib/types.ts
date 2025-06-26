export type NodeType = 'source' | 'lexer' | 'parser';

export interface Token {
  Type: string;  
  Value: string; 
}
export type NodeType = 'source' | 'lexer' ;

export interface Token {
  Type: string;  // Uppercase to match backend
  Value: string; // Value from lexer
}
# **CORE SERVICES**
This document contains a short description of functions to be called in the API
  
## Lexer functions
- Read Regex rules from user and store them in appropriate struct
  - `func ReadRegexRules(input []byte) ([]TypeRegex, error)`
  - Input example:  
  ```go
  input := []byte(`[{"Type": "KEYWORD""},{"Type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"}]`)
- Tokenise source code using regex
  - `func CreateTokens(source string, rules []TypeRegex) ([]TypeValue, []string, error) `
  - Input example:  
  ```go
  source := "int x = 3;"
  rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|int|else)\\b"},  
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
		{Type: "OPERATOR", Regex: "="},
		{Type: "NUMBER", Regex: "\\d+(\\.\\d+)?"},
		{Type: "PUNCTUATION", Regex: ";"},
	}
- Tokenise source code using DFA
  - `func CreateTokensFromDFA(source_code string, dfa Automata) ([]TypeValue, []string, error)`
  ```go
  source := "int x = 3;"
  dfa := services.Automata{
		States: []string{"START", "S1", "S2", "S3", "S4", "S5"},
		Transitions: []services.Transition{
			{From: "START", To: "S1", Label: "i"},
			{From: "S1", To: "S5", Label: "n"},
			{From: "S5", To: "S4", Label: "t"},
			{From: "S1", To: "S6", Label: "f"},
			{From: "START", To: "S2", Label: "0123456789"},
			{From: "S2", To: "S2", Label: "0123456789"},
			{From: "START", To: "S3", Label: "abcdefghijklmnopqrstuvwxyz"},
			{From: "S3", To: "S3", Label: "abcdefghijklmnopqrstuvwxyz0123456789"},
		},
		Start: "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}
- Convert DFA to regex
  - `func ConvertDFAToRegex(dfa Automata) ([]TypeRegex, error)`  
  ```go
  dfa := services.Automata{
		States: []string{"START", "S1", "S2", "S3", "S4", "S5"},
		Transitions: []services.Transition{
			{From: "START", To: "S1", Label: "i"},
			{From: "S1", To: "S5", Label: "n"},
			{From: "S5", To: "S4", Label: "t"},
			{From: "S1", To: "S6", Label: "f"},
			{From: "START", To: "S2", Label: "0123456789"},
			{From: "S2", To: "S2", Label: "0123456789"},
			{From: "START", To: "S3", Label: "abcdefghijklmnopqrstuvwxyz"},
			{From: "S3", To: "S3", Label: "abcdefghijklmnopqrstuvwxyz0123456789"},
		},
		Start: "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}
- Convert regex to NFA
  - `func ConvertRegexToNFA(regexes map[string]string) (Automata, error)`
  - Input example:
  ```go
  regexes := map[string]string {
    "IDENTIFIER": "[a-zA-Z_]\\w*",
    "NUMBER":     "\\d+(\\.\\d+)?",
    "KEYWORD":    "if|else",
	}
- Convert regex to DFA
  - `func ConvertRegexToDFA(regexes map[string]string) (Automata, error)`
  - Input example:
  ```go
  regexes := map[string]string {
    "IDENTIFIER": "[a-zA-Z_]\\w*",
    "NUMBER":     "\\d+(\\.\\d+)?",
    "KEYWORD":    "if|else",
	}
- Convert NFA to DFA
  - `func ConvertNFAToDFA(nfa Automata) (Automata, error)`
  - Input example:
  ```go
  nfa := services.Automata{
		States: []string{"START", "S1", "S2", "S3", "S4", "S5", "S6"},
		Transitions: []services.Transition{
			{From: "START", To: "S1", Label: "i"},
			{From: "START", To: "S1", Label: "f"},
			{From: "S1", To: "S5", Label: "n"},
			{From: "S5", To: "S4", Label: "t"},
			{From: "S1", To: "S6", Label: "ε"},
			{From: "START", To: "S2", Label: "0-9"},
			{From: "S2", To: "S2", Label: "0-9"},
			{From: "START", To: "S3", Label: "a-z"},
			{From: "S3", To: "S3", Label: "a-z0-9"},
			{From: "START", To: "S3", Label: "ε"},
		},
		Start: "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

## Parser functions
- Read grammar from user and ensure structure is correct
  - `func ReadGrammar(input []byte) (Grammar, error)`
  ```go
  input := []byte(`{
		"variables": ["STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"],
		"terminals": ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"],
		"start": "STATEMENT",
		"rules": [
			{ "input": "STATEMENT", "output": ["DECLARATION", "SEPARATOR"] },
			{ "input": "DECLARATION", "output": ["TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"] },
			{ "input": "EXPRESSION", "output": ["TERM", "OPERATOR", "TERM"] },
			{ "input": "TERM", "output": ["INTEGER"] },
			{ "input": "TYPE", "output": ["KEYWORD"] }
		]
	}`)
- Create syntax tree using tokens and a grammar as input
  - `func CreateSyntaxTree(tokens []TypeValue, grammar Grammar) (SyntaxTree, error)`
  ```go
  tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "blue"},
		{Type: "ASSIGNMENT", Value: "="},
		{Type: "INTEGER", Value: "13"},
		{Type: "OPERATOR", Value: "+"},
		{Type: "INTEGER", Value: "89"},
		{Type: "SEPARATOR", Value: ";"},
	}

	grammar := services.Grammar{
		Variables: []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		Terminals: []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		Start:     "STATEMENT",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
			{Input: "DECLARATION", Output: []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{Input: "EXPRESSION", Output: []string{"TERM", "OPERATOR", "TERM"}},
			{Input: "TERM", Output: []string{"INTEGER"}},
			{Input: "TYPE", Output: []string{"KEYWORD"}},
		},
	}
- Create a string representation fo the syntax tree
  - `func ConvertTreeToString(node *TreeNode, branch_indent string, is_tail bool) string`
  ```go
	node := syntax_tree.Root
	branch_indent := ""
	is_tail := true
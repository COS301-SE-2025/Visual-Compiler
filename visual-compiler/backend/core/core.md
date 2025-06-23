# **CORE SERVICES**
This document contains a short description of functions to be called in the API
  
## Lexer functions
- Read Regex rules from user and store them in appropriate struct
  - `ReadRegexRules(input []byte)`
  - Input example:  
  ```go
  input := []byte(`[{"Type": "KEYWORD""},{"Type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"}]`)
- Tokenise source code using regex
  - `CreateTokens(source string, rules []TypeRegex)`
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
  - `CreateTokensFromDFA(source string, dfa Automata)`
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
  - `ConvertDFAToRegex(dfa Automata)`  
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
  - `ConvertRegexToNFA(regexes map[string]string, nfa Automata)`
- Convert regex to DFA
  - `ConvertRegexToDFA(regexes map[string]string)`
- Convert NFA to DFA
  - `ConvertNFAToDFA(nfa Automata, dfa Automata)`


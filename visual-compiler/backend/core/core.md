# **CORE SERVICES**
This document contains a short description of functions to be called in the API
  
## Lexer functions
- Read Regex rules from user and store them in appropriate struct
  - `ReadRegexRules(input []byte)`
- Tokenise source code using regex
  - `CreateTokens(source string, rules []TypeRegex)`
- Read DFA data from user and store it in the appropriate struct
  - `ReadDFA(input []byte)`
- Tokenise source code using DFA
  - `CreateTokensFromDFA(source string, dfa Automata)`
- Convert DFA to regex
  - `ConvertDFAToRegex(dfa Automata)`  
- Convert regex to NFA
  - `ConvertRegexToNFA(regexes map[string]string, nfa Automata)`
- Convert regex to DFA
  - `ConvertRegexToDFA(regexes map[string]string)`
- Convert NFA to DFA
  - `ConvertNFAToDFA(nfa Automata, dfa Automata)`


package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Struct for the grammar
type Grammar struct {
	Variables []string      `json:"variables"`
	Terminals []string      `json:"terminals"`
	Start     string        `json:"start"`
	Rules     []ParsingRule `json:"rules"`
}

type ParsingRule struct {
	Input  string   `json:"input"`
	Output []string `json:"output"`
}

// Struct for the syntax tree
type SyntaxTree struct {
	Root *TreeNode `json:"root"`
}

type TreeNode struct {
	Symbol   string      `json:"symbol"`
	Value    string      `json:"value"`
	Children []*TreeNode `json:"children"`
}

// Struct to track parsing
type ParseState struct {
	position int
	tokens   []TypeValue
	grammar  Grammar
}

// Name: ReadGrammar
//
// Parameters: []byte
//
// Return: Grammar, error
//
// Receive an array of grammar rules and return them in the struct
func ReadGrammar(input []byte) (Grammar, error) {

	grammar := Grammar{}

	if len(input) == 0 {
		return Grammar{}, fmt.Errorf("no grammar entered")
	}

	err := json.Unmarshal(input, &grammar)

	if err != nil {
		return Grammar{}, fmt.Errorf("invalid JSON for grammar: %v", err)
	}

	valid := false

	for _, vars := range grammar.Variables {
		if vars == grammar.Start {
			valid = true
			break
		}
	}

	if !valid {
		return Grammar{}, fmt.Errorf("start symbol is not in the list of variables")
	}

	for i, term := range grammar.Terminals {
		grammar.Terminals[i] = strings.ToUpper(term)
	}

	for i, rule := range grammar.Rules {

		filtered := []string{}

		for _, out := range rule.Output {
			if out != "" && out != "ε" {
				filtered = append(filtered, out)
			}
		}

		grammar.Rules[i].Output = filtered
	}

	return grammar, nil
}

// Name: CreateSyntaxTree
//
// Parameters: []TypeValue, Grammar
//
// Return: SyntaxTree, error
//
// Recursively build the syntax tree from the tokens and the grammar
func CreateSyntaxTree(tokens []TypeValue, grammar Grammar) (SyntaxTree, error) {

	if len(tokens) == 0 {
		return SyntaxTree{}, fmt.Errorf("no tokens found")
	}

	if grammar.Start == "" {
		return SyntaxTree{}, fmt.Errorf("no start variable found")
	}

	state := &ParseState{
		position: 0,
		tokens:   tokens,
		grammar:  grammar,
	}

	root, new_position, success := parseSymbol(state, grammar.Start, 0)

	if !success || new_position != len(tokens) {
		return SyntaxTree{}, fmt.Errorf("syntax error!")
	}

	return SyntaxTree{Root: root}, nil
}

// Name: parseSymbol
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to parse a variable or a terminal starting at the given position
func parseSymbol(state *ParseState, symbol string, position int) (*TreeNode, int, bool) {

	return nil, 0, false
}

// Name: parseTerminal
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to match a terminal symbol with the current token
func parseTerminal(state *ParseState, terminal string, position int) (*TreeNode, int, bool) {

	return nil, 0, false
}

// Name: parseVariable
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to parse a variable using all applicable rules
func parseVariable(state *ParseState, variable string, position int) (*TreeNode, int, bool) {

	return nil, 0, false
}

// Name: tryRule
//
// Parameters: *ParseState, ParsingRule, int
//
// Return: *TreeNode, int, bool
//
// Attempts to apply a specific parsing rule
func tryRule(state *ParseState, rule ParsingRule, position int) (*TreeNode, int, bool) {

	return nil, 0, false
}

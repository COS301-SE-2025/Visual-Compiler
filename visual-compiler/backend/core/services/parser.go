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
	Position int
	Tokens   []TypeValue
	Grammar  Grammar
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
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := ParseSymbol(state, grammar.Start, 0)

	if !success || new_position != len(tokens) {
		return SyntaxTree{}, fmt.Errorf("syntax error")
	}

	return SyntaxTree{Root: root}, nil
}

// Name: ParseSymbol
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to parse a variable or a terminal starting at the given position
func ParseSymbol(state *ParseState, symbol string, position int) (*TreeNode, int, bool) {

	if position >= len(state.Tokens) {
		return nil, position, false
	}

	found := false

	for _, terminal := range state.Grammar.Terminals {
		if terminal == symbol {
			found = true
		}
	}

	if found {
		return ParseTerminal(state, symbol, position)
	} else {
		return ParseVariable(state, symbol, position)
	}
}

// Name: ParseTerminal
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to match a terminal symbol with the current token
func ParseTerminal(state *ParseState, terminal string, position int) (*TreeNode, int, bool) {

	if position >= len(state.Tokens) {
		return nil, position, false
	}

	token := state.Tokens[position]

	if token.Type == terminal {

		node := &TreeNode{
			Symbol:   terminal,
			Value:    token.Value,
			Children: nil,
		}

		return node, position + 1, true
	}

	return nil, position, false
}

// Name: ParseVariable
//
// Parameters: *ParseState, string, int
//
// Return: *TreeNode, int, bool
//
// Attempts to parse a variable using all applicable rules
func ParseVariable(state *ParseState, variable string, position int) (*TreeNode, int, bool) {

	for _, rule := range state.Grammar.Rules {

		if rule.Input == variable {

			node, new_position, success := TryRule(state, rule, position)

			if success {
				return node, new_position, true
			}
		}
	}

	return nil, position, false
}

// Name: tryRule
//
// Parameters: *ParseState, ParsingRule, int
//
// Return: *TreeNode, int, bool
//
// Attempts to apply a specific parsing rule
func TryRule(state *ParseState, rule ParsingRule, position int) (*TreeNode, int, bool) {

	node := &TreeNode{
		Symbol:   rule.Input,
		Value:    "",
		Children: make([]*TreeNode, 0),
	}

	current_position := position

	for _, symbol := range rule.Output {

		child_node, new_position, success := ParseSymbol(state, symbol, current_position)

		if !success {
			return nil, position, false
		}

		node.Children = append(node.Children, child_node)
		current_position = new_position
	}

	return node, current_position, true
}

// Name: ConvertTreeToString
//
// Parameters: *TreeNode, string, bool
//
// Return: string
//
// Recursively build a string of the syntax tree and returns it
func ConvertTreeToString(node *TreeNode, branch_indent string, is_tail bool) string {
	if node == nil {
		return branch_indent + "\n"
	}

	var final_tree string
	branch_char := "├── "
	tail_char := "└── "

	new_string := branch_indent
	if is_tail {
		final_tree += fmt.Sprintf("%s%sSymbol: %s, Value: %s\n", branch_indent, tail_char, node.Symbol, node.Value)
		new_string += "    "
	} else {
		final_tree += fmt.Sprintf("%s%sSymbol: %s, Value: %s\n", branch_indent, branch_char, node.Symbol, node.Value)
		new_string += "│   "
	}

	for i, child := range node.Children {
		tail_node := false
		if i == len(node.Children)-1 {
			tail_node = true
		}
		final_tree += ConvertTreeToString(child, new_string, tail_node)
	}

	return final_tree
}

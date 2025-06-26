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
	Token    string      `json:"token"`
	Value    string      `json:"value"`
	Children []*TreeNode `json:"children"`
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
		if grammar.Start == vars {
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

	return SyntaxTree{}, nil
}

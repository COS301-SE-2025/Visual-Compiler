package unit_tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestReadGrammer_NoGrammar(t *testing.T) {
	input := []byte{}
	_, err := services.ReadGrammar(input)

	if err == nil {
		t.Errorf("Error expected for no grammar")
	} else {
		if err.Error() != fmt.Errorf("no grammar entered").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadGrammer_InvalidJson(t *testing.T) {
	input := []byte("Invalid Json")
	_, err := services.ReadGrammar(input)

	if err == nil {
		t.Errorf("Error expected for invalid json")
	} else {
		if err.Error() != fmt.Errorf("invalid JSON for grammar: invalid character 'I' looking for beginning of value").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadGrammer_InvalidGrammar(t *testing.T) {
	input := []byte(`{
		"variables": ["STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"],
		"terminals": ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"],
		"start": "RANDOM",
		"rules": [
			{ "input": "STATEMENT", "output": ["DECLARATION", "SEPARATOR"] },
			{ "input": "DECLARATION", "output": ["TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"] },
			{ "input": "EXPRESSION", "output": ["TERM", "OPERATOR", "TERM"] },
			{ "input": "TERM", "output": ["INTEGER"] },
			{ "input": "TYPE", "output": ["KEYWORD"] }
		]
	}`)
	_, err := services.ReadGrammar(input)

	if err == nil {
		t.Errorf("Error expected for no start")
	} else {
		if err.Error() != fmt.Errorf("start symbol is not in the list of variables").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadGrammer_ValidGrammar(t *testing.T) {

	expected_res := services.Grammar{
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
	grammar, err := services.ReadGrammar(input)

	if err == nil {
		if grammar.Start != expected_res.Start {
			t.Errorf("Incorrect start: %v", grammar.Start)
		}
		for i, variable := range grammar.Variables {
			if variable != expected_res.Variables[i] {
				t.Errorf("Grammar incorrect: %v != %v", variable, expected_res.Variables[i])
			}
		}
		for i, terminal := range grammar.Terminals {
			if terminal != expected_res.Terminals[i] {
				t.Errorf("Grammar incorrect: %v != %v", terminal, expected_res.Terminals[i])
			}
		}
		for i, rule := range grammar.Rules {
			if rule.Input != expected_res.Rules[i].Input {
				t.Errorf("Grammar incorrect: %v != %v", rule, expected_res.Rules[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestReadGrammar_LeftRecursion(t *testing.T) {
	input := []byte(`{
		"variables": ["A", "B"],
		"terminals": ["a", "b"],
		"start": "A",
		"rules": [
			{ "input": "A", "output": ["A", "b"] },
			{ "input": "A", "output": ["a", "B"] },
			{ "input": "B", "output": ["ε"] },
		]
	}`)

	_, err := services.ReadGrammar(input)

	if err == nil {
		t.Errorf("Expected error for left recursion in grammar")
	}
}

func TestCreateSyntaxTree_NoTokens(t *testing.T) {
	tokens := []services.TypeValue{}

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
	_, err := services.CreateSyntaxTree(tokens, grammar)

	if err == nil {
		t.Errorf("Error supposed to occur for no tokens")
	} else {
		if err.Error() != fmt.Errorf("no tokens found").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestCreateSyntaxTree_NoStart(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "blue"},
		{Type: "ASSIGNMENT", Value: "="},
		{Type: "INTEGER", Value: "13"},
		{Type: "OPERATOR", Value: "+"},
		{Type: "INTEGER", Value: "89"},
		{Type: "SEPARATOR", Value: ";"}}

	grammar := services.Grammar{
		Variables: []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		Terminals: []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		Start:     "",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
			{Input: "DECLARATION", Output: []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{Input: "EXPRESSION", Output: []string{"TERM", "OPERATOR", "TERM"}},
			{Input: "TERM", Output: []string{"INTEGER"}},
			{Input: "TYPE", Output: []string{"KEYWORD"}},
		},
	}
	_, err := services.CreateSyntaxTree(tokens, grammar)

	if err == nil {
		t.Errorf("Error supposed to occur for no start variable")
	} else {
		if err.Error() != fmt.Errorf("no start variable found").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestCreateSyntaxTree_SyntaxError(t *testing.T) {
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
		Start:     "RANDOM",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
			{Input: "DECLARATION", Output: []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{Input: "EXPRESSION", Output: []string{"TERM", "OPERATOR", "TERM"}},
			{Input: "TERM", Output: []string{"INTEGER"}},
			{Input: "TYPE", Output: []string{"KEYWORD"}},
		},
	}
	_, err := services.CreateSyntaxTree(tokens, grammar)

	if err == nil {
		t.Errorf("Error supposed to occur for incorrect syntax")
	} else {
		if err.Error() != fmt.Errorf("syntax error").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestCreateSyntaxTree_Valid(t *testing.T) {

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
	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur for not tokens")
	} else {
		if syntax_tree.Root == nil {
			t.Errorf("No root")
		}
		if syntax_tree.Root.Children == nil {
			t.Errorf("No children")
		}
		if syntax_tree.Root.Symbol == "" {
			t.Errorf("No symbol")
		}
	}
}

func TestCreateSyntaxTree_TwoRules(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "INTEGER", Value: "42"},
	}

	grammar := services.Grammar{
		Variables: []string{"VARIABLE"},
		Terminals: []string{"INTEGER", "STRING"},
		Start:     "VARIABLE",
		Rules: []services.ParsingRule{
			{Input: "VARIABLE", Output: []string{"STRING"}},
			{Input: "VARIABLE", Output: []string{"INTEGER"}},
		},
	}

	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	}
	if syntax_tree.Root == nil || syntax_tree.Root.Children[0].Value != "42" {
		t.Errorf("Incorrect syntax tree structure")
	}
}

func TestCreateSyntaxTree_TwoRulesRecursive(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "STRING", Value: "red"},
		{Type: "STRING", Value: "green"},
		{Type: "STRING", Value: "blue"},
	}

	grammar := services.Grammar{
		Variables: []string{"LIST", "STATEMENT"},
		Terminals: []string{"STRING"},
		Start:     "LIST",
		Rules: []services.ParsingRule{
			{Input: "LIST", Output: []string{"STATEMENT"}},
			{Input: "LIST", Output: []string{"STATEMENT", "LIST"}},
			{Input: "STATEMENT", Output: []string{"STRING"}},			
		},
	}

	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	}
	if syntax_tree.Root == nil || len(syntax_tree.Root.Children) != 2 || syntax_tree.Root.Children[1].Symbol != "LIST" {
		t.Errorf("Incorrect syntax tree structure")
	}
}

func TestCreateSyntaxTree_Complex(t *testing.T) {

	tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "blue"},
		{Type: "ASSIGNMENT", Value: "="},
		{Type: "INTEGER", Value: "13"},
		{Type: "OPERATOR", Value: "+"},
		{Type: "INTEGER", Value: "89"},
		{Type: "SEPARATOR", Value: ";"},
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "function_do"},
		{Type: "OPEN_BRACKETS", Value: "("},
		{Type: "CLOSE_BRACKETS", Value: ")"},
	}

	grammar := services.Grammar{
		Variables: []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM", "PROGRAM"},
		Terminals: []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR", "STRING", "OPEN_BRACKETS", "CLOSE_BRACKETS"},
		Start:     "PROGRAM",
		Rules: []services.ParsingRule{
			{Input: "PROGRAM", Output: []string{"STATEMENT", "FUNCTION"}},
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
			{Input: "FUNCTION", Output: []string{"TYPE", "IDENTIFIER", "OPEN_BRACKETS", "CLOSE_BRACKETS"}},
			{Input: "DECLARATION", Output: []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{Input: "EXPRESSION", Output: []string{"TERM", "OPERATOR", "TERM"}},
			{Input: "TERM", Output: []string{"INTEGER"}},
			{Input: "TYPE", Output: []string{"KEYWORD"}},
		},
	}
	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur for not tokens: %v", err)
	} else {
		if syntax_tree.Root == nil {
			t.Errorf("No root")
		}
		if syntax_tree.Root.Children == nil {
			t.Errorf("No children")
		}
		if syntax_tree.Root.Symbol == "" {
			t.Errorf("No symbol")
		}
	}
}

func TestCreateSyntaxTree_MoreComplex(t *testing.T) {

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
		Terminals: []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR", "STRING"},
		Start:     "STATEMENT",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
			{Input: "DECLARATION", Output: []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{Input: "EXPRESSION", Output: []string{"TERM", "OPERATOR", "TERM"}},
			{Input: "TERM", Output: []string{"STRING"}},
			{Input: "TERM", Output: []string{"INTEGER"}},
			{Input: "TYPE", Output: []string{"KEYWORD"}},
		},
	}
	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur for not tokens: %v", err)
	} else {
		if syntax_tree.Root == nil {
			t.Errorf("No root")
		}
		if syntax_tree.Root.Children == nil {
			t.Errorf("No children")
		}
		if syntax_tree.Root.Symbol == "" {
			t.Errorf("No symbol")
		}
	}
}

func TestCreateSyntaxTree_EpsilonTransition(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "INTEGER", Value: "13"},
	}

	grammar := services.Grammar{
		Variables: []string{"PROGRAM"},
		Terminals: []string{"INTEGER"},
		Start:     "PROGRAM",
		Rules: []services.ParsingRule{
			{Input: "PROGRAM", Output: []string{"INTEGER", "ε"}},
		},
	}

	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	}
	if syntax_tree.Root == nil || len(syntax_tree.Root.Children) != 1 {
		t.Errorf("Incorrect syntax tree structure")
	}
}

func TestCreateSyntaxTree_TokenTerminalMismatch(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "ASSIGNMENT", Value: "="},
		{Type: "INTEGER", Value: "12"},
	}

	grammar := services.Grammar{
		Variables: []string{"STATEMENT"},
		Terminals: []string{"IDENTIFIER", "ASSIGNMENT", "NUMBER"},
		Start:     "STATEMENT",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"IDENTIFIER", "ASSIGNMENT", "NUMBER"}},
		},
	}

	_, err := services.CreateSyntaxTree(tokens, grammar)

	if err == nil {
		t.Errorf("Expected error for token types not matching grammar terminals")
	}
}

func TestParseSymbol_NoTokens(t *testing.T) {
	tokens := []services.TypeValue{}

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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseSymbol(state, grammar.Start, 0)

	if success != false {
		t.Errorf("Not supposed to be successful")
	}
	if root != nil {
		t.Errorf("Expected root to be nil")
	}
	if new_position != 0 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseSymbol_Valid(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseSymbol(state, grammar.Start, 0)

	if success != true {
		t.Errorf("Not supposed to be successful")
	}
	if root == nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 7 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseTerminal_NoTokens(t *testing.T) {
	tokens := []services.TypeValue{}

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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseTerminal(state, grammar.Start, 0)

	if success != false {
		t.Errorf("Not supposed to be successful")
	}
	if root != nil {
		t.Errorf("Expected root to be nil")
	}
	if new_position != 0 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseTerminal_Valid(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseTerminal(state, grammar.Start, 1)

	if success != false {
		t.Errorf("Not supposed to be successful")
	}
	if root != nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 1 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseTerminal_ValidChildren(t *testing.T) {
	tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
	}

	grammar := services.Grammar{
		Variables: []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		Terminals: []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		Start:     "STATEMENT",
		Rules: []services.ParsingRule{
			{Input: "STATEMENT", Output: []string{"DECLARATION", "SEPARATOR"}},
		},
	}
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseTerminal(state, "KEYWORD", 0)

	if success != true {
		t.Errorf("Not supposed to be false")
	}
	if root == nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 1 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseVariable_False(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseVariable(state, grammar.Start, 1)

	if success != false {
		t.Errorf("Not supposed to be successful")
	}
	if root != nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 1 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestParseVariable_True(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.ParseVariable(state, "STATEMENT", 0)

	if success != true {
		t.Errorf("Not supposed to be false")
	}
	if root == nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 7 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestTryRule_False(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.TryRule(state, grammar.Rules[0], 1)

	if success != false {
		t.Errorf("Not supposed to be successful")
	}
	if root != nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 1 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestTryRule_True(t *testing.T) {
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
	state := &services.ParseState{
		Position: 0,
		Tokens:   tokens,
		Grammar:  grammar,
	}

	root, new_position, success := services.TryRule(state, grammar.Rules[0], 0)

	if success != true {
		t.Errorf("Not supposed to be false")
	}
	if root == nil {
		t.Errorf("Expected root to not be nil")
	}
	if new_position != 7 {
		t.Errorf("Expected position to be : %v", new_position)
	}
}

func TestConvertTreeToString_Valid(t *testing.T) {

	expected_res := `
└──  STATEMENT
    ├──  DECLARATION
    │   ├──  TYPE
    │   │   └──  KEYWORD: int
    │   ├──  IDENTIFIER: blue
    │   ├──  ASSIGNMENT: =
    │   └──  EXPRESSION
    │       ├──  TERM
    │       │   └──  INTEGER: 13
    │       ├──  OPERATOR: +
    │       └──  TERM
    │           └──  INTEGER: 89
    └──  SEPARATOR: ;`

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
	syntax_tree, err := services.CreateSyntaxTree(tokens, grammar)

	if err != nil {
		t.Errorf("Error not supposed to occur for not tokens")
	} else {
		tree_string := services.ConvertTreeToString(syntax_tree.Root, "", true)
		if strings.TrimSpace(tree_string) != strings.TrimSpace(expected_res) {
			t.Errorf("Incorrect tree: \n%v", tree_string)
		}
	}
}

func TestConvertTreeToString_Invalid(t *testing.T) {

	expected_res := ""

	tree_string := services.ConvertTreeToString(nil, "", true)
	if strings.TrimSpace(tree_string) != strings.TrimSpace(expected_res) {
		t.Errorf("Incorrect tree: \n%v", tree_string)
	}

}

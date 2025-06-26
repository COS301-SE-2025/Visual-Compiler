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

func TestPrintTree_Valid(t *testing.T) {

	expected_res := `Symbol: STATEMENT, Value: 
  Symbol: DECLARATION, Value: 
    Symbol: TYPE, Value: 
      Symbol: KEYWORD, Value: int
    Symbol: IDENTIFIER, Value: blue
    Symbol: ASSIGNMENT, Value: =
    Symbol: EXPRESSION, Value: 
      Symbol: TERM, Value: 
        Symbol: INTEGER, Value: 13
      Symbol: OPERATOR, Value: +
      Symbol: TERM, Value: 
        Symbol: INTEGER, Value: 89
  Symbol: SEPARATOR, Value: ;
  `

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
		tree_string := services.ConvertTreeToString(syntax_tree.Root, "", "")
		if strings.TrimSpace(tree_string) != strings.TrimSpace(expected_res) {
			t.Errorf("Incorrect tree")
		}
	}
}

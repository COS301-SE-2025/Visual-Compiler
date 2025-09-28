package unit_tests

import (
	"fmt"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

// ==================== //
// TEST: ReadRegexRules //
// ==================== //

type add_invalid_test struct {
	arg1     []byte
	expected error
}

var add_invalid_tests = []add_invalid_test{
	{[]byte(`[
		{"type": "keyword","regex":"\\b(if|else)\\b},
		{"type: "identifier","regex":"[a-zA-Z_]\\w*},
	]`), fmt.Errorf(`invalid JSON for rules: invalid character '\n' in string literal`)},
	{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b"},{"type": "IDENITIFER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid JSON for rules: invalid character ']' looking for beginning of value")},
	{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b}]`), fmt.Errorf("invalid JSON for rules: unexpected end of JSON input")},
	{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b},{"type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid JSON for rules: invalid character 't' after object key:value pair")},
}

func TestReadRegexRules_InvalidCharacters(t *testing.T) {
	for _, test := range add_invalid_tests {
		_, err := services.ReadRegexRules(test.arg1)
		if err == nil {
			t.Errorf("Failed for invalid JSON object: %v", err)
		} else if err.Error() == test.expected.Error() {
			t.Logf("Test: %v", err)
		} else {
			t.Errorf("Incorrect error: %v,%v", err, test.expected)
		}
	}
}

func TestReadRegexRules_Fail(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD""},{"Type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"}]`)
	_, err := services.ReadRegexRules(c_input)
	if err == nil {
		t.Errorf("Test not supposed to pass for invalid input")
	} else if err.Error() == fmt.Errorf(`invalid JSON for rules: invalid character '"' after object key:value pair`).Error() {
		t.Logf("Test: %v", err)
	} else {
		t.Errorf("Incorrect error: %v", err)
	}
}

func TestReadRegexRules_NoInput(t *testing.T) {
	c_input := []byte{}
	_, err := services.ReadRegexRules(c_input)
	if err != nil {
		if err.Error() == fmt.Errorf("no rules specified").Error() {
			t.Logf("No rules specified")
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Test not supposed to pass for no input")
	}
}

func TestReadRegexRules_ValidInput(t *testing.T) {
	c_input := []byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b"},{"type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"}]`)
	rules, err := services.ReadRegexRules(c_input)
	if err != nil {
		t.Errorf("Failed for valid input")
	}

	expected_res := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|else)\\b"},
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
	}
	if len(rules) != len(expected_res) {
		t.Errorf("Not enough rules created")
	}

	for i, token := range rules {
		if token != expected_res[i] {
			t.Errorf("Tokenisation incorrect: %v", token)
		}
	}
}

func TestReadRegexRules_InvalidRegex(t *testing.T) {
	c_input := []byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b"},{"type": "IDENTIFIER","regex":"[a-zA-Z_"}]`)
	_, err := services.ReadRegexRules(c_input)
	if err == nil {
		t.Errorf("Expected error")
	}
	if err != nil {
		if err.Error() == fmt.Errorf("invalid regex input: error parsing regexp: missing closing ]: `[a-zA-Z_`").Error() {
			t.Logf("Error received successfully")
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Test not supposed to pass for no input")
	}
}

// ==================== //
//  TEST: CreateTokens  //
// ==================== //

func TestCreateTokens_EmptySourceCode(t *testing.T) {
	source_code := ""
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|else)\\b"},
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
	}

	_, _, err := services.CreateTokens(source_code, rules)

	if err != nil {
		if err.Error() != fmt.Errorf("source code is empty").Error() {
			t.Errorf("Incorrect error: %v", err)
		} else {
			t.Logf("Error sent for empty source code")
		}
	} else {
		t.Errorf("No error sent for no source code")
	}
}

func TestCreateTokens_NoRules(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{}

	_, _, err := services.CreateTokens(source_code, rules)

	if err != nil {
		if err.Error() != fmt.Errorf("no tokenisation rules specified").Error() {
			t.Errorf("Incorrect error: %v", err)
		} else {
			t.Logf("Error sent for empty source code")
		}
	} else {
		t.Errorf("No error sent for no source code")
	}
}


func TestCreateTokens_NoTokens(t *testing.T) {
	source_code := "int x = 13;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "(function|return)"},
	}
	expected_leftovers := []string{
		"int x = 13;",
	}

	_, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_KeywordTokens(t *testing.T) {
	source_code := "int x = 13;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "(int|str|bool)"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
	}
	expected_leftovers := []string{
		"x = 13;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_IdentifierTokens(t *testing.T) {
	source_code := "int x = 13;"
	rules := []services.TypeRegex{
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]+"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "IDENTIFIER", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
	}
	expected_leftovers := []string{
		"= 13;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_PositiveIntegerTokens(t *testing.T) {
	source_code := "12 + 13;"
	rules := []services.TypeRegex{
		{Type: "INTEGER", Regex: "[1-9][0-9]*"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "INTEGER", Value: "12"},
	}
	expected_leftovers := []string{
		"+ 13;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_NegativeIntegerTokens(t *testing.T) {
	source_code := "-13 + -12;"
	rules := []services.TypeRegex{
		{Type: "NUMBER", Regex: "[-][1-9][0-9]*"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "NUMBER", Value: "-13"},
	}
	expected_leftovers := []string{
		"+ -12;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_PositiveFloatTokens(t *testing.T) {
	source_code := "1.2 + 1.3;"
	rules := []services.TypeRegex{
		{Type: "NUMBER", Regex: "[1-9][0-9]*.[0-9]*[1-9]"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "NUMBER", Value: "1.2"},
	}
	expected_leftovers := []string{
		"+ 1.3;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_NegativeFloatTokens(t *testing.T) {
	source_code := "-1.3 + -1.2;"
	rules := []services.TypeRegex{
		{Type: "NUMBER", Regex: "[-][1-9][0-9]*.[0-9]*[1-9]"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "NUMBER", Value: "-1.3"},
	}
	expected_leftovers := []string{
		"+ -1.2;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_OperatorTokens(t *testing.T) {
	source_code := "++i;"
	rules := []services.TypeRegex{
		{Type: "OPERATOR", Regex: "\\+\\+"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "OPERATOR", Value: "++"},
	}
	expected_leftovers := []string{
		"i;",
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if leftovers[0] != expected_leftovers[0] {
			t.Errorf("Tokenisation incorrect: %v != %v", leftovers[0], expected_leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_PunctuationTokens(t *testing.T) {
	source_code := ";"
	rules := []services.TypeRegex{
		{Type: "PUNCTUATION", Regex: ";"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "PUNCTUATION", Value: ";"},
	}

	tokens, _, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_Complex(t *testing.T) {
	source_code := "int x = 13;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "(int|str|bool)"},
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]+"},
		{Type: "OPERATOR", Regex: "="},
		{Type: "INTEGER", Regex: "[1-9][0-9]*"},
		{Type: "PUNCTUATION", Regex: ";"},
	}
	expected_tokens := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "OPERATOR", Value: "="},
		{Type: "INTEGER", Value: "13"},
		{Type: "PUNCTUATION", Value: ";"},
	}

	tokens, leftovers, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_tokens[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_tokens[i])
			}
		}
		if len(leftovers) != 0 {
			t.Errorf("Tokenisation incorrect: %v", leftovers[0])
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

// ========================= //
// TEST: CreateTokensFromDFA //
// ========================= //

func TestCreateTokensFromDFA_NoSourceCode(t *testing.T) {

	source_code := ""
	dfa := services.Automata{}

	_, _, err := services.CreateTokensFromDFA(source_code, dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("source code is empty").Error() {
			t.Logf("Test Passed: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestCreateTokensFromDFA_NoTransitions(t *testing.T) {

	source_code := "int x=3;"
	dfa := services.Automata{
		States:      []string{"START", "S1", "S2", "S3", "S4", "S5"},
		Transitions: []services.Transition{},
		Start:       "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

	_, _, err := services.CreateTokensFromDFA(source_code, dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no transitions identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestCreateTokensFromDFA_NoStart(t *testing.T) {

	source_code := "int x=3;"
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
		Start: "",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

	_, _, err := services.CreateTokensFromDFA(source_code, dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no start state identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestCreateTokensFromDFA_NoAcceptingStates(t *testing.T) {

	source_code := "int x=3;"
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
		Start:     "START",
		Accepting: []services.AcceptingState{},
	}

	_, _, err := services.CreateTokensFromDFA(source_code, dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no accepting states identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestCreateTokensFromDFA_ValidDFA(t *testing.T) {
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "NUMBER", Value: "3"},
	}
	expected_res_unidentified := []string{
		"=",
		";",
	}

	source_code := "int x = 3;"
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

	tokens, unidentified_tokens, err := services.CreateTokensFromDFA(source_code, dfa)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokensFromDFA_SpacedUnidentifiedTokens(t *testing.T) {
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "NUMBER", Value: "3"},
	}
	expected_res_unidentified := []string{
		"!=",
		";",
	}

	source_code := "int x != 3;"
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

	tokens, unidentified_tokens, err := services.CreateTokensFromDFA(source_code, dfa)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokensFromDFA_ComplexDFAWithSimpleCode(t *testing.T) {
	expected_res := []services.TypeValue{
		{Type: "ID", Value: "int"},
		{Type: "ID", Value: "x"},
		{Type: "NUM", Value: "3"},
	}
	expected_res_unidentified := []string{
		"=",
		";",
	}

	source_code := "int x = 3  ;   "
	dfa := services.Automata{
		States: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
			{From: "A", To: "G", Label: "[0-9]"},
			{From: "B", To: "C", Label: "f"},
			{From: "B", To: "D", Label: "[a-eg-zA-Z_0-9]"},
			{From: "C", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "D", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "E", To: "H", Label: "[0-9]"},
			{From: "F", To: "E", Label: "."},
			{From: "F", To: "G", Label: "[0-9]"},
			{From: "G", To: "H", Label: "[.]"},
			{From: "G", To: "I", Label: "[eE]"},
			{From: "G", To: "G", Label: "[0-9]"},
			{From: "H", To: "I", Label: "[eE]"},
			{From: "H", To: "H", Label: "[0-9]"},
			{From: "I", To: "J", Label: "[+-]"},
			{From: "I", To: "K", Label: "[0-9]"},
			{From: "J", To: "K", Label: "[0-9]"},
			{From: "K", To: "D", Label: "[0-9]"},
		},
		Start: "A",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	tokens, unidentified_tokens, err := services.CreateTokensFromDFA(source_code, dfa)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokensFromDFA_ComplexDFAWithComplexCode(t *testing.T) {
	expected_res := []services.TypeValue{
		{Type: "ID", Value: "int"},
		{Type: "ID", Value: "x"},
		{Type: "NUM", Value: "3"},
		{Type: "IF", Value: "if"},
		{Type: "ID", Value: "x"},
		{Type: "NUM", Value: "5"},
		{Type: "ID", Value: "x"},
		{Type: "FLOAT", Value: "-5.928"},
	}
	expected_res_unidentified := []string{
		"=",
		";",
		"<",
		"{",
		"}",
	}

	source_code := "int x = 3; "
	source_code += "if x < 5 { x = -5.928; }"
	dfa := services.Automata{
		States: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
			{From: "A", To: "G", Label: "[0-9]"},
			{From: "B", To: "C", Label: "f"},
			{From: "B", To: "D", Label: "[a-eg-zA-Z_0-9]"},
			{From: "C", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "D", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "E", To: "H", Label: "[0-9]"},
			{From: "F", To: "E", Label: "."},
			{From: "F", To: "G", Label: "[0-9]"},
			{From: "G", To: "H", Label: "[.]"},
			{From: "G", To: "I", Label: "[eE]"},
			{From: "G", To: "G", Label: "[0-9]"},
			{From: "H", To: "I", Label: "[eE"},
			{From: "H", To: "H", Label: "[0-9]"},
			{From: "I", To: "J", Label: "[+-]"},
			{From: "I", To: "K", Label: "[0-9]"},
			{From: "J", To: "K", Label: "[0-9]"},
			{From: "K", To: "D", Label: "[0-9]"},
		},
		Start: "A",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	tokens, unidentified_tokens, err := services.CreateTokensFromDFA(source_code, dfa)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokensFromDFA_TestSimilarIDandKey(t *testing.T) {
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "integ"},
		{Type: "NUMBER", Value: "2"},
	}
	expected_res_unidentified := []string{
		"=",
		";",
	}

	source_code := "int integ = 2;"
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

	tokens, unidentified_tokens, err := services.CreateTokensFromDFA(source_code, dfa)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v != %v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

// =========================== //
//   TEST: ConvertDFAToRegex   //
// =========================== //

func TestConvertDFAToRegex_NoTransitions(t *testing.T) {

	dfa := services.Automata{
		States:      []string{"START", "S1", "S2", "S3", "S4", "S5"},
		Transitions: []services.Transition{},
		Start:       "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

	_, err := services.ConvertDFAToRegex(dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no transitions identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestConvertDFAToRegex_NoStart(t *testing.T) {

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
		Start: "",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

	_, err := services.ConvertDFAToRegex(dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no start state identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestConvertDFAToRegex_NoStates(t *testing.T) {

	dfa := services.Automata{
		States: []string{},
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

	_, err := services.ConvertDFAToRegex(dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no states identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestConvertDFAToRegex_NoAcceptingStates(t *testing.T) {

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
		Start:     "START",
		Accepting: []services.AcceptingState{},
	}

	_, err := services.ConvertDFAToRegex(dfa)

	if err != nil {
		if err.Error() == fmt.Errorf("no accepting states identified in dfa").Error() {
			t.Logf("Test Passed: %v", err)
		} else {
			t.Errorf("Incorrect error: %v", err)
		}
	} else {
		t.Errorf("Error supposed to occur")
	}
}

func TestConvertDFAToRegex_ValidDFA(t *testing.T) {
	expected_res := []services.TypeRegex{
		{Type: "IDENTIFIER", Regex: "[a-z][a-z0-9]+"},
		{Type: "KEYWORD", Regex: "int|if"},
		{Type: "NUMBER", Regex: "[0-9]+"},
	}

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

	rules, err := services.ConvertDFAToRegex(dfa)

	if err == nil {
		for _, rule := range rules {
			match_found := false
			for _, res := range expected_res {
				if rule == res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Tokenisation incorrect: %v", rule)
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestConvertDFAToRegex_ValidDFANoRegex(t *testing.T) {
	expected_res := []services.TypeRegex{}

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
			{State: "S9", Type: "ID"},
		},
	}

	rules, err := services.ConvertDFAToRegex(dfa)

	if err == nil {
		for _, rule := range rules {
			match_found := false
			for _, res := range expected_res {
				if rule == res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Tokenisation incorrect: %v", rule)
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestConvertDFAToRegex_ValidDFARanges(t *testing.T) {
	expected_res := []services.TypeRegex{
		{Type: "IDENTIFIER", Regex: "[a-z][a-z0-9]+"},
		{Type: "KEYWORD", Regex: "int|if"},
		{Type: "NUMBER", Regex: "[0-9]+"},
	}

	dfa := services.Automata{
		States: []string{"START", "S1", "S2", "S3", "S4", "S5"},
		Transitions: []services.Transition{
			{From: "START", To: "S1", Label: "i"},
			{From: "S1", To: "S5", Label: "n"},
			{From: "S5", To: "S4", Label: "t"},
			{From: "S1", To: "S6", Label: "f"},
			{From: "START", To: "S2", Label: "[0-9]"},
			{From: "S2", To: "S2", Label: "[0-9]"},
			{From: "START", To: "S3", Label: "[a-z]"},
			{From: "S3", To: "S3", Label: "[a-z0-9]"},
		},
		Start: "START",
		Accepting: []services.AcceptingState{
			{State: "S3", Type: "IDENTIFIER"},
			{State: "S4", Type: "KEYWORD"},
			{State: "S2", Type: "NUMBER"},
			{State: "S6", Type: "KEYWORD"},
		},
	}

	rules, err := services.ConvertDFAToRegex(dfa)

	if err == nil {
		for _, rule := range rules {
			match_found := false
			for _, res := range expected_res {
				if rule == res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Tokenisation incorrect: %v", rule)
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestConvertDFAToRegex_Complex(t *testing.T) {
	expected_res := []services.TypeRegex{
		{Type: "ID", Regex: "i"},
		{Type: "ID", Regex: "[a-hj-zA-Z_][a-zA-Z_0-9]+"},
		{Type: "ID", Regex: "i[a-eg-zA-Z_0-9][a-zA-Z_0-9]+"},
		{Type: "ID", Regex: "if[a-zA-Z_0-9][a-zA-Z_0-9]+"},
		{Type: "IF", Regex: "if"},
		{Type: "NUM", Regex: "[+-][0-9][0-9]+"},
		{Type: "FLOAT", Regex: `\.[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[0-9][.][0-9]+`},
		{Type: "FLOAT", Regex: `[+-]\.[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-][0-9][.][0-9]+`},
		{Type: "FLOAT", Regex: `\.[0-9]\[eE[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-][0-9][eE][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[0-9][.]\[eE[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[0-9][eE][+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `\.[0-9]\[eE[+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-]\.[0-9]\[eE[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-][0-9][.]\[eE[0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-][0-9][eE][+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[0-9][.]\[eE[+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-]\.[0-9]\[eE[+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[+-][0-9][.]\[eE[+-][0-9][0-9]+`},
		{Type: "FLOAT", Regex: `[0-9][eE][0-9][0-9]+`},
	}

	dfa := services.Automata{
		States: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
			{From: "A", To: "G", Label: "[0-9]"},
			{From: "B", To: "C", Label: "f"},
			{From: "B", To: "D", Label: "[a-eg-zA-Z_0-9]"},
			{From: "C", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "D", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "E", To: "H", Label: "[0-9]"},
			{From: "F", To: "E", Label: "."},
			{From: "F", To: "G", Label: "[0-9]"},
			{From: "G", To: "H", Label: "[.]"},
			{From: "G", To: "I", Label: "[eE]"},
			{From: "G", To: "G", Label: "[0-9]"},
			{From: "H", To: "I", Label: "[eE"},
			{From: "H", To: "H", Label: "[0-9]"},
			{From: "I", To: "J", Label: "[+-]"},
			{From: "I", To: "K", Label: "[0-9]"},
			{From: "J", To: "K", Label: "[0-9]"},
			{From: "K", To: "K", Label: "[0-9]"},
		},
		Start: "A",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	rules, err := services.ConvertDFAToRegex(dfa)

	if err == nil {
		for _, rule := range rules {
			match_found := false
			for _, res := range expected_res {
				if rule == res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Tokenisation incorrect: %v", rule)
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestSimplifyRegex_Keyword(t *testing.T) {

	regex := "if|else"

	output := services.SimplifyRegex(regex)

	if output[0] != regex {
		t.Errorf("Function not working correctly")
	}
}

func TestSimplifyRegex_RemoveDuplictaes(t *testing.T) {

	expected_res := []string{
		"\\d+",
		"[+-]\\d+",
	}

	regex := "\\d+|[+-]\\d+"

	output := services.SimplifyRegex(regex)

	if output == nil {
		t.Errorf("Function not working correctly")
	} else {
		if len(output) != 2 {
			t.Errorf("Incorrect number of items returned")
		} else {
			for i, item := range output {
				if item != expected_res[i] {
					t.Errorf("Incorrect item: %v", item)
				}
			}
		}
	}
}

func TestConvertRawRegexToRegexRules(t *testing.T) {

	expected_res := "[a-z0-9]+"
	regex := "(abcdefghijklmnopqrstuvwxyz0123456789)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[A-Z0-9]+"
	regex = "(ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[a-z]+"
	regex = "(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[A-Z]+"
	regex = "(ABCDEFGHIJKLMNOPQRSTUVWXYZ)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[0-9]+"
	regex = "(0123456789)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}
}

func TestConvertRawRegexToRegexRules_Success(t *testing.T) {

	expected_res := "[a-z][a-z0-9]+"
	regex := "abcdefghijklmnopqrstuvwxyz(abcdefghijklmnopqrstuvwxyz0123456789)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[a-z0-9]+"
	regex = "abcdefghijklmnopqrstuvwxyz0123456789(abcdefghijklmnopqrstuvwxyz0123456789)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[a-zA-Z0-9][a-z]+"
	regex = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789(abcdefghijklmnopqrstuvwxyz)*"
	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[a-z0-9][a-z]+"
	regex = "abcdefghijklmnopqrstuvwxyz0123456789(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[A-Z0-9][a-z]+"
	regex = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[a-zA-Z][a-z]+"
	regex = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[A-Z][a-z]+"
	regex = "ABCDEFGHIJKLMNOPQRSTUVWXYZ(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}

	expected_res = "[0-9][a-z]+"
	regex = "0123456789(abcdefghijklmnopqrstuvwxyz)*"

	services.ConvertRawRegexToRegexRules(&regex)

	if regex != expected_res {
		t.Errorf("Incorrect conversion: %v", regex)
	}
}

// ========================= //
//	TEST: ConvertRegexToNFA  //
// ========================= //

func TestConvertRegexToNFA_NoRegex(t *testing.T) {
	regexes := map[string]string{}

	_, err := services.ConvertRegexToNFA(regexes)

	if err == nil {
		t.Errorf("Error not returned for empty regex")
	}
}

func TestConvertRegexToNFA_Valid(t *testing.T) {
	expected_nfa := services.Automata{
		Start: "S0",
		Transitions: []services.Transition{
			{From: "S1", To: "S2", Label: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"},
			{From: "S3", To: "S4", Label: `\`},
			{From: "S2", To: "S3", Label: "ε"},
			{From: "S5", To: "S6", Label: "w"},
			{From: "S7", To: "S5", Label: "ε"},
			{From: "S7", To: "S8", Label: "ε"},
			{From: "S6", To: "S5", Label: "ε"},
			{From: "S6", To: "S8", Label: "ε"},
			{From: "S4", To: "S7", Label: "ε"},
			{From: "S0", To: "S1", Label: "ε"},
			{From: "S9", To: "S10", Label: `\`},
			{From: "S11", To: "S12", Label: "d"},
			{From: "S13", To: "S11", Label: "ε"},
			{From: "S12", To: "S11", Label: "ε"},
			{From: "S12", To: "S14", Label: "ε"},
			{From: "S10", To: "S13", Label: "ε"},
			{From: "S15", To: "S16", Label: `\`},
			{From: "S17", To: "S18", Label: "."},
			{From: "S16", To: "S17", Label: "ε"},
			{From: "S19", To: "S20", Label: `\`},
			{From: "S18", To: "S19", Label: "ε"},
			{From: "S21", To: "S22", Label: "d"},
			{From: "S23", To: "S21", Label: "ε"},
			{From: "S22", To: "S21", Label: "ε"},
			{From: "S22", To: "S24", Label: "ε"},
			{From: "S20", To: "S23", Label: "ε"},
			{From: "S14", To: "S15", Label: "ε"},
			{From: "S25", To: "S26", Label: "?"},
			{From: "S24", To: "S25", Label: "ε"},
			{From: "S0", To: "S9", Label: "ε"},
			{From: "S27", To: "S28", Label: "i"},
			{From: "S29", To: "S30", Label: "f"},
			{From: "S28", To: "S29", Label: "ε"},
			{From: "S31", To: "S32", Label: "e"},
			{From: "S33", To: "S34", Label: "l"},
			{From: "S32", To: "S33", Label: "ε"},
			{From: "S35", To: "S36", Label: "s"},
			{From: "S34", To: "S35", Label: "ε"},
			{From: "S37", To: "S38", Label: "e"},
			{From: "S36", To: "S37", Label: "ε"},
			{From: "S39", To: "S27", Label: "ε"},
			{From: "S30", To: "S40", Label: "ε"},
			{From: "S39", To: "S31", Label: "ε"},
			{From: "S38", To: "S40", Label: "ε"},
			{From: "S0", To: "S39", Label: "ε"},
		},
		Accepting: []services.AcceptingState{
			{State: "S8", Type: "IDENTIFIER"},
			{State: "S26", Type: "NUMBER"},
			{State: "S40", Type: "KEYWORD"},
		},
		States: []string{"S24", "S32", "S35", "S40", "S0", "S4", "S11", "S13", "S6", "S8", "S20", "S21", "S29", "S3", " S15", "S23", "S26", "S37", "S38", "S5", "S9", "S19", "S25", "S1", "S28", "S31", "S36", "S33", "S39", "S16", "S17", "S18", "S30", "S10", "S22", "S27", "S34", "S2", "S7", "S12", "S14", "S15"},
	}
	regexes := map[string]string{
		"IDENTIFIER": "[a-zA-Z_]\\w*",
		"NUMBER":     "\\d+(\\.\\d+)?",
		"KEYWORD":    "if|else",
	}

	nfa, err := services.ConvertRegexToNFA(regexes)

	if err != nil {
		t.Errorf("Error not supposed to occurr: %v", err)
	}

	if nfa.Start != expected_nfa.Start {
		t.Errorf("Incorrect start: %v", nfa.Start)
	}
	for _, transition := range nfa.Transitions {
		match_found := false
		for i, res := range expected_nfa.Transitions {
			if transition.Label == res.Label {
				match_found = true
				expected_nfa.Transitions = append(expected_nfa.Transitions[:i], expected_nfa.Transitions[i+1:]...)
				break
			}
		}
		if !match_found {
			t.Errorf("Incorrect transition: %v %v %v", transition.From, transition.To, transition.Label)
		}
	}
	for _, state := range nfa.States {
		match_found := false
		for _, res := range expected_nfa.States {
			if state == res {
				match_found = true
			}
		}
		if !match_found {
			t.Errorf("Incorrect state: %v", state)
		}
	}
	for _, accept := range nfa.Accepting {
		match_found := false
		for _, res := range expected_nfa.Accepting {
			if accept.Type == res.Type {
				match_found = true
			}
		}
		if !match_found {
			t.Errorf("Incorrect Accepting state: %v, %v", accept.State, accept.Type)
		}
	}
}

// ========================= //
//	TEST: ConvertRegexToDFA  //
// ========================= //

func TestConvertRegexToDFA_NoRegex(t *testing.T) {

	regexes := map[string]string{}

	_, err := services.ConvertRegexToDFA(regexes)

	if err == nil {
		t.Errorf("Error not received for empty regex")
	} else {
		if err.Error() != fmt.Errorf("could not convert regex to nfa: no regex specified").Error() {
			t.Errorf("Incorrect error recieved fro empty regex: %v", err)
		}
	}

}

func TestConvertRegexToDFA_ValidRegex(t *testing.T) {

	expected_dfa := services.Automata{
		Start:       "D0",
		Transitions: []services.Transition{},
		Accepting: []services.AcceptingState{
			{State: "D17", Type: "KEYWORD"},
			{State: "D19", Type: "NUMBER"},
			{State: "D3", Type: "IDENTIFIER"},
			{State: "D6", Type: "IDENTIFIER"},
		},
		States: []string{"D0", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10", "D11", "D12", "D13", "D14", "D15", "D16", "D17", "D18", "D19"},
	}

	regexes := map[string]string{
		"IDENTIFIER": "[a-zA-Z_]\\w*",
		"NUMBER":     "\\d+(\\.\\d+)?",
		"KEYWORD":    "\\b(if|else)\\b",
	}

	dfa, err := services.ConvertRegexToDFA(regexes)

	if err != nil {
		t.Errorf("Error received for valid regex")
	} else {
		if dfa.Start != expected_dfa.Start {
			t.Errorf("Incorrect start: %v", dfa.Start)
		}
		for _, state := range dfa.States {
			match_found := false
			for _, res := range expected_dfa.States {
				if state == res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Incorrect state: %v", state)
			}
		}

		if len(dfa.Transitions) < 10 {
			t.Error("Incorrect conversion for transitions")
		}

		for _, accept := range dfa.Accepting {
			match_found := false
			for i, res := range expected_dfa.Accepting {
				if accept.Type == res.Type {
					match_found = true
					expected_dfa.Accepting = append(expected_dfa.Accepting[:i], expected_dfa.Accepting[i+1:]...)
					break
				}
			}
			if !match_found {
				t.Errorf("Incorrect Accepting state: %v, %v", accept.State, accept.Type)
			}
		}
	}

}

// ========================= //
//	 TEST: ConvertNFAToDFA   //
// ========================= //

func TestConvertNFAToDFA_NoStart(t *testing.T) {
	nfa := services.Automata{
		States: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
			{From: "A", To: "G", Label: "[0-9]"},
			{From: "B", To: "C", Label: "f"},
			{From: "B", To: "D", Label: "[a-eg-zA-Z_0-9]"},
			{From: "C", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "D", To: "D", Label: "[a-zA-Z_0-9]"},
		},
		Start: "",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	_, err := services.ConvertNFAToDFA(nfa)

	if err.Error() != fmt.Errorf("no start state identified").Error() {
		t.Errorf("Error not received for no start")
	}
}

func TestConvertNFAToDFA_NoStates(t *testing.T) {
	nfa := services.Automata{
		States: []string{},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
			{From: "A", To: "G", Label: "[0-9]"},
			{From: "B", To: "C", Label: "f"},
			{From: "B", To: "D", Label: "[a-eg-zA-Z_0-9]"},
			{From: "C", To: "D", Label: "[a-zA-Z_0-9]"},
			{From: "D", To: "D", Label: "[a-zA-Z_0-9]"},
		},
		Start: "A",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	_, err := services.ConvertNFAToDFA(nfa)

	if err.Error() != fmt.Errorf("no states identified").Error() {
		t.Errorf("Error not received for no states")
	}
}

func TestConvertNFAToDFA_NoTransisitions(t *testing.T) {
	nfa := services.Automata{
		States:      []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{},
		Start:       "A",
		Accepting: []services.AcceptingState{
			{State: "B", Type: "ID"},
			{State: "C", Type: "IF"},
			{State: "D", Type: "ID"},
			{State: "G", Type: "NUM"},
			{State: "H", Type: "FLOAT"},
			{State: "K", Type: "FLOAT"},
		},
	}

	_, err := services.ConvertNFAToDFA(nfa)

	if err.Error() != fmt.Errorf("no transitions identified").Error() {
		t.Errorf("Error not received for no transitions")
	}
}

func TestConvertNFAToDFA_NoAcceptingStates(t *testing.T) {
	nfa := services.Automata{
		States: []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
		Transitions: []services.Transition{
			{From: "A", To: "B", Label: "i"},
			{From: "A", To: "D", Label: "[a-hj-zA-Z_]"},
			{From: "A", To: "E", Label: "."},
			{From: "A", To: "F", Label: "[+-]"},
		},
		Start:     "A",
		Accepting: []services.AcceptingState{},
	}

	_, err := services.ConvertNFAToDFA(nfa)

	if err.Error() != fmt.Errorf("no accepting states identified").Error() {
		t.Errorf("Error not received for no accepting states")
	}
}

func TestConvertNFAToDFA_ValidNFA(t *testing.T) {
	expected_dfa := services.Automata{
		Start: "D0",
		Transitions: []services.Transition{
			{From: "D0", To: "D1", Label: "a-z0-9"},
			{From: "D0", To: "D2", Label: "i"},
			{From: "D0", To: "D2", Label: "f"},
			{From: "D0", To: "D3", Label: "0-9"},
			{From: "D0", To: "D1", Label: "a-z"},
			{From: "D1", To: "D1", Label: "a-z0-9"},
			{From: "D2", To: "D4", Label: "n"},
			{From: "D3", To: "D3", Label: "0-9"},
			{From: "D4", To: "D5", Label: "t"},
		},
		Accepting: []services.AcceptingState{
			{State: "D0", Type: "IDENTIFIER"},
			{State: "D1", Type: "IDENTIFIER"},
			{State: "D2", Type: "KEYWORD"},
			{State: "D3", Type: "NUMBER"},
			{State: "D5", Type: "KEYWORD"},
		},
		States: []string{"D1", "D2", "D3", "D4", "D5"},
	}

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

	dfa, err := services.ConvertNFAToDFA(nfa)

	if err != nil {
		t.Errorf("Error not expected")
	} else {
		if dfa.Start != expected_dfa.Start {
			t.Errorf("Incorrect start: %v", dfa.Start)
		}
		for _, state := range dfa.States {
			match_found := false
			for _, res := range expected_dfa.States {
				if state != res {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Incorrect state: %v", state)
			}
		}

		if len(dfa.Transitions) < 8 {
			t.Error("Incorrect conversion for transitions")
		}
		for _, transition := range dfa.Transitions {
			match_found := false
			for _, res := range expected_dfa.Transitions {
				if transition.To != res.To && transition.From != res.From && transition.Label != res.Label {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Incorrect transition: %v %v %v", transition.From, transition.To, transition.Label)
			}
		}

		for _, accept := range dfa.Accepting {
			match_found := false
			for _, res := range expected_dfa.Accepting {
				if accept.State != res.State && accept.Type != res.Type {
					match_found = true
				}
			}
			if !match_found {
				t.Errorf("Incorrect Accepting state: %v, %v", accept.State, accept.Type)
			}
		}
	}
}
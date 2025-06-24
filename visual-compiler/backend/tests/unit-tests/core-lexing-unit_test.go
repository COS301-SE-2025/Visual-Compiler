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
	add_invalid_test{[]byte(`[
		{"type": "keyword","regex":"\\b(if|else)\\b},
		{"type: "identifier","regex":"[a-zA-Z_]\\w*},
	]`), fmt.Errorf(`invalid JSON for rules: invalid character '\n' in string literal`)},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b"},{"type": "IDENITIFER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid JSON for rules: invalid character ']' looking for beginning of value")},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b}]`), fmt.Errorf("invalid JSON for rules: unexpected end of JSON input")},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b},{"type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid JSON for rules: invalid character 't' after object key:value pair")},
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
		if err.Error() != fmt.Errorf("no rules specified").Error() {
			t.Errorf("Incorrect error: %v", err)
		} else {
			t.Logf("Error sent for empty source code")
		}
	} else {
		t.Errorf("No error sent for no source code")
	}
}

func TestCreateTokens_UnidentifiedTokensFound(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|else)\\b"},
	}
	expected_res := []string{
		"int",
		"x",
		"=",
		"3",
		";",
	}

	_, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range unidentified_tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v,!=,%v", token, expected_res[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_KeywordTokensUnidentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|else)\\b"},
	}
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
	}
	expected_res_unidentified := []string{
		"int",
		"x",
		"=",
		"3",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_KeywordTokensIdentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|int|else)\\b"},
	}
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
	}
	expected_res_unidentified := []string{
		"x",
		"=",
		"3",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_IdentifierTokensIdentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
	}
	expected_res := []services.TypeValue{
		{Type: "IDENTIFIER", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
	}
	expected_res_unidentified := []string{
		"=",
		"3",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}
func TestCreateTokens_IdentifierTokensIdentified_Complex(t *testing.T) {
	source_code := "int x = 3; int x_1 = 5;"
	rules := []services.TypeRegex{
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
	}
	expected_res := []services.TypeValue{
		{Type: "IDENTIFIER", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "IDENTIFIER", Value: "int"},
		{Type: "IDENTIFIER", Value: "x_1"},
	}
	expected_res_unidentified := []string{
		"=",
		"3",
		";",
		"=",
		"5",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_NumberTokensIdentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "NUMBER", Regex: "\\d+(\\.\\d+)?"},
	}
	expected_res := []services.TypeValue{
		{Type: "NUMBER", Value: "3"},
	}
	expected_res_unidentified := []string{
		"int",
		"x",
		"=",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_OperatorTokensIdentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "OPERATOR", Regex: "="},
	}
	expected_res := []services.TypeValue{
		{Type: "OPERATOR", Value: "="},
	}
	expected_res_unidentified := []string{
		"int",
		"x",
		"3",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_OperatorTokensUnidentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "OPERATOR", Regex: "+"},
	}
	expected_res := []services.TypeValue{}
	expected_res_unidentified := []string{
		"int",
		"x",
		"=",
		"3",
		";",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_PunctuationTokensIdentified(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "PUNCTUATION", Regex: ";|,"},
	}
	expected_res := []services.TypeValue{
		{Type: "PUNCTUATION", Value: ";"},
	}
	expected_res_unidentified := []string{
		"int",
		"x",
		"=",
		"3",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_PunctuationTokensUnidentified(t *testing.T) {
	source_code := "int x = 3; err,_ := service()"
	rules := []services.TypeRegex{
		{Type: "PUNCTUATION", Regex: ";"},
	}
	expected_res := []services.TypeValue{
		{Type: "PUNCTUATION", Value: ";"},
	}
	expected_res_unidentified := []string{
		"int",
		"x",
		"=",
		"3",
		"err",
		",",
		":=",
		"service",
		"()",
	}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

func TestCreateTokens_Complex(t *testing.T) {
	source_code := "int x = 3;"
	rules := []services.TypeRegex{
		{Type: "KEYWORD", Regex: "\\b(if|int|else)\\b"},
		{Type: "IDENTIFIER", Regex: "[a-zA-Z_]\\w*"},
		{Type: "OPERATOR", Regex: "="},
		{Type: "NUMBER", Regex: "\\d+(\\.\\d+)?"},
		{Type: "PUNCTUATION", Regex: ";"},
	}
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "OPERATOR", Value: "="},
		{Type: "NUMBER", Value: "3"},
		{Type: "PUNCTUATION", Value: ";"},
	}
	expected_res_unidentified := []string{}

	tokens, unidentified_tokens, err := services.CreateTokens(source_code, rules)

	if err == nil {
		for i, token := range tokens {
			if token != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
			}
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
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
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

	source_code := "int x = 3;"
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
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
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
		"=",
		";",
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
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
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
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res[i])
			}
		}
		for i, token := range unidentified_tokens {
			if token != expected_res_unidentified[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", token, expected_res_unidentified[i])
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
		{Type: "IDENTIFIER", Regex: `[a-zA-Z_]\\w*`},
		{Type: "KEYWORD", Regex: `\\b(int|if)\\b`},
		{Type: "NUMBER", Regex: `\\d+(\\.\\d+)?`},
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
		for i, rule := range rules {
			if rule != expected_res[i] {
				t.Errorf("Tokenisation incorrect: %v, != ,%v", rule, expected_res[i])
			}
		}
	} else {
		t.Errorf("Error not supposed to occur")
	}
}

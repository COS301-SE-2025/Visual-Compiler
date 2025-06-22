package unit_tests

import (
	"fmt"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestInitialise(t *testing.T) {
	lexing_input := "int x = 2 ;"
	services.SourceCode(lexing_input)
	source := services.GetSourceCode()
	if source != lexing_input {
		t.Errorf("Source code initialisation faied")
	}
}

type add_invalid_test struct {
	arg1     []byte
	expected error
}

var add_invalid_tests = []add_invalid_test{
	add_invalid_test{[]byte(`[
		{"type": "keyword","regex":"\\b(if|else)\\b},
		{"type: "identifier","regex":"[a-zA-Z_]\\w*},
	]`), fmt.Errorf("invalid character")},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b"},{"type": "IDENITIFER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid character")},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b}]`), fmt.Errorf("invalid character")},
	add_invalid_test{[]byte(`[{"type": "KEYWORD","regex":"\\b(if|else)\\b},{"type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"},]`), fmt.Errorf("invalid character")},
}

func TestReadRegexRules_InvalidCharacter(t *testing.T) {
	for _, test := range add_invalid_tests {
		err := services.ReadRegexRules(test.arg1)
		if err == nil {
			t.Errorf("Failed for invalid JSON object: %v", err)
		} else {
			t.Logf("Passed test for: %v", err)
		}
	}
}

func TestReadRegexRules_Fail(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD""},{"Type": "IDENTIFIER","regex":"[a-zA-Z_]\\w*"}]`)
	err := services.ReadRegexRules(c_input)
	if err == nil {
		t.Errorf("Passed for invalid input")
	}
	t.Logf("Failed for invalid input")
}

func TestReadRegexRules_Success(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD","Regex":"\\b(if|else)\\b"},{"Type": "IDENTIFIER","Regex":"[a-zA-Z_]\\w*"}]`)
	err := services.ReadRegexRules(c_input)
	if err != nil {
		t.Errorf("Failed for valid input")
	}
}

func TestCreateTokens_UnexpectedTokens(t *testing.T) {
	services.CreateTokens()
	output := services.GetTokens()
	if output == nil {
		t.Errorf("No output")
	}
}

func TestCreateTokens_NumberTokens(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD","Regex":"\\b(if|else)\\b"},{"Type": "IDENTIFIER","Regex":"[a-zA-Z_]\\w*"},{"Type":"NUMBER","Regex":"\\d+(\\.\\d+)?"}]`)
	err := services.ReadRegexRules(c_input)
	if err != nil {
		t.Errorf("Failed for valid input: %v", err)
	}
	services.CreateTokens()
	output := services.GetTokens()
	if output == nil {
		t.Errorf("No output")
	}
}

func TestCreateTokens_OperatorTokens(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD","Regex":"\\b(if|else)\\b"},{"Type": "IDENTIFIER","Regex":"[a-zA-Z_]\\w*"},{"Type":"NUMBER","Regex":"\\d+(\\.\\d+)?"},{"Type":"OPERATOR","Regex":"="}]`)
	err := services.ReadRegexRules(c_input)
	if err != nil {
		t.Errorf("Failed for valid input: %v", err)
	}
	services.CreateTokens()
	output := services.GetTokens()
	if output == nil {
		t.Errorf("No output")
	}
}

func TestUnexpectedTokens_Punctuation(t *testing.T) {
	unexpected_tokens := services.GetInvalidInput()
	expectedRes := []string{
		";",
	}
	if len(unexpected_tokens) != len(expectedRes) {
		t.Errorf("Not enough unexpected tokens found")
	}

	for i, token := range unexpected_tokens {
		if token != expectedRes[i] {
			t.Errorf("Tokenisation incorrect: %v", token)
		}
	}
}

func TestCreateTokens_PunctuationTokens(t *testing.T) {
	c_input := []byte(`[{"Type": "KEYWORD","Regex":"\\b(if|else|int)\\b"},{"Type": "IDENTIFIER","Regex":"[a-zA-Z_]\\w*"},{"Type":"NUMBER","Regex":"\\d+(\\.\\d+)?"},{"Type":"OPERATOR","Regex":"="},{"Type":"PUNCTUATION","Regex":";"}]`)
	err := services.ReadRegexRules(c_input)
	if err != nil {
		t.Errorf("Failed for valid input: %v", err)
	}
	services.CreateTokens()
	output := services.GetTokens()
	if output == nil {
		t.Errorf("No output")
	}
}

func TestFullCreateTokens_Success(t *testing.T) {
	services.CreateTokens()
	output := services.GetTokens()
	expectedRes := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "x"},
		{Type: "OPERATOR", Value: "="},
		{Type: "NUMBER", Value: "2"},
		{Type: "PUNCTUATION", Value: ";"},
	}
	if len(output) != len(expectedRes) {
		t.Errorf("Not enough tokens created")
	}

	for i, token := range output {
		if token != expectedRes[i] {
			t.Errorf("Tokenisation incorrect: %v", token)
		}
	}
}

func TestReadDFA_Success(t *testing.T) {
	user_input := []byte(`{
							"states": ["START","S1","S2","S3","S4","S5"],
							"transitions":[
								{"from": "START", "to": "S1", "label": "i"},
								{"from": "S1", "to": "S5", "label": "n"},
								{"from": "S5", "to": "S4", "label": "t"},
								{"from": "START", "to": "S2", "label": "0123456789"},
								{"from": "S2", "to": "S2", "label": "0123456789"},
								{"from": "START", "to": "S3", "label": "abcdefghijklmnopqrstuvwxyz"},
								{"from": "S3", "to": "S3", "label": "abcdefghijklmnopqrstuvwxyz0123456789"}
							],
							"start_state": "START",
							"accepting_states":[
								{"state":"S3","token_type":"IDENTIFIER"},
								{"state":"S4","token_type":"KEYWORD"},
								{"state":"S2","token_type":"NUMBER"}
							]
						}`,
	)
	err := services.ReadDFA(user_input)
	if err != nil {
		t.Errorf("Failed for valid input: %v", err)
	}
}

func TestConvertTokensFromDFA_Success(t *testing.T) {
	lexing_input := "int y = 2;"
	services.SourceCode(lexing_input)

	services.CreateTokensFromDFA()
	tokens := services.GetTokens()
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "y"},
		{Type: "NUMBER", Value: "2"},
	}
	if len(tokens) != len(expected_res) {
		t.Errorf("Not enough tokens created")
	}
	for i:=0;i<len(tokens);i++ {
		if (tokens[i] != expected_res[i]) {
			t.Errorf("Tokenisation incorrect")
		}
	}
	unidentified_tokens := services.GetInvalidInput()
	expected_res2 := []string{"=",";"}
	for i:=0;i<len(unidentified_tokens);i++ {
		if (unidentified_tokens[i] != expected_res2[i]) {
			t.Errorf("Tokenisation incorrect")
		}
	}
}

func TestConvertTokensFromDFA_TestSimilarIDandKey(t *testing.T) {
	lexing_input := "int integ = 2;"
	services.SourceCode(lexing_input)

	services.CreateTokensFromDFA()
	tokens := services.GetTokens()
	expected_res := []services.TypeValue{
		{Type: "KEYWORD", Value: "int"},
		{Type: "IDENTIFIER", Value: "integ"},
		{Type: "NUMBER", Value: "2"},
	}
	if len(tokens) != len(expected_res) {
		t.Errorf("Not enough tokens created")
	}
	for i:=0;i<len(tokens);i++ {
		if (tokens[i] != expected_res[i]) {
			t.Errorf("Tokenisation incorrect")
		}
	}
	unidentified_tokens := services.GetInvalidInput()
	expected_res2 := []string{"=",";"}
	for i:=0;i<len(unidentified_tokens);i++ {
		if (unidentified_tokens[i] != expected_res2[i]) {
			t.Errorf("Tokenisation incorrect")
		}
	}
}

func TestConvertDFAToRegex_Success(t *testing.T) {
	user_input := []byte(`{
							"states": ["START","S1","S2","S3","S4","S5"],
							"transitions":[
								{"from": "START", "to": "S1", "label": "i"},
								{"from": "S1", "to": "S5", "label": "n"},
								{"from": "S5", "to": "S4", "label": "t"},
								{"from": "START", "to": "S2", "label": "0123456789"},
								{"from": "S2", "to": "S2", "label": "0123456789"},
								{"from": "START", "to": "S3", "label": "abcdefghijklmnopqrstuvwxyz"},
								{"from": "S3", "to": "S3", "label": "abcdefghijklmnopqrstuvwxyz0123456789"}
							],
							"start_state": "START",
							"accepting_states":[
								{"state":"S3","token_type":"IDENTIFIER"},
								{"state":"S4","token_type":"KEYWORD"},
								{"state":"S2","token_type":"NUMBER"}
							]
						}`,
	)
	err := services.ReadDFA(user_input)
	if err != nil {
		t.Errorf("Test failed: %v", err)
	}

	err = services.ConvertDFAToRegex()
	if err!=nil {
		t.Errorf("Test failed: %v", err)
	}

}
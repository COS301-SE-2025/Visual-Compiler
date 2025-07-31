package unit_tests

import (
	"fmt"
	"testing"
	
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestReadTranslationRules_NoRules(t *testing.T) {
	input := []byte{}
	_, err := services.ReadTranslationRules(input)

	if err == nil {
		t.Errorf("Error expected for no translation rules")
	} else {
		if err.Error() != fmt.Errorf("no translation rules entered").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadTranslationRules_InvalidJson(t *testing.T) {
	input := []byte("Invalid Json")
	_, err := services.ReadTranslationRules(input)

	if err == nil {
		t.Errorf("Error expected for invalid json")
	} else {
		if err.Error() != fmt.Errorf("invalid JSON for translation rules: invalid character 'I' looking for beginning of value").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadTranslationRules_InvalidToken(t *testing.T) {
	input := []byte(`[
		{
			"sequence": ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER"],
			"translation": ["{INVALID}"]
		}
	]`)
	_, err := services.ReadTranslationRules(input)

	if err == nil {
		t.Errorf("Error expected for invalid token in translation")
	} else {
		expected := "token {INVALID} in translation rule not found in sequence: [KEYWORD IDENTIFIER ASSIGNMENT INTEGER]"
		if err.Error() != expected {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestReadTranslationRules_Valid(t *testing.T) {
	input := []byte(`[
		{
			"sequence": ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER"],
			"translation": ["{KEYWORD} {IDENTIFIER} = {INTEGER};"]
		},
		{
			"sequence": ["KEYWORD", "IDENTIFIER", "OPERATOR", "BOOLEAN", "BRACKETS"],
			"translation": ["if ({IDENTIFIER} {OPERATOR} {BOOLEAN}) then"]
		}
	]`)
	
	rules, err := services.ReadTranslationRules(input)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	} else {
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules but received %d", len(rules))
		}
		
		if len(rules[0].Sequence) != 4 {
			t.Errorf("Expected 4 tokens in first sequence but received %d", len(rules[0].Sequence))
		}
		
		if rules[0].Sequence[0] != "KEYWORD" {
			t.Errorf("Expected first token to be KEYWORD but received %s", rules[0].Sequence[0])
		}
		
		if len(rules[0].Translation) != 1 {
			t.Errorf("Expected 1 line in first translation but received %d", len(rules[0].Translation))
		}
	}
}
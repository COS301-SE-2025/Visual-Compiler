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
		
		if len(rules[0].Sequence) != 5 {
			t.Errorf("Expected 5 tokens in first sequence but received %d", len(rules[0].Sequence))
		}
		
		if rules[0].Sequence[0] != "KEYWORD" {
			t.Errorf("Expected first token to be KEYWORD but received %s", rules[0].Sequence[0])
		}
		
		if len(rules[0].Translation) != 1 {
			t.Errorf("Expected 1 line in first translation but received %d", len(rules[0].Translation))
		}
	}
}

func TestTranslate_EmptyTree(t *testing.T) {
	tree := services.SyntaxTree{Root: nil}
	rules := []services.TranslationRule{
		{
			Sequence:    []string{"KEYWORD"},
			Translation: []string{"{KEYWORD}"},
		},
	}

	_, err := services.Translate(tree, rules)

	if err == nil {
		t.Errorf("Error expected for empty syntax tree")
	} else {
		if err.Error() != fmt.Errorf("empty syntax tree").Error() {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestTranslate_UnmatchedToken(t *testing.T) {
	leaf1 := &services.TreeNode{Symbol: "KEYWORD", Value: "int", Children: nil}
	leaf2 := &services.TreeNode{Symbol: "UNKNOWN", Value: "unknown", Children: nil}
	
	root := &services.TreeNode{
		Symbol:   "ROOT",
		Value:    "",
		Children: []*services.TreeNode{leaf1, leaf2},
	}
	tree := services.SyntaxTree{Root: root}
	
	rules := []services.TranslationRule{
		{
			Sequence:    []string{"KEYWORD"},
			Translation: []string{"{KEYWORD}"},
		},
	}

	_, err := services.Translate(tree, rules)

	if err == nil {
		t.Errorf("Error expected for untranslated token")
	} else {
		expected := "the token (UNKNOWN: unknown) was not part of any translation"
		if err.Error() != expected {
			t.Errorf("Incorrect error received: %v", err)
		}
	}
}

func TestTranslate_Valid(t *testing.T) {
	leaf1 := &services.TreeNode{Symbol: "KEYWORD", Value: "int", Children: nil}
	leaf2 := &services.TreeNode{Symbol: "IDENTIFIER", Value: "red", Children: nil}
	leaf3 := &services.TreeNode{Symbol: "ASSIGNMENT", Value: "=", Children: nil}
	leaf4 := &services.TreeNode{Symbol: "INTEGER", Value: "13", Children: nil}
	leaf5 := &services.TreeNode{Symbol: "SEPARATOR", Value: ";", Children: nil}
	
	root := &services.TreeNode{
		Symbol:   "ROOT",
		Value:    "",
		Children: []*services.TreeNode{leaf1, leaf2, leaf3, leaf4, leaf5},
	}
	tree := services.SyntaxTree{Root: root}
	
	rules := []services.TranslationRule{
		{
			Sequence:    []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "SEPARATOR"},
			Translation: []string{"add \t rax, {INTEGER}", "mov \t [{IDENTIFIER}], rax"},
		},
	}

	result, err := services.Translate(tree, rules)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	} else {
		if len(result) != 2 {
			t.Errorf("Expected 2 lines of target code but received %d", len(result))
		}
		
		expected := []string{"add \t rax, 13", "mov \t [red], rax"}
		if result[0] != expected[0] || result[1] != expected[1] {
			t.Errorf("Expected '%s' but received '%s'", expected, result)
		}
	}
}

func TestTranslate_Repeated(t *testing.T) {
	leaf1 := &services.TreeNode{Symbol: "KEYWORD", Value: "int", Children: nil}
	leaf2 := &services.TreeNode{Symbol: "IDENTIFIER", Value: "grey", Children: nil}
	leaf3 := &services.TreeNode{Symbol: "ASSIGNMENT", Value: "=", Children: nil}
	leaf4 := &services.TreeNode{Symbol: "IDENTIFIER", Value: "black", Children: nil}
	leaf5 := &services.TreeNode{Symbol: "OPERATOR", Value: "+", Children: nil}
	leaf6 := &services.TreeNode{Symbol: "IDENTIFIER", Value: "white", Children: nil}
	leaf7 := &services.TreeNode{Symbol: "SEPARATOR", Value: ";", Children: nil}
	
	root := &services.TreeNode{
		Symbol:   "ROOT",
		Value:    "",
		Children: []*services.TreeNode{leaf1, leaf2, leaf3, leaf4, leaf5, leaf6, leaf7},
	}
	tree := services.SyntaxTree{Root: root}
	
	rules := []services.TranslationRule{
		{
			Sequence:    []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "IDENTIFIER", "OPERATOR", "IDENTIFIER", "SEPARATOR"},
			Translation: []string{"mov \t rax, [{IDENTIFIER}]", "add \t rax, [{IDENTIFIER}]", "add \t rax, [{IDENTIFIER}]"},
		},
	}

	result, err := services.Translate(tree, rules)

	if err != nil {
		t.Errorf("Error not supposed to occur: %v", err)
	} else {
		if len(result) != 3 {
			t.Errorf("Expected 3 lines of target code but received %d", len(result))
		}
		
		expected := []string{"mov \t rax, [grey]", "add \t rax, [black]", "add \t rax, [white]"}
		for i, line := range result {
			if line != expected[i] {
				t.Errorf("Expected '%s' at line %d but received '%s'", expected[i], i, line)
			}
		}
	}
}

func TestLeafNodes_NoNodes(t *testing.T) {
	result := services.LeafNodes(nil)
	
	if result != nil {
		t.Errorf("Expected no leaf nodes but received %v", result)
	}
}

func TestLeafNodes_SingleLeaf(t *testing.T) {
	leaf := &services.TreeNode{Symbol: "KEYWORD", Value: "return", Children: nil}
	
	result := services.LeafNodes(leaf)
	
	if len(result) != 1 {
		t.Errorf("Expected 1 leaf node but received %d", len(result))
	}
	
	if result[0].Symbol != "KEYWORD" {
		t.Errorf("Expected 'KEYWORD' for leaf symbol but received %s", result[0].Symbol)
	}
	
	if result[0].Value != "return" {
		t.Errorf("Expected 'return' for leaf value but received %s", result[0].Value)
	}
}

func TestLeafNodes_MultipleLeaves(t *testing.T) {
	leaf1 := &services.TreeNode{Symbol: "INTEGER", Value: "1", Children: nil}
	leaf2 := &services.TreeNode{Symbol: "OPERATOR", Value: "+", Children: nil}
	leaf3 := &services.TreeNode{Symbol: "INTEGER", Value: "2", Children: nil}
	leaf4 := &services.TreeNode{Symbol: "OPERATOR", Value: "*", Children: nil}
	leaf5 := &services.TreeNode{Symbol: "INTEGER", Value: "3", Children: nil}
	
	branch := &services.TreeNode{
		Symbol:   "BODMAS",
		Value:    "",
		Children: []*services.TreeNode{leaf3, leaf4, leaf5},
	}
	
	root := &services.TreeNode{
		Symbol:   "ROOT",
		Value:    "",
		Children: []*services.TreeNode{leaf1, leaf2, branch},
	}
	
	result := services.LeafNodes(root)
	
	if len(result) != 5 {
		t.Errorf("Expected 5 leaf nodes but received %d", len(result))
	}
	
	expected := []string{"INTEGER", "OPERATOR", "INTEGER", "OPERATOR", "INTEGER"}
	for i, node := range result {
		if node.Symbol != expected[i] {
			t.Errorf("Expected %s at position %d but received %s", expected[i], i, node.Symbol)
		}
	}
}

func TestMatchesSequence_ExactMatch(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "KEYWORD", Value: "int", Children: nil},
		{Symbol: "IDENTIFIER", Value: "blue", Children: nil},
		{Symbol: "ASSIGNMENT", Value: "=", Children: nil},
		{Symbol: "INTEGER", Value: "13", Children: nil},
		{Symbol: "SEPARATOR", Value: ";", Children: nil},
	}
	
	sequence := []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "SEPARATOR"}
	
	result := services.MatchesSequence(leaves, sequence, 0)
	
	if !result {
		t.Errorf("Expected true for exact match")
	}
}

func TestMatchesSequence_NoMatch(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "KEYWORD", Value: "int", Children: nil},
		{Symbol: "IDENTIFIER", Value: "blue", Children: nil},
		{Symbol: "ASSIGNMENT", Value: "=", Children: nil},
		{Symbol: "INTEGER", Value: "13", Children: nil},
		{Symbol: "SEPARATOR", Value: ";", Children: nil},
	}
	
	sequence := []string{"KEYWORD", "IDENTIFIER", "OPERATOR", "SEPARATOR"}
	
	result := services.MatchesSequence(leaves, sequence, 0)
	
	if result {
		t.Errorf("Expected false for no match")
	}
}

func TestMatchesSequence_OutOfBounds(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "KEYWORD", Value: "for", Children: nil},
		{Symbol: "IDENTIFIER", Value: "x", Children: nil},
		{Symbol: "KEYWORD", Value: "range", Children: nil},
		{Symbol: "INTEGER", Value: "21", Children: nil},
	}
	
	sequence := []string{"KEYWORD", "IDENTIFIER", "KEYWORD", "INTEGER", "SEPARATOR"}
	
	result := services.MatchesSequence(leaves, sequence, 0)
	
	if result {
		t.Errorf("Expected false for out of bounds sequence")
	}
}

func TestMatchesSequence_Offset(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "KEYWORD", Value: "for", Children: nil},
		{Symbol: "IDENTIFIER", Value: "x", Children: nil},
		{Symbol: "KEYWORD", Value: "range", Children: nil},
		{Symbol: "INTEGER", Value: "21", Children: nil},
	}
	
	sequence := []string{"KEYWORD", "INTEGER"}
	
	result := services.MatchesSequence(leaves, sequence, 2)
	
	if !result {
		t.Errorf("Expected true for match with offset")
	}
}

func TestUseRule_SingleResult(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "KEYWORD", Value: "print", Children: nil},
		{Symbol: "INTEGER", Value: "13", Children: nil},
	}
	
	sequence := []string{"KEYWORD", "INTEGER"}
	translation := []string{"{KEYWORD} \"{INTEGER}\""}
	
	result := services.UseRule(leaves, sequence, translation)
	
	if len(result) != 1 {
		t.Errorf("Expected 1 result but received %d", len(result))
	}
	
	expected := "print \"13\""
	if result[0] != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, result[0])
	}
}

func TestUseRule_MultipleResult(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "IDENTIFIER", Value: "i", Children: nil},
		{Symbol: "OPERATOR", Value: "++", Children: nil},
		{Symbol: "SEPARATOR", Value: ";", Children: nil},
	}
	
	sequence := []string{"IDENTIFIER", "OPERATOR", "SEPARATOR"}
	translation := []string{
		"{",
		"\t inc [{IDENTIFIER}]",
		"}",
	}
	
	result := services.UseRule(leaves, sequence, translation)
	
	if len(result) != 3 {
		t.Errorf("Expected 3 results but received %d", len(result))
	}
	
	expected := []string{
		"{",
		"\t inc [i]",
		"}",
	}
	
	for i, line := range result {
		if line != expected[i] {
			t.Errorf("Expected '%s' at line %d but received '%s'", expected[i], i, line)
		}
	}
}

func TestUseRule_RepeatedToken(t *testing.T) {
	leaves := []*services.TreeNode{
		{Symbol: "IDENTIFIER", Value: "new_blue", Children: nil},
		{Symbol: "ASSIGNMENT", Value: "=", Children: nil},
		{Symbol: "IDENTIFIER", Value: "old_blue", Children: nil},
		{Symbol: "SEPERATOR", Value: ";", Children: nil},
	}
	
	sequence := []string{"IDENTIFIER", "ASSIGNMENT", "IDENTIFIER", "SEPARATOR"}
	translation := []string{"mov [{IDENTIFIER}], [{IDENTIFIER}]"}
	
	result := services.UseRule(leaves, sequence, translation)
	
	if len(result) != 1 {
		t.Errorf("Expected 1 result but received %d", len(result))
	}
	
	expected := "mov [new_blue], [old_blue]"
	if result[0] != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, result[0])
	}
}

func TestSubstituteTemplate_NoPlaceholders(t *testing.T) {
	template := "Hello World"
	token_map := make(map[string][]*services.TokenTracker)
	
	result := services.SubstituteTemplate(template, token_map)
	
	if result != template {
		t.Errorf("Expected '%s' but received '%s'", template, result)
	}
}

func TestSubstituteTemplate_SinglePlaceholder(t *testing.T) {
	template := "print [{IDENTIFIER}]"
	token_map := map[string][]*services.TokenTracker{
		"IDENTIFIER": {{Value: "blue", Avail: true}},
	}
	
	result := services.SubstituteTemplate(template, token_map)
	
	expected := "print [blue]"
	if result != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, result)
	}
}

func TestSubstituteTemplate_MultiplePlaceholders(t *testing.T) {
	template := "print [{IDENTIFIER}] for {INTEGER}"
	token_map := map[string][]*services.TokenTracker{
		"IDENTIFIER":  {{Value: "blue", Avail: true}},
		"INTEGER":  {{Value: "13", Avail: true}},
	}
	
	result := services.SubstituteTemplate(template, token_map)
	
	expected := "print [blue] for 13"
	if result != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, result)
	}
}

func TestSubstituteTemplate_RepeatedPlaceholders(t *testing.T) {
	template := "[{IDENTIFIER}] and [{IDENTIFIER}] and [{IDENTIFIER}]"
	token_map := map[string][]*services.TokenTracker{
		"IDENTIFIER": {
			{Value: "red", Avail: true},
			{Value: "green", Avail: true},
			{Value: "blue", Avail: true},
		},
	}
	
	result := services.SubstituteTemplate(template, token_map)
	
	expected := "[red] and [green] and [blue]"
	if result != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, result)
	}
}
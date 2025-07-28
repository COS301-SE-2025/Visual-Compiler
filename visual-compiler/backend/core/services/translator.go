package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Struct for the translation rules
type TranslationRule struct {
	Sequence    []string
	Translation []string
}

// Struct to track whether the syntax tree leaves are translated
type TokenTracker struct {
	Value string
	Usage bool
}

// Name: ReadTranslationRules
//
// Parameters: []byte
//
// Return: Translator, error
//
// Receive an array of translation rules and return them in the struct
func ReadTranslationRules(input []byte) ([]TranslationRule, error) {

	translator := []TranslationRule{}

	if len(input) == 0 {
		return []TranslationRule{}, fmt.Errorf("no translation rules entered")
	}

	err := json.Unmarshal(input, &translator)

	if err != nil {
		return []TranslationRule{}, fmt.Errorf("invalid JSON for translation rules: %v", err)
	}

	delimiter := regexp.MustCompile(`[,\s]+`)

	for i, rule := range translator {

		sequence := []string{}

		for _, entry := range rule.Sequence {
			atoms := delimiter.Split(entry, -1)
			sequence = append(sequence, atoms...)
		}

		translator[i].Sequence = sequence

		token := regexp.MustCompile(`\{[^}]+\}`)

		for _, entry := range rule.Translation {

			placeholders := token.FindAllString(entry, -1)

			for _, ph := range placeholders {

				label := ph[1 : len(ph)-1]
				found := false

				for _, token := range sequence {
					if token == label {
						found = true
						break
					}
				}

				if !found {
					return []TranslationRule{}, fmt.Errorf("token {%s} in translation rule not found in sequence: %v", label, sequence)
				}
			}
		}
	}

	return translator, nil
}

// Name: LeafNodes
//
// Parameters: *TreeNode
//
// Return: []*TreeNode
//
// Combines all the leaf nodes into an array in left-to-right order
func LeafNodes(node *TreeNode) []*TreeNode {

	if node == nil {
		return nil
	}

	if len(node.Children) == 0 {
		return []*TreeNode{node}
	}

	var leaf_nodes []*TreeNode

	for _, child := range node.Children {
		branch := LeafNodes(child)
		leaf_nodes = append(leaf_nodes, branch...)
	}

	return leaf_nodes
}
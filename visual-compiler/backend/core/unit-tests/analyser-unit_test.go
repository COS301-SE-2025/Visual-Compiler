package unit_tests

import (
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestCreateSymbolTable_Valid(t *testing.T) {

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "STATEMENT",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "DECLARATION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "TYPE",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "KEYWORD",
									Value:  "int",
								},
							},
						},
						{
							Symbol: "IDENTIFIER",
							Value:  "blue",
						},
						{
							Symbol: "ASSIGNMENT",
							Value:  "=",
						},
						{
							Symbol: "EXPRESSION",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "TERM",
									Value:  "",
									Children: []*services.TreeNode{
										{
											Symbol: "INTEGER",
											Value:  "13",
										},
									},
								},
								{
									Symbol: "OPERATOR",
									Value:  "+",
								},
								{
									Symbol: "TERM",
									Value:  "",
									Children: []*services.TreeNode{
										{
											Symbol: "INTEGER",
											Value:  "89",
										},
									},
								},
							},
						},
					},
				},
				{
					Symbol: "SEPARATOR",
					Value:  ";",
				},
			},
		},
	}

	symbol_table, err := services.CreateSymbolTable(scope_rules, syntax_tree)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		for _, scope := range symbol_table.SymbolScopes {
			symbol, err := scope["blue"]
			if err != true {
				t.Errorf("Error: Blue not found")
			}
			if symbol.Name != "blue" && symbol.Scope != 0 && symbol.Type != "int" {
				t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
			}
		}

	}

}

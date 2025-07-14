package unit_tests

import (
	"fmt"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestPerformScopeCheck_Valid(t *testing.T) {

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

	symbol_table, err := services.PerfromScopeCheck(scope_rules, syntax_tree)

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

func TestPerformScopeCheck_SameSymbolNames_DifferentScope(t *testing.T) {

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "MAIN",
			Value:  "",
			Children: []*services.TreeNode{
				{
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
				{
					Symbol: "FUNCTION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "STATEMENT",
							Value:  "{",
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
									Value:  "}",
								},
							},
						},
					},
				},
			},
		},
	}

	symbol_table, err := services.PerfromScopeCheck(scope_rules, syntax_tree)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		for _, scope := range symbol_table.SymbolScopes {
			symbol, err := scope["blue"]
			if err != true {
				t.Errorf("Error: blue not found")
			}
			if symbol.Name != "blue" && symbol.Scope != 0 && symbol.Type != "int" {
				t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
			}
		}

	}

}

func TestPerformScopeCheck_UndeclaredSymbol(t *testing.T) {

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "STATEMENT",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "DEFINE",
					Value:  "",
					Children: []*services.TreeNode{
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

	_, err := services.PerfromScopeCheck(scope_rules, syntax_tree)

	if err == nil {
		t.Errorf("Error expected for undeclared variable")
	} else {
		if err.Error() != fmt.Errorf("variable not declared within it's scope: blue").Error() {
			t.Errorf("Incorrect error: %v", err)
		}
	}

}

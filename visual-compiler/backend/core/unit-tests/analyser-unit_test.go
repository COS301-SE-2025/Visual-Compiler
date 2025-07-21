package unit_tests

import (
	"fmt"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestCreateEmptySymbolTable(t *testing.T) {

	symbol_table := services.CreateEmptySymbolTable()

	if symbol_table == nil {
		t.Errorf("empty symbol table creation failed")
	} else {
		if symbol_table.SymbolScopes != nil {
			if len(symbol_table.SymbolScopes) != 1 {
				t.Errorf("empty symbol table creation failed")
			} else {
				if symbol_table.SymbolScopes[0] == nil {
					t.Errorf("empty symbol table creation failed")
				}
			}
		} else {
			t.Errorf("empty symbol table creation failed")
		}
	}

}

func TestCreateEmptySymbolTableArtefact(t *testing.T) {

	symbol_table_artefact := services.CreateEmptySymbolTableArtefact()

	if symbol_table_artefact == nil {
		t.Errorf("empty symbol table creation artefact failed")
	} else {
		if symbol_table_artefact.SymbolScopes == nil {
			t.Errorf("empty symbol table artefact creation failed")
		}
	}

}

func TestBindSymbol_SymbolExists(t *testing.T) {

	symbol := services.Symbol{
		Name: "blue", Type: "int", Scope: 0,
	}

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	err := services.BindSymbol(symbol_table, symbol)

	if err == nil {
		t.Errorf("error expected")
	} else {
		if err.Error() != fmt.Errorf("symbol already declared in scope: blue").Error() {
			t.Errorf("incorrect error returned: %v", err)
		}
	}
}

func TestBindSymbol_NoSymbolExists(t *testing.T) {
	symbol := services.Symbol{
		Name: "red", Type: "int", Scope: 0,
	}

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	err := services.BindSymbol(symbol_table, symbol)

	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestLookUpName_NameFound(t *testing.T) {

	expected_res := services.Symbol{
		Name: "blue", Type: "int", Scope: 0,
	}

	symbol_name := "blue"

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	symbol, err := services.LookupName(symbol_table, symbol_name)

	if err != nil {
		t.Errorf("error: %v", err)
	} else {
		if symbol.Name != expected_res.Name || symbol.Type != expected_res.Type || symbol.Scope != expected_res.Scope {
			t.Errorf("Incorrect symbol found: %v %v %v", symbol.Name, symbol.Type, symbol.Scope)
		}
	}
}

func TestLookUpName_NameNotFound(t *testing.T) {

	symbol_name := "red"

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	_, err := services.LookupName(symbol_table, symbol_name)

	if err == nil {
		t.Errorf("error expected")
	} else {
		if err.Error() != fmt.Errorf("symbol not found").Error() {
			t.Errorf("incorrect error: %v", err)
		}
	}
}

func TestEnterNewScope(t *testing.T) {
	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	services.EnterNewScope(symbol_table)

	if len(symbol_table.SymbolScopes) != 2 {
		t.Errorf("Incorrect table length")
	}
}

func TestExitScope_Valid(t *testing.T) {

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
				"red": {
					Name:  "red",
					Type:  "int",
					Scope: 0,
				},
			},
			{
				"red": {
					Name:  "red",
					Type:  "int",
					Scope: 1,
				},
			},
		},
	}

	err := services.ExitScope(symbol_table)

	if err != nil {
		t.Errorf("%v", err)
	}
	if len(symbol_table.SymbolScopes) != 1 {
		t.Errorf("Incorrect table length")
	}
}

func TestExitScope_GlobalScope(t *testing.T) {
	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{},
		},
	}

	err := services.ExitScope(symbol_table)

	if err == nil {
		t.Errorf("error expected")
	}
	if len(symbol_table.SymbolScopes) != 1 {
		t.Errorf("Incorrect table length")
	}
}

func TestAnalyse_EmptySyntaxTree(t *testing.T) {

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{}
	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "",
		TypeRule:      "",
		ParameterRule: "",
		FunctionRule:  "",
	}

	_, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("syntax tree is empty").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}

}

func TestAnalyse_Error(t *testing.T) {

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
							Value:  "",
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

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	_, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("declaration has no name defined").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}

}

func TestAnalyse_Valid(t *testing.T) {

	expected_res := []services.Symbol{
		{Type: "int", Name: "blue", Scope: 0},
	}

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

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		if len(symbol_table_artefact.SymbolScopes) != len(expected_res) {
			t.Errorf("not enough symbols identified")
		} else {
			for i, symbol := range symbol_table_artefact.SymbolScopes {
				if symbol.Name != expected_res[i].Name || symbol.Type != expected_res[i].Type || symbol.Scope != expected_res[i].Scope {
					t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
				}
			}
		}

	}

}

func TestAnalyse_SameSymbolNames_DifferentScope(t *testing.T) {

	expected_res := []services.Symbol{
		{Type: "int", Name: "blue", Scope: 0},
		{Type: "int", Name: "blue", Scope: 1},
	}

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
					Symbol: "START_SCOPE",
					Value:  "{",
				},
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
					Symbol: "END_SCOPE",
					Value:  "}",
				},
			},
		},
	}
	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		if len(symbol_table_artefact.SymbolScopes) != len(expected_res) {
			t.Errorf("not enough symbols identified")
		} else {
			for i, symbol := range symbol_table_artefact.SymbolScopes {
				if symbol.Name != expected_res[i].Name || symbol.Type != expected_res[i].Type || symbol.Scope != expected_res[i].Scope {
					t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
				}
			}
		}

	}

}

func TestAnalyse_UndeclaredSymbol(t *testing.T) {

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

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	_, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected for undeclared variable")
	} else {
		if err.Error() != fmt.Errorf("variable not declared within it's scope: blue").Error() {
			t.Errorf("Incorrect error: %v", err)
		}
	}

}

func TestAnalyse_NoEndScope(t *testing.T) {
	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
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
					Value:  ";",
				},
			},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	_, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("end scope symbol not found for start scope, please recheck source code").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}
}

func TestAnalyse_Valid_FunctionType(t *testing.T) {

	expected_res := []services.Symbol{
		{Type: "int", Name: "function_test", Scope: 0},
		{Type: "int", Name: "blue", Scope: 1},
		{Type: "int", Name: "red", Scope: 0},
	}

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "MAIN",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "FUNCTION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "FUNCTION_DECLARATION",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "KEYWORD",
									Value:  "func",
								},
								{
									Symbol: "IDENTIFIER",
									Value:  "function_test",
								},
								{
									Symbol: "PARAMETER",
									Value:  "string",
								},
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
							},
						},
						{
							Symbol: "START_SCOPE",
							Value:  "{",
						},
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
							Symbol: "END_SCOPE",
							Value:  "}",
						},
					},
				},
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
									Value:  "red",
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
			},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION_DECLARATION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		if len(symbol_table_artefact.SymbolScopes) != len(expected_res) {
			t.Errorf("not enough symbols identified")
		} else {
			for i, symbol := range symbol_table_artefact.SymbolScopes {
				if symbol.Name != expected_res[i].Name || symbol.Type != expected_res[i].Type || symbol.Scope != expected_res[i].Scope {
					t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
				}
			}
		}

	}

}

func TestAnalyse_Valid_NoFunctionName(t *testing.T) {

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "MAIN",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "FUNCTION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "FUNCTION_DECLARATION",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "KEYWORD",
									Value:  "func",
								},
								{
									Symbol: "IDENTIFIER",
									Value:  "",
								},
								{
									Symbol: "PARAMETER",
									Value:  "string",
								},
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
							},
						},
						{
							Symbol: "START_SCOPE",
							Value:  "{",
						},
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
							Symbol: "END_SCOPE",
							Value:  "}",
						},
					},
				},
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
									Value:  "red",
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
			},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION_DECLARATION",
	}

	_, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("function has no name defined").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}

}

func TestAnalyse_Valid_FunctionParameters(t *testing.T) {

	expected_res := []services.Symbol{
		{Type: "string", Name: "purple", Scope: 1},
		{Type: "int", Name: "func_name", Scope: 0, Parameters: []services.Symbol{{Name: "purple", Type: "string", Scope: 1}}},
		{Type: "int", Name: "blue", Scope: 1},
		{Type: "int", Name: "red", Scope: 0},
	}

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "MAIN",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "FUNCTION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "FUNCTION_DECLARATION",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "KEYWORD",
									Value:  "func",
								},
								{
									Symbol: "IDENTIFIER",
									Value:  "func_name",
								},
								{
									Symbol: "PARAMETER",
									Value:  "",
									Children: []*services.TreeNode{
										{
											Symbol: "TYPE",
											Value:  "",
											Children: []*services.TreeNode{
												{
													Symbol: "KEYWORD",
													Value:  "string",
												},
											},
										},
										{
											Symbol: "IDENTIFIER",
											Value:  "purple",
										},
									},
								},
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
							},
						},
						{
							Symbol: "START_SCOPE",
							Value:  "{",
						},
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
							Symbol: "END_SCOPE",
							Value:  "}",
						},
					},
				},
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
									Value:  "red",
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
			},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION_DECLARATION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		if len(symbol_table_artefact.SymbolScopes) != len(expected_res) {
			t.Errorf("not enough symbols identified")
		} else {

			for i, symbol := range symbol_table_artefact.SymbolScopes {
				if symbol.Name != expected_res[i].Name || symbol.Type != expected_res[i].Type || symbol.Scope != expected_res[i].Scope {
					t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
				} else {
					if len(symbol.Parameters) != len(expected_res[i].Parameters) {
						t.Errorf("not enough function parameters")
					}
					for p, param := range symbol.Parameters {
						if param.Name != expected_res[i].Parameters[p].Name || param.Type != expected_res[i].Parameters[p].Type {
							t.Errorf("Function parameter is incorrect: %v", param)
						}
					}
				}
			}

		}

	}

}

func TestAnalyse_Valid_TypeCheckAssignment(t *testing.T) {

	expected_res := []services.Symbol{
		{Type: "string", Name: "purple", Scope: 1},
		{Type: "int", Name: "func_name", Scope: 0, Parameters: []services.Symbol{{Name: "purple", Type: "string", Scope: 1}}},
		{Type: "int", Name: "blue", Scope: 1},
		{Type: "int", Name: "red", Scope: 0},
	}

	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "MAIN",
			Value:  "",
			Children: []*services.TreeNode{
				{
					Symbol: "FUNCTION",
					Value:  "",
					Children: []*services.TreeNode{
						{
							Symbol: "FUNCTION_DECLARATION",
							Value:  "",
							Children: []*services.TreeNode{
								{
									Symbol: "KEYWORD",
									Value:  "func",
								},
								{
									Symbol: "IDENTIFIER",
									Value:  "func_name",
								},
								{
									Symbol: "PARAMETER",
									Value:  "",
									Children: []*services.TreeNode{
										{
											Symbol: "TYPE",
											Value:  "",
											Children: []*services.TreeNode{
												{
													Symbol: "KEYWORD",
													Value:  "string",
												},
											},
										},
										{
											Symbol: "IDENTIFIER",
											Value:  "purple",
										},
									},
								},
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
							},
						},
						{
							Symbol: "START_SCOPE",
							Value:  "{",
						},
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
															Symbol: "int",
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
							Symbol: "END_SCOPE",
							Value:  "}",
						},
					},
				},
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
									Value:  "red",
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
			},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:   "IDENTIFIER",
		TypeRule:       "TYPE",
		ParameterRule:  "PARAMETER",
		FunctionRule:   "FUNCTION_DECLARATION",
		AssignmentRule: "ASSIGNMENT",
		OperatorRule:   "OPERATOR",
		TermRule:       "EXPRESSION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("Error: %v", err)
	} else {

		if len(symbol_table_artefact.SymbolScopes) != len(expected_res) {
			t.Errorf("not enough symbols identified")
		} else {

			for i, symbol := range symbol_table_artefact.SymbolScopes {
				if symbol.Name != expected_res[i].Name || symbol.Type != expected_res[i].Type || symbol.Scope != expected_res[i].Scope {
					t.Errorf("Symbol is incorrect: %v %v %v", symbol.Name, symbol.Scope, symbol.Type)
				} else {
					if len(symbol.Parameters) != len(expected_res[i].Parameters) {
						t.Errorf("not enough function parameters")
					}
					for p, param := range symbol.Parameters {
						if param.Name != expected_res[i].Parameters[p].Name || param.Type != expected_res[i].Parameters[p].Type {
							t.Errorf("Function parameter is incorrect: %v", param)
						}
					}
				}
			}

		}

	}

}

func TestTraverseSyntaxTree_NilNode(t *testing.T) {
	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{},
	}

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	symbol_table_artefact := &services.SymbolTableArtefact{
		SymbolScopes: []services.Symbol{
			{Name: "blue", Type: "int", Scope: 0},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	err := services.TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestTraverseSyntaxTree_NoName(t *testing.T) {
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
							Value:  "",
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

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{
				"blue": {
					Name:  "blue",
					Type:  "int",
					Scope: 0,
				},
			},
		},
	}

	symbol_table_artefact := &services.SymbolTableArtefact{
		SymbolScopes: []services.Symbol{
			{Name: "blue", Type: "int", Scope: 0},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	err := services.TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("declaration has no name defined").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}
}

func TestTraverseSyntaxTree_VariableNotDeclared(t *testing.T) {
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

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{},
		},
	}

	symbol_table_artefact := &services.SymbolTableArtefact{
		SymbolScopes: []services.Symbol{
			{Name: "blue", Type: "int", Scope: 0},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	err := services.TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("variable not declared within it's scope: blue").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}
}

func TestTraverseSyntaxTree_NoStartScope(t *testing.T) {
	scope_rules := []*services.ScopeRule{
		{Start: "{", End: "}"},
	}

	syntax_tree := services.SyntaxTree{
		Root: &services.TreeNode{
			Symbol: "STATEMENT",
			Value:  "}",
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

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{},
		},
	}

	symbol_table_artefact := &services.SymbolTableArtefact{
		SymbolScopes: []services.Symbol{
			{Name: "blue", Type: "int", Scope: 0},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	err := services.TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)

	if err == nil {
		t.Errorf("Error expected")
	} else {

		if err.Error() != fmt.Errorf("end scope symbol found without starting scope, please recheck source code").Error() {
			t.Errorf("incorrect error: %v", err)
		}

	}
}

func TestTraverseSyntaxTree_Valid(t *testing.T) {
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

	symbol_table := &services.SymbolTable{
		SymbolScopes: []map[string]services.Symbol{
			{},
		},
	}

	symbol_table_artefact := &services.SymbolTableArtefact{
		SymbolScopes: []services.Symbol{
			{Name: "blue", Type: "int", Scope: 0},
		},
	}

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	err := services.TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestStringifySymbolTable(t *testing.T) {
	expected_res := "-------------------------------------------------------------------------- \nSYMBOL TABLE\n"
	expected_res += "  Name: blue  Type: int  Scope: 0\n"
	expected_res += "--------------------------------------------------------------------------\n"

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

	type_rules := []services.TypeRule{}
	rules := services.GrammarRules{
		VariableRule:  "IDENTIFIER",
		TypeRule:      "TYPE",
		ParameterRule: "PARAMETER",
		FunctionRule:  "FUNCTION",
	}

	symbol_table_artefact, err := services.Analyse(scope_rules, syntax_tree, rules, type_rules)

	if err != nil {
		t.Errorf("%v", err)
	}

	string_artefact := services.StringifySymbolTable(symbol_table_artefact)
	if string_artefact != expected_res {
		t.Errorf("Incorrect string generated")
	}
}

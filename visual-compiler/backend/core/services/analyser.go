package services

import (
	"fmt"
)

type Symbol struct {
	Name       string
	Type       string
	Scope      int
	Parameters []Symbol
}

type SymbolTable struct {
	SymbolScopes []map[string]Symbol
}

type SymbolTableArtefact struct {
	SymbolScopes []Symbol
}

type ScopeRule struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Entered bool
}

// Name: CreateEmptySymbolTable
//
// Parameters: none
//
// Return: *SymbolTable
//
// Create an empty symbol table with an empty scope data and return the memory address to it
func CreateEmptySymbolTable() *SymbolTable {

	return &SymbolTable{
		SymbolScopes: []map[string]Symbol{
			make(map[string]Symbol),
		},
	}
}

// Name: CreateEmptySymbolTableArtefact
//
// Parameters: none
//
// Return: *SymbolTableArtefact
//
// Create an empty symbol table artefact with an empty scope data and return the memory address to it
func CreateEmptySymbolTableArtefact() *SymbolTableArtefact {

	return &SymbolTableArtefact{
		SymbolScopes: []Symbol{},
	}

}

// Name: BindSymbol
//
// Parameters: *SymbolTable
//
// Return: error
//
// Bind a name to the properties of an item in the symbol table if it does not declared in the current scope already
func BindSymbol(symbol_table *SymbolTable, symbol Symbol) error {

	current_scope := symbol_table.SymbolScopes[len(symbol_table.SymbolScopes)-1]
	_, symbol_found := current_scope[symbol.Name]

	if symbol_found {
		return fmt.Errorf("symbol already declared in scope: %v", symbol.Name)
	} else {
		current_scope[symbol.Name] = symbol
		return nil
	}

}

// Name: LookupName
//
// Parameters: *SymbolTable, string
//
// Return: Symbol,error
//
// Determines if the name of the symbol exists in the symbol table
func LookupName(symbol_table *SymbolTable, symbol_name string) (Symbol, error) {

	for index := len(symbol_table.SymbolScopes) - 1; index >= 0; index-- {
		symbol_data, symbol_found := symbol_table.SymbolScopes[index][symbol_name]
		if symbol_found {
			return symbol_data, nil
		}
	}

	return Symbol{}, fmt.Errorf("symbol not found")
}

// Name: EnterNewScope
//
// Parameters: *SymbolTable
//
// Return: none
//
// Adds new empty scope to the symbol table
func EnterNewScope(symbol_table *SymbolTable) {
	symbol_table.SymbolScopes = append(symbol_table.SymbolScopes, make(map[string]Symbol))
}

// Name: ExitScope
//
// Parameters: *SymbolTable
//
// Return: error
//
// Exit a scope and re-establish symbol table to previous state, by removing scope at the last index
func ExitScope(symbol_table *SymbolTable) error {

	if len(symbol_table.SymbolScopes) == 1 {
		return fmt.Errorf("error: could not exit global scope")
	}

	symbol_table.SymbolScopes = symbol_table.SymbolScopes[:len(symbol_table.SymbolScopes)-1]
	return nil
}

// Name: Analyse
//
// Parameters: []ScopeRule,SyntaxTree,string,string,string,string
//
// Return: SymbolTable, error
//
// Receive a syntax tree and scope rules to scope check the parse tree, and create and return a symbol table or error
func Analyse(scope_rules []*ScopeRule, syntax_tree SyntaxTree, type_rule string, variable_rule string, function_rule string, parameter_rule string) (SymbolTableArtefact, error) {

	if syntax_tree.Root == nil {
		return SymbolTableArtefact{}, fmt.Errorf("syntax tree is empty")
	}

	symbol_table := CreateEmptySymbolTable()

	symbol_table_artefact := CreateEmptySymbolTableArtefact()

	err := TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, type_rule, variable_rule, function_rule, parameter_rule)
	if err != nil {
		return *symbol_table_artefact, err
	}

	for _, rule := range scope_rules {
		if rule.Entered {
			return *symbol_table_artefact, fmt.Errorf("end scope symbol not found for start scope, please recheck source code")
		}

	}

	return *symbol_table_artefact, nil
}

// Name: HandleFunctionScope
//
// Parameters: *Symbol,*TreeNode,*SymbolTable,*SymbolTableArtefact,string,string,string,string
//
// Return: error
//
// Determine name,type and scope of variables in a function declaration
func HandleFunctionScope(new_symbol *Symbol, child *TreeNode, symbol_table *SymbolTable, symbol_table_artefact *SymbolTableArtefact, type_rule string, variable_rule string, function_rule string, parameter_rule string) error {
	for _, function_child := range child.Children {
		if function_child.Symbol == variable_rule {
			if function_child.Value == "" {
				if len(function_child.Children) > 0 {
					current_child := function_child.Children[0]
					for current_child.Value == "" {
						if len(current_child.Children) > 0 {
							current_child = current_child.Children[0]
						}
					}
					new_symbol.Name = current_child.Value
				} else {
					return fmt.Errorf("function has no name defined")
				}
			} else {
				new_symbol.Name = function_child.Value
			}
		}

		if function_child.Symbol == type_rule {
			if function_child.Value == "" {
				if len(function_child.Children) > 0 {
					current_child := function_child.Children[0]
					for current_child.Value == "" {
						if len(current_child.Children) > 0 {
							current_child = current_child.Children[0]
						}
					}
					new_symbol.Type = current_child.Value
				} else {
					new_symbol.Type = function_child.Value
				}
			} else {
				new_symbol.Type = function_child.Value
			}
		}

		if function_child.Symbol == parameter_rule {

			parameter_symbol := Symbol{}

			if len(function_child.Children) > 0 {

				for _, current_child := range function_child.Children {

					if current_child.Symbol == type_rule {
						if current_child.Value == "" && len(current_child.Children) > 0 {
							parameter_child := current_child.Children[0]
							for parameter_child.Value == "" {
								if len(parameter_child.Children) > 0 {
									parameter_child = parameter_child.Children[0]
								}
							}
							parameter_symbol.Type = parameter_child.Value
						} else {
							parameter_symbol.Type = current_child.Value
						}

					}

					if current_child.Symbol == variable_rule {

						if current_child.Value == "" {
							if len(current_child.Children) > 0 {
								parameter_child := current_child.Children[0]
								for parameter_child.Value == "" {
									if len(parameter_child.Children) > 0 {
										parameter_child = parameter_child.Children[0]
									}
								}
								parameter_symbol.Name = parameter_child.Value
							} else {
								return fmt.Errorf("declaration has no name defined")
							}
						} else {
							parameter_symbol.Name = current_child.Value
						}

					}
				}

			}

			if parameter_symbol.Name != "" && parameter_symbol.Type != "" {
				parameter_symbol.Scope = len(symbol_table.SymbolScopes)

				err := BindSymbol(symbol_table, parameter_symbol)
				if err != nil {
					return err
				}
				symbol_table_artefact.SymbolScopes = append(symbol_table_artefact.SymbolScopes, parameter_symbol)
				new_symbol.Parameters = append(new_symbol.Parameters, parameter_symbol)
			}
		}
	}

	return nil
}

// Name: HandleFunctionScope
//
// Parameters: *Symbol,*TreeNode,string,string
//
// Return: error
//
// Determines type,name and scope of variables
func HandleDeclarationScope(new_symbol *Symbol, child *TreeNode, type_rule string, variable_rule string) error {
	if child.Symbol == variable_rule {
		if child.Value == "" {
			if len(child.Children) > 0 {
				current_child := child.Children[0]
				for current_child.Value == "" {
					if len(current_child.Children) > 0 {
						current_child = current_child.Children[0]
					}
				}
				new_symbol.Name = current_child.Value
			} else {
				return fmt.Errorf("declaration has no name defined")
			}
		} else {
			new_symbol.Name = child.Value
		}
	}

	if child.Symbol == type_rule {
		if child.Value == "" {
			if len(child.Children) > 0 {
				current_child := child.Children[0]
				for current_child.Value == "" {
					if len(current_child.Children) > 0 {
						current_child = current_child.Children[0]
					}
				}
				new_symbol.Type = current_child.Value
			} else {
				new_symbol.Type = child.Value
			}
		} else {
			new_symbol.Type = child.Value
		}

		//account for system must determine type of IDENTIFIER on its own -> type inference
		/*if new_symbol.Type == "" {

		}*/
	}

	return nil
}

// Name: TraverseSyntaxTree
//
// Parameters: []ScopeRule,*TreeNode,*SymbolTable,string,string,string,string
//
// Return: error
//
// Function used to recursively traverse the syntax tree and build the symbol table
func TraverseSyntaxTree(scope_rules []*ScopeRule, current_tree_node *TreeNode, symbol_table *SymbolTable, symbol_table_artefact *SymbolTableArtefact, type_rule string, variable_rule string, function_rule string, parameter_rule string) error {

	if current_tree_node == nil {
		return nil
	}

	for _, rule := range scope_rules {
		if current_tree_node.Value == rule.Start {
			EnterNewScope(symbol_table)
			rule.Entered = true
		}
	}

	new_symbol := Symbol{}
	for _, child := range current_tree_node.Children {

		if child.Symbol == function_rule {

			err := HandleFunctionScope(&new_symbol, child, symbol_table, symbol_table_artefact, type_rule, variable_rule, function_rule, parameter_rule)

			if err != nil {
				return err
			}

		} else {
			err := HandleDeclarationScope(&new_symbol, child, type_rule, variable_rule)

			if err != nil {
				return err
			}
		}

	}

	if new_symbol.Name != "" && new_symbol.Type != "" {
		new_symbol.Scope = len(symbol_table.SymbolScopes) - 1

		err := BindSymbol(symbol_table, new_symbol)
		if err != nil {
			return err
		}
		symbol_table_artefact.SymbolScopes = append(symbol_table_artefact.SymbolScopes, new_symbol)

	}

	if new_symbol.Type == "" && new_symbol.Name != "" {
		_, err := LookupName(symbol_table, new_symbol.Name)
		if err != nil {
			return fmt.Errorf("variable not declared within it's scope: %v", new_symbol.Name)
		}
	}

	for _, child := range current_tree_node.Children {
		if child.Symbol != function_rule {
			err := TraverseSyntaxTree(scope_rules, child, symbol_table, symbol_table_artefact, type_rule, variable_rule, function_rule, parameter_rule)
			if err != nil {
				return err
			}
		}

	}

	for _, rule := range scope_rules {
		if current_tree_node.Value == rule.End {
			if rule.Entered {

				ExitScope(symbol_table)
				rule.Entered = false

			} else {
				return fmt.Errorf("end scope symbol found without starting scope, please recheck source code")
			}
		}
	}

	return nil
}

// Name: StringifySymbolTable
//
// Parameters: SymbolTableArtefact
//
// Return: string
//
// Returns a string that converts the symbol table to a string
func StringifySymbolTable(symbol_table SymbolTableArtefact) string {

	output := "-------------------------------------------------------------------------- \nSYMBOL TABLE\n"

	for _, symbol := range symbol_table.SymbolScopes {
		output += fmt.Sprintf("  Name: %v  Type: %v  Scope: %v\n", symbol.Name, symbol.Type, symbol.Scope)
	}
	output += "--------------------------------------------------------------------------\n"

	return output

}

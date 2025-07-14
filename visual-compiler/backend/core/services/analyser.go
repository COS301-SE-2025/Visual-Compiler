package services

import (
	"fmt"
)

type Symbol struct {
	Name  string
	Type  string
	Scope int
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

// Name: PerformScopeCheck
//
// Parameters: []ScopeRule,SyntaxTree
//
// Return: SymbolTable, error
//
// Receive a syntax tree and scope rules to scope check the parse tree, and create and return a symbol table or error
func PerformScopeCheck(scope_rules []*ScopeRule, syntax_tree SyntaxTree) (SymbolTableArtefact, error) {

	if syntax_tree.Root == nil {
		return SymbolTableArtefact{}, fmt.Errorf("syntax tree is empty")
	}

	symbol_table := CreateEmptySymbolTable()

	symbol_table_artefact := CreateEmptySymbolTableArtefact()

	err := TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact)
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

// Name: TraverseSyntaxTree
//
// Parameters: []ScopeRule,*TreeNode,*SymbolTable
//
// Return: error
//
// Function used to recursively traverse the syntax tree and build the symbol table
func TraverseSyntaxTree(scope_rules []*ScopeRule, current_tree_node *TreeNode, symbol_table *SymbolTable, symbol_table_artefact *SymbolTableArtefact) error {

	if current_tree_node == nil {
		return nil
	}

	for _, rule := range scope_rules {
		if current_tree_node.Value == rule.Start {
			EnterNewScope(symbol_table)
			rule.Entered = true
		}
	}

	//if current_tree_node.Symbol == "DECLARATION" {

	new_symbol := Symbol{}
	for _, child := range current_tree_node.Children {
		if child.Symbol == "TYPE" {
			if child.Value == "" && len(child.Children) > 0 {
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

			//account for system must determine type of IDENTIFIER on its own
			/*if new_symbol.Type == "" {

			}*/
		}

		if child.Symbol == "IDENTIFIER" {
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
	//}

	for _, child := range current_tree_node.Children {
		err := TraverseSyntaxTree(scope_rules, child, symbol_table, symbol_table_artefact)
		if err != nil {
			return err
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

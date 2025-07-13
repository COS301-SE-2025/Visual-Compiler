package services

import (
	"fmt"
)

type Symbol struct {
	Name  string
	Type  string
	Scope string
}

type SymbolTable struct {
	SymbolScopes []map[string]Symbol
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

// Name: CreateSymbolTable
//
// Parameters: SyntaxTree
//
// Return: SyntaxTree, error
//
// Receive a syntax tree and perform type checking and scope checking to return a symbol table
func CreateSymbolTable(syntax_tree SyntaxTree) (SyntaxTree, error) {

	if syntax_tree.Root == nil {
		return syntax_tree, fmt.Errorf("syntax tree is empty")
	}

	return syntax_tree, fmt.Errorf("error")
}

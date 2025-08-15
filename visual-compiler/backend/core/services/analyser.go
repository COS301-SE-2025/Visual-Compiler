package services

import (
	"fmt"
)

// struct to store data on a symbol for the symbol table
type Symbol struct {
	Name       string
	Type       string
	Scope      int
	Parameters []Symbol
	Assign     bool
}

// struct to store data on a symbol table
type SymbolTable struct {
	SymbolScopes []map[string]Symbol
}

// struct to store data on a symbol table artefact
type SymbolTableArtefact struct {
	SymbolScopes []Symbol
}

// struct to store data for a scope rule
type ScopeRule struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Entered bool
}

// struct to store data for valid type rules.
//
// type rules define the valid result types when using an expression
type TypeRule struct {
	ResultData string
	Assignment string
	LHSData    string
	Operator   []string
	RHSData    string
}

type GrammarRules struct {
	TypeRule       string
	VariableRule   string
	FunctionRule   string
	ParameterRule  string
	AssignmentRule string
	OperatorRule   string
	TermRule       string
}

type AssignmentData struct {
	ResultData Symbol
	Terms      []Symbol
	Operator   Symbol
	Assignment Symbol
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

// Name: HandleFunctionScope
//
// Parameters: *Symbol,*TreeNode,*SymbolTable,*SymbolTableArtefact,string,string,string,string
//
// Return: error
//
// Determine name,type and scope of variables in a function declaration
func HandleFunctionScope(new_symbol *Symbol, child *TreeNode, symbol_table *SymbolTable, symbol_table_artefact *SymbolTableArtefact, rules GrammarRules) error {
	for _, function_child := range child.Children {
		if function_child.Symbol == rules.VariableRule {
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

		if function_child.Symbol == rules.TypeRule {
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

		if function_child.Symbol == rules.ParameterRule {

			parameter_symbol := Symbol{}

			if len(function_child.Children) > 0 {

				for _, current_child := range function_child.Children {

					if current_child.Symbol == rules.TypeRule {
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

					if current_child.Symbol == rules.VariableRule {

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

// Name: HandleVariableScope
//
// Parameters: *Symbol,*TreeNode,string,string
//
// Return: error
//
// Determines type,name and scope of variables
func HandleVariableScope(new_symbol *Symbol, child *TreeNode, type_rule string, variable_rule string) error {
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

// Name: CreateAssignmentSymbol
//
// Parameters: *Symbol
//
// Return: Symbol
//
// Determines if the assignment statement is valid according to type rules
func CreateAssignmentSymbol(last_symbol *Symbol, assignment_value string) Symbol {

	last_symbol.Assign = true
	assignment_symbol := Symbol{Name: "ASSIGNMENT", Type: assignment_value}

	return assignment_symbol
}

// Name: CreateOperatorSymbol
//
// Parameters: *Symbol
//
// Return: Symbol
//
// Determines if the assignment statement is valid according to type rules
func CreateOperatorSymbol(last_symbol *Symbol, operator_value string) Symbol {

	operator_symbol := Symbol{}
	operator_symbol.Name = "OPERATOR"
	operator_symbol.Type = operator_value

	return operator_symbol
}

// Name: CreateTermSymbol
//
// Parameters: *Symbol,TreeNode,string
//
// Return: Symbol
//
// Creates a symbol for terms to be used during assignment
func CreateTermSymbol(last_symbol *Symbol, term_node TreeNode, variable_rule string) Symbol {

	term_symbol := Symbol{}
	term_symbol.Name = "TERM"

	for term_node.Value == "" && len(term_node.Children) > 0 {
		term_node = *term_node.Children[0]
	}

	if term_node.Symbol == variable_rule {
		term_symbol.Type = term_node.Value
	} else {
		term_symbol.Type = term_node.Symbol
	}

	return term_symbol
}

// Name: HandleTerms
//
// Parameters: *Symbol,TreeNode,string
//
// Return: Symbol
//
// Creates a symbol for terms to be used during assignment
func FindTerms(last_symbol *Symbol, term_node TreeNode, operator_symbol string, term_symbol string, variable_rule string) ([]Symbol, Symbol) {

	term_symbols := []Symbol{}
	operator := Symbol{}

	if term_node.Value == "" && len(term_node.Children) > 0 {
		for _, child := range term_node.Children {
			if child.Symbol == term_symbol {
				new_symbol := CreateTermSymbol(last_symbol, *child, variable_rule)
				term_symbols = append(term_symbols, new_symbol)
			} else if child.Symbol == operator_symbol {
				operator = CreateOperatorSymbol(last_symbol, child.Value)
			}
		}
	}

	return term_symbols, operator
}

// Name: HandleAssignment
//
// Parameters: AssignmentData,SymbolTable
//
// Return: error
//
// Perform scope and type check for assignment statements
func HandleAssignment(assignment_data AssignmentData, symbol_table SymbolTable, type_rules []TypeRule) error {

	if assignment_data.ResultData.Type == "" {
		return fmt.Errorf("error: no result data specified")
	}
	if assignment_data.Assignment.Type == "" {
		return fmt.Errorf("error: no assignment symbol used: %v", assignment_data.ResultData.Name)
	}
	if assignment_data.Operator.Type == "" && len(assignment_data.Terms) > 1 {
		return fmt.Errorf("error: No operator indicated for multiple terms in assignment: %v", assignment_data.ResultData.Name)
	}
	if len(assignment_data.Terms) == 0 {
		return fmt.Errorf("error: no terms identified for assignment: %v", assignment_data.ResultData.Name)
	}
	if assignment_data.Operator.Type != "" && len(assignment_data.Terms) < 2 {
		return fmt.Errorf("error: not enough terms identified for operator in assignment: %v", assignment_data.ResultData.Name)
	}

	valid_assignment := false

	for _, rule := range type_rules {
		if assignment_data.ResultData.Type == rule.ResultData {

			if assignment_data.Operator.Type != "" {

				for _, operator := range rule.Operator {
					if assignment_data.Operator.Type == operator {

						rhs_term_found := false
						lhs_term_found := false
						for _, term := range assignment_data.Terms {

							if term.Type == rule.LHSData {
								lhs_term_found = true
							}
							if term.Type == rule.RHSData {
								rhs_term_found = true
							}

							symbol, err := LookupName(&symbol_table, term.Type)
							if err == nil {
								if symbol.Type == rule.LHSData {
									lhs_term_found = true
								}
								if symbol.Type == rule.RHSData {
									rhs_term_found = true
								}
							}

						}

						if rhs_term_found && lhs_term_found {
							valid_assignment = true
						}

					}
				}
			} else {
				lhs_term_found := false

				for _, term := range assignment_data.Terms {

					if term.Type == rule.LHSData {
						lhs_term_found = true
					}

					/*current_scope := symbol_table.SymbolScopes[len(symbol_table.SymbolScopes)-1]
					symbol, term_found := current_scope[term.Type]
					fmt.Printf("%v", current_scope)*/
					symbol, err := LookupName(&symbol_table, term.Type)
					if err == nil {
						if symbol.Type == rule.LHSData {
							lhs_term_found = true
						}
					}
				}

				if lhs_term_found {
					valid_assignment = true
				}
			}
		}
	}

	if !valid_assignment {

		return fmt.Errorf("error: invalid types assigned to: %v %v", assignment_data.ResultData.Type, assignment_data.ResultData.Name)
	}

	return nil
}

// Name: TraverseSyntaxTree
//
// Parameters: []ScopeRule,*TreeNode,*SymbolTable,GrammarRules,[]TypeRule
//
// Return: error
//
// Function used to recursively traverse the syntax tree and build the symbol table.
// Performs the scope check and type check.
func TraverseSyntaxTree(scope_rules []*ScopeRule, current_tree_node *TreeNode, symbol_table *SymbolTable, symbol_table_artefact *SymbolTableArtefact, rules GrammarRules, type_rules []TypeRule) error {

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

	assignment_data := AssignmentData{}

	for _, child := range current_tree_node.Children {
		switch child.Symbol {

		case rules.FunctionRule:
			err := HandleFunctionScope(&new_symbol, child, symbol_table, symbol_table_artefact, rules)

			if err != nil {
				return err
			}

		case rules.TypeRule, rules.VariableRule:
			err := HandleVariableScope(&new_symbol, child, rules.TypeRule, rules.VariableRule)

			if err != nil {
				return err
			}

		case rules.AssignmentRule:
			assignment_symbol := CreateAssignmentSymbol(&new_symbol, child.Value)
			assignment_data.Assignment = assignment_symbol

		case rules.OperatorRule:
			operator_symbol := CreateOperatorSymbol(&new_symbol, child.Value)
			assignment_data.Operator = operator_symbol

		case rules.TermRule:
			term_symbol := CreateTermSymbol(&new_symbol, *child, rules.VariableRule)
			assignment_data.Terms = append(assignment_data.Terms, term_symbol)

		default:

			term_symbols, operator_symbol := FindTerms(&new_symbol, *child, rules.OperatorRule, rules.TermRule, rules.VariableRule)
			if len(term_symbols) > 0 {
				assignment_data.Terms = term_symbols
			}
			assignment_data.Operator = operator_symbol
		}

	}

	if new_symbol.Name != "" && new_symbol.Type != "" {
		new_symbol.Scope = len(symbol_table.SymbolScopes) - 1

		err := BindSymbol(symbol_table, new_symbol)
		if err != nil {
			return err
		}
		symbol_table_artefact.SymbolScopes = append(symbol_table_artefact.SymbolScopes, new_symbol)

		assignment_data.ResultData = new_symbol

	}

	if new_symbol.Type == "" && new_symbol.Name != "" {
		_, err := LookupName(symbol_table, new_symbol.Name)
		if err != nil {
			return fmt.Errorf("variable not declared within it's scope: %v", new_symbol.Name)
		}
	}

	if new_symbol.Assign {
		err := HandleAssignment(assignment_data, *symbol_table, type_rules)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	for _, child := range current_tree_node.Children {
		if child.Symbol != rules.FunctionRule {
			err := TraverseSyntaxTree(scope_rules, child, symbol_table, symbol_table_artefact, rules, type_rules)
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

// Name: Analyse
//
// Parameters: []ScopeRule,SyntaxTree,GrammarRules,[]TypeRule
//
// Return: SymbolTable, error
//
// Receive a syntax tree and scope rules to scope check the parse tree, and create and return a symbol table and semantically verified syntax tree or error
func Analyse(scope_rules []*ScopeRule, syntax_tree SyntaxTree, rules GrammarRules, type_rules []TypeRule) (SymbolTableArtefact, SyntaxTree, error) {

	if syntax_tree.Root == nil {
		return SymbolTableArtefact{}, SyntaxTree{}, fmt.Errorf("syntax tree is empty")
	}

	symbol_table := CreateEmptySymbolTable()

	symbol_table_artefact := CreateEmptySymbolTableArtefact()

	err := TraverseSyntaxTree(scope_rules, syntax_tree.Root, symbol_table, symbol_table_artefact, rules, type_rules)
	if err != nil {
		return *symbol_table_artefact, SyntaxTree{}, err
	}

	for _, rule := range scope_rules {
		if rule.Entered {
			return *symbol_table_artefact, SyntaxTree{}, fmt.Errorf("end scope symbol not found for start scope, please recheck source code")
		}

	}

	return *symbol_table_artefact, syntax_tree, nil
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

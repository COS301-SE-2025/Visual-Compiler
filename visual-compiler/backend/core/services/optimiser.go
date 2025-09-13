package services

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

// Name: OptimiseGoCode
//
// Parameters: string, bool, bool, bool
//
// Return: string, error
//
// Receives Go code and performs optimisation based on the selected techniques
func OptimiseGoCode(code string, constant_folding bool, dead_code bool, loop_unroling bool) (string, error) {

	if code == "" {
		return "", fmt.Errorf("source code is empty")
	}

	ast_file, file_set, err := ParseGoCode(code)
	if err != nil {
		return "", err
	}

	if constant_folding {
		err = PerformConstantFolding(ast_file, file_set)
		if err != nil {
			return "", err
		}
	}

	if dead_code {
		err = PerformDeadCodeElimination(ast_file, file_set)
		if err != nil {
			return "", err
		}
	}

	if loop_unroling {
		err = PerformLoopUnrolling(ast_file, file_set)
		if err != nil {
			return "", err
		}
	}

	code, err = StringifyAST(ast_file, file_set)
	if err != nil {
		return "", err
	}

	return code, nil
}

// Name:ParseGoCode
//
// Parameters: string
//
// Return: *ast.File, error
//
// Receives Go code and creates an abstract syntax tree
func ParseGoCode(code string) (*ast.File, *token.FileSet, error) {

	file_set := token.NewFileSet()
	ast_file, err := parser.ParseFile(file_set, "", code, parser.AllErrors)

	if err != nil {
		return nil, nil, err
	}

	return ast_file, file_set, nil
}

// Name: StringifyAST
//
// Parameters: *ast.File, *token.FileSet
//
// Return: string
//
// Converts the optimised AST to a string
func StringifyAST(ast_file *ast.File, file_set *token.FileSet) (string, error) {

	var buffer bytes.Buffer
	err := format.Node(&buffer, file_set, ast_file)
	if err != nil {
		return "", err
	}

	optimised_code := buffer.String()
	structured_code := RemoveExtraLines(optimised_code)
	return structured_code, nil
}

// Name: RemoveExtraLines
//
// Parameters: string
//
// Return: string
//
// Removes extra unneccessary blank lines from the optimised code
func RemoveExtraLines(code string) string {

	var result_code strings.Builder
	is_first_empty_line := true
	line_count := 0

	for _, character := range code {
		if character == '\n' {
			line_count++
			if line_count < 2 {
				result_code.WriteRune(character)
			} else if is_first_empty_line {
				result_code.WriteRune(character)
				is_first_empty_line = false
			}
		} else {
			line_count = 0
			result_code.WriteRune(character)
		}
	}

	return result_code.String()
}

// Name: PerformConstantFolding
//
// Parameters: *ast.File
//
// Return: string
//
// Performs constant folding on the source code
func PerformConstantFolding(ast_file *ast.File, file_set *token.FileSet) error {
	return nil
}

// Name: PerformDeadCodeElimination
//
// Parameters: *ast.File
//
// Return: string
//
// Performs dead code elimination on the source code
func PerformDeadCodeElimination(ast_file *ast.File, file_set *token.FileSet) error {
	ast_info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	type_check_conf := types.Config{Importer: nil}
	_, _ = type_check_conf.Check("", file_set, []*ast.File{ast_file}, ast_info)

	err := RemoveUnusedFunctions(ast_file, ast_info)
	if err != nil {
		return err
	}

	return nil
}

// Name: PerformLoopUnrolling
//
// Parameters: *ast.File
//
// Return: string
//
// Performs loop unrolling on the source code
func PerformLoopUnrolling(ast_file *ast.File, file_set *token.FileSet) error {
	var failed error

	ast.Inspect(ast_file, func(n ast.Node) bool {
		if block, is_block := n.(*ast.BlockStmt); is_block {
			new_stmts := make([]ast.Stmt, 0, len(block.List))

			for _, statement := range block.List {
				if for_statement, is_for := statement.(*ast.ForStmt); is_for {

					if unrolled, err := UnrollForLoop(for_statement); err != nil {
						failed = err
						return false
					} else if unrolled != nil {
						new_stmts = append(new_stmts, unrolled...)
					} else {
						new_stmts = append(new_stmts, statement)
					}

				} else {
					new_stmts = append(new_stmts, statement)
				}
			}

			block.List = new_stmts
		}
		return true
	})

	return failed
}

/* PerformConstantFolding Helper Functions */

/* PerformDeadCodeElimination Helper Functions */

// Name:RemoveUnusedFunction
//
// Parameters:  *ast.File, *types.Info
//
// Return: error
//
// Determines if the statement contains any unused functions and removes them from the AST
func RemoveUnusedFunctions(ast_file *ast.File, ast_info *types.Info) error {
	optimised_functions := []ast.Decl{}
	used_functions := make(map[string]bool)
	used_functions["main"] = true
	ast.Inspect(ast_file, func(node ast.Node) bool {
		function_call, valid_call := node.(*ast.CallExpr)
		if !valid_call {
			return true
		}

		ident, valid_ident := function_call.Fun.(*ast.Ident)
		if valid_ident {
			used_function := ast_info.Uses[ident]
			if used_function != nil {
				function, valid_function := used_function.(*types.Func)
				if valid_function {
					used_functions[function.Name()] = true
				}
			}
		}

		return true
	})

	for _, declaration := range ast_file.Decls {
		function, valid_function := declaration.(*ast.FuncDecl)
		if valid_function {
			_, is_called := used_functions[function.Name.Name]
			if is_called {
				err := EliminateFunctionDeadCode(function, ast_info)
				optimised_functions = append(optimised_functions, function)
				if err != nil {
					return err
				}
			}
		} else {
			optimised_functions = append(optimised_functions, declaration)
		}
	}

	ast_file.Decls = optimised_functions

	return nil
}

// Name: EliminateFunctionDeadCode
//
// Parameters: *ast.FuncDecl, *types.Info
//
// Return: error
//
// Performs dead code elimination on an individual function
func EliminateFunctionDeadCode(function *ast.FuncDecl, ast_info *types.Info) error {
	if function.Body == nil {
		return fmt.Errorf("empty function")
	}

	optimised_statements := []ast.Stmt{}
	unreachable := false
	unused_variables := make(map[string]string)

	for _, function_statements := range function.Body.List {

		if unreachable {
			continue
		}

		switch statement := function_statements.(type) {
		case *ast.AssignStmt: // search for unused assigned variables
			RemoveUnusedAssignedVariables(statement, &unused_variables, ast_info, function, &optimised_statements)
		case *ast.DeclStmt: //search for unused declared variables
			RemoveUnusedDeclaredVariables(statement, unused_variables, ast_info, function_statements, &optimised_statements)
		case *ast.IfStmt:
			RemoveUnusedIfStatement(unused_variables, function_statements, &optimised_statements, ast_info, function)

		case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
			optimised_statements = append(optimised_statements, function_statements)
			unreachable = true

		default:
			optimised_statements = append(optimised_statements, function_statements)
		}

	}

	function.Body.List = optimised_statements

	return nil
}

// Name:RemoveUnusedVariables
//
// Parameters: *ast.AssignStmt, map[string]string, *types.Info, *ast.FuncDecl, *[]ast.Stmt
//
// Return:
//
// Determines if the statement contains any unused assigned variables and removes them from the AST
func RemoveUnusedAssignedVariables(statement *ast.AssignStmt, unused_variables *map[string]string, ast_info *types.Info, function *ast.FuncDecl, optimised_statements *[]ast.Stmt) {
	if statement.Tok == token.DEFINE {
		new_lhs := []ast.Expr{}
		new_rhs := []ast.Expr{}
		for num, lhs := range statement.Lhs {
			ident, valid_ident := lhs.(*ast.Ident)
			if valid_ident {
				variable := ast_info.Defs[ident]
				used := false
				for _, used_variable := range ast_info.Uses {
					if used_variable == variable {
						used = true
						delete(*unused_variables, ident.Name)
						break
					}
				}
				var value_str string
				switch v := statement.Rhs[num].(type) {
				case *ast.BasicLit:
					value_str = v.Value
				case *ast.Ident:
					value_str = v.Name
				}
				(*unused_variables)[ident.Name] = value_str

				//used = false
				for _, function_statement_2 := range function.Body.List {
					if_statement, valid_if := function_statement_2.(*ast.IfStmt)
					if valid_if {

						ast.Inspect(function_statement_2, func(node ast.Node) bool {
							switch expression := node.(type) {
							case *ast.Ident:
								cond_indent, valid_cond := if_statement.Cond.(*ast.Ident)
								if valid_cond {
									ident, valid_var := node.(*ast.Ident)

									if valid_var {
										unused_variable, is_unused := (*unused_variables)[cond_indent.Name]
										if is_unused {
											switch if_statement.Cond.(type) {
											case *ast.Ident:
												if unused_variable == "true" {
													used = true
													delete((*unused_variables), ident.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}

											}
										}
									}
								}

							case *ast.BinaryExpr:
								lhs, valid_cond := expression.X.(*ast.Ident)
								if valid_cond {
									rhs, valid_var := expression.Y.(*ast.BasicLit)
									if valid_var {
										unused_variable, is_unused := (*unused_variables)[lhs.Name]
										if is_unused {
											operator := expression.Op.String()
											switch operator {
											case "==":
												if rhs.Value == unused_variable {
													used = true
													delete((*unused_variables), lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}
											case "!=":
												if rhs.Value != unused_variable {
													used = true
													delete((*unused_variables), lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}

											case ">=":
												if rhs.Value >= unused_variable {
													used = true
													delete(*unused_variables, lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}
											case "<=":
												if rhs.Value <= unused_variable {
													used = true
													delete(*unused_variables, lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}
											case ">":
												if rhs.Value > unused_variable {
													used = true
													delete(*unused_variables, lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}
											case "<":
												if rhs.Value < unused_variable {
													used = true
													delete(*unused_variables, lhs.Name)
												} else {
													used = false
													(*unused_variables)[ident.Name] = value_str
												}
											}
										}
									}
								}

							case *ast.UnaryExpr:
								lhs, valid_cond := expression.X.(*ast.Ident)
								if valid_cond {
									unused_variable, is_unused := (*unused_variables)[lhs.Name]
									if is_unused {
										operator := expression.Op.String()
										switch operator {
										case "!":
											if unused_variable == "false" {
												used = true
												delete(*unused_variables, lhs.Name)
											} else {
												used = false
												(*unused_variables)[ident.Name] = value_str
											}
										}
									}
								}

							default:
							}
							return true
						})

					}
				}
				if used {
					new_lhs = append(new_lhs, lhs)
					if num < len(statement.Rhs) {
						new_rhs = append(new_rhs, statement.Rhs[num])
					}

				}
			} else {
				new_lhs = append(new_lhs, lhs)
				if num < len(statement.Rhs) {
					new_rhs = append(new_rhs, statement.Rhs[num])
				}

			}
		}
		if len(new_lhs) > 0 {
			statement.Lhs = new_lhs
			statement.Rhs = new_rhs
			*optimised_statements = append(*optimised_statements, statement)

		}
	} else {
		*optimised_statements = append(*optimised_statements, statement)
	}
}

// Name:RemoveDeclaredVariables
//
// Parameters: *ast.AssignStmt, map[string]string, *types.Info, *ast.FuncDecl, *[]ast.Stmt
//
// Return:
//
// Determines if the statement contains any unused assigned variables and removes them from the AST
func RemoveUnusedDeclaredVariables(statement *ast.DeclStmt, unused_variables map[string]string, ast_info *types.Info, function_statements ast.Stmt, optimised_statements *[]ast.Stmt) {
	declaration_info, valid_declaration := statement.Decl.(*ast.GenDecl)
	if valid_declaration && declaration_info.Tok == token.VAR {
		new_spec := []ast.Spec{}
		for _, spec := range declaration_info.Specs {
			value, valid_spec := spec.(*ast.ValueSpec)
			if valid_spec {
				new_names := []*ast.Ident{}
				new_values := []ast.Expr{}
				for num, name := range value.Names {
					variable := ast_info.Defs[name]
					unused_variables[name.Name] = "True"
					used := false
					for _, used_variable := range ast_info.Uses {
						if used_variable == variable {
							used = true
							delete(unused_variables, name.Name)
							break
						}
					}
					if used {
						new_names = append(new_names, name)
						if num < len(value.Values) {
							new_values = append(new_values, value.Values[num])
						}
					}
				}
				if len(new_names) > 0 {
					value.Names = new_names
					value.Values = new_values
					new_spec = append(new_spec, value)
				}
			} else {
				new_spec = append(new_spec, spec)
			}
		}
		if len(new_spec) > 0 {
			declaration_info.Specs = new_spec
			*optimised_statements = append(*optimised_statements, function_statements)
		}
	} else {
		*optimised_statements = append(*optimised_statements, function_statements)
	}
}

// Name:RemoveUnusedIfStatement
//
// Parameters:  map[string]string, *ast.FuncDecl, *[]ast.Stmt, *types.Info
//
// Return:
//
// Determines if the statement contains any unused if statements and removes them from the AST
func RemoveUnusedIfStatement(unused_variables map[string]string, function_statements ast.Stmt, optimised_statements *[]ast.Stmt, ast_info *types.Info, function *ast.FuncDecl) {
	if_statement, valid_ifstatement := function_statements.(*ast.IfStmt)
	unreachable := false
	var optimised_if_body []ast.Stmt
	if valid_ifstatement {
		switch condition_variable := if_statement.Cond.(type) {

		case *ast.Ident:
			_, is_unused := unused_variables[condition_variable.Name]
			if !is_unused {
				if condition_variable.Name != "false" {
					for _, body_statement := range if_statement.Body.List {
						if unreachable {
							continue
						}

						switch statement := body_statement.(type) {
						case *ast.AssignStmt: // search for unused assigned variables
							RemoveUnusedAssignedVariables(statement, &unused_variables, ast_info, function, &optimised_if_body)
						case *ast.DeclStmt: //search for unused declared variables
							RemoveUnusedDeclaredVariables(statement, unused_variables, ast_info, function_statements, &optimised_if_body)
						case *ast.IfStmt:
							RemoveUnusedIfStatement(unused_variables, statement, &optimised_if_body, ast_info, function)

						case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
							optimised_if_body = append(optimised_if_body, body_statement)
							unreachable = true

						default:
							optimised_if_body = append(optimised_if_body, body_statement)
						}
					}
					if_statement.Body.List = optimised_if_body
					*optimised_statements = append(*optimised_statements, if_statement)
				}
			}

		case *ast.BinaryExpr:
			lhs, valid_expr := condition_variable.X.(*ast.Ident)
			if valid_expr {
				_, is_unused := unused_variables[lhs.Name]
				if !is_unused {
					for _, body_statement := range if_statement.Body.List {
						if unreachable {
							continue
						}

						switch statement := body_statement.(type) {
						case *ast.AssignStmt: // search for unused assigned variables
							RemoveUnusedAssignedVariables(statement, &unused_variables, ast_info, function, &optimised_if_body)
						case *ast.DeclStmt: //search for unused declared variables
							RemoveUnusedDeclaredVariables(statement, unused_variables, ast_info, function_statements, &optimised_if_body)
						case *ast.IfStmt:
							RemoveUnusedIfStatement(unused_variables, statement, &optimised_if_body, ast_info, function)

						case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
							optimised_if_body = append(optimised_if_body, body_statement)
							unreachable = true

						default:
							optimised_if_body = append(optimised_if_body, body_statement)
						}
					}
					if_statement.Body.List = optimised_if_body
					*optimised_statements = append(*optimised_statements, if_statement)
				}
			}

		case *ast.UnaryExpr:
			lhs, valid_expr := condition_variable.X.(*ast.Ident)
			if valid_expr {
				_, is_unused := unused_variables[lhs.Name]
				if !is_unused {
					for _, body_statement := range if_statement.Body.List {
						if unreachable {
							continue
						}

						switch statement := body_statement.(type) {
						case *ast.AssignStmt: // search for unused assigned variables
							RemoveUnusedAssignedVariables(statement, &unused_variables, ast_info, function, &optimised_if_body)
						case *ast.DeclStmt: //search for unused declared variables
							RemoveUnusedDeclaredVariables(statement, unused_variables, ast_info, function_statements, &optimised_if_body)
						case *ast.IfStmt:
							RemoveUnusedIfStatement(unused_variables, statement, &optimised_if_body, ast_info, function)

						case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
							optimised_if_body = append(optimised_if_body, body_statement)
							unreachable = true

						default:
							optimised_if_body = append(optimised_if_body, body_statement)
						}
					}
					if_statement.Body.List = optimised_if_body
					*optimised_statements = append(*optimised_statements, if_statement)
				}
			}

		default:
			*optimised_statements = append(*optimised_statements, function_statements)
		}
	}
}

/* PerformLoopUnrolling Helper Functions */

// Struct for the parameters of a for loop
type LoopInfo struct {
	VarName   string
	Start     int
	Stop      int
	Step      int
	Direction bool
	Body      *ast.BlockStmt
}

// Name: UnrollForLoop
//
// Parameters: *ast.ForStmt
//
// Return: []ast.Stmt, error
//
// Attempts to unroll a standard for loop with a constant boundary
func UnrollForLoop(for_statement *ast.ForStmt) ([]ast.Stmt, error) {

	loop_info, err := AnalyseForLoop(for_statement)

	if err != nil || loop_info == nil {
		return nil, fmt.Errorf("loop is not a standard for loop with a constant boundary")
	}

	unrolled_statements := GenerateStatements(loop_info)
	return unrolled_statements, nil
}

// Name: AnalyseForLoop
//
// Parameters: *ast.ForStmt
//
// Return: *LoopInfo, error
//
// Analyses for loop to determine if it can be unrolled
func AnalyseForLoop(for_statement *ast.ForStmt) (*LoopInfo, error) {

	if for_statement.Init == nil {
		return nil, fmt.Errorf("")
	}

	assign_statement, is_valid := for_statement.Init.(*ast.AssignStmt)
	if !is_valid || len(assign_statement.Lhs) != 1 || len(assign_statement.Rhs) != 1 {
		return nil, fmt.Errorf("")
	}

	identifier, is_valid := assign_statement.Lhs[0].(*ast.Ident)
	if !is_valid {
		return nil, fmt.Errorf("")
	}
	var_name := identifier.Name

	var start_value int
	switch rhs := assign_statement.Rhs[0].(type) {

	case *ast.BasicLit:
		if rhs.Kind != token.INT {
			return nil, fmt.Errorf("")
		}
		val, err := strconv.Atoi(rhs.Value)
		if err != nil {
			return nil, fmt.Errorf("")
		}
		start_value = val

	case *ast.UnaryExpr:
		if basic, exists := rhs.X.(*ast.BasicLit); exists && rhs.Op == token.SUB && basic.Kind == token.INT {
			val, err := strconv.Atoi(basic.Value)
			if err != nil {
				return nil, fmt.Errorf("")
			}
			start_value = -val
		} else {
			return nil, fmt.Errorf("")
		}

	default:
		return nil, fmt.Errorf("")
	}

	if for_statement.Cond == nil {
		return nil, fmt.Errorf("")
	}

	binary_expression, is_valid := for_statement.Cond.(*ast.BinaryExpr)
	if !is_valid {
		return nil, fmt.Errorf("")
	}

	condition, is_valid := binary_expression.X.(*ast.Ident)
	if !is_valid || condition.Name != var_name {
		return nil, fmt.Errorf("")
	}

	var step_direction bool
	if binary_expression.Op == token.LSS {
		step_direction = true
	} else if binary_expression.Op == token.GTR {
		step_direction = false
	} else {
		return nil, fmt.Errorf("")
	}

	var stop_value int
	switch y := binary_expression.Y.(type) {

	case *ast.BasicLit:
		val, err := strconv.Atoi(y.Value)
		if err != nil {
			return nil, fmt.Errorf("")
		}
		stop_value = val

	case *ast.UnaryExpr:
		if basic, exists := y.X.(*ast.BasicLit); exists && y.Op == token.SUB && basic.Kind == token.INT {
			val, err := strconv.Atoi(basic.Value)
			if err != nil {
				return nil, fmt.Errorf("")
			}
			stop_value = -val
		} else {
			return nil, fmt.Errorf("")
		}

	default:
		return nil, fmt.Errorf("")
	}

	if for_statement.Post == nil {
		return nil, fmt.Errorf("")
	}

	inc_statement, is_valid := for_statement.Post.(*ast.IncDecStmt)
	if !is_valid {
		return nil, nil
	}

	var step int
	if inc_statement.Tok == token.INC && step_direction {
		step = 1
	} else if inc_statement.Tok == token.DEC && !step_direction {
		step = -1
	} else {
		return nil, fmt.Errorf("")
	}

	post_identifier, is_valid := inc_statement.X.(*ast.Ident)
	if !is_valid || post_identifier.Name != var_name {
		return nil, fmt.Errorf("")
	}

	return &LoopInfo{
		VarName:   var_name,
		Start:     start_value,
		Stop:      stop_value,
		Step:      step,
		Direction: step_direction,
		Body:      for_statement.Body,
	}, nil
}

// Name: GenerateStatements
//
// Parameters: *LoopInfo
//
// Return: []ast.Stmt
//
// Generates the body statements for the unrolled loop
func GenerateStatements(info *LoopInfo) []ast.Stmt {

	var unrolled_statements []ast.Stmt

	if info.Direction {
		for i := info.Start; i < info.Stop; i++ {
			for _, statement := range info.Body.List {
				new_stmt := UnrollStatements(statement, info.VarName, i)
				unrolled_statements = append(unrolled_statements, new_stmt)
			}
		}

	} else {
		for i := info.Start; i > info.Stop; i-- {
			for _, statement := range info.Body.List {
				new_stmt := UnrollStatements(statement, info.VarName, i)
				unrolled_statements = append(unrolled_statements, new_stmt)
			}
		}
	}

	return unrolled_statements
}

// Name: UnrollStatements
//
// Parameters: ast.Stmt, string, int
//
// Return: ast.Stmt
//
// Clones the statement and substitutes the loop variable
func UnrollStatements(statement ast.Stmt, var_name string, value int) ast.Stmt {

	switch s := statement.(type) {

	case *ast.ExprStmt:
		return &ast.ExprStmt{
			X: UnrollExpressions(s.X, var_name, value),
		}

	case *ast.AssignStmt:
		new_lhs := make([]ast.Expr, len(s.Lhs))
		for i, lhs := range s.Lhs {
			new_lhs[i] = UnrollExpressions(lhs, var_name, value)
		}
		new_rhs := make([]ast.Expr, len(s.Rhs))
		for i, rhs := range s.Rhs {
			new_rhs[i] = UnrollExpressions(rhs, var_name, value)
		}
		return &ast.AssignStmt{
			Lhs: new_lhs,
			Tok: s.Tok,
			Rhs: new_rhs,
		}

	case *ast.IfStmt:
		var new_init ast.Stmt
		if s.Init != nil {
			new_init = UnrollStatements(s.Init, var_name, value)
		}
		var new_else ast.Stmt
		if s.Else != nil {
			new_else = UnrollStatements(s.Else, var_name, value)
		}
		return &ast.IfStmt{
			If:   s.If,
			Init: new_init,
			Cond: UnrollExpressions(s.Cond, var_name, value),
			Body: UnrollBlocks(s.Body, var_name, value),
			Else: new_else,
		}

	case *ast.ForStmt:
		var new_init ast.Stmt
		if s.Init != nil {
			new_init = UnrollStatements(s.Init, var_name, value)
		}
		var new_cond ast.Expr
		if s.Cond != nil {
			new_cond = UnrollExpressions(s.Cond, var_name, value)
		}
		var new_post ast.Stmt
		if s.Post != nil {
			new_post = UnrollStatements(s.Post, var_name, value)
		}
		return &ast.ForStmt{
			For:  s.For,
			Init: new_init,
			Cond: new_cond,
			Post: new_post,
			Body: UnrollBlocks(s.Body, var_name, value),
		}

	case *ast.IncDecStmt:
		return &ast.IncDecStmt{
			X:   UnrollExpressions(s.X, var_name, value),
			Tok: s.Tok,
		}

	default:
		return statement
	}
}

// Name: UnrollBlocks
//
// Parameters: *ast.BlockStmt, string, int
//
// Return: *ast.BlockStmt
//
// Clones the statements in a block and substitutes the loop variable
func UnrollBlocks(block *ast.BlockStmt, var_name string, value int) *ast.BlockStmt {

	if block == nil {
		return nil
	}

	new_statements := make([]ast.Stmt, len(block.List))

	for i, statement := range block.List {
		new_statements[i] = UnrollStatements(statement, var_name, value)
	}

	return &ast.BlockStmt{
		List: new_statements,
	}
}

// Name: UnrollExpressions
//
// Parameters: ast.Expr, string, int
//
// Return: ast.Expr
//
// Clones the expression and substitutes the loop variable
func UnrollExpressions(expr ast.Expr, var_name string, value int) ast.Expr {

	switch e := expr.(type) {

	case *ast.Ident:
		if e.Name == var_name {
			return &ast.BasicLit{
				Kind:  token.INT,
				Value: strconv.Itoa(value),
			}
		}
		return &ast.Ident{Name: e.Name}

	case *ast.CallExpr:
		new_args := make([]ast.Expr, len(e.Args))
		for i, arg := range e.Args {
			new_args[i] = UnrollExpressions(arg, var_name, value)
		}
		return &ast.CallExpr{
			Fun:  UnrollExpressions(e.Fun, var_name, value),
			Args: new_args,
		}

	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X:   UnrollExpressions(e.X, var_name, value),
			Sel: &ast.Ident{Name: e.Sel.Name},
		}

	case *ast.BinaryExpr:
		return &ast.BinaryExpr{
			X:  UnrollExpressions(e.X, var_name, value),
			Op: e.Op,
			Y:  UnrollExpressions(e.Y, var_name, value),
		}

	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			Op: e.Op,
			X:  UnrollExpressions(e.X, var_name, value),
		}

	case *ast.BasicLit:
		return &ast.BasicLit{
			Kind:  e.Kind,
			Value: e.Value,
		}

	default:
		return expr
	}
}

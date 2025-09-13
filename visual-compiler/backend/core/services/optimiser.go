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
		return "", fmt.Errorf("Source code is empty.")
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

	for _, declaration := range ast_file.Decls {
		function, valid_function := declaration.(*ast.FuncDecl)
		if valid_function {
			err := EliminateFunctionDeadCode(function, ast_info)
			if err != nil {
				return err
			}
		}
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

}

/* PerformConstantFolding Helper Functions */



/* PerformDeadCodeElimination Helper Functions */

// Name: EliminateFunctionDeadCode
//
// Parameters: *ast.FuncDecl, *types.Info
//
// Return: error
//
// Perform dead code elimination on an individual function
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
			if statement.Tok == token.DEFINE {
				new_lhs := []ast.Expr{}
				new_rhs := []ast.Expr{}
				for num, lhs := range statement.Lhs {
					ident, valid_ident := lhs.(*ast.Ident)
					if valid_ident {
						variable := ast_info.Defs[ident]
						var value_str string
						switch v := statement.Rhs[num].(type) {
						case *ast.BasicLit:
							value_str = v.Value
						case *ast.Ident:
							value_str = v.Name
						default:
						}
						unused_variables[ident.Name] = value_str
						used := false
						for _, used_variable := range ast_info.Uses {
							if used_variable == variable {
								used = true
								delete(unused_variables, ident.Name)
								break
							}
						}
						for _, function_statement_2 := range function.Body.List {
							if_statement, valid_if := function_statement_2.(*ast.IfStmt)
							if valid_if {
								{
									ast.Inspect(function_statement_2, func(node ast.Node) bool {
										switch expression := node.(type) {
										case *ast.Ident:
											cond_indent, valid_cond := if_statement.Cond.(*ast.Ident)
											if valid_cond {
												ident, valid_var := node.(*ast.Ident)
												if valid_var {
													unused_variable, is_unused := unused_variables[cond_indent.Name]
													if is_unused {

														switch if_statement.Cond.(type) {
														case *ast.Ident:
															if unused_variable == "true" {
																used = true
																delete(unused_variables, ident.Name)
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
													unused_variable, is_unused := unused_variables[lhs.Name]
													if is_unused {
														operator := expression.Op.String()
														switch operator {
														case "==":
															if rhs.Value == unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}
														case "!=":
															if rhs.Value != unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}

														case ">=":
															if rhs.Value >= unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}
														case "<=":
															if rhs.Value <= unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}
														case ">":
															if rhs.Value > unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}
														case "<":
															if rhs.Value < unused_variable {
																used = true
																delete(unused_variables, lhs.Name)
															}
														}
													}
												}
											}

										case *ast.UnaryExpr:
											lhs, valid_cond := expression.X.(*ast.Ident)
											if valid_cond {
												unused_variable, is_unused := unused_variables[lhs.Name]
												if is_unused {
													operator := expression.Op.String()
													switch operator {
													case "!":
														if unused_variable == "false" {
															used = true
															delete(unused_variables, lhs.Name)
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
					optimised_statements = append(optimised_statements, statement)
				}
			} else {
				optimised_statements = append(optimised_statements, statement)
			}

		case *ast.DeclStmt: //search for unused declared variables
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
					optimised_statements = append(optimised_statements, function_statements)
				}
			} else {
				optimised_statements = append(optimised_statements, function_statements)
			}

		case *ast.IfStmt:
			if_statement, valid_ifstatement := function_statements.(*ast.IfStmt)
			if valid_ifstatement {
				switch condition_variable := if_statement.Cond.(type) {

				case *ast.Ident:
					_, is_unused := unused_variables[condition_variable.Name]
					if !is_unused {
						if condition_variable.Name != "false" {
							optimised_statements = append(optimised_statements, function_statements)
						}
					}

				case *ast.BinaryExpr:
					lhs, valid_expr := condition_variable.X.(*ast.Ident)
					if valid_expr {
						_, is_unused := unused_variables[lhs.Name]
						if !is_unused {
							optimised_statements = append(optimised_statements, function_statements)
						}
					}

				case *ast.UnaryExpr:
					lhs, valid_expr := condition_variable.X.(*ast.Ident)
					if valid_expr {
						_, is_unused := unused_variables[lhs.Name]
						if !is_unused {
							optimised_statements = append(optimised_statements, function_statements)
						}
					}

				default:
					optimised_statements = append(optimised_statements, function_statements)
				}
			}

		default:
			optimised_statements = append(optimised_statements, function_statements)
		}

		switch function_statements.(type) {
		case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
			unreachable = true
		}
	}

	function.Body.List = optimised_statements

	return nil
}

/* PerformLoopUnrolling Helper Functions */


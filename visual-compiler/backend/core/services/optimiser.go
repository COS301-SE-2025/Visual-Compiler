package services

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
)

// Name: OptimiseGoCode
//
// Parameters: string, bool, bool, bool
//
// Return: string,error
//
// Receive Go lang code and perform optimisation based on the techniques selected
func OptimiseGoCode(code string, constant_folding bool, dead_code bool, loop_unroling bool) (string, error) {

	if code == "" {
		return "", fmt.Errorf("code is empty")
	}

	ast_file, file_set, err := ParseGoCode(code)
	if err != nil {
		return "", err
	}

	if constant_folding {
		PerformConstantFolding(ast_file, file_set)
	}
	if dead_code {
		err = PerformDeadCodeElimination(ast_file, file_set)
		if err != nil {
			return "", err
		}
	}
	if loop_unroling {
		PerformLoopUnrolling(ast_file, file_set)
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
// Return: *ast.File,error
//
// Receive Go lang code and create an Abstract Syntax Tree
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
// Converts optimised AST to a string
func StringifyAST(ast_file *ast.File, file_set *token.FileSet) (string, error) {
	var buffer bytes.Buffer
	err := format.Node(&buffer, file_set, ast_file)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// Name: PerformConstantFolding
//
// Parameters: *ast.File
//
// Return: string
//
// Perform constant folding on the source code
func PerformConstantFolding(ast_file *ast.File, file_set *token.FileSet) {

}

// Name: PerformDeadCodeElimination
//
// Parameters: *ast.File
//
// Return: string
//
// Perform dead code elimination on the source code
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

	for _, function_statements := range function.Body.List {

		if unreachable {
			continue
		}

		switch statement := function_statements.(type) {
		case *ast.AssignStmt:
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
								break
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

		case *ast.DeclStmt:
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
							used := false
							for _, used_variable := range ast_info.Uses {
								if used_variable == variable {
									used = true
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

// Name: PerformLoopUnrolling
//
// Parameters: *ast.File
//
// Return: string
//
// Perform loop unrolling on the source code
func PerformLoopUnrolling(ast_file *ast.File, file_set *token.FileSet) {

}

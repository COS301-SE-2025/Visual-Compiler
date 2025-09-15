package services

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"math"
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
func OptimiseGoCode(code string, constant_folding bool, dead_code bool, loop_unrolling bool) (string, error) {

	if code == "" {
		return "", fmt.Errorf("source code is empty")
	}

	ast_file, file_set, err := ParseGoCode(code)
	if err != nil {
		return "", err
	}

	if loop_unrolling {
		err = PerformLoopUnrolling(ast_file, file_set)
		if err != nil {
			return "", err
		}
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

	constant_folder := &Folder{
		constants: make(map[string]float64),
	}

	err := constant_folder.FoldConstants(ast_file)
	return err
}

type AstData struct {
	ast_info        *types.Info
	ast_file        *ast.File
	file_set        *token.FileSet
	variable_values map[string]interface{}
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
	type_check_conf := types.Config{Importer: importer.Default()}
	_, _ = type_check_conf.Check("", file_set, []*ast.File{ast_file}, ast_info)

	variable_values := make(map[string]interface{})

	ast.Inspect(ast_file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, lhs := range assign.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				obj := ast_info.Defs[ident]
				if obj == nil {
					continue
				}
				if i < len(assign.Rhs) {
					rhs := assign.Rhs[i]
					var_value, valid_val := ast_info.Types[rhs]
					if valid_val && var_value.Value != nil {
						variable_values[ident.Name] = var_value.Value
					}
				}
			}
		}
		return true
	})

	ast_data := AstData{ast_info, ast_file, file_set, variable_values}

	err := RemoveUnusedFunctions(ast_data)
	if err != nil {
		return err
	}

	ast_info.Types = make(map[ast.Expr]types.TypeAndValue)
	ast_info.Defs = make(map[*ast.Ident]types.Object)
	ast_info.Uses = make(map[*ast.Ident]types.Object)

	type_check_conf = types.Config{Importer: importer.Default()}
	_, _ = type_check_conf.Check("", file_set, []*ast.File{ast_file}, ast_info)
	err = RemoveUnusedFunctions(ast_data)
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
	changed := true

	for changed {
		changed = false

		ast.Inspect(ast_file, func(n ast.Node) bool {
			block, is_block := n.(*ast.BlockStmt)
			if !is_block {
				return true
			}

			new_stmts := make([]ast.Stmt, 0, len(block.List))

			for _, statement := range block.List {
				if for_statement, is_for := statement.(*ast.ForStmt); is_for {
					unrolled, err := UnrollForLoop(for_statement)
					if err != nil {
						failed = err
						return false
					}

					if unrolled != nil {
						changed = true
						new_stmts = append(new_stmts, unrolled...)
					} else {
						new_stmts = append(new_stmts, statement)
					}

				} else {
					new_stmts = append(new_stmts, statement)
				}
			}

			block.List = new_stmts
			return true
		})
	}

	return failed
}

/* PerformConstantFolding Helper Functions */

// Struct for tracking constant values by variable name
type Folder struct {
	constants map[string]float64
}

// Name: FoldConstants (for Folder)
//
// Parameters: ast.Node
//
// Return: error
//
// Implements the ast visitor interface for the constant folder
func (constant_folder *Folder) FoldConstants(node ast.Node) error {

	var err error

	ast.Inspect(node, func(n ast.Node) bool {
		if err != nil {
			return false
		}

		switch node := n.(type) {

		case *ast.AssignStmt:
			err = constant_folder.HandleAssignment(node)

		case *ast.CallExpr:
			err = constant_folder.HandleFunctionCall(node)
		}

		return err == nil
	})

	return err
}

// Name: Visit (for Folder)
//
// Parameters: ast.Node
//
// Return: ast.Visitor
//
// Implements the ast visitor interface for the constant folder
func (constant_folder *Folder) Visit(node ast.Node) ast.Visitor {

	switch n := node.(type) {

	case *ast.AssignStmt:
		constant_folder.HandleAssignment(n)

	case *ast.CallExpr:
		constant_folder.HandleFunctionCall(n)
	}

	return constant_folder
}

// Name: HandleAssignment (for Folder)
//
// Parameters: *ast.AssignStmt
//
// Return: none
//
// Folds constants in assignment statements
func (constant_folder *Folder) HandleAssignment(assign *ast.AssignStmt) error {

	if assign.Tok != token.ASSIGN && assign.Tok != token.DEFINE {
		return nil
	}

	for i, lhs := range assign.Lhs {

		if i >= len(assign.Rhs) {
			continue
		}

		rhs := assign.Rhs[i]

		var var_name string
		if identifier, exists := lhs.(*ast.Ident); exists {
			var_name = identifier.Name
		} else {
			continue
		}

		value, valid, err := constant_folder.EvaluateExpression(rhs)
		if err != nil {
			return err
		}

		if valid {
			constant_folder.constants[var_name] = value
			assign.Rhs[i] = StructureConstant(value)
		}
	}

	return nil
}

// Name: HandleFunctionCall (for Folder)
//
// Parameters: *ast.CallExpr
//
// Return: none
//
// Folds constants in arguments of function calls
func (constant_folder *Folder) HandleFunctionCall(call *ast.CallExpr) error {

	for i, arg := range call.Args {

		value, valid, err := constant_folder.EvaluateExpression(arg)
		if err != nil {
			return err
		}

		if valid {
			call.Args[i] = StructureConstant(value)
		}
	}

	return nil
}

// Name: EvaluateExpression (for Folder)
//
// Parameters: ast.Expr
//
// Return: float64, bool
//
// Evaluates an expression to a constant value
func (constant_folder *Folder) EvaluateExpression(expr ast.Expr) (float64, bool, error) {

	switch e := expr.(type) {

	case *ast.BasicLit:

		switch e.Kind {

		case token.INT:
			if val, err := strconv.Atoi(e.Value); err == nil {
				return float64(val), true, nil
			}

		case token.FLOAT:
			if val, err := strconv.ParseFloat(e.Value, 64); err == nil {
				return val, true, nil
			}
		}

	case *ast.Ident:
		if val, yes := constant_folder.constants[e.Name]; yes {
			return val, true, nil
		}

	case *ast.BinaryExpr:

		lhs, lok, lerr := constant_folder.EvaluateExpression(e.X)
		if lerr != nil {
			return 0, false, lerr
		}

		rhs, rok, rerr := constant_folder.EvaluateExpression(e.Y)
		if rerr != nil {
			return 0, false, rerr
		}

		if lok && rok {

			switch e.Op {

			case token.ADD:
				return lhs + rhs, true, nil

			case token.SUB:
				return lhs - rhs, true, nil

			case token.MUL:
				return lhs * rhs, true, nil

			case token.QUO:
				if rhs == 0 {
					return 0, false, fmt.Errorf("illegal operation: division by zero")
				}
				return lhs / rhs, true, nil

			case token.REM:
				if rhs == 0 {
					return 0, false, fmt.Errorf("illegal operation: modulo by zero")
				}
				return math.Mod(lhs, rhs), true, nil
			}
		}

	case *ast.ParenExpr:
		return constant_folder.EvaluateExpression(e.X)
	}

	return 0, false, nil
}

// Name: StructureConstant
//
// Parameters: float64
//
// Return: *ast.BasicLit
//
// Converts and formats int and float constants to string
func StructureConstant(value float64) *ast.BasicLit {

	if value == float64(int(value)) {
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: strconv.Itoa(int(value)),
		}
	}

	str_dec := strconv.FormatFloat(value, 'f', 5, 64)
	str_dec = strings.TrimRight(strings.TrimRight(str_dec, "0"), ".")

	return &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: str_dec,
	}
}

/* PerformDeadCodeElimination Helper Functions */

// Name:RemoveUnusedFunction
//
// Parameters:  *ast.File, *types.Info
//
// Return: error
//
// Determines if the statement contains any unused functions and removes them from the AST
func RemoveUnusedFunctions(ast_data AstData) error {
	optimised_functions := []ast.Decl{}
	used_functions := make(map[string]bool)
	used_functions["main"] = true
	ast.Inspect(ast_data.ast_file, func(node ast.Node) bool {
		function_call, valid_call := node.(*ast.CallExpr)
		if !valid_call {
			return true
		}

		ident, valid_ident := function_call.Fun.(*ast.Ident)
		if valid_ident {
			used_function := ast_data.ast_info.Uses[ident]
			if used_function != nil {
				function, valid_function := used_function.(*types.Func)
				if valid_function {
					used_functions[function.Name()] = true
				}
			}
		}

		return true
	})

	for _, declaration := range ast_data.ast_file.Decls {
		function, valid_function := declaration.(*ast.FuncDecl)
		if valid_function {
			_, is_called := used_functions[function.Name.Name]
			if is_called {
				err := EliminateFunctionDeadCode(function, ast_data)
				if err != nil {
					return err
				}
				optimised_functions = append(optimised_functions, function)
			}
		} else {
			optimised_functions = append(optimised_functions, declaration)
		}
	}

	ast_data.ast_file.Decls = optimised_functions

	return nil
}

// Name: EliminateFunctionDeadCode
//
// Parameters: *ast.FuncDecl, *types.Info
//
// Return: error
//
// Performs dead code elimination on an individual function
func EliminateFunctionDeadCode(function *ast.FuncDecl, ast_data AstData) error {
	if function.Body == nil {
		return fmt.Errorf("function %v has no body", function.Name)
	}

	optimised_statements := []ast.Stmt{}
	unreachable := false

	for _, function_statements := range function.Body.List {

		if unreachable {
			continue
		}

		SearchStructureBody(function_statements, &optimised_statements, ast_data, function, &unreachable)

	}

	function.Body.List = optimised_statements

	return nil
}

// Name:RemoveUnusedAssignedVariables
//
// Parameters: *ast.AssignStmt, *types.Info, *[]ast.Stmt
//
// Return:
//
// Determines if the statement contains any unused assigned variables and removes them from the AST
func RemoveUnusedAssignedVariables(statement *ast.AssignStmt, ast_data AstData, optimised_statements *[]ast.Stmt) {
	if statement.Tok == token.DEFINE {
		new_lhs := []ast.Expr{}
		new_rhs := []ast.Expr{}
		for num, lhs := range statement.Lhs {
			ident, valid_ident := lhs.(*ast.Ident)
			if valid_ident {
				variable := ast_data.ast_info.Defs[ident]
				used := false

				for _, used_variable := range ast_data.ast_info.Uses {
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
			*optimised_statements = append(*optimised_statements, statement)

		}
	} else {
		*optimised_statements = append(*optimised_statements, statement)
	}
}

// Name:RemoveUnusedDeclaredVariables
//
// Parameters: *ast.AssignStmt, *types.Info, *[]ast.Stmt
//
// Return:
//
// Determines if the statement contains any unused assigned variables and removes them from the AST
func RemoveUnusedDeclaredVariables(statement *ast.DeclStmt, ast_data AstData, optimised_statements *[]ast.Stmt) {
	declaration_info, valid_declaration := statement.Decl.(*ast.GenDecl)
	if valid_declaration && declaration_info.Tok == token.VAR {
		new_spec := []ast.Spec{}
		for _, spec := range declaration_info.Specs {
			value, valid_spec := spec.(*ast.ValueSpec)
			if valid_spec {
				new_names := []*ast.Ident{}
				new_values := []ast.Expr{}
				for num, name := range value.Names {
					variable := ast_data.ast_info.Defs[name]

					used := false
					for _, used_variable := range ast_data.ast_info.Uses {
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
			*optimised_statements = append(*optimised_statements, statement)
		}
	} else {
		*optimised_statements = append(*optimised_statements, statement)
	}
}

// Name:RemoveUnusedIfStatement
//
// Parameters:   *ast.FuncDecl, *[]ast.Stmt, AstData
//
// Return:
//
// Determines if the statement contains any unused if statements and removes them from the AST
func RemoveUnusedIfStatement(function_statements ast.Stmt, optimised_statements *[]ast.Stmt, ast_data AstData, function *ast.FuncDecl) {
	if_statement, valid_ifstatement := function_statements.(*ast.IfStmt)
	unreachable := false
	var optimised_if_body []ast.Stmt
	if valid_ifstatement {
		switch condition_variable := if_statement.Cond.(type) {

		case *ast.Ident:
			cond_variable, valid_cond := if_statement.Cond.(*ast.Ident)
			if valid_cond {
				_, exists := ast_data.ast_info.Uses[cond_variable]
				if exists {
					variable_value, def_exists := ast_data.variable_values[cond_variable.Name]
					if def_exists {
						bool_value, is_bool := variable_value.(constant.Value)
						if is_bool && bool_value.Kind() == constant.Bool && constant.BoolVal(bool_value) {
							for _, body_statement := range if_statement.Body.List {
								if unreachable {
									continue
								}

								SearchStructureBody(body_statement, &optimised_if_body, ast_data, function, &unreachable)
							}
							if_statement.Body.List = optimised_if_body
							*optimised_statements = append(*optimised_statements, if_statement)
						}
					} else {
						if cond_variable.Name == "true" {
							for _, body_statement := range if_statement.Body.List {
								if unreachable {
									continue
								}

								SearchStructureBody(body_statement, &optimised_if_body, ast_data, function, &unreachable)
							}
							if_statement.Body.List = optimised_if_body
							*optimised_statements = append(*optimised_statements, if_statement)
						}
					}
				}
			}

		case *ast.BinaryExpr:
			valid_rhs := false
			valid_lhs := false
			var rhs_value constant.Value
			var lhs_value constant.Value

			switch lhs := condition_variable.X.(type) {
			case *ast.Ident:
				_, exists := ast_data.ast_info.Uses[lhs]
				if exists {
					lhs_val, def_exists := ast_data.variable_values[lhs.Name]
					if def_exists {
						lhs_constant, is_constant := lhs_val.(constant.Value)
						if is_constant {
							valid_lhs = true
							lhs_value = lhs_constant
						}

					}
				}
			case *ast.BasicLit:
				lhs_constant, err := ConvertToConstant(lhs)
				if err != nil {
					return
				}
				lhs_value = lhs_constant
				valid_lhs = true
			}
			switch rhs := condition_variable.Y.(type) {
			case *ast.Ident:
				_, exists := ast_data.ast_info.Uses[rhs]
				if exists {
					rhs_val, def_exists := ast_data.variable_values[rhs.Name]
					if def_exists {
						rhs_constant, is_constant := rhs_val.(constant.Value)
						if is_constant {
							valid_rhs = true
							rhs_value = rhs_constant
						}
					}
				}
			case *ast.BasicLit:
				rhs_constant, err := ConvertToConstant(rhs)
				if err != nil {
					return
				}
				rhs_value = rhs_constant
				valid_rhs = true
			}

			if valid_lhs && valid_rhs {
				operator := condition_variable.Op.String()
				switch operator {
				case "==":
					if constant.Compare(lhs_value, token.NEQ, rhs_value) {
						return
					}

				case "!=":
					if constant.Compare(lhs_value, token.EQL, rhs_value) {
						return
					}

				case ">=":
					if constant.Compare(lhs_value, token.LSS, rhs_value) {
						return
					}

				case "<=":
					if constant.Compare(lhs_value, token.GTR, rhs_value) {
						return
					}

				case ">":
					if constant.Compare(lhs_value, token.LEQ, rhs_value) {
						return
					}

				case "<":
					if constant.Compare(lhs_value, token.GEQ, rhs_value) {
						return
					}

				}
			}

			for _, body_statement := range if_statement.Body.List {
				if unreachable {
					continue
				}

				SearchStructureBody(body_statement, &optimised_if_body, ast_data, function, &unreachable)
			}
			if_statement.Body.List = optimised_if_body
			*optimised_statements = append(*optimised_statements, if_statement)

		case *ast.UnaryExpr:
			operator := condition_variable.Op.String()
			cond_variable, valid_cond := condition_variable.X.(*ast.Ident)
			if valid_cond {
				_, exists := ast_data.ast_info.Uses[cond_variable]
				if exists {
					variable_value, def_exists := ast_data.variable_values[cond_variable.Name]
					if def_exists {
						bool_value, is_bool := variable_value.(constant.Value)
						switch operator {
						case "!":
							bool_value := constant.BoolVal(bool_value)
							bool_value = !bool_value
							if is_bool && bool_value {

								for _, body_statement := range if_statement.Body.List {
									if unreachable {
										continue
									}

									SearchStructureBody(body_statement, &optimised_if_body, ast_data, function, &unreachable)
								}
								if_statement.Body.List = optimised_if_body
								*optimised_statements = append(*optimised_statements, if_statement)
							}
						}

					} else {
						if cond_variable.Name == "false" {
							for _, body_statement := range if_statement.Body.List {
								if unreachable {
									continue
								}

								SearchStructureBody(body_statement, &optimised_if_body, ast_data, function, &unreachable)
							}
							if_statement.Body.List = optimised_if_body
							*optimised_statements = append(*optimised_statements, if_statement)
						}
					}
				}
			}

		default:
			*optimised_statements = append(*optimised_statements, function_statements)
		}
	}
}

func ConvertToConstant(lit *ast.BasicLit) (constant.Value, error) {
	switch lit.Kind {
	case token.INT:
		return constant.MakeFromLiteral(lit.Value, token.INT, 0), nil
	case token.FLOAT:
		return constant.MakeFromLiteral(lit.Value, token.FLOAT, 0), nil
	case token.CHAR:
		return constant.MakeFromLiteral(lit.Value, token.CHAR, 0), nil
	case token.STRING:
		return constant.MakeFromLiteral(lit.Value, token.STRING, 0), nil
	default:
		return nil, fmt.Errorf("unexpected kind: %v", lit.Kind)
	}
}

// Name:RemoveUnusedForStatement
//
// Parameters:  map[string]string, *ast.FuncDecl, *[]ast.Stmt, AstData
//
// Return:
//
// Determines if the statement contains any unused for statements and removes them from the AST
func RemoveUnusedForStatement(function_statements ast.Stmt, optimised_statements *[]ast.Stmt, ast_data AstData, function *ast.FuncDecl) {
	for_statement, valid_forstatement := function_statements.(*ast.ForStmt)
	unreachable := false
	var optimised_for_body []ast.Stmt

	if valid_forstatement {

		switch condition_type := for_statement.Cond.(type) {
		case *ast.BinaryExpr:
			_, valid_expr := condition_type.X.(*ast.Ident)
			control_restriction, valid_restriction := condition_type.Y.(*ast.BasicLit)
			control_operator := condition_type.Op.String()
			if valid_expr && valid_restriction {
				for_restriction := control_restriction.Value
				assignment_statement, valid_assignment := for_statement.Init.(*ast.AssignStmt)
				if valid_assignment {
					assign, valid_assign := assignment_statement.Rhs[0].(*ast.BasicLit)
					if valid_assign {

						switch control_operator {
						case ">":
							if assign.Value <= for_restriction {
								return
							}
						case ">=":
							if assign.Value < for_restriction {
								return
							}
						case "<":
							if assign.Value >= for_restriction {
								return
							}
						case "<=":
							if assign.Value > for_restriction {
								return
							}
						case "!=":
							if assign.Value == for_restriction {
								return
							}
						case "==":
							if assign.Value != for_restriction {
								return
							}
						default:
							return
						}
						for _, body_statement := range for_statement.Body.List {
							if unreachable {
								continue
							}

							SearchStructureBody(body_statement, &optimised_for_body, ast_data, function, &unreachable)
						}
						for_statement.Body.List = optimised_for_body
						*optimised_statements = append(*optimised_statements, for_statement)

					}
				}
			}

		default:
			*optimised_statements = append(*optimised_statements, function_statements)
		}

	}

}

// Name:RemoveUnusedSwitchStatement
//
// Parameters:  map[string]string, *ast.FuncDecl, *[]ast.Stmt, AstData
//
// Return:
//
// Determines if the statement contains any unused switch statements and removes them from the AST
func RemoveUnusedSwitchStatement(function_statements ast.Stmt, optimised_statements *[]ast.Stmt, ast_data AstData, function *ast.FuncDecl) {
	switch_statement, valid_switchstatement := function_statements.(*ast.SwitchStmt)
	unreachable := false
	var optimised_body []ast.Stmt

	if valid_switchstatement {
		if switch_statement.Body.List == nil {
			return
		}

		for _, body_statement := range switch_statement.Body.List {
			if unreachable {
				continue
			}

			switch statement := body_statement.(type) {
			case *ast.CaseClause:
				RemoveUnusedSwitchCase(statement, &optimised_body, ast_data, function)
			}
		}
		if optimised_body != nil {
			switch_statement.Body.List = optimised_body
			*optimised_statements = append(*optimised_statements, switch_statement)
		}

	}

}

// Name:RemoveUnusedSwitchCase
//
// Parameters:  map[string]string, *ast.FuncDecl, *[]ast.Stmt, AstData
//
// Return:
//
// Determines if the statement contains any unused switch case statements and removes them from the AST
func RemoveUnusedSwitchCase(function_statements ast.Stmt, optimised_statements *[]ast.Stmt, ast_data AstData, function *ast.FuncDecl) {
	switch_statement, valid_switchstatement := function_statements.(*ast.CaseClause)
	unreachable := false
	var optimised_body []ast.Stmt

	if valid_switchstatement {
		if switch_statement.Body == nil {
			return
		}

		for _, body_statement := range switch_statement.Body {
			if unreachable {
				continue
			}

			SearchStructureBody(body_statement, &optimised_body, ast_data, function, &unreachable)
		}
		switch_statement.Body = optimised_body
		*optimised_statements = append(*optimised_statements, switch_statement)

	}

}

// Name:SearchStructureBody
//
// Parameters:   ast.Stmt, *[]ast.Stmt, AstData, *ast.FuncDecl,*bool
//
// Return:
//
// Determines type of statement and performs the necessary dead code elimination
func SearchStructureBody(body_statement ast.Stmt, optimised_body *[]ast.Stmt, ast_data AstData, function *ast.FuncDecl, unreachable *bool) {
	switch statement := body_statement.(type) {
	case *ast.AssignStmt: // search for unused assigned variables
		RemoveUnusedAssignedVariables(statement, ast_data, optimised_body)
	case *ast.DeclStmt: //search for unused declared variables
		RemoveUnusedDeclaredVariables(statement, ast_data, optimised_body)
	case *ast.IfStmt:
		RemoveUnusedIfStatement(statement, optimised_body, ast_data, function)
	case *ast.ForStmt:
		RemoveUnusedForStatement(statement, optimised_body, ast_data, function)
	case *ast.SwitchStmt:
		RemoveUnusedSwitchStatement(statement, optimised_body, ast_data, function)

	case *ast.ReturnStmt, *ast.BranchStmt, *ast.GoStmt, *ast.DeferStmt:
		*optimised_body = append(*optimised_body, statement)
		*unreachable = true

	default:
		*optimised_body = append(*optimised_body, statement)
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

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

	ast, file_set, err := ParseGoCode(code)
	if err != nil {
		return "", err
	}

	if constant_folding {
		PerformConstantFolding(ast)
	}
	if dead_code {
		PerformDeadCodeElimination(ast)
	}
	if loop_unroling {
		PerformLoopUnrolling(ast)
	}

	code, err = StringifyAST(ast, file_set)
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
	ast, err := parser.ParseFile(file_set, "", code, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	return ast, file_set, nil
}

// Name: StringifyAST
//
// Parameters: *ast.File, *token.FileSet
//
// Return: string
//
// Converts optimised AST to a string
func StringifyAST(ast *ast.File, file_set *token.FileSet) (string, error) {
	var buffer bytes.Buffer
	err := format.Node(&buffer, file_set, ast)
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
func PerformConstantFolding(ast *ast.File) {

}

// Name: PerformDeadCodeElimination
//
// Parameters: *ast.File
//
// Return: string
//
// Perform dead code elimination on the source code
func PerformDeadCodeElimination(ast *ast.File) {

}

func EliminateDeadCode(fn *ast.FuncDecl, info *types.Info) {

}

// Name: PerformLoopUnrolling
//
// Parameters: *ast.File
//
// Return: string
//
// Perform loop unrolling on the source code
func PerformLoopUnrolling(ast *ast.File) {

}

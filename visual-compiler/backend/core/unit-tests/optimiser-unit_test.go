package unit_tests

import (
	"go/token"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestParseCode_Success(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "}"

	ast, _, err := services.ParseGoCode(code)
	if err != nil {
		t.Errorf("Error :%v", err)
	}
	if ast == nil {
		t.Errorf("AST is empty")
	}
}

func TestParseCode_NoPackage(t *testing.T) {
	code := "num := 2"

	_, _, err := services.ParseGoCode(code)
	if err == nil {
		t.Errorf("Error expected:  expected 'package'")
	}
}

func TestStringifyAST_Success(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	ast, file_set, err := services.ParseGoCode(code)
	if err != nil {
		t.Errorf("Error :%v", err)
	}
	if ast == nil {
		t.Errorf("AST is empty")
	}

	code, err = services.StringifyAST(ast, file_set)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if code != expected_result {
		t.Errorf("Stringify AST failed:\n %v\n %v", code, expected_result)
	}
}

func TestStringifyAST_Fail(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "}"

	ast, _, err := services.ParseGoCode(code)
	file_set := token.NewFileSet()
	if err != nil {
		t.Errorf("Error :%v", err)
	}
	if ast == nil {
		t.Errorf("AST is empty")
	}

	_, err = services.StringifyAST(ast, file_set)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

}

func TestOptimiseGo_NoInputCode(t *testing.T) {
	code := ""

	_, err := services.OptimiseGoCode(code, true, true, true)
	if err == nil {
		t.Errorf("Error expected")
	}
}

func TestPerformDeadCodeElimination_Simple(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "random_num := 5\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\treturn 13\n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_MultipleReturn(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\treturn 13\n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableAfterReturn(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "random_num := 5\n"
	code += "return 23\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\treturn 13\n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableBeforeReturn(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "random_num := 5\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableBeforeReturn_2(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "var random_num int\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ConstantVariableIfStatement(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "false_bool := false\n"
	code += "if false_bool {\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_Identifier(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "true_bool := true\n"
	code += "if true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\ttrue_bool := true\n"
	expected_result += "\tif true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "true_bool := \"halfstack\"\n"
	code += "if true_bool==\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\ttrue_bool := \"halfstack\"\n"
	expected_result += "\tif true_bool == \"halfstack\" {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "true_bool := \"halfstack\"\n"
	code += "if true_bool!=\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_Constant(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "if false{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_UnaryConstant(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "if !false{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\tif !false {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_UnaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "true_bool := true\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_UnaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\treturn 13\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v \n%v", optmised_code, expected_result)
	}
}

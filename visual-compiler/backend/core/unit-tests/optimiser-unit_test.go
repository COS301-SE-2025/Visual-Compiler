package unit_tests

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

// Tests for Lexing and Parsing

func TestParseCode_Success(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "return\n"
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

func TestOptimiseGoCode_NoPackage(t *testing.T) {
	code := "num := 2"

	_, err := services.OptimiseGoCode(code, false, true, false)
	if err == nil {
		t.Errorf("Error expected:  expected 'package'")
	}
}

func TestStringifyAST_Success(t *testing.T) {
	code := "package main\n"
	code += "func main(){\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
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
		t.Errorf("Stringify AST failed:\n%v\n%v", code, expected_result)
	}
}

func TestStringifyAST_Fail(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
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

func TestConvertToConstant_Float(t *testing.T) {

	value := &ast.BasicLit{Kind: token.FLOAT, Value: "13.2"}
	_, err := services.ConvertToConstant(value)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestConvertToConstant_Char(t *testing.T) {

	value := &ast.BasicLit{Kind: token.CHAR, Value: "C"}
	_, err := services.ConvertToConstant(value)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestConvertToConstant_Error(t *testing.T) {

	value := &ast.BasicLit{Kind: token.EQL, Value: "=="}
	_, err := services.ConvertToConstant(value)
	if err == nil {
		t.Errorf("Error expected")
	}
}

// Tests for Optimisation

func TestOptimiseGoCode_NoInputCode(t *testing.T) {
	code := ""

	_, err := services.OptimiseGoCode(code, true, true, true)
	if err == nil {
		t.Errorf("Error expected")
	}
}

func TestOptimiseGoCode_AllThreeCombined(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n\n"
	code += "func main() {\n"
	code += "\tfor blue := 1; blue < 5; blue++ {\n"
	code += "\t\tfmt.Println(\"[Block \" + fmt.Sprint(blue) + \"]\")\n"
	code += "\t\tfor red := 1; red <= blue; red++ {\n"
	code += "\t\t\tfmt.Println(\"\\t-> \" + fmt.Sprint(red))\n"
	code += "\t\t\tif blue == red {\n"
	code += "\t\t\t\tfmt.Println(blue + red)\n"
	code += "\t\t\t}\n"
	code += "\t\t}\n"
	code += "\t}\n"
	code += "}\n\n"
	code += "func nothing() int {\n"
	code += "\tvar random int = 13\n"
	code += "\treturn random\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfmt.Println(\"[Block \" + fmt.Sprint(1) + \"]\")\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(1))\n"
	expected_result += "\tif 1 == 1 {\n"
	expected_result += "\t\tfmt.Println(2)\n"
	expected_result += "\t}\n"
	expected_result += "\tfmt.Println(\"[Block \" + fmt.Sprint(2) + \"]\")\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(1))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(2))\n"
	expected_result += "\tif 2 == 2 {\n"
	expected_result += "\t\tfmt.Println(4)\n"
	expected_result += "\t}\n"
	expected_result += "\tfmt.Println(\"[Block \" + fmt.Sprint(3) + \"]\")\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(1))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(2))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(3))\n"
	expected_result += "\tif 3 == 3 {\n"
	expected_result += "\t\tfmt.Println(6)\n"
	expected_result += "\t}\n"
	expected_result += "\tfmt.Println(\"[Block \" + fmt.Sprint(4) + \"]\")\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(1))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(2))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(3))\n"
	expected_result += "\tfmt.Println(\"\\t-> \" + fmt.Sprint(4))\n"
	expected_result += "\tif 4 == 4 {\n"
	expected_result += "\t\tfmt.Println(8)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, true, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Optimisation Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

// Tests for Constant Folding

func TestPerformConstantFolding_SimpleExpression(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := 12 + 13\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tresult := 25\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_ComplexExpression(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := (12 + 13) * (28 % 10)\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tresult := 200\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_BasicIntegerArithmetic(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tnum1 := 20\n"
	code += "\tnum2 := 10\n"
	code += "\tadd := num1 + num2\n"
	code += "\tsub := num1 - num2\n"
	code += "\tmul := num1 * num2\n"
	code += "\tdiv := num1 / num2\n"
	code += "\tmod := num1 % num2\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tnum1 := 20\n"
	expected_result += "\tnum2 := 10\n"
	expected_result += "\tadd := 30\n"
	expected_result += "\tsub := 10\n"
	expected_result += "\tmul := 200\n"
	expected_result += "\tdiv := 2\n"
	expected_result += "\tmod := 0\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_BasicFloatArithmetic(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tnum1 := 121.3\n"
	code += "\tnum2 := 19.89\n"
	code += "\tadd := num1 + num2\n"
	code += "\tsub := num1 - num2\n"
	code += "\tmul := num1 * num2\n"
	code += "\tdiv := num1 / num2\n"
	code += "\tmod := num1 % num2\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tnum1 := 121.3\n"
	expected_result += "\tnum2 := 19.89\n"
	expected_result += "\tadd := 141.19\n"
	expected_result += "\tsub := 101.41\n"
	expected_result += "\tmul := 2412.657\n"
	expected_result += "\tdiv := 6.09854\n"
	expected_result += "\tmod := 1.96\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_NegativeNumbers(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := -12 - -4\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tresult := -8\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_DivisionByZero(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := 10 / 0\n"
	code += "}\n"

	_, err := services.OptimiseGoCode(code, true, false, false)
	if err == nil {
		t.Errorf("Expected error for division by zero")
	}
}

func TestPerformConstantFolding_ModuloByZero(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := 10 % 0\n"
	code += "}\n"

	_, err := services.OptimiseGoCode(code, true, false, false)
	if err == nil {
		t.Errorf("Expected error for modulo by zero")
	}
}

func TestPerformConstantFolding_FloatToIntegerConversion(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tresult := 2.0 + 8.0\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tresult := 10\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_Reassignment(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tblue := 20\n"
	code += "\tblue = blue + 2\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tblue := 20\n"
	expected_result += "\tblue = 22\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_ChainedVariables(t *testing.T) {
	code := "package main\n\n"
	code += "func main() {\n"
	code += "\tx := 1\n"
	code += "\ty := x + 1\n"
	code += "\tz := y + 1\n"
	code += "\ttotal := x + y + z + x + y + z\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tx := 1\n"
	expected_result += "\ty := 2\n"
	expected_result += "\tz := 3\n"
	expected_result += "\ttotal := 12\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_FunctionArguments(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "\tfmt.Printf(132 / 11)\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfmt.Printf(12)\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformConstantFolding_NothingFolded(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "\tfor i := 1; i < 13; i++ {\n"
	code += "\t\tif i%2 != 1 {\n"
	code += "\t\t\tfmt.Printf(random + i)\n"
	code += "\t\t}\n"
	code += "\t}\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 1; i < 13; i++ {\n"
	expected_result += "\t\tif i%2 != 1 {\n"
	expected_result += "\t\t\tfmt.Printf(random + i)\n"
	expected_result += "\t\t}\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, true, false, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Constant Folding Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

// Tests for Dead Code Elimination

func TestPerformDeadCodeElimination_Simple(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "return\n"
	code += "random_num := 5\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_MultipleReturn(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "return\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableAfterReturn(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "return\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableBeforeReturn(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnusedVariableBeforeReturn_2(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var random_num int\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UsedVariableBeforeReturn_2(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var random_num int\n"
	code += "random_num = 5\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar random_num int\n"
	expected_result += "\trandom_num = 5\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_EmptyIfStatement(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "false_bool := true\n"
	code += "if false_bool {\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ConstantVariableIfStatement(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "false_bool := false\n"
	code += "if false_bool {\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ConstantVariableIfStatement_False(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "false_bool := false\n"
	code += "if false {\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ConstantVariableIfStatement_True(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\" \n"
	code += "func main() {\n"
	code += "false_bool := false\n"
	code += "if true {\n"
	code += "random_num := 5\n"
	code += "fmt.Printf(\"%v\",random_num)\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tif true {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\tfmt.Printf(\"%v\", random_num)\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ConstantVariableIfStatement_True2(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\" \n"
	code += "func main() {\n"
	code += "false_bool := false\n"
	code += "if true {\n"
	code += "random_num := 5\n"
	code += "random_num++\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tif true {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\trandom_num++\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_Identifier(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "true_bool := true\n"
	code += "if true_bool{\n"
	code += "fmt.Printf(\"Hi\")\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := true\n"
	expected_result += "\tif true_bool {\n"
	expected_result += "\t\tfmt.Printf(\"Hi\")\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_Equal(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := \"halfstack\"\n"
	code += "if true_bool==\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := \"halfstack\"\n"
	expected_result += "\tif true_bool == \"halfstack\" {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_Equal(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := \"half\"\n"
	code += "if true_bool==\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_NotEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := \"half\"\n"
	code += "if true_bool!=\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := \"half\"\n"
	expected_result += "\tif true_bool != \"halfstack\" {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_NotEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := \"halfstack\"\n"
	code += "if true_bool!=\"halfstack\"{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_LessEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 5\n"
	code += "if true_bool<=5{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := 5\n"
	expected_result += "\tif true_bool <= 5 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}
func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_LessEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 6\n"
	code += "if true_bool<=5{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_Less(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 4\n"
	code += "if true_bool<5{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := 4\n"
	expected_result += "\tif true_bool < 5 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}
func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_Less(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 6\n"
	code += "if true_bool<5{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_GreaterEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 5\n"
	code += "if true_bool>=5{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := 5\n"
	expected_result += "\tif true_bool >= 5 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}
func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_GreaterEqual(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 4\n"
	code += "if true_bool>=5{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_Greater(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 6\n"
	code += "if true_bool > 5{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := 6\n"
	expected_result += "\tif true_bool > 5 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}
func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_Greater(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := 4\n"
	code += "if true_bool > 5{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_EqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := \"halfstack\"\n"
	code += "var2 := \"halfstack\"\n"
	code += "if var1 == var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := \"halfstack\"\n"
	expected_result += "\tvar2 := \"halfstack\"\n"
	expected_result += "\tif var1 == var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_EqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := \"halfstack\"\n"
	code += "var2 := \"half\"\n"
	code += "if var1==var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_NotEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := \"halfstack\"\n"
	code += "var2 := \"half\"\n"
	code += "if var1!=var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := \"halfstack\"\n"
	expected_result += "\tvar2 := \"half\"\n"
	expected_result += "\tif var1 != var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_NotEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := \"halfstack\"\n"
	code += "var2 := \"halfstack\"\n"
	code += "if var1!=var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_LessEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 5\n"
	code += "var2 := 5\n"
	code += "if var1<=var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := 5\n"
	expected_result += "\tvar2 := 5\n"
	expected_result += "\tif var1 <= var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_LessEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 6\n"
	code += "var2 := 5\n"
	code += "if var1<=var2{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_LessVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 4\n"
	code += "var2 := 5\n"
	code += "if var1<var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := 4\n"
	expected_result += "\tvar2 := 5\n"
	expected_result += "\tif var1 < var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_LessVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 6\n"
	code += "var2 := 5\n"
	code += "if var1<var2{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_GreaterEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 13\n"
	code += "var2 := 5\n"
	code += "if var1 >= var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := 13\n"
	expected_result += "\tvar2 := 5\n"
	expected_result += "\tif var1 >= var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_GreaterEqualVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 1\n"
	code += "var2 := 5\n"
	code += "if var1>=var2{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_GreaterVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 13\n"
	code += "var2 := 5\n"
	code += "if var1>var2{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := 13\n"
	expected_result += "\tvar2 := 5\n"
	expected_result += "\tif var1 > var2 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_BinaryExpression_GreaterVariables(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 1\n"
	code += "var2 := 5\n"
	code += "if var1>var2{\n"
	code += "random_num := 5\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_Greater2(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "var1 := 13\n"
	code += "if 13==var1 {\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tvar1 := 13\n"
	expected_result += "\tif 13 == var1 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_BinaryExpression_GreaterConstants(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "if 13==13 {\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tif 13 == 13 {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_Constant(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "if false{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_UnaryConstant(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "if !false{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\tif !false {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedIfStatement_UnaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := true\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedIfStatement_UnaryExpression(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedForStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 10; i < 10; i++ {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 0; i < 10; i++ {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 0; i < 10; i++ {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_GTR(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 15; i > 10; i-- {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 15; i > 10; i-- {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_GEQ(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 15; i >= 10; i-- {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 15; i >= 10; i-- {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_LEQ(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 0; i <= 10; i++ {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 0; i <= 10; i++ {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_NEQ(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 15; i != 10; i-- {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 15; i != 10; i-- {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_EQL(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 10; i == 10; i-- {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 10; i == 10; i-- {\n"
	expected_result += "\t\tfmt.Printf(\"%v\", i)\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedForStatement_Other(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 15; i * 10; i-- {\n"
	code += "fmt.Printf(\"%v\",i)\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_EmptyForStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 0; i < 10; i++ {\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 0; i < 10; i++ {\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_EmptySwitchStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "for i := 0; i < 10; i++ {\n"
	code += "switch i {\n"
	code += "}\n"
	code += "}\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfor i := 0; i < 10; i++ {\n"
	expected_result += "\t}\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedSwitchStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "sum := 0\n"
	code += "for i := 0; i < 10; i++ {\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "sum++\n"
	code += "default: \n"
	code += "sum+=2\n"
	code += "}\n"
	code += "}\n"
	code += "fmt.Printf(\"%v\",sum)\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tsum := 0\n"
	expected_result += "\tfor i := 0; i < 10; i++ {\n"
	expected_result += "\t\tswitch i {\n"
	expected_result += "\t\tcase 0:\n"
	expected_result += "\t\t\tsum++\n"
	expected_result += "\t\tdefault:\n"
	expected_result += "\t\t\tsum += 2\n"
	expected_result += "\t\t}\n"
	expected_result += "\t}\n"
	expected_result += "\tfmt.Printf(\"%v\", sum)\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedSwitchStatement(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "sum := 0\n"
	code += "i := 0\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "default: \n"
	code += "}\n"
	code += "fmt.Printf(\"%v\",sum)\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tsum := 0\n"
	expected_result += "\tfmt.Printf(\"%v\", sum)\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}
func TestPerformDeadCodeElimination_UnreachedSwitchStatement_InFor(t *testing.T) {
	code := "package main\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "sum := 0\n"
	code += "for i := 0; i < 10; i++ {\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "default: \n"
	code += "}\n"
	code += "}\n"
	code += "fmt.Printf(\"%v\",sum)\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tsum := 0\n"
	expected_result += "\tfor i := 0; i < 10; i++ {\n"
	expected_result += "\t}\n"
	expected_result += "\tfmt.Printf(\"%v\", sum)\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_UnreachedFunction(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}\n"
	code += "func other_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "return 13\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedFunction(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}\n"
	code += "func other_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\tother_function()\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"
	expected_result += "func other_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedFunction_Empty(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}\n"
	code += "func other_function(){\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\tother_function()\n"
	expected_result += "\t\treturn\n"
	expected_result += "\t}\n"
	expected_result += "\treturn\n"
	expected_result += "}\n"
	expected_result += "func other_function() {\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_ReachedFunction_EmptyBody(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "return\n"
	code += "}\n"
	code += "return\n"
	code += "}\n"
	code += "func other_function()"

	_, err := services.OptimiseGoCode(code, false, true, false)
	if err == nil {
		t.Errorf("Error: %v", err)
	}
}

func TestPerformDeadCodeElimination_Complex(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "}else{\n"
	code += "else_function()\n"
	code += "}\n"
	code += "}\n"
	code += "func other_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}\n"
	code += "func else_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\tother_function()\n"
	expected_result += "\t} else {\n"
	expected_result += "\t\telse_function()\n"
	expected_result += "\t}\n"
	expected_result += "}\n"
	expected_result += "func other_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"
	expected_result += "func else_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_Complex_2(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "i:=0\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "i++\n"
	code += "default:\n"
	code += "i+=2\n"
	code += "}\n"
	code += "}else{\n"
	code += "else_function()\n"
	code += "}\n"
	code += "}\n"
	code += "func other_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}\n"
	code += "func else_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\tother_function()\n"
	expected_result += "\t\ti := 0\n"
	expected_result += "\t\tswitch i {\n"
	expected_result += "\t\tcase 0:\n"
	expected_result += "\t\t\ti++\n"
	expected_result += "\t\tdefault:\n"
	expected_result += "\t\t\ti += 2\n"
	expected_result += "\t\t}\n"
	expected_result += "\t} else {\n"
	expected_result += "\t\telse_function()\n"
	expected_result += "\t}\n"
	expected_result += "}\n"
	expected_result += "func other_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"
	expected_result += "func else_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

func TestPerformDeadCodeElimination_Complex_3(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "true_bool := false\n"
	code += "if !true_bool {\n"
	code += "other_function()\n"
	code += "i:=0\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "i++\n"
	code += "default:\n"
	code += "i+=2\n"
	code += "}\n"
	code += "}else{\n"
	code += "else_function()\n"
	code += "}\n"
	code += "if true_bool {\n"
	code += "other_function()\n"
	code += "i:=0\n"
	code += "switch i {\n"
	code += "case 0:\n"
	code += "i++\n"
	code += "default:\n"
	code += "i+=2\n"
	code += "}\n"
	code += "}\n"
	code += "}\n"
	code += "func other_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}\n"
	code += "func else_function() int {\n"
	code += "true_bool := false\n"
	code += "if !true_bool{\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}\n\n"
	code += "random_n := 13\n"
	code += "return random_n\n"
	code += "}"

	expected_result := "package main\n\n"
	expected_result += "func main() {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\tother_function()\n"
	expected_result += "\t\ti := 0\n"
	expected_result += "\t\tswitch i {\n"
	expected_result += "\t\tcase 0:\n"
	expected_result += "\t\t\ti++\n"
	expected_result += "\t\tdefault:\n"
	expected_result += "\t\t\ti += 2\n"
	expected_result += "\t\t}\n"
	expected_result += "\t} else {\n"
	expected_result += "\t\telse_function()\n"
	expected_result += "\t}\n"
	expected_result += "}\n"
	expected_result += "func other_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"
	expected_result += "func else_function() int {\n"
	expected_result += "\ttrue_bool := false\n"
	expected_result += "\tif !true_bool {\n"
	expected_result += "\t\trandom_num := 5\n"
	expected_result += "\t\treturn random_num\n"
	expected_result += "\t}\n"
	expected_result += "\trandom_n := 13\n"
	expected_result += "\treturn random_n\n"
	expected_result += "}\n"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Dead Code Elimination Failed: \n%v \n%v", optmised_code, expected_result)
	}
}

// Tests for Loop Unrolling

func TestPerformLoopUnrolling_BasicForLoop(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "\tfor i := 0; i < 3; i++ {\n"
	code += "\t\tfmt.Printf(i)\n"
	code += "\t}\n"
	code += "}\n"

	expected_result := "package main\n\n"
	expected_result += "import \"fmt\"\n"
	expected_result += "func main() {\n"
	expected_result += "\tfmt.Printf(0)\n"
	expected_result += "\tfmt.Printf(1)\n"
	expected_result += "\tfmt.Printf(2)\n"
	expected_result += "}\n"

	optimised_code, err := services.OptimiseGoCode(code, false, false, true)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optimised_code != expected_result {
		t.Errorf("Loop Unrolling Failed: \n%v \n%v", optimised_code, expected_result)
	}
}

func TestPerformLoopUnrolling_NonStandardLoop1(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "\ti := 1\n"
	code += "\tfor i < 13 {\n"
	code += "\t\tfmt.Printf(i)\n"
	code += "\t\ti++\n"
	code += "\t}\n"
	code += "}\n"

	_, err := services.OptimiseGoCode(code, false, false, true)
	if err == nil {
		t.Errorf("Expected error for invalid loop structure")
	}
}

func TestPerformLoopUnrolling_NonStandardLoop2(t *testing.T) {
	code := "package main\n\n"
	code += "import \"fmt\"\n"
	code += "func main() {\n"
	code += "\ti := 1\n"
	code += "\tfor {\n"
	code += "\t\ti++\n"
	code += "\t\tif i==13 {\n"
	code += "\t\t\tbreak\n"
	code += "\t\t}\n"
	code += "\t}\n"
	code += "}\n"

	_, err := services.OptimiseGoCode(code, false, false, true)
	if err == nil {
		t.Errorf("Expected error for invalid loop structure")
	}
}
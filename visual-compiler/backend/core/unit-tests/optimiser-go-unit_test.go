package unit_tests

import (
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
)

func TestOptimiseGo_NoInputCode(t *testing.T) {
	code := ""

	_, err := services.OptimiseGoCode(code, true, true, true)
	if err == nil {
		t.Errorf("Error expected")
	}
}

func TestPerformDeadCodeElimination_Simple(t *testing.T) {
	code := "func NewFunction *SymbolTable {"
	code += "return &SymbolTable{"
	code += "SymbolScopes: []map[string]Symbol{"
	code += "	make(map[string]Symbol),"
	code += "	},"
	code += "}"
	code += "random_string:= ''"
	code += "}"

	expected_result := "func NewFunction *SymbolTable {"
	expected_result += "return &SymbolTable{"
	expected_result += "SymbolScopes: []map[string]Symbol{"
	expected_result += "	make(map[string]Symbol),"
	expected_result += "	},"
	expected_result += "}"
	expected_result += "}"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v", optmised_code)
	}
}

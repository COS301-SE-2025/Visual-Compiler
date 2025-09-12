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
	code := "package main\n"
	code += "func NewFunction() int {\n"
	code += "return 13\n"
	code += "random_num := 5\n"
	code += "return random_num\n"
	code += "}"

	expected_result := "package main\n"
	expected_result += "func NewFunction() int {\n"
	expected_result += "return 13\n"
	expected_result += "}"

	optmised_code, err := services.OptimiseGoCode(code, false, true, false)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if optmised_code != expected_result {
		t.Errorf("Optimisation failed : \n %v", optmised_code)
	}
}

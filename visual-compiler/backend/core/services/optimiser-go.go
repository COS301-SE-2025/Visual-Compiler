package services

import "fmt"

func OptimiseGoCode(code string, constant_folding bool, dead_code bool, loop_unroling bool) (string, error) {

	if code == "" {
		return "", fmt.Errorf("code is empty")
	}

	if constant_folding {
		PerformConstantFolding(&code)
	}
	if dead_code {
		PerformDeadCodeElimination(&code)
	}
	if loop_unroling {
		PerformLoopUnrolling(&code)
	}

	return code, nil
}

func PerformConstantFolding(code *string) string {
	return *code
}

func PerformDeadCodeElimination(code *string) string {
	return *code
}

func PerformLoopUnrolling(code *string) string {
	return *code
}

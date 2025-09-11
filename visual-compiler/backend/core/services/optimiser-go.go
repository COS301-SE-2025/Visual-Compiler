package services

import "fmt"

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

	if constant_folding {
		code = PerformConstantFolding(code)
	}
	if dead_code {
		code = PerformDeadCodeElimination(code)
	}
	if loop_unroling {
		code = PerformLoopUnrolling(code)
	}

	return code, nil
}

// Name: PerformConstantFolding
//
// Parameters: *string
//
// Return: string
//
// Perform constant folding on the source code
func PerformConstantFolding(code string) string {
	return code
}

// Name: PerformDeadCodeElimination
//
// Parameters: *string
//
// Return: string
//
// Perform dead code elimination on the source code
func PerformDeadCodeElimination(code string) string {

	return code
}

// Name: PerformLoopUnrolling
//
// Parameters: *string
//
// Return: string
//
// Perform loop unrolling on the source code
func PerformLoopUnrolling(code string) string {
	return code
}

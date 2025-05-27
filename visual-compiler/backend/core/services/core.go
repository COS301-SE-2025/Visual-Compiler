// Routines for compiler construction phases
package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// String variable that stores the source code
var source string

// Struct for the type and regex pairs
type TypeRegex struct {
	Type  string
	Regex string
}

// Array that stores the user's regex rules
var rules []TypeRegex

// Struct for the tokens
type TypeValue struct {
	Type  string
	Value string
}

// Array that stores the user's generated tokens
var tokens []TypeValue

// Array that stores unidentified input
var tokens_unidentified []string

// Name: SourceCode
// Parameters: <string> which will be source code
// Return: none
// Setter function which sets the source variable to the value of parameter
func SourceCode(data string) {
	source = data
}

// Name: GetSourceCode
// Parameters: none
// Return: <string>
// Getter function to get the value of the source variable
func GetSourceCode() string {
	return source
}

// Name: ReadRegexRules
// Parameters: []byte
// Return: error
// Receive an array of regex rules, validate them and store them in the rules struct
func ReadRegexRules(input []byte) error {

	rules = []TypeRegex{}

	var raw_pairs []struct {
		Type  string
		Regex string
	}

	err := json.Unmarshal(input, &raw_pairs)

	if err != nil {
		return fmt.Errorf("tokens_unidentified JSON input: %w", err)
	}

	for _, pair := range raw_pairs {

		_, err := regexp.Compile(pair.Regex)

		if err != nil {
			return fmt.Errorf("tokens_unidentified regex: %w", err)
		}

		rules = append(rules, TypeRegex{Type: strings.ToUpper(pair.Type), Regex: pair.Regex})
	}

	return nil
}

// Name: CreateTokens
// Parameters: none
// Return: none
// Loop through the source code to find all tokens that match the regex rules stored
func CreateTokens() {

	tokens = []TypeValue{}
	tokens_unidentified = []string{}

	re := regexp.MustCompile(`\s+`)
	input := re.ReplaceAllString(source, "")

	length := len(input)
	i := 0

	for i < length {

		found := false

		for _, rule := range rules {

			re := regexp.MustCompile("^" + rule.Regex)
			position := re.FindStringIndex(input[i:])

			if position != nil && position[0] == 0 {

				match := input[i : i+position[1]]
				tokens = append(tokens, TypeValue{Type: rule.Type, Value: match})
				i = i + position[1]
				found = true
				break
			}
		}

		if !found {

			tokens_unidentified = append(tokens_unidentified, string(input[i]))
			i++
		}
	}
}

// Name: GetTokens
// Parameters: none
// Return: <[]TypeValue>
// Getter function that returns the array of tokens
func GetTokens() []TypeValue {
	return tokens
}

// Name: GetInvalidInput
// Parameters: none
// Return: <[]string>
// Getter function that returns the array of unidentified input
func GetInvalidInput() []string {
	return tokens_unidentified
}

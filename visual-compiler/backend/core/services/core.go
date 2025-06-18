package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// =========== //
// SOURCE CODE //
// =========== //

// String variable that stores the source code
var source string

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

// ============ //
// TOKENISATION //
// ============ //

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

	var builder strings.Builder

	for i := 0; i < len(source); i++ {

		r := rune(source[i])

		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == '_' {

			builder.WriteRune(r)

		} else {

			builder.WriteRune(' ')
			builder.WriteRune(r)

			if i+1 < len(source) {

				i++
				r = rune(source[i])

				if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == '_') {

					builder.WriteRune(r)
				}
			}

			builder.WriteRune(' ')
		}
	}

	var words = strings.Fields(builder.String())

	for _, word := range words {

		found := false

		for _, rule := range rules {

			re := regexp.MustCompile("^" + rule.Regex + "$")

			if re.MatchString(word) {
				found = true
				tokens = append(tokens, TypeValue{Type: rule.Type, Value: word})
				break
			}
		}

		if !found {
			tokens_unidentified = append(tokens_unidentified, word)
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

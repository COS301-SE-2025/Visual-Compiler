//Package path implements routines for the phases of compiler construction
//This packages handles the lexing process for compiler construction
package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

//String variable that stores the source code
var source string

//Struct that indicates how the RegexRules
type TypeRegex struct {
	Type  string
	Regex string
}

//Rules array variable of type TypeRegex
//Stores each regex rule 
var rules []TypeRegex

//Struct that stores the Type and Value of a token
type TypeValue struct {
	Type  string
	Value string
}

//Tokens array to store each Type and the Value for each token
var tokens []TypeValue

//Unexpected tokens array to store each Type and the Value for each token that was not identified
var unexpected_tokens []string

//Name: SourceCode
//Paramters: <string> which will be source code
//Return: none
//Setter function which sets the source variable to the value of parameter
func SourceCode(data string) {

	source = data

}

//Name: GetSourceCode
//Parameters: none
//Return: <string>
//Getter function to get the value of the source variable
func GetSourceCode() string {

	return source

}

//Name: ReadRegexRules
//Parameters: []byte
//Return: error
//This function will receive an array of regex rules, validate them and store them in the rules struct
func ReadRegexRules(input []byte) error {

	rules = []TypeRegex{}

	var raw_pairs []struct {
		Type  string
		Regex string
	}

	err := json.Unmarshal(input, &raw_pairs)

	if err != nil {
		return fmt.Errorf("invalid JSON input: %w", err)
	}

	for _, pair := range raw_pairs {

		_, err := regexp.Compile(pair.Regex)

		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		rules = append(rules, TypeRegex{Type: strings.ToUpper(pair.Type), Regex: pair.Regex})
	}

	return nil
}

//Name: CreateTokens
//Parameters: none
//Return: <[]TypeValue> An array of tokens
//This function will loop through the source code to find all tokens that match the regex rules stored,
//each match will be stored in the tokens array
//Unexpected tokens will also be 
func CreateTokens() []TypeValue {

	unexpected_tokens = []string{}
	tokens = []TypeValue{}
	var words = strings.Fields(source)

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
			fmt.Printf("unexpected token type: %q\n", word)
			unexpected_tokens = append(unexpected_tokens,word)
		}
	}

	return tokens
}

//Name: GetUnexpectedTokens
//Parameters: none
//Return: <[]string> An array of string values
//Getter function
//This function will return the array of unexpected tokens found
func GetUnexpectedTokens() []string {
	return unexpected_tokens
}
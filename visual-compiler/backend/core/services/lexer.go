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
// Parameters: <string>
// Return: None
// Setter function which sets the source variable to the value of parameter
func SourceCode(data string) {
	source = data
}

// Name: GetSourceCode
// Parameters: None
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
	Type  string `json:"type"`
	Regex string `json:"regex"`
}

// Array that stores the user's regex rules
var rules []TypeRegex

// Struct for the tokens
type TypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Array that stores the user's generated tokens
var tokens []TypeValue

// Array that stores unidentified input
var tokens_unidentified []string

// Name: ReadRegexRules
// Parameters: <[]byte>
// Return: <error>
// Receive an array of regex rules, validate them and store them in the rules struct
func ReadRegexRules(input []byte) error {

	rules = []TypeRegex{}

	err := json.Unmarshal(input, &rules)

	if err != nil {
		return fmt.Errorf("invalid JSON for rules: %v", err)
	}

	for _, rule := range rules {

		_, err := regexp.Compile(rule.Regex)

		if err != nil {
			return fmt.Errorf("invalid regex input: %v", err)
		}

		rule.Type = strings.ToUpper(rule.Type)
	}

	return nil
}

// Name: CreateTokens
// Parameters: None
// Return: None
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
// Parameters: None
// Return: <[]TypeValue>
// Getter function that returns the array of tokens
func GetTokens() []TypeValue {
	return tokens
}

// Name: GetInvalidInput
// Parameters: None
// Return: <[]string>
// Getter function that returns the array of unidentified input
func GetInvalidInput() []string {
	return tokens_unidentified
}

// ================== //
// REGEX AND AUTOMATA //
// ================== //

// Struct that stores the automata data
var dfa Automata

type Transition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}
type Automata struct {
	States      []string     `json:"states"`
	Transitions []Transition `json:"transitions"`
	Start       string       `json:"start_state"`
	Accepting   []string     `json:"accepting_states"`
}

// Name: ReadDFA
// Parameters: <[]byte>
// Return: <error>
// Receive an array of DFA data, validate them and store them in the struct
func ReadDFA(input []byte) error {

	dfa = Automata{}

	err := json.Unmarshal(input, &dfa)

	if err != nil {
		return fmt.Errorf("invalid JSON for automata: %v", err)
	}

	return nil
}

// Name: ConvertDFAToRegex
// Parameters: None
// Return: None
// Convert the DFA received from the ReadDFA function to a regular expression
func ConvertDFAToRegex() {

}

// Name: ConvertRegexToDFA
// Parameters: <string>
// Return: <Automata>
// Convert the DFA received from the ReadDFA function to a regular expression
func ConvertRegexToDFA(regex string) {

	var states = []string{}
	var transitions = []Transition{}
	var start string
	var accepting = []string{}

	//Temp code to make runnable
	fmt.Printf(states[0])
	fmt.Printf(start)
	fmt.Printf(accepting[0])

	if transitions==nil {
		fmt.Errorf("No transitions")
	}
}

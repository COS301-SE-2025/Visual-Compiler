package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var source string

type TypeRegex struct {
	Type  string
	Regex string
}

var rules []TypeRegex

type TypeValue struct {
	Type  string
	Value string
}

var tokens []TypeValue

func SourceCode(data string) {

	source = data

}

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

func CreateTokens() []TypeValue {

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
		}
	}

	return tokens
}

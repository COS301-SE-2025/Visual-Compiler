package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var source_code string

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

func Initialise(data string) {
	source_code = data
}

func ReadTypeRegex(input []byte) {

	rules = []TypeRegex{}

	var raw_pairs []struct {
		Type  string `json:"type"`
		Regex string `json:"regex"`
	}

	err := json.Unmarshal(input, &raw_pairs)

	if err != nil {
		fmt.Printf("invalid input: %s", err)
	}

	for _, pair := range raw_pairs {

		_, err := regexp.Compile(pair.Regex)

		if err != nil {
			fmt.Printf("invalid regex: %s", err)
		}

		rules = append(rules, TypeRegex{Type: pair.Type, Regex: pair.Regex})
	}
}

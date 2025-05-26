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

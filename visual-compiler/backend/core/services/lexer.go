package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// =========== //
// SOURCE CODE //
// =========== //

// String variable that stores the source code
var source string

// Name: SourceCode
// Parameters: string
// Return: None
// Setter function which sets the source variable to the value of parameter
func SourceCode(data string) {
	source = data
}

// Name: GetSourceCode
// Parameters: None
// Return: string
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
// Parameters: []byte
// Return: error
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
// Return: []TypeValue
// Getter function that returns the array of tokens
func GetTokens() []TypeValue {
	return tokens
}

// Name: GetInvalidInput
// Parameters: None
// Return: []string
// Getter function that returns the array of unidentified input
func GetInvalidInput() []string {
	return tokens_unidentified
}

// ================== //
// REGEX AND AUTOMATA //
// ================== //

type Automata struct {
	States      []string         `json:"states"`
	Transitions []Transition     `json:"transitions"`
	Start       string           `json:"start_state"`
	Accepting   []AcceptingState `json:"accepting_states"`
}

type Transition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}

type AcceptingState struct {
	State string `json:"state"`
	Type  string `json:"token_type"`
}

// Structs that store the automata data
var dfa Automata
var nfa Automata

// Name: GetNFA
// Parameters: None
// Return: Automata
// Getter function to get the stored NFA
func GetNFA() Automata {
	return nfa
}

// Name: GetDFA
// Parameters: None
// Return: Automata
// Getter function to get the stored DFA
func GetDFA() Automata {
	return dfa
}

// Name: ReadDFA
// Parameters: []byte
// Return: error
// Receive an array of DFA data, validate them and store them in the struct
func ReadDFA(input []byte) error {

	dfa = Automata{}

	err := json.Unmarshal(input, &dfa)

	if err != nil {
		return fmt.Errorf("invalid JSON for automata: %v", err)
	}

	return nil
}

type candidate struct {
	value string
	token string
}

type candidate_sol struct {
	state        string
	source_index int
	sol          string
}

// Name: CreateTokensFromDFA
// Parameters: None
// Return: None
// Convert the source code, using DFA received from the ReadDFA function, to a set of tokens
func CreateTokensFromDFA() {
	tokens = []TypeValue{}
	tokens_unidentified = []string{}
	source_pos := 0

	for source_pos < len(source) {

		for source_pos < len(source) && unicode.IsSpace(rune(source[source_pos])) {
			source_pos++
		}
		if source_pos >= len(source) {
			break
		}

		solution_found := false
		best_solution := candidate{
			value: "",
			token: "",
		}

		initial_sol := candidate_sol{state: dfa.Start, source_index: source_pos, sol: ""}
		queue := []candidate_sol{initial_sol}

		for len(queue) > 0 {

			current_state := queue[0] //start searching first element in queue
			queue = queue[1:]         //remove first element from queue

			for _, accepting := range dfa.Accepting {
				if accepting.State == current_state.state {
					solution_found = true
					candidate_solution := candidate{
						value: current_state.sol,
						token: accepting.Type,
					}
					if len(candidate_solution.value) > len(best_solution.value) {
						best_solution = candidate_solution
					}
				}
			}

			if current_state.source_index < len(source) {
				current_char := string(source[current_state.source_index])
				for _, current_transition := range dfa.Transitions { //add children of current_state to queue

					if current_transition.From == current_state.state {
						if strings.Contains(current_transition.Label, current_char) {

							next_state := candidate_sol{
								state:        current_transition.To,
								source_index: current_state.source_index + 1,
								sol:          current_state.sol + current_char,
							}
							queue = append(queue, next_state)
						}
					}
				}
			}
		}

		if solution_found {
			tokens = append(tokens, TypeValue{
				Type:  best_solution.token,
				Value: best_solution.value,
			})
			source_pos += len(best_solution.value)
		} else {
			unexpected_pos := source_pos + 1
			for unexpected_pos < len(source) && !unicode.IsSpace(rune(source[unexpected_pos])) {
				unexpected_pos++
			}
			unidentified_token := source[source_pos:unexpected_pos]
			tokens_unidentified = append(tokens_unidentified, unidentified_token)
			source_pos = unexpected_pos
		}
	}
}

// Name: ConvertDFAToRegex
// Parameters: None
// Return: None
// Convert the DFA received from the ReadDFA function to a regular expression
// Generates a regex for every accepting state and converts the paths to regex rules
func ConvertDFAToRegex() error {

	if len(dfa.States) == 0 {
		return fmt.Errorf("no states identified")
	}
	if len(dfa.Transitions) == 0 {
		return fmt.Errorf("no transitions identified")
	}
	if len(dfa.Accepting) == 0 {
		return fmt.Errorf("no accepting states identified")
	}
	if dfa.Start == "" {
		return fmt.Errorf("no start state identified")
	}

	paths := make(map[string]string)

	for _, accepting := range dfa.Accepting {
		regex := buildRegexForPath(dfa.Start, accepting.State)
		if paths[accepting.Type] != "" {
			paths[accepting.Type] = paths[accepting.Type] + "|" + regex
		} else {
			paths[accepting.Type] = regex
		}
	}

	rules = []TypeRegex{}
	for token_type, regex := range paths {
		convertRawRegexToRegexRules(&regex)
		switch token_type {
		case "KEYWORD":
			convertKeywordToRegex(&regex)
		case "IDENTIFIER":
			convertIdentifierToRegex(&regex)
		case "NUMBER":
			convertNumberToRegex(&regex)
		}
		new_rule := TypeRegex{
			Type:  token_type,
			Regex: regex,
		}
		rules = append(rules, new_rule)
	}

	return nil
}

// Name: convertKeywordToRegex
// Parameters: *string
// Return: None
// Helper function to convert Keyword to proper regex
func convertKeywordToRegex(regex *string) {

	if strings.Contains(*regex, "|") {
		*regex = `\\b(` + *regex + `)\\b`
	} else {
		*regex = `\\b` + *regex + `\\b`
	}

}

// Name: convertIdentifierToRegex
// Parameters: *string
// Return: None
// Helper function to convert Identifiers to proper regex
func convertIdentifierToRegex(regex *string) {

	if matched, _ := regexp.MatchString(`^\[a-z\]\(\[a-z0-9\]\)\*$`, *regex); matched {
		*regex = `[a-zA-Z_]\\w*`
	}

	*regex = strings.ReplaceAll(*regex, "[a-z0-9]", `\\w`)
	*regex = strings.ReplaceAll(*regex, "[a-z]", "[a-zA-Z_]")

}

// Name: convertNumberToRegex
// Parameters: *string
// Return: None
// Helper function to convert Numbers to proper regex
func convertNumberToRegex(regex *string) {
	*regex = `\\d+(\\.\\d+)?`
}

// Name: convertRawRegexToRegexRules
// Parameters: *string
// Return: None
// Helper function to replace grammar with regex rules
func convertRawRegexToRegexRules(regex *string) {
	*regex = strings.ReplaceAll(*regex, "abcdefghijklmnopqrstuvwxyz0123456789", "[a-z0-9]")
	*regex = strings.ReplaceAll(*regex, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "[A-Z0-9]")
	*regex = strings.ReplaceAll(*regex, "abcdefghijklmnopqrstuvwxyz", "[a-z]")
	*regex = strings.ReplaceAll(*regex, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "[A-Z]")
	*regex = strings.ReplaceAll(*regex, "0123456789", `\d`)
}

// Name: buildRegexForPath
// Parameters: string,string
// Return: string
// Helper function to convert DFA to regex, which builds the regex string
// Uses BFS to traverse the states till it reaches accepting saves, for which it saves the current path
func buildRegexForPath(start, accept string) string {
	type state_regex struct {
		state string
		regex string
	}

	initial_candidate := state_regex{state: start, regex: ""}
	queue := []state_regex{initial_candidate}
	visited_states := map[string]bool{start: true}

	path := ""
	found := false

	for len(queue) > 0 && !found {
		current_candidate := queue[0]
		queue = queue[1:]

		for _, transition := range dfa.Transitions {
			if transition.From == current_candidate.state {
				next_state := transition.To
				if !visited_states[next_state] {
					new_regex := current_candidate.regex + regexStructure(transition.Label)
					if next_state == accept {
						path = new_regex
						found = true
						break
					}
					new_candidate := state_regex{
						state: next_state,
						regex: new_regex,
					}
					queue = append(queue, new_candidate)
					visited_states[next_state] = true
				}
			}
		}
	}

	if path == "" {
		return ""
	}

	for _, transition := range dfa.Transitions {
		if transition.From == accept && transition.To == accept {
			multiple_occurrences := "(" + regexStructure(transition.Label) + ")*"
			return path + multiple_occurrences
		}
	}

	return path
}

// Name: regexStructure
// Parameters: string
// Return: string
// Helper function escape special regex characters
func regexStructure(label string) string {

	special_characters := []string{"\\", "(", ")", "{", "}", "[", "]", ".", "^", "$", "*", "+", "?", "|"}

	structured_label := label
	for _, char := range special_characters {
		structured_label = strings.ReplaceAll(structured_label, char, "\\"+char)
	}

	return structured_label
}

type Fragment struct {
	start string
	end   string
}

type Converter struct {
	states      map[string]bool
	transitions []Transition
	count       int
}

// Name: newConverter
// Parameters: None
// Return: None
// Creates a new converter for the regex to automata conversion
func newConverter() *Converter {

	return &Converter{
		states:      make(map[string]bool),
		transitions: make([]Transition, 0),
		count:       0,
	}
}

// Name: newState (for Converter)
// Parameters: None
// Return: string
// Adds a new state to the converter and returns its name
func (c *Converter) newState() string {

	state := "S" + strconv.Itoa(c.count)

	c.count++
	c.states[state] = true

	return state
}

// Name: addTransition (for Converter)
// Parameters: string, string, string
// Return: None
// Adds a new transition to the converter
func (c *Converter) addTransition(from, to, label string) {

	c.transitions = append(c.transitions, Transition{
		From:  from,
		To:    to,
		Label: label,
	})
}

// Name: parseRegex (for Converter)
// Parameters: string
// Return: *Fragment
// Parses a regex string and returns an NFA fragment
func (c *Converter) parseRegex(regex string) *Fragment {

	fragment, _ := c.parseAlter(regex, 0)
	return fragment
}

// Name: parseAlter (for Converter)
// Parameters: string, int
// Return: *Fragment, int
// Handles the alternation regex fragments
func (c *Converter) parseAlter(regex string, position int) (*Fragment, int) {

	left, update := c.parseConcat(regex, position)

	if update >= len(regex) || regex[update] != '|' {
		return left, update
	}

	right, final := c.parseAlter(regex, update+1)

	start := c.newState()
	end := c.newState()

	c.addTransition(start, left.start, "ε")
	c.addTransition(left.end, end, "ε")

	c.addTransition(start, right.start, "ε")
	c.addTransition(right.end, end, "ε")

	return &Fragment{start: start, end: end}, final
}

// Name: parseConcat (for Converter)
// Parameters: string, int
// Return: *Fragment, int
// Handles the concatenation regex fragments
func (c *Converter) parseConcat(regex string, position int) (*Fragment, int) {

	left, update := c.parseStar(regex, position)

	for update < len(regex) && regex[update] != '|' && regex[update] != ')' {

		right, next := c.parseStar(regex, update)

		c.addTransition(left.end, right.start, "ε")
		left = &Fragment{start: left.start, end: right.end}
		update = next
	}

	return left, update
}

// Name: parseStar (for Converter)
// Parameters: string, int
// Return: *Fragment, int
// Handles the zero or more and one or more regex fragments
func (c *Converter) parseStar(regex string, position int) (*Fragment, int) {

	fragment, update := c.parseAtom(regex, position)

	if update < len(regex) {

		switch regex[update] {

		case '*':

			start := c.newState()
			end := c.newState()

			c.addTransition(start, fragment.start, "ε")
			c.addTransition(start, end, "ε")

			c.addTransition(fragment.end, fragment.start, "ε")
			c.addTransition(fragment.end, end, "ε")

			return &Fragment{start: start, end: end}, update + 1

		case '+':

			start := c.newState()
			end := c.newState()

			c.addTransition(start, fragment.start, "ε")
			c.addTransition(fragment.end, fragment.start, "ε")
			c.addTransition(fragment.end, end, "ε")

			return &Fragment{start: start, end: end}, update + 1
		}
	}

	return fragment, update
}

// Name: parseAtom (for Converter)
// Parameters: string, int
// Return: *Fragment, int
// Handles the characters and groups regex fragments
func (c *Converter) parseAtom(regex string, position int) (*Fragment, int) {

	if position >= len(regex) {

		start := c.newState()
		return &Fragment{start: start, end: start}, position
	}

	switch regex[position] {

	case '(':

		level := 1
		i := position + 1

		for i < len(regex) && level > 0 {

			if regex[i] == '(' {
				level++
			} else if regex[i] == ')' {
				level--
			}

			i++
		}

		inner := regex[position+1 : i-1]
		fragment := c.parseRegex(inner)

		return fragment, i

	case '[':

		end := strings.Index(regex[position:], "]")

		if end == -1 {

			start := c.newState()
			end := c.newState()

			c.addTransition(start, end, string(regex[position]))

			return &Fragment{start: start, end: end}, position + 1
		}

		range_str := regex[position+1 : position+end]
		fragment := c.parseRange(range_str)

		return fragment, position + end + 1

	default:

		start := c.newState()
		end := c.newState()

		c.addTransition(start, end, string(regex[position]))

		return &Fragment{start: start, end: end}, position + 1
	}
}

// Name: parseRange (for Converter)
// Parameters: string
// Return: *Fragment
// Handles the ranges like [a-z] regex fragments
func (c *Converter) parseRange(range_str string) *Fragment {

	start := c.newState()
	end := c.newState()

	i := 0

	for i < len(range_str) {

		if i+2 < len(range_str) && range_str[i+1] == '-' {

			for char := range_str[i]; char <= range_str[i+2]; char++ {

				c.addTransition(start, end, string(char))
			}

			i = i + 3

		} else {

			c.addTransition(start, end, string(range_str[i]))
			i++
		}
	}

	return &Fragment{start: start, end: end}
}

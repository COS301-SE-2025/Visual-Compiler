package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Struct for the type and regex pairs
type TypeRegex struct {
	Type  string `json:"type"`
	Regex string `json:"regex"`
}

// Struct for the tokens
type TypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Struct for the nfa and dfa data
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

// Struct for the possible tokens from dfa
type Candidate struct {
	value string
	token string
}

// Struct for the possible solutions considered when converting dfa to regex
type CandidateSol struct {
	state        string
	source_index int
	sol          string
}

// Struct for the fragments of the regex
type Fragment struct {
	start string
	end   string
}

// Struct for the regular expression to automata converter
type Converter struct {
	states      map[string]bool
	transitions []Transition
	count       int
}

// Name: ReadRegexRules
//
// Parameters: []byte
//
// Return: []TypeRegex, error
//
// Receive an array of regex rules, validate them and store them in the rules struct
func ReadRegexRules(input []byte) ([]TypeRegex, error) {

	if len(input) == 0 {
		return nil, fmt.Errorf("no rules specified")
	}

	rules := []TypeRegex{}

	err := json.Unmarshal(input, &rules)

	if err != nil {
		return nil, fmt.Errorf("invalid JSON for rules: %v", err)
	}

	for i, rule := range rules {

		_, err := regexp.Compile(rule.Regex)

		if err != nil {
			return nil, fmt.Errorf("invalid regex input: %v", err)
		}

		rules[i].Type = strings.ToUpper(rule.Type)
	}

	return rules, nil
}

// Name: CreateTokens
//
// Parameters: string, []TypeRegex
//
// Return: []TypeValue, []string, error
//
// Loop through the source code to find all tokens that match the regex rules stored
func CreateTokens(source string, rules []TypeRegex) ([]TypeValue, []string, error) {

	tokens := []TypeValue{}
	tokens_unidentified := []string{}

	if source == "" {
		return nil, nil, fmt.Errorf("source code is empty")
	}

	if len(rules) == 0 {
		return nil, nil, fmt.Errorf("no rules specified")
	}

	var builder strings.Builder

	for i := 0; i < len(source); i++ {

		r := rune(source[i])

		if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) || r == '.' || unicode.IsSpace(r) || r == '-' {

			builder.WriteRune(r)

		} else {

			builder.WriteRune(' ')
			builder.WriteRune(r)

			if i+1 < len(source) {

				i++
				r = rune(source[i])

				if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == ';' || r == '_' || r == '.' || r == '-' || r == ')') {

					builder.WriteRune(r)

				} else {
					i--
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

	return tokens, tokens_unidentified, nil
}

// Name: CreateTokensFromDFA
//
// Parameters: string, Automata
//
// Return: None
//
// Convert the source code, using DFA received from the ReadDFA function, to a set of tokens.
// Returns tokens,identified tokens, and an error if tokenisation not possible.
func CreateTokensFromDFA(source_code string, dfa Automata) ([]TypeValue, []string, error) {

	if source_code == "" {
		return nil, nil, fmt.Errorf("source code is empty")
	}

	if len(dfa.Transitions) == 0 {
		return nil, nil, fmt.Errorf("no transitions identified in dfa")
	}

	if dfa.Start == "" {
		return nil, nil, fmt.Errorf("no start state identified in dfa")
	}

	if len(dfa.Accepting) == 0 {
		return nil, nil, fmt.Errorf("no accepting states identified in dfa")
	}

	tokens := []TypeValue{}
	tokens_unidentified := []string{}
	source_pos := 0

	for source_pos < len(source_code) {

		for source_pos < len(source_code) && unicode.IsSpace(rune(source_code[source_pos])) {
			source_pos++
		}
		if source_pos >= len(source_code) {
			break
		}

		solution_found := false
		best_solution := Candidate{
			value: "",
			token: "",
		}

		initial_sol := CandidateSol{state: dfa.Start, source_index: source_pos, sol: ""}
		queue := []CandidateSol{initial_sol}

		for len(queue) > 0 {

			current_state := queue[0] //start searching first element in queue
			queue = queue[1:]         //remove first element from queue

			for _, accepting := range dfa.Accepting {
				if accepting.State == current_state.state {
					solution_found = true
					candidate_solution := Candidate{
						value: current_state.sol,
						token: accepting.Type,
					}
					if len(candidate_solution.value) > len(best_solution.value) {
						best_solution = candidate_solution
					}
				}
			}

			if current_state.source_index < len(source_code) {
				current_char := string(source_code[current_state.source_index])
				for _, current_transition := range dfa.Transitions { //add children of current_state to queue

					if current_transition.From == current_state.state {

						var match_found bool
						if strings.Contains(current_transition.Label, "[") && strings.Contains(current_transition.Label, "]") {
							match_found, _ = regexp.MatchString("^"+current_transition.Label+"$", current_char)
						} else {
							match_found = strings.Contains(current_transition.Label, current_char)
						}

						if match_found {

							next_state := CandidateSol{
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
			for unexpected_pos < len(source_code) && !unicode.IsSpace(rune(source_code[unexpected_pos])) {
				unexpected_pos++
			}
			unidentified_token := source_code[source_pos:unexpected_pos]
			tokens_unidentified = append(tokens_unidentified, unidentified_token)
			source_pos = unexpected_pos
		}
	}

	return tokens, tokens_unidentified, nil
}

// Name: ConvertDFAToRegex
//
// Parameters: None
//
// Return: None
//
// Convert the DFA received from the ReadDFA function to a regular expression.
// Generates a regex for every accepting state and converts the paths to regex rules.
func ConvertDFAToRegex(dfa Automata) ([]TypeRegex, error) {

	if len(dfa.States) == 0 {
		return nil, fmt.Errorf("no states identified in dfa")
	}

	if len(dfa.Transitions) == 0 {
		return nil, fmt.Errorf("no transitions identified in dfa")
	}

	if len(dfa.Accepting) == 0 {
		return nil, fmt.Errorf("no accepting states identified in dfa")
	}

	if dfa.Start == "" {
		return nil, fmt.Errorf("no start state identified in dfa")
	}

	paths := make(map[string][]string)

	for _, accepting := range dfa.Accepting {
		regex_paths := buildRegexForPath(dfa.Start, accepting.State, dfa)
		if len(regex_paths) > 0 {
			for _, path := range regex_paths {
				paths[accepting.Type] = append(paths[accepting.Type], string(path))
			}
		}
	}

	var rules []TypeRegex
	for token_type, regex := range paths {
		raw := strings.Join(regex, "|")

		ConvertRawRegexToRegexRules(&raw)

		multiple_rules := SimplifyRegex(raw)
		if len(multiple_rules) > 1 {
			for _, rule := range multiple_rules {
				_, err := regexp.Compile(rule)
				if rule != "()+" && err == nil {
					new_rule := TypeRegex{
						Type:  token_type,
						Regex: rule,
					}
					rules = append(rules, new_rule)
				}
			}
		} else {
			_, err := regexp.Compile(raw)
			if raw != "()+" && err == nil {
				new_rule := TypeRegex{
					Type:  token_type,
					Regex: raw,
				}
				rules = append(rules, new_rule)
			}
		}
	}

	return rules, nil
}

// Name: simplifyRegex
//
// Parameters: string
//
// Return: string
//
// Helper function to simplify regex
func SimplifyRegex(regex string) []string {

	if !strings.Contains(regex, "+") {
		var string_array []string
		string_array = append(string_array, regex)
		return string_array
	}

	parts := strings.Split(regex, "|")
	unique_parts := make(map[string]bool)
	no_duplicates := []string{}

	for _, part := range parts {
		_, match_found := unique_parts[part]
		if !match_found {
			unique_parts[part] = true
			no_duplicates = append(no_duplicates, part)
		}
	}

	return no_duplicates
}

// Name: convertRawRegexToRegexRules
//
// Parameters: *string
//
// Return: None
//
// Helper function to replace grammar with regex rules
func ConvertRawRegexToRegexRules(regex *string) {

	var strings_array []strings.Builder
	var build_string strings.Builder
	is_range := false
	end_range := false

	if strings.Contains(*regex, "(") && strings.Contains(*regex, ")") && strings.Contains(*regex, "*") {
		for _, char := range *regex {

			switch char {

			case '(':
				is_range = true
				if build_string.Len() != 0 {
					strings_array = append(strings_array, build_string)
				}
				build_string.Reset()
			case ')':
				end_range = true
			case '*':
				if is_range && end_range {
					*regex = strings.ReplaceAll(*regex, "*", "+")

					for _, curr_string := range strings_array {
						if curr_string.String() == build_string.String() {
							*regex = strings.Replace(*regex, curr_string.String(), "", 1)
						}
					}
				}
				build_string.Reset()
			default:
				build_string.WriteRune(char)

			}

		}

		u_alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		l_alphabet := "abcdefghijklmnopqrstuvwxyz"
		numbers := "0123456789"
		if strings.Contains(*regex, u_alphabet) {
			*regex = strings.ReplaceAll(*regex, u_alphabet, "A-Z")
			*regex = strings.ReplaceAll(*regex, "(", "[")
			*regex = strings.ReplaceAll(*regex, ")", "]")
		}
		if strings.Contains(*regex, l_alphabet) {
			*regex = strings.ReplaceAll(*regex, l_alphabet, "a-z")
			*regex = strings.ReplaceAll(*regex, "(", "[")
			*regex = strings.ReplaceAll(*regex, ")", "]")
		}
		if strings.Contains(*regex, numbers) {
			*regex = strings.ReplaceAll(*regex, numbers, "0-9")
			*regex = strings.ReplaceAll(*regex, "(", "[")
			*regex = strings.ReplaceAll(*regex, ")", "]")
		}

		build_string.Reset()
		for _, char := range *regex {
			switch char {
			case '[':
				if build_string.Len() != 0 {
					if !strings.Contains(build_string.String(), "[") {
						if strings.Contains(build_string.String(), "a-zA-Z0-9") {
							*regex = strings.Replace(*regex, "a-zA-Z0-9", "[a-zA-Z0-9]", 1)
						}
						if strings.Contains(build_string.String(), "a-z0-9") {
							*regex = strings.Replace(*regex, "a-z0-9", "[a-z0-9]", 1)
						}
						if strings.Contains(build_string.String(), "A-Z0-9") {
							*regex = strings.Replace(*regex, "A-Z0-9", "[A-Z0-9]", 1)
						}
						if strings.Contains(build_string.String(), "a-zA-Z") {
							*regex = strings.Replace(*regex, "a-zA-Z", "[a-zA-Z]", 1)
						}
						if strings.Contains(build_string.String(), "a-z") {
							*regex = strings.Replace(*regex, "a-z", "[a-z]", 1)
						}
						if strings.Contains(build_string.String(), "A-Z") {
							*regex = strings.Replace(*regex, "A-Z", "[A-Z]", 1)
						}
						if strings.Contains(build_string.String(), "0-9") {
							*regex = strings.Replace(*regex, "0-9", "[0-9]", 1)
						}

					}
				}
				build_string.Reset()
				build_string.WriteRune(char)
			case ']':
				build_string.WriteRune(char)
			default:
				build_string.WriteRune(char)
			}
		}
		*regex = strings.ReplaceAll(*regex, "([", "[")
		*regex = strings.ReplaceAll(*regex, "])", "]")
		*regex = strings.ReplaceAll(*regex, "[[", "[")
		*regex = strings.ReplaceAll(*regex, "]]", "]")

	}

}

// Name: buildRegexForPath
//
// Parameters: string,string
//
// Return: string
//
// Helper function to convert DFA to regex, which builds the regex string.
// Uses BFS to traverse the states till it reaches accepting saves, for which it saves the current path.
func buildRegexForPath(start, accept string, dfa Automata) []string {

	type state_regex struct {
		state   string
		regex   string
		visited map[string]bool
	}

	initial_candidate := state_regex{
		state:   start,
		regex:   "",
		visited: map[string]bool{start: true},
	}
	queue := []state_regex{initial_candidate}

	var final_path []string

	for len(queue) > 0 {
		current_candidate := queue[0]
		queue = queue[1:]

		for _, transition := range dfa.Transitions {
			if transition.From == current_candidate.state {

				if current_candidate.visited[transition.To] {
					continue
				}

				new_regex := current_candidate.regex + regexStructure(transition.Label)
				new_visited := make(map[string]bool)
				for i, visit := range current_candidate.visited {
					new_visited[i] = visit
				}
				new_visited[transition.To] = true

				if transition.To == accept {
					final_path = append(final_path, new_regex)
				} else {
					new_sol := state_regex{
						state:   transition.To,
						regex:   new_regex,
						visited: new_visited,
					}
					queue = append(queue, new_sol)
				}

			}
		}
	}

	if len(final_path) == 0 {
		return final_path
	}

	for _, transition := range dfa.Transitions {
		if transition.From == accept && transition.To == accept {
			for i, path := range final_path {
				final_path[i] = path + "(" + regexStructure(transition.Label) + ")*"
			}
			break
		}
	}

	return final_path
}

// Name: regexStructure
//
// Parameters: string
//
// Return: string
//
// Helper function escape special regex characters
func regexStructure(label string) string {

	if strings.HasPrefix(label, "[") && strings.HasSuffix(label, "]") {
		return label
	}

	special_characters := []string{"\\", "(", ")", "{", "}", "[", "]", ".", "^", "$", "*", "+", "?", "|"}

	structured_label := label
	for _, char := range special_characters {
		structured_label = strings.ReplaceAll(structured_label, char, "\\"+char)
	}

	return structured_label
}

// Name: ConvertRegexToNFA
//
// Parameters: map[string]string
//
// Return: Automata, error
//
// Converts a set of regular expressions to a single nondeterministic finite automata
func ConvertRegexToNFA(regexes map[string]string) (Automata, error) {

	if len(regexes) == 0 {
		return Automata{}, fmt.Errorf("no regex specified")
	}

	converter := newConverter()

	start_state := converter.newState()
	var accepting_states []AcceptingState

	for token_type, regex := range regexes {

		fragment := converter.parseRegex(regex)

		converter.addTransition(start_state, fragment.start, "ε")

		accepting_states = append(accepting_states, AcceptingState{
			State: fragment.end,
			Type:  token_type,
		})
	}

	states_list := make([]string, 0, len(converter.states))

	for state := range converter.states {
		states_list = append(states_list, state)
	}

	i := 0
	for i < len(converter.transitions)-1 {
		current_transition := &converter.transitions[i]
		next_transition := converter.transitions[i+1]

		if next_transition.To == current_transition.To && next_transition.From == current_transition.From {
			current_transition.Label += next_transition.Label
			converter.transitions = append(converter.transitions[:i+1], converter.transitions[i+2:]...)
		} else {
			i++
		}
	}

	nfa := Automata{}
	nfa.States = states_list
	nfa.Transitions = converter.transitions
	nfa.Start = start_state
	nfa.Accepting = accepting_states

	return nfa, nil
}

// Name: ConvertRegexToDFA
//
// Parameters: map[string]string
//
// Return: Automata, error
//
// Converts a set of regular expressions to a single deterministic finite automata
func ConvertRegexToDFA(regexes map[string]string) (Automata, error) {

	nfa, err := ConvertRegexToNFA(regexes)
	if err != nil {
		return Automata{}, fmt.Errorf("could not convert regex to nfa: %v", err)
	}

	dfa, err := ConvertNFAToDFA(nfa)
	if err != nil {
		return Automata{}, fmt.Errorf("could not convert regex to dfa: %v", err)
	}

	return dfa, nil
}

// Name: ConvertNFAToDFA
//
// Parameters: nfa Automata
//
// Return: Automata, error
//
// Converts the NFA to a DFA
func ConvertNFAToDFA(nfa Automata) (Automata, error) {

	if len(nfa.States) == 0 {
		return Automata{}, fmt.Errorf("no states identified")
	}

	if len(nfa.Transitions) == 0 {
		return Automata{}, fmt.Errorf("no transitions identified")
	}

	if len(nfa.Accepting) == 0 {
		return Automata{}, fmt.Errorf("no accepting states identified")
	}

	if nfa.Start == "" {
		return Automata{}, fmt.Errorf("no start state identified")
	}

	transition_map := make(map[string]map[string][]string)

	for _, t := range nfa.Transitions {

		if transition_map[t.From] == nil {
			transition_map[t.From] = make(map[string][]string)
		}

		transition_map[t.From][t.Label] = append(transition_map[t.From][t.Label], t.To)
	}

	start_closure := closureEpsilon([]string{nfa.Start}, transition_map)
	start_state_key := strings.Join(start_closure, ",")

	dfa_states := make(map[string][]string)
	new_transitions := make([]Transition, 0)
	state_queue := []string{start_state_key}
	processed := make(map[string]bool)

	dfa_states[start_state_key] = start_closure

	count := 0
	state_names := make(map[string]string)
	state_names[start_state_key] = "D" + strconv.Itoa(count)
	count++

	for len(state_queue) > 0 {

		current_state_key := state_queue[0]
		state_queue = state_queue[1:]

		if processed[current_state_key] {
			continue
		}
		processed[current_state_key] = true

		current_states := dfa_states[current_state_key]

		symbols := make(map[string]bool)

		for _, state := range current_states {

			if trans, exists := transition_map[state]; exists {

				for symbol := range trans {

					if symbol != "ε" {
						symbols[symbol] = true
					}
				}
			}
		}

		for symbol := range symbols {

			next_states := make([]string, 0)

			for _, current_state := range current_states {

				if trans, exists := transition_map[current_state]; exists {

					if destinations, exists := trans[symbol]; exists {
						next_states = append(next_states, destinations...)
					}
				}
			}

			if len(next_states) > 0 {

				next_closure := closureEpsilon(next_states, transition_map)
				next_state_key := strings.Join(next_closure, ",")

				if _, exists := dfa_states[next_state_key]; !exists {

					dfa_states[next_state_key] = next_closure
					state_names[next_state_key] = "D" + strconv.Itoa(count)
					count++
					state_queue = append(state_queue, next_state_key)
				}

				new_transitions = append(new_transitions, Transition{
					From:  state_names[current_state_key],
					To:    state_names[next_state_key],
					Label: symbol,
				})
			}
		}
	}

	final_states := make([]string, 0, len(state_names))

	for _, name := range state_names {
		final_states = append(final_states, name)
	}

	accepting_states := make([]AcceptingState, 0)

	for state_key, nfa_states := range dfa_states {

		dfa_state_name := state_names[state_key]

		for _, nfa_state := range nfa_states {

			for _, accepting := range nfa.Accepting {

				if accepting.State == nfa_state {

					accepting_states = append(accepting_states, AcceptingState{
						State: dfa_state_name,
						Type:  accepting.Type,
					})

					break
				}
			}
		}
	}

	dfa := Automata{}
	dfa.States = final_states
	dfa.Transitions = new_transitions
	dfa.Start = state_names[start_state_key]
	dfa.Accepting = accepting_states

	return dfa, nil
}

// Name: newConverter
//
// Parameters: None
//
// Return: None
//
// Creates a new converter for the regex to automata conversion
func newConverter() *Converter {

	return &Converter{
		states:      make(map[string]bool),
		transitions: make([]Transition, 0),
		count:       0,
	}
}

// Name: newState (for Converter)
//
// Parameters: None
//
// Return: string
//
// Adds a new state to the converter and returns its name
func (c *Converter) newState() string {

	state := "S" + strconv.Itoa(c.count)

	c.count++
	c.states[state] = true

	return state
}

// Name: addTransition (for Converter)
//
// Parameters: string, string, string
//
// Return: None
//
// Adds a new transition to the converter
func (c *Converter) addTransition(from, to, label string) {

	c.transitions = append(c.transitions, Transition{
		From:  from,
		To:    to,
		Label: label,
	})
}

// Name: parseRegex (for Converter)
//
// Parameters: string
//
// Return: *Fragment
//
// Parses a regex string and returns an NFA fragment
func (c *Converter) parseRegex(regex string) *Fragment {

	fragment, _ := c.parseAlter(regex, 0)
	return fragment
}

// Name: parseAlter (for Converter)
//
// Parameters: string, int
//
// Return: *Fragment, int
//
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
//
// Parameters: string, int
//
// Return: *Fragment, int
//
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
//
// Parameters: string, int
//
// Return: *Fragment, int
//
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
//
// Parameters: string, int
//
// Return: *Fragment, int
//
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
//
// Parameters: string
//
// Return: *Fragment
//
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

// Name: closureEpsilon
//
// Parameters: []string, map[string]map[string][]string
//
// Return: []string
//
// Removes epsilon transitions from NFA for DFA
func closureEpsilon(states []string, transition_map map[string]map[string][]string) []string {

	closure := make(map[string]bool)

	stack := make([]string, len(states))
	copy(stack, states)

	for _, s := range states {
		closure[s] = true
	}

	for len(stack) > 0 {

		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if epsilon_transition, exists := transition_map[current]["ε"]; exists {

			for _, next := range epsilon_transition {

				if !closure[next] {
					stack = append(stack, next)
					closure[next] = true
				}
			}
		}
	}

	result := make([]string, 0, len(closure))
	for state := range closure {
		result = append(result, state)
	}
	sort.Strings(result)

	return result
}

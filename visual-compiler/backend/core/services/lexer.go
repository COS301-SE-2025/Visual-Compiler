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
type AcceptingState struct {
	State string `json:"state"`
	Type string  `json:"token_type"`
}
type Automata struct {
	States      []string     `json:"states"`
	Transitions []Transition `json:"transitions"`
	Start       string       `json:"start_state"`
	Accepting   []AcceptingState `json:"accepting_states"`
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

type candidate struct {
	value string
	token string
}
type candidate_sol struct {
	state string
	source_index int
	sol string
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

		solution_found :=false
		best_solution := candidate{
									value: "",
									token: "",
								}

		initial_sol := candidate_sol{state:dfa.Start, source_index:source_pos,sol:""}
		queue := []candidate_sol{initial_sol}

		for len(queue) >0 {

			current_state:= queue[0] //start searching first element in queue
			queue = queue[1:] //remove first element from queue

			for _,accepting := range dfa.Accepting {
				if accepting.State == current_state.state {
					solution_found = true
					candidate_solution := candidate{
													value:current_state.sol,
													token: accepting.Type,
												}
					if len(candidate_solution.value) > len(best_solution.value) {
						best_solution = candidate_solution
					}
				}
			}

			if current_state.source_index < len(source) {
				current_char := string(source[current_state.source_index])
				for _,current_transition := range dfa.Transitions { //add children of current_state to queue

					if current_transition.From == current_state.state {
						if strings.Contains(current_transition.Label,current_char) {

							next_state := candidate_sol{
														state: current_transition.To,
														source_index: current_state.source_index+1,
														sol: current_state.sol + current_char,
							}
							queue = append(queue, next_state)
						}
					}
				}
			}
		}

		if solution_found == true {
			tokens = append(tokens,TypeValue{
											Type:best_solution.token,
											Value:best_solution.value,
											})
			source_pos += len(best_solution.value)
		}else {
			unexpected_pos := source_pos +1
			for unexpected_pos < len(source) && !unicode.IsSpace(rune(source[unexpected_pos])) {
				unexpected_pos++
			}
			unidentified_token := source[source_pos:unexpected_pos]
			tokens_unidentified = append(tokens_unidentified, unidentified_token)
			source_pos = unexpected_pos
		}
	}
}


type regex_candidate struct {
	state string
	path string
	visited_states map[string]bool
}

// Name: ConvertDFAToRegex
// Parameters: None
// Return: None
// Convert the DFA received from the ReadDFA function to a regular expression
func ConvertDFAToRegex() error{

	if len(dfa.States) == 0{
		return fmt.Errorf("No states identified.")
	}
	if len(dfa.Transitions) == 0{
		return fmt.Errorf("No transitions identified.")
	}
	if len(dfa.Accepting) == 0{
		return fmt.Errorf("No accepting states identified.")
	}
	if dfa.Start == ""{
		return fmt.Errorf("No start state identified")
	}

	accepting_states := make(map[string]string)
	for _, accepting := range dfa.Accepting {
		accepting_states[accepting.State] = accepting.Type
	}

	rules = []TypeRegex{}
	paths := make(map[string][]string)
	max_depth := len(dfa.States) *2

	initial_candidate := regex_candidate{
		state: dfa.Start,
		path: "",
		visited_states: make(map[string]bool),
	}
	states_queue := []regex_candidate{initial_candidate}

	token_type,is_accepting := accepting_states[dfa.Start]
	if is_accepting {
		paths[token_type] = append(paths[token_type], "")
	}

	for len(states_queue) >0 {
		current_candidate := states_queue[0]
		states_queue = states_queue[1:]

		if len(current_candidate.path) > max_depth {
			continue
		}

		token_type,is_accepting := accepting_states[current_candidate.state]
		if is_accepting && current_candidate.path!="" {
			paths[token_type] = append(paths[token_type], current_candidate.path)
		}

		for _,transition := range dfa.Transitions {
			if transition.From == current_candidate.state {

				if current_candidate.visited_states[transition.To] == false {

					update_visited := make(map[string]bool)
					for k,v :=range current_candidate.visited_states {
						update_visited[k] = v
					}

					update_visited[current_candidate.state] = true

					new_candidate := regex_candidate {
						state: transition.To,
						path: current_candidate.path + regexStructure(transition.Label),
						visited_states: update_visited,
					}

					states_queue = append(states_queue,new_candidate)
				}
			}
		}
	}

	for token_type,token_path := range paths {
		unique_paths_map := make(map[string]bool) 
		unique_paths := []string{}

		for _,path := range token_path {
			if unique_paths_map[path]==false {
				unique_paths_map[path] = true
				unique_paths = append(unique_paths,path)
			}
		}

		regex_builder := ""
		if len(unique_paths) ==1 {
			regex_builder = unique_paths[0]
		}else {
			regex_builder = "("+ strings.Join(unique_paths,"|")+")"
		}

		new_regex := TypeRegex{
			Type: token_type,
			Regex: regex_builder,
		}
		rules = append(rules,new_regex)
		fmt.Printf("%s: %s\n", token_type, regex_builder)
	}

	return nil
}

func regexStructure(label string) string {

	special_characters := []string{"\\","(",")","{","}","[","]",".","^","$","*","+","?","|"}

	structured_label := label
	for _,char := range special_characters {
		structured_label = strings.ReplaceAll(structured_label,char,"\\"+char)
	}

	return structured_label
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

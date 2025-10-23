//go:build !coverage
// +build !coverage

package ai

import (
	"context"
	"fmt"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

var (
	ai_client *openai.Client
	ai_create sync.Once
)

func ConnectAI(api_key string) *openai.Client {
	ai_create.Do(func() {
		ai_client = openai.NewClient(api_key)
		fmt.Println("Successfully connected to OpenAI")
	})

	return ai_client
}

func QuestionAnswer(question string, ctx context.Context) string {

	prompt := "You are an educational assistant. Answer any question relevant to compiler construction clearly, concisely, and at a student level of complexity. Otherwise, politely inform the user that you only explain the topics within the phases of compilation. Structure it like a chat message, therefore no headings, no newlines and no special formatting. If the user seeks system-specific guidance, such as navigation, suggest the canvas tutorial and FAQ tab."

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: question,
		},
	}

	answer, err := ai_client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    "gpt-4o-mini",
			Messages: messages,
		},
	)

	if err != nil || len(answer.Choices) == 0 {
		return "The AI Assistant is currently unavailable. Please try again later!"
	}

	return answer.Choices[0].Message.Content
}

func GenerateInput(phase string, artefact string, ctx context.Context) (string, error) {

	prompt, request := "", ""

	switch phase {

	case "source":

		prompt = `You are an educational assistant for compiler construction.

Generate a small program in an imaginary programming language, without any strings, comments or complex syntax.
Format your response strictly as plain text, without any additional text or markdown or formatting whatsoever, for example,

int blue = 13;

for _i range(12)
{
    print(blue);
}`

		request = " "

	case "lexer":

		prompt = `You are an educational assistant for compiler construction.

Generate a list of token definitions in order of priority for a lexer that will perfectly lex the user's source code.
Include punctuation, but ignore whitespaces except if they are syntactically significant.

The token definitions should follow this structure: 
1. "type": string for token type
2. "regex": string for regular expression in JavaScript style with proper escaping

Format your response strictly as a JSON array, without any additional text whatsoever, for example, 

[
    { "type": "KEYWORD", "regex": "string|integer|boolean" },
    { "type": "IDENTIFIER", "regex": "[a-zA-Z_]+" }
]`

		request = "The source code is...\n\n" + artefact

	case "parser":

		prompt = `You are an educational assistant for compiler construction.

Generate a context-free grammar for a parser that will perfectly parse the user's token stream.
Do not add any empty rules, epsilon rules or recursive rules.

The grammar should follow this structure: 
1. "variables": string of comma-separated variable names
2. "terminals": string of comma-separated terminal names
3. "start": string for starting variable
4. "rules": array of objects including
    4.1. "input": string of single variable for LHS of rule
    4.2. "output": string of space-separated variables and terminals for RHS of rule

Format your response strictly as a JSON object, without any additional text whatsoever, for example,

{
	"variables": "PROGRAM, STATEMENT, FUNCTION",
    "terminals": "KEYWORD, IDENTIFIER, SEPARATOR",
    "start": "PROGRAM",
    "rules": [
        { "input": "PROGRAM", "output": "STATEMENT SEPARATOR" },
	    { "input": "STATEMENT", "output": "KEYWORD IDENTIFIER" }
	         ]
}`

		request = "The token stream is...\n\n" + artefact

	case "analyser":

		prompt = `You are an educational assistant for compiler construction.

Generate the scope and type rules for an analyser that will appropriately analyse the user's syntax tree.
The seven grammar rules are constant, do not add more, only change the token if necessary.

Format your response strictly as a JSON object with this structure, without any additional text whatsoever, for example,

{
    "scope_rules":
    [
        { "start": "{", "end": "}" }
    ],
    "type_rules":
    [
        { "result": "int", "assignment": "=", "lhs": "INTEGER", "operator": ["+"], "rhs": "INTEGER" },
        { "result": "int", "assignment": "=", "lhs": "int", "operator": [], "rhs": "" },
    ],
    "grammar_rules":
    {
        "variable_rule": "VARIABLE",
        "type_rule": "TYPE",
        "function_rule": "FUNCTION",
        "parameter_rule": "PARAMETER",
        "assignment_rule": "ASSIGNMENT",
        "operator_rule": "OPERATOR",
        "term_rule": "TERM"
    }
}`

		request = "The syntax tree is...\n\n" + artefact

	case "translator":

		prompt = `You are an educational assistant for compiler construction.

Generate a list of translation rules for a translator that will perfectly translate the user's syntax tree.
Remember to include tokens for punctation and write out full recursions instead of using shorthands.

The translation rules should follow this structure: 
1. "sequence": string of space-separated token types
2. "translation": array of strings for target code where token placeholders appear inside braces

Format your response strictly as a JSON array, without any additional text whatsoever, for example,

[
    {
		"sequence": "KEYWORD IDENTIFIER ASSIGNMENT INTEGER SEPARATOR",
		"translation": ["add     rax, {INTEGER}", "mov     [{IDENTIFIER}], rax"]
	},
	{
		"sequence": "OPEN_BRACKET INTEGER OPERATOR INTEGER CLOSE_BRACKET",
		"translation": ["mov     rax, {INTEGER}", "add     rax, {INTEGER}"]
	}
]`

		request = "The syntax tree is...\n\n" + artefact

	case "optimiser":

		prompt = `You are an educational assistant for compiler construction.

Generate a simple, complete Go program without comments that contains examples for the following optimisation: 
1. Constant folding for basic operators
2. Dead code elimination for variables, functions, conditions and loops
3. Loop unrolling for for loops in canonical format

Format your response strictly as raw Go code starting with "package main", without any additional text or markdown whatsoever, for example,

package main

import "fmt"

func main() {
	for i := 0; i <= 2; i++ {
		fmt.Println(i)
	}
}`

		request = " "

	default:
		return "", fmt.Errorf("invalid phase keyword")
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: request,
		},
	}

	response, err := ai_client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    "gpt-4o-mini",
			Messages: messages,
		},
	)

	if err != nil || len(response.Choices) == 0 {
		return "", fmt.Errorf("assistant currently unavailable")
	}

	return response.Choices[0].Message.Content, nil
}

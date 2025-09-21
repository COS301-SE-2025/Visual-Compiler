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

func QuestionAnswer(ctx context.Context, question string) string {

	prompt := "You are an educational assistant. Answer any question relevant to compiler construction clearly, concisely, friendly and at a student level of complexity. Otherwise, inform the user that you only explain the topics within the phases of compilation. Provide simple examples where necessary and format using only standard text. If the user seeks system-specific guidance, such as navigation, suggest the canvas tutorial and FAQ tab."

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

	case "lexer":

		prompt = `You are an educational assistant.

Generate a list of token definitions in order of priority for a lexer that will perfectly lex the user's source code.

The token definitions should follow this structure: 
1. "type": string for token type
2. "regex": string for regular expression in JavaScript style with proper escaping

Format your output strictly as a JSON array, without any additional text whatsoever, for example, 

[
    { "type": "KEYWORD", "regex": "string|integer|boolean" },
    { "type": "IDENTIFIER", "regex": "[a-zA-Z_]+" }
]`

		request = "The source code is...\n\n" + artefact

	case "parser":

		prompt = `You are an educational assistant.

Generate a context-free grammar for a parser that will perfectly parse the user's token stream.

The grammar should follow this structure: 
1. "variables": string of comma-separated variable names
2. "terminals": string of comma-separated terminal names
3. "start": string for starting variable
4. "rules": array of objects including
    4.1. "input": string for LHS of rule
    4.2. "output": array of strings for variables or terminals for RHS of rule

Format your output strictly as a JSON object, without any additional text whatsoever, for example,

{
	"variables": "PROGRAM, STATEMENT, FUNCTION",
    "terminals": "KEYWORD, IDENTIFIER, SEPARATOR",
    "start": "PROGRAM",
    "rules": [
        { "input": "PROGRAM", "output": ["STATEMENT", "SEPARATOR"] },
        { "input": "PROGRAM", "output": ["FUNCTION", "SEPARATOR"] }
	    { "input": "STATEMENT", "output": ["KEYWORD", IDENTIFIER"] },
	         ]
}`

		request = "The token stream is...\n\n" + artefact

	case "analyser":

		prompt = `You are an educational assistant.

Generate the scope and type rules for an analyser that will perfectly analyse the user's syntax tree.

Format your output strictly as a JSON object with this structure, without any additional text whatsoever, for example,

{
    "scope_rules":
    [
        { "start": "{", "end": "}" }
    ],
    "type_rules":
    [
        { "result": "int", "assignment": "=", "lhs": "INTEGER", "operator": [], "rhs": "" },
        { "result": "int", "assignment": "=", "lhs": "int", "operator": [], "rhs": "" },
        { "result": "int", "assignment": "=", "lhs": "int", "operator": ["+"], "rhs": "INTEGER" }
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

		prompt = ``

		request = "\n\n" + artefact

	case "optimiser":

		prompt = ``

		request = "\n\n" + artefact

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
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

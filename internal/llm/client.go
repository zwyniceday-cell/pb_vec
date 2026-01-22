package llm

import (
	"context"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/sashabaranov/go-openai"
)

type LLMService struct {
	app    *pocketbase.PocketBase
	client *openai.Client
}

func NewLLMService(app *pocketbase.PocketBase) *LLMService {
	baseURL := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	return &LLMService{
		app:    app,
		client: openai.NewClientWithConfig(config),
	}
}

// EnsureLogCollection ensures that the llm_logs collection exists for tracking token usage.
func (s *LLMService) EnsureLogCollection() error {
	collection, _ := s.app.FindCollectionByNameOrId("llm_logs")
	if collection != nil {
		return nil
	}

	collection = core.NewCollection("", "llm_logs")
	collection.Fields.Add(&core.TextField{Name: "model", Required: true})
	collection.Fields.Add(&core.TextField{Name: "prompt"})
	collection.Fields.Add(&core.TextField{Name: "response"})
	collection.Fields.Add(&core.NumberField{Name: "prompt_tokens"})
	collection.Fields.Add(&core.NumberField{Name: "completion_tokens"})
	collection.Fields.Add(&core.NumberField{Name: "total_tokens"})
	collection.Fields.Add(&core.JSONField{Name: "metadata"})

	return s.app.Save(collection)
}

func (s *LLMService) Chat(ctx context.Context, model string, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionResponse, error) {
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	// Log usage asynchronously
	go func() {
		collection, err := s.app.FindCollectionByNameOrId("llm_logs")
		if err != nil {
			log.Printf("failed to find llm_logs collection: %v", err)
			return
		}

		record := core.NewRecord(collection)
		record.Set("model", resp.Model)
		if len(messages) > 0 {
			record.Set("prompt", messages[len(messages)-1].Content)
		}
		if len(resp.Choices) > 0 {
			record.Set("response", resp.Choices[0].Message.Content)
		}
		record.Set("prompt_tokens", resp.Usage.PromptTokens)
		record.Set("completion_tokens", resp.Usage.CompletionTokens)
		record.Set("total_tokens", resp.Usage.TotalTokens)

		if err := s.app.Save(record); err != nil {
			log.Printf("failed to save llm log: %v", err)
		}
	}()

	return &resp, nil
}

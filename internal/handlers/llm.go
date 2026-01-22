package handlers

import (
	"pocketbase_vec/internal/llm"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/sashabaranov/go-openai"
)

type LLMHandler struct {
	LLMService *llm.LLMService
}

func (h *LLMHandler) Register(se *core.ServeEvent) error {
	RegisterLLMHandlers(se.App.(*pocketbase.PocketBase), se, h.LLMService)
	return nil
}

func RegisterLLMHandlers(app *pocketbase.PocketBase, se *core.ServeEvent, llmService *llm.LLMService) {
	se.Router.POST("/api/llm/chat", func(re *core.RequestEvent) error {
		var body struct {
			Model    string                         `json:"model"`
			Messages []openai.ChatCompletionMessage `json:"messages"`
		}

		if err := re.BindBody(&body); err != nil {
			return re.BadRequestError("Invalid request body", err)
		}

		if body.Model == "" {
			body.Model = "gpt-3.5-turbo" // fallback
		}

		resp, err := llmService.Chat(re.Request.Context(), body.Model, body.Messages)
		if err != nil {
			return re.InternalServerError("LLM request failed", err)
		}

		return re.JSON(200, resp)
	})
}

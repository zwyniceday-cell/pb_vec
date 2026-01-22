package main

import (
	"log"
	"pocketbase_vec/internal/app"
	"pocketbase_vec/internal/handlers"
	"pocketbase_vec/internal/llm"
	"pocketbase_vec/internal/router"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
)

func main() {
	// Initialize the PocketBase app with sqlite-vec
	pbApp := app.NewApp()

	// Initialize LLM service
	llmService := llm.NewLLMService(pbApp)

	// Ensure collections exist on bootstrap
	pbApp.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := llmService.EnsureLogCollection(); err != nil {
			log.Printf("Error ensuring log collection: %v", err)
			return err
		}
		return e.Next()
	})

	// Register JS VM plugin (keep functionality from original)
	jsvm.MustRegister(pbApp, jsvm.Config{})

	// Register all HTTP routes
	router.RegisterRoutes(pbApp,
		&handlers.TestVecHandler{},
		// &handlers.LLMHandler{LLMService: llmService},
	)

	// Start the app
	if err := pbApp.Start(); err != nil {
		log.Fatal(err)
	}
}

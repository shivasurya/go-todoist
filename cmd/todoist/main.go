package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shivasurya/go-todoist/internal/app"
	"github.com/shivasurya/go-todoist/pkg/config"
)

func main() {
	cfg := config.New()
	if cfg.Token == "" {
		log.Fatal("TODOIST_TOKEN environment variable is required")
	}

	todoistApp, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := todoistApp.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

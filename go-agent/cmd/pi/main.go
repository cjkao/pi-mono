package main

import (
	"context"
	"fmt"
	"os"

	"github.com/badlogic/pi-mono/go-agent/pkg/agent"
	"github.com/badlogic/pi-mono/go-agent/pkg/config"
	"github.com/badlogic/pi-mono/go-agent/pkg/llm"
	"github.com/badlogic/pi-mono/go-agent/pkg/tools"
	"github.com/badlogic/pi-mono/go-agent/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func main() {
	var modelName string

	rootCmd := &cobra.Command{
		Use:   "pi",
		Short: "Pi Coding Agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(modelName)
		},
	}

	rootCmd.Flags().StringVarP(&modelName, "model", "m", "", "Model to use (e.g. gpt-4o)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(modelName string) error {
	// 1. Load config
	if modelName == "" {
		defaultModel, err := config.GetDefaultModel()
		if err != nil {
			return fmt.Errorf("failed to get default model: %w", err)
		}
		modelName = defaultModel
	}

	// 2. Setup tools
	registry := tools.NewRegistry()
	registry.Register(&tools.BashTool{})
	registry.Register(&tools.ReadTool{})
	registry.Register(&tools.WriteTool{})
	registry.Register(&tools.EditTool{})

	// 3. Build system prompt
	systemPrompt, err := agent.BuildSystemPrompt(registry)
	if err != nil {
		return fmt.Errorf("failed to build system prompt: %w", err)
	}

	// 4. Setup LLM client
	client, err := llm.NewClient(context.Background(), modelName)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// 5. Create session
	session := agent.NewSession(client, registry, systemPrompt)

	// 6. Start TUI
	p := tea.NewProgram(tui.InitialModel(session), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}

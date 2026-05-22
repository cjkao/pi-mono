package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EditTool struct{}

func (t *EditTool) Name() string {
	return "edit"
}

func (t *EditTool) Description() string {
	return "Edit a file by replacing exact text. The oldText must match exactly (including whitespace)."
}

func (t *EditTool) Parameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to edit",
			},
			"oldText": map[string]any{
				"type":        "string",
				"description": "Exact text to find and replace",
			},
			"newText": map[string]any{
				"type":        "string",
				"description": "New text to replace the old text with",
			},
		},
		"required": []string{"path", "oldText", "newText"},
	}
}

type EditArgs struct {
	Path    string `json:"path"`
	OldText string `json:"oldText"`
	NewText string `json:"newText"`
}

func (t *EditTool) Execute(ctx context.Context, args string) (string, error) {
	var parsedArgs EditArgs
	if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
		return "", err
	}

	info, err := os.Stat(parsedArgs.Path)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	mode := info.Mode()

	contentBytes, err := os.ReadFile(parsedArgs.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	content := string(contentBytes)

	// Normalize line endings for consistent matching if needed?
	// For now, assume exact match.

	if !strings.Contains(content, parsedArgs.OldText) {
		// Try normalizing line endings (simple heuristic)
		contentNorm := strings.ReplaceAll(content, "\r\n", "\n")
		oldTextNorm := strings.ReplaceAll(parsedArgs.OldText, "\r\n", "\n")

		if strings.Contains(contentNorm, oldTextNorm) {
			// Found with normalized line endings
			// We need to be careful about replacing in original content.
			// This is getting tricky without a diff library.
			// Let's stick to simple replacement for now and warn if not found.
			return "", fmt.Errorf("could not find exact text match in %s", parsedArgs.Path)
		}

		return "", fmt.Errorf("could not find exact text match in %s", parsedArgs.Path)
	}

	// Count occurrences
	count := strings.Count(content, parsedArgs.OldText)
	if count > 1 {
		return "", fmt.Errorf("found %d occurrences of text in %s. Please provide more context to make it unique.", count, parsedArgs.Path)
	}

	newContent := strings.Replace(content, parsedArgs.OldText, parsedArgs.NewText, 1)

	if err := os.WriteFile(parsedArgs.Path, []byte(newContent), mode); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully replaced text in %s.", parsedArgs.Path), nil
}

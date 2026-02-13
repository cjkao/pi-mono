package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type WriteTool struct{}

func (t *WriteTool) Name() string {
	return "write"
}

func (t *WriteTool) Description() string {
	return "Create or overwrite a file with new content."
}

func (t *WriteTool) Parameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to write",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "The full content to write to the file",
			},
		},
		"required": []string{"path", "content"},
	}
}

type WriteArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (t *WriteTool) Execute(ctx context.Context, args string) (string, error) {
	var parsedArgs WriteArgs
	if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
		return "", err
	}

	dir := filepath.Dir(parsedArgs.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(parsedArgs.Path, []byte(parsedArgs.Content), 0644); err != nil {
		return "", err
	}

	return "File written successfully.", nil
}

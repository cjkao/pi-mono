package tools

import (
	"context"
	"encoding/json"
	"os"
)

type ReadTool struct{}

func (t *ReadTool) Name() string {
	return "read"
}

func (t *ReadTool) Description() string {
	return "Read file contents."
}

func (t *ReadTool) Parameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to read",
			},
		},
		"required": []string{"path"},
	}
}

type ReadArgs struct {
	Path string `json:"path"`
}

func (t *ReadTool) Execute(ctx context.Context, args string) (string, error) {
	var parsedArgs ReadArgs
	if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
		return "", err
	}

	content, err := os.ReadFile(parsedArgs.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

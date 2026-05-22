package tools

import (
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() any
	Execute(ctx context.Context, args string) (string, error)
}

type ToolRegistry struct {
	tools map[string]Tool
}

func NewRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

func (r *ToolRegistry) Has(name string) bool {
	_, ok := r.tools[name]
	return ok
}

func (r *ToolRegistry) List() []Tool {
	var list []Tool
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

func (r *ToolRegistry) ToOpenAITools() []openai.Tool {
	var tools []openai.Tool
	for _, t := range r.tools {
		params := t.Parameters()

		// Marshal to ensure it's valid JSON schema structure
		// For OpenAI, parameters is `any` so we can pass the struct directly

		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  params,
			},
		})
	}
	return tools
}

// Helper for parsing tool arguments
func ParseArgs[T any](args string) (T, error) {
	var t T
	if err := json.Unmarshal([]byte(args), &t); err != nil {
		return t, err
	}
	return t, nil
}

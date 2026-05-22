package llm

import (
	"context"
	"fmt"

	"github.com/badlogic/pi-mono/go-agent/pkg/config"
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	model  string
}

func NewClient(ctx context.Context, model string) (*Client, error) {
	apiKey, err := config.GetApiKey("openai")
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI API key: %w", err)
	}

	config := openai.DefaultConfig(apiKey)
	client := openai.NewClientWithConfig(config)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

func (c *Client) Chat(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (openai.ChatCompletionResponse, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    tools,
		},
	)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	return resp, nil
}

func (c *Client) StreamChat(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (*openai.ChatCompletionStream, error) {
	stream, err := c.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    tools,
			Stream:   true,
		},
	)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

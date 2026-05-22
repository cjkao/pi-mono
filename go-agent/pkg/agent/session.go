package agent

import (
	"context"
	"fmt"

	"github.com/badlogic/pi-mono/go-agent/pkg/llm"
	"github.com/badlogic/pi-mono/go-agent/pkg/tools"
	"github.com/sashabaranov/go-openai"
)

type Session struct {
	client   *llm.Client
	registry *tools.ToolRegistry
	messages []openai.ChatCompletionMessage
	system   string
}

func NewSession(client *llm.Client, registry *tools.ToolRegistry, systemPrompt string) *Session {
	s := &Session{
		client:   client,
		registry: registry,
		system:   systemPrompt,
		messages: []openai.ChatCompletionMessage{},
	}
	// Initialize with system prompt
	s.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})
	return s
}

func (s *Session) AddMessage(msg openai.ChatCompletionMessage) {
	s.messages = append(s.messages, msg)
}

func (s *Session) Prompt(ctx context.Context, userInput string) (string, error) {
	// Add user message
	s.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	// Loop for tool execution
	for {
		// Call LLM
		// Note: We might want to stream here for TUI feedback, but for now blocking call
		// to get the structure right. TUI will handle streaming via a callback or channel.
		// I'll stick to blocking for the core logic first.

		resp, err := s.client.Chat(ctx, s.messages, s.registry.ToOpenAITools())
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		msg := resp.Choices[0].Message
		s.AddMessage(msg)

		// Check for tool calls
		if len(msg.ToolCalls) > 0 {
			for _, toolCall := range msg.ToolCalls {
				toolName := toolCall.Function.Name
				toolArgs := toolCall.Function.Arguments

				// Find tool
				tool, ok := s.registry.Get(toolName)
				if !ok {
					s.AddMessage(openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    fmt.Sprintf("Tool %s not found", toolName),
						ToolCallID: toolCall.ID,
					})
					continue
				}

				// Execute tool
				// In a real agent, we might want to confirm with user for some tools?
				// But standard behavior is auto-execute unless configured otherwise.
				result, err := tool.Execute(ctx, toolArgs)
				if err != nil {
					result = fmt.Sprintf("Error executing tool: %v", err)
				}

				// Add tool result
				s.AddMessage(openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
			// Loop continues to send tool results back to LLM
		} else {
			// No tool calls, return final response
			return msg.Content, nil
		}
	}
}

// StreamPrompt handles streaming responses and tool execution
// callback is called with partial content
func (s *Session) StreamPrompt(ctx context.Context, userInput string, callback func(string)) (string, error) {
	s.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	for {
		stream, err := s.client.StreamChat(ctx, s.messages, s.registry.ToOpenAITools())
		if err != nil {
			return "", fmt.Errorf("stream error: %w", err)
		}
		defer stream.Close()

		var fullContent string
		var toolCalls []openai.ToolCall

		// Map to accumulate chunks for tool calls (since they can be fragmented)
		// toolCalls index -> partial ToolCall
		toolCallMap := make(map[int]*openai.ToolCall)

		for {
			resp, err := stream.Recv()
			if err != nil {
				// EOF or error
				break
			}

			delta := resp.Choices[0].Delta

			// Handle content
			if delta.Content != "" {
				fullContent += delta.Content
				callback(delta.Content)
			}

			// Handle tool calls
			if len(delta.ToolCalls) > 0 {
				for _, tc := range delta.ToolCalls {
					index := *tc.Index // Assuming index is always present in stream

					if _, ok := toolCallMap[index]; !ok {
						toolCallMap[index] = &openai.ToolCall{
							Index: tc.Index,
							ID:    tc.ID,
							Type:  tc.Type,
							Function: openai.FunctionCall{
								Name:      tc.Function.Name,
								Arguments: "",
							},
						}
					}

					// Append arguments
					toolCallMap[index].Function.Arguments += tc.Function.Arguments
					// Append name if present (usually first chunk)
					toolCallMap[index].Function.Name += tc.Function.Name
				}
			}

			// Handle finish reason? stream.Recv returns EOF on finish usually.
		}

		// Reconstruct tool calls from map
		// Sort by index? Map iteration order is random.
		// Using a slice is better if we track indices.
		// OpenAI stream guarantees index consistency.

		// Actually, let's just use the map to build the toolCalls slice.
		// Max index?
		var maxIndex int = -1
		for i := range toolCallMap {
			if i > maxIndex {
				maxIndex = i
			}
		}

		if maxIndex >= 0 {
			toolCalls = make([]openai.ToolCall, maxIndex+1)
			for i, tc := range toolCallMap {
				toolCalls[i] = *tc
			}
		}

		// Create assistant message
		assistantMsg := openai.ChatCompletionMessage{
			Role:      openai.ChatMessageRoleAssistant,
			Content:   fullContent,
			ToolCalls: toolCalls,
		}
		s.AddMessage(assistantMsg)

		if len(toolCalls) > 0 {
			for _, toolCall := range toolCalls {
				toolName := toolCall.Function.Name
				toolArgs := toolCall.Function.Arguments

				// Notify UI about tool execution?
				callback(fmt.Sprintf("\n[Executing tool %s...]\n", toolName))

				tool, ok := s.registry.Get(toolName)
				var result string
				if !ok {
					result = fmt.Sprintf("Tool %s not found", toolName)
				} else {
					res, err := tool.Execute(ctx, toolArgs)
					if err != nil {
						result = fmt.Sprintf("Error: %v", err)
					} else {
						result = res
					}
				}

				// Notify UI about result?
				// callback(fmt.Sprintf("[Result: %s]\n", result)) // Maybe too verbose

				s.AddMessage(openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
			// Loop continues
		} else {
			return fullContent, nil
		}
	}
}

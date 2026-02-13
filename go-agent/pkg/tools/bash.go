package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

type BashTool struct{}

func (t *BashTool) Name() string {
	return "bash"
}

func (t *BashTool) Description() string {
	return "Execute a bash command in the current working directory. Returns stdout and stderr."
}

func (t *BashTool) Parameters() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The command to execute",
			},
			"timeout": map[string]any{
				"type":        "number",
				"description": "Timeout in seconds (optional)",
			},
		},
		"required": []string{"command"},
	}
}

type BashArgs struct {
	Command string  `json:"command"`
	Timeout float64 `json:"timeout,omitempty"`
}

func (t *BashTool) Execute(ctx context.Context, args string) (string, error) {
	var parsedArgs BashArgs
	if err := json.Unmarshal([]byte(args), &parsedArgs); err != nil {
		return "", err
	}

	timeout := 30 * time.Second
	if parsedArgs.Timeout > 0 {
		timeout = time.Duration(parsedArgs.Timeout * float64(time.Second))
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", parsedArgs.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "Command timed out", nil
		}
		return string(output) + "\nError: " + err.Error(), nil
	}

	return string(output), nil
}

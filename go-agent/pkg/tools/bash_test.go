package tools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestBashTool_Execute(t *testing.T) {
	tool := &BashTool{}

	args := BashArgs{
		Command: "echo 'hello world'",
	}
	argsJSON, _ := json.Marshal(args)

	output, err := tool.Execute(context.Background(), string(argsJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if output != "hello world\n" {
		t.Errorf("Execute() = %q, want %q", output, "hello world\n")
	}
}

func TestBashTool_Timeout(t *testing.T) {
	tool := &BashTool{}

	args := BashArgs{
		Command: "sleep 2; echo 'done'",
		Timeout: 0.1,
	}
	argsJSON, _ := json.Marshal(args)

	output, err := tool.Execute(context.Background(), string(argsJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if output != "Command timed out" {
		t.Errorf("Execute() = %q, want %q", output, "Command timed out")
	}
}

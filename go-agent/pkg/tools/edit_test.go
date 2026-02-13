package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEditTool_Execute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "edit_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tool := &EditTool{}

	args := EditArgs{
		Path:    filePath,
		OldText: "Line 2",
		NewText: "Line Two",
	}
	argsJSON, _ := json.Marshal(args)

	output, err := tool.Execute(context.Background(), string(argsJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	expectedOutput := "Successfully replaced text in " + filePath + "."
	if output != expectedOutput {
		t.Errorf("Execute() = %q, want %q", output, expectedOutput)
	}

	newContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	newContent := string(newContentBytes)

	expectedContent := "Line 1\nLine Two\nLine 3"
	if newContent != expectedContent {
		t.Errorf("File content = %q, want %q", newContent, expectedContent)
	}
}

func TestEditTool_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "edit_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tool := &EditTool{}

	args := EditArgs{
		Path:    filePath,
		OldText: "Line 4", // Does not exist
		NewText: "Line Four",
	}
	argsJSON, _ := json.Marshal(args)

	output, err := tool.Execute(context.Background(), string(argsJSON))
	if err == nil {
		t.Errorf("Execute() expected error, got %q", output)
	}
}

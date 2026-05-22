package agent

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/badlogic/pi-mono/go-agent/pkg/config"
	"github.com/badlogic/pi-mono/go-agent/pkg/tools"
)

func BuildSystemPrompt(registry *tools.ToolRegistry) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	readmePath, _ := config.GetReadmePath()
	docsPath, _ := config.GetDocsPath()
	examplesPath, _ := config.GetExamplesPath()

	toolDescriptions := map[string]string{
		"read":  "Read file contents",
		"bash":  "Execute bash commands (ls, grep, find, etc.)",
		"edit":  "Make surgical edits to files (find exact text and replace)",
		"write": "Create or overwrite files",
	}

	var toolsList []string
	var guidelines []string

	// Build tools list
	for _, tool := range registry.List() {
		desc, ok := toolDescriptions[tool.Name()]
		if !ok {
			desc = tool.Description()
		}
		toolsList = append(toolsList, fmt.Sprintf("- %s: %s", tool.Name(), desc))
	}

	// Guidelines
	hasBash := registry.Has("bash")
	hasEdit := registry.Has("edit")
	hasWrite := registry.Has("write")
	hasRead := registry.Has("read")

	if hasBash {
		guidelines = append(guidelines, "Use bash for file operations like ls, rg, find")
	}
	if hasRead && hasEdit {
		guidelines = append(guidelines, "Use read to examine files before editing. You must use this tool instead of cat or sed.")
	}
	if hasEdit {
		guidelines = append(guidelines, "Use edit for precise changes (old text must match exactly)")
	}
	if hasWrite {
		guidelines = append(guidelines, "Use write only for new files or complete rewrites")
	}
	if hasEdit || hasWrite {
		guidelines = append(guidelines, "When summarizing your actions, output plain text directly - do NOT use cat or bash to display what you did")
	}
	guidelines = append(guidelines, "Be concise in your responses")
	guidelines = append(guidelines, "Show file paths clearly when working with files")

	now := time.Now().Format("Monday, January 2, 2006 at 3:04:05 PM MST")

	prompt := fmt.Sprintf(`You are an expert coding assistant operating inside pi, a coding agent harness. You help users by reading files, executing commands, editing code, and writing new files.

Available tools:
%s

Guidelines:
%s

Pi documentation (read only when the user asks about pi itself, its SDK, extensions, themes, skills, or TUI):
- Main documentation: %s
- Additional docs: %s
- Examples: %s

Current date and time: %s
Current working directory: %s`,
		strings.Join(toolsList, "\n"),
		strings.Join(guidelines, "\n"),
		readmePath,
		docsPath,
		examplesPath,
		now,
		cwd,
	)

	return prompt, nil
}

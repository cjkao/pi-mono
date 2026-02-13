# Pi Coding Agent (Go)

A Go implementation of the Pi Coding Agent.

## Features

- Interactive TUI (Terminal User Interface)
- OpenAI GPT-4o integration
- Tools:
  - Bash execution
  - File reading/writing
  - Surgical file editing (exact match)

## Installation

```bash
go install ./cmd/pi
```

## Usage

```bash
export OPENAI_API_KEY=sk-...
pi
```

## Configuration

- `PI_CODING_AGENT_DIR`: Override the configuration directory (default: `~/.pi/agent`)
- `PI_MODEL`: Override the default model (default: `gpt-4o`)

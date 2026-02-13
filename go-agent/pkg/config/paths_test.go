package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAgentDir(t *testing.T) {
	// Test with env var
	expected := "/tmp/pi-agent-test"
	os.Setenv("PI_CODING_AGENT_DIR", expected)
	defer os.Unsetenv("PI_CODING_AGENT_DIR")

	dir, err := GetAgentDir()
	if err != nil {
		t.Fatalf("GetAgentDir() error = %v", err)
	}
	if dir != expected {
		t.Errorf("GetAgentDir() = %v, want %v", dir, expected)
	}

	// Test without env var (mock home dir?)
	// Hard to mock user home without patching os.UserHomeDir or running in container
	// But env var override is the main logic we added.
}

func TestGetAuthPath(t *testing.T) {
	os.Setenv("PI_CODING_AGENT_DIR", "/tmp/pi")
	defer os.Unsetenv("PI_CODING_AGENT_DIR")

	path, err := GetAuthPath()
	if err != nil {
		t.Fatalf("GetAuthPath() error = %v", err)
	}
	expected := filepath.Join("/tmp/pi", "auth.json")
	if path != expected {
		t.Errorf("GetAuthPath() = %v, want %v", path, expected)
	}
}

package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppName       = "pi"
	ConfigDirName = ".pi"
)

func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func GetAgentDir() (string, error) {
	// Check for environment variable override
	if envDir := os.Getenv("PI_CODING_AGENT_DIR"); envDir != "" {
		return envDir, nil
	}

	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDirName, "agent"), nil
}

func GetAuthPath() (string, error) {
	agentDir, err := GetAgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(agentDir, "auth.json"), nil
}

func GetModelsPath() (string, error) {
	agentDir, err := GetAgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(agentDir, "models.json"), nil
}

func GetSettingsPath() (string, error) {
	agentDir, err := GetAgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(agentDir, "settings.json"), nil
}

func GetToolsDir() (string, error) {
	agentDir, err := GetAgentDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(agentDir, "tools"), nil
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// Package Paths

func GetPackageDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func getPathWithFallback(filename string) (string, error) {
	pkgDir, err := GetPackageDir()
	if err == nil {
		path := filepath.Join(pkgDir, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Fallback to CWD
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(cwd, filename)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Fallback to go-agent subdir if in repo root (common for dev)
	path = filepath.Join(cwd, "go-agent", filename)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// If still not found, return a relative path so it doesn't crash but points to where it should be
	return filepath.Join(cwd, filename), nil
}

func GetReadmePath() (string, error) {
	return getPathWithFallback("README.md")
}

func GetDocsPath() (string, error) {
	return getPathWithFallback("docs")
}

func GetExamplesPath() (string, error) {
	return getPathWithFallback("examples")
}

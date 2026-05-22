package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type AuthConfig map[string]ProviderConfig

type ProviderConfig struct {
	ApiKey string `json:"apiKey,omitempty"`
}

func LoadAuth() (AuthConfig, error) {
	path, err := GetAuthPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(AuthConfig), nil
		}
		return nil, err
	}

	var config AuthConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse auth config: %w", err)
	}

	return config, nil
}

func GetApiKey(provider string) (string, error) {
	// 1. Check environment variable
	envKey := fmt.Sprintf("%s_API_KEY", strings.ToUpper(provider))
	if val := os.Getenv(envKey); val != "" {
		return val, nil
	}

	// Special case for OpenAI if provider is "openai" or "OPENAI"
	if strings.EqualFold(provider, "openai") {
		if val := os.Getenv("OPENAI_API_KEY"); val != "" {
			return val, nil
		}
	}

	// 2. Check auth.json
	config, err := LoadAuth()
	if err != nil {
		return "", err
	}

	// Try exact match
	if p, ok := config[provider]; ok && p.ApiKey != "" {
		return p.ApiKey, nil
	}

	// Try lowercase match
	if p, ok := config[strings.ToLower(provider)]; ok && p.ApiKey != "" {
		return p.ApiKey, nil
	}

	return "", fmt.Errorf("API key not found for provider %s", provider)
}

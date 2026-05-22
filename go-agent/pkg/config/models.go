package config

import (
	"encoding/json"
	"os"
)

type ModelsConfig struct {
	DefaultProvider string `json:"defaultProvider"`
	DefaultModel    string `json:"defaultModel"`
}

func GetDefaultModel() (string, error) {
	// 1. Check environment variable
	if model := os.Getenv("PI_MODEL"); model != "" {
		return model, nil
	}

	// Default fallback
	return "gpt-4o", nil
}

func LoadModels() (ModelsConfig, error) {
	path, err := GetModelsPath()
	if err != nil {
		return ModelsConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ModelsConfig{}, err
	}

	var config ModelsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ModelsConfig{}, err
	}

	return config, nil
}

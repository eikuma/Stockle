package config

import (
	"fmt"
	"os"
)

type AIConfig struct {
	GroqAPIKey      string
	AnthropicAPIKey string
	MaxRetries      int
	TimeoutSeconds  int
	RateLimit       int
}

func NewAIConfig() (*AIConfig, error) {
	groqKey := os.Getenv("GROQ_API_KEY")
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")

	if groqKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable is required")
	}

	if anthropicKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	return &AIConfig{
		GroqAPIKey:      groqKey,
		AnthropicAPIKey: anthropicKey,
		MaxRetries:      3,
		TimeoutSeconds:  30,
		RateLimit:       100,
	}, nil
}
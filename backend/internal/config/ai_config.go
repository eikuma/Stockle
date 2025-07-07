package config

import (
	"fmt"
	"os"
	"time"
)

type AIConfig struct {
	GroqAPIKey       string        `mapstructure:"groq_api_key"`
	AnthropicAPIKey  string        `mapstructure:"anthropic_api_key"`
	RequestTimeout   time.Duration `mapstructure:"request_timeout"`
	MaxRetries       int           `mapstructure:"max_retries"`
	RetryDelay       time.Duration `mapstructure:"retry_delay"`
	RateLimitPerMin  int           `mapstructure:"rate_limit_per_min"`
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
		RequestTimeout:  30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      1 * time.Second,
		RateLimitPerMin: 100,
	}, nil
}
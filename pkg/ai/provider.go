package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI providers
type Provider interface {
	// Generate generates a response from the AI model
	Generate(ctx context.Context, prompt string, options *Options) (*Response, error)
	
	// GenerateStructured generates a structured response (JSON)
	GenerateStructured(ctx context.Context, prompt string, schema interface{}, options *Options) (interface{}, error)
	
	// Name returns the provider name
	Name() string
}

// Options contains configuration for AI generation
type Options struct {
	Temperature   float64
	MaxTokens     int
	Model         string
	SystemPrompt  string
	StopSequences []string
}

// Response represents an AI-generated response
type Response struct {
	Content      string
	Model        string
	TokensUsed   int
	FinishReason string
}

// ProviderType represents the type of AI provider
type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderOllama    ProviderType = "ollama"
	ProviderMock      ProviderType = "mock"
)

// Config holds AI provider configuration
type Config struct {
	Provider   ProviderType
	APIKey     string
	BaseURL    string
	Model      string
	MaxTokens  int
	Temperature float64
}

// NewProvider creates a new AI provider based on configuration
func NewProvider(config *Config) (Provider, error) {
	switch config.Provider {
	case ProviderOpenAI:
		return NewOpenAIProvider(config)
	case ProviderAnthropic:
		return NewAnthropicProvider(config)
	case ProviderOllama:
		return NewOllamaProvider(config)
	case ProviderMock:
		return NewMockProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

// DefaultOptions returns default generation options
func DefaultOptions() *Options {
	return &Options{
		Temperature:   0.7,
		MaxTokens:     2000,
		Model:         "",
		SystemPrompt:  "",
		StopSequences: []string{},
	}
}

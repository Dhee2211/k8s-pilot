package ai

import (
	"context"
	"fmt"
)

// OllamaProvider implements the Provider interface for local Ollama models
type OllamaProvider struct {
	config *Config
	// Add Ollama client here when implementing
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config *Config) (Provider, error) {
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:11434"
	}
	
	if config.Model == "" {
		config.Model = "llama3"
	}
	
	return &OllamaProvider{
		config: config,
	}, nil
}

// Generate generates a response using Ollama
func (o *OllamaProvider) Generate(ctx context.Context, prompt string, options *Options) (*Response, error) {
	// TODO: Implement actual Ollama API call
	// For now, return a placeholder
	//
	// Use HTTP client to call Ollama API:
	// POST http://localhost:11434/api/generate
	// {"model": "llama3", "prompt": "..."}
	
	return nil, fmt.Errorf("Ollama provider not yet implemented - please use mock provider for now")
}

// GenerateStructured generates a structured response using Ollama
func (o *OllamaProvider) GenerateStructured(ctx context.Context, prompt string, schema interface{}, options *Options) (interface{}, error) {
	// TODO: Implement with JSON mode
	return nil, fmt.Errorf("Ollama structured generation not yet implemented")
}

// Name returns the provider name
func (o *OllamaProvider) Name() string {
	return "ollama"
}

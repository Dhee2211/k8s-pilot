package ai

import (
	"context"
	"testing"
)

func TestMockProvider(t *testing.T) {
	config := &Config{
		Provider: ProviderMock,
	}
	
	provider, err := NewMockProvider(config)
	if err != nil {
		t.Fatalf("Failed to create mock provider: %v", err)
	}
	
	if provider.Name() != "mock" {
		t.Errorf("Expected provider name 'mock', got '%s'", provider.Name())
	}
}

func TestMockProviderGenerate(t *testing.T) {
	config := &Config{
		Provider: ProviderMock,
	}
	
	provider, err := NewMockProvider(config)
	if err != nil {
		t.Fatalf("Failed to create mock provider: %v", err)
	}
	
	ctx := context.Background()
	response, err := provider.Generate(ctx, "restart pods", DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to generate response: %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.Content == "" {
		t.Error("Expected non-empty content")
	}
	
	if response.Model != "mock-v1" {
		t.Errorf("Expected model 'mock-v1', got '%s'", response.Model)
	}
}

func TestMockProviderGenerateStructured(t *testing.T) {
	config := &Config{
		Provider: ProviderMock,
	}
	
	provider, err := NewMockProvider(config)
	if err != nil {
		t.Fatalf("Failed to create mock provider: %v", err)
	}
	
	ctx := context.Background()
	result, err := provider.GenerateStructured(ctx, "restart pods", nil, DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to generate structured response: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	
	// Check if result is a map
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected result to be a map")
	}
	
	if _, hasCommands := resultMap["commands"]; !hasCommands {
		t.Error("Expected 'commands' key in result")
	}
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "mock provider",
			config: &Config{
				Provider: ProviderMock,
			},
			expectError: false,
		},
		{
			name: "unsupported provider",
			config: &Config{
				Provider: "unsupported",
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.config)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if provider == nil {
					t.Error("Expected non-nil provider")
				}
			}
		})
	}
}

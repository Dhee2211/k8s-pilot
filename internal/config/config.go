package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	AI       AIConfig       `yaml:"ai"`
	Kube     KubeConfig     `yaml:"kubernetes"`
	Policy   PolicyConfig   `yaml:"policy"`
	Logging  LoggingConfig  `yaml:"logging"`
	Plugins  []string       `yaml:"plugins"`
}

// AIConfig contains AI provider configuration
type AIConfig struct {
	Provider    string  `yaml:"provider"`
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	BaseURL     string  `yaml:"base_url"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
}

// KubeConfig contains Kubernetes configuration
type KubeConfig struct {
	Context   string `yaml:"context"`
	Namespace string `yaml:"namespace"`
}

// PolicyConfig contains policy configuration
type PolicyConfig struct {
	Enabled bool   `yaml:"enabled"`
	OPAAddr string `yaml:"opa_address"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

var globalConfig *Config

// Load loads configuration from a file
func Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	globalConfig = &cfg
	return nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		globalConfig = defaultConfig()
	}
	return globalConfig
}

// defaultConfig returns the default configuration
func defaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider:    "mock",
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		Kube: KubeConfig{
			Namespace: "default",
		},
		Policy: PolicyConfig{
			Enabled: false,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Plugins: []string{},
	}
}

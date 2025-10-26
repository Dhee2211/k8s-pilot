package plugins

// ExamplePlugin is a sample diagnostic plugin
type ExamplePlugin struct {
	initialized bool
}

// NewExamplePlugin creates a new example plugin
func NewExamplePlugin() *ExamplePlugin {
	return &ExamplePlugin{
		initialized: true,
	}
}

// Name returns the plugin name
func (p *ExamplePlugin) Name() string {
	return "example-detector"
}

// Version returns the plugin version
func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *ExamplePlugin) Description() string {
	return "Example diagnostic plugin that detects sample issues"
}

// Initialize initializes the plugin
func (p *ExamplePlugin) Initialize() error {
	p.initialized = true
	return nil
}

// Shutdown cleans up plugin resources
func (p *ExamplePlugin) Shutdown() error {
	p.initialized = false
	return nil
}

// Analyze analyzes a resource and detects issues (implements Plugin interface)
func (p *ExamplePlugin) Analyze(resource interface{}) ([]Issue, error) {
	if !p.initialized {
		return nil, nil
	}

	// Example detection logic
	issues := []Issue{
		{
			ID:          "example-001",
			Severity:    "low",
			Resource:    "example-resource",
			Description: "This is an example issue detected by the plugin",
			Metadata: map[string]interface{}{
				"type": "example",
			},
		},
	}

	return issues, nil
}

// Detect is an alias for Analyze for backward compatibility
func (p *ExamplePlugin) Detect() ([]Issue, error) {
	return p.Analyze(nil)
}

// Remediate provides remediation steps for an issue
func (p *ExamplePlugin) Remediate(issue Issue) ([]RemediationStep, error) {
	steps := []RemediationStep{
		{
			Description: "Check the resource status",
			Command:     "kubectl get <resource>",
			Safe:        true,
		},
		{
			Description: "Describe the resource for more details",
			Command:     "kubectl describe <resource>",
			Safe:        true,
		},
	}

	return steps, nil
}

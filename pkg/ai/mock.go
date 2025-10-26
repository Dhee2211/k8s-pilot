package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// MockProvider is a mock AI provider for testing
type MockProvider struct {
	config *Config
}

// NewMockProvider creates a new mock provider
func NewMockProvider(config *Config) (Provider, error) {
	return &MockProvider{config: config}, nil
}

// Generate generates a mock response
func (m *MockProvider) Generate(ctx context.Context, prompt string, options *Options) (*Response, error) {
	// Generate a simple mock response based on the prompt
	content := m.generateMockResponse(prompt)
	
	return &Response{
		Content:      content,
		Model:        "mock-v1",
		TokensUsed:   len(content) / 4, // rough estimate
		FinishReason: "stop",
	}, nil
}

// GenerateStructured generates a structured mock response
func (m *MockProvider) GenerateStructured(ctx context.Context, prompt string, schema interface{}, options *Options) (interface{}, error) {
	// Generate mock structured data based on prompt
	if strings.Contains(strings.ToLower(prompt), "kubectl") || strings.Contains(strings.ToLower(prompt), "plan") {
		return m.generateMockPlan(prompt), nil
	}
	
	return map[string]interface{}{
		"message": "Mock structured response",
		"status":  "success",
	}, nil
}

// Name returns the provider name
func (m *MockProvider) Name() string {
	return "mock"
}

// generateMockResponse generates a contextual mock response
func (m *MockProvider) generateMockResponse(prompt string) string {
	promptLower := strings.ToLower(prompt)
	
	if strings.Contains(promptLower, "restart") && strings.Contains(promptLower, "pod") {
		return `To restart pods, I recommend using a rolling restart approach:

1. kubectl rollout restart deployment <deployment-name> -n <namespace>

This will gracefully restart all pods in the deployment without downtime.

Alternative: Delete specific pods to trigger recreation:
kubectl delete pod <pod-name> -n <namespace>

The scheduler will automatically create new pods to replace them.`
	}
	
	if strings.Contains(promptLower, "diagnose") || strings.Contains(promptLower, "crashloop") {
		return `Based on the diagnostics:

Issue: CrashLoopBackOff detected
Root Cause: Application is likely failing at startup due to:
  - Missing environment variables
  - Database connection failure
  - Invalid configuration

Recommended Actions:
1. Check logs: kubectl logs <pod-name> -n <namespace>
2. Describe pod: kubectl describe pod <pod-name> -n <namespace>
3. Verify configmaps and secrets are mounted correctly
4. Check resource limits - pod may be OOMKilled`
	}
	
	if strings.Contains(promptLower, "scale") {
		return `To scale the deployment:

kubectl scale deployment <deployment-name> --replicas=<count> -n <namespace>

This will adjust the number of pod replicas to the specified count.
Use 'kubectl get deployment' to verify the scaling operation.`
	}
	
	return fmt.Sprintf("Mock AI response for prompt: %s", prompt)
}

// generateMockPlan generates a mock execution plan
func (m *MockProvider) generateMockPlan(prompt string) map[string]interface{} {
	promptLower := strings.ToLower(prompt)
	
	commands := []map[string]interface{}{}
	
	if strings.Contains(promptLower, "restart") {
		commands = append(commands, map[string]interface{}{
			"command":     "kubectl rollout restart deployment myapp -n default",
			"description": "Restart deployment pods with rolling update",
			"safe":        true,
			"dry_run":     true,
		})
	} else if strings.Contains(promptLower, "scale") {
		commands = append(commands, map[string]interface{}{
			"command":     "kubectl scale deployment myapp --replicas=5 -n default",
			"description": "Scale deployment to 5 replicas",
			"safe":        true,
			"dry_run":     true,
		})
	} else {
		commands = append(commands, map[string]interface{}{
			"command":     "kubectl get pods -n default",
			"description": "List pods in default namespace",
			"safe":        true,
			"dry_run":     false,
		})
	}
	
	plan := map[string]interface{}{
		"summary":  "Generated execution plan",
		"commands": commands,
		"warnings": []string{},
		"requires_confirmation": len(commands) > 0,
	}
	
	b, _ := json.MarshalIndent(plan, "", "  ")
	_ = b // Use b if needed for debugging
	
	return plan
}

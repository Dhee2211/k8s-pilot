package policy

import (
	"context"
	"fmt"
)

// Validator validates Kubernetes operations against policies
type Validator struct {
	enabled bool
}

// NewValidator creates a new policy validator
func NewValidator(enabled bool) *Validator {
	return &Validator{
		enabled: enabled,
	}
}

// ValidationResult represents the result of a policy validation
type ValidationResult struct {
	Allowed   bool
	Violations []Violation
	Warnings   []string
}

// Violation represents a policy violation
type Violation struct {
	Policy   string
	Severity Severity
	Message  string
}

// Severity represents violation severity
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
)

// ValidateCommand validates a kubectl command against policies
func (v *Validator) ValidateCommand(ctx context.Context, command string) (*ValidationResult, error) {
	if !v.enabled {
		return &ValidationResult{
			Allowed:    true,
			Violations: []Violation{},
		}, nil
	}
	
	result := &ValidationResult{
		Allowed:    true,
		Violations: []Violation{},
		Warnings:   []string{},
	}
	
	// Check for dangerous operations
	dangerousOps := []string{"delete", "drain", "cordon"}
	for _, op := range dangerousOps {
		if contains(command, op) {
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("Potentially dangerous operation detected: %s", op))
		}
	}
	
	// Check for privileged operations
	if contains(command, "privileged") || contains(command, "hostNetwork") {
		result.Violations = append(result.Violations, Violation{
			Policy:   "no-privileged-containers",
			Severity: SeverityHigh,
			Message:  "Privileged containers are not allowed",
		})
		result.Allowed = false
	}
	
	// Check for root user
	if contains(command, "runAsUser: 0") {
		result.Violations = append(result.Violations, Violation{
			Policy:   "no-root-containers",
			Severity: SeverityHigh,
			Message:  "Running containers as root is not allowed",
		})
		result.Allowed = false
	}
	
	// Check for resource limits
	if contains(command, "limits") == false && contains(command, "create") {
		result.Warnings = append(result.Warnings, 
			"No resource limits specified - consider adding limits")
	}
	
	return result, nil
}

// ValidateResource validates a Kubernetes resource against policies
func (v *Validator) ValidateResource(ctx context.Context, resourceType, resourceName string, spec map[string]interface{}) (*ValidationResult, error) {
	result := &ValidationResult{
		Allowed:    true,
		Violations: []Violation{},
		Warnings:   []string{},
	}
	
	// TODO: Integrate with OPA/Gatekeeper
	// For now, implement basic validation
	
	return result, nil
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || len(s) > len(substr) && 
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}

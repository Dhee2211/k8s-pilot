# Contributing to kubectl-pilot

Thank you for your interest in contributing to kubectl-pilot! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Plugin Development Guide](#plugin-development-guide)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)

## Code of Conduct

Be respectful, inclusive, and professional in all interactions.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/k8s-pilot.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `go test ./...`
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.22 or later
- Docker (for integration tests)
- kubectl
- Access to a Kubernetes cluster (Kind recommended for local development)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/k8s-pilot.git
cd k8s-pilot

# Install dependencies
go mod download

# Build the binary
go build -o kubectl-pilot .

# Run tests
go test ./...

# Run with verbose logging
./kubectl-pilot run "your command" -v
```

### Setting up a Test Environment

```bash
# Create a Kind cluster
kind create cluster --name pilot-test

# Verify connection
kubectl cluster-info --context kind-pilot-test

# Run kubectl-pilot against the test cluster
./kubectl-pilot diagnose --all-namespaces
```

## Plugin Development Guide

Plugins extend kubectl-pilot's capabilities by adding custom detectors, fixers, and diagnostics.

### Plugin Interface

All plugins must implement the `Plugin` interface:

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Initialize() error
    Shutdown() error
}
```

For diagnostic plugins, implement `DiagnosticPlugin`:

```go
type DiagnosticPlugin interface {
    Plugin
    Detect() ([]Issue, error)
    Remediate(issue Issue) ([]RemediationStep, error)
}
```

### Creating a Plugin

#### Step 1: Create Plugin Structure

```go
package myplugin

import "github.com/yourusername/k8s-pilot/pkg/plugins"

type MyDetector struct {
    initialized bool
    config      map[string]interface{}
}

func New() *MyDetector {
    return &MyDetector{
        config: make(map[string]interface{}),
    }
}
```

#### Step 2: Implement Required Methods

```go
func (p *MyDetector) Name() string {
    return "my-custom-detector"
}

func (p *MyDetector) Version() string {
    return "1.0.0"
}

func (p *MyDetector) Description() string {
    return "Detects custom issues in Kubernetes clusters"
}

func (p *MyDetector) Initialize() error {
    // Setup code here
    p.initialized = true
    return nil
}

func (p *MyDetector) Shutdown() error {
    // Cleanup code here
    p.initialized = false
    return nil
}
```

#### Step 3: Implement Detection Logic

```go
func (p *MyDetector) Detect() ([]plugins.Issue, error) {
    var issues []plugins.Issue
    
    // Your detection logic
    // Example: Check for pods without resource limits
    
    issue := plugins.Issue{
        ID:          "no-resource-limits",
        Severity:    "medium",
        Resource:    "pod/myapp",
        Description: "Pod does not have resource limits set",
        Metadata: map[string]interface{}{
            "recommendation": "Add resource limits and requests",
        },
    }
    
    issues = append(issues, issue)
    return issues, nil
}
```

#### Step 4: Implement Remediation

```go
func (p *MyDetector) Remediate(issue plugins.Issue) ([]plugins.RemediationStep, error) {
    var steps []plugins.RemediationStep
    
    switch issue.ID {
    case "no-resource-limits":
        steps = append(steps, plugins.RemediationStep{
            Description: "Add resource limits to pod specification",
            Command: `kubectl set resources deployment myapp \
              --limits=cpu=500m,memory=512Mi \
              --requests=cpu=250m,memory=256Mi`,
            Safe: true,
        })
    }
    
    return steps, nil
}
```

#### Step 5: Register Your Plugin

```go
func init() {
    manager := plugins.NewManager()
    manager.Register(New())
}
```

### Plugin Example: Memory Leak Detector

```go
package memoryleak

import (
    "context"
    "github.com/yourusername/k8s-pilot/pkg/plugins"
    "github.com/yourusername/k8s-pilot/pkg/k8s"
)

type MemoryLeakDetector struct {
    k8sClient *k8s.Client
}

func New() *MemoryLeakDetector {
    return &MemoryLeakDetector{}
}

func (p *MemoryLeakDetector) Name() string {
    return "memory-leak-detector"
}

func (p *MemoryLeakDetector) Version() string {
    return "1.0.0"
}

func (p *MemoryLeakDetector) Description() string {
    return "Detects pods with continuously increasing memory usage"
}

func (p *MemoryLeakDetector) Initialize() error {
    client, err := k8s.NewClient("")
    if err != nil {
        return err
    }
    p.k8sClient = client
    return nil
}

func (p *MemoryLeakDetector) Shutdown() error {
    return nil
}

func (p *MemoryLeakDetector) Detect() ([]plugins.Issue, error) {
    ctx := context.Background()
    pods, err := p.k8sClient.GetPods(ctx, "")
    if err != nil {
        return nil, err
    }
    
    var issues []plugins.Issue
    
    for _, pod := range pods {
        // Check metrics (simplified)
        if pod.Restarts > 5 {
            issues = append(issues, plugins.Issue{
                ID:          "suspected-memory-leak",
                Severity:    "high",
                Resource:    "pod/" + pod.Name,
                Description: "Pod has restarted multiple times, possibly due to memory issues",
            })
        }
    }
    
    return issues, nil
}

func (p *MemoryLeakDetector) Remediate(issue plugins.Issue) ([]plugins.RemediationStep, error) {
    return []plugins.RemediationStep{
        {
            Description: "Check memory usage over time",
            Command:     "kubectl top pod --all-namespaces",
            Safe:        true,
        },
        {
            Description: "Increase memory limits if needed",
            Command:     "kubectl set resources deployment <name> --limits=memory=1Gi",
            Safe:        false,
        },
    }, nil
}
```

### Plugin Testing

```go
package myplugin

import (
    "testing"
)

func TestDetector(t *testing.T) {
    detector := New()
    
    err := detector.Initialize()
    if err != nil {
        t.Fatalf("Failed to initialize: %v", err)
    }
    
    issues, err := detector.Detect()
    if err != nil {
        t.Fatalf("Detection failed: %v", err)
    }
    
    if len(issues) == 0 {
        t.Log("No issues detected (expected in test environment)")
    }
    
    err = detector.Shutdown()
    if err != nil {
        t.Fatalf("Failed to shutdown: %v", err)
    }
}
```

### Plugin Distribution

1. Create a GitHub repository for your plugin
2. Add a `go.mod` file
3. Tag releases (e.g., `v1.0.0`)
4. Users can install with: `go get github.com/yourname/your-plugin`

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/ai/
```

### Integration Tests

```bash
# Requires a Kind cluster
kind create cluster --name pilot-test

# Run integration tests
go test -tags=integration ./tests/...
```

### Manual Testing

```bash
# Test with mock provider
export K8S_PILOT_AI_PROVIDER=mock
./kubectl-pilot run "restart pods" -v

# Test diagnostics
./kubectl-pilot diagnose --all-namespaces
```

## Submitting Changes

### Pull Request Process

1. Update documentation for any new features
2. Add tests for new functionality
3. Ensure all tests pass: `go test ./...`
4. Update CHANGELOG.md
5. Create a pull request with a clear description

### Commit Message Guidelines

```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

Example:
```
feat(plugins): add memory leak detector plugin

Implements a new plugin that detects pods with continuously
increasing memory usage and provides remediation steps.

Closes #123
```

## Community

- Join our [Discussions](https://github.com/yourusername/k8s-pilot/discussions)
- Report bugs via [Issues](https://github.com/yourusername/k8s-pilot/issues)
- Contribute plugins to our [Plugin Registry](https://github.com/yourusername/k8s-pilot-plugins)

## Questions?

If you have questions about contributing, please:
1. Check the [Wiki](https://github.com/yourusername/k8s-pilot/wiki)
2. Ask in [Discussions](https://github.com/yourusername/k8s-pilot/discussions)
3. Open an issue with the `question` label

Thank you for contributing to kubectl-pilot! 🚁

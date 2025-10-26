# 🚁 kubectl-pilot

[![CI](https://github.com/yourusername/k8s-pilot/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/k8s-pilot/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/k8s-pilot)](https://goreportcard.com/report/github.com/yourusername/k8s-pilot)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> AI-powered Kubernetes operations assistant that translates natural language into safe, auditable Kubernetes actions, diagnoses cluster issues, and provides explainable fixes.

## 🎯 Features

- **Natural Language Operations**: Convert plain English to kubectl commands
- **Cluster Diagnostics**: Auto-detect CrashLoopBackOff, ImagePullBackOff, and other issues
- **Safety First**: Dry-run by default, RBAC-aware, with policy validation
- **AI-Powered Explanations**: Understand logs, events, and cluster behavior
- **Multi-Cloud Support**: Works with GKE, EKS, AKS, K3s, and Kind
- **Extensible Plugin System**: Add custom detectors and fixers
- **Audit Logging**: Track all AI suggestions and executions

## 📦 Installation

### Homebrew

```bash
brew install yourusername/tap/kubectl-pilot
```

### Krew

```bash
kubectl krew install pilot
```

### From Source

```bash
go install github.com/yourusername/k8s-pilot@latest
```

### Docker

```bash
docker pull yourusername/kubectl-pilot:latest
```

## 🚀 Quick Start

### Basic Usage

```bash
# Run a natural language command
kubectl-pilot run "restart failing pods in payments namespace"

# Diagnose cluster issues
kubectl-pilot diagnose --all-namespaces

# Explain what's happening
kubectl-pilot explain logs myapp-pod

# Get help with any resource
kubectl-pilot explain "why is my pod pending"
```

### Configuration

Create `~/.k8s-pilot.yaml`:

```yaml
ai:
  provider: anthropic  # or openai, ollama, mock
  api_key: your-api-key-here
  model: claude-sonnet-4-5-20250929

kubernetes:
  namespace: default

policy:
  enabled: true
```

Or use environment variables:

```bash
export K8S_PILOT_AI_PROVIDER=anthropic
export K8S_PILOT_AI_KEY=your-api-key
```

## 📖 Usage Examples

### Natural Language Commands

```bash
# Scale a deployment
kubectl-pilot run "scale api deployment to 5 replicas" --apply

# Restart pods
kubectl-pilot run "restart all pods in staging namespace" --dry-run

# Check resource usage
kubectl-pilot run "show pods with high memory usage"
```

### Diagnostics

```bash
# Diagnose a specific pod
kubectl-pilot diagnose pod myapp-7d4b9c8f-x9k2l

# Diagnose all pods in a namespace
kubectl-pilot diagnose pods -n production

# Full cluster health check
kubectl-pilot diagnose --all-namespaces
```

### Explanations

```bash
# Explain pod logs
kubectl-pilot explain logs nginx-pod

# Understand cluster events
kubectl-pilot explain events in production

# Learn about resources
kubectl-pilot explain "what is a StatefulSet"
```

### Plugin Management

```bash
# List installed plugins
kubectl-pilot plugin list

# Install a plugin
kubectl-pilot plugin install memory-detector

# Uninstall a plugin
kubectl-pilot plugin uninstall memory-detector
```

## 🏗️ Architecture

```
kubectl-pilot/
├── cmd/pilot/          # CLI commands
├── pkg/
│   ├── ai/            # AI provider abstraction
│   ├── k8s/           # Kubernetes client wrappers
│   ├── diagnose/      # Diagnostics engine
│   ├── plan/          # Natural language → kubectl
│   ├── policy/        # Policy validation
│   ├── explain/       # Explanation engine
│   └── plugins/       # Plugin SDK
├── internal/          # Internal utilities
└── tests/             # Test suites
```

## 🔌 Creating Plugins

```go
package main

import "github.com/yourusername/k8s-pilot/pkg/plugins"

type MyPlugin struct{}

func (p *MyPlugin) Name() string { return "my-detector" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Description() string { return "Custom detector" }

func (p *MyPlugin) Detect() ([]plugins.Issue, error) {
    // Your detection logic here
    return issues, nil
}

func (p *MyPlugin) Remediate(issue plugins.Issue) ([]plugins.RemediationStep, error) {
    // Your remediation logic here
    return steps, nil
}
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full plugin development guide.

## 🛡️ Safety Guarantees

1. **Dry-run by default**: All commands preview changes without applying
2. **Explicit confirmation**: Use `--apply` flag to execute plans
3. **RBAC-aware**: Generates commands respecting user permissions
4. **Policy validation**: Integrates with OPA/Gatekeeper
5. **Audit logging**: Records all AI suggestions and executions
6. **Secret redaction**: Automatically redacts sensitive information

## 🌐 Supported AI Providers

- **OpenAI** (GPT-4, GPT-3.5)
- **Anthropic** (Claude 4 Sonnet, Opus)
- **Ollama** (Local models: Llama 3, Mistral, etc.)
- **Mock** (For testing and development)

## 📊 Diagnostics Coverage

- CrashLoopBackOff
- ImagePullBackOff / ErrImagePull
- Probe failures (liveness, readiness, startup)
- PVC mounting issues
- Resource constraints (CPU, memory)
- Network connectivity problems
- Configuration errors

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

Key areas for contribution:
- Additional AI providers
- More diagnostic detectors
- Plugin ecosystem
- Documentation improvements
- Multi-cloud specific features

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Kubernetes](https://kubernetes.io/) for the amazing orchestration platform
- [client-go](https://github.com/kubernetes/client-go) for the Go client library
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- All contributors to the AI models that power this tool

## 📞 Support

- 📖 [Documentation](https://github.com/yourusername/k8s-pilot/wiki)
- 🐛 [Issue Tracker](https://github.com/yourusername/k8s-pilot/issues)
- 💬 [Discussions](https://github.com/yourusername/k8s-pilot/discussions)

---

**Note**: This tool uses AI to generate Kubernetes commands. Always review the generated plans before applying them to production clusters.

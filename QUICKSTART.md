# 🚁 kubectl-pilot: Quick Start Guide

This guide will help you get started with kubectl-pilot, your AI-powered Kubernetes operations assistant.

## 📁 Project Overview

kubectl-pilot is a production-ready, open-source CLI tool that:
- Translates natural language into safe kubectl commands
- Diagnoses cluster issues automatically
- Provides AI-powered explanations and fixes
- Includes a plugin SDK for extensibility

## 🏗️ Project Structure

```
k8s-pilot/
├── cmd/pilot/              # CLI commands (run, diagnose, explain, plugin)
├── pkg/
│   ├── ai/                # AI provider abstraction (OpenAI, Anthropic, Ollama, Mock)
│   ├── k8s/               # Kubernetes client wrappers
│   ├── diagnose/          # Diagnostics engine
│   ├── plan/              # Natural language → kubectl planner
│   ├── policy/            # Policy validation (OPA/Gatekeeper)
│   ├── explain/           # AI explanation engine
│   └── plugins/           # Plugin SDK with example plugin
├── internal/              # Config, logging, utilities
├── tests/                 # Unit and integration tests
├── examples/              # Sample configs and remediation cookbook
├── .github/workflows/     # CI/CD pipelines
└── Documentation files    # README, CONTRIBUTING, LICENSE
```

## 🚀 Getting Started

### Prerequisites

1. **Go 1.22+**: Install from https://go.dev/dl/
2. **kubectl**: Kubernetes CLI
3. **Access to a Kubernetes cluster** (Kind, Minikube, or cloud)

### Installation Steps

#### Option 1: Build from Source

```bash
cd k8s-pilot
make build
./kubectl-pilot --help
```

#### Option 2: Use Makefile

```bash
# Build
make build

# Build for all platforms
make build-all

# Run tests
make test

# Install locally
make install
```

#### Option 3: Docker

```bash
make docker-build
docker run --rm kubectl-pilot:latest --help
```

### Configuration

Create `~/.k8s-pilot.yaml`:

```yaml
ai:
  provider: mock          # Start with mock, then use anthropic/openai
  api_key: ""            # Add your API key when ready
  model: ""              # Specify model if needed
  
kubernetes:
  namespace: default
  
policy:
  enabled: false
  
logging:
  level: info
```

Or use environment variables:

```bash
export K8S_PILOT_AI_PROVIDER=mock
export K8S_PILOT_AI_KEY=your-api-key  # When using real AI
```

## 📖 Usage Examples

### 1. Natural Language Commands

```bash
# Preview what would happen (dry-run is default)
./kubectl-pilot run "restart failing pods in payments namespace"

# Execute the plan
./kubectl-pilot run "scale api deployment to 5 replicas" --apply

# Get information
./kubectl-pilot run "show pods with high memory usage"
```

### 2. Cluster Diagnostics

```bash
# Diagnose a specific pod
./kubectl-pilot diagnose pod myapp-pod

# Check all pods in a namespace
./kubectl-pilot diagnose pods -n production

# Full cluster health check
./kubectl-pilot diagnose --all-namespaces
```

### 3. AI Explanations

```bash
# Understand logs
./kubectl-pilot explain logs nginx-pod

# Learn about cluster events
./kubectl-pilot explain events in production

# Ask questions
./kubectl-pilot explain "why is my pod pending"
```

### 4. Plugin Management

```bash
# List plugins
./kubectl-pilot plugin list

# The example plugin is included in the project
```

## 🧪 Testing

### Run Tests

```bash
# Unit tests
make test

# With coverage
make coverage

# Integration tests (requires Kind cluster)
make run-kind          # Create test cluster
make test-integration  # Run tests
make delete-kind       # Clean up
```

### Manual Testing

```bash
# Create a test Kind cluster
kind create cluster --name pilot-test

# Test diagnostics
./kubectl-pilot diagnose --all-namespaces -v

# Test natural language
./kubectl-pilot run "list all pods" -v

# Test explanations
./kubectl-pilot explain "what is a deployment"
```

## 🔌 Developing Plugins

See `CONTRIBUTING.md` for the complete plugin development guide.

Quick example:

```go
package myplugin

import "github.com/yourusername/k8s-pilot/pkg/plugins"

type MyPlugin struct{}

func (p *MyPlugin) Name() string { return "my-detector" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Description() string { return "Custom detector" }
func (p *MyPlugin) Initialize() error { return nil }
func (p *MyPlugin) Shutdown() error { return nil }

func (p *MyPlugin) Detect() ([]plugins.Issue, error) {
    // Your detection logic
    return []plugins.Issue{}, nil
}

func (p *MyPlugin) Remediate(issue plugins.Issue) ([]plugins.RemediationStep, error) {
    // Your remediation logic
    return []plugins.RemediationStep{}, nil
}
```

## 🔧 Development Commands

```bash
# Format code
make fmt

# Run linter
make lint

# Build all platforms
make build-all

# Create release
make release-snapshot

# Run with hot reload (requires air)
make dev
```

## 🌐 AI Provider Setup

### Using Anthropic Claude

```yaml
# ~/.k8s-pilot.yaml
ai:
  provider: anthropic
  api_key: sk-ant-...
  model: claude-sonnet-4-5-20250929
```

### Using OpenAI

```yaml
ai:
  provider: openai
  api_key: sk-...
  model: gpt-4
```

### Using Ollama (Local)

```bash
# Install Ollama: https://ollama.ai
ollama pull llama3

# Configure
ai:
  provider: ollama
  base_url: http://localhost:11434
  model: llama3
```

## 📚 Key Files

- **README.md**: Comprehensive project documentation
- **CONTRIBUTING.md**: Plugin SDK guide and contribution guidelines
- **examples/config.yaml**: Sample configuration
- **examples/remediation-cookbook.md**: 20+ common Kubernetes issues and fixes
- **Makefile**: Development commands
- **.github/workflows/**: CI/CD pipelines
- **kubectl-pilot.yaml**: Krew plugin manifest

## 🛡️ Safety Features

1. **Dry-run by default**: All commands preview changes
2. **Explicit --apply flag**: Required to execute plans
3. **RBAC-aware**: Respects user permissions
4. **Audit logging**: Tracks all operations
5. **Secret redaction**: Automatically hides sensitive data

## 🐛 Troubleshooting

### Build Issues

```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version  # Should be 1.22+
```

### Runtime Issues

```bash
# Enable verbose logging
./kubectl-pilot diagnose -v

# Check Kubernetes connection
kubectl cluster-info

# Verify config
cat ~/.k8s-pilot.yaml
```

### AI Provider Issues

```bash
# Test with mock provider first
export K8S_PILOT_AI_PROVIDER=mock
./kubectl-pilot run "test command"

# Check API key
echo $K8S_PILOT_AI_KEY
```

## 📦 Distribution

### Homebrew (Future)

```bash
brew install yourusername/tap/kubectl-pilot
```

### Krew (Future)

```bash
kubectl krew install pilot
```

### From Release

Download pre-built binaries from GitHub Releases.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make test` and `make lint`
6. Submit a pull request

See `CONTRIBUTING.md` for detailed guidelines.

## 📖 Additional Resources

- **Remediation Cookbook**: `examples/remediation-cookbook.md` - 20+ common fixes
- **Example Config**: `examples/config.yaml`
- **Plugin Example**: `pkg/plugins/example_plugin.go`
- **Tests**: `tests/` directory

## 🎯 Next Steps

1. **Build the project**: `make build`
2. **Run tests**: `make test`
3. **Try it out**: `./kubectl-pilot diagnose --all-namespaces`
4. **Configure AI**: Update `~/.k8s-pilot.yaml` with your provider
5. **Read the cookbook**: Check out remediation examples
6. **Develop a plugin**: Follow the CONTRIBUTING.md guide
7. **Contribute**: Submit improvements via PR

## 🆘 Getting Help

- Read the [README.md](README.md)
- Check [CONTRIBUTING.md](CONTRIBUTING.md) for plugin development
- Review [examples/remediation-cookbook.md](examples/remediation-cookbook.md)
- Open an issue on GitHub

## ⭐ Key Features to Explore

1. **Natural Language Processing**: Try various phrasings
2. **Diagnostics Engine**: Test with different pod states
3. **Explain Mode**: Ask questions about your cluster
4. **Plugin System**: Create custom detectors
5. **Policy Validation**: Add OPA/Gatekeeper rules
6. **Multi-cloud Support**: Test across different K8s flavors

---

**Ready to go!** Start with `./kubectl-pilot --help` and explore the capabilities.

For production use, configure a real AI provider (Anthropic, OpenAI, or Ollama) and enjoy AI-powered Kubernetes operations! 🚀

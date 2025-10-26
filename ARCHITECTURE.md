# kubectl-pilot: Architecture & Implementation Guide

## 📐 System Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                        kubectl-pilot CLI                     │
├─────────────────────────────────────────────────────────────┤
│  Commands: run | diagnose | explain | plugin                │
└─────────────────┬───────────────────────────────────────────┘
                  │
    ┌─────────────┴─────────────┐
    │                           │
┌───▼────────┐          ┌──────▼──────┐
│ AI Layer   │          │  K8s Layer  │
├────────────┤          ├─────────────┤
│ • OpenAI   │          │ • client-go │
│ • Anthropic│◄────────►│ • Wrappers  │
│ • Ollama   │          │ • Utilities │
│ • Mock     │          └──────┬──────┘
└────────────┘                  │
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
    ┌────▼────┐          ┌─────▼─────┐         ┌─────▼────┐
    │ Planner │          │ Diagnose  │         │ Explain  │
    │ (NL→k8s)│          │  Engine   │         │  Engine  │
    └────┬────┘          └─────┬─────┘         └─────┬────┘
         │                     │                      │
         └─────────────────────┼──────────────────────┘
                               │
                      ┌────────▼────────┐
                      │ Policy Validator │
                      │ Plugin Manager   │
                      │ Audit Logger     │
                      └──────────────────┘
```

## 🏛️ Component Architecture

### 1. CLI Layer (`cmd/pilot/`)

**Purpose**: User interface and command routing

**Components**:

- `root.go`: Base command with global flags
- `run.go`: Natural language execution
- `diagnose.go`: Cluster diagnostics
- `explain.go`: AI explanations
- `plugin.go`: Plugin management

**Key Features**:

- Cobra-based command structure
- Global flags (dry-run, verbose, namespace)
- Configuration loading
- Error handling and user feedback

### 2. AI Abstraction Layer (`pkg/ai/`)

**Purpose**: Unified interface for multiple AI providers

**Design Pattern**: Strategy Pattern

**Interface**:

```go
type Provider interface {
    Generate(ctx, prompt, options) (*Response, error)
    GenerateStructured(ctx, prompt, schema, options) (interface{}, error)
    Name() string
}
```

**Implementations**:

1. **MockProvider**: For testing and development

   - Returns contextual responses
   - No external dependencies
   - Fast and deterministic

2. **OpenAIProvider**: GPT-4 integration

   - Stub implementation ready for completion
   - Supports function calling
   - JSON mode for structured output

3. **AnthropicProvider**: Claude integration

   - Stub implementation ready for completion
   - High-quality reasoning
   - Large context windows

4. **OllamaProvider**: Local models
   - Privacy-focused
   - No API costs
   - Runs on local hardware

**Configuration**:

```yaml
ai:
  provider: anthropic
  api_key: sk-ant-...
  model: claude-sonnet-4-5-20250929
  temperature: 0.7
  max_tokens: 2000
```

### 3. Kubernetes Client Layer (`pkg/k8s/`)

**Purpose**: Abstract and simplify Kubernetes operations

**Components**:

- `client.go`: Main client wrapper
- `pods.go`: Pod operations and helpers

**Key Methods**:

```go
// Client creation
NewClient(namespace) (*Client, error)

// Resource operations
GetPods(ctx, namespace) ([]PodInfo, error)
GetPod(ctx, name, namespace) (*Pod, error)
GetPodLogs(ctx, podName, container, namespace, tailLines) (string, error)
GetEvents(ctx, namespace) (*EventList, error)
DeletePod(ctx, name, namespace) error

// Cluster detection
DetectClusterType(ctx) (ClusterType, error)
```

**Features**:

- Automatic kubeconfig detection
- In-cluster and out-of-cluster support
- Namespace management
- Multi-cloud awareness

### 4. Plan Generator (`pkg/plan/`)

**Purpose**: Convert natural language to kubectl commands

**Flow**:

```
Natural Language Query
        │
        ▼
   AI Prompt Construction
        │
        ▼
   AI Provider (Generate)
        │
        ▼
   Response Parsing
        │
        ▼
   Plan Structure
        │
        ▼
   Display/Execute
```

**Plan Structure**:

```go
type Plan struct {
    Summary      string      // High-level description
    Commands     []Command   // kubectl commands to run
    Warnings     []string    // Safety warnings
    RequiresAuth bool        // Needs elevated permissions
    DryRun       bool        // Preview mode
}

type Command struct {
    Command     string  // kubectl command
    Description string  // What it does
    Safe        bool    // Safety indicator
    DryRun      bool    // Preview flag
}
```

**Safety Features**:

- Dry-run by default
- Explicit --apply required
- Warning detection
- RBAC-aware generation

### 5. Diagnostics Engine (`pkg/diagnose/`)

**Purpose**: Detect and analyze cluster issues

**Issue Types Detected**:

1. CrashLoopBackOff
2. ImagePullBackOff
3. Pod Eviction
4. Probe Failures (liveness, readiness)
5. PVC Issues
6. Resource Constraints

**Detection Flow**:

```
Cluster/Resource
        │
        ▼
   Data Collection
   (pods, events, logs)
        │
        ▼
   Issue Detection
   (pattern matching)
        │
        ▼
   AI Analysis
   (root cause)
        │
        ▼
   Remediation Generation
        │
        ▼
   Report with Fixes
```

**Report Structure**:

```go
type Report struct {
    Summary       string         // Overall status
    Issues        []Issue        // Detected problems
    Remediations  []Remediation  // Suggested fixes
    HealthScore   int            // 0-100 cluster health
}

type Issue struct {
    Severity    Severity  // critical, high, medium, low
    Type        IssueType // specific problem type
    Resource    string    // affected resource
    Description string    // detailed description
}

type Remediation struct {
    Title       string  // Fix name
    Description string  // What it does
    Command     string  // kubectl command
    Confidence  string  // How sure we are
    Safe        bool    // Safe to execute
}
```

### 6. Explain Engine (`pkg/explain/`)

**Purpose**: Provide AI-powered explanations

**Capabilities**:

1. **Log Analysis**: Summarize and explain pod logs
2. **Event Interpretation**: Understand cluster events
3. **Resource Explanation**: Teach K8s concepts
4. **Troubleshooting**: Answer "why" questions

**Explanation Types**:

- Logs: AI analyzes recent logs for errors/patterns
- Events: Summarizes cluster activity
- Resources: Educational content about K8s objects
- Concepts: General Kubernetes knowledge

### 7. Policy Validation (`pkg/policy/`)

**Purpose**: Validate operations against policies

**Integration Points**:

- OPA (Open Policy Agent)
- Gatekeeper
- Built-in safety rules

**Validation Checks**:

1. Privileged containers
2. Root users
3. Resource limits
4. Network policies
5. Security contexts

**Result Structure**:

```go
type ValidationResult struct {
    Allowed    bool         // Can proceed
    Violations []Violation  // Policy violations
    Warnings   []string     // Non-blocking warnings
}
```

### 8. Plugin System (`pkg/plugins/`)

**Purpose**: Extensibility through plugins

**Plugin Interface**:

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Initialize() error
    Shutdown() error
}

type DiagnosticPlugin interface {
    Plugin
    Detect() ([]Issue, error)
    Remediate(issue Issue) ([]RemediationStep, error)
}
```

**Plugin Lifecycle**:

1. Registration
2. Initialization
3. Active (detection/remediation)
4. Shutdown
5. Unregistration

**Example Plugin**: Memory leak detector included

## 🔄 Data Flow Examples

### Natural Language Execution

```
User: "restart failing pods in payments namespace"
  │
  ▼
CLI (run command)
  │
  ▼
Planner.Generate()
  │
  ├─> Build AI prompt with context
  │
  ├─> AI Provider.Generate()
  │     │
  │     └─> Returns: kubectl commands + explanations
  │
  ├─> Parse response into Plan structure
  │
  └─> Display plan to user
        │
        ▼
    If --apply:
        │
        ▼
    Plan.Execute()
        │
        └─> Run kubectl commands
            │
            └─> Return results
```

### Cluster Diagnostics

```
User: "diagnose pod myapp"
  │
  ▼
CLI (diagnose command)
  │
  ▼
Engine.DiagnoseResource()
  │
  ├─> K8s Client: Get pod details
  │
  ├─> K8s Client: Get pod events
  │
  ├─> K8s Client: Get pod logs
  │
  ├─> Analyze pod state
  │     │
  │     ├─> Check phase (Running/Failed)
  │     ├─> Check container statuses
  │     ├─> Count restarts
  │     └─> Identify issues
  │
  ├─> Generate remediation steps
  │     │
  │     └─> AI Provider: Suggest fixes
  │
  └─> Build diagnostic report
        │
        └─> Display to user
```

## 🧩 Design Patterns Used

1. **Strategy Pattern**: AI provider abstraction
2. **Factory Pattern**: Provider creation
3. **Builder Pattern**: Plan construction
4. **Template Method**: Diagnostic flow
5. **Observer Pattern**: Plugin events
6. **Singleton Pattern**: Global config

## 🔒 Security Considerations

### Built-in Security

1. **Secret Redaction**: Automatic in logs and output
2. **RBAC Awareness**: Respects user permissions
3. **Audit Logging**: All operations tracked
4. **Dry-run Default**: No accidental changes
5. **Policy Validation**: Pre-execution checks

### Security Best Practices

```go
// Secret redaction
func RedactSecrets(input string) string {
    patterns := []string{"password", "token", "secret", "key"}
    // Redact sensitive values
}

// Audit logging
logger.Audit(action, user, resource, success)
```

## 📊 Error Handling Strategy

### Error Types

1. **User Errors**: Invalid input, missing config
2. **K8s Errors**: API failures, network issues
3. **AI Errors**: Provider failures, rate limits
4. **System Errors**: File I/O, parsing failures

### Error Handling Pattern

```go
if err != nil {
    logger.Error("Operation failed: %v", err)
    return fmt.Errorf("failed to X: %w", err)  // Wrap errors
}
```

## 🎯 Performance Considerations

### Optimization Strategies

1. **Caching**: K8s API responses
2. **Concurrency**: Parallel resource queries
3. **Rate Limiting**: AI provider calls
4. **Pagination**: Large result sets
5. **Lazy Loading**: Plugin initialization

### Resource Usage

- Memory: ~50MB base + AI provider overhead
- CPU: Minimal (AI calls are async)
- Network: K8s API + AI provider
- Disk: Config files + audit logs

## 🧪 Testing Strategy

### Test Levels

1. **Unit Tests**: Individual components

   - `tests/ai_test.go`: AI provider tests
   - Mock dependencies
   - Fast execution

2. **Integration Tests**: Component interaction

   - Requires Kind cluster
   - Real K8s operations
   - Tag: `integration`

3. **E2E Tests**: Full workflows
   - Complete user scenarios
   - Production-like environment

### Test Coverage Goals

- Unit tests: >80%
- Integration tests: Core flows
- E2E tests: Happy paths

## 📈 Extensibility Points

### Adding New Features

1. **New AI Provider**:

   - Implement `Provider` interface
   - Add to factory in `provider.go`
   - Update config schema

2. **New Command**:

   - Create file in `cmd/pilot/`
   - Implement Cobra command
   - Register in `init()`

3. **New Diagnostic**:

   - Add detector to `diagnose/engine.go`
   - Define issue type
   - Create remediation logic

4. **New Plugin**:
   - Implement plugin interfaces
   - Create in separate package
   - Register with manager

## 🚀 Deployment Options

### Standalone Binary

```bash
./kubectl-pilot run "command"
```

### Kubectl Plugin

```bash
kubectl pilot run "command"
```

### Docker Container

```bash
docker run kubectl-pilot run "command"
```

### In-Cluster Deployment

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: pilot-diagnostics
spec:
  schedule: "*/15 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: pilot
              image: kubectl-pilot:latest
              command: ["./kubectl-pilot", "diagnose", "--all-namespaces"]
```

## 🔮 Future Enhancements

### Planned Features

1. **Advanced AI Integration**:

   - Streaming responses
   - Multi-turn conversations
   - Context persistence

2. **Enhanced Diagnostics**:

   - Historical analysis
   - Predictive alerts
   - Anomaly detection

3. **Extended Platform Support**:

   - GitOps integration
   - Helm chart diagnostics
   - Operator troubleshooting

4. **Community Features**:
   - Plugin marketplace
   - Shared remediation recipes
   - Community detectors

### Extensibility Roadmap

- Phase 1: Core functionality (✅ Complete)
- Phase 2: Real AI providers (🔄 Ready for implementation)
- Phase 3: Advanced diagnostics (📋 Planned)
- Phase 4: Community ecosystem (🎯 Future)

---

## 📚 Code Organization Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Injection**: Components receive dependencies
3. **Interface-Based**: Abstract implementations
4. **Testability**: All components can be tested in isolation
5. **Documentation**: Comprehensive comments and examples

## 🎓 Learning Resources

To understand the codebase:

1. Start with `cmd/pilot/root.go` - entry point
2. Review `pkg/ai/provider.go` - AI abstraction
3. Study `pkg/plan/planner.go` - core logic
4. Explore `pkg/diagnose/engine.go` - diagnostics
5. Read `CONTRIBUTING.md` - plugin development

---

**This architecture provides a solid foundation for an AI-powered Kubernetes operations tool that is safe, extensible, and production-ready.**

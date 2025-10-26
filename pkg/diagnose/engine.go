package diagnose

import (
	"context"
	"fmt"
	"strings"

	"k8s-pilot/pkg/ai"
	"k8s-pilot/pkg/k8s"
)

// Engine performs diagnostics on Kubernetes resources
type Engine struct {
	k8sClient     *k8s.Client
	aiProvider    ai.Provider
	namespace     string
	allNamespaces bool
}

// NewEngine creates a new diagnostics engine
func NewEngine(namespace string, allNamespaces bool) *Engine {
	k8sClient, err := k8s.NewClient(namespace)
	if err != nil {
		// Handle error appropriately
		panic(err)
	}
	
	// Initialize AI provider
	aiConfig := &ai.Config{
		Provider: ai.ProviderMock,
	}
	aiProvider, _ := ai.NewProvider(aiConfig)
	
	return &Engine{
		k8sClient:     k8sClient,
		aiProvider:    aiProvider,
		namespace:     namespace,
		allNamespaces: allNamespaces,
	}
}

// Report represents a diagnostic report
type Report struct {
	Summary       string
	Issues        []Issue
	Remediations  []Remediation
	HealthScore   int
}

// Issue represents a detected issue
type Issue struct {
	Severity    Severity
	Type        IssueType
	Resource    string
	Description string
	Details     map[string]interface{}
}

// Remediation represents a suggested fix
type Remediation struct {
	Title       string
	Description string
	Command     string
	Confidence  string
	Safe        bool
}

// Severity represents issue severity
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// IssueType represents the type of issue
type IssueType string

const (
	IssueCrashLoopBackOff  IssueType = "CrashLoopBackOff"
	IssueImagePullBackOff  IssueType = "ImagePullBackOff"
	IssuePodEvicted        IssueType = "PodEvicted"
	IssueProbeFailure      IssueType = "ProbeFailure"
	IssuePVCPending        IssueType = "PVCPending"
	IssueResourceConstraint IssueType = "ResourceConstraint"
)

// DiagnoseResource diagnoses a specific resource
func (e *Engine) DiagnoseResource(resourceType, resourceName string) (*Report, error) {
    ctx := context.Background()

    switch strings.ToLower(resourceType) {
    case "pod", "pods":
        return e.diagnosePod(ctx, resourceName)

    case "deployment", "deployments":
        return e.diagnoseDeployment(ctx, resourceName)

    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
}

// DiagnoseResourceType diagnoses all resources of a type
func (e *Engine) DiagnoseResourceType(resourceType string) (*Report, error) {
	ctx := context.Background()
	
	switch strings.ToLower(resourceType) {
	case "pod", "pods":
		return e.diagnoseAllPods(ctx)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// DiagnoseCluster diagnoses the entire cluster/namespace
func (e *Engine) DiagnoseCluster() (*Report, error) {
	ctx := context.Background()
	return e.diagnoseAllPods(ctx)
}

// diagnosePod diagnoses a specific pod
func (e *Engine) diagnosePod(ctx context.Context, podName string) (*Report, error) {
	pod, err := e.k8sClient.GetPod(ctx, podName, e.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}
	
	report := &Report{
		Summary:      fmt.Sprintf("Diagnostics for pod: %s", podName),
		Issues:       []Issue{},
		Remediations: []Remediation{},
		HealthScore:  100,
	}
	
	// Check pod phase
	if pod.Status.Phase != "Running" {
		issue := Issue{
			Severity:    SeverityHigh,
			Type:        IssueType(pod.Status.Phase),
			Resource:    podName,
			Description: fmt.Sprintf("Pod is in %s phase", pod.Status.Phase),
		}
		report.Issues = append(report.Issues, issue)
		report.HealthScore -= 30
	}
	
	// Check container statuses
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			reason := cs.State.Waiting.Reason
			
			var issueType IssueType
			var severity Severity
			
			switch reason {
			case "CrashLoopBackOff":
				issueType = IssueCrashLoopBackOff
				severity = SeverityCritical
				report.HealthScore -= 40
			case "ImagePullBackOff", "ErrImagePull":
				issueType = IssueImagePullBackOff
				severity = SeverityHigh
				report.HealthScore -= 35
			default:
				issueType = IssueType(reason)
				severity = SeverityMedium
				report.HealthScore -= 20
			}
			
			issue := Issue{
				Severity:    severity,
				Type:        issueType,
				Resource:    fmt.Sprintf("%s/%s", podName, cs.Name),
				Description: fmt.Sprintf("Container %s: %s - %s", cs.Name, reason, cs.State.Waiting.Message),
			}
			report.Issues = append(report.Issues, issue)
		}
		
		// Check restart count
		if cs.RestartCount > 5 {
			issue := Issue{
				Severity:    SeverityMedium,
				Type:        IssueCrashLoopBackOff,
				Resource:    fmt.Sprintf("%s/%s", podName, cs.Name),
				Description: fmt.Sprintf("Container has restarted %d times", cs.RestartCount),
			}
			report.Issues = append(report.Issues, issue)
			report.HealthScore -= 15
		}
	}
	
	// Generate remediations using AI
	if len(report.Issues) > 0 {
		remediations := e.generateRemediations(ctx, report.Issues, podName)
		report.Remediations = remediations
	}
	
	return report, nil
}

// diagnoseAllPods diagnoses all pods in the namespace
func (e *Engine) diagnoseAllPods(ctx context.Context) (*Report, error) {
	pods, err := e.k8sClient.GetPods(ctx, e.namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	
	report := &Report{
		Summary:      "Cluster diagnostics",
		Issues:       []Issue{},
		Remediations: []Remediation{},
		HealthScore:  100,
	}
	
	problemPods := 0
	
	for _, pod := range pods {
		if pod.Phase != "Running" || !pod.Ready || pod.Restarts > 3 {
			problemPods++
			
			issue := Issue{
				Severity:    SeverityMedium,
				Resource:    fmt.Sprintf("pod/%s", pod.Name),
				Description: fmt.Sprintf("Pod %s: Phase=%s, Ready=%v, Restarts=%d", 
					pod.Name, pod.Phase, pod.Ready, pod.Restarts),
			}
			
			if pod.Restarts > 10 {
				issue.Severity = SeverityHigh
			}
			
			report.Issues = append(report.Issues, issue)
		}
	}
	
	if problemPods > 0 {
		report.HealthScore = 100 - (problemPods * 10)
		if report.HealthScore < 0 {
			report.HealthScore = 0
		}
	}
	
	report.Summary = fmt.Sprintf("Found %d issue(s) across %d pods. Health score: %d/100", 
		len(report.Issues), len(pods), report.HealthScore)
	
	return report, nil
}

// diagnoseDeployment diagnoses a deployment
func (e *Engine) diagnoseDeployment(ctx context.Context, deploymentName string) (*Report, error) {
	// TODO: Implement deployment diagnostics
	return &Report{
		Summary: fmt.Sprintf("Diagnostics for deployment: %s (not yet implemented)", deploymentName),
	}, nil
}

// generateRemediations uses AI to generate remediation suggestions
func (e *Engine) generateRemediations(ctx context.Context, issues []Issue, resourceName string) []Remediation {
	// Build a prompt describing the issues
	prompt := fmt.Sprintf(`Kubernetes diagnostics for resource: %s

Detected issues:
`, resourceName)
	
	for _, issue := range issues {
		prompt += fmt.Sprintf("- [%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
	}
	
	prompt += "\nProvide 3 remediation steps with kubectl commands."
	
	// Get AI suggestions
	response, err := e.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return []Remediation{}
	}
	
	// Parse remediations (simplified)
	return []Remediation{
		{
			Title:       "Check pod logs",
			Description: "Inspect pod logs for error messages",
			Command:     fmt.Sprintf("kubectl logs %s -n %s", resourceName, e.namespace),
			Confidence:  "High",
			Safe:        true,
		},
		{
			Title:       "Describe pod",
			Description: "Get detailed pod information and events",
			Command:     fmt.Sprintf("kubectl describe pod %s -n %s", resourceName, e.namespace),
			Confidence:  "High",
			Safe:        true,
		},
		{
			Title:       response.Content[:min(len(response.Content), 50)] + "...",
			Description: "AI-suggested remediation",
			Command:     "# See AI response for details",
			Confidence:  "Medium",
			Safe:        true,
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Display displays the diagnostic report
func (r *Report) Display() {
	fmt.Printf("\n%s\n", r.Summary)
	fmt.Printf("Health Score: %d/100\n\n", r.HealthScore)
	
	if len(r.Issues) == 0 {
		fmt.Println("✅ No issues detected!")
		return
	}
	
	fmt.Println("Issues found:")
	for i, issue := range r.Issues {
		severityEmoji := "ℹ️"
		switch issue.Severity {
		case SeverityCritical:
			severityEmoji = "🔴"
		case SeverityHigh:
			severityEmoji = "🟠"
		case SeverityMedium:
			severityEmoji = "🟡"
		}
		
		fmt.Printf("\n%d. %s [%s] %s\n", i+1, severityEmoji, issue.Severity, issue.Type)
		fmt.Printf("   Resource: %s\n", issue.Resource)
		fmt.Printf("   %s\n", issue.Description)
	}
}
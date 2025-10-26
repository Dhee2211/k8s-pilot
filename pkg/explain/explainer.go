package explain

import (
	"context"
	"fmt"
	"strings"

	"k8s-pilot/pkg/ai"
	"k8s-pilot/pkg/k8s"
)

// Explainer provides AI-powered explanations
type Explainer struct {
	aiProvider ai.Provider
	k8sClient  *k8s.Client
	namespace  string
}

// NewExplainer creates a new explainer
func NewExplainer(namespace string) *Explainer {
	k8sClient, _ := k8s.NewClient(namespace)
	
	aiConfig := &ai.Config{
		Provider: ai.ProviderMock,
	}
	aiProvider, _ := ai.NewProvider(aiConfig)
	
	return &Explainer{
		aiProvider: aiProvider,
		k8sClient:  k8sClient,
		namespace:  namespace,
	}
}

// Explanation represents an AI-generated explanation
type Explanation struct {
	Query           string
	Answer          string
	RelatedCommands []string
	Tip             string
}

// Explain generates an explanation for a query
func (e *Explainer) Explain(query string) (*Explanation, error) {
	ctx := context.Background()
	
	// Determine what type of explanation is needed
	queryLower := strings.ToLower(query)
	
	var content string
	var err error
	
	if strings.Contains(queryLower, "logs") {
		content, err = e.explainLogs(ctx, query)
	} else if strings.Contains(queryLower, "events") {
		content, err = e.explainEvents(ctx, query)
	} else if strings.Contains(queryLower, "pod") || strings.Contains(queryLower, "deployment") {
		content, err = e.explainResource(ctx, query)
	} else {
		content, err = e.explainConcept(ctx, query)
	}
	
	if err != nil {
		return nil, err
	}
	
	explanation := &Explanation{
		Query:           query,
		Answer:          content,
		RelatedCommands: e.extractCommands(content),
		Tip:             e.generateTip(query),
	}
	
	return explanation, nil
}

// explainLogs explains pod logs
func (e *Explainer) explainLogs(ctx context.Context, query string) (string, error) {
	// Extract pod name from query
	words := strings.Fields(query)
	podName := ""
	for i, word := range words {
		if word == "logs" && i+1 < len(words) {
			podName = words[i+1]
			break
		}
	}
	
	if podName == "" {
		return "Please specify a pod name. Example: kubectl-pilot explain logs mypod", nil
	}
	
	// Get the actual logs
	logs, err := e.k8sClient.GetPodLogs(ctx, podName, "", e.namespace, 50)
	if err != nil {
		return fmt.Sprintf("Could not retrieve logs: %v", err), nil
	}
	
	// Use AI to summarize and explain the logs
	prompt := fmt.Sprintf(`Analyze these Kubernetes pod logs and explain what's happening:

Pod: %s
Logs:
%s

Provide:
1. A summary of what the application is doing
2. Any errors or warnings present
3. Recommendations if issues are found`, podName, logs)
	
	response, err := e.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return "", err
	}
	
	return response.Content, nil
}

// explainEvents explains Kubernetes events
func (e *Explainer) explainEvents(ctx context.Context, query string) (string, error) {
	events, err := e.k8sClient.GetEvents(ctx, e.namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get events: %w", err)
	}
	
	// Summarize recent events
	eventSummary := fmt.Sprintf("Recent events in namespace %s:\n\n", e.namespace)
	for i, event := range events.Items {
		if i >= 10 {
			break // Limit to 10 events
		}
		eventSummary += fmt.Sprintf("- [%s] %s: %s\n", 
			event.Type, event.Reason, event.Message)
	}
	
	prompt := fmt.Sprintf(`Analyze these Kubernetes events and explain what they mean:

%s

Provide a summary of cluster activity and any issues that need attention.`, eventSummary)
	
	response, err := e.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return "", err
	}
	
	return response.Content, nil
}

// explainResource explains a specific resource
func (e *Explainer) explainResource(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`User asked: "%s"

Explain this Kubernetes resource or concept clearly and concisely.
Include:
1. What the resource does
2. Common use cases
3. Best practices
4. Example kubectl commands`, query)
	
	response, err := e.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return "", err
	}
	
	return response.Content, nil
}

// explainConcept explains a general Kubernetes concept
func (e *Explainer) explainConcept(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`Explain this Kubernetes concept or question:

"%s"

Provide a clear, educational explanation that helps the user understand the concept.
Include practical examples and kubectl commands where relevant.`, query)
	
	response, err := e.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return "", err
	}
	
	return response.Content, nil
}

// extractCommands extracts kubectl commands from the explanation
func (e *Explainer) extractCommands(content string) []string {
	var commands []string
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "kubectl") {
			commands = append(commands, trimmed)
		}
	}
	
	return commands
}

// generateTip generates a helpful tip based on the query
func (e *Explainer) generateTip(query string) string {
	queryLower := strings.ToLower(query)
	
	if strings.Contains(queryLower, "logs") {
		return "Use -f flag to follow logs in real-time: kubectl logs -f <pod>"
	}
	
	if strings.Contains(queryLower, "events") {
		return "Filter events by type with --field-selector: kubectl get events --field-selector type=Warning"
	}
	
	if strings.Contains(queryLower, "pod") {
		return "Use 'kubectl describe pod' to see detailed information including events"
	}
	
	return "Use 'kubectl explain <resource>' to see detailed documentation"
}

// Display displays the explanation
func (ex *Explanation) Display() {
	fmt.Printf("\nQuery: %s\n\n", ex.Query)
	fmt.Println(ex.Answer)
}

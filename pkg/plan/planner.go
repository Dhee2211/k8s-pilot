package plan

import (
	"context"
	"fmt"
	"os"
	"strings"

	"k8s-pilot/pkg/ai"
)

// Planner generates execution plans from natural language
type Planner struct {
	aiProvider ai.Provider
	namespace  string
	dryRun     bool
}

// NewPlanner creates a new planner
func NewPlanner(namespace string, dryRun bool) *Planner {
	// Initialize AI provider (using mock for now)
	aiConfig := &ai.Config{
		Provider: ai.ProviderMock,
	}
	
	// Try to use configured provider if available
	if providerStr := os.Getenv("K8S_PILOT_AI_PROVIDER"); providerStr != "" {
		aiConfig.Provider = ai.ProviderType(providerStr)
		aiConfig.APIKey = os.Getenv("K8S_PILOT_AI_KEY")
		aiConfig.Model = os.Getenv("K8S_PILOT_AI_MODEL")
	}
	
	provider, err := ai.NewProvider(aiConfig)
	if err != nil {
		// Fallback to mock
		provider, _ = ai.NewMockProvider(aiConfig)
	}
	
	return &Planner{
		aiProvider: provider,
		namespace:  namespace,
		dryRun:     dryRun,
	}
}

// Plan represents an execution plan
type Plan struct {
	Summary      string
	Commands     []Command
	Warnings     []string
	RequiresAuth bool
	DryRun       bool
}

// Command represents a kubectl command
type Command struct {
	Command     string
	Description string
	Safe        bool
	DryRun      bool
}

// Generate generates an execution plan from natural language
func (p *Planner) Generate(query string) (*Plan, error) {
	ctx := context.Background()
	
	// Build the prompt for the AI
	prompt := p.buildPrompt(query)
	
	// Generate the plan using AI
	response, err := p.aiProvider.Generate(ctx, prompt, ai.DefaultOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}
	
	// Parse the response into a plan
	plan := p.parseResponse(response.Content, query)
	plan.DryRun = p.dryRun
	
	return plan, nil
}

// buildPrompt builds the AI prompt for plan generation
func (p *Planner) buildPrompt(query string) string {
	systemPrompt := `You are a Kubernetes expert assistant. Your task is to translate natural language queries into safe kubectl commands.

Rules:
1. Always prefer dry-run commands when possible
2. Never suggest commands that delete critical resources without warning
3. Include explanations for each command
4. Warn about potentially dangerous operations
5. Use the namespace provided in context when applicable
6. Suggest RBAC-safe alternatives when possible

Format your response as:
SUMMARY: <brief summary>
COMMANDS:
- <kubectl command> | <description> | <safe: true/false>

WARNINGS:
- <warning if any>
`
	
	namespace := p.namespace
	if namespace == "" {
		namespace = "default"
	}
	
	userPrompt := fmt.Sprintf(`Namespace: %s
Query: %s

Generate a safe execution plan for this query.`, namespace, query)
	
	return systemPrompt + "\n\n" + userPrompt
}

// parseResponse parses the AI response into a Plan
func (p *Planner) parseResponse(response string, originalQuery string) *Plan {
	plan := &Plan{
		Summary:  "Execution plan for: " + originalQuery,
		Commands: []Command{},
		Warnings: []string{},
		DryRun:   p.dryRun,
	}
	
	lines := strings.Split(response, "\n")
	inCommands := false
	inWarnings := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "SUMMARY:") {
			plan.Summary = strings.TrimSpace(strings.TrimPrefix(line, "SUMMARY:"))
			continue
		}
		
		if strings.HasPrefix(line, "COMMANDS:") {
			inCommands = true
			inWarnings = false
			continue
		}
		
		if strings.HasPrefix(line, "WARNINGS:") {
			inWarnings = true
			inCommands = false
			continue
		}
		
		if inCommands && strings.HasPrefix(line, "-") || strings.HasPrefix(line, "kubectl") {
			// Parse command line
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				cmd := strings.TrimSpace(strings.TrimPrefix(parts[0], "-"))
				desc := strings.TrimSpace(parts[1])
				safe := true
				if len(parts) >= 3 && strings.Contains(parts[2], "false") {
					safe = false
				}
				
				plan.Commands = append(plan.Commands, Command{
					Command:     cmd,
					Description: desc,
					Safe:        safe,
					DryRun:      p.dryRun,
				})
			}
		}
		
		if inWarnings && strings.HasPrefix(line, "-") {
			warning := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if warning != "" {
				plan.Warnings = append(plan.Warnings, warning)
			}
		}
	}
	
	// If no commands were parsed, create a simple one from the response
	if len(plan.Commands) == 0 {
		plan.Commands = append(plan.Commands, Command{
			Command:     fmt.Sprintf("# Generated from: %s", originalQuery),
			Description: "See AI response for details",
			Safe:        true,
			DryRun:      p.dryRun,
		})
	}
	
	return plan
}

// Display displays the plan
func (p *Plan) Display() {
	fmt.Printf("\n%s\n\n", p.Summary)
	
	if len(p.Warnings) > 0 {
		fmt.Println("⚠️  Warnings:")
		for _, warning := range p.Warnings {
			fmt.Printf("  • %s\n", warning)
		}
		fmt.Println()
	}
	
	fmt.Println("Commands to execute:")
	for i, cmd := range p.Commands {
		safetyIndicator := "✓"
		if !cmd.Safe {
			safetyIndicator = "⚠"
		}
		
		fmt.Printf("\n%d. [%s] %s\n", i+1, safetyIndicator, cmd.Description)
		fmt.Printf("   %s\n", cmd.Command)
	}
	
	if p.DryRun {
		fmt.Println("\n(Dry-run mode - no changes will be applied)")
	}
}

// Execute executes the plan
func (p *Plan) Execute() (*Result, error) {
	result := &Result{
		ExecutedCommands: []string{},
		Errors:           []string{},
	}
	
	for _, cmd := range p.Commands {
		if p.DryRun && !cmd.DryRun {
			result.ExecutedCommands = append(result.ExecutedCommands, 
				fmt.Sprintf("[DRY-RUN] %s", cmd.Command))
			continue
		}
		
		// TODO: Actually execute the kubectl command
		// For now, just record it
		result.ExecutedCommands = append(result.ExecutedCommands, cmd.Command)
	}
	
	return result, nil
}

// Result represents the result of plan execution
type Result struct {
	ExecutedCommands []string
	Errors           []string
}

// Display displays the result
func (r *Result) Display() {
	fmt.Println("\nExecuted commands:")
	for i, cmd := range r.ExecutedCommands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}
	
	if len(r.Errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, err := range r.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
	}
}

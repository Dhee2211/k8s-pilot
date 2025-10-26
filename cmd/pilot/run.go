package pilot

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s-pilot/pkg/plan"
)

var (
	applyChanges bool
)

var runCmd = &cobra.Command{
	Use:   "run [natural language query]",
	Short: "Execute a natural language Kubernetes operation",
	Long: `Translates a natural language query into a safe, auditable Kubernetes operation plan.
By default, runs in dry-run mode. Use --apply to execute the plan.

Examples:
  kubectl-pilot run "restart failing pods in payments namespace"
  kubectl-pilot run "scale deployment api to 5 replicas" --apply
  kubectl-pilot run "list pods with high memory usage"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		
		planner := plan.NewPlanner(namespace, dryRun)
		
		// Generate execution plan from natural language
		executionPlan, err := planner.Generate(query)
		if err != nil {
			return fmt.Errorf("failed to generate plan: %w", err)
		}
		
		// Display the plan
		fmt.Println("\n📋 Execution Plan:")
		fmt.Println("─────────────────")
		executionPlan.Display()
		
		if dryRun && !applyChanges {
			fmt.Println("\n✓ Dry-run complete. Use --apply to execute the plan.")
			return nil
		}
		
		if !applyChanges {
			fmt.Println("\n✓ Preview complete. Use --apply to execute.")
			return nil
		}
		
		// Execute the plan
		fmt.Println("\n⚡ Executing plan...")
		result, err := executionPlan.Execute()
		if err != nil {
			return fmt.Errorf("execution failed: %w", err)
		}
		
		fmt.Println("\n✅ Execution complete:")
		result.Display()
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVar(&applyChanges, "apply", false, "apply the generated plan (disables dry-run)")
}

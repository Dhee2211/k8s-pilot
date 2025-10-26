package pilot

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s-pilot/pkg/diagnose"
)

var (
	resourceType string
	resourceName string
	allNamespaces bool
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose [resource-type] [resource-name]",
	Short: "Diagnose Kubernetes cluster issues",
	Long: `Automatically detect and diagnose common Kubernetes issues including:
  • CrashLoopBackOff
  • ImagePullBackOff
  • Probe failures
  • PVC issues
  • Resource constraints
  • Network problems

The diagnostics engine aggregates logs, events, and resource state to provide
root-cause hypotheses and ranked remediation suggestions.

Examples:
  kubectl-pilot diagnose pod myapp-pod
  kubectl-pilot diagnose deployment myapp -n production
  kubectl-pilot diagnose --all-namespaces`,
	RunE: func(cmd *cobra.Command, args []string) error {
		engine := diagnose.NewEngine(namespace, allNamespaces)
		
		var report *diagnose.Report
		var err error
		
		if len(args) >= 2 {
			// Diagnose specific resource
			report, err = engine.DiagnoseResource(args[0], args[1])
		} else if len(args) == 1 {
			// Diagnose resource type
			report, err = engine.DiagnoseResourceType(args[0])
		} else {
			// Diagnose entire namespace/cluster
			report, err = engine.DiagnoseCluster()
		}
		
		if err != nil {
			return fmt.Errorf("diagnostics failed: %w", err)
		}
		
		// Display diagnostic report
		fmt.Println("\n🔍 Diagnostic Report:")
		fmt.Println("══════════════════════")
		report.Display()
		
		// Show recommended fixes
		if len(report.Remediations) > 0 {
			fmt.Println("\n💡 Recommended Fixes:")
			fmt.Println("────────────────────")
			for i, remedy := range report.Remediations {
				fmt.Printf("\n%d. %s (Confidence: %s)\n", i+1, remedy.Title, remedy.Confidence)
				fmt.Printf("   %s\n", remedy.Description)
				fmt.Printf("   Command: %s\n", remedy.Command)
			}
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "diagnose across all namespaces")
}

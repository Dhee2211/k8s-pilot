package pilot

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s-pilot/pkg/explain"
)

var explainCmd = &cobra.Command{
	Use:   "explain [query]",
	Short: "Get AI-powered explanations of Kubernetes resources and concepts",
	Long: `Provides AI-powered explanations of Kubernetes resources, logs, events, and concepts.
The explain mode helps you understand what's happening in your cluster and teaches
you the underlying kubectl commands.

Examples:
  kubectl-pilot explain logs mypod
  kubectl-pilot explain events in payments namespace
  kubectl-pilot explain "why is my pod pending"
  kubectl-pilot explain deployment myapp`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		
		explainer := explain.NewExplainer(namespace)
		
		// Generate explanation
		explanation, err := explainer.Explain(query)
		if err != nil {
			return fmt.Errorf("failed to generate explanation: %w", err)
		}
		
		// Display explanation
		fmt.Println("\n📚 Explanation:")
		fmt.Println("═══════════════")
		explanation.Display()
		
		// Show related commands if any
		if len(explanation.RelatedCommands) > 0 {
			fmt.Println("\n🔧 Related kubectl Commands:")
			fmt.Println("────────────────────────────")
			for _, cmd := range explanation.RelatedCommands {
				fmt.Printf("  • %s\n", cmd)
			}
		}
		
		// Show educational tips
		if explanation.Tip != "" {
			fmt.Printf("\n💡 Tip: %s\n", explanation.Tip)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}

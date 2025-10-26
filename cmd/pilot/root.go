package pilot

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s-pilot/internal/config"
	"k8s-pilot/internal/logger"
)

var (
	cfgFile   string
	dryRun    bool
	verbose   bool
	namespace string
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-pilot",
	Short: "AI-powered Kubernetes operations assistant",
	Long: `kubectl-pilot is an AI-powered CLI tool that translates natural language
into safe, auditable Kubernetes actions, diagnoses cluster issues, and provides
explainable fixes.

Features:
  • Natural language to kubectl commands
  • Cluster diagnostics and remediation
  • Dry-run by default with safety guardrails
  • RBAC-aware command generation
  • Multi-cloud support (GKE, EKS, AKS, K3s, Kind)
  • Extensible plugin system`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(verbose)
		if cfgFile != "" {
			config.Load(cfgFile)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-pilot.yaml)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", true, "preview changes without applying them")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "kubernetes namespace")
}

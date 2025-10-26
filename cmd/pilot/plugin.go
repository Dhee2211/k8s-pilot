package pilot

import (
	"fmt"

	"k8s-pilot/pkg/plugins"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage kubectl-pilot plugins",
	Long: `Install, list, and manage plugins that extend kubectl-pilot functionality.
Plugins can add custom detectors, fixers, and diagnostics capabilities.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugins.NewManager()
		pluginList := manager.List()

		if len(pluginList) == 0 {
			fmt.Println("No plugins installed.")
			return nil
		}

		fmt.Println("\n🔌 Installed Plugins:")
		fmt.Println("═══════════════════")
		for _, p := range pluginList {
			fmt.Printf("\n• %s (v%s)\n", p.Name(), p.Version())
			fmt.Printf("  %s\n", p.Description())
		}

		return nil
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [plugin-name]",
	Short: "Install a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugins.NewManager()

		fmt.Printf("Installing plugin: %s...\n", args[0])

		if err := manager.InstallByName(args[0]); err != nil {
			return fmt.Errorf("failed to install plugin: %w", err)
		}

		fmt.Println("✅ Plugin installed successfully!")
		return nil
	},
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [plugin-name]",
	Short: "Uninstall a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugins.NewManager()

		if err := manager.Uninstall(args[0]); err != nil {
			return fmt.Errorf("failed to uninstall plugin: %w", err)
		}

		fmt.Println("✅ Plugin uninstalled successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
}

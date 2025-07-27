package commands

import (
	"fmt"
	"os"

	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/spf13/cobra"
)

var (
	configMgr *config.Manager
	appConfig *config.Config
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "flow-test-go",
	Short: "A CLI tool for orchestrating AI agents with flows",
	Long: `flow-test-go is a CLI tool that orchestrates AI agents using LangGraph and OpenRouter.

It supports:
- Flow-based AI orchestration
- MCP (Model Context Protocol) server integration  
- GitHub API integration
- Multiple AI provider support via OpenRouter
- Configuration management`,
	Version: "1.0.0",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		// Initialize configuration manager
		var err error
		configMgr, err = config.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize config manager: %w", err)
		}

		// Load configuration
		appConfig, err = configMgr.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Disable help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// Add subcommands
	rootCmd.AddCommand(listCmd)
}

// GetConfigManager returns the configuration manager instance.
func GetConfigManager() *config.Manager {
	return configMgr
}

// GetConfig returns the application configuration.
func GetConfig() *config.Config {
	return appConfig
}

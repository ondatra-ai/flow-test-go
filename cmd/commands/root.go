package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/spf13/cobra"
)

var (
	configMgr *config.Manager
	appConfig *config.Config
	initMutex sync.Mutex
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
		return initializeConfig()
	},
}

// initializeConfig initializes the configuration manager and config in a thread-safe manner.
func initializeConfig() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	// Skip if already initialized
	if configMgr != nil && appConfig != nil {
		return nil
	}

	// Initialize config manager
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
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
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

// ResetGlobalState resets global state for testing purposes.
// This should only be used in tests.
func ResetGlobalState() {
	initMutex.Lock()
	defer initMutex.Unlock()
	configMgr = nil
	appConfig = nil
}

// GetConfig returns the application configuration.
func GetConfig() *config.Config {
	return appConfig
}

// CreateRootCmd creates a new root command instance for testing.
// This avoids sharing command state between tests.
func CreateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flow-test-go",
		Short: "A CLI tool for orchestrating AI agents with flows",
		Long: `flow-test-go is a CLI tool that orchestrates AI agents using LangGraph and OpenRouter.

It provides:
- Flow definition and execution
- MCP (Model Context Protocol) server integration
- AI agent orchestration
- Configuration management`,
		Version: "1.0.0",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeConfig()
		},
	}

	// Add subcommands
	// Create a new instance of listCmd to avoid sharing state
	newListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available flows",
		Long:  "List all available flows in the current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listCmd.RunE(cmd, args)
		},
	}

	// Copy flags from original listCmd
	newListCmd.Flags().BoolP("details", "d", false, "show detailed information about each flow")

	cmd.AddCommand(newListCmd)
	return cmd
}

package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/spf13/cobra"
)

// GlobalState holds the global application state.
type GlobalState struct {
	configMgr *config.Manager
	appConfig *config.Config
	initMutex sync.Mutex
}

// NewGlobalState creates a new GlobalState instance.
func NewGlobalState() *GlobalState {
	return &GlobalState{
		configMgr: nil,
		appConfig: nil,
		initMutex: sync.Mutex{},
	}
}

// createBaseCommand creates a base cobra.Command with common settings.
func createBaseCommand(state *GlobalState) *cobra.Command {
	return &cobra.Command{
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
			DisableDefaultCmd:   true,
			DisableNoDescFlag:   false,
			DisableDescriptions: false,
			HiddenDefaultCmd:    false,
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initializeConfig(state)
		},
		Aliases:                []string{},
		SuggestFor:             []string{},
		GroupID:                "",
		Example:                "",
		ValidArgs:              []string{},
		ValidArgsFunction:      nil,
		Args:                   nil,
		ArgAliases:             []string{},
		BashCompletionFunction: "",
		Deprecated:             "",
		Annotations:            map[string]string{},
		PersistentPreRun:       nil,
		PreRun:                 nil,
		PreRunE:                nil,
		Run:                    nil,
		RunE:                   nil,
		PostRun:                nil,
		PostRunE:               nil,
		PersistentPostRun:      nil,
		PersistentPostRunE:     nil,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: false,
		},
		TraverseChildren:           false,
		Hidden:                     false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}
}

// initializeConfig initializes the configuration manager and config in a thread-safe manner.
func initializeConfig(state *GlobalState) error {
	state.initMutex.Lock()
	defer state.initMutex.Unlock()

	// Skip if already initialized
	if state.configMgr != nil && state.appConfig != nil {
		return nil
	}

	// Initialize config manager
	var err error

	state.configMgr, err = config.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize config manager: %w", err)
	}

	// Load configuration
	state.appConfig, err = state.configMgr.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(state *GlobalState) {
	// Create and setup the root command with subcommands
	rootCmd := createRootCommand(state)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// createRootCommand creates and configures the root command.
func createRootCommand(state *GlobalState) *cobra.Command {
	rootCmd := createBaseCommand(state)

	// Disable help command
	rootCmd.SetHelpCommand(createDisabledHelpCommand())

	// Add subcommands
	rootCmd.AddCommand(CreateListCommand(state))

	return rootCmd
}

// createDisabledHelpCommand creates a disabled help command to hide the default help.
func createDisabledHelpCommand() *cobra.Command {
	return &cobra.Command{
		Use:                    "no-help",
		Hidden:                 true,
		Aliases:                []string{},
		SuggestFor:             []string{},
		Short:                  "",
		GroupID:                "",
		Long:                   "",
		Example:                "",
		ValidArgs:              []string{},
		ValidArgsFunction:      nil,
		Args:                   nil,
		ArgAliases:             []string{},
		BashCompletionFunction: "",
		Deprecated:             "",
		Annotations:            map[string]string{},
		Version:                "",
		PersistentPreRun:       nil,
		PersistentPreRunE:      nil,
		PreRun:                 nil,
		PreRunE:                nil,
		Run:                    nil,
		RunE:                   nil,
		PostRun:                nil,
		PostRunE:               nil,
		PersistentPostRun:      nil,
		PersistentPostRunE:     nil,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: false,
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   false,
			DisableNoDescFlag:   false,
			DisableDescriptions: false,
			HiddenDefaultCmd:    false,
		},
		TraverseChildren:           false,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
	}
}

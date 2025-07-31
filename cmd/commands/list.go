// Package commands provides CLI command implementations for the flow-test-go tool.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateListCommand creates and returns the list command.
func CreateListCommand(state *GlobalState) *cobra.Command {
	cmd := createBaseListCommand()
	cmd.RunE = func(cobraCmd *cobra.Command, args []string) error {
		return listFlows(cobraCmd, args, state)
	}

	return cmd
}

// createBaseListCommand creates the base command structure for list.
func createBaseListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available flows",
		Long: `List all available flows in the .flows/flows directory.

Examples:
  flow-test-go list`,
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

// listFlows implements the list command logic.
func listFlows(cmd *cobra.Command, _ []string, state *GlobalState) error {
	flows, err := state.configMgr.ListFlows()
	if err != nil {
		return fmt.Errorf("failed to list flows: %w", err)
	}

	if len(flows) == 0 {
		cmd.Println("📁 No flows found in .flows/flows directory")
		cmd.Println("💡 Use 'flow-test-go init' to create example flows")

		return nil
	}

	cmd.Printf("📋 Found %d flow(s):\n\n", len(flows))

	for _, flowID := range flows {
		cmd.Printf("📄 %s\n", flowID)
	}

	return nil
}

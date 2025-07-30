// Package commands provides CLI command implementations for the flow-test-go tool.
package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var showDetails bool

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available flows",
	Long: `List all available flows in the .flows/flows directory.

Examples:
  flow-test-go list
  flow-test-go list --details`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		flows, err := configMgr.ListFlows()
		if err != nil {
			return fmt.Errorf("failed to list flows: %w", err)
		}

		if len(flows) == 0 {
			cmd.Println("📁 No flows found in .flows/flows directory")
			cmd.Println("💡 Use 'flow-test-go init' to create example flows")

			return nil
		}

		cmd.Printf("📋 Found %d flow(s):\n\n", len(flows))

		// Get the details flag value from this command instance
		details, _ := cmd.Flags().GetBool("details")

		for _, flowID := range flows {
			if details {
				// Load flow to get details
				flow, err := configMgr.LoadFlow(flowID)
				if err != nil {
					cmd.Printf("❌ %s (failed to load: %v)\n", flowID, err)

					continue
				}

				cmd.Printf("📄 %s\n", flow.ID)
				cmd.Printf("   Name: %s\n", flow.Name)
				cmd.Printf("   Description: %s\n", flow.Description)
				cmd.Printf("   Steps: %d\n", len(flow.Steps))

				if len(flow.Variables) > 0 {
					var vars []string
					for k := range flow.Variables {
						vars = append(vars, k)
					}
					cmd.Printf("   Variables: %s\n", strings.Join(vars, ", "))
				}
				cmd.Println()
			} else {
				cmd.Printf("📄 %s\n", flowID)
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&showDetails, "details", false, "show detailed information about each flow")
}

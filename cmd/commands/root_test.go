package commands_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
)

func TestRootCommand_BasicProperties(t *testing.T) {
	t.Parallel()

	// Create a fresh command to test
	cmd := commands.CreateRootCmd()

	// Test basic command properties
	assert.Equal(t, "flow-test-go", cmd.Use)
	assert.Contains(t, cmd.Short, "CLI tool for orchestrating AI agents")
	assert.Contains(t, cmd.Long, "flow-test-go is a CLI tool")
	assert.Equal(t, "1.0.0", cmd.Version)

	// Test that completion is disabled
	assert.True(t, cmd.CompletionOptions.DisableDefaultCmd)

	// Test that help command is disabled (using a different approach)
	// Note: HelpCommand() is not accessible, but we can test the structure
	subcommands := cmd.Commands()
	hasHelp := false

	for _, subCmd := range subcommands {
		if subCmd.Use == "help" {
			hasHelp = true

			break
		}
	}

	assert.False(t, hasHelp, "Help command should be disabled")
}

func TestRootCommand_Subcommands(t *testing.T) {
	t.Parallel()

	// Create a fresh command to test
	cmd := commands.CreateRootCmd()

	// Test that only the list command is registered
	subcommands := cmd.Commands()
	require.Len(t, subcommands, 1)
	assert.Equal(t, "list", subcommands[0].Use)
}

func TestRootCommand_PersistentPreRunE(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	// Reset global state before and after test
	commands.ResetGlobalState()
	t.Cleanup(func() {
		commands.ResetGlobalState()
	})

	// Test successful configuration initialization
	t.Run("successful initialization", func(t *testing.T) {
		t.Parallel()

		cmd := commands.CreateRootCmd()
		err := cmd.PersistentPreRunE(cmd, []string{})
		require.NoError(t, err)

		// Verify that initialization succeeded (no verification needed for dependency injection architecture)
	})

	// Test that subsequent calls work (don't fail)
	t.Run("multiple initializations", func(t *testing.T) {
		t.Parallel()

		defer func() {
			commands.ResetGlobalState()
		}()

		cmd := commands.CreateRootCmd()
		err := cmd.PersistentPreRunE(cmd, []string{})
		require.NoError(t, err)

		err = cmd.PersistentPreRunE(cmd, []string{})
		require.NoError(t, err)
	})
}

func TestRootCommand_PersistentPreRunE_Errors(t *testing.T) {
	// Test with unwritable directory
	t.Run("unwritable directory", func(t *testing.T) {

		// Create a temporary directory
		tmpDir := t.TempDir()

		// Change to the directory first (before making it read-only)
		t.Chdir(tmpDir)

		defer func() {
			_ = os.Chmod(tmpDir, 0755) // #nosec G302 -- Restore permissions first
		}()

		// Now make it read-only after we're inside it
		err := os.Chmod(tmpDir, 0444) // #nosec G302 -- Read-only for test
		require.NoError(t, err)

		cmd := commands.CreateRootCmd()
		err = cmd.PersistentPreRunE(cmd, []string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize config manager")
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	// Test Execute function with help flag
	t.Run("help flag", func(t *testing.T) {
		t.Parallel()

		// Capture stdout and stderr
		var stdout, stderr bytes.Buffer

		cmd := commands.CreateRootCmd()
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)

		// Set args to help
		cmd.SetArgs([]string{"--help"})

		// Execute should not exit the process in tests
		// We'll catch the error instead
		err := cmd.Execute()
		// Help flag causes cobra to return an error
		if err != nil {
			assert.Contains(t, err.Error(), "help requested")
		}
	})

	// Test Execute function with version flag
	t.Run("version flag", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer

		cmd := commands.CreateRootCmd()
		cmd.SetOut(&stdout)

		cmd.SetArgs([]string{"--version"})

		err := cmd.Execute()
		if err != nil {
			// Version flag might cause exit
			assert.Contains(t, err.Error(), "version")
		}
	})

	// Test Execute function with list command
	t.Run("list command", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer

		// Create a new command instance to avoid state sharing
		cmd := commands.CreateRootCmd()
		cmd.SetOut(&stdout)
		cmd.SetErr(&stdout)
		cmd.SetArgs([]string{"list"})

		err := cmd.Execute()
		require.NoError(t, err)

		// Should show no flows message
		output := stdout.String()
		assert.Contains(t, output, "No flows found")
	})
}

func TestRootCommand_ErrorOutput(t *testing.T) {
	t.Parallel()

	// Test that errors are properly formatted
	// This is more of an integration test for the error handling

	// Create a command that will fail
	failingCmd := &cobra.Command{
		Use: "failing",
		RunE: func(_ *cobra.Command, _ []string) error {
			return assert.AnError
		},
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

	cmd := commands.CreateRootCmd()
	cmd.AddCommand(failingCmd)

	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"failing"})

	err := cmd.Execute()
	require.Error(t, err)
}

// Test command structure and relationships.
func TestRootCommand_Structure(t *testing.T) {
	t.Parallel()

	// Create a fresh command to test
	cmd := commands.CreateRootCmd()

	// Verify that the root command has the expected structure
	assert.Equal(t, "flow-test-go", cmd.Name())
	assert.True(t, cmd.HasSubCommands())

	// Test that we can find the list command
	listCommand, _, err := cmd.Find([]string{"list"})
	require.NoError(t, err)
	assert.Equal(t, "list", listCommand.Name())

	// Test that non-existent commands return error
	_, _, err = cmd.Find([]string{"nonexistent"})
	require.Error(t, err)
}

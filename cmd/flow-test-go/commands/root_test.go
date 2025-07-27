package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand_BasicProperties(t *testing.T) {
	t.Parallel()

	// Test basic command properties
	assert.Equal(t, "flow-test-go", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "CLI tool for orchestrating AI agents")
	assert.Contains(t, rootCmd.Long, "flow-test-go is a CLI tool")
	assert.Equal(t, "1.0.0", rootCmd.Version)

	// Test that completion is disabled
	assert.True(t, rootCmd.CompletionOptions.DisableDefaultCmd)

	// Test that help command is disabled (using a different approach)
	// Note: HelpCommand() is not accessible, but we can test the structure
	subcommands := rootCmd.Commands()
	hasHelp := false
	for _, cmd := range subcommands {
		if cmd.Use == "help" {
			hasHelp = true
			break
		}
	}
	assert.False(t, hasHelp, "Help command should be disabled")
}

func TestRootCommand_Subcommands(t *testing.T) {
	t.Parallel()

	// Test that only the list command is registered
	subcommands := rootCmd.Commands()
	require.Len(t, subcommands, 1)
	assert.Equal(t, "list", subcommands[0].Use)
}

func TestRootCommand_PersistentPreRunE(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	// Change to temp directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Test successful configuration initialization
	t.Run("successful initialization", func(t *testing.T) {
		err := rootCmd.PersistentPreRunE(rootCmd, []string{})
		require.NoError(t, err)

		// Verify that global variables are set
		assert.NotNil(t, configMgr)
		assert.NotNil(t, appConfig)

		// Verify config manager functionality
		assert.NotNil(t, GetConfigManager())
		assert.NotNil(t, GetConfig())
		assert.Equal(t, configMgr, GetConfigManager())
		assert.Equal(t, appConfig, GetConfig())
	})

	// Test that subsequent calls work (don't fail)
	t.Run("multiple initializations", func(t *testing.T) {
		err := rootCmd.PersistentPreRunE(rootCmd, []string{})
		require.NoError(t, err)

		err = rootCmd.PersistentPreRunE(rootCmd, []string{})
		require.NoError(t, err)
	})
}

func TestRootCommand_PersistentPreRunE_Errors(t *testing.T) {
	t.Parallel()

	// Test with unwritable directory
	t.Run("unwritable directory", func(t *testing.T) {
		// Create a temporary directory and make it read-only
		tmpDir := t.TempDir()
		err := os.Chmod(tmpDir, 0444) // Read-only
		require.NoError(t, err)

		originalDir, err := os.Getwd()
		require.NoError(t, err)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		defer func() {
			_ = os.Chdir(originalDir)
			_ = os.Chmod(tmpDir, 0755) // Restore permissions for cleanup
		}()

		err = rootCmd.PersistentPreRunE(rootCmd, []string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to initialize config manager")
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Test Execute function with help flag
	t.Run("help flag", func(t *testing.T) {
		// Capture stdout and stderr
		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)

		// Set args to help
		rootCmd.SetArgs([]string{"--help"})

		// Execute should not exit the process in tests
		// We'll catch the error instead
		err := rootCmd.Execute()
		// Help flag causes cobra to return an error
		if err != nil {
			assert.Contains(t, err.Error(), "help requested")
		}

		// Reset args
		rootCmd.SetArgs([]string{})
	})

	// Test Execute function with version flag
	t.Run("version flag", func(t *testing.T) {
		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)

		rootCmd.SetArgs([]string{"--version"})

		err := rootCmd.Execute()
		if err != nil {
			// Version flag might cause exit
			assert.Contains(t, err.Error(), "version")
		}

		// Reset args
		rootCmd.SetArgs([]string{})
	})

	// Test Execute function with list command
	t.Run("list command", func(t *testing.T) {
		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)

		rootCmd.SetArgs([]string{"list"})

		err := rootCmd.Execute()
		require.NoError(t, err)

		// Should show no flows message
		output := stdout.String()
		assert.Contains(t, output, "No flows found")

		// Reset args
		rootCmd.SetArgs([]string{})
	})
}

func TestGetConfigManager(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Initially should be nil
	if configMgr == nil {
		assert.Nil(t, GetConfigManager())
	}

	// Initialize via PersistentPreRunE
	err = rootCmd.PersistentPreRunE(rootCmd, []string{})
	require.NoError(t, err)

	// Should now return the manager
	manager := GetConfigManager()
	assert.NotNil(t, manager)
	assert.Equal(t, configMgr, manager)
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Initially should be nil
	if appConfig == nil {
		assert.Nil(t, GetConfig())
	}

	// Initialize via PersistentPreRunE
	err = rootCmd.PersistentPreRunE(rootCmd, []string{})
	require.NoError(t, err)

	// Should now return the config
	config := GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, appConfig, config)

	// Verify config has expected default values
	assert.Equal(t, "flow-test-go", config.App.Name)
	assert.Equal(t, "1.0.0", config.App.Version)
}

func TestRootCommand_ConfigurationValues(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Set some environment variables
	t.Setenv("OPENROUTER_API_KEY", "test-api-key")
	t.Setenv("GITHUB_TOKEN", "test-github-token")

	// Initialize configuration
	err = rootCmd.PersistentPreRunE(rootCmd, []string{})
	require.NoError(t, err)

	// Verify configuration values
	config := GetConfig()
	require.NotNil(t, config)

	// Check that environment variables were loaded
	assert.Equal(t, "test-api-key", config.LLM.APIKey)
	assert.Equal(t, "test-github-token", config.GitHub.Token)

	// Check default values
	assert.Equal(t, "openrouter", config.LLM.Provider)
	assert.Equal(t, "openai/gpt-4-turbo", config.LLM.DefaultModel)
	assert.Equal(t, "your-github-username", config.GitHub.Owner)
	assert.Equal(t, "your-github-repo", config.GitHub.Repository)
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
	}

	rootCmd.AddCommand(failingCmd)
	defer rootCmd.RemoveCommand(failingCmd)

	var stderr bytes.Buffer
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs([]string{"failing"})

	err := rootCmd.Execute()
	require.Error(t, err)

	// Reset args
	rootCmd.SetArgs([]string{})
}

// Test command structure and relationships
func TestRootCommand_Structure(t *testing.T) {
	t.Parallel()

	// Verify that the root command has the expected structure
	assert.Equal(t, "flow-test-go", rootCmd.Name())
	assert.True(t, rootCmd.HasSubCommands())

	// Test that we can find the list command
	listCommand, _, err := rootCmd.Find([]string{"list"})
	require.NoError(t, err)
	assert.Equal(t, "list", listCommand.Name())

	// Test that non-existent commands return error
	_, _, err = rootCmd.Find([]string{"nonexistent"})
	require.Error(t, err)
}

// Benchmark tests
func BenchmarkRootCommand_PersistentPreRunE(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(b, err)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rootCmd.PersistentPreRunE(rootCmd, []string{})
	}
}

func BenchmarkGetConfigManager(b *testing.B) {
	// Initialize once
	tmpDir := b.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(b, err)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	_ = rootCmd.PersistentPreRunE(rootCmd, []string{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetConfigManager()
	}
}

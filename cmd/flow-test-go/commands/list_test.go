package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/peterovchinnikov/flow-test-go/pkg/types"
)

func TestListCommand(t *testing.T) {
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

	// Create test manager and config
	manager, err := config.NewManager()
	require.NoError(t, err)

	appConfig, err := manager.LoadConfig()
	require.NoError(t, err)

	// Set global variables (normally set by root command)
	configMgr = manager
	appConfig = appConfig

	// Test case 1: No flows (empty directory)
	t.Run("no flows", func(t *testing.T) {
		var output bytes.Buffer
		cmd := &cobra.Command{
			Use: "list",
			RunE: func(_ *cobra.Command, _ []string) error {
				flows, err := configMgr.ListFlows()
				if err != nil {
					return err
				}

				if len(flows) == 0 {
					output.WriteString("📁 No flows found in .flows/flows directory\n")
					output.WriteString("💡 Use 'flow-test-go init' to create example flows\n")
					return nil
				}

				output.WriteString("📋 Found flows\n")
				return nil
			},
		}

		err := cmd.RunE(cmd, []string{})
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "No flows found")
		assert.Contains(t, outputStr, "Use 'flow-test-go init'")
	})

	// Test case 2: With flows, basic listing
	t.Run("with flows basic", func(t *testing.T) {
		// Create some test flows
		testFlows := []*types.FlowDefinition{
			{
				ID:          "flow1",
				Name:        "Test Flow 1",
				Description: "First test flow",
				Steps: map[string]types.Step{
					"step1": {Type: types.StepTypeEnd},
				},
			},
			{
				ID:          "flow2",
				Name:        "Test Flow 2",
				Description: "Second test flow",
				Steps: map[string]types.Step{
					"step1": {Type: types.StepTypeEnd},
					"step2": {Type: types.StepTypeEnd},
				},
			},
		}

		for _, flow := range testFlows {
			err := manager.SaveFlow(flow)
			require.NoError(t, err)
		}

		var output bytes.Buffer
		cmd := &cobra.Command{
			Use: "list",
			RunE: func(_ *cobra.Command, _ []string) error {
				flows, err := configMgr.ListFlows()
				if err != nil {
					return err
				}

				output.WriteString("📋 Found flows:\n")
				for _, flowID := range flows {
					output.WriteString("📄 " + flowID + "\n")
				}
				return nil
			},
		}

		err := cmd.RunE(cmd, []string{})
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "📋 Found flows")
		assert.Contains(t, outputStr, "📄 flow1")
		assert.Contains(t, outputStr, "📄 flow2")
	})

	// Test case 3: With flows, detailed listing
	t.Run("with flows detailed", func(t *testing.T) {
		var output bytes.Buffer
		showDetails := true

		cmd := &cobra.Command{
			Use: "list",
			RunE: func(_ *cobra.Command, _ []string) error {
				flows, err := configMgr.ListFlows()
				if err != nil {
					return err
				}

				output.WriteString("📋 Found flows:\n\n")
				for _, flowID := range flows {
					if showDetails {
						flow, err := configMgr.LoadFlow(flowID)
						if err != nil {
							output.WriteString("❌ " + flowID + " (failed to load)\n")
							continue
						}

						output.WriteString("📄 " + flow.ID + "\n")
						output.WriteString("   Name: " + flow.Name + "\n")
						output.WriteString("   Description: " + flow.Description + "\n")
						output.WriteString("   Steps: " + string(rune(len(flow.Steps))) + "\n\n")
					}
				}
				return nil
			},
		}

		err := cmd.RunE(cmd, []string{})
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "📋 Found flows")
		assert.Contains(t, outputStr, "📄 flow1")
		assert.Contains(t, outputStr, "Name: Test Flow 1")
		assert.Contains(t, outputStr, "Description: First test flow")
	})
}

func TestListCommand_Integration(t *testing.T) {
	// This test verifies the actual list command behavior
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

	// Initialize manager and set global variables
	manager, err := config.NewManager()
	require.NoError(t, err)

	appConfig, err := manager.LoadConfig()
	require.NoError(t, err)

	configMgr = manager
	appConfig = appConfig

	// Create a test flow
	testFlow := &types.FlowDefinition{
		ID:          "integration-test-flow",
		Name:        "Integration Test Flow",
		Description: "A flow for integration testing",
		Variables:   map[string]string{"test": "value"},
		Steps: map[string]types.Step{
			"step1": {
				Type: types.StepTypePrompt,
				Prompt: &types.PromptConfig{
					Template: "Hello {{.name}}",
				},
				Next: "step2",
			},
			"step2": {
				Type: types.StepTypeEnd,
			},
		},
	}

	err = manager.SaveFlow(testFlow)
	require.NoError(t, err)

	// Test the actual list command
	t.Run("actual command execution", func(t *testing.T) {
		// Capture stdout
		var output bytes.Buffer
		listCmd.SetOut(&output)

		// Execute the command
		err := listCmd.RunE(listCmd, []string{})
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "Found 1 flow")
		assert.Contains(t, outputStr, "integration-test-flow")
	})

	// Test with details flag
	t.Run("with details flag", func(t *testing.T) {
		// Reset the showDetails flag
		showDetails = true
		defer func() { showDetails = false }()

		var output bytes.Buffer
		listCmd.SetOut(&output)

		err := listCmd.RunE(listCmd, []string{})
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "integration-test-flow")
		assert.Contains(t, outputStr, "Name: Integration Test Flow")
		assert.Contains(t, outputStr, "Description: A flow for integration testing")
		assert.Contains(t, outputStr, "Steps: 2")
		assert.Contains(t, outputStr, "Variables: test")
	})
}

func TestListCommand_ErrorHandling(t *testing.T) {
	t.Parallel()

	// Test with invalid config manager
	t.Run("nil config manager", func(t *testing.T) {
		// Temporarily set configMgr to nil
		originalConfigMgr := configMgr
		configMgr = nil
		defer func() { configMgr = originalConfigMgr }()

		// This should panic or error gracefully
		assert.Panics(t, func() {
			_ = listCmd.RunE(listCmd, []string{})
		})
	})

	// Test with corrupted flow file
	t.Run("corrupted flow file", func(t *testing.T) {
		// Create a temporary directory for testing
		tmpDir := t.TempDir()
		originalDir, err := os.Getwd()
		require.NoError(t, err)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)
		defer func() {
			_ = os.Chdir(originalDir)
		}()

		// Initialize manager
		manager, err := config.NewManager()
		require.NoError(t, err)

		configMgr = manager

		// Create a corrupted flow file
		flowsDir := ".flows/flows"
		err = os.MkdirAll(flowsDir, 0750)
		require.NoError(t, err)

		corruptedFlowPath := flowsDir + "/corrupted-flow.json"
		err = os.WriteFile(corruptedFlowPath, []byte("invalid json content"), 0600)
		require.NoError(t, err)

		// Set showDetails to true to trigger loading
		showDetails = true
		defer func() { showDetails = false }()

		var output bytes.Buffer
		listCmd.SetOut(&output)

		err = listCmd.RunE(listCmd, []string{})
		require.NoError(t, err) // Should handle errors gracefully

		outputStr := output.String()
		// Should show an error for the corrupted flow
		assert.Contains(t, outputStr, "failed to load")
	})
}

// Test helper functions
func TestShowDetails_Flag(t *testing.T) {
	t.Parallel()

	// Test that the showDetails flag exists and can be set
	flag := listCmd.Flags().Lookup("details")
	require.NotNil(t, flag)
	assert.Equal(t, "bool", flag.Value.Type())
	assert.Equal(t, "false", flag.DefValue)
	assert.Contains(t, flag.Usage, "show detailed information")
}

// Benchmark tests
func BenchmarkListCommand_NoFlows(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(b, err)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	manager, err := config.NewManager()
	require.NoError(b, err)

	configMgr = manager

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = listCmd.RunE(listCmd, []string{})
	}
}

func BenchmarkListCommand_WithFlows(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	originalDir, err := os.Getwd()
	require.NoError(b, err)

	err = os.Chdir(tmpDir)
	require.NoError(b, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	manager, err := config.NewManager()
	require.NoError(b, err)

	configMgr = manager

	// Create test flows
	for i := 0; i < 10; i++ {
		flow := &types.FlowDefinition{
			ID:          "bench-flow-" + string(rune(i)),
			Name:        "Benchmark Flow",
			Description: "A flow for benchmarking",
			Steps: map[string]types.Step{
				"step1": {Type: types.StepTypeEnd},
			},
		}
		_ = manager.SaveFlow(flow)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = listCmd.RunE(listCmd, []string{})
	}
}

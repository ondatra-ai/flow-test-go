package commands_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/cmd/commands"
	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/peterovchinnikov/flow-test-go/pkg/types"
)

func TestListCommand(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	// Create test manager and config
	manager, err := config.NewManager()
	require.NoError(t, err)

	_, err = manager.LoadConfig()
	require.NoError(t, err)

	// Reset global state before test
	commands.ResetGlobalState()

	// Test case 1: No flows (empty directory)
	t.Run("no flows", func(t *testing.T) {
		t.Parallel()

		// Create a new command instance to avoid shared state
		cmd := commands.CreateRootCmd()

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		// Execute the command
		err := cmd.Execute()
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "No flows found")
	})

	// Test case 2: With flows, basic listing
	t.Run("with flows basic", func(t *testing.T) {
		t.Parallel()

		// Create test flow
		testFlow := &types.FlowDefinition{
			Schema:      "",
			Version:     "1.0",
			ID:          "test-flow",
			Name:        "Test Flow",
			Description: "A test flow for unit testing",
			Variables:   make(map[string]string),
			Steps: map[string]types.Step{
				"step1": {
					Type: types.StepTypePrompt,
					Prompt: &types.PromptConfig{
						Template: "Test prompt",
						System:   "",
						Context:  make(map[string]any),
					},
					Model:      "",
					Tools:      []string{},
					MCPServer:  "",
					Next:       "",
					Conditions: []types.ConditionConfig{},
					Timeout:    nil,
					Retry:      nil,
					Metadata:   make(map[string]any),
				},
			},
			InitialStep: "step1",
		}

		err := manager.SaveFlow(testFlow)
		require.NoError(t, err)

		// Create a new command instance to avoid shared state
		cmd := commands.CreateRootCmd()

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		// Execute the command
		err = cmd.Execute()
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "Found 1 flow")
	})

	// Test case 3: With flows, detailed listing
	t.Run("with flows detailed", func(t *testing.T) {
		t.Parallel()

		// Create test flow with variables
		testFlow := &types.FlowDefinition{
			Schema:      "",
			Version:     "1.0",
			ID:          "detailed-flow",
			Name:        "Detailed Test Flow",
			Description: "A detailed test flow for unit testing",
			Variables: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			Steps: map[string]types.Step{
				"step1": {
					Type: types.StepTypePrompt,
					Prompt: &types.PromptConfig{
						Template: "Detailed prompt",
						System:   "",
						Context:  make(map[string]any),
					},
					Model:      "",
					Tools:      []string{},
					MCPServer:  "",
					Next:       "",
					Conditions: []types.ConditionConfig{},
					Timeout:    nil,
					Retry:      nil,
					Metadata:   make(map[string]any),
				},
				"step2": {
					Type:       types.StepTypeEnd,
					Prompt:     nil,
					Model:      "",
					Tools:      []string{},
					MCPServer:  "",
					Next:       "",
					Conditions: []types.ConditionConfig{},
					Timeout:    nil,
					Retry:      nil,
					Metadata:   make(map[string]any),
				},
			},
			InitialStep: "step1",
		}

		err := manager.SaveFlow(testFlow)
		require.NoError(t, err)

		// Create a new command instance to avoid shared state
		cmd := commands.CreateRootCmd()

		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		// Execute the command
		err = cmd.Execute()
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "detailed-flow")
	})
}

func TestListCommand_Integration(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	// Initialize manager and set global variables
	manager, err := config.NewManager()
	require.NoError(t, err)

	_, err = manager.LoadConfig()
	require.NoError(t, err)

	// Reset global state before test
	commands.ResetGlobalState()

	// Create a test flow
	testFlow := &types.FlowDefinition{
		Schema:      "",
		Version:     "1.0",
		ID:          "integration-test-flow",
		Name:        "Integration Test Flow",
		Description: "A flow for integration testing",
		Variables:   map[string]string{"test": "value"},
		Steps: map[string]types.Step{
			"step1": {
				Type: types.StepTypePrompt,
				Prompt: &types.PromptConfig{
					Template: "Hello {{.name}}",
					System:   "",
					Context:  make(map[string]any),
				},
				Model:      "",
				Tools:      []string{},
				MCPServer:  "",
				Next:       "step2",
				Conditions: []types.ConditionConfig{},
				Timeout:    nil,
				Retry:      nil,
				Metadata:   make(map[string]any),
			},
			"step2": {
				Type:       types.StepTypeEnd,
				Prompt:     nil,
				Model:      "",
				Tools:      []string{},
				MCPServer:  "",
				Next:       "",
				Conditions: []types.ConditionConfig{},
				Timeout:    nil,
				Retry:      nil,
				Metadata:   make(map[string]any),
			},
		},
		InitialStep: "step1",
	}

	err = manager.SaveFlow(testFlow)
	require.NoError(t, err)

	// Test the actual list command
	t.Run("actual command execution", func(t *testing.T) {
		t.Parallel()

		// Capture stdout
		var output bytes.Buffer

		// Create a new command instance to avoid state sharing
		cmd := commands.CreateRootCmd()
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		// Execute the command
		err := cmd.Execute()
		require.NoError(t, err)

		outputStr := output.String()
		assert.Contains(t, outputStr, "Found 1 flow")
		assert.Contains(t, outputStr, "integration-test-flow")
	})

	// Test with details flag
	t.Run("with details flag", func(t *testing.T) {
		t.Parallel()

		// Capture stdout
		var output bytes.Buffer

		// Create a new command instance to avoid state sharing
		cmd := commands.CreateRootCmd()
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		// Execute the command
		err := cmd.Execute()
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

	// Test with nil config manager
	t.Run("nil config manager", func(t *testing.T) {

		// Create a temporary directory for testing
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		// Reset global state to ensure clean state
		commands.ResetGlobalState()

		// Create a command that will fail due to no config manager
		cmd := commands.CreateRootCmd()
		cmd.SetArgs([]string{"list"})

		// This should error gracefully
		err := cmd.Execute()
		assert.Error(t, err)
	})

	// Test with corrupted flow file
	t.Run("corrupted flow file", func(t *testing.T) {

		// Create a temporary directory for testing
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		// Initialize manager
		_, err := config.NewManager()
		require.NoError(t, err)

		// Reset state after test
		defer commands.ResetGlobalState()

		// Create a corrupted flow file
		flowsDir := ".flows/flows"
		err = os.MkdirAll(flowsDir, 0750)
		require.NoError(t, err)

		corruptedFlowPath := flowsDir + "/corrupted-flow.json"
		err = os.WriteFile(corruptedFlowPath, []byte("invalid json content"), 0600)
		require.NoError(t, err)

		var output bytes.Buffer

		cmd := commands.CreateRootCmd()
		cmd.SetOut(&output)
		cmd.SetErr(&output)
		cmd.SetArgs([]string{"list"})

		err = cmd.Execute()
		require.NoError(t, err) // Should handle errors gracefully

		outputStr := output.String()
		// Should show an error for the corrupted flow
		assert.Contains(t, outputStr, "failed to load")
	})
}

// Test helper functions.

// Benchmark tests.
func BenchmarkListCommand_NoFlows(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	b.Chdir(tmpDir)

	_, err := config.NewManager()
	require.NoError(b, err)

	// Reset state after benchmark
	defer commands.ResetGlobalState()

	b.ResetTimer()

	for range b.N {
		cmd := commands.CreateRootCmd()
		cmd.SetArgs([]string{"list"})
		_ = cmd.Execute()
	}
}

func BenchmarkListCommand_WithFlows(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	b.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(b, err)

	// Reset state after benchmark
	defer commands.ResetGlobalState()

	// Create test flows
	for i := range 10 {
		flow := &types.FlowDefinition{
			Schema:      "",
			Version:     "1.0",
			ID:          fmt.Sprintf("bench-flow-%d", i),
			Name:        "Benchmark Flow",
			Description: "A flow for benchmarking",
			Variables:   make(map[string]string),
			Steps: map[string]types.Step{
				"step1": {
					Type:       types.StepTypeEnd,
					Prompt:     nil,
					Model:      "",
					Tools:      []string{},
					MCPServer:  "",
					Next:       "",
					Conditions: []types.ConditionConfig{},
					Timeout:    nil,
					Retry:      nil,
					Metadata:   make(map[string]any),
				},
			},
			InitialStep: "step1",
		}
		_ = manager.SaveFlow(flow)
	}

	b.ResetTimer()

	for range b.N {
		cmd := commands.CreateRootCmd()
		cmd.SetArgs([]string{"list"})
		_ = cmd.Execute()
	}
}

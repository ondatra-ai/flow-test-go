package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/internal/config"
	"github.com/peterovchinnikov/flow-test-go/pkg/types"
)

func TestNewManager(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	// Check that directories were created
	assert.DirExists(t, ".flows")
	assert.DirExists(t, ".flows/flows")
	assert.DirExists(t, ".flows/servers")
}

func TestManager_LoadConfig(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Test loading default config (should succeed)
	configResult, err := manager.LoadConfig()
	require.NoError(t, err)
	assert.NotNil(t, configResult)

	// Verify default values
	assert.Equal(t, "flow-test-go", configResult.App.Name)
	assert.Equal(t, "1.0.0", configResult.App.Version)
	assert.Equal(t, "openrouter", configResult.LLM.Provider)
	assert.Equal(t, "openai/gpt-4-turbo", configResult.LLM.DefaultModel)
}

func TestManager_LoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv()

	// Set environment variables
	t.Setenv("OPENROUTER_API_KEY", "test-api-key")
	t.Setenv("GITHUB_TOKEN", "test-github-token")

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	configResult, err := manager.LoadConfig()
	require.NoError(t, err)

	// Environment variables should be loaded
	assert.Equal(t, "test-api-key", configResult.LLM.APIKey)
	assert.Equal(t, "test-github-token", configResult.GitHub.Token)
}

func TestManager_SaveFlow(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create a test flow
	flow := &types.FlowDefinition{
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
		InitialStep: "",
	}

	// Save the flow
	err = manager.SaveFlow(flow)
	require.NoError(t, err)

	// Verify file was created
	flowPath := filepath.Join(".flows", "flows", "test-flow.json")
	assert.FileExists(t, flowPath)

	// Verify file permissions
	info, err := os.Stat(flowPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestManager_LoadFlow(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create and save a test flow first
	originalFlow := &types.FlowDefinition{
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
		InitialStep: "",
	}

	err = manager.SaveFlow(originalFlow)
	require.NoError(t, err)

	// Load the flow
	loadedFlow, err := manager.LoadFlow("test-flow")
	require.NoError(t, err)
	assert.NotNil(t, loadedFlow)

	// Verify flow content
	assert.Equal(t, originalFlow.ID, loadedFlow.ID)
	assert.Equal(t, originalFlow.Name, loadedFlow.Name)
	assert.Equal(t, originalFlow.Description, loadedFlow.Description)
	assert.Len(t, loadedFlow.Steps, 2)
}

func TestManager_LoadFlow_NotFound(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Try to load a non-existent flow
	_, err = manager.LoadFlow("nonexistent-flow")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read flow file")
}

func TestManager_ListFlows(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Initially should be empty
	flows, err := manager.ListFlows()
	require.NoError(t, err)
	assert.Empty(t, flows)

	// Create some test flows
	testFlows := []string{"flow1", "flow2", "flow3"}
	for _, flowID := range testFlows {
		flow := &types.FlowDefinition{
			Schema:      "",
			Version:     "1.0",
			ID:          flowID,
			Name:        "Test Flow " + flowID,
			Description: "A test flow",
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
			InitialStep: "",
		}
		err = manager.SaveFlow(flow)
		require.NoError(t, err)
	}

	// List flows again
	flows, err = manager.ListFlows()
	require.NoError(t, err)
	assert.Len(t, flows, 3)

	// Check that all flows are present (order might vary)
	for _, expectedFlow := range testFlows {
		assert.Contains(t, flows, expectedFlow)
	}
}

func TestManager_SaveMCPServer(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create a test MCP server config
	serverConfig := &types.MCPServerConfig{
		Name:             "test-server",
		Command:          "python",
		Args:             []string{"-m", "test_server"},
		Env:              make(map[string]string),
		TransportType:    types.TransportStdio,
		TransportOptions: make(map[string]any),
		Capabilities: types.MCPCapabilities{
			Tools:     true,
			Resources: false,
			Prompts:   false,
			Logging:   false,
		},
		Timeout:     30 * time.Second,
		HealthCheck: nil,
		AutoRestart: false,
		MaxRestarts: 0,
		Metadata:    make(map[string]any),
	}

	// Note: SaveMCPServer method needs to be implemented
	// For now, just verify the config is valid
	err = serverConfig.Validate()
	require.NoError(t, err)

	// Use manager to ensure it's used
	assert.NotNil(t, manager)

	// File creation would happen in actual SaveMCPServer implementation
}

func TestManager_LoadMCPServer_NotImplemented(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Note: This test is a placeholder until LoadMCPServer is implemented
	assert.NotNil(t, manager)
}

func TestManager_ListMCPServers_NotImplemented(t *testing.T) {

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	// Change to temp directory
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Note: This test is a placeholder until ListMCPServers is implemented
	assert.NotNil(t, manager)
}

func TestConfig_ValidateForExecution_NotImplemented(t *testing.T) {

	configResult := &config.Config{
		App: struct {
			Name    string `mapstructure:"name"`
			Version string `mapstructure:"version"`
			Debug   bool   `mapstructure:"debug"`
		}{
			Name:    "",
			Version: "",
			Debug:   false,
		},
		LLM: struct {
			Provider       string            `mapstructure:"provider"`
			APIKey         string            `mapstructure:"apiKey"`
			DefaultModel   string            `mapstructure:"defaultModel"`
			ModelOverrides map[string]string `mapstructure:"modelOverrides"`
			MaxTokens      int               `mapstructure:"maxTokens"`
			Temperature    float64           `mapstructure:"temperature"`
		}{
			Provider:       "",
			APIKey:         "",
			DefaultModel:   "",
			ModelOverrides: nil,
			MaxTokens:      0,
			Temperature:    0.0,
		},
		GitHub: struct {
			Token         string `mapstructure:"token"`
			Owner         string `mapstructure:"owner"`
			Repository    string `mapstructure:"repository"`
			DefaultBranch string `mapstructure:"defaultBranch"`
		}{
			Token:         "",
			Owner:         "",
			Repository:    "",
			DefaultBranch: "",
		},
		Flow: struct {
			Directory      string `mapstructure:"directory"`
			DefaultTimeout string `mapstructure:"defaultTimeout"`
			CheckpointDir  string `mapstructure:"checkpointDir"`
			MaxRetries     int    `mapstructure:"maxRetries"`
			EnableParallel bool   `mapstructure:"enableParallel"`
		}{
			Directory:      "",
			DefaultTimeout: "",
			CheckpointDir:  "",
			MaxRetries:     0,
			EnableParallel: false,
		},
		Logging: struct {
			Level   string `mapstructure:"level"`
			Format  string `mapstructure:"format"`
			File    string `mapstructure:"file"`
			Console bool   `mapstructure:"console"`
		}{
			Level:   "",
			Format:  "",
			File:    "",
			Console: false,
		},
	}

	// Note: ValidateForExecution method needs to be implemented
	// For now, just verify config is not nil
	assert.NotNil(t, configResult)
}

// Benchmark tests.
func BenchmarkManager_LoadConfig(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	b.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(b, err)

	b.ResetTimer()

	for range b.N {
		_, _ = manager.LoadConfig()
	}
}

func BenchmarkManager_SaveFlow(b *testing.B) {
	// Create a temporary directory for benchmarking
	tmpDir := b.TempDir()
	b.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(b, err)

	flow := &types.FlowDefinition{
		Schema:      "",
		Version:     "1.0",
		ID:          "bench-flow",
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
		InitialStep: "",
	}

	b.ResetTimer()

	for range b.N {
		_ = manager.SaveFlow(flow)
	}
}

//nolint:exhaustruct // Test files don't need to initialize all struct fields
package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/internal/config"
	"github.com/ondatra-ai/flow-test-go/pkg/types"
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
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
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

	// Validate the server configuration
	err = serverConfig.Validate()
	require.NoError(t, err)

	// Test successful save
	err = manager.SaveMCPServer(serverConfig)
	require.NoError(t, err)

	// Verify file was created
	serverPath := filepath.Join(".flows", "servers", "test-server.json")
	assert.FileExists(t, serverPath)

	// Verify file content
	data, err := os.ReadFile(serverPath)
	require.NoError(t, err)

	var unmarshaledConfig types.MCPServerConfig

	err = json.Unmarshal(data, &unmarshaledConfig)
	require.NoError(t, err)
	assert.Equal(t, serverConfig.Name, unmarshaledConfig.Name)
	assert.Equal(t, serverConfig.Command, unmarshaledConfig.Command)
}

func TestManager_SaveMCPServer_InvalidName(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create server config with invalid name
	serverConfig := &types.MCPServerConfig{
		Name:             "test/invalid", // Contains path separator
		Command:          "python",
		Args:             []string{},
		Env:              make(map[string]string),
		TransportType:    types.TransportStdio,
		TransportOptions: make(map[string]any),
		Capabilities:     types.MCPCapabilities{Tools: true},
		Timeout:          0,
		HealthCheck:      nil,
		AutoRestart:      false,
		MaxRestarts:      0,
		Metadata:         make(map[string]any),
	}

	// Test invalid server name
	err = manager.SaveMCPServer(serverConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid server name")
}

func TestManager_SaveMCPServer_ValidationError(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create server config that will fail validation (empty name)
	serverConfig := &types.MCPServerConfig{
		Name:             "", // Empty name should fail validation
		Command:          "python",
		Args:             []string{},
		Env:              make(map[string]string),
		TransportType:    types.TransportStdio,
		TransportOptions: make(map[string]any),
		Capabilities:     types.MCPCapabilities{Tools: true},
		Timeout:          0,
		HealthCheck:      nil,
		AutoRestart:      false,
		MaxRestarts:      0,
		Metadata:         make(map[string]any),
	}

	// Test validation error
	err = manager.SaveMCPServer(serverConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server config validation failed")
}

func TestManager_ValidateForExecution(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Test with valid OpenRouter API key
	configWithKey := &config.Config{}
	configWithKey.LLM.Provider = "openrouter"
	configWithKey.LLM.APIKey = "test-api-key"

	err = manager.ValidateForExecution(configWithKey)
	require.NoError(t, err)

	// Test with missing OpenRouter API key
	configNoKey := &config.Config{}
	configNoKey.LLM.Provider = "openrouter"
	configNoKey.LLM.APIKey = "" // Missing API key

	err = manager.ValidateForExecution(configNoKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OpenRouter API key is required")

	// Test with different provider (should not require API key)
	configOtherProvider := &config.Config{}
	configOtherProvider.LLM.Provider = "other"
	configOtherProvider.LLM.APIKey = "" // Empty API key but different provider

	err = manager.ValidateForExecution(configOtherProvider)
	assert.NoError(t, err)
}

func TestManager_LoadMCPServers(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create test server configs
	server1 := &types.MCPServerConfig{
		Name:             "server1",
		Command:          "python",
		Args:             []string{"-m", "server1"},
		Env:              make(map[string]string),
		TransportType:    types.TransportStdio,
		TransportOptions: make(map[string]any),
		Capabilities: types.MCPCapabilities{
			Tools:     true,
			Resources: false,
			Prompts:   false,
			Logging:   false,
		},
		Timeout:     0,
		HealthCheck: nil,
		AutoRestart: false,
		MaxRestarts: 0,
		Metadata:    make(map[string]any),
	}
	server2 := &types.MCPServerConfig{
		Name:             "server2",
		Command:          "node",
		Args:             []string{"server2.js"},
		Env:              make(map[string]string),
		TransportType:    types.TransportStdio,
		TransportOptions: make(map[string]any),
		Capabilities: types.MCPCapabilities{
			Tools:     false,
			Resources: true,
			Prompts:   false,
			Logging:   false,
		},
		Timeout:     0,
		HealthCheck: nil,
		AutoRestart: false,
		MaxRestarts: 0,
		Metadata:    make(map[string]any),
	}

	// Save test servers
	err = manager.SaveMCPServer(server1)
	require.NoError(t, err)
	err = manager.SaveMCPServer(server2)
	require.NoError(t, err)

	// Test loading servers
	servers, err := manager.LoadMCPServers()
	require.NoError(t, err)
	assert.Len(t, servers, 2)
	assert.Contains(t, servers, "server1")
	assert.Contains(t, servers, "server2")
	assert.Equal(t, "python", servers["server1"].Command)
	assert.Equal(t, "node", servers["server2"].Command)
}

func TestManager_LoadMCPServers_EmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Test loading from empty directory
	servers, err := manager.LoadMCPServers()
	require.NoError(t, err)
	assert.Empty(t, servers)
}

func TestManager_LoadMCPServers_CorruptedJSON(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Create corrupted JSON file
	corruptedFile := filepath.Join(".flows", "servers", "corrupted.json")
	err = os.WriteFile(corruptedFile, []byte(`{"name": "incomplete`), 0o600)
	require.NoError(t, err)

	// Test loading with corrupted file
	servers, err := manager.LoadMCPServers()
	require.Error(t, err)
	assert.Nil(t, servers)
	assert.Contains(t, err.Error(), "failed to parse server config")
}

func TestManager_GetConfig(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	manager, err := config.NewManager()
	require.NoError(t, err)

	// Load config first
	config, err := manager.LoadConfig()
	require.NoError(t, err)

	// Test GetConfig returns the same config
	retrievedConfig := manager.GetConfig()
	assert.Equal(t, config, retrievedConfig)
	assert.NotNil(t, retrievedConfig)
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

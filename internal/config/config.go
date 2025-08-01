// Package config provides configuration management for the flow-test-go application.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/ondatra-ai/flow-test-go/pkg/types"
)

var (
	// ErrOpenRouterAPIKeyRequired is returned when OpenRouter API key is missing.
	ErrOpenRouterAPIKeyRequired = errors.New("OpenRouter API key is required for execution " +
		"(set OPENROUTER_API_KEY env var or llm.apiKey in config)")

	// ErrInvalidFlowID is returned when a flow ID contains path separators.
	ErrInvalidFlowID = errors.New("invalid flow ID: must not contain path separators")

	// ErrInvalidServerName is returned when a server name contains path separators.
	ErrInvalidServerName = errors.New("invalid server name: must not contain path separators")
)

// Config represents the application configuration.
type Config struct {
	// Application settings
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Debug   bool   `mapstructure:"debug"`
	} `mapstructure:"app"`

	// LLM settings
	LLM struct {
		Provider       string            `mapstructure:"provider"`
		APIKey         string            `mapstructure:"apiKey"`
		DefaultModel   string            `mapstructure:"defaultModel"`
		ModelOverrides map[string]string `mapstructure:"modelOverrides"`
		MaxTokens      int               `mapstructure:"maxTokens"`
		Temperature    float64           `mapstructure:"temperature"`
	} `mapstructure:"llm"`

	// GitHub settings
	GitHub struct {
		Token         string `mapstructure:"token"`
		Owner         string `mapstructure:"owner"`
		Repository    string `mapstructure:"repository"`
		DefaultBranch string `mapstructure:"defaultBranch"`
	} `mapstructure:"github"`

	// Flow settings
	Flow struct {
		Directory      string `mapstructure:"directory"`
		DefaultTimeout string `mapstructure:"defaultTimeout"`
		CheckpointDir  string `mapstructure:"checkpointDir"`
		MaxRetries     int    `mapstructure:"maxRetries"`
		EnableParallel bool   `mapstructure:"enableParallel"`
	} `mapstructure:"flow"`

	// Logging settings
	Logging struct {
		Level   string `mapstructure:"level"`
		Format  string `mapstructure:"format"`
		File    string `mapstructure:"file"`
		Console bool   `mapstructure:"console"`
	} `mapstructure:"logging"`
}

// Manager handles configuration loading and management.
type Manager struct {
	config     *Config
	configDir  string
	flowsDir   string
	serversDir string
}

// NewManager creates a new configuration manager.
func NewManager() (*Manager, error) {
	// Default to .flows in current directory
	configDir := ".flows"

	// Create config directories (MkdirAll is idempotent and handles concurrent creation)
	const dirPerms = 0o750

	err := os.MkdirAll(configDir, dirPerms)
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	err = os.MkdirAll(filepath.Join(configDir, "flows"), dirPerms)
	if err != nil {
		return nil, fmt.Errorf("failed to create flows directory: %w", err)
	}

	err = os.MkdirAll(filepath.Join(configDir, "servers"), dirPerms)
	if err != nil {
		return nil, fmt.Errorf("failed to create servers directory: %w", err)
	}

	return &Manager{
		config:     nil,
		configDir:  configDir,
		flowsDir:   filepath.Join(configDir, "flows"),
		serversDir: filepath.Join(configDir, "servers"),
	}, nil
}

// LoadConfig loads the application configuration.
func (cm *Manager) LoadConfig() (*Config, error) {
	// Set configuration file search paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cm.configDir)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.flow-test-go")

	// Set defaults
	cm.setDefaults()

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FLOW_TEST_GO")

	// Read config file
	var err error

	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found is not an error, we'll use defaults
	}

	// Unmarshal configuration
	config := cm.createDefaultConfig()

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load direct environment variables (fallback for common env vars)
	if config.LLM.APIKey == "" {
		if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
			config.LLM.APIKey = apiKey
		}
	}

	if config.GitHub.Token == "" {
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			config.GitHub.Token = token
		}
	}

	// Validate required settings
	err = cm.validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cm.config = config

	return config, nil
}

// LoadFlow loads a flow definition by ID.
func (cm *Manager) LoadFlow(flowID string) (*types.FlowDefinition, error) {
	// Validate flowID contains no path separators
	if strings.ContainsAny(flowID, "/\\") || strings.Contains(flowID, "..") {
		return nil, ErrInvalidFlowID
	}

	flowPath := filepath.Join(cm.flowsDir, flowID+".json")

	data, err := os.ReadFile(flowPath) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read flow file %s: %w", flowPath, err)
	}

	var flow types.FlowDefinition

	err = json.Unmarshal(data, &flow)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flow definition: %w", err)
	}

	err = flow.Validate()
	if err != nil {
		return nil, fmt.Errorf("flow validation failed: %w", err)
	}

	return &flow, nil
}

// ListFlows returns a list of available flow IDs.
func (cm *Manager) ListFlows() ([]string, error) {
	files, err := os.ReadDir(cm.flowsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read flows directory: %w", err)
	}

	var flows []string

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			flowID := file.Name()[:len(file.Name())-5] // Remove .json extension
			flows = append(flows, flowID)
		}
	}

	return flows, nil
}

// LoadMCPServers loads all MCP server configurations.
func (cm *Manager) LoadMCPServers() (map[string]*types.MCPServerConfig, error) {
	files, err := os.ReadDir(cm.serversDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read servers directory: %w", err)
	}

	servers := make(map[string]*types.MCPServerConfig)

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		serverPath := filepath.Join(cm.serversDir, file.Name())

		data, err := os.ReadFile(serverPath) // #nosec G304
		if err != nil {
			return nil, fmt.Errorf("failed to read server config %s: %w", serverPath, err)
		}

		var serverConfig types.MCPServerConfig

		err = json.Unmarshal(data, &serverConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse server config %s: %w", serverPath, err)
		}

		err = serverConfig.Validate()
		if err != nil {
			return nil, fmt.Errorf("server config validation failed for %s: %w", serverConfig.Name, err)
		}

		servers[serverConfig.Name] = &serverConfig
	}

	return servers, nil
}

// SaveFlow saves a flow definition.
func (cm *Manager) SaveFlow(flow *types.FlowDefinition) error {
	// Validate flow ID contains no path separators
	if strings.ContainsAny(flow.ID, "/\\") || strings.Contains(flow.ID, "..") {
		return ErrInvalidFlowID
	}

	err := flow.Validate()
	if err != nil {
		return fmt.Errorf("flow validation failed: %w", err)
	}

	data, err := json.MarshalIndent(flow, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal flow: %w", err)
	}

	flowPath := filepath.Join(cm.flowsDir, flow.ID+".json")

	const filePerms = 0o600

	err = os.WriteFile(flowPath, data, filePerms)
	if err != nil {
		return fmt.Errorf("failed to write flow file: %w", err)
	}

	return nil
}

// SaveMCPServer saves an MCP server configuration.
func (cm *Manager) SaveMCPServer(server *types.MCPServerConfig) error {
	// Validate server name contains no path separators
	if strings.ContainsAny(server.Name, "/\\") || strings.Contains(server.Name, "..") {
		return ErrInvalidServerName
	}

	err := server.Validate()
	if err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	data, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal server config: %w", err)
	}

	const filePerms = 0o600

	serverPath := filepath.Join(cm.serversDir, server.Name+".json")

	err = os.WriteFile(serverPath, data, filePerms)
	if err != nil {
		return fmt.Errorf("failed to write server config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration.
func (cm *Manager) GetConfig() *Config {
	return cm.config
}

// ValidateForExecution validates configuration for flow execution.
func (cm *Manager) ValidateForExecution(config *Config) error {
	// Validate LLM configuration for execution
	if config.LLM.Provider == "openrouter" && config.LLM.APIKey == "" {
		return ErrOpenRouterAPIKeyRequired
	}

	return nil
}

// createDefaultConfig creates a default configuration structure.
func (cm *Manager) createDefaultConfig() *Config {
	return &Config{
		App:     cm.createDefaultAppConfig(),
		LLM:     cm.createDefaultLLMConfig(),
		GitHub:  cm.createDefaultGitHubConfig(),
		Flow:    cm.createDefaultFlowConfig(),
		Logging: cm.createDefaultLoggingConfig(),
	}
}

func (cm *Manager) createDefaultAppConfig() struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
} {
	return struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
		Debug   bool   `mapstructure:"debug"`
	}{
		Name:    "",
		Version: "",
		Debug:   false,
	}
}

func (cm *Manager) createDefaultLLMConfig() struct {
	Provider       string            `mapstructure:"provider"`
	APIKey         string            `mapstructure:"apiKey"`
	DefaultModel   string            `mapstructure:"defaultModel"`
	ModelOverrides map[string]string `mapstructure:"modelOverrides"`
	MaxTokens      int               `mapstructure:"maxTokens"`
	Temperature    float64           `mapstructure:"temperature"`
} {
	return struct {
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
	}
}

func (cm *Manager) createDefaultGitHubConfig() struct {
	Token         string `mapstructure:"token"`
	Owner         string `mapstructure:"owner"`
	Repository    string `mapstructure:"repository"`
	DefaultBranch string `mapstructure:"defaultBranch"`
} {
	return struct {
		Token         string `mapstructure:"token"`
		Owner         string `mapstructure:"owner"`
		Repository    string `mapstructure:"repository"`
		DefaultBranch string `mapstructure:"defaultBranch"`
	}{
		Token:         "",
		Owner:         "",
		Repository:    "",
		DefaultBranch: "",
	}
}

func (cm *Manager) createDefaultFlowConfig() struct {
	Directory      string `mapstructure:"directory"`
	DefaultTimeout string `mapstructure:"defaultTimeout"`
	CheckpointDir  string `mapstructure:"checkpointDir"`
	MaxRetries     int    `mapstructure:"maxRetries"`
	EnableParallel bool   `mapstructure:"enableParallel"`
} {
	return struct {
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
	}
}

func (cm *Manager) createDefaultLoggingConfig() struct {
	Level   string `mapstructure:"level"`
	Format  string `mapstructure:"format"`
	File    string `mapstructure:"file"`
	Console bool   `mapstructure:"console"`
} {
	return struct {
		Level   string `mapstructure:"level"`
		Format  string `mapstructure:"format"`
		File    string `mapstructure:"file"`
		Console bool   `mapstructure:"console"`
	}{
		Level:   "",
		Format:  "",
		File:    "",
		Console: false,
	}
}

// setDefaults sets default configuration values.
func (cm *Manager) setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "flow-test-go")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.debug", false)

	// LLM defaults
	viper.SetDefault("llm.provider", "openrouter")
	viper.SetDefault("llm.defaultModel", "openai/gpt-4-turbo")

	const (
		defaultMaxTokens   = 4096
		defaultTemperature = 0.7
	)

	viper.SetDefault("llm.maxTokens", defaultMaxTokens)
	viper.SetDefault("llm.temperature", defaultTemperature)

	// GitHub defaults
	viper.SetDefault("github.owner", "your-github-username")
	viper.SetDefault("github.repository", "your-github-repo")
	viper.SetDefault("github.defaultBranch", "main")

	// Flow defaults
	viper.SetDefault("flow.directory", ".flows")
	viper.SetDefault("flow.defaultTimeout", "5m")
	viper.SetDefault("flow.checkpointDir", ".flows/checkpoints")

	const defaultMaxRetries = 3

	viper.SetDefault("flow.maxRetries", defaultMaxRetries)
	viper.SetDefault("flow.enableParallel", true)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
	viper.SetDefault("logging.console", true)
}

// validateConfig validates the basic structure of the configuration.
// Currently accepts all configurations as the validation logic
// has been moved to ValidateForExecution for more specific use cases.
func (cm *Manager) validateConfig(_ *Config) error {
	// Basic validation - currently all configs are accepted
	// Specific validation happens in ValidateForExecution
	return nil
}

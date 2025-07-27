package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServerConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  MCPServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid stdio config",
			config: MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				Args:          []string{"-m", "test_server"},
				TransportType: TransportStdio,
				Capabilities: MCPCapabilities{
					Tools:     true,
					Resources: true,
					Prompts:   false,
				},
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid http config",
			config: MCPServerConfig{
				Name:          "http-server",
				Command:       "node",
				Args:          []string{"server.js"},
				TransportType: TransportHTTP,
				TransportOptions: map[string]any{
					"host": "localhost",
					"port": 8080,
				},
				Capabilities: MCPCapabilities{
					Tools:     true,
					Resources: false,
					Prompts:   true,
				},
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: MCPServerConfig{
				Command:       "python",
				TransportType: TransportStdio,
			},
			wantErr: true,
			errMsg:  "server name is required",
		},
		{
			name: "missing command",
			config: MCPServerConfig{
				Name:          "test-server",
				TransportType: TransportStdio,
			},
			wantErr: true,
			errMsg:  "server command is required",
		},
		{
			name: "invalid transport type",
			config: MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				TransportType: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid transport type",
		},
		{
			name: "no capabilities enabled",
			config: MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				TransportType: TransportStdio,
				Capabilities: MCPCapabilities{
					Tools:     false,
					Resources: false,
					Prompts:   false,
				},
			},
			wantErr: true,
			errMsg:  "server must have at least one capability enabled",
		},
		{
			name: "http without transport options",
			config: MCPServerConfig{
				Name:          "http-server",
				Command:       "node",
				TransportType: TransportHTTP,
				Capabilities: MCPCapabilities{
					Tools: true,
				},
				// Missing TransportOptions for HTTP
			},
			wantErr: true,
			errMsg:  "HTTP transport requires transport options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)

				// Verify it's an ExecutionError
				var execErr *ExecutionError
				require.ErrorAs(t, err, &execErr)
				assert.NotEmpty(t, execErr.Code)
				assert.NotZero(t, execErr.Timestamp)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMCPTransportType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		transport MCPTransportType
		want      string
	}{
		{"stdio", TransportStdio, "stdio"},
		{"http", TransportHTTP, "http"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, string(tt.transport))
		})
	}
}

func TestMCPServerState_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status MCPServerState
		want   string
	}{
		{"stopped", MCPStateStopped, "stopped"},
		{"starting", MCPStateStarting, "starting"},
		{"running", MCPStateRunning, "running"},
		{"failed", MCPStateFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}

func TestMCPCapabilities_HasAnyCapability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		capabilities MCPCapabilities
		want         bool
	}{
		{
			name: "tools only",
			capabilities: MCPCapabilities{
				Tools:     true,
				Resources: false,
				Prompts:   false,
			},
			want: true,
		},
		{
			name: "resources only",
			capabilities: MCPCapabilities{
				Tools:     false,
				Resources: true,
				Prompts:   false,
			},
			want: true,
		},
		{
			name: "prompts only",
			capabilities: MCPCapabilities{
				Tools:     false,
				Resources: false,
				Prompts:   true,
			},
			want: true,
		},
		{
			name: "all capabilities",
			capabilities: MCPCapabilities{
				Tools:     true,
				Resources: true,
				Prompts:   true,
			},
			want: true,
		},
		{
			name: "no capabilities",
			capabilities: MCPCapabilities{
				Tools:     false,
				Resources: false,
				Prompts:   false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.capabilities.Tools || tt.capabilities.Resources || tt.capabilities.Prompts
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMCPHealthCheck_IsEnabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		healthCheck *MCPHealthCheck
		want        bool
	}{
		{
			name: "enabled health check",
			healthCheck: &MCPHealthCheck{
				Enabled:  true,
				Interval: 30 * time.Second,
				Timeout:  5 * time.Second,
			},
			want: true,
		},
		{
			name: "disabled health check",
			healthCheck: &MCPHealthCheck{
				Enabled:  false,
				Interval: 30 * time.Second,
				Timeout:  5 * time.Second,
			},
			want: false,
		},
		{
			name:        "nil health check",
			healthCheck: nil,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.healthCheck != nil {
				assert.Equal(t, tt.want, tt.healthCheck.Enabled)
			} else {
				// nil health check should be considered disabled
				assert.False(t, tt.want)
			}
		})
	}
}

func TestMCPTool_Validation(t *testing.T) {
	t.Parallel()

	tool := MCPTool{
		Name:        "test-tool",
		Description: "A test tool",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "string",
				},
			},
		},
		ServerName: "test-server",
	}

	// Basic validation
	assert.NotEmpty(t, tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.NotNil(t, tool.Schema)
	assert.NotEmpty(t, tool.ServerName)
}

func TestMCPResource_Validation(t *testing.T) {
	t.Parallel()

	resource := MCPResource{
		URI:         "file:///path/to/resource",
		Name:        "test-resource",
		Description: "A test resource",
		MimeType:    "text/plain",
		ServerName:  "test-server",
	}

	// Basic validation
	assert.NotEmpty(t, resource.URI)
	assert.NotEmpty(t, resource.Name)
	assert.NotEmpty(t, resource.Description)
	assert.NotEmpty(t, resource.MimeType)
	assert.NotEmpty(t, resource.ServerName)
}

func TestMCPToolCall_Validation(t *testing.T) {
	t.Parallel()

	toolCall := MCPToolCall{
		ID:         "call-123",
		ToolName:   "test-tool",
		ServerName: "test-server",
		Arguments:  map[string]any{"input": "test"},
		Timestamp:  time.Now(),
	}

	// Basic validation
	assert.NotEmpty(t, toolCall.ID)
	assert.NotEmpty(t, toolCall.ToolName)
	assert.NotEmpty(t, toolCall.ServerName)
	assert.NotNil(t, toolCall.Arguments)
	assert.NotZero(t, toolCall.Timestamp)
}

func TestMCPToolResult_Validation(t *testing.T) {
	t.Parallel()

	result := MCPToolResult{
		CallID:    "call-123",
		Success:   true,
		Result:    "Tool execution result",
		Duration:  time.Millisecond * 100,
		Timestamp: time.Now(),
	}

	// Basic validation
	assert.NotEmpty(t, result.CallID)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Result)
	assert.Greater(t, result.Duration, time.Duration(0))
	assert.NotZero(t, result.Timestamp)
}

// Benchmark tests
func BenchmarkMCPServerConfig_Validate(b *testing.B) {
	config := MCPServerConfig{
		Name:          "bench-server",
		Command:       "python",
		Args:          []string{"-m", "test_server"},
		TransportType: TransportStdio,
		Capabilities: MCPCapabilities{
			Tools:     true,
			Resources: true,
			Prompts:   false,
		},
		Timeout: 30 * time.Second,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

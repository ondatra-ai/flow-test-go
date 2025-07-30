package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/peterovchinnikov/flow-test-go/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestMCPServerConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  types.MCPServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid stdio config",
			config: types.MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				Args:          []string{"-m", "test_server"},
				TransportType: types.TransportStdio,
				Capabilities: types.MCPCapabilities{
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
			config: types.MCPServerConfig{
				Name:          "http-server",
				Command:       "node",
				Args:          []string{"server.js"},
				TransportType: types.TransportHTTP,
				TransportOptions: map[string]any{
					"host": "localhost",
					"port": 8080,
				},
				Capabilities: types.MCPCapabilities{
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
			config: types.MCPServerConfig{
				Command:       "python",
				TransportType: types.TransportStdio,
			},
			wantErr: true,
			errMsg:  "server name is required",
		},
		{
			name: "missing command",
			config: types.MCPServerConfig{
				Name:          "test-server",
				TransportType: types.TransportStdio,
			},
			wantErr: true,
			errMsg:  "server command is required",
		},
		{
			name: "invalid transport type",
			config: types.MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				TransportType: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid transport type",
		},
		{
			name: "no capabilities enabled",
			config: types.MCPServerConfig{
				Name:          "test-server",
				Command:       "python",
				TransportType: types.TransportStdio,
				Capabilities: types.MCPCapabilities{
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
			config: types.MCPServerConfig{
				Name:          "http-server",
				Command:       "node",
				TransportType: types.TransportHTTP,
				Capabilities: types.MCPCapabilities{
					Tools: true,
				},
				// Missing TransportOptions for HTTP
			},
			wantErr: true,
			errMsg:  "HTTP transport requires transport options",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.config.Validate()

			if testCase.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errMsg)

				// Verify it's an types.ExecutionError
				var execErr *types.ExecutionError
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
		transport types.MCPTransportType
		want      string
	}{
		{"stdio", types.TransportStdio, "stdio"},
		{"http", types.TransportHTTP, "http"},
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
		status types.MCPServerState
		want   string
	}{
		{"stopped", types.MCPStateStopped, "stopped"},
		{"starting", types.MCPStateStarting, "starting"},
		{"running", types.MCPStateRunning, "running"},
		{"failed", types.MCPStateFailed, "failed"},
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
		capabilities types.MCPCapabilities
		want         bool
	}{
		{
			name: "tools only",
			capabilities: types.MCPCapabilities{
				Tools:     true,
				Resources: false,
				Prompts:   false,
			},
			want: true,
		},
		{
			name: "resources only",
			capabilities: types.MCPCapabilities{
				Tools:     false,
				Resources: true,
				Prompts:   false,
			},
			want: true,
		},
		{
			name: "prompts only",
			capabilities: types.MCPCapabilities{
				Tools:     false,
				Resources: false,
				Prompts:   true,
			},
			want: true,
		},
		{
			name: "all capabilities",
			capabilities: types.MCPCapabilities{
				Tools:     true,
				Resources: true,
				Prompts:   true,
			},
			want: true,
		},
		{
			name: "no capabilities",
			capabilities: types.MCPCapabilities{
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
		healthCheck *types.MCPHealthCheck
		want        bool
	}{
		{
			name: "enabled health check",
			healthCheck: &types.MCPHealthCheck{
				Enabled:  true,
				Interval: 30 * time.Second,
				Timeout:  5 * time.Second,
			},
			want: true,
		},
		{
			name: "disabled health check",
			healthCheck: &types.MCPHealthCheck{
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			if testCase.healthCheck != nil {
				assert.Equal(t, testCase.want, testCase.healthCheck.Enabled)
			} else {
				// nil health check should be considered disabled
				assert.False(t, testCase.want)
			}
		})
	}
}

func TestMCPTool_Validation(t *testing.T) {
	t.Parallel()

	tool := types.MCPTool{
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

	resource := types.MCPResource{
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

	toolCall := types.MCPToolCall{
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

	result := types.MCPToolResult{
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

// Benchmark tests.
func BenchmarkMCPServerConfig_Validate(b *testing.B) {
	config := types.MCPServerConfig{
		Name:          "bench-server",
		Command:       "python",
		Args:          []string{"-m", "test_server"},
		TransportType: types.TransportStdio,
		Capabilities: types.MCPCapabilities{
			Tools:     true,
			Resources: true,
			Prompts:   false,
		},
		Timeout: 30 * time.Second,
	}

	b.ResetTimer()
	for range b.N {
		_ = config.Validate()
	}
}

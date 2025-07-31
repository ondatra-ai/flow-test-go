// Package types defines the data structures and types used throughout the flow-test-go application.
package types

import (
	"time"
)

// MCPServerConfig represents the configuration for an MCP server.
type MCPServerConfig struct {
	Name             string            `json:"name"                       yaml:"name"`
	Command          string            `json:"command"                    yaml:"command"`
	Args             []string          `json:"args,omitempty"             yaml:"args,omitempty"`
	Env              map[string]string `json:"env,omitempty"              yaml:"env,omitempty"`
	TransportType    MCPTransportType  `json:"transportType"              yaml:"transportType"`
	TransportOptions map[string]any    `json:"transportOptions,omitempty" yaml:"transportOptions,omitempty"`
	Capabilities     MCPCapabilities   `json:"capabilities"               yaml:"capabilities"`
	Timeout          time.Duration     `json:"timeout,omitempty"          yaml:"timeout,omitempty"`
	HealthCheck      *MCPHealthCheck   `json:"healthCheck,omitempty"      yaml:"healthCheck,omitempty"`
	AutoRestart      bool              `json:"autoRestart,omitempty"      yaml:"autoRestart,omitempty"`
	MaxRestarts      int               `json:"maxRestarts,omitempty"      yaml:"maxRestarts,omitempty"`
	Metadata         map[string]any    `json:"metadata,omitempty"         yaml:"metadata,omitempty"`
}

// MCPTransportType defines the transport mechanism for MCP communication.
type MCPTransportType string

const (
	// TransportStdio represents the stdio transport type for MCP communication.
	TransportStdio MCPTransportType = "stdio"
	// TransportHTTP represents the HTTP transport type for MCP communication.
	TransportHTTP MCPTransportType = "http"
	// TransportTCP represents the TCP transport type for MCP communication.
	TransportTCP MCPTransportType = "tcp"
)

// MCPCapabilities defines what capabilities an MCP server provides.
type MCPCapabilities struct {
	Tools     bool `json:"tools"     yaml:"tools"`
	Resources bool `json:"resources" yaml:"resources"`
	Prompts   bool `json:"prompts"   yaml:"prompts"`
	Logging   bool `json:"logging"   yaml:"logging"`
}

// MCPHealthCheck defines health check configuration for an MCP server.
type MCPHealthCheck struct {
	Enabled  bool          `json:"enabled"           yaml:"enabled"`
	Interval time.Duration `json:"interval"          yaml:"interval"`
	Timeout  time.Duration `json:"timeout"           yaml:"timeout"`
	Command  string        `json:"command,omitempty" yaml:"command,omitempty"`
}

// MCPServerStatus represents the runtime status of an MCP server.
type MCPServerStatus struct {
	Name         string          `json:"name"`
	Status       MCPServerState  `json:"status"`
	PID          int             `json:"pid,omitempty"`
	StartTime    time.Time       `json:"startTime,omitempty"`
	LastPing     time.Time       `json:"lastPing,omitempty"`
	RestartCount int             `json:"restartCount"`
	Tools        []MCPTool       `json:"tools,omitempty"`
	Resources    []MCPResource   `json:"resources,omitempty"`
	Error        *ExecutionError `json:"error,omitempty"`
	Metadata     map[string]any  `json:"metadata,omitempty"`
}

// MCPServerState represents the state of an MCP server.
type MCPServerState string

const (
	// MCPStateStarting represents a starting MCP server state.
	MCPStateStarting MCPServerState = "starting"
	// MCPStateRunning represents a running MCP server state.
	MCPStateRunning MCPServerState = "running"
	// MCPStateStopped represents a stopped MCP server state.
	MCPStateStopped MCPServerState = "stopped"
	// MCPStateFailed represents a failed MCP server state.
	MCPStateFailed MCPServerState = "failed"
	// MCPStateRestarting represents a restarting MCP server state.
	MCPStateRestarting MCPServerState = "restarting"
)

// MCPTool represents a tool available from an MCP server.
type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
	ServerName  string         `json:"serverName"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// MCPResource represents a resource available from an MCP server.
type MCPResource struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	URI         string         `json:"uri"`
	MimeType    string         `json:"mimeType,omitempty"`
	ServerName  string         `json:"serverName"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// MCPToolCall represents a call to an MCP tool.
type MCPToolCall struct {
	ID         string         `json:"id"`
	ToolName   string         `json:"toolName"`
	ServerName string         `json:"serverName"`
	Arguments  map[string]any `json:"arguments"`
	Timestamp  time.Time      `json:"timestamp"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

// MCPToolResult represents the result of an MCP tool call.
type MCPToolResult struct {
	CallID    string          `json:"callId"`
	Success   bool            `json:"success"`
	Result    any             `json:"result,omitempty"`
	Error     *ExecutionError `json:"error,omitempty"`
	Duration  time.Duration   `json:"duration"`
	Timestamp time.Time       `json:"timestamp"`
	Metadata  map[string]any  `json:"metadata,omitempty"`
}

// MCPMessage represents a message in the MCP protocol.
type MCPMessage struct {
	ID     string         `json:"id,omitempty"`
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
	Result any            `json:"result,omitempty"`
	Error  *MCPError      `json:"error,omitempty"`
}

// MCPError represents an error in the MCP protocol.
type MCPError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

// Validate validates the MCP server configuration.
func (m *MCPServerConfig) Validate() error {
	if m.Name == "" {
		return &ExecutionError{
			Code:        "INVALID_MCP_CONFIG",
			Message:     "server name is required",
			Details:     nil,
			Recoverable: false,
			Timestamp:   time.Now(),
			StackTrace:  "",
		}
	}

	if m.Command == "" {
		return &ExecutionError{
			Code:        "INVALID_MCP_CONFIG",
			Message:     "server command is required",
			Details:     nil,
			Recoverable: false,
			Timestamp:   time.Now(),
			StackTrace:  "",
		}
	}

	if m.TransportType != TransportStdio && m.TransportType != TransportHTTP {
		return &ExecutionError{
			Code:        "INVALID_MCP_CONFIG",
			Message:     "invalid transport type",
			Details:     nil,
			Recoverable: false,
			Timestamp:   time.Now(),
			StackTrace:  "",
		}
	}

	if !m.Capabilities.Tools && !m.Capabilities.Resources && !m.Capabilities.Prompts {
		return &ExecutionError{
			Code:        "INVALID_MCP_CONFIG",
			Message:     "server must have at least one capability enabled",
			Details:     nil,
			Recoverable: false,
			Timestamp:   time.Now(),
			StackTrace:  "",
		}
	}

	if m.TransportType == TransportHTTP && len(m.TransportOptions) == 0 {
		return &ExecutionError{
			Code:        "INVALID_MCP_CONFIG",
			Message:     "HTTP transport requires transport options",
			Details:     nil,
			Recoverable: false,
			Timestamp:   time.Now(),
			StackTrace:  "",
		}
	}

	return nil
}

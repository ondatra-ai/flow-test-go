// Package types defines the data structures and types used throughout the flow-test-go application.
package types

import (
	"time"
)

// FlowDefinition represents a complete flow configuration.
type FlowDefinition struct {
	Schema      string            `json:"$schema,omitempty"     yaml:"schema,omitempty"`
	Version     string            `json:"version"               yaml:"version"`
	ID          string            `json:"id"                    yaml:"id"`
	Name        string            `json:"name"                  yaml:"name"`
	Description string            `json:"description"           yaml:"description"`
	Variables   map[string]string `json:"variables,omitempty"   yaml:"variables,omitempty"`
	Steps       map[string]Step   `json:"steps"                 yaml:"steps"`
	InitialStep string            `json:"initialStep,omitempty" yaml:"initialStep,omitempty"`
}

// Step represents a single step in a flow.
type Step struct {
	Type       StepType          `json:"type"                 yaml:"type"`
	Prompt     *PromptConfig     `json:"prompt,omitempty"     yaml:"prompt,omitempty"`
	Model      string            `json:"model,omitempty"      yaml:"model,omitempty"`
	Tools      []string          `json:"tools,omitempty"      yaml:"tools,omitempty"`
	MCPServer  string            `json:"mcpServer,omitempty"  yaml:"mcpServer,omitempty"`
	Next       string            `json:"next,omitempty"       yaml:"next,omitempty"`
	Conditions []ConditionConfig `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Timeout    *time.Duration    `json:"timeout,omitempty"    yaml:"timeout,omitempty"`
	Retry      *RetryConfig      `json:"retry,omitempty"      yaml:"retry,omitempty"`
	Metadata   map[string]any    `json:"metadata,omitempty"   yaml:"metadata,omitempty"`
}

// StepType defines the type of step.
type StepType string

const (
	// StepTypePrompt represents a prompt step type.
	StepTypePrompt StepType = "prompt"
	// StepTypeCondition represents a condition step type.
	StepTypeCondition StepType = "condition"
	// StepTypeEnd represents an end step type.
	StepTypeEnd StepType = "end"
	// StepTypeGitHub represents a GitHub step type.
	StepTypeGitHub StepType = "github"
	// StepTypeTool represents a tool step type.
	StepTypeTool StepType = "tool"
)

// PromptConfig defines the configuration for a prompt step.
type PromptConfig struct {
	Template string         `json:"template"          yaml:"template"`
	System   string         `json:"system,omitempty"  yaml:"system,omitempty"`
	Context  map[string]any `json:"context,omitempty" yaml:"context,omitempty"`
}

// ConditionConfig defines a condition for conditional step execution.
type ConditionConfig struct {
	Expression string `json:"expression" yaml:"expression"`
	Next       string `json:"next"       yaml:"next"`
}

// RetryConfig defines retry behavior.
type RetryConfig struct {
	MaxAttempts int           `json:"maxAttempts"       yaml:"maxAttempts"`
	Delay       time.Duration `json:"delay"             yaml:"delay"`
	Backoff     string        `json:"backoff,omitempty" yaml:"backoff,omitempty"`
}

// ExecutionContext represents the runtime context of a flow execution.
type ExecutionContext struct {
	FlowID      string                `json:"flowId"`
	SessionID   string                `json:"sessionId"`
	CurrentStep string                `json:"currentStep"`
	Variables   map[string]any        `json:"variables"`
	StepResults map[string]StepResult `json:"stepResults"`
	StartTime   time.Time             `json:"startTime"`
	LastUpdate  time.Time             `json:"lastUpdate"`
	Status      ExecutionStatus       `json:"status"`
	Error       *ExecutionError       `json:"error,omitempty"`
	Metadata    map[string]any        `json:"metadata,omitempty"`
}

// StepResult represents the result of a step execution.
type StepResult struct {
	StepID     string          `json:"stepId"`
	Status     StepStatus      `json:"status"`
	Output     any             `json:"output,omitempty"`
	Error      *ExecutionError `json:"error,omitempty"`
	StartTime  time.Time       `json:"startTime"`
	EndTime    time.Time       `json:"endTime"`
	Duration   time.Duration   `json:"duration"`
	TokensUsed int             `json:"tokensUsed,omitempty"`
	Cost       float64         `json:"cost,omitempty"`
	Metadata   map[string]any  `json:"metadata,omitempty"`
}

// ExecutionStatus represents the status of flow execution.
type ExecutionStatus string

const (
	// StatusPending represents a pending execution status.
	StatusPending ExecutionStatus = "pending"
	// StatusRunning represents a running execution status.
	StatusRunning ExecutionStatus = "running"
	// StatusCompleted represents a completed execution status.
	StatusCompleted ExecutionStatus = "completed"
	// StatusFailed represents a failed execution status.
	StatusFailed ExecutionStatus = "failed"
	// StatusCanceled represents a canceled execution status.
	StatusCanceled ExecutionStatus = "canceled"
	// StatusPaused represents a paused execution status.
	StatusPaused ExecutionStatus = "paused"
)

// StepStatus represents the status of step execution.
type StepStatus string

const (
	// StepStatusPending represents a pending step status.
	StepStatusPending StepStatus = "pending"
	// StepStatusRunning represents a running step status.
	StepStatusRunning StepStatus = "running"
	// StepStatusCompleted represents a completed step status.
	StepStatusCompleted StepStatus = "completed"
	// StepStatusFailed represents a failed step status.
	StepStatusFailed StepStatus = "failed"
	// StepStatusSkipped represents a skipped step status.
	StepStatusSkipped StepStatus = "skipped"
)

// ExecutionError represents an error during execution.
type ExecutionError struct {
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Details     any       `json:"details,omitempty"`
	Recoverable bool      `json:"recoverable"`
	Timestamp   time.Time `json:"timestamp"`
	StackTrace  string    `json:"stackTrace,omitempty"`
}

// Validate validates the flow definition.
func (f *FlowDefinition) Validate() error {
	if err := f.validateBasicFields(); err != nil {
		return err
	}

	// Validate step references
	for stepID, step := range f.Steps {
		if err := f.validateStep(stepID, step); err != nil {
			return err
		}
	}

	return nil
}

// validateBasicFields validates the basic flow fields.
func (f *FlowDefinition) validateBasicFields() error {
	if f.ID == "" {
		return &ExecutionError{
			Code:      "INVALID_FLOW",
			Message:   "flow ID is required",
			Timestamp: time.Now(),
		}
	}

	if f.Name == "" {
		return &ExecutionError{
			Code:      "INVALID_FLOW",
			Message:   "flow name is required",
			Timestamp: time.Now(),
		}
	}

	if len(f.Steps) == 0 {
		return &ExecutionError{
			Code:      "INVALID_FLOW",
			Message:   "flow must have at least one step",
			Timestamp: time.Now(),
		}
	}

	return nil
}

// validateStep validates a single step and its references.
func (f *FlowDefinition) validateStep(stepID string, step Step) error {
	if err := f.validateStepConfiguration(stepID, step); err != nil {
		return err
	}

	if err := f.validateStepReferences(stepID, step); err != nil {
		return err
	}

	return nil
}

// validateStepConfiguration validates step-specific configuration.
func (f *FlowDefinition) validateStepConfiguration(stepID string, step Step) error {
	if step.Type == StepTypePrompt && step.Prompt == nil {
		return &ExecutionError{
			Code:      "INVALID_STEP",
			Message:   "prompt step must have prompt configuration",
			Details:   map[string]any{"stepId": stepID},
			Timestamp: time.Now(),
		}
	}

	if step.Type == StepTypeCondition && len(step.Conditions) == 0 {
		return &ExecutionError{
			Code:      "INVALID_STEP",
			Message:   "condition step must have at least one condition",
			Details:   map[string]any{"stepId": stepID},
			Timestamp: time.Now(),
		}
	}

	return nil
}

// validateStepReferences validates step references to other steps.
func (f *FlowDefinition) validateStepReferences(stepID string, step Step) error {
	// Validate next step references
	if step.Next != "" && step.Type != StepTypeEnd {
		if _, exists := f.Steps[step.Next]; !exists {
			return &ExecutionError{
				Code:      "INVALID_REFERENCE",
				Message:   "step references non-existent next step",
				Details:   map[string]any{"stepId": stepID, "nextStep": step.Next},
				Timestamp: time.Now(),
			}
		}
	}

	// Validate condition references
	for _, condition := range step.Conditions {
		if condition.Next != "" {
			if _, exists := f.Steps[condition.Next]; !exists {
				return &ExecutionError{
					Code:      "INVALID_REFERENCE",
					Message:   "condition references non-existent step",
					Details:   map[string]any{"stepId": stepID, "conditionNext": condition.Next},
					Timestamp: time.Now(),
				}
			}
		}
	}

	return nil
}

// Error implements the error interface for ExecutionError.
func (e *ExecutionError) Error() string {
	return e.Message
}

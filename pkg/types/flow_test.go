package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowDefinition_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flow    FlowDefinition
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flow",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {
						Type: StepTypePrompt,
						Prompt: &PromptConfig{
							Template: "Hello {{.name}}",
						},
						Next: "step2",
					},
					"step2": {
						Type: StepTypeEnd,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing flow ID",
			flow: FlowDefinition{
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {Type: StepTypePrompt},
				},
			},
			wantErr: true,
			errMsg:  "flow ID is required",
		},
		{
			name: "missing flow name",
			flow: FlowDefinition{
				ID:          "test-flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {Type: StepTypePrompt},
				},
			},
			wantErr: true,
			errMsg:  "flow name is required",
		},
		{
			name: "no steps",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps:       map[string]Step{},
			},
			wantErr: true,
			errMsg:  "flow must have at least one step",
		},
		{
			name: "prompt step without prompt config",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {
						Type: StepTypePrompt,
						// Missing Prompt config
					},
				},
			},
			wantErr: true,
			errMsg:  "prompt step must have prompt configuration",
		},
		{
			name: "condition step without conditions",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {
						Type: StepTypeCondition,
						// Missing Conditions
					},
				},
			},
			wantErr: true,
			errMsg:  "condition step must have at least one condition",
		},
		{
			name: "invalid next step reference",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {
						Type: StepTypePrompt,
						Prompt: &PromptConfig{
							Template: "Hello",
						},
						Next: "nonexistent", // Invalid reference
					},
				},
			},
			wantErr: true,
			errMsg:  "step references non-existent next step",
		},
		{
			name: "invalid condition next reference",
			flow: FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]Step{
					"step1": {
						Type: StepTypeCondition,
						Conditions: []ConditionConfig{
							{
								Expression: "true",
								Next:       "nonexistent", // Invalid reference
							},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "condition references non-existent step",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.flow.Validate()

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

func TestExecutionError_Error(t *testing.T) {
	t.Parallel()

	err := &ExecutionError{
		Code:      "TEST_ERROR",
		Message:   "test error message",
		Details:   map[string]any{"key": "value"},
		Timestamp: time.Now(),
	}

	errorStr := err.Error()
	assert.Contains(t, errorStr, "test error message")
}

func TestStepStatus_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status StepStatus
		want   string
	}{
		{"pending", StepStatusPending, "pending"},
		{"running", StepStatusRunning, "running"},
		{"completed", StepStatusCompleted, "completed"},
		{"failed", StepStatusFailed, "failed"},
		{"skipped", StepStatusSkipped, "skipped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}

func TestExecutionStatus_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status ExecutionStatus
		want   string
	}{
		{"pending", StatusPending, "pending"},
		{"running", StatusRunning, "running"},
		{"completed", StatusCompleted, "completed"},
		{"failed", StatusFailed, "failed"},
		{"canceled", StatusCanceled, "canceled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}

func TestPromptConfig_Validation(t *testing.T) {
	t.Parallel()

	prompt := &PromptConfig{
		Template: "Hello {{.name}}",
		System:   "You are a helpful assistant",
		Context:  map[string]any{"key": "value"},
	}

	// Basic validation - template should not be empty
	assert.NotEmpty(t, prompt.Template)
}

func TestConditionConfig_Validation(t *testing.T) {
	t.Parallel()

	condition := ConditionConfig{
		Expression: "result.success == true",
		Next:       "success_step",
	}

	// Basic validation - expression and next should not be empty
	assert.NotEmpty(t, condition.Expression)
	assert.NotEmpty(t, condition.Next)
}

func TestRetryConfig_Validation(t *testing.T) {
	t.Parallel()

	retry := &RetryConfig{
		MaxAttempts: 3,
		Delay:       time.Second * 5,
	}

	// Basic validation
	assert.Greater(t, retry.MaxAttempts, 0)
	assert.Greater(t, retry.Delay, time.Duration(0))
}

// Benchmark tests
func BenchmarkFlowDefinition_Validate(b *testing.B) {
	flow := FlowDefinition{
		ID:          "bench-flow",
		Name:        "Benchmark Flow",
		Description: "A benchmark flow",
		Steps: map[string]Step{
			"step1": {
				Type: StepTypePrompt,
				Prompt: &PromptConfig{
					Template: "Hello {{.name}}",
				},
				Next: "step2",
			},
			"step2": {
				Type: StepTypeEnd,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flow.Validate()
	}
}

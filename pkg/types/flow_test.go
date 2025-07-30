package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peterovchinnikov/flow-test-go/pkg/types"
)

func TestFlowDefinition_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		flow    types.FlowDefinition
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flow",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
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
			},
			wantErr: false,
		},
		{
			name: "missing flow ID",
			flow: types.FlowDefinition{
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {Type: types.StepTypePrompt},
				},
			},
			wantErr: true,
			errMsg:  "flow ID is required",
		},
		{
			name: "missing flow name",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {Type: types.StepTypePrompt},
				},
			},
			wantErr: true,
			errMsg:  "flow name is required",
		},
		{
			name: "no steps",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps:       map[string]types.Step{},
			},
			wantErr: true,
			errMsg:  "flow must have at least one step",
		},
		{
			name: "prompt step without prompt config",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {
						Type: types.StepTypePrompt,
						// Missing Prompt config
					},
				},
			},
			wantErr: true,
			errMsg:  "prompt step must have prompt configuration",
		},
		{
			name: "condition step without conditions",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {
						Type: types.StepTypeCondition,
						// Missing Conditions
					},
				},
			},
			wantErr: true,
			errMsg:  "condition step must have at least one condition",
		},
		{
			name: "invalid next step reference",
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {
						Type: types.StepTypePrompt,
						Prompt: &types.PromptConfig{
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
			flow: types.FlowDefinition{
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Steps: map[string]types.Step{
					"step1": {
						Type: types.StepTypeCondition,
						Conditions: []types.ConditionConfig{
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.flow.Validate()

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

func TestExecutionError_Error(t *testing.T) {
	t.Parallel()

	err := &types.ExecutionError{
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
		status types.StepStatus
		want   string
	}{
		{"pending", types.StepStatusPending, "pending"},
		{"running", types.StepStatusRunning, "running"},
		{"completed", types.StepStatusCompleted, "completed"},
		{"failed", types.StepStatusFailed, "failed"},
		{"skipped", types.StepStatusSkipped, "skipped"},
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
		status types.ExecutionStatus
		want   string
	}{
		{"pending", types.StatusPending, "pending"},
		{"running", types.StatusRunning, "running"},
		{"completed", types.StatusCompleted, "completed"},
		{"failed", types.StatusFailed, "failed"},
		{"canceled", types.StatusCanceled, "canceled"},
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

	prompt := &types.PromptConfig{
		Template: "Hello {{.name}}",
		System:   "You are a helpful assistant",
		Context:  map[string]any{"key": "value"},
	}

	// Basic validation - template should not be empty
	assert.NotEmpty(t, prompt.Template)
}

func TestConditionConfig_Validation(t *testing.T) {
	t.Parallel()

	condition := types.ConditionConfig{
		Expression: "result.success == true",
		Next:       "success_step",
	}

	// Basic validation - expression and next should not be empty
	assert.NotEmpty(t, condition.Expression)
	assert.NotEmpty(t, condition.Next)
}

func TestRetryConfig_Validation(t *testing.T) {
	t.Parallel()

	retry := &types.RetryConfig{
		MaxAttempts: 3,
		Delay:       time.Second * 5,
	}

	// Basic validation
	assert.Positive(t, retry.MaxAttempts)
	assert.Greater(t, retry.Delay, time.Duration(0))
}

// Benchmark tests.
func BenchmarkFlowDefinition_Validate(b *testing.B) {
	flow := types.FlowDefinition{
		ID:          "bench-flow",
		Name:        "Benchmark Flow",
		Description: "A benchmark flow",
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

	b.ResetTimer()
	for range b.N {
		_ = flow.Validate()
	}
}

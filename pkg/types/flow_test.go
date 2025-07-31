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
				Schema:      "",
				Version:     "1.0",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
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
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "missing flow ID",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type:       types.StepTypePrompt,
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
			},
			wantErr: true,
			errMsg:  "flow ID is required",
		},
		{
			name: "missing flow name",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type:       types.StepTypePrompt,
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
			},
			wantErr: true,
			errMsg:  "flow name is required",
		},
		{
			name: "no steps",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps:       map[string]types.Step{},
				InitialStep: "",
			},
			wantErr: true,
			errMsg:  "flow must have at least one step",
		},
		{
			name: "prompt step without prompt config",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type:       types.StepTypePrompt,
						Prompt:     nil, // Missing Prompt config
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
			},
			wantErr: true,
			errMsg:  "prompt step must have prompt configuration",
		},
		{
			name: "condition step without conditions",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type:       types.StepTypeCondition,
						Prompt:     nil,
						Model:      "",
						Tools:      []string{},
						MCPServer:  "",
						Next:       "",
						Conditions: []types.ConditionConfig{}, // Missing Conditions
						Timeout:    nil,
						Retry:      nil,
						Metadata:   make(map[string]any),
					},
				},
				InitialStep: "",
			},
			wantErr: true,
			errMsg:  "condition step must have at least one condition",
		},
		{
			name: "invalid next step reference",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type: types.StepTypePrompt,
						Prompt: &types.PromptConfig{
							Template: "Hello",
							System:   "",
							Context:  make(map[string]any),
						},
						Model:      "",
						Tools:      []string{},
						MCPServer:  "",
						Next:       "nonexistent", // Invalid reference
						Conditions: []types.ConditionConfig{},
						Timeout:    nil,
						Retry:      nil,
						Metadata:   make(map[string]any),
					},
				},
				InitialStep: "",
			},
			wantErr: true,
			errMsg:  "step references non-existent next step",
		},
		{
			name: "invalid condition next reference",
			flow: types.FlowDefinition{
				Schema:      "",
				Version:     "",
				ID:          "test-flow",
				Name:        "Test Flow",
				Description: "A test flow",
				Variables:   make(map[string]string),
				Steps: map[string]types.Step{
					"step1": {
						Type:      types.StepTypeCondition,
						Prompt:    nil,
						Model:     "",
						Tools:     []string{},
						MCPServer: "",
						Next:      "",
						Conditions: []types.ConditionConfig{
							{
								Expression: "true",
								Next:       "nonexistent", // Invalid reference
							},
						},
						Timeout:  nil,
						Retry:    nil,
						Metadata: make(map[string]any),
					},
				},
				InitialStep: "",
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
		Code:        "TEST_ERROR",
		Message:     "test error message",
		Details:     map[string]any{"key": "value"},
		Recoverable: false,
		Timestamp:   time.Now(),
		StackTrace:  "",
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
		Backoff:     "",
	}

	// Basic validation
	assert.Positive(t, retry.MaxAttempts)
	assert.Greater(t, retry.Delay, time.Duration(0))
}

// Benchmark tests.
func BenchmarkFlowDefinition_Validate(b *testing.B) {
	flow := types.FlowDefinition{
		Schema:      "",
		Version:     "1.0",
		ID:          "bench-flow",
		Name:        "Benchmark Flow",
		Description: "A benchmark flow",
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

	b.ResetTimer()

	for range b.N {
		_ = flow.Validate()
	}
}

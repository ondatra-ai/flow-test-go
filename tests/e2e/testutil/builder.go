package testutil

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const defaultTestTimeout = 30 * time.Second

// FlowTestBuilder provides a fluent API for building flow execution tests.
type FlowTestBuilder struct {
	t           *testing.T
	flowFile    string
	configDir   string
	timeout     time.Duration
	workDir     string
	expectExit  *int
	expectError string
	expectOut   string
}

// FlowTestResult contains the results of a flow test execution.
type FlowTestResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
	Duration time.Duration
}

// NewFlowTest creates a new flow test builder.
func NewFlowTest(t *testing.T) *FlowTestBuilder {
	t.Helper()

	return &FlowTestBuilder{
		t:           t,
		flowFile:    "",
		configDir:   "",
		workDir:     "",
		timeout:     defaultTestTimeout, // Default timeout
		expectExit:  nil,
		expectError: "",
		expectOut:   "",
	}
}

// WithFlow sets the flow file to execute.
func (b *FlowTestBuilder) WithFlow(flowFile string) *FlowTestBuilder {
	b.flowFile = flowFile

	return b
}

// WithConfig sets the config directory for the test.
func (b *FlowTestBuilder) WithConfig(configDir string) *FlowTestBuilder {
	b.configDir = configDir

	return b
}

// WithTimeout sets the execution timeout.
func (b *FlowTestBuilder) WithTimeout(timeout time.Duration) *FlowTestBuilder {
	b.timeout = timeout

	return b
}

// WithWorkDir sets the working directory for the test
func (b *FlowTestBuilder) WithWorkDir(workDir string) *FlowTestBuilder {
	b.workDir = workDir

	return b
}

// ExpectExitCode sets the expected exit code
func (b *FlowTestBuilder) ExpectExitCode(code int) *FlowTestBuilder {
	b.expectExit = &code

	return b
}

// ExpectSuccess is a convenience method for expecting exit code 0
func (b *FlowTestBuilder) ExpectSuccess() *FlowTestBuilder {
	return b.ExpectExitCode(0)
}

// ExpectFailure is a convenience method for expecting non-zero exit code
func (b *FlowTestBuilder) ExpectFailure() *FlowTestBuilder {
	return b.ExpectExitCode(1)
}

// ExpectOutput sets the expected output substring
func (b *FlowTestBuilder) ExpectOutput(output string) *FlowTestBuilder {
	b.expectOut = output

	return b
}

// ExpectError sets the expected error substring
func (b *FlowTestBuilder) ExpectError(errorMsg string) *FlowTestBuilder {
	b.expectError = errorMsg
	return b
}

// Run executes the flow test and returns the result
func (b *FlowTestBuilder) Run() *FlowTestResult {
	// Validate required fields
	if b.flowFile == "" {
		b.t.Fatal("Flow file is required")
	}

	// Set default work directory if not provided
	if b.workDir == "" {
		b.workDir = b.t.TempDir()
	}

	// Create flow runner
	runner := NewFlowRunner(b.t)

	// Configure runner
	runner.SetFlowFile(b.flowFile)
	runner.SetTimeout(b.timeout)
	runner.SetWorkDir(b.workDir)

	if b.configDir != "" {
		runner.SetConfigDir(b.configDir)
	}

	// Execute the flow
	result := runner.Execute()

	// Validate expectations
	b.validateResult(result)

	return result
}

// validateResult checks the test result against expectations
func (b *FlowTestBuilder) validateResult(result *FlowTestResult) {
	// Check exit code expectation
	if b.expectExit != nil {
		if result.ExitCode != *b.expectExit {
			b.t.Errorf("Expected exit code %d, got %d", *b.expectExit, result.ExitCode)
		}
	}

	// Check output expectation
	if b.expectOut != "" {
		if !contains(result.Stdout, b.expectOut) {
			b.t.Errorf("Expected output to contain %q, got: %s", b.expectOut, result.Stdout)
		}
	}

	// Check error expectation
	if b.expectError != "" {
		if !contains(result.Stderr, b.expectError) {
			b.t.Errorf("Expected error to contain %q, got: %s", b.expectError, result.Stderr)
		}
	}
}

// contains checks if a string contains a substring (case-insensitive basic check)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// TestDataPath returns the absolute path to a test data file
func TestDataPath(relativePath string) string {
	return filepath.Join("tests", "e2e", "testdata", relativePath)
}

// FlowPath returns the absolute path to a flow file
func FlowPath(flowName string) string {
	return TestDataPath(filepath.Join("flows", flowName))
}

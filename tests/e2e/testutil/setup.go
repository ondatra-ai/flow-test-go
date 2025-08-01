package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// File permissions for directories and files.
	dirPermissions  = 0o750
	filePermissions = 0o600
)

// standardFlows contains common flow definitions used across tests.
//
//nolint:gochecknoglobals // This is a standard definition for test utilities
var standardFlows = map[string]string{
	"single-step.json": `{
		"id": "single-step",
		"name": "Single Step Flow",
		"initialStep": "step1",
		"steps": {
			"step1": {
				"type": "prompt",
				"prompt": "Hello World"
			}
		}
	}`,
	"multi-step.json": `{
		"id": "multi-step",
		"name": "Multi Step Flow",
		"initialStep": "step1",
		"steps": {
			"step1": {
				"type": "prompt",
				"prompt": "Step 1",
				"nextStep": "step2"
			},
			"step2": {
				"type": "prompt",
				"prompt": "Step 2"
			}
		}
	}`,
	"with-conditions.json": `{
		"id": "with-conditions",
		"name": "Conditional Flow",
		"initialStep": "condition1",
		"steps": {
			"condition1": {
				"type": "condition",
				"condition": "true",
				"yes": "step1",
				"no": "step2"
			},
			"step1": {
				"type": "prompt",
				"prompt": "True branch"
			},
			"step2": {
				"type": "prompt",
				"prompt": "False branch"
			}
		}
	}`,
}

// GetStandardFlows returns a copy of the standard flows map.
func GetStandardFlows() map[string]string {
	flows := make(map[string]string)
	for k, v := range standardFlows {
		flows[k] = v
	}

	return flows
}

// CreateFlowsDirectory creates a standard .flows/flows directory structure
// and returns the path to the flows directory.
func CreateFlowsDirectory(t *testing.T, tempDir string) string {
	t.Helper()

	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, dirPermissions), "Should create flows directory")

	return flowsDir
}

// WriteFlowFiles writes a map of flow files to the specified directory.
func WriteFlowFiles(t *testing.T, flowsDir string, flows map[string]string) {
	t.Helper()

	for filename, content := range flows {
		filePath := filepath.Join(flowsDir, filename)
		require.NoError(t, os.WriteFile(filePath, []byte(content), filePermissions), "Should write flow file: %s", filename)
	}
}

// WriteStandardFlows writes all standard flows to the specified directory.
func WriteStandardFlows(t *testing.T, flowsDir string) {
	t.Helper()
	WriteFlowFiles(t, flowsDir, standardFlows)
}

// SetupTestWithFlows creates a temporary directory with standard flows
// and returns the temp directory and flows directory paths.
func SetupTestWithFlows(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)
	WriteStandardFlows(t, flowsDir)

	return tempDir, flowsDir
}

// SetupTestWithCustomFlows creates a temporary directory with custom flows
// and returns the temp directory and flows directory paths.
func SetupTestWithCustomFlows(t *testing.T, flows map[string]string) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)
	WriteFlowFiles(t, flowsDir, flows)

	return tempDir, flowsDir
}

// SetupEmptyTest creates a temporary directory with just the .flows structure
// but no flow files, and returns the temp directory and flows directory paths.
func SetupEmptyTest(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)

	return tempDir, flowsDir
}

// TestExecutionWrapper provides a common pattern for test execution timing and setup.
type TestExecutionWrapper struct {
	t           *testing.T
	testName    string
	start       time.Time
	binaryCheck bool
}

// NewTestExecution creates a new test execution wrapper.
func NewTestExecution(t *testing.T, testName string) *TestExecutionWrapper {
	t.Helper()

	return &TestExecutionWrapper{
		t:           t,
		testName:    testName,
		start:       time.Time{},
		binaryCheck: true,
	}
}

// WithoutBinaryCheck disables the automatic binary existence check.
func (w *TestExecutionWrapper) WithoutBinaryCheck() *TestExecutionWrapper {
	w.binaryCheck = false

	return w
}

// Start begins the test execution timing and performs setup.
func (w *TestExecutionWrapper) Start() *TestExecutionWrapper {
	w.t.Helper()

	if w.binaryCheck {
		EnsureBinaryExists(w.t)
	}

	w.start = time.Now()

	return w
}

// Complete finishes the test execution and records coverage data.
func (w *TestExecutionWrapper) Complete(result *FlowTestResult) time.Duration {
	w.t.Helper()

	duration := time.Since(w.start)
	RecordTestExecution(w.t, w.testName, "passed", duration, true)

	return duration
}

// CreateSingleFlow creates a simple flow with the given ID and name.
func CreateSingleFlow(id, name, prompt string) string {
	return `{
		"id": "` + id + `",
		"name": "` + name + `",
		"initialStep": "step1",
		"steps": {
			"step1": {
				"type": "prompt",
				"prompt": "` + prompt + `"
			}
		}
	}`
}

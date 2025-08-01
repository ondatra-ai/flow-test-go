package e2e_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_MixedFlowTypes(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with different types of flows
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create flows with different structures inline
	testFlows := map[string]string{
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

	for filename, content := range testFlows {
		filePath := filepath.Join(flowsDir, filename)
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644), "Should write test flow file")
	}

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should handle mixed flow types")

	// Should list all flow types (check stderr for output)
	assert.Contains(t, result.Stderr, "single-step", "Should list single-step flow")
	assert.Contains(t, result.Stderr, "multi-step", "Should list multi-step flow")
	assert.Contains(t, result.Stderr, "with-conditions", "Should list conditional flow")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-mixed", "passed", duration, true)

	t.Logf("List mixed flow types test completed in %v", duration)
}

func TestListCommand_SingleFlow(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with single flow
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create just one flow inline
	flowContent := `{
		"id": "test-flow",
		"name": "Test Flow",
		"initialStep": "step1",
		"steps": {
			"step1": {
				"type": "prompt",
				"prompt": "Hello World"
			}
		}
	}`

	destPath := filepath.Join(flowsDir, "test-flow.json")
	require.NoError(t, os.WriteFile(destPath, []byte(flowContent), 0o644), "Should write test flow file")

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify successful execution
	require.Equal(t, 0, result.ExitCode, "List command should handle single flow")

	// Should list the single flow (check stderr for output)
	assert.Contains(t, result.Stderr, "test-flow", "Should list the single flow")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-single", "passed", duration, true)

	t.Logf("List single flow test completed in %v", duration)
}

func TestListCommand_MissingFlowsDirectory(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory but don't create .flows/flows
	tempDir := t.TempDir()

	start := time.Now()

	// Execute list command in directory without .flows/flows
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify execution completed
	require.Equal(t, 0, result.ExitCode, "List command should handle missing flows directory")

	// Should indicate no flows found or directory doesn't exist (check stderr for output)
	assert.Contains(t, result.Stderr, "No flows found", "Should indicate no flows found")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-missing-dir", "passed", duration, true)

	t.Logf("List missing directory test completed in %v", duration)
}

func TestListCommand_ErrorHandling(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute list command and measure performance
	result := testutil.NewFlowTest(t).
		WithTimeout(10 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify reasonable performance and error handling
	require.Equal(t, 0, result.ExitCode, "List command should handle errors gracefully")
	assert.Less(t, duration, 5*time.Second, "List command should complete quickly even with errors")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-error-handling", "passed", duration, true)

	t.Logf("List error handling test completed in %v", duration)
}

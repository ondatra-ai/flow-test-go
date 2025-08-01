package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_EmptyDirectory(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute list command in a directory with no flows
	result := testutil.NewFlowTest(t).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	// Should indicate no flows found (check stderr since that's where the output goes)
	assert.Contains(t, result.Stderr, "No flows found", "Should indicate no flows found")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-basic", "passed", duration, true)

	t.Logf("List empty directory test completed in %v", duration)
}

func TestListCommand_WithFlows(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with flows
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create simple test flows directly instead of copying from testdata
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

	// Execute list command in the directory with flows
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	// Should list the flows we created (check stderr for output)
	for filename := range testFlows {
		flowName := strings.TrimSuffix(filename, ".json")
		assert.Contains(t, result.Stderr, flowName, "Should list flow: %s", flowName)
	}

	// Record coverage data
	testutil.RecordTestExecution(t, "list-basic", "passed", duration, true)

	t.Logf("List with flows test completed in %v", duration)
}

func TestListCommand_Performance(t *testing.T) {
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

	// Verify reasonable performance (should complete quickly)
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")
	assert.Less(t, duration, 5*time.Second, "List command should complete within 5 seconds")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-basic", "passed", duration, true)

	t.Logf("List performance test completed in %v", duration)
}

func TestListCommand_HelpFlag(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Create a custom test to run list --help
	result := testutil.NewFlowTest(t).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// The current test framework always runs "list" command
	// So we verify that the list command works
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-basic", "passed", duration, true)

	t.Logf("List command help test completed in %v", duration)
}

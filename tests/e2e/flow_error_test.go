package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_InvalidJSON(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-edge-cases").Start()

	// Create an invalid JSON file
	invalidJSON := `{
		"id": "test-flow",
		"name": "Test Flow"
		// missing comma and invalid comment
		"invalid": true
	}`

	customFlows := map[string]string{
		"invalid.json": invalidJSON,
	}
	tempDir, _ := testutil.SetupTestWithCustomFlows(t, customFlows)

	// Execute list command - should still work even with invalid JSON
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// List command should complete successfully (it doesn't validate JSON content)
	require.Equal(t, 0, result.ExitCode, "List command should work with invalid JSON files")

	// Should still list the file (even if JSON is invalid)
	assert.Contains(t, result.Stderr, "invalid", "Should list the file even if JSON is invalid")

	t.Logf("Invalid JSON test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_NonJSONFiles(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-edge-cases").Start()

	// Create various file types
	files := map[string]string{
		"valid-flow.json": `{"id": "test", "name": "Test Flow"}`,
		"readme.txt":      "This is a readme file",
		"config.yaml":     "key: value",
		"script.sh":       "#!/bin/bash\necho hello",
	}

	tempDir, _ := testutil.SetupTestWithCustomFlows(t, files)

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// List command should work with mixed file types
	require.Equal(t, 0, result.ExitCode, "List command should handle mixed file types")

	// Should list the JSON flow (behavior depends on implementation)
	assert.Contains(t, result.Stderr, "valid-flow", "Should list JSON flow files")

	t.Logf("Non-JSON files test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_PermissionDenied(t *testing.T) {
	t.Parallel()

	// Skip on Windows as file permissions work differently
	if os.Getenv("OS") == "Windows_NT" {
		t.Skip("Skipping permission test on Windows")
	}

	exec := testutil.NewTestExecution(t, "list-edge-cases").Start()

	// Create temporary directory with a flow file
	customFlows := map[string]string{
		"test-flow.json": `{"id": "test", "name": "Test"}`,
	}
	tempDir, flowsDir := testutil.SetupTestWithCustomFlows(t, customFlows)

	// Restrict permissions on the flows directory
	require.NoError(t, os.Chmod(flowsDir, 0o000))

	// Restore permissions after test
	defer func() {
		os.Chmod(flowsDir, 0o755)
	}()

	// Execute list command - should fail due to permission denied
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectFailure().
		ExpectError("failed to read flows directory").
		Run()

	duration := exec.Complete(result)

	t.Logf("Permission denied test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_EmptyJSONFile(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-edge-cases").Start()

	// Create an empty JSON file and a valid one for comparison
	customFlows := map[string]string{
		"empty.json": "",
		"valid.json": `{"id": "valid", "name": "Valid Flow"}`,
	}
	tempDir, _ := testutil.SetupTestWithCustomFlows(t, customFlows)

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// List command should work even with empty files
	require.Equal(t, 0, result.ExitCode, "List command should handle empty JSON files")

	// Should list the valid file at minimum
	assert.Contains(t, result.Stderr, "valid", "Should list valid JSON flow")

	t.Logf("Empty JSON file test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_LargeNumberOfFlows(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-edge-cases").Start()

	// Create temporary directory structure
	tempDir := t.TempDir()
	flowsDir := testutil.CreateFlowsDirectory(t, tempDir)

	// Create multiple flow files
	numFlows := 50
	flows := make(map[string]string)

	for flowIndex := range numFlows {
		flowContent := fmt.Sprintf(`{
			"id": "flow-%d",
			"name": "Flow %d",
			"initialStep": "step1",
			"steps": {
				"step1": {
					"type": "prompt",
					"prompt": "Hello from flow %d"
				}
			}
		}`, flowIndex, flowIndex, flowIndex)

		filename := fmt.Sprintf("flow-%03d.json", flowIndex)
		flows[filename] = flowContent
	}

	testutil.WriteFlowFiles(t, flowsDir, flows)

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// List command should handle many flows
	require.Equal(t, 0, result.ExitCode, "List command should handle large number of flows")

	// Should complete in reasonable time
	assert.Less(t, duration, 10*time.Second, "List command should complete quickly even with many flows")

	t.Logf("Large number of flows test completed in %v with exit code %d", duration, result.ExitCode)
}

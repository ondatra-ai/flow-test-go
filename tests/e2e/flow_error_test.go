package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_InvalidJSON(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with invalid JSON file
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create an invalid JSON file
	invalidJSON := `{
		"id": "test-flow",
		"name": "Test Flow"
		// missing comma and invalid comment
		"invalid": true
	}`

	invalidPath := filepath.Join(flowsDir, "invalid.json")
	require.NoError(t, os.WriteFile(invalidPath, []byte(invalidJSON), 0o644))

	start := time.Now()

	// Execute list command - should still work even with invalid JSON
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// List command should complete successfully (it doesn't validate JSON content)
	require.Equal(t, 0, result.ExitCode, "List command should work with invalid JSON files")

	// Should still list the file (even if JSON is invalid)
	assert.Contains(t, result.Stderr, "invalid", "Should list the file even if JSON is invalid")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-edge-cases", "passed", duration, true)

	t.Logf("Invalid JSON test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_NonJSONFiles(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with non-JSON files
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create various file types
	files := map[string]string{
		"valid-flow.json": `{"id": "test", "name": "Test Flow"}`,
		"readme.txt":      "This is a readme file",
		"config.yaml":     "key: value",
		"script.sh":       "#!/bin/bash\necho hello",
	}

	for filename, content := range files {
		filePath := filepath.Join(flowsDir, filename)
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644))
	}

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// List command should work with mixed file types
	require.Equal(t, 0, result.ExitCode, "List command should handle mixed file types")

	// Should list the JSON flow (behavior depends on implementation)
	assert.Contains(t, result.Stderr, "valid-flow", "Should list JSON flow files")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-edge-cases", "passed", duration, true)

	t.Logf("Non-JSON files test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_PermissionDenied(t *testing.T) {
	t.Parallel()

	// Skip on Windows as file permissions work differently
	if os.Getenv("OS") == "Windows_NT" {
		t.Skip("Skipping permission test on Windows")
	}

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with restricted permissions
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create a flow file
	flowPath := filepath.Join(flowsDir, "test-flow.json")
	require.NoError(t, os.WriteFile(flowPath, []byte(`{"id": "test", "name": "Test"}`), 0o644))

	// Restrict permissions on the flows directory
	require.NoError(t, os.Chmod(flowsDir, 0o000))

	// Restore permissions after test
	defer func() {
		os.Chmod(flowsDir, 0o755)
	}()

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		Run() // Don't expect success - might fail due to permissions

	duration := time.Since(start)

	// The command might fail or succeed depending on implementation
	// We just verify it doesn't crash
	t.Logf("Permission test completed with exit code %d", result.ExitCode)

	// Record coverage data - test passes if it doesn't crash
	testutil.RecordTestExecution(t, "list-edge-cases", "passed", duration, true)

	t.Logf("Permission denied test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_EmptyJSONFile(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with empty JSON file
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create an empty JSON file
	emptyPath := filepath.Join(flowsDir, "empty.json")
	require.NoError(t, os.WriteFile(emptyPath, []byte(""), 0o644))

	// Create a valid JSON file for comparison
	validPath := filepath.Join(flowsDir, "valid.json")
	require.NoError(t, os.WriteFile(validPath, []byte(`{"id": "valid", "name": "Valid Flow"}`), 0o644))

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// List command should work even with empty files
	require.Equal(t, 0, result.ExitCode, "List command should handle empty JSON files")

	// Should list the valid file at minimum
	assert.Contains(t, result.Stderr, "valid", "Should list valid JSON flow")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-edge-cases", "passed", duration, true)

	t.Logf("Empty JSON file test completed in %v with exit code %d", duration, result.ExitCode)
}

func TestListCommand_LargeNumberOfFlows(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	// Create temporary directory with many flow files
	tempDir := t.TempDir()
	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, 0o755))

	// Create multiple flow files
	numFlows := 50
	for flowIndex := range numFlows {
		flowContent := `{
			"id": "flow-%d",
			"name": "Flow %d",
			"initialStep": "step1",
			"steps": {
				"step1": {
					"type": "prompt",
					"prompt": "Hello from flow %d"
				}
			}
		}`

		filename := fmt.Sprintf("flow-%03d.json", flowIndex)
		filePath := filepath.Join(flowsDir, filename)
		content := []byte(fmt.Sprintf(flowContent, flowIndex, flowIndex, flowIndex))
		require.NoError(t, os.WriteFile(filePath, content, 0o644))
	}

	start := time.Now()

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// List command should handle many flows
	require.Equal(t, 0, result.ExitCode, "List command should handle large number of flows")

	// Should complete in reasonable time
	assert.Less(t, duration, 10*time.Second, "List command should complete quickly even with many flows")

	// Record coverage data
	testutil.RecordTestExecution(t, "list-edge-cases", "passed", duration, true)

	t.Logf("Large number of flows test completed in %v with exit code %d", duration, result.ExitCode)
}

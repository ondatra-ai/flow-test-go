package e2e_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_MixedFlowTypes(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-mixed").Start()

	// Create temporary directory with different types of flows
	tempDir, _ := testutil.SetupTestWithFlows(t)

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should handle mixed flow types")

	// Should list all flow types (check stderr for output)
	assert.Contains(t, result.Stderr, "single-step", "Should list single-step flow")
	assert.Contains(t, result.Stderr, "multi-step", "Should list multi-step flow")
	assert.Contains(t, result.Stderr, "with-conditions", "Should list conditional flow")

	t.Logf("List mixed flow types test completed in %v", duration)
}

func TestListCommand_SingleFlow(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-single").Start()

	// Create temporary directory with single custom flow
	customFlows := map[string]string{
		"test-flow.json": testutil.CreateSingleFlow("test-flow", "Test Flow", "Hello World"),
	}
	tempDir, _ := testutil.SetupTestWithCustomFlows(t, customFlows)

	// Execute list command
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify successful execution
	require.Equal(t, 0, result.ExitCode, "List command should handle single flow")

	// Should list the single flow (check stderr for output)
	assert.Contains(t, result.Stderr, "test-flow", "Should list the single flow")

	t.Logf("List single flow test completed in %v", duration)
}

func TestListCommand_MissingFlowsDirectory(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-missing-dir").Start()

	// Create temporary directory but don't create .flows/flows
	tempDir := t.TempDir()

	// Execute list command in directory without .flows/flows
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify execution completed
	require.Equal(t, 0, result.ExitCode, "List command should handle missing flows directory")

	// Should indicate no flows found or directory doesn't exist (check stderr for output)
	assert.Contains(t, result.Stderr, "No flows found", "Should indicate no flows found")

	t.Logf("List missing directory test completed in %v", duration)
}

func TestListCommand_ErrorHandling(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-error-handling").Start()

	// Execute list command and measure performance
	result := testutil.NewFlowTest(t).
		WithTimeout(10 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify reasonable performance and error handling
	require.Equal(t, 0, result.ExitCode, "List command should handle errors gracefully")
	assert.Less(t, duration, 5*time.Second, "List command should complete quickly even with errors")

	t.Logf("List error handling test completed in %v", duration)
}

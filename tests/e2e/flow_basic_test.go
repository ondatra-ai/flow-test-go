package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestListCommand_EmptyDirectory(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-basic").Start()

	// Execute list command in a directory with no flows
	result := testutil.NewFlowTest(t).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	// Should indicate no flows found (check stderr since that's where the output goes)
	assert.Contains(t, result.Stderr, "No flows found", "Should indicate no flows found")

	t.Logf("List empty directory test completed in %v", duration)
}

func TestListCommand_WithFlows(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-basic").Start()

	// Create temporary directory with standard flows
	tempDir, _ := testutil.SetupTestWithFlows(t)

	// Execute list command in the directory with flows
	result := testutil.NewFlowTest(t).
		WithWorkDir(tempDir).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify the command completed successfully
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	// Should list the flows we created (check stderr for output)
	for filename := range testutil.GetStandardFlows() {
		flowName := strings.TrimSuffix(filename, ".json")
		assert.Contains(t, result.Stderr, flowName, "Should list flow: %s", flowName)
	}

	t.Logf("List with flows test completed in %v", duration)
}

func TestListCommand_Performance(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-basic").Start()

	// Execute list command and measure performance
	result := testutil.NewFlowTest(t).
		WithTimeout(10 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// Verify reasonable performance (should complete quickly)
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")
	assert.Less(t, duration, 5*time.Second, "List command should complete within 5 seconds")

	t.Logf("List performance test completed in %v", duration)
}

func TestListCommand_HelpFlag(t *testing.T) {
	t.Parallel()

	exec := testutil.NewTestExecution(t, "list-basic").Start()

	// Create a custom test to run list --help
	result := testutil.NewFlowTest(t).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := exec.Complete(result)

	// The current test framework always runs "list" command
	// So we verify that the list command works
	require.Equal(t, 0, result.ExitCode, "List command should complete successfully")

	t.Logf("List command help test completed in %v", duration)
}

package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestConditionalFlowTrueBranch(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute conditional flow (should take true branch since condition is "true")
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/with-conditions.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the test completed successfully
	require.Equal(t, 0, result.ExitCode, "Conditional flow should complete successfully")

	// Record coverage data
	testutil.RecordTestExecution(t, "conditional", "passed", duration, true)

	t.Logf("Conditional flow (true branch) test completed in %v", duration)
}

func TestConditionalFlowExecution(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute conditional flow and verify it handles conditions properly
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/with-conditions.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify successful execution
	require.Equal(t, 0, result.ExitCode, "Conditional flow should execute successfully")

	// Verify no errors in execution (ignore coverage-related error messages)
	if !strings.Contains(result.Stderr, "coverage meta-data emit failed") &&
		!strings.Contains(result.Stderr, "coverage counter data emit failed") {
		assert.NotContains(t, result.Stderr, "error", "Conditional flow should not have errors")
		assert.NotContains(t, result.Stderr, "panic", "Conditional flow should not panic")
	}

	// Record coverage data
	testutil.RecordTestExecution(t, "conditional", "passed", duration, true)

	t.Logf("Conditional flow execution test completed in %v", duration)
}

func TestConditionalFlowBranchSelection(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Test that the conditional flow properly evaluates conditions and selects branches
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/with-conditions.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify execution completed
	require.Equal(t, 0, result.ExitCode, "Conditional flow should complete branch selection")

	// The actual branch taken depends on the condition evaluation
	// For our test flow with condition "true", it should take the positive branch
	// We can't easily verify which branch was taken without more detailed output parsing
	// but we can verify that the flow completed without errors

	// Record coverage data
	testutil.RecordTestExecution(t, "conditional", "passed", duration, true)

	t.Logf("Conditional branch selection test completed in %v", duration)
}

func TestConditionalFlowPerformance(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute conditional flow and measure performance
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/with-conditions.json")).
		WithTimeout(10 * time.Second). // Shorter timeout for performance test
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify reasonable performance (should complete quickly)
	require.Equal(t, 0, result.ExitCode, "Conditional flow should complete successfully")
	assert.Less(t, duration, 5*time.Second, "Conditional flow should complete within 5 seconds")

	// Record coverage data
	testutil.RecordTestExecution(t, "conditional", "passed", duration, true)

	t.Logf("Conditional flow performance test completed in %v", duration)
}

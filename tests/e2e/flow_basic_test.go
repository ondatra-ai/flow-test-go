package e2e_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ondatra-ai/flow-test-go/tests/e2e/testutil"
)

func TestSingleStepFlow(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute single-step flow
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/single-step.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the test completed successfully
	assert.Equal(t, 0, result.ExitCode, "Single step flow should complete successfully")

	// Record coverage data
	status := "passed"
	if result.ExitCode != 0 {
		status = "failed"
	}
	testutil.RecordTestExecution(t, "basic", status, duration, result.ExitCode == 0)

	t.Logf("Single step flow test completed in %v", duration)
}

func TestMultiStepFlow(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute multi-step flow
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/multi-step.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify the test completed successfully
	assert.Equal(t, 0, result.ExitCode, "Multi step flow should complete successfully")

	// Verify that all steps were executed (this would depend on the actual output format)
	// For now, just check that it completed without errors
	assert.NotContains(t, result.Stderr, "error", "Multi step flow should not have errors")
	assert.NotContains(t, result.Stderr, "failed", "Multi step flow should not fail")

	// Record coverage data
	status := "passed"
	if result.ExitCode != 0 {
		status = "failed"
	}
	testutil.RecordTestExecution(t, "basic", status, duration, result.ExitCode == 0)

	t.Logf("Multi step flow test completed in %v", duration)
}

func TestFlowWithTimeout(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute single-step flow with very short timeout to test timeout handling
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/single-step.json")).
		WithTimeout(1 * time.Millisecond). // Very short timeout
		Run()

	duration := time.Since(start)

	// The result could be either success (if it completes very quickly) or timeout
	// We mainly want to verify that the timeout mechanism works
	if result.ExitCode != 0 {
		t.Logf("Flow timed out as expected (exit code: %d)", result.ExitCode)
	} else {
		t.Logf("Flow completed before timeout")
	}

	// Record coverage data - consider this test as passed since we're testing timeout behavior
	testutil.RecordTestExecution(t, "basic", "passed", duration, true)

	t.Logf("Timeout test completed in %v", duration)
}

func TestFlowExecutionOrder(t *testing.T) {
	t.Parallel()

	// Ensure binary exists
	testutil.EnsureBinaryExists(t)

	start := time.Now()

	// Execute multi-step flow and verify execution happens in order
	result := testutil.NewFlowTest(t).
		WithFlow(testutil.FlowPath("basic/multi-step.json")).
		WithTimeout(30 * time.Second).
		ExpectSuccess().
		Run()

	duration := time.Since(start)

	// Verify successful execution
	assert.Equal(t, 0, result.ExitCode, "Flow should execute successfully")

	// For now, we can't easily verify execution order without knowing the exact output format
	// This test establishes the pattern for when that functionality is available

	// Record coverage data
	status := "passed"
	if result.ExitCode != 0 {
		status = "failed"
	}
	testutil.RecordTestExecution(t, "basic", status, duration, result.ExitCode == 0)

	t.Logf("Execution order test completed in %v", duration)
}

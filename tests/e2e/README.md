# End-to-End Tests

This directory contains end-to-end tests for the flow-test-go application that test flow execution through subprocess execution with coverage collection.

## Overview

The e2e tests focus on testing the **existing flow execution functionality** by:
- Running the application as a subprocess with various flow configurations
- Testing basic flow execution, conditional branching, and error handling
- Collecting coverage data from subprocess execution
- Providing isolated test environments

## Directory Structure

```
tests/e2e/
├── testdata/
│   ├── flows/
│   │   ├── basic/                  # Valid test flows
│   │   │   ├── single-step.json    # Single prompt flow
│   │   │   ├── multi-step.json     # Sequential steps
│   │   │   └── with-conditions.json # Conditional branching
│   │   └── error-cases/            # Invalid flows for error testing
│   │       ├── invalid-json.json   # Malformed JSON
│   │       ├── missing-initial.json # Missing initialStep
│   │       ├── invalid-step-ref.json # Invalid step references
│   │       └── circular-ref.json   # Circular references
│   └── expected/
│       └── outputs/                # Expected output patterns (future use)
├── testutil/                       # Test utilities package
│   ├── builder.go                  # FlowTestBuilder - fluent API
│   ├── runner.go                   # FlowRunner - subprocess execution
│   └── coverage.go                 # Coverage collection and aggregation
├── flow_basic_test.go              # Basic flow execution tests
├── flow_conditional_test.go        # Conditional flow tests
├── flow_error_test.go              # Error handling tests
└── README.md                       # This file
```

## Running E2E Tests

### Prerequisites

- Go 1.20+ (required for GOCOVERDIR support)
- Make

### Commands

```bash
# Run all e2e tests with coverage
make test-e2e-coverage

# Run only e2e tests (without coverage report)
make test-e2e

# Build coverage-instrumented binary manually
make build-e2e-coverage

# Generate coverage report from existing data
make coverage-e2e-report
```

### Individual Test Categories

```bash
# Run only basic flow tests
go test -v ./tests/e2e/ -run TestSingleStep

# Run only conditional flow tests
go test -v ./tests/e2e/ -run TestConditional

# Run only error handling tests
go test -v ./tests/e2e/ -run TestInvalidJSON
```

## Test Framework

### FlowTestBuilder API

The test framework provides a fluent API for constructing flow tests:

```go
func TestMyFlow(t *testing.T) {
    result := testutil.NewFlowTest(t).
        WithFlow("testdata/flows/basic/single-step.json").
        WithTimeout(30 * time.Second).
        ExpectSuccess().
        Run()

    assert.Equal(t, 0, result.ExitCode)
}
```

### Available Methods

- `WithFlow(flowFile)` - Set the flow file to execute
- `WithConfig(configDir)` - Set config directory (optional)
- `WithTimeout(duration)` - Set execution timeout
- `WithWorkDir(workDir)` - Set working directory
- `ExpectSuccess()` - Expect exit code 0
- `ExpectFailure()` - Expect non-zero exit code
- `ExpectExitCode(code)` - Expect specific exit code
- `ExpectOutput(substring)` - Expect output to contain substring
- `ExpectError(substring)` - Expect error to contain substring

### Test Result Structure

```go
type FlowTestResult struct {
    ExitCode int           // Process exit code
    Stdout   string        // Standard output
    Stderr   string        // Standard error
    Error    error         // Go error (if any)
    Duration time.Duration // Execution time
}
```

## Writing New Tests

### 1. Create Test Flow Files

Add new JSON flow files to the appropriate subdirectory:

- `testdata/flows/basic/` - Valid flows for positive testing
- `testdata/flows/error-cases/` - Invalid flows for error testing

### 2. Create Test Functions

```go
func TestMyNewFlow(t *testing.T) {
    t.Parallel() // Enable parallel execution

    // Ensure binary exists
    testutil.EnsureBinaryExists(t)

    start := time.Now()

    // Execute flow test
    result := testutil.NewFlowTest(t).
        WithFlow(testutil.FlowPath("basic/my-new-flow.json")).
        ExpectSuccess().
        Run()

    duration := time.Since(start)

    // Add assertions
    assert.Equal(t, 0, result.ExitCode)

    // Record coverage data
    status := "passed"
    if result.ExitCode != 0 {
        status = "failed"
    }
    testutil.RecordTestExecution(t, "category", status, duration, result.ExitCode == 0)
}
```

### 3. Test Categories

Organize tests by category:
- **basic** - Basic flow execution tests
- **conditional** - Conditional flow and branching tests
- **error-handling** - Error scenario tests

## Coverage Collection

### How It Works

1. Tests build the application with `-cover` flag
2. Each test runs in isolated coverage directory
3. `GOCOVERDIR` environment variable collects coverage data
4. Coverage data is aggregated using `go tool covdata`
5. Reports are generated in HTML and text formats

### Coverage Files

- `coverage/e2e/` - Individual test coverage data
- `coverage/e2e-merged/` - Aggregated coverage data
- `coverage/e2e.out` - Text coverage profile
- `coverage/e2e.html` - HTML coverage report
- `coverage/e2e-summary.txt` - Coverage summary
- `coverage/manifest.json` - Test execution metadata

### Debugging Coverage Issues

```bash
# Check if coverage data exists
ls -la coverage/e2e/

# Manually merge coverage data
go tool covdata merge -i=coverage/e2e/* -o=coverage/e2e-merged

# Check coverage profile
go tool cover -func=coverage/e2e.out
```

## Test Isolation

Each test runs in complete isolation:

- **Temporary directories**: Each test gets its own `t.TempDir()`
- **Coverage separation**: Individual coverage directories per test
- **Parallel execution**: Tests can run in parallel safely
- **Process isolation**: Each test spawns its own subprocess

## Debugging Failed Tests

### Common Issues

1. **Binary not found**: Run `make build-e2e-coverage` first
2. **Timeout errors**: Increase timeout with `WithTimeout()`
3. **Coverage issues**: Check Go version (1.20+ required)
4. **Flow file errors**: Validate JSON syntax

### Debugging Tips

```bash
# Run single test with verbose output
go test -v ./tests/e2e/ -run TestSpecificTest

# Check test logs
go test -v ./tests/e2e/ 2>&1 | grep "test completed"

# Manual binary execution
./bin/flow-test-go-e2e run --flow tests/e2e/testdata/flows/basic/single-step.json
```

## Future Enhancements

- Golden file testing for output validation
- Performance benchmarking
- Integration with CI/CD artifacts
- Test data generation utilities
- Advanced flow scenarios

## Troubleshooting

### Binary Build Issues
```bash
# Check Go version
go version  # Should be 1.20+

# Clean and rebuild
make clean
make build-e2e-coverage
```

### Coverage Collection Issues
```bash
# Verify GOCOVERDIR support
go help test | grep -i cover

# Check coverage directory permissions
ls -la coverage/e2e/
```

### Test Execution Issues
```bash
# Run tests with race detection
go test -race ./tests/e2e/...

# Run tests without parallel execution
go test -parallel 1 ./tests/e2e/...
```

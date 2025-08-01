# E2E Testing Design - REVISED

## Test Scope Clarification
The e2e tests will focus on testing the **flow execution engine** by:
- Creating various test flow JSON files
- Running the application with different flows
- Testing different execution scenarios and flags
- Validating flow execution results

## Revised Test Data Organization

### Directory Structure
```
tests/e2e/
├── testdata/
│   ├── flows/
│   │   ├── basic/
│   │   │   ├── hello-world.json      # Simple single-step flow
│   │   │   ├── multi-step.json       # Multiple sequential steps
│   │   │   └── with-conditions.json  # Conditional branching
│   │   ├── advanced/
│   │   │   ├── complex-branching.json # Complex condition logic
│   │   │   ├── with-tools.json       # Flow using MCP tools
│   │   │   └── circular-flow.json    # Flow with loops
│   │   └── error-cases/
│   │       ├── invalid-syntax.json   # Malformed JSON
│   │       ├── missing-steps.json    # Missing required fields
│   │       └── invalid-refs.json     # Invalid step references
│   ├── servers/
│   │   ├── mock-server.json          # Mock MCP server config
│   │   └── test-server.json          # Test server with tools
│   └── expected/
│       ├── hello-world.out           # Expected output files
│       └── multi-step.out
├── testutil/
│   ├── flow_runner.go                # Flow execution helper
│   ├── test_builder.go               # Test setup builder
│   └── assertions.go                 # Custom assertions
├── flow_execution_test.go            # Main flow execution tests
├── flow_validation_test.go           # Flow validation tests
├── flow_error_test.go                # Error handling tests
└── README.md
```

## Test Flow Examples

### Basic Hello World Flow
```json
{
  "id": "hello-world",
  "name": "Hello World Test Flow",
  "description": "Simple flow for testing basic execution",
  "initialStep": "greet",
  "steps": {
    "greet": {
      "type": "prompt",
      "prompt": "Say hello to the world",
      "tools": [],
      "mcpServer": "mock-server",
      "nextStep": "end"
    },
    "end": {
      "type": "prompt",
      "prompt": "Goodbye!",
      "tools": [],
      "mcpServer": "mock-server",
      "nextStep": null
    }
  }
}
```

### Flow with Conditions
```json
{
  "id": "conditional-flow",
  "name": "Conditional Flow Test",
  "description": "Test flow with branching logic",
  "initialStep": "check-condition",
  "steps": {
    "check-condition": {
      "type": "condition",
      "condition": "response.includes('yes')",
      "yes": "positive-path",
      "no": "negative-path"
    },
    "positive-path": {
      "type": "prompt",
      "prompt": "You chose yes!",
      "tools": [],
      "mcpServer": "mock-server",
      "nextStep": null
    },
    "negative-path": {
      "type": "prompt",
      "prompt": "You chose no!",
      "tools": [],
      "mcpServer": "mock-server",
      "nextStep": null
    }
  }
}
```

## Test Scenarios

### 1. Basic Flow Execution Tests
- Execute simple linear flows
- Verify output matches expected results
- Test with different verbosity levels
- Validate flow completion status

### 2. Conditional Flow Tests
- Test branching logic
- Verify correct path selection
- Test nested conditions
- Edge cases in condition evaluation

### 3. Tool Integration Tests
- Flows using MCP server tools
- Tool discovery and execution
- Error handling when tools unavailable

### 4. Error Handling Tests
- Invalid flow files
- Missing MCP servers
- Circular references
- Timeout scenarios

### 5. Performance Tests
- Large flows with many steps
- Concurrent flow execution
- Memory usage validation

## Test Execution Patterns

### Running a Flow Test
```go
func TestBasicFlowExecution(t *testing.T) {
    result := NewFlowTest(t).
        WithFlow("testdata/flows/basic/hello-world.json").
        WithConfig("testdata/.flows").
        WithTimeout(30 * time.Second).
        ExpectExitCode(0).
        ExpectOutput("Hello World").
        Run()

    assert.Contains(t, result.Output, "Flow completed successfully")
}
```

### Testing Error Cases
```go
func TestInvalidFlowHandling(t *testing.T) {
    result := NewFlowTest(t).
        WithFlow("testdata/flows/error-cases/invalid-syntax.json").
        ExpectExitCode(1).
        ExpectError("Failed to parse flow").
        Run()
}
```

## Command Line Patterns

The application will be tested with various command patterns:

```bash
# Basic flow execution
./flow-test-go run --flow testdata/flows/basic/hello-world.json

# With custom config directory
./flow-test-go run --flow test.json --config testdata/.flows

# With verbose output
./flow-test-go run --flow test.json -v

# With specific MCP server
./flow-test-go run --flow test.json --server mock-server

# Dry run mode
./flow-test-go run --flow test.json --dry-run
```

## Test Data Management

### Flow Categories
1. **Basic Flows**: Simple linear execution paths
2. **Conditional Flows**: Testing branching logic
3. **Tool Flows**: Integration with MCP servers
4. **Error Flows**: Invalid configurations
5. **Performance Flows**: Large/complex scenarios

### Expected Output Files
- Store expected outputs in `testdata/expected/`
- Use golden file testing pattern
- Support updating expected outputs with flag

### Mock MCP Server
- Create minimal mock server for testing
- Responds to tool discovery requests
- Provides predictable tool responses
- Supports error injection for testing

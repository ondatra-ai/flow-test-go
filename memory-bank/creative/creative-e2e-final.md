# E2E Testing Design - FINAL REVISION

## Test Scope - Existing Functionality Only
The e2e tests will focus on testing the **existing flow execution functionality**:
- Basic flow execution with prompts
- Conditional flow branching
- Flow validation and error handling
- NO MCP server/tool testing (not part of this task)

## Simplified Test Data Organization

### Directory Structure
```
tests/e2e/
├── testdata/
│   ├── flows/
│   │   ├── basic/
│   │   │   ├── single-step.json      # Single prompt flow
│   │   │   ├── multi-step.json       # Multiple sequential prompts
│   │   │   └── with-conditions.json  # Conditional branching
│   │   └── error-cases/
│   │       ├── invalid-json.json     # Malformed JSON
│   │       ├── missing-initial.json  # Missing initialStep
│   │       ├── invalid-step-ref.json # Reference to non-existent step
│   │       └── circular-ref.json     # Circular step references
│   └── expected/
│       └── outputs/                  # Expected output patterns
├── testutil/
│   ├── builder.go                    # Test execution builder
│   ├── runner.go                     # Flow runner helper
│   └── coverage.go                   # Coverage collection
├── flow_basic_test.go                # Basic flow execution tests
├── flow_conditional_test.go          # Conditional flow tests
├── flow_error_test.go                # Error handling tests
└── README.md
```

## Simplified Test Flow Examples

### Single Step Flow
```json
{
  "id": "single-step",
  "name": "Single Step Test",
  "description": "Minimal flow with one step",
  "initialStep": "only-step",
  "steps": {
    "only-step": {
      "type": "prompt",
      "prompt": "This is a test prompt",
      "nextStep": null
    }
  }
}
```

### Multi-Step Flow
```json
{
  "id": "multi-step",
  "name": "Multi Step Test",
  "description": "Flow with sequential steps",
  "initialStep": "step1",
  "steps": {
    "step1": {
      "type": "prompt",
      "prompt": "First step",
      "nextStep": "step2"
    },
    "step2": {
      "type": "prompt",
      "prompt": "Second step",
      "nextStep": "step3"
    },
    "step3": {
      "type": "prompt",
      "prompt": "Final step",
      "nextStep": null
    }
  }
}
```

### Conditional Flow
```json
{
  "id": "conditional",
  "name": "Conditional Test",
  "description": "Flow with branching",
  "initialStep": "check",
  "steps": {
    "check": {
      "type": "condition",
      "condition": "true",
      "yes": "true-branch",
      "no": "false-branch"
    },
    "true-branch": {
      "type": "prompt",
      "prompt": "Condition was true",
      "nextStep": null
    },
    "false-branch": {
      "type": "prompt",
      "prompt": "Condition was false",
      "nextStep": null
    }
  }
}
```

## Test Cases - Existing Functionality Only

### 1. Basic Flow Execution
- Single step flow execution
- Multi-step sequential flow
- Flow completion verification
- Exit code validation

### 2. Conditional Flow Tests
- Basic true/false conditions
- Nested conditions
- Complex condition expressions
- Condition evaluation errors

### 3. Flow Validation Tests
- Valid flow structure
- Missing required fields
- Invalid step references
- Circular references detection

### 4. Error Handling Tests
- Invalid JSON syntax
- Missing flow file
- Permission errors
- Malformed flow structure

## Implementation Priority

### Phase 1: Core Infrastructure
1. Remove placeholder_test.go
2. Create test builder/runner utilities
3. Set up coverage collection
4. Create basic test flows

### Phase 2: Basic Tests
1. Single step flow tests
2. Multi-step flow tests
3. Basic validation tests
4. Simple error cases

### Phase 3: Advanced Tests
1. Conditional flow tests
2. Complex flow structures
3. Edge cases and error scenarios
4. Performance tests (if time permits)

## Command Patterns (Hypothetical)
Since the actual CLI isn't implemented yet, tests will directly call the flow execution logic:

```go
// Direct execution for testing
result := ExecuteFlow("testdata/flows/basic/single-step.json")

// Or via test builder
result := NewFlowTest(t).
    WithFlowFile("testdata/flows/basic/multi-step.json").
    ExpectSuccess().
    Run()
```

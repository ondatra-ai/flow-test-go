# Active Context

## Current Task
- **E2E Subprocess Tests Implementation** (Level 3 - Intermediate Feature)
- **GitHub Issue**: #5 - https://github.com/ondatra-ai/flow-test-go/issues/5
- **Branch**: task-20250109-e2e-subprocess-tests
- **Status**: Creative phase COMPLETE - Final scope defined

## Task Overview - FINAL SCOPE
Implementing e2e tests for **existing flow execution functionality only**:
- Test basic flow execution (single-step, multi-step)
- Test conditional flow branching
- Test flow validation and error handling
- Coverage collection from subprocess execution
- **NO MCP server/tool testing** (out of scope)

## What We're Testing
1. **Flow Execution Engine**
   - Prompt step execution
   - Conditional step evaluation
   - Step navigation (nextStep)
   - Flow completion

2. **Flow Validation**
   - JSON schema validation
   - Required fields checking
   - Step reference validation
   - Circular reference detection

3. **Error Handling**
   - Invalid JSON files
   - Missing flow files
   - Malformed flow structure
   - Invalid step references

## Test Data Structure
```
tests/e2e/testdata/flows/
├── basic/
│   ├── single-step.json
│   ├── multi-step.json
│   └── with-conditions.json
└── error-cases/
    ├── invalid-json.json
    ├── missing-initial.json
    └── circular-ref.json
```

## Design Decisions (Final)
1. **Test Builder Pattern** - For readable test construction
2. **Direct Flow Execution** - Test internal APIs directly
3. **Golden File Testing** - For output validation
4. **Hierarchical Coverage** - Organized by test type
5. **No External Dependencies** - No MCP servers needed

## Next Step
Ready to implement e2e tests for existing flow functionality only.

Type 'IMPLEMENT' to begin the implementation phase.

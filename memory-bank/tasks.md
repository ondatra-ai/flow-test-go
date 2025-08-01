# Task: Implement End-to-End Tests with Subprocess Execution and Coverage Collection

## Description
Need to implement comprehensive end-to-end (e2e) tests that run the built flow-test-go application as a subprocess with proper test isolation and coverage collection.

## GitHub Issue
Issue #5: https://github.com/ondatra-ai/flow-test-go/issues/5

## Complexity
Level: 3
Type: Intermediate Feature

## Technology Stack
- Language: Go 1.21+
- Testing Framework: Go standard testing package + testify
- Coverage Tool: Go 1.20+ coverage with GOCOVERDIR support
- Build Tool: Make
- CI/CD: GitHub Actions

## Technology Validation Checkpoints

### Technology Validation Results (POC Completed)
- ✅ Go version 1.24.5 supports GOCOVERDIR
- ✅ Built binary with -cover flag successfully
- ✅ GOCOVERDIR collected coverage data (covcounters and covmeta files)
- ✅ Converted coverage data to standard format with go tool covdata
- ✅ Generated coverage report showing 50% coverage for POC
- [x] Go 1.20+ available (required for GOCOVERDIR coverage mode)
- [x] Build with -cover flag produces instrumented binary
- [x] GOCOVERDIR environment variable collects coverage data
- [x] go tool covdata can merge coverage files
- [x] Subprocess execution with os/exec package works

## Key Results
1. **E2E Test Infrastructure**: Create robust subprocess testing framework that runs the built application binary
2. **Test Isolation**: Each test runs in dedicated temporary directory with complete isolation between tests
3. **Coverage Collection**: Successfully collect coverage data from subprocess execution using -cover build flags
4. **Coverage Aggregation**: Combine coverage from all e2e tests into single comprehensive report
5. **CI/CD Integration**: Seamlessly integrate with existing GitHub Actions workflow
6. **Test Reliability**: All e2e tests pass consistently without flakiness
7. **Documentation**: Clear documentation for writing and maintaining e2e tests

## Requirements Analysis
- Core Requirements:
  - [x] Build application with coverage instrumentation
  - [x] Execute application as subprocess in tests
  - [x] Collect coverage data from subprocess execution
  - [x] Aggregate coverage from multiple test runs
  - [x] Integrate with existing CI/CD pipeline
- Technical Constraints:
  - [x] Go 1.20+ required for GOCOVERDIR support
  - [x] Tests must be isolated (no shared state)
  - [x] Coverage data must be preserved between runs

## Component Analysis
- Affected Components:
  - **Makefile**
    - Changes needed: Add build-e2e-coverage and test-e2e targets
    - Dependencies: None
  - **tests/e2e/**
    - Changes needed: Complete rewrite with subprocess testing framework
    - Dependencies: os/exec, testing, testify
  - **.github/workflows/ci.yml**
    - Changes needed: Update e2e-tests job to collect coverage
    - Dependencies: Coverage artifacts upload
  - **tests/e2e/README.md** (new)
    - Changes needed: Create comprehensive documentation
    - Dependencies: None

## Functional Changes (E2E Test Cases)
- **Test Case 1: Basic Flow Execution**
  - Execute simple linear flow (hello-world.json)
  - Verify flow completes successfully
  - Check output matches expected results
- **Test Case 2: Conditional Flow Testing**
  - Test flows with conditional branching
  - Verify correct path selection based on conditions
  - Test nested conditions and complex logic
- **Test Case 4: Error Scenarios**
  - Test invalid flow files (malformed JSON, missing fields)
  - Test flows with circular references
  - Test timeout handling and interruption
- **Test Case 5: Advanced Flow Features**
  - Test flows with loops and complex navigation
  - Test complex flow navigation patterns
  - Test performance with large flows
- Architecture:
  - [x] Use subprocess testing pattern with os/exec
  - [x] Create test helper package for common operations
  - [x] Use t.TempDir() for automatic cleanup
- Coverage Collection:
  - [x] Use GOCOVERDIR for subprocess coverage
  - [x] Aggregate with go tool covdata
  - [x] Generate HTML reports for visualization
- Test Organization:
  - [x] One test file per feature area
  - [x] Shared test utilities in helper package
  - [x] Parallel test execution where possible

## Implementation Strategy

### Phase 0: Technology Validation
1. [ ] Verify Go version supports GOCOVERDIR
2. [ ] Create minimal POC with coverage collection
3. [ ] Test coverage aggregation workflow

### Phase 1: Infrastructure Setup
1. [ ] Update Makefile with new targets
   - [ ] build-e2e-coverage target
   - [ ] test-e2e target
   - [ ] coverage-e2e-report target
2. [ ] Create test helper package (tests/e2e/testutil/)
   - [ ] Binary builder with coverage flags
   - [ ] Subprocess executor with timeout
   - [ ] Coverage collector utility
   - [ ] Test environment setup helper

### Phase 2: Core E2E Tests
1. [ ] Remove placeholder_test.go
2. [ ] Create test flow files in testdata/flows/
   - [ ] Basic linear flow (hello-world.json)
   - [ ] Multi-step flow with sequential execution
   - [ ] Conditional flow with branching
   - [ ] Error case flows (invalid syntax, missing fields)
3. [ ] Create flow_execution_test.go
   - [ ] Test basic flow execution
   - [ ] Test conditional branching
   - [ ] Test error handling
4. [ ] Create flow_validation_test.go
   - [ ] Test flow schema validation
   - [ ] Test invalid flow handling

### Phase 3: Coverage Integration
1. [ ] Implement coverage collection in each test
2. [ ] Create coverage aggregation script
3. [ ] Update CI workflow for coverage artifacts
4. [ ] Add coverage reporting to PR comments

### Phase 4: Documentation & Polish
1. [ ] Create comprehensive README.md
2. [ ] Add example test templates
3. [ ] Document debugging procedures
4. [ ] Performance optimization

## Testing Strategy
- Unit Tests:
  - [ ] Test helper utilities
  - [ ] Test coverage collection logic
- Integration Tests:
  - [ ] All e2e tests ARE integration tests
  - [ ] Test parallel execution
  - [ ] Test coverage aggregation
- Reliability Tests:
  - [ ] Run tests 10x to verify no flakiness
  - [ ] Test with different OS environments

## Documentation Plan
- [ ] E2E Test Developer Guide (tests/e2e/README.md)
- [ ] Coverage Collection Guide
- [ ] Troubleshooting Guide
- [ ] CI/CD Integration Notes

## Dependencies
- os/exec (standard library)
- testing (standard library)
- github.com/stretchr/testify
- Go 1.20+ toolchain

## Challenges & Mitigations
- **Challenge**: Coverage collection from subprocess
  - **Mitigation**: Use GOCOVERDIR environment variable (Go 1.20+)
- **Challenge**: Test isolation and cleanup
  - **Mitigation**: Use t.TempDir() for automatic cleanup
- **Challenge**: Flaky tests due to timing
  - **Mitigation**: Implement proper timeouts and retry logic
- **Challenge**: CI/CD coverage aggregation
  - **Mitigation**: Use go tool covdata for merging

## Creative Phases Required
- [x] Test Framework Architecture Design
- [x] Coverage Collection Strategy Design
- [x] CI/CD Integration Architecture

## Branch
- Name: task-20250109-e2e-subprocess-tests
- Created: ✅

## Status
- [x] Initialization complete
- [x] Planning complete
- [x] Creative phase complete
- [x] Technology validation complete
- [ ] Implementation complete
- [ ] Testing complete
- [ ] Reflection complete
- [ ] Archiving complete

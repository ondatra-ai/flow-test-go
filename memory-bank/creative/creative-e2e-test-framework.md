# E2E Test Framework Architecture - Creative Phase

🎨🎨🎨 ENTERING CREATIVE PHASE: TEST FRAMEWORK ARCHITECTURE 🎨🎨🎨

## 📌 CREATIVE PHASE START: Test Framework Architecture
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

### 1️⃣ PROBLEM
**Description**: Design a reusable subprocess testing framework for Go e2e tests that collects coverage data
**Requirements**:
- Build and execute Go binaries with coverage instrumentation
- Manage subprocess lifecycle (start, monitor, timeout, stop)
- Collect coverage data from subprocess execution
- Provide isolated test environments
- Support parallel test execution

**Constraints**:
- Must use Go 1.20+ GOCOVERDIR feature
- Tests must be deterministic and reliable
- Framework must be easy to use and extend

### 2️⃣ OPTIONS

**Option A**: Minimal Helper Functions - Lightweight collection of helper functions
**Option B**: TestSuite Base Class - Object-oriented test suite with lifecycle hooks
**Option C**: Builder Pattern Framework - Fluent API with command builder pattern

### 3️⃣ ANALYSIS

| Criterion | Option A | Option B | Option C |
|-----------|----------|----------|----------|
| Simplicity | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Flexibility | ⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| Maintainability | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Test Isolation | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Learning Curve | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |

**Key Insights**:
- Helper functions are simple but may lead to code duplication across tests
- TestSuite provides structure but Go doesn't have true inheritance
- Builder pattern offers best flexibility and type safety with Go's strengths

### 4️⃣ DECISION
**Selected**: Option C: Builder Pattern Framework
**Rationale**: Provides best balance of flexibility, maintainability, and Go idiomatic design. Fluent API makes tests readable while builder pattern ensures all required setup is done.

### 5️⃣ IMPLEMENTATION NOTES
- Create tests/e2e/testutil package with core components
- Implement E2ETestBuilder with fluent methods
- Use functional options pattern for configuration
- Provide sensible defaults for common scenarios
- Include built-in timeout handling and cleanup

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📌 CREATIVE PHASE END


## Detailed Architecture Design

### Core Components

1. **E2ETestBuilder** - Main builder interface
   ```go
   type E2ETestBuilder struct {
       t           *testing.T
       binaryPath  string
       coverDir    string
       workDir     string
       args        []string
       env         []string
       timeout     time.Duration
   }
   ```

2. **TestCommand** - Represents a built test command
   ```go
   type TestCommand struct {
       cmd         *exec.Cmd
       coverDir    string
       workDir     string
       stdout      *bytes.Buffer
       stderr      *bytes.Buffer
   }
   ```

3. **CoverageCollector** - Handles coverage data collection
   ```go
   type CoverageCollector struct {
       baseDir     string
       testName    string
       coverageData []string
   }
   ```

### Fluent API Example

```go
func TestFlowList(t *testing.T) {
    result := NewE2ETest(t).
        WithBinary("../bin/flow-test-go-e2e").
        WithArgs("list", "flows").
        WithTimeout(30 * time.Second).
        Build().
        Run()

    assert.Equal(t, 0, result.ExitCode)
    assert.Contains(t, result.Stdout, "Available flows:")
}
```

### Directory Structure

```
tests/e2e/
├── testutil/
│   ├── builder.go      # E2ETestBuilder implementation
│   ├── command.go      # TestCommand execution logic
│   ├── coverage.go     # CoverageCollector implementation
│   ├── cleanup.go      # Cleanup utilities
│   └── options.go      # Functional options
├── command_test.go     # Basic CLI command tests
├── flow_test.go        # Flow operation tests
├── mcp_test.go         # MCP integration tests
└── README.md           # Documentation
```

🎨 CREATIVE CHECKPOINT: Framework Structure Defined


🎨🎨🎨 ENTERING CREATIVE PHASE: COVERAGE COLLECTION STRATEGY 🎨🎨🎨

## 📌 CREATIVE PHASE START: Coverage Collection Strategy
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

### 1️⃣ PROBLEM
**Description**: Design a strategy to collect, aggregate, and report coverage data from subprocess e2e tests
**Requirements**:
- Collect coverage from multiple test runs
- Aggregate coverage data across all e2e tests
- Generate human-readable and CI-friendly reports
- Handle coverage data cleanup and archival
- Support partial test runs and re-runs

**Constraints**:
- Must use go tool covdata for processing
- Coverage data can be large for big test suites
- Must integrate with existing coverage tools

### 2️⃣ OPTIONS

**Option A**: Per-Test Coverage Files - Each test creates separate coverage file
**Option B**: Session-Based Collection - Group tests by session with timestamps
**Option C**: Hierarchical Coverage Structure - Organize by test package/suite/case

### 3️⃣ ANALYSIS

| Criterion | Option A | Option B | Option C |
|-----------|----------|----------|----------|
| Granularity | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| Aggregation Speed | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Debugging | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ |
| CI Integration | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Storage Efficiency | ⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |

**Key Insights**:
- Per-test files offer best debugging but slow aggregation
- Session-based is simple but loses test-level granularity
- Hierarchical structure provides best balance and CI integration

### 4️⃣ DECISION
**Selected**: Option C: Hierarchical Coverage Structure
**Rationale**: Provides excellent organization for CI reporting, efficient storage, and maintains test-level granularity for debugging when needed.

### 5️⃣ IMPLEMENTATION NOTES
- Use directory structure: coverage/e2e/{package}/{test_name}/
- Implement coverage manifest file for tracking
- Create aggregation script with filtering options
- Support incremental coverage updates
- Generate both text and HTML reports

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📌 CREATIVE PHASE END


## Coverage Collection Workflow

### Directory Structure
```
coverage/
├── e2e/
│   ├── command/              # Package-level grouping
│   │   ├── TestHelp/         # Test-specific coverage
│   │   ├── TestVersion/
│   │   └── TestInvalid/
│   ├── flow/
│   │   ├── TestList/
│   │   └── TestExecute/
│   ├── mcp/
│   │   └── TestServerLoad/
│   ├── manifest.json         # Coverage metadata
│   ├── e2e-merged/           # Aggregated coverage data
│   ├── e2e.out               # Final coverage profile
│   └── e2e.html              # HTML report
```

### Coverage Manifest Schema
```json
{
  "version": "1.0",
  "timestamp": "2024-01-09T10:00:00Z",
  "tests": [
    {
      "package": "command",
      "test": "TestHelp",
      "coverDir": "coverage/e2e/command/TestHelp",
      "duration": "1.2s",
      "status": "passed"
    }
  ],
  "summary": {
    "totalTests": 15,
    "passed": 15,
    "failed": 0,
    "coverage": "89.5%"
  }
}
```

### Aggregation Process
```bash
# 1. Merge all test coverage data
go tool covdata merge -i=coverage/e2e/command/*,coverage/e2e/flow/*,coverage/e2e/mcp/* -o=coverage/e2e-merged

# 2. Convert to standard format
go tool covdata textfmt -i=coverage/e2e-merged -o=coverage/e2e.out

# 3. Generate reports
go tool cover -html=coverage/e2e.out -o=coverage/e2e.html
go tool cover -func=coverage/e2e.out > coverage/e2e-summary.txt
```

🎨 CREATIVE CHECKPOINT: Coverage Strategy Detailed


🎨🎨🎨 ENTERING CREATIVE PHASE: CI/CD INTEGRATION ARCHITECTURE 🎨🎨🎨

## 📌 CREATIVE PHASE START: CI/CD Integration Architecture
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

### 1️⃣ PROBLEM
**Description**: Design GitHub Actions integration for e2e tests with coverage reporting
**Requirements**:
- Run e2e tests in CI environment
- Collect and aggregate coverage data
- Upload coverage artifacts
- Display coverage in PR comments
- Integrate with existing CI workflow

**Constraints**:
- Must work with GitHub Actions
- Should not significantly increase CI time
- Must handle flaky test retries
- Coverage artifacts should be accessible

### 2️⃣ OPTIONS

**Option A**: Sequential Workflow - Run e2e after unit tests complete
**Option B**: Parallel Jobs - Run e2e tests in parallel with other checks
**Option C**: Matrix Strategy - Run e2e tests across multiple configurations

### 3️⃣ ANALYSIS

| Criterion | Option A | Option B | Option C |
|-----------|----------|----------|----------|
| CI Speed | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| Simplicity | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Resource Usage | ⭐⭐⭐⭐ | ⭐⭐ | ⭐ |
| Debugging | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Coverage Accuracy | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**Key Insights**:
- Sequential is simple but slows down CI feedback
- Parallel jobs give fastest feedback for developers
- Matrix strategy useful for cross-platform but overkill for this project

### 4️⃣ DECISION
**Selected**: Option B: Parallel Jobs
**Rationale**: Provides fastest CI feedback while maintaining reasonable resource usage. E2E tests can run simultaneously with linting and other checks.

### 5️⃣ IMPLEMENTATION NOTES
- Create dedicated e2e-tests job in CI workflow
- Use job dependencies only where necessary
- Implement retry mechanism for flaky tests
- Upload coverage as GitHub artifact
- Add coverage comment to PRs

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📌 CREATIVE PHASE END

## CI/CD Workflow Design

### Updated GitHub Actions Workflow
```yaml
name: CI

on:
  pull_request:
  push:
    branches: [main]

jobs:
  # Existing jobs run in parallel
  unit-tests:
    # ... existing unit test configuration

  lint:
    # ... existing lint configuration

  # NEW: E2E tests with coverage
  e2e-tests:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.5'

      - name: Build with coverage
        run: make build-e2e-coverage

      - name: Run E2E tests
        id: e2e-test
        run: |
          make test-e2e || (echo "::warning::E2E tests failed, retrying..." && make test-e2e)

      - name: Generate coverage report
        if: always()
        run: |
          make coverage-e2e-report
          echo "COVERAGE=$(go tool cover -func=coverage/e2e.out | grep total | awk '{print $3}')" >> $GITHUB_ENV

      - name: Upload coverage artifacts
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: e2e-coverage-report
          path: |
            coverage/e2e.out
            coverage/e2e.html
            coverage/e2e-summary.txt
            coverage/manifest.json
          retention-days: 30

      - name: Comment PR with coverage
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            const coverage = process.env.COVERAGE;
            const comment = `### E2E Test Coverage Report

            **Coverage**: ${coverage}

            View detailed report in the [artifacts](https://github.com/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.runId})`;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });

  # All tests must pass
  all-tests:
    needs: [unit-tests, lint, e2e-tests]
    runs-on: ubuntu-latest
    steps:
      - name: All checks passed
        run: echo "All tests passed!"
```

### Makefile Additions
```makefile
# Build with coverage instrumentation
.PHONY: build-e2e-coverage
build-e2e-coverage:
	@mkdir -p bin
	go build -cover -covermode=atomic -o bin/flow-test-go-e2e ./cmd/flow-test-go

# Run e2e tests
.PHONY: test-e2e
test-e2e: build-e2e-coverage
	@mkdir -p coverage/e2e
	go test -v -timeout 5m ./tests/e2e/...

# Generate coverage report
.PHONY: coverage-e2e-report
coverage-e2e-report:
	@echo "Aggregating e2e coverage data..."
	go tool covdata merge -i=coverage/e2e/* -o=coverage/e2e-merged
	go tool covdata textfmt -i=coverage/e2e-merged -o=coverage/e2e.out
	go tool cover -html=coverage/e2e.out -o=coverage/e2e.html
	go tool cover -func=coverage/e2e.out > coverage/e2e-summary.txt
	@echo "E2E coverage report generated"
```

🎨 CREATIVE CHECKPOINT: CI/CD Integration Designed

🎨🎨🎨 EXITING CREATIVE PHASE - ALL DESIGN DECISIONS MADE 🎨🎨🎨

# System Patterns

## Core Technologies & Libraries

### 1. CLI Application (Cobra)
- **Framework**: Cobra for command structure, subcommands, and help text
- **Commands**: Root command with subcommands (execute, list, validate)
- **Configuration**: Viper for config management (env vars, files, defaults)

### 2. LangGraph Integration
- **Purpose**: State machine orchestration for flow execution
- **Features**: Built-in state management, checkpointing, graph compilation
- **Usage**: Convert JSON flow configs to LangGraph nodes/edges at runtime

### 3. OpenRouter (LLM Access)
- **Purpose**: Unified API for 200+ LLM models
- **Adapter**: Wraps OpenRouter client to implement LangChain interface
- **Models**: Dynamic selection based on task requirements

### 4. MCP Servers (Embedded)
- **Protocol**: Model Context Protocol for tool integration
- **Implementation**: Embedded in-process as tool providers, not external servers
- **Registration**: Tools registered at startup from compiled-in providers
- **Access Control**: Each flow step explicitly declares which tools it can use

### 5. Service Architecture
Services that implement functionality:
- **Tool Registry**: Manages embedded tool providers and registration
- **Context Manager**: Handles flow context and state management
- **GitHub Service**: PR/issue operations via google/go-github
- **Flow Engine**: Converts configs to executable LangGraph
- **Checkpoint Service**: Saves/restores execution state and context

## Code Organization
```
flow-test-go/
├── cmd/
│   ├── commands/          # CLI command implementations
│   │   ├── root.go        # Root command setup
│   │   ├── list.go        # List flows command
│   │   └── .flows/        # Test flow configurations
│   │       └── flows/     # Flow JSON files
│   └── flow-test-go/      # Main entry point
│       └── main.go
├── internal/
│   ├── ai/                # AI/LLM integration
│   ├── config/            # Configuration management
│   ├── flow/              # Flow engine implementation
│   ├── github/            # GitHub API client
│   └── mcp/               # MCP server management
├── pkg/
│   ├── errors/            # Custom error types
│   └── types/             # Shared type definitions
│       ├── flow.go        # Flow structures
│       └── mcp.go         # MCP types
├── tests/
│   └── e2e/               # End-to-end tests
│       ├── testdata/      # Test flow files
│       │   └── flows/     # Test flow JSONs
│       ├── testutil/      # Test utilities
│       └── *_test.go      # Test files
├── bin/                   # Compiled binaries
├── coverage/              # Test coverage reports
└── memory-bank/           # AI Memory Bank system
```

## System Architecture Diagram

```mermaid
graph TB
    subgraph "CLI Layer"
        CLI[flow-test-go CLI]
        CMD[Commands<br/>execute/list/validate/tools/resume]
    end

    subgraph "Configuration Layer"
        FLOW[Flow JSON]
        CTX_SRC[Context Sources<br/>CLI args/ENV/YAML]
        TPL[Template Engine<br/>{env.*} {context.*}]
    end

    subgraph "Flow Engine Layer"
        LG[LangGraph Engine<br/>Shared/Thread-safe]
        EXEC[Flow Executor]
        VAL[Flow Validator]
        CHK[Checkpoint Manager<br/>Optional]
    end

    subgraph "Context Management"
        CTX[Context Object<br/>Accumulating State]
        CTX_ISO[Context Isolation<br/>Per-step Access Control]
    end

    subgraph "Step Types"
        PROMPT[Prompt Step<br/>LLM + Tools]
        COND[Condition Step<br/>JS Expression]
        TOOL[Tool Step<br/>Direct Execution]
    end

    subgraph "External Integration"
        OR[OpenRouter API<br/>200+ LLM Models]
        LLM[LLM Decision Making]
    end

    subgraph "Embedded Tool Layer"
        TR[Tool Registry]
        GH[GitHub Tools]
        SLACK[Slack Tools]
        AWS[AWS Tools]
        CUSTOM[Custom Tools]
    end

    subgraph "Testing Infrastructure"
        E2E[E2E Test Runner]
        TEST_ENV[Test Environments<br/>.env Resources]
        VERIFY[Effect Verification<br/>via APIs]
    end

    subgraph "Storage"
        FS[Local Filesystem<br/>Checkpoints/Logs]
    end

    %% CLI Flow
    CLI --> CMD
    CMD --> FLOW
    CMD --> CTX_SRC

    %% Configuration Processing
    FLOW --> TPL
    CTX_SRC --> TPL
    TPL --> VAL
    VAL --> EXEC

    %% Flow Execution
    EXEC --> LG
    LG --> CTX
    CTX --> CTX_ISO

    %% Step Execution
    LG --> PROMPT
    LG --> COND
    LG --> TOOL

    %% Context Access
    CTX_ISO -.->|Read Only Declared Vars| PROMPT
    CTX_ISO -.->|Read Only Declared Vars| COND
    CTX_ISO -.->|Read Only Declared Vars| TOOL

    %% Tool Access
    PROMPT -->|Tools Array| TR
    TOOL -->|Direct Call| TR

    %% LLM Integration
    PROMPT --> OR
    OR --> LLM
    LLM -->|Tool Calls| TR

    %% Tool Registry
    TR --> GH
    TR --> SLACK
    TR --> AWS
    TR --> CUSTOM

    %% Results Flow Back
    GH -->|Write Results| CTX
    SLACK -->|Write Results| CTX
    AWS -->|Write Results| CTX
    CUSTOM -->|Write Results| CTX

    %% Checkpointing
    CTX -->|Save State| CHK
    CHK -->|Persist| FS

    %% Testing
    E2E -->|Execute Flow| CLI
    E2E -->|Verify Effects| TEST_ENV
    TEST_ENV -->|API Calls| VERIFY

    %% Styling
    classDef cli fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef engine fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef tool fill:#e8f5e9,stroke:#1b5e20,stroke-width:2px
    classDef external fill:#fff3e0,stroke:#e65100,stroke-width:2px
    classDef storage fill:#f5f5f5,stroke:#424242,stroke-width:2px
    classDef test fill:#fce4ec,stroke:#880e4f,stroke-width:2px

    class CLI,CMD cli
    class LG,EXEC,VAL,CHK,CTX,CTX_ISO engine
    class TR,GH,SLACK,AWS,CUSTOM tool
    class OR,LLM external
    class FS storage
    class E2E,TEST_ENV,VERIFY test
```

### Key Architectural Points

1. **Stateless Execution**: Each flow execution gets isolated instances (except LangGraph engine)
2. **Embedded Tools**: All MCP tools are compiled into the binary, not external processes
3. **Context Isolation**: Steps only access explicitly declared context variables
4. **Non-Interactive**: All inputs provided upfront, no mid-execution prompts
5. **E2E Testing**: Mandatory verification of actual side effects

## Flow Structure & Context Management

### Flow Execution Model
1. **Context Accumulation**: Flow maintains a growing context object during execution
2. **Controlled Access**: Each step declares which context variables it needs
3. **LLM Integration**: Prompt steps send text to OpenRouter, LLM can invoke declared tools
4. **Direct Tool Execution**: Non-LLM steps can directly execute tools without AI

### Step Types
1. **prompt**: Sends prompt to LLM with tool access
   ```json
   {
     "type": "prompt",
     "prompt": "Create a GitHub issue for the bug",
     "tools": ["github_create_issue"],
     "context": ["bug_description", "repository_name"]
   }
   ```

2. **condition**: Evaluates JavaScript expression on context
   ```json
   {
     "type": "condition",
     "condition": "context.issue_number > 0",
     "context": ["issue_number"],
     "yes": "success-step",
     "no": "failure-step"
   }
   ```

3. **tool** (or action): Direct tool execution without LLM
   ```json
   {
     "type": "tool",
     "tool": "github_get_issue",
     "parameters": {
       "issue_number": "${context.issue_number}"
     },
     "output": "issue_body",
     "context": ["issue_number"]
   }
   ```

### Context Access Pattern
- Each step must declare `"context": ["var1", "var2"]` to access specific variables
- Steps can only read declared context variables (isolation)
- Steps can write new values to context for subsequent steps
- Context is preserved in checkpoints for resume capability

### Checkpointing
- LangGraph saves both execution position and accumulated context
- Failed flows can resume from last successful step
- No need to re-execute from beginning
- Context state is fully restored on resume

## E2E Testing Strategy

### Overview
**MANDATORY**: Create E2E tests BEFORE implementing features. Tests define expected behavior.

### Testing Approach: Subprocess Execution

The application uses **subprocess testing** rather than in-process testing:
- **Subprocess**: Build the full CLI binary and execute it as a child process
- **Coverage**: Use Go's built-in GOCOVERDIR for subprocess coverage
- **Isolation**: Each test runs in a fresh environment

### Why Subprocess Testing?
1. **Real-world accuracy**: Tests the actual binary users will run
2. **Complete integration**: Tests CLI parsing, config loading, and execution
3. **Process isolation**: No shared state between tests
4. **Coverage support**: Go 1.20+ supports subprocess coverage collection

### Test Structure

#### 1. Flow Definition
Create a JSON flow file in `tests/e2e/testdata/flows/`:
```json
{
  "id": "test-flow",
  "name": "Test Flow",
  "description": "Flow for testing specific feature",
  "initialStep": "step1",
  "steps": {
    "step1": {
      "type": "prompt",
      "prompt": "Execute test action",
      "tools": ["github_create_issue"],
      "context": ["user_input", "repository_name"],
      "nextStep": "step2"
    },
    "step2": {
      "type": "condition",
      "condition": "context.issue_number > 0",
      "context": ["issue_number"],
      "yes": "success-step",
      "no": "failure-step"
    }
  }
}
```

#### 2. Test Implementation
```go
func TestFeatureName(t *testing.T) {
    // Mark start of test execution
    exec := testutil.NewTestExecution(t, "feature-name").Start()

    // Build and run the test
    result := testutil.NewFlowTest(t).
        WithFlow(testutil.FlowPath("category/test-flow.json")).
        WithTimeout(30 * time.Second).
        ExpectSuccess().
        ExpectOutput("expected output substring").
        Run()

    // Record execution time and coverage
    duration := exec.Complete(result)

    // Additional assertions
    assert.Contains(t, result.Stdout, "specific output")
    assert.Equal(t, 0, result.ExitCode)

    // For GitHub operations, verify side effects
    // e.g., check if issue was created, PR was updated
}
```

### Test Patterns

#### Basic Flow Test
```go
// Test single flow execution
result := testutil.NewFlowTest(t).
    WithFlow("testdata/flows/basic/single-step.json").
    ExpectSuccess().
    Run()
```

#### Error Handling Test
```go
// Test error scenarios
result := testutil.NewFlowTest(t).
    WithFlow("testdata/flows/error-cases/invalid-ref.json").
    ExpectFailure().
    ExpectError("invalid step reference").
    Run()
```

#### Output Validation
```go
// Validate specific outputs
result := testutil.NewFlowTest(t).
    WithFlow("testdata/flows/github/create-issue.json").
    ExpectSuccess().
    ExpectOutput("Issue created: #").
    Run()

// Parse output for detailed validation
issueNumber := extractIssueNumber(result.Stdout)
verifyIssueExists(t, issueNumber)
```

### Coverage Collection

The test framework automatically:
1. Builds binary with `-cover` flag
2. Sets `GOCOVERDIR` for each test
3. Aggregates coverage data
4. Generates HTML reports

```bash
# Run all e2e tests with coverage
make test-e2e-coverage

# View coverage report
open coverage/e2e.html
```

### Writing New Tests: Step-by-Step

1. **Define the feature** in a flow JSON file
2. **Create test file** following naming convention
3. **Use FlowTestBuilder** API for clean test structure
4. **Validate outputs** - stdout, stderr, exit codes
5. **Check side effects** - files created, APIs called
6. **Run locally** before pushing

### Common Test Scenarios

#### GitHub Integration
```go
// Test PR creation
result := testutil.NewFlowTest(t).
    WithFlow("flows/github/create-pr.json").
    WithEnv("GITHUB_TOKEN", token).
    ExpectOutput("Pull request created").
    Run()

// Verify PR exists
prNumber := extractPRNumber(result.Stdout)
pr := getGitHubPR(t, prNumber)
assert.Equal(t, "open", pr.State)
```

#### MCP Server Communication
```go
// Test tool discovery
result := testutil.NewFlowTest(t).
    WithConfig(".flows/servers/test-server.json").
    WithFlow("flows/mcp/tool-discovery.json").
    ExpectOutput("Available tools:").
    Run()
```

#### Error Recovery
```go
// Test checkpoint recovery
result1 := testutil.NewFlowTest(t).
    WithFlow("flows/long-running.json").
    WithTimeout(5 * time.Second).
    ExpectFailure(). // Times out
    Run()

// Resume from checkpoint
result2 := testutil.NewFlowTest(t).
    WithFlow("flows/long-running.json").
    WithCheckpoint(result1.CheckpointPath).
    ExpectSuccess().
    Run()
```

## Development Workflow

1. **Test First**: Write E2E test defining expected behavior
2. **Implement**: Build feature to make test pass
3. **Refactor**: Clean up while keeping tests green
4. **Document**: Update flow examples and README

## Key Standards

### Error Handling
- Tool errors classified as retryable (rate limits, network) or fatal
- Errors written to context for conditional handling
- No automatic rollback or compensation
- Clear user-facing messages

### Configuration
- JSON with schema validation
- Template syntax for `{env.*}` and `{context.*}` references
- No built-in flow composition or imports

### Performance
- Stateless architecture, no global state
- Each execution gets isolated instances
- LangGraph engine is shared but thread-safe
- No concurrent flow execution by design

## Tool Implementation

### Creating New Tools
1. Implement Tool interface in Go
2. Register with Tool Registry at startup
3. Define input/output schemas
4. Read from context, write results back
5. **MANDATORY**: Create e2e tests verifying actual side effects

### E2E Testing for Tools
- Execute flows using the tool
- Verify real effects via APIs (e.g., check Slack message was sent)
- Clean up test resources
- 95% functionality coverage requirement

### Test Resources
- Pre-existing: Provided via .env (repositories, Slack workspace)
- Dynamic: Created during tests (issues, comments)
- Cleanup: Automatic after test execution

## Operational Characteristics

### Execution Model
- Non-interactive: All inputs provided upfront
- No user prompts during execution
- Single flow at a time
- No built-in scheduling or orchestration

### Debugging & Development
- `--log-level debug` for detailed logs
- Separate stdout (user) and stderr (debug)
- Log sensitive data redacted
- Potential dry-run mode for validation

### Logging Levels
- User output → stdout
- Debug logs → stderr or file
- Step execution with timestamps
- Tool invocations with parameters

### Rate Limiting
- OpenRouter handles LLM rate limits
- Tools implement their own limits
- No global rate management
- Retryable errors for rate limits

### Validation
- Schema compliance checking
- Step reference integrity
- Tool availability verification
- Context variable consistency
- Circular reference detection

## CLI Usage

### Commands
```bash
flow-test-go execute <flow.json> [options]
flow-test-go list
flow-test-go validate <flow.json>
flow-test-go tools
flow-test-go resume <checkpoint-id>  # if checkpointing enabled
```

### Context Sources
```bash
# CLI arguments
flow-test-go execute flow.json --context "key=value"

# Environment variables (in flow.json)
"token": "{env.GITHUB_TOKEN}"

# YAML file
flow-test-go execute flow.json --context-file context.yaml
```

### Checkpointing
- System-wide decision: always on or always off
- Stores position and full context
- Local filesystem storage
- Resume capability for failures

## Design Decisions

### What's NOT Included
- No observability/metrics dashboards
- No workflow scheduling
- No flow composition/imports
- No interactive debugging
- No automatic rollback
- No mock mode for tools

### Architectural Constraints
- CLI tool only, not a service
- Flows must be re-runnable
- External orchestration for scheduling
- Test environments for development

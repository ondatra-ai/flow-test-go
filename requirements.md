# Functional and Non-Functional Requirements

## Feature 1: Flow Execution Engine

### User Story 1.1: Execute Flow from JSON Definition
**As a** developer
**I want to** execute flows defined in JSON format
**So that** I can automate complex multi-step processes

**Scenario 1.1.1: Basic flow execution**
- **Given** I have a valid flow JSON file with steps and conditions
- **When** I run `flow-test-go execute my-flow.json`
- **Then** the system should convert the JSON to a LangGraph state machine
- **And** execute each step in the defined order
- **And** return exit code 0 on success

**Scenario 1.1.2: Flow with conditions**
- **Given** I have a flow with conditional branching
- **When** the flow reaches a condition step
- **Then** the system should evaluate the expression using a Go expression evaluator
- **And** follow the "yes" or "no" path based on the result

### User Story 1.2: Step Type Support
**As a** flow author
**I want to** use different types of steps
**So that** I can handle various automation scenarios

**Scenario 1.2.1: Prompt step execution**
- **Given** I have a step with type "prompt"
- **When** the flow executes this step
- **Then** the system should send the prompt to OpenRouter
- **And** the LLM should be able to invoke tools from the declared tools array
- **And** results should be written to context

**Scenario 1.2.2: Tool step execution**
- **Given** I have a step with type "tool"
- **When** the flow executes this step
- **Then** the system should directly call the specified tool
- **And** pass the defined parameters
- **And** write results to context without LLM involvement

**Scenario 1.2.3: Condition step evaluation**
- **Given** I have a step with type "condition"
- **When** the flow executes this step
- **Then** the system should evaluate the expression against context using a Go expression evaluator
- **And** proceed to the appropriate next step

**Technical Note**: The condition expressions like `"context.issue_number > 0"` or `"result.success === true"` would be evaluated by a Go-based expression library (such as `govaluate` or similar) rather than a JavaScript engine, which aligns with the Go technology stack.

## Feature 2: Embedded Tool System

### User Story 2.1: Tool Registration and Discovery
**As a** developer
**I want to** register tools at application startup
**So that** they are available for flow execution

**Scenario 2.1.1: Tool registration**
- **Given** the application starts
- **When** the tool registry initializes
- **Then** all embedded tool providers should be registered
- **And** tools should be available for discovery

**Scenario 2.1.2: Tool access control**
- **Given** a step declares specific tools in its "tools" array
- **When** the step executes
- **Then** only the declared tools should be accessible
- **And** attempts to use undeclared tools should fail

### User Story 2.2: Tool Implementation Requirements
**As a** tool developer
**I want to** implement new tools following a standard interface
**So that** they integrate seamlessly with the flow engine

**Scenario 2.2.1: Tool interface compliance**
- **Given** I implement a new tool
- **When** I register it with the Tool Registry
- **Then** it should implement the required Tool interface
- **And** define input/output schemas
- **And** handle context read/write operations

**Scenario 2.2.2: Mandatory e2e testing**
- **Given** I implement a new tool
- **When** I create tests for it
- **Then** I must create e2e tests that verify actual side effects
- **And** the tests must use real external APIs
- **And** verify the tool's effects (e.g., message sent, issue created)
- **And** clean up test resources afterward

## Feature 3: Context Management

### User Story 3.1: Context Accumulation
**As a** flow engine
**I want to** maintain a growing context object during execution
**So that** steps can share data while maintaining isolation

**Scenario 3.1.1: Context isolation**
- **Given** a step declares `"context": ["var1", "var2"]`
- **When** the step executes
- **Then** it should only have access to var1 and var2 from context
- **And** should not be able to read other context variables
- **And** should be able to write new values to context

**Scenario 3.1.2: Context persistence**
- **Given** step1 writes `issue_number: 123` to context
- **When** step2 declares `"context": ["issue_number"]`
- **Then** step2 should be able to read the value 123
- **And** step3 without declaring issue_number should not access it

### User Story 3.2: Initial Context Population
**As a** user
**I want to** provide initial context values
**So that** flows can access external data

**Scenario 3.2.1: CLI argument context**
- **Given** I run `flow-test-go execute flow.json --context "repo=myrepo"`
- **When** the flow starts
- **Then** the context should contain `repo: "myrepo"`
- **And** be available to steps that declare it

**Scenario 3.2.2: Environment variable context**
- **Given** a flow contains `"token": "{env.GITHUB_TOKEN}"`
- **When** the flow executes
- **Then** the system should resolve the environment variable
- **And** pass the actual value to the tool

**Scenario 3.2.3: YAML file context**
- **Given** I run `flow-test-go execute flow.json --context-file context.yaml`
- **When** the flow starts
- **Then** the context should be populated from the YAML file

## Feature 4: CLI Interface

### User Story 4.1: Command Suite
**As a** user
**I want to** interact with flows through a CLI
**So that** I can manage and execute flows easily

**Scenario 4.1.1: Execute command**
- **Given** I have a valid flow file
- **When** I run `flow-test-go execute myflow.json`
- **Then** the flow should execute successfully
- **And** return appropriate exit codes

**Scenario 4.1.2: List command**
- **Given** I have flows in .flows directory
- **When** I run `flow-test-go list`
- **Then** I should see all available flows
- **And** their descriptions

**Scenario 4.1.3: Validate command**
- **Given** I have a flow file
- **When** I run `flow-test-go validate myflow.json`
- **Then** the system should check schema compliance
- **And** verify step reference integrity
- **And** validate tool availability
- **And** check context variable consistency

**Scenario 4.1.4: Tools command**
- **Given** the application has registered tools
- **When** I run `flow-test-go tools`
- **Then** I should see all available tools
- **And** their parameters and descriptions

## Feature 5: Error Handling and Recovery

### User Story 5.1: Error Classification
**As a** flow engine
**I want to** classify and handle different types of errors
**So that** flows can respond appropriately

**Scenario 5.1.1: Retryable errors**
- **Given** a tool encounters a rate limit error
- **When** the error occurs
- **Then** it should be classified as retryable
- **And** written to context for flow handling
- **And** allow conditional retry logic

**Scenario 5.1.2: Fatal errors**
- **Given** a tool encounters invalid credentials
- **When** the error occurs
- **Then** it should be classified as fatal
- **And** written to context
- **And** allow flow to handle gracefully

### User Story 5.2: Checkpointing (Optional)
**As a** system architect
**I want to** decide on checkpointing strategy
**So that** flows can recover from failures

**Scenario 5.2.1: System-wide checkpointing**
- **Given** checkpointing is enabled system-wide
- **When** any flow executes
- **Then** the system should save position and context at each step
- **And** allow resumption from the last checkpoint

**Scenario 5.2.2: No rollback behavior**
- **Given** a flow fails after creating external resources
- **When** the failure occurs
- **Then** the created resources should remain
- **And** no automatic cleanup should occur

## Feature 6: Validation System

### User Story 6.1: Flow Validation
**As a** flow author
**I want to** validate my flows before execution
**So that** I can catch errors early

**Scenario 6.1.1: Schema validation**
- **Given** I have a flow JSON file
- **When** I validate it
- **Then** the system should check required fields and types
- **And** verify JSON structure compliance

**Scenario 6.1.2: Reference integrity**
- **Given** I have a flow with step references
- **When** I validate it
- **Then** all nextStep, yes, and no references should point to valid steps
- **And** no circular references should exist

**Scenario 6.1.3: Tool availability**
- **Given** I have steps declaring tools
- **When** I validate the flow
- **Then** all declared tools should exist in the registry
- **And** unavailable tools should cause validation failure

## Feature 7: Testing Infrastructure

### User Story 7.1: E2E Testing Framework
**As a** developer
**I want to** test flows end-to-end
**So that** I can ensure they work in real environments

**Scenario 7.1.1: Subprocess execution**
- **Given** I have an e2e test
- **When** the test runs
- **Then** it should build the CLI binary with coverage
- **And** execute it as a subprocess
- **And** collect coverage data

**Scenario 7.1.2: Effect verification**
- **Given** a flow creates external resources
- **When** the e2e test runs
- **Then** it should verify the resources were actually created
- **And** use real APIs to check results
- **And** clean up test resources

**Scenario 7.1.3: Test resource management**
- **Given** tests need external resources
- **When** tests run
- **Then** pre-existing resources should be provided via .env
- **And** dynamic resources should be created and cleaned up
- **And** 95% of functionality should be covered

## Feature 8: Logging and Debugging

### User Story 8.1: Structured Logging
**As a** developer
**I want to** access detailed logs
**So that** I can debug flow execution issues

**Scenario 8.1.1: Log levels**
- **Given** I run a flow with `--log-level debug`
- **When** the flow executes
- **Then** user output should go to stdout
- **And** debug logs should go to stderr
- **And** sensitive data should be redacted

**Scenario 8.1.2: Step logging**
- **Given** a flow executes
- **When** each step runs
- **Then** it should log with timestamp and step ID
- **And** log tool invocations and results
- **And** log context changes

### User Story 8.2: Dry-Run Mode
**As a** flow author
**I want to** preview flow execution
**So that** I can understand behavior before making changes

**Scenario 8.2.1: Dry-run execution**
- **Given** I run a flow in dry-run mode
- **When** the system processes it
- **Then** it should show each step configuration
- **And** display context variable access
- **And** show which tools would be called
- **And** not execute actual tools

## Non-Functional Requirements

### NFR-1: Performance
- **Stateless Architecture**: No global state between executions
- **Isolated Instances**: Each execution gets separate tool instances
- **Minimal Startup**: Fast application initialization
- **Single Flow**: One flow execution at a time by design

### NFR-2: Security
- **No Credential Storage**: Credentials via environment or templates only
- **Context Isolation**: Steps only access declared variables
- **Input Validation**: All inputs validated before execution
- **Audit Logging**: Tool usage logged for security

### NFR-3: Reliability
- **Error Recovery**: Classified error handling
- **Validation**: Pre-execution flow validation
- **Idempotent Design**: Flows should be re-runnable
- **Test Coverage**: 95% e2e test coverage requirement

### NFR-4: Usability
- **Non-Interactive**: All inputs provided upfront
- **Clear Messages**: User-friendly error messages
- **CLI Interface**: Standard command-line patterns
- **Documentation**: Comprehensive flow examples

### NFR-5: Maintainability
- **Tool Interface**: Standard interface for new tools
- **Test First**: Mandatory tests before implementation
- **No Dependencies**: Embedded tools, no external processes
- **Schema Validation**: Structured configuration validation

### NFR-6: Architectural Constraints
- **CLI Only**: Not a service or daemon
- **No Orchestration**: No built-in scheduling
- **No Composition**: No flow imports or reuse
- **No Rollback**: No automatic compensation
- **No Observability**: No metrics or dashboards

### NFR-7: Integration
- **OpenRouter**: Single LLM provider integration
- **Real APIs**: E2E tests use actual external services
- **Template System**: Environment and context variable substitution
- **LangGraph**: Production-ready state machine engine

# Task: Go CLI AI Tool Initialization

## Description
Initialize a Go project for a CLI tool that can execute AI agents through command-line interface, integrate with MCP (Model Context Protocol) servers for tool execution, integrate with GitHub for task management and development, and support flow-based execution system for multi-step AI workflows.

## Complexity
Level: 3
Type: Intermediate Feature

## Prerequisites (MUST BE RESOLVED BEFORE IMPLEMENTATION)
- **Go Installation Required**: Go is not currently installed on the system
  - Install Go from https://go.dev/dl/
  - Verify installation with `go version`
  - Ensure Go is in system PATH

## Technology Stack
- Language: Go (latest stable version - 1.21+)
- CLI Framework: Cobra (for command structure and help generation)
- Configuration: Viper (for configuration management)
- HTTP Client: Go standard library (net/http) + GitHub Go client
- JSON Processing: encoding/json (standard library)
- Testing: Go testing package + testify
- Build Tool: Go modules
- MCP Communication: stdio/HTTP (based on transport type)

## Technology Validation Checkpoints
- [ ] **Go installation verified** - BLOCKED: Go not found in system PATH
- [ ] Go module initialization verified (`go mod init`)
- [ ] Required dependencies identified and documented
- [ ] Cobra CLI framework integration tested
- [ ] GitHub API client library validated
- [ ] MCP protocol communication verified
- [ ] JSON configuration parsing tested
- [ ] Basic CLI command structure works

## Status
- [x] Initialization complete
- [x] Planning complete
- [ ] Technology validation complete - BLOCKED by Go installation
- [x] Creative phase complete
- [ ] Implementation complete
- [ ] Testing complete
- [ ] Documentation complete

## Contracts, Schemas and Interface Updates
### Core Interfaces (to be designed in Creative phase)
- Flow Definition Schema (JSON structure)
- MCP Server Configuration Schema
- Tool Discovery Interface
- Flow Execution Engine Interface
- GitHub Integration Interface
- AI Provider Interface

### Type Definitions Needed
- FlowDefinition struct
- MCPServerConfig struct
- ToolDefinition struct
- ExecutionContext struct
- StepResult struct
- Condition evaluation types

## Functional Changes (E2E Test Cases)
### Test Case 1: Basic Flow Execution
```go
// Test executing a simple single-step flow
func TestBasicFlowExecution(t *testing.T) {
    // Given: A simple flow definition with one prompt step
    // When: User runs `flow-test-go execute simple-flow`
    // Then: AI agent executes the prompt and returns result
}
```

### Test Case 2: MCP Server Tool Discovery
```go
// Test discovering tools from MCP server
func TestMCPToolDiscovery(t *testing.T) {
    // Given: An MCP server configuration
    // When: CLI starts and queries available tools
    // Then: Tools are discovered and cached for use
}
```

### Test Case 3: GitHub Integration
```go
// Test creating GitHub issue from flow
func TestGitHubIssueCreation(t *testing.T) {
    // Given: A flow with GitHub issue creation step
    // When: Flow executes with valid GitHub credentials
    // Then: Issue is created in specified repository
}
```

### Test Case 4: Conditional Flow Execution
```go
// Test flow with conditional branching
func TestConditionalFlow(t *testing.T) {
    // Given: A flow with condition step
    // When: Condition evaluates to true/false
    // Then: Correct branch is executed
}
```

## Requirements Analysis
### Core Requirements:
- [x] CLI tool that can execute AI agents
- [x] Support for multiple AI providers (OpenAI, Anthropic, etc.)
- [x] MCP server integration for tool execution
- [x] GitHub API integration
- [x] Flow-based orchestration system
- [x] Configuration management (.flows directory)
- [x] Error handling and recovery
- [x] Logging and debugging support

### Technical Constraints:
- [x] Must support stdio and HTTP transport for MCP
- [x] Must handle async operations gracefully
- [x] Must maintain conversation context
- [x] Must support parallel MCP server execution
- [x] Must validate flow definitions before execution

## Component Analysis
### Affected Components:
1. **CLI Interface Layer** (`cmd/`)
   - Changes needed: Create command structure (execute, list, validate)
   - Dependencies: Cobra, Viper

2. **Flow Engine** (`internal/flow/`)
   - Changes needed: Create flow parser, executor, condition evaluator
   - Dependencies: JSON parsing, AI provider clients

3. **MCP Integration** (`internal/mcp/`)
   - Changes needed: Create server manager, tool discovery, communication protocol
   - Dependencies: Process management, stdio/HTTP handling

4. **GitHub Integration** (`internal/github/`)
   - Changes needed: Create API client wrapper, authentication
   - Dependencies: GitHub Go client, OAuth handling

5. **Configuration Management** (`internal/config/`)
   - Changes needed: Create config loader, validator
   - Dependencies: Viper, JSON schema validation

6. **AI Provider Integration** (`internal/ai/`)
   - Changes needed: Create provider interface, implementations
   - Dependencies: HTTP clients, API SDKs

## Design Decisions (Requiring Creative Phase)
### Architecture:
- [x] Plugin architecture for AI providers
- [x] Event-driven flow execution
- [x] Middleware pattern for step processing
- [x] Repository pattern for configuration

### Algorithms:
- [x] Flow execution state machine
- [x] Condition evaluation engine
- [x] Tool discovery and caching strategy
- [x] Error recovery mechanisms

### Data Models:
- [x] Flow definition schema design
- [x] MCP server configuration schema
- [x] Execution context structure
- [x] Tool capability representation

## Implementation Strategy
### Phase 0: Environment Setup & Technology Validation
1. [ ] **Install Go on development system**
2. [ ] Verify Go installation: `go version`
3. [ ] Initialize Go module: `go mod init github.com/peterovchinnikov/flow-test-go`
4. [ ] Install core dependencies (Cobra, Viper)
5. [ ] Create minimal CLI with basic command
6. [ ] Verify build and run process
7. [ ] Test basic configuration loading

### Phase 1: Foundation Setup
1. [ ] Project structure creation
   - [ ] Create cmd/ directory for CLI commands
   - [ ] Create internal/ for core logic
   - [ ] Create pkg/ for public APIs
   - [ ] Set up test structure
2. [ ] Basic CLI framework
   - [ ] Root command setup
   - [ ] Version command
   - [ ] Help system
3. [ ] Configuration system
   - [ ] Config file loading
   - [ ] Environment variable support
   - [ ] Default values

### Phase 2: Core Components (After Creative Phase)
1. [ ] Flow engine implementation
   - [ ] Flow parser
   - [ ] Step executor
   - [ ] Condition evaluator
2. [ ] MCP server integration
   - [ ] Server process management
   - [ ] Tool discovery
   - [ ] Communication protocol
3. [ ] AI provider abstraction
   - [ ] Provider interface
   - [ ] OpenAI implementation
   - [ ] Anthropic implementation

### Phase 3: GitHub Integration
1. [ ] GitHub client setup
   - [ ] Authentication handling
   - [ ] API wrapper
2. [ ] GitHub-specific tools
   - [ ] Issue creation
   - [ ] PR management
   - [ ] Repository operations

### Phase 4: Advanced Features
1. [ ] Error handling and recovery
2. [ ] Logging and debugging
3. [ ] Performance optimization
4. [ ] Documentation generation

## Testing Strategy
### Unit Tests:
- [ ] Flow parser tests
- [ ] Condition evaluator tests
- [ ] Configuration loader tests
- [ ] MCP communication tests

### Integration Tests:
- [ ] End-to-end flow execution
- [ ] MCP server interaction
- [ ] GitHub API operations
- [ ] Multi-step flow scenarios

## Documentation Plan
- [ ] API Documentation (godoc)
- [ ] User Guide (README.md)
- [ ] Flow Definition Guide
- [ ] MCP Server Setup Guide
- [ ] Architecture Documentation

## Dependencies
- github.com/spf13/cobra (CLI framework)
- github.com/spf13/viper (Configuration)
- github.com/google/go-github/v57 (GitHub API)
- github.com/stretchr/testify (Testing)
- Standard library packages (encoding/json, net/http, os/exec)

## Challenges & Mitigations
- **Challenge 0: Go Not Installed**: Development environment lacks Go
  - *Mitigation*: Install Go before proceeding with any implementation
- **Challenge 1: MCP Protocol Implementation**: No official Go SDK exists
  - *Mitigation*: Implement custom MCP client based on protocol specification
- **Challenge 2: Managing Multiple Async Operations**: Flows may have parallel steps
  - *Mitigation*: Use Go's goroutines and channels for concurrent execution
- **Challenge 3: Error Recovery in Multi-Step Flows**: Steps may fail mid-execution
  - *Mitigation*: Implement checkpoint/resume mechanism with state persistence
- **Challenge 4: Tool Discovery Caching**: Need efficient tool capability caching
  - *Mitigation*: Implement TTL-based cache with background refresh

## Creative Phases Required
- [x] **Flow Engine Architecture**: Design state machine and execution model
- [x] **MCP Communication Protocol**: Design abstraction for stdio/HTTP transport
- [x] **Plugin System**: Design extensible AI provider interface
- [x] **Configuration Schema**: Design JSON schemas for flows and servers
- [x] **Error Recovery Strategy**: Design fault-tolerant execution system

---

## NEXT ACTION REQUIRED
Since this is a Level 3 task with multiple components requiring design decisions, the next phase should be CREATIVE MODE to design:
1. Core architecture and interfaces
2. Flow execution state machine
3. MCP communication protocol abstraction
4. Configuration schemas
5. Error handling strategies

**NOTE**: Go installation is a prerequisite that must be resolved before the BUILD phase.

Type 'CREATIVE' to proceed with architectural design phase.

---

## ARCHITECTURAL REVISION: LangGraph Integration

After reviewing available tools, we've decided to use **LangGraph** (specifically [langgraphgo](https://github.com/tmc/langgraphgo)) as our core flow orchestration engine instead of building a custom solution.

### Updated Technology Stack
- **Flow Engine**: LangGraph/langgraphgo (replaces custom orchestrator)
- **Language**: Go (1.21+)
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **LLM Integration**: langchaingo (comes with langgraphgo)
- **MCP Communication**: Custom implementation (stdio/HTTP)
- **GitHub Integration**: google/go-github

### Updated Dependencies
- github.com/tmc/langgraphgo (core flow engine)
- github.com/tmc/langchaingo (LLM integrations)
- github.com/spf13/cobra (CLI framework)
- github.com/spf13/viper (Configuration)
- github.com/google/go-github/v57 (GitHub API)
- github.com/stretchr/testify (Testing)

### Benefits of LangGraph
1. **Built-in State Management**: No need for custom state machine
2. **Checkpointing**: Automatic state persistence and recovery
3. **Human-in-the-Loop**: Built-in support for manual intervention
4. **Production Ready**: Battle-tested in production environments
5. **LangSmith Integration**: Advanced debugging and visualization

### Updated Implementation Strategy

#### Phase 0: Technology Validation & Hello World
1. [ ] **Install Go on development system**
2. [ ] Initialize Go module: `go mod init github.com/peterovchinnikov/flow-test-go`
3. [ ] Install langgraphgo: `go get github.com/tmc/langgraphgo`
4. [ ] Create minimal LangGraph example
5. [ ] Verify graph compilation and execution

#### Phase 1: LangGraph Integration
1. [ ] Create graph builder from JSON config
2. [ ] Implement node handlers for different step types
3. [ ] Add conditional edge support
4. [ ] Test basic flow execution with LangGraph

#### Phase 2: Service Integration
1. [ ] MCP Manager integration with LangGraph nodes
2. [ ] AI Provider Manager as LangGraph tools
3. [ ] GitHub integration as graph nodes
4. [ ] Error recovery wrapper around graph execution

#### Phase 3: Advanced Features
1. [ ] Checkpointing with custom storage
2. [ ] Parallel node execution
3. [ ] Human-in-the-loop implementation
4. [ ] LangSmith integration for debugging

### Updated Creative Phases
- [x] **Flow Engine Architecture**: ~~Custom orchestrator~~ → LangGraph integration
- [x] **MCP Communication Protocol**: Unchanged - Protocol Handler Pattern
- [x] **Plugin System**: Unchanged - Full Plugin Architecture
- [x] **Configuration Schema**: Adapted for LangGraph graph building
- [x] **Error Recovery Strategy**: Integrated with LangGraph checkpointing
- [x] **LangGraph Integration**: New - Integration strategy defined

---

---

## ARCHITECTURAL REVISION 2: OpenRouter Integration

We're adopting **OpenRouter** ([go-openrouter](https://github.com/reVrost/go-openrouter)) as our unified LLM provider facade, replacing the need for individual provider integrations.

### Benefits of OpenRouter
1. **Single API**: Access to 200+ models from all major providers
2. **Automatic Failover**: Built-in redundancy across providers
3. **Cost Management**: Unified billing and usage tracking
4. **Model Flexibility**: Easy switching between models
5. **Future Proof**: New models automatically available

### Updated Technology Stack (Final)
- **Flow Engine**: LangGraph/langgraphgo
- **LLM Provider**: OpenRouter (go-openrouter)
- **Language**: Go (1.21+)
- **CLI Framework**: Cobra
- **Configuration**: Viper
- **MCP Communication**: Custom implementation (stdio/HTTP)
- **GitHub Integration**: google/go-github

### Updated Dependencies (Final)
- github.com/tmc/langgraphgo (flow orchestration)
- github.com/tmc/langchaingo (LLM interfaces)
- github.com/revrost/go-openrouter (LLM provider)
- github.com/spf13/cobra (CLI framework)
- github.com/spf13/viper (Configuration)
- github.com/google/go-github/v57 (GitHub API)
- github.com/stretchr/testify (Testing)

### Simplified Architecture
```
┌─────────────────────────────────────────────────┐
│              CLI (Cobra)                         │
├─────────────────────────────────────────────────┤
│         Configuration (JSON/Viper)               │
├─────────────────────────────────────────────────┤
│           LangGraph Engine                       │
├─────────────────────────────────────────────────┤
│  OpenRouter  │  MCP Manager  │  GitHub Client   │
│  (All LLMs)  │   (Tools)     │  (Integration)   │
└─────────────────────────────────────────────────┘
```

### Model Selection Strategy
- **Fast Tasks**: GPT-3.5-Turbo, Claude Instant
- **Balanced**: GPT-4-Turbo, Claude 3.5 Sonnet
- **Complex**: GPT-4, Claude 3.5 Opus, Deepseek V3
- **Coding**: Deepseek V3, GPT-4, Claude Sonnet

### Example Configuration
```json
{
  "llm": {
    "provider": "openrouter",
    "apiKey": "${OPENROUTER_API_KEY}",
    "defaultModel": "openai/gpt-4-turbo",
    "modelOverrides": {
      "code-analysis": "deepseek/deepseek-chat",
      "creative": "anthropic/claude-3.5-sonnet"
    }
  }
}
```

### Updated Implementation Plan

#### Phase 0: Prerequisites
1. [ ] **Install Go on development system**
2. [ ] Get OpenRouter API key from https://openrouter.ai
3. [ ] Set up environment variables

#### Phase 1: Core Setup
1. [ ] Initialize Go module
2. [ ] Install dependencies (langgraphgo, go-openrouter)
3. [ ] Create basic project structure
4. [ ] Implement OpenRouter adapter for LangGraph

#### Phase 2: Integration
1. [ ] Build LangGraph flow from JSON config
2. [ ] Integrate OpenRouter for LLM calls
3. [ ] Add MCP server management
4. [ ] Implement GitHub integration

#### Phase 3: Production Features
1. [ ] Model selection logic
2. [ ] Cost tracking
3. [ ] Error handling with fallbacks
4. [ ] Checkpointing and recovery

---

---

## Code Quality: golangci-lint Integration

We're using [golangci-lint](https://golangci-lint.run) for comprehensive code quality checks. It provides:

### Features
- ⚡ **Fast**: Runs linters in parallel
- 🥇 **Comprehensive**: 100+ linters included
- ⚙️ **Configurable**: YAML-based configuration
- �� **IDE Integration**: Works with VS Code, GoLand, Vim, etc.
- 📈 **Tuned**: Minimal false positives

### Enabled Linters
- **errcheck**: Check error handling
- **govet**: Go vet on steroids
- **gofumpt**: Stricter gofmt
- **gosec**: Security checks
- **staticcheck**: Advanced static analysis
- **revive**: Fast, configurable linter
- **gocritic**: Opinionated code critic
- **ineffassign**: Detect ineffectual assignments
- **misspell**: Spell checker
- And many more...

### Usage
```bash
# Install golangci-lint
make install-tools

# Run linting
make lint

# Run with auto-fix
make lint-fix

# Format code
make fmt

# Pre-commit checks (format, lint, test)
make pre-commit
```

### CI Integration
The project includes:
- `.golangci.yml`: Comprehensive linting configuration
- `Makefile`: Development workflow commands
- CI-friendly commands: `make ci`

### Development Workflow
1. Write code
2. Run `make fmt` to format
3. Run `make lint` to check for issues
4. Run `make test` to verify functionality
5. Run `make pre-commit` before committing

---

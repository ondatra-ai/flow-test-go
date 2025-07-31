# System Patterns

## Language & Framework
- **Primary Language**: Go (1.21+)
- **CLI Framework**: Cobra (command structure and help)
- **Configuration**: Viper (config management)
- **Build System**: Go modules
- **Testing Framework**: Go testing + testify
- **HTTP Client**: Go standard library + google/go-github

## Architecture Patterns
- **Pattern**: Command-line application with LangGraph orchestration
- **Structure**: Clean architecture with clear separation of concerns
- **LLM Integration**: OpenRouter as unified provider
- **Flow Engine**: LangGraph for state management and execution
- **Error Handling**: Checkpoint/resume with state persistence

## Integration Patterns
- **GitHub Integration**: google/go-github client library
- **MCP Integration**: Custom implementation (stdio/HTTP transport)
- **LLM Access**: OpenRouter API (single integration for all models)
- **Configuration**: JSON-based flow and server definitions
- **Tool Discovery**: TTL-based caching mechanism

## Code Organization
```
flow-test-go/
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command
│   ├── execute.go         # Execute flow command
│   ├── list.go            # List flows command
│   └── validate.go        # Validate flow command
├── internal/              # Internal packages
│   ├── graph/             # LangGraph integration
│   │   ├── builder.go     # Graph builder from config
│   │   ├── nodes.go       # Node handler implementations
│   │   ├── state.go       # State schema definitions
│   │   └── adapter.go     # Config to graph adapter
│   ├── llm/               # LLM integration (OpenRouter)
│   │   ├── client.go      # OpenRouter client wrapper
│   │   ├── adapter.go     # LangChain interface adapter
│   │   ├── models.go      # Model selection logic
│   │   └── cost.go        # Cost tracking
│   ├── mcp/               # MCP integration
│   │   ├── manager.go     # MCP server manager
│   │   ├── protocol.go    # Protocol handlers
│   │   └── tools.go       # Tool discovery/execution
│   ├── github/            # GitHub integration
│   │   ├── client.go      # API client wrapper
│   │   └── nodes.go       # GitHub-specific nodes
│   └── config/            # Configuration management
│       ├── loader.go      # Config loader
│       └── validator.go   # Schema validator
├── pkg/                   # Public packages
│   ├── types/             # Shared types
│   └── utils/             # Utilities
├── .flows/                # Configuration directory
│   ├── flows/             # Flow definitions
│   └── servers/           # MCP server configs
├── tests/                 # Test files
│   ├── unit/              # Unit tests
│   └── integration/       # Integration tests
└── examples/              # Example flows
    ├── simple-chat/       # Basic chat example
    ├── code-review/       # GitHub code review flow
    └── multi-agent/       # Multi-agent example
```

## Standards
### Error Handling
- Use custom error types for different scenarios
- Implement error wrapping with context
- Support checkpoint/resume for long-running flows
- Provide clear user-facing error messages

### Logging Standards
- Structured logging with levels (debug, info, warn, error)
- Context propagation through operations
- Separate user output from debug logs
- Log rotation for long-running processes

### Configuration Management
- Environment variable overrides
- Default configuration values
- Schema validation for all configs
- Hot reload support for server configs

### Testing Standards
- Table-driven tests for comprehensive coverage
- Mock interfaces for external dependencies
- Integration tests with real MCP servers
- E2E tests for complete flows

## Design Patterns
- **Adapter Pattern**: For OpenRouter/LangChain integration
- **Factory Pattern**: For node handler creation
- **Observer Pattern**: For flow execution events
- **Strategy Pattern**: For transport mechanisms
- **Middleware Pattern**: For step processing

## Performance Considerations
- Concurrent MCP server execution
- Tool discovery result caching
- Efficient state persistence
- Minimal memory footprint
- Fast startup time

## Security Considerations
- Secure credential storage
- API key management
- Input validation
- Safe command execution
- Audit logging

## Architectural Decisions (From Creative Phase)

### Flow Engine (LangGraph)
- **Core**: LangGraph/langgraphgo for orchestration
- **State Management**: LangGraph's built-in state graph
- **Execution**: LangGraph's compiled graph runtime
- **Checkpointing**: LangGraph's checkpoint system
- **Benefits**: Production-ready, no custom state machine needed

### LLM Integration (OpenRouter)
- **Provider**: OpenRouter as unified LLM facade
- **Models**: Access to 200+ models through single API
- **Adapter**: OpenRouter client wrapped for LangChain interface
- **Selection**: Dynamic model selection based on task type
- **Benefits**: Automatic failover, cost tracking, future-proof

### MCP Communication Protocol
- **Pattern**: Protocol Handler Pattern with Connection Pool
- **Transports**: Unified interface for stdio and HTTP
- **Features**: Health checking, connection pooling, tool discovery caching
- **Process Management**: Lifecycle management for stdio processes
- **Error Handling**: Transport-specific retry strategies

### Configuration Schema
- **Format**: JSON with schema validation
- **Features**: Variable substitution, imports, templates
- **Structure**: Modular component-based design
- **Validation**: JSON Schema v7 with semantic validation
- **Extensibility**: Support for custom step types and handlers

### Error Recovery Strategy
- **Pattern**: Comprehensive Fault Tolerance System
- **Error Classification**: Transient, RateLimit, Invalid, Fatal
- **Retry Strategies**: Exponential backoff, adaptive retry
- **Circuit Breaker**: Prevents cascade failures
- **State Persistence**: Versioned checkpoints with compression
- **Manual Intervention**: Support for fatal error recovery

## Final Architecture

### Integration Architecture
```
┌─────────────────────────────────────────────────┐
│                 CLI Interface                    │
│                   (Cobra)                        │
├─────────────────────────────────────────────────┤
│            Configuration Loader                  │
│         (JSON → LangGraph Adapter)               │
├─────────────────────────────────────────────────┤
│              LangGraph Engine                    │
│         (State Graph + Checkpointing)            │
├─────────────────────────────────────────────────┤
│     Node Handlers    │    Error Recovery         │
│  (Step Executors)    │  (Checkpoint/Retry)       │
├─────────────────────────────────────────────────┤
│  OpenRouter  │  MCP Manager  │  GitHub Client   │
│  (All LLMs)  │   (Tools)     │  (Integration)   │
└─────────────────────────────────────────────────┘
```

### Key Integration Points
1. **Graph Building**: Convert JSON config to LangGraph nodes/edges
2. **LLM Integration**: OpenRouter client implements LangChain interface
3. **Node Handlers**: Wrap services in LangGraph-compatible functions
4. **State Schema**: Define state structure for message passing
5. **Tool Integration**: Convert MCP tools to LangGraph tools

### Model Strategy (via OpenRouter)
- **Fast**: GPT-3.5-Turbo, Claude Instant, Gemini Flash
- **Balanced**: GPT-4-Turbo, Claude 3.5 Sonnet, Gemini Pro
- **Powerful**: GPT-4, Claude 3.5 Opus, Deepseek V3
- **Coding**: Deepseek V3, GPT-4, Claude Sonnet
- **Vision**: GPT-4 Vision, Claude 3.5 with vision

## Code Quality Standards (golangci-lint)

### Linting Configuration
- **Tool**: golangci-lint (https://golangci-lint.run)
- **Config**: .golangci.yml with comprehensive rules
- **Timeout**: 5 minutes for full analysis
- **Line Length**: 120 characters maximum

### Enabled Linters
- **errcheck**: Ensures all errors are handled
- **govet**: Reports suspicious constructs
- **gofumpt**: Enforces stricter formatting than gofmt
- **gosec**: Security-focused static analysis
- **staticcheck**: State-of-the-art static analysis
- **revive**: Fast, extensible linter
- **gocritic**: Opinionated linter for bugs, performance, style
- **ineffassign**: Detects ineffectual assignments
- **misspell**: Spell checker for comments and strings
- **funlen**: Function length limits (100 lines, 50 statements)
- **gocyclo**: Cyclomatic complexity (max 15)
- **dupl**: Code duplication detection (threshold 100)

### Development Workflow
1. **Pre-commit**: Run `make pre-commit` (fmt, lint, test)
2. **Linting**: `make lint` to check, `make lint-fix` to auto-fix
3. **Formatting**: `make fmt` applies gofumpt + goimports + golines
4. **CI Pipeline**: Automated checks on every PR

### Code Style Guidelines
- Use meaningful variable and function names
- Keep functions focused and under 100 lines
- Handle all errors explicitly
- Add comments for exported functions
- Use table-driven tests
- Avoid global variables where possible
- Prefer composition over inheritance

### CI/CD Integration
- GitHub Actions workflow with golangci-lint
- Multi-OS testing (Linux, macOS, Windows)
- Security scanning with gosec and trivy
- Coverage reporting with codecov
- Artifact generation for releases

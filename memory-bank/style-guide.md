# Style Guide

## Overview
This style guide provides a set of conventions and best practices for writing Go code in the `flow-test-go` project. Its purpose is to ensure code consistency, readability, and maintainability. The guidelines are derived from both the project's `.golangci.yml` linter configuration and established patterns observed within the existing codebase. Adhering to this guide will help streamline development, reduce errors, and make the codebase easier to navigate for all contributors.

## Codebase Analysis Summary
The existing codebase follows standard Go idioms and demonstrates consistent patterns.

- **Package Organization**: The project is structured with `cmd` for entry points, `internal` for core application logic, `pkg` for shared types, and `tests` for testing. Package names are lowercase and concise.
- **Naming Conventions**: `PascalCase` is used for exported identifiers like structs and functions, while `camelCase` is used for internal variables. Test functions follow the `TestFeature_Scenario` pattern.
- **Error Handling**: Errors are handled explicitly using `if err != nil` checks. Errors are wrapped with context using `fmt.Errorf("...: %w", err)`. The project defines both package-level sentinel errors and custom error structs (`ExecutionError`) for more detailed error reporting.
- **Structs and Types**: Structs are well-defined in the `pkg/types` directory. Serialization is handled using struct tags (`json`, `yaml`, `mapstructure`), which follow a consistent `camelCase` format.
- **Testing**: End-to-end tests are located in the `tests/e2e` directory and leverage the `testify` library for assertions (`assert`, `require`). Tests are run in parallel (`t.Parallel()`) and use a `testutil` package for helper functions.

## 1. Code Quality Standards
These guidelines ensure the code is robust, clean, and free of common pitfalls.

### **Avoid Unused Code**
**Rule**: All declared variables, functions, and imports must be used. The compiler often catches this, but the `unused` and `ineffassign` linters provide stricter checks.

- **Rationale**: Unused code clutters the namespace, increases cognitive load, and can hide bugs.

### **Use Constants for Repeated Values**
**Rule**: If a value (string, number, etc.) is used more than once, declare it as a constant. The `goconst` linter helps enforce this.

- **Rationale**: Constants improve readability and make it easier to update values without searching the entire codebase.

- **Correct**:
  ```go
  const defaultTimeout = 5 * time.Minute
  cfg.Timeout = defaultTimeout
  ```

- **Incorrect**:
  ```go
  cfg.Timeout = 5 * time.Minute // Magic number
  ```

### **No Naked Returns**
**Rule**: Always explicitly return values from functions. Naked returns are forbidden by the `nakedret` linter.

- **Rationale**: Explicit returns make code clearer and less error-prone, as the returned values are always visible at the `return` statement.

- **Correct**:
  ```go
  func (cm *Manager) GetConfig() *Config {
      return cm.config
  }
  ```

- **Incorrect**:
  ```go
  func (cm *Manager) GetConfig() (cfg *Config) {
      cfg = cm.config
      return // naked return
  }
  ```

## 2. Error Handling Guidelines
Consistent error handling is critical for application stability and debuggability.

### **Check Every Error**
**Rule**: Always check the error returned by a function call. This is enforced by `errcheck`.

- **Rationale**: Unchecked errors can lead to unexpected behavior, panics, and data corruption.

- **Correct**:
  ```go
  err := viper.ReadInConfig()
  if err != nil {
      return nil, fmt.Errorf("failed to read config file: %w", err)
  }
  ```

### **Wrap Errors with Context**
**Rule**: When propagating an error, wrap it with context using `fmt.Errorf` and the `%w` verb. This is enforced by `err113` and `wrapcheck`.

- **Rationale**: Wrapping preserves the original error while adding context, making it easier to trace the error's origin.

- **Correct**:
  ```go
  state.configMgr, err = config.NewManager()
  if err != nil {
      return fmt.Errorf("failed to initialize config manager: %w", err)
  }
  ```

### **Use a Custom Error Struct for Rich Errors**
**Rule**: For application-specific errors that require more metadata, define a custom error struct.

- **Rationale**: Custom error types allow for structured error data, including codes, messages, and recoverability status, which is useful for API responses or complex logic.

- **Example (`pkg/types/flow.go`)**:
  ```go
  type ExecutionError struct {
      Code        string    `json:"code"`
      Message     string    `json:"message"`
      Recoverable bool      `json:"recoverable"`
      // ...
  }

  func (e *ExecutionError) Error() string {
      return e.Message
  }
  ```

## 3. Code Structure and Organization
A well-organized codebase is easier to understand and maintain.

### **Group Imports**
**Rule**: Group imports into three blocks, separated by newlines: standard library, third-party packages, and internal project packages. This is enforced by `grouper`.

- **Example (`cmd/commands/root.go`)**:
  ```go
  import (
      "fmt"
      "os"
      "sync"

      "github.com/spf13/cobra"

      "github.com/ondatra-ai/flow-test-go/internal/config"
  )
  ```

### **Limit Function Length and Complexity**
**Rule**: Keep functions short and focused on a single responsibility. The `funlen` (max 130 lines) and `gocyclo` (max 20 complexity) linters help enforce this.

- **Rationale**: Smaller functions are easier to test, debug, and understand. Test files have relaxed rules for these linters.

### **Control Nesting Depth**
**Rule**: Avoid deep nesting of `if`, `for`, and `switch` blocks. The `nestif` linter flags excessive nesting.

- **Rationale**: Deeply nested code is hard to read. Use guard clauses or break the logic into smaller functions to reduce nesting.

- **Correct (Guard Clause)**:
  ```go
  if f.ID == "" {
      return &ExecutionError{Message: "flow ID is required"}
  }
  // ... rest of the function
  ```

## 4. Performance Guidelines
Write efficient code to ensure the application is responsive.

### **Pre-allocate Slices**
**Rule**: When the size of a slice is known beforehand, pre-allocate it using `make` with a specified capacity. This is enforced by `prealloc`.

- **Rationale**: Pre-allocation avoids repeated memory allocations and copies as the slice grows, improving performance.

- **Correct**:
  ```go
  servers := make(map[string]*types.MCPServerConfig, len(files))
  ```

### **Use `fmt.Sprintf` Efficiently**
**Rule**: Avoid unnecessary calls to `fmt.Sprintf`. The `perfsprint` linter flags inefficient usage.

- **Rationale**: Direct string conversions or other formatting methods can be faster than `Sprintf` in many cases.

## 5. Security Standards
Writing secure code is a top priority.

### **Address Security Vulnerabilities**
**Rule**: The `gosec` linter is enabled to detect common security issues. Pay attention to its warnings.

- **Rationale**: Proactively identifying and fixing security flaws prevents vulnerabilities.

- **Annotation**: If a `gosec` warning is a false positive or intentionally accepted, annotate the line with `#nosec GXXX` and a justification.
  ```go
  data, err := os.ReadFile(flowPath) // #nosec G304 (path is validated before use)
  ```

### **Avoid Forbidden Constructs**
**Rule**: Do not use functions or patterns that have been marked as forbidden by the `forbidigo` linter.

- **Rationale**: This helps phase out deprecated or unsafe APIs and enforce project-specific constraints.

## 6. Testing Standards
High-quality tests are essential for a reliable application.

### **Use `testify` for Assertions**
**Rule**: Use `github.com/stretchr/testify/assert` for non-fatal checks and `github.com/stretchr/testify/require` for fatal checks that should stop the test on failure. `testifylint` is enabled.

- **Correct**:
  ```go
  require.Equal(t, 0, result.ExitCode, "Command should complete successfully")
  assert.Contains(t, result.Stderr, "No flows found")
  ```

### **Run Tests in Parallel**
**Rule**: Enable parallel execution for tests by calling `t.Parallel()` at the beginning of each test function. `tparallel` enforces this.

- **Rationale**: Parallel tests run faster, especially as the test suite grows.

- **Example**:
  ```go
  func TestListCommand_EmptyDirectory(t *testing.T) {
      t.Parallel()
      // ... test logic
  }
  ```

### **Organize Tests with Helpers**
**Rule**: Use a `testutil` package for shared test setup, teardown, and helper functions. The `thelper` linter ensures test helpers are correctly marked.

- **Rationale**: Test helpers reduce duplication and make tests cleaner and more focused on their specific scenario.

## 7. Documentation Requirements
Code should be well-documented to be understandable.

### **Punctuate Comments**
**Rule**: All sentence-like comments must end with a period. This is enforced by `godot`.

- **Rationale**: Consistent punctuation improves readability.

- **Correct**:
  ```go
  // Package config provides configuration management.
  package config
  ```

### **Manage TODOs**
**Rule**: Use `TODO:` or `FIXME:` prefixes for comments that highlight work to be done. The `godox` linter tracks these.

- **Rationale**: This creates a centralized way to track technical debt and outstanding tasks.

## 8. Naming Conventions
Consistent naming makes code predictable and easy to read.

### **Use `camelCase` for Struct Tags**
**Rule**: Struct tags for JSON, YAML, and other formats should use `camelCase`. `tagliatelle` enforces this.

- **Rationale**: This convention is common in Go and aligns with JavaScript/JSON standards.

- **Example (`pkg/types/flow.go`)**:
  ```go
  type FlowDefinition struct {
      InitialStep string `json:"initialStep,omitempty" yaml:"initialStep,omitempty"`
  }
  ```

### **Align Struct Tags**
**Rule**: Align struct tags in vertical columns for better readability. `tagalign` helps with this.

- **Example (`pkg/types/mcp.go`)**:
  ```go
  type MCPServerConfig struct {
      Name    string   `json:"name"    yaml:"name"`
      Command string   `json:"command" yaml:"command"`
      Args    []string `json:"args"    yaml:"args"`
  }
  ```

## Linter Configuration Reference
- **Disabled Linters**: Linters like `cyclop`, `depguard`, and `revive` are commented out in `.golangci.yml`. This indicates they may be too opinionated, require specific configuration (like dependency rules for `depguard`), or are not currently a priority. They can be enabled later if needed.
- **Test Exclusions**: Rules for `gocyclo`, `errcheck`, and `funlen` are relaxed for `_test.go` files. This is a practical choice, as test files can sometimes require longer setup functions and repeated error checks that don't need to be propagated.
- **WSL (Whitespace Linter)**: The `wsl` configuration is customized to allow newlines at the start of a block but enforces them around multi-line branch statements, promoting a clean visual separation of logic blocks.

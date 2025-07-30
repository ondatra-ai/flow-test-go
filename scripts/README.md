# Scripts Directory

This directory contains utility scripts for managing GitHub Pull Requests and repository operations.

## Available Scripts

### 1. Get PR Number (`get-pr-number.go`)

Reveals the Pull Request number associated with the current git branch.

**Usage:**
```bash
# Direct execution (like a shell script)
./scripts/get-pr-number.go

# Alternative: go run
go run scripts/get-pr-number.go

# Using helper script
./scripts/run-script.sh get-pr-number

# Build and run
go build -o get-pr-number scripts/get-pr-number.go && ./get-pr-number
```

**Requirements:**
- GitHub CLI (gh) must be installed: https://cli.github.com/
- Must be run from within a git repository
- Repository must be hosted on GitHub

**Output:**
- Shows PR number, title, state, and URL
- If no open PR is found, shows closed/merged PRs
- Displays just the PR number at the end for easy scripting

### 2. List PR Conversations (`list-pr-conversations.go`)

Lists all unresolved review conversations for a specific PR.

**Usage:**
```bash
# Direct execution
./scripts/list-pr-conversations.go <PR_NUMBER>

# Alternative: go run
go run scripts/list-pr-conversations.go <PR_NUMBER>

# Using helper script
./scripts/run-script.sh list-pr-conversations 123

# Example
./scripts/list-pr-conversations.go 42
```

**Output:**
- JSON format listing all unresolved conversations
- Each conversation includes comments, authors, file locations, and metadata
- Sorted by creation date of first comment

### 3. Resolve PR Conversation (`resolve-pr-conversation.go`)

Resolves a GitHub review thread (conversation) and optionally adds a comment before resolving.

**Usage:**
```bash
# Direct execution
./scripts/resolve-pr-conversation.go <THREAD_ID> [COMMENT]

# Alternative: go run
go run scripts/resolve-pr-conversation.go <THREAD_ID> [COMMENT]

# Using helper script
./scripts/run-script.sh resolve-pr-conversation <THREAD_ID> "Fixed the issue"

# Examples
./scripts/resolve-pr-conversation.go MDExOlB1bGxSZXF1ZXN0UmV2aWV3VGhyZWFkMzg0Nzc2
./scripts/resolve-pr-conversation.go MDExOlB1bGxSZXF1ZXN0UmV2aWV3VGhyZWFkMzg0Nzc2 "Thanks for the feedback, fixed!"
```

**Parameters:**
- `THREAD_ID`: The GitHub review thread ID (from list-pr-conversations output)
- `COMMENT` (optional): Comment to add before resolving the conversation

## Helper Script

The `run-script.sh` helper script provides a convenient way to run any of the Go scripts:

```bash
# Make it executable (one time)
chmod +x scripts/run-script.sh

# Run any script
./scripts/run-script.sh <script-name> [args...]

# See available scripts
./scripts/run-script.sh
```

## Executable Go Scripts

All Go scripts in this directory use a special shebang line that allows them to be executed directly like shell scripts:

```go
//usr/bin/env go run "$0" "$@"; exit
```

This approach provides several benefits:
- **Direct execution**: Run Go scripts like `./script.go` without explicit `go run`
- **Shell-like behavior**: Pass arguments naturally
- **No compilation needed**: Scripts are interpreted on-the-fly
- **Familiar usage**: Works like traditional shell scripts

## Makefile Integration

The project Makefile includes convenient targets for script management:

```bash
# Build all scripts as binaries
make build-scripts

# Clean built binaries
make clean-scripts

# Quick execution targets
make pr-number
make pr-conversations PR=123
make resolve-conversation ID=<thread-id> COMMENT="message"
```

## Migration from TypeScript

These Go scripts replace the original TypeScript scripts with improved functionality:

| Original TypeScript | Go Replacement | Status |
|-------------------|----------------|--------|
| `get-pr-number.ts` | `get-pr-number.go` | ✅ Converted & Executable |
| `list-pr-conversations.ts` | `list-pr-conversations.go` | ✅ Converted & Executable |
| `resolve-pr-conversation.ts` | `resolve-pr-conversation.go` | ✅ Converted & Executable |

### Key Improvements

1. **Executable Scripts**: Can be run directly as `./script.go`
2. **No Dependencies**: No need for Node.js or npm packages
3. **Better Performance**: Compiled Go code runs faster
4. **Type Safety**: Compile-time type checking
5. **Single Binary**: Can be compiled to standalone executables
6. **Consistent Tooling**: Uses the same language as the main project

## Dependencies

All scripts require:
- Go 1.21 or later
- GitHub CLI (gh) installed and authenticated
- Git repository with GitHub remote

## Building Standalone Binaries

You can build standalone executables for easier distribution:

```bash
# Build all scripts
make build-scripts

# Or build individually
go build -o bin/get-pr-number scripts/get-pr-number.go
go build -o bin/list-pr-conversations scripts/list-pr-conversations.go
go build -o bin/resolve-pr-conversation scripts/resolve-pr-conversation.go

# Run the binaries
./bin/get-pr-number
./bin/list-pr-conversations 123
./bin/resolve-pr-conversation <thread-id>
```

## Error Handling

All scripts provide meaningful error messages for common issues:
- GitHub CLI not installed
- Not in a GitHub repository
- Invalid PR numbers or thread IDs
- Network connectivity issues
- Authentication problems

Exit codes:
- `0`: Success
- `1`: Error occurred

## Examples

```bash
# Get current branch PR number
./scripts/get-pr-number.go

# List conversations for PR #42
./scripts/list-pr-conversations.go 42

# Resolve a conversation with a comment
./scripts/resolve-pr-conversation.go MDExOlB1bGxSZXF1ZXN0UmV2aWV3VGhyZWFkMzg0Nzc2 "Fixed the issue!"
```
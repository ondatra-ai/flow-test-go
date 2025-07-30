# Scripts Directory

This directory contains utility scripts for managing GitHub Pull Requests and repository operations.

## Available Scripts

### 1. Get PR Number (`get-pr-number.go`)

Reveals the Pull Request number associated with the current git branch.

**Usage:**
```bash
# Direct execution
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
go run scripts/list-pr-conversations.go <PR_NUMBER>

# Using helper script
./scripts/run-script.sh list-pr-conversations 123

# Example
go run scripts/list-pr-conversations.go 42
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
go run scripts/resolve-pr-conversation.go <THREAD_ID> [COMMENT]

# Using helper script
./scripts/run-script.sh resolve-pr-conversation <THREAD_ID> "Fixed the issue"

# Examples
go run scripts/resolve-pr-conversation.go MDExOlB1bGxSZXF1ZXN0UmV2aWV3VGhyZWFkMzg0Nzc2
go run scripts/resolve-pr-conversation.go MDExOlB1bGxSZXF1ZXN0UmV2aWV3VGhyZWFkMzg0Nzc2 "Thanks for the feedback, fixed!"
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

## Migration from TypeScript

These Go scripts are direct ports of the original TypeScript scripts:

| TypeScript | Go | Status |
|------------|----| -------|
| `get-pr-number.ts` | `get-pr-number.go` | ✅ Converted |
| `list-pr-conversations.ts` | `list-pr-conversations.go` | ✅ Converted |
| `resolve-pr-conversation.ts` | `resolve-pr-conversation.go` | ✅ Converted |

### Key Differences

1. **Build Tags**: Go scripts use `//go:build ignore` to prevent them from being included in regular builds
2. **Error Handling**: Go uses explicit error handling instead of try/catch
3. **JSON Parsing**: Go uses struct tags for JSON marshaling/unmarshaling
4. **Execution**: Scripts can be run with `go run` or compiled to binaries

### Benefits of Go Version

- **No Dependencies**: No need for Node.js or npm packages
- **Better Performance**: Compiled Go code runs faster
- **Type Safety**: Compile-time type checking
- **Single Binary**: Can be compiled to standalone executables
- **Consistent Tooling**: Uses the same language as the main project

## Dependencies

All scripts require:
- Go 1.21 or later
- GitHub CLI (gh) installed and authenticated
- Git repository with GitHub remote

## Building Standalone Binaries

You can build standalone executables for easier distribution:

```bash
# Build all scripts
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
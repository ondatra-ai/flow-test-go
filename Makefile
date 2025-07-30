# Makefile for flow-test-go

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing other tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/segmentio/golines@latest
	@echo "Tools installed!"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with --fix
	@echo "Running golangci-lint with auto-fix..."
	@golangci-lint run --fix

.PHONY: lint-verbose
lint-verbose: ## Run golangci-lint with verbose output
	@echo "Running golangci-lint (verbose)..."
	@golangci-lint run -v

.PHONY: fmt
fmt: ## Format code with gofumpt
	@echo "Formatting code..."
	@gofumpt -l -w .
	@goimports -local github.com/peterovchinnikov/flow-test-go -w .
	@golines -w --max-len=120 --base-formatter=gofumpt .

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: build
build: ## Build the application
	@echo "Building flow-test-go..."
	@go build -o bin/flow-test-go ./cmd/flow-test-go

.PHONY: run
run: ## Run the application
	@echo "Running flow-test-go..."
	@go run ./cmd/flow-test-go

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@echo "Installing core dependencies..."
	@go get github.com/tmc/langgraphgo@latest
	@go get github.com/reVrost/go-openrouter@latest
	@go get github.com/spf13/cobra@latest
	@go get github.com/spf13/viper@latest
	@go get github.com/google/go-github/v57@latest
	@go mod download
	@go mod tidy

.PHONY: verify
verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

.PHONY: update
update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

.PHONY: ci
ci: deps lint test ## Run CI pipeline (deps, lint, test)
	@echo "CI pipeline completed!"

.PHONY: pre-commit
pre-commit: fmt lint test ## Run pre-commit checks
	@echo "Pre-commit checks passed!"

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t flow-test-go:latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -it flow-test-go:latest

# Development helpers
.PHONY: todo
todo: ## Show all TODO/FIXME comments
	@echo "TODO/FIXME comments:"
	@grep -r "TODO\|FIXME" --exclude-dir=.git --exclude-dir=vendor --exclude=Makefile .

.PHONY: check-go-version
check-go-version: ## Check Go version
	@echo "Current Go version:"
	@go version
	@echo ""
	@echo "Required: Go 1.21+"

# Coverage targets (excluding main.go package)
.PHONY: coverage coverage-html coverage-func coverage-full coverage-report

# Coverage targets (excluding main.go package)
coverage:
	@echo "🧪 Running tests with coverage (excluding main.go)..."
	go test ./pkg/... ./internal/... ./cmd/commands/... -coverprofile=coverage.out
	@echo "📊 Coverage Summary:"
	@go tool cover -func=coverage.out | tail -1

# Generate HTML coverage report (excluding main.go)
coverage-html: coverage
	@echo "🌐 Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"
	@echo "💡 Open with: open coverage.html"

# Show function-level coverage details (excluding main.go)
coverage-func: coverage
	@echo "📋 Function-level coverage details:"
	go tool cover -func=coverage.out

# Full coverage including main.go (for comparison)
coverage-full:
	@echo "🧪 Running tests with FULL coverage (including main.go)..."
	go test ./... -coverprofile=coverage_full.out
	@echo "📊 Full Coverage Summary:"
	go tool cover -func=coverage_full.out | tail -1

# Comprehensive coverage report
coverage-report: coverage
	@echo ""
	@echo "📊 COVERAGE REPORT (main.go excluded)"
	@echo "===================================="
	@echo ""
	@echo "🎯 Overall Coverage:"
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "📦 Package Breakdown:"
	@go test ./pkg/... ./internal/... ./cmd/commands/... -cover 2>/dev/null | grep "coverage:" || echo "No coverage data"
	@echo ""
	@echo "📁 Files Generated:"
	@echo "• coverage.out - Raw coverage data (main.go excluded)"
	@echo "• coverage.html - HTML report (run 'make coverage-html')"
	@echo ""
	@echo "💡 To include main.go in coverage: make coverage-full"

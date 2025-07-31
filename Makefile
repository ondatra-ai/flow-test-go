# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=flow-test-go
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/flow-test-go

# Build for linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v ./cmd/flow-test-go

# Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
.PHONY: test-race
test-race:
	$(GOTEST) -race -short ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Tidy modules
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Install the application
.PHONY: install
install:
	$(GOCMD) install ./cmd/flow-test-go

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  build-linux  - Build for Linux"
	@echo "  clean        - Clean build files"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage"
	@echo "  test-race    - Run tests with race detection"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  tidy         - Tidy modules"
	@echo "  deps         - Download dependencies"
	@echo "  install      - Install the application"
	@echo "  help         - Show this help"

.DEFAULT_GOAL := help
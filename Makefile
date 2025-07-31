# Build the application
.PHONY: build
build:
	@mkdir -p bin
	go build -o bin/flow-test-go ./cmd/flow-test-go

# Run tests with coverage
.PHONY: coverage
coverage:
	go test -v -coverprofile=coverage.out ./...

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build files
.PHONY: clean
clean:
	rm -rf bin/
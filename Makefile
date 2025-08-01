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

# Build with coverage instrumentation for e2e tests
.PHONY: build-e2e-coverage
build-e2e-coverage:
	@mkdir -p bin
	go build -cover -covermode=atomic -o bin/flow-test-go-e2e ./cmd/flow-test-go

# Run e2e tests
.PHONY: test-e2e
test-e2e: build-e2e-coverage
	@mkdir -p coverage/e2e
	go test -v -timeout 5m ./tests/e2e/...

# Generate e2e coverage report
.PHONY: coverage-e2e-report
coverage-e2e-report:
	@echo "Aggregating e2e coverage data..."
	@if [ -d coverage/e2e ] && [ "$$(find coverage/e2e -name 'covcounters*' -o -name 'covmeta*' | wc -l)" -gt 0 ]; then \
		go tool covdata merge -i=coverage/e2e/* -o=coverage/e2e-merged && \
		go tool covdata textfmt -i=coverage/e2e-merged -o=coverage/e2e.out && \
		go tool cover -html=coverage/e2e.out -o=coverage/e2e.html && \
		go tool cover -func=coverage/e2e.out > coverage/e2e-summary.txt && \
		echo "E2E coverage report generated"; \
	else \
		echo "No e2e coverage data found"; \
	fi

# Run e2e tests with coverage report
.PHONY: test-e2e-coverage
test-e2e-coverage: test-e2e coverage-e2e-report

# Clean build files
.PHONY: clean
clean:
	rm -rf bin/ coverage/

# KongDB Makefile
# Provides convenient commands for building, testing, and managing KongDB

.PHONY: help build test test-coverage test-race clean install lint format docs examples

# Default target
help:
	@echo "KongDB - A relational database in Go"
	@echo ""
	@echo "Available commands:"
	@echo "  build         - Build KongDB binaries"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detector"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install dependencies"
	@echo "  lint          - Run linting tools"
	@echo "  format        - Format code with gofmt"
	@echo "  docs          - Generate documentation"
	@echo "  examples      - Run example programs"
	@echo "  benchmark     - Run performance benchmarks"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run KongDB in Docker"

# Build KongDB binaries
build:
	@echo "Building KongDB..."
	@mkdir -p bin
	go build -o bin/kongdb ./cmd/kongdb
	go build -o bin/kongdb-node ./cmd/kongdb-node
	@echo "Build complete! Binaries available in bin/"

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./internal/...
	go test -v ./test/integration/...
	go test -v ./test/performance/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./internal/...
	go test -coverprofile=coverage.integration.out ./test/integration/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -race ./internal/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/
	rm -rf dist/
	rm -f coverage.out
	rm -f coverage.html
	rm -f coverage.integration.out
	@echo "Clean complete!"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Run linting tools
lint:
	@echo "Running linting tools..."
	go vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code with gofmt
format:
	@echo "Formatting code..."
	gofmt -w -s ./internal/
	gofmt -w -s ./cmd/
	gofmt -w -s ./pkg/
	gofmt -w -s ./test/
	@echo "Code formatting complete!"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server on http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Run example programs
examples:
	@echo "Running examples..."
	@if [ -d "./examples" ]; then \
		go run ./examples/basic_usage.go; \
	else \
		echo "No examples directory found"; \
	fi

# Run performance benchmarks
benchmark:
	@echo "Running performance benchmarks..."
	go test -bench=. -benchmem ./test/performance/...
	go test -bench=. -benchmem ./internal/storage/...
	go test -bench=. -benchmem ./internal/storage/bloom_filter/...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t kongdb:latest .
	@echo "Docker image built: kongdb:latest"

# Run KongDB in Docker
docker-run:
	@echo "Running KongDB in Docker..."
	docker run -p 8080:8080 -v $(PWD)/data:/app/data kongdb:latest

# Development helpers
dev-setup: install format lint
	@echo "Development environment setup complete!"

# Quick test (unit tests only)
test-quick:
	@echo "Running quick tests (unit tests only)..."
	go test -v ./internal/...

# Integration test only
test-integration:
	@echo "Running integration tests..."
	go test -v ./test/integration/...

# Performance test only
test-performance:
	@echo "Running performance tests..."
	go test -v ./test/performance/...

# Check for common issues
check:
	@echo "Running comprehensive checks..."
	@echo "1. Formatting..."
	@if [ "$$(gofmt -l ./internal/ ./cmd/ ./pkg/ ./test/ | wc -l)" -gt 0 ]; then \
		echo "Code formatting issues found. Run 'make format' to fix."; \
		exit 1; \
	fi
	@echo "2. Linting..."
	@go vet ./...
	@echo "3. Tests..."
	@go test -v ./internal/...
	@echo "All checks passed!"

# Release preparation
release-prep: clean build test-coverage lint check
	@echo "Release preparation complete!"
	@echo "Ready for release!"

# Show project statistics
stats:
	@echo "KongDB Project Statistics:"
	@echo "Lines of Go code: $$(find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1)"
	@echo "Number of Go files: $$(find . -name "*.go" -not -path "./vendor/*" | wc -l)"
	@echo "Number of test files: $$(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)"
	@echo "Number of packages: $$(go list ./... | wc -l)" 
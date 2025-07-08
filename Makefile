.PHONY: help build test lint fmt clean install-tools migrate dev examples

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the main application
	go build -o bin/go-better-auth ./cmd/...

build-examples: ## Build example applications
	cd cmd/gin-example && go build -o ../../bin/gin-example .
	cd cmd/fiber-example && go build -o ../../bin/fiber-example .

# Test targets
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

# Development
dev: fmt lint test ## Run full development cycle (format, lint, test)

# Database
migrate: ## Run database migrations
	./scripts/migrate.sh

# Dependencies
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Examples
examples: build-examples ## Build and run examples
	@echo "Examples built successfully!"
	@echo "Run: ./bin/gin-example or ./bin/fiber-example"

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Dependencies
tidy: ## Tidy go modules
	go mod tidy
	cd cmd/gin-example && go mod tidy
	cd cmd/fiber-example && go mod tidy

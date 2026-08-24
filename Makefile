.PHONY: help build test lint fmt clean install-tools migrate dev examples

# Variables
POSTGRES_DSN ?= postgres://user:password@localhost:5432/go_better_auth?sslmode=disable
SQLITE_DSN ?= ./db/go_better_auth.db

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

test-integration: ## Run integration tests
	go test -v ./tests/integration/...

test-postgres: ## Run PostgreSQL integration tests
	go test -v ./tests/integration -run TestPostgresUserRepository

test-sqlite: ## Run SQLite integration tests
	go test -v ./tests/integration -run TestSQLiteUserRepository

test-all: test test-integration ## Run all tests

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

# Migration tool
build-migrate: ## Build the migration tool
	go build -o bin/migrate ./cmd/migrate/main.go

# PostgreSQL migrations
migrate-up: build-migrate ## Run migrations up (PostgreSQL)
	./bin/migrate -driver=postgres -dsn="$(POSTGRES_DSN)" -direction=up

migrate-down: build-migrate ## Run migrations down (PostgreSQL)
	./bin/migrate -driver=postgres -dsn="$(POSTGRES_DSN)" -direction=down

# SQLite migrations
migrate-up-sqlite: build-migrate ## Run migrations up (SQLite)
	mkdir -p ./db
	./bin/migrate -driver=sqlite3 -dsn="$(SQLITE_DSN)" -direction=up

migrate-down-sqlite: build-migrate ## Run migrations down (SQLite)
	./bin/migrate -driver=sqlite3 -dsn="$(SQLITE_DSN)" -direction=down

# Create migrations
migrate-create: build-migrate ## Create new migration file
	@read -p "Enter migration name: " name; \
	./bin/migrate -driver=postgres -create="$$name"

migrate-create-sqlite: build-migrate ## Create new SQLite migration file
	@read -p "Enter migration name: " name; \
	./bin/migrate -driver=sqlite3 -create="$$name"

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

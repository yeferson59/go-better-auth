#!/bin/bash

set -e

echo "Running linting checks..."

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "golangci-lint not found. Installing..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

# Run golangci-lint
echo "Running golangci-lint..."
golangci-lint run

# Check if goimports is installed
if ! command -v goimports &> /dev/null; then
    echo "goimports not found. Installing..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# Check imports
echo "Checking imports..."
goimports -l .

echo "Linting completed successfully!"

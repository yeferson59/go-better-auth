#!/bin/bash

set -e

echo "Running database migrations..."

# Load environment variables if .env exists
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Default database URL
DB_URL=${DATABASE_URL:-"postgres://localhost:5432/go_better_auth?sslmode=disable"}

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "migrate tool not found. Installing..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# Run migrations
echo "Applying migrations to: $DB_URL"
migrate -path ./internal/infra/db/migrations -database "$DB_URL" up

echo "Migrations completed successfully!"

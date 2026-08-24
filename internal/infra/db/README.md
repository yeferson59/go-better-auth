# Database Infrastructure Layer

This directory contains the database infrastructure implementation using the Repository pattern with GORM.

## Architecture

The database layer follows the Repository pattern with the following components:

```
internal/infra/db/
├── models/             # GORM models and converters
│   ├── models.go      # Database models
│   └── converter.go   # Entity ↔ Model converters
├── postgres/          # PostgreSQL implementations
│   └── user_repository.go
├── sqlite/            # SQLite implementations
│   └── user_repository.go
└── migrations/        # Database migrations
    ├── postgres/      # PostgreSQL migrations
    └── sqlite/        # SQLite migrations
```

## Features

- **Repository Pattern**: Clean separation between business logic and data access
- **GORM Integration**: Using GORM for database operations
- **Multi-Database Support**: PostgreSQL and SQLite implementations
- **Migration System**: Database schema versioning with golang-migrate
- **Soft Deletes**: Built-in soft delete support via GORM
- **Type Safety**: Strong typing with entity conversion
- **Testing**: Comprehensive integration tests with testcontainers

## Usage

### Setting up Database Connection

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    
    "github.com/yeferson59/go-better-auth/internal/infra/db/postgres"
)

// PostgreSQL
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

// Create repository
userRepo := postgres.NewUserRepository(db)
```

### Using the Repository

```go
ctx := context.Background()

// Create user
user := &entities.User{
    Email:     "user@example.com",
    Username:  "johndoe",
    Password:  "hashed_password",
    FirstName: "John",
    LastName:  "Doe",
    IsActive:  true,
}

err := userRepo.Create(ctx, user)
if err != nil {
    log.Fatal(err)
}

// Get user by email
foundUser, err := userRepo.GetByEmail(ctx, "user@example.com")
if err != nil {
    log.Fatal(err)
}

// Update user
foundUser.FirstName = "Jane"
err = userRepo.Update(ctx, foundUser)
if err != nil {
    log.Fatal(err)
}

// List users with pagination
users, err := userRepo.List(ctx, 10, 0) // limit: 10, offset: 0
if err != nil {
    log.Fatal(err)
}

// Search users
results, err := userRepo.Search(ctx, "john", 10, 0)
if err != nil {
    log.Fatal(err)
}
```

## Database Migrations

### CLI Migration Tool

A CLI tool is provided for managing database migrations:

```bash
# Build the migration tool
go build -o migrate cmd/migrate/main.go

# Run migrations up
./migrate -driver=postgres -dsn="postgres://user:pass@localhost/dbname?sslmode=disable" -direction=up

# Run migrations down
./migrate -driver=postgres -dsn="postgres://user:pass@localhost/dbname?sslmode=disable" -direction=down

# Create new migration
./migrate -driver=postgres -create="add_user_profile_table"

# Migrate to specific version
./migrate -driver=postgres -dsn="..." -version=3

# Force migration to specific version (use with caution)
./migrate -driver=postgres -dsn="..." -force=2
```

### Migration Files

Migration files are stored in driver-specific directories:

- `migrations/postgres/` - PostgreSQL migrations
- `migrations/sqlite/` - SQLite migrations

Each migration has two files:
- `XXX_migration_name.up.sql` - Apply migration
- `XXX_migration_name.down.sql` - Rollback migration

## Testing

### Integration Tests

Integration tests are provided for both PostgreSQL and SQLite:

```bash
# Run PostgreSQL integration tests (requires Docker)
go test -v ./tests/integration -run TestPostgresUserRepository

# Run SQLite integration tests (in-memory)
go test -v ./tests/integration -run TestSQLiteUserRepository

# Run all integration tests
go test -v ./tests/integration
```

### Test Setup

#### PostgreSQL Tests
- Uses testcontainers to spin up PostgreSQL container
- Automatic cleanup after tests
- Requires Docker to be running

#### SQLite Tests
- Uses in-memory SQLite database
- Fast and isolated
- No external dependencies

## Repository Interface

All repositories implement the interfaces defined in `internal/core/interfaces/repository.go`:

```go
type UserRepository interface {
    // Basic CRUD operations
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id string) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    GetByUsername(ctx context.Context, username string) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id string) error
    
    // List operations
    List(ctx context.Context, limit, offset int) ([]*entities.User, error)
    Count(ctx context.Context) (int64, error)
    
    // Search operations
    Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)
    
    // User-specific operations
    GetByOAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error)
    GetWithRoles(ctx context.Context, id string) (*entities.User, error)
    GetWithSessions(ctx context.Context, id string) (*entities.User, error)
    
    // Bulk operations
    GetByIDs(ctx context.Context, ids []string) ([]*entities.User, error)
    UpdateLastLogin(ctx context.Context, id string) error
    
    // Status operations
    Activate(ctx context.Context, id string) error
    Deactivate(ctx context.Context, id string) error
    Verify(ctx context.Context, id string) error
}
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Roles Table

```sql
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Sessions Table

```sql
CREATE TABLE sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## Performance Considerations

1. **Indexes**: All tables have appropriate indexes for common queries
2. **Pagination**: List operations support limit/offset pagination
3. **Soft Deletes**: Uses GORM's soft delete feature
4. **Connection Pooling**: GORM handles connection pooling automatically
5. **Prepared Statements**: GORM uses prepared statements for security

## Error Handling

The repository implementations provide consistent error handling:

- Database errors are wrapped with descriptive messages
- `gorm.ErrRecordNotFound` is handled consistently (returns nil, nil)
- Context cancellation is properly handled
- Transaction rollbacks on errors

## Security

1. **SQL Injection Prevention**: Uses prepared statements
2. **Password Exclusion**: Password field excluded from JSON serialization
3. **Soft Deletes**: Deleted records are not hard-deleted
4. **Input Validation**: Relies on GORM validations and constraints

## Extension Points

To add new repositories:

1. Create GORM models in `models/models.go`
2. Add conversion methods in `models/converter.go`
3. Implement repository interface in `postgres/` and `sqlite/`
4. Add migration files
5. Write integration tests

## Dependencies

- **GORM**: ORM for database operations
- **golang-migrate**: Database migration tool
- **testcontainers**: Integration testing with real databases
- **PostgreSQL**: Primary database driver
- **SQLite**: Alternative/testing database driver

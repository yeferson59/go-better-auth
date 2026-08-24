# go-better-auth

A comprehensive, production-ready authentication library for Go applications, built with hexagonal architecture principles.

## Features

- 🔐 **Authentication**: Login, registration, password reset
- 🎭 **OAuth Integration**: Support for popular OAuth providers
- 🔑 **RBAC**: Role-based access control
- 🌐 **Multi-framework**: Gin, Fiber adapters included
- 🏗️ **Hexagonal Architecture**: Clean, testable, maintainable code
- 📦 **Modular Design**: Use only what you need
- 🔒 **Security First**: JWT tokens, secure password hashing
- 📧 **Email Integration**: Password reset, email verification

## Architecture

This project follows hexagonal architecture (ports and adapters) principles:

```
/cmd                – example apps (gin-example, fiber-example)
/internal
  /core             – domain entities, interfaces
  /auth             – use-cases (login, register, reset, oauth, rbac)
  /infra            – db, jwt, mail, logger, config
  /http             – adapters (gin, fiber)
/pkg                – exported API for users
/scripts            – dev helpers (migrate, lint)
.github             – actions, templates
```

## Quick Start

### Installation

```bash
go get github.com/yeferson59/go-better-auth
```

### Basic Usage

```go
package main

import (
    "github.com/yeferson59/go-better-auth/pkg/auth"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize the auth service
    authService := auth.NewService(auth.Config{
        JWTSecret: "your-secret-key",
        Database:  "postgres://...",
    })

    // Create Gin router
    r := gin.Default()

    // Register auth routes
    authService.RegisterRoutes(r)

    r.Run(":8080")
}
```

## Examples

Check out the example applications:

- **Gin Example**: `cmd/gin-example/` - Complete web application using Gin
- **Fiber Example**: `cmd/fiber-example/` - High-performance app using Fiber

To run examples:

```bash
make examples
./bin/gin-example
# or
./bin/fiber-example
```

## Development

### Prerequisites

- Go 1.21 or later
- PostgreSQL (for database examples)
- Make

### Setup

1. Clone the repository:
```bash
git clone https://github.com/yeferson59/go-better-auth.git
cd go-better-auth
```

2. Install development tools:
```bash
make install-tools
```

3. Run development cycle:
```bash
make dev  # format, lint, test
```

### Available Commands

```bash
make help           # Show all available commands
make build          # Build the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make fmt            # Format code
make clean          # Clean build artifacts
make migrate        # Run database migrations
make examples       # Build example applications
```

## API Documentation

### Authentication Endpoints

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh JWT token
- `POST /auth/forgot-password` - Password reset request
- `POST /auth/reset-password` - Password reset confirmation

### OAuth Endpoints

- `GET /auth/oauth/{provider}` - Initiate OAuth flow
- `GET /auth/oauth/{provider}/callback` - OAuth callback

### RBAC Endpoints

- `GET /auth/roles` - List roles
- `POST /auth/roles` - Create role
- `PUT /auth/roles/{id}` - Update role
- `DELETE /auth/roles/{id}` - Delete role

## Configuration

The library supports configuration through environment variables or config files:

```go
type Config struct {
    JWTSecret        string
    JWTExpiration    time.Duration
    Database         DatabaseConfig
    Email            EmailConfig
    OAuth            OAuthConfig
    RBAC             RBACConfig
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/yeferson59/go-better-auth/wiki)
- 🐛 [Issue Tracker](https://github.com/yeferson59/go-better-auth/issues)
- 💬 [Discussions](https://github.com/yeferson59/go-better-auth/discussions)

## Acknowledgments

- Built with ❤️ for the Go community
- Inspired by modern authentication best practices
- Following hexagonal architecture principles

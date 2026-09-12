# Security Utilities

This document describes the security utilities implemented in the Go Better Auth project. These utilities provide secure password hashing, JWT token management, environment-based configuration, and cryptographically secure token generation.

## Features

- **Password Hashing**: bcrypt-based password hashing with configurable cost
- **JWT Management**: HS256/RS256 JWT tokens with key rotation and blacklist support
- **Token Generation**: Secure random tokens using `crypto/rand`
- **Environment Configuration**: Configuration management using `caarlos0/env` and `viper`
- **Token Blacklist**: In-memory and Redis-based token blacklisting

## Components

### 1. Bcrypt Hasher (`bcrypt_hasher.go`)

The bcrypt hasher provides secure password hashing functionality.

```go
// Create hasher factory
hashConfig := interfaces.HashConfig{
    Algorithm: "bcrypt",
    Cost:      12,
}
hasherFactory := security.NewHasherFactory(hashConfig)
hasher := hasherFactory.GetDefault()

// Hash password
hashedPassword, err := hasher.Hash(ctx, "myPassword123!")
if err != nil {
    log.Fatal(err)
}

// Verify password
isValid, err := hasher.Verify(ctx, "myPassword123!", hashedPassword)
if err != nil {
    log.Fatal(err)
}
```

**Features:**

- Configurable bcrypt cost (4-31)
- Password verification
- Hash rehashing detection
- Input validation
- Context-aware operations

### 2. JWT Manager (`jwt_manager.go`)

The JWT manager handles JWT token operations with support for both HS256 and RS256 algorithms.

```go
// Create JWT manager
jwtManager, err := security.NewJWTManager(authConfig.JWT, blacklist)
if err != nil {
    log.Fatal(err)
}

// Generate JWT token
claims := map[string]any{
    "sub":   "user123",
    "email": "user@example.com",
    "roles": []string{"user", "admin"},
}
token, err := jwtManager.GenerateJWT(ctx, "user123", claims)

// Validate JWT token
tokenClaims, err := jwtManager.ValidateJWT(ctx, token)

// Refresh JWT token
accessToken, refreshToken, err := jwtManager.RefreshJWT(ctx, oldRefreshToken)

// Revoke JWT token
err = jwtManager.RevokeJWT(ctx, token)
```

**Features:**

- HS256 and RS256 algorithm support
- Token blacklisting integration
- Key rotation for RS256
- Configurable expiration times
- Automatic token refresh
- Context-aware operations

### 3. Token Blacklist (`token_blacklist.go`)

The token blacklist provides token revocation functionality.

```go
// In-memory blacklist
blacklist := security.NewInMemoryTokenBlacklist()

// Redis blacklist
blacklist := security.NewRedisTokenBlacklist(redisClient, "blacklist:")

// Add token to blacklist
err := blacklist.Add(ctx, token, expirationTime)

// Check if token is blacklisted
isBlacklisted, err := blacklist.IsBlacklisted(ctx, token)

// Remove token from blacklist
err := blacklist.Remove(ctx, token)
```

**Features:**

- In-memory and Redis storage options
- Automatic cleanup of expired tokens
- Context-aware operations
- Thread-safe operations

### 4. Token Generator (`token_generator.go`)

The token generator provides cryptographically secure token generation.

```go
// Create token generator
tokenGenerator := security.NewSecureTokenGenerator("HS256", secretKey)

// Generate random token
randomToken, err := tokenGenerator.GenerateRandom(32)

// Generate secure token
secureToken, err := tokenGenerator.GenerateSecure(64)

// Generate UUID token
uuidToken, err := tokenGenerator.GenerateUUID()

// Generate specific purpose tokens
passwordResetToken, err := tokenGenerator.GeneratePasswordResetToken()
emailVerificationToken, err := tokenGenerator.GenerateEmailVerificationToken()
oauthState, err := tokenGenerator.GenerateOAuthState()
```

**Features:**

- Cryptographically secure random generation using `crypto/rand`
- Multiple token formats (base64, hex, UUID)
- Purpose-specific token generators
- Token validation and security checks
- Backup code generation

### 5. Environment Configuration (`env_config.go`)

The configuration manager handles environment-based configuration with support for multiple sources.

```go
// Create configuration manager
configManager, err := config.NewManager()
if err != nil {
    log.Fatal(err)
}

// Get configuration
authConfig := configManager.GetConfig()

// Validate configuration
err = configManager.ValidateConfig()

// Watch for configuration changes
err = configManager.WatchConfig(func(config *interfaces.AuthConfig) {
    // Handle configuration change
})
```

**Features:**

- Environment variable parsing using `caarlos0/env`
- YAML file configuration using `viper`
- Configuration validation
- Hot reloading support
- Multiple configuration sources

## Environment Variables

The following environment variables are supported:

### JWT Configuration

- `JWT_SECRET`: JWT signing secret (required, min 32 characters)
- `JWT_ALGORITHM`: JWT algorithm (default: HS256)
- `JWT_ISSUER`: JWT issuer (default: go-better-auth)
- `JWT_AUDIENCE`: JWT audience (default: go-better-auth)
- `JWT_EXPIRATION`: JWT expiration time (default: 15m)
- `JWT_REFRESH_EXPIRATION`: JWT refresh expiration (default: 168h)
- `JWT_CLOCK_SKEW`: JWT clock skew tolerance (default: 5m)

### Session Configuration

- `SESSION_SECRET`: Session secret (required, min 32 characters)
- `SESSION_EXPIRATION`: Session expiration time (default: 24h)
- `SESSION_COOKIE_NAME`: Session cookie name (default: session)
- `SESSION_SECURE`: Use secure cookies (default: true)
- `SESSION_HTTP_ONLY`: Use HTTP-only cookies (default: true)

### Password Configuration

- `PASSWORD_ALGORITHM`: Password hashing algorithm (default: bcrypt)
- `PASSWORD_BCRYPT_COST`: Bcrypt cost factor (default: 12)
- `PASSWORD_MIN_LENGTH`: Minimum password length (default: 8)
- `PASSWORD_MAX_LENGTH`: Maximum password length (default: 128)
- `PASSWORD_REQUIRE_UPPER`: Require uppercase letters (default: true)
- `PASSWORD_REQUIRE_LOWER`: Require lowercase letters (default: true)
- `PASSWORD_REQUIRE_DIGIT`: Require digits (default: true)
- `PASSWORD_REQUIRE_SYMBOL`: Require symbols (default: false)

### Security Configuration

- `ENCRYPTION_KEY`: Encryption key (required, min 32 characters)
- `CORS_ALLOWED_ORIGINS`: CORS allowed origins (default: *)
- `CSRF_ENABLED`: Enable CSRF protection (default: true)
- `RATE_LIMIT_ENABLED`: Enable rate limiting (default: true)

## Security Best Practices

### Password Security

1. Use bcrypt with cost factor 12 or higher
2. Enforce strong password requirements
3. Implement password history to prevent reuse
4. Use secure password reset tokens

### JWT Security

1. Use strong secrets (min 32 characters)
2. Implement token blacklisting for revocation
3. Use short expiration times for access tokens
4. Implement proper token refresh mechanisms
5. Consider RS256 for distributed systems

### Token Security

1. Use `crypto/rand` for all token generation
2. Use appropriate token lengths (min 32 bytes for security tokens)
3. Implement token expiration
4. Use purpose-specific tokens
5. Validate token formats and lengths

### Configuration Security

1. Store secrets in environment variables
2. Use different secrets for different environments
3. Rotate secrets regularly
4. Validate configuration on startup
5. Use secure defaults

## Examples

See `examples/security_example.go` for a complete example of how to use all security utilities.

## Testing

Run the security utilities tests:

```bash
go test ./internal/infra/security/...
```

## Dependencies

- `golang.org/x/crypto`: For bcrypt password hashing
- `github.com/golang-jwt/jwt/v5`: For JWT token handling
- `github.com/caarlos0/env/v10`: For environment variable parsing
- `github.com/spf13/viper`: For configuration management
- `github.com/google/uuid`: For UUID generation
- `github.com/fsnotify/fsnotify`: For configuration file watching

## License

This project is licensed under the MIT License. See LICENSE file for details.

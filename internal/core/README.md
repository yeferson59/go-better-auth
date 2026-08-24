# Core Domain

This directory contains the core domain entities and interfaces for the authentication system.

## Structure

- `entities/` - Domain entities (User, Role, Permission, etc.)
- `interfaces/` - Port interfaces (repositories, services)
- `errors/` - Domain-specific errors

## Principles

- No external dependencies
- Pure business logic
- Framework-agnostic interfaces
- Clean Architecture compliance

## Core Entities

### User Entity (`entities/user.go`)

The `User` entity represents a user in the system with comprehensive profile and authentication information.

**Key Features:**
- Basic user information (ID, email, username, names)
- Authentication status (active, verified, password)
- Profile information and preferences
- Relationships with roles, sessions, and OAuth links
- Business methods for permissions and status checking

**Methods:**
- `IsEmailVerified()` - Check if user's email is verified
- `CanLogin()` - Check if user can login (active and verified)
- `HasRole(roleName)` - Check if user has specific role
- `HasPermission(permissionName)` - Check if user has specific permission
- `GetFullName()` - Get user's full name
- `UpdateLastLogin()` - Update last login timestamp

### Session Entity (`entities/session.go`)

The `Session` entity manages user authentication sessions.

**Key Features:**
- Session tokens and metadata
- Expiration and revocation handling
- IP address and user agent tracking
- Session status management

**Methods:**
- `IsActive()` - Check if session is active
- `IsExpired()` - Check if session has expired
- `IsRevoked()` - Check if session has been revoked
- `GetStatus()` - Get current session status
- `Revoke()` - Revoke the session
- `Extend(duration)` - Extend session expiration

### Role Entity (`entities/role.go`)

The `Role` entity implements role-based access control.

**Key Features:**
- Role hierarchy and system roles
- Permission management
- System role protection
- User-role relationships

**Methods:**
- `HasPermission(permissionName)` - Check if role has permission
- `AddPermission(permission)` - Add permission to role
- `RemovePermission(permissionName)` - Remove permission from role
- `IsSystemRole()` - Check if role is system-protected
- `CanDelete()` - Check if role can be deleted

### Permission Entity (`entities/permission.go`)

The `Permission` entity provides granular access control.

**Key Features:**
- Resource-action based permissions
- System permission protection
- Predefined system permissions
- Role-permission relationships

**Methods:**
- `GetFullName()` - Get full permission name (resource:action)
- `IsSystemPermission()` - Check if permission is system-protected
- `CanDelete()` - Check if permission can be deleted
- `Matches(resource, action)` - Check if permission matches criteria

### OAuthProvider Entity (`entities/oauth_provider.go`)

The `OAuthProvider` entity manages third-party authentication providers.

**Key Features:**
- Provider configuration and endpoints
- Supported providers (Google, GitHub, Facebook, etc.)
- Provider state management
- User info mapping

**Methods:**
- `IsSupported()` - Check if provider is supported
- `GetAuthURL()` - Get authorization URL
- `GetTokenURL()` - Get token exchange URL
- `Enable()` - Enable the provider
- `Disable()` - Disable the provider

### PasswordResetToken Entity (`entities/password_reset_token.go`)

The `PasswordResetToken` entity manages password reset functionality.

**Key Features:**
- Token generation and validation
- Expiration handling
- Usage tracking
- Multiple token types support

**Methods:**
- `IsExpired()` - Check if token has expired
- `IsUsed()` - Check if token has been used
- `IsValid()` - Check if token is valid
- `Use()` - Mark token as used

## Core Interfaces

### Repository Interfaces (`interfaces/repository.go`)

Comprehensive repository interfaces for data access abstraction:

- `UserRepository` - User data operations
- `SessionRepository` - Session management
- `RoleRepository` - Role and permission management
- `PermissionRepository` - Permission operations
- `OAuthProviderRepository` - OAuth provider management
- `TokenRepository` - Token operations
- `PasswordResetTokenRepository` - Password reset token management

### Security Interfaces

#### Hasher Interface (`interfaces/hasher.go`)

Password hashing and verification with multiple algorithm support:

- `Hash(password)` - Generate password hash
- `Verify(password, hash)` - Verify password against hash
- `NeedsRehash(hash)` - Check if hash needs updating

**Supported Algorithms:**
- Bcrypt (default)
- Argon2
- Scrypt

#### TokenManager Interface (`interfaces/token_manager.go`)

Comprehensive token management for JWT and session tokens:

- `GenerateJWT(userID, claims)` - Generate JWT token
- `ValidateJWT(token)` - Validate JWT token
- `RefreshJWT(refreshToken)` - Refresh JWT token
- `GenerateSessionToken(userID, metadata)` - Generate session token
- `ValidateSessionToken(token)` - Validate session token
- `GenerateRandomToken(length)` - Generate random token
- `GenerateSecureToken(length)` - Generate cryptographically secure token

### Communication Interfaces

#### Mailer Interface (`interfaces/mailer.go`)

Email sending capabilities with template support:

- `Send(email)` - Send plain text email
- `SendHTML(email)` - Send HTML email
- `SendTemplate(template, data, to, subject)` - Send templated email
- `SendBatch(emails)` - Send multiple emails
- `Verify()` - Verify email configuration

**Features:**
- Multiple email providers (SMTP, SendGrid, Mailgun, SES)
- Template engine support
- Email queuing and rate limiting
- Email validation

#### OAuth Provider Interface (`interfaces/oauth_provider.go`)

Third-party authentication provider integration:

- `GetAuthURL(state, scopes)` - Get authorization URL
- `ExchangeCode(code)` - Exchange code for tokens
- `RefreshToken(refreshToken)` - Refresh access token
- `GetUserInfo(accessToken)` - Get user information
- `RevokeToken(token)` - Revoke access token
- `ValidateToken(token)` - Validate access token

### System Interfaces

#### Logger Interface (`interfaces/logger.go`)

Structured logging with authentication-specific features:

- `Debug/Info/Warn/Error/Fatal/Panic` - Standard logging levels
- `WithFields(fields)` - Add structured fields
- `WithContext(ctx)` - Add context information
- `WithError(err)` - Add error information

**Authentication-specific logging:**
- `LogUserLogin/Logout` - User authentication events
- `LogPasswordChange/Reset` - Password events
- `LogOAuthLogin` - OAuth authentication events
- `LogSecurityEvent` - Security events
- `LogAuditEvent` - Audit events

#### Config Interface (`interfaces/config.go`)

Configuration management with comprehensive authentication settings:

- `Get/GetString/GetInt/GetBool/GetDuration` - Type-safe configuration retrieval
- `Set/SetDefault` - Configuration updates
- `UnmarshalKey/Unmarshal` - Structured configuration
- `WatchConfig` - Configuration change monitoring
- `Validate` - Configuration validation

**Configuration Sections:**
- JWT settings (secret, algorithm, expiration)
- Session configuration (storage, cookies, security)
- Password policies (complexity, history, reset)
- OAuth provider settings
- Email configuration
- Database and Redis settings
- Rate limiting configuration
- Security settings (2FA, account lockout, IP filtering)
- Feature flags

## Error Handling (`errors/errors.go`)

Comprehensive domain-specific error handling:

**Error Categories:**
- Authentication errors (invalid credentials, account locked, etc.)
- Authorization errors (forbidden, insufficient permissions)
- Password errors (too weak, reused, expired)
- OAuth errors (provider disabled, invalid state)
- Rate limiting errors
- Validation errors
- System errors
- Two-factor authentication errors

**Error Structure:**
- `Code` - Machine-readable error code
- `Message` - Human-readable error message
- `Details` - Additional error context
- `Status` - HTTP status code

**Helper Functions:**
- `IsAuthenticationError(err)` - Check error type
- `IsAuthorizationError(err)` - Check error type
- `IsValidationError(err)` - Check error type
- `IsSystemError(err)` - Check error type
- `IsRateLimitError(err)` - Check error type

## Architecture Benefits

1. **Clean Architecture**: Core domain is isolated from external dependencies
2. **Testability**: Interfaces enable easy mocking and testing
3. **Flexibility**: Implementation can be swapped without changing business logic
4. **Extensibility**: New features can be added without breaking existing code
5. **Maintainability**: Clear separation of concerns and responsibilities
6. **Security**: Comprehensive security models and validation
7. **Scalability**: Efficient data access patterns and caching support

## Usage Examples

### User Authentication

```go
// Validate user credentials
user, err := userRepo.GetByEmail(ctx, email)
if err != nil {
    return errors.NewUserNotFoundError(email)
}

if !user.CanLogin() {
    return errors.NewAccountDisabledError()
}

valid, err := hasher.Verify(ctx, password, user.Password)
if err != nil || !valid {
    return errors.NewInvalidCredentialsError()
}

// Generate session
token, err := tokenManager.GenerateSessionToken(ctx, user.ID, metadata)
if err != nil {
    return errors.NewInternalServerError("Failed to generate session")
}
```

### Permission Checking

```go
// Check user permission
if !user.HasPermission("users:read") {
    return errors.NewInsufficientPermissionsError("users:read")
}

// Check role permission
role, err := roleRepo.GetByName(ctx, "admin")
if err != nil {
    return errors.NewRoleNotFoundError("admin")
}

if !role.HasPermission("system:admin") {
    return errors.NewInsufficientPermissionsError("system:admin")
}
```

### OAuth Authentication

```go
// Get OAuth authorization URL
authURL, err := oauthProvider.GetAuthURL(ctx, state, scopes)
if err != nil {
    return errors.NewOAuthProviderNotFoundError(provider)
}

// Exchange code for tokens
tokens, err := oauthProvider.ExchangeCode(ctx, code)
if err != nil {
    return errors.NewOAuthInvalidCodeError()
}

// Get user info
userInfo, err := oauthProvider.GetUserInfo(ctx, tokens.AccessToken)
if err != nil {
    return errors.NewOAuthTokenExpiredError()
}
```

This core domain provides a solid foundation for building a robust, secure, and maintainable authentication system.

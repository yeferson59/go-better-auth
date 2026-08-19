package interfaces

import (
	"context"
	"time"
)

// TokenManager defines the interface for token management operations
type TokenManager interface {
	// JWT Token operations
	GenerateJWT(ctx context.Context, userID string, claims map[string]any) (string, error)
	ValidateJWT(ctx context.Context, token string) (*TokenClaims, error)
	RefreshJWT(ctx context.Context, refreshToken string) (string, string, error) // returns new access token and refresh token
	RevokeJWT(ctx context.Context, token string) error

	// Session Token operations
	GenerateSessionToken(ctx context.Context, userID string, metadata map[string]any) (string, error)
	ValidateSessionToken(ctx context.Context, token string) (*SessionTokenInfo, error)
	RefreshSessionToken(ctx context.Context, token string) (string, error)
	RevokeSessionToken(ctx context.Context, token string) error

	// Random Token operations
	GenerateRandomToken(ctx context.Context, length int) (string, error)
	GenerateSecureToken(ctx context.Context, length int) (string, error)

	// Token validation
	IsTokenValid(ctx context.Context, token string) (bool, error)
	GetTokenExpiration(ctx context.Context, token string) (time.Time, error)

	// Token introspection
	GetTokenClaims(ctx context.Context, token string) (map[string]any, error)
	GetTokenMetadata(ctx context.Context, token string) (map[string]any, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    string         `json:"userId"`
	Email     string         `json:"email,omitempty"`
	Username  string         `json:"username,omitempty"`
	Roles     []string       `json:"roles,omitempty"`
	Subject   string         `json:"sub"`
	Audience  string         `json:"aud"`
	Issuer    string         `json:"iss"`
	IssuedAt  time.Time      `json:"iat"`
	ExpiresAt time.Time      `json:"exp"`
	NotBefore time.Time      `json:"nbf"`
	JTI       string         `json:"jti"` // JWT ID
	Custom    map[string]any `json:"custom,omitempty"`
}

// SessionTokenInfo represents session token information
type SessionTokenInfo struct {
	UserID    string         `json:"userId"`
	SessionID string         `json:"sessionId"`
	Email     string         `json:"email,omitempty"`
	Username  string         `json:"username,omitempty"`
	Roles     []string       `json:"roles,omitempty"`
	ExpiresAt time.Time      `json:"expiresAt"`
	IssuedAt  time.Time      `json:"issuedAt"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// TokenType represents different types of tokens
type TokenType string

const (
	TokenTypeJWT               TokenType = "jwt"
	TokenTypeSession           TokenType = "session"
	TokenTypeRefresh           TokenType = "refresh"
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypeInvitation        TokenType = "invitation"
	TokenTypeAPIKey            TokenType = "api_key"
)

// TokenConfig represents token configuration
type TokenConfig struct {
	// JWT configuration
	JWTSecret            string        `json:"jwtSecret"`
	JWTAlgorithm         string        `json:"jwtAlgorithm"`
	JWTIssuer            string        `json:"jwtIssuer"`
	JWTAudience          string        `json:"jwtAudience"`
	JWTExpiration        time.Duration `json:"jwtExpiration"`
	JWTRefreshExpiration time.Duration `json:"jwtRefreshExpiration"`

	// Session token configuration
	SessionSecret     string        `json:"sessionSecret"`
	SessionExpiration time.Duration `json:"sessionExpiration"`

	// General token configuration
	TokenLength       int `json:"tokenLength"`
	SecureTokenLength int `json:"secureTokenLength"`

	// Token validation
	ClockSkew time.Duration `json:"clockSkew"`

	// Token storage
	UseRedis     bool          `json:"useRedis"`
	RedisPrefix  string        `json:"redisPrefix"`
	RedisTimeout time.Duration `json:"redisTimeout"`
}

// DefaultTokenConfig provides default token configuration
var DefaultTokenConfig = TokenConfig{
	JWTAlgorithm:         "HS256",
	JWTExpiration:        15 * time.Minute,
	JWTRefreshExpiration: 7 * 24 * time.Hour, // 7 days
	SessionExpiration:    24 * time.Hour,
	TokenLength:          32,
	SecureTokenLength:    64,
	ClockSkew:            5 * time.Minute,
	UseRedis:             false,
	RedisPrefix:          "auth:token:",
	RedisTimeout:         5 * time.Second,
}

// TokenValidator defines the interface for token validation
type TokenValidator interface {
	// Validate validates a token
	Validate(ctx context.Context, token string, tokenType TokenType) error

	// GetClaims extracts claims from a token
	GetClaims(ctx context.Context, token string) (map[string]any, error)

	// IsExpired checks if a token is expired
	IsExpired(ctx context.Context, token string) (bool, error)

	// IsRevoked checks if a token is revoked
	IsRevoked(ctx context.Context, token string) (bool, error)
}

// TokenStore defines the interface for token storage operations
type TokenStore interface {
	// Store stores a token
	Store(ctx context.Context, token string, data any, expiration time.Duration) error

	// Get retrieves token data
	Get(ctx context.Context, token string) (any, error)

	// Delete removes a token
	Delete(ctx context.Context, token string) error

	// Exists checks if a token exists
	Exists(ctx context.Context, token string) (bool, error)

	// SetExpiration sets token expiration
	SetExpiration(ctx context.Context, token string, expiration time.Duration) error

	// GetExpiration gets token expiration
	GetExpiration(ctx context.Context, token string) (time.Time, error)

	// Cleanup removes expired tokens
	Cleanup(ctx context.Context) error
}

// RefreshTokenManager defines the interface for refresh token management
type RefreshTokenManager interface {
	// Generate generates a new refresh token
	Generate(ctx context.Context, userID string, metadata map[string]any) (string, error)

	// Validate validates a refresh token
	Validate(ctx context.Context, token string) (*RefreshTokenInfo, error)

	// Revoke revokes a refresh token
	Revoke(ctx context.Context, token string) error

	// RevokeAll revokes all refresh tokens for a user
	RevokeAll(ctx context.Context, userID string) error

	// Cleanup removes expired refresh tokens
	Cleanup(ctx context.Context) error
}

// RefreshTokenInfo represents refresh token information
type RefreshTokenInfo struct {
	UserID    string         `json:"userId"`
	TokenID   string         `json:"tokenId"`
	ExpiresAt time.Time      `json:"expiresAt"`
	IssuedAt  time.Time      `json:"issuedAt"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// TokenBlacklist defines the interface for token blacklisting
type TokenBlacklist interface {
	// Add adds a token to the blacklist
	Add(ctx context.Context, token string, expiration time.Time) error

	// IsBlacklisted checks if a token is blacklisted
	IsBlacklisted(ctx context.Context, token string) (bool, error)

	// Remove removes a token from the blacklist
	Remove(ctx context.Context, token string) error

	// Cleanup removes expired blacklisted tokens
	Cleanup(ctx context.Context) error
}

// TokenGenerator defines the interface for token generation
type TokenGenerator interface {
	// GenerateJWT generates a JWT token
	GenerateJWT(claims map[string]any) (string, error)

	// GenerateRandom generates a random token
	GenerateRandom(length int) (string, error)

	// GenerateSecure generates a cryptographically secure token
	GenerateSecure(length int) (string, error)

	// GenerateUUID generates a UUID token
	GenerateUUID() (string, error)
}

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	// Set sets a key-value pair with expiration
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Get retrieves a value by key
	Get(ctx context.Context, key string) (string, error)

	// Del deletes a key
	Del(ctx context.Context, key string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Close closes the Redis connection
	Close() error
}

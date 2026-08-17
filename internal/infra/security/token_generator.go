package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

// SecureTokenGenerator implements the TokenGenerator interface
type SecureTokenGenerator struct {
	algorithm string
	secretKey []byte
}

// NewSecureTokenGenerator creates a new SecureTokenGenerator instance
func NewSecureTokenGenerator(algorithm string, secretKey []byte) *SecureTokenGenerator {
	return &SecureTokenGenerator{
		algorithm: algorithm,
		secretKey: secretKey,
	}
}

// GenerateJWT generates a JWT token with the given claims
func (g *SecureTokenGenerator) GenerateJWT(claims map[string]any) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(g.algorithm), jwt.MapClaims(claims))
	return token.SignedString(g.secretKey)
}

// GenerateRandom generates a random token of the specified length
func (g *SecureTokenGenerator) GenerateRandom(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// GenerateSecure generates a cryptographically secure token
func (g *SecureTokenGenerator) GenerateSecure(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateUUID generates a UUID token
func (g *SecureTokenGenerator) GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}
	return id.String(), nil
}

// GenerateHexToken generates a hex-encoded secure token
func (g *SecureTokenGenerator) GenerateHexToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

// GeneratePasswordResetToken generates a secure token for password reset
func (g *SecureTokenGenerator) GeneratePasswordResetToken() (string, error) {
	return g.GenerateSecure(32)
}

// GenerateEmailVerificationToken generates a secure token for email verification
func (g *SecureTokenGenerator) GenerateEmailVerificationToken() (string, error) {
	return g.GenerateSecure(32)
}

// GenerateOAuthState generates a secure state token for OAuth
func (g *SecureTokenGenerator) GenerateOAuthState() (string, error) {
	return g.GenerateSecure(32)
}

// GenerateAPIKey generates a secure API key
func (g *SecureTokenGenerator) GenerateAPIKey() (string, error) {
	return g.GenerateSecure(64)
}

// GenerateSessionToken generates a secure session token
func (g *SecureTokenGenerator) GenerateSessionToken() (string, error) {
	return g.GenerateSecure(64)
}

// GenerateCSRFToken generates a secure CSRF token
func (g *SecureTokenGenerator) GenerateCSRFToken() (string, error) {
	return g.GenerateSecure(32)
}

// GenerateNonce generates a secure nonce
func (g *SecureTokenGenerator) GenerateNonce() (string, error) {
	return g.GenerateSecure(16)
}

// GenerateInvitationToken generates a secure invitation token
func (g *SecureTokenGenerator) GenerateInvitationToken() (string, error) {
	return g.GenerateSecure(32)
}

// GenerateBackupCode generates a secure backup code
func (g *SecureTokenGenerator) GenerateBackupCode() (string, error) {
	// Generate 8-digit numeric backup code
	const digits = "0123456789"
	result := make([]byte, 8)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}

// GenerateSecurePin generates a secure PIN of specified length
func (g *SecureTokenGenerator) GenerateSecurePin(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}

// ValidateTokenLength validates if the token length is appropriate
func (g *SecureTokenGenerator) ValidateTokenLength(token string, minLength, maxLength int) error {
	if len(token) < minLength {
		return fmt.Errorf("token too short: minimum length is %d", minLength)
	}
	if len(token) > maxLength {
		return fmt.Errorf("token too long: maximum length is %d", maxLength)
	}
	return nil
}

// IsSecureToken checks if a token appears to be securely generated
func (g *SecureTokenGenerator) IsSecureToken(token string) bool {
	// Basic checks for secure token characteristics
	if len(token) < 16 {
		return false
	}

	// Check for common patterns that indicate weak tokens
	if isSequential(token) || isRepeating(token) {
		return false
	}

	return true
}

// isSequential checks if the token contains sequential characters
func isSequential(token string) bool {
	if len(token) < 3 {
		return false
	}

	sequential := 0
	for i := 1; i < len(token); i++ {
		if token[i] == token[i-1]+1 {
			sequential++
			if sequential >= 3 {
				return true
			}
		} else {
			sequential = 0
		}
	}
	return false
}

// isRepeating checks if the token contains too many repeating characters
func isRepeating(token string) bool {
	if len(token) < 3 {
		return false
	}

	charCount := make(map[rune]int)
	for _, char := range token {
		charCount[char]++
	}

	for _, count := range charCount {
		if count > len(token)/2 {
			return true
		}
	}
	return false
}

// TokenManager implements the interfaces.TokenManager interface
type TokenManager struct {
	generator interfaces.TokenGenerator
	store     interfaces.TokenStore
	blacklist interfaces.TokenBlacklist
	config    interfaces.TokenConfig
}

// NewTokenManager creates a new TokenManager instance
func NewTokenManager(generator interfaces.TokenGenerator, store interfaces.TokenStore, blacklist interfaces.TokenBlacklist, config interfaces.TokenConfig) *TokenManager {
	return &TokenManager{
		generator: generator,
		store:     store,
		blacklist: blacklist,
		config:    config,
	}
}

// GenerateRandomToken generates a random token
func (m *TokenManager) GenerateRandomToken(ctx context.Context, length int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if length <= 0 {
		length = m.config.TokenLength
	}

	return m.generator.GenerateRandom(length)
}

// GenerateSecureToken generates a secure token
func (m *TokenManager) GenerateSecureToken(ctx context.Context, length int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if length <= 0 {
		length = m.config.SecureTokenLength
	}

	return m.generator.GenerateSecure(length)
}

// IsTokenValid checks if a token is valid
func (m *TokenManager) IsTokenValid(ctx context.Context, token string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Check if token is blacklisted
	if m.blacklist != nil {
		blacklisted, err := m.blacklist.IsBlacklisted(ctx, token)
		if err != nil {
			return false, fmt.Errorf("failed to check blacklist: %w", err)
		}
		if blacklisted {
			return false, nil
		}
	}

	// Check if token exists in store. Without a store there is nothing to validate
	// against, so fail closed: returning true here would accept any string as a valid
	// token.
	if m.store == nil {
		return false, errors.New("no token store configured; cannot validate token")
	}

	exists, err := m.store.Exists(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}
	return exists, nil
}

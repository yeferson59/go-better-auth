package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

// Token type claim, used to keep an access token from being replayed as a refresh
// token and vice versa.
const (
	tokenTypeClaim   = "typ"
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// MinSecretLength is the minimum accepted length for an HMAC signing secret.
const MinSecretLength = 32

// JWTManager handles JWToken operations
type JWTManager struct {
	secretKey       []byte
	algorithm       string
	issuer          string
	blacklist       interfaces.TokenBlacklist
	rotationEnabled bool
	mu              sync.RWMutex // guards keyMap
	keyMap          map[string]*rsa.PrivateKey
	config          interfaces.JWTConfig
}

// NewJWTManager creates a JWTManager
func NewJWTManager(config interfaces.JWTConfig, blacklist interfaces.TokenBlacklist) (*JWTManager, error) {
	if config.Algorithm == "" {
		config.Algorithm = "HS256"
	}

	// HMAC algorithms are only as strong as their secret, so reject weak ones here
	// rather than relying on the caller having run config validation first.
	if strings.HasPrefix(config.Algorithm, "HS") && len(config.Secret) < MinSecretLength {
		return nil, fmt.Errorf("jwt secret must be at least %d characters, got %d", MinSecretLength, len(config.Secret))
	}

	manager := JWTManager{
		secretKey:       []byte(config.Secret),
		algorithm:       config.Algorithm,
		issuer:          config.Issuer,
		blacklist:       blacklist,
		rotationEnabled: false,
		keyMap:          make(map[string]*rsa.PrivateKey),
		config:          config,
	}

	if config.Algorithm == "RS256" {
		manager.rotationEnabled = true
		if err := manager.generateKeys(); err != nil {
			return nil, fmt.Errorf("failed to generate RSA keys: %w", err)
		}
	}

	return &manager, nil
}

// generateKeys generates a new RSA key pair, retaining the outgoing key so that tokens
// signed before the rotation stay verifiable until they expire.
func (m *JWTManager) generateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if current, exists := m.keyMap["current"]; exists {
		m.keyMap["previous"] = current
	}
	m.keyMap["current"] = privateKey
	return nil
}

// currentKey returns the key new tokens are signed with.
func (m *JWTManager) currentKey() (*rsa.PrivateKey, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key, ok := m.keyMap["current"]
	return key, ok
}

// verificationKeys returns every key a token may legitimately have been signed with,
// newest first.
func (m *JWTManager) verificationKeys() []*rsa.PrivateKey {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]*rsa.PrivateKey, 0, 2)
	if key, ok := m.keyMap["current"]; ok {
		keys = append(keys, key)
	}
	if key, ok := m.keyMap["previous"]; ok {
		keys = append(keys, key)
	}
	return keys
}

// GenerateJWT generates a JWT token
func (m *JWTManager) GenerateJWT(ctx context.Context, userID string, claims map[string]any) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Work on a copy so that issuing a token never mutates the caller's map.
	tokenClaims := make(jwt.MapClaims, len(claims)+5)
	maps.Copy(tokenClaims, claims)

	tokenClaims["iss"] = m.issuer
	tokenClaims["sub"] = userID
	tokenClaims["iat"] = time.Now().Unix()

	// Honour an explicit expiry (refresh tokens set their own); otherwise fall back to
	// the configured access-token lifetime.
	if _, ok := tokenClaims["exp"]; !ok {
		tokenClaims["exp"] = time.Now().Add(m.config.Expiration).Unix()
	}
	if _, ok := tokenClaims[tokenTypeClaim]; !ok {
		tokenClaims[tokenTypeClaim] = tokenTypeAccess
	}
	if _, ok := tokenClaims["aud"]; !ok && m.config.Audience != "" {
		tokenClaims["aud"] = m.config.Audience
	}

	var key any = m.secretKey
	if m.algorithm == "RS256" {
		privateKey, ok := m.currentKey()
		if !ok {
			return "", errors.New("no RSA signing key available")
		}
		key = privateKey
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(m.algorithm), tokenClaims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}

// GenerateRefreshToken issues a refresh token for a user. Refresh tokens carry a
// distinct type claim and the configured refresh lifetime.
func (m *JWTManager) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	return m.GenerateJWT(ctx, userID, map[string]any{
		tokenTypeClaim: tokenTypeRefresh,
		"exp":          time.Now().Add(m.config.RefreshExpiration).Unix(),
	})
}

// ValidateJWT validates a JWT token
func (m *JWTManager) ValidateJWT(ctx context.Context, tokenString string) (*interfaces.TokenClaims, error) {
	claims, err := m.validateClaims(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return &interfaces.TokenClaims{
		UserID:    getStringClaim(claims, "sub"),
		Subject:   getStringClaim(claims, "sub"),
		Issuer:    getStringClaim(claims, "iss"),
		Audience:  getStringClaim(claims, "aud"),
		IssuedAt:  getTimeClaim(claims, "iat"),
		ExpiresAt: getTimeClaim(claims, "exp"),
	}, nil
}

// validateClaims runs the full validation chain (revocation plus signature and
// registered claims) and returns the raw claims.
func (m *JWTManager) validateClaims(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Check if token is blacklisted
	if m.blacklist != nil {
		blacklisted, err := m.blacklist.IsBlacklisted(ctx, tokenString)
		if err != nil {
			return nil, fmt.Errorf("failed to check blacklist: %w", err)
		}
		if blacklisted {
			return nil, errors.New("token is blacklisted")
		}
	}

	return m.parseToken(tokenString)
}

// parseToken verifies a token's signature and registered claims. Pinning the signing
// method blocks algorithm-substitution attacks, and checking issuer and audience stops
// tokens minted by another service that happens to share the secret. For RS256 the
// previous key is also tried so a rotation does not invalidate tokens already issued.
func (m *JWTManager) parseToken(tokenString string) (jwt.MapClaims, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{m.algorithm}),
		jwt.WithLeeway(m.config.ClockSkew),
	}
	if m.issuer != "" {
		opts = append(opts, jwt.WithIssuer(m.issuer))
	}
	if m.config.Audience != "" {
		opts = append(opts, jwt.WithAudience(m.config.Audience))
	}

	var candidates []any
	switch m.algorithm {
	case "RS256":
		for _, key := range m.verificationKeys() {
			candidates = append(candidates, &key.PublicKey)
		}
		if len(candidates) == 0 {
			return nil, errors.New("no RSA verification key available")
		}
	default:
		candidates = []any{m.secretKey}
	}

	var lastErr error
	for _, key := range candidates {
		token, err := jwt.Parse(tokenString, func(*jwt.Token) (any, error) {
			return key, nil
		}, opts...)
		if err != nil {
			lastErr = err
			continue
		}
		if !token.Valid {
			lastErr = errors.New("token is invalid")
			continue
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("unable to parse claims")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("failed to parse token: %w", lastErr)
}

// RevokeJWT revokes a JWT token
func (m *JWTManager) RevokeJWT(ctx context.Context, tokenString string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if m.blacklist == nil {
		return errors.New("no blacklist configured; cannot revoke token")
	}

	// Blacklist the token until its own expiry. Using the configured access-token
	// lifetime here would drop long-lived refresh tokens out of the blacklist while
	// they were still valid, silently un-revoking them.
	expirationTime, err := m.tokenExpiry(tokenString)
	if err != nil {
		return err
	}

	return m.blacklist.Add(ctx, tokenString, expirationTime)
}

// tokenExpiry reads the exp claim without enforcing it, so an already-expired token can
// still be revoked without error.
func (m *JWTManager) tokenExpiry(tokenString string) (time.Time, error) {
	claims := jwt.MapClaims{}
	parser := jwt.NewParser()
	if _, _, err := parser.ParseUnverified(tokenString, claims); err != nil {
		return time.Time{}, fmt.Errorf("failed to read token expiration: %w", err)
	}

	if exp := getTimeClaim(claims, "exp"); !exp.IsZero() {
		return exp, nil
	}
	return time.Now().Add(m.config.Expiration), nil
}

// RefreshJWT refreshes a JWT token
func (m *JWTManager) RefreshJWT(ctx context.Context, refreshToken string) (string, string, error) {
	select {
	case <-ctx.Done():
		return "", "", ctx.Err()
	default:
	}

	// Validate refresh token
	claims, err := m.validateClaims(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Only a token minted as a refresh token may be exchanged. Without this check a
	// stolen access token could be renewed indefinitely.
	if getStringClaim(claims, tokenTypeClaim) != tokenTypeRefresh {
		return "", "", errors.New("token is not a refresh token")
	}

	userID := getStringClaim(claims, "sub")

	accessToken, err := m.GenerateJWT(ctx, userID, map[string]any{
		tokenTypeClaim: tokenTypeAccess,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := m.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := m.RevokeJWT(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// RotateKeys rotates the RSA keys for RS256 algorithm
func (m *JWTManager) RotateKeys() error {
	if !m.rotationEnabled {
		return errors.New("key rotation is not enabled")
	}

	// generateKeys retains the outgoing key as "previous" so tokens signed with it stay
	// verifiable until they expire.
	return m.generateKeys()
}

// GetPublicKey returns the current public key for RS256
func (m *JWTManager) GetPublicKey() (*rsa.PublicKey, error) {
	if !m.rotationEnabled {
		return nil, errors.New("RSA keys not available")
	}

	if currentKey, exists := m.currentKey(); exists {
		return &currentKey.PublicKey, nil
	}

	return nil, errors.New("no current key available")
}

// Helper functions for claim extraction
func getStringClaim(claims jwt.MapClaims, key string) string {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func getTimeClaim(claims jwt.MapClaims, key string) time.Time {
	if value, exists := claims[key]; exists {
		if timestamp, ok := value.(float64); ok {
			return time.Unix(int64(timestamp), 0)
		}
	}
	return time.Time{}
}

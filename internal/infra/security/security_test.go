package security

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

func TestBcryptHasher(t *testing.T) {
	ctx := context.Background()

	// Create hasher
	hasher := NewBcryptHasher(12)

	password := "testPassword123!"

	// Test password hashing
	hashedPassword, err := hasher.Hash(ctx, password)
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test password verification - correct password
	isValid, err := hasher.Verify(ctx, password, hashedPassword)
	require.NoError(t, err)
	assert.True(t, isValid)

	// Test password verification - incorrect password
	isValid, err = hasher.Verify(ctx, "wrongPassword", hashedPassword)
	require.NoError(t, err)
	assert.False(t, isValid)

	// Test hash validation
	err = hasher.ValidateHash(hashedPassword)
	assert.NoError(t, err)

	// Test invalid hash
	err = hasher.ValidateHash("invalid-hash")
	assert.Error(t, err)
}

func TestSecureTokenGenerator(t *testing.T) {
	secretKey := []byte("test-secret-key-that-is-at-least-32-characters-long")
	generator := NewSecureTokenGenerator("HS256", secretKey)

	// Test random token generation
	randomToken, err := generator.GenerateRandom(32)
	require.NoError(t, err)
	assert.Len(t, randomToken, 32)
	assert.True(t, generator.IsSecureToken(randomToken))

	// Test secure token generation
	secureToken, err := generator.GenerateSecure(64)
	require.NoError(t, err)
	assert.NotEmpty(t, secureToken)

	// Test UUID generation
	uuidToken, err := generator.GenerateUUID()
	require.NoError(t, err)
	assert.NotEmpty(t, uuidToken)
	assert.Len(t, uuidToken, 36) // UUID format: 8-4-4-4-12

	// Test password reset token
	passwordResetToken, err := generator.GeneratePasswordResetToken()
	require.NoError(t, err)
	assert.NotEmpty(t, passwordResetToken)

	// Test email verification token
	emailToken, err := generator.GenerateEmailVerificationToken()
	require.NoError(t, err)
	assert.NotEmpty(t, emailToken)

	// Test OAuth state
	oauthState, err := generator.GenerateOAuthState()
	require.NoError(t, err)
	assert.NotEmpty(t, oauthState)

	// Test backup code
	backupCode, err := generator.GenerateBackupCode()
	require.NoError(t, err)
	assert.Len(t, backupCode, 8)
	assert.Regexp(t, "^[0-9]{8}$", backupCode)

	// Test secure PIN
	pin, err := generator.GenerateSecurePin(6)
	require.NoError(t, err)
	assert.Len(t, pin, 6)
	assert.Regexp(t, "^[0-9]{6}$", pin)
}

func TestInMemoryTokenBlacklist(t *testing.T) {
	ctx := context.Background()
	blacklist := NewInMemoryTokenBlacklist()

	token := "test-token-123"
	expiration := time.Now().Add(time.Hour)

	// Test adding token to blacklist
	err := blacklist.Add(ctx, token, expiration)
	require.NoError(t, err)

	// Test checking if token is blacklisted
	isBlacklisted, err := blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.True(t, isBlacklisted)

	// Test checking non-blacklisted token
	isBlacklisted, err = blacklist.IsBlacklisted(ctx, "non-existent-token")
	require.NoError(t, err)
	assert.False(t, isBlacklisted)

	// Test removing token from blacklist
	err = blacklist.Remove(ctx, token)
	require.NoError(t, err)

	// Verify token is no longer blacklisted
	isBlacklisted, err = blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.False(t, isBlacklisted)

	// Test expired token handling
	expiredToken := "expired-token"
	pastExpiration := time.Now().Add(-time.Hour)
	err = blacklist.Add(ctx, expiredToken, pastExpiration)
	require.NoError(t, err)

	// Should return false for expired token
	isBlacklisted, err = blacklist.IsBlacklisted(ctx, expiredToken)
	require.NoError(t, err)
	assert.False(t, isBlacklisted)
}

func TestJWTManager(t *testing.T) {
	ctx := context.Background()

	config := interfaces.JWTConfig{
		Secret:            "test-jwt-secret-that-is-at-least-32-characters-long",
		Algorithm:         "HS256",
		Issuer:            "test-issuer",
		Audience:          "test-audience",
		Expiration:        15 * time.Minute,
		RefreshExpiration: 24 * time.Hour,
		ClockSkew:         5 * time.Minute,
	}

	blacklist := NewInMemoryTokenBlacklist()
	jwtManager, err := NewJWTManager(config, blacklist)
	require.NoError(t, err)

	userID := "test-user-123"
	claims := map[string]any{
		"sub":   userID,
		"email": "test@example.com",
		"roles": []string{"user"},
	}

	// Test JWT generation
	token, err := jwtManager.GenerateJWT(ctx, userID, claims)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test JWT validation
	tokenClaims, err := jwtManager.ValidateJWT(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, userID, tokenClaims.UserID)
	assert.Equal(t, config.Issuer, tokenClaims.Issuer)

	// Test JWT revocation
	err = jwtManager.RevokeJWT(ctx, token)
	require.NoError(t, err)

	// Test validation of revoked token
	_, err = jwtManager.ValidateJWT(ctx, token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "blacklisted")
}

func TestHasherFactory(t *testing.T) {
	config := interfaces.HashConfig{
		Algorithm: "bcrypt",
		Cost:      12,
	}

	factory := NewHasherFactory(config)

	// Test getting default hasher
	hasher := factory.GetDefault()
	assert.NotNil(t, hasher)

	// Test getting hasher by algorithm
	bcryptHasher, err := factory.GetByAlgorithm(interfaces.HashAlgorithmBcrypt)
	require.NoError(t, err)
	assert.NotNil(t, bcryptHasher)

	// Test unsupported algorithm
	_, err = factory.GetByAlgorithm("unsupported")
	assert.Error(t, err)

	// Test creating hasher with config
	newHasher, err := factory.Create(config)
	require.NoError(t, err)
	assert.NotNil(t, newHasher)
}

func TestTokenValidation(t *testing.T) {
	generator := NewSecureTokenGenerator("HS256", []byte("test-secret-key"))

	// Test token length validation
	err := generator.ValidateTokenLength("short", 10, 50)
	assert.Error(t, err)

	err = generator.ValidateTokenLength("this-is-a-very-long-token-that-exceeds-the-maximum-length", 10, 50)
	assert.Error(t, err)

	err = generator.ValidateTokenLength("valid-token-length", 10, 50)
	assert.NoError(t, err)

	// Test secure token validation
	assert.False(t, generator.IsSecureToken("short"))
	assert.False(t, generator.IsSecureToken("123456789abcdef"))    // sequential
	assert.False(t, generator.IsSecureToken("aaaaaaaaaaaaaaaaaa")) // repeating

	// Generate and test a proper secure token
	secureToken, err := generator.GenerateSecure(32)
	require.NoError(t, err)
	assert.True(t, generator.IsSecureToken(secureToken))
}

func BenchmarkBcryptHash(b *testing.B) {
	hasher := NewBcryptHasher(12)
	ctx := context.Background()
	password := "testPassword123!"

	b.ResetTimer()
	for b.Loop() {
		_, _ = hasher.Hash(ctx, password)
	}
}

func BenchmarkSecureTokenGeneration(b *testing.B) {
	generator := NewSecureTokenGenerator("HS256", []byte("test-secret-key"))

	b.ResetTimer()
	for b.Loop() {
		_, _ = generator.GenerateSecure(32)
	}
}

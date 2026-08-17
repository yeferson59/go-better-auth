package security

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

func testJWTConfig() interfaces.JWTConfig {
	return interfaces.JWTConfig{
		Secret:            "test-jwt-secret-that-is-at-least-32-characters-long",
		Algorithm:         "HS256",
		Issuer:            "test-issuer",
		Audience:          "test-audience",
		Expiration:        15 * time.Minute,
		RefreshExpiration: 24 * time.Hour,
		ClockSkew:         5 * time.Minute,
	}
}

func newTestJWTManager(t *testing.T, config interfaces.JWTConfig) *JWTManager {
	t.Helper()

	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() { _ = blacklist.Close() })

	manager, err := NewJWTManager(config, blacklist)
	require.NoError(t, err)
	return manager
}

// TestBlacklistConcurrentAccess exercises the read path under contention. Deleting
// expired entries while holding only a read lock used to make this crash the process
// with "concurrent map read and map write".
func TestBlacklistConcurrentAccess(t *testing.T) {
	ctx := context.Background()

	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() { _ = blacklist.Close() })

	// Seed entries that are already expired, so the read path hits the expiry branch.
	for i := range 50 {
		require.NoError(t, blacklist.Add(ctx, fmt.Sprintf("expired-%d", i), time.Now().Add(-time.Hour)))
	}

	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()
			for range 100 {
				_, err := blacklist.IsBlacklisted(ctx, fmt.Sprintf("expired-%d", i))
				assert.NoError(t, err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			for j := range 100 {
				err := blacklist.Add(ctx, fmt.Sprintf("live-%d-%d", i, j), time.Now().Add(time.Hour))
				assert.NoError(t, err)
			}
		}(i)
	}
	wg.Wait()

	// An expired entry must not count as blacklisted.
	blacklisted, err := blacklist.IsBlacklisted(ctx, "expired-0")
	require.NoError(t, err)
	assert.False(t, blacklisted)
}

// TestBlacklistDoesNotStorePlaintextTokens verifies tokens are hashed at rest.
func TestBlacklistDoesNotStorePlaintextTokens(t *testing.T) {
	ctx := context.Background()

	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() { _ = blacklist.Close() })

	const token = "header.payload.signature"
	require.NoError(t, blacklist.Add(ctx, token, time.Now().Add(time.Hour)))

	_, storedInPlaintext := blacklist.tokens[token]
	assert.False(t, storedInPlaintext, "raw token must not be used as the map key")

	blacklisted, err := blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.True(t, blacklisted)

	require.NoError(t, blacklist.Remove(ctx, token))
	blacklisted, err = blacklist.IsBlacklisted(ctx, token)
	require.NoError(t, err)
	assert.False(t, blacklisted)
}

// TestBlacklistCloseStopsCleanupRoutine guards against the cleanup goroutine leaking.
func TestBlacklistCloseStopsCleanupRoutine(t *testing.T) {
	blacklist := NewInMemoryTokenBlacklist()

	require.NoError(t, blacklist.Close())
	// Close is idempotent; a second call must not panic on a closed channel.
	require.NoError(t, blacklist.Close())
}

// TestAccessTokenCannotBeUsedAsRefreshToken covers the token-type confusion that let a
// stolen access token be renewed indefinitely.
func TestAccessTokenCannotBeUsedAsRefreshToken(t *testing.T) {
	ctx := context.Background()
	manager := newTestJWTManager(t, testJWTConfig())

	accessToken, err := manager.GenerateJWT(ctx, "user-1", map[string]any{})
	require.NoError(t, err)

	_, _, err = manager.RefreshJWT(ctx, accessToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a refresh token")
}

// TestRefreshTokenRoundTrip checks that a genuine refresh token is exchangeable and
// that the old one is revoked in the process.
func TestRefreshTokenRoundTrip(t *testing.T) {
	ctx := context.Background()
	manager := newTestJWTManager(t, testJWTConfig())

	refreshToken, err := manager.GenerateRefreshToken(ctx, "user-1")
	require.NoError(t, err)

	accessToken, newRefreshToken, err := manager.RefreshJWT(ctx, refreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, newRefreshToken)

	// The spent refresh token must no longer be accepted.
	_, _, err = manager.RefreshJWT(ctx, refreshToken)
	assert.Error(t, err)

	// The freshly minted access token is not itself exchangeable.
	_, _, err = manager.RefreshJWT(ctx, accessToken)
	assert.Error(t, err)
}

// TestRefreshTokenHonoursRefreshExpiration covers GenerateJWT clobbering the exp claim,
// which silently capped refresh tokens at the access-token lifetime.
func TestRefreshTokenHonoursRefreshExpiration(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()
	manager := newTestJWTManager(t, config)

	refreshToken, err := manager.GenerateRefreshToken(ctx, "user-1")
	require.NoError(t, err)

	claims, err := manager.ValidateJWT(ctx, refreshToken)
	require.NoError(t, err)

	// Must be far beyond the 15 minute access-token window.
	assert.True(t, claims.ExpiresAt.After(time.Now().Add(config.Expiration+time.Hour)),
		"refresh token expires at %s, which is within the access-token window", claims.ExpiresAt)
}

// TestRevokedRefreshTokenStaysRevoked covers the blacklist TTL being taken from the
// access-token lifetime, which un-revoked long-lived refresh tokens early.
func TestRevokedRefreshTokenStaysRevoked(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()

	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() { _ = blacklist.Close() })

	manager, err := NewJWTManager(config, blacklist)
	require.NoError(t, err)

	refreshToken, err := manager.GenerateRefreshToken(ctx, "user-1")
	require.NoError(t, err)
	require.NoError(t, manager.RevokeJWT(ctx, refreshToken))

	// The blacklist entry must outlive the access-token window, not expire with it.
	expiry := blacklist.tokens[blacklistKey(refreshToken)]
	assert.True(t, expiry.After(time.Now().Add(config.Expiration+time.Hour)),
		"blacklist entry expires at %s, before the refresh token itself does", expiry)

	_, err = manager.ValidateJWT(ctx, refreshToken)
	assert.Error(t, err)
}

// TestRevokeWithoutBlacklistFails covers the nil-pointer dereference on logout when the
// manager was built without a blacklist.
func TestRevokeWithoutBlacklistFails(t *testing.T) {
	ctx := context.Background()

	manager, err := NewJWTManager(testJWTConfig(), nil)
	require.NoError(t, err)

	token, err := manager.GenerateJWT(ctx, "user-1", map[string]any{})
	require.NoError(t, err)

	err = manager.RevokeJWT(ctx, token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no blacklist configured")
}

// TestForeignIssuerAndAudienceRejected covers iss and aud being read but never checked,
// which accepted tokens minted by any service sharing the secret.
func TestForeignIssuerAndAudienceRejected(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()
	manager := newTestJWTManager(t, config)

	sign := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(config.Secret))
		require.NoError(t, err)
		return signed
	}

	foreignIssuer := sign(jwt.MapClaims{
		"sub": "user-1",
		"iss": "attacker-issuer",
		"aud": config.Audience,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	_, err := manager.ValidateJWT(ctx, foreignIssuer)
	assert.Error(t, err, "token from a foreign issuer must be rejected")

	foreignAudience := sign(jwt.MapClaims{
		"sub": "user-1",
		"iss": config.Issuer,
		"aud": "some-other-service",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	_, err = manager.ValidateJWT(ctx, foreignAudience)
	assert.Error(t, err, "token for a foreign audience must be rejected")
}

// TestAlgorithmSubstitutionRejected verifies the signing method is pinned.
func TestAlgorithmSubstitutionRejected(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()
	manager := newTestJWTManager(t, config)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user-1",
		"iss": config.Issuer,
		"aud": config.Audience,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	unsigned, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = manager.ValidateJWT(ctx, unsigned)
	assert.Error(t, err)
}

// TestWeakSecretRejected covers NewJWTManager accepting an empty or short HMAC secret.
func TestWeakSecretRejected(t *testing.T) {
	config := testJWTConfig()
	config.Secret = "too-short"

	_, err := NewJWTManager(config, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least 32 characters")

	config.Secret = ""
	_, err = NewJWTManager(config, nil)
	assert.Error(t, err)
}

// TestGenerateJWTDoesNotMutateCallerClaims covers GenerateJWT writing into the map it
// was handed.
func TestGenerateJWTDoesNotMutateCallerClaims(t *testing.T) {
	ctx := context.Background()
	manager := newTestJWTManager(t, testJWTConfig())

	claims := map[string]any{"email": "user@example.com"}

	_, err := manager.GenerateJWT(ctx, "user-1", claims)
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"email": "user@example.com"}, claims)
}

// TestRS256SignAndValidate covers RS256 signing, which passed the raw private exponent
// to SignedString and could never produce a token.
func TestRS256SignAndValidate(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()
	config.Algorithm = "RS256"
	manager := newTestJWTManager(t, config)

	token, err := manager.GenerateJWT(ctx, "user-1", map[string]any{})
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateJWT(ctx, token)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)

	publicKey, err := manager.GetPublicKey()
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}

// TestRotateKeysKeepsPreviousTokensValid covers rotation discarding the outgoing key,
// which invalidated every token already in circulation.
func TestRotateKeysKeepsPreviousTokensValid(t *testing.T) {
	ctx := context.Background()
	config := testJWTConfig()
	config.Algorithm = "RS256"
	manager := newTestJWTManager(t, config)

	oldToken, err := manager.GenerateJWT(ctx, "user-1", map[string]any{})
	require.NoError(t, err)

	require.NoError(t, manager.RotateKeys())

	claims, err := manager.ValidateJWT(ctx, oldToken)
	require.NoError(t, err, "a token signed before rotation must stay valid until it expires")
	assert.Equal(t, "user-1", claims.UserID)

	newToken, err := manager.GenerateJWT(ctx, "user-2", map[string]any{})
	require.NoError(t, err)
	_, err = manager.ValidateJWT(ctx, newToken)
	assert.NoError(t, err)
}

// TestIsTokenValidFailsClosedWithoutStore covers the fail-open branch that accepted any
// string as a valid token when no store was configured.
func TestIsTokenValidFailsClosedWithoutStore(t *testing.T) {
	ctx := context.Background()

	blacklist := NewInMemoryTokenBlacklist()
	t.Cleanup(func() { _ = blacklist.Close() })

	manager := NewTokenManager(nil, nil, blacklist, interfaces.TokenConfig{})

	valid, err := manager.IsTokenValid(ctx, "not-a-real-token-at-all")
	require.Error(t, err)
	assert.False(t, valid)
}

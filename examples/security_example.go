package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/subosito/gotenv"
	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
	"github.com/yeferson59/go-better-auth/internal/infra/config"
	"github.com/yeferson59/go-better-auth/internal/infra/security"
)

func main() {
	_ = gotenv.Load()
	// Initialize configuration manager
	configManager, err := config.NewManager()
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// Get configuration
	authConfig := configManager.GetConfig()

	// Example 1: Using bcrypt hasher
	fmt.Println("=== Bcrypt Hasher Example ===")

	// Create hasher factory
	hashConfig := interfaces.HashConfig{
		Algorithm: "bcrypt",
		Cost:      12,
	}
	hasherFactory := security.NewHasherFactory(hashConfig)
	hasher := hasherFactory.GetDefault()

	// Hash a password
	ctx := context.Background()
	password := "mySecurePassword123!"
	hashedPassword, err := hasher.Hash(ctx, password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	fmt.Printf("Original password: %s\n", password)
	fmt.Printf("Hashed password: %s\n", hashedPassword)

	// Verify password
	isValid, err := hasher.Verify(ctx, password, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to verify password: %v", err)
	}
	fmt.Printf("Password verification: %v\n", isValid)

	// Check if hash needs rehashing
	needsRehash, err := hasher.NeedsRehash(ctx, hashedPassword)
	if err != nil {
		log.Fatalf("Failed to check if hash needs rehashing: %v", err)
	}
	fmt.Printf("Needs rehash: %v\n", needsRehash)

	// Example 2: Using JWT Manager
	fmt.Println("\n=== JWT Manager Example ===")

	// Create token blacklist
	inMemoryBlacklist := security.NewInMemoryTokenBlacklist()
	defer inMemoryBlacklist.Close()

	var blacklist interfaces.TokenBlacklist = inMemoryBlacklist

	// Create JWT manager
	jwtManager, err := security.NewJWTManager(authConfig.JWT, blacklist)
	if err != nil {
		log.Fatalf("Failed to create JWT manager: %v", err)
	}

	// Generate JWT token
	userID := "user123"
	claims := map[string]any{
		"sub":   userID,
		"email": "user@example.com",
		"roles": []string{"user", "admin"},
	}

	token, err := jwtManager.GenerateJWT(ctx, userID, claims)
	if err != nil {
		log.Fatalf("Failed to generate JWT token: %v", err)
	}
	fmt.Printf("Generated JWT token: %s\n", preview(token))

	// Validate JWT token
	tokenClaims, err := jwtManager.ValidateJWT(ctx, token)
	if err != nil {
		log.Fatalf("Failed to validate JWT token: %v", err)
	}
	fmt.Printf("Token claims - User ID: %s, Issuer: %s\n", tokenClaims.UserID, tokenClaims.Issuer)

	// Only a token minted as a refresh token can be exchanged; passing the access
	// token here is rejected by design.
	refreshToken, err := jwtManager.GenerateRefreshToken(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to generate refresh token: %v", err)
	}

	newAccessToken, newRefreshToken, err := jwtManager.RefreshJWT(ctx, refreshToken)
	if err != nil {
		log.Fatalf("Failed to refresh JWT token: %v", err)
	}
	fmt.Printf("New access token: %s\n", preview(newAccessToken))
	fmt.Printf("New refresh token: %s\n", preview(newRefreshToken))

	if _, _, err := jwtManager.RefreshJWT(ctx, newAccessToken); err != nil {
		fmt.Printf("Access token correctly rejected as a refresh token: %v\n", err)
	}

	// Revoke JWT token
	err = jwtManager.RevokeJWT(ctx, token)
	if err != nil {
		log.Fatalf("Failed to revoke JWT token: %v", err)
	}
	fmt.Println("JWT token revoked successfully")

	// Example 3: Using Token Generator
	fmt.Println("\n=== Token Generator Example ===")

	// Create secure token generator
	tokenGenerator := security.NewSecureTokenGenerator("HS256", []byte(authConfig.JWT.Secret))

	// Generate various types of tokens
	randomToken, err := tokenGenerator.GenerateRandom(32)
	if err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	fmt.Printf("Random token: %s\n", randomToken)

	secureToken, err := tokenGenerator.GenerateSecure(64)
	if err != nil {
		log.Fatalf("Failed to generate secure token: %v", err)
	}
	fmt.Printf("Secure token: %s\n", secureToken)

	uuidToken, err := tokenGenerator.GenerateUUID()
	if err != nil {
		log.Fatalf("Failed to generate UUID token: %v", err)
	}
	fmt.Printf("UUID token: %s\n", uuidToken)

	// Generate specific purpose tokens
	passwordResetToken, err := tokenGenerator.GeneratePasswordResetToken()
	if err != nil {
		log.Fatalf("Failed to generate password reset token: %v", err)
	}
	fmt.Printf("Password reset token: %s\n", preview(passwordResetToken))

	emailVerificationToken, err := tokenGenerator.GenerateEmailVerificationToken()
	if err != nil {
		log.Fatalf("Failed to generate email verification token: %v", err)
	}
	fmt.Printf("Email verification token: %s\n", preview(emailVerificationToken))

	oauthState, err := tokenGenerator.GenerateOAuthState()
	if err != nil {
		log.Fatalf("Failed to generate OAuth state: %v", err)
	}
	fmt.Printf("OAuth state: %s\n", preview(oauthState))

	// Generate backup codes
	backupCodes := make([]string, 10)
	for i := range 10 {
		backupCode, err := tokenGenerator.GenerateBackupCode()
		if err != nil {
			log.Fatalf("Failed to generate backup code: %v", err)
		}
		backupCodes[i] = backupCode
	}
	fmt.Printf("Backup codes: %d generated (values withheld)\n", len(backupCodes))

	// Example 4: Using Token Blacklist
	fmt.Println("\n=== Token Blacklist Example ===")

	// Add token to blacklist
	testToken := "test-token-123"
	expiration := time.Now().Add(time.Hour)
	err = blacklist.Add(ctx, testToken, expiration)
	if err != nil {
		log.Fatalf("Failed to add token to blacklist: %v", err)
	}
	fmt.Printf("Added token to blacklist: %s\n", testToken)

	// Check if token is blacklisted
	isBlacklisted, err := blacklist.IsBlacklisted(ctx, testToken)
	if err != nil {
		log.Fatalf("Failed to check if token is blacklisted: %v", err)
	}
	fmt.Printf("Token is blacklisted: %v\n", isBlacklisted)

	// Show blacklist size
	fmt.Printf("Blacklist size: %d\n", inMemoryBlacklist.Size())

	// Example 5: Environment Configuration
	fmt.Println("\n=== Environment Configuration Example ===")

	// Display some configuration values
	fmt.Printf("JWT Algorithm: %s\n", authConfig.JWT.Algorithm)
	fmt.Printf("JWT Expiration: %v\n", authConfig.JWT.Expiration)
	fmt.Printf("Password Algorithm: %s\n", authConfig.Password.Algorithm)
	fmt.Printf("Password Bcrypt Cost: %d\n", authConfig.Password.BcryptCost)
	fmt.Printf("Database Driver: %s\n", authConfig.Database.Driver)
	fmt.Printf("Server Host: %s\n", authConfig.Server.Host)
	fmt.Printf("Server Port: %d\n", authConfig.Server.Port)

	// Validate configuration
	err = configManager.ValidateConfig()
	if err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	fmt.Println("Configuration is valid!")

	fmt.Println("\n=== Security Utilities Example Complete ===")
}

// preview shortens a secret so the example can show that one was produced without
// printing anything usable to a terminal, a log file, or CI output.
func preview(secret string) string {
	shown := 8

	if len(secret) <= shown {
		return "<redacted>"
	}

	return secret[:shown] + "...<redacted>"
}

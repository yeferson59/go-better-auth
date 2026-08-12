package entities

import (
	"time"
)

// TokenType represents different types of tokens
type TokenType string

const (
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypeInvitation        TokenType = "invitation"
)

// DefaultTokenDurations defines default expiration times for different token types
var DefaultTokenDurations = map[TokenType]time.Duration{
	TokenTypePasswordReset:     15 * time.Minute,
	TokenTypeEmailVerification: 24 * time.Hour,
	TokenTypeInvitation:        7 * 24 * time.Hour, // 7 days
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Token     string     `json:"token"`
	Email     string     `json:"email"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UsedAt    *time.Time `json:"usedAt,omitempty"`
	User      *User      `json:"user,omitempty"`
}

// NewPasswordResetToken creates a new password reset token
func NewPasswordResetToken(userID, email, token string) *PasswordResetToken {
	return &PasswordResetToken{
		UserID:    userID,
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(DefaultTokenDurations[TokenTypePasswordReset]),
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if the password reset token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the password reset token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the password reset token is valid (not expired and not used)
func (t *PasswordResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// Use marks the password reset token as used
func (t *PasswordResetToken) Use() {
	now := time.Now()
	t.UsedAt = &now
}

// Token represents a generic token for various purposes
type Token struct {
	ID        string         `json:"id"`
	UserID    string         `json:"userId,omitempty"`
	Token     string         `json:"token"`
	Type      TokenType      `json:"type"`
	Email     string         `json:"email,omitempty"`
	Data      map[string]any `json:"data,omitempty"` // Additional token data
	ExpiresAt time.Time      `json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UsedAt    *time.Time     `json:"usedAt,omitempty"`
	User      *User          `json:"user,omitempty"`
}

// NewToken creates a token with the given type and duration
func NewToken(tokenType TokenType, userID, email, token string, duration time.Duration) *Token {
	return &Token{
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		Email:     email,
		ExpiresAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if the token has expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *Token) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (t *Token) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// Use marks the token as used
func (t *Token) Use() {
	now := time.Now()

	t.UsedAt = &now
}

// IsPasswordResetToken checks if this is a password reset token
func (t *Token) IsPasswordResetToken() bool {
	return t.Type == TokenTypePasswordReset
}

// IsEmailVerificationToken checks if this is an email verification token
func (t *Token) IsEmailVerificationToken() bool {
	return t.Type == TokenTypeEmailVerification
}

// IsInvitationToken checks if this is an invitation token
func (t *Token) IsInvitationToken() bool {
	return t.Type == TokenTypeInvitation
}

// GetData returns a value from the token data
func (t *Token) GetData(key string) (any, bool) {
	if t.Data == nil {
		return nil, false
	}

	value, exists := t.Data[key]

	return value, exists
}

// SetData sets a value in the token data
func (t *Token) SetData(key string, value any) {
	if t.Data == nil {
		t.Data = make(map[string]any)
	}

	t.Data[key] = value
}

// GetDefaultDuration returns the default duration for a token type
func GetDefaultDuration(tokenType TokenType) time.Duration {
	if duration, exists := DefaultTokenDurations[tokenType]; exists {
		return duration
	}

	return 15 * time.Minute
}

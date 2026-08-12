package entities

import (
	"time"
)

// SessionStatus represents the status of a session
type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusExpired SessionStatus = "expired"
	SessionStatusRevoked SessionStatus = "revoked"
)

// Session represents a user session
type Session struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	Token        string     `json:"token"`
	RefreshToken string     `json:"refreshToken,omitempty"`
	IPAddress    string     `json:"ipAddress,omitempty"`
	UserAgent    string     `json:"userAgent,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	RevokedAt    *time.Time `json:"revokedAt,omitempty"`
	User         *User      `json:"user,omitempty"`
}

// IsActive checks if the session is currently active
func (s *Session) IsActive() bool {
	return s.RevokedAt == nil && time.Now().Before(s.ExpiresAt)
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRevoked checks if the session has been revoked
func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}

// GetStatus returns the current status of the session
func (s *Session) GetStatus() SessionStatus {
	if s.IsRevoked() {
		return SessionStatusRevoked
	}

	if s.IsExpired() {
		return SessionStatusExpired
	}

	return SessionStatusActive
}

// Revoke revokes the session
func (s *Session) Revoke() {
	now := time.Now()

	s.RevokedAt = &now
	s.UpdatedAt = now
}

// Extend extends the session expiration time
func (s *Session) Extend(duration time.Duration) {
	now := time.Now()

	s.ExpiresAt = now.Add(duration)
	s.UpdatedAt = now
}

// SessionInfo represents session information for API responses
type SessionInfo struct {
	ID         string        `json:"id"`
	IPAddress  string        `json:"ipAddress,omitempty"`
	UserAgent  string        `json:"userAgent,omitempty"`
	Status     SessionStatus `json:"status"`
	ExpiresAt  time.Time     `json:"expiresAt"`
	CreatedAt  time.Time     `json:"createdAt"`
	LastUsedAt time.Time     `json:"lastUsedAt"`
}

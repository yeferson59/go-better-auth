package entities

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID          string      `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username,omitempty"`
	Password    string      `json:"-"`
	FirstName   string      `json:"firstName,omitempty"`
	LastName    string      `json:"lastName,omitempty"`
	IsActive    bool        `json:"isActive"`
	IsVerified  bool        `json:"isVerified"`
	LastLoginAt *time.Time  `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Roles       []Role      `json:"roles,omitempty"`
	Sessions    []Session   `json:"sessions,omitempty"`
	OAuthLinks  []OAuthLink `json:"oauthLinks,omitempty"`
}

// OAuthLink represents a link between a user and an OAuth provider
type OAuthLink struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	ProviderName string     `json:"providerName"`
	ProviderID   string     `json:"providerId"`
	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// IsEmailVerified checks if the user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.IsVerified
}

// CanLogin checks if the user can login (active and verified)
func (u *User) CanLogin() bool {
	return u.IsActive && u.IsVerified
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}

	return false
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permissionName string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}

	return false
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}

	return u.FirstName + " " + u.LastName
}

// UpdateLastLogin updates the user's last login time
func (u *User) UpdateLastLogin() {
	now := time.Now()

	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// UserProfile represents user profile information
type UserProfile struct {
	UserID      string            `json:"userId"`
	Avatar      string            `json:"avatar,omitempty"`
	Bio         string            `json:"bio,omitempty"`
	Location    string            `json:"location,omitempty"`
	Website     string            `json:"website,omitempty"`
	Preferences map[string]string `json:"preferences,omitempty"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

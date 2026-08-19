package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Username    string         `gorm:"uniqueIndex" json:"username"`
	Password    string         `gorm:"not null" json:"-"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	IsActive    bool           `gorm:"default:true" json:"isActive"`
	IsVerified  bool           `gorm:"default:false" json:"isVerified"`
	LastLoginAt *time.Time     `json:"lastLoginAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Roles               []Role               `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Sessions            []Session            `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	OAuthLinks          []OAuthLink          `gorm:"foreignKey:UserID" json:"oauthLinks,omitempty"`
	Tokens              []Token              `gorm:"foreignKey:UserID" json:"tokens,omitempty"`
	PasswordResetTokens []PasswordResetToken `gorm:"foreignKey:UserID" json:"passwordResetTokens,omitempty"`
}

// Role represents the roles table
type Role struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	IsSystem    bool           `gorm:"default:false" json:"isSystem"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission represents the permissions table
type Permission struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Resource    string         `gorm:"not null" json:"resource"`
	Action      string         `gorm:"not null" json:"action"`
	IsSystem    bool           `gorm:"default:false" json:"isSystem"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// Session represents the sessions table
type Session struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string         `gorm:"not null;index" json:"userId"`
	Token        string         `gorm:"uniqueIndex;not null" json:"token"`
	RefreshToken string         `gorm:"uniqueIndex" json:"refreshToken"`
	IPAddress    string         `json:"ipAddress"`
	UserAgent    string         `json:"userAgent"`
	ExpiresAt    time.Time      `gorm:"not null;index" json:"expiresAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	RevokedAt    *time.Time     `gorm:"index" json:"revokedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OAuthLink represents the oauth_links table
type OAuthLink struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string         `gorm:"not null;index" json:"userId"`
	ProviderName string         `gorm:"not null" json:"providerName"`
	ProviderID   string         `gorm:"not null" json:"providerId"`
	AccessToken  string         `json:"-"`
	RefreshToken string         `json:"-"`
	ExpiresAt    *time.Time     `json:"expiresAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OAuthProvider represents the oauth_providers table
type OAuthProvider struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name         string         `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName  string         `gorm:"not null" json:"displayName"`
	ClientID     string         `gorm:"not null" json:"clientId"`
	ClientSecret string         `gorm:"not null" json:"-"`
	RedirectURL  string         `gorm:"not null" json:"redirectUrl"`
	Scopes       string         `gorm:"type:text" json:"scopes"`
	AuthURL      string         `gorm:"not null" json:"authUrl"`
	TokenURL     string         `gorm:"not null" json:"tokenUrl"`
	UserInfoURL  string         `gorm:"not null" json:"userInfoUrl"`
	IsEnabled    bool           `gorm:"default:false" json:"isEnabled"`
	Config       string         `gorm:"type:text" json:"config"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Token represents the tokens table
type Token struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"index" json:"userId"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	Type      string         `gorm:"not null;index" json:"type"`
	Email     string         `gorm:"index" json:"email"`
	Data      string         `gorm:"type:text" json:"data"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UsedAt    *time.Time     `gorm:"index" json:"usedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// PasswordResetToken represents the password_reset_tokens table
type PasswordResetToken struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"not null;index" json:"userId"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	Email     string         `gorm:"not null;index" json:"email"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UsedAt    *time.Time     `gorm:"index" json:"usedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// UserRole represents the user_roles join table
type UserRole struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID     string    `gorm:"not null;index" json:"userId"`
	RoleID     string    `gorm:"not null;index" json:"roleId"`
	AssignedBy string    `json:"assignedBy"`
	CreatedAt  time.Time `json:"createdAt"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// RolePermission represents the role_permissions join table
type RolePermission struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RoleID       string    `gorm:"not null;index" json:"roleId"`
	PermissionID string    `gorm:"not null;index" json:"permissionId"`
	GrantedBy    string    `json:"grantedBy"`
	CreatedAt    time.Time `json:"createdAt"`

	// Relations
	Role       *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// TableName methods for custom table names
func (User) TableName() string               { return "users" }
func (Role) TableName() string               { return "roles" }
func (Permission) TableName() string         { return "permissions" }
func (Session) TableName() string            { return "sessions" }
func (OAuthLink) TableName() string          { return "oauth_links" }
func (OAuthProvider) TableName() string      { return "oauth_providers" }
func (Token) TableName() string              { return "tokens" }
func (PasswordResetToken) TableName() string { return "password_reset_tokens" }
func (UserRole) TableName() string           { return "user_roles" }
func (RolePermission) TableName() string     { return "role_permissions" }

package interfaces

import (
	"context"

	"github.com/yeferson59/go-better-auth/internal/core/entities"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Count(ctx context.Context) (int64, error)

	// Search operations
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)

	// User-specific operations
	GetByOAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error)
	GetWithRoles(ctx context.Context, id string) (*entities.User, error)
	GetWithSessions(ctx context.Context, id string) (*entities.User, error)

	// Bulk operations
	GetByIDs(ctx context.Context, ids []string) ([]*entities.User, error)
	UpdateLastLogin(ctx context.Context, id string) error

	// Status operations
	Activate(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
	Verify(ctx context.Context, id string) error
}

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, session *entities.Session) error
	GetByID(ctx context.Context, id string) (*entities.Session, error)
	GetByToken(ctx context.Context, token string) (*entities.Session, error)
	Update(ctx context.Context, session *entities.Session) error
	Delete(ctx context.Context, id string) error

	// User-specific operations
	GetByUserID(ctx context.Context, userID string) ([]*entities.Session, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]*entities.Session, error)

	// Session management
	RevokeByUserID(ctx context.Context, userID string) error
	RevokeByToken(ctx context.Context, token string) error
	RevokeExpired(ctx context.Context) error

	// Cleanup operations
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID string) error

	// Analytics
	CountActive(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *entities.Role) error
	GetByID(ctx context.Context, id string) (*entities.Role, error)
	GetByName(ctx context.Context, name string) (*entities.Role, error)
	Update(ctx context.Context, role *entities.Role) error
	Delete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, limit, offset int) ([]*entities.Role, error)
	Count(ctx context.Context) (int64, error)

	// Role-specific operations
	GetWithPermissions(ctx context.Context, id string) (*entities.Role, error)
	GetSystemRoles(ctx context.Context) ([]*entities.Role, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entities.Role, error)

	// Permission management
	AddPermission(ctx context.Context, roleID, permissionID string) error
	RemovePermission(ctx context.Context, roleID, permissionID string) error
	GetPermissions(ctx context.Context, roleID string) ([]*entities.Permission, error)

	// User-role management
	AssignToUser(ctx context.Context, roleID, userID string) error
	UnassignFromUser(ctx context.Context, roleID, userID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*entities.Role, error)
}

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, permission *entities.Permission) error
	GetByID(ctx context.Context, id string) (*entities.Permission, error)
	GetByName(ctx context.Context, name string) (*entities.Permission, error)
	Update(ctx context.Context, permission *entities.Permission) error
	Delete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, limit, offset int) ([]*entities.Permission, error)
	Count(ctx context.Context) (int64, error)

	// Permission-specific operations
	GetByResource(ctx context.Context, resource string) ([]*entities.Permission, error)
	GetByResourceAndAction(ctx context.Context, resource, action string) (*entities.Permission, error)
	GetSystemPermissions(ctx context.Context) ([]*entities.Permission, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entities.Permission, error)

	// Role-permission management
	GetByRoleID(ctx context.Context, roleID string) ([]*entities.Permission, error)
	GetByUserID(ctx context.Context, userID string) ([]*entities.Permission, error)
}

// OAuthProviderRepository defines the interface for OAuth provider data access
type OAuthProviderRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, provider *entities.OAuthProvider) error
	GetByID(ctx context.Context, id string) (*entities.OAuthProvider, error)
	GetByName(ctx context.Context, name string) (*entities.OAuthProvider, error)
	Update(ctx context.Context, provider *entities.OAuthProvider) error
	Delete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, limit, offset int) ([]*entities.OAuthProvider, error)
	GetEnabled(ctx context.Context) ([]*entities.OAuthProvider, error)
	Count(ctx context.Context) (int64, error)

	// Provider-specific operations
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	UpdateConfig(ctx context.Context, id string, config map[string]string) error
}

// TokenRepository defines the interface for token data access
type TokenRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, token *entities.Token) error
	GetByID(ctx context.Context, id string) (*entities.Token, error)
	GetByToken(ctx context.Context, token string) (*entities.Token, error)
	Update(ctx context.Context, token *entities.Token) error
	Delete(ctx context.Context, id string) error

	// Token-specific operations
	GetByUserID(ctx context.Context, userID string) ([]*entities.Token, error)
	GetByEmail(ctx context.Context, email string) ([]*entities.Token, error)
	GetByType(ctx context.Context, tokenType entities.TokenType) ([]*entities.Token, error)
	GetByTypeAndEmail(ctx context.Context, tokenType entities.TokenType, email string) (*entities.Token, error)

	// Token management
	MarkAsUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteByType(ctx context.Context, tokenType entities.TokenType) error

	// Validation
	IsTokenValid(ctx context.Context, token string, tokenType entities.TokenType) (bool, error)

	// Cleanup operations
	CleanupExpired(ctx context.Context) error
}

// PasswordResetTokenRepository defines the interface for password reset token data access
type PasswordResetTokenRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, token *entities.PasswordResetToken) error
	GetByID(ctx context.Context, id string) (*entities.PasswordResetToken, error)
	GetByToken(ctx context.Context, token string) (*entities.PasswordResetToken, error)
	Update(ctx context.Context, token *entities.PasswordResetToken) error
	Delete(ctx context.Context, id string) error

	// Token-specific operations
	GetByUserID(ctx context.Context, userID string) ([]*entities.PasswordResetToken, error)
	GetByEmail(ctx context.Context, email string) ([]*entities.PasswordResetToken, error)
	GetValidByToken(ctx context.Context, token string) (*entities.PasswordResetToken, error)

	// Token management
	MarkAsUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID string) error

	// Validation
	IsTokenValid(ctx context.Context, token string) (bool, error)

	// Cleanup operations
	CleanupExpired(ctx context.Context) error
}

// Repository aggregates all repository interfaces
type Repository struct {
	User               UserRepository
	Session            SessionRepository
	Role               RoleRepository
	Permission         PermissionRepository
	OAuthProvider      OAuthProviderRepository
	Token              TokenRepository
	PasswordResetToken PasswordResetTokenRepository
}

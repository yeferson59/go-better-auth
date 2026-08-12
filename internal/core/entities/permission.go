package entities

import (
	"time"
)

type Action string

const (
	Create Action = "create"
	Read   Action = "Read"
	Update Action = "update"
	Delete Action = "delete"
	List   Action = "list"
	All    Action = "all"
)

type Resource string

const (
	Users       Resource = "users"
	Roles       Resource = "roles"
	Permissions Resource = "permissions"
	Sessions    Resource = "sessions"
	System      Resource = "system"
)

type SystemPermission string

const (
	// User management
	UserCreate SystemPermission = "users:create"
	UserRead   SystemPermission = "users:read"
	UserUpdate SystemPermission = "users:update"
	UserDelete SystemPermission = "users:delete"
	UserList   SystemPermission = "users:list"

	// Role management
	RoleCreate SystemPermission = "roles:create"
	RoleRead   SystemPermission = "roles:read"
	RoleUpdate SystemPermission = "roles:update"
	RoleDelete SystemPermission = "roles:delete"
	RoleList   SystemPermission = "roles:list"

	// Permission management
	PermissionCreate SystemPermission = "permissions:create"
	PermissionRead   SystemPermission = "permissions:read"
	PermissionUpdate SystemPermission = "permissions:update"
	PermissionDelete SystemPermission = "permissions:delete"
	PermissionList   SystemPermission = "permissions:list"

	// Session management
	SessionRead   SystemPermission = "sessions:read"
	SessionDelete SystemPermission = "sessions:delete"
	SessionList   SystemPermission = "sessions:list"

	// System administration
	SystemAdmin SystemPermission = "system:admin"
)

// Permission represents a permission in the system
type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Resource    string    `json:"resource"` // e.g., "users", "posts", "settings"
	Action      string    `json:"action"`   // e.g., "read", "write", "delete"
	IsSystem    bool      `json:"isSystem"` // System permissions cannot be deleted
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Roles       []Role    `json:"roles,omitempty"`
}

// GetFullName returns the full permission name (resource:action)
func (p *Permission) GetFullName() string {
	return p.Resource + ":" + p.Action
}

// IsSystemPermission checks if this is a system permission
func (p *Permission) IsSystemPermission() bool {
	return p.IsSystem
}

// CanDelete checks if the permission can be deleted
func (p *Permission) CanDelete() bool {
	return !p.IsSystem
}

// Matches checks if the permission matches the given resource and action
func (p *Permission) Matches(resource, action string) bool {
	return p.Resource == resource && p.Action == action
}

// RolePermission represents the relationship between a role and a permission
type RolePermission struct {
	ID           string      `json:"id"`
	RoleID       string      `json:"roleId"`
	PermissionID string      `json:"permissionId"`
	GrantedBy    string      `json:"grantedBy,omitempty"` // ID of the user who granted this permission
	CreatedAt    time.Time   `json:"createdAt"`
	Role         *Role       `json:"role,omitempty"`
	Permission   *Permission `json:"permission,omitempty"`
}

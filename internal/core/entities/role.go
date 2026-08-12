package entities

import (
	"time"
)

type DefaultRole string

const (
	Admin     DefaultRole = "admin"
	Default   DefaultRole = "user"
	Moderator DefaultRole = "moderator"
	Guest     DefaultRole = "guest"
)

// Role represents a role in the system
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsSystem    bool         `json:"isSystem"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Permissions []Permission `json:"permissions,omitempty"`
	Users       []User       `json:"users,omitempty"`
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, permission := range r.Permissions {
		if permission.Name == permissionName {
			return true
		}
	}

	return false
}

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission Permission) bool {
	if !r.HasPermission(permission.Name) {
		r.Permissions = append(r.Permissions, permission)
		r.UpdatedAt = time.Now()

		return true
	}

	return false
}

// RemovePermission removes a permission from the role
func (r *Role) RemovePermission(permissionName string) {
	for i, permission := range r.Permissions {
		if permission.Name == permissionName {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.UpdatedAt = time.Now()

			break
		}
	}
}

// IsSystemRole checks if this is a system role
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// CanDelete checks if the role can be deleted
func (r *Role) CanDelete() bool {
	return !r.IsSystem
}

// GetPermissionNames returns a slice of permission names
func (r *Role) GetPermissionNames() []string {
	names := make([]string, len(r.Permissions))

	for i, permission := range r.Permissions {
		names[i] = permission.Name
	}

	return names
}

// UserRole represents the relationship between a user and a role
type UserRole struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	RoleID     string    `json:"roleId"`
	AssignedBy string    `json:"assignedBy,omitempty"` // ID of the user who assigned this role
	CreatedAt  time.Time `json:"createdAt"`
	User       *User     `json:"user,omitempty"`
	Role       *Role     `json:"role,omitempty"`
}

package models

import (
	"encoding/json"

	"github.com/yeferson59/go-better-auth/internal/core/entities"
)

// ToUserEntity converts a GORM User model to an entity
func (u *User) ToEntity() *entities.User {
	entity := &entities.User{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		Password:    u.Password,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		IsActive:    u.IsActive,
		IsVerified:  u.IsVerified,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	// Convert roles
	for _, role := range u.Roles {
		entity.Roles = append(entity.Roles, *role.ToEntity())
	}

	// Convert sessions
	for _, session := range u.Sessions {
		entity.Sessions = append(entity.Sessions, *session.ToEntity())
	}

	// Convert OAuth links
	for _, link := range u.OAuthLinks {
		entity.OAuthLinks = append(entity.OAuthLinks, *link.ToEntity())
	}

	return entity
}

// FromUserEntity converts an entity to a GORM User model
func (u *User) FromEntity(entity *entities.User) {
	u.ID = entity.ID
	u.Email = entity.Email
	u.Username = entity.Username
	u.Password = entity.Password
	u.FirstName = entity.FirstName
	u.LastName = entity.LastName
	u.IsActive = entity.IsActive
	u.IsVerified = entity.IsVerified
	u.LastLoginAt = entity.LastLoginAt
	u.CreatedAt = entity.CreatedAt
	u.UpdatedAt = entity.UpdatedAt
}

// ToRoleEntity converts a GORM Role model to an entity
func (r *Role) ToEntity() *entities.Role {
	entity := &entities.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	// Convert permissions
	for _, permission := range r.Permissions {
		entity.Permissions = append(entity.Permissions, *permission.ToEntity())
	}

	return entity
}

// FromRoleEntity converts an entity to a GORM Role model
func (r *Role) FromEntity(entity *entities.Role) {
	r.ID = entity.ID
	r.Name = entity.Name
	r.Description = entity.Description
	r.IsSystem = entity.IsSystem
	r.CreatedAt = entity.CreatedAt
	r.UpdatedAt = entity.UpdatedAt
}

// ToPermissionEntity converts a GORM Permission model to an entity
func (p *Permission) ToEntity() *entities.Permission {
	return &entities.Permission{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Resource:    p.Resource,
		Action:      p.Action,
		IsSystem:    p.IsSystem,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromPermissionEntity converts an entity to a GORM Permission model
func (p *Permission) FromEntity(entity *entities.Permission) {
	p.ID = entity.ID
	p.Name = entity.Name
	p.Description = entity.Description
	p.Resource = entity.Resource
	p.Action = entity.Action
	p.IsSystem = entity.IsSystem
	p.CreatedAt = entity.CreatedAt
	p.UpdatedAt = entity.UpdatedAt
}

// ToSessionEntity converts a GORM Session model to an entity
func (s *Session) ToEntity() *entities.Session {
	return &entities.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		Token:        s.Token,
		RefreshToken: s.RefreshToken,
		IPAddress:    s.IPAddress,
		UserAgent:    s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
		RevokedAt:    s.RevokedAt,
	}
}

// FromSessionEntity converts an entity to a GORM Session model
func (s *Session) FromEntity(entity *entities.Session) {
	s.ID = entity.ID
	s.UserID = entity.UserID
	s.Token = entity.Token
	s.RefreshToken = entity.RefreshToken
	s.IPAddress = entity.IPAddress
	s.UserAgent = entity.UserAgent
	s.ExpiresAt = entity.ExpiresAt
	s.CreatedAt = entity.CreatedAt
	s.UpdatedAt = entity.UpdatedAt
	s.RevokedAt = entity.RevokedAt
}

// ToOAuthLinkEntity converts a GORM OAuthLink model to an entity
func (o *OAuthLink) ToEntity() *entities.OAuthLink {
	return &entities.OAuthLink{
		ID:           o.ID,
		UserID:       o.UserID,
		ProviderName: o.ProviderName,
		ProviderID:   o.ProviderID,
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		ExpiresAt:    o.ExpiresAt,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

// FromOAuthLinkEntity converts an entity to a GORM OAuthLink model
func (o *OAuthLink) FromEntity(entity *entities.OAuthLink) {
	o.ID = entity.ID
	o.UserID = entity.UserID
	o.ProviderName = entity.ProviderName
	o.ProviderID = entity.ProviderID
	o.AccessToken = entity.AccessToken
	o.RefreshToken = entity.RefreshToken
	o.ExpiresAt = entity.ExpiresAt
	o.CreatedAt = entity.CreatedAt
	o.UpdatedAt = entity.UpdatedAt
}

// ToOAuthProviderEntity converts a GORM OAuthProvider model to an entity
func (o *OAuthProvider) ToEntity() *entities.OAuthProvider {
	entity := &entities.OAuthProvider{
		ID:           o.ID,
		Name:         entities.Provider(o.Name),
		DisplayName:  o.DisplayName,
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		RedirectURL:  o.RedirectURL,
		AuthURL:      o.AuthURL,
		TokenURL:     o.TokenURL,
		UserInfoURL:  o.UserInfoURL,
		IsEnabled:    o.IsEnabled,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}

	// Convert scopes from JSON string to slice
	if o.Scopes != "" {
		json.Unmarshal([]byte(o.Scopes), &entity.Scopes)
	}

	// Convert config from JSON string to map
	if o.Config != "" {
		json.Unmarshal([]byte(o.Config), &entity.Config)
	}

	return entity
}

// FromOAuthProviderEntity converts an entity to a GORM OAuthProvider model
func (o *OAuthProvider) FromEntity(entity *entities.OAuthProvider) {
	o.ID = entity.ID
	o.Name = string(entity.Name)
	o.DisplayName = entity.DisplayName
	o.ClientID = entity.ClientID
	o.ClientSecret = entity.ClientSecret
	o.RedirectURL = entity.RedirectURL
	o.AuthURL = entity.AuthURL
	o.TokenURL = entity.TokenURL
	o.UserInfoURL = entity.UserInfoURL
	o.IsEnabled = entity.IsEnabled
	o.CreatedAt = entity.CreatedAt
	o.UpdatedAt = entity.UpdatedAt

	// Convert scopes from slice to JSON string
	if entity.Scopes != nil {
		if scopesJSON, err := json.Marshal(entity.Scopes); err == nil {
			o.Scopes = string(scopesJSON)
		}
	}

	// Convert config from map to JSON string
	if entity.Config != nil {
		if configJSON, err := json.Marshal(entity.Config); err == nil {
			o.Config = string(configJSON)
		}
	}
}

// ToTokenEntity converts a GORM Token model to an entity
func (t *Token) ToEntity() *entities.Token {
	entity := &entities.Token{
		ID:        t.ID,
		UserID:    t.UserID,
		Token:     t.Token,
		Type:      entities.TokenType(t.Type),
		Email:     t.Email,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
		UsedAt:    t.UsedAt,
	}

	// Convert data from JSON string to map
	if t.Data != "" {
		json.Unmarshal([]byte(t.Data), &entity.Data)
	}

	return entity
}

// FromTokenEntity converts an entity to a GORM Token model
func (t *Token) FromEntity(entity *entities.Token) {
	t.ID = entity.ID
	t.UserID = entity.UserID
	t.Token = entity.Token
	t.Type = string(entity.Type)
	t.Email = entity.Email
	t.ExpiresAt = entity.ExpiresAt
	t.CreatedAt = entity.CreatedAt
	t.UsedAt = entity.UsedAt

	// Convert data from map to JSON string
	if entity.Data != nil {
		if dataJSON, err := json.Marshal(entity.Data); err == nil {
			t.Data = string(dataJSON)
		}
	}
}

// ToPasswordResetTokenEntity converts a GORM PasswordResetToken model to an entity
func (p *PasswordResetToken) ToEntity() *entities.PasswordResetToken {
	return &entities.PasswordResetToken{
		ID:        p.ID,
		UserID:    p.UserID,
		Token:     p.Token,
		Email:     p.Email,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
		UsedAt:    p.UsedAt,
	}
}

// FromPasswordResetTokenEntity converts an entity to a GORM PasswordResetToken model
func (p *PasswordResetToken) FromEntity(entity *entities.PasswordResetToken) {
	p.ID = entity.ID
	p.UserID = entity.UserID
	p.Token = entity.Token
	p.Email = entity.Email
	p.ExpiresAt = entity.ExpiresAt
	p.CreatedAt = entity.CreatedAt
	p.UsedAt = entity.UsedAt
}

// Helper functions for batch conversions
func ToUserEntities(models []*User) []*entities.User {
	entities := make([]*entities.User, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToRoleEntities(models []*Role) []*entities.Role {
	entities := make([]*entities.Role, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToPermissionEntities(models []*Permission) []*entities.Permission {
	entities := make([]*entities.Permission, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToSessionEntities(models []*Session) []*entities.Session {
	entities := make([]*entities.Session, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToTokenEntities(models []*Token) []*entities.Token {
	entities := make([]*entities.Token, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToPasswordResetTokenEntities(models []*PasswordResetToken) []*entities.PasswordResetToken {
	entities := make([]*entities.PasswordResetToken, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

func ToOAuthProviderEntities(models []*OAuthProvider) []*entities.OAuthProvider {
	entities := make([]*entities.OAuthProvider, len(models))

	for i, model := range models {
		entities[i] = model.ToEntity()
	}

	return entities
}

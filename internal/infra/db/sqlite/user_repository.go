package sqlite

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"
	"uuid"

	"gorm.io/gorm"

	"github.com/yeferson59/go-better-auth/internal/core/entities"
	apperrors "github.com/yeferson59/go-better-auth/internal/core/errors"
	"github.com/yeferson59/go-better-auth/internal/infra/db/models"
)

// UserRepository implements the UserRepository interface using SQLite
type UserRepository struct {
	db *gorm.DB
}

// Lookups return a not-found domain error rather than (nil, nil) so that a caller
// checking only err cannot mistake a missing user for a successful lookup.
//
// Those errors carry the looked-up identifier in their Details field. Authentication
// handlers must translate them into the generic invalid-credentials error before they
// reach a client, or the response becomes a user-enumeration oracle.
// NewUserRepository creates a new SQLite user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	if user.ID == "" {
		id := uuid.NewV7()

		user.ID = id.String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	var model models.User
	model.FromEntity(user)

	if err := gorm.G[models.User](r.db).Create(ctx, &model); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	*user = *model.ToEntity()

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	model, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(id)
		}

		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	model, err := gorm.G[models.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(email)
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	model, err := gorm.G[models.User](r.db).Where("username = ?", username).First(ctx)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(username)
		}

		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return model.ToEntity(), nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	user.UpdatedAt = time.Now()

	var model models.User

	model.FromEntity(user)

	rowsAffected, err := gorm.G[models.User](r.db).Updates(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if rowsAffected == 0 {
		return stderrors.New("failed to update user")
	}

	*user = *model.ToEntity()

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	rowsAffected, err := gorm.G[models.User](r.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rowsAffected == 0 {
		return stderrors.New("failed to delete user")
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	models, err := gorm.G[models.User](r.db).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*entities.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, nil
}

// Count returns the total count of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := gorm.G[models.User](r.db).Count(ctx, "id")
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// Search searches for users by query (using LIKE for SQLite)
func (r *UserRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
	searchQuery := "%" + query + "%"
	models, err := gorm.G[models.User](r.db).Where("email LIKE ? OR username LIKE ? OR first_name LIKE ? OR last_name LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery).Limit(limit).Offset(offset).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	users := make([]*entities.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, nil
}

// GetByOAuthProvider retrieves a user by OAuth provider
func (r *UserRepository) GetByOAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error) {
	var model models.User
	if err := r.db.WithContext(ctx).
		Joins("JOIN oauth_links ON users.id = oauth_links.user_id").
		Where("oauth_links.provider_name = ? AND oauth_links.provider_id = ?", provider, providerID).
		First(&model).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(providerID)
		}
		return nil, fmt.Errorf("failed to get user by OAuth provider: %w", err)
	}

	return model.ToEntity(), nil
}

// GetWithRoles retrieves a user with their roles
func (r *UserRepository) GetWithRoles(ctx context.Context, id string) (*entities.User, error) {
	model, err := gorm.G[models.User](r.db).Preload("Roles", func(_ gorm.PreloadBuilder) error { return nil }).Where("id = ?", id).First(ctx)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(id)
		}

		return nil, fmt.Errorf("failed to get user with roles: %w", err)
	}

	return model.ToEntity(), nil
}

// GetWithSessions retrieves a user with their sessions
func (r *UserRepository) GetWithSessions(ctx context.Context, id string) (*entities.User, error) {
	model, err := gorm.G[models.User](r.db).Preload("Sessions", func(_ gorm.PreloadBuilder) error { return nil }).Where("id = ?", id).First(ctx)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUserNotFoundError(id)
		}

		return nil, fmt.Errorf("failed to get user with sessions: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByIDs retrieves users by their IDs
func (r *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]*entities.User, error) {
	models, err := gorm.G[models.User](r.db).Where("id IN ?", ids).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	users := make([]*entities.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, nil
}

// UpdateLastLogin updates the user's last login time
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	rowsAffected, err := gorm.G[models.User](r.db).Where("id = ?", id).Updates(ctx, models.User{LastLoginAt: &now, UpdatedAt: now})
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	if rowsAffected == 0 {
		return stderrors.New("failed to update last login")
	}

	return nil
}

// Activate activates a user
func (r *UserRepository) Activate(ctx context.Context, id string) error {
	rowsAffected, err := gorm.G[models.User](r.db).Where("id = ?", id).Updates(ctx, models.User{IsActive: true, UpdatedAt: time.Now()})
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	if rowsAffected == 0 {
		return stderrors.New("failed to activate user")
	}

	return nil
}

// Deactivate deactivates a user
func (r *UserRepository) Deactivate(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rowsAffected, err := gorm.G[models.User](tx).Where("id = ?", id).Update(ctx, "is_active", false)
		if err != nil {
			return fmt.Errorf("failed to deactivate user: %w", err)
		}

		if rowsAffected == 0 {
			return stderrors.New("failed to deactivate user")
		}

		rowsAffected, err = gorm.G[models.User](tx).Where("id = ?", id).Update(ctx, "updated_at", time.Now())
		if err != nil {
			return fmt.Errorf("failed to deactivate user: %w", err)
		}

		if rowsAffected == 0 {
			return stderrors.New("failed to deactivate user")
		}

		return nil
	})
}

// Verify verifies a user's email
func (r *UserRepository) Verify(ctx context.Context, id string) error {
	rowsAffected, err := gorm.G[models.User](r.db).Where("id = ?", id).Updates(ctx, models.User{IsVerified: true, UpdatedAt: time.Now()})
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	if rowsAffected == 0 {
		return stderrors.New("failed to verify user")
	}

	return nil
}

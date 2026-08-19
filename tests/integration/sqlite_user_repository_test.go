package integration

import (
	"context"
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yeferson59/go-better-auth/internal/core/entities"
	apperrors "github.com/yeferson59/go-better-auth/internal/core/errors"
	"github.com/yeferson59/go-better-auth/internal/infra/db/models"
	sqliterepo "github.com/yeferson59/go-better-auth/internal/infra/db/sqlite"
)

func TestSQLiteUserRepository(t *testing.T) {
	ctx := context.Background()

	// Setup in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Session{},
		&models.OAuthLink{},
		&models.OAuthProvider{},
		&models.Token{},
		&models.PasswordResetToken{},
	)
	require.NoError(t, err)

	// Create repository
	repo := sqliterepo.NewUserRepository(db)

	// Run tests
	t.Run("Create User", func(t *testing.T) {
		user := &entities.User{
			Email:      "test@example.com",
			Username:   "testuser",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("Get User by ID", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "test2@example.com",
			Username:   "testuser2",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Get user by ID
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.Username, foundUser.Username)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "test3@example.com",
			Username:   "testuser3",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Get user by email
		foundUser, err := repo.GetByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("Get User by Username", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "test4@example.com",
			Username:   "testuser4",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Get user by username
		foundUser, err := repo.GetByUsername(ctx, user.Username)
		require.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
	})

	t.Run("Update User", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "test5@example.com",
			Username:   "testuser5",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update user
		user.FirstName = "Updated"
		user.IsVerified = true
		err = repo.Update(ctx, user)
		require.NoError(t, err)

		// Verify update
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", foundUser.FirstName)
		assert.True(t, foundUser.IsVerified)
	})

	t.Run("Delete User", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "test6@example.com",
			Username:   "testuser6",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Delete user
		err = repo.Delete(ctx, user.ID)
		require.NoError(t, err)

		// Verify deletion: a missing row is reported as a not-found domain error, not
		// as a nil user with a nil error.
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.Error(t, err)
		assert.True(t, stderrors.Is(err, &apperrors.Error{Code: apperrors.ErrUserNotFound}))
		assert.Nil(t, foundUser)
	})

	t.Run("List Users", func(t *testing.T) {
		// Create multiple users
		for i := range 5 {
			user := &entities.User{
				Email:      fmt.Sprintf("listtest%d@example.com", i),
				Username:   fmt.Sprintf("listuser%d", i),
				Password:   "hashedpassword",
				FirstName:  "Test",
				LastName:   "User",
				IsActive:   true,
				IsVerified: false,
			}
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		// List users
		users, err := repo.List(ctx, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 5)
	})

	t.Run("Count Users", func(t *testing.T) {
		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Greater(t, count, int64(0))
	})

	t.Run("Search Users", func(t *testing.T) {
		// Create user with specific data
		user := &entities.User{
			Email:      "searchme@example.com",
			Username:   "searchuser",
			Password:   "hashedpassword",
			FirstName:  "Searchable",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Search users
		users, err := repo.Search(ctx, "searchme", 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)

		found := false
		for _, u := range users {
			if u.Email == "searchme@example.com" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Update Last Login", func(t *testing.T) {
		// Create user first
		user := &entities.User{
			Email:      "logintest@example.com",
			Username:   "loginuser",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update last login
		err = repo.UpdateLastLogin(ctx, user.ID)
		require.NoError(t, err)

		// Verify update
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.NotNil(t, foundUser.LastLoginAt)
	})

	t.Run("Activate User", func(t *testing.T) {
		// Create inactive user
		user := &entities.User{
			Email:      "activatetest@example.com",
			Username:   "activateuser",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   false,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Activate user
		err = repo.Activate(ctx, user.ID)
		require.NoError(t, err)

		// Verify activation
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.True(t, foundUser.IsActive)
	})

	t.Run("Deactivate User", func(t *testing.T) {
		// Create active user
		user := &entities.User{
			Email:      "deactivatetest@example.com",
			Username:   "deactivateuser",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Deactivate user
		err = repo.Deactivate(ctx, user.ID)
		require.NoError(t, err)

		// Verify deactivation
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.False(t, foundUser.IsActive)
	})

	t.Run("Verify User", func(t *testing.T) {
		// Create unverified user
		user := &entities.User{
			Email:      "verifytest@example.com",
			Username:   "verifyuser",
			Password:   "hashedpassword",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
			IsVerified: false,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Verify user
		err = repo.Verify(ctx, user.ID)
		require.NoError(t, err)

		// Verify verification
		foundUser, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.True(t, foundUser.IsVerified)
	})
}

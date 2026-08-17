package security

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

// BcryptHasher implements the Hasher interface using bcrypt algorithm
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher instance
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{
		cost: cost,
	}
}

// Hash generates a bcrypt hash from the given password
func (h *BcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// Verify verifies if the given password matches the bcrypt hash
func (h *BcryptHasher) Verify(ctx context.Context, password, hash string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify password: %w", err)
	}

	return true, nil
}

// NeedsRehash checks if the bcrypt hash needs to be rehashed due to cost changes
func (h *BcryptHasher) NeedsRehash(ctx context.Context, hash string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return false, fmt.Errorf("failed to get hash cost: %w", err)
	}

	return cost != h.cost, nil
}

// GetCost returns the current bcrypt cost
func (h *BcryptHasher) GetCost() int {
	return h.cost
}

// SetCost sets the bcrypt cost
func (h *BcryptHasher) SetCost(cost int) {
	if cost >= bcrypt.MinCost && cost <= bcrypt.MaxCost {
		h.cost = cost
	}
}

// ValidateHash validates if the given string is a valid bcrypt hash
func (h *BcryptHasher) ValidateHash(hash string) error {
	if len(hash) != 60 {
		return fmt.Errorf("invalid bcrypt hash length: expected 60, got %d", len(hash))
	}

	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") && !strings.HasPrefix(hash, "$2x$") && !strings.HasPrefix(hash, "$2y$") {
		return fmt.Errorf("invalid bcrypt hash format")
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 4 {
		return fmt.Errorf("invalid bcrypt hash format: expected 4 parts, got %d", len(parts))
	}

	// Validate cost
	cost, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid bcrypt cost: %w", err)
	}

	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return fmt.Errorf("invalid bcrypt cost: %d (must be between %d and %d)", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	return nil
}

// HasherFactory implements the HasherFactory interface
type HasherFactory struct {
	defaultHasher interfaces.Hasher
	config        interfaces.HashConfig
}

// NewHasherFactory creates a new HasherFactory instance
func NewHasherFactory(config interfaces.HashConfig) *HasherFactory {
	var defaultHasher interfaces.Hasher

	switch interfaces.HashAlgorithm(config.Algorithm) {
	case interfaces.HashAlgorithmBcrypt:
		defaultHasher = NewBcryptHasher(config.Cost)
	default:
		// Default to bcrypt if algorithm is not recognized
		defaultHasher = NewBcryptHasher(bcrypt.DefaultCost)
	}

	return &HasherFactory{
		defaultHasher: defaultHasher,
		config:        config,
	}
}

// Create creates a new hasher with the given configuration
func (f *HasherFactory) Create(config interfaces.HashConfig) (interfaces.Hasher, error) {
	switch interfaces.HashAlgorithm(config.Algorithm) {
	case interfaces.HashAlgorithmBcrypt:
		return NewBcryptHasher(config.Cost), nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", config.Algorithm)
	}
}

// GetDefault returns the default hasher
func (f *HasherFactory) GetDefault() interfaces.Hasher {
	return f.defaultHasher
}

// GetByAlgorithm returns a hasher for the specified algorithm
func (f *HasherFactory) GetByAlgorithm(algorithm interfaces.HashAlgorithm) (interfaces.Hasher, error) {
	switch algorithm {
	case interfaces.HashAlgorithmBcrypt:
		return NewBcryptHasher(f.config.Cost), nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}
}

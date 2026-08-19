package interfaces

import (
	"context"
)

// Hasher defines the interface for password hashing operations
type Hasher interface {
	// Hash generates a hash from the given password
	Hash(ctx context.Context, password string) (string, error)

	// Verify verifies if the given password matches the hash
	Verify(ctx context.Context, password, hash string) (bool, error)

	// NeedsRehash checks if the hash needs to be rehashed (e.g., due to cost changes)
	NeedsRehash(ctx context.Context, hash string) (bool, error)
}

// HashConfig represents configuration for password hashing
type HashConfig struct {
	// Algorithm specifies the hashing algorithm to use
	Algorithm string `json:"algorithm"`

	// Cost specifies the cost parameter for algorithms like bcrypt
	Cost int `json:"cost"`

	// Memory specifies the memory parameter for algorithms like argon2
	Memory int `json:"memory"`

	// Time specifies the time parameter for algorithms like argon2
	Time int `json:"time"`

	// Threads specifies the number of threads for algorithms like argon2
	Threads int `json:"threads"`

	// SaltLength specifies the salt length for algorithms like argon2
	SaltLength int `json:"saltLength"`

	// KeyLength specifies the key length for algorithms like argon2
	KeyLength int `json:"keyLength"`
}

// HashAlgorithm represents supported hashing algorithms
type HashAlgorithm string

const (
	// HashAlgorithmBcrypt uses bcrypt algorithm
	HashAlgorithmBcrypt HashAlgorithm = "bcrypt"

	// HashAlgorithmArgon2 uses argon2 algorithm
	HashAlgorithmArgon2 HashAlgorithm = "argon2"

	// HashAlgorithmScrypt uses scrypt algorithm
	HashAlgorithmScrypt HashAlgorithm = "scrypt"
)

// DefaultHashConfigs provides default configurations for different algorithms
var DefaultHashConfigs = map[HashAlgorithm]HashConfig{
	HashAlgorithmBcrypt: {
		Algorithm: string(HashAlgorithmBcrypt),
		Cost:      12, // bcrypt cost
	},
	HashAlgorithmArgon2: {
		Algorithm:  string(HashAlgorithmArgon2),
		Memory:     64 * 1024, // 64 MB
		Time:       3,         // 3 iterations
		Threads:    4,         // 4 threads
		SaltLength: 16,        // 16 bytes salt
		KeyLength:  32,        // 32 bytes key
	},
	HashAlgorithmScrypt: {
		Algorithm: string(HashAlgorithmScrypt),
		Cost:      32768, // N parameter
		Memory:    8,     // r parameter
		Threads:   1,     // p parameter
		KeyLength: 32,    // key length
	},
}

// HasherFactory defines the interface for creating hashers
type HasherFactory interface {
	// Create creates a new hasher with the given configuration
	Create(config HashConfig) (Hasher, error)

	// GetDefault returns the default hasher
	GetDefault() Hasher

	// GetByAlgorithm returns a hasher for the specified algorithm
	GetByAlgorithm(algorithm HashAlgorithm) (Hasher, error)
}

// PasswordStrengthValidator defines the interface for password strength validation
type PasswordStrengthValidator interface {
	// Validate validates the password strength
	Validate(ctx context.Context, password string) error

	// GetRequirements returns the password requirements
	GetRequirements() PasswordRequirements
}

// PasswordRequirements defines password requirements
type PasswordRequirements struct {
	// MinLength specifies the minimum password length
	MinLength int `json:"minLength"`

	// MaxLength specifies the maximum password length
	MaxLength int `json:"maxLength"`

	// RequireUppercase specifies if uppercase letters are required
	RequireUppercase bool `json:"requireUppercase"`

	// RequireLowercase specifies if lowercase letters are required
	RequireLowercase bool `json:"requireLowercase"`

	// RequireNumbers specifies if numbers are required
	RequireNumbers bool `json:"requireNumbers"`

	// RequireSymbols specifies if symbols are required
	RequireSymbols bool `json:"requireSymbols"`

	// ForbiddenPasswords is a list of forbidden passwords
	ForbiddenPasswords []string `json:"forbiddenPasswords"`

	// ForbiddenPatterns is a list of forbidden patterns (regex)
	ForbiddenPatterns []string `json:"forbiddenPatterns"`
}

// DefaultPasswordRequirements provides default password requirements
var DefaultPasswordRequirements = PasswordRequirements{
	MinLength:        8,
	MaxLength:        128,
	RequireUppercase: true,
	RequireLowercase: true,
	RequireNumbers:   true,
	RequireSymbols:   true,
	ForbiddenPasswords: []string{
		"password",
		"123456",
		"qwerty",
		"admin",
		"letmein",
		"welcome",
		"monkey",
		"dragon",
	},
	ForbiddenPatterns: []string{
		`^(.)\1+$`, // repeated characters
		`^(012|123|234|345|456|567|678|789|890)+$`,                                                             // sequential numbers
		`^(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)+$`, // sequential letters
	},
}

// PasswordHistory defines the interface for password history management
type PasswordHistory interface {
	// Add adds a password hash to the history
	Add(ctx context.Context, userID, passwordHash string) error

	// Check checks if a password has been used before
	Check(ctx context.Context, userID, passwordHash string) (bool, error)

	// GetHistory gets the password history for a user
	GetHistory(ctx context.Context, userID string, limit int) ([]string, error)

	// Cleanup cleans up old password history entries
	Cleanup(ctx context.Context, userID string, keepCount int) error

	// Delete deletes all password history for a user
	Delete(ctx context.Context, userID string) error
}

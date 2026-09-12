package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

// blacklistKey derives the storage key for a token. Tokens are hashed so that live
// credentials are never held in memory (or in Redis) in plaintext, and so that the key
// size stays constant regardless of how large the token is.
func blacklistKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// InMemoryTokenBlacklist implements the TokenBlacklist interface using in-memory storage
type InMemoryTokenBlacklist struct {
	tokens    map[string]time.Time
	mu        sync.RWMutex
	done      chan struct{}
	closeOnce sync.Once
}

// NewInMemoryTokenBlacklist creates a new in-memory token blacklist.
// The caller is responsible for calling Close to stop the background cleanup goroutine.
func NewInMemoryTokenBlacklist() *InMemoryTokenBlacklist {
	blacklist := &InMemoryTokenBlacklist{
		tokens: make(map[string]time.Time),
		done:   make(chan struct{}),
	}

	// Start cleanup routine
	go blacklist.cleanupRoutine()

	return blacklist
}

// Close stops the background cleanup goroutine. It is safe to call more than once.
func (b *InMemoryTokenBlacklist) Close() error {
	b.closeOnce.Do(func() {
		close(b.done)
	})
	return nil
}

// Add adds a token to the blacklist
func (b *InMemoryTokenBlacklist) Add(ctx context.Context, token string, expiration time.Time) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens[blacklistKey(token)] = expiration
	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (b *InMemoryTokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	expiration, exists := b.tokens[blacklistKey(token)]
	if !exists {
		return false, nil
	}

	// An expired entry is reported as not blacklisted and left in place for
	// cleanupRoutine to purge. Deleting it here would be a map write under a read
	// lock, which races with concurrent readers and crashes the process.
	if time.Now().After(expiration) {
		return false, nil
	}

	return true, nil
}

// Remove removes a token from the blacklist
func (b *InMemoryTokenBlacklist) Remove(ctx context.Context, token string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.tokens, blacklistKey(token))
	return nil
}

// Cleanup removes expired blacklisted tokens
func (b *InMemoryTokenBlacklist) Cleanup(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	for token, expiration := range b.tokens {
		if now.After(expiration) {
			delete(b.tokens, token)
		}
	}

	return nil
}

// cleanupRoutine runs periodic cleanup of expired tokens
func (b *InMemoryTokenBlacklist) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = b.Cleanup(context.Background())
		case <-b.done:
			return
		}
	}
}

// Size returns the number of blacklisted tokens
func (b *InMemoryTokenBlacklist) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.tokens)
}

// Clear removes all tokens from the blacklist
func (b *InMemoryTokenBlacklist) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens = make(map[string]time.Time)
}

// RedisTokenBlacklist implements the TokenBlacklist interface using Redis
type RedisTokenBlacklist struct {
	client interfaces.RedisClient
	prefix string
}

// NewRedisTokenBlacklist creates a new Redis-based token blacklist
func NewRedisTokenBlacklist(client interfaces.RedisClient, prefix string) *RedisTokenBlacklist {
	if prefix == "" {
		prefix = "blacklist:"
	}

	return &RedisTokenBlacklist{
		client: client,
		prefix: prefix,
	}
}

// Add adds a token to the blacklist
func (b *RedisTokenBlacklist) Add(ctx context.Context, token string, expiration time.Time) error {
	key := b.prefix + blacklistKey(token)
	duration := time.Until(expiration)

	if duration <= 0 {
		return fmt.Errorf("expiration time must be in the future")
	}

	return b.client.Set(ctx, key, "1", duration)
}

// IsBlacklisted checks if a token is blacklisted
func (b *RedisTokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := b.prefix + blacklistKey(token)

	exists, err := b.client.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check token existence: %w", err)
	}

	return exists, nil
}

// Remove removes a token from the blacklist
func (b *RedisTokenBlacklist) Remove(ctx context.Context, token string) error {
	key := b.prefix + blacklistKey(token)
	return b.client.Del(ctx, key)
}

// Cleanup removes expired blacklisted tokens (Redis handles this automatically)
func (b *RedisTokenBlacklist) Cleanup(_ context.Context) error {
	// Redis automatically removes expired keys, so this is a no-op
	return nil
}

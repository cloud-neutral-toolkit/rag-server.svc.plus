package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"rag-server/internal/cache"
)

// TokenCache stores validated access tokens in a cache store.
type TokenCache struct {
	store  cache.Store
	prefix string
}

// NewTokenCache creates a token cache backed by the provided store.
func NewTokenCache(store cache.Store) *TokenCache {
	return &TokenCache{store: store, prefix: "auth:token:"}
}

// Get retrieves cached claims for an access token.
func (c *TokenCache) Get(ctx context.Context, token string) (*Claims, bool, error) {
	if c == nil || c.store == nil {
		return nil, false, nil
	}
	key := c.key(token)
	value, ok, err := c.store.Get(ctx, key)
	if err != nil || !ok {
		return nil, ok, err
	}
	var claims Claims
	if err := json.Unmarshal([]byte(value), &claims); err != nil {
		_ = c.store.Delete(ctx, key)
		return nil, false, err
	}
	return &claims, true, nil
}

// Set caches claims for an access token.
func (c *TokenCache) Set(ctx context.Context, token string, claims *Claims, ttl time.Duration) error {
	if c == nil || c.store == nil {
		return nil
	}
	if claims == nil {
		return nil
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return c.store.Set(ctx, c.key(token), string(payload), ttl)
}

func (c *TokenCache) key(token string) string {
	hash := sha256.Sum256([]byte(token))
	return c.prefix + hex.EncodeToString(hash[:])
}

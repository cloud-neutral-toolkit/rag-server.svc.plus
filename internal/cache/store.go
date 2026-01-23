package cache

import (
	"context"
	"time"
)

// Store defines the minimal cache interface needed by services.
type Store interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	PurgeExpired(ctx context.Context) (int64, error)
}

// Options configure cache storage behavior.
type Options struct {
	Table      string
	DefaultTTL time.Duration
}

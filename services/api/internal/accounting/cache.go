package accounting

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache is a minimal read-through cache for expensive report queries. Redis-backed so
// multiple API instances share it. A nil cache disables caching (safe fallback).
type Cache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCache(rdb *redis.Client) *Cache {
	if rdb == nil {
		return nil
	}
	return &Cache{rdb: rdb, ttl: 30 * time.Second}
}

// GetOrCompute returns a cached value or computes+caches it.
func (c *Cache) GetOrCompute(ctx context.Context, key string, compute func() (any, error)) (any, error) {
	if c == nil {
		return compute()
	}
	if raw, err := c.rdb.Get(ctx, key).Bytes(); err == nil {
		var boxed []byte
		_ = json.Unmarshal(raw, &boxed)
		if boxed != nil {
			var out any
			if err := json.Unmarshal(boxed, &out); err == nil {
				return out, nil
			}
		}
	}
	val, err := compute()
	if err != nil {
		return nil, err
	}
	if b, err := json.Marshal(val); err == nil {
		_ = c.rdb.Set(ctx, key, b, c.ttl).Err()
	}
	return val, nil
}

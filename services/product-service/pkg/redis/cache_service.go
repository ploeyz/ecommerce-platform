package redis

import (
	"context"
	"fmt"
	"time"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService(client *redis.Client) *CacheService{
	return &CacheService{client: client}
}

// Set stores a value in cache with TTL
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (c *CacheService) Get(ctx context.Context,key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a key from cache
func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// DeletePattern deletes all keys matching a pattern
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
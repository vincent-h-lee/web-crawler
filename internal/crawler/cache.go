package crawler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	HasRecentlyCrawled(ctx context.Context, u string) (bool, error)
	SetUrl(ctx context.Context, u string) error
}

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache(cache *redis.Client) Cache {
	return &RedisCache{cache}
}

func (c *RedisCache) HasRecentlyCrawled(ctx context.Context, u string) (bool, error) {
	_, err := c.cache.Get(ctx, u).Result()

	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// the set value in cache doesn't matter
	return true, nil
}

func (c *RedisCache) SetUrl(ctx context.Context, u string) error {
	// TODO pass in duration
	duration := 24 * time.Hour
	c.cache.Set(ctx, u, time.Now().Unix(), duration)
	return nil
}

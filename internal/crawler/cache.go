package crawler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	cache *redis.Client
}

func NewCache(cache *redis.Client) *Cache {
	return &Cache{cache}
}

func (c *Cache) HasRecentlyCrawled(ctx context.Context, u string) bool {
	hit := c.cache.Get(ctx, u)

	return hit != nil
}

func (c *Cache) SetUrl(ctx context.Context, u string) error {
	duration, err := time.ParseDuration("7d")
	if err != nil {
		return err
	}
	c.cache.Set(ctx, u, time.Now().Unix(), duration)
	return nil
}

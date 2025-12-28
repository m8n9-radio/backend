package cache

import (
	"context"
	"fmt"
	"hub/internal/config"
	"hub/internal/logger"
	"log"

	"github.com/go-redis/redis/v8"
)

type (
	Cache interface {
		Get(ctx context.Context, key CacheKey) (string, error)
		Set(ctx context.Context, key CacheKey, value string) error
	}
	cache struct {
		prefix string
		rdb    *redis.Client
	}
)

type CacheKey string

const (
	CacheKeyStartStream CacheKey = "start_stream"
)

func NewCache(cfg config.Config, logger logger.Logger) Cache {
	ctx := context.Background()
	url, prefix := cfg.RedisConnection()

	options, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	rdb := redis.NewClient(options)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &cache{
		rdb:    rdb,
		prefix: prefix,
	}
}

func (c *cache) Get(ctx context.Context, key CacheKey) (string, error) {
	return c.rdb.Get(ctx, fmt.Sprintf("%s%s", c.prefix, key)).Result()
}

func (c *cache) Set(ctx context.Context, key CacheKey, value string) error {
	return c.rdb.Set(ctx, fmt.Sprintf("%s%s", c.prefix, key), value, 0).Err()
}

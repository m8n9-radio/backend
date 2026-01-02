package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hub/internal/config"
	"hub/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss        = errors.New("cache miss")
	ErrCacheUnavailable = errors.New("cache unavailable")
)

type (
	Cache interface {
		Get(ctx context.Context, key string, dest interface{}) error
		Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
		Delete(ctx context.Context, key string) error
		Exists(ctx context.Context, key string) (bool, error)
		Client() *redis.Client
	}

	cache struct {
		client *redis.Client
		prefix string
		logger *logger.Logger
	}
)

func NewCache(cfg config.Config, log *logger.Logger) (Cache, error) {
	url, prefix := cfg.RedisConnection()

	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis URL: %w", err)
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Warnf("Redis connection failed: %v", err)
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Info("Redis cache connected")

	return &cache{
		client: client,
		prefix: prefix,
		logger: log,
	}, nil
}

func (c *cache) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := c.prefixKey(key)

	data, err := c.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := c.prefixKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, fullKey, data, ttl).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	fullKey := c.prefixKey(key)
	return c.client.Del(ctx, fullKey).Err()
}

func (c *cache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := c.prefixKey(key)
	count, err := c.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *cache) Client() *redis.Client {
	return c.client
}

func (c *cache) prefixKey(key string) string {
	return fmt.Sprintf("%s%s", c.prefix, key)
}

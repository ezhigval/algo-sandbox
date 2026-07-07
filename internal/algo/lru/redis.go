package lru

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache approximates LRU using sorted set scores as access time.
type RedisCache struct {
	client *redis.Client
	prefix string
	cap    int
	ttl    time.Duration
}

func NewRedis(client *redis.Client, capacity int, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: "lru:",
		cap:    capacity,
		ttl:    ttl,
	}
}

func (c *RedisCache) keySet() string { return c.prefix + "keys" }

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.HGet(ctx, c.keySet(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	score := float64(time.Now().UnixNano())
	_ = c.client.ZAdd(ctx, c.prefix+"order", redis.Z{Score: score, Member: key})
	return val, nil
}

func (c *RedisCache) Put(ctx context.Context, key, value string) error {
	pipe := c.client.Pipeline()
	pipe.HSet(ctx, c.keySet(), key, value)
	pipe.ZAdd(ctx, c.prefix+"order", redis.Z{Score: float64(time.Now().UnixNano()), Member: key})
	pipe.Expire(ctx, c.keySet(), c.ttl)
	pipe.Expire(ctx, c.prefix+"order", c.ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	count, err := c.client.ZCard(ctx, c.prefix+"order").Result()
	if err != nil {
		return err
	}

	for count > int64(c.cap) {
		removed, err := c.client.ZPopMin(ctx, c.prefix+"order").Result()
		if err != nil || len(removed) == 0 {
			break
		}
		evictKey := fmt.Sprint(removed[0].Member)
		_ = c.client.HDel(ctx, c.keySet(), evictKey)
		count--
	}

	return nil
}

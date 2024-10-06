package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{client}
}

func (rc *RedisCache) Set(ctx context.Context, key string, value any) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(value); err != nil {
		return err
	}

	cmd := rc.client.Set(ctx, key, buf.Bytes(), time.Hour)

	return cmd.Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string, res any) error {
	b, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("failed to fetch session data from Redis, err=%q", err)
	}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(res); err != nil {
		return fmt.Errorf("failed to decode go binary, err=%q", err)
	}
	return nil
}

func (rc *RedisCache) Exists(ctx context.Context, key string) bool {
	cmd := rc.client.Exists(ctx, key)
	res, err := cmd.Result()
	if err != nil {
		log.Printf("!!! WARNING: an error occured while checking the exists for key %q, err=%q", key, err)
	}

	return res != 0 && err == nil
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	cmd := rc.client.Del(ctx, key)

	return cmd.Err()
}

package auth

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisDBSessionStorage = 1
)

type RedisSessionStorage struct {
	client *redis.Client
}

func NewRedisSessionStorage(client *redis.Client) SessionStorage {
	return &RedisSessionStorage{client}
}

func (rss *RedisSessionStorage) StoreSession(ctx context.Context, key string, s Session) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(s); err != nil {
		return err
	}

	ttl := s.ExpiresAt.Sub(time.Now().UTC())
	cmd := rss.client.Set(ctx, key, buf.Bytes(), ttl)

	return cmd.Err()
}

func (rss *RedisSessionStorage) GetSession(ctx context.Context, key string) (Session, error) {

	sessBytes, err := rss.client.Get(ctx, key).Bytes()
	if err != nil {
		return Session{}, fmt.Errorf("failed to fetch session data from Redis, err=%q", err)
	}

	buf := bytes.NewBuffer(sessBytes)
	dec := gob.NewDecoder(buf)

	var session Session
	if err := dec.Decode(&session); err != nil {
		return session, fmt.Errorf("failed to decode go binary, err=%q", err)
	}
	return session, nil
}

package auth

import (
	"context"
	"time"
)

type Session struct {
	Id        string
	Secret    string
	ExpiresAt *time.Time
}

func (s Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now().UTC())
}

type SessionStorage interface {
	StoreSession(ctx context.Context, key string, s Session) error
	GetSession(ctx context.Context, key string) (Session, error)
}

package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, res any) error
	Exists(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
}

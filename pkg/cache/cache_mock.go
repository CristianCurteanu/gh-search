package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CacheMock struct {
	mock.Mock
}

func (cm *CacheMock) Set(ctx context.Context, key string, value any) error {
	args := cm.Called(ctx, key, value)

	return args.Error(0)
}

func (cm *CacheMock) Get(ctx context.Context, key string, res any) error {
	args := cm.Called(ctx, key, res)

	return args.Error(0)
}

func (cm *CacheMock) Exists(ctx context.Context, key string) bool {
	args := cm.Called(ctx, key)

	return args.Bool(0)
}

func (cm *CacheMock) Delete(ctx context.Context, key string) error {
	args := cm.Called(ctx, key)

	return args.Error(0)
}

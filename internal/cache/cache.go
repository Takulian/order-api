package cache

import (
	"context"
	"time"
)

type OrderCache interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

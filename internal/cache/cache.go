package cache

import (
	"context"
	"order-api/internal/model"
	"time"
)

type OrderCache interface {
	GetAll(ctx context.Context, key string) ([]model.Order, error)
	SetAll(ctx context.Context, key string, orders []model.Order, ttl time.Duration) error

	GetByID(ctx context.Context, id int, key string) (model.Order, error)
	SetByID(ctx context.Context, key string, order model.Order, ttl time.Duration) error

	Del(ctx context.Context, key string) error
}

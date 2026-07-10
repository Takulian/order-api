package cache

import (
	"context"
	"encoding/json"
	"order-api/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) GetAll(ctx context.Context, key string) ([]model.Order, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var orders []model.Order

	if err := json.Unmarshal([]byte(data), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *RedisCache) SetAll(ctx context.Context, key string, orders []model.Order, ttl time.Duration) error {
	data, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) GetByID(ctx context.Context, id int, key string) (model.Order, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return model.Order{}, err
	}
	var order model.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *RedisCache) SetByID(ctx context.Context, key string, order model.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

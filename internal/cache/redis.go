package cache

import (
	"context"
	"encoding/json"
	"errors"
	"order-api/internal/breaker"
	"order-api/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	cb     *breaker.CircuitBreaker
}

func NewRedisCache(client *redis.Client) *RedisCache {
	cb := breaker.NewCircuitBreaker(breaker.Settings{
		Name:                "redis-cache",
		FailureThreshold:    3,
		Openduration:        10 * time.Second,
		HalfOpenMaxRequests: 1,
	})

	return &RedisCache{
		client: client,
		cb:     cb,
	}
}

func (r *RedisCache) GetAll(ctx context.Context, key string) ([]model.Order, error) {
	var orders []model.Order
	err := r.cb.Execute(func() error {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(data), &orders)
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *RedisCache) SetAll(ctx context.Context, key string, orders []model.Order, ttl time.Duration) error {
	data, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	return r.cb.Execute(func() error {
		return r.client.Set(ctx, key, data, ttl).Err()
	})
}

func (r *RedisCache) GetByID(ctx context.Context, id int, key string) (model.Order, error) {
	var order model.Order
	var notfound bool
	err := r.cb.Execute(func() error {
		data, err := r.client.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			notfound = true
			return nil
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(data), &order)
	})
	if notfound == true {
		return model.Order{}, redis.Nil
	}

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *RedisCache) SetByID(ctx context.Context, key string, order model.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.cb.Execute(func() error {
		return r.client.Set(ctx, key, data, ttl).Err()
	})
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.cb.Execute(func() error {
		return r.client.Del(ctx, key).Err()
	})
}

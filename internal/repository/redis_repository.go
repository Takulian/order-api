package repository

import (
	"context"
	"errors"
	"fmt"
	"order-api/internal/model"
	"sort"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func orderKey(id int) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepository) GetAll() ([]model.Order, error) {

	keys, err := r.client.Keys(
		r.ctx,
		"order:*",
	).Result()

	if err != nil {
		return nil, err
	}

	orders := make([]model.Order, 0, len(keys))

	for _, key := range keys {

		data, err := r.client.HGetAll(
			r.ctx,
			key,
		).Result()

		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(data["id"])
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(data["quantity"])
		if err != nil {
			return nil, err
		}

		order := model.Order{
			ID:       id,
			Customer: data["customer"],
			Product:  data["product"],
			Quantity: quantity,
			Status:   data["status"],
		}

		orders = append(orders, order)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	return orders, nil
}

func (r *RedisRepository) GetByID(id int) (model.Order, error) {
	data, err := r.client.HGetAll(
		r.ctx,
		orderKey(id),
	).Result()

	if err != nil {
		return model.Order{}, err
	}

	if len(data) == 0 {
		return model.Order{}, errors.New("order not found")
	}

	order := model.Order{
		ID:       id,
		Customer: data["customer"],
		Product:  data["product"],
		Status:   data["status"],
	}

	order.Quantity, _ = strconv.Atoi(
		data["quantitiy"],
	)

	return order, nil
}

func (r *RedisRepository) Create(order model.Order) (model.Order, error) {

	id, err := r.client.Incr(
		r.ctx,
		"next_order_id",
	).Result()

	if err != nil {
		return model.Order{}, err
	}

	order.ID = int(id)

	err = r.client.HSet(
		r.ctx,
		orderKey(order.ID),

		"id", order.ID,
		"customer", order.Customer,
		"product", order.Product,
		"quantity", order.Quantity,
		"status", order.Status,
	).Err()

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *RedisRepository) Update(id int, order model.Order) (model.Order, error) {

	exists, err := r.client.Exists(
		r.ctx,
		orderKey(id),
	).Result()

	if err != nil {
		return model.Order{}, err
	}

	if exists == 0 {
		return model.Order{}, errors.New("order not found")
	}

	err = r.client.HSet(
		r.ctx,
		orderKey(order.ID),

		"id", order.ID,
		"customer", order.Customer,
		"product", order.Product,
		"quantity", order.Quantity,
		"status", order.Status,
	).Err()

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *RedisRepository) Delete(id int) error {

	deleted, err := r.client.Del(
		r.ctx,
		orderKey(id),
	).Result()

	if err != nil {
		return err
	}

	if deleted == 0 {
		return errors.New("order not found")
	}

	return nil
}

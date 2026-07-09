package service

import (
	"order-api/internal/dto"
	"order-api/internal/model"
	"order-api/internal/repository"

	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	repository repository.OrderRepository
	cache      *redis.Client
}

func NewOrderService(
	repository repository.OrderRepository,
	cache *redis.Client,
) *OrderService {
	return &OrderService{
		repository: repository,
		cache:      cache,
	}
}

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	const cacheKey = "orders"
	cachedOrders, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var orders []model.Order
		err := json.Unmarshal([]byte(cachedOrders), &orders)
		if err == nil {
			return orders, nil
		}
	}

	orders, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(orders)
	s.cache.Set(ctx, cacheKey, bytes, 0)

	return orders, nil
}

func (s *OrderService) GetByID(id int) (model.Order, error) {
	return s.repository.GetByID(id)
}

func (s *OrderService) Create(ctx context.Context, req dto.CreateOrderRequest) (model.Order, error) {
	if req.Customer == "" {
		return model.Order{}, ErrCustomerRequired
	}
	if req.Product == "" {
		return model.Order{}, ErrProductRequired
	}
	if req.Quantity <= 0 {
		return model.Order{}, ErrInvalidQuantity
	}

	order, err := s.repository.Create(model.Order{
		Customer: req.Customer,
		Product:  req.Product,
		Quantity: req.Quantity,
		Status:   "Pending",
	})
	if err != nil {
		return model.Order{}, err
	}

	s.cache.Del(ctx, "orders")

	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id int, req dto.UpdateOrderRequest) (model.Order, error) {
	if req.Customer == "" {
		return model.Order{}, ErrCustomerRequired
	}
	if req.Product == "" {
		return model.Order{}, ErrProductRequired
	}
	if req.Quantity <= 0 {
		return model.Order{}, ErrInvalidQuantity
	}
	order, err := s.GetByID(id)
	if err != nil {
		return model.Order{}, err
	}

	order.Customer = req.Customer
	order.Product = req.Product
	order.Quantity = req.Quantity

	updatedOrder, err := s.repository.Update(id, order)
	if err != nil {
		return model.Order{}, err
	}
	s.cache.Del(ctx, "orders")
	return updatedOrder, nil
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	s.cache.Del(ctx, "orders")
	return nil
}

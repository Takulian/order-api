package service

import (
	"fmt"
	"order-api/internal/cache"
	"order-api/internal/dto"
	"order-api/internal/model"
	"order-api/internal/repository"
	"time"

	"context"
)

type OrderService struct {
	repository repository.OrderRepository
	cache      cache.OrderCache
}

func NewOrderService(
	repository repository.OrderRepository,
	cache cache.OrderCache,
) *OrderService {
	return &OrderService{
		repository: repository,
		cache:      cache,
	}
}

const cacheKey = "orders"

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order

	if err := s.cache.Get(ctx, cacheKey, &orders); err == nil {
		return orders, nil
	}

	orders, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, orders, 5*time.Minute)

	return orders, nil
}

func (s *OrderService) GetByID(ctx context.Context, id int) (model.Order, error) {
	cacheKey := fmt.Sprintf("order:%d", id)
	var order model.Order
	if err := s.cache.Get(ctx, cacheKey, &order); err == nil {
		return order, nil
	}
	order, err := s.repository.GetByID(id)
	if err != nil {
		return model.Order{}, err
	}

	_ = s.cache.Set(ctx, cacheKey, order, 5*time.Minute)

	return order, nil
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

	_ = s.cache.Del(ctx, cacheKey)

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
	order, err := s.GetByID(ctx, id)
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
	_ = s.cache.Del(ctx, cacheKey)
	_ = s.cache.Del(ctx, fmt.Sprintf("order:%d", id))

	return updatedOrder, nil
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	_ = s.cache.Del(ctx, cacheKey)
	return nil
}

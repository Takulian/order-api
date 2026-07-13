package service

import (
	"fmt"
	"log/slog"
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
	logger     *slog.Logger
}

func NewOrderService(
	repository repository.OrderRepository,
	cache cache.OrderCache,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		repository: repository,
		cache:      cache,
		logger:     logger,
	}
}

const cacheKey = "orders"

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order

	if orders, err := s.cache.GetAll(ctx, cacheKey); err == nil {
		return orders, nil
	}

	orders, err := s.repository.GetAll()
	if err != nil {
		s.logger.ErrorContext(ctx, "gagal ambil data", "error", err)
		return nil, err
	}

	if err := s.cache.SetAll(ctx, cacheKey, orders, 5*time.Minute); err != nil {
		s.logger.ErrorContext(ctx, "gagal simpan cache", "error", err)
	}

	return orders, nil
}

func (s *OrderService) GetByID(ctx context.Context, id int) (model.Order, error) {
	idKey := fmt.Sprintf("orders:%d", id)
	if order, err := s.cache.GetByID(ctx, id, idKey); err == nil {
		return order, nil
	}
	order, err := s.repository.GetByID(id)
	if err != nil {
		s.logger.ErrorContext(ctx, "gagal ambil data", "error", err)
		return model.Order{}, err
	}
	if err := s.cache.SetByID(ctx, idKey, order, 5*time.Minute); err != nil {
		s.logger.ErrorContext(ctx, "gagal simpan cache", "error", err)
	}

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
		s.logger.ErrorContext(ctx, "gagal membuat order", "error", err)
		return model.Order{}, err
	}

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus cache", "error", err)
	}

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
		s.logger.ErrorContext(ctx, "gagal ambil data", "error", err)
		return model.Order{}, err
	}

	order.Customer = req.Customer
	order.Product = req.Product
	order.Quantity = req.Quantity

	updatedOrder, err := s.repository.Update(id, order)
	if err != nil {
		s.logger.ErrorContext(ctx, "gagal update data", "error", err)
		return model.Order{}, err
	}
	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus cache", "error", err)
	}
	if err := s.cache.Del(ctx, fmt.Sprintf("orders:%d", id)); err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus cache", "error", err)
	}

	return updatedOrder, nil
}

func (s *OrderService) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus data", "error", err)
		return err
	}
	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus cache", "error", err)
	}
	if err := s.cache.Del(ctx, fmt.Sprintf("orders:%d", id)); err != nil {
		s.logger.ErrorContext(ctx, "gagal hapus cache", "error", err)
	}

	return nil
}

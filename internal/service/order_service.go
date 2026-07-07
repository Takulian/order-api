package service

import (
	"order-api/internal/dto"
	"order-api/internal/model"
	"order-api/internal/repository"
)

type OrderService struct{}

var OrderSvc = &OrderService{}

func (s *OrderService) GetAll() []model.Order {
	return repository.OrderRepo.GetAll()
}

func (s *OrderService) Create(req dto.CreateOrderRequest) model.Order {
	order := model.Order{
		ID:       model.NextID,
		Customer: req.Customer,
		Product:  req.Product,
		Quantity: req.Quantity,
		Status:   "Pending",
	}

	model.NextID++

	return repository.OrderRepo.Create(order)
}

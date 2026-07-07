package service

import (
	"order-api/internal/dto"
	"order-api/internal/model"
	"order-api/internal/repository"
)

type OrderService struct{
	repository *repository.OrderRepository
}

func NewOrderService(repository *repository.OrderRepository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s *OrderService) GetAll() []model.Order {
	return s.repository.GetAll()
}

func (s *OrderService) GetByID(id int) (model.Order, error) {
	return s.repository.GetByID(id)
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

	return s.repository.Create(order)
}

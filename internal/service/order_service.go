package service

import (
	"order-api/internal/dto"
	"order-api/internal/model"
)

type OrderRepository interface {
	GetAll() []model.Order
	GetByID(id int) (model.Order, error)
	Create(order model.Order) model.Order
	Update(id int, order model.Order) (model.Order, error)
	Delete(id int) error
}

type OrderService struct{
	repository OrderRepository
}

func NewOrderService(repository OrderRepository) *OrderService {
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
		Customer: req.Customer,
		Product:  req.Product,
		Quantity: req.Quantity,
		Status:   "Pending",
	}

	return s.repository.Create(order)
}

func (s *OrderService) Update(id int, req dto.UpdateOrderRequest) (model.Order, error) {
	order, err := s.repository.GetByID(id)
	if err != nil {
		return model.Order{}, err
	}

	order.Customer = req.Customer
	order.Product = req.Product
	order.Quantity = req.Quantity

	return s.repository.Update(id, order)
}

func (s *OrderService) Delete(id int) error {
	return s.repository.Delete(id)
}


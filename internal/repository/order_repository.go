package repository

import (
	"errors"
	"order-api/internal/model"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) GetAll() []model.Order {
	return model.Orders
}

func (r *OrderRepository) GetByID(id int) (model.Order, error) {
	for _, order := range model.Orders {
		if order.ID == id {
			return order, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *OrderRepository) Create(order model.Order) model.Order {

	model.Orders = append(model.Orders, order)

	return order
}

func (r *OrderRepository) Update(id int, updatedOrder model.Order) (model.Order, error) {
	for i, order := range model.Orders {
		if order.ID == id {
			model.Orders[i] = updatedOrder
			return updatedOrder, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *OrderRepository) Delete(id int) error {
	for i, order := range model.Orders {
		if order.ID == id {
			model.Orders = append(model.Orders[:i], model.Orders[i+1:]...)
			return nil
		}
	}
	return errors.New("order not found")
}
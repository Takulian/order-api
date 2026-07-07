package repository

import (
	"order-api/internal/model"
)

type OrderRepository struct{}

func (r *OrderRepository) GetAll() []model.Order {
	return model.Orders
}

func (r *OrderRepository) Create(order model.Order) model.Order {

	model.Orders = append(model.Orders, order)

	return order
}

var OrderRepo = &OrderRepository{}
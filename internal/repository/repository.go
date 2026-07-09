package repository

import (
	"order-api/internal/model"
)

type OrderRepository interface {
	GetAll() ([]model.Order, error)
	GetByID(id int) (model.Order, error)
	Create(order model.Order) (model.Order, error)
	Update(id int, updatedOrder model.Order) (model.Order, error)
	Delete(id int) error
}

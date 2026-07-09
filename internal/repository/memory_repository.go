package repository

import (
	"errors"
	"order-api/internal/model"
)

type MemoryRepository struct {
	orders []model.Order
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: []model.Order{
			{
				ID:       1,
				Customer: "Andi",
				Product:  "Laptop",
				Quantity: 1,
				Status:   "Pending",
			},
			{
				ID:       2,
				Customer: "Budi",
				Product:  "Mouse",
				Quantity: 2,
				Status:   "Paid",
			},
		},
		nextID: 3,
	}
}

func (r *MemoryRepository) GetAll() ([]model.Order, error) {
	return r.orders, nil
}

func (r *MemoryRepository) GetByID(id int) (model.Order, error) {
	for _, order := range r.orders {
		if order.ID == id {
			return order, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *MemoryRepository) Create(order model.Order) (model.Order, error) {
	order.ID = r.nextID
	r.nextID++
	r.orders = append(r.orders, order)

	return order, nil
}

func (r *MemoryRepository) Update(id int, updatedOrder model.Order) (model.Order, error) {
	for i, order := range r.orders {
		if order.ID == id {
			r.orders[i] = updatedOrder
			return updatedOrder, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *MemoryRepository) Delete(id int) error {
	for i, order := range r.orders {
		if order.ID == id {
			r.orders = append(r.orders[:i], r.orders[i+1:]...)
			return nil
		}
	}
	return errors.New("order not found")
}

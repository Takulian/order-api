package repository

import (
	"errors"
	"order-api/internal/model"
)

type OrderRepository struct{
	orders []model.Order
	nextID int
}


func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
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

func (r *OrderRepository) GetAll() []model.Order {
	return r.orders
}

func (r *OrderRepository) GetByID(id int) (model.Order, error) {
	for _, order := range r.orders {
		if order.ID == id {
			return order, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *OrderRepository) Create(order model.Order) model.Order {
	order.ID = r.nextID
	r.nextID++
	r.orders = append(r.orders, order)

	return order
}

func (r *OrderRepository) Update(id int, updatedOrder model.Order) (model.Order, error) {
	for i, order := range r.orders {
		if order.ID == id {
			r.orders[i] = updatedOrder
			return updatedOrder, nil
		}
	}
	return model.Order{}, errors.New("order not found")
}

func (r *OrderRepository) Delete(id int) error {
	for i, order := range r.orders {
		if order.ID == id {
			r.orders = append(r.orders[:i], r.orders[i+1:]...)
			return nil
		}
	}
	return errors.New("order not found")
}
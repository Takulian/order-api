package repository

import (
	"database/sql"
	"errors"
	"order-api/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetAll() ([]model.Order, error) {
	query := `SELECT id, customer, product, quantity, status FROM orders ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.Customer,
			&order.Product,
			&order.Quantity,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *PostgresRepository) GetByID(id int) (model.Order, error) {
	query := `SELECT id, customer, product, quantity, status FROM orders WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var order model.Order
	err := row.Scan(
		&order.ID,
		&order.Customer,
		&order.Product,
		&order.Quantity,
		&order.Status,
	)
	if err != nil {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *PostgresRepository) Create(order model.Order) (model.Order, error) {
	query := `INSERT INTO orders (customer, product, quantity, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, order.Customer, order.Product, order.Quantity, order.Status).Scan(&order.ID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (r *PostgresRepository) Update(id int, order model.Order) (model.Order, error) {
	query := `UPDATE orders SET customer = $1, product = $2, quantity = $3, status = $4 WHERE id = $5`
	result, err := r.db.Exec(query, order.Customer, order.Product, order.Quantity, order.Status, order.ID)
	if err != nil {
		return model.Order{}, errors.New("failed to update order")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Order{}, errors.New("failed to retrieve rows affected")
	}
	if rowsAffected == 0 {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *PostgresRepository) Delete(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

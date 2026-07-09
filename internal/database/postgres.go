package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB() (*sql.DB, error) {
	conn := "postgres://myuser:mykisah@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

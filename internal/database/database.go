package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser() error
}

type PostgresqlStorage struct {
	db *sql.DB
}

func (storage *PostgresqlStorage) CreateUser() error {
	return nil
}

func NewPostgresStorage() (*PostgresqlStorage, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresqlStorage{
		db: db,
	}, nil
}

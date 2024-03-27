package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type contextKey string

const DBContextKey contextKey = "db"

func (c contextKey) String() string {
	return string(c)
}

type PostgresqlStorage struct {
	DB *sql.DB
}

func NewPostgresStorage() (*PostgresqlStorage, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresqlStorage{
		DB: db,
	}, nil
}

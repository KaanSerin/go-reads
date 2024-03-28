package database

import (
	"database/sql"
	"fmt"
	"os"

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
	dbUrl := os.Getenv("DB_URL")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")

	var connStr string
	if dbUrl != "" {
		connStr = dbUrl
		fmt.Println("Trying DB connection using ", dbUrl)
	} else {
		connStr = fmt.Sprintf("host=%s port=%s sslmode=disable dbname=%s user=%s password=%s", dbHost, dbPort, dbName, dbUser, dbPass)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresqlStorage{
		DB: db,
	}, nil
}

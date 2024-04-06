package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type contextKey string

const DBContextKey contextKey = "db"

func (c contextKey) String() string {
	return string(c)
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresqlStorage struct {
	db *sql.DB
}

type Storage interface {
	GetUsers() ([]*User, error)
	GetUserById(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	CreateUser(string, string, string, string) (*User, error)
	DeleteUserById(int) error
}

func (storage *PostgresqlStorage) GetUserById(id int) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, created_at from users where id = $1", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *PostgresqlStorage) GetUsers() ([]*User, error) {
	var users []*User = []*User{}

	query := "select id, first_name, last_name, email, created_at from users"
	rows, err := storage.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (storage *PostgresqlStorage) GetUserByEmail(email string) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, created_at from users where email = $1", email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *PostgresqlStorage) CreateUser(first_name, last_name, email, password string) (*User, error) {
	_, err := storage.db.Exec(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
		first_name,
		last_name,
		email,
		password)

	if err != nil {
		return nil, err
	}

	return storage.GetUserByEmail(email)
}

func (storage *PostgresqlStorage) DeleteUserById(id int) error {
	_, err := storage.db.Exec(fmt.Sprintf("DELETE FROM users where id = %d", id))
	return err
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
		db: db,
	}, nil
}

func GetPgStorageFromRequest(r *http.Request) (*PostgresqlStorage, error) {
	db := r.Context().Value(DBContextKey).(*PostgresqlStorage)
	return db, nil
}

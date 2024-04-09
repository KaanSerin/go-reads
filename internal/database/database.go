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
	RoleId    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresqlStorage struct {
	db *sql.DB
}

type Storage interface {
	GetUsers() ([]*User, error)
	GetUserById(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	CreateUser(firstName, lastName, email, password string) (*User, error)
	DeleteUserById(int) error
	GetRoleById(int) (*Role, error)
}

func (storage *PostgresqlStorage) GetUserById(id int) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, role_id, created_at from users where id = $1", id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.RoleId,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (storage *PostgresqlStorage) GetUsers() ([]*User, error) {
	var users []*User = []*User{}

	query := "select id, first_name, last_name, email, role_id, created_at from users"
	rows, err := storage.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.RoleId, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (storage *PostgresqlStorage) GetUserByEmail(email string) (*User, error) {
	var user *User = &User{}

	err := storage.db.QueryRow(
		"SELECT id, first_name, last_name, email, password, created_at from users where email = $1", email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
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

	user, err := storage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (storage *PostgresqlStorage) DeleteUserById(id int) error {
	_, err := storage.db.Exec(fmt.Sprintf("DELETE FROM users where id = %d", id))
	return err
}

type UpdateUserDto struct {
	ID        int    `json:"id" validate:"nonzero"`
	FirstName string `json:"first_name" validate:"nonzero"`
	LastName  string `json:"last_name" validate:"nonzero"`
	Email     string `json:"email" validate:"nonzero"`
}

func (storage *PostgresqlStorage) UpdateUserById(id int, payload *UpdateUserDto) (*User, error) {
	user, err := storage.GetUserById(id)
	if err != nil {
		return nil, err
	}

	_, err = storage.db.Exec("UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4", payload.FirstName, payload.LastName, payload.Email, id)
	if err != nil {
		return nil, err
	}

	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Email = payload.Email
	return user, nil
}

func (storage *PostgresqlStorage) GetRoleById(id int) (*Role, error) {
	var role *Role = &Role{}
	err := storage.db.QueryRow("SELECT id, name, created_at FROM roles WHERE id = $1", id).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func NewPostgresStorage() (*PostgresqlStorage, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
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

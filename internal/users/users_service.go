package users

import (
	"net/http"
	"time"

	"github.com/kaanserin/go-reads/internal/database"
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUsers(r *http.Request) ([]*User, error) {
	storage := r.Context().Value(database.DBContextKey).(*database.PostgresqlStorage)

	users := []*User{}

	query := "select id, first_name, last_name, email, created_at from users"
	rows, err := storage.DB.Query(query)
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

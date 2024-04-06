package users

import (
	"github.com/kaanserin/go-reads/internal/database"
)

func GetUsers(storage database.Storage) ([]*database.User, error) {
	return storage.GetUsers()
}

func GetUserById(storage database.Storage, id int) (*database.User, error) {
	return storage.GetUserById(id)
}

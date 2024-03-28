package auth

import (
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/utils"
)

func SignUp(createUserDto CreateUserDto, storage database.Storage) (*database.User, error) {
	sameUser, err := storage.GetUserByEmail(createUserDto.Email)
	if sameUser != nil && err != nil {
		return nil, err
	}

	if sameUser != nil {
		return nil, &utils.CustomError{
			Message: "User with same email already exists",
		}
	}

	return storage.CreateUser(createUserDto.FirstName, createUserDto.LastName,
		createUserDto.Email, createUserDto.Password)
}

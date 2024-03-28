package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/utils"
	"gopkg.in/validator.v2"
)

func GetAuthRouter(r chi.Router) {
	r.Post("/sign_up", utils.MakeHandlerFunc(signUpHandler))
}

type CreateUserDto struct {
	FirstName string `validate:"nonzero"`
	LastName  string `validate:"nonzero"`
	Email     string `validate:"nonzero"`
	Password  string `validate:"nonzero"`
}

func signUpHandler(w http.ResponseWriter, r *http.Request) error {
	var createUserDto CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&createUserDto)
	if err != nil {
		return err
	}

	if errs := validator.Validate(createUserDto); errs != nil {
		return errs
	}

	db, err := database.GetPgStorageFromRequest(r)
	if err != nil {
		return err
	}

	user, err := SignUp(createUserDto, db)
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, 200, user)
}

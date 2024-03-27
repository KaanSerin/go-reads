package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	utils "github.com/kaanserin/go-reads/internal/utils"
)

var makeHandlerFunc = utils.MakeHandlerFunc

// Router
func GetUsersRouter(r chi.Router) {
	r.Get("/", makeHandlerFunc(getUsers))
}

// Handler Functions
func getUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := GetUsers(r)
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, 200, users)
}

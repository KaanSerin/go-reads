package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaanserin/go-reads/internal/database"
	utils "github.com/kaanserin/go-reads/internal/utils"
)

var makeHandlerFunc = utils.MakeHandlerFunc

// Router
func GetUsersRouter(r chi.Router) {
	r.Get("/", makeHandlerFunc(getUsers))
}

// Handler Functions
func getUsers(w http.ResponseWriter, r *http.Request) error {
	db, err := database.GetPgStorageFromRequest(r)
	if err != nil {
		return err
	}

	users, err := GetUsers(db)
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, http.StatusOK, users)
}

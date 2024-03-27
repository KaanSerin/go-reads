package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	users "github.com/kaanserin/go-reads/internal/users"
)

func CreateNewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Register routes here
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Users routes
	r.Route("/users", users.GetUsersRouter)

	return r
}

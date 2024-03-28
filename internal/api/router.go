package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaanserin/go-reads/internal/auth"
	users "github.com/kaanserin/go-reads/internal/users"
)

func CreateNewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Register routes here
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// Users routes
	r.Route("/users", users.GetUsersRouter)

	r.Route("/auth", auth.GetAuthRouter)

	return r
}

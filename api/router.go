package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateNewRouter(server *ApiServer) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", server.CreateUser)
	})

	return r
}

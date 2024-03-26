package api

import (
	"fmt"
	"net/http"

	database "github.com/kaanserin/go-reads/internal/database"
)

type ApiServer struct {
	listenAddr string
	store      database.Storage
}

func NewApiServer(listenAddr string, store database.Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (server *ApiServer) ListenAndServe() error {
	router := CreateNewRouter(server)
	fmt.Printf("Server running on %s\n", server.listenAddr)
	err := http.ListenAndServe(server.listenAddr, router)
	return err
}

func (server *ApiServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := server.store.CreateUser(); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("User created successfully"))
}

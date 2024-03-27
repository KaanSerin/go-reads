package api

import (
	"context"
	"net"
	"net/http"

	"github.com/kaanserin/go-reads/internal/database"
)

func NewServer(listenAddr string) (*http.Server, error) {
	router := CreateNewRouter()

	db, err := database.NewPostgresStorage()
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(context.Background(), database.DBContextKey, db)
		},
	}

	return server, nil
}

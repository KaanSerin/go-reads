package main

import (
	"log"

	api "github.com/kaanserin/go-reads/api"
	database "github.com/kaanserin/go-reads/internal/database"
)

func main() {
	port := ":8080"

	storage, err := database.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(port, storage)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

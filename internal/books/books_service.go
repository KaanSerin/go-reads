package books

import (
	"net/http"

	"github.com/kaanserin/go-reads/internal/database"
)

func GetBooks(r *http.Request) ([]*database.Book, error) {
	storage, err := database.GetPgStorageFromRequest(r)
	if err != nil {
		return nil, err
	}

	return storage.GetBooks()
}

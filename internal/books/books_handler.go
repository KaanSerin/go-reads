package books

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/utils"
)

func AddBooksRoutes(r *gin.Engine) {
	booksGroup := r.Group("books")

	booksGroup.GET("/", utils.MakeHandlerFunc(getBooksForSubject))
}

func getBooksForSubject(c *gin.Context) error {
	books, err := GetBooks(c.Request)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, books)
	return nil
}

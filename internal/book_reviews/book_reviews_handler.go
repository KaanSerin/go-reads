package bookreviews

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/middleware"
	"github.com/kaanserin/go-reads/internal/utils"
)

func AddBookReviewsRoutes(c *gin.Engine) {
	router := c.Group("/book_reviews")

	router.Use(middleware.Authentication())

	router.GET("/", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(getBookReviews))
}

func getBookReviews(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	bookReviews, err := storage.GetBookReviews(c.Request)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, bookReviews)
	return nil
}

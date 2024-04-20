package bookreviews

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/middleware"
	"github.com/kaanserin/go-reads/internal/utils"
)

func AddBookReviewsRoutes(c *gin.Engine) {
	router := c.Group("/book_reviews")

	router.Use(middleware.Authentication())

	router.GET("/", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(getBookReviews))
	router.GET("/:id", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(getBookReviewById))
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

func getBookReviewById(c *gin.Context) error {
	idParam, _ := c.Params.Get("id")
	if idParam == "" {
		return &utils.CustomError{
			Message: "No id param in given",
		}
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, utils.CustomError{
			Message: "Id is not a number",
		})

		return nil
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	bookReview, err := storage.GetBookReviewById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, bookReview)
	return nil
}

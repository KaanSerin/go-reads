package bookreviews

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/middleware"
	"github.com/kaanserin/go-reads/internal/utils"
	"gopkg.in/validator.v2"
)

func AddBookReviewsRoutes(c *gin.Engine) {
	router := c.Group("/book_reviews")

	router.Use(middleware.Authentication())

	router.GET("/", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(getBookReviews))
	router.POST("/", utils.MakeHandlerFunc(createBookReview))
	router.GET("/:id", utils.MakeHandlerFunc(getBookReviewById))
	router.DELETE("/:id", utils.MakeHandlerFunc(deleteBookReviewById))
	router.PUT("/:id", utils.MakeHandlerFunc(updateBookReview))
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

func createBookReview(c *gin.Context) error {
	var createBookReviewDto *database.CreateBookReviewDto = &database.CreateBookReviewDto{}
	json.NewDecoder(c.Request.Body).Decode(createBookReviewDto)

	if err := validator.Validate(createBookReviewDto); err != nil {
		return err
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	userVal, _ := c.Get("user")
	user := userVal.(*database.User)
	createBookReviewDto.UserID = user.ID

	bookReview, err := storage.CreateBookReview(createBookReviewDto)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, bookReview)
	return nil
}

func deleteBookReviewById(c *gin.Context) error {
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

	userTmp, _ := c.Get("user")
	user := userTmp.(*database.User)
	bookReview, err := storage.GetBookReviewById(id)
	if err != nil {
		return err
	}

	if bookReview.UserID != user.ID {
		return &utils.CustomError{
			Message: "Forbidden",
		}
	}

	if err := storage.DeleteBookReviewById(id); err != nil {
		return err
	}

	c.JSON(http.StatusOK, utils.CustomError{
		Message: "Book review deleted successfully",
	})

	return nil
}

func updateBookReview(c *gin.Context) error {
	var updateBookReviewDto *database.UpdateBookReviewDto = &database.UpdateBookReviewDto{}
	json.NewDecoder(c.Request.Body).Decode(updateBookReviewDto)

	err := validator.Validate(updateBookReviewDto)
	if err != nil {
		return err
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CustomError{
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

	userTmp, _ := c.Get("user")
	user := userTmp.(*database.User)
	if bookReview.UserID != user.ID {
		return &utils.CustomError{
			Message: "Forbidden",
		}
	}

	bookReview, err = storage.UpdateBookReview(id, *updateBookReviewDto)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, bookReview)
	return nil
}

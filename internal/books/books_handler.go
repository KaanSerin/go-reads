package books

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

func AddBooksRoutes(r *gin.Engine) {
	booksGroup := r.Group("books")

	booksGroup.Use(middleware.Authentication())

	booksGroup.GET("/", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(getBooks))
	booksGroup.GET("/:id", utils.MakeHandlerFunc(getBookById))
	booksGroup.PUT("/:id", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(updateBookById))
	booksGroup.DELETE("/:id", middleware.AuthorizeAdmin(), utils.MakeHandlerFunc(deleteBookById))
}

func getBooks(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	books, err := storage.GetBooks(c.Request)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, books)
	return nil
}

func getBookById(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}

	books, err := storage.GetBookById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, books)
	return nil
}

func updateBookById(c *gin.Context) error {
	var updateBookDto *database.UpdateBookDto = &database.UpdateBookDto{}
	json.NewDecoder(c.Request.Body).Decode(updateBookDto)

	errs := validator.Validate(updateBookDto)
	if errs != nil {
		return errs
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return &utils.CustomError{
			Message: "Please enter a valid integer for id",
		}
	}

	book, err := storage.UpdateBookById(id, updateBookDto)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, book)
	return nil
}

func deleteBookById(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return &utils.CustomError{
			Message: "Please enter a valid integer for id",
		}
	}

	err = storage.DeleteBookById(id)
	if err != nil {
		return err
	}

	c.JSON(200, &utils.CustomError{
		Message: "Book deleted successfully",
	})

	return nil
}

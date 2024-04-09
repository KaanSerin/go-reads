package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/middleware"
	utils "github.com/kaanserin/go-reads/internal/utils"
	"gopkg.in/validator.v2"
)

var makeHandlerFunc = utils.MakeHandlerFunc

// Router
func AddUserRoutes(g *gin.Engine) {
	users := g.Group("/users")

	users.Use(middleware.Authentication())
	users.Use(middleware.AuthorizeAdmin())

	users.GET("/", makeHandlerFunc(getUsers))
	users.GET("/:id", makeHandlerFunc(getUserById))
	users.PUT("/:id", makeHandlerFunc(updateUser))
	users.DELETE("/:id", makeHandlerFunc(deleteUserById))
}

// Handler Functions
func getUsers(c *gin.Context) error {
	db, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	users, err := GetUsers(db)
	if err != nil {
		return err
	}

	c.JSON(200, users)
	return nil
}

func getUserById(c *gin.Context) error {
	db, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := GetUserById(db, id)
	if err != nil {
		return err
	}

	c.JSON(200, user)
	return nil
}

func deleteUserById(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

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

	user, err := storage.GetUserById(id)
	if user == nil {
		c.JSON(404, utils.CustomError{
			Message: "User not found",
		})

		return nil
	}

	if err != nil {
		return err
	}

	err = storage.DeleteUserById(id)
	if err != nil {
		return err
	}

	c.JSON(200, utils.MessageResponse{
		Message: "User deleted successfully",
	})

	return nil
}

func updateUser(c *gin.Context) error {
	var updatePayload *database.UpdateUserDto
	err := json.NewDecoder(c.Request.Body).Decode(&updatePayload)
	if err != nil {
		return err
	}

	if err := validator.Validate(updatePayload); err != nil {
		return err
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	idParam, _ := c.Params.Get("id")
	if idParam == "" {
		return &utils.CustomError{
			Message: "No id param in given",
		}
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CustomError{
			Message: "Id is not a number",
		})

		return nil
	}

	user, err := storage.UpdateUserById(id, updatePayload)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}

package users

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/database"
	utils "github.com/kaanserin/go-reads/internal/utils"
)

var makeHandlerFunc = utils.MakeHandlerFunc

// Router
func AddUserRoutes(g *gin.Engine) {
	users := g.Group("/users")
	users.GET("/", makeHandlerFunc(getUsers))
	users.GET("/:id", makeHandlerFunc(getUserById))
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

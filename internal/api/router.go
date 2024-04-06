package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kaanserin/go-reads/internal/auth"
	"github.com/kaanserin/go-reads/internal/users"
)

func CreateNewRouter() *gin.Engine {
	r := gin.Default()

	// Register routes here
	users.AddUserRoutes(r)
	auth.AddAuthRoutes(r)

	return r
}

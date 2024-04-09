package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/utils"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("APP_KEY")), nil
		})

		fmt.Println(token.Valid)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		claimStrings := token.Claims.(jwt.MapClaims)
		jti := claimStrings["jti"].(string)
		id, err := strconv.Atoi(jti)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		storage, err := database.GetPgStorageFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		user, err := storage.GetUserById(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		c.Set("user", user)
		c.Next()
	}
}

func AuthorizeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userTmp, exists := c.Get("user")
		if !exists || userTmp == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		storage, err := database.GetPgStorageFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		var user = userTmp.(*database.User)
		role, err := storage.GetRoleById(user.RoleId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		if role.Name != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CustomError{
				Message: "Unauthorized",
			})
		}

		c.Next()
	}
}

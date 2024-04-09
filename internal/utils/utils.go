package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONResponse(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

type functionWithError = func(c *gin.Context) error

func MakeHandlerFunc(fn functionWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := fn(c)
		if err != nil {
			c.JSON(400, CustomError{Message: err.Error()})
		}
	}
}

type CustomError struct {
	Message string `json:"message"`
}

func (c *CustomError) Error() string {
	return c.Message
}

type MessageResponse struct {
	Message string `json:"message"`
}

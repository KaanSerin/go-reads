package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kaanserin/go-reads/internal/database"
	"github.com/kaanserin/go-reads/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

var makeHandlerFunc = utils.MakeHandlerFunc

type CreateUserDto struct {
	FirstName string `validate:"nonzero"`
	LastName  string `validate:"nonzero"`
	Email     string `validate:"nonzero"`
	Password  string `validate:"nonzero"`
}

type AuthUserResponse struct {
	User        *database.User `json:"user"`
	AccessToken string         `json:"accessToken"`
}

// Register Handlers
func AddAuthRoutes(c *gin.Engine) {
	router := c.Group("/auth")
	router.POST("/sign_up", makeHandlerFunc(signUpHandler))
}

// Handlers
func signUpHandler(c *gin.Context) error {
	var createUserDto CreateUserDto
	err := json.NewDecoder(c.Request.Body).Decode(&createUserDto)
	if err != nil {
		return err
	}

	if errs := validator.Validate(createUserDto); errs != nil {
		return errs
	}

	db, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(createUserDto.Password)
	if err != nil {
		return err
	}

	createUserDto.Password = hashedPassword

	user, err := SignUp(createUserDto, db)
	if err != nil {
		return err
	}

	claims := &jwt.RegisteredClaims{
		ID:        user.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	appKey := os.Getenv("APP_KEY")
	fmt.Printf("app key %s\n", appKey)

	accessToken, err := token.SignedString([]byte(appKey))
	if err != nil {
		c.JSON(200, nil)
	}

	fmt.Printf("Access token %s\n", accessToken)

	c.JSON(200, AuthUserResponse{
		User:        user,
		AccessToken: accessToken,
	})
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

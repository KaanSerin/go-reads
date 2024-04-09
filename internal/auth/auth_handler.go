package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	router.POST("/sign_in", makeHandlerFunc(signInHandler))
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

	accessToken, err := getAccessTokenStringForUser(user)
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

type SignInDto struct {
	Email    string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

func signInHandler(c *gin.Context) error {
	var signIn SignInDto
	if err := json.NewDecoder(c.Request.Body).Decode(&signIn); err != nil {
		return err
	}

	if err := validator.Validate(signIn); err != nil {
		return err
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	user, err := storage.GetUserByEmail(signIn.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, &utils.CustomError{
			Message: "Invalid email or password",
		})

		return nil
	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signIn.Password))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, &utils.CustomError{
			Message: "Invalid email or password",
		})
		return nil
	}

	accessToken, err := getAccessTokenStringForUser(user)
	if err != nil {
		return err
	}

	user.Password = ""

	c.JSON(http.StatusOK, AuthUserResponse{
		User:        user,
		AccessToken: accessToken,
	})
	return nil
}

func getAccessTokenStringForUser(user *database.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ID:        fmt.Sprint(user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	appKey := os.Getenv("APP_KEY")
	return token.SignedString([]byte(appKey))
}

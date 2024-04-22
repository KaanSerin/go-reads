package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	users.GET("/", middleware.AuthorizeAdmin(), makeHandlerFunc(getUsers))
	users.GET("/profile", makeHandlerFunc(getUserProfile))
	users.PUT("/profile", makeHandlerFunc(updateUserProfile))
	users.POST("/profile_image", makeHandlerFunc(updateUserProfileImage))
	users.GET("/:id", makeHandlerFunc(getUserById))
	users.PUT("/:id", middleware.AuthorizeAdmin(), makeHandlerFunc(updateUser))
	users.DELETE("/:id", middleware.AuthorizeAdmin(), makeHandlerFunc(deleteUserById))
}

// Handler Functions
func getUsers(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	users, err := storage.GetUsers(c.Request)
	if err != nil {
		return err
	}

	c.JSON(200, users)
	return nil
}

func getUserById(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := storage.GetUserById(id)
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

func getUserProfile(c *gin.Context) error {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
	return nil
}

func updateUserProfile(c *gin.Context) error {
	var updatePayload *database.UpdateUserDto
	err := json.NewDecoder(c.Request.Body).Decode(&updatePayload)
	if err != nil {
		return err
	}

	userTmp, _ := c.Get("user")
	user := userTmp.(*database.User)
	if user.ID != updatePayload.ID {
		return &utils.CustomError{
			Message: "Forbidden",
		}
	}

	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	user, err = storage.UpdateUserById(user.ID, updatePayload)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}

func updateUserProfileImage(c *gin.Context) error {
	storage, err := database.GetPgStorageFromRequest(c.Request)
	if err != nil {
		return err
	}

	userTmp, _ := c.Get("user")
	user := userTmp.(*database.User)

	imageFile, fileHeaders, err := c.Request.FormFile("image")
	if err != nil {
		return err
	}
	defer imageFile.Close()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	objectKey := fmt.Sprintf("profile/user/%d/profile_image%s", user.ID, filepath.Ext(fileHeaders.Filename))

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   imageFile,
	})

	if err != nil {
		return err
	}

	if err := storage.UpdateUserProfileImageUrl(user.ID, objectKey); err != nil {
		return err
	}

	user.ProfileImageUrl = objectKey
	c.JSON(http.StatusOK, user)

	return nil
}

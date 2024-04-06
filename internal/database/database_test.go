package database

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestPostgresStorage(t *testing.T) {
	godotenv.Load("../../.env")
	var storage Storage

	// Test with storage of choice
	storage, err := NewPostgresStorage()
	if err != nil {
		t.Fatal(err)
	}

	createUserPayload := User{
		FirstName: "TestFirstName",
		LastName:  "TestLastName",
		Email:     "testmail@mail.com",
		Password:  "pass",
	}

	t.Run("TestCreateUserWithValidInputsAndDelete", func(t *testing.T) {
		// Create a user
		user, err := storage.CreateUser(createUserPayload.FirstName, createUserPayload.LastName, createUserPayload.Email, createUserPayload.Password)
		if err != nil {
			t.Error(err)
		}

		// Check if the user is created
		if user == nil {
			t.Error("User is nil")
			return
		}

		// Check if the user has an ID
		if user.ID == 0 {
			t.Error("User ID is empty")
		}

		// Check user fields
		if user.FirstName != createUserPayload.FirstName {
			t.Errorf("User first name is not equal to %s", createUserPayload.FirstName)
		}

		if user.LastName != createUserPayload.LastName {
			t.Errorf("User last name is not equal to %s", createUserPayload.LastName)
		}

		if user.Email != createUserPayload.Email {
			t.Errorf("User email %s is not equal to %s", user.Email, createUserPayload.Email)
		}

		if user.CreatedAt.IsZero() {
			t.Error("User created at is empty")
		}

		// Delete the user for cleanup
		err = storage.DeleteUserById(user.ID)
		if err != nil {
			t.Error(err)
		}
	})
}
